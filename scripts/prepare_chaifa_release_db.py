#!/usr/bin/env python3
"""
制备柴发独立部署用 ism.db：
  1) 从 Sqlite3 备份还原
  2) 重置 admin 密码为 bcrypt(MD5(123456))
  3) 监控树按楼层重组 Zone
  4) 克隆中航信最新三页运行模板（首页/设备列表/点位列表）为大屏
  5) 设置项目默认首页大屏

大屏请用:
  python3 scripts/clone_hx_3page_dashboard_to_chaifa.py --src-db ... --dst-db ...
本脚本内建旧版多页大屏已废弃，改为调用上述克隆逻辑。

用法:
  python3 scripts/prepare_chaifa_release_db.py \\
    --sql Sqlite3_Backup_2026-07-13_10-37-19.zip \\
    --out releases/.../ism_server_user/data/db/ism.db
"""

from __future__ import annotations

import argparse
import base64
import hashlib
import json
import re
import sqlite3
import subprocess
import sys
import tempfile
import uuid
import zipfile
from collections import defaultdict
from datetime import datetime
from pathlib import Path

try:
    import bcrypt
except ImportError:
    print("需要 bcrypt: pip install bcrypt", file=sys.stderr)
    sys.exit(1)

NOW = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
PROJECT_NAME = "后沙峪改造-柴发部分"
DISPLAY_NAME = "柴发楼监控大屏"
NS = uuid.UUID("a1b2c3d4-e5f6-7890-abcd-ef1234567890")


def new_id() -> str:
    return str(uuid.uuid4())


def page_id(key: str) -> str:
    return str(uuid.uuid5(NS, f"chaifa-{key}"))


def restore_sql(sql_or_zip: Path, db_path: Path) -> None:
    db_path.parent.mkdir(parents=True, exist_ok=True)
    if db_path.exists():
        db_path.unlink()
    with tempfile.TemporaryDirectory(prefix="chaifa_sql_") as tmp:
        tmpdir = Path(tmp)
        work = sql_or_zip
        if sql_or_zip.suffix.lower() == ".zip":
            with zipfile.ZipFile(sql_or_zip) as zf:
                names = [n for n in zf.namelist() if n.endswith(".sql")]
                if not names:
                    raise RuntimeError("zip 内无 .sql")
                zf.extract(names[0], tmpdir)
                work = tmpdir / names[0]
        r = subprocess.run(["sqlite3", str(db_path)], input=work.read_bytes(), capture_output=True)
        if r.returncode != 0:
            raise RuntimeError(r.stderr.decode("utf-8", errors="replace")[:2000])


def reset_admin_password(conn: sqlite3.Connection) -> None:
    md5 = hashlib.md5(b"123456").hexdigest()
    hashed = bcrypt.hashpw(md5.encode(), bcrypt.gensalt(rounds=12)).decode()
    conn.execute(
        "UPDATE user SET password=?, updated_at=? WHERE username='admin'",
        (hashed, NOW),
    )
    # project_user 若有 admin 一并同步
    try:
        conn.execute(
            "UPDATE project_user SET password=?, updated_at=? WHERE username='admin'",
            (hashed, NOW),
        )
    except sqlite3.Error:
        pass
    print(f"  admin 密码已重置为 bcrypt(MD5(123456))")


def restructure_floors(conn: sqlite3.Connection, project_uuid: str) -> dict[str, list[dict]]:
    """RootZone 下按楼层建 Zone，设备挂到对应楼层。返回 {楼层名: [devices]}"""
    floors_order = ["柴发楼1层", "柴发楼2层", "柴发楼3层", "柴发楼4层"]
    # 兼容「一层」
    alias = {"柴发楼一层": "柴发楼1层"}

    cur = conn.execute(
        """
        SELECT uuid, sid, pid, name, type, muid, extra_data, device_type
        FROM monitor_list
        WHERE project_uuid=? AND deleted_at IS NULL
        ORDER BY sid
        """,
        (project_uuid,),
    )
    rows = [dict(zip([c[0] for c in cur.description], r)) for r in cur.fetchall()]
    root = next((r for r in rows if r["type"] == 0 and r["pid"] == 0), None)
    if not root:
        raise RuntimeError("未找到 RootZone")
    root_sid = int(root["sid"])

    devices = [r for r in rows if r["type"] == 1]
    by_floor: dict[str, list[dict]] = defaultdict(list)
    for d in devices:
        name = d["name"] or ""
        floor = None
        for key in floors_order:
            if key in name:
                floor = key
                break
        if not floor:
            for old, new in alias.items():
                if old in name:
                    floor = new
                    break
        if not floor:
            # 从「柴发楼N层」或设备名推断
            m = re.search(r"柴发楼([1-4一二三四])层", name)
            if m:
                digit = {"1": "1", "2": "2", "3": "3", "4": "4", "一": "1", "二": "2", "三": "3", "四": "4"}[
                    m.group(1)
                ]
                floor = f"柴发楼{digit}层"
            else:
                floor = "柴发楼其他"
        by_floor[floor].append(d)

    # 删除旧 Zone（非 Root）
    conn.execute(
        """
        DELETE FROM monitor_list
        WHERE project_uuid=? AND type=0 AND pid!=0 AND deleted_at IS NULL
        """,
        (project_uuid,),
    )

    # 分配楼层 sid：从 3000 起
    next_sid = 3000
    floor_sids: dict[str, int] = {}
    for floor in floors_order + [f for f in by_floor if f not in floors_order]:
        if floor not in by_floor:
            continue
        sid = next_sid
        next_sid += 1
        floor_sids[floor] = sid
        zuuid = page_id(f"zone-{floor}")
        conn.execute(
            """
            INSERT INTO monitor_list (
              created_at, updated_at, sid, pid, name, type, timeout, project_uuid,
              "interval", failed_times, described, device_type, muid, uuid,
              extra_data, status, is_enable, offline_clear, offline_default_value
            ) VALUES (?,?,?,?,?,0,5,?,5,5,?,0,'',?,'{}',0,1,0,'0')
            """,
            (NOW, NOW, sid, root_sid, floor, project_uuid, f"{floor}区域", zuuid),
        )

    # 设备挂到楼层
    for floor, devs in by_floor.items():
        psid = floor_sids.get(floor, root_sid)
        for d in devs:
            conn.execute(
                "UPDATE monitor_list SET pid=?, updated_at=? WHERE uuid=?",
                (psid, NOW, d["uuid"]),
            )

    # 重命名 RootZone
    conn.execute(
        "UPDATE monitor_list SET name=?, updated_at=? WHERE uuid=?",
        ("柴发楼", NOW, root["uuid"]),
    )
    print(f"  楼层 Zone: { {k: len(v) for k, v in by_floor.items()} }")
    return {k: by_floor[k] for k in floors_order if k in by_floor}


def default_animate() -> dict:
    return {
        "condition": {
            "deviceSN": "",
            "selectVideoType": 0,
            "isBandDevice": False,
            "bandType": 1,
            "dataID": "",
            "dataName": "",
            "operator": "",
            "OperatorValue": "",
            "OperatorMaxValue": "",
        },
        "isExpression": False,
        "animateList": [],
        "animateElement": [],
        "selected": [],
    }


def make_text(text, x, y, w=300, h=36, font_size=16, color="#e0e8f0",
              font_weight="normal", text_align="left", z=10, action=None, name=""):
    cid = new_id()
    cell = {
        "shape": "view-svg-text",
        "id": cid,
        "x": x,
        "y": y,
        "width": w,
        "height": h,
        "zIndex": z,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "name": name or text[:24],
                "identifier": cid,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "text": text,
                    "visible": 1,
                    "fontSize": font_size,
                    "fontWeight": font_weight,
                    "foreColor": color,
                    "textAlign": text_align,
                    "diy": [],
                    "zIndex": z,
                },
                "animate": default_animate(),
                "action": [action] if action else [],
                "active": [],
                "dataBind": [],
            }
        },
    }
    return cell


def make_border(x, y, w, h, z=1, shape="dv-border-box1"):
    cid = new_id()
    return {
        "shape": shape,
        "id": cid,
        "x": x,
        "y": y,
        "width": w,
        "height": h,
        "zIndex": z,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": shape,
                "name": "border",
                "identifier": cid,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "visible": 1,
                    "diy": [],
                    "zIndex": z,
                },
                "animate": default_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def link_action(model_id: str, target_page: str) -> dict:
    return {
        "type": "click",
        "action": "link",
        "isPopUp": False,
        "link": {
            "linkType": "Inside",
            "Inside": {
                "displayUUID": model_id,
                "pageUUID": target_page,
                "displayType": 1,
            },
        },
    }


def pick_columns(conn: sqlite3.Connection, device_uuid: str, limit: int = 10) -> list[str]:
    """为单设备挑一批有代表性的测点名。"""
    preferred_suffixes = [
        "AB线电压",
        "BC线电压",
        "CA线电压",
        "A相电流",
        "B相电流",
        "C相电流",
        "总有功功率",
        "频率",
    ]
    cur = conn.execute(
        """
        SELECT name FROM device_real_data
        WHERE device_uuid=? AND deleted_at IS NULL
        ORDER BY id
        """,
        (device_uuid,),
    )
    names = [r[0] for r in cur.fetchall() if r[0]]
    picked = []
    for suf in preferred_suffixes:
        for n in names:
            if n.endswith(suf) and n not in picked:
                picked.append(n)
                break
        if len(picked) >= limit:
            break
    if len(picked) < 4:
        for n in names:
            if any(k in n for k in ("电压", "电流", "功率", "频率")) and n not in picked:
                picked.append(n)
            if len(picked) >= limit:
                break
    return picked[:limit]


def make_real_table(title: str, devices: list[dict], columns: list[str],
                    x=40, y=140, w=1840, h=860) -> dict:
    cid = new_id()
    device_names = [d["name"] for d in devices]
    device_codes = [d["name"] for d in devices]
    bindings = []
    for code in device_codes:
        for col in columns:
            bindings.append(f"{code}->{col}")
    return {
        "shape": "ism-view-real-table",
        "id": cid,
        "x": x,
        "y": y,
        "width": w,
        "height": h,
        "zIndex": 20,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "ism-view-real-table",
                "name": f"{title}-实时数据",
                "identifier": cid,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "visible": 1,
                    "backColor": "#0d1b2a",
                    "foreColor": "#e0e8f0",
                    "fontSize": 13,
                    "fontFamily": "Arial",
                    "diy": [
                        {
                            "key": "columnHeaders",
                            "name": "configComponent.viewRealTable.columnHeaders",
                            "type": 9,
                            "value": ",".join(columns),
                        },
                        {
                            "key": "rowDeviceNames",
                            "name": "configComponent.viewRealTable.rowDeviceNames",
                            "type": 9,
                            "value": ",".join(device_names),
                        },
                        {
                            "key": "rowDeviceCodes",
                            "name": "configComponent.viewRealTable.rowDeviceCodes",
                            "type": 9,
                            "value": ",".join(device_codes),
                        },
                        {
                            "key": "rowBindings",
                            "name": "configComponent.viewRealTable.rowBindings",
                            "type": 9,
                            "value": ",".join(bindings),
                        },
                        {
                            "key": "waitTime",
                            "name": "configComponent.AlarmList.waitTime",
                            "type": 2,
                            "value": 5,
                        },
                    ],
                    "zIndex": 20,
                },
                "animate": default_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def dark_layer() -> dict:
    return {
        "width": 1920,
        "height": 1080,
        "backColor": "#0a0e17",
        "backgroundImage": "",
        "animate": "fadeIn",
        "autoSize": 1,
    }


def encode_components(cells: list) -> str:
    raw = json.dumps({"cells": cells}, ensure_ascii=False, separators=(",", ":"))
    return base64.b64encode(raw.encode("utf-8")).decode("ascii")


def encode_layer(layer: dict) -> str:
    return base64.b64encode(
        json.dumps(layer, ensure_ascii=False, separators=(",", ":")).encode("utf-8")
    ).decode("ascii")


def upsert_page(conn: sqlite3.Connection, model_id: str, page_name: str, page_uuid: str,
                cells: list, is_home: int = 0) -> None:
    layer_b64 = encode_layer(dark_layer())
    comp_b64 = encode_components(cells)
    conn.execute(
        "DELETE FROM display_model_layer WHERE model_id=? AND page_id=?",
        (model_id, page_uuid),
    )
    # 清掉 DisplayModelAdd 可能产生的 demo 页
    conn.execute(
        "DELETE FROM display_model_layer WHERE model_id=? AND page_name='demo'",
        (model_id,),
    )
    conn.execute(
        """
        INSERT INTO display_model_layer (
          created_at, updated_at, model_id, page_name, page_id,
          is_home, page_type, layer, components, is_login
        ) VALUES (?,?,?,?,?,?,1,?,?,0)
        """,
        (NOW, NOW, model_id, page_name, page_uuid, is_home, layer_b64, comp_b64),
    )


def build_dashboard(conn: sqlite3.Connection, project_uuid: str,
                    floors: dict[str, list[dict]]) -> str:
    model_id = page_id("display-main")
    # display_models
    conn.execute(
        "DELETE FROM display_model_layer WHERE model_id IN "
        "(SELECT display_model_uid FROM display_models WHERE project_uuid=?)",
        (project_uuid,),
    )
    conn.execute("DELETE FROM display_models WHERE project_uuid=?", (project_uuid,))
    conn.execute(
        """
        INSERT INTO display_models (
          created_at, updated_at, name, project_uuid, description,
          display_model_uid, display_image, display_user_list, display_type
        ) VALUES (?,?,?,?,?,?, '','',1)
        """,
        (NOW, NOW, DISPLAY_NAME, project_uuid, "柴发楼独立部署默认大屏", model_id),
    )

    floor_pages = {fname: page_id(f"floor-{fname}") for fname in floors}
    device_pages = {}
    for fname, devs in floors.items():
        for d in devs:
            device_pages[d["uuid"]] = page_id(f"dev-{d['uuid']}")

    # ── Overview ──
    cells = [
        make_border(20, 20, 1880, 1040, z=0, shape="dv-border-box8"),
        make_text(
            "后沙峪 · 柴发楼监控总览",
            60,
            40,
            900,
            48,
            font_size=32,
            color="#00d4ff",
            font_weight="bold",
            name="title",
        ),
        make_text(
            f"设备 {sum(len(v) for v in floors.values())} 台  ·  楼层 {len(floors)}  ·  Modbus TCP",
            60,
            95,
            800,
            28,
            font_size=14,
            color="#8fa3b8",
            name="subtitle",
        ),
    ]

    # KPI cards per floor
    card_w = 420
    gap = 20
    start_x = 60
    y = 150
    for i, (fname, devs) in enumerate(floors.items()):
        x = start_x + i * (card_w + gap)
        if i > 3:
            break
        cells.append(make_border(x, y, card_w, 160, z=2, shape="dv-border-box13"))
        cells.append(
            make_text(
                fname,
                x + 24,
                y + 24,
                300,
                36,
                font_size=20,
                color="#00d4ff",
                font_weight="bold",
                action=link_action(model_id, floor_pages[fname]),
                name=f"kpi-{fname}",
            )
        )
        cells.append(
            make_text(
                f"{len(devs)} 台采集端口",
                x + 24,
                y + 70,
                300,
                36,
                font_size=28,
                color="#ffffff",
                font_weight="bold",
                action=link_action(model_id, floor_pages[fname]),
            )
        )
        cells.append(
            make_text(
                "点击进入楼层 →",
                x + 24,
                y + 115,
                200,
                28,
                font_size=13,
                color="#5eead4",
                action=link_action(model_id, floor_pages[fname]),
            )
        )

    # Device grid
    cells.append(make_border(60, 340, 1800, 700, z=1, shape="dv-border-box1"))
    cells.append(
        make_text(
            "设备一览（点击进入设备详情）",
            90,
            360,
            500,
            32,
            font_size=18,
            color="#00d4ff",
            font_weight="bold",
        )
    )
    all_devs = []
    for fname, devs in floors.items():
        for d in sorted(devs, key=lambda x: x["name"]):
            all_devs.append((fname, d))

    col_w, row_h = 420, 48
    cols = 4
    ox, oy = 90, 410
    for idx, (fname, d) in enumerate(all_devs):
        c = idx % cols
        r = idx // cols
        x = ox + c * (col_w + 16)
        y = oy + r * (row_h + 8)
        if y + row_h > 1000:
            break
        short = d["name"].replace("柴发楼", "")
        cells.append(
            make_text(
                f"▸ {short}",
                x,
                y,
                col_w,
                row_h,
                font_size=15,
                color="#c8d6e5",
                action=link_action(model_id, device_pages[d["uuid"]]),
                name=d["name"],
            )
        )

    upsert_page(conn, model_id, "监控总览", model_id, cells, is_home=1)

    # ── Floor pages ──
    for fname, devs in floors.items():
        fcells = [
            make_border(20, 20, 1880, 1040, z=0, shape="dv-border-box8"),
            make_text(
                fname,
                60,
                40,
                600,
                44,
                font_size=28,
                color="#00d4ff",
                font_weight="bold",
            ),
            make_text(
                "← 返回总览",
                1600,
                45,
                200,
                32,
                font_size=14,
                color="#5eead4",
                action=link_action(model_id, model_id),
            ),
        ]
        # device buttons
        for i, d in enumerate(sorted(devs, key=lambda x: x["name"])):
            x = 60 + (i % 4) * 450
            y = 110 + (i // 4) * 56
            fcells.append(
                make_text(
                    f"▸ {d['name']}",
                    x,
                    y,
                    430,
                    44,
                    font_size=15,
                    color="#e0e8f0",
                    action=link_action(model_id, device_pages[d["uuid"]]),
                )
            )

        # real table: use shared short suffixes if possible; else full names from first device
        sample = sorted(devs, key=lambda x: x["name"])[0]
        cols = pick_columns(conn, sample["uuid"], limit=8)
        # Prefer suffix-only headers when all devices share same suffix set
        short_cols = []
        for c in cols:
            for suf in ("AB线电压", "BC线电压", "CA线电压", "A相电流", "B相电流", "C相电流", "总有功功率", "频率"):
                if c.endswith(suf):
                    short_cols.append(suf)
                    break
            else:
                short_cols.append(c)
        # Bindings need full names — rebuild with per-device match
        # Simpler: use full names from first device as columns; table may show blank for others
        # Better: use suffix matching in bindings device->suffix and columns=suffixes
        if len(set(short_cols)) == len(short_cols) and all(
            any(c.endswith(s) for c in cols) for s in short_cols
        ):
            # Build per-device full name map for bindings via suffix
            table_cols = short_cols
            # For ViewRealTable, binding is deviceName->dataPointName (exact name in device_real_data)
            # So we need exact names — take first matching point per device per suffix
            bind_cols = []
            for suf in table_cols:
                # use first device's full name ending with suf as column header (common pattern in ISM)
                full = next((c for c in cols if c.endswith(suf)), suf)
                bind_cols.append(full)
            # Actually for multi-device different feeders, short suffix won't match.
            # Use device detail pages for tables; floor page just navigation.
            fcells.append(
                make_text(
                    f"本层 {len(devs)} 台设备 · 点击上方端口进入实时数据表",
                    60,
                    280,
                    1000,
                    32,
                    font_size=14,
                    color="#8fa3b8",
                )
            )
            # show compact list of first points as hint
            hint = "、".join(short_cols[:6]) if short_cols else ""
            if hint:
                fcells.append(
                    make_text(
                        f"典型测点: {hint}",
                        60,
                        320,
                        1600,
                        28,
                        font_size=13,
                        color="#64748b",
                    )
                )
        else:
            fcells.append(
                make_text(
                    f"本层 {len(devs)} 台设备，点击进入详情查看实时数据",
                    60,
                    280,
                    1000,
                    32,
                    font_size=14,
                    color="#8fa3b8",
                )
            )

        upsert_page(conn, model_id, fname, floor_pages[fname], fcells, is_home=0)

    # ── Device pages ──
    for fname, devs in floors.items():
        for d in sorted(devs, key=lambda x: x["name"]):
            cols = pick_columns(conn, d["uuid"], limit=12)
            if not cols:
                cols = ["(无测点)"]
            dcells = [
                make_border(20, 20, 1880, 1040, z=0, shape="dv-border-box8"),
                make_text(
                    d["name"],
                    60,
                    36,
                    1200,
                    40,
                    font_size=24,
                    color="#00d4ff",
                    font_weight="bold",
                ),
                make_text(
                    f"← 返回{fname}",
                    1500,
                    40,
                    280,
                    32,
                    font_size=14,
                    color="#5eead4",
                    action=link_action(model_id, floor_pages[fname]),
                ),
            ]
            if cols and cols[0] != "(无测点)":
                dcells.append(make_real_table(d["name"], [d], cols, x=40, y=100, w=1840, h=920))
            else:
                dcells.append(
                    make_text("该设备暂无实时数据点", 60, 120, 400, 32, font_size=16, color="#f87171")
                )
            upsert_page(conn, model_id, d["name"], device_pages[d["uuid"]], dcells, is_home=0)

    # Home dashboard config (project scope + global fallback)
    cfg = json.dumps(
        {"dashboardUuid": model_id, "projectUuid": project_uuid},
        ensure_ascii=False,
        separators=(",", ":"),
    )
    for scope in (project_uuid, "ism.system"):
        conn.execute(
            "DELETE FROM system_data_model WHERE project_uuid=? AND uuid=?",
            (scope, "ism.SystemHomeDashboard"),
        )
        conn.execute(
            """
            INSERT INTO system_data_model (
              created_at, updated_at, name, uuid, auth, type, data_unit,
              conversion_expression, is_alarm, alarm_level, alarm_message,
              alarm_clear_message, is_record, record_type, record_interval,
              record_data_charge, value, project_uuid
            ) VALUES (?,?,?,?, 'ReadWrite', 0, '', '', 0, 0, '', '', 0, 0, 0, '', ?, ?)
            """,
            (
                NOW,
                NOW,
                "ProjectHomeDashboard" if scope == project_uuid else "SystemHomeDashboard",
                "ism.SystemHomeDashboard",
                cfg,
                scope,
            ),
        )

    # Bind project configUid so tree can jump to dashboard
    conn.execute(
        "UPDATE project_lists SET updated_at=? WHERE uuid=?",
        (NOW, project_uuid),
    )
    # devices_model / monitor configUid optional
    conn.execute(
        "UPDATE monitor_list SET configuration_uid=?, page_uuid=?, updated_at=? "
        "WHERE project_uuid=? AND type=1",
        (model_id, model_id, NOW, project_uuid),
    )

    n_pages = conn.execute(
        "SELECT COUNT(*) FROM display_model_layer WHERE model_id=? AND deleted_at IS NULL",
        (model_id,),
    ).fetchone()[0]
    print(f"  大屏 model_id={model_id}  页面数={n_pages}")
    return model_id


def vacuum_hint(db_path: Path) -> None:
    # optional compact
    pass


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--sql", required=True, help="Sqlite3_Backup_*.sql 或 .zip")
    ap.add_argument("--out", required=True, help="输出 ism.db 路径")
    ap.add_argument("--project-name", default=PROJECT_NAME)
    args = ap.parse_args()

    sql_path = Path(args.sql).expanduser().resolve()
    out_path = Path(args.out).expanduser().resolve()

    print(f"[1/5] 还原备份 → {out_path}")
    restore_sql(sql_path, out_path)

    conn = sqlite3.connect(str(out_path))
    try:
        print("[2/5] 重置 admin 密码")
        reset_admin_password(conn)

        row = conn.execute(
            "SELECT uuid, name FROM project_lists WHERE deleted_at IS NULL ORDER BY id LIMIT 1"
        ).fetchone()
        if not row:
            raise RuntimeError("无项目")
        project_uuid, pname = row
        if args.project_name and args.project_name != pname:
            conn.execute(
                "UPDATE project_lists SET name=?, updated_at=? WHERE uuid=?",
                (args.project_name, NOW, project_uuid),
            )
            pname = args.project_name
        print(f"  项目: {pname} ({project_uuid})")

        # creator → admin
        admin_uuid = conn.execute(
            "SELECT uuid FROM user WHERE username='admin'"
        ).fetchone()[0]
        conn.execute(
            "UPDATE project_lists SET creator_uuid=?, creator=?, updated_at=? WHERE uuid=?",
            (admin_uuid, "admin", NOW, project_uuid),
        )

        print("[3/5] 重组楼层 Zone")
        floors = restructure_floors(conn, project_uuid)

        print("[4/5] 克隆中航信最新三页大屏")
        conn.commit()
        conn.close()
        conn = None
        # 从主库（或 --hx-src）克隆三页模板
        hx_src = Path(__file__).resolve().parents[1] / "ism_server_user" / "data" / "db" / "ism.db"
        import subprocess as _sp
        r = _sp.run(
            [
                sys.executable,
                str(Path(__file__).resolve().parents[1] / "scripts" / "clone_hx_3page_dashboard_to_chaifa.py"),
                "--src-db",
                str(hx_src),
                "--dst-db",
                str(out_path),
            ],
            check=False,
        )
        if r.returncode != 0:
            raise RuntimeError("clone_hx_3page_dashboard_to_chaifa.py 失败")
        conn = sqlite3.connect(str(out_path))
        model_id = conn.execute(
            "SELECT display_model_uid FROM display_models WHERE project_uuid=? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
            (project_uuid,),
        ).fetchone()[0]

        print("[5/5] 完成")
        print(f"  DB: {out_path} ({out_path.stat().st_size / 1024 / 1024:.1f} MB)")
        print(f"  登录: admin / 123456")
        print(f"  大屏: /#/AppRun/{model_id}")
    finally:
        if conn is not None:
            conn.close()

    # compact
    subprocess.run(["sqlite3", str(out_path), "VACUUM;"], check=False)
    return 0


if __name__ == "__main__":
    sys.exit(main())

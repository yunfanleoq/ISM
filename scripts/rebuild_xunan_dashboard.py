#!/usr/bin/env python3
"""循安电力监控大屏重设计 — 修复合规 detail 结构并绑定设备实时数据。"""
import base64
import json
import sqlite3
import uuid
from datetime import datetime
from pathlib import Path

DB_PATH = Path(__file__).resolve().parents[1] / "ism_server_user" / "data" / "db" / "ism.db"
PROJECT_UUID = "3ec5821f-b512-2adb-3e1c-473720d0a93e"
MODEL_ID = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
CANVAS_W, CANVAS_H = 1920, 1080

# 15 个配电室 page_id — 与 ism-front-end-v2/src/config/xunanDashboardPages.js 保持同步
ROOM_PAGES = [
    ("1A配电室", "a212682e-70a1-bcd4-c91f-49a5c6786f5f", ["1A", "1A1"]),
    ("1B配电室", "406a4f48-5c75-49b1-0019-0e6f80e584cd", ["1B", "1B1"]),
    ("2A1配电室", "07c7fa5b-f71a-093b-25d2-7e9f2e2eb658", ["2A1"]),
    ("2A2配电室", "8f1193e9-dbbd-f183-aaab-0ecb2373412c", ["2A2"]),
    ("2A3配电室", "4fb22468-e607-93da-5909-5d729b1ff731", ["2A3"]),
    ("2A4配电室", "5badaa89-bbb1-9159-4d9e-807e3c305d4d", ["2A4"]),
    ("2B1配电室", "5c2be637-be7d-b3ec-bb42-7a36f290cdf3", ["2B1"]),
    ("2B2配电室", "9387c9e9-c49b-4bb0-be5a-6c4c305c2b91", ["2B2"]),
    ("2B3配电室", "914a33f2-66ac-533e-4fc9-cff7e5447fad", ["2B3"]),
    ("2B4配电室", "8dc72a27-4d54-2f24-d257-447088b36d43", ["2B4"]),
    ("3A1配电室", "b82cbaa7-8942-3b6f-424a-439d24702bc9", ["3A1"]),
    ("3A2配电室", "ffa295da-576d-19e1-4b36-0845ea14faf1", ["3A2"]),
    ("3A3配电室", "86e5c739-59cb-f418-c337-59e4191439d5", ["3A3"]),
    ("3A4配电室", "e39454ed-e4de-8dac-6263-e0bfd04ca8cb", ["3A4"]),
    ("4A1配电室", "2f468fd5-def4-8109-53e4-ca460b74a520", ["4A1"]),
]

# 首页 page_id 必须与 MODEL_ID 一致 —— ISMRunTreeNav 根节点 GoPage 使用 modelId
MAIN_PAGE_ID = MODEL_ID
ALARM_PAGE_ID = "d5b1c66d-cb33-df7d-a25a-e9e1daff8dbf"

# 变电所编码 → 配电室 page_id（与 ISMRunTreeNav.vue XUNAN_ROOM_PAGE_MAP 保持同步）
ROOM_PAGE_BY_KEY = {}
for _room_name, _page_id, _prefixes in ROOM_PAGES:
    for _pfx in _prefixes:
        ROOM_PAGE_BY_KEY[_pfx.upper()] = _page_id


def new_id():
    return str(uuid.uuid4())


def default_animate():
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


def make_text_cell(text, x, y, w=300, h=36, font_size=16, font_weight="normal",
                   color="#e0e8f0", text_align="left", name="", z_index=10):
    cid = new_id()
    if not name:
        name = text[:24]
    return {
        "shape": "view-svg-text",
        "id": cid,
        "x": x, "y": y,
        "width": w, "height": h,
        "zIndex": z_index,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "name": name,
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
                    "zIndex": z_index,
                },
                "animate": default_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def make_border_box(x, y, w, h, z_index=1):
    cid = new_id()
    return {
        "shape": "dv-border-box1",
        "id": cid,
        "x": x, "y": y,
        "width": w, "height": h,
        "zIndex": z_index,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-border-box1",
                "name": "border",
                "identifier": cid,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "visible": 1,
                    "diy": [],
                    "zIndex": z_index,
                },
                "animate": default_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def make_real_table(room_name, devices, x=40, y=120, w=1840, h=880):
    """生成绑定设备的实时数据表格。"""
    cid = new_id()
    device_names = [d["name"] for d in devices[:12]]
    device_codes = [d["name"] for d in devices[:12]]
    columns = ["AB线电压", "BC线电压", "CA线电压", "A相电流", "B相电流", "C相电流", "频率", "有功功率"]
    bindings = []
    for code in device_codes:
        for col in columns:
            bindings.append(f"{code}->{col}")
    return {
        "shape": "ism-view-real-table",
        "id": cid,
        "x": x, "y": y,
        "width": w, "height": h,
        "zIndex": 20,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "ism-view-real-table",
                "name": f"{room_name}-实时数据",
                "identifier": cid,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "visible": 1,
                    "backColor": "#0d1b2a",
                    "foreColor": "#e0e8f0",
                    "fontSize": 14,
                    "fontFamily": "Arial",
                    "diy": [
                        {"key": "columnHeaders", "name": "configComponent.viewRealTable.columnHeaders",
                         "type": 9, "value": ",".join(columns)},
                        {"key": "rowDeviceNames", "name": "configComponent.viewRealTable.rowDeviceNames",
                         "type": 9, "value": ",".join(device_names)},
                        {"key": "rowDeviceCodes", "name": "configComponent.viewRealTable.rowDeviceCodes",
                         "type": 9, "value": ",".join(device_codes)},
                        {"key": "rowBindings", "name": "configComponent.viewRealTable.rowBindings",
                         "type": 9, "value": ",".join(bindings)},
                        {"key": "waitTime", "name": "configComponent.AlarmList.waitTime",
                         "type": 7, "value": 1000, "min": 100, "max": 10000},
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


def make_layer():
    return {
        "width": CANVAS_W,
        "height": CANVAS_H,
        "backColor": "#0a0e17",
        "backgroundImage": "",
        "widthHeightRatio": 0,
    }


def b64_pack(cells):
    payload = json.dumps({"cells": cells}, ensure_ascii=False, separators=(",", ":"))
    return base64.b64encode(payload.encode("utf-8")).decode("ascii")


def b64_layer(layer):
    payload = json.dumps(layer, ensure_ascii=False, separators=(",", ":"))
    return base64.b64encode(payload.encode("utf-8")).decode("ascii")


def load_devices(cur):
    cur.execute(
        """SELECT uuid, name, muid, status FROM monitor_list
           WHERE project_uuid=? AND deleted_at IS NULL AND type=1
           ORDER BY name""",
        (PROJECT_UUID,),
    )
    all_devs = [{"uuid": r[0], "name": r[1], "muid": r[2], "status": r[3]} for r in cur.fetchall()]
    by_room = {}
    for room_name, _, prefixes in ROOM_PAGES:
        matched = []
        for d in all_devs:
            n = d["name"]
            if any(p in n for p in prefixes) or room_name.replace("配电室", "") in n:
                matched.append(d)
        if not matched:
            matched = [d for d in all_devs if room_name[:2] in d["name"]][:8]
        by_room[room_name] = matched
    return all_devs, by_room


def build_main_page(cur, all_devs, by_room):
    online = sum(1 for d in all_devs if d["status"] == 1)
    offline = len(all_devs) - online
    cells = [
        make_border_box(0, 0, CANVAS_W, 80, 0),
        make_text_cell("循安电力监控系统", 40, 18, 600, 44, 28, "bold", "#00d4ff", "left", "主标题", 5),
        make_text_cell(
            f"设备总数 {len(all_devs)}  |  在线 {online}  |  离线 {offline}  |  刷新 5s",
            700, 28, 700, 32, 14, "normal", "#90a4be", "left", "统计", 5,
        ),
        make_text_cell(
            datetime.now().strftime("更新时间 %Y-%m-%d %H:%M"),
            1500, 28, 380, 32, 14, "normal", "#607d8b", "right", "时间", 5,
        ),
        make_border_box(20, 100, 1880, 960, 1),
        make_text_cell("配电室导航", 40, 110, 200, 32, 18, "bold", "#4fc3f7", "left", "导航标题", 6),
    ]
    cols, card_w, card_h, gap = 5, 350, 100, 20
    start_x, start_y = 40, 160
    for i, (room_name, page_id, _) in enumerate(ROOM_PAGES):
        col, row = i % cols, i // cols
        x = start_x + col * (card_w + gap)
        y = start_y + row * (card_h + gap)
        dev_count = len(by_room.get(room_name, []))
        cells.append(make_border_box(x, y, card_w, card_h, 2))
        cells.append(make_text_cell(
            room_name, x + 16, y + 12, card_w - 32, 32, 18, "bold", "#81c784", "left",
            f"nav-{room_name}", 8,
        ))
        cells.append(make_text_cell(
            f"{dev_count} 台设备", x + 16, y + 48, card_w - 32, 24, 13, "normal", "#78909c", "left",
            f"nav-count-{room_name}", 8,
        ))
        btn_id = new_id()
        cells.append({
            "shape": "view-svg-text",
            "id": btn_id,
            "x": x + 16, "y": y + 72,
            "width": card_w - 32, "height": 24,
            "zIndex": 9, "visible": True,
            "position": {"x": x + 16, "y": y + 72},
            "size": {"width": card_w - 32, "height": 24},
            "data": {
                "detail": {
                    "type": "view-svg-text",
                    "name": f"进入{room_name}",
                    "identifier": btn_id,
                    "style": {
                        "position": {"x": x + 16, "y": y + 72, "w": card_w - 32, "h": 24},
                        "text": "进入监控 →",
                        "visible": 1,
                        "fontSize": 13,
                        "foreColor": "#00d4ff",
                        "textAlign": "left",
                        "diy": [],
                        "zIndex": 9,
                    },
                    "animate": default_animate(),
                    "action": [{
                        "actionType": 1,
                        "actionName": "跳转页面",
                        "pageId": page_id,
                        "pageName": room_name,
                    }],
                    "active": [],
                    "dataBind": [],
                }
            },
        })
    return cells


def build_room_page(room_name, devices):
    cells = [
        make_border_box(0, 0, CANVAS_W, 80, 0),
        make_text_cell(room_name, 40, 18, 500, 44, 26, "bold", "#00d4ff", "left", f"{room_name}-标题", 5),
        make_text_cell(
            f"绑定设备 {len(devices)} 台", 560, 28, 400, 32, 14, "normal", "#90a4be", "left", "设备数", 5,
        ),
        make_text_cell(
            "← 返回主界面", 40, 50, 200, 24, 13, "normal", "#607d8b", "left", "返回", 6,
        ),
    ]
    if devices:
        cells.append(make_real_table(room_name, devices))
    else:
        cells.append(make_text_cell(
            "暂无绑定设备，请在监控树中配置", 40, 200, 800, 40, 16, "normal", "#ff9800", "left", "空提示", 10,
        ))
    return cells


def sync_home_page_id(cur):
    """首页 page_id 对齐 MODEL_ID，否则运行时树 RootZone 钻探报「找不到页面」。"""
    cur.execute(
        """SELECT page_id FROM display_model_layer
           WHERE model_id=? AND is_home=1 AND deleted_at IS NULL LIMIT 1""",
        (MODEL_ID,),
    )
    row = cur.fetchone()
    if not row:
        print("[sync] ⚠️ 未找到首页")
        return 0
    old_id = row[0]
    if old_id == MODEL_ID:
        print(f"[sync] 首页 page_id 已是 MODEL_ID")
        return 0
    cur.execute(
        """UPDATE display_model_layer SET page_id=?, updated_at=?
           WHERE model_id=? AND is_home=1 AND deleted_at IS NULL""",
        (MODEL_ID, datetime.now().strftime("%Y-%m-%d %H:%M:%S"), MODEL_ID),
    )
    print(f"[sync] 首页 page_id: {old_id} → {MODEL_ID}")
    return cur.rowcount


def verify_page_ids(cur):
    """校验大屏全部 page_id 在 display_model_layer 中存在。"""
    cur.execute(
        """SELECT page_id, page_name FROM display_model_layer
           WHERE model_id=? AND deleted_at IS NULL""",
        (MODEL_ID,),
    )
    db_ids = {r[0] for r in cur.fetchall()}
    required = {MODEL_ID, ALARM_PAGE_ID}
    required.update(pid for _, pid, _ in ROOM_PAGES)
    missing = sorted(required - db_ids)
    orphan = sorted(db_ids - required)
    ok = len(missing) == 0
    print(f"[verify] 必需 page_id: {len(required)}, DB 实际: {len(db_ids)}, 缺失: {len(missing)}")
    for mid in missing:
        print(f"  ✗ 缺失: {mid}")
    if orphan:
        print(f"[verify] DB 额外页面 ({len(orphan)}):")
        for oid in orphan:
            cur.execute(
                "SELECT page_name FROM display_model_layer WHERE page_id=? AND model_id=?",
                (oid, MODEL_ID),
            )
            name = cur.fetchone()
            print(f"  · {oid} ({name[0] if name else '?'})")
    if ok:
        print("[verify] ✅ 100% page_id 校验通过")
    else:
        raise SystemExit("[verify] ❌ page_id 不完整，请先补齐 display_model_layer")
    return ok


def verify_nav_page_map():
    """校验 ROOM_PAGE_BY_KEY 中每个 page_id 都在 ROOM_PAGES 内。"""
    room_ids = {pid for _, pid, _ in ROOM_PAGES}
    bad = {k: v for k, v in ROOM_PAGE_BY_KEY.items() if v not in room_ids}
    if bad:
        raise SystemExit(f"[verify] ROOM_PAGE_BY_KEY 含非法 page_id: {bad}")
    print(f"[verify] 导航映射 {len(ROOM_PAGE_BY_KEY)} 个变电所编码 → {len(room_ids)} 个配电室页")


def update_page(cur, page_id, cells):
    comp_b64 = b64_pack(cells)
    layer_b64 = b64_layer(make_layer())
    cur.execute(
        """UPDATE display_model_layer
           SET components=?, layer=?, updated_at=?
           WHERE page_id=? AND model_id=? AND deleted_at IS NULL""",
        (comp_b64, layer_b64, datetime.now().strftime("%Y-%m-%d %H:%M:%S"), page_id, MODEL_ID),
    )
    return cur.rowcount


def main():
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    verify_nav_page_map()
    sync_home_page_id(cur)
    verify_page_ids(cur)

    all_devs, by_room = load_devices(cur)

    print(f"[1] 项目设备: {len(all_devs)} 台")
    updated = 0

    main_cells = build_main_page(cur, all_devs, by_room)
    n = update_page(cur, MAIN_PAGE_ID, main_cells)
    print(f"[2] 主界面: {n} 行更新, cells={len(main_cells)}")
    updated += n

    for room_name, page_id, _ in ROOM_PAGES:
        cells = build_room_page(room_name, by_room.get(room_name, []))
        n = update_page(cur, page_id, cells)
        print(f"[3] {room_name}: {n} 行, devices={len(by_room.get(room_name, []))}, cells={len(cells)}")
        updated += n

    conn.commit()
    verify_page_ids(cur)
    conn.close()
    print(f"\n✅ 大屏重设计完成，共更新 {updated} 个页面")
    print(f"MODEL_ID={MODEL_ID}")
    print(f"AppRun: http://localhost:7080/#/AppRun/{MODEL_ID}")


if __name__ == "__main__":
    main()

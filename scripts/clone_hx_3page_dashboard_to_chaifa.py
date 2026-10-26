#!/usr/bin/env python3
"""
将中航信最新「三页运行模板」克隆到柴发项目大屏。

三页（与 ISMRender / ISMRunTreeNav 约定一致）:
  home          → 首页模板
  deviceList    → 设备列表模板
  datapointList → 点位列表模板

用法:
  python3 scripts/clone_hx_3page_dashboard_to_chaifa.py \\
    --src-db ism_server_user/data/db/ism.db \\
    --dst-db dev-envs/chaifa-local/ism_server_user/data/db/ism.db
"""

from __future__ import annotations

import argparse
import base64
import json
import sqlite3
import uuid
from datetime import datetime
from pathlib import Path

NOW = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
HX_MODEL = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
DISPLAY_NAME = "柴发楼监控大屏"
BRAND_FROM = "中航信数据中心电力监控系统"
BRAND_TO = "后沙峪 · 柴发楼监控系统"
NS = uuid.UUID("c0ffee00-cba1-fa00-0000-000000000001")


def page_id(key: str) -> str:
    return str(uuid.uuid5(NS, f"chaifa-tpl-{key}"))


def b64_decode_maybe(s: str) -> str:
    if not s:
        return s
    try:
        return base64.b64decode(s).decode("utf-8")
    except Exception:
        return s


def b64_encode(s: str) -> str:
    return base64.b64encode(s.encode("utf-8")).decode("ascii")


def rewrite_brand(components_raw: str) -> str:
    """components 可能是 base64 或明文 JSON；替换品牌文案后仍按原形态写回。"""
    was_b64 = False
    text = components_raw
    try:
        decoded = base64.b64decode(components_raw).decode("utf-8")
        json.loads(decoded)
        text = decoded
        was_b64 = True
    except Exception:
        text = components_raw

    text2 = text.replace(BRAND_FROM, BRAND_TO)
    # 其它常见旧标题
    for old in ("中航信电力监控系统", "中航信数据中心", "航信机房"):
        text2 = text2.replace(old, BRAND_TO)

    if was_b64:
        return b64_encode(text2)
    return text2


def ensure_template_cols(conn: sqlite3.Connection) -> None:
    cols = {r[1] for r in conn.execute("PRAGMA table_info(display_model_layer)").fetchall()}
    if "template_kind" not in cols:
        conn.execute(
            "ALTER TABLE display_model_layer ADD COLUMN template_kind varchar(64) DEFAULT ''"
        )
    if "template_model_uuid" not in cols:
        conn.execute(
            "ALTER TABLE display_model_layer ADD COLUMN template_model_uuid varchar(250) DEFAULT ''"
        )


def load_hx_pages(src: sqlite3.Connection, hx_model: str = HX_MODEL) -> list[dict]:
    rows = src.execute(
        """
        SELECT page_name, page_id, is_home, page_type, layer, components,
               COALESCE(template_kind,'') AS template_kind,
               COALESCE(is_login,0) AS is_login
        FROM display_model_layer
        WHERE model_id=? AND deleted_at IS NULL
        ORDER BY is_home DESC, id
        """,
        (hx_model,),
    ).fetchall()
    if len(rows) < 3:
        raise RuntimeError(f"源库中航信大屏页数不足: {len(rows)} (期望 3)")
    pages = []
    for r in rows:
        kind = r[6] or ""
        if not kind:
            # 兜底按页名
            if "首页" in (r[0] or ""):
                kind = "home"
            elif "设备列表" in (r[0] or ""):
                kind = "deviceList"
            elif "点位列表" in (r[0] or ""):
                kind = "datapointList"
        pages.append(
            {
                "page_name": r[0],
                "src_page_id": r[1],
                "is_home": int(r[2] or 0),
                "page_type": int(r[3] or 1),
                "layer": r[4] or "",
                "components": rewrite_brand(r[5] or ""),
                "template_kind": kind,
                "is_login": int(r[7] or 0),
            }
        )
    need = {"home", "deviceList", "datapointList"}
    # 源大屏可能含大量 zone/room/floor 运行页（template_kind 空），只保留三页模板
    pages = [p for p in pages if p["template_kind"] in need]
    kinds = {p["template_kind"] for p in pages}
    if not need.issubset(kinds):
        raise RuntimeError(f"源模板 kind 不全: {kinds}")
    # 每种 kind 只取一条（优先 is_home / 已按 ORDER BY）
    picked = {}
    for p in pages:
        if p["template_kind"] not in picked:
            picked[p["template_kind"]] = p
    return [picked[k] for k in ("home", "deviceList", "datapointList")]


def apply_to_dst(dst: sqlite3.Connection, pages: list[dict], project_uuid: str) -> str:
    ensure_template_cols(dst)

    # 新 model_id：首页 page_id 必须 = model_id（运行态约定）
    model_id = page_id("home")
    home_page = next(p for p in pages if p["template_kind"] == "home")
    # 强制首页 page_id = model_id
    kind_to_pid = {
        "home": model_id,
        "deviceList": page_id("deviceList"),
        "datapointList": page_id("datapointList"),
    }

    # 清旧大屏
    old_uids = [
        r[0]
        for r in dst.execute(
            "SELECT display_model_uid FROM display_models WHERE project_uuid=? AND deleted_at IS NULL",
            (project_uuid,),
        ).fetchall()
    ]
    for uid in old_uids:
        dst.execute("DELETE FROM display_model_layer WHERE model_id=?", (uid,))
    dst.execute("DELETE FROM display_models WHERE project_uuid=?", (project_uuid,))

    dst.execute(
        """
        INSERT INTO display_models (
          created_at, updated_at, name, project_uuid, description,
          display_model_uid, display_image, display_user_list, display_type
        ) VALUES (?,?,?,?,?,?, '','',1)
        """,
        (NOW, NOW, DISPLAY_NAME, project_uuid, "克隆自中航信最新三页运行模板", model_id),
    )

    for p in pages:
        kind = p["template_kind"]
        if kind not in kind_to_pid:
            continue
        pid = kind_to_pid[kind]
        is_home = 1 if kind == "home" else 0
        # 首页名保持「首页模板」等，便于编辑器识别
        dst.execute(
            """
            INSERT INTO display_model_layer (
              created_at, updated_at, model_id, page_name, page_id,
              is_home, page_type, layer, components, is_login,
              template_kind, template_model_uuid
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?, '')
            """,
            (
                NOW,
                NOW,
                model_id,
                p["page_name"],
                pid,
                is_home,
                p["page_type"],
                p["layer"],
                p["components"],
                p["is_login"],
                kind,
            ),
        )

    # 默认首页指针
    cfg = json.dumps(
        {"dashboardUuid": model_id, "projectUuid": project_uuid},
        ensure_ascii=False,
        separators=(",", ":"),
    )
    for scope in (project_uuid, "ism.system"):
        dst.execute(
            "DELETE FROM system_data_model WHERE project_uuid=? AND uuid=?",
            (scope, "ism.SystemHomeDashboard"),
        )
        dst.execute(
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

    # 设备挂到新大屏（树点击可跳）
    dst.execute(
        "UPDATE monitor_list SET configuration_uid=?, page_uuid=?, updated_at=? "
        "WHERE project_uuid=? AND type=1",
        (model_id, model_id, NOW, project_uuid),
    )

    return model_id


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--src-db", required=True, help="含中航信三页大屏的 ism.db")
    ap.add_argument("--dst-db", required=True, help="柴发 ism.db")
    ap.add_argument("--hx-model", default=HX_MODEL)
    args = ap.parse_args()

    src_path = Path(args.src_db).resolve()
    dst_path = Path(args.dst_db).resolve()
    if not src_path.exists() or not dst_path.exists():
        raise SystemExit("src/dst db 不存在")

    src = sqlite3.connect(str(src_path))
    dst = sqlite3.connect(str(dst_path))
    try:
        pages = load_hx_pages(src, args.hx_model)
        print("源三页:")
        for p in pages:
            print(f"  {p['template_kind']:14s}  {p['page_name']}  home={p['is_home']}")

        proj = dst.execute(
            "SELECT uuid, name FROM project_lists WHERE deleted_at IS NULL ORDER BY id LIMIT 1"
        ).fetchone()
        if not proj:
            raise SystemExit("目标库无项目")
        project_uuid, pname = proj
        print(f"目标项目: {pname} ({project_uuid})")

        model_id = apply_to_dst(dst, pages, project_uuid)
        dst.commit()

        rows = dst.execute(
            """
            SELECT page_name, template_kind, is_home, length(components)
            FROM display_model_layer WHERE model_id=? AND deleted_at IS NULL
            ORDER BY is_home DESC, id
            """,
            (model_id,),
        ).fetchall()
        print(f"\n写入完成 model_id={model_id}")
        for r in rows:
            print(f"  {r[1]:14s}  {r[0]}  is_home={r[2]}  components_len={r[3]}")
        print(f"AppRun: /#/AppRun/{model_id}")
    finally:
        src.close()
        dst.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

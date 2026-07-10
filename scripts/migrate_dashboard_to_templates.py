#!/usr/bin/env python3
"""
将 NCC 默认大屏从「一节点一页」迁移为「全层级模板页」。

效果（约）:
  399 页 → ~5 + N 页（N = 物模型覆盖模板数，默认对设备数≥3 的模型建覆盖）
  组件总量随页数同比下降（每页 cells 保留一份模板）

做法:
  1. 确保 display_model_layer 有 template_kind / template_model_uuid
  2. 从现有页各抽一份样本 → 写入稳定 page_id 的模板页并打标
  3. 模板页绑点相对化：deviceSN 清空、保留 dataName、dataID 清空（运行时解析）
  4. 标题类静态设备名改为 {{nav.name}}（可识别的种子）
  5. 软删除其余 zone/room/building/floor/device 复制页（保留模板页与首页/oneline）

用法:
  # 预览
  python3 scripts/migrate_dashboard_to_templates.py

  # 写库
  python3 scripts/migrate_dashboard_to_templates.py --apply

  NCC_MODEL_ID=... NCC_PROJECT_UUID=... python3 scripts/migrate_dashboard_to_templates.py --apply
"""
from __future__ import annotations

import argparse
import base64
import json
import os
import re
import sqlite3
import sys
import uuid
from copy import deepcopy
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_MODEL = os.environ.get("NCC_MODEL_ID", "b8b4c094-faa9-a22a-1d0d-037539b27a6c")
DEFAULT_PROJECT = os.environ.get("NCC_PROJECT_UUID", "3ec5821f-b512-2adb-3e1c-473720d0a93e")
DEFAULT_SQLITE = ROOT / "ism_server_user" / "data" / "db" / "ism.db"

# 稳定模板 page_id（可反复跑幂等）
TPL_PAGE = {
    "home": uuid.uuid5(uuid.NAMESPACE_DNS, "ism-tpl-home").hex,
    "zone": uuid.uuid5(uuid.NAMESPACE_DNS, "ism-tpl-zone").hex,
    "room": uuid.uuid5(uuid.NAMESPACE_DNS, "ism-tpl-room").hex,
    "cabinet": uuid.uuid5(uuid.NAMESPACE_DNS, "ism-tpl-cabinet").hex,
    "device": uuid.uuid5(uuid.NAMESPACE_DNS, "ism-tpl-device-default").hex,
}


def device_override_page_id(muid: str) -> str:
    return uuid.uuid5(uuid.NAMESPACE_DNS, f"ism-tpl-device-{muid}").hex


def now_sql():
    return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")


def decode_components(raw: str):
    if raw is None or raw == "":
        return {"cells": []}
    s = raw
    # 可能是 base64，也可能已是 JSON
    try:
        if not s.lstrip().startswith("{"):
            s = base64.b64decode(s).decode("utf-8", errors="replace")
        data = json.loads(s)
        if isinstance(data, dict) and "cells" in data:
            return data
        if isinstance(data, list):
            return {"cells": data}
        return {"cells": []}
    except Exception:
        return {"cells": []}


def encode_components(obj) -> str:
    raw = json.dumps(obj, ensure_ascii=False, separators=(",", ":"))
    return base64.b64encode(raw.encode("utf-8")).decode("ascii")


def walk_conditions(obj, visit):
    if obj is None:
        return
    if isinstance(obj, list):
        for v in obj:
            walk_conditions(v, visit)
        return
    if not isinstance(obj, dict):
        return
    cond = obj.get("condition")
    if isinstance(cond, dict):
        visit(cond)
    for k, v in obj.items():
        if k == "condition":
            continue
        if isinstance(v, (dict, list)):
            walk_conditions(v, visit)


def relativize_components(components: dict, *, as_device_tpl: bool) -> dict:
    """清空写死的 deviceSN/dataID，保留 dataName；设备标题改为 {{nav.name}}。"""
    out = deepcopy(components)
    cells = out.get("cells") or []

    def visit(cond: dict):
        # 已绑设备的测点 → 相对化
        if cond.get("dataName") or cond.get("dataID") or cond.get("deviceSN"):
            if cond.get("isBandDevice") is True:
                # 固定设备绑点保留（跨设备写死的全局点）
                return
            cond["deviceSN"] = ""
            cond["DeviceName"] = ""
            # 保留 dataName；清空 dataID 让运行时按模型解析
            if cond.get("dataName"):
                cond["dataID"] = ""
            cond["isBandDevice"] = False

    for cell in cells:
        detail = (((cell or {}).get("data") or {}).get("detail")) or {}
        walk_conditions(detail.get("active"), visit)
        walk_conditions(detail.get("animate"), visit)
        # 图表等可能把 active 放在 diy
        style = detail.get("style") or {}
        walk_conditions(style.get("diy"), visit)

        if as_device_tpl:
            text = style.get("text")
            if isinstance(text, str) and text:
                # 设备详情标题常见形态
                if text.startswith("🔧 ") or "设备" in text[:8]:
                    # 仅替换明显是设备名的标题行
                    seed = str(cell.get("id") or "")
                    if "-title" in seed or seed.endswith("title"):
                        style["text"] = "{{nav.name}}"
                # 基本参数「设备名称」值
                seed = str(cell.get("id") or "")
                if re.search(r"-bpv-0$", seed):
                    style["text"] = "{{nav.name}}"

    out["cells"] = cells
    return out


def ensure_columns(conn: sqlite3.Connection):
    cols = {r[1] for r in conn.execute("PRAGMA table_info(display_model_layer)").fetchall()}
    if "template_kind" not in cols:
        conn.execute(
            "ALTER TABLE display_model_layer ADD COLUMN template_kind varchar(64) DEFAULT ''"
        )
        print("+ added column template_kind")
    if "template_model_uuid" not in cols:
        conn.execute(
            "ALTER TABLE display_model_layer ADD COLUMN template_model_uuid varchar(250) DEFAULT ''"
        )
        print("+ added column template_model_uuid")
    conn.commit()


def fetch_pages(conn, model_id: str):
    rows = conn.execute(
        """
        SELECT id, page_name, page_id, is_home, is_login, page_type, layer, components,
               COALESCE(template_kind,'') AS template_kind,
               COALESCE(template_model_uuid,'') AS template_model_uuid
        FROM display_model_layer
        WHERE model_id=? AND deleted_at IS NULL
        ORDER BY id
        """,
        (model_id,),
    ).fetchall()
    keys = [
        "id", "page_name", "page_id", "is_home", "is_login", "page_type",
        "layer", "components", "template_kind", "template_model_uuid",
    ]
    return [dict(zip(keys, r)) for r in rows]


def classify(page_name: str, is_home: int) -> str:
    if is_home == 1 or page_name == "main":
        return "home"
    if page_name.startswith("zone-"):
        return "zone"
    if page_name.startswith("room-"):
        return "room"
    if page_name.startswith("building-"):
        return "cabinet"
    if page_name.startswith("floor-"):
        return "floor"
    if page_name.startswith("device-"):
        return "device"
    if page_name.startswith("oneline"):
        return "oneline"
    return "other"


def soft_delete(conn, model_id: str, page_ids: list[str], apply: bool):
    if not page_ids:
        return 0
    if not apply:
        return len(page_ids)
    ts = now_sql()
    # sqlite 参数上限，分批
    n = 0
    for i in range(0, len(page_ids), 200):
        chunk = page_ids[i : i + 200]
        placeholders = ",".join("?" * len(chunk))
        conn.execute(
            f"""UPDATE display_model_layer SET deleted_at=?, updated_at=?
               WHERE model_id=? AND page_id IN ({placeholders}) AND deleted_at IS NULL""",
            [ts, ts, model_id, *chunk],
        )
        n += len(chunk)
    return n


def upsert_template(
    conn,
    model_id: str,
    *,
    page_id: str,
    page_name: str,
    kind: str,
    model_uuid: str,
    layer: str,
    components_b64: str,
    is_home: int,
    apply: bool,
):
    exist = conn.execute(
        "SELECT id FROM display_model_layer WHERE model_id=? AND page_id=? AND deleted_at IS NULL",
        (model_id, page_id),
    ).fetchone()
    if not apply:
        return "would-update" if exist else "would-insert"
    ts = now_sql()
    if exist:
        conn.execute(
            """UPDATE display_model_layer
               SET page_name=?, components=?, layer=?, is_home=?,
                   template_kind=?, template_model_uuid=?, updated_at=?, deleted_at=NULL
               WHERE model_id=? AND page_id=?""",
            (
                page_name, components_b64, layer, is_home,
                kind, model_uuid or "", ts, model_id, page_id,
            ),
        )
        return "update"
    conn.execute(
        """INSERT INTO display_model_layer
           (created_at, updated_at, deleted_at, model_id, page_name, page_id,
            is_home, is_login, page_type, layer, components, template_kind, template_model_uuid)
           VALUES (?,?,NULL,?,?,?,?,0,1,?,?,?,?)""",
        (
            ts, ts, model_id, page_name, page_id, is_home,
            layer or "", components_b64, kind, model_uuid or "",
        ),
    )
    return "insert"


def device_sid_from_name(page_name: str):
    m = re.match(r"^device-(\d+)$", page_name or "")
    return int(m.group(1)) if m else None


def load_device_muid_map(conn, project_uuid: str) -> dict[int, str]:
    rows = conn.execute(
        """SELECT sid, muid FROM monitor_list
           WHERE project_uuid=? AND type=1 AND deleted_at IS NULL""",
        (project_uuid,),
    ).fetchall()
    return {int(sid): (muid or "") for sid, muid in rows}


def pick_sample(pages: list[dict], kind: str):
    cands = [p for p in pages if classify(p["page_name"], p["is_home"]) == kind]
    if not cands:
        return None
    # 选 components 最长的（内容更完整）
    cands.sort(key=lambda p: len(p["components"] or ""), reverse=True)
    return cands[0]


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--apply", action="store_true", help="写库；默认 dry-run")
    ap.add_argument("--db", default=str(DEFAULT_SQLITE))
    ap.add_argument("--model", default=DEFAULT_MODEL)
    ap.add_argument("--project", default=DEFAULT_PROJECT)
    ap.add_argument(
        "--min-override",
        type=int,
        default=3,
        help="设备数≥该值的物模型单独建覆盖模板（默认 3）",
    )
    ap.add_argument(
        "--keep-oneline",
        action="store_true",
        default=True,
        help="保留 oneline 页（默认保留）",
    )
    args = ap.parse_args()
    apply = args.apply

    db_path = Path(args.db)
    if not db_path.exists():
        print(f"DB not found: {db_path}", file=sys.stderr)
        sys.exit(1)

    conn = sqlite3.connect(str(db_path))
    ensure_columns(conn)
    pages = fetch_pages(conn, args.model)
    print(f"model={args.model} pages={len(pages)} apply={apply}")

    by_kind: dict[str, list] = {}
    for p in pages:
        k = classify(p["page_name"], p["is_home"])
        by_kind.setdefault(k, []).append(p)
    for k, lst in sorted(by_kind.items(), key=lambda x: -len(x[1])):
        print(f"  {k}: {len(lst)}")

    sid_muid = load_device_muid_map(conn, args.project)
    # 设备页 → muid
    muid_pages: dict[str, list] = {}
    for p in by_kind.get("device", []):
        sid = device_sid_from_name(p["page_name"])
        muid = sid_muid.get(sid, "") if sid is not None else ""
        muid_pages.setdefault(muid or "_unknown", []).append(p)

    print("device by muid (top):")
    for muid, lst in sorted(muid_pages.items(), key=lambda x: -len(x[1]))[:12]:
        print(f"  {muid[:36] or '(empty)'}: {len(lst)}")

    # —— 构建模板 ——
    templates_built = []

    def build_from_sample(kind: str, sample: dict, page_id: str, page_name: str, model_uuid: str = ""):
        comps = decode_components(sample["components"])
        comps = relativize_components(comps, as_device_tpl=(kind == "device"))
        b64 = encode_components(comps)
        # home 保持 is_home=1；其它模板 is_home=0
        is_home = 1 if kind == "home" else 0
        # 若样本不是首页，home 模板仍标记 is_home
        action = upsert_template(
            conn,
            args.model,
            page_id=page_id,
            page_name=page_name,
            kind=kind,
            model_uuid=model_uuid,
            layer=sample["layer"] or "",
            components_b64=b64,
            is_home=is_home,
            apply=apply,
        )
        cell_n = len(comps.get("cells") or [])
        templates_built.append((kind, page_name, page_id, cell_n, action, model_uuid))
        print(f"  [{action}] {kind} tpl {page_name} cells={cell_n} page_id={page_id[:12]}…")

    # home：优先 is_home 页
    home = pick_sample(pages, "home") or next((p for p in pages if p["is_home"] == 1), None)
    if home:
        # 首页：保留原 page_id（导航 root 用 modelId），只打标，不换 id
        comps = decode_components(home["components"])
        # 首页不强行相对化全部绑点（可能有全局 KPI）
        b64 = encode_components(comps)
        action = upsert_template(
            conn,
            args.model,
            page_id=home["page_id"],
            page_name=home["page_name"] or "main",
            kind="home",
            model_uuid="",
            layer=home["layer"] or "",
            components_b64=b64,
            is_home=1,
            apply=apply,
        )
        templates_built.append(("home", home["page_name"], home["page_id"], len(comps.get("cells") or []), action, ""))
        print(f"  [{action}] home keep page_id={home['page_id'][:12]}… cells={len(comps.get('cells') or [])}")
        home_keep_id = home["page_id"]
    else:
        home_keep_id = None
        print("! no home page found")

    for kind, name in (("zone", "模板-区域"), ("room", "模板-机房"), ("cabinet", "模板-机柜")):
        sample = pick_sample(pages, kind)
        if not sample:
            print(f"! skip {kind}: no sample")
            continue
        build_from_sample(kind, sample, TPL_PAGE[kind], name)

    # 设备：通用模板 = 数量最多的 muid 样本；覆盖 = 其它 count>=min_override
    ranked = sorted(
        ((m, lst) for m, lst in muid_pages.items() if m and m != "_unknown"),
        key=lambda x: -len(x[1]),
    )
    if ranked:
        default_muid, default_list = ranked[0]
        default_list.sort(key=lambda p: len(p["components"] or ""), reverse=True)
        build_from_sample(
            "device",
            default_list[0],
            TPL_PAGE["device"],
            "模板-设备-通用",
            model_uuid="",  # 通用不写 muid
        )
        for muid, lst in ranked:
            if len(lst) < args.min_override:
                continue
            if muid == default_muid:
                # 同时为默认模型建覆盖，便于精确匹配
                pass
            lst.sort(key=lambda p: len(p["components"] or ""), reverse=True)
            build_from_sample(
                "device",
                lst[0],
                device_override_page_id(muid),
                f"模板-设备-{muid[:8]}",
                model_uuid=muid,
            )
    else:
        sample = pick_sample(pages, "device")
        if sample:
            build_from_sample("device", sample, TPL_PAGE["device"], "模板-设备-通用")

    # —— 软删冗余 ——
    keep_ids = {t[2] for t in templates_built}
    if home_keep_id:
        keep_ids.add(home_keep_id)

    # 保留 oneline / other（报警等）
    for p in pages:
        k = classify(p["page_name"], p["is_home"])
        if k in ("oneline", "other"):
            keep_ids.add(p["page_id"])

    delete_ids = []
    for p in pages:
        k = classify(p["page_name"], p["is_home"])
        if k in ("zone", "room", "cabinet", "floor", "device"):
            if p["page_id"] not in keep_ids:
                delete_ids.append(p["page_id"])
        # 若某样本页被选为模板但 page_id 不是稳定 TPL_PAGE，样本原页也会被删（内容已拷到模板）

    n_del = soft_delete(conn, args.model, delete_ids, apply)
    print(f"soft-delete candidates: {n_del}")

    if apply:
        conn.commit()
        # 清掉其它页上残留的错误 template 标记（仅保留 keep 模板）
        conn.execute(
            """UPDATE display_model_layer SET template_kind='', template_model_uuid=''
               WHERE model_id=? AND deleted_at IS NULL
                 AND page_id NOT IN ({})
                 AND COALESCE(template_kind,'')<>''""".format(
                ",".join("?" * len(keep_ids)) if keep_ids else "''"
            ),
            (args.model, *keep_ids) if keep_ids else (args.model,),
        )
        conn.commit()

    # 汇总
    alive = conn.execute(
        "SELECT COUNT(*) FROM display_model_layer WHERE model_id=? AND deleted_at IS NULL",
        (args.model,),
    ).fetchone()[0]
    tpl = conn.execute(
        """SELECT template_kind, COALESCE(template_model_uuid,''), page_name, page_id
           FROM display_model_layer
           WHERE model_id=? AND deleted_at IS NULL AND COALESCE(template_kind,'')<>''
           ORDER BY template_kind, template_model_uuid""",
        (args.model,),
    ).fetchall()
    print("\n=== RESULT ===")
    print(f"alive pages: {alive} (was {len(pages)})")
    print(f"templates tagged: {len(tpl)}")
    for row in tpl:
        print(f"  kind={row[0]} model={row[1][:12] or '-'} name={row[2]} id={row[3][:12]}…")
    if not apply:
        print("\n(dry-run) 加 --apply 才会写库")
    conn.close()


if __name__ == "__main__":
    main()

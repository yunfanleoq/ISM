#!/usr/bin/env python3
"""将设备/点位列表模板收敛为与运行态一致的卡片矩阵。

默认 dry-run；--apply 时先备份 SQLite 数据库，再事务更新 deviceList/datapointList 模板。
运行态会继续用 navContext 把模板预览数据替换为真实设备/测点卡片。
"""
from __future__ import annotations

import argparse
import base64
import copy
import json
import shutil
import sqlite3
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
TABLE_SHAPE = "ism-view-real-table"
MODEL_DEFAULT = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
SAMPLE_POINTS = (
    ("AB线电压", "V"),
    ("BC线电压", "V"),
    ("CA线电压", "V"),
    ("频率", "Hz"),
    ("A相电流", "A"),
    ("B相电流", "A"),
    ("C相电流", "A"),
    ("中性线电流", "A"),
    ("总有功功率", "kW"),
    ("总无功功率", "kvar"),
    ("总视在功率", "kVA"),
    ("总功率因数", "—"),
    ("正有功电度", "kWh"),
    ("A相电流谐波畸变率", "%"),
    ("B相电流谐波畸变率", "%"),
    ("C相电流谐波畸变率", "%"),
    ("中性线电流谐波畸变率", "%"),
    ("故障状态", "—"),
    ("合分闸状态", "—"),
    ("通讯状态", "—"),
)


def decode_components(value: str) -> dict:
    return json.loads(base64.b64decode(value).decode("utf-8"))


def encode_components(value: dict) -> str:
    raw = json.dumps(value, ensure_ascii=False, separators=(",", ":")).encode("utf-8")
    return base64.b64encode(raw).decode("ascii")


def detail_of(cell: dict) -> dict:
    return cell.get("data", {}).get("detail", {})


def style_of(cell: dict) -> dict:
    return detail_of(cell).get("style", {})


def is_table(cell: dict) -> bool:
    return cell.get("shape") == TABLE_SHAPE or detail_of(cell).get("type") == TABLE_SHAPE


def keep_common_header(cell: dict) -> bool:
    style = style_of(cell)
    position = style.get("position", {})
    y = float(position.get("y", cell.get("y", 9999)) or 0)
    if y > 55:
        return False
    text = str(style.get("text", "") or "").strip()
    x = float(position.get("x", cell.get("x", 0)) or 0)
    # 顶部多段历史面包屑依赖样本设备；只保留“全局总览”公共入口。
    if text and x >= 740 and "全局总览" not in text:
        return False
    return True


def set_diy(style: dict, values: dict[str, object]) -> None:
    old = {
        item.get("key"): item
        for item in style.get("diy", [])
        if isinstance(item, dict) and item.get("key")
    }
    diy = []
    for key, value in values.items():
        item = copy.deepcopy(old.get(key, {}))
        item.update({"name": item.get("name") or key, "type": item.get("type", 9), "key": key, "value": value})
        diy.append(item)
    style["diy"] = diy


def upsert_diy(style: dict, values: dict[str, object]) -> None:
    diy = [copy.deepcopy(item) for item in style.get("diy", []) if isinstance(item, dict)]
    by_key = {item.get("key"): item for item in diy if item.get("key")}
    for key, value in values.items():
        item = by_key.get(key)
        if item is None:
            item = {"name": key, "type": 9, "key": key}
            diy.append(item)
        item["value"] = value
    style["diy"] = diy


def build_sample_table(source: dict) -> dict:
    table = copy.deepcopy(source)
    cell_id = "datapoint-template-card-grid"
    x, y, width, height = 16, 152, 1888, 896
    table.update(
        {
            "shape": TABLE_SHAPE,
            "id": cell_id,
            "x": x,
            "y": y,
            "width": width,
            "height": height,
            "zIndex": 8,
            "visible": True,
            "position": {"x": x, "y": y},
            "size": {"width": width, "height": height},
        }
    )
    data = table.setdefault("data", {})
    data.update({"editMode": False, "IsToolBox": False})
    detail = data.setdefault("detail", {})
    detail.update(
        {
            "type": TABLE_SHAPE,
            "identifier": cell_id,
            "name": "点位卡片列表示例",
            "active": [],
            "action": [],
            "dataBind": [],
        }
    )
    style = detail.setdefault("style", {})
    style["position"] = {"x": x, "y": y, "w": width, "h": height}
    style["visible"] = 1
    style["zIndex"] = 8
    style["borderColor"] = "#2fd5f2"
    style["borderWidth"] = 1
    style["borderRadius"] = 10
    point_names = [name for name, _ in SAMPLE_POINTS]
    point_units = [unit for _, unit in SAMPLE_POINTS]
    set_diy(
        style,
        {
            "rowSource": "navDatapoints",
            "columnHeaders": "实时值",
            "rowDeviceNames": "\n".join(point_names),
            "rowDeviceCodes": "\n".join(point_units),
            "rowBindings": ";".join(f"示例设备->{name}" for name in point_names),
            "ShowCount": 20,
            "themeName": "dark",
            "panelAccentColor": "#4ae6ff",
            "deviceIconAccent": "#52e8ff",
            "pointIconAccent": "#bca5ff",
            "AutoPageEnabled": 0,
            "AutoPageInterval": 5,
            "AutoPageResumeDelay": 60,
            "navTotalDatapoints": str(len(SAMPLE_POINTS)),
            "navDatapointPageIndex": "0",
            "navDatapointPageSize": "20",
            "navDatapointTotalPages": "1",
        },
    )
    animate = detail.setdefault("animate", {})
    animate.setdefault("selected", [])
    animate.setdefault("animateElement", [])
    animate.setdefault("animateList", [])
    return table


def normalize_device_card_table(source: dict) -> dict:
    table = copy.deepcopy(source)
    x, y, width, height = 16, 152, 1888, 896
    table.update(
        {
            "x": x,
            "y": y,
            "width": width,
            "height": height,
            "position": {"x": x, "y": y},
            "size": {"width": width, "height": height},
        }
    )
    detail = detail_of(table)
    style = detail.setdefault("style", {})
    style["position"] = {"x": x, "y": y, "w": width, "h": height}
    style["borderColor"] = "#2fd5f2"
    style["borderWidth"] = 1
    style["borderRadius"] = 10
    upsert_diy(
        style,
        {
            "rowSource": "navChildren",
            "ShowCount": 49,
            "themeName": "dark",
            "panelAccentColor": "#4ae6ff",
            "deviceIconAccent": "#52e8ff",
            "pointIconAccent": "#bca5ff",
            "AutoPageEnabled": 0,
            "AutoPageInterval": 5,
            "AutoPageResumeDelay": 60,
        },
    )
    detail["name"] = "设备卡片列表"
    return table


def build_text_cell(
    source: dict,
    cell_id: str,
    x: int,
    y: int,
    width: int,
    text: str,
    role: str,
    *,
    color: str = "#b9cce6",
    font_size: int = 13,
) -> dict:
    cell = copy.deepcopy(source)
    height = 30
    cell.update(
        {
            "shape": "view-svg-text",
            "id": cell_id,
            "x": x,
            "y": y,
            "width": width,
            "height": height,
            "zIndex": 10,
            "visible": True,
            "position": {"x": x, "y": y},
            "size": {"width": width, "height": height},
        }
    )
    data = cell.setdefault("data", {})
    data.update({"editMode": False, "IsToolBox": False})
    detail = data.setdefault("detail", {})
    detail.update(
        {
            "type": "view-svg-text",
            "identifier": cell_id,
            "name": text,
            "active": [],
            "action": [],
            "dataBind": [],
        }
    )
    style = detail.setdefault("style", {})
    style.update(
        {
            "position": {"x": x, "y": y, "w": width, "h": height},
            "visible": 1,
            "text": text,
            "foreColor": color,
            "fontSize": font_size,
            "backColor": "transparent",
            "zIndex": 10,
        }
    )
    set_diy(style, {"labelRole": role})
    animate = detail.setdefault("animate", {})
    animate.setdefault("selected", [])
    animate.setdefault("animateElement", [])
    animate.setdefault("animateList", [])
    return cell


def build_sample_chrome(source: dict) -> list[dict]:
    return [
        build_text_cell(source, "datapoint-template-back", 16, 72, 120, "‹ 返回上一级", "deviceListBack", color="#7ee8ff"),
        build_text_cell(source, "datapoint-template-device", 152, 72, 500, "设备：示例设备", "deviceInfoName", color="#dffaff", font_size=15),
        build_text_cell(source, "datapoint-template-meta", 668, 72, 770, "编号 sample-001 · 模型 sample-model · 状态 离线", "deviceInfoMeta", color="#7894aa", font_size=11),
        build_text_cell(source, "datapoint-template-prev", 1480, 72, 88, "‹ 上一页", "detailPagePrev"),
        build_text_cell(source, "datapoint-template-info", 1572, 72, 176, "第 1/1 页 · 共 20 个测点", "detailPageInfo"),
        build_text_cell(source, "datapoint-template-next", 1752, 72, 88, "下一页 ›", "detailPageNext"),
    ]


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--db", default=str(ROOT / "ism_server_user/data/db/ism.db"))
    parser.add_argument("--model", default=MODEL_DEFAULT)
    parser.add_argument("--backup-dir", default=str(ROOT / "backups"))
    parser.add_argument("--apply", action="store_true")
    args = parser.parse_args()

    db_path = Path(args.db).expanduser().resolve()
    conn = sqlite3.connect(str(db_path))
    conn.row_factory = sqlite3.Row
    point_page = conn.execute(
        """SELECT id,page_name,page_id,components
           FROM display_model_layer
           WHERE model_id=? AND template_kind='datapointList' AND deleted_at IS NULL
           ORDER BY id DESC LIMIT 1""",
        (args.model,),
    ).fetchone()
    device_page = conn.execute(
        """SELECT id,page_name,page_id,components
           FROM display_model_layer
           WHERE model_id=? AND template_kind='deviceList' AND deleted_at IS NULL
           ORDER BY id DESC LIMIT 1""",
        (args.model,),
    ).fetchone()
    if point_page is None or device_page is None:
        parser.error("找不到当前模型的 datapointList/deviceList 模板")

    point_components = decode_components(point_page["components"])
    device_components = decode_components(device_page["components"])
    source_table = next((cell for cell in device_components.get("cells", []) if is_table(cell)), None)
    if source_table is None:
        parser.error("deviceList 模板中没有可复用的实时卡片容器")

    old_cells = point_components.get("cells", [])
    header_cells = [copy.deepcopy(cell) for cell in old_cells if keep_common_header(cell)]
    text_source = next(
        (cell for cell in header_cells if cell.get("shape") == "view-svg-text" and style_of(cell).get("text")),
        None,
    )
    if text_source is None:
        parser.error("点位模板中没有可复用的文字组件")
    point_components["cells"] = (
        header_cells
        + build_sample_chrome(text_source)
        + [build_sample_table(source_table)]
    )
    old_device_cells = device_components.get("cells", [])
    device_components["cells"] = [
        normalize_device_card_table(cell) if is_table(cell) else cell
        for cell in old_device_cells
    ]
    print(
        f"{point_page['page_name']} ({point_page['page_id']}): "
        f"{len(old_cells)} cells -> {len(point_components['cells'])} cells"
    )
    print(
        f"{device_page['page_name']} ({device_page['page_id']}): "
        f"{len(old_device_cells)} cells，设备表已切换为卡片矩阵"
    )
    if not args.apply:
        print("dry-run 完成；加 --apply 才会备份并写库。")
        conn.close()
        return 0

    backup_dir = Path(args.backup_dir).expanduser().resolve()
    backup_dir.mkdir(parents=True, exist_ok=True)
    stamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    backup_path = backup_dir / f"{db_path.stem}-before-datapoint-template-{stamp}{db_path.suffix}"
    shutil.copy2(db_path, backup_path)
    print(f"备份: {backup_path}")

    now = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
    try:
        conn.execute("BEGIN IMMEDIATE")
        conn.execute(
            "UPDATE display_model_layer SET components=?,updated_at=? WHERE id=?",
            (encode_components(point_components), now, point_page["id"]),
        )
        conn.execute(
            "UPDATE display_model_layer SET components=?,updated_at=? WHERE id=?",
            (encode_components(device_components), now, device_page["id"]),
        )
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()
    print("设备/点位列表模板已更新。")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

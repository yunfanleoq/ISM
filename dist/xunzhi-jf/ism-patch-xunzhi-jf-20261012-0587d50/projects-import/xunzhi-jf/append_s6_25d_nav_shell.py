#!/usr/bin/env python3
"""S6：2.5D 导航壳（占位页草稿 + 平面页入口）；非产品级 2.5D 引擎。"""
import json
import uuid
from pathlib import Path

PKG_PATH = Path(__file__).resolve().parent / "xunzhi-jf_complete_project_package.json"
DISPLAY_UID = "a1b2c3d4e5f6478990abcdef12345678"
PAGE_PLANE = "b2c3d4e5f6a7488990abcdef123456789"
PAGE_25D = "c3d4e5f6a7b84990abcdef25d000001"
SHELL_JSON = Path(__file__).resolve().parent / "xunzhi-jf_25d_shell_page.json"
META_KEY = "s6_25d_nav_shell_v1"


def make_animate():
    return {
        "selected": [],
        "animateElement": [],
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
    }


def nav_button_cell():
    cell_id = uuid.uuid4().hex
    return {
        "shape": "view-svg-text",
        "id": cell_id,
        "x": 1520,
        "y": 40,
        "width": 320,
        "height": 48,
        "zIndex": 25,
        "visible": True,
        "position": {"x": 1520, "y": 40},
        "size": {"width": 320, "height": 48},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "name": "nav-25d-shell",
                "identifier": cell_id,
                "style": {
                    "position": {"x": 1520, "y": 40, "w": 320, "h": 48},
                    "text": "2.5D 导航（壳）→ 占位页",
                    "visible": 1,
                    "fontSize": 15,
                    "foreColor": "#ffb74d",
                    "diy": [],
                },
                "animate": make_animate(),
                "action": [
                    {
                        "type": "click",
                        "action": "link",
                        "link": {
                            "linkType": "Inside",
                            "isPopUp": False,
                            "Inside": {
                                "displayUUID": DISPLAY_UID,
                                "pageUUID": PAGE_25D,
                                "displayType": 1,
                            },
                        },
                    }
                ],
                "active": [],
                "dataBind": [],
            }
        },
    }


def shell_page_doc():
    return {
        "pageId": PAGE_25D,
        "pageName": "xunzhi-jf-25d-shell",
        "route_hint": "/ism/display/render/{displayUid}?page=xunzhi-jf-25d-shell",
        "layer": {"width": 1920, "height": 1080, "backColor": "#0d1117"},
        "components": {
            "cells": [
                {
                    "shape": "view-svg-text",
                    "data": {
                        "detail": {
                            "type": "view-svg-text",
                            "name": "25d-placeholder",
                            "style": {
                                "position": {"x": 120, "y": 200, "w": 1680, "h": 120},
                                "text": "2.5D 场景占位 · 现场图/引擎未接入 · 请使用平面组态 xunzhi-jf-env-plane 兜底监视",
                                "visible": 1,
                                "fontSize": 22,
                                "foreColor": "#90caf9",
                                "diy": [],
                            },
                            "animate": make_animate(),
                        }
                    },
                }
            ]
        },
        "import_note": "ImportProject 当前仅写入单页 displayLayer；本页需组态编辑器手工创建或后续多页导入扩展。",
    }


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    if pkg.get("_meta", {}).get(META_KEY):
        print(f"Skip: {META_KEY}")
        return
    cells = pkg["displayLayer"]["components"].get("cells", [])
    if not any(c.get("data", {}).get("detail", {}).get("name") == "nav-25d-shell" for c in cells):
        cells.append(nav_button_cell())
        pkg["displayLayer"]["components"]["cells"] = cells
    shell = shell_page_doc()
    pkg.setdefault("draftPages", {})["xunzhi-jf-25d-shell"] = shell
    SHELL_JSON.write_text(json.dumps(shell, ensure_ascii=False, indent=2), encoding="utf-8")
    meta = pkg.get("_meta", {})
    meta[META_KEY] = True
    meta["s6_25d_page_id"] = PAGE_25D
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Updated nav shell + {SHELL_JSON}")


if __name__ == "__main__":
    main()

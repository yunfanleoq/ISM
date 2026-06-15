import os
import json
import re
from collections import Counter, defaultdict

def main(ctx):
    root = "/Users/yunfanleo/cursorProjects/ISM源码"
    obsidian_dir = os.path.join(root, "obsidian-notes")
    code_graph_dir = os.path.join(root, "code-graph")
    os.makedirs(obsidian_dir, exist_ok=True)
    os.makedirs(code_graph_dir, exist_ok=True)

    # ============ 前端文件扫描 ============
    fe_src = os.path.join(root, "ism-front-end", "src")
    fe_stats = {
        "total_files": 0,
        "total_lines": 0,
        "by_extension": Counter(),
        "by_directory": Counter(),
        "components": [],
        "pages": [],
        "services": [],
        "store_modules": []
    }
    
    for dirpath, dirnames, filenames in os.walk(fe_src):
        rel = os.path.relpath(dirpath, fe_src)
        for f in filenames:
            ext = os.path.splitext(f)[1].lower()
            if ext in ('.vue', '.js', '.less', '.css', '.json', '.html'):
                fe_stats["total_files"] += 1
                fe_stats["by_extension"][ext] += 1
                if rel != '.':
                    top_dir = rel.split(os.sep)[0]
                    fe_stats["by_directory"][top_dir] += 1
                fp = os.path.join(dirpath, f)
                try:
                    with open(fp, 'r', encoding='utf-8', errors='ignore') as fh:
                        lines = len(fh.readlines())
                        fe_stats["total_lines"] += lines
                except:
                    pass
                if ext == '.vue':
                    if 'components' in rel:
                        fe_stats["components"].append(os.path.join(rel, f))
                    elif 'pages' in rel:
                        fe_stats["pages"].append(os.path.join(rel, f))
                elif ext == '.js' and 'services' in rel:
                    fe_stats["services"].append(os.path.join(rel, f))
                elif ext == '.js' and 'store' in rel:
                    fe_stats["store_modules"].append(os.path.join(rel, f))

    # ============ 后端文件扫描 ============
    be_dir = os.path.join(root, "ism_server_user")
    be_stats = {
        "total_files": 0,
        "total_lines": 0,
        "by_directory": Counter(),
        "imports": Counter(),
        "protocols": []
    }
    
    for dirpath, dirnames, filenames in os.walk(be_dir):
        rel = os.path.relpath(dirpath, be_dir)
        if any(skip in rel for skip in ['vendor', 'upgrade_windows', 'dmgorm2']):
            continue
        for f in filenames:
            if f.endswith('.go'):
                be_stats["total_files"] += 1
                if rel != '.':
                    top_dir = rel.split(os.sep)[0]
                    be_stats["by_directory"][top_dir] += 1
                fp = os.path.join(dirpath, f)
                try:
                    with open(fp, 'r', encoding='utf-8', errors='ignore') as fh:
                        content = fh.read()
                        lines = content.count('\n') + 1
                        be_stats["total_lines"] += lines
                        for m in re.findall(r'"([^"]+)"', content):
                            if m.startswith('github.com') or m.startswith('gorm.io'):
                                parts = m.split('/')
                                key = parts[0] + '/' + parts[1] if len(parts) > 1 else m
                                be_stats["imports"][key] += 1
                except:
                    pass
                if 'protocol' in rel and f not in ('enter.go', 'common.go'):
                    proto = rel.split(os.sep)[1] if len(rel.split(os.sep)) > 1 else 'other'
                    be_stats["protocols"].append(proto)

    # save code graph
    with open(os.path.join(code_graph_dir, "code-graph.json"), 'w', encoding='utf-8') as f:
        json.dump({
            "frontend": {
                "total_files": fe_stats["total_files"],
                "total_lines": fe_stats["total_lines"],
                "by_extension": dict(fe_stats["by_extension"]),
                "by_directory": dict(fe_stats["by_directory"]),
                "component_count": len(fe_stats["components"]),
                "page_count": len(fe_stats["pages"]),
                "service_count": len(fe_stats["services"]),
                "store_module_count": len(fe_stats["store_modules"])
            },
            "backend": {
                "total_files": be_stats["total_files"],
                "total_lines": be_stats["total_lines"],
                "by_directory": dict(be_stats["by_directory"]),
                "top_imports": dict(be_stats["imports"].most_common(30)),
                "protocols": list(set(be_stats["protocols"]))
            }
        }, f, ensure_ascii=False, indent=2)

    # ============ 生成 Obsidian 笔记 ============
    # 1. Home / Index
    index_md = """---
tags: [ISM, SCADA, 代码分析, 组态软件]
date: 2026-06-12
---

# ISM Web组态软件 源码全景笔记

> 本笔记由 AI 代码分析工具自动生成，基于 ISM 源码静态扫描。  
> 项目路径：`/Users/yunfanleo/cursorProjects/ISM源码`  
> 版本：V3.01.RC07

## 项目概述

ISM Web组态软件是一款完整的 **工业物联网 Web 组态 SCADA 平台**，支持 15+ 种工业协议采集、实时组态可视化、视频联动、告警通知、历史数据分析。

- **前端**：Vue 2 + Ant Design Vue，基于 vue-antd-admin 二次开发
- **后端**：Go + Beego v2 + GORM，支持 6 种数据库
- **核心特色**：自研低代码组态编辑器 + 运行时渲染引擎

## 快速导航

| 模块 | 笔记 | 说明 |
|------|------|------|
| [[01-项目概览]] | 项目概述 | 整体架构、技术栈、版本信息 |
| [[02-前端架构]] | 前端架构 | Vue 项目结构、组件、路由、Store、Services |
| [[03-后端架构]] | 后端架构 | Go 项目结构、Controllers、Models、Tasks |
| [[04-工业协议]] | 工业协议 | 10+ 种协议实现详解 |
| [[05-数据模型]] | 数据模型 | 核心表结构、数据流转 |
| [[06-组态系统]] | 组态系统 | 低代码编辑器、组件库、渲染引擎 |
| [[07-启动指南]] | 启动指南 | 环境初始化、前后端启动步骤 |
| [[08-代码统计]] | 代码统计 | 文件数量、代码行数、依赖分析 |

## 关键数据速览

- 前端文件：""" + str(fe_stats['total_files']) + """ 个，约 """ + str(fe_stats['total_lines']) + """ 行代码
- 后端文件：""" + str(be_stats['total_files']) + """ 个 Go 文件，约 """ + str(be_stats['total_lines']) + """ 行代码
- 工业协议：15+ 种（Modbus, OPC UA, S7, SNMP, DLT645, IEC104, IEC61850, BACnet, MQTT, HJ212, CJT188, RESTFul, GB28181, 视频流, ISM组网）
- 数据库：SQLite / MySQL / PostgreSQL / 达梦(DM) / SQL Server / ClickHouse
- 时序数据库：InfluxDB / TDengine / LevelDB

---
*Generated by AI Code Analysis*
"""
    with open(os.path.join(obsidian_dir, "00-首页.md"), 'w', encoding='utf-8') as f:
        f.write(index_md)

    return {
        "obsidian_dir": obsidian_dir,
        "code_graph_dir": code_graph_dir,
        "frontend_files": fe_stats['total_files'],
        "frontend_lines": fe_stats['total_lines'],
        "backend_files": be_stats['total_files'],
        "backend_lines": be_stats['total_lines'],
        "notes_created": 1
    }

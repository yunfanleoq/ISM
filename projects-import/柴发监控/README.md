# 柴发监控项目导入包

> **重要（2026-07-14）**：柴发与中航信/循安部署在 **两套独立物理服务器** 上。  
> 客户柴发机请直接使用完整发版包（前后端+库+大屏）：  
> `releases/ism-release-sqlite-chaifa-*.zip`  
> 本目录 JSON 仅作「已有同机 ISM 上再导一个项目」的备选，**不是**柴发机的主交付物。

从 `Sqlite3_Backup_2026-07-13_10-37-19` 提取的项目导入 JSON（备选）。

## 包内容统计

| 项 | 数量 |
|---|---|
| 项目名称 | 后沙峪改造-柴发部分 |
| 设备模型 | 16 |
| 寄存器组 | 176 |
| 数据点 | 12959 |
| 设备实例 | 16（柴发楼 1~4 层各端口） |
| 组态大屏 | 无（原备份无 display） |
| 告警触发器 | 无 |

设备示例：`柴发楼1层_4001端口` … `柴发楼4层_505端口`，Modbus TCP 地址保留备份中的现场 IP/端口（如 `172.31.97.8:4001`）。

## 文件清单

| 文件 | 说明 |
|---|---|
| `后沙峪改造-柴发部分_ISM项目包.json` | 正式导入包（约 6.2MB），POST `/ImportProject` |
| `import_chaifa_project.py` | 一键导入脚本（登录 + 调用导入 API） |
| `README.md` | 本说明 |

## 客户导入步骤（推荐）

前提：ISM 后端已在运行（默认 `http://127.0.0.1:8081`），管理员账号可用（默认 `admin` / `123456`）。

```bash
# 进入本目录
cd 柴发监控

# 一键导入（新建独立项目）
python3 import_chaifa_project.py \
  --base-url http://127.0.0.1:8081 \
  --user admin \
  --password 123456
```

成功后前端「项目列表」应能同时看到：

1. 原有项目（如循安 / 配电室等）
2. **后沙峪改造-柴发部分**（新建）

两个项目互不共享设备树与数据模型。

### 手动 curl 导入

```bash
# 1) 登录拿 token（密码为 MD5，123456 → e10adc3949ba59abbe56e057f20f883e）
TOKEN=$(curl -s -X POST http://127.0.0.1:8081/login \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

# 2) 导入（勿加 Bearer 前缀；数据点多，建议超时 ≥ 30 分钟）
curl -sS -X POST http://127.0.0.1:8081/ImportProject \
  -H "Authorization: $TOKEN" \
  -H 'Content-Type: application/json' \
  --data-binary @后沙峪改造-柴发部分_ISM项目包.json \
  --max-time 1800
```

期望响应：`code: 0`，并返回新建 `project_uuid`。

## 注意事项

1. **仅管理员**可调用 `/ImportProject`。
2. 导入会 **新建项目**，不会覆盖已有项目。
3. 设备通信 IP/端口来自原现场备份；客户现场若网段不同，导入后需在设备「扩展参数」中改 IP/端口。
4. 本包不含组态大屏；需要大屏可导入后在编辑器中新建，或另行交付大屏包。
5. 若项目列表看不到新项目：确认用 `admin`（user 表）登录；必要时调用 `POST /ProjectFixCreator` 修正归属。

## 重新从备份生成

开发侧若需从原始 zip 再生成：

```bash
python3 scripts/sqlite_backup_to_project_package.py \
  --sql Sqlite3_Backup_2026-07-13_10-37-19.zip \
  --out projects-import/柴发监控/后沙峪改造-柴发部分_ISM项目包.json \
  --project-name "后沙峪改造-柴发部分"
```

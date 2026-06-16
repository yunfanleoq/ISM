# 云端 / 离线 SQLite 部署

## 数据库文件

- 路径：`ism_server_user/data/db/ism.db`（相对后端工作目录为 `data/db/ism.db`）
- 由 OceanBase/MySQL 导出：`python3 scripts/export_db_to_sqlite.py`
- 历史大表默认**不导数据**（空表结构）：`devices_alarm_list`、`devices_history_data_list`、`system_journal`

## 后端配置

在 `ism_server_user/conf/app.conf` 中：

```ini
dbtype=1
```

`dbtype=1` 使用 GORM SQLite 驱动，连接串为 `data/db/ism.db`。请在 `ism_server_user` 目录下启动 `ism_server`，保证相对路径正确。

## 启动顺序

1. 启动后端：`cd ism_server_user && ./ism_server`（或等价二进制）
2. 启动前端 dev：`cd ism-front-end-v2 && npx vue-cli-service serve --port 7080`（需 `/api` 代理，见 `vue.config.js`）
3. 默认账号：`admin` / `123456`（前端 MD5 后提交，库内为 bcrypt(MD5)）

## 重新导出

本地仍使用 OceanBase（`dbtype=4`）时，可随时运行导出脚本更新仓库内快照库后再 commit。

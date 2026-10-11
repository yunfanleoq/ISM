# S2 骨架 · dcim-asset（契约 v2 · 不阻塞 W1/W2）

> U 位全生命周期属同窗全量规划；本迭代交付 **Mock 数据 + JSON 存储 + 抽样脚本 + Go 只读/写占位 API**，无 SpringBoot、无正式 `dcim_rack_u` 表迁移。

## 实体（规划）

```
Rack(rack_id) 1─* UPosition(u_index) 1─0..1 Asset(asset_id)
```

| 字段 | 说明 |
|------|------|
| rack_id | 关联 **MonitorList** 机房/机柜节点 `uuid`（与 ImportProject `monitorTree` 一致） |
| u_index | 1..42 |
| asset_id | 可选，上架后绑定 |
| status | `empty` / `occupied` / `reserved` |

### 与 MonitorList 挂接

| MonitorList 节点 | uuid（寻智 W1 包） | 用途 |
|------------------|-------------------|------|
| `F3-模块机房-01` | `88888888888888888888888888888888` | **rack_id**；U 位图、容量统计挂此节点 |
| 设备子节点（如 `xunzhi-jf-TH-01`） | 各设备 uuid | 监测采集；**不**直接写 U 位占用 |

组态/监视树选中机柜节点时，前端后续可调用 `GET /api/xunzhi-jf/asset/rack/{rack_id}/u-positions`（经 `vue.config.js` 代理到后端 `/xunzhi-jf/...`）拉 U 位占用；设备实时值仍走 WebSocket RealData。

契约表名（正式版）：`dcim_asset`、`dcim_rack_u`、`dcim_asset_ticket` — 见 `docs/ISM-寻智机房环境监测-增量架构契约.md`。

## 仓内 Mock 资产

| 路径 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/xunzhi-jf_dcim_asset_mock.json` | 种子：机柜 A01、资产 `as-xz-0001` |
| `data/xunzhi-jf/dcim_asset_store.json` | 运行态（`init` / `demo` 生成，勿提交敏感数据） |

## 抽样路径（上架 → 占用 → 下架）

```bash
# 一步自证（从种子重置 + 绑 U10-11 + 断言 + 下架；无需后端）
bash scripts/xunzhi_jf_asset_u_demo.sh
# 期望末尾：DEMO PASS · ASSERT OK（U10-U11 occupied / empty）

# 或分步
python3 scripts/xunzhi_jf_asset_store.py init
python3 scripts/xunzhi_jf_asset_store.py bind --asset-id as-xz-0001 --u-start 10 --u-height 2
python3 scripts/xunzhi_jf_asset_store.py list
python3 scripts/xunzhi_jf_asset_store.py unbind --asset-id as-xz-0001
```

后端已起且 store 已 init 时（需登录 Token，与其它 API 一致）：

```bash
RACK=88888888888888888888888888888888
curl -s "http://127.0.0.1:8081/xunzhi-jf/asset/rack/${RACK}/u-positions" -H "Authorization: <token>"
curl -s -X POST http://127.0.0.1:8081/xunzhi-jf/asset/u-position/bind \
  -H 'Content-Type: application/json' -H "Authorization: <token>" \
  -d '{"asset_id":"as-xz-0001","rack_id":"'"${RACK}"'","u_start":10,"u_height":2}'
```

## HTTP 占位（Go · Beego）

| 方法 | 路径（前端 `/api` 前缀） | 说明 |
|------|--------------------------|------|
| GET | `/api/xunzhi-jf/asset/rack/:rackId/u-positions` | U 位列表 |
| POST | `/api/xunzhi-jf/asset/u-position/bind` | 上架绑定 |
| POST | `/api/xunzhi-jf/asset/u-position/unbind` | 下架 |

实现：`ism_server_user/controllers/xunzhiJfAssetCtl.go`、`models/xunzhiJfAsset.go`。

**禁止**：SpringBoot 平行系统；绕过 `dcim-asset` 在采集侧双写 U 位（契约依赖方向）。

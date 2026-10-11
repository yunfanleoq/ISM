# 寻智机房环境 · xunzhi-jf（S0 竖切 + W2 扩展）

Modbus TCP：**温湿度 + 漏水**（W1）+ **列头柜 + UPS**（W2 点表已并入同一导入包）。

链路：**协议 → DeviceRealData → WebSocket（`conf/app.conf` 的 `wsport`，本仓 dev 常见 10225）→ 组态**。

## 测试复测入口（推荐）

```bash
./scripts/xunzhi_jf_smoke_env.sh
./scripts/xunzhi_jf_collect_evidence.sh
bash scripts/xunzhi_jf_alarm_demo.sh
```

告警：导入包含 **`alarmTriggers`**（漏水 `=`1、温度 `>`35）及告警点 **`alarmOnValue: 1`**。项目已存在时用 `scripts/xunzhi_jf_sync_alarm_db.py`，勿仅重复 ImportProject（会 4001）。

证据目录：`docs/ISM-寻智-W1复测证据/` · 试打样包：`dist/xunzhi-jf/ism-patch-xunzhi-jf-*`

### WebSocket 端口

| 说明 | 端口 |
|------|------|
| 契约/前端默认示例 | 10215 |
| **本仓 `ism_server_user/conf/app.conf`** | **10225（联调以实测为准）** |

业主本机冒烟：`:10225` 可连 RealData，`:10215` 拒连。文档已对齐，**未改现网配置**。

## 资产路径

| 文件 | 说明 |
|------|------|
| `xunzhi-jf_complete_project_package.json` | `POST /ImportProject`（含 W2 设备时需重新导入或新库） |
| `xunzhi-jf_device_points.json` | W1 点表 |
| `xunzhi-jf_w2_device_points.json` | W2 配电/UPS 附录 |
| `generate_xunzhi_jf_package.py` | 生成 W1 基线包 |
| `append_w2_power_ups.py` | 向包追加 PWR/UPS（slave 3/4） |

组态：**寻智机房环境平面监视-xunzhi-jf** / 页 **`xunzhi-jf-env-plane`**

## 账号与端口

| 服务 | 默认 |
|------|------|
| HTTP API | `http://127.0.0.1:8081` |
| 前端 dev | `http://127.0.0.1:7080` |
| Mock Modbus | `127.0.0.1:1502` |
| 登录 | `admin` / `123456`（前端 MD5 后传 `e10adc3949ba59abbe56e057f20f883e`） |

## 禁止

不覆盖麒麟生产、`patches/ism-*-kylin`、`releases/` 现网包。

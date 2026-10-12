# W2 竖切 · 配电/UPS 采集 + 平面可演示

> **范围**：试演示竖切（点表 + Mock slave 3/4 + 平面实时表/表盘），非完整配电产品。S0–S6 骨架不变。

## 仓内资产

| 路径 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/append_w2_power_ups.py` | Import 包：列头柜 + UPS 物模型/设备（slave 3/4） |
| `projects-import/xunzhi-jf/append_w2_power_viz.py` | 平面组态：W2 实时表 + 2 个表盘 |
| `projects-import/xunzhi-jf/reflow_xunzhi_jf_env_plane_layout.py` | 演示平面收口（含 W2 区块） |
| `projects-import/xunzhi-jf/xunzhi-jf_w2_device_points.json` | W2 点表附录 |
| `scripts/xunzhi_jf_modbus_mock.py` | Mock HR3/4（ABCD 原码） |
| `scripts/xunzhi_jf_power_demo.sh` | W2 探针脚本（P-COL-02） |

**字节序**：物模型 `DataFormat` / 点表 `ByteOrder` = **ABCD**（与 Mock 一致）。

## 生成导入包（在已有 W1 管线中插入 W2 viz）

```bash
python3 projects-import/xunzhi-jf/append_w2_power_ups.py   # 若包已含 W2 可跳过
python3 projects-import/xunzhi-jf/append_w2_power_viz.py
python3 projects-import/xunzhi-jf/reflow_xunzhi_jf_env_plane_layout.py
python3 scripts/import_xunzhi_jf_project.py
```

## 本机双口演示（7081 / 8082 / 1503）

1. Mock：`python3 scripts/xunzhi_jf_modbus_mock.py`（或 `XUNZHI_MOCK_PORT=1503`）
2. 后端：`ism_server_user` 监听 **8082**（见 `docs/ISM-寻智-本机双开端口.md`）
3. 前端 dev **7081**，登录 `admin` / `123456`（前端 MD5）
4. 重导项目后打开 **`/#/xunzhi/jf/viz/plane`**（`noRunTree=1` 平面闸）
5. 目视 **③ W2** 区：`xunzhi-jf-PWR-01` / `xunzhi-jf-UPS-01` 实时表 + 总有功功率 / 电池电压表盘

探针：

```bash
bash scripts/xunzhi_jf_power_demo.sh
# 期望后端日志含 cycle_log_power_ups / probe=P-COL-02
```

## 禁止

- SpringBoot、`patches/*kylin`、`releases/`、`/tmp/keep_xunzhi_be.sh`

*文档随竖切迭代，未宣称正式验收。*

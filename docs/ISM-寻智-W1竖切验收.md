# 寻智机房环境 · W1 竖切验收（温湿度 + 漏水）

> 分支：`feature/xunzhi-jf-env-monitor` · 探针：`P-COL-01 cycle_log_th_leak`  
> 契约：`docs/ISM-寻智机房环境监测-增量架构契约.md`（同步自项目经理材料）

## 交付物路径

- 导入包：`projects-import/xunzhi-jf/xunzhi-jf_complete_project_package.json`
- 点表：`projects-import/xunzhi-jf/xunzhi-jf_device_points.json`
- Mock：`scripts/xunzhi_jf_modbus_mock.py`
- 组态页名：`xunzhi-jf-env-plane`（模型名：寻智机房环境平面监视-xunzhi-jf）

## 调用链自证

1. Mock/现场 Modbus → `ism_server` Modbus 采集
2. 写入 `DeviceRealData` + `cycle_log_th_leak` 日志
3. WebSocket `:10215` 推送 `RealData`
4. 组态 `ism-view-real-table` 显示 `设备->测点` 绑定值

## 操作步骤

见 `projects-import/xunzhi-jf/README.md`。

### 周期日志（P-COL-01）

```bash
grep 'cycle_log_th_leak' /tmp/ism_be.log | tail -3
```

期望字段：`gather_interval_ms=5000`，`cycle_ms` 远小于 `10000`，`points=3`（每设备一轮）。

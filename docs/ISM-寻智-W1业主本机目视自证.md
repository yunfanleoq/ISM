# W1 业主本机目视自证（前端 7080）

> 路径：`~/cursorProjects/ISM源码` · 分支 `feature/xunzhi-jf-env-monitor`  
> **非对外验收结论**；与 API/WS 证据一并归档。

## 前置

```bash
cd ~/cursorProjects/ISM源码
git fetch && git checkout feature/xunzhi-jf-env-monitor && git pull
./scripts/xunzhi_jf_smoke_env.sh --with-fe   # Mock + 后端 + 导入 + 前端
# 或业主已验证流程：Mock + CGO_ENABLED=1 编译 ism_server + import
```

## 目视步骤（F3.1 / F3.3）

1. 浏览器打开 `http://127.0.0.1:7080`，登录 `admin` / `123456`
2. 切换项目 **寻智机房环境-xunzhi-jf**
3. **项目树**：展开 **F3-模块机房-01**，可见 `xunzhi-jf-TH-01`、`xunzhi-jf-LEAK-01`（F3.1）；W2 导入后另见 `xunzhi-jf-PWR-01`、`xunzhi-jf-UPS-01`
4. 打开组态 **寻智机房环境平面监视-xunzhi-jf** → 页 **`xunzhi-jf-env-plane`**
5. 确认表格温湿度/漏水列数值随 Mock 变化刷新（F3.3）
6. 截图或录屏保存到 `docs/ISM-寻智-W1复测证据/` 或 `artifacts/`

## WebSocket（I3）

- 端口以 **`ism_server_user/conf/app.conf` 的 `wsport`** 为准（业主冒烟 **10225**，非契约示例 10215）
- DevTools → Network → WS → 筛选 `RealData` 帧；或运行：

```bash
./scripts/xunzhi_jf_collect_evidence.sh
```

## 告警闭环（关阻塞 · 与目视独立）

```bash
python3 scripts/xunzhi_jf_sync_alarm_db.py   # 项目已存在、未整包重导时先跑
python3 scripts/xunzhi_jf_modbus_mock.py --leak-on &   # 或 smoke 已起 Mock 则重启带 --leak-on
bash scripts/xunzhi_jf_alarm_demo.sh
```

整包重导：删项目后 `python3 scripts/import_xunzhi_jf_project.py`（导入包含 `alarmTriggers` 与 `alarmOnValue`）。

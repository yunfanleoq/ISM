# S1 envmon 可演示面（告警 + 多机房树草稿）

> 分支 `feature/xunzhi-jf-env-monitor` · PR #5 · **非正式验收**  
> 业主路径：`~/cursorProjects/ISM源码`

## 交付物

| 能力 | 入口 | API |
|------|------|-----|
| **envmon-alarm** | 组态页 `xunzhi-jf-env-plane` 底部 **`view-device-alarm-list`** | `GetAlarmEventFeed`（四级着色） |
| **告警确认** | 菜单 **告警管理 → 实时告警** | `POST /AlarmOpt`（`UpdateCurrentAlarm` type=1） |
| **envmon-viz** | 项目树 **F3-模块机房-02-Mock** + 组态标题条 | `POST /monitortree` |
| **脚本闭环** | `bash scripts/xunzhi_jf_alarm_demo.sh` | 与目视独立 |

导入包由 `generate_xunzhi_jf_package.py` → `append_w2_power_ups.py` → **`append_s1_envmon_demo.py`** 生成。

## 前置

```bash
git pull origin feature/xunzhi-jf-env-monitor
# 新库或删项目后：
python3 projects-import/xunzhi-jf/generate_xunzhi_jf_package.py
python3 projects-import/xunzhi-jf/append_w2_power_ups.py
python3 projects-import/xunzhi-jf/append_s1_envmon_demo.py
python3 scripts/import_xunzhi_jf_project.py
# 已存在项目：仅 S1 组态/树需重导或组态编辑器手工挂告警列表组件
```

起 Mock + 后端 + 前端（7080），见 `docs/ISM-寻智-W1业主本机目视自证.md`。

## 演示步骤 A · envmon-alarm（组态）

1. 登录 `http://127.0.0.1:7080`，项目 **寻智机房环境-xunzhi-jf**
2. 打开大屏 **寻智机房环境平面监视-xunzhi-jf** → 页 **`xunzhi-jf-env-plane`**
3. 页面下方 **实时告警滚动表**（`envmon-alarm`）：列含设备名、告警名、时间、**等级（提示/次要/重要/严重/致命）**、状态、内容
4. 造数：`python3 scripts/xunzhi_jf_modbus_mock.py --leak-on`（或 `xunzhi_jf_alarm_demo.sh` 前置 sync）
5. 等待数秒，列表应出现 **漏水** 相关行（依赖告警入列修复 + `xunzhi_jf_sync_alarm_db.py`）

## 演示步骤 B · 四级 + 确认（平台页）

1. 左侧 **告警管理 → 实时告警**（路由 `currentAlarm`）
2. 筛选当前项目，列表展示 **AlarmLevel** 四级样式
3. 选中一条 → **确认/清除**（与 `AlarmOpt` 一致；与 `xunzhi_jf_alarm_demo.sh` 相同 API）
4. 返回组态告警表，状态应变为已恢复或列表减少（视 `GetAlarmEventFeed` 事件流）

## 演示步骤 C · envmon-viz（多机房树）

1. **项目树** 展开：根 → **F3-模块机房-01**（TH/LEAK 实装）与 **F3-模块机房-02-Mock**（草稿）
2. **02-Mock** 下 **xunzhi-jf-TH-02-Mock** 为占位节点（`IsEnable=0`，不参与采集）
3. 组态顶部蓝字说明条对照两机房命名

## W2 探针（同文档链）

重导含 W2 设备后：`grep cycle_log_power_ups` 后端日志见 P-COL-02（见 `docs/ISM-寻智-W1竖切验收.md`）。

## 禁止

不覆盖 `patches/*kylin`、`releases/`、现网麒麟包。

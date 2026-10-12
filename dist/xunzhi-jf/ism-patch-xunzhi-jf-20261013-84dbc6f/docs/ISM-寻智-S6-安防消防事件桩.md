# S6 桩 · 安防 / 消防事件入列（Mock · 非完整平台）

> **诚实口径**：演示 **消防 + 安防** 两类事件写入内存队列、 dequeue 一条进「已处理」；**无** 视频联动、**无** 国标协议接入、**无** 与 ISM 原生告警引擎合并。**未宣称** 完整安防消防平台。

## 场景

| type | 样例 code | 说明 |
|------|-----------|------|
| `fire` | `SMOKE_ZONE_A` | 烟感（Mock） |
| `security` | `DOOR_FORCED` | 门禁异常（Mock） |

与 MonitorList 机房节点 `88888888888888888888888888888888` 关联，便于后续组态/工单挂接。

## 路径

| 文件 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/xunzhi-jf_security_event_mock.json` | 空队列种子 |
| `scripts/xunzhi_jf_security_event_store.py` | ingest / dequeue / **demo** |
| `scripts/xunzhi_jf_security_demo.sh` | 一步自证 |
| `data/xunzhi-jf/security_event_store.json` | 运行态（gitignore） |

## 自证

```bash
bash scripts/xunzhi_jf_security_demo.sh
```

### 期望输出（关键行）

```
ASSERT OK: 2 events queued (fire + security)
ASSERT OK: 1 processed, 1 remaining in queue
DEMO PASS · Store: .../data/xunzhi-jf/security_event_store.json
```

后续可扩展：事件 → S4 `create-from-alarm` 风格工单（**未实现**）。

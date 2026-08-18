# ISM 现场验收清单 — 实时通道 / 大屏刷新 / 历史存储

> 适用：OceanBase 正式包（如 `ism-release-oceanbase-20260714-1746-3589`）  
> 目标：5～10 分钟内判定「通道丢包、大屏刷新、历史时区与点数」三项是否正常。  
> 默认端口：前端 **7090**、后端 **8091**、OceanBase **2881**、TDengine **6041**。

---

## 0. 前置

```bash
cd /opt/ism/ism-release-oceanbase-*/   # 以实际目录为准
grep -E '^dbtype=|^httpport=|^realdatachanelcache=' ism_server_user/conf/app.conf
grep -E '^historyrecorddbtype=|^oncewrite' ism_server_user/conf/historyData.conf
# 期望: dbtype=4  httpport=8091  realdatachanelcache>=20000  historyrecorddbtype=2
```

登录：`admin` / `123456`。密码链路为 MD5 后 bcrypt，验收接口用 MD5(`123456`)=`e10adc3949ba59abbe56e057f20f883e`。

---

## 1. 实时通道是否拥堵丢包（问题 2）

```bash
# 最近是否在丢
grep -c 'RealDataChanel full' logs/ism_server.log 2>/dev/null || \
  grep -c 'RealDataChanel full' ism_server_user/logs/ism.log 2>/dev/null

# 看最近几条（应很少或仅有节流汇总）
grep 'RealDataChanel full' logs/ism_server.log 2>/dev/null | tail -5
```

| 结果 | 判定 |
|------|------|
| 启服 5 分钟内 count≈0 或仅偶发 | **通过** |
| 持续刷屏 `full, drop` | **未通过** → 先清告警再加大缓存 |

应急：

```bash
# 清活跃告警（减轻推送风暴）
python3 scripts/clear_all_alarms.py --dry-run
python3 scripts/clear_all_alarms.py

# 可选：加大通道（改完重启后端）
# realdatachanelcache=50000
```

---

## 2. 大屏是否「采集到才更新」（问题 1）

期望：卡片数值随采集/WS 推送更新，**不是**整页无规律翻页乱滚。

1. 浏览器打开 `http://<IP>:7090/#/login`，进入大屏。
2. 选一个已知在变的测点（或 Modbus 模拟值），盯住 **同一卡片** 30～60 秒。
3. 同时看后端是否还在 `RealDataChanel full`（与 §1 联动）。

| 现象 | 判定 |
|------|------|
| 同一卡片值随采集变化、页码不自动跳 | **通过** |
| 整页测点像轮播/乱跳，且 §1 在丢包 | 先修通道 |
| 整页自动翻页、页码自己变 | 确认前端为新 dist（`ViewRealTable` 已强制关自动翻页）；硬刷新 Ctrl+F5 |

后端登录验活：

```bash
curl -s -X POST "http://127.0.0.1:8091/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
# 期望含 "code":1000
```

前端代理：

```bash
curl -s -o /dev/null -w '%{http_code}\n' "http://127.0.0.1:7090/"
# 期望 200
```

---

## 3. 历史存储：时区 + 点数（问题 3）

### 3.1 系统时区

```bash
date
timedatectl 2>/dev/null | head -8
# 期望 Asia/Shanghai 或 CST（UTC+8）
```

### 3.2 是否开了历史记录

OceanBase：

```bash
# 包内 python 或系统 python3 + pymysql
python3 - <<'PY'
import configparser
from pathlib import Path
try:
    import pymysql
except ImportError:
    raise SystemExit('pip3 install pymysql')
conf_path = Path('ism_server_user/conf/app.conf')
text = '[app]\n' + conf_path.read_text(encoding='utf-8', errors='ignore')
p = configparser.ConfigParser(); p.optionxform=str; p.read_string(text)
c = {k.lower(): v for k,v in p.items('app')}
conn = pymysql.connect(host=c.get('oceanbasehost','127.0.0.1'),
    port=int(c.get('oceanbaseport','2881')),
    user=c.get('oceanbaseuser','root'),
    password=c.get('oceanbasepwd',''),
    database=c.get('oceanbasedbname','ism'), charset='utf8mb4')
cur = conn.cursor()
cur.execute('SELECT is_record, COUNT(*) FROM device_real_data WHERE deleted_at IS NULL GROUP BY is_record')
print('device_real_data is_record:', cur.fetchall())
conn.close()
PY
```

| 结果 | 判定 |
|------|------|
| `is_record=1` 行数 ≈ 0 | **未通过**（几乎不会入库）→ 用下方脚本开启 |
| `is_record=1` 覆盖业务需要的测点 | 继续查 TDengine |

批量开启（示例：列头柜相关，60 秒间隔）：

```bash
python3 scripts/enable_history_record.py --dbtype 4 \
  --device-like '%列头%' --interval 60 --dry-run

python3 scripts/enable_history_record.py --dbtype 4 \
  --device-like '%列头%' --interval 60
# 然后重启 ism_server
```

### 3.3 TDengine 是否在写、时间是否偏 8 小时

```bash
# REST 探活
curl -s -u root:taosdata -d 'show databases;' http://127.0.0.1:6041/rest/sql

# 抽样最近记录（库名/表名以现场为准；稳定名为 ISMHistoryDb）
curl -s -u root:taosdata \
  -d "select last_row(*) from ISMHistoryDb.HistoryDatas;" \
  http://127.0.0.1:6041/rest/sql
```

| 现象 | 判定 |
|------|------|
| 有多条、多 data_name，时间≈当前北京时间（±1 分钟） | **通过** |
| 有数据但显示比墙钟快/慢约 8 小时 | 时区链路未对齐（确认二进制含 `FormatTDengineTimestamp` + 系统时区） |
| 长期只有 1 个 data_name / 几乎无新行 | 回到 3.2 开 `is_record`，并确认 `historyrecorddbtype=2` |

### 3.4 定时历史存档逻辑（实时库快照）

`RecordType=1`（定时存储）由后台任务 **到点从实时内存库取最新值** 写入历史，**不依赖**当轮采集是否成功。

- 任务：`DealWithTimedRealtimeHistorySnapshot`
- 存档时间：按 `record_interval`（秒）截断，如 600 → `10:50:00`、`11:00:00`
- 采集路径对 type=1 不再写历史（避免双写）

```bash
grep -E 'timed history snapshot|no realtime value' ism_server_user/logs/ism.log*
```

### 3.5 缺档诊断日志（`loglevel=3` 可见）

路径：`ism_server_user/logs/ism.log`（非 Docker / 非 `tdengine/log`）。

```bash
grep -E 'timed history snapshot|history archive gap|history queue full|history pipeline pressure|write history data batch failed|history (change|pct)-store skip' \
  ism_server_user/logs/ism.log*
```

| 关键词 | 含义 |
|--------|------|
| `timed history snapshot` | 定时快照 tick / 点数加载 / 无实时值 |
| `history archive gap` | （旧采集驱动路径）间隔翻倍空档 |
| `history queue full` | 入队丢样（含点位名，按点节流） |
| `history pipeline pressure` | drops/flushFail > 0 |
| `write history data batch failed` | 写历史库失败 |
| `history change-store skip` / `history pct-store skip` | 变化/百分比存档解析失败 |

---

## 4. 一票否决（任一条失败则本轮验收不通过）

1. 启服后持续 `RealDataChanel full, drop`
2. 大屏页码自行翻页或数值因丢包无规律乱跳（§1 未过时）
3. `is_record=1` 为 0 却要求看历史曲线
4. `historyrecorddbtype` 不是 2，或 TDengine 6041 不通

---

## 5. 相关文件

| 文件 | 用途 |
|------|------|
| `docs/ISM-告警风暴与RealDataChanel满修复.md` | 通道拥堵根因与清告警 |
| `scripts/clear_all_alarms.py` | 清活跃告警 |
| `scripts/enable_history_record.py` | 批量开历史存储 |
| `ism_server_user/protocol/common/tdengine_time.go` | 历史写入 UTC 格式化 |

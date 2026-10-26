# 历史库备份：命令

只在要停备份、核对张数、或把容器目录落到宿主时打开。先读 `SKILL.md` 的禁止项。

## 停备份

只杀 taosdump，以及宿主机上包着它的 `docker exec`。不要 `docker stop` / `docker kill` 容器，不要杀 `taosd`，不要停 `ism_server`。

```bash
pgrep -af taosdump
pkill -f '[t]aosdump'
docker exec tdengine sh -c 'ps -ef | grep "[t]aosdump"'
docker exec tdengine sh -c 'kill $(pidof taosdump)'
```

容器名默认 `tdengine`，现场改过则用 `TD_CONTAINER`。容器里没有 `pidof` 时，用上一行 `ps` 里的 pid 执行 `kill <pid>`。`pkill` 的模式必须是 `taosdump`，不能是 `taos` 或 `taosd`。

进程被取消后，宿主还没有 `dbs.sql` 就留下容器 `/tmp/ISM_TDengine_Backup_*`，不要 `rm -rf`。

## 查子表（只读）

总数以 taosdump 日志为准：`child tables (119166)`。需要自己数时只读，不要把结果拿去删表：

```sql
SELECT COUNT(*) FROM information_schema.ins_tables
 WHERE db_name = 'ismhistorydb'
   AND stable_name = 'templetehistorydatas';
```

超级表是 `ISMHistoryDb.TempleteHistoryDatas`，子表名是 `hd_<测点uuid>`（uuid 里非字母数字会变成 `_`）。停录（`is_record=0`）的点留着。禁止 `DROP STABLE`。孤儿子表只在已有完整备份、且业务库确认该测点不存在时，由人逐张删。

## 进度

整数除法。总张数 `T`，已导 `D`。

| 阶段 | 百分比 | 界面 |
|---|---|---|
| 表结构 | `min(30, D * 30 / T)` | `正在备份中 表结构 D/T` |
| 数据 | `min(99, 30 + D * 69 / T)` | `正在备份中 数据 D/T` |
| 拷到宿主 | 99 | `正在备份中 正在拷到备份目录` |
| 文件已落盘且通过完成条件 | 100 | `备份完成` |

`792 * 69 / 119166 = 0`，所以 792/119166 仍是 30%。`1727 * 69 / 119166` 仍是 0，1728 张才到 31%。停在 30% 继续等。

满 5 分钟张数没变：只在文案后加「已超过 5 分钟没有新进度」，不杀进程。连续 30 分钟没有任何新输出，或从启动算满 48 小时，才停。

中断文案用拷贝前的进度。表结构：`导出中断，表结构还没导完（已看到 D/T 张）`。数据：`导出中断，数据 D/T`。不要写成「数据还没开始」。

## 流式 tar

导出过程中不拷（`DumpHostSyncInterval = 0`）。结束后只拷一次，上限 2 小时。avro 已是 snappy，用 `-cf`，不要 `-czf`，不要再接 gzip。

容器没有 `tar`（`docker exec <容器> tar --help` 失败）就失败，不要改成逐文件 `docker cp`。

`<名>` 两边必须相同，例如 `ISM_TDengine_Backup_2026-10-22_09-00-00`：

```bash
docker exec tdengine tar -C /tmp -cf - <名> \
  | tar -xf - -C data/hisdbbackup
```

对应代码参数：容器 `exec <容器> tar -C /tmp -cf - <名>`，宿主 `tar -xf - -C <备份目录的父目录>`。包内是 `<名>/dbs.sql`，解到父目录后，`dbs.sql` 在 `data/hisdbbackup/<名>/dbs.sql`。

不要：

- `tar -C /tmp/<名> .`：成员 `.` 在部分 tar 上会让解压以非 0 退出，文件已经落上也会被当成失败。
- 把包解进 `<名>/` 自己里面：会多套一层，`dbs.sql` 不在备份根上，检查认不出。
- 导出过程中每 30 秒 `docker cp` 十几万小文件。

宿主还没有 `dbs.sql` 时，容器 `/tmp/<名>` 是唯一副本，不要删。宿主已经有 `dbs.sql` 但数据没导完：写 `backup.incomplete`，保留宿主目录，这时才可以删容器副本。

## 怎样才叫备份完成

三者同时成立：

1. 备份目录根上有 `dbs.sql`（或其下一层 `taosdump-*` 目录里有，现有检查只认这一层）。
2. 数据阶段打到最后一张（`ReachedEnd`，阶段仍是 data）。表结构导完不算。
3. 目录里除了 `.sql` 和 `backup.incomplete`，还有大小大于 0 的数据文件。

只有 `dbs.sql`：写 `backup.incomplete`，列表显示「不完整」，目录留下。没有标记、且能看到 `dbs.sql` 的旧目录仍显示「备份完成」。同时只允许一个历史库备份；第二次直接拒绝。

备份目录：`data/hisdbbackup/ISM_TDengine_Backup_<2006-01-02_15-04-05>`。容器内是 `/tmp/` 加同一个名字。容器内 taosdump 走 `127.0.0.1:6030`，不要改成 REST 6041。

## 前端超时（不要改）

| 调用 | 超时 |
|---|---|
| `HisDbBackup` | 51 小时（盖住后端 48 小时 + 最多 2 小时拷贝） |
| `HisDbRestore` / `HisDbDownload` / `HisDbBackupUpload` | 6 小时 |
| `DbRestore` / `DbDown`（业务库） | 1 小时 |
| `DbBackup`（业务库） | 10 分钟 |

历史库的 48 小时不要套到业务库上。

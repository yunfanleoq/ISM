# ISM 20260819 麒麟补丁

对照现场《20260819问题项》5 条：导航渲染崩溃与保存文案、历史快照空值跳过/阻塞丢档、脚本 native 停更与手动执行无日志、虚拟设备全量导出缺列且导入不更新、备份上传与还原语义。

适用：麒麟 V10 / OceanBase 一体包（`ism_server` 静态链接，不依赖目标机 glibc）。

## 禁止

- **禁止**用本补丁整包覆盖 OceanBase 数据目录、`data/db/ism.db` 或其它业务库文件。
- **禁止**把「备份上传」当成还原。上传只把 SQL 放到备份目录，**不会自动覆盖业务库**。
- 打补丁后若首页「异常点数」从 27 万回到 21 万，说明覆盖了业务数据；用昨天备份在「数据库管理 → 还原」里 **点还原** 恢复。这是预期，不是上传没生效。

## 安装

```bash
# 在补丁解压目录
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-YYYYMMDD
```

脚本只替换：

- `ism_server_user/ism_server`
- `web/dist`（若本包含前端）

**不会**复制 `data/`、OceanBase datadir、`ism.db`。

若本包只有后端二进制（云端内存不足未编前端），请在本机编 `ism-front-end-v2` 的 `dist` 后拷到目标 `web/dist`，并 Ctrl+F5。

## 验收清单（5 条 + 备份语义）

1. **组态保存 / 导航**：新增页面可保存；运行态点击已绑定页面的菜单能跳转，不再整页空白或「运行时异常」。
2. **历史库备份**：历史库配置为 TDengine 时，点「备份 TDengine 历史库」。有本机 `taosdump` 或 docker 容器 `tdengine`（`TD_CONTAINER` 可改名）时可备份；没有工具时界面提示明确原因。备份目录 `data/hisdbbackup/`。
3. **定时入库**：`is_record=1 AND record_type=1` 点位按间隔落档。2208 点 / 600s 一小时约 6 条（允许 ±1）。日志为汇总 `wrote / fromDB / reusedLast / noRealtimeValue`，不再按点打 `no realtime value yet, skip`。无实时值时复用上一档，不因 30s 节流整轮 skip。
4. **脚本**：`native-bitunpack` 与 `anko-onchange` 都能写目标点。纯 native 时每秒 SettleAll，避免源点不变导致输出停在旧值（如 -1）。手动执行 `ExecSysScript` 成功/失败都有 Info/Error 日志；能编译的走 native。
5. **虚拟设备 Excel**：全量导出含「报警触发值(0,1)」（在「告警消除消息」后）。改 Excel 再全量导入后库内值已更新。单设备导入重名走编辑（upsert），不再只 Add 导致 `SNMP_MODEL_EXIST` 静默丢改。
6. **备份上传 ≠ 还原**：上传成功提示「文件已加入还原列表，不会自动覆盖业务库；请在下方点还原。」还原 Tab 顶部有 Alert。

## 双路径说明（问题 3）

- `native-bitunpack`：脚本能编译成纯 `BitGet + SetDeviceData` → 源点变化时拆位；另有 1s 周期 Settle，避免启动后源未再变时输出冻结。
- `anko-onchange`：编译失败 → 解释执行，带 delay 周期跑。

本包不关掉 native（7/24 性能优化仍要），也不恢复 Anko 全量轮询。

## 回滚

停服务后把备份的 `ism_server` / `web/dist` 拷回即可。不要用补丁目录覆盖数据库。

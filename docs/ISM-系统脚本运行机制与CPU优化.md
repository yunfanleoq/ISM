# ISM 系统脚本运行机制与 CPU 根治

## 问题

「按位解析_*」等自动脚本原先按 `Delay`（常为 100ms）用 Anko **整段重解析执行**。本地约 20 条、单条可达百 KB 级脚本时，形成约 200 次/秒解析 + 海量 `SetDeviceData` DB 访问，CPU 飙高。

## 根治架构

1. **原生 BitUnpack**：纯 `BitGet` + `SetDeviceData` 脚本在加载时编译为内存规则，**不再进入 Anko 定时循环**。
2. **变化驱动**：`StoreDeviceRealValue` 检测到源点值变化时执行 BitUnpack，并唤醒依赖该点的通用脚本。
3. **SetDeviceData 快路径**：内存值未变则零 DB；元数据进程内缓存。
4. **通用脚本**：AST 一次解析 + `vm.Run`；依赖变化或 Delay 兜底（下限 1000ms）执行。

## 运行日志

- `native-bitunpack: <脚本名> rules=<N>` — 已原生化，无 Anko 协程
- `anko-onchange: <脚本名> deps=<N> delay=<ms>` — 变化驱动 + 兜底定时

## 运维

- 升级后无需改 UI / 迁表；编辑脚本仍走原 CRUD，会触发 `ScriptCloseChan` 全量重载并重新编译。
- 若某「按位」脚本未原生化，检查是否含分支、循环或其它 API（非纯模式）。
- 冷启动会对全部 BitUnpack 源点做一次全量结算。
- 本地库验证：20 条自动「按位解析」类脚本均可原生化，合计约 1.4 万条规则；启动日志应出现 `script scheduler: native-bitunpack=20 anko-onchange=0`（若仅有这类脚本）。

## 关键代码

- `ism_server_user/task/ISMScript/bitunpack/`
- `ism_server_user/task/ISMScript/ISMScript.go` / `execScript.go`
- `ism_server_user/task/ISMScript/func/deviceData.go`
- `ism_server_user/protocol/common/realDataCache.go`

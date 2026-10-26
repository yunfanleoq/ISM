# 数据包 / 拓扑 输入契约

本技能对「新项目」的输入。航信机房放在 `liu-chang-1A-dev/`，两个 JSON 文件互补：
一个描述**设备清单**（建树用），一个描述**数据模型/寄存器/数据点**（灌点用）。
导入脚本 `scripts/import_hx_project.py` 严格按下面字段读取。

---

## 1. `ism_data_models.json` —— 设备清单 + 模板

`import_hx_project.py` 读 `devices[]` 建设备树与归柜。顶层键：`templates` / `devices` / `sourceNodePoints` / `models`。

### `devices[]`（必需，导入主驱动）

```json
{
  "name": "1A1_U11_S18_1",      // 设备名：编码了层级与归柜（见下「命名规约」）
  "templateType": "A20",        // 模型类型：A20 / A40 / 施耐德UPS（决定 muid 与寄存器组）
  "aiStartAddr": 1000,          // AI 起始地址（写进 extra/采集，非必须用于建树）
  "diStartAddr": 8501,          // DI 起始地址
  "aiName": "AB线电压",          // 仅参考
  "diName": "输入状态2（故障状态）"
}
```

- 设备总数 = `len(devices)`（航信机房 = 76）。导入后 `monitor_list` 中本项目 `type=1` 行数应等于它。
- `templateType` 必须能映射到 `MODELS_DEF` / `TEMPLATE_MAP` 里的模型名（脚本 Step5/Step8）。

### 命名规约（层级靠名字编码，无专用字段）

`1A1_U11_S18_1` 按 `_` 分段：

| 段 | 例 | 用途 | 脚本里 |
|---|---|---|---|
| parts[0] | `1A1` / `1A3` / `UPS...` | 归**机柜** | `get_cab_for_device()`：`'UPS' in name→UPS柜`，`'1A3_'→1A3柜`，`'1A1_'→1A1柜` |
| parts[2] | `S18` / `D14` / `U11` | **设备组**虚拟分组键 | build 脚本 `floor_key = name.split('_')[2]`；前端 `floorKey()` |
| parts[3] | `1` | 同组内序号 | 仅展示 |

> **接新数据包最需要确认的空白点**：如果命名不是这种规律，必须同时改三处取段逻辑：
> `import_hx_project.py:get_cab_for_device()`、`build_ncc_dashboard.py` 的 `name.split('_')[2]`、
> `ISMRunTreeNav.vue:floorKey()`。否则归柜/设备组/page_id 全乱。

---

## 2. `<项目>_complete_project_package.json` —— 数据模型包

`import_hx_project.py` Step7 读它灌 Modbus 数据点。顶层键：
`project` / `deviceModels` / `registerGroups` / `registerPoints` / `monitorTree` / `alarmTriggers` / `displayModel` / `displayLayer` / `statistics`。

| 段 | 用途 | 关键字段 |
|---|---|---|
| `deviceModels[]` | 模型名↔包内 UUID 映射源 | `name`、`uuid` |
| `registerGroups[]` | 寄存器组 | `muid`(指向 deviceModel)、`name`(如 AI数据/DI数据)、`uuid` |
| `registerPoints[]` | 数据点（92 条）| `muid`、`registerGroupUuid`、`name`、`registerAddress`、`type`、`ByteOrder`、`unit`、`record*` |
| `displayModel`/`displayLayer` | （可选）画面骨架，可作为新项目 `MODEL_ID` 来源 | — |
| `monitorTree`/`alarmTriggers`/`statistics` | 参考/校验 | — |

### 关键映射逻辑（脚本如何用包）

包里每次生成的 UUID 是**随机**的，导入脚本**不能**用硬编码旧 UUID 映射，而是：

1. 按**模型名**把包内 `deviceModels[].uuid` → DB 中同名 `devices_model.uuid`（`NAME_MAP` + `pkg2db_muid`）。
2. 寄存器组按 `(模型名, 组名)` → DB 真实 group uuid（`pkg_rg2db`）。
3. 数据点用上面两个映射改写 `muid`/`registerGroupUuid` 后调 `/modbusModelRegisterAdd`。
4. 灌完清掉系统自动生成的 `register0/register1...` 默认命名行（`/modbusModelRegisterDel`）。

> 这就是 `ism-excel-import` 技能坑点 #26 的根治方式：**按名映射，不按固定 UUID**。

---

## 3. 数据模型铁律（来自 ism-excel-import，灌点必守）

- **每种 `devices_model` 必须独占自己的 `register_group`**，严禁跨模型混用同一组（否则 Go 后端字节偏移错乱出天文数字）。
- `register_group.registerCount` = 该模型 `max(offset)+1`（覆盖所有数据点地址）。航信：A20=30、A40=41、UPS=44。
- 数据点 `type` 按 parse_mode 推断：`177→Unsigned short`、`179→Long`、`73→Short`、`71/1→Unsigned short`，默认 `Float`。
- `device_real_data` 必含 NOT NULL：`type=1, device_type=2, oid=<与 uuid 同值>`。
- 导入后跑 `ism-excel-import` 的「导入后必做校验清单」（寄存器组隔离/数据点归属/地址冲突/绑定唯一性/死数据/名称重复）。

---

## 4. 接新项目时给用户的 checklist（哪些必须补充）

- [ ] 确认 dbtype 与连接参数（OceanBase 端口/账号/库名）。
- [ ] 提供上述两个 JSON（或等价的设备清单 + 数据模型包）。
- [ ] **明确设备命名规律**：机柜归属段、设备组分组段分别是第几段？（决定要不要改三处取段逻辑）
- [ ] 提供层级树形状：机房→配电室/楼层→机柜 的中间 zone 名称与顺序（脚本 `add_zone()` 链）。
- [ ] 每种设备模型的寄存器布局（offset/类型/系数/单位），用于 `MODELS_DEF`/`RG_DEFS`/数据点。
- [ ] 新大屏的 `MODEL_ID`（在前端「应用管理」新建空大屏拿 displayUUID，或用包里的 `displayModel`）。

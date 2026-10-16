# 寻智 xunzhi-jf · 菜单入口（契约 v1.0 §8）

> **状态**：演示挂入口 · **非正式验收**  
> **分支**：`feature/xunzhi-jf-env-monitor` · PR [#5](https://github.com/yunfanleoq/ISM/pull/5)  
> **权威**：`docs/机房环境监测-菜单树与路由权限挂点.md` §8 扁平表（已废弃 PM 暂定一级十项树）

## 实现位置

| 层 | 路径 |
|----|------|
| §8 数据源 | `ism-front-end-v2/src/xunzhiJf/menuFlat.js` |
| 路由注册 | `ism-front-end-v2/src/xunzhiJf/registerRoutes.js` → `router.map.js` |
| 登录菜单 | `ism_server_user/conf/MenuConfig.json`（及 Opt/User 菜单） |
| 同步脚本 | `scripts/sync_xunzhi_jf_menu_config.mjs` |

修改 §8 后执行：`node scripts/sync_xunzhi_jf_menu_config.mjs`（**仅 upsert** `XunzhiJfSection`，不替换整份 MenuConfig）

## 演示优先接通

| menu_code | route |
|-----------|-------|
| `jf.monitor.overview` | `/xunzhi/jf/monitor/overview` |
| `jf.viz.plane` | `/xunzhi/jf/viz/plane`（组态 `xunzhi-jf-env-plane`） |
| `jf.alarm.active` | `/xunzhi/jf/alarm/active` |

## 一级分组（侧栏「循智·机房环境」下）

场地监测、告警中心、可视化、资产管理、容量管理、运维工单、能耗管理、安防消防、报表中心、系统管理。

完整 **menu_code → route** 见权威文档 §3 / §8。

## 本机验证（7080）

```bash
./scripts/check_mem_before_compile.sh   # 须 PASS
# 起 ism_server + vue-cli-service serve :7080（见 ism-service-startup 技能）
```

登录 admin → 侧栏 **循智·机房环境** → 展开各分组 → 点击「实时总览 / 环境平面大屏 / 活动告警」可打开。

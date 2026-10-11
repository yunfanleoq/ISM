# 双机软件就绪 · 演练计划（文档 · 非热备产品）

> **不触碰** 麒麟现网、`patches/*kylin`、`releases/` 生产包。本文 **非** 双机热备产品交付，仅为 **软件级备份/恢复口径 + 演练 checklist**。

## 范围

| 包含 | 不包含 |
|------|--------|
| dev/试打样环境演练步骤 | 生产 keepalived/共享存储实装 |
| `data/db/ism.db` 备份与回滚验证 | 麒麟现场切换窗口 |
| 配置 `conf/app.conf` 导出 | 商务「热备验收」签字 |

## 演练计划（建议顺序）

1. **冻结试打样 tip**：记录 `git rev-parse HEAD` 与 `dist/xunzhi-jf/ism-patch-*` 目录名。
2. **备份**：停 `ism_server` → 复制 `data/db/ism.db` → `backups/ism-YYYYMMDD-HHMM.db`；复制 `ism_server_user/conf/app.conf`。
3. **破坏演练**：删改测试项目或导入错误包（仅试打样库）。
4. **恢复**：停服务 → 还原 `ism.db` + `app.conf` → 启服务 → `curl` 登录 + `bash scripts/xunzhi_jf_smoke_env.sh`。
5. **记录**：演练人、耗时、是否 PASS（写入 `docs/ISM-寻智-W1复测证据/` 可选）。

## 与寻智竖切关系

- 动环竖切（W1–W3）、S1–S6 Mock **不依赖** 双机；双机就绪不阻塞 Mock demo。
- wsport 以还原后 `app.conf` 为准（dev 常见 **10225**）。

## 禁止

- 在生产/麒麟目录直接覆盖二进制或数据库。
- 将本文档表述为「已完成双机热备交付」。

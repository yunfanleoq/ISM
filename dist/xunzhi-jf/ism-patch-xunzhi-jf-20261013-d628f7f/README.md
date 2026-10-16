# ism-patch-xunzhi-jf-20261013-d628f7f（试打样 · 非麒麟现网包）

- **Git tip**: `d628f7f`（分支 `feature/xunzhi-jf-env-monitor`，完整 SHA 见 `GIT_TIP.txt`）
- **禁止**：覆盖 `patches/ism-*-kylin`、`releases/` 或现场麒麟生产目录
- **禁止**：使用 `/tmp/keep_xunzhi_be.sh` 等方式保活双口后端（易锁错/写坏 `ism.db`）

## 双口演示（`ISM源码-xunzhi-jf` / 本包文档）

| 服务 | 端口 |
|------|------|
| 前端 | **7081** |
| `ism_server` | **8082** |
| Modbus Mock | **1503**（`XUNZHI_MOCK_PORT`） |

平面闸（须 **同一 git tip 的前后端源码**，本 tar 仅增量脚本+导入包）：
- `noRunTree=1` 真删监视树格；`/#/xunzhi/jf/viz/plane` missPlane 修复
- ProjectUuid 按项目名「寻智机房环境-xunzhi-jf」动态解析（勿写死 uuid）
- 量纲 **DataFormat=ABCD**（导入包内）

## Apply（到已有 ISM 工作区）

```bash
# 1) 工作区先 checkout 到与本包相同的 git tip
git fetch && git checkout feature/xunzhi-jf-env-monitor && git checkout d628f7f

# 2) 复制导入包与脚本（不覆盖 data/db）
bash apply-patch.sh /path/to/ISM-workspace

# 3) 按 docs/ISM-寻智-本机双开端口.md 起 Mock → 8082 后端 → 7081 前端 → 登录 → 重导
python3 scripts/import_xunzhi_jf_project.py
python3 scripts/xunzhi_jf_seed_admin.py --dbtype 1   # 双口 SQLite 登录 1003 时
```

## 导入（已在 ISM 根目录且已 apply）

```bash
python3 scripts/import_xunzhi_jf_project.py
```

## 组态

项目「寻智机房环境-xunzhi-jf」→ 页 `xunzhi-jf-env-plane`；侧栏「监控大屏」或 `/#/xunzhi/jf/viz/plane`

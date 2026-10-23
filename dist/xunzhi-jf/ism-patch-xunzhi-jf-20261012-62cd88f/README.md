# ism-patch-xunzhi-jf-20261012-62cd88f（试打样 · 非麒麟现网包）

- **Git tip**: `62cd88f`（分支 `feature/xunzhi-jf-env-monitor`）
- **禁止**：覆盖 `patches/ism-*-kylin`、`releases/` 或现场麒麟生产目录
- **安装**：解压到任意路径（如 `/opt/ism-patch-xunzhi-jf-dev/`），按 `scripts/xunzhi_jf_smoke_env.sh` 起 Mock+后端+导入

## 内容
- 竖切 ImportProject 包（温湿度+漏水；若含 W2 则含配电/UPS 扩展点表）
- 冒烟脚本与证据采集
- 文档：wsport 以目标机 `conf/app.conf` 为准（dev 常见 10225）

## 导入
```bash
python3 scripts/import_xunzhi_jf_project.py
```

## 组态
项目「寻智机房环境-xunzhi-jf」→「寻智机房环境平面监视-xunzhi-jf」→ 页 `xunzhi-jf-env-plane`

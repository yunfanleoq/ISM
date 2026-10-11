#!/usr/bin/env bash
# 试打样补丁包（独立目录，不触碰 patches/*kylin / releases/）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
SHA7=$(git rev-parse --short=7 HEAD)
# Asia/Shanghai 日期
DATE=$(TZ=Asia/Shanghai date +%Y%m%d)
NAME="ism-patch-xunzhi-jf-${DATE}-${SHA7}"
DEST="dist/xunzhi-jf/${NAME}"
rm -rf "$DEST"
mkdir -p "$DEST/projects-import/xunzhi-jf" "$DEST/scripts" "$DEST/docs"

cp projects-import/xunzhi-jf/xunzhi-jf_complete_project_package.json "$DEST/projects-import/xunzhi-jf/"
cp projects-import/xunzhi-jf/xunzhi-jf_device_points.json "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true
cp projects-import/xunzhi-jf/xunzhi-jf_w2_device_points.json "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true
cp projects-import/xunzhi-jf/xunzhi-jf_w3_device_points.json "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true
cp scripts/xunzhi_jf_*.sh scripts/xunzhi_jf_*.py "$DEST/scripts/" 2>/dev/null || true
chmod +x "$DEST/scripts/"*.sh 2>/dev/null || true
cp scripts/import_xunzhi_jf_project.py "$DEST/scripts/" 2>/dev/null || true
cp docs/ISM-寻智-W1竖切验收.md docs/ISM-寻智-W1业主本机目视自证.md docs/ISM-寻智-S1-envmon演示.md docs/ISM-寻智-S1-图表仪表盘清单草稿.md docs/ISM-寻智-S2-资产U位骨架.md docs/ISM-寻智-S3-NSPC容量骨架.md docs/ISM-寻智-S4-工单骨架.md docs/ISM-寻智-S5-能耗PUE碳骨架.md docs/ISM-寻智-S6-安防消防事件桩.md docs/ISM-寻智-W1复测证据/README.md "$DEST/docs/" 2>/dev/null || true
cp projects-import/xunzhi-jf/append_w2_power_ups.py projects-import/xunzhi-jf/append_w3_ac_h2.py "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true
cp projects-import/xunzhi-jf/append_s1_envmon_demo.py projects-import/xunzhi-jf/append_s1_chart_thicken.py projects-import/xunzhi-jf/append_s1_gauge_plus.py "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true
cp projects-import/xunzhi-jf/xunzhi-jf_dcim_asset_mock.json projects-import/xunzhi-jf/xunzhi-jf_dcim_capacity_mock.json projects-import/xunzhi-jf/xunzhi-jf_dcim_workorder_mock.json projects-import/xunzhi-jf/xunzhi-jf_dcim_energy_mock.json projects-import/xunzhi-jf/xunzhi-jf_security_event_mock.json "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true

cat >"$DEST/README.md" <<EOF
# ${NAME}（试打样 · 非麒麟现网包）

- **Git tip**: \`${SHA7}\`（分支 \`feature/xunzhi-jf-env-monitor\`）
- **禁止**：覆盖 \`patches/ism-*-kylin\`、\`releases/\` 或现场麒麟生产目录
- **安装**：解压到任意路径（如 \`/opt/ism-patch-xunzhi-jf-dev/\`），按 \`scripts/xunzhi_jf_smoke_env.sh\` 起 Mock+后端+导入

## 内容
- 竖切 ImportProject 包（温湿度+漏水；若含 W2 则含配电/UPS 扩展点表）
- 冒烟脚本与证据采集
- 文档：wsport 以目标机 \`conf/app.conf\` 为准（dev 常见 10225）

## 导入
\`\`\`bash
python3 scripts/import_xunzhi_jf_project.py
\`\`\`

## 组态
项目「寻智机房环境-xunzhi-jf」→「寻智机房环境平面监视-xunzhi-jf」→ 页 \`xunzhi-jf-env-plane\`
EOF

cat >"$DEST/apply-patch.sh" <<'EOS'
#!/usr/bin/env bash
# 仅复制脚本与项目包到目标 ISM 源码树旁，不替换现网二进制/数据库
set -euo pipefail
TARGET="${1:-}"
if [[ -z "$TARGET" ]]; then
  echo "用法: bash apply-patch.sh /path/to/ISM-workspace"
  exit 1
fi
PATCH_DIR="$(cd "$(dirname "$0")" && pwd)"
cp -r "$PATCH_DIR/projects-import/xunzhi-jf" "$TARGET/projects-import/"
cp "$PATCH_DIR/scripts/"* "$TARGET/scripts/" 2>/dev/null || true
echo "已复制增量文件到 $TARGET （未动 data/db、未覆盖 kylin 包）"
EOS
chmod +x "$DEST/apply-patch.sh"

(
  cd dist/xunzhi-jf
  tar -czf "${NAME}.tar.gz" "${NAME}"
)
echo "Created $DEST and dist/xunzhi-jf/${NAME}.tar.gz"

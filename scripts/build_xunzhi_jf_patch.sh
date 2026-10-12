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
cp docs/ISM-寻智-演示脚本索引.md docs/ISM-寻智-本机双开端口.md docs/ISM-寻智-W1竖切验收.md docs/ISM-寻智-W1业主本机目视自证.md docs/ISM-寻智-S1-envmon演示.md docs/ISM-寻智-S1-图表仪表盘清单草稿.md docs/ISM-寻智-S2-资产U位骨架.md docs/ISM-寻智-S3-NSPC容量骨架.md docs/ISM-寻智-S4-工单骨架.md docs/ISM-寻智-S5-能耗PUE碳骨架.md docs/ISM-寻智-S6-安防消防事件桩.md docs/ISM-寻智-S6-2.5D导航壳.md docs/ISM-寻智-电池1920抽样.md docs/ISM-寻智-双机软件就绪演练.md docs/ISM-寻智-W1复测证据/README.md "$DEST/docs/" 2>/dev/null || true
mkdir -p "$DEST/docs/ISM-寻智-W1复测证据"
cp docs/ISM-寻智-W1复测证据/机房环境监测-复测-aa65d820.md docs/ISM-寻智-W1复测证据/机房环境监测-复测-大屏.md "$DEST/docs/ISM-寻智-W1复测证据/" 2>/dev/null || true
git rev-parse HEAD >"$DEST/GIT_TIP.txt"
cp projects-import/xunzhi-jf/append_w2_power_ups.py projects-import/xunzhi-jf/append_w3_ac_h2.py projects-import/xunzhi-jf/append_w4_battery_sample.py projects-import/xunzhi-jf/append_s6_25d_nav_shell.py projects-import/xunzhi-jf/append_s6_chart_catalog_push.py "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true
cp projects-import/xunzhi-jf/append_s1_envmon_demo.py projects-import/xunzhi-jf/append_s1_chart_thicken.py projects-import/xunzhi-jf/append_s1_gauge_plus.py "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true
cp projects-import/xunzhi-jf/xunzhi-jf_dcim_asset_mock.json projects-import/xunzhi-jf/xunzhi-jf_dcim_capacity_mock.json projects-import/xunzhi-jf/xunzhi-jf_dcim_workorder_mock.json projects-import/xunzhi-jf/xunzhi-jf_dcim_energy_mock.json projects-import/xunzhi-jf/xunzhi-jf_security_event_mock.json projects-import/xunzhi-jf/xunzhi-jf_25d_shell_page.json projects-import/xunzhi-jf/xunzhi-jf_battery_device_points.json "$DEST/projects-import/xunzhi-jf/" 2>/dev/null || true

cat >"$DEST/README.md" <<EOF
# ${NAME}（试打样 · 非麒麟现网包）

- **Git tip**: \`${SHA7}\`（分支 \`feature/xunzhi-jf-env-monitor\`，完整 SHA 见 \`GIT_TIP.txt\`）
- **禁止**：覆盖 \`patches/ism-*-kylin\`、\`releases/\` 或现场麒麟生产目录
- **禁止**：使用 \`/tmp/keep_xunzhi_be.sh\` 等方式保活双口后端（易锁错/写坏 \`ism.db\`）

## 双口演示（\`ISM源码-xunzhi-jf\` / 本包文档）

| 服务 | 端口 |
|------|------|
| 前端 | **7081** |
| \`ism_server\` | **8082** |
| Modbus Mock | **1503**（\`XUNZHI_MOCK_PORT\`） |

平面闸（须 **同一 git tip 的前后端源码**，本 tar 仅增量脚本+导入包）：
- \`noRunTree=1\` 真删监视树格；\`/#/xunzhi/jf/viz/plane\` missPlane 修复
- ProjectUuid 按项目名「寻智机房环境-xunzhi-jf」动态解析（勿写死 uuid）
- 量纲 **DataFormat=ABCD**（导入包内）

## Apply（到已有 ISM 工作区）

\`\`\`bash
# 1) 工作区先 checkout 到与本包相同的 git tip
git fetch && git checkout feature/xunzhi-jf-env-monitor && git checkout ${SHA7}

# 2) 复制导入包与脚本（不覆盖 data/db）
bash apply-patch.sh /path/to/ISM-workspace

# 3) 按 docs/ISM-寻智-本机双开端口.md 起 Mock → 8082 后端 → 7081 前端 → 登录 → 重导
python3 scripts/import_xunzhi_jf_project.py
python3 scripts/xunzhi_jf_seed_admin.py --dbtype 1   # 双口 SQLite 登录 1003 时
\`\`\`

## 导入（已在 ISM 根目录且已 apply）

\`\`\`bash
python3 scripts/import_xunzhi_jf_project.py
\`\`\`

## 组态

项目「寻智机房环境-xunzhi-jf」→ 页 \`xunzhi-jf-env-plane\`；侧栏「监控大屏」或 \`/#/xunzhi/jf/viz/plane\`
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

# ISM 现场诊断脚本（FinalShell 上传用）

将本目录下 3 个脚本上传到测试机对应路径：

```
/opt/ISM/ism-release-oceanbase-20260708/scripts/
├── run_full_field_check.sh    ← 推荐，一键全流程
├── collect_diagnose_log.sh    ← 仅采集诊断
└── check_login_and_user.sh    ← 仅查登录/user表
```

## FinalShell 上传步骤

1. 左侧文件树进入 `/opt/ISM/ism-release-oceanbase-20260708/scripts/`
2. 右键 → **上传** → 选择本机这 3 个 `.sh` 文件
3. 终端执行赋权：

```bash
cd /opt/ISM/ism-release-oceanbase-20260708
chmod +x scripts/run_full_field_check.sh scripts/collect_diagnose_log.sh scripts/check_login_and_user.sh
```

## 推荐：一键排查（重启 + 等待 + 全量日志）

```bash
cd /opt/ISM/ism-release-oceanbase-20260708
bash scripts/run_full_field_check.sh
```

约 2~3 分钟完成后，下载：

```
/opt/ISM/ism-release-oceanbase-20260708/logs/ism_field_check_*.tar.gz
```

## 不重启，只采集当前状态

```bash
bash scripts/run_full_field_check.sh --no-start
```

## 单独脚本

| 脚本 | 用途 |
|------|------|
| `collect_diagnose_log.sh` | 环境 + 端口 + API + 数据库 |
| `collect_diagnose_log.sh --start --wait 120` | 先启动再采集 |
| `check_login_and_user.sh` | 登录 1003 专项（user 表 + curl login） |

## 注意

- **只执行命令，不要把脚本输出文字粘贴回终端**
- 日志在 `logs/` 目录，用 FinalShell 右键下载发回分析

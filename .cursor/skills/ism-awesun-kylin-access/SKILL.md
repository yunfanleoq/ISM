---
name: ism-awesun-kylin-access
description: >
  Mac 经向日葵 AweSun MCP 访问航信麒麟 V10 现场（跳板 Windows → 内网麒麟）。
  含 forward 端口转发、cmd2 不可靠说明、cpolar 固定 SSH/7090 隧道配置。
  触发词：向日葵 MCP、forward、portforward、麒麟 SSH、192.168.110、cpolar、
  跳板 DESKTOP-06IOM7I、Mac 连麒麟、航信现场部署。
disable-model-invocation: false
---

# ISM · Mac 经向日葵访问麒麟 V10

主文档：`docs/ISM-麒麟V10部署方案.md` **§7**。帆软同类流程参考：`~/.cursor/skills/fanruan-report-jump-deploy/SKILL.md`。

---

## 固定事实（航信电力现场，以 ls 为准）

| 项 | 典型值 |
|----|--------|
| 跳板 Windows | `DESKTOP-06IOM7I`，`remote_id=1749886310`，识别码 `1334266437` |
| 麒麟 SSH | `root@192.168.110.56:22`（内网 IP 以 Xshell 为准） |
| OB 一体包根目录 | `/opt/ISM/ism-20261014` |
| 前端 / 后端端口 | **7090** / **8091**（勿与客户 7080/8081、柴发 8082 混） |
| Mac 直连麒麟内网 IP | **不通**（须 forward 或 cpolar） |
| MCP 配置 | `~/.cursor/mcp.json` → `AWESUN_API_URL=http://127.0.0.1:8908` + `AWESUN_API_TOKEN`（Mac 向日葵 MCP 页复制） |

**禁止**：覆盖现场 `conf/`、业务库、客户 `ism.cpolar.cn` 隧道指向错误端口。

---

## 路径 A：向日葵 MCP forward（每次会话）

### 闭环

```
- [ ] 1. 跳板 Windows 在线，Mac 向日葵 + MCP 已启用
- [ ] 2. control_connect type=forward remote_id=1749886310
- [ ] 3. control_portforward target_addresses（见下）
- [ ] 4. Mac 上用返回的 127.0.0.1:端口 → ssh / curl
- [ ] 5. 可选：ssh-copy-id 到 root@127.0.0.1:<映射22端口>，便于后续 BatchMode
```

### portforward 示例（麒麟 IP 以现场为准）

```json
[
  "192.168.110.56:22",
  "192.168.110.56:7090",
  "192.168.110.56:8091",
  "192.168.110.56:9200"
]
```

返回的 `local_address` **每次 forward 会变**，禁止写死旧端口。

### Mac 验证

```bash
curl -sI http://127.0.0.1:<7090映射>/hx-dcim/index.html | head -3
curl -s -X POST http://127.0.0.1:<8091映射>/login \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
```

### cmd2 / control_command

**勿依赖** MCP `control_command` 代发 Windows CMD（Mac 客户端常见 **API 超时**）。替代：

- forward 后在 Mac `ssh` / `scp`；
- 或桌面 MCP 在 Xshell 粘贴命令；
- 或文件 MCP 上传到跳板 `Downloads` 再 scp 到麒麟。

---

## 路径 B：麒麟 cpolar（Mac 不依赖 Windows）

适用：长期调试、Agent 在 Mac 上 `ssh`/`scp` 直连 `*.cpolar.cn`。

详细步骤见本目录 **`references/cpolar-kylin-field.md`**；通用安装见 `docs/ISM-测试环境部署指引.md` **§11**。

### Agent 能否代配 cpolar？

| 步骤 | Agent |
|------|--------|
| 在麒麟 **检查** cpolar 是否已装、`9200` 是否监听 | ✅ 经 forward + 用户 Xshell 执行只读命令 |
| 写入 **authtoken**、预留 **固定 TCP/子域名** | ❌ 须登录 <https://www.cpolar.com>（你的账号） |
| 编辑 `cpolar.yml` 并 `systemctl restart cpolar` | ✅ 用户 SSH 到麒麟后粘贴脚本，或配好密钥后 Agent ssh |

**铁律**：authtoken、预留地址勿进 Git、勿在对话里明文贴 root 密码。

### 现场一体包推荐隧道（与 7090/8091 对齐）

- **TCP**：本地 `22` → 后台「保留 TCP」地址（例 `*.tcp.cpolar.cn:端口`）
- **HTTP**：本地 **`7090`** → 预留二级子域名（**新名字**，勿占用客户 `ism` / 测试机 `largescreen`）

前端 dev 代理已含 `/api` 时，外网通常 **只暴露 7090** 即可；若要直连 API 可再加 `8091` 的 http 隧道（注意安全）。

### Mac 经 forward 打开 cpolar Web UI

forward 映射 `192.168.110.56:9200` 后，Mac 浏览器访问 **`http://127.0.0.1:<映射端口>/`**，用 cpolar 账号登录，可在 UI 里增隧道（改完仍建议同步写入 `cpolar.yml` 并 `systemctl restart cpolar`）。

---

## 路径 C：跳板 Windows → scp（无 forward 时）

1. `control_connect` `type=file` 上传到 `C:\Users\xunan\Downloads\`
2. 远程 CMD / Xshell：`scp` 到 `root@192.168.110.56:/tmp/`
3. 麒麟执行 `artifacts/deploy_kylin_hx_dcim.sh` 或文档 §6

---

## Mac 打 tar 给麒麟

```bash
COPYFILE_DISABLE=1 tar -czf artifacts/hangxin-YYYYMMDD-主题.tar.gz -C /tmp hangxin-YYYYMMDD-主题
```

Linux 上 `LIBARCHIVE.xattr.com.apple.provenance` 警告可忽略。

---

## 故障排查

| 现象 | 处理 |
|------|------|
| forward 连不上 | 跳板是否在线；MCP Token 是否最新；Reload MCP |
| 映射端口 curl 失败 | IP/端口是否仍为 110.56 与 7090/8091；ISM 是否 `start-all.sh` |
| cpolar 外网 SSH 失败 | `remote_addr` 与后台 TCP 一致；`systemctl status cpolar` |
| HTTP 隧道 404 | 隧道 `addr` 是否为 **7090**；`hx-dcim/index.html` 是否 200 |

---

## Performance Notes

- forward 会话占用向日葵连接，用完可 `control_disconnect`。
- 48MB+ 包优先 **file 传 Windows** 再 scp 麒麟，或 forward 后 Mac **scp 直连映射 22**。

## Troubleshooting

见上表；OceanBase/TD 容器勿随意 `docker stop`，见 `ism-field-lessons`。

# 麒麟 V10 现场 cpolar 配置（航信 OB 一体包 · 7090/8091）

> 机器示例：`192.168.110.56`，`/opt/ISM/ism-20261014`。  
> **勿**把隧道指到客户生产 `7080/8081` 或占用其 `ism.cpolar.cn`。

---

## 1. 先检查（在麒麟 SSH 里执行）

```bash
which cpolar && cpolar version
systemctl status cpolar --no-pager
ss -lntp | grep -E ':9200|:7090|:8091|:22'
curl -s -o /dev/null -w "9200=%{http_code}\n" http://127.0.0.1:9200/
grep -E '^tunnels:|addr:|remote_addr:|subdomain:' /usr/local/etc/cpolar/cpolar.yml 2>/dev/null | head -30
# authtoken 行勿外传、勿贴到聊天
```

已安装且 `active` → 跳到 **§3** 增改隧道；否则 **§2**。

---

## 2. 安装与 authtoken（仅首次）

```bash
uname -m   # 须 x86_64
curl -L https://www.cpolar.com/static/downloads/install-release-cpolar.sh | sudo bash
```

1. 浏览器打开 <https://www.cpolar.com> 注册/登录。  
2. 仪表盘 → **验证** → 复制 **authtoken**。  
3. 麒麟上：

```bash
sudo cpolar authtoken YOUR_AUTHTOKEN
```

---

## 3. 在 cpolar 官网预留地址（固定隧道需套餐）

登录 <https://dashboard.cpolar.com>：

| 用途 | 后台操作 | 记下 |
|------|----------|------|
| SSH | **预留 → 保留 TCP** | 例 `12.tcp.cpolar.cn:12345` |
| ISM Web | **预留 → 保留二级子域名** | 例 `hangxin-kylin` → `https://hangxin-kylin.cpolar.cn` |

名称 **不要** 与客户 `ism`、测试机 `largescreen` 冲突。

---

## 4. 编辑 cpolar.yml

```bash
sudo cp /usr/local/etc/cpolar/cpolar.yml /usr/local/etc/cpolar/cpolar.yml.bak.$(date +%F_%H%M)
sudo vi /usr/local/etc/cpolar/cpolar.yml
```

在已有 `authtoken` 与 `tunnels:` 下 **追加或修改**（地址换成你在后台预留的）：

```yaml
tunnels:
  kylin-ssh:
    addr: 22
    proto: tcp
    region: cn_vip
    remote_addr: 12.tcp.cpolar.cn:12345

  hangxin-web-7090:
    addr: 7090
    proto: http
    region: cn_vip
    subdomain: hangxin-kylin
```

说明：

- `addr: 7090` 对应当前一体包 **Python 静态前端**（含 `/hx-dcim`、`/api` 代理到 8091 时走同源）。
- 若暂无固定子域名，删掉 `subdomain` 行，启动后到 dashboard **状态页** 看临时 HTTPS 地址（免费版会变）。

```bash
sudo systemctl enable cpolar
sudo systemctl restart cpolar
sudo systemctl status cpolar --no-pager
journalctl -u cpolar -n 30 --no-pager
```

---

## 5. Web UI 方式（可选）

Mac 已建向日葵 forward 映射 `9200` 时，浏览器打开 `http://127.0.0.1:<本地映射端口>/`，用 cpolar 账号登录，图形界面添加 **TCP 22** 与 **HTTP 7090** 隧道。  
**仍建议** 把最终配置同步进 `cpolar.yml` 并 restart，便于备份与复现。

---

## 6. Mac 外网验证（不经过 Windows）

```bash
# SSH（端口、域名以 dashboard 为准）
ssh -p 12345 root@12.tcp.cpolar.cn

# Web
curl -sI https://hangxin-kylin.cpolar.cn/hx-dcim/index.html | head -5
curl -s -X POST https://hangxin-kylin.cpolar.cn/api/login \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
```

若 HTTPS 经 cpolar 反代 `/api` 异常，可仅暴露 SSH，Mac 上 `ssh` 进去后 `curl http://127.0.0.1:7090/...` 验收。

---

## 7. 安全建议

- root 密码仅现场维护；建议 **ssh-copy-id** 后改 `PasswordAuthentication no`（可选，需团队同意）。
- 勿在公网隧道暴露未鉴权管理口；7090 仍走 ISM 登录。
- authtoken 泄露时在官网重置并 `cpolar authtoken` 更新。

---

## 8. Agent 后续自动化（可选）

用户在 Mac 执行一次（forward 映射 22 后，输入 root 密码）：

```bash
ssh-copy-id -p <本地22映射> root@127.0.0.1
```

之后 Agent 可用 `ssh -p <映射> root@127.0.0.1 '...'` 维护 cpolar/ISM（映射端口随 forward 变）。

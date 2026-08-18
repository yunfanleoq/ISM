---
name: ISM Docker日志与TDengine磁盘满
overview: 麒麟现场 /var 被 Docker 容器日志撑满导致部署失败；清日志恢复后需补 Docker 日志轮转与 ISM 部署脚本加固。
todos:
  - id: daemon-json
    content: 现场 /etc/docker/daemon.json 增加 log-opts（max-size + max-file）
    status: completed
  - id: install-script
    content: install_docker_kylin_sp3.sh 安装时自动写入 daemon.json 日志策略
    status: completed
  - id: start-all-log-opt
    content: start-all.sh 的 docker run 增加 --log-opt max-size/max-file
    status: completed
  - id: compose-logging
    content: docker-compose.tdengine.yml / oceanbase 增加 logging 段
    status: completed
  - id: tdengine-storage
    content: 排查 TDengine 历史库存储报错根因，避免日志风暴复发
    status: pending
---

# ISM Docker 日志与 TDengine 磁盘满

详见 `docs/ISM-Docker日志与TDengine磁盘满.md`。

## 结论

- 非缺组件，是 **/var 磁盘满**（Docker json 日志 + 历史库报错日志风暴）
- 部署包 `ism-release-oceanbase-20260817-0001-c851` 安装 Docker 时 **未配置日志轮转**
- 现场已恢复；需在源码/下一版发布包中固化日志限制

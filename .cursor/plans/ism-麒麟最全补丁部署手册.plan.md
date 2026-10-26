# ISM 麒麟 V10 · 最全补丁包部署手册

> 当前最终包：`releases/ism-patch-kylin-ultimate-20260709-1822-f25f.zip`  
> 构建脚本：`scripts/build_kylin_ultimate_patch.sh`

## 一键应用

```bash
unzip ism-patch-kylin-ultimate-20260709-1822-f25f.zip
cd ism-patch-kylin-ultimate-20260709-1822-f25f
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
```

## 本包已含

- MQTT Broker 缺配置 panic 修复 + 默认关闭内置 Broker
- 写盘收敛 / 强制分页 / 大屏分批
- OceanBase + TDengine 离线镜像
- 完整 conf（含 mqtt_broken_config.json、videoConfig.json）
- 部署手册与诊断脚本

登录：admin / 123456  
端口：前端 7090 / 后端 8091 / OB 2881 / TD 6041

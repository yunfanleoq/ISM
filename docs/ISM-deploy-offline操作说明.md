# deploy-offline.sh 调用方式

```bash
# 解压整包后进入根目录（含 deploy-offline.sh）
cd /path/to/ism-release-oceanbase-*
sudo bash deploy-offline.sh
```

成功后访问：`http://<本机IP>:7090/#/login`，账号 `admin` / `123456`。

说明：`deploy-min.sh` 只是薄封装，最终同样调用 `deploy-offline.sh`。

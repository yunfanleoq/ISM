# 柴发监控项目导入包（后沙峪改造-柴发部分）

从 Sqlite3 备份提取的 **独立 ISM 项目导入文件**。

导入后会在客户系统中 **新建** 项目「后沙峪改造-柴发部分」，与现有项目相互独立。

## 统计

- 设备模型: 16
- 寄存器组: 176
- 数据点: 12959
- 设备: 16（柴发楼1~4层各端口）
- 组态大屏: 无（原备份无）

## 导入（推荐）

```bash
cd chaifa-project-import
python3 import_chaifa_project.py \
  --base-url http://127.0.0.1:8081 \
  --user admin \
  --password 123456
```

成功后项目列表同时出现原项目与「后沙峪改造-柴发部分」。

## 手动 curl

```bash
TOKEN=$(curl -s -X POST http://127.0.0.1:8081/login \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

curl -sS -X POST http://127.0.0.1:8081/ImportProject \
  -H "Authorization: $TOKEN" \
  -H 'Content-Type: application/json' \
  --data-binary @houshayu-chaifa-ISM-project-package.json \
  --max-time 1800
```

期望 `code: 0`。数据点较多，导入可能需数分钟。

## 注意

1. 仅管理员可导入；Authorization 不要加 Bearer。
2. 新建项目，不覆盖已有项目。
3. 设备 IP/端口来自原现场；客户网段不同需导入后修改设备扩展参数。
4. 本包不含组态大屏。

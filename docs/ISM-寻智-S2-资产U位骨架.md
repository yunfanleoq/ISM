# S2 骨架 · dcim-asset（契约 v2 · 不阻塞 W1/W2）

> U 位全生命周期属同窗全量规划；本文件仅 **数据模型/接口占位**，无业务闭环实现。

## 实体（规划）

```
Rack(rack_id) 1─* UPosition(u_index) 1─0..1 Asset(asset_id)
```

| 字段 | 说明 |
|------|------|
| rack_id | 关联 MonitorList 机房/机柜节点 |
| u_index | 1..42 |
| asset_id | 可选，上架后绑定 |
| status | empty / occupied / reserved |

## 接口占位（未实现 HTTP）

- `GET /api/xunzhi-jf/asset/rack/{id}/u-positions` — 列表
- `POST /api/xunzhi-jf/asset/u-position/bind` — 上架绑定

实现路径：复用 ISM `StaticData` 或后续独立表；**禁止** SpringBoot 平行系统。

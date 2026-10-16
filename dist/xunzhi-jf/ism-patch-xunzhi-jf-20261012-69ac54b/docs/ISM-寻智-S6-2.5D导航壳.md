# S6 · 2.5D 导航壳（占位 · 非产品级引擎）

> **无现场图 / 无 2.5D 引擎** 时的降级交付：**平面组态** `xunzhi-jf-env-plane` 仍为监视兜底；本壳仅提供 **导航入口 + 占位页 JSON 草稿**。

## 交付物

| 路径 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/append_s6_25d_nav_shell.py` | 平面页右上角「2.5D 导航」链到占位 `pageUUID` |
| `projects-import/xunzhi-jf/xunzhi-jf_25d_shell_page.json` | 占位页元数据（`xunzhi-jf-25d-shell`） |
| 导入包 `_meta.s6_25d_page_id` | `c3d4e5f6a7b84990abcdef25d000001` |

## 降级说明

1. **ImportProject** 当前只写入 **单页** `displayLayer`；占位页在 `draftPages` / 独立 JSON，需组态编辑器 **手工建页** 或后续多页导入扩展后链通。
2. 点击导航若目标页不存在，回 **平面监视** 操作。
3. **禁止** 宣称已交付产品级 2.5D 渲染引擎或现场实景模型。

## 管线

```bash
python3 projects-import/xunzhi-jf/append_s6_25d_nav_shell.py
```

复导入项目后可见平面页导航文案。

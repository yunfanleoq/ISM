---
name: RC08问题交付核对
overview: 从 docx(3) 问题清单对照源码与补丁包，完成收口并给出可发客户结论（有条件可发）。
todos:
  - id: extract-matrix
    content: 写出 docx(3) 八项对照矩阵文档（问题/修复记录/源码/交付包/判定）
    status: completed
  - id: fix-bacnet-route
    content: 补齐 bacnetModelList 路由并对齐 IEC61850 列表接口行为
    status: completed
  - id: fix-name-split-sticker
    content: 闭环最后_拆分设备/点位 + 去掉在线设备 KPI 贴图
    status: completed
  - id: bump-version-rebuild
    content: VERSION 对齐 RC08bate 并产出含 RealDataChanel 根治的补丁包
    status: completed
  - id: smoke-then-ship
    content: 按最小冒烟清单验收；结论为有条件可发客户
    status: completed
isProject: false
---

# RC08bate 问题清单 vs 修复记录 — 交付核对计划

## 交付产物（2026-07-22）

- Zip：`releases/ism-patch-rc08bate-docx3-20260722-2232-9394.zip`
- 目录：`releases/ism-patch-rc08bate-docx3-20260722-2232-9394/`
- 对照文档：`docs/ISM-RC08bate-docx3-交付核对.md`

## 结论

**有条件可发客户**：补丁含 A1–A3 / B4–B8 收口；API 与包内静态冒烟通过；现场需补 10 分钟通道观察与 UI 人眼确认。

详见对照文档冒烟表与客户注意事项。

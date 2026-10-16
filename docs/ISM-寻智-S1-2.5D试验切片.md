# S1 · 2.5D 试验切片（Three.js · 非正式验收）

> **范围**：寻智机房 Vue2 大屏壳内的可开关 **2.5D 观感演示**，不替换 `#/xunzhi/jf/viz/plane` 平面组态，不接真实 BIM/CAD。

## 如何打开

1. 登录 ISM，进入项目 **寻智机房环境-xunzhi-jf**（若直链无效，请先走一次侧栏首页以对齐 `ProjectUuid`）。
2. 侧栏：**可视化** → **2.5D 导航**。
3. 浏览器地址：`/#/xunzhi/jf/viz/nav25d`。

## 页面说明

| 控件 | 作用 |
|------|------|
| **2.5D 试验** | 默认：Three.js 斜俯视机柜阵列，拖拽旋转、滚轮缩放 |
| **平面组态兜底** | 切换为原组态占位页 iframe（`xunzhi-jf-25d-shell`） |
| **机房/楼层** | A/B/C 三区假数据切换（仅改阵列布局与告警示意块） |
| **打开环境平面大屏** | 跳转正式监视 `#/xunzhi/jf/viz/plane` |

页顶横幅标明 **试验切片 · 非正式验收**。

## 技术要点

- 依赖：`ism-front-end-v2` 已声明 npm 包 `three`（与 vue-cli/webpack 一致，走 `import * as THREE from 'three'`）。
- 组件：`ism-front-end-v2/src/xunzhiJf/components/XunzhiJfViz25d.vue`
- 路由页：`ism-front-end-v2/src/xunzhiJf/pages/VizNav25dPage.vue`（`jf.viz.nav25d`）

## 限制（勿对外宣称验收）

- 机柜为 **程序化方盒**，无 GLB、无现场布置图。
- 告警为 **假数据色块**，未绑南向点表。
- 销毁路由/切换回平面时会 **dispose** renderer、几何与材质，避免 WebGL 泄漏。

## 相关文档

- 平面监视：`docs/ISM-寻智-S1-envmon演示.md`
- 原 S6 占位壳说明：`docs/ISM-寻智-S6-2.5D导航壳.md`（组态 JSON 草稿；本切片为 Vue 内 Three 增量）

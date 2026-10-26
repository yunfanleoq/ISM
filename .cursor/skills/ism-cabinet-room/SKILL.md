---
name: ism-cabinet-room
description: >
  做机房、机柜、动环画面，让机柜看起来像真的。柜体按外框、立柱、前门、门缝、通风、状态灯拆开；组态层用设备图、U 位和测点弹窗，不要用纯色方块凑。
  触发词：做机柜、机柜、机房、动环、像真的、柜体、U位、U 位、前门、设备贴图、机房透视、列头柜、冷通道、热通道。
  用户说做机柜、机房、动环、看起来像真的、柜门、设备贴图、U 位时，先读本技能再改画面。
disable-model-invocation: false
---

# 机房机柜：以后照做

大屏钻探、备份、`page_id` 在 `ism-scada-dashboard`。本页只写柜体怎么拼才像真的。航信和寻智的端口、项目名不要顺手改，那条在 `ism-field-lessons`。

先选定一种画面，不要把两套零件塞进同一个格子：

| 要的效果 | 照这个做 |
|---|---|
| 一整间机房，能旋转、缩放，柜子立在地上 | 寻智 `XunzhiJfViz25d.vue` 的 `addCabinet` |
| 组态平面上的一只柜，要绑测点、悬停弹窗 | `view-svg-dc-rack`，设备正面用 `ism-view-png-image` |
| 一个真设备里拆出「A列头 / B列尾」 | `virtualCabinet.js`，不往 `monitor_list` 插假设备 |

寻智那份在 `/Users/yunfanleo/cursorProjects/ISM源码-xunzhi-jf/ism-front-end-v2/src/xunzhiJf/components/XunzhiJfViz25d.vue`。航信仓库里没有这个文件。不要把寻智进程或库拷进航信。

## 柜体怎么拆

真柜大约宽 0.6、高 2、深 1。代码里就是这组常量：`CAB_W = 0.6`，`CAB_H = 2`，`CAB_D = 1`。列与列 `GAP_COL = 0.32`。面对面的冷通道 `GAP_ROW_COLD = 1.05`，背靠背的热通道 `GAP_ROW_HOT = 0.22`。奇偶行 `faceSign` 相反，门朝冷通道。

`addCabinet` 按这个顺序加零件。少一件，远看就又变成一个盒子：

1. **外框** `shell`：内宽、内深的盒子，高是 `CAB_H * 0.98`。
2. **立柱** `posts`：四角各一根，边长 `postT = 0.045`，高 `CAB_H`。
3. **前门** `door`：厚 `doorTh = 0.035`，宽 `innerW * 0.9`，高 `CAB_H * 0.92`，贴在朝向那一面的外侧。
4. **门缝** `seam`：门上两根竖缝，宽 `0.018`。
5. **通风** `vent`：门内上方一条横槽。
6. **状态灯** `led`：正常用绿，告警用红。告警再加一盏 `lamp`。
7. 房间：`addAisleFloors` 铺冷/热通道色带，`addCeilingStrips` 铺顶灯带，`addRoomPillars` 放四角柱。

这套没有 U 位网格，没有标牌网格，也没有设备照片。材料是 `MeshStandardMaterial` 的颜色，加上 `metalness` / `roughness`。像真靠的是比例、分件和透视。

航信 `IndustrialObjects.js` 的 `case 'serverrack'` 是反例：一个 `0.4 × 1 × 0.5` 的灰盒子，里面五个 `0.35 × 0.12` 的黑块，再加绿点。又矮又扁，没有门、没有立柱、没有通道。`SceneTemplates.js` 的数据中心模板摆的就是这种 `serverrack`。新机房不要再生成这个 case。

## 透视和缩放

机房相机用 `XunzhiJfViz25d` 里这组限制，否则不是钻进地板，就是柜子被压扁：

- `PerspectiveCamera` 视场 `48`，位置 `(5.5, 5.2, 6.5)`，看向 `(0, 0.8, 0)`。
- `OrbitControls`：`maxPolarAngle = Math.PI * 0.48`，距离 `2.5` 到 `22`。
- 雾 `Fog`，近 `22`、远 `50`。
- 环境光、主光、补光都要。只留环境光会发灰、发平。

组态图片的翻转在 `ViewPngImage.vue` 和 `ViewSvgImage.vue`，不是机房相机：

- `detail.style.transform === -1099` → `rotateY(180deg)`，左右翻。
- `-1098` → `rotateX(180deg)`，上下翻。
- 外层 CSS `perspective: 1000px` 只配合这个翻转，不会把平面图变成机房。
- `animate.selected` 含 `Zoom` 时，SVG `animateTransform` 在 `0.9;1;0.9` 之间缩放。

`view-svg-dc-rack` 的 SVG 是 `preserveAspectRatio: 'xMidYMid meet'`，格子被拉高时线框仍保持比例。`ism-view-png-image` 的图是 `preserveAspectRatio="none meet"`，会按格子宽高拉变形。设备图的 `position.w` / `position.h` 按原图比例设，不要拉成正方形。

## 组态平面上的一只柜

| 零件 | 组件名 | 代码 |
|---|---|---|
| 外框和层线 | `view-svg-dc-rack` | `dcEnvSymbols.js` 里 `id: 'dc-rack'`。`dcenv/19-dc-rack.vue` 只有一行：`createPowerSymbol('dc-rack')` |
| 设备图 | `ism-view-png-image` | `ISMComponents/Images/ViewPngImage.vue` |
| 标牌 | `view-svg-text` | 文字。居中规则在 `ism-scada-dashboard`，这里不重复 |
| 测点弹窗 | `BoundPointsPopup` | mixin 是 `mixins/imageBoundPointsMixin.js` |

`dc-rack` 的 viewBox 是 `0 0 80 96`，默认宽 `64`、高 `96`。几何是外框矩形、左右两根立柱、六条横线、三个点和一颗灯。这是线框，不是照片。柜门里的设备用 `ism-view-png-image` 叠上去，不要再画纯色矩形当服务器。

图片地址：

- 字段是 `detail.style.imageURL`，或 `detail.style.diy` 里 `key === 'imageURL'`。
- 解析函数是 `pickDisplayImageUrl` / `resolveDisplayAssetUrl`（`utils/displayAssetUrl.js`）。
- 系统图放在 `ism-front-end-v2/public/static/ISM/systemImage/`。`ISMBase.vue` 里电气图写成 `/static/ISM/systemImage/electric/` 加文件名。
- 自定义图走 `/static/customPel/`、`/static/upload/`、`/static/DiyUpload/`。这三条在运行时前面加 `/api`。
- `public/static/MapStorage/BitMap/末端设备/Rack.png` 是末端设备图，不是服务器机柜。不要拿来当柜体贴图。

U 位要在平面上逐层可点时：按 42U 把柜门内高分成 42 份，每台设备占若干份，每份一个 `ism-view-png-image`。寻智演示页 `DemoShowPage.vue` 的 `uSlots` 是从 42 数到 1，样式类是 `.cabinet` / `.u-slot`。那是 HTML 列表，不是组态 cell，不要并进 `view-svg-dc-rack`。

标牌用 `view-svg-text`。不要新造 nameplate 组件。

## 格子必填字段

生成 cell 用 `utils/ismSceneNormalizer.js` 的 `createISMNode`，它会补上下面这些默认值。手写 JSON 时自己带上。缺了会空白，或在 `includes` 上崩，`$el` 变成注释节点。

- `detail.animate` 整个对象都要在。`ViewPngImage.vue` 的 `initComponents` 直接读 `option.animate.isExpression`，没有 `animate` 会抛错。
- `detail.animate.selected` 必须是数组 `[]`。旧组件把它赋给 `animateType` 再 `includes(...)`。
- `detail.animate.animateElement` 必须是数组 `[]`。`ViewPngImage.vue` 读 `option.animate.animateElement.length`，没有这个数组会抛错。`createPowerSymbol` 用了 `|| []`，PNG 图没有这条保护。
- `detail.style.visible = 1`
- `detail.style.diy` 必须是数组。`initComponents` 会遍历它。
- `detail.style.position` 为 `{x, y, w, h}`
- 标牌再加 `detail.style.text`、`foreColor`、`fontSize`
- `display_model_layer.components` 是 `{"cells":[...]}`，不能是裸数组

动环符号的动画模板在 `createPowerSymbol.js` 的 `animateConfig()`：`selected: []`，`animateElement` 里有 blink、millcolorGrad、animateSpin。

## 测点绑定和弹窗

图片、矢量图、文字的悬停弹窗走 `imageBoundPointsMixin.js`，不是大屏的 `DeviceHoverTooltip`。

- `detail.dataBind` 是数组。项里有 `dataID` 或 `dataName` 才弹（`canShowBoundPoints`）。
- `editMode` 或 `IsToolBox` 时不弹。
- 移入 `show`，移动 `move`，移出 `hide(false)`，点击 `togglePin`。
- `ViewPngImage`、`ViewSvgImage`、`ViewSvgText` 里已经有 `<bound-points-popup ref="boundPointsPopup">`。为了弹窗不要去改 `ViewSvgText.vue`。

列头柜在 `utils/virtualCabinet.js`：

- 测点名最后一个 `_` 前面是设备名（`splitNameByLastUnderscore`）。
- 分组数 `>= 2` 才启用虚拟柜（`shouldUseVirtualCabinetLayer`）。
- 不写 `monitor_list`。
- 节点 `kind: 'virtualCabinet'`。

## 改的时候不要动

这些已经验证过。这次需求不是它们，就不要改：

- `ISMBase.vue` 里 `require.context` 扫描 `ISMComponents` 的注册。新 `.vue` 放进那个目录会被扫进去，整页报 `comp.default.data is not a function`。新型号写进 `dcEnvSymbols.js`，用现成的一行 `createPowerSymbol`。
- `ViewSvgText.vue`
- `store/ISM/actions.js` 的 `parseRawPageLayerFields`、`applyNavContextToPageConfig`
- 航信端口 `7080` / `8081`，寻智端口 `7081` / `8082`，以及两边的项目名

## 新的一间机房

1. 先定是透视机房，还是组态平面柜。
2. 透视机房：按上面的分件、宽高深、通道间距和相机限制做。每柜独立网格，奇偶行门朝向相反。
3. 组态柜：外框用 `view-svg-dc-rack`，或一张按原图比例放的 `ism-view-png-image`。设备层用图，标牌用 `view-svg-text`。每个格子带上完整的 `animate`（`selected` 和 `animateElement` 都是数组）。
4. 测点放进 `dataBind`，用现成弹窗。列头柜按最后一个 `_` 拆。
5. 不要用 `serverrack` 色块，不要把 `末端设备/Rack.png` 当服务器柜。

# 中航信画面热修包（2026-10-12）

对应代码：`235275ae`（分支 `cursor/fix-scada-nav-and-point-popup-64b6`）

本包只含麒麟静态 `ism_server`、`web/dist/index.html`、`web/dist/static/js`、`web/dist/static/css`、`web/dist/css`。  
不含地图、上传文件、三维模型、自定义图元和数据库。套用时也不会删这些目录。

## 这次改了什么

1. 左侧菜单（首页 / 产品中心 / 配电室5A1 / 配电室5A2）按绑定的子页切换，不再闪一下后提示「已在目标页面」。
2. 总览画面按钮「打开网页 / 内部页面 / 一次图界面 / 配电室5A1 / 弹窗」打开对应子页，不再提示「弹窗目标与当前页相同，请绑定具体子页」。
3. 一次图悬停测点卡片停在原地，点击卡片可以钉住，标题栏可以关闭，同一时间只有一张。实时值按设备 UUID 查询，A/B/C 相电流应与数据仓库一致。

## 现场套用

主包目录一般是 `/opt/ISM/ism-release-oceanbase-20260709`。若现场目录名是 `20260708`，把下面路径换成实际目录。

```bash
unzip ism-patch-kylin-ui-20261012-235275ae.zip
cd ism-patch-kylin-ui-20261012-235275ae
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
```

脚本会先停服务，把当前 `ism_server`、`index.html`、`static/js`、`static/css` 备份到主包下的 `backup/ui-20261012-<时间>/`，再替换并启动。  
不重置 admin 密码，不跑数据库脚本。

## 复测（必须 Ctrl+F5）

浏览器打开 `http://172.31.4.1:7090`，强制刷新后再测：

1. 运行态 `#/AppRun/b620dc75-1da2-351b-b89e-05b5cb608886`  
   点左侧「配电室5A1」「配电室5A2」，画面应切到对应一次图，不应再提示「已在目标页面」。
2. 总览「中航信数据中心电力监控系统」上，已配置弹窗的按钮应弹出配电室5A1，不应再提示「弹窗目标与当前页相同」。
3. 一次图上悬停「配电室5A1_3」一类图元：卡片不跟着鼠标跑；点卡片后钉住；点 × 关闭；不会叠多张。  
   核对 A/B/C 相电流和数据仓库（例如 A 相 2.80 A、B 相 2.250 A、C 相 19.50 A，以当时实时值为准）。

## 回退

在主包目录执行，把 `<时间>` 换成 `backup/` 下实际目录名：

```bash
bash stop-all.sh
cp backup/ui-20261012-<时间>/ism_server_user/ism_server ism_server_user/ism_server
cp backup/ui-20261012-<时间>/web/dist/index.html web/dist/index.html
rm -rf web/dist/static/js web/dist/static/css
mkdir -p web/dist/static/js web/dist/static/css
cp -a backup/ui-20261012-<时间>/web/dist/static/js/. web/dist/static/js/
cp -a backup/ui-20261012-<时间>/web/dist/static/css/. web/dist/static/css/
bash start-all.sh
```

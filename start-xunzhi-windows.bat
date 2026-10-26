﻿@echo off
chcp 65001 >nul
cd /d "%~dp0"
echo.
echo 寻智机房动环 — Windows 启动
echo 请先把整个文件夹从优盘复制到本机硬盘，再双击本文件。
echo 不要在优盘上直接编译。
echo.

where node >nul 2>&1
if errorlevel 1 (
  echo 未找到 Node.js。请安装 Node.js 20 或更高，并加入 PATH。
  echo https://nodejs.org/
  pause
  exit /b 1
)

if not exist "artifacts" mkdir artifacts

if not exist "ism_server_user\ism_server.exe" (
  echo 没有现成的 ism_server.exe，改为在本机编译。需要 Go 1.22 和 MinGW-w64 的 gcc。
  where go >nul 2>&1
  if errorlevel 1 (
    echo 未找到 Go。https://go.dev/dl/
    pause
    exit /b 1
  )
  where gcc >nul 2>&1
  if errorlevel 1 (
    echo 未找到 gcc。https://github.com/brechtsanders/winlibs_mingw/releases
    pause
    exit /b 1
  )
  pushd ism_server_user
  set CGO_ENABLED=1
  go build -mod=vendor -o ism_server.exe .
  if errorlevel 1 (
    echo 后端编译失败。
    popd
    pause
    exit /b 1
  )
  popd
)

echo 启动后端 8082 ...
start "ism-backend-8082" cmd /c "cd /d %~dp0ism_server_user && ism_server.exe >> ..\artifacts\ism_be.log 2>&1"

if not exist "ism-front-end-v2\node_modules\.bin\vue-cli-service.cmd" (
  echo 第一次运行，正在安装前端依赖。需要联网，可能要十几分钟。
  pushd ism-front-end-v2
  call npm install
  if errorlevel 1 (
    echo npm install 失败。确认这台电脑能访问 npm，并且文件夹不在优盘上。
    popd
    pause
    exit /b 1
  )
  popd
)

echo 启动前端 7081，第一次编译要几分钟 ...
start "ism-frontend-7081" cmd /c "cd /d %~dp0ism-front-end-v2 && set NODE_OPTIONS=--max-old-space-size=8192 --openssl-legacy-provider && npx vue-cli-service serve --port 7081 >> ..\artifacts\ism_fe.log 2>&1"

echo.
echo 后端  http://127.0.0.1:8082
echo 前端  http://127.0.0.1:7081
echo 登录  admin / 123456
echo 本体  http://127.0.0.1:7081/#/xunzhi/jf/ontology
echo 智能体 http://127.0.0.1:7081/#/xunzhi/jf/agent
echo.
echo 用 Cursor 打开复制到硬盘上的这个文件夹即可继续改代码。
echo 日志在 artifacts\ism_be.log 和 artifacts\ism_fe.log
pause

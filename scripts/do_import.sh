#!/usr/bin/env bash
# 一键启动后端 + 模拟器 + 导入项目
set -e
PROJECT_ROOT="/Users/yunfanleo/cursorProjects/ISM源码"

# 先杀掉旧的
pkill -9 ism_server 2>/dev/null || true
sleep 1

# 启动后端
cd "$PROJECT_ROOT/ism_server_user"
./ism_server > /tmp/ism_backend_import.log 2>&1 &
BACKEND_PID=$!
echo "backend PID=$BACKEND_PID"

# 等待后端就绪
for i in $(seq 1 15); do
  sleep 1
  if curl -s -o /dev/null http://127.0.0.1:8081/ 2>/dev/null; then
    echo "backend ready after ${i}s"
    break
  fi
done

# 确认模拟器在跑
if ! lsof -i :502 -sTCP:LISTEN > /dev/null 2>&1; then
  cd "$PROJECT_ROOT"
  python3 scripts/modbus_simulator.py > /tmp/modbus_sim_import.log 2>&1 &
  echo "simulator started"
  sleep 3
fi

# 执行导入
cd "$PROJECT_ROOT"
python3 scripts/import_1a_project.py
echo "IMPORT DONE"

#!/bin/bash
# ISM 系统自动重启守护脚本
# 监控 modbus_simulator v7 和 ism_server，任一死亡自动重启

SIM_LOG="/tmp/sim_watch.log"
BE_LOG="/tmp/ism_watch.log"
ISM_DIR="/Users/yunfanleo/cursorProjects/ISM源码"
SIM_PY="$ISM_DIR/scripts/modbus_simulator.py"
BE_BIN="$ISM_DIR/ism_server_user/ism_server"

SIM_PID=""
BE_PID=""
RESTART_COUNT=0

log() { echo "[$(date '+%H:%M:%S')] $1"; }

start_sim() {
    log "Starting simulator..."
    cd "$ISM_DIR" && nohup python3 -u "$SIM_PY" > "$SIM_LOG" 2>&1 &
    SIM_PID=$!
    sleep 8
    if kill -0 $SIM_PID 2>/dev/null && lsof -i :502 2>/dev/null | grep -q LISTEN; then
        log "Simulator started OK (PID=$SIM_PID)"
        return 0
    else
        log "ERROR: Simulator failed to start"
        return 1
    fi
}

start_be() {
    log "Starting backend..."
    cd "$ISM_DIR/ism_server_user" && nohup ./ism_server > "$BE_LOG" 2>&1 &
    BE_PID=$!
    sleep 5
    if kill -0 $BE_PID 2>/dev/null; then
        log "Backend started OK (PID=$BE_PID)"
        return 0
    else
        log "ERROR: Backend failed to start"
        return 1
    fi
}

check_db() {
    local nz=$(sqlite3 "$ISM_DIR/ism_server_user/data/db/ism.db" \
        "SELECT COUNT(*) FROM device_real_data WHERE cast(value AS REAL)!=0 AND value!='';" 2>/dev/null)
    echo "${nz:-0}"
}

# ─── 主循环 ────────────────────────────────────────────

log "=== ISM Watchdog started ==="

while true; do
    # 检查并重启模拟器
    if [ -z "$SIM_PID" ] || ! kill -0 $SIM_PID 2>/dev/null; then
        RESTART_COUNT=$((RESTART_COUNT + 1))
        log "Simulator died (restart #$RESTART_COUNT)"
        start_sim || { log "FATAL: cannot start sim"; sleep 10; continue; }
        # 模拟器重启后，重启后端以建立新连接
        kill $BE_PID 2>/dev/null
        BE_PID=""
    fi

    # 检查并重启后端
    if [ -z "$BE_PID" ] || ! kill -0 $BE_PID 2>/dev/null; then
        log "Backend died (restart #$RESTART_COUNT)"
        start_be || { log "FATAL: cannot start BE"; sleep 10; continue; }
    fi

    # 每 60 秒报告一次
    if [ $(( $(date +%s) % 60 )) -lt 5 ]; then
        NZ=$(check_db)
        log "UP: sim=$SIM_PID be=$BE_PID  non_zero=$NZ  restarts=$RESTART_COUNT"
    fi

    sleep 10
done

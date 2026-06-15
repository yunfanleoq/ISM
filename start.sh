#!/bin/bash
# ============================================================
# ISM Web组态软件 - 一键启动脚本
# 用途: 同时启动后端 (Go Server)、模拟器、前端 (Vue Dev Server)
# 用法: ./start.sh [backend|simulator|frontend|all]
# ============================================================

set -e

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$PROJECT_ROOT/ism_server_user"
FRONTEND_DIR="$PROJECT_ROOT/ism-front-end-v2"
SIMULATOR_SCRIPT="$PROJECT_ROOT/scripts/modbus_simulator.py"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info()  { echo -e "${BLUE}[INFO]${NC}  $1"; }
log_ok()    { echo -e "${GREEN}[OK]${NC}    $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 后端 PID 文件
BACKEND_PID_FILE="$PROJECT_ROOT/.backend.pid"

# 清理函数 - 退出时停止所有服务
cleanup() {
    echo ""
    log_info "正在停止所有服务..."

    # 停止后端
    if [ -f "$BACKEND_PID_FILE" ]; then
        BACKEND_PID=$(cat "$BACKEND_PID_FILE")
        if kill -0 "$BACKEND_PID" 2>/dev/null; then
            kill "$BACKEND_PID" 2>/dev/null
            log_ok "后端服务已停止 (PID: $BACKEND_PID)"
        fi
        rm -f "$BACKEND_PID_FILE"
    fi

    # 停止前端 (vue-cli-service 的子进程)
    if [ ! -z "$FRONTEND_PID" ]; then
        # 先杀子进程，再杀主进程
        pkill -P $FRONTEND_PID 2>/dev/null
        kill $FRONTEND_PID 2>/dev/null
        log_ok "前端服务已停止 (PID: $FRONTEND_PID)"
    fi

    log_info "所有服务已停止。"
    exit 0
}

# 注册退出信号
trap cleanup SIGINT SIGTERM EXIT

# -----------------------------------------------------------
# 检查并启动 OceanBase（当 dbtype=4 时）
# -----------------------------------------------------------
check_oceanbase() {
    local DBTYPE=$(grep "^dbtype=" "$BACKEND_DIR/conf/app.conf" 2>/dev/null | cut -d= -f2)
    
    if [ "$DBTYPE" != "4" ]; then
        return 0
    fi

    log_info "检测到 dbtype=4，检查 OceanBase 容器状态..."

    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^oceanbase$"; then
        log_ok "OceanBase 容器已在运行"
        return 0
    fi

    if docker ps -a --format '{{.Names}} {{.Status}}' 2>/dev/null | grep -q "^oceanbase"; then
        log_info "OceanBase 容器已停止，正在启动..."
        docker start oceanbase > /dev/null 2>&1
    else
        log_error "OceanBase 容器不存在！请先创建容器："
        echo "  docker run -d --name oceanbase --ulimit nofile=65536:65536 --ulimit nproc=65536:65536 \\"
        echo "    -p 2881:2881 -e MODE=mini -e OB_MEMORY_LIMIT=8G \\"
        echo "    -e OB_DATAFILE_SIZE=10G -e OB_LOG_DISK_SIZE=5G \\"
        echo "    -e OB_CLUSTER_NAME=ism_cluster -e OB_TENANT_NAME=ism_tenant \\"
        echo "    -e OB_TENANT_PASSWORD=ism2024! oceanbase/oceanbase-ce:latest"
        return 1
    fi

    # 等待 OceanBase 就绪
    log_info "等待 OceanBase 就绪..."
    for i in $(seq 1 60); do
        if docker exec oceanbase obclient -h 127.0.0.1 -P 2881 -uroot@ism_tenant -p'ism2024!' -e "SELECT 1" > /dev/null 2>&1; then
            log_ok "OceanBase 已就绪 (端口 2881)"
            return 0
        fi
        sleep 2
    done
    log_error "OceanBase 启动超时，请检查容器日志: docker logs oceanbase"
    return 1
}
start_backend() {
    log_info "正在启动后端服务..."

    if [ ! -d "$BACKEND_DIR" ]; then
        log_error "后端目录不存在: $BACKEND_DIR"
        exit 1
    fi

    cd "$BACKEND_DIR"

    # 检查二进制文件是否存在
    if [ ! -f "./ism_server" ]; then
        log_error "后端二进制文件不存在: ./ism_server"
        log_info "请先编译后端: cd $BACKEND_DIR && go build -o ism_server"
        exit 1
    fi

    # 确保必要的目录存在
    mkdir -p data/auth static/HistoryData static/reportTemplete static/RecordVideo logs

    # 启动后端 (后台运行)
    ./ism_server &
    BACKEND_PID=$!
    echo $BACKEND_PID > "$BACKEND_PID_FILE"

    log_ok "后端服务已启动 (PID: $BACKEND_PID)"

    # 等待后端就绪
    log_info "等待后端服务就绪..."
    for i in $(seq 1 30); do
        if curl -s http://127.0.0.1:8081/ > /dev/null 2>&1; then
            log_ok "后端服务就绪 (http://127.0.0.1:8081)"
            break
        fi
        if [ $i -eq 30 ]; then
            log_warn "后端服务可能启动较慢，请稍候..."
        fi
        sleep 1
    done

    cd "$PROJECT_ROOT"
}

# -----------------------------------------------------------
# 启动前端
# -----------------------------------------------------------
start_frontend() {
    log_info "正在启动前端开发服务器..."

    if [ ! -d "$FRONTEND_DIR" ]; then
        log_error "前端目录不存在: $FRONTEND_DIR"
        exit 1
    fi

    cd "$FRONTEND_DIR"

    # 检查 node_modules 是否存在
    if [ ! -d "node_modules" ]; then
        log_warn "node_modules 不存在，正在安装依赖..."
        npm install --legacy-peer-deps
    fi

    # 检查 cross-env 是否安装（Node 16 兼容版本）
    if ! npx --no-install cross-env --version > /dev/null 2>&1; then
        log_warn "cross-env 未安装，正在安装..."
        npm install cross-env@7 --save-dev --legacy-peer-deps
    fi

    log_ok "前端依赖已就绪"

    # 启动前端 (新进程组，方便一起杀掉)
    log_info "前端开发服务器启动中 (端口: 7080)..."
    log_info "浏览器将自动打开: http://localhost:7080"
    echo ""

    # Node 26 + 20G 内存 + OpenSSL 兼容（webpack 4 需要）
    # 两个参数缺一不可：缺内存会 OOM 崩溃，缺 openssl-legacy-provider 会在 95% 报错
    NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" npm run serve &
    FRONTEND_PID=$!

    # 等待前端进程
    wait $FRONTEND_PID 2>/dev/null
}

# -----------------------------------------------------------
# 启动 Modbus 模拟器
# -----------------------------------------------------------
start_simulator() {
    log_info "正在启动 Modbus TCP 模拟器..."

    if [ ! -f "$SIMULATOR_SCRIPT" ]; then
        log_error "模拟器脚本不存在: $SIMULATOR_SCRIPT"
        exit 1
    fi

    # 检查 pymodbus 是否安装
    if ! python3 -c "import pymodbus" 2>/dev/null; then
        log_warn "pymodbus 未安装，正在安装..."
        pip3 install "pymodbus<3.0" --break-system-packages
    fi

    # 检查是否已经运行
    if lsof -i :502 -sTCP:LISTEN > /dev/null 2>&1; then
        log_warn "端口 502 已被占用，跳过模拟器启动"
        return
    fi

    cd "$PROJECT_ROOT"
    nohup python3 "$SIMULATOR_SCRIPT" > /tmp/modbus_sim.log 2>&1 &
    SIMULATOR_PID=$!

    log_ok "Modbus 模拟器已启动 (PID: $SIMULATOR_PID, 端口: 502)"

    # 等待模拟器就绪
    log_info "等待模拟器就绪..."
    for i in $(seq 1 10); do
        if lsof -i :502 -sTCP:LISTEN > /dev/null 2>&1; then
            log_ok "Modbus 模拟器就绪 (127.0.0.1:502)"
            return
        fi
        sleep 1
    done
    log_warn "模拟器启动可能较慢，请检查 /tmp/modbus_sim.log"
}

# -----------------------------------------------------------
# 主逻辑
# -----------------------------------------------------------
MODE="${1:-all}"

echo ""
echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}   ISM Web组态软件 - 启动脚本${NC}"
echo -e "${BLUE}   Version: V3.01.RC07${NC}"
echo -e "${BLUE}============================================================${NC}"
echo ""

case "$MODE" in
    backend)
        check_oceanbase || exit 1
        start_backend
        log_info "后端运行中... 按 Ctrl+C 停止"
        wait
        ;;
    simulator)
        start_simulator
        log_info "模拟器运行中..."
        wait
        ;;
    frontend)
        start_frontend
        ;;
    all)
        # 先启动模拟器
        start_simulator
        sleep 1

        # 检查数据库
        check_oceanbase || exit 1

        # 启动后端
        start_backend

        # 等待后端完全就绪
        sleep 2

        # 再启动前端
        start_frontend
        ;;
    *)
        log_error "未知参数: $MODE"
        echo "用法: ./start.sh [backend|simulator|frontend|all]"
        echo "  backend    - 仅启动后端服务 (端口 8081)"
        echo "  simulator  - 仅启动 Modbus TCP 模拟器 (端口 502) + HTTP API (端口 5040)"
        echo "  frontend   - 仅启动前端开发服务器 (端口 7080)"
        echo "  all        - 启动模拟器 + 后端 + 前端 (默认)"
        echo ""
        echo "监控页面:"
        echo "  模拟器界面: http://localhost:7080/#/SimulatorMonitor"
        echo "  API 测试:   curl http://localhost:5040/api/slaves"
        exit 1
        ;;
esac

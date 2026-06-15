#!/usr/bin/env python3
"""
ISM 系统 24 小时稳定性守护进程 v3
- 使用 subprocess.Popen 管理模拟器和后端
- 任一进程死亡自动重启
- 记录日志和数据库状态
"""
import subprocess, time, sys, os, signal, socket, sqlite3
from datetime import datetime

# 忽略无关信号
for sig in [signal.SIGPIPE, signal.SIGCHLD]:
    try: signal.signal(sig, signal.SIG_IGN)
    except: pass

ISM_DIR = "/Users/yunfanleo/cursorProjects/ISM源码"
BE_DIR = f"{ISM_DIR}/ism_server_user"
SIM_LOG = "/tmp/ism_wd_sim.log"
BE_LOG = "/tmp/ism_wd_be.log"
WD_LOG = "/tmp/ism_wd.log"
DB = f"{BE_DIR}/data/db/ism.db"

start_time = datetime.now()
restart_count = 0


def wlog(msg):
    t = datetime.now().strftime("%m-%d %H:%M:%S")
    line = f"[{t}] {msg}"
    try:
        with open(WD_LOG, "a") as f:
            f.write(line + "\n")
    except: pass
    print(line, flush=True)


def uptime_h():
    return f"{(datetime.now() - start_time).total_seconds()/3600:.1f}h"


def port_ok():
    try:
        s = socket.socket(); s.settimeout(3)
        s.connect(("127.0.0.1", 502)); s.close()
        return True
    except: return False


def db_status():
    try:
        c = sqlite3.connect(DB, timeout=5)
        cur = c.cursor()
        cur.execute("SELECT COUNT(*) FROM device_real_data WHERE cast(value AS REAL)!=0 AND value!=''")
        nz = cur.fetchone()[0]
        cur.execute("SELECT MAX(updated_at) FROM device_real_data")
        lat = cur.fetchone()[0] or "N/A"
        c.close()
        return nz, lat
    except Exception as e:
        return 0, f"ERR:{e}"


def killall():
    global sim_p, be_p
    for p in [sim_p, be_p]:
        if p:
            try: p.kill()
            except: pass
    sim_p = be_p = None
    time.sleep(2)
    os.system("lsof -ti :502 2>/dev/null | xargs kill -9 2>/dev/null")
    time.sleep(2)


def start_sim():
    global sim_p
    try:
        sim_p = subprocess.Popen(
            [sys.executable, "-u", "scripts/modbus_simulator.py"],
            cwd=ISM_DIR,
            stdout=open(SIM_LOG, "a"),
            stderr=subprocess.STDOUT,
        )
    except Exception as e:
        wlog(f"SIM start error: {e}")
        return False

    time.sleep(8)
    if sim_p.poll() is not None:
        wlog(f"SIM died (rc={sim_p.returncode})")
        return False
    if not port_ok():
        wlog("SIM port not open")
        sim_p.kill()
        return False
    wlog(f"SIM OK (PID={sim_p.pid})")
    return True


def start_be():
    global be_p
    try:
        be_p = subprocess.Popen(
            ["./ism_server"],
            cwd=BE_DIR,
            stdout=open(BE_LOG, "a"),
            stderr=subprocess.STDOUT,
        )
    except Exception as e:
        wlog(f"BE start error: {e}")
        return False

    time.sleep(8)
    if be_p.poll() is not None:
        wlog(f"BE died (rc={be_p.returncode})")
        return False
    wlog(f"BE OK (PID={be_p.pid})")
    return True


def full_restart():
    global restart_count
    restart_count += 1
    wlog(f"=== RESTART #{restart_count} (uptime {uptime_h()}) ===")
    killall()
    if not start_sim():
        wlog("FATAL: SIM start failed")
        return False
    if not start_be():
        wlog("FATAL: BE start failed")
        sim_p.kill()
        return False
    return True


# ─── Main ─────────────────────────────────────────────

wlog("=" * 50)
wlog(f"ISM Watchdog v3 | Start: {start_time.strftime('%Y-%m-%d %H:%M:%S')}")
wlog("=" * 50)

sim_p = None
be_p = None

# 初始启动
if not full_restart():
    wlog("Retrying in 30s...")
    time.sleep(30)
    if not full_restart():
        wlog("FATAL: Cannot start system. Exiting.")
        sys.exit(1)

last_report = time.time()
last_db = time.time()
REPORT = 300  # 5 分钟
DB_INT = 60   # 1 分钟

while True:
    try:
        now = time.time()

        # 进程监控
        if sim_p and sim_p.poll() is not None:
            wlog(f"SIM DEAD (rc={sim_p.returncode})")
            full_restart()
        elif be_p and be_p.poll() is not None:
            wlog(f"BE DEAD (rc={be_p.returncode})")
            full_restart()

        # 数据库监控
        if now - last_db > DB_INT:
            nz, lat = db_status()
            last_db = now
            if now - last_report > REPORT:
                wlog(f"STATUS [{uptime_h()}] "
                     f"nz={nz} @{lat} restarts={restart_count}")
                last_report = now

        time.sleep(10)

    except KeyboardInterrupt:
        wlog(f"Shutdown. Total restarts: {restart_count}, uptime: {uptime_h()}")
        killall()
        sys.exit(0)
    except Exception as e:
        wlog(f"Watchdog error: {type(e).__name__}: {e}")
        time.sleep(15)

#!/usr/bin/env python3
"""启动 ISM 系统：模拟器 → 后端"""
import subprocess, time, sys, os

def start():
    # 启动模拟器
    sim = subprocess.Popen(
        [sys.executable, "-u", "scripts/modbus_simulator.py"],
        cwd="/Users/yunfanleo/cursorProjects/ISM源码",
        stdout=open("/tmp/sim_sys.log", "a"),
        stderr=subprocess.STDOUT,
    )
    print(f"Sim PID: {sim.pid}")
    time.sleep(8)
    if sim.poll() is not None:
        print(f"Sim died immediately (rc={sim.returncode})")
        return

    # 启动后端
    be = subprocess.Popen(
        ["./ism_server"],
        cwd="/Users/yunfanleo/cursorProjects/ISM源码/ism_server_user",
        stdout=open("/tmp/be_sys.log", "a"),
        stderr=subprocess.STDOUT,
    )
    print(f"BE PID: {be.pid}")
    time.sleep(5)
    if be.poll() is not None:
        print(f"BE died immediately (rc={be.returncode})")
        sim.kill()
        return

    print("System started. Waiting for exit...")
    # 等待其中任何一个退出
    import select
    while True:
        if sim.poll() is not None:
            print(f"Sim exited with rc={sim.returncode}")
            be.kill()
            break
        if be.poll() is not None:
            print(f"BE exited with rc={be.returncode}")
            sim.kill()
            break
        time.sleep(2)

if __name__ == "__main__":
    start()

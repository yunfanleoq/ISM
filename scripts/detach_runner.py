#!/usr/bin/env python3
"""
完全脱离终端的启动器 - 使用 os.fork() 实现真 daemon
"""
import os, sys, time, subprocess, signal

SCRIPT = sys.argv[1] if len(sys.argv) > 1 else "scripts/watchdog.py"
PIDFILE = "/tmp/ism_daemon.pid"

# 忽略无关信号
for sig in [signal.SIGPIPE, signal.SIGCHLD]:
    try: signal.signal(sig, signal.SIG_IGN)
    except: pass


def daemonize():
    """经典的两次 fork 守护进程化"""
    # 第一次 fork - 脱离父进程
    pid = os.fork()
    if pid > 0:
        # 父进程退出
        sys.exit(0)

    # 创建新会话，脱离控制终端
    os.setsid()

    # 第二次 fork - 彻底脱离会话领导地位
    pid = os.fork()
    if pid > 0:
        sys.exit(0)

    # 重定向标准文件描述符到 /dev/null
    fd = os.open("/dev/null", os.O_RDWR)
    os.dup2(fd, 0)
    os.dup2(fd, 1)
    os.dup2(fd, 2)
    os.close(fd)

    # 写 PID 文件
    with open(PIDFILE, "w") as f:
        f.write(str(os.getpid()))

    # 设定工作目录
    os.chdir("/Users/yunfanleo/cursorProjects/ISM源码")


if __name__ == "__main__":
    daemonize()
    # 现在我们已经完全脱离终端，执行目标脚本
    os.execv(sys.executable, [sys.executable, "-u", SCRIPT])

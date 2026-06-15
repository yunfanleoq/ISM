"""测试模拟器并发连接能力"""
import socket, struct, time, sys, os

HOST = "127.0.0.1"
PORT = 502

# 检查模拟器 PID
sim_pid = int(sys.argv[1]) if len(sys.argv) > 1 else 0
print(f"Testing simulator PID={sim_pid}")

error_count = 0
ok_count = 0

for i in range(30):
    s = socket.socket()
    s.settimeout(5)
    try:
        s.connect((HOST, PORT))
        # Read HR from address 0, 10 registers, slave 1
        req = struct.pack('>HHHBBHH', i, 0, 6, 1, 3, 0, 10)
        s.sendall(req)
        resp = s.recv(512)
        print(f"  iter={i}: got {len(resp)} bytes OK")
        ok_count += 1
        s.close()
    except Exception as e:
        print(f"  iter={i}: ERROR {e}")
        error_count += 1
    time.sleep(0.05)

print(f"\nResult: {ok_count} OK, {error_count} errors")

# 检查模拟器是否还活着
try:
    os.kill(sim_pid, 0)
    print(f"SIM PID {sim_pid}: ALIVE")
except:
    print(f"SIM PID {sim_pid}: DEAD")

#!/usr/bin/env python3
"""一键导入组态大屏：启动后端 → 登录 → 创建 DisplayModel → 保存页面"""
import json, base64, time, subprocess, sys
import requests
from requests.adapters import HTTPAdapter

BASE = 'http://[::1]:8081'
PROJECT = 'e308e378-ed94-1264-197b-33e535b812b8'

# Create session with no connection reuse
session = requests.Session()
session.headers.update({'Connection': 'close'})

def wait_server(pid, timeout=15):
    for i in range(timeout * 2):
        try:
            # Don't use session for health check - use raw requests
            r = requests.get(f'{BASE}/', timeout=2, headers={'Connection': 'close'})
            r.close()
            print(f"  Server ready (attempt {i})")
            return True
        except:
            time.sleep(0.5)
    return False

def login():
    for attempt in range(5):
        try:
            r = requests.post(f'{BASE}/login', json={'username': 'admin', 'password': '123456'}, timeout=5, headers={'Connection': 'close'})
            r.close()
            print(f"  登录响应[{attempt}]: status={r.status_code} text={r.text[:80]}")
            if r.status_code != 200:
                time.sleep(1)
                continue
            d = r.json()
            if d.get('code') == 1000:
                token = d['data']['token']
                print(f"  登录成功 token={token[:20]}...")
                return token
            print(f"  登录失败 code={d.get('code')} msg={d.get('message','')}")
        except Exception as e:
            print(f"  登录异常[{attempt}]: {e}")
            time.sleep(1)
    return None

# 1. 启动后端
print("=== 启动后端 ===")
proc = subprocess.Popen(
    ['./ism_server'], 
    cwd='/Users/yunfanleo/cursorProjects/ISM源码/ism_server_user',
    stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL
)
print(f"  PID={proc.pid}")
time.sleep(1)  # Give server more time

# 2. 等待服务就绪
if not wait_server(proc.pid, timeout=20):
    print("  服务启动超时!")
    proc.terminate()
    sys.exit(1)

# 3. 登录
token = login()
if not token:
    sys.exit(1)

h = {'Content-Type': 'application/json', 'Authorization': token, 'Connection': 'close'}

# 4. 创建 DisplayModel
r = requests.post(f'{BASE}/displayModelAdd', json={
    'name': '1A配电室监控总览',
    'description': '1A配电室 Modbus 模拟监控大屏',
    'display_type': 1,
    'project_uuid': PROJECT,
}, headers=h, timeout=5)
d = r.json()
muid = d.get('data', {}).get('display_model_uid', '')
print(f"  displayModelAdd: code={d.get('code')} uid={muid[:20] if muid else 'NONE'}")

# If already exists, find from list
if not muid:
    r = requests.post(f'{BASE}/displayModelList', json={'DisplayType': 1}, headers=h, timeout=5)
    for m in r.json().get('list', []):
        if '1A配电室' in m.get('name', ''):
            muid = m['display_model_uid']
            print(f"  使用已有模型: {muid[:20]}")
            break

if not muid:
    print("  无法获取模型 UUID!")
    sys.exit(1)

# 5. 添加页面
r = requests.post(f'{BASE}/DisplayModelPageAdd', json={
    'modelUuid': muid,
    'name': '总览',
    'size': '1920x1080',
    'pageType': 1,
    'isLogin': 0,
}, headers=h, timeout=5)
print(f"  DisplayModelPageAdd: code={r.json().get('code')}")

# 6. 构建组件
components = [
    {
        'type': 'DvBorderBox1', 'identifier': 'hdr', 'name': 'hdr',
        'detail': {
            'style': {'width': 1920, 'height': 80, 'left': 0, 'top': 0},
            'config': {'title': '1A配电室 - 监控总览 (76台设备)', 'titleColor': '#00d4ff'}
        }
    },
    {
        'type': 'DeviceTree', 'identifier': 'tree', 'name': 'tree',
        'detail': {
            'style': {'width': 280, 'height': 1000, 'left': 0, 'top': 80},
            'config': {'title': '设备拓扑'}
        }
    },
    {
        'type': 'alarmList', 'identifier': 'alarm', 'name': 'alarm',
        'detail': {
            'style': {'width': 330, 'height': 500, 'left': 1590, 'top': 80},
            'config': {'title': '实时告警', 'count': 10}
        }
    },
    {
        'type': 'RealDataTable', 'identifier': 'tbl', 'name': 'tbl',
        'detail': {
            'style': {'width': 1300, 'height': 700, 'left': 290, 'top': 90},
            'config': {
                'columns': [
                    {'title': '设备编号', 'dataIndex': 'name'},
                    {'title': 'AB线电压(V)', 'dataIndex': 'vab'},
                    {'title': 'A相电流(A)', 'dataIndex': 'ia'},
                    {'title': '功率因数', 'dataIndex': 'pf'},
                ],
                'rowNum': 76,
                'headerBGC': '#0a1628',
                'evenRowBGC': '#0d1f35',
                'align': ['center'],
            }
        }
    },
]

comp_b64 = base64.b64encode(json.dumps(components, ensure_ascii=False).encode()).decode()
layer = json.dumps({
    'backColor': '#0a1628', 'backgroundImage': '',
    'widthHeightRatio': '', 'width': 1920, 'height': 1080, 'autoSize': 0,
}, ensure_ascii=False)

# 7. 保存层数据
r = requests.post(f'{BASE}/saveDisplayModelLayerData', json={
    'modelUuid': muid,
    'pageName': '总览',
    'layer': layer,
    'components': comp_b64,
    'is_home': 1,
    'page_type': 1,
}, headers=h, timeout=10)
print(f"  saveDisplayModelLayerData: code={r.json().get('code')}")

# 8. 验证
r = requests.post(f'{BASE}/displayModelList', json={'DisplayType': 1}, headers=h, timeout=5)
apps = r.json().get('list', [])
sim_apps = [m for m in apps if '1A配电室监控' in m.get('name', '')]
print(f"\n✅ 导入完成！模拟器项目组态应用数: {len(sim_apps)}")
for m in apps:
    print(f"  {m.get('name', '?')}")

print(f"\n后端 PID={proc.pid} 仍在运行，按 Ctrl+C 停止")
try:
    proc.wait()
except KeyboardInterrupt:
    pass

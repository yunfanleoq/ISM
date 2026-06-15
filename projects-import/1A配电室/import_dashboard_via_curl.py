#!/usr/bin/env python3
"""一键导入组态大屏（使用 curl 绕过 Python HTTP 库兼容问题）"""
import json, base64, time, subprocess, sys, os

BASE = 'http://[::1]:8081'
PROJECT = 'e308e378-ed94-1264-197b-33e535b812b8'
SERVER_DIR = '/Users/yunfanleo/cursorProjects/ISM源码/ism_server_user'

def curl(method, path, data=None, token=None):
    """Execute curl command and return parsed JSON"""
    cmd = ['curl', '-s', '--max-time', '5', '--globoff',
           f'{BASE}{path}', '-X', method,
           '-H', 'Content-Type: application/json',
           '-H', f'ProjectUuid: {PROJECT}']
    if token:
        cmd += ['-H', f'Authorization: {token}']
    if data:
        cmd += ['-d', json.dumps(data, ensure_ascii=False)]
    
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
    if result.stderr:
        print(f"  curl stderr: {result.stderr[:100]}")
    if not result.stdout.strip():
        print(f"  curl empty response")
        return None
    try:
        return json.loads(result.stdout)
    except:
        print(f"  curl non-JSON: {result.stdout[:100]}")
        return None

# 1. Kill old server
print("=== 清理旧进程 ===")
subprocess.run(['pkill', '-f', 'ism_server'], capture_output=True)
time.sleep(1)

# 2. Start server
print("=== 启动后端 ===")
proc = subprocess.Popen(
    ['./ism_server'],
    cwd=SERVER_DIR,
    stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
)
print(f"  PID={proc.pid}")

# 3. Wait for server
print("=== 等待服务就绪 ===")
for i in range(20):
    try:
        r = curl('GET', '/')
        if r is not None or subprocess.run(['lsof', '-nP', '-a', '-p', str(proc.pid), '-iTCP:8081', '-sTCP:LISTEN'], 
                                           capture_output=True).returncode == 0:
            print(f"  Ready (attempt {i})")
            break
    except:
        pass
    time.sleep(0.5)

# 4. Login
print("=== 登录 ===")
d = curl('POST', '/login', data={'username': 'admin', 'password': '123456'})
if not d or d.get('code') != 1000:
    print(f"  登录失败: {d}")
    sys.exit(1)
token = d['data']['token']
print(f"  成功 token={token[:20]}...")

# 5. Create DisplayModel
print("=== 创建 DisplayModel ===")
d = curl('POST', '/displayModelAdd', token=token, data={
    'name': '1A配电室监控总览',
    'description': '1A配电室 Modbus 模拟监控大屏',
    'display_type': 1,
    'project_uuid': PROJECT,
})
print(f"  displayModelAdd result: {d}")
muid = d.get('data', {}).get('display_model_uid', '') or d.get('data', {}).get('displayUid', '') if d else ''
if muid:
    print(f"  新建成功 uid={muid[:20]}")
else:
    # Already exists - check list
    print("  检查已有模型...")
    d2 = curl('POST', '/displayModelList', token=token, data={'DisplayType': 1})
    print(f"  displayModelList result: {d2}")
    if d2 and d2.get('code') == 0:
        for m in d2.get('list', []):
            if '1A配电室' in m.get('name', ''):
                muid = m.get('displayUid', '') or m.get('display_model_uid', '')
                print(f"  使用已有 uid={muid[:20]}")
                break

if not muid:
    print(f"  创建失败: {d}")
    sys.exit(1)

# 6. Add page
print("=== 添加页面 ===")
d = curl('POST', '/DisplayModelPageAdd', token=token, data={
    'modelUuid': muid,
    'name': '总览',
    'size': '1920x1080',
    'pageType': 1,
    'isLogin': 0,
})
print(f"  code={d.get('code') if d else 'NONE'}")

# 7. Build components and save
print("=== 保存组件 ===")
components = [
    {'type': 'DvBorderBox1', 'identifier': 'hdr', 'name': 'hdr',
     'detail': {'style': {'width': 1920, 'height': 80, 'left': 0, 'top': 0},
                'config': {'title': '1A配电室 - 监控总览 (76台设备)', 'titleColor': '#00d4ff'}}},
    {'type': 'DeviceTree', 'identifier': 'tree', 'name': 'tree',
     'detail': {'style': {'width': 280, 'height': 1000, 'left': 0, 'top': 80},
                'config': {'title': '设备拓扑'}}},
    {'type': 'alarmList', 'identifier': 'alarm', 'name': 'alarm',
     'detail': {'style': {'width': 330, 'height': 500, 'left': 1590, 'top': 80},
                'config': {'title': '实时告警', 'count': 10}}},
    {'type': 'RealDataTable', 'identifier': 'tbl', 'name': 'tbl',
     'detail': {'style': {'width': 1300, 'height': 700, 'left': 290, 'top': 90},
                'config': {'columns': [
                    {'title': '设备编号', 'dataIndex': 'name'},
                    {'title': 'AB线电压(V)', 'dataIndex': 'vab'},
                    {'title': 'A相电流(A)', 'dataIndex': 'ia'},
                    {'title': '功率因数', 'dataIndex': 'pf'},
                ], 'rowNum': 76, 'headerBGC': '#0a1628', 'evenRowBGC': '#0d1f35', 'align': ['center']}}},
]
comp_b64 = base64.b64encode(json.dumps(components, ensure_ascii=False).encode()).decode()
layer = json.dumps({
    'backColor': '#0a1628', 'backgroundImage': '',
    'widthHeightRatio': '', 'width': 1920, 'height': 1080, 'autoSize': 0,
}, ensure_ascii=False)

d = curl('POST', '/saveDisplayModelLayerData', token=token, data={
    'modelUuid': muid,
    'pageName': '总览',
    'layer': layer,
    'components': comp_b64,
    'is_home': 1,
    'page_type': 1,
})
print(f"  code={d.get('code') if d else 'NONE'}")

# 8. Verify
print("\n=== 验证 ===")
d = curl('POST', '/displayModelList', token=token, data={'DisplayType': 1})
count = len(d.get('list', []))
print(f"  组态应用总数: {count}")
for m in d.get('list', []):
    print(f"    {m.get('name', '?')}")

print(f"\n✅ 大屏导入完成！请刷新 http://localhost:7080")
print(f"   查看 1A配电室-模拟器 项目 → 应用列表")
print(f"\n服务 PID={proc.pid}，Ctrl+C 停止")
if proc.poll() is None:
    print(f"\n  后端仍在运行 PID={proc.pid}")
else:
    print(f"\n  后端已退出 (exit={proc.returncode})")

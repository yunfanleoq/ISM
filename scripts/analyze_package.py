import json
from collections import defaultdict

with open('/Users/yunfanleo/cursorProjects/ISM源码/liu-chang-1A-dev/1A_complete_project_package.json', 'r') as f:
    data = json.load(f)

# Group by muid (device model)
points_by_muid = defaultdict(list)
for rp in data['registerPoints']:
    points_by_muid[rp['muid']].append(rp)

device_map = {d['uuid']: d['name'] for d in data['deviceModels']}

total = 0
for muid, points in points_by_muid.items():
    dev_name = device_map.get(muid, 'Unknown')
    print(f"\n=== {dev_name} (muid={muid}) === {len(points)} points")
    # Further group by what should be AI vs DI based on type
    ai_points = [p for p in points if p['type'] != 'Bool']
    di_points = [p for p in points if p['type'] == 'Bool']
    print(f"  AI (non-Bool): {len(ai_points)}")
    for p in sorted(ai_points, key=lambda x: x['registerAddress']):
        print(f"    addr={p['registerAddress']:3d} | {p['name']:25s} | type={p['type']:6s} | unit={p['unit']:4s} | expr={p['conversionExpression']:6s}")
    print(f"  DI (Bool): {len(di_points)}")
    for p in sorted(di_points, key=lambda x: x['registerAddress']):
        print(f"    addr={p['registerAddress']:3d} | {p['name']:25s} | type={p['type']:6s}")
    total += len(points)

print(f"\nTotal points in package: {total}")

# Database register groups
print("\n=== Database Register Groups ===")
# We'll just list them based on user's info

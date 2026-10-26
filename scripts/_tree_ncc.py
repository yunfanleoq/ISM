#!/usr/bin/env python3
import pymysql
from collections import defaultdict

conn = pymysql.connect(host='127.0.0.1', port=2881,
                       user='root@ism_tenant', password='ism2024!', database='ism')
cur = conn.cursor()
PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'

cur.execute("""SELECT uuid, name, sid, pid, type, muid, status
               FROM monitor_list WHERE project_uuid=%s AND deleted_at IS NULL
               ORDER BY pid, type, name""", (PROJECT_UUID,))
rows = cur.fetchall()
by_sid = {r[2]: r for r in rows}
children = defaultdict(list)
for r in rows:
    children[r[3]].append(r)

# find roots (pid not a known sid)
sids = set(by_sid)
roots = [r for r in rows if r[3] not in sids]
print(f"total nodes={len(rows)} type0={sum(1 for r in rows if r[4]==0)} type1={sum(1 for r in rows if r[4]==1)}")
print(f"roots ({len(roots)}):")

def walk(node, depth):
    uuid, name, sid, pid, t, muid, status = node
    kids = children.get(sid, [])
    t1 = sum(1 for k in kids if k[4]==1)
    t0 = sum(1 for k in kids if k[4]==0)
    print(f"{'  '*depth}[{ 'GRP' if t==0 else 'DEV'}] {name} sid={sid} pid={pid} status={status} kids(t0={t0},t1={t1})")
    if depth < 3:  # don't print all leaf devices
        for k in kids:
            if k[4]==0:
                walk(k, depth+1)
        # print first 3 device leaves
        devs = [k for k in kids if k[4]==1]
        for k in devs[:3]:
            print(f"{'  '*(depth+1)}[DEV] {k[1]} sid={k[2]} status={k[6]} muid={k[5]}")
        if len(devs) > 3:
            print(f"{'  '*(depth+1)}... +{len(devs)-3} more devices")

for r in roots:
    walk(r, 0)
conn.close()

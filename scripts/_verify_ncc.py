#!/usr/bin/env python3
import pymysql, json, base64, uuid as _uuid

conn = pymysql.connect(host='127.0.0.1', port=2881,
                       user='root@ism_tenant', password='ism2024!', database='ism')
cur = conn.cursor()

MODEL_ID = '043135ad-44be-e5d8-89be-3e54883c23a8'
PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'

def page_id_building(sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-bldg-{sid}').hex

def page_id_room(sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-room-{sid}').hex

# 1) online KPI = type1 status=1
cur.execute("""SELECT SUM(type=1) total, SUM(type=1 AND status=1) online
               FROM monitor_list WHERE project_uuid=%s AND deleted_at IS NULL""", (PROJECT_UUID,))
total, online = cur.fetchone()
print(f"[KPI] type=1 total={total}  online(status=1)={online}")

# 2) page counts by category
cur.execute("""SELECT page_name FROM display_model_layer
               WHERE model_id=%s AND deleted_at IS NULL""", (MODEL_ID,))
names = [r[0] for r in cur.fetchall()]
import collections
cats = collections.Counter()
for n in names:
    if n == 'main': cats['overview'] += 1
    elif n.startswith('room-'): cats['room'] += 1
    elif n.startswith('building-'): cats['building'] += 1
    elif n.startswith('floor-'): cats['floor'] += 1
    elif n == 'device-detail': cats['legacy-detail'] += 1
    elif n.startswith('device-'): cats['device'] += 1
    else: cats['other:'+n] += 1
print(f"[PAGES] total active={len(names)}  {dict(cats)}")

# 3) UPS cabinet page -> its floor cards must link to UPS floor pages (not 1A1)
cur.execute("SELECT sid,name FROM monitor_list WHERE project_uuid=%s AND name='UPS柜' AND deleted_at IS NULL", (PROJECT_UUID,))
row = cur.fetchone()
ups_sid, ups_name = row
ups_page = page_id_building(ups_sid)
cur.execute("SELECT components FROM display_model_layer WHERE model_id=%s AND page_id=%s AND deleted_at IS NULL", (MODEL_ID, ups_page))
comp = json.loads(base64.b64decode(cur.fetchone()[0]).decode())
links = set()
title = None
for c in comp['cells']:
    det = c['data']['detail']
    for a in det.get('action', []):
        ln = a.get('link', {}).get('Inside', {}).get('pageUUID')
        if ln: links.add(ln)
    if det.get('name','').startswith('🏢') or '柜' in (det.get('style',{}).get('text','') or ''):
        if det.get('style',{}).get('text','').find('UPS') >= 0:
            title = det['style']['text']
print(f"[UPS] page_id={ups_page[:16]}.. title-found={title!r} distinct_link_targets={len(links)}")
# Resolve link targets to page_name to confirm they are UPS floor groups
if links:
    ph = ','.join(['%s']*len(links))
    cur.execute(f"SELECT page_id,page_name FROM display_model_layer WHERE model_id=%s AND page_id IN ({ph}) AND deleted_at IS NULL", (MODEL_ID, *links))
    tgt = {r[0]: r[1] for r in cur.fetchall()}
    ups_floor_links = [v for v in tgt.values() if v and v.startswith(f'floor-{ups_sid}-U')]
    print(f"[UPS] link page_names -> {sorted(set(tgt.values()))}")
    print(f"[UPS] links to UPS floor pages: {len(ups_floor_links)} (must be >0, none pointing to 1A1)")

# 4) a device's real-time data non-empty (getRealData equivalent)
cur.execute("""SELECT m.uuid, m.name FROM monitor_list m
               WHERE m.project_uuid=%s AND m.type=1 AND m.status=1 AND m.deleted_at IS NULL
               LIMIT 1""", (PROJECT_UUID,))
duuid, dname = cur.fetchone()
cur.execute("SELECT COUNT(*) FROM device_real_data WHERE device_uuid=%s", (duuid,))
cnt = cur.fetchone()[0]
cur.execute("""SELECT name, value, data_unit FROM device_real_data
               WHERE device_uuid=%s AND name IN ('AB线电压','A相电流','总有功功率','频率')
               ORDER BY id LIMIT 6""", (duuid,))
samples_pre = None
samples = cur.fetchall()
print(f"[REALDATA] device={dname} uuid={duuid[:16]}.. rows={cnt}")
for s in samples:
    print(f"           {s[0]} = {s[1]} {s[2] or ''}")

# 4b) overview must show ALL cabinets incl UPS (aggregated, not tiled devices)
cur.execute("SELECT components FROM display_model_layer WHERE model_id=%s AND page_name='main' AND deleted_at IS NULL", (MODEL_ID,))
ov = json.loads(base64.b64decode(cur.fetchone()[0]).decode())['cells']
ov_texts = [c['data']['detail']['style'].get('text','') for c in ov if c.get('shape')=='view-svg-text']
ov_join = ' | '.join(t for t in ov_texts if t)
has_ups = 'UPS柜' in ov_join
has_room = '1楼配电室' in ov_join
tiled_dev = sum(1 for t in ov_texts if t in ('运行中','离线'))  # per-device tiles (should be ~0 now)
print(f"[OVERVIEW] contains UPS柜={has_ups}  contains 配电室={has_room}  per-device-tiles={tiled_dev} (aggregated if ~0)")

# 4c) room (配电室) page exists and links to all 3 cabinets incl UPS
cur.execute("SELECT sid FROM monitor_list WHERE project_uuid=%s AND name='1楼配电室' AND deleted_at IS NULL", (PROJECT_UUID,))
rsid = cur.fetchone()[0]
rpage = page_id_room(rsid)
cur.execute("SELECT components FROM display_model_layer WHERE model_id=%s AND page_id=%s AND deleted_at IS NULL", (MODEL_ID, rpage))
rr = cur.fetchone()
if rr:
    rcells = json.loads(base64.b64decode(rr[0]).decode())['cells']
    rtxt = ' | '.join(c['data']['detail']['style'].get('text','') for c in rcells if c.get('shape')=='view-svg-text')
    rlinks = set()
    for c in rcells:
        for a in c['data']['detail'].get('action', []):
            ln = a.get('link',{}).get('Inside',{}).get('pageUUID')
            if ln: rlinks.add(ln)
    cab_links = sum(1 for tgt in (page_id_building(s) for s in ['1083784365','170900667','1603862331']) if tgt in rlinks)
    print(f"[ROOM PAGE] exists=True  cabinets_shown(1A1/1A3/UPS in text)={'1A1_U11柜' in rtxt},{'1A3_U11柜' in rtxt},{'UPS柜' in rtxt}  cabinet_links={cab_links}/3")
else:
    print("[ROOM PAGE] MISSING")

# 5) boundary / overflow check across every active page (cells must fit 1920x1080)
cur.execute("""SELECT page_name, components, layer FROM display_model_layer
               WHERE model_id=%s AND deleted_at IS NULL""", (MODEL_ID,))
bad = []
layer_bad = []
for pname, comp_b64, layer in cur.fetchall():
    try:
        lj = json.loads(layer)
        if int(lj.get('width', 0)) != 1920 or int(lj.get('height', 0)) != 1080 or int(lj.get('autoSize', -1)) != 1:
            layer_bad.append((pname, lj.get('width'), lj.get('height'), lj.get('autoSize')))
        cells = json.loads(base64.b64decode(comp_b64).decode())['cells']
    except Exception as e:
        bad.append((pname, f'decode-error:{e}')); continue
    for c in cells:
        r = c.get('x', 0) + c.get('width', 0)
        b = c.get('y', 0) + c.get('height', 0)
        if r > 1920 or b > 1080 or c.get('x', 0) < 0 or c.get('y', 0) < 0:
            bad.append((pname, c.get('shape'), c.get('x'), c.get('y'), c.get('width'), c.get('height')))
print(f"[BOUNDS] pages checked, out-of-canvas cells={len(bad)}  layer-misconfig={len(layer_bad)}")
for b in bad[:10]:
    print("   OVERFLOW", b)
for l in layer_bad[:10]:
    print("   LAYER", l)

conn.close()

/* 现场测点按需取数。进画面只订开关和连线，悬停再订这一回路。 */
(function (w) {
  'use strict';

  var byName = {};
  var byCircuit = {};
  var book = {};
  var pageCache = {};
  var treeCache = null;
  var values = {};
  var active = {};
  var ws = null;

  function storeGet(key) {
    try {
      var s = sessionStorage.getItem(key);
      if (s) return s;
    } catch (e) {}
    var m = document.cookie.match(new RegExp('(?:^|; )' + key + '=([^;]*)'));
    if (!m) return '';
    try { return decodeURIComponent(m[1]); } catch (e2) { return m[1]; }
  }

  function headers() {
    var h = { 'Content-Type': 'application/json' };
    var token = storeGet('Authorization');
    var project = storeGet('ProjectUuid');
    if (token) h.Authorization = token;
    if (project) h.ProjectUuid = project;
    return h;
  }

  function post(path, body) {
    return fetch('/api' + path, {
      method: 'POST',
      headers: headers(),
      body: JSON.stringify(body || {})
    }).then(function (r) { return r.json(); });
  }

  function isCap(name) { return String(name || '').indexOf('电容补偿') >= 0; }

  function circuitKey(name) {
    if (!name) return '';
    if (name.slice(-4) === '开关状态') name = name.slice(0, -4);
    else {
      var i = name.lastIndexOf('_');
      if (i > 0 && /电压|电流|功率|电能|频率|谐波|状态|电度|因数|畸变|THD|通讯/.test(name.slice(i + 1))) {
        name = name.slice(0, i);
      }
    }
    if (name.slice(-1) === '_') name = name.slice(0, -1);
    return name;
  }

  function addCircuit(key, pt) {
    if (!key) return;
    (byCircuit[key] = byCircuit[key] || []).push(pt);
  }

  function tie(key, pt) {
    addCircuit(key, pt);
    var m = key.match(/^(配电室\d[AB]\d_.+)_([1-9])$/);
    if (m) addCircuit(m[1], pt);
  }

  function slot(folder, room, cab, dev) {
    var id = folder + '|' + room + '|' + cab;
    if (!book[id]) {
      book[id] = {
        id: id, folder: folder, room: room, cab: cab, dev: dev || '',
        main: null, mainAlarm: null, branches: {}
      };
    }
    if (dev && !book[id].dev) book[id].dev = dev;
    return book[id];
  }

  function branch(sl, num) {
    return sl.branches[num] || (sl.branches[num] = {
      mode: '', dev: '', key: '', p: '', alarm: null, localStatus: null,
      phaseKey: {}, phaseP: {}, phaseDev: {}
    });
  }

  function roomOf(dev) {
    var m = String(dev || '').match(/机房模块(\d[AB]\d)/);
    if (!m) return null;
    var code = m[1];
    var b = code.indexOf('B') >= 0;
    return { code: code, b: b, folder: b ? '数据机房' : '机房模块', room: (b ? '数据机房' : '机房模块') + code };
  }

  function branchKey(info, cab, num, phase) {
    var tail = '_支路' + num + (phase || '');
    if (info && info.b) return '数据机房' + info.code + '_' + cab + tail;
    return cab + tail;
  }

  /* 列头柜按测点归组。带 B 的开关状态记到数据机房告警设备上。返回 stop 时不再走通用回路名。 */
  function noteCabinet(pt) {
    var name = pt.name || '';
    var dev = pt.device || '';
    var am = dev.match(/^数据机房(\d[AB]\d)_([A-Z]列[头尾])$/);
    if (am) {
      var sl = slot('数据机房', '数据机房' + am[1], am[2], dev);
      if (name === '主开关断路器状态' || name === '主开关状态') {
        sl.mainAlarm = pt;
        addCircuit(dev + '_主路', pt);
        return 'stop';
      }
      var bm = name.match(/^支路(\d+)$/);
      if (bm) {
        var b0 = branch(sl, bm[1]);
        b0.alarm = pt;
        if (!b0.mode) b0.mode = '1';
        b0.key = dev + '_支路' + bm[1];
        addCircuit(b0.key, pt);
        ['A', 'B', 'C'].forEach(function (ph) { addCircuit(dev + '_支路' + bm[1] + ph, pt); });
        return 'stop';
      }
      return;
    }

    var longName = name.match(/^机房模块(\d[AB]\d)_([A-Z]列[头尾])_(\d+)([ABC])?相?(?:支路)?/);
    if (longName) {
      var code = longName[1];
      var bRoom = code.indexOf('B') >= 0;
      var sl2 = slot(bRoom ? '数据机房' : '机房模块', (bRoom ? '数据机房' : '机房模块') + code, longName[2], dev);
      var b2 = branch(sl2, longName[3]);
      var ph = longName[4] || '';
      var info = { b: bRoom, code: code };
      if (ph) {
        b2.mode = '3';
        b2.phaseKey[ph] = branchKey(info, longName[2], longName[3], ph);
        b2.phaseDev[ph] = dev;
        addCircuit(b2.phaseKey[ph], pt);
      } else {
        if (b2.mode !== '3') b2.mode = '1';
        b2.key = branchKey(info, longName[2], longName[3], '');
        b2.dev = dev;
        addCircuit(b2.key, pt);
      }
      return 'stop';
    }

    var info2 = roomOf(dev);
    if (!info2 || dev.indexOf('机房模块') !== 0) return;
    var cabM = name.match(/^([A-Z]列[头尾])_(主路开关状态|主开关断路器状态|主开关状态)$/);
    if (cabM) {
      var slM = slot(info2.folder, info2.room, cabM[1], dev);
      if (info2.b) slM.mainAlarm = slM.mainAlarm || pt;
      else slM.main = pt;
      return;
    }
    var brM = name.match(/^([A-Z]列[头尾])_支路(?:电流越上上限告警|电流越上限告警|有功功率|有功电能|有功电度|电流)?(\d+)([ABC])?(X)?(?:开关状态)?$/);
    if (!brM) return;
    var num = brM[2] + (brM[4] || '');
    var phase = brM[4] ? '' : (brM[3] || '');
    var sl3 = slot(info2.folder, info2.room, brM[1], dev);
    var b3 = branch(sl3, num);
    b3.dev = dev;
    if (phase) {
      b3.mode = '3';
      b3.phaseKey[phase] = branchKey(info2, brM[1], brM[2], phase);
      b3.phaseDev[phase] = dev;
      if (/开关状态$/.test(name)) b3.phaseP[phase] = name;
      if (info2.b) addCircuit(b3.phaseKey[phase], pt);
    } else {
      if (b3.mode !== '3') b3.mode = '1';
      b3.key = branchKey(info2, brM[1], num, '');
      if (/开关状态$/.test(name)) {
        b3.p = name;
        b3.localStatus = pt;
        if (info2.b && !b3.alarm) b3.alarm = pt;
      }
      if (info2.b) addCircuit(b3.key, pt);
    }
  }

  function indexPoint(pt) {
    var name = pt.name || '';
    if (!name || isCap(name) || isCap(pt.device)) return;
    if (noteCabinet(pt) === 'stop') {
      (byName[name] = byName[name] || []).push(pt);
      return;
    }
    (byName[name] = byName[name] || []).push(pt);
    var m = name.match(/^(.*列[头尾])_主路/);
    if (m) { addCircuit(m[1] + '_主路', pt); return; }
    m = name.match(/^(.*列[头尾])_主回路/);
    if (m) { addCircuit(m[1] + '_主路', pt); return; }
    m = name.match(/^(.*列[头尾])_主开关/);
    if (m) { addCircuit(m[1] + '_主路', pt); return; }
    m = name.match(/^(.*列[头尾])_支路(?:电流越上上限告警|电流越上限告警|有功功率|有功电能|有功电度|电流)?(\d+[ABC]?)(?:开关状态)?$/);
    if (m) { addCircuit(m[1] + '_支路' + m[2], pt); return; }
    tie(circuitKey(name), pt);
  }

  function loadField() {
    return fetch('/hx-dcim/data/hx-field-index.json').then(function (r) {
      if (!r.ok) throw new Error('no index');
      return r.json();
    }).then(function (j) {
      byName = {};
      byCircuit = {};
      book = {};
      pageCache = {};
      treeCache = null;
      Object.keys(j.byDevice || {}).forEach(function (k) {
        (j.byDevice[k] || []).forEach(indexPoint);
      });
    }).catch(function () {
      byName = {};
      byCircuit = {};
      book = {};
      treeCache = null;
    });
  }

  function pointOf(name, dev) {
    var list = byName[name] || [];
    if (dev) {
      for (var i = 0; i < list.length; i++) if (list[i].device === dev) return list[i];
    }
    return list[0] || null;
  }

  function remembered(pt) {
    if (!pt || !pt.uuid) return null;
    return values[pt.uuid] || null;
  }

  function nodeValue(n) {
    if (!n || !n.p) return null;
    var pt = pointOf(n.p, n.dev);
    var known = remembered(pt);
    if (known) return known.value;
    return pt ? pt.value : null;
  }

  function showValue(p) {
    if (!p || p.value === '' || p.value == null) return '—';
    return String(p.value) + (p.unit ? ' ' + p.unit : '');
  }

  function normVal(value) {
    var s = String(value == null ? '' : value).trim();
    if (s === '合' || s === '合闸' || s === '闭合' || s === 'true') return '1';
    if (s === '分' || s === '分闸' || s === 'false') return '0';
    return s;
  }

  function matchOp(value, op, expect) {
    var s = normVal(value);
    var n = parseFloat(s);
    var e = parseFloat(expect);
    var both = !isNaN(n) && !isNaN(e) && String(expect).trim() !== '';
    if (op === '>') return both && n > e;
    if (op === '>=') return both && n >= e;
    if (op === '<') return both && n < e;
    if (op === '<=') return both && n <= e;
    if (both) return n === e;
    return s === String(expect);
  }

  function remember(rows) {
    (rows || []).forEach(function (p) {
      var id = p.uuid || p.Uuid;
      if (!id) return;
      values[id] = {
        name: p.name || p.DataName || p.Name || '',
        value: p.value != null ? p.value : p.Value,
        unit: p.unit || p.DataUnit || p.data_unit || ''
      };
    });
    if (w.HXLive && w.HXLive.onChange) w.HXLive.onChange();
  }

  function allow(rows) {
    active = {};
    (rows || []).forEach(function (p) {
      var id = p.uuid || p.Uuid;
      if (id) active[id] = 1;
    });
  }

  function fetchUuids(rows) {
    var seen = {};
    var uniq = [];
    (rows || []).forEach(function (p) {
      var id = p.uuid || p.Uuid;
      if (id && !seen[id]) { seen[id] = 1; uniq.push(id); }
    });
    if (!uniq.length) return Promise.resolve([]);
    var jobs = [];
    for (var i = 0; i < uniq.length; i += 60) {
      jobs.push(post('/getRealDataByUuid', { uuid: uniq.slice(i, i + 60) }).catch(function () { return null; }));
    }
    return Promise.all(jobs).then(function (all) {
      var got = [];
      all.forEach(function (j) {
        if (j && j.realData) got = got.concat(j.realData);
      });
      remember(got);
      return got;
    });
  }

  function drivers(pg) {
    var seen = {};
    var rows = [];
    (pg.nodes || []).forEach(function (n) {
      if (n.t !== 'I' && n.t !== 'S') return;
      if (!n.p) return;
      var pt = pointOf(n.p, n.dev);
      if (!pt || !pt.uuid || seen[pt.uuid]) return;
      seen[pt.uuid] = 1;
      rows.push(pt);
    });
    return rows;
  }

  function pointsFor(key, dev) {
    var seen = {};
    var rows = [];
    var open = String(key || '').indexOf('数据机房') === 0;
    (byCircuit[key] || []).forEach(function (p) {
      if (!p.uuid || seen[p.uuid] || isCap(p.name)) return;
      if (!open && dev && p.device && p.device !== dev) return;
      seen[p.uuid] = 1;
      rows.push(p);
    });
    return rows;
  }

  function hasCircuit(key) {
    return !!(key && byCircuit[key] && byCircuit[key].length);
  }

  function cabRank(name) {
    var m = String(name || '').match(/^([A-Z])列(头|尾)$/);
    if (!m) return 9999;
    return (m[2] === '尾' ? 100 : 0) + m[1].charCodeAt(0);
  }

  function numRank(s) {
    var m = String(s).match(/^(\d+)(.*)$/);
    if (!m) return 100000;
    return parseInt(m[1], 10) * 10 + (m[2] ? 1 : 0);
  }

  function packRooms(map) {
    return Object.keys(map).sort().map(function (room) {
      var cabs = map[room].slice().sort(function (a, b) { return cabRank(a.cab) - cabRank(b.cab); });
      return {
        name: room,
        kind: 'column',
        cabs: cabs.map(function (sl) {
          return { name: sl.cab, pageId: 'live:' + sl.id, live: sl.id };
        })
      };
    });
  }

  function tree() {
    if (treeCache) return treeCache;
    var mod = {};
    var data = {};
    Object.keys(book).forEach(function (id) {
      var sl = book[id];
      if (!Object.keys(sl.branches).length) return;
      var box = sl.folder === '数据机房' ? data : mod;
      (box[sl.room] = box[sl.room] || []).push(sl);
    });
    treeCache = { modules: packRooms(mod), dataRooms: packRooms(data) };
    return treeCache;
  }

  /* 机房模块2A1 点表没有无字母支路。原图上的 307 在数据机房2B1 同名列头，单独补到这一柜。 */
  function fillPlainBranches(sl) {
    if (!sl || sl.room !== '机房模块2A1') return;
    var donor = book['数据机房|数据机房2B1|' + sl.cab];
    if (!donor) return;
    Object.keys(donor.branches).forEach(function (num) {
      if (sl.branches[num] || !/^\d+$/.test(String(num))) return;
      var src = donor.branches[num];
      if (!src || src.mode === '3') return;
      var b = branch(sl, num);
      b.mode = '1';
      b.key = src.key || '';
      b.p = src.p || '';
      b.dev = src.dev || donor.dev || '';
      b.alarm = src.alarm || null;
      b.localStatus = src.localStatus || null;
    });
  }

  function cabinetSpec(liveId) {
    var sl = book[liveId];
    if (!sl) return null;
    fillPlainBranches(sl);
    var code = (sl.room.match(/\d[AB]\d/) || [''])[0];
    var mainPt = sl.mainAlarm || sl.main;
    var branches = Object.keys(sl.branches).filter(function (num) {
      var b = sl.branches[num];
      var blob = String(num) + (b.key || '') + (b.dev || '');
      return blob.indexOf('电容补偿') < 0;
    }).sort(function (a, b) { return numRank(a) - numRank(b); }).map(function (num) {
      var b = sl.branches[num];
      var phases = b.mode === '3' ? ['A', 'B', 'C'].filter(function (ph) { return b.phaseKey[ph]; }) : [];
      /* 末尾没有 A/B/C 的是单相。数据机房2B1 的 307 只有一个支路点，单独占一格。 */
      if (phases.length < 2) phases = [];
      var status = b.alarm || b.localStatus;
      return {
        num: num,
        phases: phases,
        key: phases.length ? '' : (b.key || (sl.cab + '_支路' + num)),
        p: status ? status.name : (b.p || ''),
        dev: status ? status.device : (b.dev || sl.dev),
        phaseKey: b.phaseKey,
        phaseP: b.phaseP,
        phaseDev: b.phaseDev
      };
    });
    return {
      id: 'live:' + sl.id,
      name: sl.room + '_' + sl.cab,
      room: sl.room,
      cab: sl.cab,
      code: code,
      incoming: (sl.cab || code || '') + '进线',
      mainKey: mainPt && String(mainPt.device || '').indexOf('数据机房') === 0 ? (mainPt.device + '_主路') : (sl.cab + '_主路'),
      mainPoint: mainPt ? mainPt.name : '',
      mainDev: mainPt ? mainPt.device : sl.dev,
      branches: branches
    };
  }

  function cabinetPage(id) {
    if (!id || id.indexOf('live:') !== 0) return null;
    if (pageCache[id]) return pageCache[id];
    var spec = cabinetSpec(id.slice(5));
    if (!spec || !w.HXDraw || !w.HXDraw.cabinet) return null;
    pageCache[id] = w.HXDraw.cabinet(spec);
    return pageCache[id];
  }

  function connect() {
    post('/GetSystemParams', {}).then(function (j) {
      var list = j && (j.list || (j.data && j.data.list));
      if (!list && j && j.WSPort) list = j;
      if (!list) return;
      var port = list.WSPort;
      var addr = list.WSAddress;
      if (!addr || addr === 'local') addr = location.hostname;
      if (location.protocol === 'https:') port = 443;
      var token = storeGet('Authorization');
      var project = storeGet('ProjectUuid');
      if (!token || !project || !port) return;
      var url = (location.protocol === 'https:' ? 'wss' : 'ws') + '://' + addr + ':' + port + '/ws?token=' + encodeURIComponent(token) + '&project=' + encodeURIComponent(project);
      if (ws) { try { ws.close(); } catch (e) {} }
      ws = new WebSocket(url);
      ws.onmessage = function (ev) {
        var msg;
        try { msg = JSON.parse(ev.data); } catch (e) { return; }
        if (!msg || msg.Cmd !== 'RealData') return;
        var hit = [];
        (msg.Data || []).forEach(function (p) {
          var id = p.Uuid || p.uuid;
          if (id && active[id]) hit.push(p);
        });
        if (hit.length) remember(hit);
      };
    }).catch(function () {});
  }

  w.HXLive = {
    ready: loadField(),
    headers: headers,
    pointOf: pointOf,
    remembered: remembered,
    nodeValue: nodeValue,
    showValue: showValue,
    matchOp: matchOp,
    allow: allow,
    fetch: fetchUuids,
    drivers: drivers,
    pointsFor: pointsFor,
    hasCircuit: hasCircuit,
    tree: tree,
    cabinetPage: cabinetPage,
    uuidValue: function (id) { return values[id] || null; },
    clear: function () { active = {}; },
    connect: connect,
    onChange: null
  };
})(window);

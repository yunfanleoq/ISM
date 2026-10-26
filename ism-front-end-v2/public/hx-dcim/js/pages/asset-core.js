/* =========================================================
   资产管理 · 口径模型（依据《资产管理 · 口径定稿 v1.0》）
   ---------------------------------------------------------
   第 07 条 空间利用率 = (Σ设备U + 设备台数 × 空U系数) ÷ 42，上限 100%
   第 08 条 综合使用率 = max(空间, 电力, 接口)，取极高不平均
   第 12 条 四级色阶：灰未启用 / 蓝<70 / 黄70-90 / 红>90；电力红线单设 85%
   第 04 条 永久资产编码 + IP 一对多子表 + 主 SN 查重
   第 01 条 使用单位由设备细表汇总，机柜按「主用单位」落主分类
   第 02 条「涉及机房数」只统计已部署（有在架设备）的机房
   ========================================================= */
(function (A) {
  'use strict';

  var AM = {};
  A.AM = AM;

  /* 从资产管理进入某一层 / 某一间的立体画面；hash 末段 assetmgr 让返回回到资产管理。
     若这一次是从总览进的资产管理，来源 ~dashboard 要跟着走，回到资产页还能再回总览。 */
  AM.openSpace = function (kind, name) {
    if (!kind || !name) return;
    A.closeDrawer();
    var tail = '/assetmgr';
    if (A.origin === 'dashboard') tail += '/~dashboard';
    A.go('building/' + kind + '/' + encodeURIComponent(name) + tail);
  };

  /* ---------- 空U系数（按设备类型取，不从机柜取） ---------- */
  AM.COEF_DEF = [
    { type: '服务器',       coef: 1,   note: '1U / 2U 机架式，标准占位' },
    { type: '存储设备',     coef: 1,   note: '存储节点 / 磁盘阵列' },
    { type: '交换机',       coef: 0.5, note: '网络设备，共享理线空间' },
    { type: '网络安全设备', coef: 0.5, note: '防火墙 / 网闸 / 审计' },
    { type: '其他设备',     coef: 0.5, note: 'KVM / 终端 / 采集器' },
    { type: '光纤配线柜',   coef: 0,   note: '不占设备安装空间' },
    { type: '电池柜',       coef: 0,   note: '不占设备安装空间' }
  ];

  var COEF = {};
  AM.COEF_DEF.forEach(function (c) { COEF[c.type] = c.coef; });

  AM.coef = function (type) {
    var v = COEF[type];
    return v == null ? 1 : v;
  };
  /* 系数未改动时，柜内混合系数 = 柜级基准（与容量总览同源）；
     改动后按「类型占比 × 系数偏差」叠加增量，默认值与基准恒等，
     因此恢复默认会自动回到同一个数，不需要额外状态位。 */
  AM.setCoef = function (type, v) {
    var n = Number(v);
    COEF[type] = isFinite(n) ? Math.max(0, Math.min(4, n)) : AM.coef(type);
    AM.invalidate();
  };
  AM.coefDelta = function (type) {
    var def = 0;
    AM.COEF_DEF.forEach(function (c) { if (c.type === type) def = c.coef; });
    return AM.coef(type) - def;
  };
  AM.isDefault = function () {
    return AM.COEF_DEF.every(function (c) { return AM.coef(c.type) === c.coef; });
  };
  AM.resetCoef = function () {
    AM.COEF_DEF.forEach(function (c) { COEF[c.type] = c.coef; });
    AM.invalidate();
  };

  /* ---------- 阈值（第 12 条） ---------- */
  AM.TH = {
    space: { warn: 70, red: 90, label: '空间' },
    power: { warn: 70, red: 85, label: '电力' },
    iface: { warn: 70, red: 90, label: '接口' }
  };
  AM.OVR = {};                               /* 机房级覆盖 {room:{power:{warn,red}}} */
  AM.thOf = function (dim, room) {
    var base = AM.TH[dim];
    var o = room && AM.OVR[room] && AM.OVR[room][dim];
    return o ? { warn: o.warn, red: o.red, label: base.label } : base;
  };
  AM.setTh = function (dim, warn, red, room) {
    var w = Number(warn), r = Number(red);
    if (!isFinite(w) || !isFinite(r)) return;
    if (room) {
      AM.OVR[room] = AM.OVR[room] || {};
      AM.OVR[room][dim] = { warn: w, red: r };
    } else {
      AM.TH[dim].warn = w;
      AM.TH[dim].red = r;
    }
    AM.invalidate();
  };

  /* ---------- 色阶 ---------- */
  AM.LV_TXT = { off: '未启用', blue: '正常', yellow: '关注', red: '满载' };
  AM.lv = function (rate, dev) {
    if (!dev) return 'off';
    if (rate > AM.TH.space.red) return 'red';
    if (rate >= AM.TH.space.warn) return 'yellow';
    return 'blue';
  };
  AM.lvDim = function (rate, dim, room) {
    var t = AM.thOf(dim, room);
    if (rate > t.red) return 'red';
    if (rate >= t.warn) return 'yellow';
    return 'blue';
  };

  /* ---------- 卡在哪一项（第 08 条：回答「为什么不能装了」） ---------- */
  AM.bottleneck = function (c) {
    var s = c.spaceRate / AM.thOf('space', c.room).red;
    var p = c.powRate / AM.thOf('power', c.room).red;
    var i = c.ifRate / AM.thOf('iface', c.room).red;
    var mx = Math.max(s, p, i);
    var dim = mx === s ? 'space' : (mx === p ? 'power' : 'iface');
    var val = dim === 'space' ? c.spaceRate : (dim === 'power' ? c.powRate : c.ifRate);
    var th = AM.thOf(dim, c.room);
    return {
      dim: dim,
      label: th.label,
      value: val,
      red: th.red,
      over: val > th.red,
      ratio: mx
    };
  };

  /* ---------- 工具 ---------- */
  function pad(n, w) {
    var s = String(n);
    while (s.length < w) s = '0' + s;
    return s;
  }
  AM.pad = pad;

  function fnv(s) {
    var h = 2166136261;
    for (var i = 0; i < s.length; i++) { h ^= s.charCodeAt(i); h = (h * 16777619) >>> 0; }
    return h;
  }
  AM.fnv = fnv;

  function r1(n) { return Math.round(Number(n || 0) * 10) / 10; }
  function r2(n) { return Math.round(Number(n || 0) * 100) / 100; }

  /* 可部署台数：空间 / 电力 / 接口三面取小（第 08 条） */
  function deployable(c, size) {
    var freeU = 42 - c.usedU - (c.dev || 0) * (c.coef || 0);
    var byU = Math.floor(freeU / (size + 1));
    var byKw = Math.floor(Math.max(0, (c.rated || 4) - (c.rtKw || 0)) / 0.6);
    var byIf = Math.max(0, 20 - (c.ifUsed || 0));
    return Math.max(0, Math.min(byU, byKw, byIf, 20));
  }
  AM.deployable = deployable;

  var NET_SEG = { '公安网': 0, '网安专网': 1, '互联网': 2, '视频专网': 3, '涉密网': 4 };

  AM.M = null;
  AM.invalidate = function () { AM.M = null; };

  /* =========================================================
     建模
     ========================================================= */
  AM.build = function () {
    if (AM.M) return AM.M;
    var D = window.DCIM;

    /* --- 机柜骨架 --- */
    var cabIndex = {}, cabByRoom = {}, roomOrder = [], roomInfo = {};
    D.rooms.forEach(function (r) {
      var arr = [];
      (r.cabinets || []).forEach(function (c) {
        var short = String(c.name).replace('屏蔽机房', '').replace('机房', '');
        var o = {
          name: c.name, short: short, room: r.name, floor: r.floor,
          col: c.col, seq: c.seq,
          dev: c.dev || 0, usedU: c.usedU || 0, rated: c.rated || 4, rtKw: c.rtKw || 0,
          coefRaw: Number(c.coef || 1), coef: Number(c.coef || 1),
          powRate: r1(c.powRate), ifUsed: c.ifUsed || 0, ifRate: r1(c.ifRate),
          dlist: [], coefSum: 0, mix: {}, orgCount: {},
          spaceRate: 0, comp: 0, lv: 'off', freeU: 42, can2: 0, can4: 0,
          units: [], mainUnit: '—',
          placed: c.x != null || !!c.id
        };
        cabIndex[c.name] = o;
        arr.push(o);
      });
      if (arr.length) {
        cabByRoom[r.name] = arr;
        roomOrder.push(r.name);
        roomInfo[r.name] = { name: r.name, floor: r.floor, type: r.type, cabs: arr.length };
      }
    });

    /* --- 设备 --- */
    var list = [], byCode = {};
    D.assets.forEach(function (a, i) {
      var cabName = a.cabinet || '';
      var cab = cabIndex[cabName];
      var located = !!cab;

      var st = !located ? '未定位'
        : a.status === '下架待处理' ? '已下架'
        : a.status === '使用中' ? '在用'
        : a.status === '备用' ? '备用' : '维修中';

      var seg = String(cabName).split('-');
      /* 第 04 条：编码永久且唯一。未定位设备没有柜位，但机房仍然知道，
         用房号兜底，避免出现 BJDC-000-0-0-xxxx 这种无意义编码 */
      var roomDigits = String(a.room || '').replace(/[^0-9]/g, '');
      var roomNo = String(seg[0] || '').replace(/[^0-9]/g, '') || roomDigits || '000';
      var floorNum = parseInt(a.floor, 10) || 1;
      var h1 = (i >> 8) & 255, h2 = i & 255;
      var seg2 = NET_SEG[a.net];
      if (seg2 == null) seg2 = 5;

      var ips = [{ t: '管理', ip: '10.' + floorNum + '.' + h1 + '.' + h2, main: true, note: '带外管理口' }];
      if (a.cat === '服务器' || a.cat === '交换机' || a.cat === '网络安全设备' || a.cat === '其他设备' || a.cat === '其他') {
        ips.push({ t: '业务', ip: '172.' + (16 + seg2) + '.' + h1 + '.' + h2, main: false, note: a.net });
      }
      if (a.cat === '存储设备' || a.cat === '服务器') {
        ips.push({ t: '存储', ip: '192.168.' + h1 + '.' + h2, main: false, note: '存储平面' });
      }

      var h = fnv(String(a.sn || '') + i);
      var snState = (i % 47 === 13) ? 'missing' : (i % 9 === 4 ? 'dual' : 'ok');
      var cat = a.cat === '其他' ? '其他设备' : a.cat;

      var d = {
        i: i,
        code: 'BJDC-' + pad(roomNo, 3) + '-' + (cab ? cab.col : (Number(seg[1]) || 0)) + '-' +
          (cab ? cab.seq : (Number(seg[2]) || 0)) + '-' + pad(i + 1, 4),
        room: a.room, floor: a.floor,
        cab: cabName, cabShort: String(cabName).replace('屏蔽机房', '').replace('机房', ''),
        col: Number(seg[1]) || 0, seq: Number(seg[2]) || 0,
        uStart: a.uStart || 1, uEnd: a.uEnd || a.uStart || 1, uH: a.uH || 1,
        sn: a.sn || '', snAlt: snState === 'dual' ? ('SN2-' + String(a.sn || '').slice(-8) + pad(i % 90 + 10, 2)) : '',
        snState: snState,
        cat: cat, type: a.type || cat,
        brand: a.brand || '—', model: a.model || '—',
        kw: a.kw || 0, psu: a.psu || 1, weight: a.weight || 0,
        net: a.net || '—', biz: a.biz || '—', org: a.org || '—', owner: a.owner || '—',
        st: st, located: located, ips: ips, mainIp: ips[0].ip, coef: AM.coef(cat),
        hist: []
      };
      list.push(d);
      byCode[d.code] = d;

      if (cab) {
        cab.dlist.push(d);
        cab.coefSum += d.coef;
        cab.mix[d.cat] = (cab.mix[d.cat] || 0) + 1;
        cab.orgCount[d.org] = (cab.orgCount[d.org] || 0) + 1;
      }
    });

    /* --- 位置轨迹（第 06 条：不审批，强制留痕） --- */
    list.forEach(function (d) {
      var h = fnv(d.code);
      d.hist.push({
        t: '2024-' + pad(1 + h % 12, 2) + '-' + pad(1 + (h >> 3) % 27, 2) + ' 09:' + pad(h % 60, 2),
        act: '首次上架', from: '设备入库', to: d.cabShort + ' · U' + d.uStart,
        by: '系统导入', why: '初始建档'
      });
      if (d.st === '未定位' && d.cab) {
        d.hist.push({
          t: '2026-09-' + pad(2 + (h >> 5) % 26, 2) + ' 1' + (h % 8) + ':' + pad((h >> 2) % 60, 2),
          act: '下架回库', from: d.cabShort + ' · U' + d.uStart, to: '未定位（保留原位置）',
          by: '张伟', why: '整机置换'
        });
      } else if (h % 4 === 0) {
        var arr = cabByRoom[d.room] || [];
        var oc = arr.length ? arr[(h >> 7) % arr.length] : null;
        d.hist.push({
          t: '2025-' + pad(1 + (h >> 4) % 12, 2) + '-' + pad(1 + (h >> 6) % 27, 2) + ' 1' + (h % 9) + ':' + pad((h >> 8) % 60, 2),
          act: '位置变更',
          from: oc ? (oc.short + ' · U' + (1 + h % 40)) : '—',
          to: d.cabShort + ' · U' + d.uStart,
          by: '李强', why: h % 3 === 0 ? '机柜扩容' : '同机房优化'
        });
      }
      d.hist.sort(function (x, y) { return x.t < y.t ? 1 : -1; });
    });

    /* --- 机柜口径重算（公式自洽，与容量页同源） --- */
    var cabs = [], gEffU = 0, gUsedU = 0;
    roomOrder.forEach(function (rn) {
      cabByRoom[rn].forEach(function (c) {
        /* 第 07 条：柜内有效占用 = Σ设备U + 设备台数 × 空U系数
           系数取设备类型；柜级基准负责与存量口径对齐，类型偏差负责灵敏响应 */
        var delta = 0;
        if (c.dev) {
          Object.keys(c.mix).forEach(function (t) {
            delta += c.mix[t] / c.dev * AM.coefDelta(t);
          });
        }
        c.coefMix = r2(delta);
        c.coef = r2(c.coefRaw + delta);
        var eff = c.usedU + c.dev * c.coef;
        c.spaceRate = c.dev ? Math.min(100, r1(eff / 42 * 100)) : 0;
        c.units = Object.keys(c.orgCount).sort(function (x, y) { return c.orgCount[y] - c.orgCount[x]; });
        c.mainUnit = c.units[0] || '—';
        c.comp = Math.min(100, r1(Math.max(c.spaceRate, c.powRate, c.ifRate)));
        c.lv = AM.lv(c.comp, c.dev);
        c.freeU = Math.max(0, 42 - c.usedU);
        c.can2 = deployable(c, 2);
        c.can4 = deployable(c, 4);
        gEffU += eff;
        gUsedU += c.usedU;
        cabs.push(c);
      });
    });

    /* --- 全局汇总 ---
       第 12 条：未启用机柜（灰）不计入可分配容量，与容量总览同源
       总容量 = 已启用机柜数 × 42U；未启用柜位单列，不参与空间率分母 */
    var usedCabN = cabs.filter(function (c) { return c.dev > 0 || c.placed; }).length;
    var capU = usedCabN * 42;
    var capAllU = cabs.length * 42;
    var cnt = { '在用': 0, '备用': 0, '维修中': 0, '已下架': 0, '未定位': 0 };
    var snMissing = 0, snDual = 0, locatedN = 0;
    var byOrg = {}, byBiz = {}, byCat = {}, byRoom = {};

    list.forEach(function (d) {
      cnt[d.st] = (cnt[d.st] || 0) + 1;
      if (d.snState === 'missing') snMissing++;
      if (d.snState === 'dual') snDual++;
      if (d.located) locatedN++;

      if (!byOrg[d.org]) byOrg[d.org] = { name: d.org, dev: 0, cabs: 0, touch: 0, rooms: {}, mainCabs: 0 };
      byOrg[d.org].dev++;
      if (d.located) byOrg[d.org].rooms[d.room] = 1;

      if (!byBiz[d.biz]) byBiz[d.biz] = { name: d.biz, dev: 0, rooms: {}, cabs: {}, locatedDev: 0 };
      byBiz[d.biz].dev++;
      if (d.located) {
        byBiz[d.biz].rooms[d.room] = 1;
        byBiz[d.biz].cabs[d.cab] = 1;
        byBiz[d.biz].locatedDev++;
      }

      if (!byCat[d.cat]) byCat[d.cat] = { name: d.cat, dev: 0, u: 0, coef: 0 };
      byCat[d.cat].dev++;
      byCat[d.cat].u += d.uH;
      byCat[d.cat].coef += d.coef;

      if (!byRoom[d.room]) byRoom[d.room] = { name: d.room, floor: d.floor, dev: 0, locatedDev: 0, cabs: 0, compSum: 0 };
      byRoom[d.room].dev++;
      if (d.located) byRoom[d.room].locatedDev++;
    });

    /* 机柜归属：主用单位落主分类；涉及柜按柜内单位去重 */
    cabs.forEach(function (c) {
      if (c.mainUnit !== '—' && byOrg[c.mainUnit]) byOrg[c.mainUnit].mainCabs++;
      c.units.forEach(function (u) {
        if (byOrg[u]) byOrg[u].touch++;
      });
    });
    roomOrder.forEach(function (rn) {
      var arr = cabByRoom[rn];
      for (var i = 0; i < arr.length; i++) {
        var c = arr[i];
        if (byRoom[rn]) {
          byRoom[rn].cabs++;
          byRoom[rn].compSum += c.comp;
        }
      }
    });

    var orgs = Object.keys(byOrg).map(function (k) {
      var o = byOrg[k];
      o.roomN = Object.keys(o.rooms).length;
      return o;
    }).sort(function (x, y) { return y.dev - x.dev; });

    var bizs = Object.keys(byBiz).map(function (k) {
      var b = byBiz[k];
      b.roomN = Object.keys(b.rooms).length;
      b.cabN = Object.keys(b.cabs).length;
      return b;
    }).sort(function (x, y) { return y.dev - x.dev; });

    var cats = Object.keys(byCat).map(function (k) { return byCat[k]; })
      .sort(function (x, y) { return y.dev - x.dev; });

    var roomList = roomOrder.map(function (rn) {
      var r = byRoom[rn];
      r.cabN = cabByRoom[rn].length;
      r.comp = r.cabN ? r1(r.compSum / r.cabN) : 0;
      r.usedCabs = cabByRoom[rn].filter(function (c) { return c.dev > 0 || c.placed; }).length;
      r.freeCabs = r.cabN - r.usedCabs;
      return r;
    });

    var deployedRooms = roomList.filter(function (r) { return r.locatedDev > 0; });

    var M = {
      D: D,
      list: list, byCode: byCode,
      cabs: cabs, cabIndex: cabIndex, cabByRoom: cabByRoom,
      roomOrder: roomOrder, roomInfo: roomInfo,
      rooms: roomList, deployedRooms: deployedRooms,
      orgs: orgs, bizs: bizs, cats: cats,
      cnt: cnt,
      snMissing: snMissing, snDual: snDual,
      locatedN: locatedN,
      total: list.length,
      totalCabs: cabs.length,
      usedCabs: usedCabN,
      freeCabs: cabs.length - usedCabN,
      capU: capU,
      capAllU: capAllU,
      usedU: gUsedU,
      effU: gEffU,
      spaceRate: capU ? r2(gEffU / capU * 100) : 0,
      dist: { off: 0, blue: 0, yellow: 0, red: 0 }
    };

    cabs.forEach(function (c) { M.dist[c.lv] = (M.dist[c.lv] || 0) + 1; });

    AM.M = M;
    return M;
  };

  /* ---------- 便捷查询 ---------- */
  AM.cabOf = function (name) {
    return AM.build().cabIndex[name] || null;
  };
  AM.assetOf = function (code) {
    return AM.build().byCode[code] || null;
  };
  AM.orgCabs = function (org) {
    var M = AM.build();
    return M.cabs.filter(function (c) { return c.mainUnit === org; });
  };
  AM.roomCabs = function (room) {
    return AM.build().cabByRoom[room] || [];
  };
  AM.stTag = function (st) {
    var map = { '在用': 'green', '备用': 'blue', '维修中': 'orange', '已下架': 'grey', '未定位': 'red' };
    return map[st] || 'grey';
  };
})(window.App);

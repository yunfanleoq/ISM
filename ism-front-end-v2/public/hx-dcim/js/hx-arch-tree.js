/* 架构层级机构.html 导出的工程树：配电室/机房 → 楼层 → 房间 → 列（TA/WA…） */
(function (w) {
  'use strict';

  var spec = null;
  var loadP = null;

  function esc(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function load() {
    if (spec) return Promise.resolve(spec);
    if (loadP) return loadP;
    loadP = fetch('/hx-dcim/data/hx-arch-tree.json?v=20261025treehint').then(function (r) {
      if (!r.ok) throw new Error('no arch tree');
      return r.json();
    }).then(function (j) {
      spec = j;
      w.HXArch.branches = spec.branches || [];
      return spec;
    }).catch(function () {
      spec = { version: 0, branches: [] };
      w.HXArch.branches = [];
      return spec;
    });
    return loadP;
  }

  /** 设计稿「数据机房2A1」等在测点侧常挂在「机房模块2A1」。 */
  function liveRoomName(designName) {
    var n = String(designName || '');
    var m = n.match(/^数据机房(\d[AB]\d)$/);
    if (!m) return n;
    if (m[1].indexOf('B') >= 0) return n;
    return '机房模块' + m[1];
  }

  /** TA/TB… → A列头；WA/WB… → A列尾 */
  function cabFromDesign(label) {
    var m = String(label || '').match(/^([TW])([A-K])$/);
    if (!m) return '';
    return m[2] + '列' + (m[1] === 'T' ? '头' : '尾');
  }

  function roomCode(name) {
    var m = String(name || '').match(/(\d[AB]\d)/);
    return m ? m[1] : '';
  }

  function isPowerLeaf(name) {
    return /^配电室\d/.test(name) || /^ECC\d/.test(name);
  }

  /** 配电室叶子 ↔ 同编码数据机房 / 机房模块（设计稿与测点命名） */
  function powerDrillTargets(powerLeafName) {
    var code = roomCode(powerLeafName);
    if (!code) return [];
    var dataDesign = '数据机房' + code;
    return [
      { role: 'module', design: '机房模块' + code, live: '机房模块' + code, label: '机房模块' + code },
      { role: 'data', design: dataDesign, live: liveRoomName(dataDesign), label: dataDesign + '（列头柜）' }
    ];
  }

  function getNodeAt(pathNames) {
    var cur = { children: (spec && spec.branches) || [] };
    var i, j, name, kids, hit;
    for (i = 0; i < pathNames.length; i++) {
      name = pathNames[i];
      kids = cur.children || [];
      hit = null;
      for (j = 0; j < kids.length; j++) {
        if (kids[j].name === name) { hit = kids[j]; break; }
      }
      if (!hit) return null;
      cur = hit;
    }
    return cur;
  }

  function findLeafPath(targetName) {
    var found = null;
    function walk(node, anc) {
      if (found) return;
      var name = node.name || '';
      var kids = node.children || [];
      if (name === targetName) {
        if (!kids.length) found = { ancestors: anc.slice(), name: name };
        return;
      }
      anc.push(name);
      kids.forEach(function (ch) { walk(ch, anc); });
      anc.pop();
    }
    ((spec && spec.branches) || []).forEach(function (b) { walk(b, []); });
    return found;
  }

  /** 从配电室下钻：优先有测点的机房模块，否则数据机房首列 TA */
  function suggestFromPower(powerLeafName) {
    var targets = powerDrillTargets(powerLeafName);
    var t, path, resolved, node, col, i;
    for (i = 0; i < targets.length; i++) {
      t = targets[i];
      path = findLeafPath(t.design);
      if (!path) continue;
      if (t.role === 'module') {
        resolved = resolveLeaf(path.name, path.ancestors);
        if (resolved && resolved.ok) {
          return { item: { name: path.name, ancestors: path.ancestors, resolved: resolved }, target: t };
        }
      }
      node = getNodeAt(path.ancestors.concat([path.name]));
      if (node && node.children && node.children.length) {
        col = node.children[0];
        path = { ancestors: path.ancestors.concat([path.name]), name: col.name };
        resolved = resolveLeaf(path.name, path.ancestors);
        if (resolved && resolved.ok) {
          return { item: { name: path.name, ancestors: path.ancestors, resolved: resolved }, target: t };
        }
      }
    }
    return null;
  }

  function crumbLabels(ancestors, leafName, cabName) {
    var parts = (ancestors || []).slice();
    if (leafName) parts.push(leafName);
    if (cabName) parts.push(cabName);
    return parts;
  }

  function isColumnLeaf(name) {
    return /^[TW][A-K]$/.test(name);
  }

  function isModuleLeaf(name) {
    return /^机房模块\d/.test(name);
  }

  function catalogIndex() {
    var idx = { byName: {}, pairs: {} };
    var cat = w.HXScreen && w.HXScreen.catalog ? w.HXScreen.catalog() : [];
    cat.forEach(function (f) {
      (f.rooms || []).forEach(function (room) {
        idx.byName[room.name] = room;
        (room.cabs || []).forEach(function (cab) {
          idx.pairs[room.name + '\0' + cab.name] = { room: room, cab: cab };
        });
      });
    });
    return idx;
  }

  function resolveLeaf(name, ancestors) {
    var idx = catalogIndex();
    var path = (ancestors || []).concat([name]);
    var i, roomName, cabName, room, pair;

    if (isPowerLeaf(name)) {
      room = idx.byName[name];
      if (room) return { room: room, cab: null, ok: true, kind: 'power', design: name };
      return {
        ok: false,
        reason: '一次图未解开',
        design: name,
        kind: 'power',
        links: powerDrillTargets(name)
      };
    }

    if (isColumnLeaf(name)) {
      roomName = '';
      for (i = path.length - 2; i >= 0; i--) {
        if (/^数据机房\d/.test(path[i])) {
          roomName = liveRoomName(path[i]);
          break;
        }
      }
      cabName = cabFromDesign(name);
      if (!roomName || !cabName) return { ok: false, reason: '列标识无法解析', design: name };
      pair = idx.pairs[roomName + '\0' + cabName];
      if (pair) return { room: pair.room, cab: pair.cab, ok: true };
      room = idx.byName[roomName];
      if (room) return { ok: false, reason: '现场索引无' + cabName, design: name, room: room };
      return { ok: false, reason: '机房未绑点', design: roomName };
    }

    if (isModuleLeaf(name)) {
      room = idx.byName[name];
      if (room) return { room: room, cab: room.cabs && room.cabs[0] ? room.cabs[0] : null, ok: !!room };
      return { ok: false, reason: '模块未绑点', design: name };
    }

    if (/^数据机房\d/.test(name)) {
      roomName = liveRoomName(name);
      room = idx.byName[roomName];
      if (room) {
        return {
          room: room,
          cab: room.cabs && room.cabs[0] ? room.cabs[0] : null,
          ok: true
        };
      }
      return { ok: false, reason: '机房未绑点', design: name };
    }

    return { ok: false, reason: '非叶子节点', design: name };
  }

  function branchFilter(mode) {
    if (mode === 'power') return function (b) { return b.name === '配电室'; };
    if (mode === 'room') return function (b) { return b.name === '机房'; };
    if (mode === 'all' || mode === 'drill') return function () { return true; };
    return function () { return true; };
  }

  function renderTree(host, opts) {
    opts = opts || {};
    var mode = opts.mode || '';
    var onPick = opts.onPick || function () {};
    var active = opts.active || null;
    var branches = (spec && spec.branches) || [];
    branches = branches.filter(branchFilter(mode));
    if (!branches.length) {
      host.innerHTML = '<div class="hx-pt">架构树未加载。</div>';
      return { flat: [], openPrefer: null };
    }

    var flat = [];
    var html = '';

    var ui = (spec && spec.ui) || {};
    var hideSoftPower = !!ui.hidePowerWithoutScene;

    function renderNode(node, depth, ancestors) {
      var name = node.name || '';
      var kids = node.children || [];
      var anc = ancestors.slice();
      var resolved = kids.length ? null : resolveLeaf(name, anc);
      var pad = depth > 0 ? ' hx-sub' + Math.min(depth, 4) : '';
      var isOn = active && active.archPath && active.archPath.join('\0') === anc.concat([name]).join('\0');
      if (!isOn && active && active.room && resolved) {
        isOn = resolved.ok && active.room.name === resolved.room.name &&
          (active.cab == null ? !resolved.cab : (resolved.cab && active.cab && active.cab.name === resolved.cab.name));
      }

      if (!kids.length) {
        var soft = resolved && !resolved.ok;
        if (hideSoftPower && soft && resolved && resolved.kind === 'power') return;
        var fi = flat.length;
        flat.push({ name: name, ancestors: anc, resolved: resolved });
        var title = esc(name);
        var cls = 'hx-node' + pad + (soft ? ' is-soft' : '') + (isOn ? ' is-on' : '');
        html += '<button type="button" class="' + cls + '" data-i="' + fi + '">' + title + '</button>';
        return;
      }

      html += '<div class="hx-folder' + pad + '">' + esc(name) + '</div>';
      anc.push(name);
      kids.forEach(function (ch) { renderNode(ch, depth + 1, anc); });
    }

    branches.forEach(function (b) { renderNode(b, 0, []); });

    host.innerHTML = html;
    host.querySelectorAll('.hx-node').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var item = flat[Number(btn.getAttribute('data-i'))];
        if (!item || !item.resolved) return;
        onPick(item.resolved.room, item.resolved.cab, item);
        host.querySelectorAll('.hx-node').forEach(function (b) { b.classList.remove('is-on'); });
        btn.classList.add('is-on');
      });
    });

    return { flat: flat, host: host };
  }

  function activateFlatItem(flatOut, index) {
    var host = flatOut && flatOut.host;
    if (!host || index < 0) return;
    var btn = host.querySelector('.hx-node[data-i="' + index + '"]');
    if (!btn) return;
    host.querySelectorAll('.hx-node').forEach(function (b) { b.classList.remove('is-on'); });
    btn.classList.add('is-on');
    try { btn.scrollIntoView({ block: 'nearest' }); } catch (e) { btn.scrollIntoView(true); }
  }

  function pickByDesign(flatOut, designName, preferCab) {
    var flat = (flatOut && flatOut.flat) || [];
    var i, item, r;
    for (i = 0; i < flat.length; i++) {
      item = flat[i];
      if (item.name !== designName) continue;
      r = item.resolved;
      if (!r) continue;
      if (r.ok) {
        if (preferCab && r.cab && r.cab.name !== preferCab) continue;
        activateFlatItem(flatOut, i);
        return item;
      }
    }
    return null;
  }

  function findPrefer(flat, preferRoom, preferCab) {
    var i, item, r;
    for (i = 0; i < flat.length; i++) {
      item = flat[i];
      r = item.resolved;
      if (!r || !r.ok || !r.room) continue;
      if (r.room.name !== preferRoom) continue;
      if (preferCab) {
        if (r.cab && r.cab.name === preferCab) return i;
      } else if (!r.cab || preferCab === null) {
        return i;
      }
    }
    for (i = 0; i < flat.length; i++) {
      r = flat[i].resolved;
      if (r && r.ok && r.room && r.room.name === preferRoom) return i;
    }
    return -1;
  }

  w.HXArch = {
    branches: [],
    ready: load(),
    load: load,
    liveRoomName: liveRoomName,
    cabFromDesign: cabFromDesign,
    roomCode: roomCode,
    powerDrillTargets: powerDrillTargets,
    suggestFromPower: suggestFromPower,
    findLeafPath: findLeafPath,
    crumbLabels: crumbLabels,
    resolveLeaf: resolveLeaf,
    renderTree: renderTree,
    findPrefer: findPrefer,
    activateFlatItem: activateFlatItem,
    pickByDesign: pickByDesign
  };
})(window);

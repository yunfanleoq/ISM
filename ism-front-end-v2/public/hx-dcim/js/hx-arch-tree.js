/* 架构层级机构.html 导出的工程树：配电室/机房 → 楼层 → 房间 → 列（TA/WA…）；
 * 网络拓扑图来自 设备目录树.html 的 treeData（通讯管理机 / COM / 设备）。 */
(function (w) {
  'use strict';

  var spec = null;
  var netSpec = null;
  var loadP = null;
  var CACHE = '20261029topo04';

  function esc(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function appendNetBranch(branches) {
    var list = (branches || []).slice();
    var i;
    for (i = 0; i < list.length; i++) {
      if (list[i] && list[i].name === '网络拓扑图') return list;
    }
    if (netSpec && netSpec.branch) list.push(netSpec.branch);
    return list;
  }

  function load() {
    if (spec) return Promise.resolve(spec);
    if (loadP) return loadP;
    loadP = Promise.all([
      fetch('/hx-dcim/data/hx-arch-tree.json?v=' + CACHE).then(function (r) {
        if (!r.ok) throw new Error('no arch tree');
        return r.json();
      }),
      fetch('/hx-dcim/data/hx-net-topo-tree.json?v=' + CACHE).then(function (r) {
        if (!r.ok) return null;
        return r.json();
      }).catch(function () { return null; })
    ]).then(function (pair) {
      spec = pair[0] || { version: 0, branches: [] };
      netSpec = pair[1];
      w.HXArch.branches = appendNetBranch(spec.branches || []);
      w.HXArch.netRules = (netSpec && netSpec.rules) || null;
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

  function allBranches() {
    return appendNetBranch((spec && spec.branches) || []);
  }

  function getNodeAt(pathNames) {
    var cur = { children: allBranches() };
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
    allBranches().forEach(function (b) { walk(b, []); });
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

  function isTopoPath(ancestors) {
    return (ancestors || [])[0] === '网络拓扑图';
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

  /** 组态页名 → 若 HXScreen 已登记可接；否则占位 */
  function findTopoPage(pageHint) {
    if (!pageHint || !w.HXScreen || !w.HXScreen.findPageByName) return null;
    return w.HXScreen.findPageByName(pageHint) || null;
  }

  function resolveLeaf(name, ancestors) {
    var idx = catalogIndex();
    var path = (ancestors || []).concat([name]);
    var i, roomName, cabName, room, pair;
    var anc = ancestors || [];
    var node, pageHint, page;

    if (isTopoPath(anc)) {
      node = getNodeAt(path);
      /* 通讯管理机是末级叶子：点开即画该机 COM 通讯图，不再展开 COM */
      if (node && node.kind === 'mgr') {
        pageHint = node.pageHint || '';
        return {
          ok: true,
          kind: 'topo',
          design: name,
          pageHint: pageHint,
          room: { name: name, kind: 'topo', mgr: name },
          cab: null,
          mgr: name,
          com: '',
          mgrNode: node,
          rules: (netSpec && netSpec.rules) || null
        };
      }
      pageHint = (node && node.pageHint) || '';
      if (!pageHint && anc[1]) {
        var gwNode = getNodeAt([anc[0], anc[1]]);
        if (gwNode && gwNode.pageHint) pageHint = gwNode.pageHint;
      }
      /* 设备名优先：配电室3A1_… → 3A1配电室通讯图（管理机名可能是 3A2） */
      var dm = String(name).match(/^配电室([0-9A-Za-z]+)_/);
      if (dm) pageHint = dm[1] + '配电室通讯图';
      else {
        dm = String(name).match(/^(?:机房模块|SJJF)([0-9A-Za-z]+)/i);
        if (dm) pageHint = '机房模块' + dm[1] + '通讯图';
      }
      if (!pageHint && name === '通讯图入口' && anc[1]) {
        var mgrHead = String(anc[1]).split('(')[0].trim();
        if (/^ECC/i.test(mgrHead)) pageHint = mgrHead + '通讯图';
        else pageHint = mgrHead.replace(/配电室$/, '') + '配电室通讯图';
      }
      page = findTopoPage(pageHint);
      if (page && page.id) {
        return {
          ok: true,
          kind: 'topo',
          design: name,
          pageHint: pageHint,
          pageId: page.id,
          room: { name: pageHint || name, kind: 'topo', pageId: page.id },
          cab: null,
          mgr: anc[1] || '',
          com: /^COM\d+/i.test(anc[2] || '') ? anc[2] : '',
          rules: (netSpec && netSpec.rules) || null
        };
      }
      return {
        ok: false,
        kind: 'topo',
        reason: pageHint ? (pageHint + ' · 待做') : '通讯图待做',
        design: name,
        pageHint: pageHint,
        mgr: anc[1] || '',
        com: /^COM\d+/i.test(anc[2] || '') ? anc[2] : '',
        rules: (netSpec && netSpec.rules) || null
      };
    }

    if (isPowerLeaf(name)) {
      room = idx.byName[name];
      if (room) return { room: room, cab: null, ok: true, kind: 'power', design: name };
      var leafNode = getNodeAt(path);
      if (leafNode && leafNode.pageId) {
        room = { name: name, kind: 'power', pageId: leafNode.pageId };
        return { room: room, cab: null, ok: true, kind: 'power', design: name };
      }
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
    if (mode === 'power') {
      return function (b) { return b.name === '配电室' || b.name === '网络拓扑图'; };
    }
    if (mode === 'room') return function (b) { return b.name === '机房'; };
    if (mode === 'all' || mode === 'drill') return function () { return true; };
    return function () { return true; };
  }

  function markOpenPath(open, pathNames) {
    var acc = [];
    var i;
    if (!pathNames || !pathNames.length) return;
    for (i = 0; i < pathNames.length; i++) {
      acc.push(pathNames[i]);
      open[acc.join('\0')] = 1;
    }
  }

  /** 默认可展开：配电室 + 当前楼层；机房仅在有 prefer/active 时展开到目标；网络拓扑默认折叠 */
  function defaultOpenSet(preferRoom, active) {
    var open = {};
    var path = null;
    var leaf;
    var design;
    var code;
    open['配电室'] = 1;

    if (active && active.ancestors && active.name) {
      path = active.ancestors.concat([active.name]);
    } else if (active && active.archPath) {
      path = active.archPath;
    }

    if (path && path.length) {
      markOpenPath(open, path.slice(0, -1));
      return open;
    }

    design = preferRoom || '';
    if (/^配电室/.test(design) || /^ECC\d/.test(design)) {
      leaf = findLeafPath(design);
      if (leaf) markOpenPath(open, leaf.ancestors);
      else {
        code = roomCode(design);
        if (/^ECC/.test(design)) open['配电室\0ECC'] = 1;
        else if (code) open['配电室\0F' + code.charAt(0)] = 1;
      }
      return open;
    }

    if (/^机房模块|^数据机房/.test(design)) {
      leaf = findLeafPath(design) || findLeafPath(design.replace(/^机房模块/, '数据机房'));
      if (leaf) {
        markOpenPath(open, leaf.ancestors.concat([leaf.name]));
      } else {
        code = roomCode(design);
        open['机房'] = 1;
        /* 机房楼层名是 2F / 3F，不是 F2 */
        if (code) open['机房\0' + code.charAt(0) + 'F'] = 1;
      }
      delete open['配电室'];
      return open;
    }

    /* drill 默认：配电室展开 + F5（与 prefer 配电室5A1 一致） */
    open['配电室\0F5'] = 1;
    return open;
  }

  function pathKey(anc, name) {
    return anc.concat([name]).join('\0');
  }

  function renderTree(host, opts) {
    opts = opts || {};
    var mode = opts.mode || '';
    var onPick = opts.onPick || function () {};
    var active = opts.active || null;
    var preferRoom = opts.preferRoom || '';
    var branches = allBranches().filter(branchFilter(mode));
    if (!branches.length) {
      host.innerHTML = '<div class="hx-pt">架构树未加载。</div>';
      return { flat: [], openPrefer: null };
    }

    var flat = [];
    var openSet = defaultOpenSet(preferRoom, active);
    var ui = (spec && spec.ui) || {};
    var hideSoftPower = !!ui.hidePowerWithoutScene;
    host.innerHTML = '';

    function bindLeaf(btn) {
      btn.addEventListener('click', function () {
        var item = flat[Number(btn.getAttribute('data-i'))];
        if (!item || !item.resolved) return;
        onPick(item.resolved.room, item.resolved.cab, item);
        host.querySelectorAll('.hx-node').forEach(function (b) { b.classList.remove('is-on'); });
        btn.classList.add('is-on');
        expandToNode(host, btn);
      });
    }

    function appendLeaf(parentEl, node, depth, ancestors) {
      var name = node.name || '';
      var anc = ancestors.slice();
      var resolved = resolveLeaf(name, anc);
      var soft = resolved && !resolved.ok;
      if (hideSoftPower && soft && resolved && resolved.kind === 'power') return;
      var fi = flat.length;
      flat.push({ name: name, ancestors: anc, resolved: resolved, node: node });
      var pad = depth > 0 ? ' hx-sub' + Math.min(depth, 4) : '';
      var isOn = active && active.ancestors && active.name === name &&
        (active.ancestors || []).join('\0') === anc.join('\0');
      if (!isOn && active && active.room && resolved) {
        isOn = resolved.ok && resolved.room && active.room.name === resolved.room.name &&
          (active.cab == null ? !resolved.cab : (resolved.cab && active.cab && active.cab.name === resolved.cab.name));
      }
      var isMgr = node.kind === 'mgr' || (resolved && resolved.kind === 'topo' && resolved.mgr === name && !resolved.com);
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'hx-node' + pad + (soft ? ' is-soft' : '') + (isOn ? ' is-on' : '') +
        (resolved && resolved.kind === 'topo' ? ' is-topo' : '') +
        (isMgr ? ' is-mgr' : '');
      btn.setAttribute('data-i', String(fi));
      if (node.value) {
        var lab = document.createElement('span');
        lab.className = 'hx-node-lab';
        lab.textContent = name;
        var meta = document.createElement('span');
        meta.className = 'hx-node-meta';
        meta.textContent = node.value;
        btn.appendChild(lab);
        btn.appendChild(meta);
      } else {
        btn.textContent = name;
      }
      parentEl.appendChild(btn);
      bindLeaf(btn);
    }

    function fillKids(foldEl, node, depth, ancestors) {
      if (foldEl.getAttribute('data-filled') === '1') return;
      foldEl.setAttribute('data-filled', '1');
      var anc = ancestors.slice();
      (node.children || []).forEach(function (ch) {
        appendNode(foldEl, ch, depth, anc);
      });
    }

    function appendNode(parentEl, node, depth, ancestors) {
      var name = node.name || '';
      var kids = node.children || [];
      var anc = ancestors.slice();
      /* 通讯管理机：末级菜单，不展开 COM；点开右侧画通讯图 */
      if (!kids.length || node.kind === 'mgr') {
        appendLeaf(parentEl, node, depth, anc);
        return;
      }
      var key = pathKey(anc, name);
      var opened = !!openSet[key];
      var pad = depth > 0 ? ' hx-sub' + Math.min(depth, 4) : '';
      var block = document.createElement('div');
      block.className = 'hx-fold-block' + (opened ? '' : ' is-collapsed');
      block.setAttribute('data-path', key);
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'hx-folder' + pad;
      btn.setAttribute('data-fold', '1');
      btn.setAttribute('aria-expanded', opened ? 'true' : 'false');
      var tog = document.createElement('span');
      tog.className = 'hx-folder-tog';
      tog.setAttribute('aria-hidden', 'true');
      var lab = document.createElement('span');
      lab.className = 'hx-folder-lab';
      lab.textContent = name;
      btn.appendChild(tog);
      btn.appendChild(lab);
      if (node.value) {
        var meta = document.createElement('span');
        meta.className = 'hx-folder-meta';
        meta.textContent = node.value;
        btn.appendChild(meta);
      }
      var fold = document.createElement('div');
      fold.className = 'hx-fold';
      block.appendChild(btn);
      block.appendChild(fold);
      parentEl.appendChild(block);

      var childAnc = anc.concat([name]);
      btn.addEventListener('click', function (ev) {
        ev.preventDefault();
        var collapsed = block.classList.toggle('is-collapsed');
        btn.setAttribute('aria-expanded', collapsed ? 'false' : 'true');
        if (!collapsed) fillKids(fold, node, depth + 1, childAnc);
      });

      /* 默认展开的立刻填子节点；折叠的懒渲染，避免拓扑树一次挂几千节点 */
      if (opened) fillKids(fold, node, depth + 1, childAnc);
    }

    branches.forEach(function (b) { appendNode(host, b, 0, []); });

    return { flat: flat, host: host };
  }

  function expandToNode(host, btn) {
    if (!host || !btn) return;
    var el = btn.parentElement;
    while (el && el !== host) {
      if (el.classList && el.classList.contains('hx-fold-block')) {
        el.classList.remove('is-collapsed');
        var tog = el.firstElementChild;
        if (tog && tog.getAttribute('data-fold')) tog.setAttribute('aria-expanded', 'true');
      }
      el = el.parentElement;
    }
  }

  function activateFlatItem(flatOut, index) {
    var host = flatOut && flatOut.host;
    if (!host || index < 0) return;
    var btn = host.querySelector('.hx-node[data-i="' + index + '"]');
    if (!btn) return;
    host.querySelectorAll('.hx-node').forEach(function (b) { b.classList.remove('is-on'); });
    btn.classList.add('is-on');
    expandToNode(host, btn);
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

  /** 管理机下实际 COM 列（跳过「通讯图入口」等非 COM） */
  function mgrComPorts(mgrName) {
    var node = getNodeAt(['网络拓扑图', mgrName]);
    if (!node) return [];
    return (node.children || []).filter(function (ch) {
      return ch && ch.kind === 'com' && /^COM\d+/i.test(ch.name || '');
    }).slice().sort(function (a, b) {
      var na = parseInt(String(a.name).replace(/\D/g, ''), 10) || 0;
      var nb = parseInt(String(b.name).replace(/\D/g, ''), 10) || 0;
      return na - nb;
    });
  }

  w.HXArch = {
    branches: [],
    netRules: null,
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
    getNodeAt: getNodeAt,
    mgrComPorts: mgrComPorts,
    renderTree: renderTree,
    findPrefer: findPrefer,
    activateFlatItem: activateFlatItem,
    pickByDesign: pickByDesign,
    expandToNode: expandToNode
  };
})(window);

/* display_model_layer → hx-scenes 页面（与 scripts/hangxin_export_hx_scenes.py 同构） */
(function (w) {
  'use strict';

  var POWER_PAGE_RE = /^配电室\d[AB]\d$/;
  var COLUMN_PAGE_RE = /^(\d[AB]\d)数据机房_([A-K]列头)$/;

  function decodeBlob(raw) {
    if (!raw) return {};
    var text = String(raw).trim();
    if (!text) return {};
    if (text.charAt(0) === '{' || text.charAt(0) === '[') {
      try { return JSON.parse(text); } catch (e) { return {}; }
    }
    try {
      return JSON.parse(atob(text));
    } catch (e2) {
      return {};
    }
  }

  function keyFromPoint(pointName, nodeType) {
    var p = String(pointName || '');
    if (!p) return '';
    if (nodeType === 'S' || nodeType === 'I') {
      if (p.slice(-4) === '开关状态') return p.slice(0, -4);
      var i = p.lastIndexOf('_');
      if (i > 0) return p.slice(0, i);
    }
    return p;
  }

  function skipCapacitor(det) {
    try {
      return JSON.stringify(det || {}).indexOf('电容补偿') >= 0;
    } catch (e) {
      return false;
    }
  }

  function nodePos(c, style) {
    var pos = (c.position || style.position || {});
    var x = Number(pos.x || 0);
    var y = Number(pos.y || 0);
    var w = Number((c.size || {}).width || pos.w || 0);
    var h = Number((c.size || {}).height || pos.h || 0);
    var z = Number(style.zIndex || 0);
    return { x: x, y: y, w: w, h: h, z: z };
  }

  function convertCells(cells, layer) {
    var nodes = [];
    if (!cells || !cells.length) return nodes;
    cells.forEach(function (c) {
      var shape = c.shape || '';
      var det = ((c.data || {}).detail) || {};
      if (skipCapacitor(det)) return;
      var style = det.style || {};
      var np = nodePos(c, style);
      var cond = ((det.animate || {}).condition) || {};
      var x = np.x;
      var y = np.y;
      var ww = np.w;
      var hh = np.h;
      var z = np.z;

      if (shape === 'view-svg-line') {
        nodes.push({
          t: 'L', x: x, y: y, w: ww, h: hh, a: Number(c.angle || 0),
          c: style.foreColor || '#ffffff', sw: Number(style.borderWidth || 2), z: z, k: ''
        });
      } else if (shape === 'view-svg-text') {
        nodes.push({
          t: 'T', x: x, y: y, w: ww, h: hh, s: style.text || '',
          c: style.foreColor || '#fff', fs: style.fontSize || 18, z: z
        });
      } else if (shape === 'dv-border-box12') {
        nodes.push({ t: 'B', x: x, y: y, w: ww, h: hh, z: z });
      } else if (shape === 'view-svg-sld-out') {
        nodes.push({ t: 'O', x: x, y: y, w: ww, h: hh, z: z });
      } else if (shape === 'view-svg-sld-load-switch') {
        var ps = cond.dataName || '';
        if (!ps) return;
        nodes.push({
          t: 'S', x: x, y: y, w: ww, h: hh, c: '#00ff00',
          k: keyFromPoint(ps, 'S'), p: ps, dev: cond.DeviceName || '',
          op: cond.operator || '==',
          val: String(cond.OperatorValue != null ? cond.OperatorValue : '0'),
          z: z, decor: 1
        });
      } else if (shape === 'ism-view-png-image') {
        var pi = cond.dataName || '';
        if (!pi) return;
        var img = style.imageURL || '';
        var st = img.indexOf('.state0') >= 0 ? '0' : '1';
        nodes.push({
          t: 'I', x: x, y: y, w: ww, h: hh, st: st,
          k: keyFromPoint(pi, 'I'), p: pi, dev: cond.DeviceName || '',
          op: cond.operator || '==',
          val: String(cond.OperatorValue != null ? cond.OperatorValue : ''),
          z: z
        });
      }
    });
    return nodes;
  }

  function pageBounds(nodes, layer) {
    if (!nodes.length) return { w: 3150, h: 1400, x0: 0 };
    var maxX = 0;
    var maxY = 0;
    var minX = nodes[0].x;
    nodes.forEach(function (n) {
      maxX = Math.max(maxX, n.x + n.w);
      maxY = Math.max(maxY, n.y + n.h);
      minX = Math.min(minX, n.x);
    });
    var x0 = minX > 50 ? Math.max(0, Math.floor(minX) - 20) : 0;
    var w = Math.max(Math.floor(maxX) + 40, Number(layer.width || 0) || 3150);
    var h = Math.max(Math.floor(maxY) + 40, Number(layer.height || 0) || 1400);
    return { w: w, h: h, x0: x0 };
  }

  function classifyPage(pageName) {
    var m = COLUMN_PAGE_RE.exec(pageName);
    if (m) return { kind: 'column', group: m[1] + '数据机房', cab: m[2] };
    if (POWER_PAGE_RE.test(pageName)) return { kind: 'power', group: '配电一次图', cab: '' };
    return null;
  }

  function buildPageFromLayerRow(row) {
    var pageName = row.PageName || row.page_name || '';
    var cls = classifyPage(pageName);
    if (!cls) return null;
    var comp = decodeBlob(row.components || row.Components);
    var layer = decodeBlob(row.layer || row.Layer);
    var cells = comp.cells;
    if (!cells || !cells.length) return null;
    var nodes = convertCells(cells, layer);
    if (!nodes.length) return null;
    var b = pageBounds(nodes, layer);
    return {
      id: row.PageId || row.page_id,
      name: pageName,
      kind: cls.kind,
      group: cls.group,
      cab: cls.cab,
      home: (row.IsHome || row.is_home) ? 1 : 0,
      w: b.w,
      h: b.h,
      x0: b.x0,
      bg: layer.backgroundImage || '/static/ISM/page/background/background07.state0.jpg',
      nodes: nodes
    };
  }

  w.HXSceneImport = {
    powerPageRe: POWER_PAGE_RE,
    buildPageFromLayerRow: buildPageFromLayerRow,
    classifyPage: classifyPage
  };
})(window);

import fs from 'fs';
import vm from 'vm';
import path from 'path';

export function bootHx(root) {
  const layoutCode = fs.readFileSync(path.join(root, 'js/hx-5a1-layout.js'), 'utf8');
  const roomCode = fs.readFileSync(path.join(root, 'js/hx-power-room-layout.js'), 'utf8');
  const screenCode = fs.readFileSync(path.join(root, 'js/hx-screen.js'), 'utf8');
  const layout = JSON.parse(fs.readFileSync(path.join(root, 'data/hx-5a1-cad-layout.json'), 'utf8'));
  const scenes = JSON.parse(fs.readFileSync(path.join(root, 'data/hx-scenes.json'), 'utf8'));
  const roomDir = path.join(root, 'data/hx-power-rooms');
  const roomIndex = JSON.parse(fs.readFileSync(path.join(roomDir, 'index.json'), 'utf8'));
  const rooms = {};
  (roomIndex.rooms || []).forEach((r) => {
    rooms[r.code] = JSON.parse(fs.readFileSync(path.join(roomDir, r.file), 'utf8'));
  });
  const sandbox = { window: {}, console, w: null };
  sandbox.w = sandbox.window;
  sandbox.window.HXLive = { onChange: null, connect: () => {}, ready: Promise.resolve() };
  sandbox.window.HXArch = { ready: Promise.resolve() };
  vm.createContext(sandbox);
  vm.runInContext(layoutCode, sandbox);
  vm.runInContext(roomCode, sandbox);
  const cut = screenCode.indexOf('  L().onChange = paint;');
  const lib =
    screenCode.slice(0, cut) +
    '\nw.runPrep = function(pg){ if(!pg||pg._prepped)return pg; fixPowerLinesFromCadHints(pg); pg._prepped=1; return pg; };\n' +
    'w.__knifeInTap = knifeInTap;\n' +
    'w.__pickPrimary = pickPrimaryKnives;\n' +
    '})(window);\n';
  vm.runInContext(lib, sandbox);
  sandbox.window.HX5A1.setLayout(layout);
  sandbox.window.HXRooms.setLayouts(rooms);
  return { sandbox, layout, scenes, rooms };
}

export function prepRoom(ctx, name) {
  const pg0 = ctx.scenes.pages.find((p) => p.name === name);
  const pg = JSON.parse(JSON.stringify(pg0));
  ctx.sandbox.window.runPrep(pg);
  return pg;
}

export function lineEnds(n) {
  const rad = ((Number(n.a) || 0) * Math.PI) / 180;
  const cos = Math.cos(rad);
  const sin = Math.sin(rad);
  const cx = n.x + n.w / 2;
  const cy = n.y + n.h / 2;
  const hw = n.w / 2;
  return {
    cx,
    cy,
    x1: cx - hw * cos,
    y1: cy - hw * sin,
    x2: cx + hw * cos,
    y2: cy + hw * sin
  };
}

export function segments(nodes) {
  return nodes
    .filter((n) => n.t === 'L')
    .map((n) => {
      const e = lineEnds(n);
      const horiz = Math.abs(e.y1 - e.y2) < 1.2;
      return Object.assign({ horiz, net: n.net || '', c: n.c, n }, e);
    });
}

export function touches(seg, x, y, tol) {
  if (seg.horiz) {
    const cy = (seg.y1 + seg.y2) / 2;
    if (Math.abs(cy - y) > tol) return false;
    const x0 = Math.min(seg.x1, seg.x2);
    const x1 = Math.max(seg.x1, seg.x2);
    return x >= x0 - tol && x <= x1 + tol;
  }
  const cx = (seg.x1 + seg.x2) / 2;
  if (Math.abs(cx - x) > tol) return false;
  const y0 = Math.min(seg.y1, seg.y2);
  const y1 = Math.max(seg.y1, seg.y2);
  return y >= y0 - tol && y <= y1 + tol;
}

export function cluster(segs, tol) {
  const parent = segs.map((_, i) => i);
  function find(i) {
    while (parent[i] !== i) {
      parent[i] = parent[parent[i]];
      i = parent[i];
    }
    return i;
  }
  function union(a, b) {
    parent[find(a)] = find(b);
  }
  for (let i = 0; i < segs.length; i++) {
    for (let j = i + 1; j < segs.length; j++) {
      if (
        touches(segs[j], segs[i].x1, segs[i].y1, tol) ||
        touches(segs[j], segs[i].x2, segs[i].y2, tol) ||
        touches(segs[i], segs[j].x1, segs[j].y1, tol) ||
        touches(segs[i], segs[j].x2, segs[j].y2, tol)
      ) {
        union(i, j);
      }
    }
  }
  return { find };
}

export function idsTouching(segs, x, y, tol) {
  const ids = [];
  segs.forEach((s, i) => {
    if (touches(s, x, y, tol)) ids.push(i);
  });
  return ids;
}

export function linked(segs, uf, ax, ay, bx, by, tol) {
  const A = idsTouching(segs, ax, ay, tol);
  const B = idsTouching(segs, bx, by, tol);
  if (!A.length || !B.length) return false;
  const roots = new Set(A.map((i) => uf.find(i)));
  return B.some((i) => roots.has(uf.find(i)));
}

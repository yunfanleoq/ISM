/* =========================================================
   自监控：四台平台服务器的监控台
   角色来自 /bjga/ops 的 ha；资源只读本机 host，其余标未接入
   ========================================================= */
(function (A) {
  'use strict';

  var LOCAL = '2号';
  var cpuHist = [];

  A.register('selfmon', {
    render: function (root) {
      if (window.__smTimer) {
        clearInterval(window.__smTimer);
        window.__smTimer = 0;
      }
      paint(root, readHa(), null, []);
      pull(root);
      window.__smTimer = setInterval(function () {
        if (!root.isConnected) {
          clearInterval(window.__smTimer);
          window.__smTimer = 0;
          return;
        }
        pull(root);
      }, 4000);
    }
  });

  function readHa() {
    var extra = (window.DCIM && window.DCIM.extra) || {};
    return extra.ha || [];
  }

  function pull(root) {
    var t0 = (window.performance && performance.now) ? performance.now() : Date.now();
    A.api('/bjga/ops').then(function (j) {
      if (!j || j.code !== 1000) return;
      var rtt = Math.round(((window.performance && performance.now) ? performance.now() : Date.now()) - t0);
      if (!window.DCIM) window.DCIM = {};
      if (!window.DCIM.extra) window.DCIM.extra = {};
      if (j.ha) window.DCIM.extra.ha = j.ha;
      var sig = roleSig(j.ha || []);
      if (root.getAttribute('data-sig') !== sig || !root.querySelector('[data-srv]')) {
        paint(root, j.ha || [], j.host || null, j.audits || [], rtt);
      } else {
        applyLive(root, j.host || null, j.audits || [], rtt);
      }
    });
  }

  function roleSig(ha) {
    return (ha || []).map(function (r) { return r.name + ':' + r.role; }).join('|');
  }

  function paint(root, ha, host, audits, rtt) {
    ha = order(ha);
    var local = ha.filter(function (s) { return s.name === LOCAL; })[0];
    var primary = ha.filter(function (s) { return s.role === '主'; })[0];
    var standby = ha.filter(function (s) { return s.role === '备'; })[0];
    var off = ha.filter(function (s) { return s.name !== LOCAL; }).length;
    root.setAttribute('data-sig', roleSig(ha));
    root.innerHTML =
      '<div class="page-head">' +
        '<h1 class="page-title">自监控</h1>' +
        '<span class="page-crumb">' + A.crumb('平台运行状态') + '</span>' +
        '<div class="page-tools"><span class="muted" id="smAt" style="font-size:12px;font-family:var(--num)"></span>' +
          '<button class="btn is-ghost is-sm" data-act="refresh">立即刷新</button>' +
          '<button class="btn is-primary is-sm" data-act="diag">一键诊断</button></div>' +
      '</div>' +
      '<p class="page-lead">看的是四台平台服务器，不是机房里的动环设备。角色标在每台卡片上。只有正在跑服务的这一台有实时资源，其余显示未接入。</p>' +
      '<div class="card" style="flex:0 0 auto"><div class="card-body" style="padding:2px 0">' +
        A.UI.stats([
          { v: ha.length || 4, u: '台', l: '平台服务器' },
          { v: '—', l: '本机 CPU', id: 'stCpu' },
          { v: '—', l: '本机内存', id: 'stMem' },
          { v: off, u: '台', l: '未接入采集' },
          { v: '—', l: '本机协程', id: 'stGo' },
          { v: primary ? primary.name : '—', l: '当前承担业务' }
        ]) +
      '</div></div>' +
      '<div class="row grow">' +
        '<section class="card grow">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">服务器资源</h2>' +
            '<div class="card-extra"><span class="muted" style="font-size:12px">4 秒刷新</span></div></div>' +
          '<div class="card-body" style="overflow:auto"><div class="srv-grid">' +
            ha.map(function (s) { return card(s, s.name === LOCAL); }).join('') +
          '</div></div>' +
        '</section>' +
        '<section class="card" style="flex:0 0 420px">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">热备</h2>' +
            '<div class="card-extra"><button class="btn is-primary is-sm" id="smHa" type="button">演练切换</button></div></div>' +
          '<div class="card-body col" style="gap:12px">' +
            haStrip(primary, standby) +
            '<div><div class="side-title" style="margin-bottom:8px">本机进程</div>' +
              '<div id="smProc"></div></div>' +
          '</div>' +
        '</section>' +
      '</div>' +
      '<section class="card" style="flex:0 0 220px">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">动态指令记录</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">审计里的真实操作</span></div></div>' +
        '<div class="card-body" id="smCmd"></div>' +
      '</section>';

    bind(root, ha, local);
    applyLive(root, host, audits, rtt);
  }

  function haStrip(primary, standby) {
    return '<div class="srv" style="gap:8px">' +
      '<div class="srv-hd"><i class="dot" style="background:var(--green)"></i><span class="srv-nm">' +
        A.esc(primary ? primary.name : '—') + '</span><span class="tag is-green" style="margin-left:auto">主 · 承担业务</span></div>' +
      '<div class="srv-hd"><i class="dot" style="background:var(--blue)"></i><span class="srv-nm">' +
        A.esc(standby ? standby.name : '—') + '</span><span class="tag is-blue" style="margin-left:auto">备 · 故障后顶上</span></div>' +
      '<div style="font-size:12px;color:var(--t-muted)">热备心跳未接入采集。切换只对调 2 号和 3 号，1 号和 4 号不动。</div></div>';
  }

  function card(s, local) {
    var tone = roleTone(s, local);
    var bars = local
      ? resRow('CPU', 'cpu') + resRow('内存', 'mem') + resRow('磁盘', 'disk')
      : '<div class="srv-row"><span style="color:var(--t-label)">资源</span><span style="color:var(--t-muted)">未接入</span><span></span></div>' +
        '<div class="srv-row"><span style="color:var(--t-label)">心跳</span><span style="color:var(--t-muted)">未接入</span><span></span></div>';
    return '<div class="srv' + (local ? '' : ' is-off') + '" data-srv="' + A.esc(s.name) + '">' +
      '<div class="srv-hd">' +
        '<i class="dot" style="background:' + tone.color + '"></i>' +
        '<span class="srv-nm">' + A.esc(s.name) + '</span>' +
        '<span class="tag ' + tone.tag + '" style="margin-left:auto">' + A.esc(roleText(s)) + '</span>' +
      '</div>' +
      '<div style="font-size:12px;color:var(--t-muted)">规格 ' + A.esc(s.spec || '未列') +
        (local ? ' · <b style="color:var(--cyan)">本机实时</b>' : ' · 未接入采集') + '</div>' +
      bars +
      (local
        ? '<div style="display:flex;justify-content:space-between;align-items:center">' +
            '<span style="font-size:12px;color:var(--t-muted)">采样 <b id="smRtt" style="font-family:var(--num);color:var(--cyan)">—</b> ms</span>' +
            '<span id="smSpark"></span></div>'
        : '') +
    '</div>';
  }

  function resRow(label, key) {
    return '<div class="srv-row"><span style="color:var(--t-label)">' + label + '</span>' +
      '<span class="bar-track"><i class="bar-fill" data-bar="' + key + '" style="width:0%;background:var(--cyan)"></i></span>' +
      '<span data-pct="' + key + '" style="font-family:var(--num);color:var(--cyan);text-align:right">—</span></div>';
  }

  function applyLive(root, host, audits, rtt) {
    var at = A.q('#smAt', root);
    if (at) at.textContent = host && host.at ? ('采样 ' + host.at) : '';
    if (!host) {
      setStat(root, 'stCpu', '—');
      setStat(root, 'stMem', '—');
      setStat(root, 'stGo', '—');
      return;
    }
    cpuHist.push(host.cpu);
    if (cpuHist.length > 20) cpuHist.shift();
    setBar(root, 'cpu', host.cpu, host.cpu > 80);
    setBar(root, 'mem', host.mem, host.mem > 85);
    setBar(root, 'disk', host.disk, host.disk > 90);
    setStat(root, 'stCpu', A.fix(host.cpu, 1) + '%');
    setStat(root, 'stMem', A.fix(host.mem, 1) + '%');
    setStat(root, 'stGo', String(host.goroutine));
    var rttEl = A.q('#smRtt', root);
    if (rttEl && rtt != null) rttEl.textContent = String(rtt);
    var spark = A.q('#smSpark', root);
    if (spark && cpuHist.length) {
      var data = cpuHist.length < 2 ? [host.cpu, host.cpu] : cpuHist.slice();
      spark.innerHTML = A.Charts.spark({ w: 120, h: 30, data: data, color: '#00D8FF' });
    }
    var proc = A.q('#smProc', root);
    if (proc) {
      proc.innerHTML = A.UI.table({
        cols: [
          { t: '进程', k: 'n', w: 120 },
          { t: '协程', k: 'g', align: 'r', w: 70 },
          { t: '内存', k: 'm', align: 'r', w: 90 },
          { t: '已运行', k: 'u', w: 90 }
        ],
        rows: [{
          n: 'ism_server',
          g: host.goroutine,
          m: host.rssMb + ' MB',
          u: fmtUp(host.uptimeSec)
        }]
      });
    }
    var cmd = A.q('#smCmd', root);
    if (cmd) {
      var rows = (audits || []).map(function (r) {
        return { t: r.at, c: (r.action || '') + (r.object ? ' · ' + r.object : ''), u: r.actor || '', d: r.after || '' };
      });
      cmd.innerHTML = rows.length ? A.UI.table({
        cols: [
          { t: '时间', k: 't', w: 170, mono: true },
          { t: '指令', k: 'c', w: 280 },
          { t: '操作人', k: 'u', w: 120 },
          { t: '结果', k: 'd', w: 160 }
        ],
        rows: rows
      }) : '<div class="muted">暂无指令记录</div>';
    }
  }

  function setBar(root, key, v, warn) {
    var col = warn ? 'var(--orange)' : 'var(--cyan)';
    var bar = root.querySelector('[data-bar="' + key + '"]');
    var pct = root.querySelector('[data-pct="' + key + '"]');
    var n = Math.max(0, Math.min(100, Number(v) || 0));
    if (bar) {
      bar.style.width = n + '%';
      bar.style.background = col;
    }
    if (pct) {
      pct.style.color = col;
      pct.textContent = A.fix(n, 1) + '%';
    }
  }

  function setStat(root, id, text) {
    var map = { stCpu: 1, stMem: 2, stGo: 4 };
    var boxes = root.querySelectorAll('.stat-box .v');
    var el = boxes[map[id]];
    if (el) el.textContent = text;
  }

  function bind(root, ha, local) {
    var haBtn = A.q('#smHa', root);
    if (haBtn) haBtn.addEventListener('click', function () {
      if (!window.confirm('确认执行主备切换演练？只对调 2 号和 3 号，1 号南向和 4 号备份不动。')) return;
      A.api('/bjga/ops/ha/switch', { method: 'POST', body: {} }).then(function (res) {
        if (!res || res.code !== 1000) { A.toast('切换失败'); return; }
        if (window.DCIM && window.DCIM.extra) window.DCIM.extra.ha = res.ha;
        A.toast('主备角色已对调');
        pull(root);
      });
    });
    A.q('[data-act="refresh"]', root).addEventListener('click', function () {
      pull(root);
      A.toast('已读取本机实时');
    });
    A.q('[data-act="diag"]', root).addEventListener('click', function () {
      var cpu = root.querySelector('[data-pct="cpu"]');
      A.toast('本机 ' + (cpu ? cpu.textContent : '—') + '。1 号、3 号、4 号未接入采集');
    });
    A.qa('.srv', root).forEach(function (c) {
      c.addEventListener('click', function () {
        var name = c.getAttribute('data-srv');
        if (!name) return;
        var s = ha.filter(function (x) { return x.name === name; })[0];
        if (!s) return;
        var live = name === LOCAL;
        A.drawer(s.name + ' · ' + roleText(s), A.UI.kv([
          ['机器', s.name],
          ['角色', roleText(s)],
          ['规格', s.spec || '未列'],
          ['资源', live ? '本机实时，见卡片' : '未接入采集'],
          ['心跳', '未接入采集'],
          ['对端', s.peer || '不在热备对']
        ]));
      });
    });
  }

  function order(ha) {
    var rank = { '1号': 1, '2号': 2, '3号': 3, '4号': 4 };
    return (ha || []).slice().sort(function (a, b) {
      return (rank[a.name] || 9) - (rank[b.name] || 9);
    });
  }

  function fmtUp(sec) {
    sec = Number(sec) || 0;
    if (sec < 60) return sec + '秒';
    if (sec < 3600) return Math.floor(sec / 60) + '分';
    return Math.floor(sec / 3600) + '时' + Math.floor((sec % 3600) / 60) + '分';
  }

  function roleText(s) {
    if (s.kind === 'south' || s.role === '南向') return '南向采集';
    if (s.kind === 'backup' || s.role === '备份') return '备份机';
    if (s.role === '主') return '热备主';
    if (s.role === '备') return '热备备';
    return s.role || '';
  }

  function roleTone(s, local) {
    if (!local && (s.kind === 'south' || s.role === '南向' || s.kind === 'backup' || s.role === '备份')) {
      return { tag: 'is-grey', color: 'var(--t-muted)' };
    }
    if (s.kind === 'south' || s.role === '南向') return { tag: 'is-blue', color: 'var(--blue)' };
    if (s.role === '主') return { tag: 'is-green', color: 'var(--green)' };
    if (s.role === '备') return { tag: 'is-blue', color: 'var(--blue)' };
    return { tag: 'is-grey', color: 'var(--t-muted)' };
  }
})(window.App);

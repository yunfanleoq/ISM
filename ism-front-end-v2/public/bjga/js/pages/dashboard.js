/* =========================================================
   驾驶舱大屏 —— 严格沿用设计稿布局 / 配色 / 间距
   ========================================================= */
(function (A) {
  'use strict';

  var BASE = { pue: 1.50, upsPower: 1158.97, upsEnergy: 1685.6, itPower: 1158.97, cabRate: 97.96 };

  function applyLive() {
    var live = window.DCIM && window.DCIM.live;
    if (!live) return;
    if (live.pue != null) BASE.pue = Math.round(Number(live.pue) * 100) / 100;
    if (live.upsPower != null) BASE.upsPower = Number(live.upsPower);
    if (live.upsEnergy != null) BASE.upsEnergy = Number(live.upsEnergy);
    if (live.itPower != null) BASE.itPower = Number(live.itPower);
    if (live.cabRate != null) BASE.cabRate = Number(live.cabRate);
  }

  var POWER_HOURS = ['00:00', '02:00', '04:00', '06:00', '08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00', '22:00'];
  var ENV_HOURS = ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'];

  function liveSnap() {
    var devs = (window.DCIM && window.DCIM.live && window.DCIM.live.devices) || [];
    function vals(kind, name) {
      var out = [];
      devs.forEach(function (d) {
        if (d.kind !== kind) return;
        (d.points || []).forEach(function (p) {
          if (p.name === name && p.value != null && p.value !== '') out.push(String(p.value));
        });
      });
      return out;
    }
    function avg(kind, name) {
      var xs = vals(kind, name).map(Number).filter(function (n) { return isFinite(n); });
      if (!xs.length) return null;
      return xs.reduce(function (a, b) { return a + b; }, 0) / xs.length;
    }
    function first(kind, name) {
      var xs = vals(kind, name);
      return xs.length ? xs[0] : null;
    }
    return {
      temp: avg('th', '温度'),
      humi: avg('th', '湿度'),
      kw: first('pwr', '总有功功率'),
      leak: first('leak', '漏水信号'),
      leakM: first('leak', '漏水位置米数')
    };
  }

  /* 驾驶舱告警优先用未清除采集告警 */
  function todayLevels() {
    var D = window.DCIM, out = { 紧急: 0, 严重: 0, 重要: 0, 次要: 0, 预警: 0 };
    var list = (D && D.live && Array.isArray(D.live.alarms)) ? D.live.alarms : (D && D.alarms);
    if (!list || !list.length) return out;
    var latest = new Date(String(list[0].time || '').slice(0, 10).replace(/-/g, '/')).getTime();
    var from = latest - 6 * 86400000;
    list.forEach(function (a) {
      if (!a.time) return;
      var t = new Date(a.time.slice(0, 10).replace(/-/g, '/')).getTime();
      if (t < from) return;
      if (out[a.level] != null) out[a.level] += 1;
    });
    return out;
  }

  var timer = null;

  A.register('dashboard', {
    unmount: function () { if (timer) { clearInterval(timer); timer = null; } },

    render: function (root) {
      applyLive();
      var snap = liveSnap();
      var kwShow = BASE.itPower;
      var powerTitle = 'IT 总功率';
      var powerLabel = '当前功率';
      BASE._kwLive = false;
      if (snap.kw != null && isFinite(Number(snap.kw))) {
        kwShow = Number(snap.kw);
        BASE._kwLive = true;
        powerTitle = '列头柜有功功率';
        powerLabel = '采集功率';
      }
      var assetN = (window.DCIM && window.DCIM.meta && window.DCIM.meta.assetTotal) || 4293;
      var lv = todayLevels();
      var lvTotal = lv.紧急 + lv.严重 + lv.重要 + lv.次要 + lv.预警;
      var pending = lv.紧急 + lv.严重;

      root.innerHTML =
        '<main class="content">' +

          /* ---------- 左栏 ---------- */
          '<aside class="col col-left">' +

            '<section class="card">' +
              '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">能耗 PUE</h2></div>' +
              '<div class="card-body pue-body">' +
                '<div class="gauge">' +
                  '<div class="gauge-svg" id="pueGauge"></div>' +
                  '<div class="gauge-val" id="pueValue">' + A.fix(BASE.pue, 2) + '</div>' +
                  '<div class="gauge-label">PUE</div>' +
                '</div>' +
                '<div class="pue-stats">' +
                  '<div class="stat"><span class="stat-label">UPS 输出总功率</span>' +
                    '<span class="stat-row"><b class="stat-num" id="upsPower">' + A.fix(BASE.upsPower, 2) + '</b><span class="stat-unit">kW</span></span></div>' +
                  '<div class="stat"><span class="stat-label">UPS 输出总电度</span>' +
                    '<span class="stat-row"><b class="stat-num" id="upsEnergy">' + A.fix(BASE.upsEnergy, 1) + '</b><span class="stat-unit">万度</span></span></div>' +
                '</div>' +
              '</div>' +
            '</section>' +

            '<section class="card">' +
              '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">部署设备数量</h2></div>' +
              '<div class="summary">' +
                '<div class="sum-item"><span class="sum-label">单位</span><b class="sum-num">5</b></div>' +
                '<div class="sum-item"><span class="sum-label">系统</span><b class="sum-num">33</b></div>' +
                '<div class="sum-item"><span class="sum-label">设备</span><b class="sum-num">' + assetN + '</b></div>' +
              '</div>' +
              '<ul class="dev-list">' +
                devRow('服务器', 3063, srvIcon()) +
                devRow('交换机', 652, swIcon()) +
                devRow('存储设备', 176, stoIcon()) +
                devRow('网络安全设备', 283, secIcon()) +
                devRow('其他设备', 119, othIcon()) +
              '</ul>' +
            '</section>' +

            '<section class="card">' +
              '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">机柜使用情况</h2></div>' +
              '<div class="card-body cab-body">' +
                '<div class="ring">' +
                  '<div class="ring-svg" id="cabRing"></div>' +
                  '<div class="ring-val" id="cabRate">' + A.fix(BASE.cabRate, 2) + '%</div>' +
                  '<div class="ring-label">机柜使用率(%)</div>' +
                '</div>' +
                '<div class="cab-stats">' +
                  '<div class="stat"><span class="stat-label">机柜总数</span>' +
                    '<span class="stat-row"><b class="stat-num sm c-cyan">588</b><span class="stat-unit">个</span></span></div>' +
                  '<div class="stat"><span class="stat-label">已启用机柜</span>' +
                    '<span class="stat-row"><b class="stat-num sm c-green">576</b><span class="stat-unit">个</span></span></div>' +
                  '<div class="stat"><span class="stat-label">未启用机柜</span>' +
                    '<span class="stat-row"><b class="stat-num sm c-orange">12</b><span class="stat-unit">个</span></span></div>' +
                '</div>' +
              '</div>' +
            '</section>' +

          '</aside>' +

          /* ---------- 中央 ---------- */
          '<section class="col col-center">' +
            '<div class="stage" id="stage" data-floor="main">' +
              '<div class="stage-bg" id="stageBg"></div>' +
              '<h1 class="stage-title">数据中心动环管理系统</h1>' +
              '<div class="stage-sub">驾驶舱 · ' + A.esc(window.DCIM.meta.project) + '</div>' +
              '<i class="corner tl"></i><i class="corner tr"></i><i class="corner br"></i><i class="corner bl"></i>' +
              '<button class="btn is-ghost" type="button" id="btnTour" ' +
                'style="position:absolute;right:14px;top:14px;z-index:6;font-size:12px;padding:6px 12px;background:rgba(3,14,32,.9);backdrop-filter:blur(8px)">进入空间巡检 ›</button>' +
              '<div class="floor-bar" id="floorBar">' +
                '<button class="fbtn is-active" type="button" data-floor="main">主界面</button>' +
                '<button class="fbtn" type="button" data-floor="1F">1F</button>' +
                '<button class="fbtn" type="button" data-floor="2F">2F</button>' +
                '<button class="fbtn" type="button" data-floor="3F">3F</button>' +
                '<button class="fbtn" type="button" data-floor="4F">4F</button>' +
              '</div>' +
            '</div>' +
            '<div class="floor-strip" id="floorStrip">' +
              window.DCIM.floors.slice().sort(function (a, b) {
                return a.floor < b.floor ? -1 : 1;
              }).map(function (f) {
                var hot = f.alarmN >= 500;
                var mrooms = window.DCIM.rooms.filter(function (r) {
                  return r.floor === f.floor && (r.type === '主机房' || r.type === '屏蔽机房');
                });
                f.main = mrooms.length;
                f.shield = mrooms.length > 0 && mrooms.every(function (r) { return r.type === '屏蔽机房'; });
                return '<button class="fs-item" type="button" data-floor="' + f.floor + '">' +
                  '<span class="fs-head"><span class="fs-tag">' + f.floor + '</span>' +
                    '<span class="fs-name">' + A.esc(f.main ? (f.shield ? '屏蔽机房层 · ' : '主机房层 · ') + f.main + ' 间' : '配套与弱电层') + '</span></span>' +
                  '<span class="fs-kv">机房 <b>' + f.rooms + '</b> 间 · 机柜 <b>' + A.num(f.cabCount) + '</b> 架</span>' +
                  '<span class="fs-al">告警 <b' + (hot ? ' style="color:var(--red)"' : '') + '>' +
                    A.num(f.alarmN) + '</b> 条</span>' +
                '</button>';
              }).join('') +
            '</div>' +
            '<div class="metrics">' +
              metric('9745', '㎡', '数据中心总面积') +
              metric('588', '个', '机柜总数') +
              metric('21', '台', 'UPS 总数') +
              metric('71', '台', '空调总数') +
            '</div>' +
          '</section>' +

          /* ---------- 右栏 ---------- */
          '<aside class="col col-right">' +

            '<section class="card">' +
              '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + powerTitle + '</h2></div>' +
              '<div class="power-row">' +
                '<span class="power-label">' + powerLabel + '</span>' +
                '<b class="power-num" id="itPower">' + A.fix(kwShow, 2) + '</b>' +
                '<span class="power-unit">kW</span>' +
              '</div>' +
              '<div class="chart" id="powerChart"></div>' +
            '</section>' +

            '<section class="card">' +
              '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">机房温湿度</h2></div>' +
              '<div class="legend">' +
                '<span class="legend-item"><i class="legend-dot" style="background:#FFA24D"></i><span class="legend-text">温度 <b id="liveTemp">' +
                  (snap.temp == null ? '—' : A.fix(snap.temp, 1)) + '</b> ℃</span></span>' +
                '<span class="legend-item"><i class="legend-dot" style="background:#2E9CF6"></i><span class="legend-text">湿度 <b id="liveHumi">' +
                  (snap.humi == null ? '—' : A.fix(snap.humi, 1)) + '</b> %</span></span>' +
              '</div>' +
              '<div id="liveStrip" style="padding:0 4px 8px;font-size:13px;color:#9FD7FF">漏水信号 <b id="liveLeak">' +
                A.esc(snap.leak == null ? '—' : String(snap.leak)) + '</b> · 位置 <b id="liveLeakM">' +
                A.esc(snap.leakM == null ? '—' : String(snap.leakM)) + '</b> m · 列头柜 <b id="liveKw">' +
                A.esc(snap.kw == null ? '—' : String(snap.kw)) + '</b> kW</div>' +
              '<div class="chart" id="envChart"></div>' +
            '</section>' +

            '<section class="card">' +
              '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">告警统计（近 7 日）</h2></div>' +
              '<div class="alert-grid">' +
                alertItem('紧急告警', lv.紧急, '#FF4D5E') +
                alertItem('严重告警', lv.严重, '#FF7849') +
                alertItem('重要告警', lv.重要, '#FFA24D') +
                alertItem('次要告警', lv.次要, '#FFD24D') +
                alertItem('预警', lv.预警, '#00D8FF') +
              '</div>' +
              '<div class="alert-note"><span>近 7 日 ' + A.num(lvTotal) + ' 条</span>' +
                '<span class="alert-note-sep">·</span>' +
                '<span>紧急/严重 ' + A.num(pending) + ' 条</span></div>' +
            '</section>' +

          '</aside>' +
        '</main>';

      /* 图表 */
      A.q('#pueGauge', root).innerHTML = A.Charts.gauge({ value: BASE.pue, min: 1, max: 2 });
      A.q('#cabRing', root).innerHTML = A.Charts.ring({ size: 150, r: 58, sw: 12, rate: BASE.cabRate / 100 });
      A.q('#powerChart', root).innerHTML = A.Charts.line({
        w: 396, h: 202, labels: POWER_HOURS, yLabels: false,
        min: BASE._kwLive ? 0 : 1060,
        max: BASE._kwLive ? Math.max(kwShow * 2, 20) : 1260,
        series: [{ color: '#00D8FF', area: true, data: BASE._kwLive ? POWER_HOURS.map(function () { return kwShow; }) : A.wave(BASE.itPower, 12, 24, 7) }]
      });
      A.q('#envChart', root).innerHTML = A.Charts.line({
        w: 396, h: 214, labels: ENV_HOURS, yLabels: false,
        series: [
          { color: '#FFA24D', area: true, norm: [18, 28], data: snap.temp == null ? A.wave(23.5, 6, 1.0, 11) : ENV_HOURS.map(function () { return snap.temp; }) },
          { color: '#2E9CF6', norm: [40, 60], data: snap.humi == null ? A.wave(48, 6, 2.2, 13) : ENV_HOURS.map(function () { return snap.humi; }) }
        ]
      });
    },

    mount: function (root) {
      applyLive();
      var pueG = A.q('#pueGauge', root), pueV = A.q('#pueValue', root);
      var upsP = A.q('#upsPower', root), upsE = A.q('#upsEnergy', root), itP = A.q('#itPower', root);
      var state = { pue: BASE.pue, power: BASE.itPower };

      function draw() {
        pueG.innerHTML = A.Charts.gauge({ value: state.pue, min: 1, max: 2 });
        pueV.textContent = state.pue.toFixed(2);
        upsP.textContent = state.power.toFixed(2);
        if (!BASE._kwLive) itP.textContent = state.power.toFixed(2);
      }
      draw();
      upsE.textContent = A.fix(BASE.upsEnergy, 1);

      /* PUE 仍是演示计算，允许轻微波动；采集功率不被改写 */
      timer = setInterval(function () {
        state.pue = A.clamp(BASE.pue + (Math.random() - 0.5) * 0.02, 1.2, 1.9);
        if (!BASE._kwLive) state.power = BASE.itPower + (Math.random() - 0.5) * 6;
        draw();
      }, 4000);
      setTimeout(function () {
        if (!timer) return;
      }, 5000);

      /* 楼层切换 */
      var stage = A.q('#stage', root), floorBar = A.q('#floorBar', root);
      var stageBg = A.q('#stageBg', root);

      /* 楼体线框：按项目实景绘制 */
      if (stageBg && window.BLD) stageBg.innerHTML = window.BLD.svg({ floor: 'main' });

      var floorStrip = A.q('#floorStrip', root);

      function setFloor(f) {
        stage.setAttribute('data-floor', f);
        var idx = /^[1-9]/.test(f) ? Number(f.charAt(0)) : 0;
        A.qa('.fbtn', floorBar).forEach(function (el) {
          el.classList.toggle('is-active', el.getAttribute('data-floor') === f);
        });
        A.qa('.fs-item', floorStrip).forEach(function (el) {
          el.classList.toggle('is-active', el.getAttribute('data-floor') === f);
        });
        A.qa('.bld-svg .win[data-f]', stageBg).forEach(function (w) {
          w.classList.toggle('is-hl', Number(w.getAttribute('data-f')) === idx);
        });
        A.qa('.bld-svg .flr[data-f]', stageBg).forEach(function (g) {
          g.classList.toggle('is-hl', Number(g.getAttribute('data-f')) === idx);
        });
      }

      floorBar.addEventListener('click', function (e) {
        var b = e.target.closest('.fbtn');
        if (b) setFloor(b.getAttribute('data-floor'));
      });

      /* 楼层概况卡：点击定位到该层，再次点击回到主界面 */
      floorStrip.addEventListener('click', function (e) {
        var it = e.target.closest('.fs-item');
        if (!it) return;
        var f = it.getAttribute('data-floor');
        setFloor(stage.getAttribute('data-floor') === f ? 'main' : f);
      });

      /* 3D 楼宇 / 楼层按钮 → 空间巡检 */
      var bg = A.q('#stageBg', root);
      if (bg) {
        bg.style.cursor = 'pointer';
        bg.addEventListener('click', function () { A.go('building'); });
      }
      var btnTour = A.q('#btnTour', root);
      if (btnTour) btnTour.addEventListener('click', function () { A.go('building'); });

      /* 告警卡片下钻 */
      var alertCard = A.q('.alert-grid', root);
      if (alertCard) {
        alertCard.style.cursor = 'pointer';
        alertCard.addEventListener('click', function () { A.go('facility'); });
      }
    }
  });

  /* ---------- 片段 ---------- */
  function devRow(name, count, icon) {
    return '<li class="dev-row"><span class="dev-icon">' + icon + '</span>' +
      '<span class="dev-name">' + name + '</span><b class="dev-count">' + A.num(count) + '</b></li>';
  }
  function metric(v, u, l) {
    return '<div class="metric"><span class="metric-row"><b class="metric-num">' + v + '</b>' +
      '<span class="metric-unit">' + u + '</span></span><span class="metric-label">' + l + '</span></div>';
  }
  function alertItem(n, v, c) {
    return '<div class="alert-item"><i class="alert-dot" style="--c:' + c + '"></i>' +
      '<span class="alert-name">' + n + '</span><b class="alert-num" style="color:' + c + '">' + v + '</b></div>';
  }
  function srvIcon() {
    return '<svg width="18" height="18" viewBox="0 0 18 18" fill="none"><rect x="2" y="2.5" width="14" height="5" rx="1" stroke="#00D8FF" stroke-width="1.3"/><rect x="2" y="10.5" width="14" height="5" rx="1" stroke="#00D8FF" stroke-width="1.3"/><circle cx="5" cy="5" r="0.9" fill="#00D8FF"/><circle cx="5" cy="13" r="0.9" fill="#00D8FF"/></svg>';
  }
  function swIcon() {
    return '<svg width="18" height="18" viewBox="0 0 18 18" fill="none"><path d="M2 6H16M2 12H16" stroke="#00D8FF" stroke-width="1.3" stroke-linecap="round"/><path d="M13 3.5L16 6L13 8.5M5 9.5L2 12L5 14.5" stroke="#00D8FF" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>';
  }
  function stoIcon() {
    return '<svg width="18" height="18" viewBox="0 0 18 18" fill="none"><ellipse cx="9" cy="4.5" rx="6.5" ry="2.5" stroke="#00D8FF" stroke-width="1.3"/><path d="M2.5 4.5V13.5C2.5 14.9 5.4 16 9 16C12.6 16 15.5 14.9 15.5 13.5V4.5" stroke="#00D8FF" stroke-width="1.3"/><path d="M2.5 9C2.5 10.4 5.4 11.5 9 11.5C12.6 11.5 15.5 10.4 15.5 9" stroke="#00D8FF" stroke-width="1.3"/></svg>';
  }
  function secIcon() {
    return '<svg width="18" height="18" viewBox="0 0 18 18" fill="none"><path d="M9 1.5L15.5 4V8.5C15.5 12.5 12.8 15.5 9 16.5C5.2 15.5 2.5 12.5 2.5 8.5V4L9 1.5Z" stroke="#00D8FF" stroke-width="1.3" stroke-linejoin="round"/><path d="M6.5 9L8.3 10.8L11.8 7" stroke="#00D8FF" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>';
  }
  function othIcon() {
    return '<svg width="18" height="18" viewBox="0 0 18 18" fill="none"><rect x="2.5" y="4" width="13" height="10" rx="1" stroke="#00D8FF" stroke-width="1.3"/><path d="M2.5 7.5H15.5" stroke="#00D8FF" stroke-width="1.3"/><circle cx="5" cy="11.5" r="0.9" fill="#00D8FF"/></svg>';
  }
})(window.App);

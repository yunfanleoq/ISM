/* 活动告警：四级口径、轮询 /bjga/live、房间定位弹层 */
(function (A) {
  'use strict';

  var TIER_OF = { 紧急: '紧急', 严重: '紧急', 重要: '重要', 次要: '一般', 预警: '提示' };
  var TIER_COLOR = { 紧急: '#FF4D5E', 重要: '#FF7849', 一般: '#FFA24D', 提示: '#FFD24D' };
  var TIER_RANK = { 紧急: 4, 重要: 3, 一般: 2, 提示: 1 };

  A.alarmTier = function (level) {
    return TIER_OF[level] || TIER_OF[String(level || '')] || '提示';
  };
  A.alarmTierColor = function (tier) {
    return TIER_COLOR[tier] || TIER_COLOR['提示'];
  };

  function dcim() { return window.DCIM; }

  function isActive(a) {
    return a && a.recover !== '已恢复';
  }

  function alarmList() {
    var D = dcim();
    var live = D && D.live && D.live.alarms;
    if (Array.isArray(live) && live.length) return live;
    return (D && D.alarms) || [];
  }

  A.activeAlarms = function () {
    return alarmList().filter(isActive);
  };

  A.resolveAlarmPlace = function (a) {
    if (!a) return { room: '', floor: '' };
    if (a.room) return { room: a.room, floor: a.floor || '' };
    var obj = String(a.obj || '').trim();
    var rooms = (dcim() && dcim().rooms) || [];
    var i, r, key;
    for (i = 0; i < rooms.length; i++) {
      r = rooms[i];
      if (!obj) break;
      if (r.name === obj || (obj.length >= 2 && r.name.indexOf(obj) >= 0) || r.name.indexOf(obj.slice(0, 3)) >= 0) {
        return { room: r.name, floor: r.floor };
      }
    }
    var th = (dcim() && dcim().thDevices) || [];
    var best = '', bestLen = 0, place = null;
    for (i = 0; i < th.length; i++) {
      key = th[i].name || '';
      if (!key || !obj) continue;
      if (obj.indexOf(key) >= 0 || key.indexOf(obj) >= 0) {
        if (key.length > bestLen) {
          bestLen = key.length;
          best = key;
          place = th[i];
        }
      }
    }
    if (place) return { room: place.room || '', floor: place.floor || '' };
    return { room: '', floor: '' };
  };

  function fp(a) {
    return [a.time, a.obj, a.what, a.level].join('|');
  }

  var seen = {};
  var pollTimer = null;
  var pollMs = 2000;
  var locateOpen = false;

  function mergeLive(live) {
    if (!live || !dcim()) return;
    dcim().live = dcim().live || {};
    if (live.devices) dcim().live.devices = live.devices;
    if (live.alarms) dcim().live.alarms = live.alarms;
    Object.keys(live).forEach(function (k) {
      if (k !== 'devices' && k !== 'alarms') dcim().live[k] = live[k];
    });
  }

  function pickNewAlarms(list) {
    var fresh = [];
    list.forEach(function (a) {
      if (!isActive(a)) return;
      var id = fp(a);
      if (!seen[id]) {
        seen[id] = 1;
        fresh.push(a);
      }
    });
    return fresh;
  }

  function seedSeen(list) {
    list.forEach(function (a) {
      if (isActive(a)) seen[fp(a)] = 1;
    });
  }

  function highestTier(alarms) {
    var best = null, rank = 0;
    alarms.forEach(function (a) {
      var t = A.alarmTier(a.level);
      var r = TIER_RANK[t] || 0;
      if (r >= rank) { rank = r; best = a; }
    });
    return best;
  }

  function floorPlanHtml(floor, highlightRoom, tierColor) {
    var rooms = (dcim().rooms || []).filter(function (r) { return r.floor === floor; });
    if (!rooms.length) {
      return '<div class="al-plan-empty">暂无 ' + A.esc(floor) + ' 平面数据</div>';
    }
    return '<div class="al-plan-grid">' + rooms.map(function (r) {
      var on = r.name === highlightRoom;
      return '<div class="al-plan-room' + (on ? ' is-alarm' : '') + '"' +
        (on ? ' style="--al-c:' + tierColor + '"' : '') + ' data-room="' + A.esc(r.name) + '">' +
        '<div class="al-plan-nm">' + A.esc(r.name) + '</div>' +
        '<div class="al-plan-meta">' + A.esc(r.type || '') + ' · ' + r.cabCount + ' 柜</div>' +
        (on ? '<div class="al-plan-tag">告警定位</div>' : '') +
      '</div>';
    }).join('') + '</div>';
  }

  A.closeAlarmLocate = function () {
    locateOpen = false;
    A.qa('.alarm-locate-mask').forEach(function (n) { n.remove(); });
  };

  A.showAlarmLocate = function (alarm) {
    if (!alarm) return;
    var place = A.resolveAlarmPlace(alarm);
    var tier = A.alarmTier(alarm.level);
    var color = A.alarmTierColor(tier);
    var floor = place.floor || (dcim().rooms[0] && dcim().rooms[0].floor) || '3F';
    var room = place.room || '（未匹配机房）';

    A.closeAlarmLocate();
    locateOpen = true;
    var mask = A.h(
      '<div class="alarm-locate-mask" role="dialog" aria-label="告警房间定位">' +
        '<div class="alarm-locate-panel">' +
          '<header class="alarm-locate-head">' +
            '<div><span class="al-tier" style="color:' + color + ';border-color:' + color + '">' + A.esc(tier) + '</span>' +
            '<b class="alarm-locate-title">' + A.esc(alarm.what || '活动告警') + '</b></div>' +
            '<button type="button" class="alarm-locate-x" aria-label="关闭">✕</button>' +
          '</header>' +
          '<div class="alarm-locate-sub">' +
            '<span>对象 <b>' + A.esc(alarm.obj || '-') + '</b></span>' +
            '<span>机房 <b>' + A.esc(room) + '</b></span>' +
            '<span>楼层 <b>' + A.esc(floor) + '</b></span>' +
            '<span>时间 <b class="mono">' + A.esc(alarm.time || '-') + '</b></span>' +
          '</div>' +
          '<div class="alarm-locate-body">' +
            '<div class="alarm-locate-cap"><i class="card-bar"></i>' + A.esc(floor) + ' · 电子地图（演示）</div>' +
            floorPlanHtml(floor, place.room, color) +
          '</div>' +
          '<footer class="alarm-locate-foot">' +
            '<button type="button" class="btn is-ghost" data-act="fac">打开运行监控</button>' +
            '<button type="button" class="btn is-primary" data-act="ok">知道了</button>' +
          '</footer>' +
        '</div>' +
      '</div>');

    var host = document.getElementById('screen') || document.body;
    host.appendChild(mask);
    mask.addEventListener('click', function (e) {
      if (e.target === mask) A.closeAlarmLocate();
    });
    A.q('.alarm-locate-x', mask).addEventListener('click', A.closeAlarmLocate);
    A.q('[data-act="ok"]', mask).addEventListener('click', A.closeAlarmLocate);
    A.q('[data-act="fac"]', mask).addEventListener('click', function () {
      A.closeAlarmLocate();
      A.go('facility/alarm');
    });
    var hi = A.q('.al-plan-room.is-alarm', mask);
    if (hi && hi.scrollIntoView) {
      setTimeout(function () { hi.scrollIntoView({ block: 'center', behavior: 'smooth' }); }, 80);
    }
  };

  function onPoll(live) {
    mergeLive(live);
    var active = A.activeAlarms();
    var fresh = pickNewAlarms(active);
    if (fresh.length && !locateOpen) {
      var top = highestTier(fresh);
      if (top) A.showAlarmLocate(top);
    }
    try {
      window.dispatchEvent(new CustomEvent('dcim:alarms', { detail: { active: active, fresh: fresh } }));
    } catch (e) {}
  }

  A.startAlarmPoll = function (ms) {
    pollMs = ms || pollMs;
    if (pollTimer) return;
    seedSeen(A.activeAlarms());
    function tick() {
      if (!A.api) return;
      A.api('/bjga/live').then(function (j) {
        if (!j || j.code !== 1000 || !j.live) return;
        onPoll(j.live);
      }).catch(function () {});
    }
    tick();
    pollTimer = setInterval(tick, pollMs);
  };

  A.stopAlarmPoll = function () {
    if (pollTimer) { clearInterval(pollTimer); pollTimer = null; }
  };

  A.tierCounts = function (list) {
    var out = { 紧急: 0, 重要: 0, 一般: 0, 提示: 0 };
    (list || []).forEach(function (a) {
      var t = A.alarmTier(a.level);
      if (out[t] != null) out[t]++;
    });
    return out;
  };
})(window.App);

/* =========================================================
   系统管理：用户 / 角色 / 组织 / 页面配置 / 消息 / 备份 / 日志 / 字典
   ========================================================= */
(function (A) {
  'use strict';

  var MENU = [
    { id: 'user', name: '用户管理' }, { id: 'role', name: '角色权限' },
    { id: 'org', name: '组织部门' }, { id: 'page', name: '页面配置' },
    { id: 'msg', name: '消息中心' }, { id: 'backup', name: '备份还原' },
    { id: 'log', name: '系统日志' }, { id: 'dict', name: '字典管理' },
    { id: 'api', name: '接口配置' }
  ];
  var state = { m: 'user' };

  var DATA = {
    user: {
      title: '用户管理',
      cols: [
        { t: '账号', k: 'a', w: 150, mono: true }, { t: '姓名', k: 'n', w: 120 },
        { t: '所属部门', k: 'd', w: 160 }, { t: '角色', k: 'r', w: 140 },
        { t: '手机号', k: 'p', w: 140, mono: true }, { t: '最近登录', k: 'l', w: 170, mono: true },
        { t: '状态', k: 's', w: 100, html: function (r) { return r.s === '启用' ? A.UI.tag('启用', 'green') : A.UI.tag('停用', 'grey'); } },
        { t: '操作', k: 'x', w: 160, html: function () { return '<button class="btn is-ghost is-sm" data-op="edit">编辑</button> <button class="btn is-ghost is-sm" data-op="reset">重置密码</button>'; } }
      ],
      rows: [
        { a: 'admin', n: '系统管理员', d: '技术保障科', r: '超级管理员', p: '138****0001', l: '2026-10-20 09:12:44', s: '启用' },
        { a: 'ops01', n: '杨XX', d: '运维中心', r: '运维值班员', p: '138****2312', l: '2026-10-19 22:40:07', s: '启用' },
        { a: 'ops02', n: '李XX', d: '运维中心', r: '运维值班员', p: '138****5543', l: '2026-10-19 18:02:31', s: '启用' },
        { a: 'ops03', n: '张XX', d: '运维中心', r: '运维工程师', p: '139****7788', l: '2026-10-18 08:31:15', s: '启用' },
        { a: 'mgr01', n: '王XX', d: '科信部门', r: '部门管理员', p: '137****9012', l: '2026-10-17 15:22:09', s: '启用' },
        { a: 'view01', n: '刘XX', d: '后勤保障部', r: '只读用户', p: '136****3345', l: '2026-10-12 10:05:52', s: '停用' },
        { a: 'audit01', n: '陈XX', d: '网络安全保卫', r: '审计员', p: '135****6677', l: '2026-10-16 14:18:26', s: '启用' }
      ]
    },
    role: {
      title: '角色权限',
      cols: [
        { t: '角色名称', k: 'n', w: 180 }, { t: '标识', k: 'c', w: 160, mono: true },
        { t: '用户数', k: 'u', align: 'r', w: 90 }, { t: '权限项', k: 'p', align: 'r', w: 100 },
        { t: '数据范围', k: 'd', w: 180 },
        { t: '状态', k: 's', w: 100, html: function (r) { return A.UI.tag(r.s, 'green'); } }
      ],
      rows: [
        { n: '超级管理员', c: 'ROLE_SUPER', u: 1, p: 286, d: '全部数据', s: '启用' },
        { n: '部门管理员', c: 'ROLE_MGR', u: 3, p: 164, d: '本部门及下级', s: '启用' },
        { n: '运维值班员', c: 'ROLE_OPS', u: 12, p: 98, d: '所属机房', s: '启用' },
        { n: '运维工程师', c: 'ROLE_ENG', u: 8, p: 132, d: '全部机房', s: '启用' },
        { n: '审计员', c: 'ROLE_AUDIT', u: 2, p: 42, d: '只读 + 日志', s: '启用' },
        { n: '只读用户', c: 'ROLE_VIEW', u: 26, p: 18, d: '授权机房', s: '启用' }
      ]
    },
    org: {
      title: '组织部门',
      cols: [
        { t: '单位名称', k: 'n', w: 200 }, { t: '部门', k: 'd', w: 180 },
        { t: '负责人', k: 'o', w: 110 }, { t: '设备数', k: 'c', align: 'r', w: 100 },
        { t: '机房数', k: 'r', align: 'r', w: 100 }, { t: '联系电话', k: 'p', w: 140, mono: true }
      ],
      rows: []
    },
    page: {
      title: '页面配置',
      cols: [
        { t: '配置项', k: 'n', w: 220 }, { t: '当前值', k: 'v', w: 260 },
        { t: '可选值', k: 'o', w: 240 }, { t: '说明', k: 'd', w: 300 }
      ],
      rows: [
        { n: '登录页皮肤', v: '深蓝科技风', o: '深蓝科技风 / 浅色简约 / 深空灰', d: '控制登录页与首页背景主题' },
        { n: '首页默认页', v: '总览', o: '总览 / 设施监控 / 资产总览', d: '用户登录后进入的默认页面' },
        { n: '导航显示样式', v: '横向页签', o: '横向页签 / 左侧菜单', d: '顶部模块导航的呈现方式' },
        { n: '大屏分辨率基准', v: '1920 × 1080', o: '1920×1080 / 2560×1440 / 3840×2160', d: '大屏等比缩放基准' },
        { n: '自动刷新间隔', v: '5 秒', o: '3 / 5 / 10 / 30 秒', d: '实时数据的刷新周期' },
        { n: '告警声音提示', v: '开启', o: '开启 / 关闭', d: '紧急告警是否播放提示音' },
        { n: '温度单位', v: '摄氏度 ℃', o: '摄氏度 / 华氏度', d: '全站温度显示单位' }
      ]
    },
    msg: {
      title: '消息中心',
      cols: [
        { t: '消息类型', k: 't', w: 140 }, { t: '标题', k: 'n', w: 320 },
        { t: '接收人', k: 'r', w: 120 }, { t: '发送时间', k: 's', w: 170, mono: true },
        { t: '状态', k: 'c', w: 100, html: function (r) { return r.c === '已读' ? A.UI.tag('已读', 'grey') : A.UI.tag('未读', 'orange'); } }
      ],
      rows: [
        { t: '告警通知', n: '417 空调间 6 号空调通信中断', r: '值班员', s: '2026-10-19 17:35:42', c: '已读' },
        { t: '告警通知', n: '211 电池室氢气浓度超限', r: '运维中心', s: '2026-10-19 09:12:04', c: '未读' },
        { t: '系统通知', n: '历史库归档任务已完成（9 月数据）', r: '系统管理员', s: '2026-10-18 02:00:11', c: '已读' },
        { t: '工单通知', n: '工单 WO-20261018-014 已派发', r: '运维工程师', s: '2026-10-18 10:24:36', c: '已读' },
        { t: '巡检提醒', n: '今日 2F 机房巡检任务待执行', r: '值班员', s: '2026-10-20 08:00:00', c: '未读' }
      ]
    },
    backup: {
      title: '备份还原',
      cols: [
        { t: '备份名称', k: 'n', w: 260 }, { t: '类型', k: 't', w: 110 },
        { t: '大小', k: 's', w: 110 }, { t: '耗时', k: 'd', w: 100 },
        { t: '创建时间', k: 'c', w: 170, mono: true },
        { t: '状态', k: 'x', w: 100, html: function () { return A.UI.tag('成功', 'green'); } }
      ],
      rows: [
        { n: 'DCIM_FULL_20261019.dat', t: '全量备份', s: '4.82 GB', d: '6分12秒', c: '2026-10-19 23:00:00' },
        { n: 'DCIM_CFG_20261019.dat', t: '配置备份', s: '18.4 MB', d: '8秒', c: '2026-10-19 23:05:10' },
        { n: 'DCIM_FULL_20261012.dat', t: '全量备份', s: '4.76 GB', d: '5分58秒', c: '2026-10-12 23:00:00' },
        { n: 'DCIM_HIS_20261001.dat', t: '历史备份', s: '186 GB', d: '41分02秒', c: '2026-10-01 01:00:00' }
      ]
    },
    log: {
      title: '系统日志',
      cols: [
        { t: '时间', k: 't', w: 170, mono: true }, { t: '级别', k: 'l', w: 90, html: function (r) {
          return r.l === 'ERROR' ? A.UI.tag('ERROR', 'red') : r.l === 'WARN' ? A.UI.tag('WARN', 'orange') : A.UI.tag('INFO', 'green');
        } }, { t: '模块', k: 'm', w: 140 }, { t: '操作人', k: 'u', w: 110 }, { t: '内容', k: 'c', w: 520 }
      ],
      rows: [
        { t: '2026-10-20 09:12:44', l: 'INFO', m: '用户认证', u: 'admin', c: '管理员登录成功，来源 10.18.32.11' },
        { t: '2026-10-20 08:00:02', l: 'INFO', m: '定时任务', u: 'system', c: '日报生成任务启动，共 6 个报表模板' },
        { t: '2026-10-19 23:00:14', l: 'INFO', m: '备份任务', u: 'system', c: '全量备份完成 DCIM_FULL_20261019.dat' },
        { t: '2026-10-19 21:14:08', l: 'ERROR', m: '采集驱动', u: 'system', c: '驱动 BatteryMon_211 通讯中断，已启动重连' },
        { t: '2026-10-19 17:35:42', l: 'WARN', m: '告警服务', u: 'ops01', c: '产生重要告警：417 空调间 6 号空调通信中断' },
        { t: '2026-10-19 14:02:19', l: 'INFO', m: '资产管理', u: 'ops03', c: '更新设备档案 FWQ2026100042，修改责任人' },
        { t: '2026-10-19 10:31:55', l: 'WARN', m: '容量分析', u: 'system', c: '413 机房空间使用率超过 85%，建议扩容评估' }
      ]
    },
    dict: {
      title: '字典管理',
      cols: [
        { t: '字典名称', k: 'n', w: 180 }, { t: '字典编码', k: 'c', w: 180, mono: true },
        { t: '条目数', k: 'i', align: 'r', w: 90 }, { t: '说明', k: 'd', w: 380 }
      ],
      rows: [
        { n: '设备类别', c: 'DICT_DEV_CAT', i: 5, d: '服务器 / 交换机 / 存储 / 网络安全 / 其他' },
        { n: '告警等级', c: 'DICT_ALM_LV', i: 5, d: '紧急 / 严重 / 重要 / 次要 / 预警' },
        { n: '使用状态', c: 'DICT_USE_ST', i: 4, d: '使用中 / 备用 / 故障维修 / 下架待处理' },
        { n: '联网类型', c: 'DICT_NET', i: 4, d: '互联网 / 公安网 / 视频专网 / 涉密网' },
        { n: '机房楼层', c: 'DICT_FLOOR', i: 4, d: '1F / 2F / 3F / 4F' },
        { n: '告警分类', c: 'DICT_ALM_CAT', i: 6, d: '供电配电 / 监控采集 / 环境空调 / 安防 / 消防 / 其他' }
      ]
    },
    api: {
      title: '接口配置',
      cols: [
        { t: '接口名称', k: 'n', w: 200 }, { t: '地址', k: 'u', w: 320, mono: true },
        { t: '方法', k: 'm', w: 90 }, { t: '鉴权', k: 'a', w: 100 },
        { t: '调用次数/日', k: 'c', align: 'r', w: 120 }, { t: '状态', k: 's', w: 100, html: function (r) {
          return r.s === '正常' ? A.UI.tag('正常', 'green') : A.UI.tag('异常', 'red'); } }
      ],
      rows: [
        { n: '实时数据查询', u: '/api/v1/realtime/values', m: 'POST', a: 'Token', c: 428610, s: '正常' },
        { n: '历史数据查询', u: '/api/v1/history/query', m: 'POST', a: 'Token', c: 18624, s: '正常' },
        { n: '告警推送', u: '/api/v1/alarm/push', m: 'POST', a: 'Token', c: 12084, s: '正常' },
        { n: '资产同步', u: '/api/v1/asset/sync', m: 'PUT', a: 'Token', c: 412, s: '正常' },
        { n: '第三方平台对接', u: '/api/v1/open/dcim', m: 'POST', a: 'AK/SK', c: 0, s: '异常' }
      ]
    }
  };

  function orgRows() {
    var list = (window.DCIM && window.DCIM.assetStats && window.DCIM.assetStats.byOrg) || [];
    return list.slice(0, 8).map(function (o, i) {
      return {
        n: o.name,
        d: ['技术保障科', '运维中心', '指挥中心', '后勤保障部'][i % 4],
        o: ['杨XX', '李XX', '王XX', '张XX'][i % 4],
        c: o.value,
        r: 2 + (i % 5),
        p: '5312000' + (i % 9)
      };
    });
  }

  function readStore(store) {
    var keys, i, raw, obj;
    try { keys = Object.keys(store); } catch (e) { return null; }
    for (i = 0; i < keys.length; i++) {
      if (/pass|token|authorization/i.test(keys[i])) continue;
      try { raw = store.getItem(keys[i]); } catch (e2) { continue; }
      if (!raw || raw.charAt(0) !== '{') continue;
      try { obj = JSON.parse(raw); } catch (e3) { continue; }
      if (!obj || typeof obj.name !== 'string' || !obj.name) continue;
      if (!(obj.Role || obj.Uuid)) continue;
      if (obj.password || obj.Password) continue;
      return { name: obj.name, role: obj.Role || '' };
    }
    return null;
  }

  function loginUser() {
    var found = null;
    try { found = readStore(window.localStorage); } catch (e) {}
    if (found) return found;
    try {
      if (window.parent && window.parent !== window) found = readStore(window.parent.localStorage);
    } catch (e2) {}
    return found;
  }

  function barName() {
    var el = document.querySelector('.user-name');
    return el ? String(el.textContent || '').trim() : '';
  }

  function userText() {
    var bar = barName() || '（顶栏未标注）';
    var u = loginUser();
    var text = '看本屏右上角，当前标注为「' + bar + '」。';
    if (u && u.name) {
      text += '平台登录姓名是「' + u.name + '」' + (u.role ? '，角色 ' + u.role : '') + '，和外层系统顶栏头像旁是同一身份。';
    } else {
      text += '平台登录姓名在外层系统顶栏头像旁边。';
    }
    return text + '这里只展示，不能改密码。';
  }

  function projectText() {
    var m = (window.DCIM && window.DCIM.meta) || {};
    if (!m.project && !m.system) return '项目名称尚未装载';
    return [m.project, m.system].filter(Boolean).join(' · ');
  }

  function settingsHtml() {
    return '<p class="sys-lead">本页是系统设置，查看当前项目、登录身份、数据刷新和热备。不在这里改密码。</p>' +
      '<section class="card sys-board">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">系统设置</h2></div>' +
        '<div class="card-body"><div class="sys-grid">' +
          '<div class="sys-item"><div class="k">当前项目</div><div class="v">' + A.esc(projectText()) + '</div></div>' +
          '<div class="sys-item"><div class="k">登录用户</div><div class="v">' + A.esc(userText()) + '</div></div>' +
          '<div class="sys-item"><div class="k">数据刷新</div><div class="v">打开本大屏时请求一次 /bjga/bootstrap，装入项目、告警和采集。总览上的 PUE 使用这一次算出的分项，不再每隔几秒改数字。要看最新采集，重新进入页面或刷新浏览器。自监控里的「立即刷新」只更新那一页的主机状态。</div></div>' +
          '<div class="sys-item"><div class="k">审计入口</div><div class="v"><a class="sys-link" href="#/audit">进入安全审计</a></div></div>' +
          '<div class="sys-item is-wide"><div class="k">热备状态（只读）</div><div class="v" id="sysHa">读取中…</div></div>' +
        '</div></div></section>';
  }

  function paintHa(root) {
    var box = A.q('#sysHa', root);
    if (!box || !A.api) {
      if (box) box.innerHTML = '<a class="sys-link" href="#/ops/ha">到运维中心查看</a>';
      return;
    }
    A.api('/bjga/ops').then(function (j) {
      var ha = (j && j.code === 1000 && j.ha) || [];
      if (!box.isConnected) return;
      if (!ha.length) {
        box.innerHTML = '<a class="sys-link" href="#/ops/ha">到运维中心查看</a>';
        return;
      }
      box.innerHTML = '<table class="tbl sys-ha"><thead><tr><th>节点</th><th>角色</th><th>对端</th><th>心跳 ms</th><th>状态</th></tr></thead><tbody>' +
        ha.map(function (r) {
          return '<tr><td>' + A.esc(r.name) + '</td><td>' + A.esc(r.role) + '</td><td>' + A.esc(r.peer) +
            '</td><td class="mono">' + A.esc(r.delayMs) + '</td><td>' + A.esc(r.status) + '</td></tr>';
        }).join('') + '</tbody></table>';
    }).catch(function () {
      if (box.isConnected) box.innerHTML = '<a class="sys-link" href="#/ops/ha">到运维中心查看</a>';
    });
  }

  A.register('admin', {
    render: function (root) {
      root.innerHTML =
        '<div class="page-head">' +
          '<h1 class="page-title">系统管理</h1>' +
          '<span class="page-crumb">' + A.crumb('平台配置') + '</span>' +
        '</div>' +
        settingsHtml() +
        '<div class="admin grow">' +
          '<aside class="side" style="flex:0 0 200px">' +
            '<div class="side-title">管理项</div>' +
            '<div class="side-scroll admin-menu">' +
              MENU.map(function (m) {
                return '<div class="admin-mi' + (m.id === state.m ? ' is-active' : '') + '" data-m="' + m.id + '">' +
                  m.name + '</div>';
              }).join('') +
            '</div>' +
          '</aside>' +
          '<section class="card grow">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title" id="adTitle"></h2>' +
              '<div class="card-extra">' +
                '<button class="btn is-primary is-sm" data-ad="new">+ 新增</button>' +
                '<button class="btn is-ghost is-sm" data-ad="exp">导出</button>' +
              '</div></div>' +
            '<div class="card-body" id="adBody"></div>' +
          '</section>' +
        '</div>';
      this.paint(root);
      paintHa(root);
    },

    paint: function (root) {
      var d = DATA[state.m];
      if (state.m === 'org') d = { title: d.title, cols: d.cols, rows: orgRows() };
      A.q('#adTitle', root).textContent = d.title;
      A.q('#adBody', root).innerHTML = A.UI.table({ cols: d.cols, rows: d.rows });

      A.qa('.admin-mi', root).forEach(function (m) {
        m.addEventListener('click', function () {
          state.m = m.getAttribute('data-m');
          A.qa('.admin-mi', root).forEach(function (x) { x.classList.toggle('is-active', x === m); });
          A.pages.admin.paint(root);
        });
      });
      A.q('[data-ad="new"]', root).addEventListener('click', function () {
        A.toast('演示环境：新增记录已提交（未持久化）');
      });
      A.q('[data-ad="exp"]', root).addEventListener('click', function () {
        A.toast(d.title + ' 已导出');
      });
      A.qa('[data-op]', root).forEach(function (b) {
        b.addEventListener('click', function (e) {
          e.stopPropagation();
          A.toast(b.getAttribute('data-op') === 'edit' ? '打开编辑弹窗（演示）' : '密码重置链接已发送');
        });
      });
    }
  });
})(window.App);

/* =========================================================
   工程配置：工程组态 / 融合组态编辑器 / 驱动 / 告警 / 存储 / 联动 / 函数
   ========================================================= */
(function (A) {
  'use strict';

  var SUBS = [
    { id: 'proj', name: '工程组态' }, { id: 'fusion', name: '融合组态' },
    { id: 'driver', name: '驱动管理' }, { id: 'alarm', name: '告警管理' },
    { id: 'storage', name: '存储配置' }, { id: 'link', name: '联动策略' },
    { id: 'func', name: '函数管理' }
  ];
  var state = { sub: 'fusion', page: 0 };
  var nodes = [];
  var seq = 0;

  /* ---------------- 图元库 ---------------- */
  var LIB = [
    { g: '基础图元', items: [
      { t: '矩形', k: 'rect', ic: '<rect x="3" y="6" width="22" height="14" rx="2" stroke="#00D8FF" stroke-width="1.4" fill="none"/>' },
      { t: '圆形', k: 'circle', ic: '<circle cx="14" cy="13" r="8" stroke="#00D8FF" stroke-width="1.4" fill="none"/>' },
      { t: '管线', k: 'pipe', ic: '<path d="M3 13h6l3-6h6l3 6h6" stroke="#00D8FF" stroke-width="1.4" fill="none"/>' },
      { t: '机柜', k: 'rack', ic: '<rect x="6" y="4" width="16" height="18" rx="2" stroke="#00D8FF" stroke-width="1.4" fill="none"/><path d="M6 9h16M6 14h16" stroke="#00D8FF" stroke-width="1"/>' }
    ] },
    { g: '图表', items: [
      { t: '趋势曲线', k: 'line', ic: '<path d="M3 17l5-6 4 4 5-8 5 5" stroke="#00D8FF" stroke-width="1.6" fill="none"/>' },
      { t: '柱状图', k: 'bar', ic: '<rect x="5" y="10" width="4" height="8" fill="#00D8FF"/><rect x="12" y="6" width="4" height="12" fill="#00D8FF"/><rect x="19" y="13" width="4" height="5" fill="#00D8FF"/>' },
      { t: '饼图', k: 'pie', ic: '<circle cx="14" cy="13" r="8" stroke="#00D8FF" stroke-width="1.4" fill="none"/><path d="M14 13V5a8 8 0 018 8" fill="#00D8FF" fill-opacity=".5"/>' }
    ] },
    { g: '仪表', items: [
      { t: '仪表盘', k: 'gauge', ic: '<path d="M5 19a9 9 0 1118 0" stroke="#00D8FF" stroke-width="1.6" fill="none"/><path d="M14 13l5-4" stroke="#00D8FF" stroke-width="1.6"/>' },
      { t: '数值', k: 'value', ic: '<text x="5" y="18" font-size="12" fill="#00D8FF" font-family="monospace">23.5</text>' },
      { t: '指示灯', k: 'lamp', ic: '<circle cx="14" cy="13" r="7" fill="#00E39A" fill-opacity=".5" stroke="#00E39A" stroke-width="1.4"/>' }
    ] },
    { g: '其他', items: [
      { t: '文本', k: 'text', ic: '<path d="M5 7h18M10 7v13" stroke="#00D8FF" stroke-width="1.5" fill="none"/>' },
      { t: '图片', k: 'image', ic: '<rect x="4" y="6" width="20" height="15" rx="2" stroke="#00D8FF" stroke-width="1.4" fill="none"/><circle cx="10" cy="11" r="2" fill="#00D8FF"/>' },
      { t: '视频', k: 'video', ic: '<rect x="4" y="7" width="20" height="13" rx="2" stroke="#00D8FF" stroke-width="1.4" fill="none"/><path d="M13 11l5 2.5-5 2.5z" fill="#00D8FF"/>' },
      { t: '表格', k: 'table', ic: '<rect x="4" y="6" width="20" height="15" stroke="#00D8FF" stroke-width="1.4" fill="none"/><path d="M4 11h20M4 16h20M12 6v15" stroke="#00D8FF" stroke-width="1"/>' }
    ] }
  ];

  var DEFAULT_SIZE = {
    rect: [140, 80], circle: [100, 100], pipe: [180, 44], rack: [96, 150],
    line: [340, 176], bar: [260, 160], pie: [190, 190],
    gauge: [196, 168], value: [164, 96], lamp: [76, 76],
    text: [160, 36], image: [200, 130], video: [260, 150], table: [320, 180]
  };

  function nodeBody(n) {
    var s = DEFAULT_SIZE[n.k] || [140, 80];
    var w = n.w - 4, h = n.h - 24;
    switch (n.k) {
      case 'gauge':
        return A.Charts.gauge({ value: 1.5, min: 1, max: 2 });
      case 'line':
        return A.Charts.line({ w: w, h: h, yLabels: false, labels: ['00:00', '06:00', '12:00', '18:00'],
          series: [{ color: '#00D8FF', area: true, data: A.wave(60, 24, 8, n.id) }] });
      case 'bar':
        return A.Charts.bars({ w: w, h: h, labels: false,
          data: [{ label: 'A', value: 62 }, { label: 'B', value: 48 }, { label: 'C', value: 75 }, { label: 'D', value: 33 }] });
      case 'pie':
        return A.Charts.donut({ size: Math.min(w, h), r: Math.min(w, h) / 2 - 14, sw: 16,
          data: [{ value: 55, color: '#00D8FF' }, { value: 30, color: '#00E39A' }, { value: 15, color: '#FFA24D' }] });
      case 'value':
        return '<div style="font-family:var(--num);font-size:30px;font-weight:700;color:var(--cyan)">' + A.fix(23.5, 1) + '<span style="font-size:13px;color:var(--t-label)"> ℃</span></div>';
      case 'lamp':
        return '<div style="width:100%;height:100%;display:grid;place-items:center">' +
          '<i style="width:26px;height:26px;border-radius:50%;background:var(--green);box-shadow:0 0 14px rgba(0,227,154,.7)"></i></div>';
      case 'rack':
        return '<div style="width:70%;height:88%;margin:auto;display:flex;flex-direction:column-reverse;gap:2px">' +
          Array.apply(null, Array(18)).map(function (_, i) {
            return '<i style="flex:1;background:' + (i < 12 ? 'rgba(0,216,255,.55)' : 'rgba(30,66,116,.5)') + ';border-radius:1px"></i>';
          }).join('') + '</div>';
      case 'table':
        return '<div style="width:100%;height:100%;font-size:11px;color:var(--t-soft)">' +
          A.UI.table({ cols: [{ t: '机柜', k: 'a' }, { t: '温度', k: 'b' }, { t: '状态', k: 'c' }],
            rows: [{ a: 'A-01', b: '23.1℃', c: '正常' }, { a: 'A-02', b: '24.6℃', c: '正常' }] }) + '</div>';
      case 'pipe':
        return '<svg viewBox="0 0 ' + w + ' ' + h + '" width="100%" height="100%"><path d="M4 ' + (h / 2) +
          ' H' + (w / 3) + ' V' + (h / 4) + ' H' + (w * 2 / 3) + ' V' + (h * 3 / 4) + ' H' + (w - 4) +
          '" stroke="#00D8FF" stroke-width="3" fill="none" stroke-linejoin="round"/></svg>';
      case 'text':
        return '<div style="font-size:14px;color:var(--t-soft)">文本标签</div>';
      case 'circle':
        return '<div style="width:100%;height:100%;display:grid;place-items:center"><i style="width:70%;height:80%;border-radius:50%;box-shadow:inset 0 0 0 2px rgba(0,216,255,.7)"></i></div>';
      case 'image':
      case 'video':
        return '<div style="color:var(--t-muted);font-size:12px">' + (n.k === 'image' ? '图片占位' : '视频窗口') + '</div>';
      default:
        return '<div style="width:76%;height:66%;margin:auto;box-shadow:inset 0 0 0 2px rgba(0,216,255,.6);border-radius:3px"></div>';
    }
  }

  function seed() {
    nodes = [
      { id: 1, k: 'gauge', t: 'PUE 仪表', x: 36, y: 32 },
      { id: 2, k: 'line', t: '功率趋势', x: 258, y: 32 },
      { id: 3, k: 'rack', t: '机柜 A-01', x: 640, y: 32 },
      { id: 4, k: 'value', t: '冷通道温度', x: 36, y: 232 },
      { id: 5, k: 'lamp', t: '运行状态', x: 224, y: 246 },
      { id: 6, k: 'pie', t: '容量占比', x: 336, y: 220 },
      { id: 7, k: 'table', t: '设备清单', x: 570, y: 230 }
    ];
    nodes.forEach(function (n) { var s = DEFAULT_SIZE[n.k]; n.w = s[0]; n.h = s[1]; });
    seq = 7;
  }
  seed();

  /* ---------------- 视图 ---------------- */
  function viewFusion(root) {
    return '<div class="editor">' +
      '<aside class="editor-col ed-left">' +
        '<div class="side-title">页面结构</div>' +
        '<div class="side-scroll tree">' +
          pageTree() +
        '</div>' +
        '<div class="side-title" style="margin-top:6px">变量</div>' +
        '<div class="side-scroll tree-var">' + varList() + '</div>' +
      '</aside>' +

      '<section class="editor-col ed-center">' +
        '<div class="card" style="flex:0 0 auto;padding:10px 14px">' +
          '<div class="card-body row ac" style="padding:0;gap:8px">' +
            '<span class="field-label">画布：2F 平面图</span>' +
            '<button class="btn is-ghost is-sm" data-ed="clear">清空</button>' +
            '<button class="btn is-ghost is-sm" data-ed="reset">还原</button>' +
            '<button class="btn is-ghost is-sm" data-ed="save">保存</button>' +
            '<span class="muted" style="margin-left:auto;font-size:12px">点击图元库添加 · 拖动摆放 · 点击选中</span>' +
          '</div>' +
        '</div>' +
        '<div class="canvas" id="canvas">' +
          (nodes.length ? '' : '<div class="canvas-hint">从右侧图元库选择组件开始组态</div>') +
          nodes.map(function (n) {
            return '<div class="canvas-node" data-n="' + n.id + '" style="left:' + n.x + 'px;top:' + n.y +
              'px;width:' + n.w + 'px;height:' + n.h + 'px">' +
              '<div class="cn-h">' + A.esc(n.t) + '</div><div class="cn-b">' + nodeBody(n) + '</div></div>';
          }).join('') +
        '</div>' +
      '</section>' +

      '<aside class="editor-col ed-right">' +
        '<div class="side-title">图元库</div>' +
        '<div class="side-scroll">' +
          LIB.map(function (g) {
            return '<div style="margin-bottom:14px"><div style="font-size:12px;color:var(--t-muted);margin-bottom:8px">' + g.g + '</div>' +
              '<div class="widget-lib">' + g.items.map(function (it) {
                return '<div class="widget" data-w="' + it.k + '" data-wt="' + A.esc(it.t) + '" draggable="true">' +
                  '<svg width="28" height="26" viewBox="0 0 28 26" fill="none">' + it.ic + '</svg>' + it.t + '</div>';
              }).join('') + '</div></div>';
          }).join('') +
        '</div>' +
        '<div class="side-title" style="margin-top:4px">属性</div>' +
        '<div id="propPane" style="flex:0 0 auto;min-height:120px;font-size:12px;color:var(--t-muted)">' +
          '未选中图元</div>' +
      '</aside>' +
    '</div>';

    function pageTree() {
      var t = [
        { n: '数据中心', c: '1' },
        { n: '1F 配电区', c: '2', i: 1 },
        { n: '2F 主机房', c: '8', i: 1 },
        { n: '205 平面图', c: '', i: 2, act: 1 },
        { n: '211 电池室', c: '', i: 2 },
        { n: '3F 屏蔽区', c: '3', i: 1 },
        { n: '4F 机房区', c: '4', i: 1 }
      ];
      return t.map(function (x) {
        return '<div class="tree-node' + (x.act ? ' is-active' : '') + '" style="padding-left:' + (6 + x.i * 14) + 'px">' +
          '<span class="nm">' + A.esc(x.n) + '</span>' + (x.c ? '<span class="ct">' + x.c + '</span>' : '') + '</div>';
      }).join('');
    }

    function varList() {
      var vs = [
        ['PUE_Value', '1.50'], ['UPS_OutPower', '1158.97'], ['IT_TotalPower', '1158.97'],
        ['Room205_Temp', '23.5'], ['Room205_Humi', '48.2'], ['Cab_A01_Load', '62'],
        ['AC_01_Return', '24.8'], ['Batt_01_Volt', '13.44'], ['Leak_205_01', '0']
      ];
      return vs.map(function (v) {
        return '<div class="vn"><i class="dot" style="width:6px;height:6px"></i><span class="nm">' + v[0] + '</span>' +
          '<span class="vt">' + v[1] + '</span></div>';
      }).join('');
    }
  }

  function viewTable(cols, rows, note) {
    return '<div class="card grow">' +
      '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(note) + '</h2>' +
        '<div class="card-extra"><button class="btn is-primary is-sm" data-act="new">+ 新增</button></div></div>' +
      '<div class="card-body">' + A.UI.table({ cols: cols, rows: rows }) + '</div></div>';
  }

  A.register('config', {
    render: function (root) {
      if (A.subArg && SUBS.some(function (s) { return s.id === A.subArg; })) state.sub = A.subArg;
      root.innerHTML =
        '<div class="page-head">' +
          '<h1 class="page-title">工程配置</h1>' +
          '<span class="page-crumb">' + A.crumb('组态与策略') +
        '</div>' +
        '<div class="subtabs">' + SUBS.map(function (s) {
          return '<button class="subtab' + (s.id === state.sub ? ' is-active' : '') + '" data-sub="' + s.id + '">' + s.name + '</button>';
        }).join('') + '</div>' +
        '<div id="cfgBody" class="grow" style="display:flex;flex-direction:column;min-height:0"></div>';
      this.paint(root);
    },

    paint: function (root) {
      var body = A.q('#cfgBody', root), html = '';

      if (state.sub === 'fusion') {
        html = viewFusion(root);
      } else if (state.sub === 'proj') {
        html = viewTable([
          { t: '工程名称', k: 'n', w: 220 }, { t: '所属园区', k: 'p', w: 180 },
          { t: '通道数', k: 'ch', align: 'r', w: 90 }, { t: '设备数', k: 'dv', align: 'r', w: 90 },
          { t: '变量数', k: 'vr', align: 'r', w: 100 }, { t: '画面数', k: 'sc', align: 'r', w: 90 },
          { t: '运行状态', k: 's', w: 110, html: function () { return A.UI.tag('运行中', 'green'); } },
          { t: '最近发布', k: 't', w: 170, mono: true }
        ], window.DCIM.rooms.slice(0, 10).map(function (r, i) {
          return { n: r.name + ' 监控工程', p: window.DCIM.meta.project, ch: 4, dv: r.dev,
                   vr: r.dev * 6, sc: 6 + (i % 4), t: '2026-10-18 09:1' + (i % 9) };
        }), '工程列表');
      } else if (state.sub === 'driver') {
        html = viewTable([
          { t: '驱动名称', k: 'n', w: 200 }, { t: '协议类型', k: 'p', w: 140 },
          { t: '通道', k: 'c', w: 100 }, { t: '设备数', k: 'd', align: 'r', w: 100 },
          { t: '采集周期', k: 'i', align: 'r', w: 110 }, { t: '通讯状态', k: 's', w: 120, html: function (r) {
            return r.s === '正常' ? A.UI.tag('正常', 'green') : A.UI.tag('中断', 'red');
          } }, { t: '最后通讯', k: 't', w: 170, mono: true }
        ], [
          { n: 'ModbusTCP_205', p: 'Modbus TCP', c: 'CH01', d: 128, i: '1000ms', s: '正常', t: '2026-10-19 23:58:12' },
          { n: 'SNMP_网络', p: 'SNMP v2c', c: 'CH02', d: 652, i: '5000ms', s: '正常', t: '2026-10-19 23:59:40' },
          { n: 'BACnet_空调', p: 'BACnet/IP', c: 'CH03', d: 71, i: '2000ms', s: '正常', t: '2026-10-19 23:59:31' },
          { n: 'UPS_智能协议', p: '厂商私有', c: 'CH04', d: 21, i: '1000ms', s: '正常', t: '2026-10-19 23:59:55' },
          { n: 'BatteryMon_211', p: 'Modbus RTU', c: 'CH05', d: 30, i: '3000ms', s: '中断', t: '2026-10-19 21:14:08' },
          { n: 'LeakRS485', p: 'Modbus RTU', c: 'CH06', d: 46, i: '2000ms', s: '正常', t: '2026-10-19 23:59:02' }
        ], '驱动管理');
      } else if (state.sub === 'alarm') {
        html = viewTable([
          { t: '规则名称', k: 'n', w: 220 }, { t: '适用对象', k: 'o', w: 160 },
          { t: '等级', k: 'l', w: 90, html: function (r) { return '<span class="lv-' + A.esc(r.l) + '">' + A.esc(r.l) + '</span>'; } },
          { t: '触发条件', k: 'c', w: 260 }, { t: '阈值', k: 'v', w: 100 },
          { t: '延时', k: 'd', align: 'r', w: 90 }, { t: '通知方式', k: 'm', w: 140 },
          { t: '启用', k: 'e', w: 90, html: function (r) { return r.e ? A.UI.tag('已启用', 'green') : A.UI.tag('停用', 'grey'); } }
        ], [
          { n: '机房高温告警', o: '温湿度测点', l: '重要', c: '温度 > 上限', v: '27.0℃', d: '60s', m: '短信+平台', e: 1 },
          { n: '机房低温告警', o: '温湿度测点', l: '次要', c: '温度 < 下限', v: '18.0℃', d: '60s', m: '平台', e: 1 },
          { n: '湿度越限', o: '温湿度测点', l: '次要', c: '湿度超出 40~70%', v: '40/70%', d: '120s', m: '平台', e: 1 },
          { n: 'UPS 过载', o: 'UPS 回路', l: '紧急', c: '负载率 > 阈值', v: '85%', d: '10s', m: '短信+电话', e: 1 },
          { n: '蓄电池电压低', o: '蓄电池单体', l: '严重', c: '单体电压 < 阈值', v: '12.0V', d: '30s', m: '短信+平台', e: 1 },
          { n: '漏水检测', o: '漏水绳', l: '紧急', c: '检测到漏水', v: '1', d: '0s', m: '短信+电话', e: 1 },
          { n: '氢气浓度高', o: '氢气传感器', l: '严重', c: '浓度 > 阈值', v: '25%LEL', d: '5s', m: '短信+电话', e: 1 },
          { n: '设备通信中断', o: '全部设备', l: '重要', c: '连续 3 次采集失败', v: '3', d: '0s', m: '平台', e: 1 }
        ], '告警规则');
      } else if (state.sub === 'storage') {
        html = viewTable([
          { t: '存储组', k: 'n', w: 180 }, { t: '数据类型', k: 't', w: 120 },
          { t: '点位数量', k: 'p', align: 'r', w: 110 }, { t: '存储周期', k: 'c', w: 110 },
          { t: '压缩', k: 'z', w: 90 }, { t: '保留时长', k: 'k', w: 110 },
          { t: '已用空间', k: 'u', align: 'r', w: 120 }, { t: '状态', k: 's', w: 100, html: function () { return A.UI.tag('正常', 'green'); } }
        ], [
          { n: '实时库_RTDB', t: '实时值', p: 12846, c: '1s', z: '否', k: '30天', u: '2.4 GB' },
          { n: '历史库_01', t: '模拟量', p: 9632, c: '5min', z: '是', k: '3年', u: '186 GB' },
          { n: '历史库_02', t: '状态量', p: 3214, c: '变位存', z: '是', k: '3年', u: '42 GB' },
          { n: '告警库', t: '事件', p: 12108, c: '事件触发', z: '否', k: '5年', u: '3.8 GB' },
          { n: '报表库', t: '聚合', p: 428, c: '日/月', z: '是', k: '10年', u: '1.1 GB' }
        ], '存储策略');
      } else if (state.sub === 'link') {
        html = viewTable([
          { t: '策略名称', k: 'n', w: 220 }, { t: '触发条件', k: 'c', w: 300 },
          { t: '执行动作', k: 'a', w: 300 }, { t: '优先级', k: 'p', align: 'r', w: 90 },
          { t: '执行次数', k: 'x', align: 'r', w: 100 }, { t: '启用', k: 'e', w: 90, html: function (r) { return r.e ? A.UI.tag('已启用', 'green') : A.UI.tag('停用', 'grey'); } }
        ], [
          { n: '高温启动备用空调', c: '机房温度 > 27℃ 持续 120s', a: '启动备用空调 + 推送运维', p: 1, x: 128, e: 1 },
          { n: 'UPS 过载切旁路', c: 'UPS 负载率 > 90%', a: '切换至旁路并短信通知', p: 0, x: 6, e: 1 },
          { n: '漏水关闭加湿', c: '检测到漏水', a: '关闭加湿器 + 关闭水阀', p: 0, x: 2, e: 1 },
          { n: '氢气超限排风', c: '氢气浓度 > 25%LEL', a: '启动事故排风机', p: 0, x: 1, e: 1 },
          { n: '夜间照明联动', c: '22:00 - 06:00 且无人', a: '关闭机房照明至 30%', p: 9, x: 412, e: 1 },
          { n: '消防切电', c: '消防主机报警', a: '切断非关键回路供电', p: 0, x: 0, e: 0 }
        ], '联动策略');
      } else {
        html = viewTable([
          { t: '函数名称', k: 'n', w: 200 }, { t: '类型', k: 't', w: 100 },
          { t: '参数', k: 'p', w: 220 }, { t: '返回值', k: 'r', w: 100 },
          { t: '用途', k: 'u', w: 300 }, { t: '调用次数', k: 'c', align: 'r', w: 110 }
        ], [
          { n: 'CalcPUE', t: '计算', p: 'totalPower, itPower', r: 'Number', u: '计算 PUE = 总功率 / IT 功率', c: 86420 },
          { n: 'AvgValue', t: '统计', p: 'tag, period', r: 'Number', u: '计算点位周期平均值', c: 152308 },
          { n: 'MaxValue', t: '统计', p: 'tag, period', r: 'Number', u: '计算点位周期最大值', c: 98412 },
          { n: 'AlarmFilter', t: '逻辑', p: 'level, room', r: 'Array', u: '按等级与机房过滤告警', c: 23184 },
          { n: 'PowerSum', t: '聚合', p: 'room', r: 'Number', u: '按机房汇总额定功率', c: 41220 },
          { n: 'TimeDiff', t: '工具', p: 't1, t2', r: 'String', u: '计算两个时间点的时长', c: 120984 }
        ], '函数管理');
      }

      body.innerHTML = html;
      this.bind(root);
    },

    bind: function (root) {
      var self = this;

      A.qa('.subtab', root).forEach(function (b) {
        b.addEventListener('click', function () {
          state.sub = b.getAttribute('data-sub');
          A.qa('.subtab', root).forEach(function (x) { x.classList.toggle('is-active', x === b); });
          self.paint(root);
        });
      });

      if (state.sub !== 'fusion') {
        A.qa('[data-act="new"]', root).forEach(function (b) {
          b.addEventListener('click', function () { A.toast('演示环境：新增配置已提交（未持久化）'); });
        });
        return;
      }

      var canvas = A.q('#canvas', root);
      var propPane = A.q('#propPane', root);

      /* 图元库 → 画布 */
      A.qa('.widget', root).forEach(function (w) {
        w.addEventListener('click', function () {
          var k = w.getAttribute('data-w'), t = w.getAttribute('data-wt');
          var s = DEFAULT_SIZE[k] || [140, 80];
          seq += 1;
          nodes.push({ id: seq, k: k, t: t, x: 60 + (nodes.length % 4) * 40, y: 60 + (nodes.length % 5) * 30, w: s[0], h: s[1] });
          self.paint(root);
          A.toast('已添加图元：' + t);
        });
        w.addEventListener('dragstart', function (e) {
          e.dataTransfer.setData('text/plain', w.getAttribute('data-w') + '|' + w.getAttribute('data-wt'));
        });
      });

      canvas.addEventListener('dragover', function (e) { e.preventDefault(); });
      canvas.addEventListener('drop', function (e) {
        e.preventDefault();
        var d = (e.dataTransfer.getData('text/plain') || '').split('|');
        if (!d[0]) return;
        var rect = canvas.getBoundingClientRect();
        var scale = rect.width / canvas.offsetWidth || 1;
        var s = DEFAULT_SIZE[d[0]] || [140, 80];
        seq += 1;
        nodes.push({
          id: seq, k: d[0], t: d[1] || d[0],
          x: Math.max(0, Math.round((e.clientX - rect.left) / scale)),
          y: Math.max(0, Math.round((e.clientY - rect.top) / scale)),
          w: s[0], h: s[1]
        });
        self.paint(root);
      });

      /* 选中 + 拖动 */
      var drag = null;
      A.qa('.canvas-node', root).forEach(function (n) {
        n.addEventListener('mousedown', function (e) {
          var id = Number(n.getAttribute('data-n'));
          var nd = nodes.filter(function (x) { return x.id === id; })[0];
          if (!nd) return;
          A.qa('.canvas-node', root).forEach(function (x) { x.classList.remove('is-sel'); });
          n.classList.add('is-sel');
          var sc = A.q('#screen');
          var m = (sc.style.transform.match(/scale\(([\d.]+)\)/) || [0, 1]);
          var k = Number(m[1]) || 1;
          drag = { nd: nd, ox: e.clientX / k - nd.x, oy: e.clientY / k - nd.y };
          propPane.innerHTML = A.UI.kv([
            ['图元类型', nd.t], ['标识', 'WGT_' + nd.k + '_' + nd.id],
            ['位置', 'X ' + nd.x + ' / Y ' + nd.y], ['尺寸', nd.w + ' × ' + nd.h],
            ['数据绑定', '未绑定变量'],
            ['提示', '点击右侧变量可绑定到该图元']
          ]);
          e.preventDefault();
        });
      });

      function onMove(e) {
        if (!drag) return;
        var sc = A.q('#screen');
        var m = (sc.style.transform.match(/scale\(([\d.]+)\)/) || [0, 1]);
        var k = Number(m[1]) || 1;
        var cv = A.q('#canvas');
        drag.nd.x = Math.max(0, Math.min(cv.offsetWidth - drag.nd.w, Math.round(e.clientX / k - drag.ox)));
        drag.nd.y = Math.max(0, Math.min(cv.offsetHeight - drag.nd.h, Math.round(e.clientY / k - drag.oy)));
        var el = A.q('.canvas-node.is-sel');
        if (el) { el.style.left = drag.nd.x + 'px'; el.style.top = drag.nd.y + 'px'; }
      }
      function onUp() {
        if (drag) {
          var el = A.q('.canvas-node.is-sel');
          if (el) A.q('.cn-h', el) && null;
          drag = null;
        }
      }
      document.addEventListener('mousemove', onMove);
      document.addEventListener('mouseup', onUp);
      this._off = function () {
        document.removeEventListener('mousemove', onMove);
        document.removeEventListener('mouseup', onUp);
      };

      /* 画布工具栏 */
      A.qa('[data-ed]', root).forEach(function (b) {
        b.addEventListener('click', function () {
          var d = b.getAttribute('data-ed');
          if (d === 'clear') { nodes = []; self.paint(root); A.toast('画布已清空'); }
          else if (d === 'reset') { seed(); self.paint(root); A.toast('已还原示例组态'); }
          else A.toast('组态已保存（演示）');
        });
      });

      /* 变量绑定 */
      A.qa('.tree-var .vn', root).forEach(function (v) {
        v.addEventListener('click', function () {
          var el = A.q('.canvas-node.is-sel');
          if (!el) { A.toast('请先在画布中选中一个图元'); return; }
          var nm = A.q('.nm', v).textContent;
          var id = Number(el.getAttribute('data-n'));
          var nd = nodes.filter(function (x) { return x.id === id; })[0];
          if (nd) { nd.t = nd.t + ' · ' + nm; self.paint(root); A.toast('已绑定变量 ' + nm); }
        });
      });
    },

    unmount: function () { if (this._off) { this._off(); this._off = null; } }
  });
})(window.App);

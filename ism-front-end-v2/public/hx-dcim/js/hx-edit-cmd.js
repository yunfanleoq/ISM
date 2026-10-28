/* 一次图编辑：智能折线几何，以及「加竖刀」这类指令。 */
(function (w) {
  'use strict';

  function bendOf(it) {
    if (it.bend === 'hv' || it.bend === 'vh') return it.bend;
    return Math.abs(+it.x2 - +it.x1) >= Math.abs(+it.y2 - +it.y1) ? 'hv' : 'vh';
  }

  function orthoPts(it) {
    var x1 = +it.x1;
    var y1 = +it.y1;
    var x2 = +it.x2;
    var y2 = +it.y2;
    if (Math.abs(x2 - x1) < 1.5 || Math.abs(y2 - y1) < 1.5) return [x1, y1, x2, y2];
    if (bendOf(it) === 'hv') {
      var xm = it.xm != null ? +it.xm : (x1 + x2) / 2;
      return [x1, y1, xm, y1, xm, y2, x2, y2];
    }
    var ym = it.ym != null ? +it.ym : (y1 + y2) / 2;
    return [x1, y1, x1, ym, x2, ym, x2, y2];
  }

  function orthoD(pts) {
    var d = 'M' + pts[0] + ' ' + pts[1];
    var i;
    for (i = 2; i < pts.length; i += 2) d += ' L' + pts[i] + ' ' + pts[i + 1];
    return d;
  }

  function orthoMid(it) {
    var p = orthoPts(it);
    if (p.length < 6) {
      return { x: (p[0] + p[p.length - 2]) / 2, y: (p[1] + p[p.length - 1]) / 2 };
    }
    return { x: (p[2] + p[4]) / 2, y: (p[3] + p[5]) / 2 };
  }

  function parse(raw) {
    var s = String(raw || '').replace(/\s+/g, ' ').trim();
    if (!s) return { err: '请输入指令' };
    if (/^(删除|删掉|去掉)$/.test(s)) return { act: 'del' };
    var kind = '';
    if (/智能折线|折线/.test(s)) kind = 'ortho';
    else if (/母线/.test(s)) kind = 'bus';
    else if (/横线/.test(s)) kind = 'hline';
    else if (/竖线/.test(s)) kind = 'vline';
    else if (/自定义|新组件/.test(s)) kind = 'custom';
    else if (/补偿/.test(s)) kind = 'cap';
    else if (/接地/.test(s)) kind = 'gnd';
    else if (/柴油|发电机|油机/.test(s)) kind = 'gen';
    else if (/变压器/.test(s)) kind = 'xfmr';
    else if (/逆变器/.test(s)) kind = 'inv';
    else if (/静态开关|STS/i.test(s)) kind = 'ssw';
    else if (/UPS|ups/.test(s)) kind = 'ups';
    else if (/竖刀|刀闸|闸刀|开关|刀/.test(s)) kind = 'knife';
    else if (/文字|文本|标注/.test(s)) kind = 'text';
    if (!kind) {
      return { err: '可以说：加竖刀、加变压器、加补偿柜、加UPS、加逆变器、加静态开关、加自定义 接地、加文字 备用、加横线、加竖线、加智能折线' };
    }
    var count = 1;
    var cm = s.match(/(\d+)\s*(把|个|条)/);
    if (cm) count = Math.min(12, parseInt(cm[1], 10) || 1);
    var text = '文字';
    if (kind === 'text') {
      var tm = s.match(/(?:文字|文本|标注)\s+(.+)$/);
      if (tm) text = tm[1].trim();
    }
    if (kind === 'custom') {
      var cmn = s.match(/(?:自定义|新组件)\s+(.+)$/);
      if (cmn) text = cmn[1].trim();
    }
    return { act: 'add', kind: kind, count: count, text: text };
  }

  w.HXEditCmd = {
    bendOf: bendOf,
    orthoPts: orthoPts,
    orthoD: orthoD,
    orthoMid: orthoMid,
    parse: parse
  };
})(window);

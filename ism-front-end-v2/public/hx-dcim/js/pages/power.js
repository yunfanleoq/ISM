/* 航信配电一次图。不用演示单线，不画电容补偿柜。 */
(function (A) {
  'use strict';

  A.register('power', {
    render: function (root) {
      root.innerHTML = '<div id="hxPower" class="grow" style="min-height:0"></div>';
      var host = document.getElementById('hxPower');
      host.innerHTML =
        '<div class="hx-pt" id="hxPowerWait">' +
          '正在加载配电专业画面…<br>' +
          '<span class="muted">测点索引约 25MB，弱网或首次打开可能需要 1～2 分钟；超时后会先展示一次图。</span>' +
        '</div>';
      /* drill：左侧整棵架构树，配电室→一次图，机房/列→列头柜，顶栏面包屑 */
      var run = function () {
        var wait = document.getElementById('hxPowerWait');
        if (wait) wait.remove();
        if (window.HXScreen) window.HXScreen.mount(host, 'drill');
      };
      if (window.HXScreen && window.HXScreen.ready) {
        window.HXScreen.ready.then(run).catch(function () { run(); });
      } else run();
    },
    unmount: function () { if (window.HXScreen) window.HXScreen.leave(); }
  });
})(window.App);

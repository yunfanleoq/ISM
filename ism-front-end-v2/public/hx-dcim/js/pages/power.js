/* 航信配电一次图。不用演示单线，不画电容补偿柜。 */
(function (A) {
  'use strict';

  A.register('power', {
    render: function (root) {
      root.innerHTML = '<div id="hxPower" class="grow" style="min-height:0"></div>';
      var host = document.getElementById('hxPower');
      /* drill：左侧整棵架构树，配电室→一次图，机房/列→列头柜，顶栏面包屑 */
      var run = function () { if (window.HXScreen) window.HXScreen.mount(host, 'drill'); };
      if (window.HXScreen && window.HXScreen.ready) window.HXScreen.ready.then(run);
      else run();
    },
    unmount: function () { if (window.HXScreen) window.HXScreen.leave(); }
  });
})(window.App);

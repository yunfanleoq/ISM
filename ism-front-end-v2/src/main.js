import Vue from 'vue'
import App from './App.vue'
import {initRouter} from './router'
import './theme/index.less'
import Antd from 'ant-design-vue'
import Viser from 'viser-vue'
import store from './store'
import 'animate.css/source/animate.css'
import Plugins from '@/plugins'
import {initI18n} from '@/utils/i18n'
import bootstrap from '@/bootstrap'
import DisableDevtool from 'disable-devtool';
import 'moment/locale/zh-cn'
import dataV from '@jiaminghi/data-view'
import './font/font.css'
import VueResource from 'vue-resource'
import { Icon } from 'ant-design-vue';
import  "@/font/iconfont.css"
const IconFont = Icon.createFromIconfontCN({
  scriptUrl: 'static/js/iconfont/iconfont.js?'+Math.round(Math.random()*100),
});

import BaiduMap from 'vue-baidu-map'

 
// DisableDevtool({
//   url: 'http://www.ismctl.com',
//   timeOutUrl: 'http://www.ismctl.com',
//   disableMenu: true
// });

Vue.use(BaiduMap, {
  ak: '04BeCXvIjhMuZCmsrqVsBhXW1G8p3DKp'
})
// main.js
const router = initRouter(store.state.setting.asyncRoutes)
const i18n = initI18n('CN', 'US')

Vue.use(Antd)
Vue.config.productionTip = false
Vue.use(Viser)
Vue.use(Plugins)
Vue.use(dataV)
Vue.use(VueResource)
Vue.component('icon-font',IconFont)

Vue.prototype.$EventBus = new Vue()

bootstrap({router, store, i18n, message: Vue.prototype.$message})

new Vue({
  router,
  store,
  i18n,
  render: h => h(App),
}).$mount('#app')

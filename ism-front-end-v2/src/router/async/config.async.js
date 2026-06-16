import routerMap from './router.map'
import {parseRoutes} from '@/utils/routerUtil'

// 异步路由配置
const routesConfig = [
  'login',
  'loginPhone',
  'auth',
  'e404',
  'e403',
  'Exp500',
  'Project',
  'UserDisplayList',
  'SimulatorMonitor',
  'DisPlayEditor',
  'DisPlay3DEditor',
  'DisPlay3DRunApp',
  'AmisEditor',
  'AppRun',
  'DisPlayRunApp',
  'ShareApp',
  'AppLogin',
  'root',
]

const options = {
  routes: parseRoutes(routesConfig, routerMap)
}

export default options

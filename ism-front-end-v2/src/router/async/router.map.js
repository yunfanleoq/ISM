// 视图组件
const view = {
  tabs: () => import('@/layouts/tabs'),
  blank: () => import('@/layouts/BlankView'),
  page: () => import('@/layouts/PageView')
}
import TabsView from '@/layouts/tabs/TabsView'
// import PageView from '@/layouts/PageView'
import BlankView from '@/layouts/BlankView'

// 路由组件注册
const routerMap = {
  login: {
    path: '/login',
    component: () => import('@/pages/login')
  },
  root: {
    path: '/',
    redirect: '/dashboard',
    component: () => import('@/layouts/AdminLayout')
  },
  auth:{
    path: '/auth',
    name: '登录授权',
    component: () => import('@/pages/login/Auth')
  },
  loginPhone:{
    path: '/loginPhone',
    name: '手机登录',
    component: () => import('@/pages/login/LoginPhone')
  },
  e404:{
    path: '*',
    name: '404',
    component: () => import('@/pages/exception/404'),
  },
  e403:{
    path: '/403',
    name: '403',
    component: () => import('@/pages/exception/403'),
  },
  Exp500:{
    path: '500',
    name: 'Exp500',
    component: () => import('@/pages/exception/500')
  },
  Project:{
    path: '/Project',
    name: '项目',
    meta: {
      icon:"database",
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/project/index')
  },
  UserDisplayList:{
    path: '/UserDisplayList/:uuid',
    name: '普通用户应用列表',
    meta: {
      icon:"database",
    },
    component: () => import('@/pages/project/userDisplayList')
  },
  DisPlayEditor:{
    path: '/DisPlayEditor/:uid',
    name: '组态编辑',
    meta: {
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/ISMDisPlay/ISMDisPlayEditor')
  },
  AmisEditor:{
    path: '/AmisEditor/:uid',
    name: '页面编辑',
    meta: {
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/amisEditor/AmisEditor')
  },
  DisPlay3DEditor:{
    path: '/DisPlay3DEditor/:uid',
    name: '3D场景编辑',
    meta: {
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/ISM3DEditor/ISM3DEditor')
  },
  DisPlay3DRunApp:{
    path: '/DisPlay3DRunApp/:uid',
    name: '3D场景预览',
    meta: {
    },
    component: () => import('@/pages/ISM3DEditor/components/Preview/ScenePreview')
  },
  AppRun:{
    path: '/AppRun/:uid',
    name: '组态预览',
    meta: {

    },
    component: () => import('@/pages/ISMDisPlay/pageView')
  },
  ShareApp:{
    path: '/ShareApp/:uid/:token',
    name: '组态分享',
    meta: {

    },
    component: () => import('@/pages/ISMDisPlay/ShareApp')
  },
  AppLogin:{
    path: '/AppLogin/:uid',
    name: '组态登录',
    meta: {

    },
    component: () => import('@/pages/ISMDisPlay/pageViewLogin')
  },
  root: {
    path: '/',
    name: '首页',
    redirect: '/login',
    component: TabsView
  },
  dashboard:{
    path: 'dashboard',
    name: '数据中台',
    meta: {
      icon: 'dashboard',
    },
    component: () => import('@/pages/project/dashboard')
  },
  SCADAMonitor:{
    path: '/SCADAMonitor',
    name: '电力监控大屏',
    meta: {
      icon: 'dashboard',
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/SCADAMonitor/index')
  },
  DataWarehouse:{
    path: 'DataWarehouse',
    name: '数据仓库',
    meta: {
      icon: 'container',
    },
    component: () => import('@/pages/deviceLibrary/monitor')
  },
  DeviceLibraryConfig:{
    path: 'DeviceLibraryConfig',
    name: '设备管理',
    meta: {
      icon: 'hdd',
      authority: {
        role: ['Admin','Operator','User']
      }
    },
    component: () => import('@/pages/deviceLibrary/deviceConfig')
  },
  DeviceModel:{
    path: 'DeviceModel',
    name: '数据模型',
    meta: {
      icon: 'dropbox',
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: BlankView,
  },
  SnmpModel:{
  path: 'SnmpModel',
      name: 'SNMP设备',
    meta: {

},
  component: () => import('@/pages/dataModel/snmp/SnmpModel'),
},
  SnmpAdd:{
  path: 'SnmpAdd',
      name: '添加SNMP数据模型',
    meta: {

  invisible:true
},
  component: () => import('@/pages/dataModel/snmp/SnmpModelAdd'),
},
  SnmpDetail:{
  path: 'SnmpDetail/:uid/:tab',
      name: 'SNMP数据模型详情',
    meta: {
  invisible:true
},
  component: () => import('@/pages/dataModel/snmp/SnmpModelDetail'),
},

  ModbusModel:{
  path: 'ModbusModel',
      name: 'Modbus设备',
    meta: {

},
  component: () => import('@/pages/dataModel/modbus/ModbusModel'),
},
  ModbusAdd:{
  path: 'ModbusAdd',
      name: '添加modbus数据模型',
    meta: {
  invisible:true
},
  component: () => import('@/pages/dataModel/modbus/ModbusModelAdd'),
},
  ModbusRegister:{
  path: 'ModbusRegister/:uid',
      name: 'Modbus寄存器配置',
    meta: {
  invisible:true
},
  component: () => import('@/pages/dataModel/modbus/ModbusModelRegister'),
},
  ModbusDetail:{
  path: 'ModbusDetail/:uid',
      name: 'Modbus模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/modbus/ModbusModelDetail'),
},


  DLT645Model:{
  path: 'DLT645Model',
      name: 'DLT645设备',
    meta: {

},
  component: () => import('@/pages/dataModel/DLT645/DLT645Model'),
},
  DLT645ModelAdd:{
  path: 'DLT645ModelAdd',
      name: '添加DLT645数据模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/DLT645/DLT645ModelAdd'),
},
  DLT645ModelData:{
  path: 'DLT645ModelData/:uid',
      name: 'DLT645数据配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/DLT645/DLT645Data'),
},
  DLT645ModelDetail:{
  path: 'DLT645ModelDetail/:uid',
      name: 'DLT645模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/DLT645/DLT645ModelDetail'),
},

  CJT188Model:{
    path: 'CJT188Model',
    name: 'CJT188仪表',
    meta: {
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/CJT188/CJT188Model'),
  },
  CJT188ModelAdd:{
    path: 'CJT188ModelAdd',
    name: '添加CJT188数据模型',
    meta: {
      invisible:true,
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/CJT188/CJT188ModelAdd'),
  },
  CJT188ModelData:{
    path: 'CJT188ModelData/:uid',
    name: 'CJT188数据配置',
    meta: {
      invisible:true,
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/CJT188/CJT188Data'),
  },
  CJT188ModelDetail:{
    path: 'CJT188ModelDetail/:uid',
    name: 'CJT188模型详情',
    meta: {
      invisible:true,
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/CJT188/CJT188ModelDetail'),
  },

  OPCUAModel:{
  path: 'OPCUAModel',
      name: 'OPCUA设备',
    meta: {

},
  component: () => import('@/pages/dataModel/opcua/OpcuaModel'),
},
  OpcuaAdd:{
  path: 'OpcuaAdd',
      name: '添加OPCUA数据模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/opcua/OpcuaModelAdd'),
},
  OpcuaNodeid:{
  path: 'OpcuaNodeid/:uid',
      name: 'OPCUA Nodeid配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/opcua/OpcuaNodeid'),
},
  OpcuaDetail:{
  path: 'OpcuaDetail/:uid',
      name: 'OPCUA模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/opcua/OpcuaModelDetail'),
},

  MqttModel:{
  path: 'MqttModel',
      name: 'MQTT设备',
    meta: {

},
  component: () => import('@/pages/dataModel/mqtt/mqttModel'),
},
  MqttAdd:{
  path: 'MqttAdd',
      name: '添加MQTT数据模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/mqtt/mqttModelAdd'),
},
  MqttNodeid:{
  path: 'MqttNodeid/:uid',
      name: 'MQTT配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/mqtt/mqttNodeid'),
},
  MqttDetail:{
  path: 'MqttDetail/:uid',
      name: 'MQTT模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/mqtt/mqttModelDetail'),
},
  SiemensS7Model:{
  path: 'SiemensS7Model',
      name: '西门子S7设备',
    meta: {

},
  component: () => import('@/pages/dataModel/SimS7/SimS7Model'),
},
  SimS7Add:{
  path: 'SimS7Add',
      name: '添加S7模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/SimS7/S7ModelAdd'),
},
  S7DataList:{
  path: 'S7DataList/:uid',
      name: 'S7数据配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/SimS7/S7ModelDataList'),
},
  SimS7Detail:{
  path: 'SimS7Detail/:uid',
      name: 'S7模型详细',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/SimS7/S7ModelDetail'),
},

  IEC104Model:{
  path: 'IEC104Model',
      name: 'IEC104设备',
    meta: {

},
  component: () => import('@/pages/dataModel/IEC104/IEC104Model'),
},
  IEC104ModelAdd:{
  path: 'IEC104ModelAdd',
      name: '添加IEC104数据模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/IEC104/IEC104ModelAdd'),
},
  IEC104ModelData:{
  path: 'IEC104ModelData/:uid',
      name: 'IEC104数据配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/IEC104/IEC104Data'),
},
  IEC104ModelDetail:{
  path: 'IEC104ModelDetail/:uid',
      name: 'IEC104模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/IEC104/IEC104ModelDetail'),
},

  IEC61850Model:{
  path: 'IEC61850Model',
      name: 'IEC61850设备',
    meta: {

},
  component: () => import('@/pages/dataModel/IEC61850/IEC61850Model'),
},
  IEC61850Add:{
  path: 'IEC61850Add',
      name: '添加IEC61850数据模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/IEC61850/IEC61850ModelAdd'),
},
  IEC61850Nodeid:{
  path: 'IEC61850Nodeid/:uid',
      name: 'IEC61850 Nodeid配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/IEC61850/IEC61850Nodeid'),
},
  IEC61850Detail:{
  path: 'IEC61850Detail/:uid',
      name: 'IEC61850模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/IEC61850/IEC61850ModelDetail'),
},

  DNP3Model:{
  path: 'DNP3Model',
      name: 'DNP3设备',
    meta: {
      authority: {
        role: ['Admin','Operator']
      }
},
  component: () => import('@/pages/dataModel/DNP3/DNP3Model'),
},
  DNP3ModelAdd:{
  path: 'DNP3ModelAdd',
      name: '添加DNP3数据模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/DNP3/DNP3ModelAdd'),
},
  DNP3ModelData:{
  path: 'DNP3ModelData/:uid',
      name: 'DNP3数据配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/DNP3/DNP3Data'),
},
  DNP3ModelDetail:{
  path: 'DNP3ModelDetail/:uid',
      name: 'DNP3模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/DNP3/DNP3ModelDetail'),
},

  BACnetModel:{
    path: 'BACnetModel',
    name: 'BACnet设备',
    meta: {
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/bacnet/bacnetModel'),
  },
  BACnetAdd:{
    path: 'BACnetAdd',
    name: '添加BACnet数据模型',
    meta: {
      invisible:true,
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/bacnet/bacnetModelAdd'),
  },
  BACnetNodeid:{
    path: 'BACnetNodeid/:uid',
    name: 'BACnet配置',
    meta: {
      invisible:true,
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/bacnet/bacnetNodeid'),
  },
  BACnetDetail:{
    path: 'BACnetDetail/:uid',
    name: 'BACnet模型详情',
    meta: {
      invisible:true,
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/dataModel/bacnet/bacnetModelDetail'),
  },

  RestFulModel:{
  path: 'RestFulModel',
      name: 'RESTFul设备',
    meta: {

},
  component: () => import('@/pages/dataModel/RESTFulData/RESTFulModel'),
},
  RestFulData:{
  path: 'RestFulData/:uid',
      name: 'RESTFul模型数据',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/RESTFulData/RESTFulDataList'),
},
  HJ212Model:{
  path: 'HJ212Model',
      name: '环保HJ212-2017',
    meta: {

},
  component: () => import('@/pages/dataModel/hj212/hj212Model'),
},
  HJ212Add:{
  path: 'HJ212Add',
      name: '添加环保数据模型',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/hj212/hj212ModelAdd'),
},
  HJ212Nodeid:{
  path: 'HJ212Nodeid/:uid',
      name: '环保协议数据配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/hj212/hj212Nodeid'),
},
  HJ212Detail:{
  path: 'HJ212Detail/:uid',
      name: '环保模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/hj212/hj212ModelDetail'),
},
  VirtualDevice:{
  path: 'VirtualDevice',
      name: '虚拟设备',
    meta: {

},
  component: () => import('@/pages/dataModel/VirtualDevice/VirtualDeviceModel'),
},
  VirtualDeviceData:{
  path: 'VirtualDeviceData/:uid',
      name: '虚拟模型数据',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/VirtualDevice/VirtualDeviceDataList'),
},
  DeviceCustomData:{
  path: 'DeviceCustomData',
      name: '自定义数据',
    meta: {

},
  component: () => import('@/pages/dataModel/CustomData/CustomDataModel'),
},
  SystemData:{
  path: 'SystemData',
      name: '系统变量',
    meta: {

},
  component: () => import('@/pages/dataModel/static/StaticModel'),
},
  StaticDataAdd:{
  path: 'StaticDataAdd',
      name: '添加静态数据模型',
    meta: {

  invisible:true
},
  component: () => import('@/pages/dataModel/static/StaticModelAdd'),
},
  StaticDataDetail:{
  path: 'StaticDataDetail/:uid',
      name: '静态数据模型详情',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/dataModel/static/StaticModelDetail'),
},
  Application:{
    path: 'Application',
    name: '应用管理',
    meta: {
      icon:"pic-right",
      authority: {
        role: ['Admin','Operator']
      }
    },

    component: () => import('@/pages/disPlayModel/displayModelList')
  },
  TemplateMarket:{
    path: 'TemplateMarket',
    name: '模板市场',
    meta: {
      icon:"appstore",
      invisible: false
    },
    component: () => import('@/pages/disPlayModel/TemplateMarket')
  },
  DigitalTwin:{
    path: 'DigitalTwin',
    name: '数字孪生',
    meta: {
      icon:"appstore",
      authority: {
        role: ['Admin','Operator']
      }
    },

    component: () => import('@/pages/ISM3DApp/ISM3DAppList')
  },
  ISM3DApp:{
    path: 'DigitalTwin',
    name: '数字孪生',
    meta: {
      icon:"appstore",
      authority: {
        role: ['Admin','Operator']
      }
    },

    component: () => import('@/pages/ISM3DApp/ISM3DAppList')
  },
  VideoManager:{
    path: 'VideoManager',
    name: '视频管理',
    meta: {
      icon:"video-camera",
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: BlankView,
  },
  videoScreen:{
  path: 'videoScreen',
      name: '分屏展示',
    meta: {
  icon:"table",
},
  component: () => import('@/pages/video/videoShowMore')
},
  videoList:{
  path: 'videoList',
      name: '视频列表',
    meta: {
  icon:"unordered-list"
},
  component: () => import('@/pages/video/videoManager')
},
  GB28281List:{
  path: 'GB28281List',
      name: '国标视频',
    meta: {
  icon:"youtube"
},
  component: () => import('@/pages/video/GB28281VideoManager')
},
  RealTimeAlarm:{
    path: 'Real-timeAlarm',
    name: '实时告警',
    meta: {
      icon:"bell"
    },
    component: () => import('@/pages/alarm/currentAlarm/currentAlarm')
},
  AlarmStrategy:{
    path: 'AlarmStrategy',
    name: '告警策略',
    meta: {
      icon: 'alert',
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: BlankView,
  },
  ModelTrigger:{
  path: 'ModelTrigger',
      name: '模型触发器',
    meta: {
  icon:"thunderbolt"
},
  component: () => import('@/pages/alarm/trigger/trigger'),
},
  AlarmRestoreMask:{
  path: 'AlarmRestoreMask',
      name: '告警恢复',
    meta: {
  icon:"undo",

},
  component: () => import('@/pages/alarm/ShieldAlarm/shieldAlarm'),
},
  TaskPlan:{
    path: 'TaskPlan',
    name: '任务计划',
    meta: {
      icon:"schedule"
    },
    component: () => import('@/pages/taskplan/taskplan'),
  },
  ISMScripts:{
    path: 'ISMScripts',
    name: '系统脚本',
    meta: {
      icon:"code",
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: () => import('@/pages/ISMScripts/scriptsList'),
  },
  DataPush:{
    path: 'DataPush',
    name: '数据接口',
    meta: {
      icon: 'radar-chart',
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: BlankView,
  },
  DataTemplete:{
  path: 'DataTemplete',
      name: '数据模版',
    meta: {
  icon: 'file-text'
},
  component: () => import('@/pages/DataPush/DataTemplete'),
},
  IEC104DataTemplete:{
  path: 'IEC104DataTemplete',
      name: 'IEC104数据模版',
    meta: {
  icon: 'file-text'
},
  component: () => import('@/pages/DataPush/iec104/IEC104DataTemplete'),
},
  IEC104TempleteData:{
  path: 'IEC104TempleteData/:uid',
      name: 'IEC104数据模版配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/DataPush/iec104/IEC104PusbData'),
},
  ModbusDataTemplete:{
  path: 'ModbusDataTemplete',
      name: 'Modbus 数据模版',
    meta: {
  icon: 'file-text'
},
  component: () => import('@/pages/DataPush/modbusTcp/ModbusTcpDataTemplete'),
},
  ModbusTempleteData:{
  path: 'ModbusTempleteData/:uid',
      name: 'Modbus 数据模版配置',
    meta: {
  invisible:true,

},
  component: () => import('@/pages/DataPush/modbusTcp/ModbusTcpPusbData'),
},
  DataInterface:{
  path: 'DataInterface',
      name: '接口管理',
    meta: {
  icon: 'hdd'
},
  component: () => import('@/pages/DataPush/DataInterface'),
},
  Reporting:{
    path: 'Reporting',
    name: '数据报表',
    meta: {
      icon: 'file-done'
    },
    component: BlankView,
  },
  AlarmHistory:{
  path: 'AlarmHistory',
      name: '告警报表',
    meta: {

},
  component: () => import('@/pages/reporting/alarmReport/alarmHistory'),
},
  DataHistory:{
  path: 'DataHistory',
      name: '历史数据',
    meta: {

},
  component: () => import('@/pages/reporting/dataHistoryReport/dataHistoryReport'),
},
  DataHistoryReport:{
    path: 'DataHistoryReport',
    name: '历史报表',
    meta: {

    },
    component: () => import('@/pages/reporting/dataHistoryReport/dataHistoryQueryReport'),
  },
  SqlReportTemplate:{
    path: 'SqlReportTemplate',
    name: 'SQL报表',
    meta: {

    },
    component: () => import('@/pages/reporting/sqlReport/SqlReportTemplate'),
  },
  DiyReport:{
  path: 'DiyReport',
      name: '自定义报表',
    meta: {

},
  component: () => import('@/pages/reporting/diyReport/diyReport'),
},
  DiyReportTemplete:{
  path: 'DiyReportTemplete',
      name: '报表模板',
    meta: {
      authority: {
        role: ['Admin','Operator']
      }
},
  component: () => import('@/pages/reporting/diyReport/diyReportTemplete'),
},
  ReportTempleteContent:{
  path: 'ReportTempleteContent/:uuid',
      name: '报表模板内容修改',
    meta: {

  invisible:true
},
  component: () => import('@/pages/reporting/diyReport/diyReportContent'),
},
  Network:{
    path: 'Network',
    name: '网络中心',
    meta: {
      icon: 'global',
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: BlankView,
  },
  SystemNetwork:{
  path: 'SystemNetwork',
      name: '网络',
    meta: {
  icon:"ie",
},
  component: () => import('@/pages/network/network')
},
  ISMNetwork:{
  path: 'ISMNetwork',
      name: '组网',
    meta: {
  icon:"slack",
},
  component: () => import('@/pages/network/ismnode')
},
  DataBase:{
    path: 'DataBase',
    name: '数据库管理',
    meta: {
      icon: 'database',
      authority: {
        role: ['Admin','Operator']
      }
    },
    component: BlankView,
  },
  DbManager:{
  path: 'DbManager',
      name: '实时数据库',
    meta: {
  icon:"slack",
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/db/DbManager')
},
  HistoryManager:{
  path: 'HistoryManager',
      name: '历史数据库',
    meta: {
  icon:"code-sandbox",
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/db/HisManager')
},
  Setting:{
    path: 'Setting',
    name: '设置中心',
    meta: {
      icon: 'setting'
    },
    component: BlankView,
  },
  Account:{
  path: 'Account',
      name: '个人设置',
    meta: {
  icon:"user",
},
  component: () => import('@/pages/account/settings/Index')
},
  UserAdd:{
  path: 'UserAdd',
      name: '用户添加',
    meta: {
  icon:"user-add",
      invisible:true,
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/account/settings/UserAdd')
},
  UserManager:{
  path: 'UserManager',
      name: '用户管理',
    meta: {
  icon:"solution",
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/account/settings/UserList')
},
  AlarmTipsSetting:{
  path: 'AlarmTipsSetting',
      name: '告警通知',
    meta: {
  icon:"sound"
},
  component: () => import('@/pages/Setting/AlarmTips')
},
  SystemParams:{
  path: 'SystemParams',
      name: '系统参数',
    meta: {
  icon:"tool",
},
  component: () => import('@/pages/Setting/SystemParams')
},
  AccessToken:{
  path: 'AccessToken',
      name: 'API令牌',
    meta: {
  icon:"safety",
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/account/ApiToken')
},
  RoleManager:{
  path: 'RoleManager',
      name: '角色管理',
    meta: {
  icon:"team",
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/account/RoleManager')
},
  Journal:{
    path: 'Journal',
    name: '系统日志',
    meta: {
      icon: 'file-text'
    },
    component: BlankView,
  },
  OperationJournal:{
  path: 'OperationJournal',
      name: '操作日志',
    meta: {

},
  component: () => import('@/pages/journal/OperationLog/OperationLog'),
},
  SimulatorMonitor:{
    path: '/SimulatorMonitor',
    name: '模拟器界面',
    meta: {
      icon: 'api',
      authority: {
        role: ['Admin', 'Operator']
      },
    },
    component: () => import('@/pages/simulator/SimulatorMonitor'),
  },
  Help:{
    path: '/Help',
    name: '帮助中心',
    meta: {
      icon: 'question-circle',
      authority: {
        role: ['Admin']
      },
    },
    component: BlankView,
  },
  AboutSystem:{
  path: 'AboutSystem',
      name: '关于',
    meta: {
  icon:"profile",
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/Setting/SystemUpgrade')
},
  SystemAuth:{
  path: 'SystemAuth',
      name: '系统授权',
    meta: {
  icon:"key",
      authority: {
    role: ['Admin']
  }
},
  component: () => import('@/pages/Setting/SystemAuth')
},
}
export default routerMap
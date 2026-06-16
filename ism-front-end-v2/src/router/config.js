import TabsView from '@/layouts/tabs/TabsView'
// import PageView from '@/layouts/PageView'
import BlankView from '@/layouts/BlankView'
// 路由配置
const options = {
  routes: [
    {
      path: '/login',
      name: '登录',
      component: () => import('@/pages/login/Login')
    },
    {
      path: '/auth',
      name: '登录授权',
      component: () => import('@/pages/login/Auth')
    },
    {
          path: '/loginPhone',
          name: '手机登录',
          component: () => import('@/pages/login/LoginPhone')
    },
    {
      path: '*',
      name: '404',
      component: () => import('@/pages/exception/404'),
    },
    {
      path: '/403',
      name: '403',
      component: () => import('@/pages/exception/403'),
    },
    {
      path: '500',
      name: 'Exp500',
      component: () => import('@/pages/exception/500')
    },
    {
      path: '/SimulatorMonitor',
      name: '模拟器界面',
      meta: {
        icon: 'api',
        authority: {
          role: ['Admin', 'Operator']
        }
      },
      component: () => import('@/pages/simulator/SimulatorMonitor')
    },
    {
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
    {
      path: '/UserDisplayList/:uuid',
      name: '普通用户应用列表',
      meta: {
        icon:"database",
      },
      component: () => import('@/pages/project/userDisplayList')
    },
    {
      path: '/Help',
      name: '帮助',
      component: () => import('@/pages/help/help')
    },
    {
      path: '/SCADAMonitor',
      name: '电力监控大屏',
      meta: {
        icon: 'dashboard',
        authority: {
          role: ['Admin', 'Operator']
        }
      },
      component: () => import('@/pages/SCADAMonitor/index')
    },
    {
      path: '/DisPlayEditor/:uid',
      name: '组态编辑',
      meta: {
        authority: {
          role: ['Admin','Operator']
        }
      },
      component: () => import('@/pages/ISMDisPlay/ISMDisPlayEditor')
    },
    {
      path: '/AmisEditor/:uid',
      name: '页面编辑',
      meta: {
        authority: {
          role: ['Admin','Operator']
        }
      },
      component: () => import('@/pages/amisEditor/AmisEditor')
    },
    {
      path: '/DisPlay3DEditor/:uid',
      name: '3D场景编辑',
      meta: {
        authority: {
          role: ['Admin','Operator']
        }
      },
      component: () => import('@/pages/ISM3DEditor/ISM3DEditor')
    },
    {
      path: '/DisPlay3DPreview/:uid',
      name: '3D场景预览',
      meta: {
      },
      component: () => import('@/pages/ISM3DEditor/components/Preview/ScenePreview')
    },
    {
      path: '/AppRun/:uid',
      name: '组态预览',
      meta: {

      },
      component: () => import('@/pages/ISMDisPlay/pageView')
    },
    {
      path: '/ISMDisPlay/DisPlayRunApp',
      name: '组态预览(兼容)',
      redirect: to => {
        const displayUUID = to.query.displayUUID || to.query.uid
        if (displayUUID) {
          const query = { ...to.query }
          delete query.displayUUID
          return { path: `/AppRun/${displayUUID}`, query }
        }
        return '/404'
      }
    },
    {
      path: '/ShareApp/:uid/:token',
      name: '组态分享',
      meta: {

      },
      component: () => import('@/pages/ISMDisPlay/ShareApp')
    },
    {
      path: '/AppLogin/:uid',
      name: '组态登录',
      meta: {

      },
      component: () => import('@/pages/ISMDisPlay/pageViewLogin')
    },
    {
      path: '/',
      name: '首页',
      component: TabsView,
      redirect: '/login',
      children: [
        {
          path: 'dashboard',
          name: '资源统计',
          meta: {
            icon: 'dashboard',
          },
          component: () => import('@/pages/project/dashboard')
        },
        {
          path: 'DataWarehouse',
          name: '数据仓库',
          meta: {
            icon: 'container',
          },
          component: () => import('@/pages/deviceLibrary/monitor')
        },
        {
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
        {
          path: 'DeviceModel',
          name: '数据模型',
          meta: {
            icon: 'dropbox',
            authority: {
              role: ['Admin','Operator']
            }
          },
          component: BlankView,
          children: [
            {
              path: 'SnmpModel',
              name: 'SNMP设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/snmp/SnmpModel'),
            },
            {
              path: 'SnmpAdd',
              name: '添加SNMP数据模型',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                },
                invisible:true
              },
              component: () => import('@/pages/dataModel/snmp/SnmpModelAdd'),
            },
            {
              path: 'SnmpDetail/:uid/:tab',
              name: 'SNMP数据模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/snmp/SnmpModelDetail'),
            },

            {
              path: 'ModbusModel',
              name: 'Modbus设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/modbus/ModbusModel'),
            },
            {
              path: 'ModbusAdd',
              name: '添加modbus数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/modbus/ModbusModelAdd'),
            },
            {
              path: 'ModbusRegister/:uid',
              name: 'Modbus寄存器配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/modbus/ModbusModelRegister'),
            },
            {
              path: 'ModbusDetail/:uid',
              name: 'Modbus模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/modbus/ModbusModelDetail'),
            },


            {
              path: 'DLT645Model',
              name: 'DLT645设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DLT645/DLT645Model'),
            },
            {
              path: 'DLT645ModelAdd',
              name: '添加DLT645数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DLT645/DLT645ModelAdd'),
            },
            {
              path: 'DLT645ModelData/:uid',
              name: 'DLT645数据配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DLT645/DLT645Data'),
            },
            {
              path: 'DLT645ModelDetail/:uid',
              name: 'DLT645模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DLT645/DLT645ModelDetail'),
            },

            {
              path: 'CJT188Model',
              name: 'CJT188仪表',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/CJT188/CJT188Model'),
            },
            {
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
            {
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
            {
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

            {
              path: 'OPCUAModel',
              name: 'OPCUA设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/opcua/OpcuaModel'),
            },
            {
              path: 'OpcuaAdd',
              name: '添加OPCUA数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/opcua/OpcuaModelAdd'),
            },
            {
              path: 'OpcuaNodeid/:uid',
              name: 'OPCUA Nodeid配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/opcua/OpcuaNodeid'),
            },
            {
              path: 'OpcuaDetail/:uid',
              name: 'OPCUA模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/opcua/OpcuaModelDetail'),
            },

            {
              path: 'MqttModel',
              name: 'MQTT设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/mqtt/mqttModel'),
            },
            {
              path: 'MqttAdd',
              name: '添加MQTT数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/mqtt/mqttModelAdd'),
            },
            {
              path: 'MqttNodeid/:uid',
              name: 'MQTT配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/mqtt/mqttNodeid'),
            },
            {
              path: 'MqttDetail/:uid',
              name: 'MQTT模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/mqtt/mqttModelDetail'),
            },
            {
              path: 'SiemensS7Model',
              name: '西门子S7设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/SimS7/SimS7Model'),
            },
            {
              path: 'SimS7Add',
              name: '添加S7模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/SimS7/S7ModelAdd'),
            },
            {
              path: 'S7DataList/:uid',
              name: 'S7数据配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/SimS7/S7ModelDataList'),
            },
            {
              path: 'SimS7Detail/:uid',
              name: 'S7模型详细',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/SimS7/S7ModelDetail'),
            },

            {
              path: 'IEC104Model',
              name: 'IEC104设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC104/IEC104Model'),
            },
            {
              path: 'IEC104ModelAdd',
              name: '添加IEC104数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC104/IEC104ModelAdd'),
            },
            {
              path: 'IEC104ModelData/:uid',
              name: 'IEC104数据配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC104/IEC104Data'),
            },
            {
              path: 'IEC104ModelDetail/:uid',
              name: 'IEC104模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC104/IEC104ModelDetail'),
            },

            {
              path: 'IEC61850Model',
              name: 'IEC61850设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC61850/IEC61850Model'),
            },
            {
              path: 'IEC61850Add',
              name: '添加IEC61850数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC61850/IEC61850ModelAdd'),
            },
            {
              path: 'IEC61850Nodeid/:uid',
              name: 'IEC61850 Nodeid配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC61850/IEC61850Nodeid'),
            },
            {
              path: 'IEC61850Detail/:uid',
              name: 'IEC61850模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/IEC61850/IEC61850ModelDetail'),
            },
            {
              path: 'DNP3Model',
              name: 'DNP3设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DNP3/DNP3Model'),
            },
            {
              path: 'DNP3ModelAdd',
              name: '添加DNP3数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DNP3/DNP3ModelAdd'),
            },
            {
              path: 'DNP3ModelData/:uid',
              name: 'DNP3数据配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DNP3/DNP3Data'),
            },
            {
              path: 'DNP3ModelDetail/:uid',
              name: 'DNP3模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/DNP3/DNP3ModelDetail'),
            },
            {
              path: 'BACnetModel',
              name: 'BACnet设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/bacnet/bacnetModel'),
            },
            {
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
            {
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
            {
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
            {
              path: 'RestFulModel',
              name: 'RESTFul设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/RESTFulData/RESTFulModel'),
            },
            {
              path: 'RestFulData/:uid',
              name: 'RESTFul模型数据',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/RESTFulData/RESTFulDataList'),
            },
            {
              path: 'HJ212Model',
              name: '环保HJ212-2017',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/hj212/hj212Model'),
            },
            {
              path: 'HJ212Add',
              name: '添加环保数据模型',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/hj212/hj212ModelAdd'),
            },
            {
              path: 'HJ212Nodeid/:uid',
              name: '环保协议数据配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/hj212/hj212Nodeid'),
            },
            {
              path: 'HJ212Detail/:uid',
              name: '环保模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/hj212/hj212ModelDetail'),
            },
            {
              path: 'VirtualDevice',
              name: '虚拟设备',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/VirtualDevice/VirtualDeviceModel'),
            },
            {
              path: 'VirtualDeviceData/:uid',
              name: '虚拟模型数据',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/VirtualDevice/VirtualDeviceDataList'),
            },
            {
              path: 'DeviceCustomData',
              name: '自定义数据',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/CustomData/CustomDataModel'),
            },
            {
              path: 'SystemData',
              name: '系统变量',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/static/StaticModel'),
            },
            {
              path: 'StaticDataAdd',
              name: '添加静态数据模型',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                },
                invisible:true
              },
              component: () => import('@/pages/dataModel/static/StaticModelAdd'),
            },
            {
              path: 'StaticDataDetail/:uid',
              name: '静态数据模型详情',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/dataModel/static/StaticModelDetail'),
            }
          ]
        },
        {
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
        {
          path: 'TemplateMarket',
          name: '模板市场',
          meta: {
            icon: 'appstore',
            authority: {
              role: ['Admin','Operator']
            }
          },
          component: () => import('@/pages/disPlayModel/TemplateMarket')
        },
        {
          path: 'AIProjectGenerator',
          name: 'AI项目生成器',
          meta: {
            icon:"robot",
            authority: {
              role: ['Admin','Operator']
            }
          },
          component: () => import('@/pages/disPlayModel/displayModelList')
        },
        {
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
        {
          path: 'VideoManager',
          name: '视频管理',
          meta: {
            icon:"video-camera"
          },
          component: BlankView,
          children: [
            {
              path: 'videoScreen',
              name: '分屏展示',
              meta: {
                icon:"table",
              },
              component: () => import('@/pages/video/videoShowMore')
            },
            {
              path: 'videoList',
              name: '视频列表',
              meta: {
                icon:"unordered-list"
              },
              component: () => import('@/pages/video/videoManager')
            },
            {
              path: 'GB28281List',
              name: '国标视频',
              meta: {
                icon:"youtube"
              },
              component: () => import('@/pages/video/GB28281VideoManager')
            }
          ]
        },
        {
          path: 'Real-timeAlarm',
          name: '实时告警',
          meta: {
            icon:"bell"
          },
          component: () => import('@/pages/alarm/currentAlarm/currentAlarm')
        },
        {
          path: 'AlarmStrategy',
          name: '告警策略',
          meta: {
            icon: 'alert'
          },
          component: BlankView,
          children: [
            {
              path: 'ModelTrigger',
              name: '模型触发器',
              meta: {
                icon:"thunderbolt"
              },
              component: () => import('@/pages/alarm/trigger/trigger'),
            },
            {
              path: 'AlarmRestoreMask',
              name: '告警恢复',
              meta: {
                icon:"undo",
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/alarm/ShieldAlarm/shieldAlarm'),
            }
          ]
        },
        {
          path: 'TaskPlan',
          name: '任务计划',
          meta: {
            icon:"schedule"
          },
          component: () => import('@/pages/taskplan/taskplan'),
        },
        {
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
        {
          path: 'DataPush',
          name: '数据接口',
          meta: {
            icon: 'radar-chart'
          },
          component: BlankView,
          children: [
            {
              path: 'DataTemplete',
              name: '数据模版',
              meta: {
                icon: 'file-text'
              },
              component: () => import('@/pages/DataPush/DataTemplete'),
            },
            {
              path: 'IEC104DataTemplete',
              name: 'IEC104数据模版',
              meta: {
                icon: 'file-text'
              },
              component: () => import('@/pages/DataPush/iec104/IEC104DataTemplete'),
            },
            {
              path: 'IEC104TempleteData/:uid',
              name: 'IEC104数据模版配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/DataPush/iec104/IEC104PusbData'),
            },
            {
              path: 'ModbusDataTemplete',
              name: 'Modbus 数据模版',
              meta: {
                icon: 'file-text'
              },
              component: () => import('@/pages/DataPush/modbusTcp/ModbusTcpDataTemplete'),
            },
            {
              path: 'ModbusTempleteData/:uid',
              name: 'Modbus 数据模版配置',
              meta: {
                invisible:true,
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/DataPush/modbusTcp/ModbusTcpPusbData'),
            },
            {
              path: 'DataInterface',
              name: '接口管理',
              meta: {
                icon: 'hdd'
              },
              component: () => import('@/pages/DataPush/DataInterface'),
            },
          ]
        },
        {
          path: 'Reporting',
          name: '数据报表',
          meta: {
            icon: 'file-done'
          },
          component: BlankView,
          children: [
            {
              path: 'AlarmHistory',
              name: '告警报表',
              meta: {

              },
              component: () => import('@/pages/reporting/alarmReport/alarmHistory'),
            },
            {
              path: 'DataHistory',
              name: '历史数据',
              meta: {

              },
              component: () => import('@/pages/reporting/dataHistoryReport/dataHistoryReport'),
            },
            {
              path: 'DataHistoryReport',
              name: '历史报表',
              meta: {

              },
              component: () => import('@/pages/reporting/dataHistoryReport/dataHistoryQueryReport'),
            },
            {
              path: 'SqlReportTemplate',
              name: 'SQL报表',
              meta: {

              },
              component: () => import('@/pages/reporting/sqlReport/SqlReportTemplate'),
            },
            {
              path: 'DiyReport',
              name: '自定义报表',
              meta: {

              },
              component: () => import('@/pages/reporting/diyReport/diyReport'),
            },
            {
              path: 'DiyReportTemplete',
              name: '报表模板',
              meta: {

              },
              component: () => import('@/pages/reporting/diyReport/diyReportTemplete'),
            },
            {
              path: 'ReportTempleteContent/:uuid',
              name: '报表模板内容修改',
              meta: {
                authority: {
                  role: ['Admin','Operator']
                },
                invisible:true
              },
              component: () => import('@/pages/reporting/diyReport/diyReportContent'),
            },
          ]
        },
        // {
        //   path: 'AdvanceReporting',
        //   name: '高级报表',
        //   meta: {
        //     icon: 'file-done'
        //   },
        //   component: BlankView,
        //   children: [
        //     {
        //       path: 'AlarmHistory',
        //       name: '报表测试',
        //       meta: {
        //
        //       },
        //       component: () => import('@/pages/AdvanceReporting/text'),
        //     }
        //   ]
        // },
        {
          path: 'Network',
          name: '网络中心',
          meta: {
            icon: 'global'
          },
          component: BlankView,
          children: [
            {
              path: 'SystemNetwork',
              name: '网络',
              meta: {
                icon:"ie",
              },
              component: () => import('@/pages/network/network')
            },
            {
              path: 'ISMNetwork',
              name: '组网',
              meta: {
                icon:"slack",
              },
              component: () => import('@/pages/network/ismnode')
            }
          ]
        },
        {
          path: 'DataBase',
          name: '数据库管理',
          meta: {
            icon: 'database'
          },
          component: BlankView,
          children: [
            {
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
            {
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
          ]
        },
        {
          path: 'Setting',
          name: '设置中心',
          meta: {
            icon: 'setting'
          },
          component: BlankView,
          children: [
            {
              path: 'Account',
              name: '个人设置',
              meta: {
                icon:"user",
              },
              component: () => import('@/pages/account/settings/Index')
            },
            {
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
            {
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
            {
              path: 'AlarmTipsSetting',
              name: '告警通知',
              meta: {
                icon:"sound",
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/Setting/AlarmTips')
            },
            {
              path: 'SystemParams',
              name: '系统参数',
              meta: {
                icon:"tool",
                authority: {
                  role: ['Admin','Operator']
                }
              },
              component: () => import('@/pages/Setting/SystemParams')
            },
            {
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
            {
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
          ]
        },
        {
          path: 'Journal',
          name: '系统日志',
          meta: {
            icon: 'file-text'
          },
          component: BlankView,
          children: [
            {
              path: 'OperationJournal',
              name: '操作日志',
              meta: {

              },
              component: () => import('@/pages/journal/OperationLog/OperationLog'),
            },
          ]
        },
        {
          path: '/Help',
          name: '帮助中心',
          meta: {
            icon: 'question-circle',
            authority: {
              role: ['Admin']
            },
          },
          component: BlankView,
          children:[
            {
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
            {
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
            // {
            //   path: '/HelpDocument',
            //   name: '帮助文档',
            //   meta: {
            //     icon:"read",
            //   },
            //   component: () => import('@/pages/help/help')
            // },
          ]
        },
      ]
    },
  ]
}

export default options
/**
 * IndustrialObjects - 3D 对象工厂 + 组件库定义
 * 参照 3d-configuration-editor 完整实现
 */
import * as THREE from 'three'
import { CSS2DObject } from 'three/examples/jsm/renderers/CSS2DRenderer.js'
import { createDTObject } from './DigitalTwinObjects.js'

// ===== 颜色映射 =====
const COMP_COLORS = {
  blue: '#4da6ff', green: '#4dffa6', orange: '#ffaa4d',
  purple: '#c17bff', pink: '#ff7bb5', cyan: '#4dfff0', yellow: '#fffb4d',
  red: '#ff4d4f', gray: '#8c8c8c', brown: '#a67c52'
}

// ===== 组件库定义 =====
const COMPONENT_LIBRARY = [
  {
    name: '基础几何', items: [
      { type: 'box',        label: '立方体',   icon: 'fas fa-cube',          color: 'blue',   typeName: 'BoxGeometry' },
      { type: 'sphere',     label: '球体',     icon: 'fas fa-circle',        color: 'green',  typeName: 'SphereGeometry' },
      { type: 'cylinder',   label: '圆柱',     icon: 'fas fa-database',      color: 'orange', typeName: 'CylinderGeometry' },
      { type: 'cone',       label: '圆锥',     icon: 'fas fa-play',          color: 'purple', typeName: 'ConeGeometry' },
      { type: 'torus',      label: '圆环',     icon: 'fas fa-dot-circle',    color: 'pink',   typeName: 'TorusGeometry' },
      { type: 'plane',      label: '平面',     icon: 'fas fa-square',        color: 'cyan',   typeName: 'PlaneGeometry' },
    ]
  },
  {
    name: '泵阀管件', items: [
      { type: 'pipe',         label: '管道',       icon: 'fas fa-grip-lines',         color: 'orange', typeName: 'Pipe' },
      { type: 'pipeelbow',    label: '弯头',       icon: 'fas fa-angle-right',        color: 'orange', typeName: 'PipeElbow' },
      { type: 'teepipe',      label: '三通',       icon: 'fas fa-code-branch',        color: 'orange', typeName: 'TeePipe' },
      { type: 'crosspipe',    label: '四通',       icon: 'fas fa-crosshairs',         color: 'orange', typeName: 'CrossPipe' },
      { type: 'reducerpipe',  label: '异径管',     icon: 'fas fa-arrows-alt-v',       color: 'orange', typeName: 'Reducer' },
      { type: 'pipecap',      label: '管帽',       icon: 'fas fa-arrow-down',         color: 'orange', typeName: 'PipeCap' },
      { type: 'flange',       label: '法兰',       icon: 'fas fa-circle-notch',       color: 'gray',   typeName: 'Flange' },
      { type: 'valve',        label: '阀门',       icon: 'fas fa-sliders-h',          color: 'purple', typeName: 'Valve' },
      { type: 'ballvalve',    label: '球阀',       icon: 'fas fa-circle',             color: 'purple', typeName: 'BallValve' },
      { type: 'butterflyvalve', label: '蝶阀',     icon: 'fas fa-paper-plane',        color: 'purple', typeName: 'ButterflyValve' },
      { type: 'checkvalve',   label: '止回阀',     icon: 'fas fa-arrow-circle-up',    color: 'purple', typeName: 'CheckValve' },
      { type: 'safetyvalve',  label: '安全阀',     icon: 'fas fa-exclamation-triangle', color: 'red', typeName: 'SafetyValve' },
      { type: 'gatevalve',    label: '闸阀',       icon: 'fas fa-arrows-alt-h',       color: 'purple', typeName: 'GateValve' },
      { type: 'pump',         label: '泵',         icon: 'fas fa-fan',                color: 'green',  typeName: 'Pump' },
      { type: 'centrifugalpump', label: '离心泵',  icon: 'fas fa-spinner',            color: 'green',  typeName: 'CentrifugalPump' },
      { type: 'gearpump',     label: '齿轮泵',     icon: 'fas fa-cogs',               color: 'green',  typeName: 'GearPump' },
      { type: 'motor',        label: '电机',       icon: 'fas fa-cog',                color: 'red',    typeName: 'Motor' },
      { type: 'fan',          label: '风机',       icon: 'fas fa-wind',               color: 'cyan',   typeName: 'Fan' },
      { type: 'coolingFanUnit', label: '双风机冷却机组', icon: 'fas fa-fan',           color: 'cyan',   typeName: 'CoolingFanUnit' },
      { type: 'compressor',   label: '压缩机',     icon: 'fas fa-compress-arrows-alt', color: 'blue',  typeName: 'Compressor' },
    ]
  },
  {
    name: '容器储罐', items: [
      { type: 'tank',           label: '储罐',       icon: 'fas fa-oil-can',           color: 'blue',   typeName: 'Tank' },
      { type: 'pressurevessel', label: '压力容器',   icon: 'fas fa-dumbbell',          color: 'red',    typeName: 'PressureVessel' },
      { type: 'reactor',        label: '反应釜',     icon: 'fas fa-flask',             color: 'purple', typeName: 'Reactor' },
      { type: 'boiler',         label: '锅炉',       icon: 'fas fa-fire',              color: 'orange', typeName: 'Boiler' },
      { type: 'buffertank',     label: '缓冲罐',     icon: 'fas fa-arrow-alt-circle-down', color: 'blue', typeName: 'BufferTank' },
      { type: 'mixingtank',     label: '混合罐',     icon: 'fas fa-blender',           color: 'purple', typeName: 'MixingTank' },
      { type: 'airreceiver',    label: '储气罐',     icon: 'fas fa-cloud-upload-alt',  color: 'blue',   typeName: 'AirReceiver' },
      { type: 'silo',           label: '筒仓',       icon: 'fas fa-warehouse',         color: 'gray',   typeName: 'Silo' },
      { type: 'hopper',         label: '料斗',       icon: 'fas fa-filter',            color: 'brown',  typeName: 'Hopper' },
    ]
  },
  {
    name: '换热分离', items: [
      { type: 'heatex',         label: '换热器',     icon: 'fas fa-th',                color: 'blue',   typeName: 'HeatExchanger' },
      { type: 'condenser',      label: '冷凝器',     icon: 'fas fa-snowflake',         color: 'cyan',   typeName: 'Condenser' },
      { type: 'evaporator',     label: '蒸发器',     icon: 'fas fa-temperature-high',  color: 'orange', typeName: 'Evaporator' },
      { type: 'coolingtower',   label: '冷却塔',     icon: 'fas fa-water',             color: 'blue',   typeName: 'CoolingTower' },
      { type: 'separator',      label: '分离器',     icon: 'fas fa-divide',            color: 'cyan',   typeName: 'Separator' },
      { type: 'mixer',          label: '搅拌器',     icon: 'fas fa-blender',           color: 'orange', typeName: 'Mixer' },
      { type: 'cyclone',        label: '旋风分离器', icon: 'fas fa-tornado',           color: 'cyan',   typeName: 'Cyclone' },
      { type: 'filterunit',     label: '过滤器',     icon: 'fas fa-asterisk',          color: 'green',  typeName: 'FilterUnit' },
      { type: 'dustcollector',  label: '除尘器',     icon: 'fas fa-wind',              color: 'gray',   typeName: 'DustCollector' },
      { type: 'scrubber',       label: '洗涤塔',     icon: 'fas fa-shower',            color: 'blue',   typeName: 'Scrubber' },
      { type: 'dryer',          label: '干燥器',     icon: 'fas fa-sun',               color: 'orange', typeName: 'Dryer' },
      { type: 'furnace',        label: '加热炉',     icon: 'fas fa-fire-alt',          color: 'red',    typeName: 'Furnace' },
      { type: 'chimney',        label: '烟囱',       icon: 'fas fa-smoke',             color: 'gray',   typeName: 'Chimney' },
    ]
  },
  {
    name: '输送存储', items: [
      { type: 'conveyor',       label: '传送带',     icon: 'fas fa-stream',            color: 'gray',   typeName: 'Conveyor' },
      { type: 'screwconveyor',  label: '螺旋输送机', icon: 'fas fa-arrow-right',       color: 'gray',   typeName: 'ScrewConveyor' },
      { type: 'bucketelevator', label: '斗式提升机', icon: 'fas fa-arrow-up',          color: 'gray',   typeName: 'BucketElevator' },
      { type: 'bin',            label: '料仓',       icon: 'fas fa-box',               color: 'gray',   typeName: 'Bin' },
      { type: 'funnel',         label: '漏斗',       icon: 'fas fa-vial',              color: 'brown',  typeName: 'Funnel' },
    ]
  },
  {
    name: '仪器仪表', items: [
      { type: 'gauge',          label: '仪表盘',     icon: 'fas fa-tachometer-alt',    color: 'yellow', typeName: 'Gauge' },
      { type: 'sensor',         label: '传感器',     icon: 'fas fa-microchip',         color: 'green',  typeName: 'Sensor' },
      { type: 'thermometer',    label: '温度计',     icon: 'fas fa-thermometer-half',  color: 'red',    typeName: 'Thermometer' },
      { type: 'pressuregauge',  label: '压力表',     icon: 'fas fa-compress',          color: 'yellow', typeName: 'PressureGauge' },
      { type: 'flowmeter',      label: '流量计',     icon: 'fas fa-tint',              color: 'blue',   typeName: 'FlowMeter' },
      { type: 'levelgauge',     label: '液位计',     icon: 'fas fa-ruler-vertical',    color: 'cyan',   typeName: 'LevelGauge' },
    ]
  },
  {
    name: '管道支架', items: [
      { type: 'pipesupport',    label: '管架',       icon: 'fas fa-grip-lines-vertical', color: 'gray', typeName: 'PipeSupport' },
      { type: 'pipebracket',    label: '管托',       icon: 'fas fa-level-up-alt',      color: 'gray',   typeName: 'PipeBracket' },
    ]
  },
  {
    name: '其他设备', items: [
      { type: 'flare',          label: '火炬',       icon: 'fas fa-burn',              color: 'red',    typeName: 'Flare' },
      { type: 'silencer',       label: '消音器',     icon: 'fas fa-volume-mute',       color: 'gray',   typeName: 'Silencer' },
    ]
  },
  {
    name: '智慧城市', items: [
      { type: 'building',       label: '建筑楼',     icon: 'fas fa-building',          color: 'blue',   typeName: 'Building' },
      { type: 'apartment',      label: '住宅楼',     icon: 'fas fa-home',              color: 'orange', typeName: 'Apartment' },
      { type: 'factorybld',     label: '厂房',       icon: 'fas fa-industry',          color: 'gray',   typeName: 'FactoryBld' },
      { type: 'streetlight',    label: '路灯',       icon: 'fas fa-street-view',       color: 'yellow', typeName: 'StreetLight' },
      { type: 'trafficlight',   label: '信号灯',     icon: 'fas fa-traffic-light',     color: 'red',    typeName: 'TrafficLight' },
      { type: 'bridge',         label: '桥梁',       icon: 'fas fa-archway',           color: 'gray',   typeName: 'Bridge' },
      { type: 'tunnel',         label: '隧道',       icon: 'fas fa-archway',           color: 'gray',   typeName: 'Tunnel' },
      { type: 'road',           label: '道路',       icon: 'fas fa-road',              color: 'gray',   typeName: 'Road' },
      { type: 'parkingspot',    label: '停车位',     icon: 'fas fa-parking',           color: 'blue',   typeName: 'ParkingSpot' },
      { type: 'busstop',        label: '公交站',     icon: 'fas fa-bus',               color: 'green',  typeName: 'BusStop' },
      { type: 'bench',          label: '长椅',       icon: 'fas fa-chair',             color: 'brown',  typeName: 'Bench' },
      { type: 'trashbin',       label: '垃圾桶',     icon: 'fas fa-trash',             color: 'green',  typeName: 'TrashBin' },
      { type: 'citytree',       label: '树木',       icon: 'fas fa-tree',              color: 'green',  typeName: 'CityTree' },
      { type: 'fountain',       label: '喷泉',       icon: 'fas fa-water',             color: 'cyan',   typeName: 'Fountain' },
      { type: 'surveilcam',     label: '摄像头',     icon: 'fas fa-video',             color: 'purple', typeName: 'SurveilCam' },
      { type: 'billboard',      label: '广告牌',     icon: 'fas fa-ad',                color: 'yellow', typeName: 'Billboard' },
      { type: 'solarpanel',     label: '太阳能板',   icon: 'fas fa-solar-panel',       color: 'blue',   typeName: 'SolarPanel' },
      { type: 'windturbine',    label: '风力发电机', icon: 'fas fa-fan',               color: 'cyan',   typeName: 'WindTurbine' },
      { type: 'telecomtower',   label: '通信塔',     icon: 'fas fa-broadcast-tower',   color: 'orange', typeName: 'TelecomTower' },
      { type: 'manhole',        label: '井盖',       icon: 'fas fa-circle',            color: 'gray',   typeName: 'Manhole' },
      { type: 'hydrant',        label: '消防栓',     icon: 'fas fa-fire-extinguisher', color: 'red',    typeName: 'Hydrant' },
      { type: 'bollard',        label: '路桩',       icon: 'fas fa-circle',            color: 'yellow', typeName: 'Bollard' },
      { type: 'flagpole',       label: '旗杆',       icon: 'fas fa-flag',              color: 'red',    typeName: 'Flagpole' },
    ]
  },
  {
    name: '工业流水线', items: [
      { type: 'robotarm',       label: '机械臂',     icon: 'fas fa-robot',             color: 'orange', typeName: 'RobotArm' },
      { type: 'conveyorbelt',   label: '传送带段',   icon: 'fas fa-stream',            color: 'gray',   typeName: 'ConveyorBelt' },
      { type: 'assemblystation', label: '装配站',    icon: 'fas fa-tools',             color: 'blue',   typeName: 'AssemblyStation' },
      { type: 'weldingrobot',   label: '焊接机器人', icon: 'fas fa-fire',              color: 'yellow', typeName: 'WeldingRobot' },
      { type: 'inspectioncam',  label: '检测相机',   icon: 'fas fa-camera',            color: 'purple', typeName: 'InspectionCam' },
      { type: 'packmachine',    label: '包装机',     icon: 'fas fa-box',               color: 'green',  typeName: 'PackMachine' },
      { type: 'palletizer',     label: '码垛机',     icon: 'fas fa-layer-group',       color: 'orange', typeName: 'Palletizer' },
      { type: 'agv',            label: 'AGV小车',    icon: 'fas fa-truck-moving',      color: 'blue',   typeName: 'AGV' },
      { type: 'feeder',         label: '上料机',     icon: 'fas fa-arrow-alt-circle-up', color: 'green', typeName: 'Feeder' },
      { type: 'stampingpress',  label: '冲压机',     icon: 'fas fa-hammer',            color: 'red',    typeName: 'StampingPress' },
      { type: 'cncmachine',     label: 'CNC机床',    icon: 'fas fa-cogs',              color: 'blue',   typeName: 'CNCMachine' },
      { type: 'injectionmold',  label: '注塑机',     icon: 'fas fa-syringe',           color: 'purple', typeName: 'InjectionMold' },
      { type: 'workbench',      label: '工作台',     icon: 'fas fa-table',             color: 'brown',  typeName: 'Workbench' },
      { type: 'turntable',      label: '转台',       icon: 'fas fa-sync-alt',          color: 'gray',   typeName: 'Turntable' },
      { type: 'lifter',         label: '升降机',     icon: 'fas fa-arrow-up',          color: 'orange', typeName: 'Lifter' },
      { type: 'scanner3d',      label: '3D扫描仪',   icon: 'fas fa-qrcode',            color: 'cyan',   typeName: 'Scanner3D' },
      { type: 'labelprinter',   label: '贴标机',     icon: 'fas fa-tag',               color: 'yellow', typeName: 'LabelPrinter' },
      { type: 'qualitygate',    label: '质检门',     icon: 'fas fa-check-double',      color: 'green',  typeName: 'QualityGate' },
      { type: 'buffertable',    label: '缓存台',     icon: 'fas fa-pause-circle',      color: 'cyan',   typeName: 'BufferTable' },
      { type: 'safetyfence',    label: '安全围栏',   icon: 'fas fa-shield-alt',        color: 'yellow', typeName: 'SafetyFence' },
    ]
  },
  {
    name: '数字孪生', items: [
      { type: 'serverrack',     label: '服务器机柜', icon: 'fas fa-server',            color: 'gray',   typeName: 'ServerRack' },
      { type: 'datacenter',     label: '数据中心',   icon: 'fas fa-cloud',             color: 'blue',   typeName: 'DataCenter' },
      { type: 'iotdevice',      label: 'IoT设备',    icon: 'fas fa-microchip',         color: 'green',  typeName: 'IoTDevice' },
      { type: 'gateway',        label: '网关',       icon: 'fas fa-network-wired',     color: 'purple', typeName: 'Gateway' },
      { type: 'controlpanel',   label: '控制面板',   icon: 'fas fa-sliders-h',         color: 'orange', typeName: 'ControlPanel' },
      { type: 'dashboardscreen', label: '仪表大屏',  icon: 'fas fa-desktop',           color: 'blue',   typeName: 'DashboardScreen' },
      { type: 'antenna',        label: '天线',       icon: 'fas fa-satellite-dish',    color: 'gray',   typeName: 'Antenna' },
      { type: 'radardish',      label: '雷达',       icon: 'fas fa-satellite',         color: 'cyan',   typeName: 'RadarDish' },
      { type: 'router',         label: '路由器',     icon: 'fas fa-wifi',              color: 'blue',   typeName: 'Router' },
      { type: 'switchgear',     label: '交换机',     icon: 'fas fa-exchange-alt',      color: 'gray',   typeName: 'Switchgear' },
      { type: 'ups',            label: 'UPS电源',    icon: 'fas fa-battery-full',      color: 'green',  typeName: 'UPS' },
      { type: 'patchpanel',     label: '配线架',     icon: 'fas fa-plug',              color: 'gray',   typeName: 'PatchPanel' },
      { type: 'fiberbox',       label: '光纤盒',     icon: 'fas fa-fiber',             color: 'yellow', typeName: 'FiberBox' },
      { type: 'rtu',            label: 'RTU设备',    icon: 'fas fa-hdd',               color: 'purple', typeName: 'RTU' },
      { type: 'plc',            label: 'PLC控制器',  icon: 'fas fa-calculator',        color: 'orange', typeName: 'PLC' },
      { type: 'hmipanel',       label: 'HMI面板',    icon: 'fas fa-tablet-alt',        color: 'blue',   typeName: 'HMIPanel' },
      { type: 'beacon',         label: '信标',       icon: 'fas fa-dot-circle',        color: 'red',    typeName: 'Beacon' },
      { type: 'sensorarray',    label: '传感器阵列', icon: 'fas fa-th-large',          color: 'green',  typeName: 'SensorArray' },
    ]
  },
  {
    name: '能源电力', items: [
      { type: 'transformer',    label: '变压器',     icon: 'fas fa-bolt',              color: 'red',    typeName: 'Transformer' },
      { type: 'powerpole',      label: '电线杆',     icon: 'fas fa-grip-lines',        color: 'gray',   typeName: 'PowerPole' },
      { type: 'enginegen',      label: '发电机',     icon: 'fas fa-cog',               color: 'orange', typeName: 'EngineGen' },
      { type: 'solararray',     label: '太阳能阵列', icon: 'fas fa-sun',               color: 'blue',   typeName: 'SolarArray' },
      { type: 'substation',     label: '变电站',     icon: 'fas fa-charging-station',  color: 'gray',   typeName: 'Substation' },
      { type: 'batterystorage', label: '储能电池',   icon: 'fas fa-car-battery',       color: 'green',  typeName: 'BatteryStorage' },
      { type: 'powermeter',     label: '电表',       icon: 'fas fa-tachometer-alt',    color: 'yellow', typeName: 'PowerMeter' },
      { type: 'circuitbreaker', label: '断路器',     icon: 'fas fa-toggle-on',         color: 'red',    typeName: 'CircuitBreaker' },
      { type: 'insulator',      label: '绝缘子',     icon: 'fas fa-grip-lines-vertical', color: 'gray', typeName: 'Insulator' },
      { type: 'cabletray',      label: '电缆桥架',   icon: 'fas fa-ellipsis-h',        color: 'gray',   typeName: 'CableTray' },
      { type: 'junctionbox',    label: '接线盒',     icon: 'fas fa-box-open',          color: 'gray',   typeName: 'JunctionBox' },
      { type: 'lightningrod',   label: '避雷针',     icon: 'fas fa-bolt',              color: 'yellow', typeName: 'LightningRod' },
    ]
  },
  {
    name: '仓储物流', items: [
      { type: 'warehouserack',  label: '货架',       icon: 'fas fa-th-list',           color: 'orange', typeName: 'WarehouseRack' },
      { type: 'forklift',       label: '叉车',       icon: 'fas fa-truck',             color: 'yellow', typeName: 'Forklift' },
      { type: 'towercrane',     label: '塔吊',       icon: 'fas fa-building',          color: 'orange', typeName: 'TowerCrane' },
      { type: 'container',      label: '集装箱',     icon: 'fas fa-cube',              color: 'blue',   typeName: 'Container' },
      { type: 'loadingdock',    label: '装卸台',     icon: 'fas fa-warehouse',         color: 'gray',   typeName: 'LoadingDock' },
      { type: 'pallet',         label: '托盘',       icon: 'fas fa-pallet',            color: 'brown',  typeName: 'Pallet' },
      { type: 'shelving',       label: '置物架',     icon: 'fas fa-border-all',        color: 'brown',  typeName: 'Shelving' },
      { type: 'barcodereader',  label: '条码阅读器', icon: 'fas fa-barcode',           color: 'red',    typeName: 'BarcodeReader' },
      { type: 'sorter',         label: '分拣机',     icon: 'fas fa-sort-amount-down',  color: 'green',  typeName: 'Sorter' },
      { type: 'stackercrane',   label: '堆垛机',     icon: 'fas fa-arrows-alt-v',      color: 'orange', typeName: 'StackerCrane' },
      { type: 'dockleveler',    label: '装卸平台',   icon: 'fas fa-equals',            color: 'gray',   typeName: 'DockLeveler' },
      { type: 'conveyorjunc',   label: '传送交汇',   icon: 'fas fa-code-branch',       color: 'gray',   typeName: 'ConveyorJunc' },
    ]
  },
  {
    name: '交通设施', items: [
      { type: 'car',            label: '小汽车',     icon: 'fas fa-car',               color: 'blue',   typeName: 'Car' },
      { type: 'suv',            label: 'SUV',        icon: 'fas fa-car',               color: 'green',  typeName: 'SUV' },
      { type: 'bus',            label: '公交车',     icon: 'fas fa-bus',               color: 'yellow', typeName: 'Bus' },
      { type: 'truck',          label: '卡车',       icon: 'fas fa-truck',             color: 'orange', typeName: 'Truck' },
      { type: 'firetruck',      label: '消防车',     icon: 'fas fa-fire',              color: 'red',    typeName: 'FireTruck' },
      { type: 'ambulance',      label: '救护车',     icon: 'fas fa-ambulance',         color: 'cyan',   typeName: 'Ambulance' },
      { type: 'policecar',      label: '警车',       icon: 'fas fa-shield-alt',        color: 'purple', typeName: 'PoliceCar' },
      { type: 'motorcycle',     label: '摩托车',     icon: 'fas fa-motorcycle',        color: 'red',    typeName: 'Motorcycle' },
      { type: 'bicycle',        label: '自行车',     icon: 'fas fa-bicycle',           color: 'green',  typeName: 'Bicycle' },
      { type: 'train',          label: '火车',       icon: 'fas fa-train',             color: 'green',  typeName: 'Train' },
      { type: 'metro',          label: '地铁',       icon: 'fas fa-subway',            color: 'blue',   typeName: 'Metro' },
      { type: 'airplane',       label: '飞机',       icon: 'fas fa-plane',             color: 'cyan',   typeName: 'Airplane' },
      { type: 'helicopter',     label: '直升机',     icon: 'fas fa-helicopter',        color: 'gray',   typeName: 'Helicopter' },
      { type: 'drone',          label: '无人机',     icon: 'fas fa-fighter-jet',       color: 'purple', typeName: 'Drone' },
      { type: 'ship',           label: '轮船',       icon: 'fas fa-ship',              color: 'blue',   typeName: 'Ship' },
      { type: 'sailboat',       label: '帆船',       icon: 'fas fa-sailboat',          color: 'cyan',   typeName: 'Sailboat' },
      { type: 'trafficcone',    label: '锥形桶',     icon: 'fas fa-exclamation-triangle', color: 'orange', typeName: 'TrafficCone' },
      { type: 'barrier',        label: '护栏',       icon: 'fas fa-arrows-alt-h',      color: 'yellow', typeName: 'Barrier' },
      { type: 'speedbump',      label: '减速带',     icon: 'fas fa-grip-lines',        color: 'yellow', typeName: 'SpeedBump' },
      { type: 'zebracrossing',  label: '斑马线',     icon: 'fas fa-road',              color: 'gray',   typeName: 'ZebraCrossing' },
    ]
  },
  {
    name: '室内设施', items: [
      { type: 'desk',           label: '办公桌',     icon: 'fas fa-desktop',           color: 'brown',  typeName: 'Desk' },
      { type: 'officechair',    label: '办公椅',     icon: 'fas fa-chair',             color: 'gray',   typeName: 'OfficeChair' },
      { type: 'sofa',           label: '沙发',       icon: 'fas fa-couch',             color: 'orange', typeName: 'Sofa' },
      { type: 'coffeetable',    label: '茶几',       icon: 'fas fa-coffee',            color: 'brown',  typeName: 'CoffeeTable' },
      { type: 'bookshelf',      label: '书架',       icon: 'fas fa-book',              color: 'brown',  typeName: 'Bookshelf' },
      { type: 'cabinet',        label: '柜子',       icon: 'fas fa-archive',           color: 'brown',  typeName: 'Cabinet' },
      { type: 'bed',            label: '床',         icon: 'fas fa-bed',               color: 'blue',   typeName: 'Bed' },
      { type: 'floorlamp',      label: '落地灯',     icon: 'fas fa-lightbulb',         color: 'yellow', typeName: 'FloorLamp' },
      { type: 'ceilinglight',   label: '吸顶灯',     icon: 'fas fa-lightbulb',         color: 'yellow', typeName: 'CeilingLight' },
      { type: 'monitor',        label: '显示器',     icon: 'fas fa-tv',                color: 'blue',   typeName: 'Monitor' },
      { type: 'keyboard',       label: '键盘',       icon: 'fas fa-keyboard',          color: 'gray',   typeName: 'Keyboard' },
      { type: 'computer',       label: '电脑主机',   icon: 'fas fa-desktop',           color: 'gray',   typeName: 'Computer' },
      { type: 'printer',        label: '打印机',     icon: 'fas fa-print',             color: 'gray',   typeName: 'Printer' },
      { type: 'whiteboard',     label: '白板',       icon: 'fas fa-chalkboard',        color: 'cyan',   typeName: 'Whiteboard' },
      { type: 'meetingtable',   label: '会议桌',     icon: 'fas fa-users',             color: 'brown',  typeName: 'MeetingTable' },
      { type: 'waterdispenser', label: '饮水机',     icon: 'fas fa-tint',              color: 'blue',   typeName: 'WaterDispenser' },
      { type: 'microwave',      label: '微波炉',     icon: 'fas fa-utensils',          color: 'gray',   typeName: 'Microwave' },
      { type: 'fridge',         label: '冰箱',       icon: 'fas fa-snowflake',         color: 'cyan',   typeName: 'Fridge' },
      { type: 'acunit',         label: '空调',       icon: 'fas fa-wind',              color: 'cyan',   typeName: 'ACUnit' },
      { type: 'wallclock',      label: '时钟',       icon: 'fas fa-clock',             color: 'yellow', typeName: 'WallClock' },
    ]
  },
  {
    name: '消防安防', items: [
      { type: 'fireextinguisher', label: '灭火器',   icon: 'fas fa-fire-extinguisher', color: 'red',    typeName: 'FireExtinguisher' },
      { type: 'firealarm',      label: '火警铃',     icon: 'fas fa-bell',              color: 'red',    typeName: 'FireAlarm' },
      { type: 'sprinkler',      label: '喷淋头',     icon: 'fas fa-shower',            color: 'red',    typeName: 'Sprinkler' },
      { type: 'smokedetector',  label: '烟感',       icon: 'fas fa-smoking',           color: 'gray',   typeName: 'SmokeDetector' },
      { type: 'emergencysign',  label: '应急出口',   icon: 'fas fa-door-open',         color: 'green',  typeName: 'EmergencySign' },
      { type: 'firehose',       label: '消防水带',   icon: 'fas fa-water',             color: 'red',    typeName: 'FireHose' },
      { type: 'firehydrantbox', label: '消防栓箱',   icon: 'fas fa-first-aid',         color: 'red',    typeName: 'FireHydrantBox' },
      { type: 'securitygate',   label: '安检门',     icon: 'fas fa-door-closed',       color: 'gray',   typeName: 'SecurityGate' },
      { type: 'accesscontrol',  label: '门禁',       icon: 'fas fa-id-card',           color: 'purple', typeName: 'AccessControl' },
      { type: 'cctv',           label: '监控球机',   icon: 'fas fa-video',             color: 'purple', typeName: 'CCTV' },
      { type: 'alarmbell',      label: '警铃',       icon: 'fas fa-bell',              color: 'red',    typeName: 'AlarmBell' },
      { type: 'firstaidkit',    label: '急救箱',     icon: 'fas fa-medkit',            color: 'green',  typeName: 'FirstAidKit' },
    ]
  },
  {
    name: '绿化景观', items: [
      { type: 'bush',           label: '灌木',       icon: 'fas fa-leaf',              color: 'green',  typeName: 'Bush' },
      { type: 'hedge',          label: '树篱',       icon: 'fas fa-tree',              color: 'green',  typeName: 'Hedge' },
      { type: 'flowerbed',      label: '花坛',       icon: 'fas fa-seedling',          color: 'pink',   typeName: 'FlowerBed' },
      { type: 'grasspatch',     label: '草坪',       icon: 'fas fa-square',            color: 'green',  typeName: 'GrassPatch' },
      { type: 'rock',           label: '岩石',       icon: 'fas fa-mountain',          color: 'gray',   typeName: 'Rock' },
      { type: 'fence',          label: '围栏',       icon: 'fas fa-border-none',       color: 'brown',  typeName: 'Fence' },
      { type: 'gardenwall',     label: '围墙',       icon: 'fas fa-th-large',          color: 'gray',   typeName: 'GardenWall' },
      { type: 'gate',           label: '大门',       icon: 'fas fa-door-open',         color: 'gray',   typeName: 'Gate' },
      { type: 'lamppost',       label: '庭院灯',     icon: 'fas fa-lightbulb',         color: 'yellow', typeName: 'LampPost' },
      { type: 'pavilion',       label: '凉亭',       icon: 'fas fa-umbrella-beach',    color: 'brown',  typeName: 'Pavilion' },
      { type: 'pond',           label: '池塘',       icon: 'fas fa-water',             color: 'cyan',   typeName: 'Pond' },
      { type: 'sculpture',      label: '雕塑',       icon: 'fas fa-monument',          color: 'gray',   typeName: 'Sculpture' },
    ]
  },
  {
    name: '医疗设施', items: [
      { type: 'hospitalbed',    label: '病床',       icon: 'fas fa-bed',               color: 'cyan',   typeName: 'HospitalBed' },
      { type: 'operatingtable', label: '手术台',     icon: 'fas fa-procedures',        color: 'gray',   typeName: 'OperatingTable' },
      { type: 'mrimachine',     label: 'MRI设备',    icon: 'fas fa-magnet',            color: 'purple', typeName: 'MRIMachine' },
      { type: 'ctscanner',      label: 'CT扫描仪',   icon: 'fas fa-x-ray',             color: 'blue',   typeName: 'CTScanner' },
      { type: 'ivstand',        label: '输液架',     icon: 'fas fa-grip-lines-vertical', color: 'gray', typeName: 'IVStand' },
      { type: 'wheelchair',     label: '轮椅',       icon: 'fas fa-wheelchair',        color: 'blue',   typeName: 'Wheelchair' },
      { type: 'stretcher',      label: '担架',       icon: 'fas fa-procedures',        color: 'orange', typeName: 'Stretcher' },
      { type: 'medicinecabinet', label: '药品柜',    icon: 'fas fa-capsules',          color: 'green',  typeName: 'MedicineCabinet' },
      { type: 'examinationlamp', label: '检查灯',    icon: 'fas fa-lightbulb',         color: 'yellow', typeName: 'ExaminationLamp' },
      { type: 'sinkstation',    label: '洗手台',     icon: 'fas fa-faucet',            color: 'cyan',   typeName: 'SinkStation' },
      { type: 'hospitalmonitor', label: '监护仪',    icon: 'fas fa-heartbeat',         color: 'green',  typeName: 'HospitalMonitor' },
      { type: 'defibrillator',  label: '除颤仪',     icon: 'fas fa-bolt',              color: 'red',    typeName: 'Defibrillator' },
    ]
  },
  {
    name: '农业设施', items: [
      { type: 'greenhouse',     label: '温室大棚',   icon: 'fas fa-warehouse',         color: 'green',  typeName: 'Greenhouse' },
      { type: 'tractor',        label: '拖拉机',     icon: 'fas fa-tractor',           color: 'red',    typeName: 'Tractor' },
      { type: 'harvester',      label: '收割机',     icon: 'fas fa-tractor',           color: 'orange', typeName: 'Harvester' },
      { type: 'irrigationpivot', label: '喷灌机',    icon: 'fas fa-water',             color: 'blue',   typeName: 'IrrigationPivot' },
      { type: 'grainsilo',      label: '粮仓',       icon: 'fas fa-warehouse',         color: 'brown',  typeName: 'GrainSilo' },
      { type: 'sprayer',        label: '喷雾器',     icon: 'fas fa-cloud-rain',        color: 'green',  typeName: 'Sprayer' },
      { type: 'plantingbed',    label: '种植床',     icon: 'fas fa-seedling',          color: 'green',  typeName: 'PlantingBed' },
      { type: 'windbreak',      label: '防风网',     icon: 'fas fa-wind',              color: 'gray',   typeName: 'Windbreak' },
      { type: 'weatherstation', label: '气象站',     icon: 'fas fa-cloud-sun',         color: 'cyan',   typeName: 'WeatherStation' },
      { type: 'chickencoop',    label: '鸡舍',       icon: 'fas fa-egg',               color: 'brown',  typeName: 'ChickenCoop' },
      { type: 'watertrough',    label: '饮水槽',     icon: 'fas fa-tint',              color: 'blue',   typeName: 'WaterTrough' },
      { type: 'seedlingtray',   label: '育苗盘',     icon: 'fas fa-th-large',          color: 'green',  typeName: 'SeedlingTray' },
    ]
  },
  {
    name: '水利水务', items: [
      { type: 'watertower',     label: '水塔',       icon: 'fas fa-tint',              color: 'blue',   typeName: 'WaterTower' },
      { type: 'treatmentbasin', label: '处理池',     icon: 'fas fa-swimming-pool',     color: 'cyan',   typeName: 'TreatmentBasin' },
      { type: 'pipegallery',    label: '管廊',       icon: 'fas fa-grip-lines',        color: 'gray',   typeName: 'PipeGallery' },
      { type: 'sluicegate',     label: '闸门',       icon: 'fas fa-arrows-alt-v',      color: 'gray',   typeName: 'SluiceGate' },
      { type: 'pumpstation',    label: '泵站',       icon: 'fas fa-building',          color: 'blue',   typeName: 'PumpStation' },
      { type: 'sedimentationtank', label: '沉淀池',  icon: 'fas fa-circle',            color: 'cyan',   typeName: 'SedimentationTank' },
      { type: 'aerationtank',   label: '曝气池',     icon: 'fas fa-water',             color: 'green',  typeName: 'AerationTank' },
      { type: 'overflowweir',   label: '溢流堰',     icon: 'fas fa-minus',             color: 'gray',   typeName: 'OverflowWeir' },
      { type: 'intaketower',    label: '取水塔',     icon: 'fas fa-chess-rook',        color: 'blue',   typeName: 'IntakeTower' },
      { type: 'flowchannel',    label: '水渠',       icon: 'fas fa-arrows-alt-h',      color: 'cyan',   typeName: 'FlowChannel' },
      { type: 'manholecover',   label: '窨井盖',     icon: 'fas fa-circle',            color: 'gray',   typeName: 'ManholeCover' },
      { type: 'drainpipe',      label: '排水管',     icon: 'fas fa-tint-slash',        color: 'gray',   typeName: 'DrainPipe' },
    ]
  },
  {
    name: '运动休闲', items: [
      { type: 'basketballhoop', label: '篮球架',     icon: 'fas fa-basketball-ball',   color: 'orange', typeName: 'BasketballHoop' },
      { type: 'soccergoal',     label: '足球门',     icon: 'fas fa-futbol',            color: 'cyan',   typeName: 'SoccerGoal' },
      { type: 'tenniscourt',    label: '网球场',     icon: 'fas fa-baseball-ball',     color: 'green',  typeName: 'TennisCourt' },
      { type: 'swimmingpool',   label: '游泳池',     icon: 'fas fa-swimming-pool',     color: 'blue',   typeName: 'SwimmingPool' },
      { type: 'slide',          label: '滑梯',       icon: 'fas fa-child',             color: 'yellow', typeName: 'Slide' },
      { type: 'swingset',       label: '秋千',       icon: 'fas fa-child',             color: 'orange', typeName: 'SwingSet' },
      { type: 'exercisebike',   label: '健身车',     icon: 'fas fa-bicycle',           color: 'purple', typeName: 'ExerciseBike' },
      { type: 'treadmill',      label: '跑步机',     icon: 'fas fa-running',           color: 'gray',   typeName: 'Treadmill' },
      { type: 'stadium',        label: '体育场',     icon: 'fas fa-flag-checkered',    color: 'green',  typeName: 'Stadium' },
      { type: 'scoreboard',     label: '记分牌',     icon: 'fas fa-clock',             color: 'red',    typeName: 'Scoreboard' },
      { type: 'bleacher',       label: '看台座位',   icon: 'fas fa-th-list',           color: 'blue',   typeName: 'Bleacher' },
      { type: 'pingpongtable',  label: '乒乓球台',   icon: 'fas fa-table-tennis',      color: 'green',  typeName: 'PingPongTable' },
    ]
  },
  {
    name: '实验室', items: [
      { type: 'labbench',       label: '实验台',     icon: 'fas fa-flask',             color: 'gray',   typeName: 'LabBench' },
      { type: 'fumehood',       label: '通风橱',     icon: 'fas fa-wind',              color: 'cyan',   typeName: 'FumeHood' },
      { type: 'microscope',     label: '显微镜',     icon: 'fas fa-microscope',        color: 'purple', typeName: 'Microscope' },
      { type: 'centrifuge',     label: '离心机',     icon: 'fas fa-spinner',           color: 'blue',   typeName: 'Centrifuge' },
      { type: 'incubator',      label: '培养箱',     icon: 'fas fa-box',               color: 'orange', typeName: 'Incubator' },
      { type: 'autoclave',      label: '高压灭菌器', icon: 'fas fa-temperature-high',  color: 'red',    typeName: 'Autoclave' },
      { type: 'reagentshelf',   label: '试剂架',     icon: 'fas fa-th-list',           color: 'brown',  typeName: 'ReagentShelf' },
      { type: 'glovebox',       label: '手套箱',     icon: 'fas fa-box-open',          color: 'yellow', typeName: 'GloveBox' },
      { type: 'spectrometer',   label: '光谱仪',     icon: 'fas fa-chart-line',        color: 'purple', typeName: 'Spectrometer' },
      { type: 'labscale',       label: '天平',       icon: 'fas fa-balance-scale',     color: 'gray',   typeName: 'LabScale' },
    ]
  },
  {
    name: '工业扩展', items: [
      { type: 'distillationtower', label: '精馏塔',     icon: 'fas fa-industry',          color: 'blue',   typeName: 'DistillationTower' },
      { type: 'absorbercolumn',    label: '吸收塔',     icon: 'fas fa-stream',            color: 'cyan',   typeName: 'AbsorberColumn' },
      { type: 'platefilter',       label: '板框过滤机', icon: 'fas fa-filter',            color: 'green',  typeName: 'PlateFilter' },
      { type: 'bagfilter',         label: '袋式过滤器', icon: 'fas fa-wind',              color: 'gray',   typeName: 'BagFilter' },
      { type: 'dosingunit',        label: '加药装置',   icon: 'fas fa-vial',              color: 'purple', typeName: 'DosingUnit' },
      { type: 'skidunit',          label: '撬装单元',   icon: 'fas fa-border-all',        color: 'orange', typeName: 'SkidUnit' },
      { type: 'cipstation',        label: 'CIP清洗站',  icon: 'fas fa-shower',            color: 'blue',   typeName: 'CIPStation' },
      { type: 'airblower',         label: '罗茨风机',   icon: 'fas fa-wind',              color: 'cyan',   typeName: 'AirBlower' },
      { type: 'vacuumunit',        label: '真空机组',   icon: 'fas fa-compress',          color: 'gray',   typeName: 'VacuumUnit' },
      { type: 'chillerunit',       label: '冷水机组',   icon: 'fas fa-snowflake',         color: 'cyan',   typeName: 'ChillerUnit' },
      { type: 'coolingcoil',       label: '冷却盘管',   icon: 'fas fa-wave-square',       color: 'blue',   typeName: 'CoolingCoil' },
      { type: 'weighhopper',       label: '称重料斗',   icon: 'fas fa-balance-scale',     color: 'brown',  typeName: 'WeighHopper' },
      { type: 'magneticseparator', label: '磁选机',     icon: 'fas fa-magnet',            color: 'purple', typeName: 'MagneticSeparator' },
      { type: 'vibratingscreen',   label: '振动筛',     icon: 'fas fa-grip-lines',        color: 'orange', typeName: 'VibratingScreen' },
      { type: 'rollermill',        label: '辊压机',     icon: 'fas fa-circle-notch',      color: 'gray',   typeName: 'RollerMill' },
      { type: 'controlcabinet',    label: '控制柜',     icon: 'fas fa-server',            color: 'gray',   typeName: 'ControlCabinet' },
      { type: 'mccpanel',          label: 'MCC柜',      icon: 'fas fa-bolt',              color: 'red',    typeName: 'MCCPanel' },
      { type: 'analyzerhouse',     label: '分析小屋',   icon: 'fas fa-chart-line',        color: 'blue',   typeName: 'AnalyzerHouse' },
      { type: 'pipebridge',        label: '管廊桥架',   icon: 'fas fa-grip-lines',        color: 'gray',   typeName: 'PipeBridge' },
      { type: 'maintenanceplatform', label: '检修平台', icon: 'fas fa-th-large',          color: 'yellow', typeName: 'MaintenancePlatform' },
    ]
  },
  {
    name: '灯光', items: [
      { type: 'point_light',  label: '点光源',  icon: 'fas fa-lightbulb',  color: 'orange', typeName: 'PointLight' },
      { type: 'spot_light',   label: '聚光灯',  icon: 'fas fa-satellite',  color: 'yellow', typeName: 'SpotLight' },
      { type: 'ambient_light',   label: '环境光',  icon: 'fas fa-sun',  color: 'yellow', typeName: 'AmbientLight' },
      { type: 'directional_light',   label: '平行光',  icon: 'fas fa-location-arrow',  color: 'cyan', typeName: 'DirectionalLight' },
    ]
  },
  {
    name: '视觉特效', items: [
      // = 光效类 =
      { type: 'scanLine',       label: '扫描线',     icon: 'fas fa-wave-square',    color: 'cyan',   typeName: 'ScanLine' },
      { type: 'flyLine',        label: '飞线',       icon: 'fas fa-bezier-curve',   color: 'cyan',   typeName: 'FlyLine' },
      { type: 'glowSphere',     label: '辉光球',     icon: 'fas fa-circle',         color: 'green',  typeName: 'GlowSphere' },
      { type: 'haloRing',       label: '光晕环',     icon: 'fas fa-dot-circle',     color: 'green',  typeName: 'HaloRing' },
      { type: 'ringScan',       label: '雷达扫描',   icon: 'fas fa-sync-alt',       color: 'cyan',   typeName: 'RingScan' },
      { type: 'energyWall',     label: '光墙',       icon: 'fas fa-shield-alt',     color: 'cyan',   typeName: 'EnergyWall' },
      { type: 'energyShield',   label: '能量护盾',   icon: 'fas fa-atom',           color: 'purple', typeName: 'EnergyShield' },
      { type: 'magicGlow',      label: '魔法光晕',   icon: 'fas fa-magic',          color: 'purple', typeName: 'MagicGlow' },
      { type: 'starFX',         label: '星星',       icon: 'fas fa-star',           color: 'yellow', typeName: 'StarFX' },
      // = 粒子流动 =
      { type: 'particleCloud',  label: '粒子云',     icon: 'fas fa-cloud',          color: 'cyan',   typeName: 'ParticleCloud' },
      { type: 'fireFX',         label: '火焰',       icon: 'fas fa-fire',           color: 'red',    typeName: 'FireFX' },
      { type: 'smokeFX',        label: '烟雾',       icon: 'fas fa-smog',           color: 'gray',   typeName: 'SmokeFX' },
      { type: 'dustFX',         label: '灰尘',       icon: 'fas fa-feather-alt',    color: 'brown',  typeName: 'DustFX' },
      { type: 'fireflyFX',      label: '萤火虫',     icon: 'fas fa-certificate',    color: 'green',  typeName: 'FireflyFX' },
      // = 天气自然 =
      { type: 'snowFX',         label: '雪花',       icon: 'fas fa-snowflake',      color: 'cyan',   typeName: 'SnowFX' },
      { type: 'rainFX',         label: '雨滴',       icon: 'fas fa-tint',           color: 'blue',   typeName: 'RainFX' },
      { type: 'bubbleFX',       label: '气泡',       icon: 'fas fa-film',           color: 'cyan',   typeName: 'BubbleFX' },
      { type: 'leafFX',         label: '落叶',       icon: 'fas fa-leaf',           color: 'green',  typeName: 'LeafFX' },
      // = 爆发冲击 =
      { type: 'sparkFX',        label: '火花',       icon: 'fas fa-bolt',           color: 'yellow', typeName: 'SparkFX' },
      { type: 'electricArc',    label: '电弧',       icon: 'fas fa-poo-storm',      color: 'yellow', typeName: 'ElectricArc' },
      { type: 'explosionFX',    label: '爆炸',       icon: 'fas fa-bomb',           color: 'red',    typeName: 'ExplosionFX' },
      { type: 'splashFX',       label: '水花',       icon: 'fas fa-water',          color: 'blue',   typeName: 'SplashFX' },
      // = 能量场景 =
      { type: 'pulseColumn',    label: '脉冲柱',     icon: 'fas fa-chart-bar',      color: 'green',  typeName: 'PulseColumn' },
      { type: 'heatmapCloud',   label: '热力云',     icon: 'fas fa-cloud-sun',      color: 'red',    typeName: 'HeatmapCloud' },
      { type: 'flowRibbon',     label: '流动光带',   icon: 'fas fa-wave-square',    color: 'orange', typeName: 'FlowRibbon' },
      { type: 'gridFloor',      label: '网格地面',   icon: 'fas fa-border-all',     color: 'cyan',   typeName: 'GridFloor' },
      { type: 'gradientBlock',  label: '渐变区块',   icon: 'fas fa-cube',           color: 'blue',   typeName: 'GradientBlock' },
    ]
  },
  {
    name: '基础组件', items: [
      { type: 'gltf',        label: '3D模型',   description: '导入模型',       icon: 'fas fa-spinner',  color: 'blue',   typeName: 'GLTFModel' },
      { type: 'text3d',      label: '3D标签',   description: '可交互标注',     icon: 'fas fa-tag',      color: 'pink', defaultColor: '#ffffff', typeName: 'Text3D', textContent: '3D标签' },
      { type: 'image3d',     label: '3D图片',   description: '插入图片',       icon: 'far fa-image',    color: 'cyan',   typeName: 'Image3D' },
      { type: 'textPlain3d', label: '3D文字',   description: '生成文字模型',   icon: 'fas fa-font',     color: 'purple', defaultColor: '#ffffff', typeName: 'Text3D', textContent: '3D文字' },
      { type: 'dataText',    label: '实时数据', description: '绑定数据文本',   icon: 'fas fa-broadcast-tower', color: 'yellow', defaultColor: '#ffffff', typeName: 'DataText' },
      { type: 'uiLabel',     label: 'UI标签',   description: '屏幕固定标签',   icon: 'fas fa-bullseye', color: 'green',  typeName: 'UILabel' },
      { type: 'uiImage',     label: 'UI图片',   description: '屏幕固定图片',   icon: 'far fa-image',    color: 'cyan',   typeName: 'UIImage' },
      { type: 'video3d',     label: '3D视频',   description: '插入视频',       icon: 'fas fa-video',    color: 'red',    typeName: 'Video3D' },
      { type: 'webEmbed',    label: '网页嵌入', description: '嵌入网页',       icon: 'fas fa-globe',    color: 'blue',   typeName: 'WebEmbed' },
    ]
  },
  {
    name: '外部模型', items: [
      { type: 'gltf',       label: '加载模型', icon: 'fas fa-file-import',   color: 'blue',   typeName: 'GLTFModel' },
    ]
  },
  // ===== 数字孪生高端特效（新增） =====
  {
    name: '数字孪生特效', items: [
      { type: 'dtBuilding',        label: '发光建筑',   icon: 'fas fa-building',        color: 'cyan',   typeName: 'DTBuilding' },
      { type: 'dtBuildingTall',    label: '高层幕墙建筑', icon: 'fas fa-city',         color: 'cyan',   typeName: 'DTBuildingTall' },
      { type: 'dtBuildingWide',    label: '扁平厂房',   icon: 'fas fa-industry',        color: 'blue',   typeName: 'DTBuildingWide' },
      { type: 'dtBuildingComplex', label: '复合建筑群', icon: 'fas fa-building',        color: 'cyan',   typeName: 'DTBuildingComplex' },
      { type: 'dtGroundGrid',      label: '发光地面网格', icon: 'fas fa-th',            color: 'cyan',   typeName: 'DTGroundGrid' },
      { type: 'dtScanRing',        label: '扫描环',     icon: 'fas fa-bullseye',       color: 'cyan',   typeName: 'DTScanRing' },
      { type: 'dtScanLine',        label: '扫描光束',   icon: 'fas fa-satellite-dish',  color: 'cyan',   typeName: 'DTScanLine' },
      { type: 'dtDataPanel',       label: '数据面板',   icon: 'fas fa-desktop',         color: 'blue',   typeName: 'DTDataPanel' },
      { type: 'dtFlowLine',        label: '流光数据线', icon: 'fas fa-wave-square',    color: 'cyan',   typeName: 'DTFlowLine' },
      { type: 'dtParticleField',   label: '氛围粒子',   icon: 'fas fa-stars',           color: 'cyan',   typeName: 'DTParticleField' },
      { type: 'dtHologram',        label: '全息投影',   icon: 'fas fa-eye',            color: 'purple', typeName: 'DTHologram' },
      { type: 'dtRoad',            label: '智慧道路',   icon: 'fas fa-road',            color: 'gray',   typeName: 'DTRoad' },
      { type: 'dtPark',            label: '智慧绿地',   icon: 'fas fa-tree',            color: 'green',  typeName: 'DTPark' },
      { type: 'dtStreetLightDT',   label: '智慧路灯',   icon: 'fas fa-lightbulb',      color: 'yellow', typeName: 'DTStreetLight' },
      { type: 'dtBaseStation',     label: '通信基站',   icon: 'fas fa-broadcast-tower', color: 'orange', typeName: 'DTBaseStation' },
    ]
  },
]

// 将 hex 颜色转换为 rgba 格式（支持 8 位 hex 透明度）
function hexToRgba(hex) {
  if (!hex || typeof hex !== 'string') return 'rgba(0,0,0,0)'
  var s = hex.trim()
  if (s === 'transparent') return 'rgba(0,0,0,0)'
  if (s[0] !== '#') return 'rgba(0,0,0,0)'
  // #RRGGBBAA
  if (/^#[0-9a-fA-F]{8}$/.test(s)) {
    var r = parseInt(s.slice(1, 3), 16)
    var g = parseInt(s.slice(3, 5), 16)
    var b = parseInt(s.slice(5, 7), 16)
    var a = parseInt(s.slice(7, 9), 16) / 255
    return 'rgba(' + r + ',' + g + ',' + b + ',' + a + ')'
  }
  // #RRGGBB
  if (/^#[0-9a-fA-F]{6}$/.test(s)) {
    var r = parseInt(s.slice(1, 3), 16)
    var g = parseInt(s.slice(3, 5), 16)
    var b = parseInt(s.slice(5, 7), 16)
    return 'rgba(' + r + ',' + g + ',' + b + ',1)'
  }
  // #RGB
  if (/^#[0-9a-fA-F]{3}$/.test(s)) {
    var r = parseInt(s[1] + s[1], 16)
    var g = parseInt(s[2] + s[2], 16)
    var b = parseInt(s[3] + s[3], 16)
    return 'rgba(' + r + ',' + g + ',' + b + ',1)'
  }
  return 'rgba(0,0,0,0)'
}

// ===== 文字 Sprite 工厂 =====
function colorToCss(value, fallback) {
  fallback = fallback || '#ffffff'
  if (typeof value === 'number') return '#' + value.toString(16).padStart(6, '0')
  if (typeof value === 'string' && value) {
    if (value.charAt(0) === '#') {
      if (/^#[0-9a-fA-F]{3}$/.test(value)) {
        return '#' + value[1] + value[1] + value[2] + value[2] + value[3] + value[3]
      }
      if (/^#[0-9a-fA-F]{6}/.test(value)) return value.slice(0, 7)
    }
  }
  return fallback
}

function hexAlphaToOpacity(value, fallback) {
  fallback = fallback === undefined ? 1 : fallback
  if (typeof value !== 'string' || !/^#[0-9a-fA-F]{8}$/.test(value)) return fallback
  return parseInt(value.slice(7, 9), 16) / 255
}

function roundedRectPath(ctx, x, y, w, h, r) {
  r = Math.max(0, Math.min(r, w / 2, h / 2))
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.lineTo(x + w - r, y)
  ctx.quadraticCurveTo(x + w, y, x + w, y + r)
  ctx.lineTo(x + w, y + h - r)
  ctx.quadraticCurveTo(x + w, y + h, x + w - r, y + h)
  ctx.lineTo(x + r, y + h)
  ctx.quadraticCurveTo(x, y + h, x, y + h - r)
  ctx.lineTo(x, y + r)
  ctx.quadraticCurveTo(x, y, x + r, y)
  ctx.closePath()
}

function buildTextSpriteCanvas(text, color, fontSize, bgColor, options) {
  options = options || {}
  var canvas = document.createElement('canvas')
  var ctx = canvas.getContext('2d', { alpha: true })

  fontSize = fontSize || 16
  var drawScale = 3
  var drawFontSize = fontSize * drawScale
  text = text === undefined || text === null ? '' : String(text)
  text = text.replace(/\r\n/g, '\n').replace(/\\n/g, '\n')
  var paddingX = options.paddingX !== undefined ? options.paddingX : drawFontSize * 0.5
  var paddingY = options.paddingY !== undefined ? options.paddingY : drawFontSize * 0.32
  var borderWidth = options.showBorder ? (options.borderWidth !== undefined ? options.borderWidth : 1) : 0
  var radius = options.radius !== undefined ? options.radius : Math.max(3, drawFontSize * 0.08)
  var pixelRatio = Math.min(window.devicePixelRatio || 2, 2)

  var fontWeight = options.fontWeight || 500
  var fontFamily = options.fontFamily && options.fontFamily !== '系统默认'
    ? options.fontFamily
    : '"Microsoft YaHei", "PingFang SC", Arial, sans-serif'
  ctx.font = fontWeight + ' ' + drawFontSize + 'px ' + fontFamily
  var lines = text.split('\n')
  if (!lines.length) lines = ['']
  var lineHeight = options.lineHeight !== undefined ? options.lineHeight * drawScale : drawFontSize * 1.2
  var textWidths = lines.map(function(line) {
    return ctx.measureText(line || ' ').width
  })
  var textWidth = textWidths.length ? Math.max.apply(null, textWidths) : 0
  var textHeight = lineHeight * lines.length
  var width = Math.ceil(textWidth + paddingX * 2 + borderWidth * 2)
  var height = Math.ceil(textHeight + paddingY * 2 + borderWidth * 2)

  canvas.width = width * pixelRatio
  canvas.height = height * pixelRatio
  canvas.style.width = width + 'px'
  canvas.style.height = height + 'px'

  ctx.scale(pixelRatio, pixelRatio)
  ctx.imageSmoothingEnabled = true
  ctx.imageSmoothingQuality = 'high'
  ctx.clearRect(0, 0, width, height)

  if (!options.noBackground) {
    var bgOpacity = options.backgroundOpacity !== undefined ? options.backgroundOpacity : hexAlphaToOpacity(bgColor, 0.2)
    var bgCss = colorToCss(bgColor, '#000000') + Math.round(Math.max(0, Math.min(1, bgOpacity)) * 255).toString(16).padStart(2, '0')
    var panelX = borderWidth / 2
    var panelY = borderWidth / 2
    var panelW = width - borderWidth
    var panelH = height - borderWidth

    ctx.save()
    ctx.shadowColor = 'transparent'
    ctx.shadowBlur = 0
    ctx.shadowOffsetY = 0
    roundedRectPath(ctx, panelX, panelY, panelW, panelH, radius)
    ctx.fillStyle = hexToRgba(bgCss)
    ctx.fill()
    ctx.restore()

    if (borderWidth > 0) {
      ctx.lineWidth = borderWidth
      ctx.strokeStyle = hexToRgba(colorToCss(options.borderColor, '#ffffff') + 'b3')
      roundedRectPath(ctx, panelX, panelY, panelW, panelH, radius)
      ctx.stroke()
    }
  }

  ctx.font = fontWeight + ' ' + drawFontSize + 'px ' + fontFamily
  ctx.textAlign = 'center'
  ctx.textBaseline = 'alphabetic'
  ctx.fillStyle = colorToCss(color, '#ffffff')
  ctx.shadowColor = 'transparent'
  ctx.shadowBlur = 0
  var textMetrics = ctx.measureText(lines[0] || ' ')
  var ascent = textMetrics.actualBoundingBoxAscent !== undefined ? textMetrics.actualBoundingBoxAscent : drawFontSize * 0.8
  var descent = textMetrics.actualBoundingBoxDescent !== undefined ? textMetrics.actualBoundingBoxDescent : drawFontSize * 0.2
  var firstLineY = (height - textHeight) / 2 + ascent + Math.max(0, (lineHeight - ascent - descent) / 2)
  for (var i = 0; i < lines.length; i++) {
    ctx.fillText(lines[i], width / 2, firstLineY + i * lineHeight)
  }

  return canvas
}

// ===== 文字 Sprite 工厂 =====
function createTextSprite(text, color, fontSize, bgColor, options) {
  options = options || {}
  var canvas = buildTextSpriteCanvas(text, color, fontSize, bgColor, options)

  var texture = new THREE.CanvasTexture(canvas)
  texture.needsUpdate = true
  texture.colorSpace = THREE.SRGBColorSpace
  texture.minFilter = THREE.LinearFilter
  texture.magFilter = THREE.LinearFilter
  texture.anisotropy = 16
  texture.generateMipmaps = false

  var spriteMat = new THREE.SpriteMaterial({
    map: texture,
    color: 0xffffff,
    transparent: true,
    opacity: options.opacity !== undefined ? options.opacity : 1,
    depthTest: options.depthTest === true,
    depthWrite: false,
    alphaTest: 0.0001,
    sizeAttenuation: options.fixedSize ? false : true,
    side: THREE.DoubleSide
  })

  var sprite = new THREE.Sprite(spriteMat)
  var spriteHeight = options.fixedSize ? (options.fixedScale || 0.055) : (options.worldScale || 0.28)
  sprite.scale.set((canvas.width / canvas.height) * spriteHeight, spriteHeight, 1)
  sprite.userData.isTextSprite = true
  sprite.userData.faceCamera = options.faceCamera !== false
  sprite.userData.baseTextScale = sprite.scale.clone()
  sprite.userData.textData = {
    text: text,
    color: color,
    fontSize: fontSize,
    bgColor: bgColor || '#00000000',
    options: Object.assign({}, options)
  }

  return sprite
}

// ===== 翠鸟风格 3D 文字工厂 =====
// 优先使用真正的3D几何体文字（需要字体文件），失败则 fallback 到 Canvas 贴图方案
// 返回 { mesh, canvas, texture, isTrue3D }
function createText3DMesh(text, color, fontSize, bgColor, options) {
  options = options || {}
  // 优先尝试真正的3D几何体文字
  if (options.true3D !== false && getTrue3DFont(options.fontFamily)) {
    var trueMesh = createTrue3DTextMesh(text, color, fontSize, options)
    if (trueMesh) {
      return { mesh: trueMesh, canvas: null, texture: null, isTrue3D: true }
    }
  }
  // Fallback: Canvas 贴图方案（无背景）
  var drawOptions = Object.assign({}, options)
  var canvas = buildTextSpriteCanvas(text, color, fontSize, bgColor, drawOptions)
  var texture = new THREE.CanvasTexture(canvas)
  texture.needsUpdate = true
  texture.colorSpace = THREE.SRGBColorSpace
  texture.minFilter = THREE.LinearFilter
  texture.magFilter = THREE.LinearFilter
  texture.anisotropy = 16
  texture.generateMipmaps = false

  // 世界空间尺寸: 基于 fontSize 计算，16号字 ≈ 0.8 单位高
  var spriteHeight = (fontSize / 16) * 0.8
  var w = (canvas.width / canvas.height) * spriteHeight
  var h = spriteHeight
  var depth = spriteHeight * 0.28

  // BoxGeometry (6个材质组: +X,-X,+Y,-Y,+Z,-Z)
  var geo = new THREE.BoxGeometry(w, h, depth, 1, 1, 1)

  // 文字颜色提取作为侧面颜色
  var bgOpacity = options.backgroundOpacity !== undefined ? Number(options.backgroundOpacity) : hexAlphaToOpacity(bgColor, 0)
  if (!isFinite(bgOpacity)) bgOpacity = 0
  var sideHex = bgOpacity > 0 ? colorToCss(bgColor, '#000000') : colorToCss(color, '#ffffff')
  var sideInt = parseInt(sideHex.replace('#', ''), 16)
  var sideColor = new THREE.Color(sideInt)
  var sideOpacity = bgOpacity > 0 ? bgOpacity : (options.opacity !== undefined ? options.opacity : 1)

  var matFront = new THREE.MeshBasicMaterial({
    map: texture,
    color: 0xffffff,
    transparent: true,
    opacity: options.opacity !== undefined ? options.opacity : 1,
    alphaTest: 0.01,
    depthTest: true,
    depthWrite: true,
    toneMapped: false,
    side: THREE.FrontSide
  })

  var matSide = new THREE.MeshBasicMaterial({
    color: sideColor,
    transparent: true,
    opacity: sideOpacity,
    depthTest: true,
    depthWrite: true,
    toneMapped: false,
    side: THREE.FrontSide
  })

  var matBack = new THREE.MeshBasicMaterial({
    color: sideColor,
    transparent: true,
    opacity: sideOpacity,
    depthTest: true,
    depthWrite: true,
    toneMapped: false,
    side: THREE.FrontSide
  })

  var mesh = new THREE.Mesh(geo, [matSide, matSide, matSide, matSide, matFront, matBack])
  mesh.castShadow = true
  mesh.receiveShadow = true

  mesh.userData.is3DText = true
  mesh.userData.isTextSprite = true
  mesh.userData.faceCamera = false
  mesh.userData.textData = {
    text: text,
    color: color,
    fontSize: fontSize,
    bgColor: bgColor || '#1f1f1f',
    options: Object.assign({}, options)
  }
  mesh.name = '3DText_' + (text || '').slice(0, 10)

  return { mesh: mesh, canvas: canvas, texture: texture }
}

function updateTextSprite(sprite, text, color, fontSize, bgColor, options) {
  if (!sprite || !sprite.userData.isTextSprite) return

  fontSize = fontSize || sprite.userData.textData.fontSize
  color = color !== undefined ? color : sprite.userData.textData.color
  bgColor = bgColor !== undefined ? bgColor : sprite.userData.textData.bgColor
  options = Object.assign({}, sprite.userData.textData.options || {}, options || {})
  var canvas = buildTextSpriteCanvas(text, color, fontSize, bgColor, options)

  var texture = new THREE.CanvasTexture(canvas)
  texture.needsUpdate = true
  texture.colorSpace = THREE.SRGBColorSpace
  texture.minFilter = THREE.LinearFilter
  texture.magFilter = THREE.LinearFilter
  texture.anisotropy = 16
  texture.generateMipmaps = false

  sprite.material.map.dispose()
  sprite.material.map = texture
  sprite.material.color.set(0xffffff)
  sprite.material.opacity = options.opacity !== undefined ? options.opacity : 1
  sprite.material.depthTest = options.depthTest === true
  sprite.material.depthWrite = false
  sprite.material.sizeAttenuation = options.fixedSize ? false : true
  sprite.material.needsUpdate = true
  sprite.material.alphaTest = 0.0001
  sprite.material.side = THREE.DoubleSide
  sprite.material.needsUpdate = true

  var spriteHeight = options.fixedSize ? (options.fixedScale || 0.055) : (options.worldScale || 0.28)
  sprite.scale.set((canvas.width / canvas.height) * spriteHeight, spriteHeight, 1)
  sprite.userData.baseTextScale = sprite.scale.clone()
  sprite.userData.faceCamera = options.faceCamera !== false

  sprite.userData.textData = {
    text: text,
    color: color,
    fontSize: fontSize,
    bgColor: bgColor,
    options: Object.assign({}, options)
  }
}

// ===== 翠鸟风格 3D 文字 Badge 更新 =====
function updateText3DMesh(mesh, text, color, fontSize, bgColor, options) {
  if (!mesh || !mesh.userData.is3DText) return

  options = options || {}
  fontSize = fontSize || (mesh.userData.textData && mesh.userData.textData.fontSize)
  color = color !== undefined ? color : (mesh.userData.textData && mesh.userData.textData.color)
  bgColor = bgColor !== undefined ? bgColor : (mesh.userData.textData && mesh.userData.textData.bgColor)
  options = Object.assign({}, (mesh.userData.textData && mesh.userData.textData.options) || {}, options || {})
  var shouldUseCanvasText = options.true3D === false

  // Canvas fallback 升级：如果字体现在已就绪，将 BoxGeometry+贴图升级为 true 3D
  if (!shouldUseCanvasText && !mesh.userData.isTrue3DText && getTrue3DFont()) {
    var upgradeMesh = createTrue3DTextMesh(text, color, fontSize, options)
    if (upgradeMesh) {
      // 销毁旧的 BoxGeometry 和全部材质
      if (mesh.geometry) mesh.geometry.dispose()
      var oldMats = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
      oldMats.forEach(function(m) {
        if (m) {
          if (m.map) m.map.dispose()
          m.dispose()
        }
      })
      // 替换几何体和材质
      mesh.geometry = upgradeMesh.geometry
      mesh.material = upgradeMesh.material
      // 设置 true 3D 标记
      mesh.userData.isTrue3DText = true
      mesh.userData.isTextSprite = true
      mesh.userData.faceCamera = false
      mesh.userData.textData = {
        text: text,
        color: color,
        fontSize: fontSize,
        bgColor: '#00000000',
        options: Object.assign({}, options)
      }
      return
    }
  }

  // 真正的3D几何体文字：原地重建几何体，保持 mesh 引用
  if (mesh.userData.isTrue3DText && shouldUseCanvasText) {
    var fallbackResult = createText3DMesh(text, color, fontSize, bgColor, Object.assign({}, options, { true3D: false }))
    if (!fallbackResult || !fallbackResult.mesh) return
    if (mesh.geometry) mesh.geometry.dispose()
    var currentMats = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
    currentMats.forEach(function(m) {
      if (m) {
        if (m.map) m.map.dispose()
        m.dispose()
      }
    })
    mesh.geometry = fallbackResult.mesh.geometry
    mesh.material = fallbackResult.mesh.material
    mesh.userData.isTrue3DText = false
    mesh.userData.is3DText = true
    mesh.userData.isTextSprite = true
    mesh.userData.faceCamera = false
    mesh.userData.textData = {
      text: text,
      color: color,
      fontSize: fontSize,
      bgColor: bgColor,
      options: Object.assign({}, options)
    }
    return
  }

  if (mesh.userData.isTrue3DText) {
    var tempMesh = createTrue3DTextMesh(text, color, fontSize, options)

    if (tempMesh) {
      // 重建几何体成功：替换几何体，丢弃临时 mesh 的材质（使用原 mesh 材质）
      if (mesh.geometry) {
        mesh.geometry.dispose()
      }
      mesh.geometry = tempMesh.geometry
      var tempMats = Array.isArray(tempMesh.material) ? tempMesh.material : [tempMesh.material]
      tempMats.forEach(function (m) { m && m.dispose() })
    } else {
      // 重建失败（字体未就绪）：仅更新材质颜色，保留现有几何体
      console.warn('[3DText] updateText3DMesh: createTrue3DTextMesh returned null, updating material color only. fontFamily=' + (options && options.fontFamily))
    }

    // 无论几何体是否重建，都更新材质颜色（关键：使前景色修改生效）
    var col = new THREE.Color(color || '#ffffff')
    var mats = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
    mats.forEach(function (m) {
      if (m) {
        m.color.copy(col)
        m.emissive.copy(col.clone().multiplyScalar(0.25))
        m.emissiveIntensity = 0.6
        m.toneMapped = false
        if (options && options.opacity !== undefined) {
          m.transparent = options.opacity < 1
          m.opacity = options.opacity
        }
        m.needsUpdate = true
      }
    })

    mesh.userData.textData = {
      text: text,
      color: color,
      fontSize: fontSize,
      bgColor: bgColor,
      options: Object.assign({}, options)
    }
    return
  }

  // 旧 Canvas 贴图方案（兼容代码）
  var drawOptions = Object.assign({}, options)
  var canvas = buildTextSpriteCanvas(text, color, fontSize, bgColor, drawOptions)
  var newTex = new THREE.CanvasTexture(canvas)
  newTex.needsUpdate = true
  newTex.colorSpace = THREE.SRGBColorSpace
  newTex.minFilter = THREE.LinearFilter
  newTex.magFilter = THREE.LinearFilter
  newTex.anisotropy = 16
  newTex.generateMipmaps = false

  var matList = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
  var matFront = matList[4] || matList[0]
  if (matFront && !matFront.isMeshBasicMaterial) {
    var oldFront = matFront
    matFront = new THREE.MeshBasicMaterial({
      map: oldFront.map || null,
      color: 0xffffff,
      transparent: true,
      opacity: oldFront.opacity !== undefined ? oldFront.opacity : 1,
      alphaTest: oldFront.alphaTest !== undefined ? oldFront.alphaTest : 0.01,
      depthTest: oldFront.depthTest !== undefined ? oldFront.depthTest : true,
      depthWrite: oldFront.depthWrite !== undefined ? oldFront.depthWrite : true,
      toneMapped: false,
      side: THREE.FrontSide
    })
    oldFront.dispose()
    if (Array.isArray(mesh.material)) {
      mesh.material[4] = matFront
    } else {
      mesh.material = matFront
    }
  }
  if (matFront && matFront.map) {
    matFront.map.dispose()
    matFront.map = newTex
    matFront.needsUpdate = true
  }
  if (matFront) {
    if (matFront.color && matFront.color.set) matFront.color.set(0xffffff)
    matFront.opacity = options.opacity !== undefined ? options.opacity : 1
    matFront.toneMapped = false
    matFront.needsUpdate = true
  }

  // 文字颜色变化时同步更新侧面颜色
  var updateBgOpacity = options.backgroundOpacity !== undefined ? Number(options.backgroundOpacity) : hexAlphaToOpacity(bgColor, 0)
  if (!isFinite(updateBgOpacity)) updateBgOpacity = 0
  var updateSideHex = updateBgOpacity > 0 ? colorToCss(bgColor, '#000000') : colorToCss(color, '#ffffff')
  var updateSideInt = parseInt(updateSideHex.replace('#', ''), 16)
  var sideColor = new THREE.Color(updateSideInt)

  var materials = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
  materials.forEach(function(m, index) {
    if (m && m !== matFront) {
      if (!m.isMeshBasicMaterial) {
        var oldSide = m
        m = new THREE.MeshBasicMaterial({
          color: sideColor,
          transparent: true,
          opacity: oldSide.opacity !== undefined ? oldSide.opacity : 1,
          depthTest: oldSide.depthTest !== undefined ? oldSide.depthTest : true,
          depthWrite: oldSide.depthWrite !== undefined ? oldSide.depthWrite : true,
          toneMapped: false,
          side: THREE.FrontSide
        })
        oldSide.dispose()
        if (Array.isArray(mesh.material)) mesh.material[index] = m
      }
      m.color.copy(sideColor)
      m.opacity = updateBgOpacity > 0 ? updateBgOpacity : (options.opacity !== undefined ? options.opacity : 1)
      m.toneMapped = false
      m.needsUpdate = true
    }
  })

  mesh.userData.textData = {
    text: text,
    color: color,
    fontSize: fontSize,
    bgColor: bgColor,
    options: Object.assign({}, options)
  }
}

// ===== 流动管道工厂 =====
function normalizePipeColor(value, fallback) {
  if (!value || typeof value !== 'string') return fallback
  var s = value.trim()
  if (/^#[0-9a-fA-F]{3}$/.test(s)) {
    return '#' + s[1] + s[1] + s[2] + s[2] + s[3] + s[3]
  }
  if (/^#[0-9a-fA-F]{6}$/.test(s)) return s
  if (/^#[0-9a-fA-F]{8}$/.test(s)) return s.slice(0, 7)
  return fallback
}

function pipeHexToRgb(hex) {
  var s = normalizePipeColor(hex, '#ffffff').slice(1)
  return {
    r: parseInt(s.slice(0, 2), 16),
    g: parseInt(s.slice(2, 4), 16),
    b: parseInt(s.slice(4, 6), 16)
  }
}

function pipeRgbToCss(rgb) {
  return 'rgb(' + Math.round(rgb.r) + ',' + Math.round(rgb.g) + ',' + Math.round(rgb.b) + ')'
}

function mixPipeColor(color, targetColor, amount) {
  var source = pipeHexToRgb(color)
  var target = pipeHexToRgb(targetColor)
  return pipeRgbToCss({
    r: source.r + (target.r - source.r) * amount,
    g: source.g + (target.g - source.g) * amount,
    b: source.b + (target.b - source.b) * amount
  })
}

function getFlowPipeLength(pipePoints) {
  if (!pipePoints || pipePoints.length < 2) return 1
  var length = 0
  for (var i = 1; i < pipePoints.length; i++) {
    length += pipePoints[i - 1].distanceTo(pipePoints[i])
  }
  return Math.max(length, 0.1)
}

function applyFlowPipeTextureRepeat(material, pipeLength, dashLength) {
  if (!material || !material.map) return
  var bandWorldPitch = getFlowDashLength(dashLength) * 2.55
  var repeatX = Math.max(1, pipeLength / bandWorldPitch)
  material.map.repeat.set(repeatX, 1)
  material.map.needsUpdate = true
}

function getFlowDashLength(value) {
  var length = parseFloat(value)
  if (isNaN(length)) return 3
  return Math.max(0.05, Math.min(3, length))
}

function createFlowPipeTexture(color, highlightColor, dashLength) {
  var canvas = document.createElement('canvas')
  canvas.height = 96
  var dashWidth = Math.round(96 * getFlowDashLength(dashLength))
  var gapWidth = Math.round(Math.max(36, dashWidth * 1.55))
  var pitch = dashWidth + gapWidth
  canvas.width = pitch * 6
  var ctx = canvas.getContext('2d')
  var pipeBodyColor = normalizePipeColor(color, '#f4ead8')
  var bandColor = normalizePipeColor(highlightColor, '#ff6a00')
  var bodyLight = mixPipeColor(pipeBodyColor, '#ffffff', 0.32)
  var bodyDark = mixPipeColor(pipeBodyColor, '#000000', 0.34)
  var bodyShadow = mixPipeColor(pipeBodyColor, '#000000', 0.18)
  var bandLight = mixPipeColor(bandColor, '#ffffff', 0.28)
  var bandDark = mixPipeColor(bandColor, '#000000', 0.28)

  ctx.clearRect(0, 0, canvas.width, canvas.height)

  var bodyGradient = ctx.createLinearGradient(0, 0, 0, canvas.height)
  bodyGradient.addColorStop(0, bodyLight)
  bodyGradient.addColorStop(0.16, pipeBodyColor)
  bodyGradient.addColorStop(0.5, pipeBodyColor)
  bodyGradient.addColorStop(0.84, bodyDark)
  bodyGradient.addColorStop(1, bodyLight)
  ctx.fillStyle = bodyGradient
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  ctx.fillStyle = mixPipeColor(pipeBodyColor, '#ffffff', 0.48)
  ctx.fillRect(0, 10, canvas.width, 4)
  ctx.fillStyle = bodyShadow
  ctx.fillRect(0, 78, canvas.width, 3)

  for (var x = 0; x < canvas.width; x += pitch) {
    var bandGradient = ctx.createLinearGradient(0, 0, 0, canvas.height)
    bandGradient.addColorStop(0, bandLight)
    bandGradient.addColorStop(0.18, bandColor)
    bandGradient.addColorStop(0.5, bandColor)
    bandGradient.addColorStop(0.82, bandDark)
    bandGradient.addColorStop(1, bandLight)
    ctx.fillStyle = bandGradient
    ctx.fillRect(x, 0, dashWidth, canvas.height)

    ctx.fillStyle = mixPipeColor(bandColor, '#ffffff', 0.42)
    ctx.fillRect(x + 5, 8, dashWidth - 10, 6)
    ctx.fillStyle = bandDark
    ctx.fillRect(x, 0, 3, canvas.height)
    ctx.fillRect(x + dashWidth - 3, 0, 3, canvas.height)
  }

  var texture = new THREE.CanvasTexture(canvas)
  texture.wrapS = THREE.RepeatWrapping
  texture.wrapT = THREE.RepeatWrapping
  texture.repeat.set(1, 1)
  texture.colorSpace = THREE.SRGBColorSpace
  texture.needsUpdate = true
  return texture
}

function createFlowPipeMaterial(color, highlightColor, dashLength) {
  return new THREE.MeshStandardMaterial({
    map: createFlowPipeTexture(color, highlightColor, dashLength),
    color: 0xffffff,
    metalness: 0.12,
    roughness: 0.34,
    transparent: false,
    opacity: 1
  })
}

// ===== 粒子纹理工具函数 =====
var _circleTexCache = null
function createCircleTexture(colorHex) {
  if (_circleTexCache) return _circleTexCache
  var size = 64
  var canvas = document.createElement('canvas')
  canvas.width = size; canvas.height = size
  var ctx = canvas.getContext('2d')
  ctx.clearRect(0, 0, size, size)
  var gradient = ctx.createRadialGradient(size/2, size/2, 0, size/2, size/2, size/2)
  gradient.addColorStop(0, 'rgba(255,255,255,1)')
  gradient.addColorStop(0.15, 'rgba(255,255,255,0.9)')
  gradient.addColorStop(0.4, 'rgba(255,255,255,0.4)')
  gradient.addColorStop(0.7, 'rgba(255,255,255,0.05)')
  gradient.addColorStop(1, 'rgba(255,255,255,0)')
  ctx.fillStyle = gradient; ctx.fillRect(0, 0, size, size)
  _circleTexCache = new THREE.CanvasTexture(canvas)
  _circleTexCache.needsUpdate = true
  return _circleTexCache
}

var _starTexCache = null
function createStarTexture() {
  if (_starTexCache) return _starTexCache
  var size = 64
  var canvas = document.createElement('canvas')
  canvas.width = size; canvas.height = size
  var ctx = canvas.getContext('2d')
  ctx.clearRect(0, 0, size, size)
  var cx = size/2, cy = size/2, outerR = 30, innerR = 12, spikes = 4
  ctx.beginPath()
  for (var i = 0; i < spikes * 2; i++) {
    var r = i % 2 === 0 ? outerR : innerR
    var angle = (i * Math.PI) / spikes - Math.PI / 2
    var x = cx + Math.cos(angle) * r
    var y = cy + Math.sin(angle) * r
    if (i === 0) ctx.moveTo(x, y)
    else ctx.lineTo(x, y)
  }
  ctx.closePath()
  var gradient = ctx.createRadialGradient(cx, cy, innerR * 0.5, cx, cy, outerR)
  gradient.addColorStop(0, 'rgba(255,255,240,1)')
  gradient.addColorStop(0.3, 'rgba(255,255,200,0.9)')
  gradient.addColorStop(0.7, 'rgba(255,255,150,0.3)')
  gradient.addColorStop(1, 'rgba(255,255,100,0)')
  ctx.fillStyle = gradient; ctx.fill()
  _starTexCache = new THREE.CanvasTexture(canvas)
  _starTexCache.needsUpdate = true
  return _starTexCache
}

function createRoundedPipeCurve(pipePoints, pipeRadius) {
  var curve = new THREE.CurvePath()
  if (!pipePoints || pipePoints.length < 2) return curve

  var current = pipePoints[0].clone()
  for (var i = 1; i < pipePoints.length - 1; i++) {
    var prev = pipePoints[i - 1]
    var corner = pipePoints[i]
    var next = pipePoints[i + 1]
    var prevLen = corner.distanceTo(prev)
    var nextLen = corner.distanceTo(next)
    if (prevLen <= 0.0001 || nextLen <= 0.0001) continue

    var bendRadius = Math.min(Math.max(pipeRadius * 3.2, 0.18), prevLen * 0.35, nextLen * 0.35)
    var before = corner.clone().add(prev.clone().sub(corner).normalize().multiplyScalar(bendRadius))
    var after = corner.clone().add(next.clone().sub(corner).normalize().multiplyScalar(bendRadius))

    if (current.distanceTo(before) > 0.0001) {
      curve.add(new THREE.LineCurve3(current, before))
    }
    curve.add(new THREE.QuadraticBezierCurve3(before, corner, after))
    current = after
  }

  var last = pipePoints[pipePoints.length - 1]
  if (current.distanceTo(last) > 0.0001) {
    curve.add(new THREE.LineCurve3(current, last))
  }
  return curve
}

function isFlowPipeSegmented(pipe) {
  return !!(pipe && pipe.userData && pipe.userData.isFlowPipe && pipe.userData.pipePoints && pipe.userData.flowMaterial && pipe.userData.roundedTube)
}

function createFlowPipe(points, color, radius, flowSpeed, highlightColor, dashLength) {
  if (!points || points.length < 2) return null

  var pipePoints = points.map(function(p) { return new THREE.Vector3(p.x, p.y, p.z) })
  var tubeRadius = radius || 0.1
  var baseColor = color || '#f4ead8'
  var pipeHighlightColor = highlightColor || '#ff6a00'
  var flowDashLength = getFlowDashLength(dashLength)
  var material = createFlowPipeMaterial(baseColor, pipeHighlightColor, flowDashLength)
  applyFlowPipeTextureRepeat(material, getFlowPipeLength(pipePoints), flowDashLength)
  var curve = createRoundedPipeCurve(pipePoints, tubeRadius)
  var geometry = new THREE.TubeGeometry(curve, Math.max(32, (points.length - 1) * 32), tubeRadius, 16, false)
  var tube = new THREE.Mesh(geometry, material)
  tube.userData.isFlowPipe = true
  tube.userData.flowSpeed = flowSpeed || 1.0
  tube.userData.flowDirection = 'forward'
  tube.userData.pipePoints = pipePoints
  tube.userData.flowMaterial = material
  tube.userData.flowOffset = 0
  tube.userData.flowEnabled = true
  tube.userData.flowDashLength = flowDashLength
  tube.userData.roundedTube = true
  return tube
}

function updateFlowPipeAnimation(pipe, deltaTime, flowCondition) {
  if (!pipe || !pipe.userData.isFlowPipe) return

  // 检查流动条件
  var shouldFlow = true
  if (flowCondition !== undefined) {
    shouldFlow = flowCondition
  } else if (pipe.userData.flowEnabled !== undefined) {
    shouldFlow = pipe.userData.flowEnabled
  }

  if (!shouldFlow) return

  // 根据流动方向更新偏移
  var direction = pipe.userData.flowDirection || 'forward'
  var speed = pipe.userData.flowSpeed || 1.0
  
  if (direction === 'forward') {
    pipe.userData.flowOffset -= speed * deltaTime * 0.5
    if (pipe.userData.flowOffset < 0) {
      pipe.userData.flowOffset += 1
    }
  } else {
    pipe.userData.flowOffset += speed * deltaTime * 0.5
    if (pipe.userData.flowOffset > 1) {
      pipe.userData.flowOffset -= 1
    }
  }

  if (pipe.userData.flowMaterial && pipe.userData.flowMaterial.map) {
    pipe.userData.flowMaterial.map.offset.x = pipe.userData.flowOffset
  }
}

function updateFlowPipe(pipe, obj) {
  if (!pipe || !pipe.userData.isFlowPipe || !obj) return

  if (obj.flowSpeed !== undefined) {
    pipe.userData.flowSpeed = obj.flowSpeed
  }

  if (obj.flowDirection !== undefined) {
    pipe.userData.flowDirection = obj.flowDirection
  }

  var nextColor = obj.color || '#f4ead8'
  var nextHighlightColor = obj.highlightColor || '#ff6a00'
  var nextRadius = obj.radius !== undefined ? obj.radius : 0.1
  var nextDashLength = getFlowDashLength(obj.flowDashLength)
  var textureKey = nextColor + '|' + nextHighlightColor + '|' + nextDashLength

  if (!pipe.userData.flowMaterial) {
    pipe.userData.flowMaterial = createFlowPipeMaterial(nextColor, nextHighlightColor, nextDashLength)
    pipe.material = pipe.userData.flowMaterial
  }

  if (pipe.userData.lastTextureKey !== textureKey) {
    pipe.userData.lastTextureKey = textureKey
    pipe.userData.flowDashLength = nextDashLength
    var nextTexture = createFlowPipeTexture(nextColor, nextHighlightColor, nextDashLength)
    if (pipe.userData.flowMaterial.map) {
      pipe.userData.flowMaterial.map.dispose()
    }
    pipe.userData.flowMaterial.map = nextTexture
    applyFlowPipeTextureRepeat(pipe.userData.flowMaterial, getFlowPipeLength(pipe.userData.pipePoints), nextDashLength)
    pipe.userData.flowMaterial.needsUpdate = true
    pipe.material = pipe.userData.flowMaterial
  }

  if (pipe.userData.lastRadius !== nextRadius && pipe.userData.pipePoints) {
    pipe.userData.lastRadius = nextRadius
    if (pipe.geometry) pipe.geometry.dispose()
    var nextCurve = createRoundedPipeCurve(pipe.userData.pipePoints, nextRadius)
    pipe.geometry = new THREE.TubeGeometry(nextCurve, Math.max(32, (pipe.userData.pipePoints.length - 1) * 32), nextRadius, 16, false)
    applyFlowPipeTextureRepeat(pipe.userData.flowMaterial, getFlowPipeLength(pipe.userData.pipePoints), nextDashLength)
  }
}

/*

  // 更新流动速度
  if (obj.flowSpeed !== undefined) {
    pipe.userData.flowSpeed = obj.flowSpeed
  }

  // 更新流动方向
  if (obj.flowDirection !== undefined) {
    pipe.userData.flowDirection = obj.flowDirection
  }

  // 检查属性是否真正变化
  var needsTextureUpdate = false
  var needsGeometryUpdate = false

  if (pipe.userData.lastColor !== obj.color) {
    pipe.userData.lastColor = obj.color
    needsTextureUpdate = true
  }

  if (pipe.userData.lastHighlightColor !== obj.highlightColor) {
    pipe.userData.lastHighlightColor = obj.highlightColor
    needsTextureUpdate = true
  }

  if (pipe.userData.lastRadius !== obj.radius) {
    pipe.userData.lastRadius = obj.radius
    needsGeometryUpdate = true
  }

  // 只有属性真正变化时才重新创建纹理和几何体
  if (needsTextureUpdate || needsGeometryUpdate) {
    var color = obj.color || '#f4ead8'
    var highlightColor = obj.highlightColor || '#ff6a00'
    var radius = obj.radius !== undefined ? obj.radius : 0.1

    if (needsTextureUpdate) {
      var flowTexture = createFlowPipeTexture(color, highlightColor)
      if (pipe.userData.flowMaterial && pipe.userData.flowMaterial.map) {
        pipe.userData.flowMaterial.map.dispose()
      }
      pipe.userData.flowMaterial.map = flowTexture
      pipe.userData.flowMaterial.needsUpdate = true
      needsTextureUpdate = false
    }

    // 重新创建纹理
    if (needsTextureUpdate) {
      var canvas = document.createElement('canvas')
      canvas.width = 512
      canvas.height = 64
      var ctx = canvas.getContext('2d')

      var gradient = ctx.createLinearGradient(0, 0, 512, 0)
      var baseColor = color
      
      for (var i = 0; i <= 10; i++) {
        var position = i / 10
        if (position < 0.2 || position > 0.8) {
          gradient.addColorStop(position, baseColor)
        } else {
          gradient.addColorStop(position, highlightColor)
        }
      }

      ctx.fillStyle = gradient
      ctx.fillRect(0, 0, 512, 64)

      var texture = new THREE.CanvasTexture(canvas)
      texture.wrapS = THREE.RepeatWrapping
      texture.wrapT = THREE.RepeatWrapping
      texture.repeat.set(4, 1)

      // 更新材质
      if (pipe.material && pipe.material.map) {
        pipe.material.map.dispose()
      }
      pipe.material.map = texture
      pipe.material.needsUpdate = true
    }

    // 如果半径变化，需要重新创建几何体
    if (needsGeometryUpdate && pipe.userData.pipePoints) {
      if (pipe.geometry) pipe.geometry.dispose()
      var roundedCurve = createRoundedPipeCurve(pipe.userData.pipePoints, radius)
      pipe.geometry = new THREE.TubeGeometry(roundedCurve, Math.max(32, (pipe.userData.pipePoints.length - 1) * 32), radius, 16, false)
    }
  }
}

// ===== Three.js 对象工厂 =====
*/

function normColorHex(value) {
  if (!value || typeof value !== 'string') return '#13c2c2'
  var s = value.trim()
  if (s === 'transparent') return '#13c2c2'
  if (/^#[0-9a-fA-F]{3}$/.test(s)) {
    return '#' + s[1] + s[1] + s[2] + s[2] + s[3] + s[3]
  }
  if (/^#[0-9a-fA-F]{6}$/.test(s)) return s
  if (/^#[0-9a-fA-F]{8}$/.test(s)) return s.slice(0, 7)
  return '#13c2c2'
}

function create2DComponentPlane(obj) {
  var sourceStyle = obj && obj.source2D && obj.source2D.style ? obj.source2D.style : null
  var width = 2
  var height = 1.2
  if (sourceStyle && sourceStyle.position) {
    var w = parseFloat(sourceStyle.position.w)
    var h = parseFloat(sourceStyle.position.h)
    if (!isNaN(w) && !isNaN(h) && w > 0 && h > 0) {
      width = Math.max(1, Math.min(4, w / 80))
      height = Math.max(0.6, Math.min(3, h / 80))
    }
  }

  var bgColor = sourceStyle && sourceStyle.backColor && sourceStyle.backColor !== 'transparent'
    ? sourceStyle.backColor
    : '#f7fbff'
  var borderColor = sourceStyle && sourceStyle.borderColor
    ? sourceStyle.borderColor
    : '#13c2c2'
  var textColor = sourceStyle && sourceStyle.foreColor
    ? sourceStyle.foreColor
    : '#13c2c2'

  var group = new THREE.Group()
  var panel = new THREE.Mesh(
    new THREE.PlaneGeometry(width, height),
    new THREE.MeshBasicMaterial({
      color: new THREE.Color(normColorHex(bgColor)),
      transparent: true,
      opacity: obj && obj.opacity !== undefined ? obj.opacity : 0.92,
      side: THREE.DoubleSide
    })
  )
  panel.rotation.x = -Math.PI / 2

  var points = [
    new THREE.Vector3(-width / 2, 0.01, -height / 2),
    new THREE.Vector3(width / 2, 0.01, -height / 2),
    new THREE.Vector3(width / 2, 0.01, height / 2),
    new THREE.Vector3(-width / 2, 0.01, height / 2),
    new THREE.Vector3(-width / 2, 0.01, -height / 2)
  ]
  var border = new THREE.Line(
    new THREE.BufferGeometry().setFromPoints(points),
    new THREE.LineBasicMaterial({ color: new THREE.Color(normColorHex(borderColor)) })
  )

  var label = obj && (obj.label2D || obj.textContent || obj.name) ? (obj.label2D || obj.textContent || obj.name) : '2D Component'
  var labelColor = parseInt(normColorHex(textColor).replace('#', ''), 16)
  var textSprite = createTextSprite(label, labelColor, 32, '#ffffffcc')
  textSprite.position.set(0, 0.08, 0)
  textSprite.scale.set(Math.min(width, textSprite.scale.x), 0.35, 1)

  group.add(panel, border, textSprite)
  group.userData.isMeshGroup = true
  group.userData.is2DComponent = true
  return group
}

function createMediaPlaceholderTexture(label, options) {
  options = options || {}
  var canvas = document.createElement('canvas')
  canvas.width = 512
  canvas.height = 320
  var ctx = canvas.getContext('2d')
  var bg = options.background || '#111827'
  var fg = options.foreground || '#ffffff'
  ctx.fillStyle = bg
  ctx.fillRect(0, 0, canvas.width, canvas.height)
  ctx.strokeStyle = options.border || 'rgba(255,255,255,0.6)'
  ctx.lineWidth = 8
  ctx.strokeRect(10, 10, canvas.width - 20, canvas.height - 20)
  ctx.fillStyle = fg
  ctx.font = 'bold 48px Arial, sans-serif'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText(label || 'Media', canvas.width / 2, canvas.height / 2)
  var texture = new THREE.CanvasTexture(canvas)
  texture.colorSpace = THREE.SRGBColorSpace
  texture.needsUpdate = true
  return texture
}

function disposeMediaPlane(mesh) {
  if (!mesh || !mesh.userData) return
  if (mesh.userData.mediaVideo) {
    try {
      mesh.userData.mediaVideo.pause()
      mesh.userData.mediaVideo.removeAttribute('src')
      mesh.userData.mediaVideo.load()
    } catch (e) {}
    mesh.userData.mediaVideo = null
  }
  if (mesh.userData.mediaTexture) {
    mesh.userData.mediaTexture.dispose()
    mesh.userData.mediaTexture = null
  }
}

function getMediaSource(obj) {
  if (!obj) return ''
  return obj.mediaUrl || obj.videoUrl || obj.imageUrl || obj.textureData || ''
}

function createMediaTexture(obj, type) {
  var source = getMediaSource(obj)
  if (type === 'video3d' && source) {
    var video = document.createElement('video')
    video.src = source
    video.crossOrigin = 'anonymous'
    video.loop = obj.mediaLoop !== false
    video.muted = obj.mediaMuted !== false
    video.playsInline = true
    video.autoplay = obj.mediaAutoplay !== false
    video.preload = 'auto'
    if (video.autoplay) {
      var playPromise = video.play()
      if (playPromise && playPromise.catch) playPromise.catch(function() {})
    }
    var videoTexture = new THREE.VideoTexture(video)
    videoTexture.colorSpace = THREE.SRGBColorSpace
    return { texture: videoTexture, video: video, key: source }
  }
  if (source) {
    var texture = new THREE.TextureLoader().load(source)
    texture.colorSpace = THREE.SRGBColorSpace
    return { texture: texture, video: null, key: source }
  }
  return {
    texture: createMediaPlaceholderTexture(type === 'video3d' ? '视频' : '图片', {
      background: type === 'video3d' ? '#111827' : '#164e63',
      foreground: '#ffffff'
    }),
    video: null,
    key: ''
  }
}

function createMediaPlane(obj, type) {
  obj = obj || {}
  var aspect = Number(obj.mediaAspect)
  if (!isFinite(aspect) || aspect <= 0) aspect = type === 'video3d' ? 16 / 9 : 4 / 3
  var width = Number(obj.mediaWidth)
  if (!isFinite(width) || width <= 0) width = type === 'video3d' ? 1.6 : 1.4
  var height = width / aspect
  var media = createMediaTexture(obj, type)
  var material = new THREE.MeshBasicMaterial({
    map: media.texture,
    transparent: obj.opacity !== undefined ? obj.opacity < 1 : false,
    opacity: obj.opacity !== undefined ? obj.opacity : 1,
    side: THREE.DoubleSide
  })
  var mesh = new THREE.Mesh(new THREE.PlaneGeometry(width, height), material)
  mesh.userData.isMediaPlane = true
  mesh.userData.mediaType = type
  mesh.userData.mediaKey = media.key
  mesh.userData.mediaTexture = media.texture
  mesh.userData.mediaVideo = media.video
  return mesh
}

function updateMediaPlane(mesh, obj, type) {
  if (!mesh || !mesh.userData || !mesh.userData.isMediaPlane || !mesh.material) return
  var source = getMediaSource(obj)
  if (mesh.userData.mediaKey !== source || mesh.userData.mediaType !== type) {
    disposeMediaPlane(mesh)
    var media = createMediaTexture(obj, type)
    mesh.material.map = media.texture
    mesh.userData.mediaType = type
    mesh.userData.mediaKey = media.key
    mesh.userData.mediaTexture = media.texture
    mesh.userData.mediaVideo = media.video
  }
  var aspect = Number(obj.mediaAspect)
  if (!isFinite(aspect) || aspect <= 0) aspect = type === 'video3d' ? 16 / 9 : 4 / 3
  var width = Number(obj.mediaWidth)
  if (!isFinite(width) || width <= 0) width = type === 'video3d' ? 1.6 : 1.4
  var height = width / aspect
  if (mesh.geometry) mesh.geometry.dispose()
  mesh.geometry = new THREE.PlaneGeometry(width, height)
  mesh.material.opacity = obj.opacity !== undefined ? obj.opacity : 1
  mesh.material.transparent = mesh.material.opacity < 1
  mesh.material.needsUpdate = true
}

function createIndustrialExtensionMesh(type, mat) {
  const steel = new THREE.MeshStandardMaterial({ color: 0x777777, metalness: 0.65, roughness: 0.35 })
  const dark = new THREE.MeshStandardMaterial({ color: 0x333333, metalness: 0.35, roughness: 0.6 })
  const glowGreen = new THREE.MeshStandardMaterial({ color: 0x4dffa6, emissive: 0x1c8a55, emissiveIntensity: 0.35 })
  const glowBlue = new THREE.MeshStandardMaterial({ color: 0x66ccff, emissive: 0x0b5a77, emissiveIntensity: 0.35 })

  switch (type) {
    case 'stationyard': {
      const group = new THREE.Group()
      const baseMat = new THREE.MeshStandardMaterial({ color: 0x26323b, emissive: 0x080f16, emissiveIntensity: 0.18, metalness: 0.22, roughness: 0.58 })
      const glassMat = new THREE.MeshStandardMaterial({ color: 0x557d86, emissive: 0x102c35, emissiveIntensity: 0.22, transparent: true, opacity: 0.24, metalness: 0.12, roughness: 0.38 })
      const roadMat = new THREE.MeshStandardMaterial({ color: 0x30383d, emissive: 0x070b0f, emissiveIntensity: 0.08, roughness: 0.68 })
      const lineMat = new THREE.MeshStandardMaterial({ color: 0x86d8e6, emissive: 0x1f7f95, emissiveIntensity: 0.32, transparent: true, opacity: 0.58 })
      const base = new THREE.Mesh(new THREE.BoxGeometry(5.6, 0.08, 3.9), baseMat)
      base.position.y = 0.04
      const glass = new THREE.Mesh(new THREE.BoxGeometry(4.65, 0.035, 2.75), glassMat)
      glass.position.set(0, 0.105, -0.08)
      const road = new THREE.Mesh(new THREE.BoxGeometry(5.8, 0.045, 0.46), roadMat)
      road.position.set(0, 0.13, 1.58)
      const centerLine = new THREE.Mesh(new THREE.BoxGeometry(5.2, 0.018, 0.025), lineMat)
      centerLine.position.set(0, 0.158, 1.58)
      const ring1 = new THREE.Mesh(new THREE.TorusGeometry(2.18, 0.012, 8, 96), lineMat)
      ring1.rotation.x = Math.PI / 2
      ring1.position.y = 0.14
      ring1.scale.z = 0.58
      const ring2 = ring1.clone()
      ring2.scale.set(0.72, 0.72, 0.42)
      group.add(base, glass, road, centerLine, ring1, ring2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'stationhall': {
      const group = new THREE.Group()
      const wallMat = new THREE.MeshStandardMaterial({ color: 0xdfe8e8, metalness: 0.22, roughness: 0.48 })
      const roofMat = new THREE.MeshStandardMaterial({ color: 0x314951, emissive: 0x061820, emissiveIntensity: 0.22, metalness: 0.45, roughness: 0.3 })
      const winMat = new THREE.MeshStandardMaterial({ color: 0xa8d4df, emissive: 0x1b5465, emissiveIntensity: 0.28, metalness: 0.12, roughness: 0.32 })
      const trimMat = new THREE.MeshStandardMaterial({ color: 0x92d7e5, emissive: 0x26798a, emissiveIntensity: 0.36, transparent: true, opacity: 0.62 })
      const body = new THREE.Mesh(new THREE.BoxGeometry(2.6, 0.68, 0.9), wallMat)
      body.position.y = 0.38
      const roof = new THREE.Mesh(new THREE.BoxGeometry(2.78, 0.12, 1.02), roofMat)
      roof.position.y = 0.78
      const annex = new THREE.Mesh(new THREE.BoxGeometry(0.82, 0.55, 0.78), wallMat)
      annex.position.set(1.78, 0.32, 0.02)
      const annexRoof = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.1, 0.88), roofMat)
      annexRoof.position.set(1.78, 0.63, 0.02)
      for (let i = 0; i < 7; i++) {
        const win = new THREE.Mesh(new THREE.BoxGeometry(0.18, 0.22, 0.025), winMat)
        win.position.set(-1.05 + i * 0.34, 0.44, -0.466)
        group.add(win)
      }
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.24, 0.36, 0.03), new THREE.MeshStandardMaterial({ color: 0x46545a, metalness: 0.3 }))
      door.position.set(0.95, 0.24, -0.47)
      const sign = new THREE.Mesh(new THREE.BoxGeometry(0.72, 0.09, 0.025), glowBlue)
      sign.position.set(0, 0.74, -0.48)
      const trimFront = new THREE.Mesh(new THREE.BoxGeometry(2.82, 0.025, 0.035), trimMat)
      trimFront.position.set(0, 0.87, -0.53)
      const trimBack = trimFront.clone()
      trimBack.position.z = 0.53
      group.add(body, roof, annex, annexRoof, door, sign, trimFront, trimBack)
      group.userData.isMeshGroup = true
      return group
    }
    case 'sitefence': {
      const group = new THREE.Group()
      const postMat = new THREE.MeshStandardMaterial({ color: 0x6b7b83, metalness: 0.45, roughness: 0.45 })
      const railMat = new THREE.MeshStandardMaterial({ color: 0x9aa7ad, metalness: 0.5, roughness: 0.35 })
      for (let i = 0; i < 7; i++) {
        const post = new THREE.Mesh(new THREE.CylinderGeometry(0.025, 0.025, 0.55, 8), postMat)
        post.position.set(-1.5 + i * 0.5, 0.28, 0)
        group.add(post)
      }
      for (let y = 0; y < 3; y++) {
        const rail = new THREE.Mesh(new THREE.BoxGeometry(3.1, 0.025, 0.025), railMat)
        rail.position.set(0, 0.16 + y * 0.14, 0)
        group.add(rail)
      }
      group.userData.isMeshGroup = true
      return group
    }
    case 'coolingrack4': {
      const group = new THREE.Group()
      const frameMat = new THREE.MeshStandardMaterial({ color: 0xcfd8d6, metalness: 0.35, roughness: 0.36 })
      const glassMat = new THREE.MeshStandardMaterial({ color: 0x2f7d5c, emissive: 0x0b2f22, emissiveIntensity: 0.16, transparent: true, opacity: 0.66, roughness: 0.36 })
      const darkMat = new THREE.MeshStandardMaterial({ color: 0x27333a, metalness: 0.42, roughness: 0.35 })
      const base = new THREE.Mesh(new THREE.BoxGeometry(3.0, 0.12, 1.08), darkMat)
      base.position.y = 0.06
      const body = new THREE.Mesh(new THREE.BoxGeometry(2.9, 0.82, 0.86), glassMat)
      body.position.y = 0.52
      const top = new THREE.Mesh(new THREE.BoxGeometry(3.05, 0.08, 1.02), frameMat)
      top.position.y = 0.96
      group.add(base, body, top)
      for (let i = 0; i < 5; i++) {
        const post = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.95, 0.04), frameMat)
        post.position.set(-1.45 + i * 0.72, 0.52, -0.48)
        const postRear = post.clone()
        postRear.position.z = 0.48
        group.add(post, postRear)
      }
      for (let i = 0; i < 4; i++) {
        const fan = new THREE.Mesh(new THREE.TorusGeometry(0.18, 0.018, 8, 30), darkMat)
        fan.rotation.x = Math.PI / 2
        fan.position.set(-1.08 + i * 0.72, 1.04, 0)
        const hub = new THREE.Mesh(new THREE.CylinderGeometry(0.035, 0.035, 0.025, 16), darkMat)
        hub.position.copy(fan.position)
        hub.rotation.x = Math.PI / 2
        group.add(fan, hub)
      }
      for (let i = 0; i < 4; i++) {
        const pump = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.34, 18), steel)
        pump.rotation.z = Math.PI / 2
        pump.position.set(-1.08 + i * 0.72, 0.24, -0.52)
        group.add(pump)
      }
      const pipeA = new THREE.Mesh(new THREE.CylinderGeometry(0.035, 0.035, 3.25, 16), steel)
      pipeA.rotation.z = Math.PI / 2
      pipeA.position.set(0, 0.2, -0.66)
      const pipeB = pipeA.clone()
      pipeB.position.z = 0.66
      group.add(pipeA, pipeB)
      group.userData.isMeshGroup = true
      group.userData.hasRotorAnimation = true
      return group
    }
    case 'datapanel': {
      const group = new THREE.Group()
      const panelMat = new THREE.MeshStandardMaterial({ color: 0x101f2a, emissive: 0x071521, emissiveIntensity: 0.32, transparent: true, opacity: 0.72, metalness: 0.2, roughness: 0.34 })
      const screenMat = new THREE.MeshStandardMaterial({ color: 0x1b3a46, emissive: 0x124f61, emissiveIntensity: 0.26, transparent: true, opacity: 0.56 })
      const borderMat = new THREE.MeshStandardMaterial({ color: 0x8bd7e8, emissive: 0x267f9a, emissiveIntensity: 0.42, transparent: true, opacity: 0.68 })
      const panel = new THREE.Mesh(new THREE.BoxGeometry(1.05, 0.62, 0.04), panelMat)
      panel.position.y = 0.42
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.92, 0.46, 0.02), screenMat)
      screen.position.set(0, 0.43, -0.025)
      const top = new THREE.Mesh(new THREE.BoxGeometry(1.12, 0.018, 0.025), borderMat)
      top.position.set(0, 0.745, -0.04)
      const bottom = top.clone()
      bottom.position.y = 0.095
      const left = new THREE.Mesh(new THREE.BoxGeometry(0.018, 0.62, 0.025), borderMat)
      left.position.set(-0.56, 0.42, -0.04)
      const right = left.clone()
      right.position.x = 0.56
      for (let i = 0; i < 4; i++) {
        const bar = new THREE.Mesh(new THREE.BoxGeometry(0.12 + i * 0.05, 0.025, 0.012), glowGreen)
        bar.position.set(-0.28, 0.3 + i * 0.08, -0.04)
        group.add(bar)
      }
      const stand = new THREE.Mesh(new THREE.CylinderGeometry(0.035, 0.04, 0.42, 10), steel)
      stand.position.y = 0.18
      group.add(panel, screen, top, bottom, left, right, stand)
      group.userData.isMeshGroup = true
      return group
    }
    case 'neonflowline': {
      const group = new THREE.Group()
      const lineMat = new THREE.MeshStandardMaterial({ color: 0x9adceb, emissive: 0x2b8297, emissiveIntensity: 0.42, transparent: true, opacity: 0.62 })
      const nodeMat = new THREE.MeshStandardMaterial({ color: 0xdff8ff, emissive: 0x67b7c8, emissiveIntensity: 0.36, transparent: true, opacity: 0.68 })
      const line = new THREE.Mesh(new THREE.BoxGeometry(1.4, 0.018, 0.018), lineMat)
      line.position.y = 0.025
      const nodeA = new THREE.Mesh(new THREE.SphereGeometry(0.045, 16, 10), nodeMat)
      nodeA.position.set(-0.7, 0.03, 0)
      const nodeB = nodeA.clone()
      nodeB.position.x = 0.7
      group.add(line, nodeA, nodeB)
      group.userData.isMeshGroup = true
      return group
    }
    case 'titlebillboard': {
      const group = new THREE.Group()
      const bgMat = new THREE.MeshStandardMaterial({ color: 0x0b2230, emissive: 0x063041, emissiveIntensity: 0.62, transparent: true, opacity: 0.8, metalness: 0.18, roughness: 0.22 })
      const edgeMat = new THREE.MeshStandardMaterial({ color: 0x66ccff, emissive: 0x13c2c2, emissiveIntensity: 0.9, transparent: true, opacity: 0.85 })
      const panel = new THREE.Mesh(new THREE.BoxGeometry(1.7, 0.34, 0.035), bgMat)
      panel.position.y = 0.9
      const top = new THREE.Mesh(new THREE.BoxGeometry(1.78, 0.025, 0.025), edgeMat)
      top.position.set(0, 1.085, -0.03)
      const bottom = top.clone()
      bottom.position.y = 0.715
      const wingL = new THREE.Mesh(new THREE.BoxGeometry(0.36, 0.018, 0.02), edgeMat)
      wingL.position.set(-1.08, 0.9, -0.03)
      wingL.rotation.z = -0.28
      const wingR = wingL.clone()
      wingR.position.x = 1.08
      wingR.rotation.z = 0.28
      group.add(panel, top, bottom, wingL, wingR)
      group.userData.isMeshGroup = true
      return group
    }
    case 'screenframe': {
      const group = new THREE.Group()
      const edgeMat = new THREE.MeshStandardMaterial({ color: 0x8ab8d5, emissive: 0x2d6183, emissiveIntensity: 0.48, transparent: true, opacity: 0.56 })
      const glassMat = new THREE.MeshStandardMaterial({ color: 0x182f43, emissive: 0x071421, emissiveIntensity: 0.24, transparent: true, opacity: 0.16, roughness: 0.42 })
      const panel = new THREE.Mesh(new THREE.BoxGeometry(7.6, 0.08, 4.6), glassMat)
      panel.position.y = 0.04
      const top = new THREE.Mesh(new THREE.BoxGeometry(7.8, 0.035, 0.035), edgeMat)
      top.position.set(0, 0.12, -2.35)
      const bottom = top.clone()
      bottom.position.z = 2.35
      const left = new THREE.Mesh(new THREE.BoxGeometry(0.035, 0.035, 4.7), edgeMat)
      left.position.set(-3.9, 0.12, 0)
      const right = left.clone()
      right.position.x = 3.9
      group.add(panel, top, bottom, left, right)
      group.userData.isMeshGroup = true
      return group
    }
    case 'tophudbar': {
      const group = new THREE.Group()
      const bgMat = new THREE.MeshStandardMaterial({ color: 0x1b3348, emissive: 0x0b1d2c, emissiveIntensity: 0.38, transparent: true, opacity: 0.74 })
      const edgeMat = new THREE.MeshStandardMaterial({ color: 0x9fc7dd, emissive: 0x316d8a, emissiveIntensity: 0.45, transparent: true, opacity: 0.72 })
      const center = new THREE.Mesh(new THREE.BoxGeometry(2.6, 0.36, 0.035), bgMat)
      center.position.y = 1.05
      const wingL = new THREE.Mesh(new THREE.BoxGeometry(1.55, 0.22, 0.03), bgMat)
      wingL.position.set(-2.1, 1.0, 0)
      wingL.rotation.z = -0.08
      const wingR = wingL.clone()
      wingR.position.x = 2.1
      wingR.rotation.z = 0.08
      const underline = new THREE.Mesh(new THREE.BoxGeometry(5.4, 0.025, 0.025), edgeMat)
      underline.position.set(0, 0.82, -0.02)
      group.add(center, wingL, wingR, underline)
      group.userData.isMeshGroup = true
      return group
    }
    case 'sidehudpanel': {
      const group = new THREE.Group()
      const bgMat = new THREE.MeshStandardMaterial({ color: 0x172735, emissive: 0x091723, emissiveIntensity: 0.32, transparent: true, opacity: 0.7 })
      const edgeMat = new THREE.MeshStandardMaterial({ color: 0x90bfd6, emissive: 0x315f78, emissiveIntensity: 0.42, transparent: true, opacity: 0.66 })
      const lineMat = new THREE.MeshStandardMaterial({ color: 0xb8e4ef, emissive: 0x3a7d8f, emissiveIntensity: 0.32, transparent: true, opacity: 0.6 })
      const panel = new THREE.Mesh(new THREE.BoxGeometry(1.05, 1.72, 0.035), bgMat)
      panel.position.y = 0.86
      const top = new THREE.Mesh(new THREE.BoxGeometry(1.12, 0.025, 0.025), edgeMat)
      top.position.set(0, 1.74, -0.03)
      const bottom = top.clone()
      bottom.position.y = 0.02
      const left = new THREE.Mesh(new THREE.BoxGeometry(0.025, 1.75, 0.025), edgeMat)
      left.position.set(-0.56, 0.86, -0.03)
      const right = left.clone()
      right.position.x = 0.56
      group.add(panel, top, bottom, left, right)
      for (let i = 0; i < 6; i++) {
        const line = new THREE.Mesh(new THREE.BoxGeometry(0.58 - (i % 3) * 0.08, 0.022, 0.012), lineMat)
        line.position.set(-0.06, 1.42 - i * 0.22, -0.045)
        group.add(line)
        const dot = new THREE.Mesh(new THREE.SphereGeometry(0.025, 12, 8), lineMat)
        dot.position.set(-0.38, line.position.y, -0.045)
        group.add(dot)
      }
      group.userData.isMeshGroup = true
      return group
    }
    case 'bottomnavhud': {
      const group = new THREE.Group()
      const bgMat = new THREE.MeshStandardMaterial({ color: 0x1b2c3a, emissive: 0x0a1723, emissiveIntensity: 0.34, transparent: true, opacity: 0.78 })
      const itemMat = new THREE.MeshStandardMaterial({ color: 0x3e7892, emissive: 0x1b5269, emissiveIntensity: 0.36, transparent: true, opacity: 0.58 })
      const edgeMat = new THREE.MeshStandardMaterial({ color: 0xa7d7e4, emissive: 0x3b7f91, emissiveIntensity: 0.36, transparent: true, opacity: 0.66 })
      const bar = new THREE.Mesh(new THREE.BoxGeometry(4.9, 0.34, 0.04), bgMat)
      bar.position.y = 0.18
      group.add(bar)
      for (let i = 0; i < 5; i++) {
        const item = new THREE.Mesh(new THREE.BoxGeometry(0.62, 0.18, 0.035), itemMat)
        item.position.set(-1.6 + i * 0.8, 0.2, -0.035)
        group.add(item)
      }
      const top = new THREE.Mesh(new THREE.BoxGeometry(5.2, 0.022, 0.02), edgeMat)
      top.position.set(0, 0.39, -0.04)
      group.add(top)
      group.userData.isMeshGroup = true
      return group
    }
    case 'distillationtower':
    case 'absorbercolumn': {
      const group = new THREE.Group()
      const shell = new THREE.Mesh(new THREE.CylinderGeometry(0.22, 0.26, 2.4, 32), mat.clone())
      shell.position.y = 1.2
      const skirt = new THREE.Mesh(new THREE.CylinderGeometry(0.32, 0.38, 0.35, 24), mat.clone())
      skirt.position.y = 0.18
      const top = new THREE.Mesh(new THREE.SphereGeometry(0.22, 24, 12, 0, Math.PI * 2, 0, Math.PI / 2), mat.clone())
      top.position.y = 2.4
      for (let i = 0; i < 6; i++) {
        const tray = new THREE.Mesh(new THREE.CylinderGeometry(0.235, 0.235, 0.018, 24), steel)
        tray.position.y = 0.55 + i * 0.28
        group.add(tray)
      }
      const nozzleA = new THREE.Mesh(new THREE.CylinderGeometry(0.045, 0.045, 0.42, 12), mat.clone())
      nozzleA.rotation.z = Math.PI / 2
      nozzleA.position.set(0.36, 1.2, 0)
      const nozzleB = nozzleA.clone()
      nozzleB.position.set(-0.36, 1.7, 0)
      const ladder = new THREE.Mesh(new THREE.BoxGeometry(0.035, 1.8, 0.035), steel)
      ladder.position.set(0.29, 1.25, 0)
      group.add(shell, skirt, top, nozzleA, nozzleB, ladder)
      group.userData.isMeshGroup = true
      return group
    }
    case 'platefilter': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.35, 0.08, 0.42), steel)
      base.position.y = 0.06
      for (let i = 0; i < 9; i++) {
        const plate = new THREE.Mesh(new THREE.BoxGeometry(0.055, 0.58, 0.5), mat.clone())
        plate.position.set(-0.36 + i * 0.09, 0.42, 0)
        group.add(plate)
      }
      const rail1 = new THREE.Mesh(new THREE.BoxGeometry(1.25, 0.04, 0.04), steel)
      rail1.position.set(0, 0.74, 0.28)
      const rail2 = rail1.clone()
      rail2.position.z = -0.28
      group.add(base, rail1, rail2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'bagfilter': {
      const group = new THREE.Group()
      const shell = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.8, 0.55), mat.clone())
      shell.position.y = 0.75
      const hopper = new THREE.Mesh(new THREE.ConeGeometry(0.42, 0.55, 4), mat.clone())
      hopper.rotation.y = Math.PI / 4
      hopper.position.y = 0.22
      const duct = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.7, 12), mat.clone())
      duct.rotation.z = Math.PI / 2
      duct.position.set(0.55, 0.85, 0)
      group.add(shell, hopper, duct)
      group.userData.isMeshGroup = true
      return group
    }
    case 'dosingunit':
    case 'cipstation':
    case 'skidunit': {
      const group = new THREE.Group()
      const skid = new THREE.Mesh(new THREE.BoxGeometry(1.55, 0.08, 0.78), steel)
      skid.position.y = 0.04
      const tank1 = new THREE.Mesh(new THREE.CylinderGeometry(0.18, 0.18, 0.55, 24), mat.clone())
      tank1.position.set(-0.35, 0.38, -0.12)
      const tank2 = tank1.clone()
      tank2.position.x = 0.07
      const pump = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.14, 0.22, 16), glowGreen)
      pump.position.set(0.45, 0.22, 0.15)
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.28, 0.42, 0.12), dark)
      panel.position.set(0.58, 0.34, -0.22)
      group.add(skid, tank1, tank2, pump, panel)
      group.userData.isMeshGroup = true
      return group
    }
    case 'airblower':
    case 'vacuumunit': {
      const group = new THREE.Group()
      const lobe1 = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.5, 24), mat.clone())
      lobe1.rotation.z = Math.PI / 2
      lobe1.position.y = 0.3
      const lobe2 = lobe1.clone()
      lobe2.position.y = 0.05
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.85, 0.12, 0.6), mat.clone())
      base.position.y = -0.12
      const motor = new THREE.Mesh(new THREE.CylinderGeometry(0.16, 0.16, 0.42, 20), steel)
      motor.rotation.z = Math.PI / 2
      motor.position.set(0.55, 0.1, 0)
      group.add(lobe1, lobe2, base, motor)
      group.userData.isMeshGroup = true
      return group
    }
    case 'chillerunit': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.55, 0.7), mat.clone())
      body.position.y = 0.32
      for (let i = 0; i < 2; i++) {
        const fanRing = new THREE.Mesh(new THREE.TorusGeometry(0.17, 0.02, 8, 24), dark)
        fanRing.rotation.x = Math.PI / 2
        fanRing.position.set(-0.28 + i * 0.56, 0.62, 0)
        group.add(fanRing)
      }
      const coil = new THREE.Mesh(new THREE.BoxGeometry(1.1, 0.35, 0.03), glowBlue)
      coil.position.set(0, 0.32, -0.37)
      group.add(body, coil)
      group.userData.isMeshGroup = true
      return group
    }
    case 'coolingFanUnit': {
      const group = new THREE.Group()

      // 材质定义
      const bodyMat = new THREE.MeshStandardMaterial({ color: 0xc5cdd1, metalness: 0.28, roughness: 0.38 })
      const finMat = new THREE.MeshStandardMaterial({ color: 0xa8b0b4, metalness: 0.35, roughness: 0.32 })
      const darkMat = new THREE.MeshStandardMaterial({ color: 0x2a343a, metalness: 0.55, roughness: 0.28 })
      const frameMat = new THREE.MeshStandardMaterial({ color: 0x505a60, metalness: 0.48, roughness: 0.3 })
      const bladeMat = new THREE.MeshStandardMaterial({
        color: 0x4ec3e0,
        metalness: 0.3,
        roughness: 0.25,
        side: THREE.DoubleSide
      })
      const copperMat = new THREE.MeshStandardMaterial({ color: 0xd47028, metalness: 0.62, roughness: 0.18 })
      const hubMat = new THREE.MeshStandardMaterial({ color: 0x3a444c, metalness: 0.6, roughness: 0.22 })

      // 底座框架
      const baseFrame = new THREE.Mesh(new THREE.BoxGeometry(2.06, 0.05, 1.06), darkMat)
      baseFrame.position.y = 0.025
      group.add(baseFrame)

      // 主体箱体
      const body = new THREE.Mesh(new THREE.BoxGeometry(1.96, 0.58, 0.96), bodyMat)
      body.position.y = 0.32
      group.add(body)

      // 竖向散热片（前后两面）—— 图片中箱体表面最明显的特征
      for (let i = 0; i < 14; i++) {
        const x = -0.84 + i * 0.128
        const finFront = new THREE.Mesh(new THREE.BoxGeometry(0.008, 0.50, 0.015), finMat)
        finFront.position.set(x, 0.32, -0.49)
        const finRear = finFront.clone()
        finRear.position.z = 0.49
        group.add(finFront, finRear)
      }

      // 顶部面板
      const topPlate = new THREE.Mesh(new THREE.BoxGeometry(2.02, 0.045, 1.02), bodyMat)
      topPlate.position.y = 0.632
      group.add(topPlate)

      // 顶部边框装饰
      const topTrimFront = new THREE.Mesh(new THREE.BoxGeometry(2.02, 0.025, 0.02), frameMat)
      topTrimFront.position.set(0, 0.66, -0.51)
      const topTrimRear = topTrimFront.clone()
      topTrimRear.position.z = 0.51
      const topTrimLeft = new THREE.Mesh(new THREE.BoxGeometry(0.02, 0.025, 1.02), frameMat)
      topTrimLeft.position.set(-1.01, 0.66, 0)
      const topTrimRight = topTrimLeft.clone()
      topTrimRight.position.x = 1.01
      group.add(topTrimFront, topTrimRear, topTrimLeft, topTrimRight)

      // 两个大型顶部风扇
      const fanCenters = [-0.50, 0.50]
      fanCenters.forEach((x) => {
        // 风扇底座微凸台
        const fanBase = new THREE.Mesh(new THREE.CylinderGeometry(0.44, 0.44, 0.02, 48), frameMat)
        fanBase.position.set(x, 0.665, 0)
        group.add(fanBase)

        // 外圈金属环
        const outerRing = new THREE.Mesh(new THREE.TorusGeometry(0.42, 0.015, 10, 64), frameMat)
        outerRing.rotation.x = Math.PI / 2
        outerRing.position.set(x, 0.685, 0)

        // 内圈保护网
        const innerRing = new THREE.Mesh(new THREE.TorusGeometry(0.32, 0.005, 6, 48), darkMat)
        innerRing.rotation.x = Math.PI / 2
        innerRing.position.set(x, 0.688, 0)

        // 转子组（动画旋转部分）
        const rotor = new THREE.Group()
        rotor.position.set(x, 0.69, 0)

        // 6片蓝色扇叶
        for (let i = 0; i < 6; i++) {
          const bladeShape = new THREE.Shape()
          bladeShape.moveTo(0.04, -0.02)
          bladeShape.bezierCurveTo(0.09, -0.055, 0.20, -0.065, 0.28, -0.02)
          bladeShape.bezierCurveTo(0.23, 0.025, 0.12, 0.05, 0.04, 0.025)
          bladeShape.lineTo(0.04, -0.02)
          const blade = new THREE.Mesh(new THREE.ShapeGeometry(bladeShape), bladeMat)
          blade.rotation.x = -Math.PI / 2
          blade.rotation.y = i * Math.PI / 3
          blade.rotation.z = -0.20
          rotor.add(blade)
        }

        // 中心毂
        const hub = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.03, 24), hubMat)
        rotor.add(hub)

        // 毂上小螺丝装饰
        for (let i = 0; i < 3; i++) {
          const screw = new THREE.Mesh(new THREE.CylinderGeometry(0.008, 0.008, 0.015, 8), hubMat)
          const angle = i * Math.PI * 2 / 3
          screw.position.set(Math.cos(angle) * 0.025, 0.015, Math.sin(angle) * 0.025)
          rotor.add(screw)
        }

        group.add(outerRing, innerRing, rotor)
        group.userData.rotors = group.userData.rotors || []
        group.userData.rotors.push(rotor)
      })

      // 四个支撑脚
      const footGeo = new THREE.BoxGeometry(0.12, 0.06, 0.12)
      const footPositions = [[-0.78, 0.03, -0.38], [0.78, 0.03, -0.38], [-0.78, 0.03, 0.38], [0.78, 0.03, 0.38]]
      footPositions.forEach(([fx, fy, fz]) => {
        const foot = new THREE.Mesh(footGeo, darkMat)
        foot.position.set(fx, fy, fz)
        group.add(foot)
      })

      // 侧面铜管接口
      const pipeSide = new THREE.Mesh(new THREE.CylinderGeometry(0.022, 0.022, 0.35, 12), copperMat)
      pipeSide.rotation.z = Math.PI / 2
      pipeSide.position.set(1.05, 0.15, 0)
      const pipeSide2 = pipeSide.clone()
      pipeSide2.position.set(-1.05, 0.15, 0.2)

      // 正面阀门接头
      const valve = new THREE.Mesh(new THREE.CylinderGeometry(0.035, 0.035, 0.05, 12), copperMat)
      valve.rotation.x = Math.PI / 2
      valve.position.set(0, 0.15, -0.52)

      group.add(pipeSide, pipeSide2, valve)

      group.userData.isMeshGroup = true
      group.userData.hasRotorAnimation = true
      return group
    }
    case 'coolingcoil': {
      const group = new THREE.Group()
      for (let i = 0; i < 5; i++) {
        const tube = new THREE.Mesh(new THREE.TorusGeometry(0.22 + i * 0.035, 0.012, 8, 36), glowBlue)
        tube.rotation.x = Math.PI / 2
        tube.position.y = 0.08 + i * 0.06
        group.add(tube)
      }
      group.userData.isMeshGroup = true
      return group
    }
    case 'weighhopper': {
      const group = new THREE.Group()
      const hopper = new THREE.Mesh(new THREE.ConeGeometry(0.42, 0.65, 4), mat.clone())
      hopper.rotation.y = Math.PI / 4
      hopper.position.y = 0.48
      const neck = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.25, 12), mat.clone())
      neck.position.y = 0.05
      for (let sx = -1; sx <= 1; sx += 2) {
        for (let sz = -1; sz <= 1; sz += 2) {
          const leg = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.7, 0.04), steel)
          leg.position.set(sx * 0.34, 0.35, sz * 0.34)
          group.add(leg)
        }
      }
      group.add(hopper, neck)
      group.userData.isMeshGroup = true
      return group
    }
    case 'magneticseparator':
    case 'vibratingscreen':
    case 'rollermill': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.1, 0.08, 0.55), steel)
      base.position.y = 0.08
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.42, 0.5), mat.clone())
      body.position.y = 0.4
      const drumA = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.62, 24), steel)
      drumA.rotation.x = Math.PI / 2
      drumA.position.set(-0.2, 0.43, 0)
      const drumB = drumA.clone()
      drumB.position.x = 0.2
      const hopper = new THREE.Mesh(new THREE.ConeGeometry(0.24, 0.32, 4), mat.clone())
      hopper.rotation.y = Math.PI / 4
      hopper.position.y = 0.82
      group.add(base, body, drumA, drumB, hopper)
      group.userData.isMeshGroup = true
      return group
    }
    case 'controlcabinet':
    case 'mccpanel': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.48, 1.15, 0.32), mat.clone())
      body.position.y = 0.58
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.12, 0.015), glowBlue)
      screen.position.set(-0.1, 0.92, -0.18)
      const lamp1 = new THREE.Mesh(new THREE.SphereGeometry(0.025, 12, 8), new THREE.MeshStandardMaterial({ color: 0xff4d4f, emissive: 0xff2222, emissiveIntensity: 0.5 }))
      lamp1.position.set(0.12, 0.92, -0.18)
      const lamp2 = new THREE.Mesh(new THREE.SphereGeometry(0.025, 12, 8), glowGreen)
      lamp2.position.set(0.19, 0.92, -0.18)
      group.add(body, screen, lamp1, lamp2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'analyzerhouse': {
      const group = new THREE.Group()
      const room = new THREE.Mesh(new THREE.BoxGeometry(1.1, 0.75, 0.75), mat.clone())
      room.position.y = 0.42
      const roof = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.08, 0.85), steel)
      roof.position.y = 0.84
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.22, 0.45, 0.015), dark)
      door.position.set(-0.32, 0.3, -0.385)
      const analyzer = new THREE.Mesh(new THREE.BoxGeometry(0.22, 0.38, 0.18), new THREE.MeshStandardMaterial({ color: 0xeeeeee }))
      analyzer.position.set(0.28, 0.35, -0.18)
      group.add(room, roof, door, analyzer)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pipebridge':
    case 'maintenanceplatform': {
      const group = new THREE.Group()
      const deck = new THREE.Mesh(new THREE.BoxGeometry(1.8, 0.05, 0.72), steel)
      deck.position.y = 0.62
      for (let x = -0.8; x <= 0.8; x += 0.4) {
        const rail = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.36, 0.04), steel)
        rail.position.set(x, 0.82, -0.34)
        group.add(rail)
        const rail2 = rail.clone()
        rail2.position.z = 0.34
        group.add(rail2)
      }
      for (let x = -0.75; x <= 0.75; x += 1.5) {
        for (let z = -0.3; z <= 0.3; z += 0.6) {
          const leg = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.62, 0.05), steel)
          leg.position.set(x, 0.31, z)
          group.add(leg)
        }
      }
      group.add(deck)
      group.userData.isMeshGroup = true
      return group
    }
    default:
      return null
  }
}

function createThreeMesh(type, color) {
  const mat = new THREE.MeshStandardMaterial({
    color: color || 0x1890ff,
    metalness: 0.3,
    roughness: 0.5,
    transparent: true,
    opacity: 1.0,
  })

  let geo
  switch (type) {
    case 'box':       geo = new THREE.BoxGeometry(1,1,1); break
    case 'sphere':    geo = new THREE.SphereGeometry(0.6,32,32); break
    case 'cylinder':  geo = new THREE.CylinderGeometry(0.5,0.5,1.2,32); break
    case 'cone':      geo = new THREE.ConeGeometry(0.5,1.2,32); break
    case 'torus':     geo = new THREE.TorusGeometry(0.6,0.2,16,64); break
    case 'plane':     geo = new THREE.PlaneGeometry(2,2); break

    // ===== 工业 — 用组合模拟 =====
    case 'tank': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.5,0.5,1.5,32), mat.clone())
      const top  = new THREE.Mesh(new THREE.SphereGeometry(0.5,32,16,0,Math.PI*2,0,Math.PI/2), mat.clone())
      top.position.y = 0.75
      const bot  = new THREE.Mesh(new THREE.SphereGeometry(0.5,32,16,0,Math.PI*2,Math.PI/2,Math.PI/2), mat.clone())
      bot.position.y = -0.75
      group.add(body, top, bot)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pump': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.2,0.4,0.8), mat.clone())
      base.position.y = -0.3
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.3,0.35,0.7,32), mat.clone())
      body.position.y = 0.25
      const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.1,0.1,0.5,12), mat.clone())
      pipe.rotation.z = Math.PI/2
      pipe.position.set(0.5,0.2,0)
      group.add(base, body, pipe)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pipe': {
      geo = new THREE.CylinderGeometry(0.12, 0.12, 2.5, 16)
      break
    }
    case 'valve': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.25,0.25,0.3,32), mat.clone())
      const stem = new THREE.Mesh(new THREE.CylinderGeometry(0.06,0.06,0.6,12), mat.clone())
      stem.position.y = 0.45
      const hand = new THREE.Mesh(new THREE.BoxGeometry(0.6,0.06,0.08), mat.clone())
      hand.position.y = 0.75
      group.add(body, stem, hand)
      group.userData.isMeshGroup = true
      return group
    }
    case 'motor': {
      const group = new THREE.Group()
      const bodyMat = mat.clone()
      const housing = new THREE.Mesh(new THREE.BoxGeometry(1.6, 0.7, 0.7), bodyMat)
      const endCap1 = new THREE.Mesh(new THREE.CylinderGeometry(0.32, 0.32, 0.1, 24), bodyMat)
      endCap1.rotation.z = Math.PI/2; endCap1.position.x = 0.85
      const endCap2 = new THREE.Mesh(new THREE.CylinderGeometry(0.32, 0.32, 0.1, 24), bodyMat)
      endCap2.rotation.z = Math.PI/2; endCap2.position.x = -0.85
      const shaft = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.5, 12), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.8, roughness:0.2}))
      shaft.rotation.z = Math.PI/2; shaft.position.x = 1.1
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.15, 0.9), bodyMat)
      base.position.y = -0.42
      const foot1 = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.9), bodyMat)
      foot1.position.set(0.6, -0.55, 0)
      const foot2 = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.9), bodyMat)
      foot2.position.set(-0.6, -0.55, 0)
      group.add(housing, endCap1, endCap2, shaft, base, foot1, foot2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'conveyor': {
      const group = new THREE.Group()
      const frameColor = new THREE.MeshStandardMaterial({color:0x555555, metalness:0.6, roughness:0.4})
      const beltColor = new THREE.MeshStandardMaterial({color:0x222222, metalness:0.1, roughness:0.9})
      const sideL = new THREE.Mesh(new THREE.BoxGeometry(3, 0.08, 0.06), frameColor)
      sideL.position.set(0, 0.1, 0.35)
      const sideR = new THREE.Mesh(new THREE.BoxGeometry(3, 0.08, 0.06), frameColor)
      sideR.position.set(0, 0.1, -0.35)
      const roller1 = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.7, 16), frameColor)
      roller1.rotation.x = Math.PI/2; roller1.position.set(1.4, 0.1, 0)
      const roller2 = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.7, 16), frameColor)
      roller2.rotation.x = Math.PI/2; roller2.position.set(-1.4, 0.1, 0)
      const belt = new THREE.Mesh(new THREE.BoxGeometry(2.6, 0.02, 0.6), beltColor)
      belt.position.y = 0.22
      const leg1 = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.4, 0.06), frameColor)
      leg1.position.set(1.2, -0.15, 0.3)
      const leg2 = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.4, 0.06), frameColor)
      leg2.position.set(1.2, -0.15, -0.3)
      const leg3 = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.4, 0.06), frameColor)
      leg3.position.set(-1.2, -0.15, 0.3)
      const leg4 = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.4, 0.06), frameColor)
      leg4.position.set(-1.2, -0.15, -0.3)
      group.add(sideL, sideR, roller1, roller2, belt, leg1, leg2, leg3, leg4)
      group.userData.isMeshGroup = true
      return group
    }
    case 'chimney': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.35, 2.5, 24), mat.clone())
      body.position.y = 1.25
      const rim = new THREE.Mesh(new THREE.TorusGeometry(0.22, 0.04, 8, 24), mat.clone())
      rim.position.y = 2.5; rim.rotation.x = Math.PI/2
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.5, 0.2, 24), mat.clone())
      base.position.y = 0.1
      group.add(body, rim, base)
      group.userData.isMeshGroup = true
      return group
    }
    case 'heatex': {
      const group = new THREE.Group()
      const shell = new THREE.Mesh(new THREE.CylinderGeometry(0.45, 0.45, 2, 24), mat.clone())
      shell.rotation.z = Math.PI/2
      const cap1 = new THREE.Mesh(new THREE.SphereGeometry(0.45, 24, 16, 0, Math.PI*2, 0, Math.PI/2), mat.clone())
      cap1.rotation.z = -Math.PI/2; cap1.position.x = 1
      const cap2 = new THREE.Mesh(new THREE.SphereGeometry(0.45, 24, 16, 0, Math.PI*2, 0, Math.PI/2), mat.clone())
      cap2.rotation.z = Math.PI/2; cap2.position.x = -1
      const tubeMat = new THREE.MeshStandardMaterial({color:0xcc8844, metalness:0.6, roughness:0.3})
      for (let i = 0; i < 4; i++) {
        const angle = (i / 4) * Math.PI * 2
        const tube = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.05, 1.8, 8), tubeMat)
        tube.rotation.z = Math.PI/2
        tube.position.set(0, Math.cos(angle) * 0.25, Math.sin(angle) * 0.25)
        group.add(tube)
      }
      group.add(shell, cap1, cap2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'fan': {
      const group = new THREE.Group()
      const housing = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.5, 0.4, 24, 1, true), mat.clone())
      housing.rotation.x = Math.PI/2
      const bladeMat = new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.5, roughness:0.3, side: THREE.DoubleSide})
      for (let i = 0; i < 4; i++) {
        const blade = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.02, 0.15), bladeMat)
        blade.rotation.y = (i / 4) * Math.PI * 2
        group.add(blade)
      }
      const hub = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.15, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8, roughness:0.2}))
      hub.rotation.x = Math.PI/2
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.35, 0.45, 0.3, 24), mat.clone())
      inlet.rotation.x = Math.PI/2; inlet.position.z = 0.35
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.12, 0.8), mat.clone())
      base.position.y = -0.35
      const leg1 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.25, 0.08), mat.clone())
      leg1.position.set(0.3, -0.2, 0.3)
      const leg2 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.25, 0.08), mat.clone())
      leg2.position.set(-0.3, -0.2, 0.3)
      const leg3 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.25, 0.08), mat.clone())
      leg3.position.set(0.3, -0.2, -0.3)
      const leg4 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.25, 0.08), mat.clone())
      leg4.position.set(-0.3, -0.2, -0.3)
      group.add(housing, hub, inlet, base, leg1, leg2, leg3, leg4)
      group.userData.isMeshGroup = true
      return group
    }
    case 'reactor': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.6, 0.6, 1.5, 32), mat.clone())
      body.position.y = 0.5
      const bottom = new THREE.Mesh(new THREE.SphereGeometry(0.6, 32, 16, 0, Math.PI*2, Math.PI/2, Math.PI/2), mat.clone())
      bottom.position.y = -0.25
      const lid = new THREE.Mesh(new THREE.CylinderGeometry(0.65, 0.65, 0.1, 32), mat.clone())
      lid.position.y = 1.3
      const motor = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.3, 0.3), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.5, roughness:0.4}))
      motor.position.y = 1.65
      const shaft = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 1.2, 8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.8, roughness:0.2}))
      shaft.position.y = 0.7
      const blade1 = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.04, 0.08), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.6, roughness:0.3}))
      blade1.position.y = 0.1
      const blade2 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.04, 0.8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.6, roughness:0.3}))
      blade2.position.y = 0.1
      for (let i = 0; i < 3; i++) {
        const angle = (i / 3) * Math.PI * 2
        const leg = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.5, 0.08), mat.clone())
        leg.position.set(Math.cos(angle) * 0.5, -0.55, Math.sin(angle) * 0.5)
        group.add(leg)
      }
      group.add(body, bottom, lid, motor, shaft, blade1, blade2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'mixer': {
      const group = new THREE.Group()
      const container = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.15, 0.8, 24), mat.clone())
      container.position.y = 0.2
      const shaft = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 1.2, 8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.8, roughness:0.2}))
      shaft.position.y = 0.8
      const blade1 = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.03, 0.06), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.6, roughness:0.3}))
      blade1.position.y = 0.1
      const blade2 = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.03, 0.6), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.6, roughness:0.3}))
      blade2.position.y = 0.1
      const handle = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.4, 8), new THREE.MeshStandardMaterial({color:0x333333, metalness:0.5, roughness:0.4}))
      handle.rotation.z = Math.PI/2; handle.position.set(0.3, 1.4, 0)
      group.add(container, shaft, blade1, blade2, handle)
      group.userData.isMeshGroup = true
      return group
    }
    case 'coolingtower': {
      const group = new THREE.Group()
      const towerMat = mat.clone()
      const sections = 8
      for (let i = 0; i < sections; i++) {
        const t = i / sections
        const r = 0.6 + 0.3 * Math.sin(t * Math.PI)
        const section = new THREE.Mesh(new THREE.CylinderGeometry(
          0.6 + 0.3 * Math.sin((i+1)/sections * Math.PI),
          r, 2/sections, 24
        ), towerMat)
        section.position.y = t * 2
        group.add(section)
      }
      const pool = new THREE.Mesh(new THREE.CylinderGeometry(0.95, 0.95, 0.15, 24), new THREE.MeshStandardMaterial({color:0x4488aa, metalness:0.1, roughness:0.6}))
      pool.position.y = -0.07
      group.add(pool)
      group.userData.isMeshGroup = true
      return group
    }
    case 'flange': {
      const group = new THREE.Group()
      const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.6, 16), mat.clone())
      const disc1 = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.25, 0.06, 24), mat.clone())
      disc1.position.y = 0.3
      const disc2 = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.25, 0.06, 24), mat.clone())
      disc2.position.y = -0.3
      const boltMat = new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8, roughness:0.2})
      for (let i = 0; i < 6; i++) {
        const angle = (i / 6) * Math.PI * 2
        const b1 = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.08, 8), boltMat)
        b1.position.set(Math.cos(angle)*0.19, 0.3, Math.sin(angle)*0.19)
        const b2 = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.08, 8), boltMat)
        b2.position.set(Math.cos(angle)*0.19, -0.3, Math.sin(angle)*0.19)
        group.add(b1, b2)
      }
      group.add(pipe, disc1, disc2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'gauge': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.CylinderGeometry(0.3, 0.3, 0.08, 32), new THREE.MeshStandardMaterial({color:0x333333, metalness:0.6, roughness:0.3}))
      frame.rotation.x = Math.PI/2
      const dial = new THREE.Mesh(new THREE.CircleGeometry(0.27, 32), new THREE.MeshStandardMaterial({color:0xffffff, roughness:0.8, side: THREE.DoubleSide}))
      dial.position.z = 0.04; dial.rotation.x = -Math.PI/2
      const needle = new THREE.Mesh(new THREE.BoxGeometry(0.02, 0.22, 0.01), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.3}))
      needle.position.set(0.06, 0.05, 0); needle.rotation.z = -Math.PI/4
      const center = new THREE.Mesh(new THREE.CylinderGeometry(0.025, 0.025, 0.03, 12), new THREE.MeshStandardMaterial({color:0x222222, metalness:0.8}))
      center.rotation.x = Math.PI/2; center.position.z = 0.05
      const conn = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.2, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      conn.position.z = -0.14
      group.add(frame, dial, needle, center, conn)
      group.userData.isMeshGroup = true
      return group
    }
    case 'sensor': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.3, 16), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.4, roughness:0.5}))
      body.position.y = 0.15
      const probe = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.25, 8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.8, roughness:0.2}))
      probe.position.y = -0.12
      const led = new THREE.Mesh(new THREE.SphereGeometry(0.025, 8, 8), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.8}))
      led.position.set(0.1, 0.25, 0)
      const cable = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.4, 6), new THREE.MeshStandardMaterial({color:0x222222}))
      cable.position.y = 0.5
      group.add(body, probe, led, cable)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pipeelbow': {
      const group = new THREE.Group()
      const r = 0.1
      const p0 = new THREE.Vector3(0, -0.45, 0)
      const p1 = new THREE.Vector3(0, 0.22, 0)
      const p2 = new THREE.Vector3(0.67, 0.22, 0)
      const curve = createRoundedPipeCurve([p0, p1, p2], r)
      const elbow = new THREE.Mesh(new THREE.TubeGeometry(curve, 36, r, 16, false), mat.clone())
      group.add(elbow)
      group.userData.isMeshGroup = true
      return group
    }
    case 'hopper': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.CylinderGeometry(0.6, 0.3, 0.6, 24), mat.clone())
      top.position.y = 0.5
      const bottom = new THREE.Mesh(new THREE.CylinderGeometry(0.3, 0.1, 0.5, 24), mat.clone())
      bottom.position.y = -0.05
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.15, 12), mat.clone())
      outlet.position.y = -0.37
      for (let i = 0; i < 4; i++) {
        const angle = (i / 4) * Math.PI * 2 + Math.PI/4
        const leg = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.8, 0.06), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5, roughness:0.4}))
        leg.position.set(Math.cos(angle)*0.5, -0.6, Math.sin(angle)*0.5)
        group.add(leg)
      }
      group.add(top, bottom, outlet)
      group.userData.isMeshGroup = true
      return group
    }
    case 'silo': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.7, 0.7, 2, 24), mat.clone())
      body.position.y = 1
      const roof = new THREE.Mesh(new THREE.ConeGeometry(0.75, 0.5, 24), mat.clone())
      roof.position.y = 2.25
      const cone = new THREE.Mesh(new THREE.CylinderGeometry(0.7, 0.15, 0.8, 24), mat.clone())
      cone.position.y = 0.1
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.2, 12), mat.clone())
      outlet.position.y = -0.4
      group.add(body, roof, cone, outlet)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pressurevessel': {
      const group = new THREE.Group()
      const shell = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.5, 2.2, 24), mat.clone())
      shell.rotation.z = Math.PI/2
      const cap1 = new THREE.Mesh(new THREE.SphereGeometry(0.5, 24, 16), mat.clone())
      cap1.position.x = 1.1
      const cap2 = new THREE.Mesh(new THREE.SphereGeometry(0.5, 24, 16), mat.clone())
      cap2.position.x = -1.1
      const nozzle1 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 12), mat.clone())
      nozzle1.position.set(0.3, 0.55, 0)
      const nozzle2 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 12), mat.clone())
      nozzle2.position.set(-0.3, 0.55, 0)
      const saddle1 = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.3, 0.8), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5, roughness:0.4}))
      saddle1.position.set(0.6, -0.55, 0)
      const saddle2 = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.3, 0.8), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5, roughness:0.4}))
      saddle2.position.set(-0.6, -0.55, 0)
      group.add(shell, cap1, cap2, nozzle1, nozzle2, saddle1, saddle2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'separator': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.4, 1.8, 24), mat.clone())
      body.position.y = 0.9
      const top = new THREE.Mesh(new THREE.SphereGeometry(0.4, 24, 16, 0, Math.PI*2, 0, Math.PI/2), mat.clone())
      top.position.y = 1.8
      const bot = new THREE.Mesh(new THREE.SphereGeometry(0.4, 24, 16, 0, Math.PI*2, Math.PI/2, Math.PI/2), mat.clone())
      bot.position.y = 0
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 12), mat.clone())
      inlet.rotation.z = Math.PI/2; inlet.position.set(0.5, 1.4, 0)
      const outlet1 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 12), mat.clone())
      outlet1.rotation.z = Math.PI/2; outlet1.position.set(0.5, 0.6, 0)
      const drain = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.25, 8), mat.clone())
      drain.position.y = -0.25
      const baffle = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.04, 0.04), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      baffle.position.y = 1.1
      group.add(body, top, bot, inlet, outlet1, drain, baffle)
      group.userData.isMeshGroup = true
      return group
    }
    case 'boiler': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(1.4, 1.6, 1), mat.clone())
      body.position.y = 0.8
      const chimney = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.6, 12), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.5, roughness:0.4}))
      chimney.position.set(0.3, 1.9, 0)
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.4, 0.02), new THREE.MeshStandardMaterial({color:0x222222, metalness:0.4}))
      door.position.set(0, 0.5, 0.51)
      const pipe1 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.5, 8), mat.clone())
      pipe1.rotation.z = Math.PI/2; pipe1.position.set(-0.95, 1.2, 0)
      const pipe2 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.5, 8), mat.clone())
      pipe2.rotation.z = Math.PI/2; pipe2.position.set(-0.95, 0.5, 0)
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.6, 0.1, 1.2), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      base.position.y = 0.05
      group.add(body, chimney, door, pipe1, pipe2, base)
      group.userData.isMeshGroup = true
      return group
    }
    case 'compressor': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.6, 0.7), mat.clone())
      body.position.y = 0.4
      const cylinder = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.5, 16), mat.clone())
      cylinder.position.set(0, 0.95, 0)
      const head = new THREE.Mesh(new THREE.CylinderGeometry(0.22, 0.22, 0.06, 16), mat.clone())
      head.position.y = 1.22
      const wheel = new THREE.Mesh(new THREE.TorusGeometry(0.25, 0.05, 8, 24), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7, roughness:0.2}))
      wheel.position.set(-0.6, 0.4, 0); wheel.rotation.y = Math.PI/2
      const pipe1 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.4, 8), mat.clone())
      pipe1.rotation.z = Math.PI/2; pipe1.position.set(0.7, 0.9, 0)
      const pipe2 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.4, 8), mat.clone())
      pipe2.rotation.z = Math.PI/2; pipe2.position.set(0.7, 0.5, 0)
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.4, 0.12, 0.9), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      base.position.y = 0.06
      group.add(body, cylinder, head, wheel, pipe1, pipe2, base)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 阀门管件新增 =====
    case 'teepipe': {
      const group = new THREE.Group()
      const vPipe = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 1.2, 12), mat.clone())
      const hPipe = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.8, 12), mat.clone())
      hPipe.rotation.z = Math.PI/2; hPipe.position.y = 0.4
      const collars = [
        {pos: [0, -0.6, 0]}, {pos: [0, 0.6, 0]}, {pos: [0.4, 0.4, 0]}
      ]
      collars.forEach(function(c) {
        const collar = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.14, 0.08, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
        collar.position.set(c.pos[0], c.pos[1], c.pos[2])
        group.add(collar)
      })
      group.add(vPipe, hPipe)
      group.userData.isMeshGroup = true
      return group
    }
    case 'crosspipe': {
      const group = new THREE.Group()
      const vPipe = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 1.2, 12), mat.clone())
      const hPipe1 = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.8, 12), mat.clone())
      hPipe1.rotation.z = Math.PI/2; hPipe1.position.y = 0.4
      const hPipe2 = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.8, 12), mat.clone())
      hPipe2.rotation.x = Math.PI/2; hPipe2.position.y = 0.4
      const collars = [
        {pos: [0, -0.6, 0]}, {pos: [0, 0.6, 0]}, {pos: [0.4, 0.4, 0]}, {pos: [-0.4, 0.4, 0]}, {pos: [0, 0.4, 0.4]}, {pos: [0, 0.4, -0.4]}
      ]
      collars.forEach(function(c) {
        const collar = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.14, 0.08, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
        collar.position.set(c.pos[0], c.pos[1], c.pos[2])
        group.add(collar)
      })
      group.add(vPipe, hPipe1, hPipe2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'reducerpipe': {
      const group = new THREE.Group()
      const cone = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.08, 1.0, 16), mat.clone())
      const topFlange = new THREE.Mesh(new THREE.CylinderGeometry(0.18, 0.18, 0.06, 16), mat.clone())
      topFlange.position.y = 0.5
      const botFlange = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.06, 16), mat.clone())
      botFlange.position.y = -0.5
      group.add(cone, topFlange, botFlange)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pipecap': {
      const group = new THREE.Group()
      const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.5, 16), mat.clone())
      pipe.position.y = 0.15
      const cap = new THREE.Mesh(new THREE.SphereGeometry(0.13, 16, 12, 0, Math.PI*2, 0, Math.PI/2), mat.clone())
      cap.position.y = 0.45
      const rim = new THREE.Mesh(new THREE.TorusGeometry(0.11, 0.02, 8, 16), mat.clone())
      rim.position.y = 0.4; rim.rotation.x = Math.PI/2
      group.add(pipe, cap, rim)
      group.userData.isMeshGroup = true
      return group
    }
    case 'ballvalve': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.22, 0.22, 0.6, 24), mat.clone())
      body.rotation.z = Math.PI/2
      const ball = new THREE.Mesh(new THREE.SphereGeometry(0.16, 16, 16), mat.clone())
      const neck = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.05, 0.3, 8), mat.clone())
      neck.position.y = 0.3
      const handle = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.5, 8), new THREE.MeshStandardMaterial({color:0xcc0000, metalness:0.5, roughness:0.3}))
      handle.rotation.z = Math.PI/2; handle.position.set(0, 0.45, 0)
      const handleBall1 = new THREE.Mesh(new THREE.SphereGeometry(0.04, 8, 8), new THREE.MeshStandardMaterial({color:0xcc0000}))
      handleBall1.position.set(0.25, 0.45, 0)
      const handleBall2 = new THREE.Mesh(new THREE.SphereGeometry(0.04, 8, 8), new THREE.MeshStandardMaterial({color:0xcc0000}))
      handleBall2.position.set(-0.25, 0.45, 0)
      group.add(body, ball, neck, handle, handleBall1, handleBall2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'butterflyvalve': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.3, 24), mat.clone())
      const disc = new THREE.Mesh(new THREE.CylinderGeometry(0.18, 0.18, 0.04, 24), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.5}))
      disc.rotation.z = Math.PI*0.15
      const stem = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.5, 8), mat.clone())
      stem.position.y = 0.3
      const handwheel = new THREE.Mesh(new THREE.TorusGeometry(0.14, 0.03, 8, 16), mat.clone())
      handwheel.position.y = 0.55; handwheel.rotation.x = Math.PI/2
      const hub = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.05, 0.08, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      hub.position.y = 0.56
      group.add(body, disc, stem, handwheel, hub)
      group.userData.isMeshGroup = true
      return group
    }
    case 'checkvalve': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.5, 24), mat.clone())
      body.rotation.z = Math.PI/2
      const arrowHead = new THREE.Mesh(new THREE.ConeGeometry(0.1, 0.25, 8), new THREE.MeshStandardMaterial({color:0x44ff44, metalness:0.3}))
      arrowHead.rotation.z = Math.PI/2; arrowHead.position.set(0.35, 0.08, 0)
      const arrowShaft = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.3, 8), new THREE.MeshStandardMaterial({color:0x44ff44}))
      arrowShaft.rotation.z = Math.PI/2; arrowShaft.position.set(0.15, 0.08, 0)
      const cap = new THREE.Mesh(new THREE.CylinderGeometry(0.22, 0.22, 0.05, 24), mat.clone())
      cap.position.x = -0.26
      group.add(body, arrowHead, arrowShaft, cap)
      group.userData.isMeshGroup = true
      return group
    }
    case 'safetyvalve': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.18, 0.18, 0.4, 24), mat.clone())
      const spring = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.3, 12), new THREE.MeshStandardMaterial({color:0xff4444, metalness:0.4}))
      spring.position.y = 0.35
      const top = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.08, 16), mat.clone())
      top.position.y = 0.55
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.3, 8), mat.clone())
      outlet.rotation.z = Math.PI/2; outlet.position.set(0.25, 0.2, 0)
      const pin = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.15, 6), new THREE.MeshStandardMaterial({color:0xff4444}))
      pin.position.y = 0.6
      group.add(body, spring, top, outlet, pin)
      group.userData.isMeshGroup = true
      return group
    }
    case 'gatevalve': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.5, 0.35), mat.clone())
      const bonnet = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.2, 0.3, 16), mat.clone())
      bonnet.position.y = 0.4
      const handwheel = new THREE.Mesh(new THREE.TorusGeometry(0.2, 0.04, 8, 24), mat.clone())
      handwheel.position.y = 0.8; handwheel.rotation.x = Math.PI/2
      const stem = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.35, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      stem.position.y = 0.65
      const pipeL = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.3, 12), mat.clone())
      pipeL.rotation.z = Math.PI/2; pipeL.position.x = -0.35
      const pipeR = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.3, 12), mat.clone())
      pipeR.rotation.z = Math.PI/2; pipeR.position.x = 0.35
      group.add(body, bonnet, handwheel, stem, pipeL, pipeR)
      group.userData.isMeshGroup = true
      return group
    }
    case 'centrifugalpump': {
      const group = new THREE.Group()
      const volute = new THREE.Mesh(new THREE.SphereGeometry(0.35, 24, 16, 0, Math.PI*2, 0, Math.PI*0.6), mat.clone())
      volute.scale.set(1, 0.7, 1)
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.3, 12), mat.clone())
      inlet.position.set(0, 0.2, 0.2)
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 8), mat.clone())
      outlet.rotation.z = Math.PI/2; outlet.position.set(0.4, -0.05, 0)
      const motorHousing = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.5, 16), mat.clone())
      motorHousing.rotation.z = Math.PI/2; motorHousing.position.x = -0.5
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.1, 0.6), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      base.position.y = -0.25
      const coupling = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.15, 16), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      coupling.rotation.z = Math.PI/2; coupling.position.x = -0.15
      group.add(volute, inlet, outlet, motorHousing, base, coupling)
      group.userData.isMeshGroup = true
      return group
    }
    case 'gearpump': {
      const group = new THREE.Group()
      const housing = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.4, 0.5), mat.clone())
      const flangeL = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.15, 16), mat.clone())
      flangeL.rotation.z = Math.PI/2; flangeL.position.x = -0.35
      const flangeR = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.15, 16), mat.clone())
      flangeR.rotation.z = Math.PI/2; flangeR.position.x = 0.35
      const shaft = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.3, 8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.7}))
      shaft.position.set(-0.55, 0, 0)
      const gear1 = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.14, 0.08, 24), mat.clone())
      gear1.rotation.x = Math.PI/2; gear1.position.set(0.1, 0, 0.15)
      const gear2 = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.14, 0.08, 24), mat.clone())
      gear2.rotation.x = Math.PI/2; gear2.position.set(0.1, 0, -0.15)
      group.add(housing, flangeL, flangeR, shaft, gear1, gear2)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 容器储罐新增 =====
    case 'buffertank': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.45, 0.45, 1.4, 24), mat.clone())
      body.position.y = 0.7
      const top = new THREE.Mesh(new THREE.SphereGeometry(0.45, 24, 16, 0, Math.PI*2, 0, Math.PI/2), mat.clone())
      top.position.y = 1.4
      const bot = new THREE.Mesh(new THREE.SphereGeometry(0.45, 24, 16, 0, Math.PI*2, Math.PI/2, Math.PI/2), mat.clone())
      bot.position.y = 0
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.25, 8), mat.clone())
      inlet.position.set(0.4, 1.2, 0)
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.25, 8), mat.clone())
      outlet.position.set(0.4, 0.3, 0)
      for (let i = 0; i < 3; i++) {
        const leg = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.4, 0.06), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
        leg.position.set(Math.cos(i/3*Math.PI*2)*0.45, -0.3, Math.sin(i/3*Math.PI*2)*0.45)
        group.add(leg)
      }
      group.add(body, top, bot, inlet, outlet)
      group.userData.isMeshGroup = true
      return group
    }
    case 'mixingtank': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.4, 1.0, 24), mat.clone())
      body.position.y = 0.5
      const bottom = new THREE.Mesh(new THREE.SphereGeometry(0.4, 24, 16, 0, Math.PI*2, Math.PI/2, Math.PI/2), mat.clone())
      bottom.position.y = 0
      const topLid = new THREE.Mesh(new THREE.CylinderGeometry(0.42, 0.42, 0.08, 24), mat.clone())
      topLid.position.y = 1.04
      const motor = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.3, 16), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.5}))
      motor.position.y = 1.4
      const shaft = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.9, 8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.7}))
      shaft.position.y = 0.5
      const blade = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.03, 0.12), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.5}))
      blade.position.y = 0.25
      group.add(body, bottom, topLid, motor, shaft, blade)
      group.userData.isMeshGroup = true
      return group
    }
    case 'airreceiver': {
      const group = new THREE.Group()
      const shell = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.4, 1.6, 24), mat.clone())
      shell.rotation.z = Math.PI/2
      const cap1 = new THREE.Mesh(new THREE.SphereGeometry(0.4, 24, 16), mat.clone())
      cap1.position.x = 0.8
      const cap2 = new THREE.Mesh(new THREE.SphereGeometry(0.4, 24, 16), mat.clone())
      cap2.position.x = -0.8
      const nozzle1 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.25, 8), mat.clone())
      nozzle1.position.set(0.2, 0.45, 0)
      const nozzle2 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.25, 8), mat.clone())
      nozzle2.position.set(-0.2, 0.45, 0)
      const saddle1 = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.35, 0.7), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      saddle1.position.set(0.4, -0.5, 0)
      const saddle2 = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.35, 0.7), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      saddle2.position.set(-0.4, -0.5, 0)
      group.add(shell, cap1, cap2, nozzle1, nozzle2, saddle1, saddle2)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 换热分离新增 =====
    case 'condenser': {
      const group = new THREE.Group()
      const shell = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.4, 1.8, 24), mat.clone())
      shell.rotation.z = Math.PI/2
      const cap1 = new THREE.Mesh(new THREE.CylinderGeometry(0.42, 0.42, 0.2, 24), mat.clone())
      cap1.position.x = 0.95
      const cap2 = new THREE.Mesh(new THREE.CylinderGeometry(0.42, 0.42, 0.2, 24), mat.clone())
      cap2.position.x = -0.95
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 8), mat.clone())
      inlet.position.set(0, 0.48, 0.6)
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.3, 8), mat.clone())
      outlet.position.set(0, -0.35, -0.6)
      for (let i = 0; i < 6; i++) {
        const angle = (i / 6) * Math.PI * 2
        const tube = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 1.6, 8), new THREE.MeshStandardMaterial({color:0xcc8844, metalness:0.6}))
        tube.rotation.z = Math.PI/2
        tube.position.set(Math.cos(angle)*0.25, Math.sin(angle)*0.25, 0)
        group.add(tube)
      }
      group.add(shell, cap1, cap2, inlet, outlet)
      group.userData.isMeshGroup = true
      return group
    }
    case 'evaporator': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.4, 1.5, 24), mat.clone())
      body.position.y = 0.75
      const top = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.4, 0.3, 24), mat.clone())
      top.position.y = 1.65
      const vaporOutlet = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 8), mat.clone())
      vaporOutlet.position.y = 1.85
      const feedInlet = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.3, 8), mat.clone())
      feedInlet.rotation.z = Math.PI/2; feedInlet.position.set(0.5, 1.0, 0)
      const bottom = new THREE.Mesh(new THREE.SphereGeometry(0.4, 24, 16, 0, Math.PI*2, Math.PI/2, Math.PI/2), mat.clone())
      bottom.position.y = 0
      const drain = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.25, 8), mat.clone())
      drain.position.y = -0.3
      for (let i = 0; i < 3; i++) {
        const leg = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.35, 0.06), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
        leg.position.set(Math.cos(i/3*Math.PI*2)*0.4, -0.28, Math.sin(i/3*Math.PI*2)*0.4)
        group.add(leg)
      }
      group.add(body, top, vaporOutlet, feedInlet, bottom, drain)
      group.userData.isMeshGroup = true
      return group
    }
    case 'cyclone': {
      const group = new THREE.Group()
      const cylinder = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.25, 0.6, 24), mat.clone())
      cylinder.position.y = 0.9
      const cone = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.06, 0.9, 24), mat.clone())
      cone.position.y = 0.15
      const inlet = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.4, 0.3), mat.clone())
      inlet.position.set(0.35, 0.7, 0)
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.25, 8), mat.clone())
      outlet.position.y = -0.35
      const topPipe = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.4, 12), mat.clone())
      topPipe.position.y = 1.4
      group.add(cylinder, cone, inlet, outlet, topPipe)
      group.userData.isMeshGroup = true
      return group
    }
    case 'filterunit': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.3, 0.3, 1.0, 24), mat.clone())
      body.position.y = 0.5
      const topFlange = new THREE.Mesh(new THREE.CylinderGeometry(0.35, 0.35, 0.08, 24), mat.clone())
      topFlange.position.y = 1.04
      const botFlange = new THREE.Mesh(new THREE.CylinderGeometry(0.35, 0.35, 0.08, 24), mat.clone())
      botFlange.position.y = -0.04
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 8), mat.clone())
      inlet.rotation.z = Math.PI/2; inlet.position.set(0.4, 0.6, 0)
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 8), mat.clone())
      outlet.rotation.z = Math.PI/2; outlet.position.set(0.4, 0.2, 0)
      const element1 = new THREE.Mesh(new THREE.TorusGeometry(0.28, 0.02, 6, 24), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      element1.position.y = 0.7; element1.rotation.x = Math.PI/2
      const element2 = new THREE.Mesh(new THREE.TorusGeometry(0.28, 0.02, 6, 24), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      element2.position.y = 0.4; element2.rotation.x = Math.PI/2
      group.add(body, topFlange, botFlange, inlet, outlet, element1, element2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'dustcollector': {
      const group = new THREE.Group()
      const box = new THREE.Mesh(new THREE.BoxGeometry(1.0, 1.0, 0.8), mat.clone())
      box.position.y = 0.5
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.4, 12), mat.clone())
      inlet.rotation.z = Math.PI/2; inlet.position.set(-0.7, 0.5, 0)
      const hopper = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.1, 0.6, 24), mat.clone())
      hopper.position.y = -0.1
      const dustOutlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.2, 8), mat.clone())
      dustOutlet.position.y = -0.5
      const bag1 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.6, 8), new THREE.MeshStandardMaterial({color:0x4488ff, metalness:0.1, roughness:0.8}))
      bag1.position.set(-0.2, 0.8, 0.2)
      const bag2 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.6, 8), new THREE.MeshStandardMaterial({color:0x4488ff, metalness:0.1, roughness:0.8}))
      bag2.position.set(0.2, 0.8, 0.2)
      const bag3 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.6, 8), new THREE.MeshStandardMaterial({color:0x4488ff, metalness:0.1, roughness:0.8}))
      bag3.position.set(0, 0.8, -0.2)
      group.add(box, inlet, hopper, dustOutlet, bag1, bag2, bag3)
      group.userData.isMeshGroup = true
      return group
    }
    case 'scrubber': {
      const group = new THREE.Group()
      const towerBody = new THREE.Mesh(new THREE.CylinderGeometry(0.35, 0.35, 2.2, 24), mat.clone())
      towerBody.position.y = 1.1
      const gasInlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.3, 8), mat.clone())
      gasInlet.rotation.z = Math.PI/2; gasInlet.position.set(0.45, 0.5, 0)
      const gasOutlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.3, 8), mat.clone())
      gasOutlet.rotation.z = Math.PI/2; gasOutlet.position.set(0.45, 1.8, 0)
      const top = new THREE.Mesh(new THREE.CylinderGeometry(0.38, 0.38, 0.1, 24), mat.clone())
      top.position.y = 2.25
      const bot = new THREE.Mesh(new THREE.CylinderGeometry(0.38, 0.38, 0.1, 24), mat.clone())
      bot.position.y = 0.05
      const packing1 = new THREE.Mesh(new THREE.TorusGeometry(0.34, 0.015, 6, 24), new THREE.MeshStandardMaterial({color:0x66cc66, metalness:0.3}))
      packing1.position.y = 0.9; packing1.rotation.x = Math.PI/2
      const packing2 = new THREE.Mesh(new THREE.TorusGeometry(0.34, 0.015, 6, 24), new THREE.MeshStandardMaterial({color:0x66cc66, metalness:0.3}))
      packing2.position.y = 1.3; packing2.rotation.x = Math.PI/2
      group.add(towerBody, gasInlet, gasOutlet, top, bot, packing1, packing2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'dryer': {
      const group = new THREE.Group()
      const drum = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.4, 1.8, 24), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.3, roughness:0.5}))
      drum.rotation.z = Math.PI/2
      const rim1 = new THREE.Mesh(new THREE.TorusGeometry(0.42, 0.06, 8, 24), mat.clone())
      rim1.position.x = 0.6; rim1.rotation.y = Math.PI/2
      const rim2 = new THREE.Mesh(new THREE.TorusGeometry(0.42, 0.06, 8, 24), mat.clone())
      rim2.position.x = -0.6; rim2.rotation.y = Math.PI/2
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.3, 8), mat.clone())
      inlet.position.set(0.9, 0.2, 0)
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.3, 8), mat.clone())
      outlet.position.set(-0.9, -0.2, 0)
      const leg1 = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.5, 0.1), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      leg1.position.set(0.5, -0.35, 0.3)
      const leg2 = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.5, 0.1), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      leg2.position.set(-0.5, -0.35, 0.3)
      const leg3 = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.5, 0.1), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      leg3.position.set(0.5, -0.35, -0.3)
      const leg4 = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.5, 0.1), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      leg4.position.set(-0.5, -0.35, -0.3)
      group.add(drum, rim1, rim2, inlet, outlet, leg1, leg2, leg3, leg4)
      group.userData.isMeshGroup = true
      return group
    }
    case 'furnace': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(1.5, 1.0, 1.2), mat.clone())
      body.position.y = 0.5
      const chimney = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 1.0, 16), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.5, roughness:0.4}))
      chimney.position.set(0.3, 1.3, 0)
      const fireDoor = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.35, 0.03), new THREE.MeshStandardMaterial({color:0x222222, metalness:0.4}))
      fireDoor.position.set(0, 0.3, 0.61)
      const flame = new THREE.Mesh(new THREE.ConeGeometry(0.12, 0.25, 8), new THREE.MeshStandardMaterial({color:0xff6600, emissive:0xff4400, emissiveIntensity:0.6}))
      flame.position.set(0, 0.25, 0.65)
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.7, 0.1, 1.4), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      base.position.y = 0.05
      group.add(body, chimney, fireDoor, flame, base)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 输送存储新增 =====
    case 'screwconveyor': {
      const group = new THREE.Group()
      const trough = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 2.5, 16, 1, true, 0, Math.PI), mat.clone())
      trough.rotation.z = Math.PI/2
      trough.rotation.y = Math.PI/2
      const shaft = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 2.5, 8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.7}))
      shaft.rotation.z = Math.PI/2
      for (let i = 0; i < 10; i++) {
        const flight = new THREE.Mesh(new THREE.BoxGeometry(0.02, 0.35, 0.15), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.6}))
        flight.position.set(i * 0.25 - 1.1, 0, 0)
        flight.rotation.x = (i % 2) * Math.PI/6
        group.add(flight)
      }
      const motor = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.3, 0.3), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.5}))
      motor.position.x = -1.55
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.3, 0.08), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      base.position.set(-1.1, -0.25, 0)
      const base2 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.3, 0.08), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      base2.position.set(1.1, -0.25, 0)
      group.add(trough, shaft, motor, base, base2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'bucketelevator': {
      const group = new THREE.Group()
      const frameL = new THREE.Mesh(new THREE.BoxGeometry(0.06, 2.0, 0.06), mat.clone())
      frameL.position.set(-0.3, 1, 0)
      const frameR = new THREE.Mesh(new THREE.BoxGeometry(0.06, 2.0, 0.06), mat.clone())
      frameR.position.set(0.3, 1, 0)
      const cross1 = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.04, 0.04), mat.clone())
      cross1.position.set(0, 0.5, 0)
      const cross2 = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.04, 0.04), mat.clone())
      cross2.position.set(0, 1.5, 0)
      const topRoller = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.5, 16), mat.clone())
      topRoller.rotation.z = Math.PI/2; topRoller.position.set(0, 2.0, 0)
      const botRoller = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.5, 16), mat.clone())
      botRoller.rotation.z = Math.PI/2; botRoller.position.set(0, 0.1, 0)
      for (let i = 0; i < 4; i++) {
        const bucket = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.12, 0.15, 3), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.4}))
        bucket.position.set(0, i * 0.5 + 0.25, 0)
        group.add(bucket)
      }
      const hopper = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.08, 0.3, 12), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.3}))
      hopper.position.set(0, -0.15, -0.2)
      group.add(frameL, frameR, cross1, cross2, topRoller, botRoller, hopper)
      group.userData.isMeshGroup = true
      return group
    }
    case 'bin': {
      const group = new THREE.Group()
      const topBox = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.6, 0.8), mat.clone())
      topBox.position.y = 0.8
      const cone = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.1, 0.7, 24), mat.clone())
      cone.position.y = 0.15
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.2, 12), mat.clone())
      outlet.position.y = -0.25
      const lid = new THREE.Mesh(new THREE.BoxGeometry(0.85, 0.06, 0.85), mat.clone())
      lid.position.y = 1.13
      for (let i = 0; i < 4; i++) {
        const leg = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.5, 0.06), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
        leg.position.set(Math.cos(i/4*Math.PI*2+Math.PI/4)*0.4, -0.4, Math.sin(i/4*Math.PI*2+Math.PI/4)*0.4)
        group.add(leg)
      }
      group.add(topBox, cone, outlet, lid)
      group.userData.isMeshGroup = true
      return group
    }
    case 'funnel': {
      const group = new THREE.Group()
      const wide = new THREE.Mesh(new THREE.CylinderGeometry(0.6, 0.35, 0.4, 24), mat.clone())
      wide.position.y = 0.3
      const narrow = new THREE.Mesh(new THREE.CylinderGeometry(0.35, 0.06, 0.5, 24), mat.clone())
      narrow.position.y = -0.15
      const neck = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.3, 12), mat.clone())
      neck.position.y = -0.55
      group.add(wide, narrow, neck)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 仪器仪表新增 =====
    case 'thermometer': {
      const group = new THREE.Group()
      const dial = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.06, 32), new THREE.MeshStandardMaterial({color:0xffffff, roughness:0.8}))
      dial.rotation.x = Math.PI/2
      const rim = new THREE.Mesh(new THREE.TorusGeometry(0.2, 0.025, 8, 32), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      rim.rotation.x = Math.PI/2
      const face = new THREE.Mesh(new THREE.CircleGeometry(0.18, 32), new THREE.MeshStandardMaterial({color:0xffffff, side: THREE.DoubleSide}))
      face.position.z = 0.03; face.rotation.x = -Math.PI/2
      const needle = new THREE.Mesh(new THREE.BoxGeometry(0.015, 0.16, 0.01), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.3}))
      needle.position.set(0.04, 0.04, 0); needle.rotation.z = -Math.PI/6
      const stem = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.5, 8), new THREE.MeshStandardMaterial({color:0xcc0000, metalness:0.4}))
      stem.position.y = -0.25
      group.add(dial, rim, face, needle, stem)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pressuregauge': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.08, 32), new THREE.MeshStandardMaterial({color:0x333333, metalness:0.5}))
      body.rotation.x = Math.PI/2
      const face = new THREE.Mesh(new THREE.CircleGeometry(0.18, 32), new THREE.MeshStandardMaterial({color:0x222222, side: THREE.DoubleSide}))
      face.position.z = 0.04; face.rotation.x = -Math.PI/2
      const needle = new THREE.Mesh(new THREE.BoxGeometry(0.015, 0.15, 0.01), new THREE.MeshStandardMaterial({color:0xff4444, emissive:0xff4444, emissiveIntensity:0.5}))
      needle.position.set(0.05, 0.03, 0); needle.rotation.z = -Math.PI/3
      const center = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.03, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8}))
      center.rotation.x = Math.PI/2; center.position.z = 0.05
      const conn = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.05, 0.15, 8), mat.clone())
      conn.position.y = -0.1
      group.add(body, face, needle, center, conn)
      group.userData.isMeshGroup = true
      return group
    }
    case 'flowmeter': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.5, 24), mat.clone())
      const display = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.3, 0.15), new THREE.MeshStandardMaterial({color:0x222222, metalness:0.3}))
      display.position.y = 0.4
      const screen = new THREE.Mesh(new THREE.PlaneGeometry(0.2, 0.2), new THREE.MeshStandardMaterial({color:0x003366, emissive:0x003366, emissiveIntensity:0.4, side: THREE.DoubleSide}))
      screen.position.set(0, 0.4, 0.08)
      const flange1 = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.06, 16), mat.clone())
      flange1.position.z = 0.28
      const flange2 = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.06, 16), mat.clone())
      flange2.position.z = -0.28
      group.add(body, display, screen, flange1, flange2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'levelgauge': {
      const group = new THREE.Group()
      const tube = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 1.2, 16), new THREE.MeshStandardMaterial({color:0x88ccff, metalness:0.1, roughness:0.3, transparent:true, opacity:0.6}))
      const metalFrame = new THREE.Mesh(new THREE.CylinderGeometry(0.07, 0.07, 1.2, 16, 1, true), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      for (let i = 1; i <= 5; i++) {
        const mark = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.04, 0.12), new THREE.MeshStandardMaterial({color:0x333333}))
        mark.position.y = 1.2 * i / 6 - 0.2
        group.add(mark)
      }
      const topCap = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.06, 16), mat.clone())
      topCap.position.y = 0.43
      const botCap = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.06, 16), mat.clone())
      botCap.position.y = -0.43
      group.add(tube, metalFrame, topCap, botCap)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 管道支架新增 =====
    case 'pipesupport': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.08, 0.3), mat.clone())
      base.position.y = -0.35
      const column = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.5, 12), mat.clone())
      column.position.y = -0.1
      const cradle = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.14, 0.15, 16, 1, true, 0, Math.PI), mat.clone())
      cradle.rotation.x = Math.PI/2; cradle.position.set(0, 0.2, 0)
      const bolt1 = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.15, 6), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      bolt1.position.set(0.12, 0.18, 0)
      const bolt2 = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.15, 6), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      bolt2.position.set(-0.12, 0.18, 0)
      group.add(base, column, cradle, bolt1, bolt2)
      group.userData.isMeshGroup = true
      return group
    }
    case 'pipebracket': {
      const group = new THREE.Group()
      const wallPlate = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.05, 0.3), mat.clone())
      const arm = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.08, 0.4), mat.clone())
      arm.position.z = -0.2; arm.position.y = -0.02
      const brace = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.05, 0.35), mat.clone())
      brace.rotation.z = Math.PI/4; brace.position.set(0, -0.25, -0.2)
      const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.15, 12), mat.clone())
      pipe.rotation.x = Math.PI/2; pipe.position.set(0, 0, -0.35)
      group.add(wallPlate, arm, brace, pipe)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 其他设备新增 =====
    case 'flare': {
      const group = new THREE.Group()
      const tower = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.18, 2.5, 16), mat.clone())
      tower.position.y = 1.25
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.25, 0.2, 16), mat.clone())
      base.position.y = 0.1
      const tip = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.2, 8), mat.clone())
      tip.position.y = 2.6
      const flame1 = new THREE.Mesh(new THREE.ConeGeometry(0.12, 0.4, 8), new THREE.MeshStandardMaterial({color:0xff6600, emissive:0xff4400, emissiveIntensity:0.8, transparent:true, opacity:0.9}))
      flame1.position.y = 2.9
      const flame2 = new THREE.Mesh(new THREE.ConeGeometry(0.08, 0.3, 8), new THREE.MeshStandardMaterial({color:0xffcc00, emissive:0xffaa00, emissiveIntensity:1.0, transparent:true, opacity:0.8}))
      flame2.position.y = 3.05
      const flame3 = new THREE.Mesh(new THREE.ConeGeometry(0.05, 0.2, 8), new THREE.MeshStandardMaterial({color:0xffffff, emissive:0xffffff, emissiveIntensity:0.6, transparent:true, opacity:0.7}))
      flame3.position.y = 3.15
      for (let i = 0; i < 3; i++) {
        const guyWire = new THREE.Mesh(new THREE.CylinderGeometry(0.008, 0.008, 1.8, 4), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.5}))
        guyWire.position.set(Math.cos(i/3*Math.PI*2)*1.3, 1.2, Math.sin(i/3*Math.PI*2)*1.3)
        guyWire.rotation.z = Math.PI*0.25
        group.add(guyWire)
      }
      group.add(tower, base, tip, flame1, flame2, flame3)
      group.userData.isMeshGroup = true
      return group
    }
    case 'silencer': {
      const group = new THREE.Group()
      const shell = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 1.0, 24), mat.clone())
      const inlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.25, 12), mat.clone())
      inlet.position.y = 0.5
      const outlet = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.25, 12), mat.clone())
      outlet.position.y = -0.5
      const baffle1 = new THREE.Mesh(new THREE.CylinderGeometry(0.19, 0.19, 0.04, 24), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      baffle1.position.y = 0.2
      const baffle2 = new THREE.Mesh(new THREE.CylinderGeometry(0.19, 0.19, 0.04, 24), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      baffle2.position.y = -0.2
      const perfTube = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.8, 12, 1, true), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.4, roughness:0.6}))
      group.add(shell, inlet, outlet, baffle1, baffle2, perfTube)
      group.userData.isMeshGroup = true
      return group
    }

    // ===== 智慧城市 =====
    case 'building': {
      const group = new THREE.Group()
      const bmat = new THREE.MeshStandardMaterial({color:0x4da6ff, metalness:0.3, roughness:0.7})
      const wmat = new THREE.MeshStandardMaterial({color:0x88ccff, emissive:0x224466, emissiveIntensity:0.3, metalness:0.2})
      const mainBody = new THREE.Mesh(new THREE.BoxGeometry(0.8, 2.4, 0.8), bmat)
      mainBody.position.y = 1.2
      // 窗户
      for (let f = 0; f < 4; f++) {
        for (let r = 0; r < 3; r++) {
          const win = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.2, 0.02), wmat)
          win.position.y = 0.5 + r * 0.6
          if (f === 0) { win.position.z = 0.41; win.position.x = (r - 1) * 0.22 }
          if (f === 1) { win.position.z = -0.41; win.position.x = (r - 1) * 0.22 }
          if (f === 2) { win.position.x = 0.41; win.position.z = (r - 1) * 0.22 }
          if (f === 3) { win.position.x = -0.41; win.position.z = (r - 1) * 0.22 }
          group.add(win)
        }
      }
      const roof = new THREE.Mesh(new THREE.ConeGeometry(0.65, 0.4, 4), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      roof.position.y = 2.6; roof.rotation.y = Math.PI/4
      group.add(mainBody, roof)
      group.userData.isMeshGroup = true; return group
    }
    case 'apartment': {
      const group = new THREE.Group()
      const bmat = new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.2, roughness:0.8})
      const body = new THREE.Mesh(new THREE.BoxGeometry(1.2, 1.8, 0.7), bmat)
      body.position.y = 0.9
      const roof = new THREE.Mesh(new THREE.ConeGeometry(0.85, 0.35, 4), new THREE.MeshStandardMaterial({color:0x8c5523}))
      roof.position.y = 1.975; roof.rotation.y = Math.PI/4
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.4, 0.05), new THREE.MeshStandardMaterial({color:0x6b3a1f}))
      door.position.set(0, 0.2, 0.36)
      group.add(body, roof, door)
      group.userData.isMeshGroup = true; return group
    }
    case 'factorybld': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(2, 1.5, 1.2), new THREE.MeshStandardMaterial({color:0x8c8c8c, metalness:0.4, roughness:0.6}))
      body.position.y = 0.75
      const roof1 = new THREE.Mesh(new THREE.BoxGeometry(1.1, 0.15, 1.3), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.5}))
      roof1.position.set(0, 1.575, 0)
      const roof2 = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.3, 0.9), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.5}))
      roof2.position.set(0, 1.75, 0)
      group.add(body, roof1, roof2)
      group.userData.isMeshGroup = true; return group
    }
    case 'streetlight': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.08, 3, 12), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      pole.position.y = 1.5
      const arm = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.04, 0.04), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      arm.position.set(0.3, 3, 0)
      const lamp = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.15, 0.2, 16), new THREE.MeshStandardMaterial({color:0xffffcc, emissive:0xfffb4d, emissiveIntensity:0.6}))
      lamp.position.set(0.6, 2.88, 0)
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.16, 0.15, 16), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.7}))
      base.position.y = 0.075
      group.add(pole, arm, lamp, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'trafficlight': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.06, 2.2, 12), new THREE.MeshStandardMaterial({color:0x333333, metalness:0.7}))
      pole.position.y = 1.1
      const box = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.5, 0.12), new THREE.MeshStandardMaterial({color:0x222222, metalness:0.5}))
      box.position.y = 2.35
      const redLight = new THREE.Mesh(new THREE.SphereGeometry(0.05, 12, 12), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.7}))
      redLight.position.set(0, 2.48, 0.07)
      const yellowLight = new THREE.Mesh(new THREE.SphereGeometry(0.05, 12, 12), new THREE.MeshStandardMaterial({color:0xffaa00, emissive:0x884400, emissiveIntensity:0.3}))
      yellowLight.position.set(0, 2.35, 0.07)
      const greenLight = new THREE.Mesh(new THREE.SphereGeometry(0.05, 12, 12), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x004400, emissiveIntensity:0.3}))
      greenLight.position.set(0, 2.22, 0.07)
      group.add(pole, box, redLight, yellowLight, greenLight)
      group.userData.isMeshGroup = true; return group
    }
    case 'bridge': {
      const group = new THREE.Group()
      const deck = new THREE.Mesh(new THREE.BoxGeometry(3, 0.15, 0.8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      const pillar1 = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.18, 1.5, 16), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.4}))
      pillar1.position.set(-1, -0.75, 0)
      const pillar2 = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.18, 1.5, 16), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.4}))
      pillar2.position.set(1, -0.75, 0)
      const cable1 = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1.8, 8), new THREE.MeshStandardMaterial({color:0x333333}))
      cable1.position.set(-0.5, 0.9, 0); cable1.rotation.z = Math.PI/4
      const cable2 = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1.8, 8), new THREE.MeshStandardMaterial({color:0x333333}))
      cable2.position.set(0.5, 0.9, 0); cable2.rotation.z = -Math.PI/4
      group.add(deck, pillar1, pillar2, cable1, cable2)
      group.userData.isMeshGroup = true; return group
    }
    case 'tunnel': {
      const group = new THREE.Group()
      const arch = new THREE.Mesh(new THREE.TorusGeometry(0.6, 0.08, 8, 16, Math.PI), new THREE.MeshStandardMaterial({color:0x777777, metalness:0.5}))
      arch.rotation.z = Math.PI
      const wall1 = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.6, 1.5), new THREE.MeshStandardMaterial({color:0x888888}))
      wall1.position.set(-0.65, 0, 0)
      const wall2 = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.6, 1.5), new THREE.MeshStandardMaterial({color:0x888888}))
      wall2.position.set(0.65, 0, 0)
      group.add(arch, wall1, wall2)
      group.userData.isMeshGroup = true; return group
    }
    case 'road': {
      const group = new THREE.Group()
      const surface = new THREE.Mesh(new THREE.BoxGeometry(3, 0.05, 1), new THREE.MeshStandardMaterial({color:0x444444, roughness:0.9}))
      surface.position.y = 0.025
      const line1 = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.02, 0.04), new THREE.MeshStandardMaterial({color:0xffffff}))
      line1.position.set(-0.6, 0.04, 0)
      const line2 = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.02, 0.04), new THREE.MeshStandardMaterial({color:0xffffff}))
      line2.position.set(0.6, 0.04, 0)
      group.add(surface, line1, line2)
      group.userData.isMeshGroup = true; return group
    }
    case 'parkingspot': {
      const group = new THREE.Group()
      const ground = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.03, 2.5), new THREE.MeshStandardMaterial({color:0x555555, roughness:0.9}))
      const lineLeft = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.02, 2.5), new THREE.MeshStandardMaterial({color:0xffffff}))
      lineLeft.position.x = -0.6
      const lineRight = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.02, 2.5), new THREE.MeshStandardMaterial({color:0xffffff}))
      lineRight.position.x = 0.6
      const lineFront = new THREE.Mesh(new THREE.BoxGeometry(1.26, 0.02, 0.03), new THREE.MeshStandardMaterial({color:0xffffff}))
      lineFront.position.z = 1.25
      group.add(ground, lineLeft, lineRight, lineFront)
      group.userData.isMeshGroup = true; return group
    }
    case 'busstop': {
      const group = new THREE.Group()
      const floor = new THREE.Mesh(new THREE.BoxGeometry(1.5, 0.05, 0.6), new THREE.MeshStandardMaterial({color:0x888888}))
      const back = new THREE.Mesh(new THREE.BoxGeometry(1.5, 1.2, 0.05), new THREE.MeshStandardMaterial({color:0x4488cc, transparent:true, opacity:0.6}))
      back.position.set(0, 0.65, -0.3)
      const leftPillar = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 1.2, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      leftPillar.position.set(-0.7, 0.6, -0.3)
      const rightPillar = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 1.2, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      rightPillar.position.set(0.7, 0.6, -0.3)
      const roof = new THREE.Mesh(new THREE.BoxGeometry(1.6, 0.06, 0.7), new THREE.MeshStandardMaterial({color:0x555555}))
      roof.position.y = 1.25
      const bench = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.08, 0.2), new THREE.MeshStandardMaterial({color:0xa67c52}))
      bench.position.set(0, 0.27, -0.15)
      group.add(floor, back, leftPillar, rightPillar, roof, bench)
      group.userData.isMeshGroup = true; return group
    }
    case 'bench': {
      const group = new THREE.Group()
      const seat = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.06, 0.25), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.7}))
      seat.position.y = 0.35
      for (let i = -1; i <= 1; i++) {
        const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.35, 8), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
        leg.position.set(i * 0.45, 0.175, 0)
        group.add(leg)
      }
      group.add(seat)
      group.userData.isMeshGroup = true; return group
    }
    case 'trashbin': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.12, 0.6, 16), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.2, roughness:0.6}))
      body.position.y = 0.3
      const lid = new THREE.Mesh(new THREE.CylinderGeometry(0.16, 0.14, 0.08, 16), new THREE.MeshStandardMaterial({color:0x339955, metalness:0.2}))
      lid.position.y = 0.64
      group.add(body, lid)
      group.userData.isMeshGroup = true; return group
    }
    case 'citytree': {
      const group = new THREE.Group()
      const trunk = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.09, 0.8, 8), new THREE.MeshStandardMaterial({color:0x8c5523, roughness:0.9}))
      trunk.position.y = 0.4
      const foliage1 = new THREE.Mesh(new THREE.SphereGeometry(0.35, 12, 12), new THREE.MeshStandardMaterial({color:0x44aa44}))
      foliage1.position.y = 0.85
      const foliage2 = new THREE.Mesh(new THREE.SphereGeometry(0.28, 12, 12), new THREE.MeshStandardMaterial({color:0x55bb55}))
      foliage2.position.set(0.15, 1.05, 0.1)
      const foliage3 = new THREE.Mesh(new THREE.SphereGeometry(0.26, 12, 12), new THREE.MeshStandardMaterial({color:0x339933}))
      foliage3.position.set(-0.12, 1.08, -0.08)
      group.add(trunk, foliage1, foliage2, foliage3)
      group.userData.isMeshGroup = true; return group
    }
    case 'fountain': {
      const group = new THREE.Group()
      const basin = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.55, 0.2, 24), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.5}))
      basin.position.y = 0.1
      const water = new THREE.Mesh(new THREE.CylinderGeometry(0.45, 0.45, 0.05, 24), new THREE.MeshStandardMaterial({color:0x4dfff0, transparent:true, opacity:0.6, emissive:0x2266aa, emissiveIntensity:0.3}))
      water.position.y = 0.18
      const pillar = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.5, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.5}))
      pillar.position.y = 0.35
      const spray = new THREE.Mesh(new THREE.ConeGeometry(0.3, 0.4, 12, 1, true), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.3, side:THREE.DoubleSide}))
      spray.position.y = 0.7
      group.add(basin, water, pillar, spray)
      group.userData.isMeshGroup = true; return group
    }
    case 'surveilcam': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.05, 1.8, 12), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      pole.position.y = 0.9
      const arm = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.04, 0.04), new THREE.MeshStandardMaterial({color:0x444444}))
      arm.position.set(0.15, 1.82, 0); arm.rotation.z = Math.PI/6
      const camBody = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.08, 0.08), new THREE.MeshStandardMaterial({color:0x222222}))
      camBody.position.set(0.3, 1.88, 0)
      const lens = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.06, 12), new THREE.MeshStandardMaterial({color:0x111111}))
      lens.position.set(0.3, 1.88, 0.06); lens.rotation.x = Math.PI/2
      group.add(pole, arm, camBody, lens)
      group.userData.isMeshGroup = true; return group
    }
    case 'billboard': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.1, 2.5, 12), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.7}))
      pole.position.y = 1.25
      const board = new THREE.Mesh(new THREE.BoxGeometry(1.5, 0.7, 0.05), new THREE.MeshStandardMaterial({color:0xfffb4d, emissive:0x333300, emissiveIntensity:0.5}))
      board.position.y = 2.7
      const frame = new THREE.Mesh(new THREE.BoxGeometry(1.6, 0.8, 0.06), new THREE.MeshStandardMaterial({color:0x333333, metalness:0.3, wireframe:true}))
      frame.position.y = 2.7
      group.add(pole, board, frame)
      group.userData.isMeshGroup = true; return group
    }
    case 'solarpanel': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(1, 0.04, 0.6), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.7}))
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.02, 0.5), new THREE.MeshStandardMaterial({color:0x114488, metalness:0.3, emissive:0x001122, emissiveIntensity:0.2}))
      panel.position.y = 0.04
      const stand = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.05, 0.5, 8), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.7}))
      stand.position.y = -0.25
      group.add(frame, panel, stand)
      group.userData.isMeshGroup = true; return group
    }
    case 'windturbine': {
      const group = new THREE.Group()
      const tower = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.1, 3, 16), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.6}))
      tower.position.y = 1.5
      const nacelle = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.2, 0.2), new THREE.MeshStandardMaterial({color:0xdddddd, metalness:0.5}))
      nacelle.position.y = 3.1
      const hub = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.1, 12), new THREE.MeshStandardMaterial({color:0xffffff}))
      hub.position.set(0.2, 3.1, 0); hub.rotation.z = Math.PI/2
      for (let i = 0; i < 3; i++) {
        const blade = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.03, 0.08), new THREE.MeshStandardMaterial({color:0xeeeeee, metalness:0.5}))
        blade.position.set(0.45, 3.1, 0)
        blade.rotation.z = (i * 2 * Math.PI / 3)
        blade.rotation.x = Math.PI/2
        group.add(blade)
      }
      group.add(tower, nacelle, hub)
      group.userData.isMeshGroup = true; return group
    }
    case 'telecomtower': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.1, 0.5), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      for (let i = 0; i < 3; i++) {
        const seg = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 1, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
        seg.position.y = 0.5 + i * 1
        group.add(seg)
        const cross = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.05, 0.3), new THREE.MeshStandardMaterial({color:0x777777, metalness:0.7}))
        cross.position.y = 1 + i * 1
        group.add(cross)
      }
      const antenna = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.02, 0.5, 8), new THREE.MeshStandardMaterial({color:0xff0000}))
      antenna.position.y = 3.25
      group.add(base, antenna)
      group.userData.isMeshGroup = true; return group
    }
    case 'manhole': {
      const group = new THREE.Group()
      const cover = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.25, 0.04, 24), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8, roughness:0.4}))
      const ring = new THREE.Mesh(new THREE.TorusGeometry(0.24, 0.02, 8, 24), new THREE.MeshStandardMaterial({color:0x777777, metalness:0.8}))
      ring.rotation.x = Math.PI/2; ring.position.y = 0.03
      group.add(cover, ring)
      group.userData.isMeshGroup = true; return group
    }
    case 'hydrant': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.12, 0.6, 16), new THREE.MeshStandardMaterial({color:0xff0000, metalness:0.3}))
      body.position.y = 0.3
      const top = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.1, 0.08, 16), new THREE.MeshStandardMaterial({color:0xcc0000, metalness:0.3}))
      top.position.y = 0.64
      const cap = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.04, 16), new THREE.MeshStandardMaterial({color:0x990000}))
      cap.position.y = 0.7
      const nozzle1 = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.15, 8), new THREE.MeshStandardMaterial({color:0xcc0000}))
      nozzle1.position.set(0.08, 0.55, 0); nozzle1.rotation.z = Math.PI/2
      const nozzle2 = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.15, 8), new THREE.MeshStandardMaterial({color:0xcc0000}))
      nozzle2.position.set(-0.08, 0.55, 0); nozzle2.rotation.z = Math.PI/2
      group.add(body, top, cap, nozzle1, nozzle2)
      group.userData.isMeshGroup = true; return group
    }
    case 'bollard': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.1, 0.6, 16), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.3, roughness:0.5}))
      body.position.y = 0.3
      const top = new THREE.Mesh(new THREE.SphereGeometry(0.09, 12, 12, 0, Math.PI*2, 0, Math.PI/2), new THREE.MeshStandardMaterial({color:0xcccc00, metalness:0.3}))
      top.position.y = 0.6
      const stripe = new THREE.Mesh(new THREE.TorusGeometry(0.09, 0.02, 8, 16), new THREE.MeshStandardMaterial({color:0xff0000}))
      stripe.position.y = 0.35; stripe.rotation.x = Math.PI/2
      group.add(body, top, stripe)
      group.userData.isMeshGroup = true; return group
    }
    case 'flagpole': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.05, 2.5, 12), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      pole.position.y = 1.25
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.14, 0.15, 16), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      const flag = new THREE.Mesh(new THREE.PlaneGeometry(0.4, 0.25), new THREE.MeshStandardMaterial({color:0xff0000, side:THREE.DoubleSide}))
      flag.position.set(0.25, 2.5, 0)
      group.add(pole, base, flag)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 工业流水线 =====
    case 'robotarm': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.25, 0.15, 16), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.7}))
      const turntable = new THREE.Mesh(new THREE.CylinderGeometry(0.18, 0.18, 0.1, 24), new THREE.MeshStandardMaterial({color:0xdd8822, metalness:0.7}))
      turntable.position.y = 0.13
      const lowerArm = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.5, 0.12), new THREE.MeshStandardMaterial({color:0xff8800, metalness:0.6}))
      lowerArm.position.y = 0.5; lowerArm.position.z = -0.06
      const elbow = new THREE.Mesh(new THREE.SphereGeometry(0.08, 12, 12), new THREE.MeshStandardMaterial({color:0xdd6600, metalness:0.7}))
      elbow.position.set(0, 0.78, -0.1)
      const upperArm = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.35, 0.08), new THREE.MeshStandardMaterial({color:0xff8800, metalness:0.6}))
      upperArm.position.set(0, 0.95, -0.2); upperArm.rotation.x = Math.PI/3
      const wrist = new THREE.Mesh(new THREE.SphereGeometry(0.05, 8, 8), new THREE.MeshStandardMaterial({color:0xdd6600, metalness:0.7}))
      wrist.position.set(0, 1.2, -0.35)
      const gripper1 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.1, 0.02), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.8}))
      gripper1.position.set(0.04, 1.26, -0.35)
      const gripper2 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.1, 0.02), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.8}))
      gripper2.position.set(-0.04, 1.26, -0.35)
      group.add(base, turntable, lowerArm, elbow, upperArm, wrist, gripper1, gripper2)
      group.userData.isMeshGroup = true; return group
    }
    case 'conveyorbelt': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(1.5, 0.08, 0.3), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      frame.position.y = 0.2
      const belt = new THREE.Mesh(new THREE.BoxGeometry(1.4, 0.02, 0.22), new THREE.MeshStandardMaterial({color:0x333333, roughness:0.9}))
      belt.position.y = 0.25
      const roller1 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.24, 16), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      roller1.position.set(-0.7, 0.25, 0); roller1.rotation.z = Math.PI/2
      const roller2 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.24, 16), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      roller2.position.set(0.7, 0.25, 0); roller2.rotation.z = Math.PI/2
      for (let i = 0; i < 4; i++) {
        const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.2, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
        leg.position.set(-0.5 + i * 0.35, 0.1, i%2 ? -0.12 : 0.12)
        group.add(leg)
      }
      group.add(frame, belt, roller1, roller2)
      group.userData.isMeshGroup = true; return group
    }
    case 'assemblystation': {
      const group = new THREE.Group()
      const table = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.06, 0.5), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.5}))
      table.position.y = 0.5
      for (let i = -1; i <= 1; i+=2) {
        const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.5, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
        leg.position.set(i * 0.3, 0.25, 0.18)
        group.add(leg)
        const leg2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.5, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
        leg2.position.set(i * 0.3, 0.25, -0.18)
        group.add(leg2)
      }
      const toolbox = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.25), new THREE.MeshStandardMaterial({color:0xff4444, metalness:0.5}))
      toolbox.position.set(0.3, 0.6, -0.1)
      const light = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.03, 0.12), new THREE.MeshStandardMaterial({color:0xffffcc, emissive:0xfffb4d, emissiveIntensity:0.5}))
      light.position.set(0, 0.7, 0.1)
      group.add(table, toolbox, light)
      group.userData.isMeshGroup = true; return group
    }
    case 'weldingrobot': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.1, 0.25), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.7}))
      const arm1 = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.06, 0.4, 12), new THREE.MeshStandardMaterial({color:0xffaa00, metalness:0.6}))
      arm1.position.y = 0.25
      const arm2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.3, 12), new THREE.MeshStandardMaterial({color:0xffaa00, metalness:0.6}))
      arm2.position.set(0.15, 0.45, 0); arm2.rotation.z = -Math.PI/4
      const torch = new THREE.Mesh(new THREE.ConeGeometry(0.02, 0.1, 8), new THREE.MeshStandardMaterial({color:0xff8800, emissive:0xff4400, emissiveIntensity:0.8}))
      torch.position.set(0.3, 0.6, 0)
      const spark = new THREE.Mesh(new THREE.SphereGeometry(0.03, 8, 8), new THREE.MeshStandardMaterial({color:0xffffff, emissive:0xffff00, emissiveIntensity:1}))
      spark.position.set(0.32, 0.55, 0)
      group.add(base, arm1, arm2, torch, spark)
      group.userData.isMeshGroup = true; return group
    }
    case 'inspectioncam': {
      const group = new THREE.Group()
      const mount = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.05, 0.6, 12), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      mount.position.y = 0.3
      const camHead = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.1, 0.15), new THREE.MeshStandardMaterial({color:0x222222}))
      camHead.position.y = 0.65
      const lens = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.05, 0.08, 12), new THREE.MeshStandardMaterial({color:0x111111}))
      lens.position.set(0, 0.65, 0.1); lens.rotation.x = Math.PI/2
      const ring = new THREE.Mesh(new THREE.TorusGeometry(0.05, 0.01, 8, 16), new THREE.MeshStandardMaterial({color:0xff4444, emissive:0xff0000, emissiveIntensity:0.5}))
      ring.position.set(0, 0.65, 0.12)
      group.add(mount, camHead, lens, ring)
      group.userData.isMeshGroup = true; return group
    }
    case 'packmachine': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.7, 0.6, 0.5), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.4}))
      body.position.y = 0.3
      const inlet = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.1, 0.1), new THREE.MeshStandardMaterial({color:0x888888}))
      inlet.position.set(0, 0.5, -0.26)
      const outlet = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.1, 0.1), new THREE.MeshStandardMaterial({color:0x888888}))
      outlet.position.set(0, 0.1, -0.26)
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.02), new THREE.MeshStandardMaterial({color:0x222222}))
      panel.position.set(0.3, 0.4, 0.26)
      group.add(body, inlet, outlet, panel)
      group.userData.isMeshGroup = true; return group
    }
    case 'palletizer': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.05, 0.5), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.6}))
      const pillar1 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pillar1.position.set(-0.2, 0.4, -0.2)
      const pillar2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pillar2.position.set(0.2, 0.4, -0.2)
      const pillar3 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pillar3.position.set(-0.2, 0.4, 0.2)
      const pillar4 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pillar4.position.set(0.2, 0.4, 0.2)
      const head = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.05, 0.3), new THREE.MeshStandardMaterial({color:0xff8800, metalness:0.6}))
      head.position.y = 0.83
      const gripper = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.04, 0.15), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8}))
      gripper.position.y = 0.78
      group.add(frame, pillar1, pillar2, pillar3, pillar4, head, gripper)
      group.userData.isMeshGroup = true; return group
    }
    case 'agv': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.12, 0.35), new THREE.MeshStandardMaterial({color:0x4da6ff, metalness:0.5}))
      body.position.y = 0.12
      for (let i = -1; i <= 1; i+=2) {
        const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.04, 16), new THREE.MeshStandardMaterial({color:0x222222}))
        wheel.position.set(i * 0.2, 0.06, 0.2); wheel.rotation.z = Math.PI/2
        group.add(wheel)
        const wheel2 = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.04, 16), new THREE.MeshStandardMaterial({color:0x222222}))
        wheel2.position.set(i * 0.2, 0.06, -0.2); wheel2.rotation.z = Math.PI/2
        group.add(wheel2)
      }
      const light = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.04, 0.02), new THREE.MeshStandardMaterial({color:0xfffb4d, emissive:0xfffb4d, emissiveIntensity:0.6}))
      light.position.set(0, 0.12, 0.19)
      group.add(body, light)
      group.userData.isMeshGroup = true; return group
    }
    case 'feeder': {
      const group = new THREE.Group()
      const hopper = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.2, 0.25, 16, 1, true), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5, side:THREE.DoubleSide}))
      hopper.position.y = 0.5
      const chute = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.12, 0.3, 12, 1, true), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.5, side:THREE.DoubleSide}))
      chute.position.y = 0.2
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.22, 0.24, 0.08, 16), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      group.add(hopper, chute, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'stampingpress': {
      const group = new THREE.Group()
      const frame1 = new THREE.Mesh(new THREE.BoxGeometry(0.07, 0.8, 0.07), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.6}))
      frame1.position.x = -0.2
      const frame2 = new THREE.Mesh(new THREE.BoxGeometry(0.07, 0.8, 0.07), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.6}))
      frame2.position.x = 0.2
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.1, 0.3), new THREE.MeshStandardMaterial({color:0xcc0000, metalness:0.6}))
      top.position.y = 0.45
      const ram = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.3, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.8}))
      ram.position.y = 0.3
      const bed = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.06, 0.3), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      bed.position.y = 0.03
      group.add(frame1, frame2, top, ram, bed)
      group.userData.isMeshGroup = true; return group
    }
    case 'cncmachine': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.6, 0.5), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.5}))
      body.position.y = 0.3
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.4, 0.02), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.5}))
      door.position.set(0, 0.3, 0.26)
      const spindle = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.2, 12), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.8}))
      spindle.position.set(0, 0.55, -0.2)
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.2, 0.02), new THREE.MeshStandardMaterial({color:0x222222}))
      panel.position.set(-0.35, 0.4, 0.26)
      group.add(body, door, spindle, panel)
      group.userData.isMeshGroup = true; return group
    }
    case 'injectionmold': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.5, 0.4), new THREE.MeshStandardMaterial({color:0xc17bff, metalness:0.4}))
      body.position.y = 0.25
      const barrel = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.08, 0.4, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      barrel.position.set(0.35, 0.4, 0); barrel.rotation.z = Math.PI/2
      const hopper = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.12, 0.2, 12, 1, true), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.5, side:THREE.DoubleSide}))
      hopper.position.set(0.5, 0.5, 0)
      group.add(body, barrel, hopper)
      group.userData.isMeshGroup = true; return group
    }
    case 'workbench': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.BoxGeometry(1, 0.05, 0.5), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.7}))
      top.position.y = 0.6
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.6, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
          leg.position.set(i * 0.4, 0.3, j * 0.18)
          group.add(leg)
        }
      }
      const vise = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.1, 0.12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      vise.position.set(0.3, 0.68, 0)
      group.add(top, vise)
      group.userData.isMeshGroup = true; return group
    }
    case 'turntable': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.3, 0.1, 24), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      const table = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.25, 0.06, 32), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.7}))
      table.position.y = 0.08
      const ring = new THREE.Mesh(new THREE.TorusGeometry(0.24, 0.015, 8, 32), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      ring.position.y = 0.05; ring.rotation.x = Math.PI/2
      group.add(base, table, ring)
      group.userData.isMeshGroup = true; return group
    }
    case 'lifter': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.06, 0.4), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.6}))
      const pole1 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pole1.position.set(-0.15, 0.4, -0.15)
      const pole2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pole2.position.set(0.15, 0.4, -0.15)
      const pole3 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pole3.position.set(-0.15, 0.4, 0.15)
      const pole4 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pole4.position.set(0.15, 0.4, 0.15)
      const platform = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.04, 0.35), new THREE.MeshStandardMaterial({color:0xff8800, metalness:0.6}))
      platform.position.y = 0.42
      group.add(base, pole1, pole2, pole3, pole4, platform)
      group.userData.isMeshGroup = true; return group
    }
    case 'scanner3d': {
      const group = new THREE.Group()
      const stand = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.06, 0.8, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      stand.position.y = 0.4
      const head = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.1, 0.15), new THREE.MeshStandardMaterial({color:0x4dfff0, metalness:0.3, emissive:0x003333, emissiveIntensity:0.5}))
      head.position.y = 0.85
      const lens = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.06, 12), new THREE.MeshStandardMaterial({color:0x111111}))
      lens.position.set(0, 0.85, -0.1); lens.rotation.x = Math.PI/2
      group.add(stand, head, lens)
      group.userData.isMeshGroup = true; return group
    }
    case 'labelprinter': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.25, 0.2), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.3}))
      body.position.y = 0.125
      const roll = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.12, 12), new THREE.MeshStandardMaterial({color:0xffffff}))
      roll.position.set(0, 0.2, 0.12); roll.rotation.x = Math.PI/2
      const slot = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.02, 0.05), new THREE.MeshStandardMaterial({color:0x222222}))
      slot.position.set(0, 0.1, -0.13)
      group.add(body, roll, slot)
      group.userData.isMeshGroup = true; return group
    }
    case 'qualitygate': {
      const group = new THREE.Group()
      const pillar1 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 1, 0.08), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.5}))
      pillar1.position.set(-0.3, 0.5, 0)
      const pillar2 = new THREE.Mesh(new THREE.BoxGeometry(0.08, 1, 0.08), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.5}))
      pillar2.position.set(0.3, 0.5, 0)
      const topBar = new THREE.Mesh(new THREE.BoxGeometry(0.7, 0.06, 0.1), new THREE.MeshStandardMaterial({color:0x33cc66, metalness:0.5}))
      topBar.position.y = 1.03
      const light = new THREE.Mesh(new THREE.SphereGeometry(0.05, 8, 8), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.8}))
      light.position.set(0, 0.95, 0.05)
      group.add(pillar1, pillar2, topBar, light)
      group.userData.isMeshGroup = true; return group
    }
    case 'buffertable': {
      const group = new THREE.Group()
      const platform = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.04, 0.4), new THREE.MeshStandardMaterial({color:0x4dfff0, metalness:0.3}))
      platform.position.y = 0.35
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.03, 0.35, 8), new THREE.MeshStandardMaterial({color:0x777777, metalness:0.7}))
          leg.position.set(i * 0.23, 0.175, j * 0.14)
          group.add(leg)
        }
      }
      group.add(platform)
      group.userData.isMeshGroup = true; return group
    }
    case 'safetyfence': {
      const group = new THREE.Group()
      const rail1 = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1, 8), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.7}))
      rail1.position.set(0, 0.6, 0)
      const rail2 = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1, 8), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.7}))
      rail2.position.set(0, 0.3, 0)
      for (let i = -3; i <= 3; i++) {
        const post = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.7, 8), new THREE.MeshStandardMaterial({color:0xffcc00, metalness:0.7}))
        post.position.set(i * 0.15, 0.45, 0)
        group.add(post)
      }
      group.add(rail1, rail2)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 数字孪生 =====
    case 'serverrack': {
      const group = new THREE.Group()
      const cabinet = new THREE.Mesh(new THREE.BoxGeometry(0.4, 1, 0.5), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      cabinet.position.y = 0.5
      for (let i = 0; i < 5; i++) {
        const server = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.12, 0.45), new THREE.MeshStandardMaterial({color:0x222222, metalness:0.5}))
        server.position.y = 0.08 + i * 0.18
        const blinker = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.02, 0.02), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.7}))
        blinker.position.set(0.16, 0.1 + i * 0.18, 0.23)
        group.add(server, blinker)
      }
      group.add(cabinet)
      group.userData.isMeshGroup = true; return group
    }
    case 'datacenter': {
      const group = new THREE.Group()
      const room = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.8, 0.8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.3, transparent:true, opacity:0.5}))
      room.position.y = 0.4
      const floor = new THREE.Mesh(new THREE.BoxGeometry(1.1, 0.03, 0.7), new THREE.MeshStandardMaterial({color:0x666666}))
      const rack1 = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.5, 0.3), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      rack1.position.set(-0.35, 0.28, 0)
      const rack2 = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.5, 0.3), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      rack2.position.set(0.35, 0.28, 0)
      group.add(room, floor, rack1, rack2)
      group.userData.isMeshGroup = true; return group
    }
    case 'iotdevice': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.08, 0.1), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.3, emissive:0x003300, emissiveIntensity:0.3}))
      const antenna = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.02, 0.12, 8), new THREE.MeshStandardMaterial({color:0x888888}))
      antenna.position.y = 0.1
      const led = new THREE.Mesh(new THREE.SphereGeometry(0.015, 6, 6), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:1}))
      led.position.set(0.06, 0, 0.06)
      group.add(body, antenna, led)
      group.userData.isMeshGroup = true; return group
    }
    case 'gateway': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.12), new THREE.MeshStandardMaterial({color:0xc17bff, metalness:0.4}))
      for (let i = 0; i < 4; i++) {
        const port = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.03, 0.02), new THREE.MeshStandardMaterial({color:0xffaa00}))
        port.position.set(-0.07 + i * 0.05, 0, 0.07)
        group.add(port)
      }
      const ant = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.02, 0.08, 8), new THREE.MeshStandardMaterial({color:0x444444}))
      ant.position.set(0, 0.12, 0)
      group.add(body, ant)
      group.userData.isMeshGroup = true; return group
    }
    case 'controlpanel': {
      const group = new THREE.Group()
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.3, 0.05), new THREE.MeshStandardMaterial({color:0x333333, metalness:0.5}))
      panel.position.y = 0.15
      const btn1 = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.04, 0.01), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.5}))
      btn1.position.set(-0.1, 0.22, 0.03)
      const btn2 = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.04, 0.01), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.5}))
      btn2.position.set(0, 0.22, 0.03)
      const btn3 = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.04, 0.01), new THREE.MeshStandardMaterial({color:0xfffb4d, emissive:0x884400, emissiveIntensity:0.3}))
      btn3.position.set(0.1, 0.22, 0.03)
      const knob = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.02, 16), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      knob.position.set(0, 0.1, 0.03); knob.rotation.x = Math.PI/2
      group.add(panel, btn1, btn2, btn3, knob)
      group.userData.isMeshGroup = true; return group
    }
    case 'dashboardscreen': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.1, 0.3, 12), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      base.position.y = 0.15
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.4, 0.04), new THREE.MeshStandardMaterial({color:0x4da6ff, emissive:0x001133, emissiveIntensity:0.4}))
      screen.position.y = 0.55
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.66, 0.46, 0.06), new THREE.MeshStandardMaterial({color:0x222222, metalness:0.5, wireframe:true}))
      frame.position.y = 0.55
      group.add(base, screen, frame)
      group.userData.isMeshGroup = true; return group
    }
    case 'antenna': {
      const group = new THREE.Group()
      const mast = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.06, 1.2, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.7}))
      mast.position.y = 0.6
      const dish = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.25, 0.08, 24), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.5}))
      dish.position.y = 1.24
      const feed = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.15, 8), new THREE.MeshStandardMaterial({color:0x666666}))
      feed.position.y = 1.15
      group.add(mast, dish, feed)
      group.userData.isMeshGroup = true; return group
    }
    case 'radardish': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.14, 0.15, 16), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.06, 0.6, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pole.position.y = 0.38
      const dish = new THREE.Mesh(new THREE.SphereGeometry(0.25, 24, 16, 0, Math.PI*2, 0, Math.PI/3), new THREE.MeshStandardMaterial({color:0x4dfff0, metalness:0.5, side:THREE.DoubleSide}))
      dish.position.set(0, 0.6, -0.05); dish.rotation.x = Math.PI/2
      group.add(base, pole, dish)
      group.userData.isMeshGroup = true; return group
    }
    case 'router': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.06, 0.15), new THREE.MeshStandardMaterial({color:0x4da6ff, metalness:0.3}))
      for (let i = 0; i < 3; i++) {
        const ant = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.02, 0.12, 8), new THREE.MeshStandardMaterial({color:0x888888}))
        ant.position.set(-0.06 + i * 0.06, 0.1, 0)
        group.add(ant)
      }
      const lights = []
      for (let i = 0; i < 4; i++) {
        const led = new THREE.Mesh(new THREE.BoxGeometry(0.015, 0.01, 0.01), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.6}))
        led.position.set(-0.06 + i * 0.04, 0.035, 0.08)
        group.add(led)
      }
      group.add(body)
      group.userData.isMeshGroup = true; return group
    }
    case 'switchgear': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.08, 0.2), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      for (let i = 0; i < 8; i++) {
        const port = new THREE.Mesh(new THREE.BoxGeometry(0.025, 0.04, 0.01), new THREE.MeshStandardMaterial({color:0x444444}))
        port.position.set(-0.14 + i * 0.04, 0, 0.11)
        group.add(port)
      }
      const led = new THREE.Mesh(new THREE.SphereGeometry(0.01, 4, 4), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.8}))
      led.position.set(0.14, 0, 0.05)
      group.add(body, led)
      group.userData.isMeshGroup = true; return group
    }
    case 'ups': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.4, 0.25), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.4}))
      body.position.y = 0.2
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.12, 0.02), new THREE.MeshStandardMaterial({color:0x222222}))
      panel.position.set(0, 0.3, 0.13)
      const indicator = new THREE.Mesh(new THREE.SphereGeometry(0.02, 6, 6), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.8}))
      indicator.position.set(0, 0.35, 0.14)
      group.add(body, panel, indicator)
      group.userData.isMeshGroup = true; return group
    }
    case 'patchpanel': {
      const group = new THREE.Group()
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.06, 0.08), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      for (let i = 0; i < 12; i++) {
        const port = new THREE.Mesh(new THREE.BoxGeometry(0.02, 0.04, 0.02), new THREE.MeshStandardMaterial({color:0x333333}))
        port.position.set(-0.14 + (i%6) * 0.055, 0, 0.05)
        group.add(port)
      }
      group.add(panel)
      group.userData.isMeshGroup = true; return group
    }
    case 'fiberbox': {
      const group = new THREE.Group()
      const box = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.1, 0.1), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.3}))
      const fiber1 = new THREE.Mesh(new THREE.CylinderGeometry(0.005, 0.005, 0.15, 6), new THREE.MeshStandardMaterial({color:0xffaa00}))
      fiber1.position.set(0, 0, 0.12)
      const fiber2 = new THREE.Mesh(new THREE.CylinderGeometry(0.005, 0.005, 0.13, 6), new THREE.MeshStandardMaterial({color:0xffaa00}))
      fiber2.position.set(-0.03, 0, 0.14)
      group.add(box, fiber1, fiber2)
      group.userData.isMeshGroup = true; return group
    }
    case 'rtu': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.12), new THREE.MeshStandardMaterial({color:0xc17bff, metalness:0.4}))
      const term1 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.04, 0.02), new THREE.MeshStandardMaterial({color:0x888888}))
      term1.position.set(-0.06, 0.08, 0.07)
      const term2 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.04, 0.02), new THREE.MeshStandardMaterial({color:0x888888}))
      term2.position.set(0.06, 0.08, 0.07)
      const status = new THREE.Mesh(new THREE.SphereGeometry(0.015, 6, 6), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.7}))
      status.position.set(-0.07, 0, 0.07)
      group.add(body, term1, term2, status)
      group.userData.isMeshGroup = true; return group
    }
    case 'plc': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.2, 0.12), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.4}))
      body.position.y = 0.1
      for (let i = 0; i < 6; i++) {
        const out = new THREE.Mesh(new THREE.BoxGeometry(0.015, 0.025, 0.01), new THREE.MeshStandardMaterial({color:0x444444}))
        out.position.set(-0.1 + i * 0.04, 0.18, 0.07)
        group.add(out)
      }
      const runLed = new THREE.Mesh(new THREE.SphereGeometry(0.015, 6, 6), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.8}))
      runLed.position.set(0.1, 0.18, 0.07)
      group.add(body, runLed)
      group.userData.isMeshGroup = true; return group
    }
    case 'hmipanel': {
      const group = new THREE.Group()
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.18, 0.03), new THREE.MeshStandardMaterial({color:0x4da6ff, emissive:0x001133, emissiveIntensity:0.4}))
      screen.position.y = 0.2
      const bezel = new THREE.Mesh(new THREE.BoxGeometry(0.26, 0.19, 0.04), new THREE.MeshStandardMaterial({color:0x333333, metalness:0.5, wireframe:true}))
      bezel.position.y = 0.2
      const mount = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.25, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      mount.position.y = 0.075
      group.add(screen, bezel, mount)
      group.userData.isMeshGroup = true; return group
    }
    case 'beacon': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.04, 16), new THREE.MeshStandardMaterial({color:0xcccccc}))
      const top = new THREE.Mesh(new THREE.SphereGeometry(0.06, 12, 8, 0, Math.PI*2, 0, Math.PI/2), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.8}))
      top.position.y = 0.04
      group.add(body, top)
      group.userData.isMeshGroup = true; return group
    }
    case 'sensorarray': {
      const group = new THREE.Group()
      const bar = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.04, 0.06), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.5}))
      for (let i = 0; i < 4; i++) {
        const sensor = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.02, 0.06, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.5}))
        sensor.position.set(-0.1 + i * 0.07, 0.05, 0)
        group.add(sensor)
      }
      group.add(bar)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 能源电力 =====
    case 'transformer': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.5, 0.3), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.4}))
      body.position.y = 0.25
      const bushing1 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.15, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa}))
      bushing1.position.set(-0.1, 0.58, 0)
      const bushing2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.15, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa}))
      bushing2.position.set(0.1, 0.58, 0)
      const fins = []
      for (let i = 0; i < 4; i++) {
        const fin = new THREE.Mesh(new THREE.BoxGeometry(0.38, 0.35, 0.02), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
        fin.position.set(0, 0.25, -0.12 + i * 0.08)
        group.add(fin)
      }
      group.add(body, bushing1, bushing2)
      group.userData.isMeshGroup = true; return group
    }
    case 'powerpole': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.1, 3, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.3}))
      pole.position.y = 1.5
      const cross1 = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.05, 0.06), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      cross1.position.y = 2.8
      const cross2 = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.05, 0.06), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      cross2.position.y = 2.4
      for (let i = 0; i < 3; i++) {
        const insulator = new THREE.Mesh(new THREE.CylinderGeometry(0.025, 0.03, 0.08, 8), new THREE.MeshStandardMaterial({color:0xcccccc}))
        insulator.position.set(-0.3 + i * 0.3, 2.85, 0)
        group.add(insulator)
      }
      group.add(pole, cross1, cross2)
      group.userData.isMeshGroup = true; return group
    }
    case 'enginegen': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.4, 0.3), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.5}))
      body.position.y = 0.2
      const head = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.15, 0.2), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      head.position.set(-0.35, 0.2, 0)
      const exhaust = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.2, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      exhaust.position.set(0.2, 0.5, 0)
      group.add(body, head, exhaust)
      group.userData.isMeshGroup = true; return group
    }
    case 'solararray': {
      const group = new THREE.Group()
      for (let r = 0; r < 2; r++) {
        for (let c = 0; c < 3; c++) {
          const panel = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.02, 0.25), new THREE.MeshStandardMaterial({color:0x114488, metalness:0.3, emissive:0x001122, emissiveIntensity:0.2}))
          panel.position.set((c - 1) * 0.4, 0.01, (r - 0.5) * 0.3)
          const frame = new THREE.Mesh(new THREE.BoxGeometry(0.38, 0.015, 0.28), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7, wireframe:true}))
          frame.position.copy(panel.position)
          group.add(panel, frame)
        }
      }
      const stand = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.03, 0.8), new THREE.MeshStandardMaterial({color:0x777777, metalness:0.7}))
      stand.position.set(-0.2, -0.05, 0.15)
      group.add(stand)
      group.userData.isMeshGroup = true; return group
    }
    case 'substation': {
      const group = new THREE.Group()
      const fence = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.4, 0.4), new THREE.MeshStandardMaterial({color:0x888888, wireframe:true, transparent:true, opacity:0.6}))
      fence.position.y = 0.2
      const tx = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.3, 0.15), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.4}))
      tx.position.set(-0.15, 0.15, 0)
      const tx2 = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.3, 0.15), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.4}))
      tx2.position.set(0.15, 0.15, 0)
      group.add(fence, tx, tx2)
      group.userData.isMeshGroup = true; return group
    }
    case 'batterystorage': {
      const group = new THREE.Group()
      const container = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.25, 0.3), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.3}))
      container.position.y = 0.125
      for (let i = 0; i < 4; i++) {
        const cell = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.15, 0.15), new THREE.MeshStandardMaterial({color:0x33aa55, metalness:0.2, emissive:0x002200, emissiveIntensity:0.2}))
        cell.position.set(-0.15 + i * 0.1, 0.15, 0)
        group.add(cell)
      }
      group.add(container)
      group.userData.isMeshGroup = true; return group
    }
    case 'powermeter': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.15, 0.08), new THREE.MeshStandardMaterial({color:0xfffb4d}))
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.06, 0.01), new THREE.MeshStandardMaterial({color:0x222222, emissive:0x111111, emissiveIntensity:0.3}))
      screen.position.set(0, 0.03, 0.05)
      group.add(body, screen)
      group.userData.isMeshGroup = true; return group
    }
    case 'circuitbreaker': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.15, 0.08), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.3}))
      const lever = new THREE.Mesh(new THREE.BoxGeometry(0.015, 0.05, 0.04), new THREE.MeshStandardMaterial({color:0xcccccc}))
      lever.position.set(0, 0.08, 0.06)
      group.add(body, lever)
      group.userData.isMeshGroup = true; return group
    }
    case 'insulator': {
      const group = new THREE.Group()
      for (let i = 0; i < 3; i++) {
        const disc = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.07, 0.03, 16), new THREE.MeshStandardMaterial({color:0xdddddd}))
        disc.position.y = i * 0.07
        group.add(disc)
      }
      const rod = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.2, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      rod.position.y = 0.07
      group.add(rod)
      group.userData.isMeshGroup = true; return group
    }
    case 'cabletray': {
      const group = new THREE.Group()
      const tray = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.03, 0.15), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      tray.position.y = 0.015
      const side1 = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.05, 0.02), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      side1.position.set(0, 0.05, 0.065)
      const side2 = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.05, 0.02), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      side2.position.set(0, 0.05, -0.065)
      group.add(tray, side1, side2)
      group.userData.isMeshGroup = true; return group
    }
    case 'junctionbox': {
      const group = new THREE.Group()
      const box = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.08, 0.08), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.4}))
      const lid = new THREE.Mesh(new THREE.BoxGeometry(0.11, 0.01, 0.09), new THREE.MeshStandardMaterial({color:0x888888}))
      lid.position.y = 0.045
      group.add(box, lid)
      group.userData.isMeshGroup = true; return group
    }
    case 'lightningrod': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.1, 0.1, 12), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.5}))
      const rod = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.03, 0.8, 8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.8}))
      rod.position.y = 0.45
      const tip = new THREE.Mesh(new THREE.ConeGeometry(0.03, 0.1, 8), new THREE.MeshStandardMaterial({color:0xffffff, metalness:0.8, emissive:0x333333, emissiveIntensity:0.3}))
      tip.position.y = 0.9
      group.add(base, rod, tip)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 仓储物流 =====
    case 'warehouserack': {
      const group = new THREE.Group()
      for (let level = 0; level < 3; level++) {
        const shelf = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.03, 0.35), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.5}))
        shelf.position.y = 0.2 + level * 0.35
        group.add(shelf)
      }
      for (let i = -1; i <= 1; i+=2) {
        const pillar = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1.2, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
        pillar.position.set(i * 0.23, 0.6, 0)
        group.add(pillar)
      }
      group.userData.isMeshGroup = true; return group
    }
    case 'forklift': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.15, 0.25), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.4}))
      body.position.set(-0.05, 0.12, 0)
      const mast = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.5, 0.06), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      mast.position.set(0.22, 0.35, 0)
      const fork1 = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.02, 0.03), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8}))
      fork1.position.set(0.28, 0.08, 0.05)
      const fork2 = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.02, 0.03), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8}))
      fork2.position.set(0.28, 0.08, -0.05)
      for (let i = -1; i <= 1; i+=2) {
        const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.03, 12), new THREE.MeshStandardMaterial({color:0x222222}))
        wheel.position.set(i * 0.12, 0.04, 0.15); wheel.rotation.z = Math.PI/2
        group.add(wheel)
      }
      group.add(body, mast, fork1, fork2)
      group.userData.isMeshGroup = true; return group
    }
    case 'towercrane': {
      const group = new THREE.Group()
      const mast = new THREE.Mesh(new THREE.BoxGeometry(0.15, 2, 0.15), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.6}))
      mast.position.y = 1
      const jib = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.08, 0.08), new THREE.MeshStandardMaterial({color:0xff8800, metalness:0.6}))
      jib.position.set(0.5, 2.04, 0)
      const counterJib = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.06, 0.06), new THREE.MeshStandardMaterial({color:0xff8800, metalness:0.6}))
      counterJib.position.set(-0.2, 2.03, 0)
      const cab = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.15), new THREE.MeshStandardMaterial({color:0xffcc00, metalness:0.5}))
      cab.position.y = 2.05
      const hook = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.4, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8}))
      hook.position.set(0.7, 1.8, 0)
      group.add(mast, jib, counterJib, cab, hook)
      group.userData.isMeshGroup = true; return group
    }
    case 'container': {
      const group = new THREE.Group()
      const box = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.5, 1.2), new THREE.MeshStandardMaterial({color:0x4da6ff, metalness:0.4, roughness:0.6}))
      box.position.y = 0.25
      const door1 = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.4, 0.02), new THREE.MeshStandardMaterial({color:0x3377cc, metalness:0.4}))
      door1.position.set(-0.12, 0.25, 0.61)
      const door2 = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.4, 0.02), new THREE.MeshStandardMaterial({color:0x3377cc, metalness:0.4}))
      door2.position.set(0.12, 0.25, 0.61)
      const corner = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.55, 0.05), new THREE.MeshStandardMaterial({color:0x225599, metalness:0.5}))
      corner.position.set(-0.28, 0.25, -0.58)
      const corner2 = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.55, 0.05), new THREE.MeshStandardMaterial({color:0x225599, metalness:0.5}))
      corner2.position.set(0.28, 0.25, -0.58)
      const corner3 = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.55, 0.05), new THREE.MeshStandardMaterial({color:0x225599, metalness:0.5}))
      corner3.position.set(-0.28, 0.25, 0.58)
      const corner4 = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.55, 0.05), new THREE.MeshStandardMaterial({color:0x225599, metalness:0.5}))
      corner4.position.set(0.28, 0.25, 0.58)
      group.add(box, door1, door2, corner, corner2, corner3, corner4)
      group.userData.isMeshGroup = true; return group
    }
    case 'loadingdock': {
      const group = new THREE.Group()
      const platform = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.3, 0.5), new THREE.MeshStandardMaterial({color:0x888888, roughness:0.8}))
      platform.position.y = 0.15
      const ramp = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.15, 0.4), new THREE.MeshStandardMaterial({color:0x999999}))
      ramp.position.set(0, 0.075, -0.45)
      const bumper = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.1, 0.4), new THREE.MeshStandardMaterial({color:0xfffb4d}))
      bumper.position.set(0, 0.3, 0)
      group.add(platform, ramp, bumper)
      group.userData.isMeshGroup = true; return group
    }
    case 'pallet': {
      const group = new THREE.Group()
      const deck = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.03, 0.5), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.8}))
      deck.position.y = 0.08
      for (let i = -1; i <= 1; i+=2) {
        const stringer = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.06, 0.05), new THREE.MeshStandardMaterial({color:0x8c5523, roughness:0.8}))
        stringer.position.set(0, 0.03, i * 0.15)
        group.add(stringer)
        const block = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.06, 0.05), new THREE.MeshStandardMaterial({color:0x6b3a1f}))
        block.position.set(i * 0.19, 0.03, 0)
        group.add(block)
      }
      group.add(deck)
      group.userData.isMeshGroup = true; return group
    }
    case 'shelving': {
      const group = new THREE.Group()
      for (let s = 0; s < 3; s++) {
        const board = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.02, 0.25), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.7}))
        board.position.y = 0.15 + s * 0.3
        group.add(board)
      }
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.9, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
          leg.position.set(i * 0.22, 0.45, j * 0.1)
          group.add(leg)
        }
      }
      group.userData.isMeshGroup = true; return group
    }
    case 'barcodereader': {
      const group = new THREE.Group()
      const handle = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.03, 0.2, 8), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.4}))
      handle.position.y = 0.15
      const head = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.06, 0.04), new THREE.MeshStandardMaterial({color:0xcc0000}))
      head.position.y = 0.27
      const lens = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.04, 8), new THREE.MeshStandardMaterial({color:0xffffff, emissive:0xff0000, emissiveIntensity:0.6}))
      lens.position.set(0, 0.27, 0.04); lens.rotation.x = Math.PI/2
      group.add(handle, head, lens)
      group.userData.isMeshGroup = true; return group
    }
    case 'sorter': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.06, 0.3), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.5}))
      frame.position.y = 0.25
      for (let i = 0; i < 3; i++) {
        const chute = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.03, 0.15), new THREE.MeshStandardMaterial({color:0x33aa55, metalness:0.4}))
        chute.position.set(-0.2 + i * 0.2, 0.3, 0.2)
        group.add(chute)
      }
      const scanner = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.08, 0.06), new THREE.MeshStandardMaterial({color:0x4dfff0, emissive:0x003333, emissiveIntensity:0.3}))
      scanner.position.set(0, 0.35, -0.05)
      group.add(frame, scanner)
      group.userData.isMeshGroup = true; return group
    }
    case 'stackercrane': {
      const group = new THREE.Group()
      const rail = new THREE.Mesh(new THREE.BoxGeometry(1, 0.04, 0.04), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.6}))
      rail.position.y = 0.02
      const pillar1 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 1.2, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pillar1.position.set(-0.45, 0.6, 0)
      const pillar2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 1.2, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pillar2.position.set(0.45, 0.6, 0)
      const carriage = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.1, 0.1), new THREE.MeshStandardMaterial({color:0xff8800, metalness:0.5}))
      carriage.position.set(0, 1.1, 0)
      const fork = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.02, 0.15), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.8}))
      fork.position.set(0.08, 1.05, 0)
      group.add(rail, pillar1, pillar2, carriage, fork)
      group.userData.isMeshGroup = true; return group
    }
    case 'dockleveler': {
      const group = new THREE.Group()
      const plate = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.03, 0.5), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      plate.position.y = 0.015
      const lip = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.02, 0.1), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.6}))
      lip.position.set(0, -0.005, 0.3)
      const cyl = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.1, 8), new THREE.MeshStandardMaterial({color:0xffaa00}))
      cyl.position.set(-0.25, 0.02, -0.2); cyl.rotation.x = Math.PI/2
      const cyl2 = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.1, 8), new THREE.MeshStandardMaterial({color:0xffaa00}))
      cyl2.position.set(0.25, 0.02, -0.2); cyl2.rotation.x = Math.PI/2
      group.add(plate, lip, cyl, cyl2)
      group.userData.isMeshGroup = true; return group
    }
    case 'conveyorjunc': {
      const group = new THREE.Group()
      const main = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.04, 0.2), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      main.position.y = 0.15
      const branch = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.04, 0.4), new THREE.MeshStandardMaterial({color:0x777777, metalness:0.5}))
      branch.position.set(0, 0.15, 0.25)
      const roller = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.16, 12), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      roller.position.set(0, 0.15, 0.1); roller.rotation.z = Math.PI/2
      group.add(main, branch, roller)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 交通设施 =====
    case 'car': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.2, 1), new THREE.MeshStandardMaterial({color:0x4da6ff, metalness:0.5}))
      body.position.y = 0.14
      const cabin = new THREE.Mesh(new THREE.BoxGeometry(0.42, 0.14, 0.4), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.6}))
      cabin.position.set(0, 0.3, 0.1)
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.07, 0.07, 0.06, 16), new THREE.MeshStandardMaterial({color:0x111111}))
          wheel.position.set(i * 0.22, 0.07, j * 0.35); wheel.rotation.z = Math.PI/2
          group.add(wheel)
        }
      }
      group.add(body, cabin)
      group.userData.isMeshGroup = true; return group
    }
    case 'suv': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.55, 0.25, 1.1), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.5}))
      body.position.y = 0.16
      const cabin = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.2, 0.5), new THREE.MeshStandardMaterial({color:0xaaffcc, transparent:true, opacity:0.5}))
      cabin.position.set(0, 0.4, 0.1)
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.08, 16), new THREE.MeshStandardMaterial({color:0x111111}))
          wheel.position.set(i * 0.24, 0.08, j * 0.38); wheel.rotation.z = Math.PI/2
          group.add(wheel)
        }
      }
      group.add(body, cabin)
      group.userData.isMeshGroup = true; return group
    }
    case 'bus': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.7, 0.5, 2), new THREE.MeshStandardMaterial({color:0xfffb4d, metalness:0.3}))
      body.position.y = 0.28
      for (let r = 0; r < 4; r++) {
        const win = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.2, 0.02), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.5}))
        win.position.set(0.3, 0.5, -0.7 + r * 0.45)
        group.add(win)
        const win2 = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.2, 0.02), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.5}))
        win2.position.set(-0.3, 0.5, -0.7 + r * 0.45)
        group.add(win2)
      }
      for (let i = -1; i <= 1; i+=2) {
        const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.08, 16), new THREE.MeshStandardMaterial({color:0x222222}))
        wheel.position.set(i * 0.32, 0.1, -0.6); wheel.rotation.z = Math.PI/2
        group.add(wheel)
        const wheel2 = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.08, 16), new THREE.MeshStandardMaterial({color:0x222222}))
        wheel2.position.set(i * 0.32, 0.1, 0.6); wheel2.rotation.z = Math.PI/2
        group.add(wheel2)
      }
      group.add(body)
      group.userData.isMeshGroup = true; return group
    }
    case 'truck': {
      const group = new THREE.Group()
      const cab = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.35, 0.4), new THREE.MeshStandardMaterial({color:0xffaa4d, metalness:0.4}))
      cab.position.set(0, 0.2, -0.5)
      const cargo = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.3, 1.2), new THREE.MeshStandardMaterial({color:0xdd8822, metalness:0.4}))
      cargo.position.set(0, 0.18, 0.3)
      for (let i = -1; i <= 1; i+=2) {
        const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.06, 16), new THREE.MeshStandardMaterial({color:0x111111}))
        wheel.position.set(i * 0.22, 0.08, -0.5); wheel.rotation.z = Math.PI/2
        group.add(wheel)
      }
      for (let i = -1; i <= 1; i+=2) {
        const wheel2 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.07, 16), new THREE.MeshStandardMaterial({color:0x111111}))
        wheel2.position.set(i * 0.28, 0.08, 0.3); wheel2.rotation.z = Math.PI/2
        group.add(wheel2)
        const wheel3 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.07, 16), new THREE.MeshStandardMaterial({color:0x111111}))
        wheel3.position.set(i * 0.28, 0.08, 0.8); wheel3.rotation.z = Math.PI/2
        group.add(wheel3)
      }
      group.add(cab, cargo)
      group.userData.isMeshGroup = true; return group
    }
    case 'firetruck': {
      const group = new THREE.Group()
      const cab = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.35, 0.4), new THREE.MeshStandardMaterial({color:0xff0000, metalness:0.5}))
      cab.position.set(0, 0.2, -0.5)
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.35, 1), new THREE.MeshStandardMaterial({color:0xcc0000, metalness:0.5}))
      body.position.set(0, 0.2, 0.2)
      const ladder = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.04, 0.8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      ladder.position.set(0, 0.45, 0.2)
      const light = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.04, 0.04), new THREE.MeshStandardMaterial({color:0xff4444, emissive:0xff0000, emissiveIntensity:0.8}))
      light.position.set(0, 0.3, -0.7)
      for (let i = -1; i <= 1; i+=2) {
        const w = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.06, 16), new THREE.MeshStandardMaterial({color:0x111111}))
        w.position.set(i * 0.22, 0.08, -0.5); w.rotation.z = Math.PI/2
        group.add(w)
        const w2 = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.07, 16), new THREE.MeshStandardMaterial({color:0x111111}))
        w2.position.set(i * 0.28, 0.08, 0.7); w2.rotation.z = Math.PI/2
        group.add(w2)
      }
      group.add(cab, body, ladder, light)
      group.userData.isMeshGroup = true; return group
    }
    case 'ambulance': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.55, 0.35, 1.1), new THREE.MeshStandardMaterial({color:0xffffff, metalness:0.3}))
      body.position.y = 0.2
      const stripe = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.1, 0.8), new THREE.MeshStandardMaterial({color:0xff0000}))
      stripe.position.set(0, 0.25, 0)
      const cross = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.2, 0.02), new THREE.MeshStandardMaterial({color:0xff0000}))
      cross.position.set(0, 0.35, 0.56)
      const crossH = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.05, 0.02), new THREE.MeshStandardMaterial({color:0xff0000}))
      crossH.position.set(0, 0.35, 0.56)
      const light = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.04, 0.03), new THREE.MeshStandardMaterial({color:0x4dfff0, emissive:0x0044ff, emissiveIntensity:0.8}))
      light.position.set(0, 0.3, -0.56)
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const w = new THREE.Mesh(new THREE.CylinderGeometry(0.07, 0.07, 0.06, 16), new THREE.MeshStandardMaterial({color:0x111111}))
          w.position.set(i * 0.24, 0.07, j * 0.38); w.rotation.z = Math.PI/2
          group.add(w)
        }
      }
      group.add(body, stripe, cross, crossH, light)
      group.userData.isMeshGroup = true; return group
    }
    case 'policecar': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.2, 1), new THREE.MeshStandardMaterial({color:0xffffff, metalness:0.4}))
      body.position.y = 0.14
      const stripe = new THREE.Mesh(new THREE.BoxGeometry(0.44, 0.06, 0.6), new THREE.MeshStandardMaterial({color:0x000088}))
      stripe.position.set(0, 0.18, 0)
      const lightBar = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.03, 0.05), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.6}))
      lightBar.position.set(0, 0.32, -0.25)
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const w = new THREE.Mesh(new THREE.CylinderGeometry(0.07, 0.07, 0.06, 16), new THREE.MeshStandardMaterial({color:0x111111}))
          w.position.set(i * 0.22, 0.07, j * 0.35); w.rotation.z = Math.PI/2
          group.add(w)
        }
      }
      group.add(body, stripe, lightBar)
      group.userData.isMeshGroup = true; return group
    }
    case 'motorcycle': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.15, 0.5), new THREE.MeshStandardMaterial({color:0xff4d4f, metalness:0.5}))
      body.position.y = 0.15
      const seat = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.05, 0.2), new THREE.MeshStandardMaterial({color:0x222222}))
      seat.position.set(0, 0.25, -0.05)
      const handle = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.2, 8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      handle.position.set(0, 0.3, 0.1); handle.rotation.x = Math.PI/3
      const wheel1 = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.05, 16), new THREE.MeshStandardMaterial({color:0x111111}))
      wheel1.position.set(0, 0.1, 0.22); wheel1.rotation.z = Math.PI/2
      const wheel2 = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.05, 16), new THREE.MeshStandardMaterial({color:0x111111}))
      wheel2.position.set(0, 0.1, -0.22); wheel2.rotation.z = Math.PI/2
      group.add(body, seat, handle, wheel1, wheel2)
      group.userData.isMeshGroup = true; return group
    }
    case 'bicycle': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.5, 8), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.6}))
      frame.position.set(0, 0.25, 0); frame.rotation.z = Math.PI/6
      const frame2 = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.4, 8), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.6}))
      frame2.position.set(0.05, 0.15, 0.1); frame2.rotation.z = -Math.PI/5
      const seat = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.02, 0.1), new THREE.MeshStandardMaterial({color:0x222222}))
      seat.position.set(0.05, 0.35, -0.05)
      const handle = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.15, 8), new THREE.MeshStandardMaterial({color:0xcccccc}))
      handle.position.set(0, 0.38, 0.15); handle.rotation.z = Math.PI/4
      const w1 = new THREE.Mesh(new THREE.TorusGeometry(0.12, 0.01, 8, 16), new THREE.MeshStandardMaterial({color:0x333333}))
      w1.position.set(0, 0.12, 0.2)
      const w2 = new THREE.Mesh(new THREE.TorusGeometry(0.12, 0.01, 8, 16), new THREE.MeshStandardMaterial({color:0x333333}))
      w2.position.set(0, 0.12, -0.15)
      group.add(frame, frame2, seat, handle, w1, w2)
      group.userData.isMeshGroup = true; return group
    }
    case 'train': {
      const group = new THREE.Group()
      const head = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.5, 0.6), new THREE.MeshStandardMaterial({color:0x4dffa6, metalness:0.5}))
      head.position.set(0, 0.25, 0.8)
      const nose = new THREE.Mesh(new THREE.ConeGeometry(0.25, 0.4, 4), new THREE.MeshStandardMaterial({color:0x33aa55, metalness:0.5}))
      nose.position.set(0, 0.5, 1.3); nose.rotation.x = -Math.PI/2
      const car1 = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.45, 0.8), new THREE.MeshStandardMaterial({color:0x44cc66, metalness:0.5}))
      car1.position.set(0, 0.225, -0.1)
      const car2 = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.45, 0.8), new THREE.MeshStandardMaterial({color:0x44cc66, metalness:0.5}))
      car2.position.set(0, 0.225, -1)
      for (let i = -1; i <= 1; i+=2) {
        for (let pos = -0.7; pos <= 1.1; pos += 0.9) {
          const w = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.04, 12), new THREE.MeshStandardMaterial({color:0x333333}))
          w.position.set(i * 0.22, 0.08, pos); w.rotation.z = Math.PI/2
          group.add(w)
        }
      }
      group.add(head, nose, car1, car2)
      group.userData.isMeshGroup = true; return group
    }
    case 'metro': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.45, 0.4, 2.2), new THREE.MeshStandardMaterial({color:0x4da6ff, metalness:0.5}))
      body.position.y = 0.2
      for (let r = 0; r < 5; r++) {
        const win = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.15, 0.02), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.5}))
        win.position.set(0.18, 0.3, -0.9 + r * 0.4)
        group.add(win)
        const win2 = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.15, 0.02), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.5}))
        win2.position.set(-0.18, 0.3, -0.9 + r * 0.4)
        group.add(win2)
      }
      for (let i = -1; i <= 1; i+=2) {
        for (let p = -0.8; p <= 0.8; p += 0.55) {
          const w = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.05, 0.04, 12), new THREE.MeshStandardMaterial({color:0x333333}))
          w.position.set(i * 0.2, 0.05, p); w.rotation.z = Math.PI/2
          group.add(w)
        }
      }
      group.add(body)
      group.userData.isMeshGroup = true; return group
    }
    case 'airplane': {
      const group = new THREE.Group()
      const fuselage = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.06, 1.8, 12), new THREE.MeshStandardMaterial({color:0xffffff, metalness:0.3}))
      fuselage.rotation.x = Math.PI/2
      const nose = new THREE.Mesh(new THREE.SphereGeometry(0.06, 8, 8, 0, Math.PI*2, 0, Math.PI/2), new THREE.MeshStandardMaterial({color:0xffffff}))
      nose.position.set(0, 0, 0.9); nose.rotation.x = Math.PI
      const wing = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.02, 0.2), new THREE.MeshStandardMaterial({color:0xdddddd, metalness:0.3}))
      wing.position.set(0, 0, -0.1)
      const tailV = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.2, 0.15), new THREE.MeshStandardMaterial({color:0xdddddd}))
      tailV.position.set(0, 0.15, -0.7)
      const tailH = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.02, 0.1), new THREE.MeshStandardMaterial({color:0xdddddd}))
      tailH.position.set(0, 0.05, -0.75)
      group.add(fuselage, nose, wing, tailV, tailH)
      group.userData.isMeshGroup = true; return group
    }
    case 'helicopter': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.5), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      body.position.y = 0.1
      const tail = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.03, 0.4, 8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.5}))
      tail.position.set(0, 0.1, -0.45); tail.rotation.x = Math.PI/2
      const rotor = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.01, 0.06), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.7}))
      rotor.position.y = 0.25
      const rotor2 = new THREE.Mesh(new THREE.BoxGeometry(0.01, 0.01, 0.3), new THREE.MeshStandardMaterial({color:0xaaaaaa}))
      rotor2.position.set(0, 0.1, -0.6)
      group.add(body, tail, rotor, rotor2)
      group.userData.isMeshGroup = true; return group
    }
    case 'drone': {
      const group = new THREE.Group()
      const center = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.04, 0.1), new THREE.MeshStandardMaterial({color:0xc17bff, metalness:0.5}))
      for (let i = 0; i < 4; i++) {
        const arm = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.01, 0.3, 8), new THREE.MeshStandardMaterial({color:0x888888}))
        arm.rotation.set(0, i * Math.PI/2, Math.PI/2)
        group.add(arm)
        const prop = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.005, 0.02), new THREE.MeshStandardMaterial({color:0xdddddd}))
        prop.position.set(0, 0.02, 0.15); prop.rotation.y = i * Math.PI/2
        group.add(prop)
      }
      const cam = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.03, 0.04, 8), new THREE.MeshStandardMaterial({color:0x111111}))
      cam.position.y = -0.04
      group.add(center, cam)
      group.userData.isMeshGroup = true; return group
    }
    case 'ship': {
      const group = new THREE.Group()
      const hull = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.2, 1.5), new THREE.MeshStandardMaterial({color:0x4da6ff, metalness:0.4}))
      hull.position.y = 0.1
      const deck = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.05, 1), new THREE.MeshStandardMaterial({color:0xdddddd}))
      deck.position.y = 0.22
      const cabin = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.3, 0.5), new THREE.MeshStandardMaterial({color:0xffffff}))
      cabin.position.set(0, 0.4, 0.2)
      const funnel = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.05, 0.3, 12), new THREE.MeshStandardMaterial({color:0xff4d4f}))
      funnel.position.set(0, 0.65, -0.1)
      group.add(hull, deck, cabin, funnel)
      group.userData.isMeshGroup = true; return group
    }
    case 'sailboat': {
      const group = new THREE.Group()
      const hull = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.1, 0.6), new THREE.MeshStandardMaterial({color:0x4dfff0, metalness:0.3}))
      hull.position.y = 0.05
      const mast = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.02, 0.8, 8), new THREE.MeshStandardMaterial({color:0xa67c52}))
      mast.position.y = 0.4
      const sail = new THREE.Mesh(new THREE.PlaneGeometry(0.4, 0.5), new THREE.MeshStandardMaterial({color:0xffffff, side:THREE.DoubleSide}))
      sail.position.set(0.02, 0.5, 0); sail.rotation.y = Math.PI/12
      group.add(hull, mast, sail)
      group.userData.isMeshGroup = true; return group
    }
    case 'trafficcone': {
      const group = new THREE.Group()
      const cone = new THREE.Mesh(new THREE.ConeGeometry(0.1, 0.3, 12), new THREE.MeshStandardMaterial({color:0xffaa4d, roughness:0.6}))
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.15, 0.04, 16), new THREE.MeshStandardMaterial({color:0xff8800}))
      const stripe = new THREE.Mesh(new THREE.TorusGeometry(0.07, 0.01, 6, 12), new THREE.MeshStandardMaterial({color:0xffffff}))
      stripe.position.y = 0.15; stripe.rotation.x = Math.PI/2
      group.add(cone, base, stripe)
      group.userData.isMeshGroup = true; return group
    }
    case 'barrier': {
      const group = new THREE.Group()
      const bar = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.04, 0.1), new THREE.MeshStandardMaterial({color:0xfffb4d}))
      bar.position.y = 0.5
      const stripe1 = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.02, 0.01), new THREE.MeshStandardMaterial({color:0xff0000}))
      stripe1.position.set(-0.3, 0.52, 0.05)
      const stripe2 = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.02, 0.01), new THREE.MeshStandardMaterial({color:0xff0000}))
      stripe2.position.set(0.3, 0.52, 0.05)
      const post1 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.5, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      post1.position.set(-0.55, 0.25, 0)
      const post2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.5, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      post2.position.set(0.55, 0.25, 0)
      group.add(bar, stripe1, stripe2, post1, post2)
      group.userData.isMeshGroup = true; return group
    }
    case 'speedbump': {
      const group = new THREE.Group()
      const bump = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.04, 0.15), new THREE.MeshStandardMaterial({color:0xfffb4d}))
      const stripe1 = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.01, 0.03), new THREE.MeshStandardMaterial({color:0xff0000}))
      stripe1.position.z = 0.04
      const stripe2 = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.01, 0.03), new THREE.MeshStandardMaterial({color:0xff0000}))
      stripe2.position.z = -0.04
      group.add(bump, stripe1, stripe2)
      group.userData.isMeshGroup = true; return group
    }
    case 'zebracrossing': {
      const group = new THREE.Group()
      for (let i = 0; i < 6; i++) {
        const stripe = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.01, 0.8), new THREE.MeshStandardMaterial({color:0xffffff}))
        stripe.position.set(-0.25 + i * 0.1, 0.005, 0)
        group.add(stripe)
      }
      group.userData.isMeshGroup = true; return group
    }

    // ===== 室内设施 =====
    case 'desk': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.03, 0.5), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.6}))
      top.position.y = 0.45
      for (let i = -1; i <= 1; i+=2) {
        const leg = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.45, 0.04), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
        leg.position.set(i * 0.35, 0.225, 0.2)
        group.add(leg)
        const leg2 = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.45, 0.04), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
        leg2.position.set(i * 0.35, 0.225, -0.2)
        group.add(leg2)
      }
      group.add(top)
      group.userData.isMeshGroup = true; return group
    }
    case 'officechair': {
      const group = new THREE.Group()
      const seat = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.04, 24), new THREE.MeshStandardMaterial({color:0x444444}))
      seat.position.y = 0.32
      const back = new THREE.Mesh(new THREE.BoxGeometry(0.22, 0.25, 0.03), new THREE.MeshStandardMaterial({color:0x444444}))
      back.position.set(0, 0.45, -0.12)
      const pillar = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.3, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pillar.position.y = 0.15
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.22, 0.03, 16), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      for (let i = 0; i < 5; i++) {
        const spoke = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.01, 0.2, 6), new THREE.MeshStandardMaterial({color:0x999999}))
        spoke.position.y = 0.02; spoke.rotation.set(0, i * Math.PI*2/5, Math.PI/2)
        group.add(spoke)
      }
      group.add(seat, back, pillar, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'sofa': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.2, 0.5), new THREE.MeshStandardMaterial({color:0xffaa4d, roughness:0.7}))
      base.position.y = 0.2
      const back = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.4, 0.1), new THREE.MeshStandardMaterial({color:0xff8800, roughness:0.7}))
      back.position.set(0, 0.45, -0.2)
      const armL = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.3, 0.4), new THREE.MeshStandardMaterial({color:0xff8800, roughness:0.7}))
      armL.position.set(-0.58, 0.3, 0)
      const armR = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.3, 0.4), new THREE.MeshStandardMaterial({color:0xff8800, roughness:0.7}))
      armR.position.set(0.58, 0.3, 0)
      group.add(base, back, armL, armR)
      group.userData.isMeshGroup = true; return group
    }
    case 'coffeetable': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.03, 0.4), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.6}))
      top.position.y = 0.3
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.3, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
          leg.position.set(i * 0.23, 0.15, j * 0.14)
          group.add(leg)
        }
      }
      group.add(top)
      group.userData.isMeshGroup = true; return group
    }
    case 'bookshelf': {
      const group = new THREE.Group()
      for (let s = 0; s < 4; s++) {
        const shelf = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.02, 0.25), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.6}))
        shelf.position.y = 0.15 + s * 0.25
        group.add(shelf)
        if (s < 3) {
          for (let b = 0; b < 4; b++) {
            const book = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.18, 0.12), new THREE.MeshStandardMaterial({color:[0x4da6ff, 0xff4d4f, 0x4dffa6, 0xfffb4d][b]}))
            book.position.set(-0.15 + b * 0.08, 0.3 + s * 0.25, 0.05)
            group.add(book)
          }
        }
      }
      const side1 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.9, 0.25), new THREE.MeshStandardMaterial({color:0x8c5523}))
      side1.position.x = -0.24
      const side2 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.9, 0.25), new THREE.MeshStandardMaterial({color:0x8c5523}))
      side2.position.x = 0.24
      group.add(side1, side2)
      group.userData.isMeshGroup = true; return group
    }
    case 'cabinet': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.8, 0.35), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.6}))
      body.position.y = 0.4
      const doorL = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.7, 0.02), new THREE.MeshStandardMaterial({color:0x8c5523}))
      doorL.position.set(-0.13, 0.4, 0.18)
      const doorR = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.7, 0.02), new THREE.MeshStandardMaterial({color:0x8c5523}))
      doorR.position.set(0.13, 0.4, 0.18)
      group.add(body, doorL, doorR)
      group.userData.isMeshGroup = true; return group
    }
    case 'bed': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.15, 1.2), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.6}))
      frame.position.y = 0.15
      const mattress = new THREE.Mesh(new THREE.BoxGeometry(0.7, 0.08, 1.1), new THREE.MeshStandardMaterial({color:0xffffff}))
      mattress.position.y = 0.25
      const headboard = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.3, 0.04), new THREE.MeshStandardMaterial({color:0x8c5523}))
      headboard.position.set(0, 0.35, -0.55)
      const pillow = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.04, 0.15), new THREE.MeshStandardMaterial({color:0xeeeeee}))
      pillow.position.set(0, 0.3, -0.35)
      group.add(frame, mattress, headboard, pillow)
      group.userData.isMeshGroup = true; return group
    }
    case 'floorlamp': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.14, 0.04, 16), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1.2, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      pole.position.y = 0.6
      const shade = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.2, 0.25, 16, 1, true), new THREE.MeshStandardMaterial({color:0xffffcc, emissive:0xfffb4d, emissiveIntensity:0.4, side:THREE.DoubleSide}))
      shade.position.y = 1.25
      group.add(base, pole, shade)
      group.userData.isMeshGroup = true; return group
    }
    case 'ceilinglight': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.04, 24), new THREE.MeshStandardMaterial({color:0xdddddd}))
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.02, 0.3), new THREE.MeshStandardMaterial({color:0xffffcc, emissive:0xfffb4d, emissiveIntensity:0.3}))
      panel.position.y = -0.04
      group.add(base, panel)
      group.userData.isMeshGroup = true; return group
    }
    case 'monitor': {
      const group = new THREE.Group()
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.25, 0.02), new THREE.MeshStandardMaterial({color:0x222222, emissive:0x111133, emissiveIntensity:0.4}))
      screen.position.y = 0.2
      const bezel = new THREE.Mesh(new THREE.BoxGeometry(0.37, 0.27, 0.03), new THREE.MeshStandardMaterial({color:0x111111, wireframe:true}))
      bezel.position.y = 0.2
      const stand = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.2, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      stand.position.y = 0.08
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.1, 0.03, 16), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      group.add(screen, bezel, stand, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'keyboard': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.02, 0.1), new THREE.MeshStandardMaterial({color:0x333333}))
      for (let i = 0; i < 4; i++) {
        for (let j = 0; j < 12; j++) {
          const key = new THREE.Mesh(new THREE.BoxGeometry(0.014, 0.005, 0.014), new THREE.MeshStandardMaterial({color:0x555555}))
          key.position.set(-0.09 + j * 0.014, 0.015, -0.03 + i * 0.018)
          group.add(key)
        }
      }
      group.add(body)
      group.userData.isMeshGroup = true; return group
    }
    case 'computer': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.12, 0.25, 0.22), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.5}))
      body.position.y = 0.125
      const powerLed = new THREE.Mesh(new THREE.SphereGeometry(0.01, 4, 4), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.8}))
      powerLed.position.set(0.05, 0.15, 0.12)
      group.add(body, powerLed)
      group.userData.isMeshGroup = true; return group
    }
    case 'printer': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.2, 0.2), new THREE.MeshStandardMaterial({color:0xdddddd, metalness:0.3}))
      body.position.y = 0.1
      const lid = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.02, 0.18), new THREE.MeshStandardMaterial({color:0xcccccc}))
      lid.position.set(0, 0.22, 0)
      const slot = new THREE.Mesh(new THREE.BoxGeometry(0.18, 0.01, 0.04), new THREE.MeshStandardMaterial({color:0x222222}))
      slot.position.set(0, 0.07, -0.12)
      group.add(body, lid, slot)
      group.userData.isMeshGroup = true; return group
    }
    case 'whiteboard': {
      const group = new THREE.Group()
      const board = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.5, 0.02), new THREE.MeshStandardMaterial({color:0xffffff}))
      board.position.y = 0.25
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.82, 0.52, 0.03), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5, wireframe:true}))
      frame.position.y = 0.25
      const stand = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.8, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
      stand.position.set(0, 0, -0.2)
      group.add(board, frame, stand)
      group.userData.isMeshGroup = true; return group
    }
    case 'meetingtable': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.04, 0.6), new THREE.MeshStandardMaterial({color:0xa67c52, roughness:0.6}))
      top.position.y = 0.5
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.5, 12), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.7}))
          leg.position.set(i * 0.5, 0.25, j * 0.22)
          group.add(leg)
        }
      }
      group.add(top)
      group.userData.isMeshGroup = true; return group
    }
    case 'waterdispenser': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.5, 0.2), new THREE.MeshStandardMaterial({color:0xffffff}))
      body.position.y = 0.25
      const bottle = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.25, 16), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.5}))
      bottle.position.y = 0.55
      const tap1 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.04, 0.04), new THREE.MeshStandardMaterial({color:0xff4444}))
      tap1.position.set(0.07, 0.3, 0.08)
      const tap2 = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.04, 0.04), new THREE.MeshStandardMaterial({color:0x4444ff}))
      tap2.position.set(-0.07, 0.3, 0.08)
      group.add(body, bottle, tap1, tap2)
      group.userData.isMeshGroup = true; return group
    }
    case 'microwave': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.2, 0.2), new THREE.MeshStandardMaterial({color:0xdddddd, metalness:0.3}))
      body.position.y = 0.1
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.15, 0.02), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.6}))
      door.position.set(0, 0.1, 0.1)
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.04, 0.01), new THREE.MeshStandardMaterial({color:0x222222}))
      panel.position.set(0.1, 0.1, 0.1)
      group.add(body, door, panel)
      group.userData.isMeshGroup = true; return group
    }
    case 'fridge': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.8, 0.35), new THREE.MeshStandardMaterial({color:0xdddddd, metalness:0.4}))
      body.position.y = 0.4
      const door1 = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.4, 0.02), new THREE.MeshStandardMaterial({color:0xcccccc}))
      door1.position.set(0, 0.55, 0.18)
      const door2 = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.3, 0.02), new THREE.MeshStandardMaterial({color:0xcccccc}))
      door2.position.set(0, 0.2, 0.18)
      group.add(body, door1, door2)
      group.userData.isMeshGroup = true; return group
    }
    case 'acunit': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.2, 0.12), new THREE.MeshStandardMaterial({color:0xeeeeee}))
      body.position.y = 0.1
      const vent = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.08, 0.01), new THREE.MeshStandardMaterial({color:0xcccccc}))
      vent.position.set(0, 0.1, 0.07)
      group.add(body, vent)
      group.userData.isMeshGroup = true; return group
    }
    case 'wallclock': {
      const group = new THREE.Group()
      const face = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.02, 24), new THREE.MeshStandardMaterial({color:0xffffff}))
      const hourHand = new THREE.Mesh(new THREE.BoxGeometry(0.01, 0.06, 0.005), new THREE.MeshStandardMaterial({color:0x111111}))
      hourHand.position.set(0, 0.04, 0.015)
      const minHand = new THREE.Mesh(new THREE.BoxGeometry(0.008, 0.08, 0.005), new THREE.MeshStandardMaterial({color:0x111111}))
      minHand.position.set(-0.01, 0.05, 0.015); minHand.rotation.z = Math.PI/4
      const dot = new THREE.Mesh(new THREE.SphereGeometry(0.01, 6, 6), new THREE.MeshStandardMaterial({color:0xff0000}))
      dot.position.z = 0.015
      group.add(face, hourHand, minHand, dot)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 消防安防 =====
    case 'fireextinguisher': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.07, 0.35, 16), new THREE.MeshStandardMaterial({color:0xff0000, metalness:0.3}))
      body.position.y = 0.175
      const neck = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.08, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      neck.position.y = 0.36
      const handle = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.02, 0.04), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      handle.position.set(0, 0.42, 0.02)
      const nozzle = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.08, 8), new THREE.MeshStandardMaterial({color:0x444444}))
      nozzle.position.set(0.04, 0.42, 0.04); nozzle.rotation.x = Math.PI/4
      group.add(body, neck, handle, nozzle)
      group.userData.isMeshGroup = true; return group
    }
    case 'firealarm': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.07, 0.07, 0.03, 16), new THREE.MeshStandardMaterial({color:0xff0000}))
      const dome = new THREE.Mesh(new THREE.SphereGeometry(0.06, 12, 8, 0, Math.PI*2, 0, Math.PI/2), new THREE.MeshStandardMaterial({color:0xcc0000, emissive:0xff0000, emissiveIntensity:0.3}))
      dome.position.y = 0.03
      group.add(base, dome)
      group.userData.isMeshGroup = true; return group
    }
    case 'sprinkler': {
      const group = new THREE.Group()
      const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.08, 8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      const head = new THREE.Mesh(new THREE.CylinderGeometry(0.025, 0.035, 0.03, 12), new THREE.MeshStandardMaterial({color:0xff0000, metalness:0.3}))
      head.position.y = -0.04
      const deflector = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.01, 16), new THREE.MeshStandardMaterial({color:0x888888}))
      deflector.position.y = -0.06
      group.add(pipe, head, deflector)
      group.userData.isMeshGroup = true; return group
    }
    case 'smokedetector': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.02, 16), new THREE.MeshStandardMaterial({color:0xffffff}))
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.05, 0.06, 0.03, 16), new THREE.MeshStandardMaterial({color:0xdddddd}))
      body.position.y = 0.02
      const led = new THREE.Mesh(new THREE.SphereGeometry(0.01, 4, 4), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.8}))
      led.position.set(0.04, 0.035, 0)
      group.add(base, body, led)
      group.userData.isMeshGroup = true; return group
    }
    case 'emergencysign': {
      const group = new THREE.Group()
      const board = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.1, 0.01), new THREE.MeshStandardMaterial({color:0x00cc00, emissive:0x003300, emissiveIntensity:0.4}))
      const arrow = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.02, 0.005), new THREE.MeshStandardMaterial({color:0xffffff}))
      arrow.position.set(0.05, 0, 0.005)
      const person = new THREE.Mesh(new THREE.BoxGeometry(0.015, 0.04, 0.005), new THREE.MeshStandardMaterial({color:0xffffff}))
      person.position.set(-0.03, 0, 0.005)
      group.add(board, arrow, person)
      group.userData.isMeshGroup = true; return group
    }
    case 'firehose': {
      const group = new THREE.Group()
      const reel = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.1, 16), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      const hose = new THREE.Mesh(new THREE.TorusGeometry(0.1, 0.015, 8, 16), new THREE.MeshStandardMaterial({color:0xffffff}))
      hose.position.y = 0.05; hose.rotation.x = Math.PI/2
      group.add(reel, hose)
      group.userData.isMeshGroup = true; return group
    }
    case 'firehydrantbox': {
      const group = new THREE.Group()
      const box = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.4, 0.1), new THREE.MeshStandardMaterial({color:0xff0000, metalness:0.3}))
      box.position.y = 0.2
      const glass = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.3, 0.02), new THREE.MeshStandardMaterial({color:0xffcccc, transparent:true, opacity:0.5}))
      glass.position.set(0, 0.2, 0.06)
      group.add(box, glass)
      group.userData.isMeshGroup = true; return group
    }
    case 'securitygate': {
      const group = new THREE.Group()
      const frameL = new THREE.Mesh(new THREE.BoxGeometry(0.06, 1, 0.06), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      frameL.position.set(-0.25, 0.5, 0)
      const frameR = new THREE.Mesh(new THREE.BoxGeometry(0.06, 1, 0.06), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.6}))
      frameR.position.set(0.25, 0.5, 0)
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.56, 0.05, 0.15), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      top.position.y = 1.03
      group.add(frameL, frameR, top)
      group.userData.isMeshGroup = true; return group
    }
    case 'accesscontrol': {
      const group = new THREE.Group()
      const panel = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.12, 0.03), new THREE.MeshStandardMaterial({color:0x444444}))
      panel.position.y = 0.06
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.04, 0.01), new THREE.MeshStandardMaterial({color:0x222222}))
      screen.position.set(0, 0.08, 0.02)
      const keypad = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.03, 0.005), new THREE.MeshStandardMaterial({color:0x666666}))
      keypad.position.set(0, 0.04, 0.02)
      const led = new THREE.Mesh(new THREE.SphereGeometry(0.006, 4, 4), new THREE.MeshStandardMaterial({color:0x00ff00, emissive:0x00ff00, emissiveIntensity:0.8}))
      led.position.set(0.02, 0.1, 0.02)
      group.add(panel, screen, keypad, led)
      group.userData.isMeshGroup = true; return group
    }
    case 'cctv': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.06, 0.08, 16), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.7}))
      const dome = new THREE.Mesh(new THREE.SphereGeometry(0.06, 12, 8, 0, Math.PI*2, 0, Math.PI/2), new THREE.MeshStandardMaterial({color:0x222222, transparent:true, opacity:0.7}))
      dome.position.y = 0.08
      const ring = new THREE.Mesh(new THREE.TorusGeometry(0.06, 0.01, 8, 16), new THREE.MeshStandardMaterial({color:0xff4444}))
      ring.position.y = 0.08
      group.add(base, dome, ring)
      group.userData.isMeshGroup = true; return group
    }
    case 'alarmbell': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.03, 16), new THREE.MeshStandardMaterial({color:0xff0000}))
      const bell = new THREE.Mesh(new THREE.SphereGeometry(0.06, 12, 8, 0, Math.PI*2, 0, Math.PI/2), new THREE.MeshStandardMaterial({color:0xcc0000, metalness:0.5}))
      bell.position.y = 0.03
      group.add(base, bell)
      group.userData.isMeshGroup = true; return group
    }
    case 'firstaidkit': {
      const group = new THREE.Group()
      const box = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.1, 0.1), new THREE.MeshStandardMaterial({color:0xffffff}))
      const cross = new THREE.Mesh(new THREE.BoxGeometry(0.02, 0.06, 0.01), new THREE.MeshStandardMaterial({color:0xff0000}))
      cross.position.set(0, 0, 0.06)
      const crossH = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.02, 0.01), new THREE.MeshStandardMaterial({color:0xff0000}))
      crossH.position.set(0, 0, 0.06)
      group.add(box, cross, crossH)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 绿化景观 =====
    case 'bush': {
      const group = new THREE.Group()
      const f1 = new THREE.Mesh(new THREE.SphereGeometry(0.2, 10, 10), new THREE.MeshStandardMaterial({color:0x33aa33}))
      f1.position.set(0, 0.1, 0)
      const f2 = new THREE.Mesh(new THREE.SphereGeometry(0.18, 10, 10), new THREE.MeshStandardMaterial({color:0x44bb44}))
      f2.position.set(0.12, 0.12, 0.08)
      const f3 = new THREE.Mesh(new THREE.SphereGeometry(0.16, 10, 10), new THREE.MeshStandardMaterial({color:0x228822}))
      f3.position.set(-0.1, 0.08, -0.06)
      group.add(f1, f2, f3)
      group.userData.isMeshGroup = true; return group
    }
    case 'hedge': {
      const group = new THREE.Group()
      const main = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.3, 0.25), new THREE.MeshStandardMaterial({color:0x33aa33}))
      main.position.y = 0.15
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.75, 0.1, 0.2), new THREE.MeshStandardMaterial({color:0x44bb44}))
      top.position.y = 0.33
      group.add(main, top)
      group.userData.isMeshGroup = true; return group
    }
    case 'flowerbed': {
      const group = new THREE.Group()
      const border = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.06, 0.5), new THREE.MeshStandardMaterial({color:0x8c5523, roughness:0.8, wireframe:true}))
      border.position.y = 0.03
      const soil = new THREE.Mesh(new THREE.BoxGeometry(0.44, 0.03, 0.44), new THREE.MeshStandardMaterial({color:0x6b3a1f}))
      soil.position.y = 0.015
      for (let i = 0; i < 6; i++) {
        const flower = new THREE.Mesh(new THREE.SphereGeometry(0.03, 6, 6), new THREE.MeshStandardMaterial({color:[0xff7bb5, 0xfffb4d, 0xff4d4f, 0xc17bff, 0xffaa4d, 0xff4d4f][i]}))
        flower.position.set(-0.12 + (i%3) * 0.14, 0.05, -0.12 + Math.floor(i/3) * 0.28)
        group.add(flower)
      }
      group.add(border, soil)
      group.userData.isMeshGroup = true; return group
    }
    case 'grasspatch': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.01, 0.6), new THREE.MeshStandardMaterial({color:0x44bb44}))
      for (let i = 0; i < 20; i++) {
        const blade = new THREE.Mesh(new THREE.CylinderGeometry(0.003, 0.005, 0.05 + Math.random() * 0.05, 4), new THREE.MeshStandardMaterial({color:0x33aa33}))
        blade.position.set(-0.25 + Math.random() * 0.5, 0.01, -0.25 + Math.random() * 0.5)
        group.add(blade)
      }
      group.add(base)
      group.userData.isMeshGroup = true; return group
    }
    case 'rock': {
      const group = new THREE.Group()
      const r1 = new THREE.Mesh(new THREE.IcosahedronGeometry(0.2, 0), new THREE.MeshStandardMaterial({color:0x888888, roughness:0.9}))
      r1.position.set(0, 0.05, 0)
      const r2 = new THREE.Mesh(new THREE.IcosahedronGeometry(0.14, 0), new THREE.MeshStandardMaterial({color:0x999999, roughness:0.9}))
      r2.position.set(0.15, 0.03, 0.08)
      const r3 = new THREE.Mesh(new THREE.IcosahedronGeometry(0.1, 0), new THREE.MeshStandardMaterial({color:0x777777, roughness:0.9}))
      r3.position.set(-0.12, 0.02, -0.06)
      group.add(r1, r2, r3)
      group.userData.isMeshGroup = true; return group
    }
    case 'fence': {
      const group = new THREE.Group()
      const rail1 = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 1, 8), new THREE.MeshStandardMaterial({color:0xa67c52}))
      rail1.position.set(0, 0.55, 0)
      const rail2 = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 1, 8), new THREE.MeshStandardMaterial({color:0xa67c52}))
      rail2.position.set(0, 0.25, 0)
      for (let i = -4; i <= 4; i++) {
        const post = new THREE.Mesh(new THREE.BoxGeometry(0.02, 0.7, 0.02), new THREE.MeshStandardMaterial({color:0x8c5523}))
        post.position.set(i * 0.12, 0.4, 0)
        group.add(post)
      }
      group.add(rail1, rail2)
      group.userData.isMeshGroup = true; return group
    }
    case 'gardenwall': {
      const group = new THREE.Group()
      for (let i = 0; i < 6; i++) {
        const brick = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.1, 0.15), new THREE.MeshStandardMaterial({color:0xcc9966, roughness:0.8}))
        brick.position.set(-0.5 + i * 0.2, 0.05, 0)
        group.add(brick)
      }
      for (let i = 0; i < 6; i++) {
        const brick = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.1, 0.15), new THREE.MeshStandardMaterial({color:0xcc9966, roughness:0.8}))
        brick.position.set(-0.4 + i * 0.2, 0.15, 0)
        group.add(brick)
      }
      group.userData.isMeshGroup = true; return group
    }
    case 'gate': {
      const group = new THREE.Group()
      const pillarL = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.8, 0.1), new THREE.MeshStandardMaterial({color:0x888888}))
      pillarL.position.set(-0.35, 0.4, 0)
      const pillarR = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.8, 0.1), new THREE.MeshStandardMaterial({color:0x888888}))
      pillarR.position.set(0.35, 0.4, 0)
      const gateL = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.6, 0.03), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      gateL.position.set(-0.15, 0.3, 0)
      const gateR = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.6, 0.03), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.6}))
      gateR.position.set(0.15, 0.3, 0)
      group.add(pillarL, pillarR, gateL, gateR)
      group.userData.isMeshGroup = true; return group
    }
    case 'lamppost': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.08, 0.1, 12), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 1.2, 12), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.7}))
      pole.position.y = 0.6
      const lamp = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.1, 0.15), new THREE.MeshStandardMaterial({color:0xffffcc, emissive:0xfffb4d, emissiveIntensity:0.6}))
      lamp.position.y = 1.2
      const cap = new THREE.Mesh(new THREE.CylinderGeometry(0.07, 0.1, 0.06, 4), new THREE.MeshStandardMaterial({color:0x333333}))
      cap.position.y = 1.28
      group.add(base, pole, lamp, cap)
      group.userData.isMeshGroup = true; return group
    }
    case 'pavilion': {
      const group = new THREE.Group()
      const floor = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.55, 0.08, 16), new THREE.MeshStandardMaterial({color:0x8c5523}))
      for (let i = 0; i < 6; i++) {
        const pillar = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 1, 8), new THREE.MeshStandardMaterial({color:0xcc0000}))
        pillar.position.set(Math.cos(i*Math.PI/3)*0.45, 0.5, Math.sin(i*Math.PI/3)*0.45)
        group.add(pillar)
      }
      const roof = new THREE.Mesh(new THREE.ConeGeometry(0.55, 0.3, 6), new THREE.MeshStandardMaterial({color:0x444444}))
      roof.position.y = 1.1
      const tip = new THREE.Mesh(new THREE.SphereGeometry(0.04, 6, 6), new THREE.MeshStandardMaterial({color:0xfffb4d}))
      tip.position.y = 1.35
      group.add(floor, roof, tip)
      group.userData.isMeshGroup = true; return group
    }
    case 'pond': {
      const group = new THREE.Group()
      const water = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.45, 0.04, 24), new THREE.MeshStandardMaterial({color:0x4dfff0, transparent:true, opacity:0.6, emissive:0x003333, emissiveIntensity:0.2}))
      water.position.y = 0.02
      const border = new THREE.Mesh(new THREE.TorusGeometry(0.42, 0.03, 8, 24), new THREE.MeshStandardMaterial({color:0x888888}))
      border.position.y = 0.04; border.rotation.x = Math.PI/2
      group.add(water, border)
      group.userData.isMeshGroup = true; return group
    }
    case 'sculpture': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.18, 0.15, 16), new THREE.MeshStandardMaterial({color:0x888888}))
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.12, 0.4, 12), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.5}))
      body.position.y = 0.25
      const top = new THREE.Mesh(new THREE.SphereGeometry(0.08, 12, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.6}))
      top.position.y = 0.5
      group.add(base, body, top)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 医疗设施 =====
    case 'hospitalbed': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.08, 1.8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.6}))
      const mattress = new THREE.Mesh(new THREE.BoxGeometry(0.76, 0.1, 1.76), new THREE.MeshStandardMaterial({color:0xe8f4f8}))
      mattress.position.y = 0.1
      const headboard = new THREE.Mesh(new THREE.BoxGeometry(0.82, 0.4, 0.06), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.5}))
      headboard.position.set(0, 0.3, 0.88)
      const footboard = new THREE.Mesh(new THREE.BoxGeometry(0.82, 0.25, 0.04), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.5}))
      footboard.position.set(0, 0.2, -0.88)
      for (let i = -1; i <= 1; i+=2) {
        for (let j = -1; j <= 1; j+=2) {
          const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.04, 12), new THREE.MeshStandardMaterial({color:0x333333}))
          wheel.rotation.z = Math.PI/2; wheel.position.set(i*0.36, -0.06, j*0.75); group.add(wheel)
        }
      }
      group.add(frame, mattress, headboard, footboard)
      group.userData.isMeshGroup = true; return group
    }
    case 'operatingtable': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.06, 2), new THREE.MeshStandardMaterial({color:0xe8f0f0}))
      top.position.y = 0.9
      const col = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.1, 0.85, 16), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      col.position.y = 0.42
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.08, 0.5), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.4}))
      base.position.y = 0.04
      group.add(top, col, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'mrimachine': {
      const group = new THREE.Group()
      const bore = new THREE.Mesh(new THREE.CylinderGeometry(0.35, 0.35, 1.2, 32, 1, true), new THREE.MeshStandardMaterial({color:0xdddddd, metalness:0.7}))
      bore.rotation.x = Math.PI/2
      const housing = new THREE.Mesh(new THREE.BoxGeometry(1, 0.9, 1.2), new THREE.MeshStandardMaterial({color:0xeeeeff, metalness:0.3}))
      housing.position.y = 0.5
      const base = new THREE.Mesh(new THREE.BoxGeometry(1.1, 0.1, 1.3), new THREE.MeshStandardMaterial({color:0x888888}))
      group.add(bore, housing, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'ctscanner': {
      const group = new THREE.Group()
      const ring = new THREE.Mesh(new THREE.TorusGeometry(0.4, 0.06, 16, 32), new THREE.MeshStandardMaterial({color:0xffffff, metalness:0.5}))
      ring.rotation.x = Math.PI/2; ring.position.y = 0.55
      const bed = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.06, 1.5), new THREE.MeshStandardMaterial({color:0xe0e8f0}))
      bed.position.y = 0.55; bed.position.z = 0.3
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.15, 0.8), new THREE.MeshStandardMaterial({color:0x999999}))
      group.add(ring, bed, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'ivstand': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 1.3, 12), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.5}))
      pole.position.y = 0.7
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.18, 0.2, 0.06, 16), new THREE.MeshStandardMaterial({color:0x888888}))
      base.position.y = 0.05
      const hook = new THREE.Mesh(new THREE.TorusGeometry(0.08, 0.02, 8, 16), new THREE.MeshStandardMaterial({color:0xbbbbbb, metalness:0.6}))
      hook.position.y = 1.35
      group.add(pole, base, hook)
      group.userData.isMeshGroup = true; return group
    }
    case 'wheelchair': {
      const group = new THREE.Group()
      const seat = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.06, 0.35), new THREE.MeshStandardMaterial({color:0x3366cc}))
      seat.position.y = 0.35
      const back = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.45, 0.04), new THREE.MeshStandardMaterial({color:0x3366cc}))
      back.position.set(0, 0.55, -0.16)
      for (let s = -1; s <= 1; s+=2) {
        const bigWheel = new THREE.Mesh(new THREE.TorusGeometry(0.16, 0.03, 8, 16), new THREE.MeshStandardMaterial({color:0x444444}))
        bigWheel.position.set(s*0.22, 0.16, -0.08); group.add(bigWheel)
        const smallWheel = new THREE.Mesh(new THREE.TorusGeometry(0.07, 0.02, 8, 12), new THREE.MeshStandardMaterial({color:0x444444}))
        smallWheel.position.set(s*0.16, 0.07, 0.18); group.add(smallWheel)
      }
      group.add(seat, back)
      group.userData.isMeshGroup = true; return group
    }
    case 'stretcher': {
      const group = new THREE.Group()
      const bed = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.06, 1.8), new THREE.MeshStandardMaterial({color:0xff8844}))
      bed.position.y = 0.6
      const leg1 = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.55, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      leg1.position.set(0.2, 0.28, 0.7)
      const leg2 = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.55, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      leg2.position.set(-0.2, 0.28, 0.7)
      const leg3 = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.55, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      leg3.position.set(0.2, 0.28, -0.7)
      const leg4 = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.55, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      leg4.position.set(-0.2, 0.28, -0.7)
      group.add(bed, leg1, leg2, leg3, leg4)
      group.userData.isMeshGroup = true; return group
    }
    case 'medicinecabinet': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.8, 0.35), new THREE.MeshStandardMaterial({color:0xeeeeee, metalness:0.3}))
      body.position.y = 0.4
      for (let i = 0; i < 3; i++) {
        const shelf = new THREE.Mesh(new THREE.BoxGeometry(0.46, 0.02, 0.33), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.4}))
        shelf.position.y = 0.15 + i * 0.24; group.add(shelf)
      }
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.24, 0.72, 0.02), new THREE.MeshStandardMaterial({color:0xffffff, metalness:0.2, transparent:true, opacity:0.7}))
      door.position.set(0.13, 0.4, 0.18); door.visible = false
      group.add(body, door)
      group.userData.isMeshGroup = true; return group
    }
    case 'examinationlamp': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.14, 0.08, 16), new THREE.MeshStandardMaterial({color:0x888888}))
      const arm1 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.5, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.5}))
      arm1.position.y = 0.3; arm1.rotation.z = 0.3
      const arm2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.4, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.5}))
      arm2.position.set(0.15, 0.6, 0); arm2.rotation.z = -0.5
      const head = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.1, 0.1, 16), new THREE.MeshStandardMaterial({color:0xffffff, emissive:0xffffcc, emissiveIntensity:0.3}))
      head.position.set(0.1, 0.8, 0)
      group.add(base, arm1, arm2, head)
      group.userData.isMeshGroup = true; return group
    }
    case 'sinkstation': {
      const group = new THREE.Group()
      const counter = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.08, 0.4), new THREE.MeshStandardMaterial({color:0xcccccc}))
      counter.position.y = 0.7
      const basin = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.13, 0.1, 24, 1, true), new THREE.MeshStandardMaterial({color:0xeeeeff, metalness:0.5}))
      basin.position.y = 0.65
      const faucet = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.25, 12), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      faucet.position.set(0, 0.85, -0.12)
      const spout = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.1, 8), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.7}))
      spout.rotation.x = Math.PI/2; spout.position.set(0, 0.95, -0.02)
      group.add(counter, basin, faucet, spout)
      group.userData.isMeshGroup = true; return group
    }
    case 'hospitalmonitor': {
      const group = new THREE.Group()
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.28, 0.2, 0.04), new THREE.MeshStandardMaterial({color:0x111133}))
      screen.position.y = 0.85
      const bezel = new THREE.Mesh(new THREE.BoxGeometry(0.32, 0.24, 0.06), new THREE.MeshStandardMaterial({color:0xcccccc}))
      bezel.position.y = 0.85
      const stand = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.7, 12), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.4}))
      stand.position.y = 0.4
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.14, 0.06, 16), new THREE.MeshStandardMaterial({color:0x888888}))
      base.position.y = 0.08
      group.add(screen, bezel, stand, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'defibrillator': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.2, 0.25), new THREE.MeshStandardMaterial({color:0xff6633}))
      body.position.y = 0.15
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.1, 0.01), new THREE.MeshStandardMaterial({color:0x003322}))
      screen.position.set(0, 0.22, 0.13)
      for (let i = -1; i <= 1; i+=2) {
        const paddle = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.04, 0.1), new THREE.MeshStandardMaterial({color:0x888888}))
        paddle.position.set(i*0.08, 0.08, 0.15); group.add(paddle)
      }
      group.add(body, screen)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 农业设施 =====
    case 'greenhouse': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.1, 1.2), new THREE.MeshStandardMaterial({color:0x888888}))
      const arch = new THREE.Mesh(new THREE.CylinderGeometry(0.4, 0.4, 1.2, 16, 1, false, 0, Math.PI), new THREE.MeshStandardMaterial({color:0xaaddff, metalness:0.1, transparent:true, opacity:0.5}))
      arch.rotation.z = Math.PI/2; arch.rotation.y = Math.PI/2; arch.position.y = 0.35
      const ridge = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1.2, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      ridge.position.y = 0.75; ridge.rotation.x = Math.PI/2
      group.add(base, arch, ridge)
      group.userData.isMeshGroup = true; return group
    }
    case 'tractor': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.2, 0.55), new THREE.MeshStandardMaterial({color:0xcc3333}))
      body.position.y = 0.35
      const cabin = new THREE.Mesh(new THREE.BoxGeometry(0.32, 0.25, 0.3), new THREE.MeshStandardMaterial({color:0xaaddff, metalness:0.1, transparent:true, opacity:0.6}))
      cabin.position.set(0, 0.55, -0.05)
      for (let s = -1; s <= 1; s+=2) {
        const bigWheel = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.08, 16), new THREE.MeshStandardMaterial({color:0x333333}))
        bigWheel.rotation.z = Math.PI/2; bigWheel.position.set(s*0.18, 0.15, 0.15); group.add(bigWheel)
        const smallWheel = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.06, 16), new THREE.MeshStandardMaterial({color:0x333333}))
        smallWheel.rotation.z = Math.PI/2; smallWheel.position.set(s*0.18, 0.1, -0.22); group.add(smallWheel)
      }
      group.add(body, cabin)
      group.userData.isMeshGroup = true; return group
    }
    case 'harvester': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.25, 0.5), new THREE.MeshStandardMaterial({color:0x33aa33}))
      body.position.y = 0.3
      const header = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.1, 0.2), new THREE.MeshStandardMaterial({color:0xffaa00}))
      header.position.set(0, 0.25, 0.35)
      for (let s = -1; s <= 1; s+=2) {
        const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.14, 0.06, 16), new THREE.MeshStandardMaterial({color:0x333333}))
        wheel.rotation.z = Math.PI/2; wheel.position.set(s*0.18, 0.14, 0.05); group.add(wheel)
      }
      group.add(body, header)
      group.userData.isMeshGroup = true; return group
    }
    case 'irrigationpivot': {
      const group = new THREE.Group()
      const center = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.1, 0.5, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      center.position.y = 0.25
      const arm = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 1.5, 12), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.5}))
      arm.rotation.z = Math.PI/2; arm.position.set(0.7, 0.55, 0)
      for (let i = 0; i < 3; i++) {
        const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.45, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
        leg.position.set(0.2 + i*0.5, 0.25, 0); group.add(leg)
        const wheel = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.04, 12), new THREE.MeshStandardMaterial({color:0x333333}))
        wheel.rotation.z = Math.PI/2; wheel.position.set(0.2 + i*0.5, 0.06, 0); group.add(wheel)
      }
      group.add(center, arm)
      group.userData.isMeshGroup = true; return group
    }
    case 'grainsilo': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.25, 0.3, 0.8, 24), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.6}))
      body.position.y = 0.4
      const roof = new THREE.Mesh(new THREE.ConeGeometry(0.32, 0.2, 24), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.5}))
      roof.position.y = 0.9
      group.add(body, roof)
      group.userData.isMeshGroup = true; return group
    }
    case 'sprayer': {
      const group = new THREE.Group()
      const tank = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.4, 16), new THREE.MeshStandardMaterial({color:0xff8800}))
      tank.position.y = 0.25
      const boom = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 1.2, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      boom.rotation.z = Math.PI/2; boom.position.y = 0.3
      for (let i = 0; i < 4; i++) {
        const nozzle = new THREE.Mesh(new THREE.ConeGeometry(0.03, 0.06, 8), new THREE.MeshStandardMaterial({color:0xffcc00}))
        nozzle.position.set(-0.5 + i*0.33, 0.2, 0); group.add(nozzle)
      }
      group.add(tank, boom)
      group.userData.isMeshGroup = true; return group
    }
    case 'plantingbed': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.2, 0.3), new THREE.MeshStandardMaterial({color:0x8B4513}))
      frame.position.y = 0.1
      const soil = new THREE.Mesh(new THREE.BoxGeometry(0.72, 0.06, 0.22), new THREE.MeshStandardMaterial({color:0x4a3000}))
      soil.position.y = 0.2
      group.add(frame, soil)
      group.userData.isMeshGroup = true; return group
    }
    case 'windbreak': {
      const group = new THREE.Group()
      for (let i = 0; i < 5; i++) {
        const post = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.03, 0.6, 8), new THREE.MeshStandardMaterial({color:0x666666}))
        post.position.set(-0.6 + i*0.3, 0.3, 0); group.add(post)
      }
      const meshPanel = new THREE.Mesh(new THREE.PlaneGeometry(1, 0.45), new THREE.MeshStandardMaterial({color:0x448844, transparent:true, opacity:0.4, side:THREE.DoubleSide}))
      meshPanel.position.y = 0.45
      group.add(meshPanel)
      group.userData.isMeshGroup = true; return group
    }
    case 'weatherstation': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.7, 12), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.5}))
      pole.position.y = 0.35
      const sensorBox = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.1, 0.1), new THREE.MeshStandardMaterial({color:0xeeeeee}))
      sensorBox.position.y = 0.72
      const anemometer = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.2, 8), new THREE.MeshStandardMaterial({color:0xff6600}))
      anemometer.rotation.z = Math.PI/2; anemometer.position.y = 0.8
      const cup = new THREE.Mesh(new THREE.SphereGeometry(0.04, 8, 8), new THREE.MeshStandardMaterial({color:0xff6600}))
      cup.position.set(0.12, 0.8, 0)
      group.add(pole, sensorBox, anemometer, cup)
      group.userData.isMeshGroup = true; return group
    }
    case 'chickencoop': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.3, 0.8), new THREE.MeshStandardMaterial({color:0xcc8844}))
      base.position.y = 0.15
      for (let side of [0.3, -0.3]) {
        const roofPlane = new THREE.Mesh(new THREE.BoxGeometry(0.45, 0.04, 0.85), new THREE.MeshStandardMaterial({color:0x884422}))
        roofPlane.position.set(side, 0.3, 0); roofPlane.rotation.z = side > 0 ? 0.4 : -0.4; group.add(roofPlane)
      }
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.2, 0.02), new THREE.MeshStandardMaterial({color:0x663300}))
      door.position.set(0, 0.15, 0.41)
      group.add(base, door)
      group.userData.isMeshGroup = true; return group
    }
    case 'watertrough': {
      const group = new THREE.Group()
      const trough = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.15, 0.25), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      trough.position.y = 0.12
      const water = new THREE.Mesh(new THREE.BoxGeometry(0.56, 0.02, 0.21), new THREE.MeshStandardMaterial({color:0x4488ff, metalness:0.1, transparent:true, opacity:0.7}))
      water.position.y = 0.18
      group.add(trough, water)
      group.userData.isMeshGroup = true; return group
    }
    case 'seedlingtray': {
      const group = new THREE.Group()
      const tray = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.04, 0.3), new THREE.MeshStandardMaterial({color:0x222222}))
      tray.position.y = 0.02
      for (let x = -1; x <= 1; x++) {
        for (let z = -1; z <= 1; z++) {
          const cell = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.03, 0.04, 8), new THREE.MeshStandardMaterial({color:0x3a2010}))
          cell.position.set(x*0.1, 0.06, z*0.08); group.add(cell)
          const sprout = new THREE.Mesh(new THREE.ConeGeometry(0.015, 0.04, 6), new THREE.MeshStandardMaterial({color:0x33aa33}))
          sprout.position.set(x*0.1, 0.1, z*0.08); group.add(sprout)
        }
      }
      group.add(tray)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 水利水务 =====
    case 'watertower': {
      const group = new THREE.Group()
      for (let i = 0; i < 4; i++) {
        const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.05, 0.7, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
        const a = i * Math.PI/2; leg.position.set(Math.cos(a)*0.2, 0.35, Math.sin(a)*0.2); group.add(leg)
      }
      const tank = new THREE.Mesh(new THREE.CylinderGeometry(0.28, 0.28, 0.35, 24), new THREE.MeshStandardMaterial({color:0x4488cc, metalness:0.5}))
      tank.position.y = 0.85
      const top = new THREE.Mesh(new THREE.CylinderGeometry(0.28, 0.22, 0.08, 24), new THREE.MeshStandardMaterial({color:0x3377bb, metalness:0.5}))
      top.position.y = 1.05
      group.add(tank, top)
      group.userData.isMeshGroup = true; return group
    }
    case 'treatmentbasin': {
      const group = new THREE.Group()
      const basin = new THREE.Mesh(new THREE.BoxGeometry(1.2, 0.15, 0.8), new THREE.MeshStandardMaterial({color:0x778899}))
      basin.position.y = 0.1
      const water = new THREE.Mesh(new THREE.BoxGeometry(1.1, 0.02, 0.7), new THREE.MeshStandardMaterial({color:0x4488cc, transparent:true, opacity:0.6}))
      water.position.y = 0.16
      group.add(basin, water)
      group.userData.isMeshGroup = true; return group
    }
    case 'pipegallery': {
      const group = new THREE.Group()
      const tunnel = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.4, 1.2), new THREE.MeshStandardMaterial({color:0x888888}))
      tunnel.position.y = 0.3
      for (let i = 0; i < 3; i++) {
        const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 1.1, 16), new THREE.MeshStandardMaterial({color:0xcc8844, metalness:0.4}))
        pipe.rotation.x = Math.PI/2; pipe.position.set(-0.12 + i*0.12, 0.2, 0); group.add(pipe)
      }
      group.add(tunnel)
      group.userData.isMeshGroup = true; return group
    }
    case 'sluicegate': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.6, 0.1), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      frame.position.set(0.2, 0.3, 0)
      const frame2 = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.6, 0.1), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      frame2.position.set(-0.2, 0.3, 0)
      const gate = new THREE.Mesh(new THREE.BoxGeometry(0.45, 0.5, 0.05), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.6}))
      gate.position.y = 0.25
      group.add(frame, frame2, gate)
      group.userData.isMeshGroup = true; return group
    }
    case 'pumpstation': {
      const group = new THREE.Group()
      const building = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.35, 0.5), new THREE.MeshStandardMaterial({color:0x99aabb}))
      building.position.y = 0.18
      const roof = new THREE.Mesh(new THREE.ConeGeometry(0.4, 0.15, 4), new THREE.MeshStandardMaterial({color:0x667788}))
      roof.rotation.y = Math.PI/4; roof.position.y = 0.42
      const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.5, 12), new THREE.MeshStandardMaterial({color:0xcc8844, metalness:0.4}))
      pipe.rotation.x = Math.PI/2; pipe.position.set(0.3, 0.1, 0)
      group.add(building, roof, pipe)
      group.userData.isMeshGroup = true; return group
    }
    case 'sedimentationtank': {
      const group = new THREE.Group()
      const tank = new THREE.Mesh(new THREE.BoxGeometry(1, 0.2, 0.6), new THREE.MeshStandardMaterial({color:0x8899aa}))
      tank.position.y = 0.15
      const water = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.04, 0.5), new THREE.MeshStandardMaterial({color:0x557788, transparent:true, opacity:0.6}))
      water.position.y = 0.24
      const baffle = new THREE.Mesh(new THREE.BoxGeometry(0.04, 0.18, 0.55), new THREE.MeshStandardMaterial({color:0x777777}))
      baffle.position.set(-0.2, 0.12, 0)
      group.add(tank, water, baffle)
      group.userData.isMeshGroup = true; return group
    }
    case 'aerationtank': {
      const group = new THREE.Group()
      const tank = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.25, 0.6), new THREE.MeshStandardMaterial({color:0x778899}))
      tank.position.y = 0.15
      const water = new THREE.Mesh(new THREE.BoxGeometry(0.7, 0.03, 0.5), new THREE.MeshStandardMaterial({color:0x336688, transparent:true, opacity:0.6}))
      water.position.y = 0.26
      for (let i = 0; i < 4; i++) {
        const bubble = new THREE.Mesh(new THREE.SphereGeometry(0.03, 8, 8), new THREE.MeshStandardMaterial({color:0xffffff, transparent:true, opacity:0.5}))
        bubble.position.set(-0.25 + i*0.16, 0.3 + Math.random()*0.05, (Math.random()-0.5)*0.3); group.add(bubble)
      }
      group.add(tank, water)
      group.userData.isMeshGroup = true; return group
    }
    case 'overflowweir': {
      const group = new THREE.Group()
      const wall = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.3, 0.08), new THREE.MeshStandardMaterial({color:0x777777}))
      wall.position.y = 0.2
      const channel = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.1, 0.15), new THREE.MeshStandardMaterial({color:0x667788}))
      const cascade = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.08, 0.2), new THREE.MeshStandardMaterial({color:0x5577aa}))
      cascade.position.set(0, -0.05, 0.12)
      group.add(wall, channel, cascade)
      group.userData.isMeshGroup = true; return group
    }
    case 'intaketower': {
      const group = new THREE.Group()
      const tower = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.2, 0.8, 16), new THREE.MeshStandardMaterial({color:0x888899, metalness:0.3}))
      tower.position.y = 0.4
      const cap = new THREE.Mesh(new THREE.ConeGeometry(0.2, 0.15, 16), new THREE.MeshStandardMaterial({color:0x778899, metalness:0.3}))
      cap.position.y = 0.87
      for (let i = 0; i < 3; i++) {
        const window = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.08, 0.02), new THREE.MeshStandardMaterial({color:0x333344}))
        window.position.set(0, 0.3 + i*0.2, 0.18); group.add(window)
      }
      group.add(tower, cap)
      group.userData.isMeshGroup = true; return group
    }
    case 'flowchannel': {
      const group = new THREE.Group()
      const leftWall = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.2, 1), new THREE.MeshStandardMaterial({color:0x777777}))
      leftWall.position.set(-0.25, 0.1, 0)
      const rightWall = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.2, 1), new THREE.MeshStandardMaterial({color:0x777777}))
      rightWall.position.set(0.25, 0.1, 0)
      const bottom = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.04, 1), new THREE.MeshStandardMaterial({color:0x667788}))
      const water = new THREE.Mesh(new THREE.BoxGeometry(0.42, 0.02, 0.9), new THREE.MeshStandardMaterial({color:0x4488cc, transparent:true, opacity:0.5}))
      water.position.y = 0.06
      group.add(leftWall, rightWall, bottom, water)
      group.userData.isMeshGroup = true; return group
    }
    case 'manholecover': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.CylinderGeometry(0.22, 0.25, 0.06, 24), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.5}))
      const cover = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 0.03, 24), new THREE.MeshStandardMaterial({color:0x555555, metalness:0.6}))
      cover.position.y = 0.03
      for (let i = 0; i < 2; i++) {
        const rib = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.015, 0.03), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.5}))
        rib.position.y = 0.04; rib.rotation.y = i * Math.PI/2; group.add(rib)
      }
      group.add(frame, cover)
      group.userData.isMeshGroup = true; return group
    }
    case 'drainpipe': {
      const group = new THREE.Group()
      const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.8, 16, 1, true), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      pipe.position.y = 0.4
      const flange = new THREE.Mesh(new THREE.TorusGeometry(0.09, 0.02, 8, 16), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.5}))
      flange.rotation.x = Math.PI/2; flange.position.y = 0.1
      const flange2 = new THREE.Mesh(new THREE.TorusGeometry(0.09, 0.02, 8, 16), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.5}))
      flange2.rotation.x = Math.PI/2; flange2.position.y = 0.7
      group.add(pipe, flange, flange2)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 运动休闲 =====
    case 'basketballhoop': {
      const group = new THREE.Group()
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.05, 0.9, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      pole.position.y = 0.45
      const backboard = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.28, 0.04), new THREE.MeshStandardMaterial({color:0xffffff, metalness:0.1, transparent:true, opacity:0.85}))
      backboard.position.y = 0.85
      const rim = new THREE.Mesh(new THREE.TorusGeometry(0.12, 0.02, 8, 16), new THREE.MeshStandardMaterial({color:0xff4400}))
      rim.rotation.x = Math.PI/2; rim.position.set(0, 0.75, 0.12)
      group.add(pole, backboard, rim)
      group.userData.isMeshGroup = true; return group
    }
    case 'soccergoal': {
      const group = new THREE.Group()
      const posts = [
        [0.3, 0.3, -0.2], [-0.3, 0.3, -0.2], [0, 0.55, -0.2]
      ]
      for (const [x, y, z] of posts) {
        const post = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.55, 8), new THREE.MeshStandardMaterial({color:0xffffff}))
        if (x === 0) post.rotation.x = Math.PI/2
        post.position.set(x, y, z); group.add(post)
      }
      const net = new THREE.Mesh(new THREE.PlaneGeometry(0.6, 0.3), new THREE.MeshStandardMaterial({color:0xffffff, transparent:true, opacity:0.2, side:THREE.DoubleSide}))
      net.position.set(0, 0.4, -0.15)
      group.add(net)
      group.userData.isMeshGroup = true; return group
    }
    case 'tenniscourt': {
      const group = new THREE.Group()
      const court = new THREE.Mesh(new THREE.PlaneGeometry(1.2, 0.7), new THREE.MeshStandardMaterial({color:0x3366aa, side:THREE.DoubleSide}))
      court.rotation.x = -Math.PI/2; court.position.y = 0.01
      const net = new THREE.Mesh(new THREE.PlaneGeometry(1, 0.12), new THREE.MeshStandardMaterial({color:0x333333, side:THREE.DoubleSide}))
      net.position.y = 0.06
      const netTop = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.01, 1, 8), new THREE.MeshStandardMaterial({color:0xffffff}))
      netTop.rotation.z = Math.PI/2; netTop.position.y = 0.12
      group.add(court, net, netTop)
      group.userData.isMeshGroup = true; return group
    }
    case 'swimmingpool': {
      const group = new THREE.Group()
      const basin = new THREE.Mesh(new THREE.BoxGeometry(1, 0.15, 0.5), new THREE.MeshStandardMaterial({color:0x8899bb}))
      basin.position.y = 0.08
      const water = new THREE.Mesh(new THREE.BoxGeometry(0.9, 0.02, 0.4), new THREE.MeshStandardMaterial({color:0x4499ff, metalness:0.2, transparent:true, opacity:0.7, emissive:0x1155aa, emissiveIntensity:0.2}))
      water.position.y = 0.14
      group.add(basin, water)
      group.userData.isMeshGroup = true; return group
    }
    case 'slide': {
      const group = new THREE.Group()
      const platform = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.08, 0.5), new THREE.MeshStandardMaterial({color:0xff8844}))
      platform.position.y = 0.6
      for (let i = 0; i < 4; i++) {
        const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.56, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
        leg.position.set((i%2===0?-1:1)*0.2, 0.28, (i<2?-1:1)*0.2); group.add(leg)
      }
      const ramp = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.04, 0.5), new THREE.MeshStandardMaterial({color:0xffcc44}))
      ramp.position.set(0.3, 0.3, 0); ramp.rotation.z = 0.5
      group.add(platform, ramp)
      group.userData.isMeshGroup = true; return group
    }
    case 'swingset': {
      const group = new THREE.Group()
      for (let s = -1; s <= 1; s+=2) {
        const leg1 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.7, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
        leg1.position.set(s*0.35, 0.35, -0.15); leg1.rotation.z = s*0.2; group.add(leg1)
        const leg2 = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.7, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
        leg2.position.set(s*0.35, 0.35, 0.15); leg2.rotation.z = s*0.2; group.add(leg2)
      }
      const bar = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.75, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      bar.rotation.z = Math.PI/2; bar.position.y = 0.65
      const swing1 = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.04, 0.08), new THREE.MeshStandardMaterial({color:0xff8844}))
      swing1.position.set(-0.12, 0.4, 0)
      const swing2 = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.04, 0.08), new THREE.MeshStandardMaterial({color:0xff8844}))
      swing2.position.set(0.12, 0.4, 0)
      group.add(bar, swing1, swing2)
      group.userData.isMeshGroup = true; return group
    }
    case 'exercisebike': {
      const group = new THREE.Group()
      const frame = new THREE.Mesh(new THREE.BoxGeometry(0.08, 0.3, 0.5), new THREE.MeshStandardMaterial({color:0x444444, metalness:0.6}))
      frame.position.y = 0.2
      const seat = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.06, 16), new THREE.MeshStandardMaterial({color:0x222222}))
      seat.position.set(0, 0.35, 0.15)
      const handlebar = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.25, 8), new THREE.MeshStandardMaterial({color:0x666666, metalness:0.5}))
      handlebar.rotation.x = Math.PI/2; handlebar.position.set(0, 0.32, -0.15)
      const wheel = new THREE.Mesh(new THREE.TorusGeometry(0.12, 0.02, 8, 16), new THREE.MeshStandardMaterial({color:0x888888}))
      wheel.position.y = 0.2
      group.add(frame, seat, handlebar, wheel)
      group.userData.isMeshGroup = true; return group
    }
    case 'treadmill': {
      const group = new THREE.Group()
      const deck = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.06, 0.7), new THREE.MeshStandardMaterial({color:0x333333}))
      const belt = new THREE.Mesh(new THREE.BoxGeometry(0.26, 0.02, 0.65), new THREE.MeshStandardMaterial({color:0x444444}))
      belt.position.y = 0.04
      const console = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.2, 0.08), new THREE.MeshStandardMaterial({color:0x888888}))
      console.position.set(0, 0.18, -0.35)
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.18, 0.1, 0.01), new THREE.MeshStandardMaterial({color:0x112244}))
      screen.position.set(0, 0.2, -0.4)
      group.add(deck, belt, console, screen)
      group.userData.isMeshGroup = true; return group
    }
    case 'stadium': {
      const group = new THREE.Group()
      const field = new THREE.Mesh(new THREE.PlaneGeometry(1.2, 0.8), new THREE.MeshStandardMaterial({color:0x33aa55, side:THREE.DoubleSide}))
      field.rotation.x = -Math.PI/2; field.position.y = 0.01
      const track = new THREE.Mesh(new THREE.TorusGeometry(0.4, 0.06, 8, 32), new THREE.MeshStandardMaterial({color:0xcc6633}))
      track.rotation.x = Math.PI/2; track.position.y = 0.02
      group.add(field, track)
      group.userData.isMeshGroup = true; return group
    }
    case 'scoreboard': {
      const group = new THREE.Group()
      const board = new THREE.Mesh(new THREE.BoxGeometry(0.6, 0.35, 0.06), new THREE.MeshStandardMaterial({color:0x111122}))
      board.position.y = 0.55
      const stand = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.06, 0.4, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
      stand.position.y = 0.2
      const display1 = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.12, 0.01), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.4}))
      display1.position.set(-0.12, 0.6, 0.04)
      const display2 = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.12, 0.01), new THREE.MeshStandardMaterial({color:0xff0000, emissive:0xff0000, emissiveIntensity:0.4}))
      display2.position.set(0.12, 0.6, 0.04)
      group.add(board, stand, display1, display2)
      group.userData.isMeshGroup = true; return group
    }
    case 'bleacher': {
      const group = new THREE.Group()
      for (let i = 0; i < 4; i++) {
        const step = new THREE.Mesh(new THREE.BoxGeometry(0.8, 0.12, 0.25), new THREE.MeshStandardMaterial({color:0x888899}))
        step.position.set(0, 0.06 + i*0.15, i*0.2); group.add(step)
      }
      group.userData.isMeshGroup = true; return group
    }
    case 'pingpongtable': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.04, 0.8), new THREE.MeshStandardMaterial({color:0x336644}))
      top.position.y = 0.55
      for (let s = -1; s <= 1; s+=2) {
        for (let t = -1; t <= 1; t+=2) {
          const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.5, 8), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.4}))
          leg.position.set(s*0.2, 0.25, t*0.3); group.add(leg)
        }
      }
      const net = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.08, 0.02), new THREE.MeshStandardMaterial({color:0xffffff}))
      net.position.y = 0.6
      group.add(top, net)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 实验室 =====
    case 'labbench': {
      const group = new THREE.Group()
      const top = new THREE.Mesh(new THREE.BoxGeometry(0.7, 0.06, 0.4), new THREE.MeshStandardMaterial({color:0xeeeeee}))
      top.position.y = 0.55
      for (let s = -1; s <= 1; s+=2) {
        for (let t = -1; t <= 1; t+=2) {
          const leg = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.5, 8), new THREE.MeshStandardMaterial({color:0x999999, metalness:0.4}))
          leg.position.set(s*0.28, 0.25, t*0.15); group.add(leg)
        }
      }
      group.add(top)
      group.userData.isMeshGroup = true; return group
    }
    case 'fumehood': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.35, 0.3), new THREE.MeshStandardMaterial({color:0xeeeeee}))
      base.position.y = 0.18
      const hood = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.2, 0.28), new THREE.MeshStandardMaterial({color:0xddddff, metalness:0.2, transparent:true, opacity:0.6}))
      hood.position.y = 0.42
      const duct = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.15, 12), new THREE.MeshStandardMaterial({color:0xaaaaaa, metalness:0.4}))
      duct.position.y = 0.58
      group.add(base, hood, duct)
      group.userData.isMeshGroup = true; return group
    }
    case 'microscope': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.05, 0.18), new THREE.MeshStandardMaterial({color:0x666666}))
      const arm = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.25, 12), new THREE.MeshStandardMaterial({color:0x888888, metalness:0.5}))
      arm.position.y = 0.16; arm.position.x = -0.05
      const tube = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.15, 12), new THREE.MeshStandardMaterial({color:0x444444}))
      tube.position.y = 0.3; tube.rotation.x = 0.2
      const lens = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.03, 0.06, 12), new THREE.MeshStandardMaterial({color:0x222222}))
      lens.position.set(-0.02, 0.23, 0.08)
      group.add(base, arm, tube, lens)
      group.userData.isMeshGroup = true; return group
    }
    case 'centrifuge': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.13, 0.2, 24), new THREE.MeshStandardMaterial({color:0xeeeeee, metalness:0.3}))
      body.position.y = 0.12
      const lid = new THREE.Mesh(new THREE.CylinderGeometry(0.12, 0.12, 0.04, 24), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.4}))
      lid.position.y = 0.24
      const knob = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.04, 8), new THREE.MeshStandardMaterial({color:0x333333}))
      knob.position.y = 0.28
      group.add(body, lid, knob)
      group.userData.isMeshGroup = true; return group
    }
    case 'incubator': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.35, 0.3), new THREE.MeshStandardMaterial({color:0xeeeeee, metalness:0.3}))
      body.position.y = 0.18
      const door = new THREE.Mesh(new THREE.BoxGeometry(0.25, 0.28, 0.02), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.2, transparent:true, opacity:0.7}))
      door.position.z = 0.16
      const display = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.06, 0.01), new THREE.MeshStandardMaterial({color:0x003322, emissive:0x004422, emissiveIntensity:0.3}))
      display.position.set(0, 0.34, 0.16)
      group.add(body, door, display)
      group.userData.isMeshGroup = true; return group
    }
    case 'autoclave': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.15, 0.4, 24), new THREE.MeshStandardMaterial({color:0xdddddd, metalness:0.5}))
      body.position.y = 0.22
      const lid = new THREE.Mesh(new THREE.CylinderGeometry(0.16, 0.14, 0.06, 24), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.6}))
      lid.position.y = 0.45
      const gauge = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.03, 16), new THREE.MeshStandardMaterial({color:0xeeeeee}))
      gauge.position.set(0, 0.48, 0.12)
      group.add(body, lid, gauge)
      group.userData.isMeshGroup = true; return group
    }
    case 'reagentshelf': {
      const group = new THREE.Group()
      for (let i = 0; i < 4; i++) {
        const shelf = new THREE.Mesh(new THREE.BoxGeometry(0.45, 0.02, 0.2), new THREE.MeshStandardMaterial({color:0xccbbaa}))
        shelf.position.y = 0.08 + i*0.12; group.add(shelf)
      }
      for (let s = -1; s <= 1; s+=2) {
        const side = new THREE.Mesh(new THREE.BoxGeometry(0.02, 0.4, 0.2), new THREE.MeshStandardMaterial({color:0xccbbaa}))
        side.position.set(s*0.22, 0.2, 0); group.add(side)
      }
      for (let j = 0; j < 12; j++) {
        const bottle = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.012, 0.05, 8), new THREE.MeshStandardMaterial({color:0x88ccff, transparent:true, opacity:0.6}))
        bottle.position.set(-0.15 + (j%4)*0.1, 0.13 + Math.floor(j/4)*0.12, 0); group.add(bottle)
      }
      group.userData.isMeshGroup = true; return group
    }
    case 'glovebox': {
      const group = new THREE.Group()
      const box = new THREE.Mesh(new THREE.BoxGeometry(0.4, 0.35, 0.3), new THREE.MeshStandardMaterial({color:0xeeffff, metalness:0.1, transparent:true, opacity:0.6}))
      box.position.y = 0.2
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.42, 0.08, 0.32), new THREE.MeshStandardMaterial({color:0x888888}))
      for (let s = -1; s <= 1; s+=2) {
        const glove = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 0.15, 8), new THREE.MeshStandardMaterial({color:0xffcc88}))
        glove.position.set(s*0.15, 0.15, -0.22); group.add(glove)
      }
      group.add(box, base)
      group.userData.isMeshGroup = true; return group
    }
    case 'spectrometer': {
      const group = new THREE.Group()
      const body = new THREE.Mesh(new THREE.BoxGeometry(0.35, 0.25, 0.3), new THREE.MeshStandardMaterial({color:0xeeeeff, metalness:0.2}))
      body.position.y = 0.18
      const detector = new THREE.Mesh(new THREE.CylinderGeometry(0.06, 0.06, 0.15, 12), new THREE.MeshStandardMaterial({color:0x444444}))
      detector.position.set(0.12, 0.32, 0)
      const screen = new THREE.Mesh(new THREE.BoxGeometry(0.15, 0.1, 0.01), new THREE.MeshStandardMaterial({color:0x112233}))
      screen.position.set(-0.08, 0.22, 0.16)
      group.add(body, detector, screen)
      group.userData.isMeshGroup = true; return group
    }
    case 'labscale': {
      const group = new THREE.Group()
      const base = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.06, 0.18), new THREE.MeshStandardMaterial({color:0xeeeeee}))
      base.position.y = 0.07
      const pan = new THREE.Mesh(new THREE.CylinderGeometry(0.08, 0.08, 0.02, 24), new THREE.MeshStandardMaterial({color:0xcccccc, metalness:0.6}))
      pan.position.y = 0.12
      const display = new THREE.Mesh(new THREE.BoxGeometry(0.1, 0.06, 0.03), new THREE.MeshStandardMaterial({color:0x003322, emissive:0x003322, emissiveIntensity:0.3}))
      display.position.set(0, 0.12, -0.1)
      const shield = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.1, 0.15, 16, 1, true), new THREE.MeshStandardMaterial({color:0xeeeeff, metalness:0.1, transparent:true, opacity:0.4}))
      shield.position.y = 0.2
      group.add(base, pan, display, shield)
      group.userData.isMeshGroup = true; return group
    }

    // ===== 视觉特效 =====
    case 'scanLine': {
      const g = new THREE.Group()
      // 中心柱
      const pole = new THREE.Mesh(new THREE.CylinderGeometry(0.04,0.04,0.8,8), new THREE.MeshStandardMaterial({color:0x4488cc,emissive:0x224466,emissiveIntensity:0.4}))
      pole.position.y = 0.4; g.add(pole)
      // 旋转扫描臂
      const arm = new THREE.Group()
      const armBody = new THREE.Mesh(new THREE.BoxGeometry(0.6,0.012,0.02), new THREE.MeshStandardMaterial({color:0x00ffff,emissive:0x00ffff,emissiveIntensity:0.8}))
      armBody.position.x = 0.3; arm.add(armBody)
      arm.position.y = 0.75; arm.name = 'scanArm'
      g.add(arm)
      // 发光底座环
      const baseRing = new THREE.Mesh(new THREE.TorusGeometry(0.2,0.02,8,32), new THREE.MeshStandardMaterial({color:0x00ccff,emissive:0x00ccff,emissiveIntensity:0.6,transparent:true,opacity:0.7}))
      baseRing.rotation.x = -Math.PI/2; baseRing.position.y = 0.05; g.add(baseRing)
      // 扫描面
      const scanPlane = new THREE.Mesh(new THREE.CylinderGeometry(0.55,0.55,0.005,64), new THREE.MeshStandardMaterial({color:0x00ffff,emissive:0x00ffff,emissiveIntensity:0.3,transparent:true,opacity:0.25,depthWrite:false}))
      scanPlane.position.y = 0.01; g.add(scanPlane)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'flyLine': {
      const g = new THREE.Group()
      // 起点球
      const startSphere = new THREE.Mesh(new THREE.SphereGeometry(0.08,16,16), new THREE.MeshStandardMaterial({color:0x00ff88,emissive:0x00ff88,emissiveIntensity:0.8}))
      startSphere.position.set(-0.6,0.3,0); g.add(startSphere)
      // 终点球
      const endSphere = new THREE.Mesh(new THREE.SphereGeometry(0.08,16,16), new THREE.MeshStandardMaterial({color:0xff6600,emissive:0xff6600,emissiveIntensity:0.8}))
      endSphere.position.set(0.6,0.3,0); g.add(endSphere)
      // 贝塞尔曲线管道
      const cp1 = new THREE.Vector3(-0.3,0.8,0)
      const cp2 = new THREE.Vector3(0.3,0.8,0)
      const curve = new THREE.CubicBezierCurve3(new THREE.Vector3(-0.6,0.3,0),cp1,cp2,new THREE.Vector3(0.6,0.3,0))
      const tubeGeo = new THREE.TubeGeometry(curve,48,0.025,8,false)
      const tube = new THREE.Mesh(tubeGeo, new THREE.MeshStandardMaterial({color:0x00ccff,emissive:0x004488,emissiveIntensity:0.5,transparent:true,opacity:0.7}))
      g.add(tube)
      // 流动光点
      const dot = new THREE.Mesh(new THREE.SphereGeometry(0.04,8,8), new THREE.MeshStandardMaterial({color:0xffffff,emissive:0xffffff,emissiveIntensity:1,depthWrite:false}))
      dot.name = 'flowDot'; g.add(dot)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true; g.userData.flyCurve = curve
      return g
    }
    case 'particleCloud': {
      const g = new THREE.Group()
      const count = 200
      const geo = new THREE.BufferGeometry()
      const positions = new Float32Array(count*3)
      const colors = new Float32Array(count*3)
      for (let i=0;i<count;i++) {
        const angle = Math.random()*Math.PI*2
        const radius = 0.2+Math.random()*0.5
        positions[i*3] = Math.cos(angle)*radius
        positions[i*3+1] = (Math.random()-0.5)*0.8
        positions[i*3+2] = Math.sin(angle)*radius
        colors[i*3] = 0.2+Math.random()*0.8
        colors[i*3+1] = 0.5+Math.random()*0.5
        colors[i*3+2] = 1.0
      }
      geo.setAttribute('position',new THREE.BufferAttribute(positions,3))
      geo.setAttribute('color',new THREE.BufferAttribute(colors,3))
      const mat = new THREE.PointsMaterial({size:0.04,vertexColors:true,transparent:true,opacity:0.8,depthWrite:false,blending:THREE.AdditiveBlending})
      const particles = new THREE.Points(geo,mat)
      particles.name = 'particles'; g.add(particles)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true; g.userData.originalPositions = positions
      return g
    }
    case 'glowSphere': {
      const g = new THREE.Group()
      // 核心球
      const core = new THREE.Mesh(new THREE.SphereGeometry(0.3,32,32), new THREE.MeshStandardMaterial({color:0x00ff88,emissive:0x00ff88,emissiveIntensity:1,roughness:0.2}))
      g.add(core)
      // 外层辉光
      const glowGeo = new THREE.SphereGeometry(0.42,32,32)
      const glowMat = new THREE.MeshBasicMaterial({color:0x00ff88,transparent:true,opacity:0.15,depthWrite:false,blending:THREE.AdditiveBlending})
      g.add(new THREE.Mesh(glowGeo,glowMat))
      // 更外层
      const outerGeo = new THREE.SphereGeometry(0.55,32,32)
      const outerMat = new THREE.MeshBasicMaterial({color:0x00ff88,transparent:true,opacity:0.05,depthWrite:false,blending:THREE.AdditiveBlending})
      g.add(new THREE.Mesh(outerGeo,outerMat))
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'haloRing': {
      const g = new THREE.Group()
      // 实体环
      const ring = new THREE.Mesh(new THREE.TorusGeometry(0.4,0.03,16,64), new THREE.MeshStandardMaterial({color:0x00ffcc,emissive:0x00ffcc,emissiveIntensity:0.8,metalness:0.3,roughness:0.2}))
      g.add(ring)
      // 外层光晕
      const glowGeo = new THREE.TorusGeometry(0.4,0.1,16,64)
      const glowMat = new THREE.MeshBasicMaterial({color:0x00ffcc,transparent:true,opacity:0.12,depthWrite:false,blending:THREE.AdditiveBlending})
      g.add(new THREE.Mesh(glowGeo,glowMat))
      // 旋转粒子环
      const particleCount = 80
      const pGeo = new THREE.BufferGeometry()
      const pPositions = new Float32Array(particleCount*3)
      for (let i=0;i<particleCount;i++) {
        const a = (i/particleCount)*Math.PI*2
        pPositions[i*3] = Math.cos(a)*0.4
        pPositions[i*3+1] = Math.sin(a)*0.4
        pPositions[i*3+2] = 0
      }
      pGeo.setAttribute('position',new THREE.BufferAttribute(pPositions,3))
      const pMat = new THREE.PointsMaterial({size:0.03,color:0xffffff,transparent:true,opacity:0.7,depthWrite:false,blending:THREE.AdditiveBlending})
      const particles = new THREE.Points(pGeo,pMat)
      particles.name = 'ringParticles'; g.add(particles)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'ringScan': {
      const g = new THREE.Group()
      // 地面大环
      const baseRing = new THREE.Mesh(new THREE.TorusGeometry(0.55,0.01,8,64), new THREE.MeshStandardMaterial({color:0x00ffff,emissive:0x00ffff,emissiveIntensity:0.6,transparent:true,opacity:0.8}))
      baseRing.rotation.x = -Math.PI/2; g.add(baseRing)
      // 扫描扇形
      const scanGeo = new THREE.RingGeometry(0.05,0.55,64,1,0,Math.PI/4)
      const scanMat = new THREE.MeshBasicMaterial({color:0x00ffff,transparent:true,opacity:0.35,side:THREE.DoubleSide,depthWrite:false,blending:THREE.AdditiveBlending})
      const scanSector = new THREE.Mesh(scanGeo,scanMat)
      scanSector.rotation.x = -Math.PI/2; scanSector.name = 'scanSector'
      g.add(scanSector)
      // 中心点
      const center = new THREE.Mesh(new THREE.SphereGeometry(0.06,16,16), new THREE.MeshStandardMaterial({color:0x00ffff,emissive:0x00ffff,emissiveIntensity:1}))
      center.position.y = 0.02; g.add(center)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'gridFloor': {
      const g = new THREE.Group()
      const size = 1.2; const divisions = 12
      const grid = new THREE.GridHelper(size,divisions,0x00aaff,0x004466)
      grid.material.transparent = true; grid.material.opacity = 0.4; grid.material.depthWrite = false
      g.add(grid)
      // 边缘发光框
      const edgeGeo = new THREE.EdgesGeometry(new THREE.PlaneGeometry(size,size))
      const edgeLine = new THREE.LineSegments(edgeGeo, new THREE.LineBasicMaterial({color:0x00ccff,transparent:true,opacity:0.6}))
      edgeLine.rotation.x = -Math.PI/2; edgeLine.position.y = 0.005
      g.add(edgeLine)
      g.userData.isMeshGroup = true
      return g
    }
    case 'electricArc': {
      const g = new THREE.Group()
      const segments = 12
      const points = [new THREE.Vector3(0,0.6,0)]
      for (let i=1;i<segments;i++) {
        const t = i/segments
        const x = (Math.random()-0.5)*0.15
        const z = (Math.random()-0.5)*0.15
        points.push(new THREE.Vector3(x,0.6*(1-t),z))
      }
      points.push(new THREE.Vector3(0,0,0))
      const curve = new THREE.CatmullRomCurve3(points)
      const tubeGeo = new THREE.TubeGeometry(curve,24,0.015,8,false)
      const tube = new THREE.Mesh(tubeGeo, new THREE.MeshBasicMaterial({color:0x44aaff,transparent:true,opacity:0.7,depthWrite:false}))
      tube.name = 'arcTube'; g.add(tube)
      // 发光核心线
      const coreGeo = new THREE.TubeGeometry(curve,24,0.005,4,false)
      const core = new THREE.Mesh(coreGeo, new THREE.MeshBasicMaterial({color:0xffffff,transparent:true,opacity:0.9,depthWrite:false,blending:THREE.AdditiveBlending}))
      core.name = 'arcCore'; g.add(core)
      // 端点球
      for (const pt of [points[0],points[points.length-1]]) {
        const s = new THREE.Mesh(new THREE.SphereGeometry(0.04,8,8), new THREE.MeshBasicMaterial({color:0x88ccff,blending:THREE.AdditiveBlending}))
        s.position.copy(pt); g.add(s)
      }
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    // ===== 新增酷炫特效 =====
    case 'energyWall': {
      // 光墙：竖直半透明渐变墙体，Shader实现上密下疏透明度
      const g = new THREE.Group()
      const width = 1.2, height = 1.0, depth = 0.06
      const geo = new THREE.BoxGeometry(width, height, depth, 1, 32, 1)
      // 自定义顶点UV：v坐标沿Y轴从0→1
      const pos = geo.attributes.position
      const vUv = new Float32Array(pos.count * 2)
      for (let i = 0; i < pos.count; i++) {
        const y = pos.getY(i)
        vUv[i*2] = 0.5 + pos.getX(i) / width   // u: 0~1
        vUv[i*2+1] = (y + height/2) / height       // v: 0(底)~1(顶)
      }
      geo.setAttribute('uv2', new THREE.BufferAttribute(vUv, 2))
      geo.setAttribute('uv', new THREE.BufferAttribute(vUv, 2))
      const mat = new THREE.ShaderMaterial({
        transparent: true,
        depthWrite: false,
        side: THREE.DoubleSide,
        uniforms: {
          color1: { value: new THREE.Color(0x00ccff) },
          time:   { value: 0 }
        },
        vertexShader: `
          varying vec2 vUv;
          void main(){
            vUv = uv;
            gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.0);
          }`,
        fragmentShader: `
          uniform vec3 color1;
          uniform float time;
          varying vec2 vUv;
          void main(){
            // 上密下疏：v越接近1（顶部）透明度越低（更不透明）→ 反之底部更透明
            float alpha = mix(0.05, 0.75, vUv.y);
            // 横向扫描光带
            float scan = smoothstep(0.0, 0.04, abs(fract(vUv.x * 3.0 + time * 0.5) - 0.5)) * 0.5 + 0.5;
            vec3 col = color1 * (0.8 + 0.2 * scan);
            gl_FragColor = vec4(col, alpha * (0.7 + 0.3 * sin(vUv.y * 6.283 + time)));
          }`
      })
      const mesh = new THREE.Mesh(geo, mat)
      mesh.name = 'wallBody'
      g.add(mesh)
      // 边框线
      const edges = new THREE.EdgesGeometry(geo)
      const edgeLine = new THREE.LineSegments(edges, new THREE.LineBasicMaterial({color:0x00ccff, transparent:true, opacity:0.6}))
      g.add(edgeLine)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'pulseColumn': {
      // 数据脉冲柱：底部粗→顶部细的柱体 + 顶部脉冲光晕
      const g = new THREE.Group()
      // 柱体（圆锥台形）
      const pillarGeo = new THREE.CylinderGeometry(0.08, 0.15, 0.9, 16, 1, true)
      const pillarMat = new THREE.ShaderMaterial({
        transparent: true, depthWrite: false, side: THREE.DoubleSide,
        uniforms: {
          color1: { value: new THREE.Color(0x00ff88) },
          time:   { value: 0 }
        },
        vertexShader: `
          varying vec2 vUv;
          void main(){
            vUv = uv;
            gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.0);
          }`,
        fragmentShader: `
          uniform vec3 color1;
          uniform float time;
          varying vec2 vUv;
          void main(){
            // 呼吸：透明度随时间和高度变化
            float breath = 0.5 + 0.5 * sin(time * 2.0);
            float alpha = mix(0.9, 0.15, vUv.y) * breath;
            vec3 col = color1 * (1.0 + 0.3 * sin(vUv.y * 12.566 + time * 3.0));
            gl_FragColor = vec4(col, alpha);
          }`
      })
      const pillar = new THREE.Mesh(pillarGeo, pillarMat)
      pillar.position.y = 0.45; g.add(pillar)
      // 顶部光晕（Additive）
      const glowGeo = new THREE.SphereGeometry(0.18, 16, 16)
      const glowMat = new THREE.MeshBasicMaterial({color:0x00ff88, transparent:true, opacity:0.25, blending:THREE.AdditiveBlending, depthWrite:false})
      const glow = new THREE.Mesh(glowGeo, glowMat)
      glow.position.y = 0.92; glow.name = 'topGlow'; g.add(glow)
      // 底部光环
      const ringGeo = new THREE.TorusGeometry(0.18, 0.015, 8, 48)
      const ringMat = new THREE.MeshBasicMaterial({color:0x00ff88, transparent:true, opacity:0.7, depthWrite:false})
      const ring = new THREE.Mesh(ringGeo, ringMat)
      ring.rotation.x = -Math.PI/2; ring.position.y = 0.01; g.add(ring)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'heatmapCloud': {
      // 热力图点云：按高度/强度着色的大屏风格粒子云
      const g = new THREE.Group()
      const count = 300
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const colors = new Float32Array(count * 3)
      const sizes = new Float32Array(count)
      for (let i = 0; i < count; i++) {
        const angle = Math.random() * Math.PI * 2
        const r = Math.random() * 0.6
        const y = (Math.random() - 0.3) * 0.8
        pos[i*3]   = Math.cos(angle) * r
        pos[i*3+1] = y
        pos[i*3+2] = Math.sin(angle) * r
        // 热力图颜色映射：蓝→青→绿→黄→红
        const t = Math.random()
        let r1, g1, b1
        if (t < 0.25) { r1=t*4; g1=0; b1=1 }
        else if (t < 0.5) { r1=0; g1=(t-0.25)*4; b1=1-(t-0.25)*4 }
        else if (t < 0.75) { r1=0; g1=1; b1=(t-0.5)*4 }
        else { r1=(t-0.75)*4; g1=1-(t-0.75)*4; b1=0 }
        colors[i*3]   = r1; colors[i*3+1] = g1; colors[i*3+2] = b1
        sizes[i] = 0.02 + Math.random() * 0.04
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('color',    new THREE.BufferAttribute(colors, 3))
      geo.setAttribute('size',     new THREE.BufferAttribute(sizes, 1))
      const mat = new THREE.ShaderMaterial({
        transparent: true, depthWrite: false, blending: THREE.AdditiveBlending,
        uniforms: { time: { value: 0 }, scale: { value: 1 } },
        vertexShader: `
          attribute float size;
          uniform float time;
          varying vec3 vColor;
          void main(){
            vColor = color;
            // 粒子上下浮动
            vec3 p = position;
            p.y += sin(time * 1.5 + position.x * 8.0) * 0.03;
            vec4 mv = modelViewMatrix * vec4(p,1.0);
            gl_PointSize = size * (300.0 / -mv.z);
            gl_Position = projectionMatrix * mv;
          }`,
        fragmentShader: `
          varying vec3 vColor;
          void main(){
            // 圆形粒子
            float d = length(gl_PointCoord - vec2(0.5));
            if(d > 0.5) discard;
            float alpha = smoothstep(0.5, 0.1, d);
            gl_FragColor = vec4(vColor * 1.5, alpha * 0.85);
          }`
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'heatPoints'; g.add(points)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'energyShield': {
      // 能量护盾：锐利菲涅尔边缘光 + 扫描脉冲（来自掘金文章）
      const g = new THREE.Group()
      const radius = 0.5, thickness = 0.06
      const geo = new THREE.SphereGeometry(radius, 32, 32)
      const mat = new THREE.ShaderMaterial({
        transparent: true, depthWrite: false, side: THREE.DoubleSide,
        blending: THREE.AdditiveBlending,
        uniforms: {
          color1: { value: new THREE.Color(0x00ddff) },
          time:  { value: 0 }
        },
        vertexShader: `
          varying vec3 vWorldNormal;
          varying vec3 vViewDir;
          void main(){
            vec4 wp = modelMatrix * vec4(position,1.0);
            vWorldNormal = normalize(mat3(modelMatrix) * normal);
            vViewDir = normalize(cameraPosition - wp.xyz);
            gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.0);
          }`,
        fragmentShader: `
          uniform vec3 color1;
          uniform float time;
          varying vec3 vWorldNormal;
          varying vec3 vViewDir;
          void main(){
            float fresnel = 1.0 - abs(dot(vViewDir, vWorldNormal));
            fresnel = pow(fresnel, 6.0);          // 锐利边缘
            fresnel = max(fresnel - 0.08, 0.0);  // 截断中心雾化
            fresnel *= 2.5;                        // 亮度补偿
            // 扫描脉冲：沿Y轴扫过的亮环
            float scanRing = smoothstep(0.08, 0.0, abs(sin(time * 1.5 + (position.y / 1.0 + 0.5) * 3.14159)) - 0.92);
            vec3 col = color1 * (fresnel + scanRing * 1.5);
            float alpha = fresnel * 0.85 + scanRing * 0.5;
            gl_FragColor = vec4(col, alpha);
          }`
      })
      const shield = new THREE.Mesh(geo, mat)
      shield.scale.set(1, 1, 1)  // 可拉伸成椭球护盾
      shield.name = 'shieldBody'; g.add(shield)
      // 外轮廓线框
      const wireGeo = new THREE.WireframeGeometry(geo)
      const wire = new THREE.LineSegments(wireGeo, new THREE.LineBasicMaterial({color:0x00ddff, transparent:true, opacity:0.15}))
      g.add(wire)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'gradientBlock': {
      // 建筑渐变区块：长方体 + 高度渐变 Shader（模拟智慧城市建筑）
      const g = new THREE.Group()
      const w = 0.5, d = 0.5, h = 0.8
      const geo = new THREE.BoxGeometry(w, h, d, 1, 8, 1)
      // 手动设置每个顶点的uv.y = 高度比例
      const posAttr = geo.attributes.position
      const uvArr = new Float32Array(posAttr.count * 2)
      for (let i = 0; i < posAttr.count; i++) {
        const y = posAttr.getY(i)
        uvArr[i*2]   = 0.5 + posAttr.getX(i) / w * 0.5
        uvArr[i*2+1] = (y + h/2) / h   // 0(底) ~ 1(顶)
      }
      geo.setAttribute('uv', new THREE.BufferAttribute(uvArr, 2))
      const mat = new THREE.ShaderMaterial({
        transparent: false, depthWrite: true,
        uniforms: {
          colorLow:  { value: new THREE.Color(0x003366) },
          colorHigh: { value: new THREE.Color(0x00ccff) },
          time: { value: 0 }
        },
        vertexShader: `
          varying vec2 vUv;
          void main(){
            vUv = uv;
            gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.0);
          }`,
        fragmentShader: `
          uniform vec3 colorLow;
          uniform vec3 colorHigh;
          uniform float time;
          varying vec2 vUv;
          void main(){
            // 高度渐变
            vec3 col = mix(colorLow, colorHigh, vUv.y);
            // 顶部高亮光带
            float topGlow = smoothstep(0.92, 1.0, vUv.y) * 0.6;
            // 扫描光（横向）
            float scan = smoothstep(0.0, 0.05, abs(fract(vUv.y * 5.0 - time * 0.3) - 0.5)) * 0.3;
            col += vec3(0.2, 0.6, 1.0) * (topGlow + scan);
            gl_FragColor = vec4(col, 1.0);
          }`
      })
      const block = new THREE.Mesh(geo, mat)
      block.position.y = h / 2; block.name = 'blockBody'; g.add(block)
      // 顶部发光面
      const topGeo = new THREE.PlaneGeometry(w, d)
      const topMat = new THREE.MeshBasicMaterial({color:0x00ccff, transparent:true, opacity:0.35, blending:THREE.AdditiveBlending, depthWrite:false})
      const top = new THREE.Mesh(topGeo, topMat)
      top.rotation.x = -Math.PI/2; top.position.y = h; g.add(top)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'flowRibbon': {
      // 流动光带：宽曲线带 + 沿曲线流动的UV偏移（增强版飞线）
      const g = new THREE.Group()
      // 创建正弦波浪曲线
      const points = []
      for (let i = 0; i <= 60; i++) {
        const t = i / 60
        const x = (t - 0.5) * 1.6
        const y = 0.3 + Math.sin(t * Math.PI * 3) * 0.15
        const z = Math.cos(t * Math.PI * 2) * 0.2
        points.push(new THREE.Vector3(x, y, z))
      }
      const curve = new THREE.CatmullRomCurve3(points, false, 'catmullrom', 0.3)
      // 用 TubeGeometry 模拟光带
      const tubeGeo = new THREE.TubeGeometry(curve, 60, 0.02, 6, false)
      const mat = new THREE.ShaderMaterial({
        transparent: true, depthWrite: false, blending: THREE.AdditiveBlending,
        uniforms: {
          color1: { value: new THREE.Color(0xff6600) },
          color2: { value: new THREE.Color(0x00ccff) },
          time:  { value: 0 }
        },
        vertexShader: `
          varying vec2 vUv;
          void main(){
            vUv = uv;
            gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.0);
          }`,
        fragmentShader: `
          uniform vec3 color1;
          uniform vec3 color2;
          uniform float time;
          varying vec2 vUv;
          void main(){
            // 沿曲线方向流动的光带
            float flow = fract(vUv.x * 3.0 - time * 0.8);
            float pulse = smoothstep(0.0, 0.3, flow) * smoothstep(1.0, 0.7, flow);
            vec3 col = mix(color1, color2, vUv.x);
            col += vec3(1.0) * pulse * 0.7;
            float alpha = (0.5 + 0.5 * pulse) * 0.8;
            gl_FragColor = vec4(col, alpha);
          }`
      })
      const tube = new THREE.Mesh(tubeGeo, mat)
      tube.name = 'ribbonBody'; g.add(tube)
      // 起点/终点光球
      for (const [i, offset] of [[0, 0], [points.length-1, 1]]) {
        const c = i === 0 ? 0xff6600 : 0x00ccff
        const s = new THREE.Mesh(new THREE.SphereGeometry(0.04, 12, 12), new THREE.MeshBasicMaterial({color:c, blending:THREE.AdditiveBlending, depthWrite:false}))
        s.position.copy(points[offset]); g.add(s)
      }
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      g.userData.flowCurve = curve
      return g
    }

    // ===== 新增视觉特效（翠鸟风格） =====
    case 'fireFX': {
      // 火焰：Shuriken粒子风格的火焰柱体
      const g = new THREE.Group()
      const count = 150
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const colors = new Float32Array(count * 3)
      for (let i = 0; i < count; i++) {
        const angle = Math.random() * Math.PI * 2
        const radius = 0.05 + Math.random() * 0.25
        const y = Math.random() * 0.9
        pos[i*3]   = Math.cos(angle) * radius * (1 - y * 1.3)
        pos[i*3+1] = y
        pos[i*3+2] = Math.sin(angle) * radius * (1 - y * 1.3)
        const t = y
        colors[i*3]   = 1.0
        colors[i*3+1] = 0.2 + t * 0.5
        colors[i*3+2] = t * 0.15
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
      const mat = new THREE.PointsMaterial({
        size: 0.06, vertexColors: true, transparent: true, opacity: 0.85,
        depthWrite: false, blending: THREE.AdditiveBlending,
        map: createCircleTexture(0xffffff)
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'fireParticles'; g.add(points)
      // 火芯
      const coreGeo = new THREE.ConeGeometry(0.1, 0.4, 8, 1, true)
      const coreMat = new THREE.MeshBasicMaterial({color: 0xff6622, transparent: true, opacity: 0.5, depthWrite: false, blending: THREE.AdditiveBlending})
      const core = new THREE.Mesh(coreGeo, coreMat)
      core.position.y = 0.15; core.name = 'fireCore'; g.add(core)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'smokeFX': {
      // 烟雾：柔和上升的灰色粒子团
      const g = new THREE.Group()
      const count = 120
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const colors = new Float32Array(count * 3)
      for (let i = 0; i < count; i++) {
        const angle = Math.random() * Math.PI * 2
        const radius = 0.08 + Math.random() * 0.35
        const y = Math.random() * 1.0
        pos[i*3]   = Math.cos(angle) * radius
        pos[i*3+1] = y
        pos[i*3+2] = Math.sin(angle) * radius
        const gray = 0.4 + Math.random() * 0.5
        colors[i*3] = gray; colors[i*3+1] = gray; colors[i*3+2] = gray
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
      const mat = new THREE.PointsMaterial({
        size: 0.1, vertexColors: true, transparent: true, opacity: 0.55,
        depthWrite: false,
        map: createCircleTexture(0xffffff)
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'smokeParticles'; g.add(points)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'dustFX': {
      // 灰尘：地面扩散的微小粉尘
      const g = new THREE.Group()
      const count = 200
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const colors = new Float32Array(count * 3)
      for (let i = 0; i < count; i++) {
        const angle = Math.random() * Math.PI * 2
        const radius = 0.05 + Math.random() * 0.6
        const y = Math.random() * 0.25
        pos[i*3]   = Math.cos(angle) * radius
        pos[i*3+1] = y
        pos[i*3+2] = Math.sin(angle) * radius
        const gray = 0.6 + Math.random() * 0.35
        colors[i*3] = gray; colors[i*3+1] = gray * 0.85; colors[i*3+2] = gray * 0.7
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      const mat = new THREE.PointsMaterial({
        size: 0.025, color: 0xc0a878, transparent: true, opacity: 0.6,
        depthWrite: false, blending: THREE.NormalBlending
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'dustParticles'; g.add(points)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'fireflyFX': {
      // 萤火虫：明灭闪烁的绿黄光点
      const g = new THREE.Group()
      const count = 30
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const sizes = new Float32Array(count)
      for (let i = 0; i < count; i++) {
        pos[i*3]   = (Math.random() - 0.5) * 1.5
        pos[i*3+1] = 0.2 + Math.random() * 1.0
        pos[i*3+2] = (Math.random() - 0.5) * 1.5
        sizes[i] = 3 + Math.random() * 6
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('size', new THREE.BufferAttribute(sizes, 1))
      const mat = new THREE.PointsMaterial({
        size: 0.06, color: 0xccff33, transparent: true, opacity: 0.9,
        depthWrite: false, blending: THREE.AdditiveBlending,
        map: createCircleTexture(0xffffff)
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'fireflyParticles'; g.add(points)
      // 存储每只萤火虫的相位和速度
      g.userData.fireflyPhases = new Float32Array(count)
      g.userData.fireflySpeeds = new Float32Array(count)
      for (let i = 0; i < count; i++) {
        g.userData.fireflyPhases[i] = Math.random() * Math.PI * 2
        g.userData.fireflySpeeds[i] = 0.3 + Math.random() * 0.8
      }
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'snowFX': {
      // 雪花：白色飘落粒子
      const g = new THREE.Group()
      const count = 80
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      for (let i = 0; i < count; i++) {
        pos[i*3]   = (Math.random() - 0.5) * 1.6
        pos[i*3+1] = 0.2 + Math.random() * 1.2
        pos[i*3+2] = (Math.random() - 0.5) * 1.6
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      const mat = new THREE.PointsMaterial({
        size: 0.05, color: 0xffffff, transparent: true, opacity: 0.85,
        depthWrite: false,
        map: createCircleTexture(0xffffff)
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'snowParticles'; g.add(points)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'rainFX': {
      // 雨滴：下落的水蓝线条
      const g = new THREE.Group()
      const count = 100
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 6)
      for (let i = 0; i < count; i++) {
        const x = (Math.random() - 0.5) * 1.5
        const y = Math.random() * 1.5
        const z = (Math.random() - 0.5) * 1.5
        pos[i*6]   = x
        pos[i*6+1] = y + 0.12
        pos[i*6+2] = z
        pos[i*6+3] = x
        pos[i*6+4] = y
        pos[i*6+5] = z
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      const mat = new THREE.LineBasicMaterial({
        color: 0x4488cc, transparent: true, opacity: 0.5, depthWrite: false
      })
      const lines = new THREE.LineSegments(geo, mat)
      lines.name = 'rainLines'; g.add(lines)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'bubbleFX': {
      // 气泡：透明的浮动球体
      const g = new THREE.Group()
      const count = 20
      for (let i = 0; i < count; i++) {
        const r = 0.03 + Math.random() * 0.08
        const s = new THREE.Mesh(
          new THREE.SphereGeometry(r, 12, 12),
          new THREE.MeshPhongMaterial({
            color: 0x88ccff, emissive: 0x224466, emissiveIntensity: 0.15,
            transparent: true, opacity: 0.35, shininess: 60, depthWrite: false
          })
        )
        s.position.set((Math.random()-0.5)*1.2, 0.05 + Math.random()*0.9, (Math.random()-0.5)*1.2)
        s.name = 'bubble_' + i; g.add(s)
      }
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'leafFX': {
      // 落叶：旋转飘落的扁平方形/椭圆模拟叶片
      const g = new THREE.Group()
      const count = 18
      const leafColors = [0xd4a64a, 0xc47a38, 0x8b5e3c, 0xb8860a, 0x9b6b2a]
      for (let i = 0; i < count; i++) {
        const size = 0.04 + Math.random() * 0.07
        const leafGeo = new THREE.PlaneGeometry(size, size * 1.8)
        const leafMat = new THREE.MeshStandardMaterial({
          color: leafColors[Math.floor(Math.random() * leafColors.length)],
          side: THREE.DoubleSide, transparent: true, opacity: 0.85, roughness: 0.7, depthWrite: false
        })
        const leaf = new THREE.Mesh(leafGeo, leafMat)
        leaf.position.set((Math.random()-0.5)*1.4, 0.2 + Math.random()*1.1, (Math.random()-0.5)*1.4)
        leaf.rotation.set(Math.random()*Math.PI, Math.random()*Math.PI, Math.random()*Math.PI)
        leaf.name = 'leaf_' + i; g.add(leaf)
      }
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'sparkFX': {
      // 火花：沿锥形方向喷射的亮白粒子
      const g = new THREE.Group()
      const count = 120
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const colors = new Float32Array(count * 3)
      for (let i = 0; i < count; i++) {
        const angle = Math.random() * Math.PI * 2
        const spread = 0.1 + Math.random() * 0.5
        const y = Math.random() * 0.7
        pos[i*3]   = Math.cos(angle) * spread * (0.3 + y * 0.7)
        pos[i*3+1] = y
        pos[i*3+2] = Math.sin(angle) * spread * (0.3 + y * 0.7)
        const t = Math.random()
        colors[i*3]   = 1.0
        colors[i*3+1] = 0.7 + t * 0.25
        colors[i*3+2] = t * 0.5
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
      const mat = new THREE.PointsMaterial({
        size: 0.03, vertexColors: true, transparent: true, opacity: 0.9,
        depthWrite: false, blending: THREE.AdditiveBlending,
        map: createCircleTexture(0xffffff)
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'sparkParticles'; g.add(points)
      // 中心光球
      const core = new THREE.Mesh(new THREE.SphereGeometry(0.06, 8, 8),
        new THREE.MeshBasicMaterial({color: 0xffffff, blending: THREE.AdditiveBlending, depthWrite: false}))
      g.add(core)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'explosionFX': {
      // 爆炸：膨胀的环形冲击波 + 碎屑粒子
      const g = new THREE.Group()
      // 冲击波环
      const ringGeo = new THREE.TorusGeometry(0.15, 0.04, 12, 48)
      const ringMat = new THREE.MeshBasicMaterial({
        color: 0xff6633, transparent: true, opacity: 0.8, depthWrite: false, blending: THREE.AdditiveBlending
      })
      const ring = new THREE.Mesh(ringGeo, ringMat)
      ring.rotation.x = Math.PI / 2; ring.name = 'explosionRing'; g.add(ring)
      // 内环
      const inRing = new THREE.Mesh(new THREE.TorusGeometry(0.08, 0.03, 8, 32),
        new THREE.MeshBasicMaterial({color: 0xffcc00, transparent: true, opacity: 0.6, depthWrite: false, blending: THREE.AdditiveBlending}))
      inRing.rotation.x = Math.PI / 2; inRing.name = 'explosionInnerRing'; g.add(inRing)
      // 碎屑粒子
      const count = 80
      const geo = new THREE.BufferGeometry()
      const posArr = new Float32Array(count * 3)
      const colArr = new Float32Array(count * 3)
      for (let i = 0; i < count; i++) {
        const a = Math.random() * Math.PI * 2
        const b = Math.random() * Math.PI
        const r = 0.08 + Math.random() * 0.1
        posArr[i*3]   = Math.sin(b) * Math.cos(a) * r
        posArr[i*3+1] = Math.sin(b) * Math.sin(a) * r
        posArr[i*3+2] = Math.cos(b) * r
        colArr[i*3]   = 1.0; colArr[i*3+1] = 0.4 + Math.random() * 0.4; colArr[i*3+2] = 0.0
      }
      geo.setAttribute('position', new THREE.BufferAttribute(posArr, 3))
      geo.setAttribute('color', new THREE.BufferAttribute(colArr, 3))
      const pMat = new THREE.PointsMaterial({size: 0.04, vertexColors: true, transparent: true, opacity: 0.9, depthWrite: false, blending: THREE.AdditiveBlending})
      const pts = new THREE.Points(geo, pMat)
      pts.name = 'explosionParticles'; g.add(pts)
      // 中心光点
      const center = new THREE.Mesh(new THREE.SphereGeometry(0.07, 8, 8),
        new THREE.MeshBasicMaterial({color: 0xffffff, blending: THREE.AdditiveBlending, depthWrite: false}))
      g.add(center)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'splashFX': {
      // 水花：从中心向上/向外喷射的蓝色粒子
      const g = new THREE.Group()
      const count = 100
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const colors = new Float32Array(count * 3)
      for (let i = 0; i < count; i++) {
        const angle = Math.random() * Math.PI * 2
        const upBias = Math.random() < 0.6
        const dist = 0.03 + Math.random() * 0.45
        const y = upBias ? Math.random() * 0.6 : Math.random() * 0.2
        pos[i*3]   = Math.cos(angle) * dist
        pos[i*3+1] = y
        pos[i*3+2] = Math.sin(angle) * dist
        colors[i*3]   = 0.2 + Math.random() * 0.2
        colors[i*3+1] = 0.5 + Math.random() * 0.4
        colors[i*3+2] = 0.85 + Math.random() * 0.15
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
      const mat = new THREE.PointsMaterial({
        size: 0.04, vertexColors: true, transparent: true, opacity: 0.8,
        depthWrite: false, blending: THREE.AdditiveBlending,
        map: createCircleTexture(0xffffff)
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'splashParticles'; g.add(points)
      // 底部水环
      const waterRing = new THREE.Mesh(new THREE.TorusGeometry(0.25, 0.02, 8, 32),
        new THREE.MeshBasicMaterial({color: 0x4488ff, transparent: true, opacity: 0.5, depthWrite: false}))
      waterRing.rotation.x = -Math.PI/2; waterRing.position.y = 0.01; g.add(waterRing)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'starFX': {
      // 星星：闪烁的星形光点
      const g = new THREE.Group()
      const count = 40
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const sizes = new Float32Array(count)
      for (let i = 0; i < count; i++) {
        pos[i*3]   = (Math.random() - 0.5) * 1.8
        pos[i*3+1] = 0.2 + Math.random() * 1.2
        pos[i*3+2] = (Math.random() - 0.5) * 1.8
        sizes[i] = 2 + Math.random() * 4
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('size', new THREE.BufferAttribute(sizes, 1))
      const mat = new THREE.PointsMaterial({
        size: 0.05, color: 0xffffcc, transparent: true, opacity: 0.9,
        depthWrite: false, blending: THREE.AdditiveBlending,
        map: createStarTexture()
      })
      const points = new THREE.Points(geo, mat)
      points.name = 'starParticles'; g.add(points)
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }
    case 'magicGlow': {
      // 魔法光晕：翠鸟风格——大量白色微小粒子散布在空间中
      const g = new THREE.Group()
      const count = 350
      const geo = new THREE.BufferGeometry()
      const pos = new Float32Array(count * 3)
      const sizes = new Float32Array(count)
      const phases = new Float32Array(count)
      for (let i = 0; i < count; i++) {
        // 在球形空间内随机分布
        const r = 0.05 + Math.random() * 1.2
        const theta = Math.random() * Math.PI * 2
        const phi = Math.acos(2 * Math.random() - 1)
        pos[i*3]   = r * Math.sin(phi) * Math.cos(theta)
        pos[i*3+1] = r * Math.sin(phi) * Math.sin(theta) * 0.6 + 0.3 // 略微偏上
        pos[i*3+2] = r * Math.cos(phi)
        sizes[i] = 0.008 + Math.random() * 0.018
        phases[i] = Math.random() * Math.PI * 2
      }
      geo.setAttribute('position', new THREE.BufferAttribute(pos, 3))
      geo.setAttribute('size', new THREE.BufferAttribute(sizes, 1))
      const mat = new THREE.PointsMaterial({
        size: 0.015, color: 0xffffff, transparent: true, opacity: 0.75,
        depthWrite: false, blending: THREE.AdditiveBlending,
        map: createCircleTexture(0xffffff)
      })
      const particles = new THREE.Points(geo, mat)
      particles.name = 'magicGlowParticles'; g.add(particles)
      // 存储动画数据
      g.userData.magicPhases = phases
      g.userData.magicBasePos = pos.slice()
      g.userData.isMeshGroup = true; g.userData.hasAnimation = true
      return g
    }

    // ===== 灯光 =====
    case 'point_light': {
      const group = new THREE.Group()
      const sphere = new THREE.Mesh(new THREE.SphereGeometry(0.15,16,16),
        new THREE.MeshStandardMaterial({ color: 0xfffb4d, emissive: 0xfffb4d, emissiveIntensity: 0.8 }))
      const light = new THREE.PointLight(0xfffb4d, 1.5, 8)
      group.add(sphere, light)
      group.userData.isMeshGroup = true
      group.userData.isLight = true
      return group
    }
    case 'spot_light': {
      const group = new THREE.Group()
      const cone = new THREE.Mesh(new THREE.ConeGeometry(0.2,0.5,16),
        new THREE.MeshStandardMaterial({ color: 0x88ccff, emissive: 0x4499ff, emissiveIntensity: 0.5 }))
      const light = new THREE.SpotLight(0x88ccff, 1.5, 10, Math.PI/6)
      light.position.y = -0.5
      group.add(cone, light)
      group.userData.isMeshGroup = true
      group.userData.isLight = true
      return group
    }
    case 'ambient_light': {
      const group = new THREE.Group()
      const sphere = new THREE.Mesh(new THREE.SphereGeometry(0.2, 20, 20),
        new THREE.MeshStandardMaterial({ color: 0xfff1b8, emissive: 0xffd666, emissiveIntensity: 0.6 }))
      const light = new THREE.AmbientLight(0xffffff, 0.8)
      group.add(sphere, light)
      group.userData.isMeshGroup = true
      group.userData.isLight = true
      return group
    }
    case 'directional_light': {
      const group = new THREE.Group()
      const arrow = new THREE.ArrowHelper(new THREE.Vector3(-1, -1, -0.5).normalize(), new THREE.Vector3(0, 0, 0), 1.2, 0x88ccff, 0.35, 0.18)
      const light = new THREE.DirectionalLight(0xffffff, 0.9)
      light.position.set(1, 1, 0.5)
      group.add(arrow, light)
      group.userData.isMeshGroup = true
      group.userData.isLight = true
      return group
    }

    // ===== 基础组件 =====
    case 'arrow': {
      const group = new THREE.Group()
      const shaft = new THREE.Mesh(new THREE.CylinderGeometry(0.07,0.07,1.2,12), mat.clone())
      shaft.position.y = 0.6
      const head = new THREE.Mesh(new THREE.ConeGeometry(0.2,0.5,16), mat.clone())
      head.position.y = 1.45
      group.add(shaft, head)
      group.userData.isMeshGroup = true
      return group
    }
    case 'text3d': {
      return createTextSprite('3D标签', 0xffffff, 16, '#1f1f1f', { backgroundOpacity: 0.48, fixedSize: true, showBorder: false })
    }
    case 'textPlain3d': {
      return createTextSprite('3D文字', 0xffffff, 16, '#1f1f1f', { backgroundOpacity: 0.48, fixedSize: true, showBorder: false })
    }
    case 'dataText': {
      return createTextSprite('--', 0xffffff, 16, '#1f1f1f', { backgroundOpacity: 0.48, fixedSize: true, showBorder: false })
    }
    case 'image3d':
      return createMediaPlane({ color: color, opacity: 0.92 }, 'image3d')
    case 'uiImage': {
      return createMediaPlane({ color: color, opacity: 0.92 }, 'image3d')
    }
    case 'uiLabel': {
      return createTextSprite('UI标签', color || 0x4dffa6, 32, '#000000', { backgroundOpacity: 0.35, fixedSize: true, showBorder: true, borderColor: '#4dffa6' })
    }
    case 'video3d': {
      return createMediaPlane({ opacity: 0.92, mediaAspect: 16 / 9, mediaWidth: 1.6 }, 'video3d')
    }
    case 'webEmbed': {
      const group = new THREE.Group()
      const page = new THREE.Mesh(
        new THREE.PlaneGeometry(1.5, 0.95),
        new THREE.MeshBasicMaterial({ color: 0xffffff, transparent: true, opacity: 0.72, side: THREE.DoubleSide })
      )
      const label = createTextSprite('网页', 0x1f2937, 24, '#ffffffcc', { backgroundOpacity: 0.8, showBorder: true, borderColor: '#d1d5db' })
      label.position.z = 0.02
      label.scale.set(0.7, 0.7, 1)
      group.add(page, label)
      group.userData.isMeshGroup = true
      return group
    }
    case 'cssLabel': {
      // CSS2DObject 需要场景中有 CSS2DRenderer 才能渲染
      var div = document.createElement('div')
      div.style.cssText = 'padding:4px 10px;background:rgba(19,194,194,0.14);border:1px solid #13c2c2;border-radius:12px;color:#08979c;font-size:13px;font-family:"Microsoft YaHei",sans-serif;white-space:nowrap;pointer-events:none'
      div.textContent = '数据标签'
      var cssObj = new CSS2DObject(div)
      cssObj.userData.isCSSLabel = true
      cssObj.userData.isMeshGroup = true // 兼容编辑器选择
      return cssObj
    }

    // ===== GLTF 外部模型（占位） =====
    case '2dComponent': {
      return create2DComponentPlane({
        color: '#' + (color || 0x13c2c2).toString(16).padStart(6, '0')
      })
    }
    case 'gltf': {
      const group = new THREE.Group()
      const placeholder = new THREE.Mesh(new THREE.BoxGeometry(1, 1, 1),
        new THREE.MeshStandardMaterial({ color: 0x1890ff, transparent: true, opacity: 0.3, wireframe: true }))
      const label = createTextSprite('拖入模型文件', 0x1890ff, 36, '#00000000')
      label.position.y = 0.9
      group.add(placeholder, label)
      group.userData.isMeshGroup = true
      group.userData.isGLTFPlaceholder = true
      return group
    }

    // ===== 数字孪生高端特效 =====
    case 'dtBuilding':
    case 'dtBuildingTall':
    case 'dtBuildingWide':
    case 'dtBuildingComplex':
    case 'dtGroundGrid':
    case 'dtScanRing':
    case 'dtScanLine':
    case 'dtDataPanel':
    case 'dtFlowLine':
    case 'dtParticleField':
    case 'dtHologram':
    case 'dtRoad':
    case 'dtPark':
    case 'dtStreetLightDT':
    case 'dtBaseStation': {
      const dtObj = createDTObject(type, color)
      if (dtObj) return dtObj
      break
    }

    default: {
      const extensionMesh = createIndustrialExtensionMesh(type, mat)
      if (extensionMesh) return extensionMesh
      geo = new THREE.BoxGeometry(1,1,1)
    }
  }
  return new THREE.Mesh(geo, mat)
}

// ===== 真正的3D几何体文字（需要 opentype.js + 中文字体文件）=====
// 通过 webpack 引用 src/font/ 下的字体文件，无需手动复制到 public/

var _opentypeReady = false
var _chineseFont3D = null
var _fontLoadPromise = null // 防止重复加载

// 字体注册表：fontFamily → opentype Font 对象
var _fontRegistry = {}
// fontFamily → 字体文件路径映射（供按需加载）
var FONT_FAMILY_PATHS = {
  '系统默认': '/fonts/NotoSansSC-Regular.ttf',
  'Microsoft YaHei': '/fonts/NotoSansSC-Regular.ttf',
  'Arial': '/fonts/NotoSansSC-Regular.ttf',
  // 中文字体
  '黑体': '/fonts/simhei.ttf',
  '楷体': '/fonts/simkai.ttf',
  '隶书': '/fonts/SIMLI.TTF',
  '宋体': '/fonts/simsun.ttf',
  'SimSun': '/fonts/simsun.ttf',
  // 数字字体
  '数字字体-1': '/fonts/digitalism.TTF',
  '数字字体-2': '/fonts/DS-DIGIB.TTF',
  '数字字体-3': '/fonts/DS-DIGI.TTF',
  '数字字体-4': '/fonts/DS-DIGII.TTF',
  '数字字体-5': '/fonts/DS-DIGIT.TTF',
  '数字字体-6': '/fonts/FakeHope.TTF',
  '数字字体-7': '/fonts/FakeHopeFilled.TTF',
  '数字字体-8': '/fonts/QuartzRegular.TTF',
  '数字字体-9': '/fonts/Technology.TTF',
  '数字字体-10': '/fonts/Technology-Bold.TTF',
  '数字字体-11': '/fonts/Technology-BoldItalic.TTF',
  '数字字体-12': '/fonts/Technology-Italic.TTF'
}
// 按需加载的 Promise 缓存：fontFamily → Promise
var _fontLoadPromises = {}

var FONT_CANDIDATES = [
  '/fonts/NotoSansSC-Regular.ttf',
]

function tryLoadChineseFontFor3D() {
  // 已有正在进行中或已完成的 Promise，直接返回
  if (_fontLoadPromise) return _fontLoadPromise

  if (typeof window === 'undefined' || !window.opentype) {
    console.warn('[3DText] opentype.js not loaded. True 3D text disabled.')
    _fontLoadPromise = Promise.resolve(null)
    return _fontLoadPromise
  }

  _opentypeReady = true
  // 逐个尝试加载候选字体，返回 Promise
  _fontLoadPromise = _loadFontFromCandidates(0)
  return _fontLoadPromise
}

function _loadFontFromCandidates(index) {
  if (index >= FONT_CANDIDATES.length) {
    console.warn('[3DText] No valid Chinese font found in candidates.')
    return Promise.resolve(null)
  }

  var fontUrl = FONT_CANDIDATES[index]
  return fetch(fontUrl)
    .then(function(res) {
      if (!res.ok) throw new Error('HTTP ' + res.status)
      return res.arrayBuffer()
    })
    .then(function(buffer) {
      var font = window.opentype.parse(buffer)
      if (font && font.outlinesFormat) {
        _chineseFont3D = font
        // 初始加载的是系统默认字体，只注册到 fallback 族名
        // 其他字体（黑体/楷体等）将通过 getTrue3DFont 按需加载
        _fontRegistry['系统默认'] = font
        _fontRegistry['Microsoft YaHei'] = font
        _fontRegistry['Arial'] = font
        console.log('[3DText] Chinese font loaded for true 3D geometry from:', fontUrl, 'registered as: 系统默认')
        // 字体加载完成后触发全局事件，通知场景刷新 Canvas fallback 文字
        if (typeof window !== 'undefined') {
          window.dispatchEvent(new CustomEvent('ism3d-font-ready'))
        }
        return font
      } else {
        throw new Error('Invalid font file')
      }
    })
    .catch(function(err) {
      console.warn('[3DText] Failed to load font from ' + fontUrl + ':', err.message)
      // 尝试下一个候选字体
      return _loadFontFromCandidates(index + 1)
    })
}

/**
 * 根据字体族名获取已加载的 opentype 字体对象
 * @param {string} fontFamily - 字体族名，如 'SimSun'、'Microsoft YaHei'
 * @returns {opentype.Font|null}
 */
function getTrue3DFont(fontFamily) {
  if (!fontFamily || fontFamily === '系统默认') return _chineseFont3D || _fontRegistry['系统默认']
  // 先查注册表
  if (_fontRegistry[fontFamily]) return _fontRegistry[fontFamily]
  // 若未加载，启动按需加载（异步），同步返回默认字体
  if (!_fontLoadPromises[fontFamily] && typeof window !== 'undefined' && window.opentype) {
    var fontPath = FONT_FAMILY_PATHS[fontFamily] || FONT_FAMILY_PATHS['系统默认']
    _fontLoadPromises[fontFamily] = fetch(fontPath)
      .then(function(res) { return res.arrayBuffer() })
      .then(function(buffer) {
        var font = window.opentype.parse(buffer)
        if (font && font.outlinesFormat) {
          _fontRegistry[fontFamily] = font
          console.log('[3DText] Font loaded on demand: ' + fontFamily + ' from ' + fontPath)
          // 字体加载完成后触发事件，让场景中所有使用该字体的文字重建几何体
          if (typeof window !== 'undefined') {
            window.dispatchEvent(new CustomEvent('ism3d-font-ready'))
          }
        }
        return font
      })
      .catch(function(err) {
        console.warn('[3DText] Failed to load font ' + fontFamily + ':', err.message)
        _fontLoadPromises[fontFamily] = null // 允许下次重试
        return null
      })
  }
  return _chineseFont3D || _fontRegistry['系统默认'] || null
}

/**
 * 创建真正的3D几何体文字（挤出几何体）
 * 需要 window.opentype 已加载，且中文字体已加载
 * @returns {THREE.Mesh|null} 成功返回Mesh，失败返回null
 */
function createTrue3DTextMesh(text, color, fontSize, options) {
  options = options || {}
  var font = getTrue3DFont(options.fontFamily)
  if (!font || !text) return null

  var fontFamilyName = options.fontFamily || '系统默认'
  console.log('[True3D] createTrue3DTextMesh: text="' + text + '" fontFamily=' + fontFamilyName + ' color=' + color + ' fontSize=' + fontSize)
  var worldHeight = (fontSize / 16) * 0.8
  var depth = worldHeight * 0.2

  // 逐个字符处理，用字形索引获取独立路径
  var glyphs = font.stringToGlyphs(text)
  var allShapes = []
  var xOffset = 0

  var scale = worldHeight / font.unitsPerEm

  for (var gi = 0; gi < glyphs.length; gi++) {
    var glyph = glyphs[gi]
    var glyphPath = glyph.getPath(xOffset, 0, font.unitsPerEm)
    var shapes = _glyphPathToShapes(glyphPath.commands, scale)
    for (var si = 0; si < shapes.length; si++) {
      allShapes.push(shapes[si])
    }
    // advanceWidth 用设计单位，按 scale 缩放到世界坐标
    xOffset += (glyph.advanceWidth || font.unitsPerEm)
  }

  if (allShapes.length === 0) {
    console.warn('[True3D] No shapes generated for text: ' + text)
    return null
  }

  // ExtrudeGeometry - bevel 参数相对于字符大小
  var geo = new THREE.ExtrudeGeometry(allShapes, {
    depth: depth,
    bevelEnabled: true,
    bevelThickness: worldHeight * 0.02,
    bevelSize: worldHeight * 0.01,
    bevelSegments: 3,
    curveSegments: 8,
    steps: 1
  })

  // 居中
  geo.computeBoundingBox()
  var box = geo.boundingBox
  geo.translate(
    -(box.min.x + box.max.x) / 2,
    -(box.min.y + box.max.y) / 2,
    -depth / 2
  )

  var col = new THREE.Color(color || '#ffffff')
  // emissive 让 true 3D 文字亮度与 Canvas 贴图方案接近
  var emissiveCol = col.clone().multiplyScalar(0.25)
  var mat = new THREE.MeshStandardMaterial({
    color: col,
    emissive: emissiveCol,
    emissiveIntensity: 0.6,
    metalness: 0.1,
    roughness: 0.4,
    side: THREE.DoubleSide,
    toneMapped: false
  })

  var mesh = new THREE.Mesh(geo, mat)
  mesh.castShadow = true
  mesh.receiveShadow = false

  mesh.userData.is3DText = true
  mesh.userData.isTrue3DText = true
  mesh.userData.isTextSprite = true
  mesh.userData.faceCamera = false
  mesh.userData.textData = {
    text: text,
    color: color,
    fontSize: fontSize,
    bgColor: '#00000000',
    options: Object.assign({}, options)
  }
  mesh.name = 'True3DText_' + text.slice(0, 10)

  return mesh
}

/**
 * 字符间距常量：advanceWidth 系数
 * 不同字体设计单位不同（1024、2048 等），这里按比例统一缩放
 */
/**
 * 单个字形路径 → Three.js Shape[]
 * 策略：按子路径（M...Z）切分，按面积排序：最大=外轮廓，其余=hole
 */
function _glyphPathToShapes(commands, scale) {
  var contours = []
  var cur = null
  var sx0 = 0, sy0 = 0
  var lx = 0, ly = 0

  function flush() {
    if (cur && cur.cmds.length > 0) {
      // 计算 bounding box 面积（用于外轮廓/内轮廓判定）
      var minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity
      for (var i = 0; i < cur.pts.length; i++) {
        var p = cur.pts[i]
        if (p.x < minX) minX = p.x
        if (p.x > maxX) maxX = p.x
        if (p.y < minY) minY = p.y
        if (p.y > maxY) maxY = p.y
      }
      cur.area = (maxX - minX) * (maxY - minY)
      cur.minX = minX; cur.minY = minY; cur.maxX = maxX; cur.maxY = maxY
      contours.push(cur)
    }
    cur = null
  }

  for (var i = 0; i < commands.length; i++) {
    var cmd = commands[i]
    switch (cmd.type) {
      case 'M':
        flush()
        cur = { pts: [], cmds: [] }
        var mx = cmd.x * scale
        var my = -cmd.y * scale  // 翻转 Y：opentype.js Y 向下，Three.js Y 向上
        cur.pts.push({ x: mx, y: my })
        cur.cmds.push({ t: 'M', x: mx, y: my })
        sx0 = mx; sy0 = my
        lx = mx; ly = my
        break
      case 'L':
        if (!cur) break
        var lnx = cmd.x * scale
        var lny = -cmd.y * scale
        cur.pts.push({ x: lnx, y: lny })
        cur.cmds.push({ t: 'L', x: lnx, y: lny })
        lx = lnx; ly = lny
        break
      case 'Q':
        if (!cur || cmd.x1 == null) break
        var qp1 = { x: cmd.x1 * scale, y: -cmd.y1 * scale }
        var qp2 = { x: cmd.x * scale, y: -cmd.y * scale }
        cur.cmds.push({ t: 'Q', cp: qp1, end: qp2 })
        // 也采样几个点供面积计算
        for (var t = 0.1; t < 1; t += 0.2) {
          var t1 = 1 - t
          cur.pts.push({
            x: t1 * t1 * lx + 2 * t1 * t * qp1.x + t * t * qp2.x,
            y: t1 * t1 * ly + 2 * t1 * t * qp1.y + t * t * qp2.y
          })
        }
        cur.pts.push({ x: qp2.x, y: qp2.y })
        lx = qp2.x; ly = qp2.y
        break
      case 'C':
        if (!cur || cmd.x1 == null || cmd.x2 == null) break
        var cp1 = { x: cmd.x1 * scale, y: -cmd.y1 * scale }
        var cp2 = { x: cmd.x2 * scale, y: -cmd.y2 * scale }
        var cp3 = { x: cmd.x * scale, y: -cmd.y * scale }
        cur.cmds.push({ t: 'C', cp1: cp1, cp2: cp2, end: cp3 })
        for (var t = 0.1; t < 1; t += 0.2) {
          var t1 = 1 - t
          var a = t1 * t1 * t1, b = 3 * t1 * t1 * t, c = 3 * t1 * t * t, d = t * t * t
          cur.pts.push({
            x: a * lx + b * cp1.x + c * cp2.x + d * cp3.x,
            y: a * ly + b * cp1.y + c * cp2.y + d * cp3.y
          })
        }
        cur.pts.push({ x: cp3.x, y: cp3.y })
        lx = cp3.x; ly = cp3.y
        break
      case 'Z':
        if (cur) {
          cur.cmds.push({ t: 'Z' })
          if (cur.pts.length > 0) {
            cur.pts.push({ x: cur.pts[0].x, y: cur.pts[0].y })
          }
        }
        lx = sx0; ly = sy0
        break
    }
  }
  flush()

  if (contours.length === 0) return []

  // 按面积从大到小排序：最大 = 外轮廓，其余 = hole
  contours.sort(function (a, b) { return b.area - a.area })

  var result = []
  var used = new Array(contours.length).fill(false)

  for (var oi = 0; oi < contours.length; oi++) {
    if (used[oi]) continue
    var outer = contours[oi]
    used[oi] = true

    var shape = new THREE.Shape()
    for (var ci = 0; ci < outer.cmds.length; ci++) {
      var c = outer.cmds[ci]
      switch (c.t) {
        case 'M': shape.moveTo(c.x, c.y); break
        case 'L': shape.lineTo(c.x, c.y); break
        case 'Q': shape.quadraticCurveTo(c.cp.x, c.cp.y, c.end.x, c.end.y); break
        case 'C': shape.bezierCurveTo(c.cp1.x, c.cp1.y, c.cp2.x, c.cp2.y, c.end.x, c.end.y); break
        case 'Z': shape.closePath(); break
      }
    }

    // 找所有面积更小的轮廓，判断其中心点是否在外轮廓内
    for (var hi = oi + 1; hi < contours.length; hi++) {
      if (used[hi]) continue
      var hole = contours[hi]
      var cx = (hole.minX + hole.maxX) / 2
      var cy = (hole.minY + hole.maxY) / 2
      if (_pointInContour({ x: cx, y: cy }, outer.pts)) {
        used[hi] = true
        var holePath = new THREE.Path()
        for (var ci = 0; ci < hole.cmds.length; ci++) {
          var c = hole.cmds[ci]
          switch (c.t) {
            case 'M': holePath.moveTo(c.x, c.y); break
            case 'L': holePath.lineTo(c.x, c.y); break
            case 'Q': holePath.quadraticCurveTo(c.cp.x, c.cp.y, c.end.x, c.end.y); break
            case 'C': holePath.bezierCurveTo(c.cp1.x, c.cp1.y, c.cp2.x, c.cp2.y, c.end.x, c.end.y); break
            case 'Z': holePath.closePath(); break
          }
        }
        shape.holes.push(holePath)
      }
    }
    result.push(shape)
  }

  return result
}

/** 射线法：判断点是否在多边形内（用于 hole 归属） */
function _pointInContour(pt, poly) {
  var x = pt.x, y = pt.y, inside = false
  for (var i = 0, j = poly.length - 1; i < poly.length; j = i++) {
    var xi = poly[i].x, yi = poly[i].y
    var xj = poly[j].x, yj = poly[j].y
    var inter = ((yi > y) !== (yj > y)) && (x < (xj - xi) * (y - yi) / (yj - yi) + xi)
    if (inter) inside = !inside
  }
  return inside
}


export {
  COMP_COLORS,
  COMPONENT_LIBRARY,
  createThreeMesh,
  create2DComponentPlane,
  createMediaPlane,
  updateMediaPlane,
  disposeMediaPlane,
  createTextSprite,
  updateTextSprite,
  createText3DMesh,
  updateText3DMesh,
  createTrue3DTextMesh,
  getTrue3DFont,
  tryLoadChineseFontFor3D,
  createFlowPipe,
  updateFlowPipeAnimation,
  updateFlowPipe,
  isFlowPipeSegmented,
  CSS2DObject
}

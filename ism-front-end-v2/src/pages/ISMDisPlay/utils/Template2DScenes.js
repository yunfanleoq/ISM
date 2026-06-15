import { createISMNode } from './ismSceneNormalizer'

const L = (cn, us, hk) => ({
  CN: cn,
  US: us,
  HK: hk || cn
})

function textNode(name, x, y, width, height, style = {}) {
  return createISMNode({
    shape: 'view-svg-text',
    name,
    x,
    y,
    width,
    height,
    zIndex: typeof style.zIndex === 'number' ? style.zIndex : 8,
    style: {
      text: style.text || name,
      ...style
    }
  })
}

function componentNode(shape, name, x, y, width, height, style = {}) {
  return createISMNode({
    shape,
    name,
    x,
    y,
    width,
    height,
    zIndex: typeof style.zIndex === 'number' ? style.zIndex : 6,
    style: {
      text: name,
      backColor: 'transparent',
      foreColor: '#d7ecff',
      fontSize: 14,
      zIndex: typeof style.zIndex === 'number' ? style.zIndex : 6,
      ...style
    }
  })
}

function panel(name, x, y, width, height, title, options = {}) {
  const zIndex = options.zIndex || 1
  const border = options.border || 'dv-border-box8'
  const titleColor = options.titleColor || '#d6e8f7'
  return [
    componentNode(border, name, x, y, width, height, { zIndex }),
    textNode(`${name}标题`, x + 18, y + 14, width - 36, 30, {
      text: title,
      foreColor: titleColor,
      fontSize: 18,
      fontWeight: 600,
      textAlign: 'left',
      zIndex: zIndex + 1
    })
  ]
}

function metric(name, x, y, width, label, value, sub, color = '#8fd8ff') {
  return textNode(name, x, y, width, 84, {
    text: `${label}\n${value}\n${sub}`,
    backColor: 'rgba(13, 35, 52, 0.72)',
    foreColor: color,
    fontSize: 17,
    fontWeight: 600,
    textAlign: 'center',
    borderWidth: 1,
    BorderEdges: 4,
    borderColor: 'rgba(120, 190, 220, 0.28)',
    zIndex: 8
  })
}

function miniMetric(name, x, y, width, label, value, color = '#8fd8ff') {
  return textNode(name, x, y, width, 52, {
    text: `${label}\n${value}`,
    backColor: 'rgba(10, 30, 46, 0.68)',
    foreColor: color,
    fontSize: 15,
    fontWeight: 600,
    textAlign: 'center',
    borderWidth: 1,
    BorderEdges: 4,
    borderColor: 'rgba(120, 190, 220, 0.22)',
    zIndex: 8
  })
}

function statusRow(name, x, y, width, label, value, color = '#7ad7a5') {
  return textNode(name, x, y, width, 34, {
    text: `${label}    ${value}`,
    backColor: 'rgba(8, 26, 40, 0.64)',
    foreColor: color,
    fontSize: 15,
    textAlign: 'left',
    borderWidth: 1,
    BorderEdges: 3,
    borderColor: 'rgba(120, 190, 220, 0.18)',
    zIndex: 8
  })
}

function label(name, x, y, width, text, color = '#9db8cc') {
  return textNode(name, x, y, width, 26, {
    text,
    foreColor: color,
    fontSize: 14,
    textAlign: 'left',
    zIndex: 8
  })
}

function tableBlock(name, x, y, width, height, title, rows, color = '#8fd8ff') {
  return textNode(name, x, y, width, height, {
    text: `${title}\n${rows.join('\n')}`,
    backColor: 'rgba(7, 24, 38, 0.64)',
    foreColor: color,
    fontSize: 15,
    textAlign: 'left',
    borderWidth: 1,
    BorderEdges: 4,
    borderColor: 'rgba(120, 190, 220, 0.26)',
    zIndex: 8
  })
}

function trendBlock(name, x, y, width, height, title, rows, color = '#8fd8ff') {
  return textNode(name, x, y, width, height, {
    text: `${title}\n${rows.join('\n')}`,
    backColor: 'rgba(7, 24, 38, 0.64)',
    foreColor: color,
    fontSize: 17,
    textAlign: 'left',
    borderWidth: 1,
    BorderEdges: 4,
    borderColor: 'rgba(120, 190, 220, 0.26)',
    zIndex: 8
  })
}

function buildFactoryDashboard() {
  const cells = []
  cells.push(textNode('页面标题', 540, 28, 840, 54, {
    text: '智慧工厂生产运营看板',
    foreColor: '#edf7ff',
    fontSize: 32,
    fontWeight: 700,
    zIndex: 12
  }))
  cells.push(componentNode('view-svg-time', '系统时间', 1468, 34, 300, 38, { foreColor: '#9db8cc', fontSize: 16, zIndex: 12 }))
  cells.push(...panel('设备概览容器', 40, 108, 380, 330, '设备概览'))
  cells.push(trendBlock('设备状态实时', 72, 166, 316, 220, '设备总览  128台', [
    '在线 121    离线 3    告警 4',
    '稼动率 94.5%    MTBF 186h',
    'CNC-08  运行  主轴负载 72%',
    'AGV-02  离线  最后心跳 2分钟前',
    '空压站  运行  压力 0.72MPa'
  ], '#8fd8ff'))
  cells.push(...panel('产线运行容器', 450, 108, 1020, 330, '产线运行'))
  cells.push(metric('产线A指标', 498, 176, 210, '产线A', '98.6%', '稳定运行', '#7ad7a5'))
  cells.push(metric('产线B指标', 736, 176, 210, '产线B', '96.2%', '稳定运行', '#7ad7a5'))
  cells.push(metric('产线C指标', 974, 176, 210, '产线C', '维护中', '#e6c06a'))
  cells.push(metric('今日产量指标', 1212, 176, 210, '今日产量', '36,920', '达成率 94%', '#8fd8ff'))
  cells.push(statusRow('工位状态1', 508, 306, 410, '冲压 / 焊装 / 涂装', '在线', '#7ad7a5'))
  cells.push(statusRow('工位状态2', 958, 306, 410, '总装 / 质检 / 包装', '待检 12 件', '#e6c06a'))
  cells.push(miniMetric('当前班次', 508, 356, 190, '当前班次', '白班 08:00-20:00', '#8fd8ff'))
  cells.push(miniMetric('订单进度', 718, 356, 190, '订单 J20260610', '74.8%', '#7ad7a5'))
  cells.push(miniMetric('返修率', 928, 356, 190, '返修率', '0.42%', '#7ad7a5'))
  cells.push(miniMetric('换型倒计时', 1138, 356, 230, '换型倒计时', '01:24:36', '#e6c06a'))
  cells.push(...panel('实时数据容器', 1500, 108, 380, 330, '实时数据'))
  cells.push(tableBlock('生产实时数据', 1530, 166, 320, 220, '点位              当前值', [
    'A线节拍           42 s',
    'A线良率           98.6%',
    'B线节拍           45 s',
    '涂装温度          68 C',
    '总装待检          12 件',
    '空压压力          0.72 MPa'
  ]))
  cells.push(...panel('趋势容器', 40, 468, 920, 460, '关键趋势'))
  cells.push(trendBlock('生产趋势模拟', 72, 528, 856, 330, '产量 / 良率 / 能耗趋势', [
    '08:00  产量 4,820   良率 98.1%   能耗 1,420 kWh',
    '10:00  产量 9,640   良率 98.4%   能耗 2,860 kWh',
    '12:00  产量 14,300  良率 97.9%   能耗 4,180 kWh',
    '14:00  产量 22,780  良率 98.6%   能耗 7,940 kWh',
    '16:00  产量 31,260  良率 98.2%   能耗 10,880 kWh',
    '趋势判断：产量高于计划 3.4%，单位能耗低于昨日 2.1%'
  ]))
  cells.push(label('趋势说明', 72, 872, 620, '静态示例数据；应用后可替换为实时曲线组件并绑定设备点位'))
  cells.push(...panel('告警容器', 990, 468, 420, 460, '告警与待办', { titleColor: '#e6c06a' }))
  cells.push(tableBlock('生产告警列表', 1020, 528, 360, 276, '等级  对象        内容', [
    '高    AGV-02     离线 2 分钟',
    '中    三号产线    烘房温度偏高',
    '中    质检站      待检 12 件',
    '低    空压机      24小时后保养'
  ], '#e6c06a'))
  cells.push(statusRow('待办1', 1020, 826, 360, 'AGV-02', '离线 2 分钟', '#e6c06a'))
  cells.push(statusRow('待办2', 1020, 868, 360, '三号产线', '温度偏高', '#e6c06a'))
  cells.push(statusRow('待办3', 1020, 910, 360, '质检站', '待检积压 12 件', '#8fd8ff'))
  cells.push(...panel('报表容器', 1440, 468, 440, 460, '生产明细'))
  cells.push(tableBlock('生产明细模拟', 1470, 528, 380, 330, '产线  计划  完成  良率  状态', [
    'A线   15000  14820  98.6%  正常',
    'B线   12000  11840  96.2%  正常',
    'C线   10000  8260   94.1%  维护',
    '包装  36000  35220  99.0%  正常',
    '质检  36000  35180  98.8%  待检'
  ]))
  cells.push(label('报表说明', 1470, 872, 340, '可配置设备、数据列、刷新间隔和表格主题'))
  return {
    layer: { width: 1920, height: 1080, backColor: '#07111d', backgroundImage: '', autoSize: 0, Padding: 1 },
    components: { cells }
  }
}

function buildPowerRoom() {
  const cells = []
  cells.push(textNode('页面标题', 600, 28, 720, 54, {
    text: '配电室运行监控看板',
    foreColor: '#edf7ff',
    fontSize: 32,
    fontWeight: 700,
    zIndex: 12
  }))
  cells.push(...panel('运行概览容器', 40, 108, 380, 820, '运行概览'))
  cells.push(metric('总负荷', 76, 178, 300, '总负荷', '1860 kW', '负载率 72%', '#8fd8ff'))
  cells.push(metric('功率因数', 76, 292, 300, '功率因数', '0.96', '补偿正常', '#7ad7a5'))
  cells.push(metric('母线电压', 76, 406, 300, '母线电压', '10.4 kV', '电压稳定', '#7ad7a5'))
  cells.push(metric('变压器温度', 76, 520, 300, '变压器温度', '68 C', '持续关注', '#e6c06a'))
  cells.push(trendBlock('配电设备状态', 76, 652, 300, 210, '设备状态', [
    '进线柜 2/2  正常',
    '变压器 2/2  运行',
    '馈线柜 18/18 合闸',
    '电容柜 1组投入  2组备用',
    'UPS 负载率 42%'
  ], '#8fd8ff'))
  cells.push(...panel('一次系统容器', 450, 108, 1020, 820, '一次系统图'))
  cells.push(componentNode('view-svg-electric1', '主变', 590, 210, 150, 150, { zIndex: 8 }))
  cells.push(componentNode('view-svg-electric2', '高压柜', 850, 210, 150, 150, { zIndex: 8 }))
  cells.push(componentNode('view-svg-electric3', '低压柜', 1110, 210, 150, 150, { zIndex: 8 }))
  cells.push(miniMetric('主变负载', 570, 370, 190, '1#主变', '负载 68%', '#8fd8ff'))
  cells.push(miniMetric('高压柜状态', 830, 370, 190, '高压柜', '合闸 / 远方', '#7ad7a5'))
  cells.push(miniMetric('低压柜状态', 1090, 370, 190, '低压柜', '电流 1240A', '#7ad7a5'))
  cells.push(componentNode('ViewCanvasMoveLineArrow', '母线一', 650, 440, 560, 50, { zIndex: 7 }))
  cells.push(componentNode('view-svg-electric4', '母联柜', 720, 520, 150, 150, { zIndex: 8 }))
  cells.push(componentNode('view-svg-electric5', '电容柜', 1000, 520, 150, 150, { zIndex: 8 }))
  cells.push(miniMetric('母联状态', 690, 680, 220, '母联柜', '热备用', '#e6c06a'))
  cells.push(miniMetric('补偿状态', 970, 680, 220, '电容补偿', '投入 2 组', '#7ad7a5'))
  cells.push(componentNode('ViewCanvasMoveLineArrow', '馈线', 800, 752, 360, 50, { zIndex: 7 }))
  cells.push(statusRow('系统说明', 540, 824, 840, '10kV I段 / II段', '母联、馈线、补偿回路可绑定实时数据', '#8fd8ff'))
  cells.push(...panel('环境容器', 1500, 108, 380, 390, '环境数据'))
  cells.push(tableBlock('配电环境数据', 1530, 168, 320, 240, '监测项          当前值', [
    '室内温度        28.6 C',
    '室内湿度        42%',
    'SF6浓度         0 ppm',
    '烟感状态        正常',
    '门禁状态        已关闭',
    '水浸状态        正常'
  ]))
  cells.push(...panel('告警容器', 1500, 528, 380, 400, '告警记录', { titleColor: '#e6c06a' }))
  cells.push(tableBlock('配电告警记录', 1530, 588, 320, 250, '等级  回路        内容', [
    '中    1#变压器   温度偏高',
    '低    母联柜     通信延迟',
    '低    电容柜     投切次数偏多',
    '提示  UPS        电池自检通过'
  ], '#e6c06a'))
  cells.push(statusRow('告警摘要1', 1530, 858, 320, '变压器温度', '偏高', '#e6c06a'))
  cells.push(statusRow('巡检摘要', 1530, 900, 320, '今日巡检', '12/12 已完成', '#7ad7a5'))
  return {
    layer: { width: 1920, height: 1080, backColor: '#08111b', backgroundImage: '', autoSize: 0, Padding: 1 },
    components: { cells }
  }
}

function buildEnergyOverview() {
  const cells = []
  cells.push(textNode('页面标题', 560, 28, 800, 54, {
    text: '园区能源管理驾驶舱',
    foreColor: '#edf7ff',
    fontSize: 32,
    fontWeight: 700,
    zIndex: 12
  }))
  cells.push(...panel('能源概览容器', 40, 108, 520, 360, '能源概览'))
  cells.push(metric('电耗指标', 82, 178, 140, '电耗', '8,620', 'kWh', '#8fd8ff'))
  cells.push(metric('水耗指标', 242, 178, 140, '水耗', '326', 'm3', '#7ad7a5'))
  cells.push(metric('气耗指标', 402, 178, 120, '气耗', '1,082', 'Nm3', '#e6c06a'))
  cells.push(trendBlock('节能目标进度', 82, 306, 440, 92, '月度节能目标', [
    '目标  38,000 kWh    已节约  26,420 kWh    完成率 69.5%'
  ], '#7ad7a5'))
  cells.push(miniMetric('峰时占比', 82, 408, 140, '峰时占比', '36%', '#e6c06a'))
  cells.push(miniMetric('谷时占比', 242, 408, 140, '谷时占比', '44%', '#7ad7a5'))
  cells.push(miniMetric('光伏抵扣', 402, 408, 120, '光伏抵扣', '5.6%', '#8fd8ff'))
  cells.push(...panel('分时趋势容器', 590, 108, 880, 360, '分时趋势'))
  cells.push(trendBlock('能源分时趋势', 620, 168, 820, 238, '分时能耗趋势', [
    '00:00  电 310 kWh  水 11 m3  气 46 Nm3',
    '06:00  电 420 kWh  水 16 m3  气 58 Nm3',
    '12:00  电 760 kWh  水 31 m3  气 96 Nm3',
    '18:00  电 690 kWh  水 28 m3  气 82 Nm3',
    '峰值 12:00，建议冷机错峰提前 30 分钟启动'
  ]))
  cells.push(label('趋势说明', 620, 416, 760, '静态示例数据；应用后可替换为历史趋势组件并配置时间范围'))
  cells.push(...panel('能源排行容器', 1500, 108, 380, 360, '能耗排行'))
  cells.push(tableBlock('能耗排行模拟', 1530, 168, 320, 238, '区域        今日能耗   环比', [
    '生产楼      3,860 kWh  +4.2%',
    '动力站      2,140 kWh  -1.8%',
    '研发楼      1,260 kWh  +0.6%',
    '仓储区      760 kWh    -3.1%',
    '办公楼      604 kWh    -2.4%'
  ]))
  cells.push(statusRow('排行结论', 1530, 416, 320, '重点区域', '生产楼 / 动力站', '#e6c06a'))
  cells.push(...panel('实时数据容器', 40, 498, 520, 430, '实时数据'))
  cells.push(tableBlock('能源实时数据模拟', 72, 558, 456, 286, '表计            当前值', [
    '总电表          8620 kWh',
    '冷站电表        2160 kWh',
    '空压站电表      940 kWh',
    '园区水表        326 m3',
    '燃气总表        1082 Nm3',
    '光伏发电        486 kWh'
  ]))
  cells.push(statusRow('实时数据说明', 72, 864, 456, '数据源', '电表、水表、燃气表', '#8fd8ff'))
  cells.push(...panel('分析图表容器', 590, 498, 880, 430, '指标分析'))
  cells.push(trendBlock('能源指标分析模拟', 620, 558, 820, 286, '指标分析', [
    '单位面积能耗  42.6 kWh/m2    同比 -3.8%',
    '单位产值能耗  0.86 kWh/万元   同比 -5.1%',
    '碳排估算      4.82 tCO2       环比 -2.7%',
    '异常点位      3 个            已派单 2 个',
    '节能空间      冷站群控 / 空压机联控 / 照明时段'
  ]))
  cells.push(statusRow('分析建议', 620, 864, 820, '建议', '错峰启动冷机，预计节能 6.8%', '#7ad7a5'))
  cells.push(...panel('告警建议容器', 1500, 498, 380, 430, '告警与建议', { titleColor: '#e6c06a' }))
  cells.push(tableBlock('能源告警模拟', 1530, 558, 320, 216, '等级  对象        内容', [
    '中    冷站        负载偏高',
    '中    空压站      夜间能耗异常',
    '低    办公楼      照明未关闭',
    '提示  光伏        发电低于预测'
  ], '#e6c06a'))
  cells.push(statusRow('建议1', 1530, 802, 320, '空调负载', '偏高', '#e6c06a'))
  cells.push(statusRow('建议2', 1530, 844, 320, '峰谷策略', '建议启用', '#7ad7a5'))
  cells.push(statusRow('建议3', 1530, 886, 320, '照明策略', '办公楼延时关闭', '#8fd8ff'))
  return {
    layer: { width: 1920, height: 1080, backColor: '#07131a', backgroundImage: '', autoSize: 0, Padding: 1 },
    components: { cells }
  }
}

export const TEMPLATE_2D_SCENES = [
  {
    id: 'factory-dashboard',
    name: L('智慧工厂监控', 'Smart Factory Dashboard', '智慧工廠監控'),
    category: L('工业监控', 'Industrial', '工業監控'),
    description: L('包含设备状态、产线指标、实时趋势、告警列表和生产明细的成品页面骨架', 'Production page skeleton with device status, line KPIs, realtime trend, alarms, and details.', '包含設備狀態、產線指標、即時趨勢、告警列表和生產明細的成品頁面骨架'),
    icon: 'fas fa-industry',
    build: buildFactoryDashboard
  },
  {
    id: 'power-room',
    name: L('配电室综合监控', 'Power Room Monitoring', '配電室綜合監控'),
    category: L('能源电力', 'Energy', '能源電力'),
    description: L('包含运行指标、一次系统图、电气元件、环境数据和告警记录', 'Power-room page with KPIs, one-line diagram, electrical elements, environment data, and alarms.', '包含運行指標、一次系統圖、電氣元件、環境資料和告警記錄'),
    icon: 'fas fa-bolt',
    build: buildPowerRoom
  },
  {
    id: 'energy-overview',
    name: L('园区能源驾驶舱', 'Energy Overview', '園區能源駕駛艙'),
    category: L('能源分析', 'Analytics', '能源分析'),
    description: L('包含分项能耗、历史趋势、实时数据、指标分析和节能建议', 'Energy page with category KPIs, history trend, realtime data, analytics, and saving advice.', '包含分項能耗、歷史趨勢、即時資料、指標分析和節能建議'),
    icon: 'fas fa-chart-line',
    build: buildEnergyOverview
  }
]

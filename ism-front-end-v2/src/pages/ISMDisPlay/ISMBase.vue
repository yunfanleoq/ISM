

<script>

// const context = require.context('./', true, /\.vue$/);
// const install = (Vue) => {
//   context.keys().forEach((key) => {
//     const component = context(key).default;
//     Vue.component(component.name, component);
//   });
// };

import {mapMutations} from "vuex";
import { register } from '@antv/x6-vue-shape'
const toolStandardBoxList = {
  title: "displayConfig.ToolBox.Base.title",
  icon: "icon-standard-application",
  opened: false,
  items:[]
}
const toolVideoBoxList = {
  title: "displayConfig.ToolBox.Video.title",
  icon: "icon-standard-application",
  opened: false,
  items:[]
}
const toolLoginBoxList = {
  title: "displayConfig.ToolBox.login.title",
  icon: "icon-loginBox",
  opened: false,
  items:[]
}

const toolCanvasBoxList = {
  title: "displayConfig.ToolBox.Diagram.title",
  icon: "icon-standard-picture-empty",
  opened: false,
  items:[]
}

const toolMapBoxList = {
  title: "displayConfig.ToolBox.Map",
  icon: "icon-map",
  opened: false,
  items:[]
}
const toolHistoryChartsBoxList = {
  title: "displayConfig.ToolBox.HistoryCharts",
  icon: "icon-standard-chart-curve",
  opened: false,
  items:[]
}

const toolChartsBoxList = {
  title: "displayConfig.ToolBox.Charts.title",
  icon: "icon-standard-chart-bar",
  opened: false,
  items:[]
}

const MesStandardList={
  title: "displayConfig.ToolBox.MesStandard.title",
  icon: "icon-standard-chart-bar",
  opened: false,
  items:[]
}



const toolSvgArrowsBoxList={
  title: "displayConfig.ToolBox.Arrows",
  "icon": "icon-standard-arrows",
  "opened": false,
  "isSequence": true,
  items:[]
}
const toolSvgElectricBoxList={
  title: "displayConfig.ToolBox.Electric",
  "icon": "icon-electric",
  "opened": false,
  "isSequence": true,
  vueCount:0,
  items:[]
}
const toolBigScreenBoxContainerList={
  title: "displayConfig.ToolBox.bigScreen.Container",
  icon: "icon-standard-big-screen",
  opened: false,
  items:[]
}

const toolDeviceContainerList={
  title: "displayConfig.ToolBox.device.Container",
  icon: "icon-device",
  opened: false,
  items:[]
}

let res_components = {}
const componentsStandard = require.context('./ISMComponents/standard/', true, /\.vue$/)
const componentsVideo = require.context('./ISMComponents/video/', true, /\.vue$/)
const componentsLogin = require.context('./ISMComponents/login/', true, /\.vue$/)
const componentsCanvas = require.context('./ISMComponents/canvas/', true, /\.vue$/)
const componentsCharts = require.context('./ISMComponents/charts/', true, /\.vue$/)
const componentsBigScreen = require.context('./ISMComponents/bigScreen/', true, /\.vue$/)
const componentsSvgArrows = require.context('./ISMComponents/svg/arrows/', true, /\.vue$/)
const componentsSvgElectric = require.context('./ISMComponents/ComponentClassification/electric/', true, /\.vue$/)
const componentsImages = require.context('./ISMComponents/Images/', true, /\.vue$/)
const componentPiping = require.context('../../../public/static/ISM/Conduit', true, /\.png$/)
const componentPageNavigationList = require.context('../../../public/static/ISM/page/navigation/', true, /\.png$/)
const componentPageContainerList = require.context('../../../public/static/ISM/page/container/', true, /\.png$/)
const componentPageDecorateList = require.context('../../../public/static/ISM/page/decorate/', true, /\.png$/)
const componentPageBackgroundBarList = require.context('../../../public/static/ISM/page/BackgroundBar/', true, /\.png$/)

const componentPageScreenedging1List = require.context('../../../public/static/ISM/page/screenedging1/', true, /\.png$/)
const componentPageScreenedging2List = require.context('../../../public/static/ISM/page/screenedging2/', true, /\.svg$/)
const componentPageScreenedging3List = require.context('../../../public/static/ISM/page/screenedging3/', true, /\.png$/)

const componentPage3DIconList = require.context('../../../public/static/ISM/page/3DIcon/', true, /\.svg$/)
const componentPageBigpicList = require.context('../../../public/static/ISM/page/bigpic/', true, /\.png$/)
const componentPageLlustrationList = require.context('../../../public/static/ISM/page/Illustration/', true, /\.png$/)

//图片
const componentHVACList = require.context('../../../public/static/ISM/HVAC/', true, /\.png$/)
const componentElectricMachineryList = require.context('../../../public/static/ISM/ElectricMachinery/', true, /\.png$/)
const componentFanList = require.context('../../../public/static/ISM/Fan/', true, /\.png$/)
const componentMercuryList = require.context('../../../public/static/ISM/Mercury/', true, /\.png$/)
const componentBlenderList = require.context('../../../public/static/ISM/Blender/', true, /\.png$/)

const componentElectricPngList = require.context('../../../public/static/ISM/systemImage/electric', true, /\.svg$/)

const deviceComponents = require.context('./ISMComponents/device/', true, /\.vue$/)
const mapComponents = require.context('./ISMComponents/map/', true, /\.vue$/)

const historyChartsComponents = require.context('./ISMComponents/historyCharts/', true, /\.vue$/)
const MesStandardComponentsDir = require.context('./ISMComponents/Mes/standard/', true, /\.vue$/)

// 防御性读取组态组件的 base 样式。
// require.context 批量扫描注册组态组件时，若个别组件因循环依赖 / 导出残缺，
// 使得 comp.default.data 不是函数，原本 `comp.default.data().base` 会直接抛
// TypeError，进而让整个 forEach 中断、ISMBase 模块加载失败，最终路由异步组件
// 解析失败（"comp.default.data is not a function" 运行时回归）。
// 这里对单个组件 try/catch 跳过，保证一个坏组件不会拖垮整张大屏的加载。
function safeBaseOf(comp, filePath) {
  try {
    if (comp && comp.default && typeof comp.default.data === 'function') {
      return comp.default.data().base
    }
    console.error('[ISMBase] 跳过无效组态组件(缺少 data 函数): ' + filePath)
  } catch (e) {
    console.error('[ISMBase] 组态组件扫描异常已跳过: ' + filePath + ' -> ' + (e && e.message))
  }
  return undefined
}
//标准控件
componentsStandard.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsStandard(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolStandardBoxList.items.push(componentsInfo)
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})

//视频控件
componentsVideo.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsVideo(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolVideoBoxList.items.push(componentsInfo)
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})
//登录元素
componentsLogin.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsLogin(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolLoginBoxList.items.push(componentsInfo)
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})
// Canvas控件
componentsCanvas.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsCanvas(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolCanvasBoxList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else{
    console.error(filePath+"缺少base基本样式")
  }

})
// 图表控件
componentsCharts.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsCharts(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolChartsBoxList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})
// Svg组件
componentsSvgArrows.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsSvgArrows(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolSvgArrowsBoxList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }
})
//电力
toolSvgElectricBoxList.vueCount=0
componentsSvgElectric.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsSvgElectric(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolSvgElectricBoxList.items.push(componentsInfo)
    toolSvgElectricBoxList.vueCount++
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }
})
// PNG图片组件
componentsImages.keys().forEach(filePath => {
  let comp = componentsImages(filePath)
  register({
    shape: comp.default.name,
    component: comp.default,
  })
  res_components[comp.default.name] = comp.default
})
//大屏容器组件
componentsBigScreen.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = componentsBigScreen(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolBigScreenBoxContainerList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})


deviceComponents.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = deviceComponents(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolDeviceContainerList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})

MesStandardComponentsDir.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = MesStandardComponentsDir(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    MesStandardList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})

historyChartsComponents.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = historyChartsComponents(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolHistoryChartsBoxList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})

mapComponents.keys().forEach(filePath => {
  const keyArr = filePath.split('/')
  const fileName = keyArr.pop()
  const compKey = fileName.replace(/\.vue$/g, '')
  let comp = mapComponents(filePath)
  let componentsInfo = safeBaseOf(comp, filePath)
  if(typeof componentsInfo!="undefined")
  {
    register({
      shape: comp.default.name,
      component: comp.default,
    })
    componentsInfo.info.type=comp.default.name
    toolMapBoxList.items.push(componentsInfo)
    res_components[componentsInfo.info.type] = comp.default
  }
  else
  {
    console.error(filePath+"缺少base基本样式")
  }

})
//导入管道的配置文件
import Conduit from './ISMComponents/Conduit/config.json';
import PageNavigation from './ISMComponents/page/navigation.json';
import {GetUserCustomPel} from "@/services/system";
import {setDiyComponentsList} from "@/store/ISM/mutations";
export default {
  name: 'ISMBase',
  components: res_components,
  data() {
    return {
      toolBoxList:[]
    }
  },
  methods: {
    ...mapMutations('ISMDisPlayEditorTool', [
      'setToolBoxList',
      'setPageCanVasList',
      'setMesComponentsList',
      'setDiyComponentsList'
    ]),
    GetUserComponents(){


    },
    GetUserCustomPel(){
      let _t = this
      GetUserCustomPel().then(function (res){
        if(res.data.code==0)
        {
          const ISMDiyComponentsList = []
          if((res.data.list!=null)&&(res.data.list.length>0))
          {
            for(let i=0;i<res.data.list.length;i++)
            {
              if((res.data.list[i].FilePath!=null)&&(res.data.list[i].FilePath.length>0))
              {
                let toolBigScreenCustomPel={
                  title: res.data.list[i].DirName,
                  icon: "icon-custom",
                  opened: false,
                  items:[]
                }
                for(let k=0;k<res.data.list[i].FilePath.length;k++)
                {
                  let showPngVue = {
                    "text": "configComponent.image.Text",
                    "icon": "icon-xitongshezhi_tuxiang-copy",
                    "isFontIcon": true,
                    "info": {
                      "type": "ism-view-png-image",
                      "action": [],
                      "dataBind":[],
                      "animate": {
                        "selected": [],
                        "condition":{
                          deviceSN:"",
                          selectVideoType:0,
                          isBandDevice:false,
                          bandType:1,
                          dataID: "",
                          dataName: "",
                          operator:"",
                          OperatorValue:"",
                          OperatorMaxValue:"",
                        },
                        "isExpression": false,
                        "animateList": [
                          {
                            id: "Forbidden",
                            name: "component.public.Forbidden",
                          },
                          {
                            id: "blink",
                            name: "component.public.animateBlink",
                          },
                          {
                            id: "Zoom",
                            name: "component.public.Zoom",
                          },
                          {
                            id: "animateSpin",
                            name: "component.public.animateSpin",
                          },
                        ],
                        "animateElement": [
                          {
                            id: "blink",
                            elementList:[
                              {
                                "name":"component.public.animateSpeed",
                                "type":7,
                                "value":1,
                                "min":0.1,
                                "key":"blinkSpeed",
                              },
                            ]
                          },
                          {
                            id: "millcolorGrad",
                            elementList:[
                              {
                                "name": "component.public.startColor",
                                "type": 2,
                                "value": "#74f808",
                                "key": "startColor",
                              },
                              {
                                "name": "component.public.stopColor",
                                "type": 2,
                                "value": "#f30b0b",
                                "key": "stopColor",
                              },
                              {
                                "name":"component.public.animateSpeed",
                                "type":7,
                                "value":1,
                                "min":0.1,
                                "key":"animateSpeed",
                              },
                            ]
                          },
                          {
                            id: "animateSpin",
                            elementList:[
                              {
                                "name":"component.public.animateSpinSpeed",
                                "type":7,
                                "value":1,
                                "min":0.1,
                                "key":"spinSpeed",
                              },
                              {
                                name:"configComponent.bigScreen.border.border89Direction",
                                type:6,
                                value:0,
                                enumList:[
                                  {
                                    value:0,
                                    option:"configComponent.bigScreen.border.border89DirectionForward"
                                  },
                                  {
                                    value:1,
                                    option:"configComponent.bigScreen.border.border89DirectionNegative"
                                  }
                                ],
                                min:1,
                                key:"spinDirection",
                              }
                            ]
                          },
                        ],
                      },
                      "style": {
                        "position": {
                          "x": 0,
                          "y": 0,
                          "w": 32,
                          "h": 32
                        },
                        "visible":1,
                        "backColor": "transparent",
                        "zIndex": -1,
                        "transform": 0,
                        imageURL:"",
                        "diy":[

                        ]
                      }
                    }
                  }
                  showPngVue.text = ""
                  showPngVue.icon = res.data.list[i].FilePath[k]
                  showPngVue.info.style.imageURL =  res.data.list[i].FilePath[k]
                  toolBigScreenCustomPel.items.push(showPngVue)
                }
                ISMDiyComponentsList.push(toolBigScreenCustomPel)
              }
            }
            _t.setDiyComponentsList(ISMDiyComponentsList);
          }
        }

      })
    }
  },
  created() {

    const toolBigScreenBoxDecorateList={
      title: "displayConfig.ToolBox.bigScreen.Decorate",
      icon: "icon-Decorate",
      opened: false,
      items:[]
    }
    const toolBigScreenBoxBackgroundBarList={
      title: "displayConfig.ToolBox.bigScreen.BackgroundBar",
      icon: "icon-BackgroundBar",
      opened: false,
      items:[]
    }
    const toolBigScreenBoxScreenedging1List={
      title: "displayConfig.ToolBox.bigScreen.Edging1",
      icon: "icon-edding",
      opened: false,
      items:[]
    }
    const toolBigScreenBoxScreenedging2List={
      title: "displayConfig.ToolBox.bigScreen.Edging2",
      icon: "icon-edding",
      opened: false,
      items:[]
    }
    const toolBigScreenBoxScreenedging3List={
      title: "displayConfig.ToolBox.bigScreen.Edging3",
      icon: "icon-edding",
      opened: false,
      items:[]
    }
    const toolBigScreenBox3DIconList={
      title: "displayConfig.ToolBox.bigScreen.3D_MODEL",
      icon: "icon-3DMODEL",
      opened: false,
      items:[]
    }
    const toolBigScreenBoxTeIconList={
      title: "displayConfig.ToolBox.bigScreen.TeIcon",
      icon: "icon-TeIcon",
      opened: false,
      items:[]
    }
    const toolBigScreenBoxllustrationList={
      title: "displayConfig.ToolBox.bigScreen.Illustration",
      icon: "icon-Illustration",
      opened: false,
      items:[]
    }


    const toolBoxList = []
    const PageCanVasList = []
    const ISMMesComponentsList = []
    const ISMDiyComponentsList = []
    toolSvgElectricBoxList.items.splice(toolSvgElectricBoxList.vueCount,toolSvgElectricBoxList.items.length-toolSvgElectricBoxList.vueCount)
    componentElectricPngList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 160,
              "h": 160
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[
            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/systemImage/electric/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/systemImage/electric/"+fileName
      toolSvgElectricBoxList.items.push(showPngVue)
    })

    toolBoxList.push(toolStandardBoxList)
    toolBoxList.push(toolVideoBoxList)
    toolBoxList.push(toolLoginBoxList)
    toolBoxList.push(toolDeviceContainerList)
    toolBoxList.push(toolCanvasBoxList)
    toolBoxList.push(toolChartsBoxList)
    toolBoxList.push(toolHistoryChartsBoxList)
    toolBoxList.push(toolMapBoxList)

    toolBoxList.push(toolSvgArrowsBoxList)
    toolBoxList.push(toolSvgElectricBoxList)

    //==================管道开始========================
    const toolConduitBoxList={
      title: Conduit.groupTitle,
      "icon": Conduit.icon,
      "opened": Conduit.opened,
      "isSequence": Conduit.isSequence,
      items:[]
    }

    componentPiping.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 32,
              "h": 32
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[
            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/Conduit/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/Conduit/"+fileName
      toolConduitBoxList.items.push(showPngVue)
    })
    toolBoxList.push(toolConduitBoxList);
    // ==================管道结束========================

    //==================页面元素导航开始========================
    const toolPageNavigationList={
      title: PageNavigation.groupTitle,
      "icon": PageNavigation.icon,
      "opened": PageNavigation.opened,
      "isSequence": PageNavigation.isSequence,
      items:[]
    }

    componentPageNavigationList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 900,
              "h": 100
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[
            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/navigation/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/navigation/"+fileName

      toolPageNavigationList.items.push(showPngVue)
    })
    PageCanVasList.push(toolPageNavigationList);
    // ==================页面元素导航开始========================

    //==================页面元素容器开始========================
    // 大屏幕展示
    componentPageContainerList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[
            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/container/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/container/"+fileName
      toolBigScreenBoxContainerList.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxContainerList)
    // ==================页面元素导航开始========================

    //==================页面元素装饰开始========================
    componentPageDecorateList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[
            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/decorate/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/decorate/"+fileName
      toolBigScreenBoxDecorateList.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxDecorateList)
    // ==================页面元素装饰开始========================

    //==================页面元素背景开始========================
    componentPageBackgroundBarList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 240,
              "h": 50
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/BackgroundBar/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/BackgroundBar/"+fileName
      toolBigScreenBoxBackgroundBarList.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxBackgroundBarList)

    componentPageScreenedging1List.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[
            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/screenedging1/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/screenedging1/"+fileName
      toolBigScreenBoxScreenedging1List.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxScreenedging1List)

    componentPageScreenedging2List.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/screenedging2/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/screenedging2/"+fileName
      toolBigScreenBoxScreenedging2List.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxScreenedging2List)

    componentPageScreenedging3List.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/screenedging3/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/screenedging3/"+fileName
      toolBigScreenBoxScreenedging3List.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxScreenedging3List)

    componentPage3DIconList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-3D_MODEL",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/3DIcon/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/3DIcon/"+fileName
      toolBigScreenBox3DIconList.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBox3DIconList)

    componentPageBigpicList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-3D_MODEL",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/bigpic/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/bigpic/"+fileName
      toolBigScreenBoxTeIconList.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxTeIconList)
    componentPageLlustrationList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-3D_MODEL",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 300,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/page/Illustration/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/page/Illustration/"+fileName
      toolBigScreenBoxllustrationList.items.push(showPngVue)
    })
    PageCanVasList.push(toolBigScreenBoxllustrationList)

    // ==================页面元素背景开始========================

    //==================HVAC========================
    const toolBoxHVACList={
      title: "displayConfig.ToolBox.HVAC",
      icon: "icon-hvac",
      opened: false,
      items:[]
    }
    componentHVACList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 163,
              "h": 226
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/HVAC/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/HVAC/"+fileName
      toolBoxHVACList.items.push(showPngVue)
    })
    toolBoxList.push(toolBoxHVACList)

    //==================电机========================
    const toolBoxElectricMachineryList={
      title: "displayConfig.ToolBox.ElectricMachinery",
      icon: "icon-dianji",
      opened: false,
      items:[]
    }
    componentElectricMachineryList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 160,
              "h": 160
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/ElectricMachinery/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/ElectricMachinery/"+fileName
      toolBoxElectricMachineryList.items.push(showPngVue)
    })
    toolBoxList.push(toolBoxElectricMachineryList)

    //==================风机========================
    const toolBoxFanList={
      title: "displayConfig.ToolBox.Fan",
      icon: "icon-fan",
      opened: false,
      items:[]
    }
    componentFanList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 160,
              "h": 160
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/Fan/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/Fan/"+fileName
      toolBoxFanList.items.push(showPngVue)
    })
    toolBoxList.push(toolBoxFanList)

    //==================汞========================
    const toolBoxMercuryList={
      title: "displayConfig.ToolBox.Mercury",
      icon: "icon-gong",
      opened: false,
      items:[]
    }
    componentMercuryList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 160,
              "h": 160
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[

            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/Mercury/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/Mercury/"+fileName
      toolBoxMercuryList.items.push(showPngVue)
    })
    toolBoxList.push(toolBoxMercuryList)
    //==================搅拌机========================
    const toolBoxBlenderList={
      title: "displayConfig.ToolBox.Blender",
      icon: "icon-blender",
      opened: false,
      items:[]
    }
    componentBlenderList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      let showPngVue = {
        "text": "configComponent.image.Text",
        "icon": "icon-xitongshezhi_tuxiang-copy",
        "isFontIcon": true,
        "info": {
          "type": "ism-view-png-image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 160,
              "h": 160
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            imageURL:"",
            "diy":[
            ]
          }
        }
      }
      showPngVue.text = ""
      showPngVue.icon = "/static/ISM/Blender/"+fileName
      showPngVue.info.style.imageURL = "/static/ISM/Blender/"+fileName
      toolBoxBlenderList.items.push(showPngVue)
    })
    toolBoxList.push(toolBoxBlenderList)

    //==================ISM Mes系统标准组件========================
    ISMMesComponentsList.push(MesStandardList);

    this.toolBoxList = toolBoxList
    this.GetUserCustomPel()
    this.setToolBoxList(toolBoxList);
    this.setPageCanVasList(PageCanVasList);
    this.setMesComponentsList(ISMMesComponentsList)
  }
}
</script>

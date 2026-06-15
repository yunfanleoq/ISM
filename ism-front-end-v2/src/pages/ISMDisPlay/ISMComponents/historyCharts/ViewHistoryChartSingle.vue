<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject :style="styleVar" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
              <a-card :class="['history-chart-single-card', controlThemeClass]" :style="CardStyle" :bodyStyle="CardBodyStyle" :bordered="false"  :ref="detail.identifier+'card'">
                  <div >
                    <a-form layout="horizontal">
                      <div :class="advanced ? null: 'fold'" >
                        <a-row >
                          <a-col :md="8" :sm="24" >
                            <a-form-item
                                :label="$t('reporting.AlarmHistory.DeviceList')"
                                :labelCol= "{ span: 5 }"
                                :wrapperCol="{span: 18, offset: 1}"
                            >
                              <a-tree-select
                                  show-search
                                  tree-node-filter-prop="title"
                                  v-model="SelectDevice"
                                  :style="{'width': '100%','height':'32px','overflow': 'auto'}"
                                  allow-clear
                                  @change="SelectTreeDevice"
                                  :getPopupContainer="getPopupContainer"
                                  :dropdownClassName="dropdownThemeClass"
                                  :dropdown-style="dropdownStyle"
                                  :tree-data="deviceTreeData"
                                  :replace-fields="{ value: 'key',title:'text'}"
                                  :placeholder="$t('reporting.AlarmHistory.DeviceList')"
                                  tree-default-expand-all
                              >
                              </a-tree-select>
                            </a-form-item>
                          </a-col>
                          <a-col :md="8" :sm="24" >
                            <a-form-item
                                :label="$t('reporting.AlarmHistory.DataList')"
                                :labelCol="{span: 5}"
                                :wrapperCol="{span: 18, offset: 1}"
                            >
                              <a-select
                                  :style="{'width': '100%','height':'32px','overflow': 'auto'}"
                                  @change="DeviceDataCharge"
                                  @dropdownVisibleChange="GetDeviceModelDataList"
                                  @popupScroll="handlePopupScroll"
                                  @search="handleSearch" allowClear show-search optionFilterProp="children" mode="multiple" style="width: 100%" :token-separators="[',']" v-model="SelectAlarmData" :dropdownClassName="dropdownThemeClass" :dropdown-style="dropdownStyle" :getPopupContainer="getPopupContainer">
                                <a-select-option v-for="(alarmItem,itemIndex) in frontDataZ" :key="itemIndex" :value=alarmItem.uuid>
                                  {{ $t(alarmItem.name) }}
                                </a-select-option>
                              </a-select>

                            </a-form-item>
                          </a-col>
                          <a-col :md="8" :sm="24" >
                            <a-form-item
                                :label="$t('reporting.AlarmHistory.DateType')"
                                :labelCol="{span: 5}"
                                :wrapperCol="{span: 18, offset: 1}"
                            >
                              <a-radio-group v-model="SelectDateType" @change="chargeDateType">
                                <a-radio-button value="Day">
                                  {{$t('reporting.AlarmHistory.DateDayType')}}
                                </a-radio-button>
                                <a-radio-button value="Weekly">
                                  {{$t('reporting.AlarmHistory.DateWeeklyType')}}
                                </a-radio-button>
                                <a-radio-button value="Month">
                                  {{$t('reporting.AlarmHistory.DateMonthType')}}
                                </a-radio-button>
                                <a-radio-button value="Diy">
                                  {{$t('reporting.AlarmHistory.DateDiyType')}}
                                </a-radio-button>
                              </a-radio-group>
                            </a-form-item>
                          </a-col>
                        </a-row>
                        <a-row >
                          <a-col :md="8" :sm="24" >
                            <a-form-item
                                :label="$t('reporting.AlarmHistory.SelectDate')"
                                :labelCol="{span: 5}"
                                :wrapperCol="{span: 18, offset: 1}"
                            >
                              <a-date-picker :defaultValue="moment()" style="width: 100%" @change="onDateChange"  size="default" :dropdownClassName="dropdownThemeClass" :getCalendarContainer="getPopupContainer" :placeholder="$t('reporting.AlarmHistory.DateDayType')" v-if="SelectDateType=='Day'"/>
                              <a-month-picker :defaultValue="moment()" style="width: 100%" @change="onDateChange" size="default" :dropdownClassName="dropdownThemeClass" :getCalendarContainer="getPopupContainer" :placeholder="$t('reporting.AlarmHistory.DateMonthType')" v-if="SelectDateType=='Month'"/>
                              <a-week-picker :defaultValue="moment()" style="width: 100%" @change="onWeeklyDateChange" size="default" :dropdownClassName="dropdownThemeClass" :getCalendarContainer="getPopupContainer" :placeholder="$t('reporting.AlarmHistory.DateWeeklyType')" v-if="SelectDateType=='Weekly'"/>
                              <a-range-picker :default-value="[moment().add(-1, 'day'),moment()]" :showTime="true" @change="onDateChange" size="default" :dropdownClassName="dropdownThemeClass" :getCalendarContainer="getPopupContainer" v-if="SelectDateType=='Diy'"/>
                            </a-form-item>
                          </a-col>
                          <a-col :md="3" :sm="24" v-if="SelectDateType=='Diy'">
                            <span style="float: right; margin-top: 3px;">
                              <a-button class="history-chart-single-query" :disabled="messageShowLoad" @click="QueryHistoryDataList">{{$t('reporting.AlarmHistory.Query')}}</a-button>
                            </span>
                          </a-col>
                          <a-col v-else :md="2" :sm="24" >
                            <span style="float: right; margin-top: 3px;">
                              <a-button class="history-chart-single-query" :disabled="messageShowLoad" @click="QueryHistoryDataList">{{$t('reporting.AlarmHistory.Query')}}</a-button>
                            </span>
                          </a-col>
                        </a-row>
                      </div>
                    </a-form>
                  </div>
              </a-card>
                <div
                  class="view-chart-real-data"
                  :ref="detail.identifier"
                  :style="{'width': detail.style.position.w + 'px', 'height': EchartsHeight + 'px'}"
                ></div>

      </foreignObject>
      <!--      闪烁-->
      <animate v-if="isStart&&animateType.includes('blink')&&!IsToolBox" attributeName="opacity"
               values="0.1;1;0.1" :dur="blinkSpeed+'s'"
               repeatCount="indefinite"/>
      <!--渐变-->
      <animate v-if="isStart&&animateType.includes('millcolorGrad')&&!IsToolBox" attributeName="fill"
               :values="startColor+';'+stopColor+';'+startColor" :dur="animateSpeed+'s'"
               repeatCount="indefinite"/>
      <!--缩放      -->
      <animateTransform v-if="isStart&&animateType.includes('Zoom')&&!IsToolBox" attributeName="transform"   begin="0s" dur="0.6s" type="scale" values="0.9;1;0.9" repeatCount="indefinite"/>
      <!--      顺时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="0 0 0" to="360 0 0" repeatCount="indefinite" />
      <!--      逆时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="360 0 0" to="0 0 0" repeatCount="indefinite" />
  </g>
</svg>

</template>

<script>
import * as echarts from 'echarts'
import BaseView from '../View';
import moment from 'moment/moment'
import chalk from '@/theme/echarts/chalk'
import essos from '@/theme/echarts/essos'
import dark from '@/theme/echarts/dark'
import infographic from '@/theme/echarts/infographic'
import macarons from '@/theme/echarts/macarons'
import roma from '@/theme/echarts/roma'
import shine from '@/theme/echarts/shine'
import vintage from '@/theme/echarts/vintage'
import purplePassion from '@/theme/echarts/purple-passion'
import walden from '@/theme/echarts/walden'
import westeros from '@/theme/echarts/westeros'
import wonderland from '@/theme/echarts/wonderland'
import {GetChartDataHistoryList, GetDataHistoryList} from "@/services/report";

import {getMonitorTree} from "@/services/device"
import {GetDeviceModelDataList} from "@/services/device"
import 'moment/locale/zh-cn';
import  'moment/locale/en-ie';
import  'moment/locale/zh-tw';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-history-chart-single',
  i18n: require('@/i18n/language'),
  inject: ['getNode'],
  props: {

  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        this.syncEchartsThemeFromDetail(newVal)
        if(this.editMode){
          this.initComponents(newVal)
        }
        this.onResize()
      },
      deep: true
    },
    SelectDevice: {
      handler(newVal, oldVal) {
        let that = this
        setTimeout(function (){
          that.initComponents(that.detail)
        },300)

      },
      deep: true
    },
    SelectAlarmData: {
      handler(newVal, oldVal) {
        let that = this
        setTimeout(function (){
          that.initComponents(that.detail)
        },300)

      },
      deep: true
    },
  },
  data() {
    return {
      detail:{},
      IsToolBox:false,
      editMode:true,
      moment,
      dataZ:[],
      valueData: '',
      treePageSize: 100,
      scrollPage: 1,
      frontDataZ:[],
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      ChartTimelyRefreshTimer:null,
      ChartTimelyRefresh:5*60*1000,
      fillOpacity:1,
      SelectDateType: 'Day',
      TimeInterval:60,
      strokeOpacity:1,
      animateType:"blink",
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      ValueType:5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      IsShowDate:0,
      messageShowLoad:false,
      advanced: true,
      SelectDevice:"",
      SelectDeviceName:"",
      SelectDateRange:moment().format("YYYY-MM-DD"),
      SelectAlarmData:[],
      deviceTreeData:[],
      AlarmDataTree:[],
      SelectDeviceList:[],
      SelectDeviceDataList:[],
      cardHeight:0,
      EchartsHeight:0,
      base:{
        "text": "configComponent.RealDataChart.deviceSingleTitle",
        "icon": "icon-fsux_tubiao_huiguiquxiantu",
        "isFontIcon": true,
        "info": {
          "type": "real-data-chart",
          "action": [],
          "active": [],
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
              "w": 700,
              "h": 580
            },
            "backColor": "#ffffff",
            "foreColor": "#ffffff",
            "fontSize": 14,
            fontFamily: "Arial",
            "zIndex": 1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.ChartPublic.YMax",
                "type":1,
                "value":0,
                "key":"YMax",
              },
              {
                "name":"configComponent.ChartPublic.YMin",
                "type":1,
                "value":0,
                "key":"YMin",
              },
              {
                "name":"configComponent.ChartPublic.IsShowDate",
                type:6,
                value:0,
                enumList:[
                  {
                    value:0,
                    option:"False"
                  },
                  {
                    value:1,
                    option:"True"
                  },
                ],
                "key":"IsShowDate",
              },
              {
                "name":"configComponent.ChartPublic.EchartsWidth",
                "type":1,
                "value":2,
                "key":"EchartsWidth",
              },
              {
                "name":"configComponent.ChartPublic.EchartsXRotate",
                "type":1,
                "value":0,
                "key":"EchartsXRotate",
              },
              {
                "name":"configComponent.ChartPublic.EchartsXFormat",
                type:4,
                value:"YYYY-MM-DD HH:mm:ss",
                "key":"EchartsXFormat",
              },
              {
                "name":"configComponent.ChartPublic.EchartsXTheme",
                type:6,
                value:"dark",
                enumList:[
                  {
                    value:"chalk",
                    option:"chalk"
                  },
                  {
                    value:"essos",
                    option:"essos"
                  },
                  {
                    value:"dark",
                    option:"dark"
                  },
                  {
                    value:"infographic",
                    option:"infographic"
                  },
                  {
                    value:"macarons",
                    option:"macarons"
                  },
                  {
                    value:"roma",
                    option:"roma"
                  },
                  {
                    value:"shine",
                    option:"shine"
                  },
                  {
                    value:"vintage",
                    option:"vintage"
                  },
                  {
                    value:"purplePassion",
                    option:"purplePassion"
                  },
                  {
                    value:"walden",
                    option:"walden"
                  },
                  {
                    value:"westeros",
                    option:"westeros"
                  },
                  {
                    value:"wonderland",
                    option:"wonderland"
                  }
                ],
                "key":"EchartsXTheme",
              },
              {
                "name":"configComponent.ChartPublic.IsShowZoom",
                type:6,
                value:1,
                enumList:[
                  {
                    value:0,
                    option:"False"
                  },
                  {
                    value:1,
                    option:"True"
                  },
                ],
                "key":"IsShowZoom",
              },
            ]
          }
        }
      },
      EchartsXFormat:1,
      EchartsTheme:"dark",
      date: [],
      yieldRate: [],
      EchartsWidth:1,
      yieldIndex: [],
      timelySave:null,
      getList:[],
      option : {
        title: {
          text: "",
          textStyle:{
            color:"",
            fontFamily:"",
            fontSize:"",
          }
        },
        tooltip: {
          trigger: 'axis'
        },
        legend: {
          data: ['Email', 'Union Ads', 'Video Ads', 'Direct', 'Search Engine'],
          textStyle:{
            color:"",
            fontFamily:"",
            fontSize:"",
          }
        },
        grid: {
          left: '1%',
          right: '4%',
          bottom: '10%',
          containLabel: true
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: 'none'
            },
            magicType: { type: ['line', 'bar'] },
            restore: {},
            saveAsImage: {}
          }
        },
        xAxis: {
          type: 'time',
          boundaryGap: false,
          axisLine: {
            show: false, //是否显示坐标刻度
            lineStyle: {
              color: '#eeeeee'
            }
          },
          axisLabel: {
            show: true,
            rotate:70,
            margin: 15,
            interval: 30,
            formatter: '{yy}-{MM}-{dd}\n{HH}:{mm}:{ss}' , // 仅显示时分
            textStyle: {
              color: '#fff'
            }
          },
          data: []
        },
        yAxis: {
          type: 'value',
          min:100,
          max:1000,
          boundaryGap: [0, '100%'],
          axisLabel: {
            show: true,
            textStyle: {
              color: '#fff'
            }
          },
        },
        dataZoom: [
          // Y轴滑块
          {
            type: 'slider',
            yAxisIndex: 0,
            width: 15,
            filterMode: 'filter',
            start: 0,
            end: 100
          },
          // 内部控制器（支持鼠标滚轮）
          {
            type: 'inside',
            yAxisIndex: 0,
            zoomOnMouseWheel: 'shift', // 按住shift滚轮缩放
            moveOnMouseWheel: true,    // 直接滚轮平移
            moveOnMouseMove: true      // 鼠标拖动平移
          },
          // 内部控制器（支持鼠标滚轮）
          {
            type: 'inside',
            xAxisIndex: 0,
            zoomOnMouseWheel: 'shift', // 按住shift滚轮缩放
            moveOnMouseWheel: true,    // 直接滚轮平移
            moveOnMouseMove: true      // 鼠标拖动平移
          },
          {
            type: 'inside',
            xAxisIndex: 0,
            start: 0,
            end: 100,
            zoomLock: false
          },
          {
            type: 'slider',
            height: 15,
            xAxisIndex: 0,
          },
        ],
        series: [
          {
            name: 'Email',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Union Ads',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Video Ads',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Direct',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Search Engine',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          }
        ]
      },
      seriesMap:[]
    }
  },
  computed: {
    styleVar() {
      return {
        "--foreColor": this.foreColor ,
        '--backColor':this.backColor,
        "--selectedColor": this.selectedColor ,
        '--hoverColor': this.hoverColor,
        '--selectedTextColor': this.selectedTextColor,
        '--TextFontSize': this.TextFontSize+'px',
        '--hoverTextColor': this.hoverTextColor,
        '--SearchColor': this.SearchColor,
        '--SearchBackColor': this.SearchBackColor,
        '--SearchBorderColor': this.SearchBorderColor,
        '--EchartsHeight': this.EchartsHeight+'px',
        '--dateSelectColor': this.dateSelectColor,
        '--dateSelectBackColor': this.dateSelectBackColor,
        '--dateSelectBorderColor': this.dateSelectBorderColor,

        '--tableHeaderColor': this.tableHeaderColor,
        '--tableHeaderBackColor': this.tableHeaderBackColor,
        '--tableSplitColor': this.tableSplitColor,
        '--tableHoverColor':this.tableHoverColor
      };
    },
    CardStyle:function () {
      let styles = [];
      const customBg = this.detail && this.detail.style ? this.detail.style.backColor : ''
      const backgroundColor = customBg && customBg !== 'transparent' ? customBg : this.chartThemePalette.panelBg
      styles.push(`background-color: ${backgroundColor}`);
      if(this.detail.style.foreColor) {
        styles.push(`color: ${this.detail.style.foreColor}`);
      }
      let style = styles.join(';');
      return style;
    },
    CardBodyStyle() {
      const customBg = this.detail && this.detail.style ? this.detail.style.backColor : ''
      const backgroundColor = customBg && customBg !== 'transparent' ? customBg : this.chartThemePalette.panelBg
      return {
        padding: '8px 10px',
        background: backgroundColor,
      }
    },
    chartThemePalette() {
      const paletteMap = {
        dark: {
          panelBg: '#100C2A',
          accent: '#4992ff',
          accentSoft: 'rgba(73, 146, 255, 0.32)',
          mode: 'dark',
        },
        chalk: {
          panelBg: '#293441',
          accent: '#87f7cf',
          accentSoft: 'rgba(135, 247, 207, 0.28)',
          mode: 'dark',
        },
        purplePassion: {
          panelBg: '#5B5C6E',
          accent: '#ffb980',
          accentSoft: 'rgba(255, 185, 128, 0.26)',
          mode: 'dark',
        },
        vintage: {
          panelBg: '#FEF8EF',
          accent: '#d87c7c',
          accentSoft: 'rgba(216, 124, 124, 0.24)',
          mode: 'light',
        },
        essos: {
          panelBg: '#F6EED4',
          accent: '#893448',
          accentSoft: 'rgba(137, 52, 72, 0.18)',
          mode: 'light',
        },
        infographic: {
          panelBg: '#FFFFFF',
          accent: '#c1232b',
          accentSoft: 'rgba(193, 35, 43, 0.16)',
          mode: 'light',
        },
        macarons: {
          panelBg: '#FFFFFF',
          accent: '#2ec7c9',
          accentSoft: 'rgba(46, 199, 201, 0.18)',
          mode: 'light',
        },
        roma: {
          panelBg: '#FFFFFF',
          accent: '#e01f54',
          accentSoft: 'rgba(224, 31, 84, 0.16)',
          mode: 'light',
        },
        shine: {
          panelBg: '#FFFFFF',
          accent: '#c12e34',
          accentSoft: 'rgba(193, 46, 52, 0.16)',
          mode: 'light',
        },
        walden: {
          panelBg: '#FCFCFC',
          accent: '#3f9bc9',
          accentSoft: 'rgba(63, 155, 201, 0.18)',
          mode: 'light',
        },
        westeros: {
          panelBg: '#FFFFFF',
          accent: '#516b91',
          accentSoft: 'rgba(81, 107, 145, 0.18)',
          mode: 'light',
        },
        wonderland: {
          panelBg: '#FFFFFF',
          accent: '#4ea397',
          accentSoft: 'rgba(78, 163, 151, 0.18)',
          mode: 'light',
        },
      }
      return paletteMap[this.EchartsTheme] || paletteMap.dark
    },
    controlThemeClass() {
      return this.chartThemePalette.mode === 'light' ? 'history-chart-single-theme-light' : 'history-chart-single-theme-dark'
    },
    dropdownThemeClass() {
      return this.chartThemePalette.mode === 'light' ? 'history-chart-single-dropdown-light' : 'history-chart-single-dropdown-dark'
    },
    dropdownStyle() {
      return {
        maxHeight: '400px',
        overflow: 'auto',
        '--history-chart-single-accent': this.chartThemePalette.accent,
        '--history-chart-single-accent-soft': this.chartThemePalette.accentSoft,
      }
    },
  },
  methods: {
    syncEchartsThemeFromDetail(option) {
      if (!option || !option.style || !Array.isArray(option.style.diy)) {
        return
      }
      const themeOption = option.style.diy.find(item => item.key === 'EchartsXTheme')
      if (themeOption && themeOption.value) {
        this.EchartsTheme = themeOption.value
      }
    },
    getPopupContainer(triggerNode) {
      return triggerNode && triggerNode.parentNode ? triggerNode.parentNode : document.body
    },
    waitChartContainerReady(view, callback, retryCount = 0) {
      if (!view) {
        return
      }
      if (view.clientWidth > 0 && view.clientHeight > 0) {
        callback()
        return
      }
      if (retryCount >= 10) {
        view.style.width = this.detail.style.position.w + 'px'
        view.style.height = this.EchartsHeight + 'px'
        callback()
        return
      }
      const raf = typeof requestAnimationFrame === 'function' ? requestAnimationFrame : (fn) => setTimeout(fn, 16)
      raf(() => this.waitChartContainerReady(view, callback, retryCount + 1))
    },
    chargeDateType(e){
      let type = e.target.value
      if(type=="Day")
      {
        this.SelectDateRange = moment().format("YYYY-MM-DD")
      }
      else  if(type=="Weekly")
      {
        const startDate = moment().day(1).format('YYYY-MM-DD'); // 周一日期
        const endDate = moment().day(7).format('YYYY-MM-DD'); // 周日日期
        this.SelectDateRange = [startDate,endDate]
      }
      else  if(type=="Month")
      {
        this.SelectDateRange = moment().format("YYYY-MM")
      }
      else{
        this.SelectDateRange =[moment().add(-1, 'day'),moment()]
      }
    },
    handleSearch (val) {
      this.valueData = val
      if (!val) {
        this.GetDeviceModelDataList()
      } else {
        this.frontDataZ = []
        this.scrollPage = 1
        this.dataZ.forEach(item => {
          if (item.name.indexOf(val) >= 0) {
            this.frontDataZ.push(item)
          }
        })
        this.allDataZ = this.frontDataZ
        this.frontDataZ = this.frontDataZ.slice(0, 100)
      }
    },
//下拉框下滑事件
    handlePopupScroll (e) {
      const { target } = e
      const scrollHeight = target.scrollHeight - target.scrollTop
      const clientHeight = target.clientHeight
      // 下拉框不下拉的时候
      if (scrollHeight === 0 && clientHeight === 0) {
        this.scrollPage = 1
      } else {
        // 当下拉框滚动条到达底部的时候
        if (scrollHeight < clientHeight + 5) {
          this.scrollPage = this.scrollPage + 1
          const scrollPage = this.scrollPage// 获取当前页
          const treePageSize = this.treePageSize * (scrollPage || 1)// 新增数据量
          const newData = [] // 存储新增数据
          let max = '' // max 为能展示的数据的最大条数
          if (this.dataZ.length > treePageSize) {
            // 如果总数据的条数大于需要展示的数据
            max = treePageSize
          } else {
            // 否则
            max = this.dataZ.length
          }
          // 判断是否有搜索
          if (this.valueData) {
            this.allDataZ.forEach((item, index) => {
              if (index < max) { // 当data数组的下标小于max时
                newData.push(item)
              }
            })
          } else {
            this.dataZ.forEach((item, index) => {
              if (index < max) { // 当data数组的下标小于max时
                newData.push(item)
              }
            })
          }

          this.frontDataZ = newData // 将新增的数据赋值到要显示的数组中
        }
      }
    },
    refreshChart(){
      if(this.echartsView==null)
      {
        return
      }
      this.date = []
      this.initComponents(this.detail)
    },
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      let i=0
      let refObj = this.detail.identifier
      let view = this.$refs[refObj]
      let refCardObj = this.detail.identifier+'card'
      let cardObj = this.$refs[refCardObj]
      if (!view || !cardObj || !cardObj.$el) {
        return
      }
      this.cardHeight = cardObj.$el.clientHeight
      this.fillOpacity=option.style.opacity
      this.EchartsHeight = this.detail.style.position.h-this.cardHeight
      view.style.width = this.detail.style.position.w + 'px'
      view.style.height = this.EchartsHeight + 'px'
      if (view.clientWidth === 0 || view.clientHeight === 0) {
        this.waitChartContainerReady(view, () => this.initComponents(option))
        return
      }
      this.option.title.textStyle.color = option.style.foreColor
      this.option.title.textStyle.fontSize = option.style.fontSize
      this.option.title.textStyle.fontFamily  = option.style.fontFamily

      this.option.legend.textStyle.color = option.style.foreColor
      this.option.legend.textStyle.fontSize = option.style.fontSize
      this.option.legend.textStyle.fontFamily  = option.style.fontFamily
      this.option.xAxis.axisLabel.interval=30
      this.option.xAxis.axisLabel.textStyle.color = option.style.foreColor
      this.option.xAxis.axisLabel.textStyle.fontSize = option.style.fontSize
      this.option.xAxis.axisLabel.textStyle.fontFamily  = option.style.fontFamily

      this.option.yAxis.axisLabel.textStyle.color = option.style.foreColor
      this.option.yAxis.axisLabel.textStyle.fontSize = option.style.fontSize
      this.option.yAxis.axisLabel.textStyle.fontFamily  = option.style.fontFamily

      view.style.height = this.EchartsHeight+'px'

      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.title.text=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="IsShowZoom")
        {
          if(option.style.diy[i].value!=1)
          {
            this.option.dataZoom=[]
            this.option.grid.bottom='1%'
          }
          else
          {
            this.option.dataZoom=[
              // Y轴滑块
              {
                type: 'slider',
                yAxisIndex: 0,
                width: 15,
                filterMode: 'filter',
                start: 0,
                end: 100
              },
              // 内部控制器（支持鼠标滚轮）
              {
                type: 'inside',
                yAxisIndex: 0,
                zoomOnMouseWheel: 'shift', // 按住shift滚轮缩放
                moveOnMouseWheel: true,    // 直接滚轮平移
                moveOnMouseMove: true      // 鼠标拖动平移
              },
              // 内部控制器（支持鼠标滚轮）
              {
                type: 'inside',
                xAxisIndex: 0,
                zoomOnMouseWheel: 'shift', // 按住shift滚轮缩放
                moveOnMouseWheel: true,    // 直接滚轮平移
                moveOnMouseMove: true      // 鼠标拖动平移
              },
              {
                type: 'inside',
                xAxisIndex: 0,
                start: 0,
                end: 100,
                zoomLock: false
              },
              {
                type: 'slider',
                height: 15,
                xAxisIndex: 0,
              },
            ]
            this.option.grid.bottom='10%'
          }
        }
        else if(option.style.diy[i].key=="TimelyRefresh")
        {
          this.ChartTimelyRefresh = parseInt(option.style.diy[i].value)*60*1000
        }
        else if(option.style.diy[i].key=="YMax")
        {
          if(option.style.diy[i].value==0)
          {
            this.option.yAxis.max = 'dataMax'
          }
          else {
            this.option.yAxis.max = parseFloat(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="YMin")
        {
          if(option.style.diy[i].value==0)
          {
            this.option.yAxis.min = 'dataMin'
          }
          else {
            this.option.yAxis.min = parseFloat(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="ChartRange")
        {
          this.ChartRange = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="TimeInterval")
        {
          this.TimeInterval = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="ValueType")
        {
          this.ValueType = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="IsShowDate")
        {
          this.IsShowDate = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsWidth")
        {
          this.EchartsWidth = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsXRotate")
        {
          this.option.xAxis.axisLabel.rotate = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsXFormat")
        {
          this.EchartsXFormat = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="EchartsXTheme")
        {
          this.EchartsTheme = option.style.diy[i].value
          if(this.echartsView) {
            this.echartsView.dispose()
            this.echartsView = null
          }
          this.echartsView = echarts.init(view, this.EchartsTheme);
        }
      }
      if (!this.echartsView) {
        this.echartsView = echarts.init(view, this.EchartsTheme);
      }
      i=0
      this.animateType = option.animate.selected
      if(option.animate.isExpression)
      {
        this.isStart = false
      }
      else
      {
        this.isStart = true
      }
      for( i=0;i<option.animate.animateElement.length;i++)
      {
        if(option.animate.animateElement[i].id=="millcolorGrad")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++)
          {
            if(option.animate.animateElement[i].elementList[k].key=="startColor")
            {
              this.startColor=option.animate.animateElement[i].elementList[k].value
            }
            else if(option.animate.animateElement[i].elementList[k].key=="stopColor")
            {
              this.stopColor=option.animate.animateElement[i].elementList[k].value
            }
            else if(option.animate.animateElement[i].elementList[k].key=="animateSpeed")
            {
              this.animateSpeed=option.animate.animateElement[i].elementList[k].value
            }
          }
        }
        else if(option.animate.animateElement[i].id=="blink")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
            if (option.animate.animateElement[i].elementList[k].key == "blinkSpeed") {
              this.blinkSpeed = option.animate.animateElement[i].elementList[k].value
            }
          }
        }
        else if(option.animate.animateElement[i].id=="animateSpin")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
            if (option.animate.animateElement[i].elementList[k].key == "spinSpeed") {
              this.animateSpinSpeed = option.animate.animateElement[i].elementList[k].value
            }
            else if (option.animate.animateElement[i].elementList[k].key == "spinDirection") {
              this.spinDirection = option.animate.animateElement[i].elementList[k].value
            }
          }
        }
      }

      this.echartsView.resize()
      this.echartsView.setOption(this.option);
    },
    timelyTimeGetHistory(){
      let _t = this
      if (!this.echartsView) {
        return
      }
      let ChartRangeTime = 0
      if(this.ChartRange==0)
      {
        ChartRangeTime = 60
      }
      else if(this.ChartRange==1){
        ChartRangeTime = 60*5
      }
      else if(this.ChartRange==2){
        ChartRangeTime = 60*15
      }
      else if(this.ChartRange==3){
        ChartRangeTime = 60*24
      }
      else if(this.ChartRange==4){
        ChartRangeTime = 60*24*3
      }
      else if(this.ChartRange==5){
        ChartRangeTime = 60*24*7
      }
      else if(this.ChartRange==6){
        ChartRangeTime = 60*24*15
      }
      else if(this.ChartRange==7){
        ChartRangeTime = 60*24*30
      }
      const params ={
        TimeIn:this.TimeInterval,
        HistoryTime:ChartRangeTime,
        List: this.getList
      }
      if(this.getList.length==0)
      {
        _t.echartsView.setOption(_t.option,true)
        return
      }
      GetChartDataHistoryList(params).then(function (res){
        if(res.data.code==0&&res.data.data!=null)
        {
          for(let i=0;i<_t.option.series.length;i++)
          {
            _t.option.series[i].data =[]
            _t.seriesMap[_t.option.series[i].dataID] = []
          }
          _t.option.xAxis.data = []
          for(let i=0;i<res.data.data.length;i++)
          {
            let DateTime = moment(res.data.data[i].HistoryRecordDateTime).format(_t.EchartsXFormat);
            _t.option.xAxis.data.push(DateTime)

            for(let k=0;k<res.data.data[i].dataList.length;k++)
            {
              let item = res.data.data[i].dataList[k]
              if(item.average=="-") {
                item.average=0
              }
              if(item.diff=="-") {
                item.diff=0
              }
               if(item.max=="-") {
                item.max=0
              }
               if(item.min=="-") {
                item.min=0
              }
               if(item.sum=="-") {
                item.sum=0
              }
              if(_t.ValueType==0)
              {
                const seriesKey = item.DeviceUuid + item.ModelDataUuid
                if (_t.seriesMap[seriesKey]) {
                  _t.seriesMap[seriesKey].push(item.Value ? item.Value : 0)
                }
              }
              else if(_t.ValueType==1)
              {
                const seriesKey = item.DeviceUuid + item.ModelDataUuid
                if (_t.seriesMap[seriesKey]) {
                  _t.seriesMap[seriesKey].push(item.max ? item.max : 0)
                }
              }
              else if(_t.ValueType==2)
              {
                const seriesKey = item.DeviceUuid + item.ModelDataUuid
                if (_t.seriesMap[seriesKey]) {
                  _t.seriesMap[seriesKey].push(item.min ? item.min : 0)
                }
              }
              else if(_t.ValueType==3)
              {
                const seriesKey = item.DeviceUuid + item.ModelDataUuid
                if (_t.seriesMap[seriesKey]) {
                  _t.seriesMap[seriesKey].push(item.diff ? item.diff : 0)
                }
              }
              else if(_t.ValueType==4)
              {
                const seriesKey = item.DeviceUuid + item.ModelDataUuid
                if (_t.seriesMap[seriesKey]) {
                  _t.seriesMap[seriesKey].push(item.sum ? item.sum : 0)
                }
              }
              else if(_t.ValueType==5)
              {
                const seriesKey = item.DeviceUuid + item.ModelDataUuid
                if (_t.seriesMap[seriesKey]) {
                  _t.seriesMap[seriesKey].push(item.average ? item.average : 0)
                }
              }
            }
          }
          for(let i=0;i<_t.option.series.length;i++)
          {
            _t.option.series[i].data = _t.seriesMap[_t.option.series[i].dataID] || [];
          }

        }
        _t.echartsView.setOption(_t.option,true)
      })
    },
    onResize() {
      if (this.echartsView) {
        this.echartsView.resize();
      }
    },
    updateView() {
      this.setOption(this.option);
    },
    GetDeviceModelDataList(){
      let _t = this
      this.AlarmDataTree=[]
      _t.dataZ=[]
      _t.frontDataZ=[]
      if(!this.SelectDevice)
      {
        return
      }
      const params ={
        SelectDevice:[this.SelectDevice],
        getType:2
      }
      GetDeviceModelDataList(params).then(function (res){
        if(res.data.code==0)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            if(typeof res.data.list[i].DataList!="undefined"&&res.data.list[i].DataList!=null) {
              for (let j = 0; j < res.data.list[i].DataList.length; j++) {
                _t.dataZ.push(res.data.list[i].DataList[j])
                _t.AlarmDataTree.push(res.data.list[i].DataList[j])
              }
            }
          }
          _t.frontDataZ = _t.dataZ.slice(0, 100)
        }
      })
    },
    onDateChange(date, dateString){
      this.SelectDateRange = dateString
    },
    onWeeklyDateChange(date, dateString){
      const startDate = moment(date).day(1).format('YYYY-MM-DD'); // 周一日期
      const endDate = moment(date).day(7).format('YYYY-MM-DD'); // 周日日期
      this.SelectDateRange = [startDate,endDate]
    },
    getMonitorTree(){
      let _t = this
      this.deviceTreeData=[]
      getMonitorTree().then(function (res){
        if(res.data.code==0)
        {
          _t.deviceTreeData =res.data.list
        }
      })
    },
    DeviceDataCharge(value,option){
      this.SelectDeviceDataList = []
      const optionList = Array.isArray(option) ? option : []
      for(let i=0;i<optionList.length;i++)
      {
        let temp_obj = {}
        temp_obj.title = optionList[i].componentOptions.children[0].text
        temp_obj.DataUUID = optionList[i].componentOptions.propsData.value
        this.SelectDeviceDataList.push(temp_obj)
      }
    },
    SelectTreeDevice(value,node,extera){
      this.SelectDevice = value || ""
      this.SelectDeviceName = ""
      this.SelectDeviceList = []
      this.SelectAlarmData = []
      this.SelectDeviceDataList = []
      if(this.SelectDevice)
      {
        const deviceName = extera && extera.triggerNode ? extera.triggerNode.title : ""
        this.SelectDeviceName = deviceName
        this.SelectDeviceList.push({
          deviceUUID: this.SelectDevice,
          deviceName: this.SelectDeviceName,
        })
        this.GetDeviceModelDataList()
      }
    },
    QueryHistoryDataList(){
      let _t = this
      _t.dataSource = []
      const params = {
        deviceList:this.SelectDevice ? [this.SelectDevice] : [],
        dataList:this.SelectAlarmData,
        dateType:this.SelectDateType,
        dateRange:this.SelectDateRange,
      }
      if(this.SelectDeviceDataList.length==0)
      {
        this.$message.error(this.$t("reporting.AlarmHistory.SelectDataError"))
        return
      }
      if(this.SelectDevice.length==0)
      {
        this.$message.error(this.$t("reporting.AlarmHistory.SelectDeviceError"))
        return
      }
      for(let i=0;i<this.option.series.length;i++)
      {
        this.option.series[i].data = []
        this.option.series[i].smooth= true
        this.seriesMap[this.option.series[i].dataID] = []
      }
      this.option.xAxis.data = [];
      this.option.legend.data=[]
      this.option.series=[]
      this.getList =[]
      for(let k=0;k<this.SelectDeviceDataList.length;k++)
      {
        let listObj = {}
        const deviceName = this.SelectDeviceName || this.SelectDevice
        this.option.legend.data.push(deviceName+"-"+this.SelectDeviceDataList[k].title)
        let series= {
          name: deviceName+"-"+this.SelectDeviceDataList[k].title,
          type: 'line',
          dataID: this.SelectDevice+this.SelectDeviceDataList[k].DataUUID,
          data: [],
          itemStyle: {
            normal: {
              lineStyle: {
                width:this.EchartsWidth
              }
            }
          },
          symbolSize: this.EchartsWidth,
          markPoint: {
            symbol: 'circle',
            symbolSize: 10,
            label: {
              show: true,
              color: '#ffffff',
              fontSize: 14,
              fontWeight: 'bold',
              backgroundColor: 'rgba(16, 20, 36, 0.82)',
              borderColor: 'rgba(255,255,255,0.35)',
              borderWidth: 1,
              borderRadius: 4,
              padding: [3, 6],
              formatter: function(param) {
                if (!param) return ''
                const label = param.name === 'Max' ? 'Max:' : (param.name === 'Min' ? 'Min:' : '')
                return label + (param.value !== undefined ? param.value : '')
              }
            },
            data: []
          },
          markLine: {
            symbol: ['none', 'arrow'],
            lineStyle: {
              type: 'dashed',
              width: 1
            },
            label: {
              show: true,
              color: '#ffffff',
              fontSize: 13,
              fontWeight: 'bold',
              backgroundColor: 'rgba(16, 20, 36, 0.82)',
              borderColor: 'rgba(255,255,255,0.35)',
              borderWidth: 1,
              borderRadius: 4,
              padding: [3, 6],
              formatter: function(param) {
                return param && param.value !== undefined ? ('Avg:' + param.value) : ''
              }
            },
            data: []
          }
        }
        listObj.DeviceUuid = this.SelectDevice
        listObj.ModelDataUuid = this.SelectDeviceDataList[k].DataUUID
        this.getList.push(listObj)
        this.option.series.push(series)
        this.seriesMap[this.SelectDevice+this.SelectDeviceDataList[k].DataUUID]=[]
      }

      if(params.dateRange==""||(params.dateRange[0]==""))
      {
        this.$message.error(this.$t("reporting.DataHistory.SelectDateError"))
        return
      }
      _t.echartsView.clear()
      this.messageShowLoad=true
      GetDataHistoryList(params).then(function (res){
        if(res.data.code==0)
        {
            _t.dataSource =res.data.list
            for(let i=0;i<_t.option.series.length;i++)
            {
              _t.option.series[i].data =[]
              _t.seriesMap[_t.option.series[i].dataID] = []
            }
            _t.option.xAxis.data = []
            for(let i=0;i<_t.dataSource.length;i++)
            {
              let singleArray={
                value:[]
              }
              let DateTime = moment(_t.dataSource[i].RecordTime).format(_t.EchartsXFormat);
              // _t.option.xAxis.data.push(DateTime)
              let item = _t.dataSource[i]
              singleArray.value[0]=DateTime
              singleArray.value[1]=item.DataValue ? item.DataValue : 0
              const seriesKey = item.DeviceUuid + item.ModelDataUuid
              if (_t.seriesMap[seriesKey]) {
                _t.seriesMap[seriesKey].push(singleArray)
              }
            }
            _t.option.xAxis.data = Array.from(new Set(_t.option.xAxis.data))
            for(let i=0;i<_t.option.series.length;i++)
            {
              _t.option.series[i].data = _t.seriesMap[_t.option.series[i].dataID] || [];
              const seriesValues = _t.option.series[i].data
                .map(item => Array.isArray(item.value) ? parseFloat(item.value[1]) : parseFloat(item))
                .filter(item => !isNaN(item))
              if (seriesValues.length) {
                const maxVal = Math.max.apply(null, seriesValues)
                const minVal = Math.min.apply(null, seriesValues)
                const avgVal = seriesValues.reduce((total, current) => total + current, 0) / seriesValues.length
                _t.option.series[i].markPoint.data = [
                  { type: 'max', name: 'Max', label: { position: 'insideTop' } },
                  { type: 'min', name: 'Min', label: { position: 'insideBottom' } }
                ]
                _t.option.series[i].markLine.data = [
                  { yAxis: avgVal.toFixed(2), label: { position: 'middle' } }
                ]
              } else {
                _t.option.series[i].markPoint.data = []
                _t.option.series[i].markLine.data = []
              }
            }

            _t.echartsView.setOption(_t.option,true)
        }
        _t.messageShowLoad=false
      })
    },
    // 获取当前时间
  },
  beforeDestroy () {
    clearInterval(this.ChartTimelyRefreshTimer)
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.dispose()
    }
  },
  created(){
    let _t = this
    this.getMonitorTree()
    this.GetDeviceModelDataList()
    const node = this.getNode()
    node.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    node.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
      _t.$nextTick(() => {
        const view = _t.$refs[_t.detail.identifier]
        const cardObj = _t.$refs[_t.detail.identifier + 'card']
        if (view && cardObj && cardObj.$el) {
          _t.cardHeight = cardObj.$el.clientHeight
          _t.EchartsHeight = current.height - _t.cardHeight
          view.style.width = current.width + 'px'
          view.style.height = _t.EchartsHeight + 'px'
        }
        _t.onResize()
      })
    });
    this.detail = node.getData().detail
    this.syncEchartsThemeFromDetail(this.detail)
    this.editMode = node.getData().editMode
    this.showDeviceUuid = node.getData().showDeviceUuid
    this.IsToolBox = node.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
    })
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据

      _t.$EventBus.$on(activeEvent, (data) => {

      })
      _t.$EventBus.$on(animateEvent, (data) => {
        _t.isStart = data
      })
    });
  }
}
</script>

<style lang="less">
.history-chart-single-card {
  .ant-card-body {
    padding: 8px 10px;
    background: transparent;
  }

  .ant-form-item {
    margin-bottom: 8px;
  }

  .ant-form-item-label > label {
    font-weight: 500;
  }

  .ant-select-selection,
  .ant-calendar-picker-input,
  .ant-input,
  .ant-radio-button-wrapper,
  .history-chart-single-query {
    height: 32px;
    border-radius: 4px;
    box-shadow: none;
  }

  .history-chart-single-query {
    min-width: 58px;
    padding: 0 14px;
    font-weight: 500;
  }

}

.history-chart-single-theme-dark {
  color: #f4f8ff;

  .ant-card-body,
  .ant-form,
  .ant-row,
  .fold {
    background: transparent;
  }

  .ant-form-item-label > label,
  .ant-radio-button-wrapper,
  .ant-select-selection,
  .ant-calendar-picker-input,
  .ant-input,
  .ant-calendar-picker-icon,
  .ant-select-arrow {
    color: #f4f8ff;
  }

  .ant-form-item-label > label {
    text-shadow: 0 1px 1px rgba(0, 0, 0, 0.32);
  }

  .ant-select-selection,
  .ant-calendar-picker-input,
  .ant-input,
  .ant-radio-button-wrapper,
  .history-chart-single-query {
    background: rgba(39, 55, 82, 0.98);
    border-color: rgba(150, 181, 225, 0.82);
  }

  .ant-select-selection__rendered,
  .ant-calendar-picker-input input,
  .ant-input::placeholder,
  .ant-select-selection__placeholder {
    color: rgba(244, 248, 255, 0.84);
  }

  .ant-select-selection__choice {
    color: #f4f8ff;
    background: rgba(12, 24, 42, 0.92);
    border-color: rgba(150, 181, 225, 0.62);
  }

  .ant-select-selection__placeholder,
  .ant-select-selection__choice__remove {
    color: rgba(244, 248, 255, 0.78);
  }

  .ant-radio-button-wrapper-checked,
  .ant-select-focused .ant-select-selection,
  .ant-select-selection:hover,
  .ant-calendar-picker-input:hover,
  .ant-input:hover,
  .ant-input:focus,
  .history-chart-single-query:hover,
  .history-chart-single-query:focus {
    color: #ffffff;
    background: rgba(59, 90, 135, 1);
    border-color: rgba(190, 218, 255, 0.96);
  }

  .history-chart-single-query {
    color: #f4f8ff;
  }

  .history-chart-single-query[disabled] {
    color: rgba(244, 248, 255, 0.48);
    background: rgba(31, 43, 64, 0.74);
    border-color: rgba(112, 139, 176, 0.42);
  }
}

.history-chart-single-theme-light {
  color: #1f2937;

  .ant-form-item-label > label,
  .ant-radio-button-wrapper,
  .ant-select-selection,
  .ant-calendar-picker-input,
  .ant-input,
  .ant-calendar-picker-icon,
  .ant-select-arrow {
    color: #1f2937;
  }

  .ant-select-selection,
  .ant-calendar-picker-input,
  .ant-input,
  .ant-radio-button-wrapper,
  .history-chart-single-query {
    background: rgba(255, 255, 255, 0.92);
    border-color: rgba(148, 163, 184, 0.75);
  }

  .ant-select-selection__choice {
    color: #1f2937;
    background: rgba(241, 245, 249, 0.92);
    border-color: rgba(148, 163, 184, 0.58);
  }

  .ant-select-selection__placeholder,
  .ant-select-selection__choice__remove {
    color: rgba(31, 41, 55, 0.62);
  }

  .ant-radio-button-wrapper-checked,
  .history-chart-single-query:hover,
  .history-chart-single-query:focus {
    color: #1d4ed8;
    background: #ffffff;
    border-color: rgba(59, 130, 246, 0.75);
  }
}

.history-chart-single-dropdown-dark {
  color: #f4f8ff;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 12px;
    background: #172238;
    pointer-events: none;
    z-index: 10;
  }

  .ant-select-dropdown,
  .ant-select-dropdown-menu,
  .ant-select-tree,
  .ant-select-tree-dropdown,
  .ant-calendar,
  .ant-calendar-panel,
  .ant-calendar-input-wrap,
  .ant-calendar-time-picker,
  .ant-calendar-time-picker-inner,
  .ant-calendar-time-picker-select {
    color: #f4f8ff;
    background: #172238;
  }

  .ant-select-search__field,
  .ant-calendar-input {
    color: #f4f8ff;
    background: rgba(39, 55, 82, 0.98);
    border-color: rgba(150, 181, 225, 0.82);
  }

  .ant-select-dropdown-search,
  .ant-select-tree-dropdown-search,
  .ant-select-dropdown-search .ant-select-search,
  .ant-select-tree-dropdown-search .ant-select-search,
  .ant-select-dropdown-search .ant-select-search__field__wrap,
  .ant-select-tree-dropdown-search .ant-select-search__field__wrap {
    background: #172238 !important;
  }

  .ant-select-dropdown-search .ant-select-search__field__wrap,
  .ant-select-tree-dropdown-search .ant-select-search__field__wrap {
    border-radius: 4px;
    border: 1px solid rgba(150, 181, 225, 0.82) !important;
  }

  .ant-select-dropdown-search,
  .ant-select-tree-dropdown-search {
    padding: 6px;
    margin: 0 !important;
    border: 0 !important;
    border-bottom: 1px solid rgba(150, 181, 225, 0.28) !important;
    box-shadow: none !important;
  }

  .ant-select-dropdown-search .ant-select-search__field,
  .ant-select-tree-dropdown-search .ant-select-search__field {
    height: 30px;
    border-radius: 4px;
    color: #f4f8ff !important;
    background: rgba(39, 55, 82, 0.98) !important;
    border: 0 !important;
    outline: none !important;
    box-shadow: none !important;
  }

  .ant-select-dropdown-search .ant-select-search__field:hover,
  .ant-select-dropdown-search .ant-select-search__field:focus,
  .ant-select-tree-dropdown-search .ant-select-search__field:hover,
  .ant-select-tree-dropdown-search .ant-select-search__field:focus {
    border-color: rgba(190, 218, 255, 0.96) !important;
    box-shadow: 0 0 0 1px rgba(116, 166, 230, 0.28) !important;
  }

  .ant-select-tree,
  .ant-select-dropdown-menu {
    margin-top: 0 !important;
    padding-top: 0 !important;
    border-top: 0 !important;
    box-shadow: none !important;
  }

  .ant-select-tree > li:first-child {
    padding-top: 6px !important;
  }

  .ant-select-dropdown-content,
  .ant-select-tree-dropdown-content {
    background: #172238 !important;
    box-shadow: none !important;
  }

  .ant-select-dropdown,
  .ant-select-tree-dropdown {
    background: #172238 !important;
    border: 0 !important;
    box-shadow: none !important;
    padding: 0 !important;
    overflow: hidden !important;
  }

  .ant-select-tree-dropdown > div,
  .ant-select-tree-dropdown .ant-select-tree,
  .ant-select-tree-dropdown ul,
  .ant-select-tree-dropdown li,
  .ant-select-tree-dropdown .ant-select-tree-list-holder-inner > * {
    background: #172238 !important;
  }

  .ant-select-tree-list-holder,
  .ant-select-tree-list,
  .ant-select-tree-list-scrollbar,
  .ant-select-tree-list-holder-inner {
    background: #172238 !important;
  }

  .ant-select-dropdown-menu,
  .ant-select-dropdown-menu-root,
  .ant-select-tree,
  .ant-select-tree-list-holder::before,
  .ant-select-tree-list-holder::after,
  .ant-empty,
  .ant-empty-image,
  .ant-empty-description {
    background: #172238 !important;
    border: 0 !important;
  }

  .ant-select-tree-list-scrollbar-thumb,
  .ant-select-tree-list-scrollbar-track {
    background: #172238 !important;
  }

  .ant-select-tree-node-content-wrapper,
  .ant-select-tree-title,
  .ant-select-dropdown-menu-item,
  .ant-calendar-date,
  .ant-calendar-month-panel-month,
  .ant-calendar-year-panel-year,
  .ant-calendar-decade-panel-decade {
    color: #f4f8ff !important;
  }

  .ant-select-tree-node-content-wrapper:hover,
  .ant-select-tree-node-content-wrapper:hover .ant-select-tree-title,
  .ant-select-dropdown-menu-item:hover,
  .ant-select-dropdown-menu-item-active,
  .ant-calendar-date:hover,
  .ant-calendar-month-panel-month:hover,
  .ant-calendar-year-panel-year:hover {
    color: #ffffff !important;
    background: #4c6fa3 !important;
  }

  .ant-select-tree-node-selected,
  .ant-select-tree-node-selected .ant-select-tree-title,
  .ant-select-dropdown-menu-item-selected,
  .ant-calendar-selected-day .ant-calendar-date,
  .ant-calendar-today .ant-calendar-date {
    color: #ffffff !important;
    background: #4c6fa3 !important;
  }

  .ant-select-tree li .ant-select-tree-node-content-wrapper {
    transition: background 0.15s ease, color 0.15s ease;
    border-radius: 4px;
    padding-right: 6px;
  }

  .ant-select-tree-node-selected {
    border-left: 3px solid var(--history-chart-single-accent);
    border-radius: 3px;
  }

  .ant-select-tree-switcher,
  .ant-select-tree-switcher-icon,
  .ant-select-tree-iconEle,
  .ant-calendar-header,
  .ant-calendar-header a,
  .ant-calendar-footer,
  .ant-calendar-time-picker-btn {
    color: rgba(244, 248, 255, 0.9) !important;
  }

  .ant-select-tree li .ant-select-tree-node-content-wrapper {
    color: #f4f8ff !important;
  }

  .ant-select-search__field::placeholder,
  .ant-calendar-input::placeholder {
    color: rgba(244, 248, 255, 0.76);
  }

  .ant-calendar-header,
  .ant-calendar-input-wrap,
  .ant-calendar-footer,
  .ant-calendar-time-picker-inner {
    border-color: rgba(150, 181, 225, 0.32);
  }

  ::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }

  ::-webkit-scrollbar-track {
    background: rgba(12, 24, 42, 0.65);
  }

  ::-webkit-scrollbar-thumb {
    background: rgba(150, 181, 225, 0.68);
    border-radius: 999px;
  }
}

.history-chart-single-dropdown-light {
  .ant-select-tree-node-selected,
  .ant-select-dropdown-menu-item-selected,
  .ant-calendar-selected-day .ant-calendar-date,
  .ant-calendar-today .ant-calendar-date {
    color: #0f172a;
    background: var(--history-chart-single-accent-soft);
  }

  .ant-select-tree-node-selected {
    border-left: 3px solid var(--history-chart-single-accent);
    border-radius: 3px;
  }
}

.view-chart-real-data {
  width: 100%;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject :style="{overflow:'visible','z-index':5}" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
          <baidu-map  @ready="mapHandler"  :scrollWheelZoom="true" :inertialDragging="false" :center="center" :zoom="zoom" :style="{'z-index':1000,width:detail.style.position.w+'px' ,height:detail.style.position.h+'px'}">
             <div v-if="mapUpdate">
               <bm-geolocation anchor="BMAP_ANCHOR_BOTTOM_RIGHT" :showAddressBar="true" :autoLocation="true"></bm-geolocation>
               <template v-for="(item,index) in DeviceListMap" >
                 <template v-if="item.Status==1">
                   <bm-marker  v-if="item.AlarmCount>0" :animation="'BMAP_ANIMATION_BOUNCE'"  :key="index"     @click="lookDetail(index)"  :position="{lng: item.longitude, lat: item.latitude}" :dragging="false"  :icon="{url: MakerUrlAlarm, size: {width: MakerWidth, height: MakerHeight}}">
                     <bm-label :content="item.name" @click="lookDetail(index)" :labelStyle="{border:'','background-color':detail.style.backColor,'font-weight':detail.style.fontWeight,'font-family':detail.style.fontFamily,color: detail.style.foreColor, fontSize : detail.style.fontSize+'px'}" :offset="{width: LabelOffsetX, height: LabelOffsetY}"/>
                     <bm-info-window   :autoPan="true" :title="item.name"  :position="{lng: item.longitude, lat: item.latitude}"   :show="item.showFlag" @close="infoWindowClose(index)">

                        告警数量 : {{ item.AlarmCount }}<br/>
                       状态 : {{ item.Status==2?"未激活":item.Status?"在线":"离线" }}

                     </bm-info-window>
               </bm-marker>
                   <bm-marker  v-else :key="index"     @click="lookDetail(index)"  :position="{lng: item.longitude, lat: item.latitude}" :dragging="false"  :icon="{url:MakerUrl , size: {width: MakerWidth, height: MakerHeight}}">
                     <bm-label :content="item.name" @click="lookDetail(index)" :labelStyle="{border:'','background-color':detail.style.backColor,'font-weight':detail.style.fontWeight,'font-family':detail.style.fontFamily,color: detail.style.foreColor, fontSize : detail.style.fontSize+'px'}" :offset="{width: LabelOffsetX, height: LabelOffsetY}"/>
                     <bm-info-window   :autoPan="true" :title="item.name"  :position="{lng: item.longitude, lat: item.latitude}"   :show="item.showFlag" @close="infoWindowClose(index)">

                        告警数量 : {{ item.AlarmCount }}<br/>
                       状态 : {{ item.Status==2?"未激活":item.Status?"在线":"离线" }}

                     </bm-info-window>
               </bm-marker>
                 </template>

                 <bm-marker v-else-if="item.Status==0||item.Status==2"   :key="index"     @click="lookDetail(index)"  :position="{lng: item.longitude, lat: item.latitude}" :dragging="false"  :icon="{url:MakerUrlOffline, size: {width: MakerWidth, height: MakerHeight}}">
                     <bm-label :content="item.name" @click="lookDetail(index)" :labelStyle="{border:'','background-color':detail.style.backColor,'font-weight':detail.style.fontWeight,'font-family':detail.style.fontFamily,color: detail.style.foreColor, fontSize : detail.style.fontSize+'px'}" :offset="{width: LabelOffsetX, height: LabelOffsetY}"/>
                   <bm-info-window   :autoPan="true" :title="item.name"  :position="{lng: item.longitude, lat: item.latitude}"   :show="item.showFlag" @close="infoWindowClose(index)">

                        告警数量 : {{ item.AlarmCount }}<br/>
                     状态 : {{ item.Status==2?"未激活":item.Status?"在线":"离线" }}

                     </bm-info-window>
               </bm-marker>

               </template>
             </div>
          </baidu-map>
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
import svgView from '../View';
import {GetDeviceInfo} from '@/services/system';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-device-map',
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        BMap:null,
		BMapObj:null,
        mapUpdate:false,
        LabelOffsetX:50,
        LabelOffsetY:30,
        MakerUrl:"",
        MakerWidth:100,
        infoWindowShow:false,
        MakerHeight:100,
        MakerLng:113.870713,
        MakerLat:22.640808,
        center: {lng: 113.870713, lat: 22.640808},
        zoom: 3,
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
        fillOpacity:1,
        strokeOpacity:1,
        animateType:"blink",
        startColor:"#74f808",
        stopColor:"#74f808",
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        imageURL:"",
        DeviceTimer:null,

        DeviceListMap:[],
        MakerUrlOffline:"",
        MakerUrlAlarm:"",
        base:{
          text: "configComponent.DeviceMap.title",
          "icon": "icon-ditu",
          "isFontIcon": true,
          "info": {
            "type": "image",
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
                "h": 400
              },
              "visible":1,
              "backColor": "transparent",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              text: "标签",
              fontWeight:400,
              textAlign: "center",
              fontSize: 30,
              fontFamily: "Arial",
              "diy":[
                {
                  "name":"component.public.fillOpacity",
                  "type":7,
                  "value":1,
                  "min":0,
                  "max":1,
                  "key":"fillOpacity",
                },
                {
                  "name":"configComponent.Map.PlaneMap.MakerUrlOnline",
                  "type":5,
                  "value":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABkAAAAZCAYAAADE6YVjAAAAAXNSR0IArs4c6QAAAYVJREFUSEu9VVFugzAMtUFj2g069tP+0Z0CTtJykrYnoTtJe4qVv/Ez1htMYyKeUspEE5ukm9RISAgSP79nPwfhBgtvgAFeIJPvJA0oWABQqpMigAoBK4Xq5XhX7l2JjoJMPpMphlggwCk4twhgTy3lx4eykvaIIF32uHNl2f9XSJnEigXRDIIQ33wBfoFamnGMWJDHZr6TJNLyIMAUuudi6X8f0SEzv1sgMgvc1tFr3geIm+cCgJZmQMWwsUDir2QFiGvjcFVHh5kZMG7mWlKD0WUy+owFwkpFtK7vy40J8tQkSwIsht85ybxAECh/j8qtxcSTtS0Xr3WlWsrMzmHlYlhbIJwEZ5efTKffwxBSBbjgOpBjzbYwX1Av17ANMmZG7XbLC2NQkuvFsSK0sowhdCDbwsMovkCS0/tYzlHvUR+2DsNknSDnMSPWZ2z6ejPRGyXZJJM6B6RUWRvInlFXX1rcgUF9nHW4qibDzV19gpVq1Wbsuv2zXF5+FzY5u+s/wfuzP87puxohMqhuAAAAAElFTkSuQmCC",
                  "key":"MakerUrl",
                },
                {
                  "name":"configComponent.Map.PlaneMap.MakerUrlOffline",
                  "type":5,
                  "value":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABkAAAAZCAYAAADE6YVjAAAAAXNSR0IArs4c6QAAAWtJREFUSEu9lW2KwyAQhhUv0t6kOcluDuFAfrX9FdBDbPckzU2aiwwusyi4OqO2LBUCoZp5fN/5qFZvWPoNDDUE8d6fQggfSqlTvNSulNq11t/W2q130SZkXdeDMeYrC87F2xBxXpaFwOwSIfH2994t077WepJUsZCo4DEKSOcQ8cgpYiHOOVKQ/C9ZlINDfKo9AJjKHytIQ8UNAOYUwDlHufosA3JqKoj3/hxCuBQf7wBwLAM658hSUpWvP5ehjQrCWaW1vlhrrwyElJCifG2lZUMQpdQMALcSMqqag3Be74g4lZXD2cWp5iCcBSTit+noxRhDlZdPgFxkpVoqYS6hI23DFkirGalXysppgqSub40VrpRFiFSBbAnnUYTq4UBV2eaHuqNeaLg8BpuHpyBxzIj5aU3fBOoqoYMN29gm7Q5IKbMMqJpR0rdDSrLJm/qnm4encpIfjvk5I+K19Xf7sl0j7f4vdr0K+gEyx7MatVVAEQAAAABJRU5ErkJggg==",
                  "key":"MakerUrlOffline",
                },
                {
                  "name":"configComponent.Map.PlaneMap.MakerUrlAlarm",
                  "type":5,
                  "value":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABkAAAAZCAYAAADE6YVjAAAAAXNSR0IArs4c6QAAAYJJREFUSEu9lVFuwjAMhu0w0A4xqvIEvQU9CXAS4CRlJ6G3AJ6oyg4xlTWeUsbUJnaToYlKvJQkX/7f/l2EJzz4BAYEQQ5vL3Ol1IIA5s2liApELLTW78nHV+67aC/kEL/GWNcZIN4O5x6inAaDVVJ8FtISEWJuj0rtfbe8/09ap5IqFtIo0PocCvgFKTXhFLGQ43i4Fy0iygExBgDz6z5E+exyTe3XDkRSgQC7aVmt7gecolFGAEuHw6hxIKdotCaAjbW5mJXVxD7wGI2MpR1F9mXMHgfCWYUAm2lZbW3IYTxcImLWec9YFgQholVyue5sSKhqzi7O64KUSu3OEexyVLuF5yy4pbwJXaOmrucIsOA6kFPNtzBT0MDMsA3SF0aTdjcLPTQp9eJYEYoqIqQOZFu4fUowSEj6/SzvqOc6KCSo7TVeyM+YEevTN32DlZiFkm1SSL0DUqqsDeJmlLTXa1d7Y6s+bB7+BWLqo7Rea6W2fZ/bh+0KTDy77E92PQr6BjRCyRpHUeOCAAAAAElFTkSuQmCC",
                  "key":"MakerUrlAlarm",
                },
                {
                  "name":"configComponent.Map.PlaneMap.MakerWidth",
                  "type":7,
                  "value":23,
                  "min":0,
                  "max":1000,
                  "key":"MakerWidth",
                },
                {
                  "name":"configComponent.Map.PlaneMap.MakerHeight",
                  "type":7,
                  "value":25,
                  "min":0,
                  "max":1000,
                  "key":"MakerHeight",
                },
                {
                  "name":"configComponent.Map.PlaneMap.LabelOffsetX",
                  "type":7,
                  "value":-20,
                  "min":0,
                  "max":1000,
                  "key":"LabelOffsetX",
                },
                {
                  "name":"configComponent.Map.PlaneMap.LabelOffsetY",
                  "type":7,
                  "value":25,
                  "min":0,
                  "max":1000,
                  "key":"LabelOffsetY",
                },
              ]
            }
          }
        }
      }
    },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          this.initComponents(newVal);
        },
        deep: true
      }
    },
    methods: {
      mapHandler ({BMap, map}) {
        this.BMap = map
		    this.BMapObj = BMap
        this.initComponents(this.detail);
        this.zoom = 15
      },
      infoWindowClose(index){
        this.$set(this.DeviceListMap[index], "showFlag", false);
      },
      lookDetail(index) {
        this.$set(this.DeviceListMap[index], "showFlag", true);
      },
      GetDeviceInfo(){
        let _t = this
		_t.DeviceListMap = []
        GetDeviceInfo().then(function (res){
          if(res.data.code==0)
          {
            let points = []
            for(let i = 0;i<res.data.list.length;i++)
            {
              res.data.list[i].showFlag = false

              if((_t.BMap!=null)&&(res.data.list[i].longitude!="")&&(res.data.list[i].latitude!="")) {
                points.push(new _t.BMapObj.Point(res.data.list[i].longitude, res.data.list[i].latitude))
              }
            }
            _t.DeviceListMap = res.data.list
            if(_t.BMap!=null)
            {
              _t.BMap.setViewport(points)
            }

          }
        })
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.mapUpdate = false
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="strokeWidth")
          {
            this.strokeWidth=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeFill")
          {
            this.fill=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeColor")
          {
            this.strokeColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="fillOpacity")
          {
            this.fillOpacity=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeOpacity")
          {
            this.strokeOpacity=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerUrl")
          {
            this.MakerUrl=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerUrlOffline")
          {
            this.MakerUrlOffline=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerUrlAlarm")
          {
            this.MakerUrlAlarm=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerWidth")
          {
            this.MakerWidth=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerHeight")
          {
            this.MakerHeight=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="LabelOffsetX")
          {
            this.LabelOffsetX=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="LabelOffsetY")
          {
            this.LabelOffsetY=option.style.diy[i].value
          }
        }
        this.center.lng = this.MakerLng
        this.center.lat = this.MakerLat
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
		this.GetDeviceInfo()
        this.mapUpdate = true
      }
    },
    beforeDestroy(){
      clearInterval(this.DeviceTimer)
      // 清理 EventBus 监听
      if (this._activeEventHandler) {
        this.$EventBus.$off(this.detail.identifier + 'activeEvent', this._activeEventHandler)
      }
      if (this._animateEventHandler) {
        this.$EventBus.$off(this.detail.identifier + 'animateEvent', this._animateEventHandler)
      }
      if (this._cellEditModeHandler) {
        this.$EventBus.$off('cell-editMode', this._cellEditModeHandler)
      }
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
        if(!_t.editMode)
        {
          this.DeviceTimer = setInterval(function (){
            _t.GetDeviceInfo()
          },60000)
        }

        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, _t._activeEventHandler = (data) => {


        })
        _t.$EventBus.$on(animateEvent, _t._animateEventHandler = (data) => {
          _t.isStart = data
        })

      });
    },
    created(){
      let _t = this
      this.GetNodeObj = this.getNode()
      this.GetNodeObj.on('change:data', ({ current }) => {
        if(current) {
          _t.detail = current.detail
        }
      })
      this.GetNodeObj.on('change:size', ({ current }) => {
        _t.detail.style.position.w = current.width
        _t.detail.style.position.h = current.height
      });
      this.detail = this.GetNodeObj.getData().detail
      this.editMode = this.GetNodeObj.getData().editMode
      this.showDeviceUuid = this.GetNodeObj.getData().showDeviceUuid
      this.IsToolBox = this.GetNodeObj.getData().IsToolBox
      _t.$EventBus.$on('cell-editMode', _t._cellEditModeHandler = (data) => {
        _t.editMode = data.edit
        _t.IsToolBox = data.toolbox
		_t.initComponents(_t.detail)
      })
    }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}

.bm-view {
  z-index:-1000
}
/*地图标题*/
.BMap_bubble_title {
  color:white;
  font-size:13px;
  font-weight:bold;
  text-align:left;
  padding-left:5px;
  padding-top:5px;
  border-bottom:1px solid gray;
  background-color:#0066b3;
}
/* 消息内容 */
.BMap_bubble_content {
  background-color:white;
  padding-left:5px;
  padding-top:5px;
  padding-bottom:10px;
}
/* 内容 */
.BMap_pop div:nth-child(9) {
  top:35px !important;
  border-radius:7px;
}
/* 左上角删除按键 */
.BMap_pop img {
  top:43px !important;
  left:215px !important;
}
.BMap_top {
  display:none;
}
.BMap_bottom {
  display:none;
}
.BMap_center {
  display:none;
}
/* 隐藏边角 */
.BMap_pop div:nth-child(1) div {
  display:none;
}
.BMap_pop div:nth-child(3) {
  display:none;
}
.BMap_pop div:nth-child(7) {
  display:none;
}
</style>

<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject :style="{overflow:'visible','z-index':5}" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
          <baidu-map  @ready="mapHandler"  :scrollWheelZoom="true" :inertialDragging="false" :center="center" :zoom="zoom" :style="{'z-index':0,width:detail.style.position.w+'px' ,height:detail.style.position.h+'px'}">
             <div v-if="mapUpdate">
               <bm-geolocation anchor="BMAP_ANCHOR_BOTTOM_RIGHT" :showAddressBar="true" :autoLocation="true"></bm-geolocation>
               <bm-marker :position="{lng: MakerLng, lat: MakerLat}" :dragging="false" animation="BMAP_ANIMATION_DROP" :icon="{url: MakerUrl, size: {width: MakerWidth, height: MakerHeight}}">
                <bm-label :content="detail.style.text"  :labelStyle="{border:'','background-color':detail.style.backColor,'font-weight':detail.style.fontWeight,'font-family':detail.style.fontFamily,color: detail.style.foreColor, fontSize : detail.style.fontSize+'px'}" :offset="{width: LabelOffsetX, height: LabelOffsetY}"/>
               </bm-marker>
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
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-map-1',
    inject: ['getNode'],
    data() {
      return {
        detail:{},
        IsToolBox:false,
        editMode:true,
        mapUpdate:false,
        LabelOffsetX:50,
        LabelOffsetY:30,
        MakerUrl:"",
        MakerWidth:100,
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
        base:{
          text: "configComponent.Map.PlaneMap.title",
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
                "w": 800,
                "h": 560
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
                  "name":"configComponent.Map.PlaneMap.MakerLng",
                  "type":7,
                  "value":113.870713,
                  "min":0,
                  "max":1000,
                  "key":"MakerLng",
                },
                {
                  "name":"configComponent.Map.PlaneMap.MakerLat",
                  "type":7,
                  "value":22.640808,
                  "min":0,
                  "max":1000,
                  "key":"MakerLat",
                },
                {
                  "name":"configComponent.Map.PlaneMap.MakerUrl",
                  "type":5,
                  "value":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABcAAAAZCAYAAADaILXQAAAACXBIWXMAAC4jAAAuIwF4pT92AAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAArlJREFUeNqllt1LU3EYx49Nj5EmLdfWZtHM2XHM7eyssjIhh60ta8WMdmEhgrdR3nTRVQODqFisS/sDgiIoiYgwyKAiki6ySE2DmmZMZ9uqIfTCr+/vdAZy9lt78eJzsbfP8+x5nj3PuFAoxKl5bnKVAx2wPTVJXSNGZ/eDjY6u24ZmYchgr2V9hgVLLIBLYOLVZomMmZ3kDRgzS2QUj/H8C3AWbCpYrmR7AsxObZNIrL2DpIJBkurpIQmfj3xtaSGLdjv5bHWQiQYnDfIS7CtUfhr8/NK+n/wYGCDf+/tJwu2WhSw+CiINEAPu/8rxBh9ILnT6SToSIUm/n8RzSBkBpoGFKccLevAoKraQdDhMEh5PXulyphvlABHAs+Q069+p3l6SDASKElPizXYqXwImlvzilEUk3/r6ihZneFsvN7ibJb8XFXeShNdbslwpzVWWfGhmeytZFMWS5fSbw3OFKf8klC7OJ78fFRwrkmcmhiUfeG9ZWeZ0PcATZMk9r7dIJYsX/o1iItcomsBItMmxkpIMgspcP/8ASNMsihHPYok9M7nm7hrsDo7jKoAGrGItrmvF1J7uHoh/XdcJ5yBsAg2gDqwDPChbLjeCUfo1C5FPYu3e1NseGjT8EYh8wAs6wG5QD1arD4UVzM/kqf8HJIDLNCny1Sch8avoBHvB+qw1ibMWRIA/MVuOOiPwY6MzfrxKf4YhphwEO0BN9g7muPI7hubw+FZnlngeARF86bzWfJkhPQwOgF1gg9xchlzjqlxrxiEenlI1mAYc1Ak3VEJaBg9oAxZQTR3MAy13Gc04VVO3Z9gojmcaTO/mLb1tpFZTcQyvH1Ia6AYuZUqqlHEsy3n9lwXgL2jr254Ypbl3yBgNHDdq+KPKRLSCRqCVp0LJNO9fC3X9qQANJsqIWZV6rlFnyeIvmIDJtcarlTEAAAAASUVORK5CYII=",
                  "key":"MakerUrl",
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
        this.initComponents(this.detail);
        this.zoom = 15
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
          else if(option.style.diy[i].key=="MakerLat")
          {
            this.MakerLat=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerWidth")
          {
            this.MakerWidth=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerHeight")
          {
            this.MakerHeight=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MakerLng")
          {
            this.MakerLng=option.style.diy[i].value
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

        this.mapUpdate = true
      }
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {


        })
        _t.$EventBus.$on(animateEvent, (data) => {
          _t.isStart = data
        })
      });
    },
    created(){
      let _t = this
      const node = this.getNode()
      node.on('change:data', ({ current }) => {
        if(current) {
          _t.detail = current.detail
        }
      })
      node.on('change:size', ({ current }) => {
        _t.detail.style.position.w = current.width
        _t.detail.style.position.h = current.height
      });
      this.detail = node.getData().detail
      this.editMode = node.getData().editMode
      this.showDeviceUuid = node.getData().showDeviceUuid
      this.IsToolBox = node.getData().IsToolBox
      _t.$EventBus.$on('cell-editMode', (data) => {
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
</style>

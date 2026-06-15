<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
       <div :style="animatedStyle" v-show="detail.style.visible==1 ||isStart? true:false">
          <div :class="{
                'animated':true,[`${detail.style.animate}`]: true,
                'color-animation':isStart&&animateType.includes('millcolorGrad')&&!IsToolBox,
                'blink-animation':isStart&&animateType.includes('blink')&&!IsToolBox,
                'scale-animation':isStart&&animateType.includes('Zoom')&&!IsToolBox,
                'rotate-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,
                'rotate-anti-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1
              }"
               :style="{
                                      width: detail.style.position.w + 'px',
                                      height: detail.style.position.h + 'px',
                                      'background-color': detail.style.backColor,
                                      'border-radius':detail.style.BorderEdges+'px',
                                      opacity:detail.style.opacity,
                                      borderWidth: detail.style.borderWidth + 'px',
                                      borderStyle: detail.style.borderStyle,
                                      borderColor: detail.style.borderColor,
                                      transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
                                  }">
          <a-layout :style="styleVar">
            <a-layout-header  :bodyStyle="{padding:'2px',}" :panelStyle="{width:'100%',height:'100%'}">
              <a-form layout="inline" :label-col="{ span: 8}" :wrapper-col="{ span: 16 }">
                <a-form-item :label="$t('systemHistoryVideoModel.deviceList')" >
                  <a-select v-model="selectDevice" @change="changeDeviceEvent" style="width: 200px" :dropdownStyle="{'z-index': 9999999}">
                    <a-select-option  v-for="(deviceValue,index) in deviceList" :key="index" :value="deviceValue">
                      {{ deviceValue }}
                    </a-select-option>
                  </a-select>
                </a-form-item>

                <a-form-item :label="$t('systemHistoryVideoModel.DateList')" >
                  <a-select v-model="selectDate" @change="changeDeviceDateEvent" style="width: 200px" :dropdownStyle="{'z-index': 9999999}">
                    <a-select-option  v-for="(dateValue,index) in deviceDateList" :key="index" :value="dateValue">
                      {{ dateValue }}
                    </a-select-option>
                  </a-select>
                </a-form-item>

                <a-form-item label=" ">
                  <a-button  @click="getSystemHistoryImagesList"   type="primary">  {{$t('configComponent.video.VideoRefresh')}} </a-button>
                </a-form-item>

              </a-form>
            </a-layout-header>
            <a-layout-content style="background: #fff;"  :bodyStyle="{padding:'10px','min-height':'300px'}" :panelStyle="{width:'100%',height:'100%'}">
              <swiper :options="mySwiperOption">
                <swiper-slide v-for=" item in deviceFileList" :key="item.key">
                  <img :src="item.Path" :width="width+'px'" :height="(height-50)+'px'" />
                </swiper-slide>
                <div class="swiper-button-prev" slot="button-prev"></div>
                <div class="swiper-button-next" slot="button-next"></div>
                <div class="swiper-pagination" slot="pagination"></div>
              </swiper>
            </a-layout-content>
          </a-layout>
        </div>
       </div>
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
import { historyImagesList} from "@/services/video";

import 'swiper/dist/css/swiper.css'
import { swiper, swiperSlide } from 'vue-awesome-swiper'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-history-images-list',
  inject: ['getNode'],
  i18n: require('../../../../i18n/language'),
  props: {

  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
          if(this.editMode) {
            this.initComponents(newVal);
          }
      },
      deep: true
    }
  },
  components: {
    swiper,
    swiperSlide
  },
  data() {
    return {
      detail:null,
      IsToolBox:false,
      editMode:true,
      isStart:false,
      animateType:"blink",
      spinDirection:0,
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      blinkSpeed:0.5,
      startColor:"#74f808",
      stopColor:"#74f808",
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      fillOpacity:1,
      strokeOpacity:1,
      deviceList:[],
      selectDevice:"",
      deviceDateList:[],
      deviceFileList:[],
      selectDate:"",
      mySwiperOption: {
        pagination: {
          el: '.swiper-pagination',
          type: 'fraction'
        },
        navigation: {
          nextEl: '.swiper-button-next',
          prevEl: '.swiper-button-prev'
        },
        //自动播放
        autoplay: {
          delay: 1000,
          disableOnInteraction: false
        },
        //循环
        loop:true
      },
      videoComponentData:{},
      token:"",
      videoStyle:{
        height:300+"px",
        width:'auto'
      },
      foreColor:"#000000",
      backColor:"#ffffff",
      tableSplitColor:"#000",
      tableHoverColor:"#fff",
      tableHeaderColor:"",
      tableHeaderBackColor:"",
      SearchColor:"#000000",
      SearchBorderColor:"#108ec4",
      SearchBackColor:"#ffffff",
      dateSelectColor:"",
      dateSelectBackColor:"",
      dateSelectBorderColor:"",
      VideoStatus:0,
      VideoTitle:"",
      videoType:1,
      video_url:"",
      temp_video_url:"",
      width:600,
      height:600,
      base:{
        "text": "configComponent.video.historyImagesList",
        "icon": "icon-zhuapaitupian",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "active": [],
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
              "w": 680,
              "h": 300
            },
            "visible":1,
            "backColor": "#ffffff",
            "foreColor": "#000000",
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.DeviceTree.SearchColor",
                "type":2,
                "value":"#000000",
                "key":"SearchColor",
              },
              {
                "name":"configComponent.DeviceTree.SearchBackColor",
                "type":2,
                "value":"#ffffff",
                "key":"SearchBackColor",
              },
              {
                "name":"configComponent.DeviceTree.SearchBorderColor",
                "type":2,
                "value":"#cbc6c6",
                "key":"SearchBorderColor",
              },
            ]
          }
        }
      }
    }
  },
  methods: {
    changeDeviceEvent(device){
      this.selectDevice = device
      this.selectDate=""
      let tempDate = []
      for(let i=0;i<this.videoList.length;i++)
      {
        if(this.videoList[i].Device==device)
        {
          tempDate.push(this.videoList[i].Date)
        }
      }
      this.deviceDateList = tempDate.filter((value, index, array) => array.indexOf(value) === index);
    },
    changeDeviceDateEvent(date){
      this.selectDate = date
      this.deviceFileList=[]
      for(let i=0;i<this.videoList.length;i++)
      {
        if(this.videoList[i].Device==this.selectDevice&&(this.videoList[i].Date==date))
        {
          let img = {}
          img.key =  this.videoList[i].FileName
          img.title =  this.videoList[i].FileName
          img.Path =  this.videoList[i].Path
          this.deviceFileList.push(img)
        }
      }
    },

    onSelectVideo(token){
      this.video_url = token.url
      this.VideoTitle = token.name
      this.visible =false
    },
    changeDeviceModel(ev,event){

    },
    getSystemHistoryImagesList(){
      let _t = this
      _t.imagesList=[]
      let tableData={}
      _t.videoList=[]
      this.selectDevice=""
      this.selectDate=""
      this.deviceList=[]
      this.deviceDateList=[]
      this.deviceFileList=[]

      historyImagesList().then(function (res){
        if(res.data.files!=null)
        {
          let tableData={}
          let deviceListTemp = []
          let list = res.data.files
          for(let i=0;i<list.length;i++)
          {
            let splitArray =  list[i].Path.split("/");
            tableData.key = list[i].Path
            tableData.Name = splitArray[4]
            tableData.Device = splitArray[2]
            tableData.Date = splitArray[3]
            tableData.FileName = splitArray[4]
            tableData.Path = list[i].Path
            tableData.Size = (list[i].Size/(1024*1024)).toFixed(2)
            _t.videoList.push(tableData)
            deviceListTemp.push(tableData.Device)
            tableData={}
          }
          _t.deviceList = deviceListTemp.filter((value, index, array) => array.indexOf(value) === index);
        }
        if(_t.deviceList.length>0)
        {
          _t.changeDeviceEvent(_t.deviceList[0])
          if(_t.deviceDateList.length>0)
          {
            _t.changeDeviceDateEvent(_t.deviceDateList[0])
          }
        }
      })

    },
    initComponents(option){
      this.width = option.style.position.w
      this.height = option.style.position.h
      if(this.IsToolBox)
      {
        return
      }
      let i=0
      this.backColor=this.detail.style.backColor
      this.foreColor=this.detail.style.foreColor
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="hoverColor")
        {
          this.hoverColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="selectedColor")
        {
          this.selectedColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="TextFontSize")
        {
          this.TextFontSize=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="selectedTextColor")
        {
          this.selectedTextColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="SearchColor")
        {
          this.SearchColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="SearchBackColor")
        {
          this.SearchBackColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="SearchBorderColor")
        {
          this.SearchBorderColor=option.style.diy[i].value
        }

        else if(option.style.diy[i].key=="dateSelectColor")
        {
          this.dateSelectColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="dateSelectBackColor")
        {
          this.dateSelectBackColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="dateSelectBorderColor")
        {
          this.dateSelectBorderColor=option.style.diy[i].value
        }

        else if(option.style.diy[i].key=="tableHeaderColor")
        {
          this.tableHeaderColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="tableHeaderBackColor")
        {
          this.tableHeaderBackColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="tableSplitColor")
        {
          this.tableSplitColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="tableHoverColor")
        {
          this.tableHoverColor=option.style.diy[i].value
        }

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
    }
  },
  mounted() {
    let _t = this
    this.getSystemHistoryImagesList()
    this.$nextTick(function(){
      this.initComponents(this.detail);
      if(true){
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {

        })
        _t.$EventBus.$on(animateEvent, (data) => {
          _t.isStart = data
        })
      }
    });
  },
  destroyed(){
    this.video_url=""
  },
  computed: {
    styleVar() {
      return {
        "height": this.detail.style.position.h+'px',
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

        '--dateSelectColor': this.dateSelectColor,
        '--dateSelectBackColor': this.dateSelectBackColor,
        '--dateSelectBorderColor': this.dateSelectBorderColor,

        '--tableHeaderColor': this.tableHeaderColor,
        '--tableHeaderBackColor': this.tableHeaderBackColor,
        '--tableSplitColor': this.tableSplitColor,
        '--tableHoverColor':this.tableHoverColor
      };
    },
    animatedStyle(){
      return {
        "--blinkSpeed":this.blinkSpeed+'s',
        "--stopColor":this.stopColor,
        "--startColor":this.startColor,
        "--animateSpeed":this.animateSpeed+'s',
        "--animateSpinSpeed":this.animateSpinSpeed+'s'
      }
    },
    textAlign: function(){
      if(this.detail.style.textAlign == undefined) {
        return "center";
      } else {
        return this.detail.style.textAlign;
      }
    },
    lineHeight: function() {
      if(this.detail.style.lineHeight == undefined) {
        return this.detail.style.position.h;
      }
      return this.detail.style.lineHeight;
    },
    CardStyle:function () {
      let styles = [];
      if(this.detail.style.backColor) {
        styles.push(`background-color: ${this.detail.style.backColor}`);
      }
      if(this.detail.style.foreColor) {
        styles.push(`color: ${this.detail.style.foreColor}`);
      }
      let style = styles.join(';');
      return style;
    },
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
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
      _t.initComponents(_t.detail)
    })
    this.initComponents(this.detail);
  }

}
</script>

<style lang="less" scoped>
.ant-layout-header {
  height: 50px;
  padding: 0 2px;
  line-height: 50px;
  background: var(--backColor);
}

::v-deep .ant-form-item-label > label {
  color: var(--foreColor)
}
::v-deep .ant-card {
  background-color: var(--backColor);
}
::v-deep  .ant-table-placeholder {
  background:var(--backColor);
  border-top: 1px solid var(--backColor);
  border-bottom: 1px solid var(--backColor);
  border-radius: 0 0 4px 4px;
}
::v-deep  .ant-btn-primary {
  color: var(--foreColor);
  background-color: var(--SearchBackColor);
  border-color: var(--SearchBorderColor);
}
::v-deep .ant-input .ant-radio-button-wrapper {
  background-color: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}
::v-deep  .ant-select-selection {
  background-color: var(--SearchBackColor)
}
::v-deep  .ant-select-selection {
  border: 1px solid var(--SearchBorderColor);
}

</style>

<template>
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
       <a-select
           placeholder="Please select"
           class="eia-select-class"
           :dropdownStyle="styleVar"
           v-model="StatusValue"
           dropdown-class-name="eia-dropdown-class"
           @change="handleChange"
       >
          <a-select-option :value="item.value"  v-for="(item,index) in StatusList" :key="index">
            {{item.Text}}
          </a-select-option>
        </a-select>

  </div>
    </div>
</template>

<script>
import svgView from '../View';
import {setData} from "@/services/device";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-com-box-status',
    i18n: require('../../../../i18n/language'),
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        Text:"",
        foreColor:"#000000",
        backColor:"#ffffff",
        hoverForeColor:"#000000a6",
        hoverBackColor:"#DAEAF6",
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
        StatusList:[],
        StatusValue:0,
        base:{
          text: "configComponent.ComBox.title",
          "icon": "icon-erjixialakuang",
          "isFontIcon": true,
          "info": {
            "type": "image",
            "action": [],
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
            "active": [
              {
                id:"ControlStatus",
                name:"configComponent.status.ControlStatus",
                result:0,
                isStatus:true,
                isSwitch:false,
                isImageStatus:false,
                isTextStatus:false,
                isExpression:false,
                isComBox:true,
                condition:{
                  deviceSN:"",
                  isBandDevice:false,
                  bandType:1,
                  dataID: "",
                  dataName: "",
                  IsManual:false,
                  StatusList:[
                    {

                      "Text":"Text",
                      "value":0
                    }
                  ]
                },
              }
            ],
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 170,
                "h": 50
              },
              "visible":1,
              "backColor": "#ffffff",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              fontSize: 30,
              textAlign: "center",
              fontFamily: "Arial",
              italic:0,
              fontWeight:400,
              "diy":[
                {
                  "name":"configComponent.ComBox.hoverForce",
                  "type": 2,
                  "value": "#000000",
                  "key":"hoverForce",
                },
                {
                  "name":"configComponent.ComBox.hoverBack",
                  "type": 2,
                  "value": "#6ACBFF",
                  "key":"hoverBack",
                },
              ]
            }
          }
        }
      }
    },
    computed: {
      animatedStyle(){
        return {
          "--tableBorderColor": this.tableBorderColor ,
          "--blinkSpeed":this.blinkSpeed+'s',
          "--stopColor":this.stopColor,
          "--startColor":this.startColor,
          "--animateSpeed":this.animateSpeed+'s',
          "--animateSpinSpeed":this.animateSpinSpeed+'s',
          '--fontWeight':this.detail.style.fontWeight,
          '--fontSize':this.detail.style.fontSize+ 'px',
          '--fontFamily':this.detail.style.fontFamily,
          '--fontStyle':this.detail.style.italic?'oblique':'normal',
          "--height": this.detail.style.position.h+'px',
          "--width": this.detail.style.position.w+'px',
          "--lineHeight": (this.lineHeight-5)+'px' ,
          "--foreColor": this.foreColor ,
          '--backColor':this.backColor,
          '--textAlign':this.detail.style.textAlign,
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
      styleVar() {
        return {
          "--hoverForeColor": this.hoverForeColor ,
          '--hoverBackColor':this.hoverBackColor,
          '--fontWeight':this.detail.style.fontWeight,
          '--fontSize':this.detail.style.fontSize+ 'px',
          '--fontFamily':this.detail.style.fontFamily,
          '--fontStyle':this.detail.style.italic?'oblique':'normal',
          "--height": this.detail.style.position.h+'px',
          "--width": this.detail.style.position.w+'px',
          "--lineHeight": (this.lineHeight-5)+'px' ,
          "--foreColor": this.foreColor ,
          '--backColor':this.backColor,
          '--textAlign':this.detail.style.textAlign,
        };
      },
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
    methods: {
      ManualSetData(value){
        let _t = this

        this.settingLoading=true
        let params = {
          deviceUuid:this.detail.active[0].condition.deviceSN,
          dataUuid:this.detail.active[0].condition.dataID,
          value:value.toString(),
        };

        setData(params).then(function (res){
          if(res.data.code==0)
          {
            _t.$message.success(_t.$t("readData.SetSuccess"))
          }
          else
          {
            _t.$message.error(_t.$t("readData.SetFailed"))
          }
        }).catch(function (error) {
          _t.settingLoading = false
          _t.$message.error(_t.$t("readData.SetFailed"))
        }).finally(function (error) {
          _t.settingLoading = false
        })
      },
      handleChange(value) {
        let  _this = this
        if((!this.editMode))
        {
          _this.$confirm({
            title: _this.$t('displayConfig.Properties.SecondConfirm'),
            content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
            onOk() {
              _this.ManualSetData(value)
            },
            onCancel() {

            }
          });

        }
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.backColor=this.detail.style.backColor
        this.foreColor=this.detail.style.foreColor
        for(let i=0;i<option.active[0].condition.StatusList.length;i++)
        {
          option.active[0].condition.StatusList[i].value = parseFloat(option.active[0].condition.StatusList[i].value)
        }
        this.StatusValue = option.active[0].condition.StatusList[0].value
        this.StatusList = option.active[0].condition.StatusList
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="hoverBack")
          {
            this.hoverBackColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="hoverForce")
          {
            this.hoverForeColor=option.style.diy[i].value
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
      this.$nextTick(function(){
       this.initComponents(this.detail);
          let activeEvent = this.detail.identifier+"activeEvent"//动作数据
          let animateEvent = this.detail.identifier+"animateEvent"//动作数据

          _t.$EventBus.$on(activeEvent, (data) => {
            if(data.ID == "ControlStatus")
            {
              let temp  =  parseFloat(data.result)
              if(!isNaN(temp))
              {
                _t.StatusValue=temp
              }
            }
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
      })
      this.initComponents(this.detail);
    }
}
</script>
<style lang="less">

/* 使用animation关键帧 */
.color-animation {
  animation: colorChange var(--animateSpeed) linear infinite;
}

@keyframes colorChange {
  0% { background-color: var(--startColor); }
  100% { background-color:  var(--stopColor); }
}
/* 使用animation关键帧 */
.blink-animation {
  animation: blink var(--blinkSpeed) linear infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
/*缩放*/
.scale-animation {
  animation: pulse 0.6s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.5); }
}
/*顺时针旋转*/
.rotate-animation {
  animation: clockwiseRotate var(--animateSpinSpeed) linear infinite;
  transform-origin: center;
}

@keyframes clockwiseRotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
/*逆时针旋转*/
.rotate-anti-animation {
  animation: counterClockwiseRotate var(--animateSpinSpeed) linear infinite;
  transform-origin: center;
}

@keyframes counterClockwiseRotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(-360deg); }
}
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
.eia-dropdown-class {
  .ant-select-dropdown-menu-item {
    font-family: var(--fontFamily);
    font-size: var(--fontSize);
    text-align: var(--textAlign);
    line-height: var(--lineHeight);
    font-weight:var(--fontWeight);
    background-color:  var(--backColor);
    color: var(--foreColor);
  }
  .ant-select-dropdown-menu-item:hover{
    color: var(--hoverForeColor);
    background-color: var(--hoverBackColor);
  }
  .ant-select-dropdown-menu{
    background-color:  var(--backColor);
  }
}
.eia-select-class {
  .ant-select-selection {
    width: var(--width);
    font-family: var(--fontFamily);
    font-size: var(--fontSize);
    font-weight:var(--fontWeight);
    color: var(--foreColor);
    font-style: var(--fontStyle);
    display: flex;
    justify-content: var(--textAlign);;
    background-color: var(--backColor);
    border: 1px solid var(--backColor);
  }
  .ant-select-selection--single {
    height: var(--height);
  }
  .ant-select-selection__placeholder {
    color: var(--foreColor);
  }
  .ant-select-selection__rendered {
    line-height: var(--height);
  }
}
</style>

<template>
  <div style="color: red; font-size: 24px">
    {{ detail && detail.style ? detail.style.text : 'NO DETAIL' }}
  </div>
</template>

<script>

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-text',
    inject: ['getNode'],
    data() {
      return {
        Text:"",
        DivOpacity:1,
        animateType:"blink",
        startColor:"#74f808",
        stopColor:"#74f808",
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        italic:false,
        imageURL:"",
        detail: { style: {}, animate: { selected: [], animateElement: [] } },
        IsToolBox:false,
        editMode:true,
        base:{
          text: "configComponent.label.Text",
          "icon": "icon-icon_svg_wenben",
          "isFontIcon": true,
          "info": {
            "type": "view-svg-text",
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
                "w": 100,
                "h": 40
              },
              "visible":1,
              "backColor": "transparent",
              "foreColor": "#000000",
              fontWeight:400,
              "zIndex": -1,
              "transform": 0,
              text: "标签",
              textAlign: "center",
              fontSize: 30,
              fontFamily: "Arial",
              letterSpacing:0,
              italic:0,
              "diy":[

              ]
            }
          }
        }
      }
    },
    computed: {
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
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.DivOpacity = option.style.opacity
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
          else if(option.style.diy[i].key=="imageURL")
          {
            this.imageURL=option.style.diy[i].value
          }
        }
        i=0
        this.animateType = option.animate.selected || []
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
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }

        })
        _t.$EventBus.$on(animateEvent, (data) => {
          console.log("data",data)
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          console.log("data222222",data)
          _t.isStart = data
        })

      });
    },
    created(){
    let _t = this
    try {
      this.GetNodeObj = this.getNode()
      if (!this.GetNodeObj) {
        console.warn('[ViewSvgText] getNode returned null')
        return
      }
      const nodeData = this.GetNodeObj.getData() || {}
      this.detail = nodeData.detail || { style: { diy: [] }, animate: { selected: [], animateElement: [] } }
      this.editMode = nodeData.editMode || false
      this.showDeviceUuid = nodeData.showDeviceUuid || ''
      this.IsToolBox = nodeData.IsToolBox || false
      this.GetNodeObj.on('change:data', ({ current }) => {
        if(current && current.detail) {
          _t.detail = current.detail
        }
      })
      this.GetNodeObj.on('change:size', ({ current }) => {
        if (_t.detail && _t.detail.style && _t.detail.style.position) {
          _t.detail.style.position.w = current.width
          _t.detail.style.position.h = current.height
        }
      });
      _t.$EventBus.$on('cell-editMode', (data) => {
        _t.editMode = data.edit
        _t.IsToolBox = data.toolbox
      })
      this.initComponents(this.detail);
    } catch(e) {
      console.error('[ViewSvgText] created error:', e)
    }
  }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}

</style>

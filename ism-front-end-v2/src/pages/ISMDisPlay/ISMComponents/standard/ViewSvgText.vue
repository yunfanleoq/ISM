<template>
  <div xmlns="http://www.w3.org/1999/xhtml"
       v-show="detail.style.visible==1||isStart"
       @click="onTextClick"
       :style="{
         width: (detail.style.position && detail.style.position.w ? detail.style.position.w : 100) + 'px',
         height: (detail.style.position && detail.style.position.h ? detail.style.position.h : 40) + 'px',
         cursor: hasClickAction ? 'pointer' : 'default',
         overflow: 'hidden',
         display: 'flex',
         alignItems: 'center',
         fontSize: (detail.style.fontSize || 14) + 'px',
         fontFamily: detail.style.fontFamily || 'Microsoft YaHei, PingFang SC, sans-serif',
         fontWeight: detail.style.fontWeight || 400,
         color: detail.style.foreColor || '#e2e8f0',
         textAlign: textAlign,
         whiteSpace: 'nowrap',
         textOverflow: 'ellipsis'
       }">
    {{ detail.style.text }}
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
        animateType:[],
        startColor:"#74f808",
        stopColor:"#74f808",
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        italic:false,
        imageURL:"",
        detail: { style: { diy: [], position: { w: 100, h: 40 }, visible: 1, text: '' }, animate: { selected: [], animateElement: [] }, action: [] },
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
      hasClickAction() {
        const actions = (this.detail && this.detail.action) || []
        return actions.some(a => a && a.type === 'click' && a.action === 'link')
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
      onTextClick() {
        if (this.IsToolBox || this.editMode) {
          return
        }
        const actions = (this.detail && this.detail.action) || []
        const clickAction = actions.find(
          a => a && a.type === 'click' && a.action === 'link' && a.link
        )
        if (!clickAction || !clickAction.link) {
          return
        }
        const link = clickAction.link
        this.$EventBus.$emit('GoPage', {
          IsPopUp: link.isPopUp,
          autoClose: link.autoClose,
          linkType: link.linkType,
          ModelId: link.Inside && link.Inside.displayUUID,
          PageUuid: link.Inside && link.Inside.pageUUID,
          width: link.width,
          height: link.height,
          External: link.External,
          title: link.title,
          OpenExternalType: link.OpenExternalType
        })
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        if (!option || !option.style) return
        this.DivOpacity = option.style.opacity
        const diy = option.style.diy || []
        let i=0
        for( i=0;i<diy.length;i++)
        {
          if(diy[i].key=="strokeWidth")
          {
            this.strokeWidth=diy[i].value
          }
          else if(diy[i].key=="strokeFill")
          {
            this.fill=diy[i].value
          }
          else if(diy[i].key=="strokeColor")
          {
            this.strokeColor=diy[i].value
          }
          else if(diy[i].key=="fillOpacity")
          {
            this.fillOpacity=diy[i].value
          }
          else if(diy[i].key=="strokeOpacity")
          {
            this.strokeOpacity=diy[i].value
          }
          else if(diy[i].key=="imageURL")
          {
            this.imageURL=diy[i].value
          }
        }
        i=0
        this.animateType = (option.animate && option.animate.selected) || []
        if(option.animate && option.animate.isExpression)
        {
          this.isStart = false
        }
        else
        {
          this.isStart = true
        }
      }
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
       this.initComponents(this.detail);
        if (!this.detail || !this.detail.identifier) return
        let activeEvent = this.detail.identifier+"activeEvent"
        let animateEvent = this.detail.identifier+"animateEvent"

        _t.$EventBus.$on(activeEvent, (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          if (data && data.result !== undefined) {
            let value = data.result
            if (typeof value === 'boolean') value = value ? '1' : '0'
            _t.detail.style.text = String(value)
          }
        })
        _t.$EventBus.$on(animateEvent, (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          _t.isStart = data
        })

      });
    },
    created(){
    let _t = this
    try {
      this.GetNodeObj = this.getNode()
      if (!this.GetNodeObj) {
        return
      }
      const nodeData = this.GetNodeObj.getData() || {}
      this.detail = nodeData.detail || this.detail
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
  transform-origin: center center;
}

</style>

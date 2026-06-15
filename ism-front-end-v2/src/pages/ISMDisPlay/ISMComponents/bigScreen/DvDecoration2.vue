<template>

 <dv-decoration-2 :dur="dur"  :reverse="reverse" :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-decoration-2>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-decoration2',
  inject: ['getNode'],
    props: {

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
    })
  },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          this.setOption(newVal);
        },
        deep: true
      }
    },
    data() {
        return {
          detail:null,
          IsToolBox:false,
          editMode:true,
          base:{
            text: "configComponent.bigScreen.embellish.dvDecoration2",
            icon: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAUcAAABgCAYAAACOl7BRAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAFsSURBVHhe7dehDcJQAEVRxqnA0a0YggFIEFWdo5rdPpKUXoGGI84I7ybvNJ3nAcCeOAIEcQQI4ggQxBEgiCNAEEeAII4AQRwBgjgCBHEECOIIEMQRIIgjQBBHgCCOAEEcAYI4AgRxBAjiCBDEESCII0AQR4AgjgBBHAGCOAIEcQQI4ggQxBEgiCNAEEeAII4AQRwBgjgCBHEECOIIEMQRIIgjQBBHgCCOAEEcAYI4AgRxBAjiCBDEESCII0AQR4AgjgDhp+N4ud6AD+v2HPfHkpvhTRzhz4jjd9xqgCCOAEEcAYI4AgRxBAjiCBDEESCII0AQR4AgjgBBHAGCOAIEcQQI4ggQxBEgiCNAEEeAII4AQRwBgjgCBHEECOIIEMQRIIgjQBBHgCCOAEEcAYI4AgRxBAjiCBDEESCII0AQR4AgjgBBHAGCOAIEcQQI4ggQxBEgiCNAEEeAII4AQRwBgjgCBHEECOIIcDCPF1ojQSN+bOWxAAAAAElFTkSuQmCC",
            isFontIcon: true,
            info: {
              type: "text",
              action: [],
              dataBind:
                [
                ],
              style: {
                position: {
                  x: 0,
                  y: 0,
                  w: 300,
                  h: 10
                },
                backColor: "transparent",
                zIndex: -1,
                transform: 0,
                diy:[
                  {
                    name:"configComponent.bigScreen.border.border89cur",
                    type:1,
                    value:10,
                    min:1,
                    key:"border89cur",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.horizontal"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.vertical"
                      }
                    ],
                    min:1,
                    key:"border89Direction",
                  },
                ]
              }
            }
          },
          dur:1,
          reverse:false,
          secondColor:"#ffffff",
          mainColor:"#ffffff",
        }
    },
    methods: {
      setOption(option){
        for(let i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="border89cur")
          {
            this.dur=parseInt(option.style.diy[i].value)
          }
          if(option.style.diy[i].key=="border89Direction")
          {
            const value = parseInt(option.style.diy[i].value)
            if(value)
            {
              this.reverse=true
            }
            else
            {
              this.reverse=false
            }

          }
        }
      }
    },
    mounted() {
    this.$nextTick(function(){
      this.setOption(this.detail);
    });
  }
}
</script>

<style lang="less">
.view-text {
    height: 100%;
    width: 100%;
}
</style>

<template>

 <dv-border-box-8 :dur="dur" :reverse="reverse" :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-8>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box8',
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
            text: "configComponent.bigScreen.border.border8title",
            icon: "data:image/png;base64,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",
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
                  h: 300
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
                    value:1,
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
                    key:"border89Direction",
                  }
                ]
              }
            }
          },
          dur:1,
          reverse:false
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

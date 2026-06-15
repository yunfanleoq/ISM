<template>

 <dv-border-box-4 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-4>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box4',
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
    data() {
        return {
          detail:null,
          IsToolBox:false,
          editMode:true,
          base:{
            text: "configComponent.bigScreen.border.border4.title",
            icon: "data:image/png;base64,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",
            isFontIcon: true,
            info: {
              type: "text",
              action: [],
              dataBind:
                [
                ],
              style: {
                "visible":1,
                position: {
                  x: 0,
                  y: 0,
                  w: 300,
                  h: 300
                },
                backColor: "transparent",
                zIndex: -1,
                transform: 0,
              }
            }
          }
        }
    },
    methods: {

    }
}
</script>

<style lang="less">
.view-text {
    height: 100%;
    width: 100%;
}
</style>

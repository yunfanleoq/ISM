<template>

 <dv-border-box-6 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-6>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box6',
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
            text: "configComponent.bigScreen.border.border6title",
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

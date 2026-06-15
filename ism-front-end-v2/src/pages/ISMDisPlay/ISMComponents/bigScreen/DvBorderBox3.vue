<template>

 <dv-border-box-3 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-3>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box3',
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
            text: "configComponent.bigScreen.border.border3title",
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

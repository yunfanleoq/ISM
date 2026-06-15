<template>

 <dv-border-box-4 :reverse="true" :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-4>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box4-reverse',
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
            text: "configComponent.bigScreen.border.border4.reverseTitle",
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

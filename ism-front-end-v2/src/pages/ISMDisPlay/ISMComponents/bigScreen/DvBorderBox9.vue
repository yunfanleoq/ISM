<template>

 <dv-border-box-9 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-9>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box9',
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
            text: "configComponent.bigScreen.border.border9title",
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

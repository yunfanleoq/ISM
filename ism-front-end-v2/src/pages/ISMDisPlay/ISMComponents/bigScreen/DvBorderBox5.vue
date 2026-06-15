<template>

 <dv-border-box-5 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-5>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box5',
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
            text: "configComponent.bigScreen.border.border5.title",
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

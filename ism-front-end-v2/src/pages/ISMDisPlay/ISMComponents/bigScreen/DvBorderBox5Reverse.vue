<template>

 <dv-border-box-5 :reverse="true" :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-5>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box5-reverse',
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
            text: "configComponent.bigScreen.border.border5.reverseTitle",
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

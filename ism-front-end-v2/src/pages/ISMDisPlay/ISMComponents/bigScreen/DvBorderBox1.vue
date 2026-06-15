<template>

 <dv-border-box-1 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-1>
</template>

<script>

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box1',
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
            text: "configComponent.bigScreen.border.border1title",
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

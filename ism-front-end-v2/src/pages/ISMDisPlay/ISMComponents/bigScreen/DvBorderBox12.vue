<template>
  <dv-border-box-12 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

  </dv-border-box-12>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box12',
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
            text: "configComponent.bigScreen.border.border12title",
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
          },
          title:"title",
          titleWidth:30
        }
    },
    methods: {

    },
    mounted() {

  }
}
</script>

<style lang="less">
.view-text {
    height: 100%;
    width: 100%;
}
</style>

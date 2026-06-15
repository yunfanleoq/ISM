<template>

 <dv-border-box-13 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}">

 </dv-border-box-13>
</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-border-box13',
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
            text: "configComponent.bigScreen.border.border13title",
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

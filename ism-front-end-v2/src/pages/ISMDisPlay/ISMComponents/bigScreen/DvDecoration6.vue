<template>
<div>
  <dv-decoration-6 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}"/>
</div>


</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-decoration6',
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
            text: "configComponent.bigScreen.embellish.dvDecoration6",
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
                  w: 270,
                  h: 50
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

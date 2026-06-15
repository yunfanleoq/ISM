<template>
<div>
  <dv-decoration-8 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}"/>
</div>


</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-decoration8',
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
            text: "configComponent.bigScreen.embellish.dvDecoration7",
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

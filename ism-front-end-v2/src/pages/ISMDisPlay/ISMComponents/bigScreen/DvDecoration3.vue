<template>
<div>
  <dv-decoration-3 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}"/>
</div>


</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-decoration3',
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
            text: "configComponent.bigScreen.embellish.dvDecoration3",
            icon: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAATAAAABNCAYAAAAhMDhYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAL1SURBVHhe7dm/b01hHMdx/4ze5rYl1UuDpjTVNOnQIhZFMPmdEpZ20jAYDAyIxMhgZLBJjGLwH/gHxP9w3IfoeR7nDCei6pu8htfQ5tOn994nfQ+nu0bHBxVARAIGhCVgQFgCBoQlYEBYAgaEJWBAWAIGhCVgQFgCBoQlYEBYAgaEJWBAWAIGhCVgQFgCBoQlYEBYAgaEJWBAWAIGhCVgQFgCBoQlYEBYAgaEJWBAWAIGhCVgQFgCBoQlYEBYAgaEteMBmzt7rdr49K3h1N0nxW56+Uw1e+5Gw9jUzNamN7S7xchEfU6yfOdB6+9cvLxe7E5uPq1WH78unH70utgcWjnXetbqw5fFrjc2fB39qYbRsaFst/7xa+Os628+F5v0On/fJOl95btLrz607vZMz21tfnxmew40jAzlZ6X7aDsr3V++u/3+S2OTvpdvut75xedvW3eD+ZWtTXovbZv03vOzut55+qx/36Q7yTdd7zx93bZLP5/v/vWdp8+vbZM+7/yskaG2v6fe+P5it5MELCNgAiZg9VkC1oGACdgvAlYTsG4ELCNgAiZg9VkC1oGACdgvAlYTsG4ELCNgAiZg9VkCBrCNBAwIS8CAsAQMCCvMQ/z/9YHu5OxSdWLzWcPc+VvFbqS/r9rdn2zojW/fQ/xeehi/d7oh33R9oPs3H+L3BzPV5MKJhomDx4rd/3rn/clDw9d7vGHi8EKxW7p5v7rw4l3DYPhe89385Y1q/krpyPm1YvPzH1RTDelBe77r8hB/Yvpotbh2r2F29WpxVtc730kClhGwmoDVBKwkYBkBEzABE7A/JWAZAasJWE3ASgKWETABEzAB+1MClhGwmoDVBKwkYAB/gYABYQkYEJaAAWEJGBCWgAFhCRgQloABYQkYEJaAAWEJGBCWgAFhCRgQloABYQkYEJaAAWEJGBCWgAFhCRgQloABYQkYEJaAAWEJGBCWgAFhCRgQloABYQkYEJaAAWEJGBCWgAFhCRgQ1KD6Dg8sV1+x0pejAAAAAElFTkSuQmCC",
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

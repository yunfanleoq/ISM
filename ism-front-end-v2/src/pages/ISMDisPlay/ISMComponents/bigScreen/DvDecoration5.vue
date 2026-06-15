<template>
<div>
  <dv-decoration-5 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}"/>
</div>


</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-decoration5',
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
            text: "configComponent.bigScreen.embellish.dvDecoration5",
            icon: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAXIAAAB2CAYAAADVwhirAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAM9SURBVHhe7drNitMAFIZhL0dQBHEYKAiCG7fC7AX3A26F2bqZG9Ar8DorZ1EI4fQnbZ3kS5/FA10EJtCTd5KTvnr99m4LQC4hBwgn5ADhhBwgnJADhBNygHBCDhBOyAHCCTlAOCEHCCfkAOGEHCCckAOEE3KAcEIOEE7IAcIJOUA4IQcIJ+QA4YQcIJyQA4QTcoBwQg4QTsgBwgk5QDghBwgn5ADhhBwgnJADhBNygHBCDhBOyAHCCTlAOCEHCCfkAOGEHCCckAOEE3KAcEIOEE7IAcIJOUA4IQcIJ+QA4RYR8ndfvm4/fH9cnPcP39rzhbXrroclqFZ053vrFhHy+oI+/f67SHVu3TnDWrke8wj5ER+f/7gL4GbUrHfXwVIIec9qZY8K+G546vObzef23GEtasaH0dw8PbfXxpzcVPW87NyjhnoY8xrq7jhYi5rx3bzX7HfHsExCfsD4MbPuCLrjIF3N9nDW3flmEfIjhgNuX84auWHJJ+RH1Ipl/MhpX85adHvx7jiWTchPMN6X1+e6azmH36Zzbd2cncpefB2E/ETjx89L1AXU/Q2YahjiS1kb5hLyCSrA3QUwlV0713DNiLu5yCbkE9SKpVYju8fSqcbrGbt2zlFzM474/Y+nduZO1f0dcgj5Cxrv2r1YYqqaoWHAi/cuCPkL81MvztW9p7Giowj5DCreuwvRvpxTiDiHCPkMxjtO+3IOqdXJMOD++TMm5DMZ78vr8/DlE5R6iTmOeDdP3DYhn1H3uAz7eDnOPkI+s7rr6i5aGBJxDhHymdWK5ZLfpnMbutmBHSEHCCfkAOGEHCCckAOEE3KAcEK+Apufv9qfrLE89V113yFcQshXQMhzCDn/g5ADhBNygHBCDhBOyAHCCTlAOCEHCCfkAOGEHCCckAOEE3KAcEIOEE7IAcIJOUA4IQcIJ+QA4YQcIJyQA4QTcoBwQg4QTsgBwgk5QDghBwgn5ADhhBwgnJADhBNygHBCDhBOyAHCCTlAOCEHCCfkAOGEHCCckAOEE3KAcEIOEE7IAcIJOUA4IQcIJ+QA4YQcIJyQA4QTcoBwQg4QTsgBwgk5QDghBwgn5ADhhBwg2t32H/r2dD0gIlKTAAAAAElFTkSuQmCC",
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

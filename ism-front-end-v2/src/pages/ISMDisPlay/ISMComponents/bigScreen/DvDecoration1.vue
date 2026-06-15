<template>
<div>
  <dv-decoration-1 :style="{width:detail.style.position.w+'px',height: detail.style.position.h+'px'}"/>
</div>


</template>

<script>

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'dv-decoration1',
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
            text: "configComponent.bigScreen.embellish.dvDecoration1",
            icon: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQwAAAAoCAYAAAALxsQHAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAALYSURBVHhe7d1PaxNBHMZxQU2bNNXWVkmVJCb9A4KCRIIeChXTVgWPavUgFUGh4KlosXjwBSievPguqvgKfA/pW/DoNbeRWZplZ/mpM0Oz2T/fhQ8U+pTtZjYPncywPXHyVFkBgA0Kw1Nl6qxBymA87vTuGqRMVkxVZwwTk9NiLim5LIzTpUpI+v5xmD13QdUbrUBtoS5mMB4HB9/V8NBfS5ms0CXRXlwJURgjQGEUG4UxOkxJPDElSS+mJKNDYQCwRmEAsEZhALBGYQCwRmEAsEZhILcm377Kj5dPxGtMGoWRQ67Liq32ivPPlCtnDKWJqpgbh+E1NL99VdO/fuaCS2FEx1GTMnGlBz2DlNEojBwaDAZH25aU2tl5LWaidGFEj/WN+2IuSpfExUuNUJoKY7hx69nvQ/HNl0UuhaHHfHjoe0HKxJW/fAjPVe3/EDMahZFDFAaFMTwoDPyX65+kTEnSjykJMGLih4dZxYeeALKGwgBgjcIAYI3CAGCNwgBgzbswms0lg5SJc13u0Q8MmZurGaRcVPwceslQymWFfrJXlJTJivjYSJnjkNR5XPksX8/PLxj0e0LKJcW7MLafv1C7u28Ca7f/vm4b5bqhSL84vfXNkH7BpFyU3nQUPbJeGJ1OV62urgWyfi0+G4p8pPURfXr8oofNBjl9z29s3gtRGP9AYVAYPiiM0WFKknJMSdwldR5X+l50/d10YURltjAAFA+FAcAahQHAGoUBwBqFAcAahQHAmndhzMyeN0iZONclJZ9lqPi/MEzTg12SUqvVDVJmHFzHEn4ut5YNUibOdmy8C+Pqteuq270VaDQXxUyc68Ytn40uuiQazXaoiIXx8NGWerz1NLC0fEXMjENSG7eKrnPjptrffx/Y23snZuI+fvp8NDJK9fuHYkajMHKIwii2VBYGU5L0YkpSbKmckgAoHgoDgDUKA4A1CgOANQoDgKWy+gPfSq8UGGcZxAAAAABJRU5ErkJggg==",
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

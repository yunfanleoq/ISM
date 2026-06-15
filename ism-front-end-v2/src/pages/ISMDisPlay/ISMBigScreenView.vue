<template>
<div @mousewheel="handleScroll" style="overflow: auto; width: 100%; height: 100%; position: relative;">
    <div style="position: absolute; top: 0px; left: 0px; width:100%; height: 100%; background-color: rgb(8, 0, 0);" class="ism-render" :style="layerStyle" v-if="configData.layer">
      <template v-for="(component,index) in configData.components">
        <div class="ism-render-wrapper"
             :key="index"
             @click="doClickComponent(component)"
             @dblclick="doDbClickComponent(component)"
             :class="{'ism-render-wrapper-clickable': component.action.length > 0 }"
             v-show="component.style.visible == undefined? true:component.style.visible"
             :style="{
                        left: component.style.position.x + 'px',
                        top: component.style.position.y + 'px',
                        width: component.style.position.w + 'px',
                        height: component.style.position.h + 'px',
                        backgroundColor: component.style.backColor,
                        zIndex: component.style.zIndex,
                        borderWidth: component.style.borderWidth + 'px',
                        borderStyle: component.style.borderStyle,
                        borderColor: component.style.borderColor,
                        transform: component.style.transform? `rotate(${component.style.transform}deg)`:'rotate(0deg)',
                    }">
          <component v-bind:is="component.type" :editMode="false" :showDeviceUuid="showDeviceUuid" :detail="component" ref="spirit" />
        </div>
      </template>
    </div>
  </div>

</template>

<script>
import ISMBase from './ISMBase';
import store from "@/store";
import { mapActions, mapGetters, mapState, mapMutations } from 'vuex'

export default {
    name: 'ISMRender',
    extends: ISMBase,
    components: {

    },
    props: {
      showUuid: {
        type: String,
        required: true
      },
      showDeviceUuid: {
        type: String,
        required: true
      }
    },
    watch: {
      showUuid: {
        handler(newVal, oldVal) {
          if(newVal!="")
          {
            this.getLayerDataStruct(newVal);
            this.showUuid = newVal
          }
        },
        deep: true
      }
    },
    computed: {
      ...mapState({
        configData: state => store.state.ISMDisPlayEditorTool.LayerData,
        selectedValue:state => store.state.ISMDisPlayEditorTool.selectedValue,
      }),
        layerStyle:function () {
            let scale = this.selectedValue / 100;
            let styles = [`transform:scale(${scale})`];
            if(this.configData.layer.backColor) {
                styles.push(`background-color: ${this.configData.layer.backColor}`);
            }
            if(this.configData.layer.backgroundImage) {
                styles.push(`background-image: url("${this.configData.layer.backgroundImage}")`);
            }
            if(this.configData.layer.width > 0) {
                // styles.push(`width: ${this.configData.layer.width}px`);
            }
            if(this.configData.layer.height > 0) {
                // styles.push(`height: ${this.configData.layer.height}px`);
            }
            var style = styles.join(';');
            return style;
        }
    },
    data() {
        return {
          selectedValueTemp:100
        }
    },
    methods: {
      ...mapMutations('ISMDisPlayEditorTool',[
        'setlayerZoom',
      ]),
      ...mapActions('ISMDisPlayEditorTool',[
        'getLayerDataStruct'
      ]),
      getScale () {
        const { width, height } = this
        const wh = window.innerHeight / this.configData.layer.height
        const ww = window.innerWidth / this.configData.layer.width
        return ww < wh ? ww : wh
      },
      zoomOut(){
        this.selectedValueTemp = this.selectedValueTemp-5
        if(this.selectedValueTemp <5)
        {
          this.selectedValueTemp=5
        }
        this.setlayerZoom(this.selectedValueTemp)
      },
      zoomIn(){
        this.selectedValueTemp=this.selectedValueTemp+5
        if(this.selectedValueTemp >200)
        {
          this.selectedValueTemp=200
        }
        this.setlayerZoom(this.selectedValueTemp)
      },
      handleScroll: function (e) {
        if (e.wheelDelta || e.detail) {
          if (e.wheelDelta > 0 || e.detail < 0) {     //当鼠标滚轮向上滚动时
            this.zoomIn()
          }
          if (e.wheelDelta < 0 || e.detail > 0) {     //当鼠标滚轮向下滚动时

            this.zoomOut()
          }
        }
      },
      doClickComponent(component){
          for(let i = 0; i < component.action.length; i++) {
              let action = component.action[i];
              if(action.type == 'click') {
                   this.handleComponentActuib(action);
              }
          }
      },
      doDbClickComponent(component){
          for(let i = 0; i < component.action.length; i++) {
              let action = component.action[i];
              if(action.type == 'dblclick') {
                   this.handleComponentActuib(action);
              }
          }
      },
      handleComponentActuib(action){
          let _this = this;
          if(action.action == 'visible'){
              if(action.showItems.length > 0) {
                  action.showItems.forEach(identifier => {
                      _this.showComponent(identifier,true);
                  });
              }
              if(action.hideItems.length > 0) {
                  action.hideItems.forEach(identifier => {
                      _this.showComponent(identifier,false);
                  });
              }
          } else if(action.action == 'service') {
              _this.sendFun(action);
          }
      },
      showComponent(identifier,visible) {
            let spirits = this.$refs['spirit'];
            for(var i = 0; i < spirits.length; i++){
                var spirit = spirits[i];
                if(spirit.detail.identifier == identifier) {
                    spirit.detail.style.visible = visible;
                    break;
                }
            }
        },
    },
    mounted() {
      if(this.showUuid!="")
      {
        this.getLayerDataStruct(this.showUuid);
      }
      // window.addEventListener('scroll', this.handleScroll, true);  // 监听（绑定）滚轮滚动事件
      this.selectedValueTemp = this.selectedValue
    },
    destroyed: function () {
      // window.removeEventListener('scroll', this.handleScroll);   //  离开页面清除（移除）滚轮滚动事件
    }
}
</script>

<style lang="less">
    .ism-render {
        overflow: auto;
        z-index: 99;
        //background-color: white;
        //background-clip: padding-box;
        //background-origin: padding-box;
        //background-repeat: no-repeat;
        //background-size: 100% 100%;
        //position: relative;
        //height: 100%;

        .ism-render-wrapper {
            position: absolute;
        }

        .ism-render-wrapper-clickable {
            cursor: pointer;
        }
    }
</style>


<template>
  <div style="border-right:0px solid #95B8E7; overflow: hidden;
            overflow-y: auto;">
    <a-modal id="showImg" :visible="showImg" :footer="null" @cancel="viewImg">
      <img alt="example" style="width: 100%" :src="imgSrc" />
    </a-modal>

    <a-modal id="iconIsShow" :visible="iconIsShow" :footer="null" @cancel="iconIsShow=false" width="400px">
      <icon-font :type="iconIsShowSrc" :style="{ fontSize: '256px' }"/>
    </a-modal>

    <a-collapse accordion :bordered="false" default-active-key="1"  expand-icon-position="left">
      <template v-for="(group,index1) in mestoolbox">
        <a-collapse-panel  :header="$t(group.title)" :key="index1" :style="customStyle">
          <div class="ism-toolbox " style="border-right:0px solid #95B8E7;height: 300px; overflow: hidden;
            overflow-y: auto;">
            <div class="toolbox-group scrollbar">
              <template v-for="(value,index) in validResultsList(group.items)">
                <div
                    class="toolbox-item"
                    v-bind:key="index"
                    draggable="true"
                    @dragstart="onDragstart($event,value)"
                >
                  <template >
                    <div class="toolbox-item-icon">
                      <div v-if="value.icon.indexOf('icon-')!=-1">
                        <icon-font :type="value.icon" :style="{ fontSize: '32px' }" @click="clickIconImg(value.icon)"/>
                      </div>
                      <div v-else-if="value.icon.indexOf('svg-')!=-1">
                        <component :IsToolBox=true v-bind:is="value.info.type" :detail="value.info" showDeviceUuid=""  :ref="'comp' + index"/>
                      </div>
                      <div v-else>
                        <img alt="icon" :src="value.icon"  @click="clickImg($event)" style="width:32px;height: 32px"/>
                      </div>

                    </div>
                    <div class="toolbox-item-text" v-if="group.isSequence"> {{$t(group.title)}}{{index+1}}</div>
                    <div class="toolbox-item-text" v-else> {{$t(value.text)}}</div>
                  </template>
                </div>
              </template>
            </div>
          </div>
          <a-icon slot="extra" v-if="group.items.length>100" type="double-right" @click="ShowMoreItem(group)" />
        </a-collapse-panel>
      </template>
    </a-collapse>

    <a-modal :visible="itemMore"
            :title="$t(ShowMoreItemArray.title)"
            :dialogStyle="{width:'300px',height:'500px','z-index':-1,left: '352px', top:'70px', position:'absolute'}"
            :iconCls="ShowMoreItemArray.icon"
             :footer="null"
             @cancel="itemMore=false"
             v-drag-modal
             :destroyOnClose="true"
             :maskClosable="false"
             :maskStyle="{}"
            :mask="false">
      <div class="ism-toolbox " style="border-right:0px solid #95B8E7;" >
        <div class="toolbox-group scrollbar" v-if="(Array.isArray(ShowMoreItemArray.items))&&(ShowMoreItemArray.items.length>0)">
          <template v-for="(value,index) in MoreItemList(ShowMoreItemArray.items)" >
            <div
                class="toolbox-item"
                v-bind:key="index"
                draggable="true"
                @dragstart="onDragstart($event,value)"
            >
              <template >
                <div class="toolbox-item-icon">
                  <div v-if="value.icon.indexOf('icon-')!=-1">
                    <icon-font :type="value.icon" :style="{ fontSize: '32px' }"/>
                  </div>
                  <div v-else-if="value.icon.indexOf('svg-')!=-1">
                    <component :IsToolBox=true v-bind:is="value.info.type" :detail="value.info" showDeviceUuid=""  :ref="'comp' + index"/>
                  </div>
                  <div v-else>
                    <img alt="icon" :src="value.icon" @click="clickImg($event)" style="width:32px;height: 32px"/>
                  </div>

                </div>
                <div class="toolbox-item-text" v-if="ShowMoreItemArray.isSequence"> {{$t(ShowMoreItemArray.title)}}{{index+21}}</div>
                <div class="toolbox-item-text" v-else> {{$t(value.text)}}</div>
              </template>
            </div>
          </template>
        </div>
      </div>
    </a-modal>

  </div>
</template>

<script>
import ISMBase from './ISMBase';
import {mapState} from "vuex";
import store from "../../store";

export default {
  name: "MesToolBox",
  i18n: require('../../i18n/language'),
  extends: ISMBase,
  data() {
    return {
      iconIsShowSrc:"",
      iconIsShow:false,
      showImg:false,
      itemMore:false,
      ShowMoreItemArray:[],
      customStyle: 'background: #fff;border-radius: 4px;margin-bottom: 1px;border: 0;overflow: hidden',
      ShowMoreModel:false,
      imgSrc: '',
      selectedIndex:[0],
    }
  },
  components: {

  },
  props: [],
  computed: {
    ...mapState({
      mestoolbox: state => store.state.ISMDisPlayEditorTool.MesComponentsList,
    }),
    validResultsList(){
      return item => {
        return item.slice(0,100)
      }
    },
    MoreItemList(){
      return item => {
        return item.slice(100,item.length)
      }
    }
  },
  created(){
    this.$nextTick(function () {
      this.itemMore = false
    });
  },
  methods: {
    SelectAccordion(ev){
      this.itemMore = false
      this.ShowMoreModel = false
    },
    ShowMoreItem(items){
      let _t = this
      this.ShowMoreItemArray = items
      this.ShowMoreModel = true
      this.itemMore = true
    },
    clickImg(e) {
      this.showImg = true;
      // 获取当前图片地址
      this.imgSrc = e.currentTarget.src;
    },
    clickIconImg(e) {
      this.iconIsShow = true;
      console.log(e)
      // 获取当前图片地址
      this.iconIsShowSrc = e;
    },
    viewImg(){
      this.showImg = false;
    },
    onDragstart (event, info) {
      this.$emit('drag-start', info.info,event)
    },
  },
  mounted () {

  }
}
</script>

<style lang="less" scoped>
::v-deep .ant-modal-mask {
  background: transparent;
  pointer-events: none;
}
::v-deep .ant-modal-wrap {
  pointer-events: none;
}
::v-deep .ant-modal {
  pointer-events: all;
}
.window-mask{
  z-index:0
}
::v-deep .ant-collapse-borderless > .ant-collapse-item > .ant-collapse-content > .ant-collapse-content-box {
  padding-top: 0px;
}

::v-deep .ant-collapse-content > .ant-collapse-content-box {
  padding: 0px;
}
#showImg .ant-modal-content {
  background-color: transparent !important;
}
.more-right {
  position: absolute;
  right: 24px;
  height: 30px;
  top: 5px;
}
.ism-toolbox {
  background-color: white;
  overflow-y: auto;
  overflow-x: hidden;
  .toolbox-group {
    display: flex;
    flex-wrap: wrap;
    padding: 5px;
    //margin-top: 5px;
    //margin-bottom: 5px;
    justify-content: flex-start;
    align-content: space-between;
    .toolbox-item {
      width: 70px;
      //margin: 1px 5px;
      padding: 6px;
      color: #777;
      border: transparent solid 1px;
      &.base {
        width: 64px;
      }
      .toolbox-item-icon {
        text-align: center;
      }

      .toolbox-item-text {
        margin-top: 2px;
        overflow: hidden;
        white-space: nowrap;
        /* ⽂字超出宽度则显⽰ellipsis省略号 */
        text-overflow:ellipsis;
        text-align: center;
      }
    }

    .toolbox-item:hover {
      border: #ccc solid 1px;
      background: #ccc;
      color: #3388ff;
      border-radius: 6px;
      cursor: pointer;
    }
  }
}
.item-icon {
  height: 30px;
  .iconfont {
    font-size: 30px;
    margin-left: 14px;
  }
}
::-webkit-scrollbar {
  /*滚动条整体样式*/
  width : 5px;  /*高宽分别对应横竖滚动条的尺寸*/
  height: 9px;
}
::-webkit-scrollbar-thumb {
  /*滚动条里面小方块*/
  border-radius   : 10px;
  background-color: skyblue;
  background-image: -webkit-linear-gradient(
      45deg,
      rgba(255, 255, 255, 0.2) 25%,
      transparent 25%,
      transparent 50%,
      rgba(255, 255, 255, 0.2) 50%,
      rgba(255, 255, 255, 0.2) 75%,
      transparent 75%,
      transparent
  );
}
::-webkit-scrollbar-track {
  /*滚动条里面轨道*/
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}
</style>

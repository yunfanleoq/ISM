<template>
  <div id="toolbox">

    <input type="file" ref="filElem" name="file" @change="showFile($event)" style="display: none"/>
    <a-card  :bodyStyle="{padding:'6px',border:0}" style="height: 50px;">
        <div  style="display: inline">
          <div style="float: left">
          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.ToolBar.Import') }}
            </template>
            <a-button type="link"  @click="choiceModelJson($event)"  >
              <icon-font type="icon-daoru"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>
          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.ToolBar.Export') }}
            </template>
            <a-button type="link"  @click="exportJson">
            <icon-font type="icon-peizhidaochu"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>
<!--          <a-tooltip>-->
<!--            <template slot="title">-->
<!--              {{ $t('displayConfig.ToolBar.ExportSvg') }}-->
<!--            </template>-->
<!--            <a-button type="link"  @click="exportSvg">-->
<!--              <icon-font type="icon-daochusvg"  style="font-size: 22px"/>-->
<!--            </a-button>-->
<!--          </a-tooltip>-->
          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.ToolBar.Save') }}
            </template>
            <a-button type="link"  @click="doSaveLayerData($route.params.uid)"  >
              <icon-font type="icon-baocun"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.PicTemplete') }}
              </template>
              <a-button type="link"  @click="showSystemImageModel(0)">
                <icon-font type="icon-tupian"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Template2D.title') }}
              </template>
              <a-button type="link" @click="show2DTemplatePicker = true">
                <a-icon type="appstore" style="font-size: 20px"/>
              </a-button>
            </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.View') }}
              </template>
              <a-button type="link"  @click="pageView"  >
                <icon-font type="icon-yunhang"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                AI 2D Scene
              </template>
              <a-button type="link" @click="openAiSceneModal">
                <a-icon type="thunderbolt" style="font-size: 20px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.Delete') }}
              </template>
              <a-button type="link" @click="DelNode"  >
                <icon-font type="icon-shanchu"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.Format') }}
              </template>
              <a-button type="link" @click="FormatDo"  >
                <icon-font type="icon-geshishua"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.ToolBar.Undo') }}
            </template>
            <a-button type="link" @click="doWithUndo"  >
              <icon-font type="icon-jurassic_laststep"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>

          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.ToolBar.Redo') }}
            </template>
            <a-button type="link" @click="doWithRedo"  >
              <icon-font type="icon-jurassic_nextstep"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>

          <a-divider type="vertical" style="height:1.5em;top:-0.2em;background-color: #95B8E7;"/>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.ZoomOut') }}
              </template>
              <a-button type="link" @click="zoomOut"  >
                <icon-font type="icon-suoxiao"  style="font-size: 20px"/>
              </a-button>
            </a-tooltip>

          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.ToolBar.Zoom') }}
            </template>
            <a-select v-model="selectedValueTemp" size="small" default-value="a1" style="width: 200px;line-height:300px" @change="onZoomChange($event)"  :style="{height:'22px',width:'80px'}">
              <a-select-option v-for="(item ,index) in zoomList" :key="index" :value="item.value">
                {{item.text}}
              </a-select-option>
            </a-select>
          </a-tooltip>

          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.ToolBar.ZoomIn') }}
            </template>
            <a-button type="link" @click="zoomIn"  >
              <icon-font type="icon-fangda"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>
          <a-divider type="vertical" style="height:1.5em;top:-0.2em;background-color: #95B8E7;"/>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.MoveUpOneIndex') }}
              </template>
              <a-button type="link" @click="SetIndexUp"  >
                <icon-font type="icon-shangyiyiceng1"  style="font-size: 20px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.MoveDownOneIndex') }}
              </template>
              <a-button type="link" @click="SetIndexDown"  >
                <icon-font type="icon-xiayiyiceng2"  style="font-size: 20px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.MoveFront') }}
              </template>
              <a-button type="link" @click="SetIndexFront"  >
                <icon-font type="icon-dingceng"  style="font-size: 20px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.ToolBar.MoveBack') }}
              </template>
              <a-button type="link" @click="SetIndexBack"  >
                <icon-font type="icon-zhiyudiceng1"  style="font-size: 20px"/>
              </a-button>
            </a-tooltip>
            <a-divider type="vertical" style="height:1.5em;top:-0.2em;background-color: #95B8E7;"/>
          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.Canvas.menu.Revolve') }}
            </template>
            <a-button type="link" @click="revolve"  >
              <icon-font type="icon-shunshizhenxuanzhuan"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>
          <a-tooltip>
            <template slot="title">
              {{ $t('displayConfig.Canvas.menu.Reverse') }}
            </template>
            <a-button type="link" @click="reverse"  >
              <icon-font type="icon-nishizhenxuanzhuan"  style="font-size: 22px"/>
            </a-button>
          </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.FlipVertical') }}
              </template>
              <a-button type="link" @click="FlipVertical"  >
                <icon-font type="icon-shuipingfanzhuan"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.FlipHorizontally') }}
              </template>
              <a-button type="link" @click="FlipHorizontally"  >
                <icon-font type="icon-chuizhifanzhuan"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
          <a-divider type="vertical" style="height:1.5em;top:-0.2em;background-color: #95B8E7;"/>
<!--            <a-tooltip>-->
<!--              <template slot="title">-->
<!--                {{ $t('displayConfig.Canvas.menu.VerticalEquidistant') }}-->
<!--              </template>-->
<!--              <a-input size="small"  v-model="Isometric_colu" style="width: 50px;height: 22px"/>-->
<!--            </a-tooltip>-->
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.VerticalEquidistant') }}
              </template>
              <a-button type="link" @click="HeaderSetCommentsAlign('Vertical')"  >
                <icon-font type="icon-chuizhidengjufenbu"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
<!--            <a-tooltip>-->
<!--              <template slot="title">-->
<!--                {{ $t('displayConfig.Canvas.menu.HorizontalEquidistant') }}-->
<!--              </template>-->
<!--              <a-input size="small"  v-model="Isometric_row" style="width: 50px;height: 22px"/>-->
<!--            </a-tooltip>-->
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.HorizontalEquidistant') }}
              </template>
              <a-button type="link" @click="HeaderSetCommentsAlign('Horizontal')"  >
                <icon-font type="icon-shuipingdengjufenbu"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.topAlign') }}
              </template>
              <a-button type="link" @click="HeaderSetCommentsAlign('t')"  >
                <icon-font type="icon-jurassic_verticalalign-top"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.bottomAlign') }}
              </template>
              <a-button type="link" @click="HeaderSetCommentsAlign('b')"  >
                <icon-font type="icon-jurassic_verticalalign-bottom"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.rightAlign') }}
              </template>

              <a-button type="link" @click="HeaderSetCommentsAlign('r')"  >
                <icon-font type="icon-jurassic_horizalign-right"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>

            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.leftAlign') }}
              </template>
              <a-button type="link" @click="HeaderSetCommentsAlign('l')"  >
                <icon-font type="icon-jurassic_horizalign-left"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
          <a-divider type="vertical" style="height:1.5em;top:-0.2em;background-color: #95B8E7;"/>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.AllSelect') }}
              </template>

              <a-button type="link" @click="SelectAllComments()"  >
                <icon-font type="icon-total_selection"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.Lock') }}
              </template>

              <a-button type="link" @click="HeaderLockItem(true)"  >
                <icon-font type="icon-suoding"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.UnLockAll') }}
              </template>
              <a-button type="link" @click="HeaderUnlockItem(false)"  >
                <icon-font type="icon-jiesuo"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>

            <a-divider type="vertical" style="height:1.5em;top:-0.2em;background-color: #95B8E7;"/>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.GroupTogether') }}
              </template>
              <a-button type="link" @click="HeaderCreateGroup('group')"  >
                <icon-font type="icon-zuhe"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>
            <a-tooltip>
              <template slot="title">
                {{ $t('displayConfig.Canvas.menu.UnGroup') }}
              </template>
              <a-button type="link" @click="HeaderSplitGroup('ungroup')"  >
                <icon-font type="icon-jiezu"  style="font-size: 22px"/>
              </a-button>
            </a-tooltip>

          </div>
<!--          <div class="charge-tips"  @click="doSaveLayerData($route.params.uid)" style="" v-if="isCharge"><a-icon type="exclamation-circle" style="margin-right: 5px"/>{{$t('displayConfig.ToolBar.isChargeTips')}}</div>-->
        </div>
    </a-card>

    <a-modal v-drag-modal :destroyOnClose="true"  :footer="null" v-model="visible">
      <json-viewer :value="configData" :expand-depth="4"  sort ></json-viewer>
    </a-modal>


    <a-modal v-model="aboutVisible" :footer="null"  :title="$t('displayConfig.ToolBar.HelpAbout')">
        <div class="help-content">
          <p > {{$t('displayConfig.ToolBar.HelpContent')}}</p>
        </div>
    </a-modal>

    <a-modal
      v-model="aiSceneVisible"
      title="AI 2D Scene"
      width="640px"
      :confirmLoading="aiSceneLoading"
      okText="Generate"
      cancelText="Cancel"
      @ok="handleAiSceneGenerate"
    >
      <a-form :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }">
        <a-form-item label="Prompt">
          <a-textarea
            v-model="aiScenePrompt"
            :rows="5"
            placeholder="Example: create a power room monitoring screen with transformer, switch cabinets, real-time data, alarms, and trend chart"
          />
        </a-form-item>
        <a-form-item label="Mode">
          <a-radio-group v-model="aiSceneMode">
            <a-radio-button value="replace">Replace</a-radio-button>
            <a-radio-button value="append">Append</a-radio-button>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="Canvas">
          <a-input-group compact>
            <a-input-number v-model="aiSceneWidth" :min="300" :max="7680" style="width: 120px" />
            <a-input-number v-model="aiSceneHeight" :min="300" :max="4320" style="width: 120px" />
          </a-input-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <ISM2DTemplatePicker
      :visible="show2DTemplatePicker"
      @close="show2DTemplatePicker = false"
      @select="handle2DTemplateSelect"
    />

    <system-image-model @onSelectImage="onSelectImage" :networkImageUrl="imageComponentData" ref="HeaderSystemImageModel"></system-image-model>
  </div>
</template>

<script>
import store from "@/store";

import JsonViewer from 'vue-json-viewer'
import {
  deepCopy
} from "@/assets/libs/utils";
import { uuid } from 'vue-uuid';
import { DataUri } from "@antv/x6";
import systemImageModel from '@/components/systemImageModel/systemImageModel'
import { mapActions, mapGetters, mapState, mapMutations } from 'vuex'
import {GetSystemPageTempleteContent} from "@/services/system";
import {ClearOpFlag, SetFormatPainterCell} from "@/store/ISM/mutations";
import {updateAllLayerDataStruct, updateLayerDataStruct} from "@/store/ISM/actions";
import {sm4EncryptBase64,sm4DecryptBase64} from "@/utils/smUtils";
import {GenerateISM2DScene} from "@/services/aiScene";
import {buildLocalAiSceneGraph} from "@/pages/ISMDisPlay/utils/aiSceneLocalGenerator";
import {normalizeISMScene} from "@/pages/ISMDisPlay/utils/ismSceneNormalizer";
import ISM2DTemplatePicker from './ISM2DTemplatePicker'
export default {
  name: "toolBox",
  i18n: require('../../i18n/language'),
  data() {
    return {
      isFormatPainterActiveTemp: false,
      imageComponentData:"",
      visible:false,
      aboutVisible:false,
      isCharge:false,
      isAll:1,
      isCreateTimer:null,
      _prevComponents: null,
      zoomList: [
        {
          value: 25,
          text: "25%"
        },
        {
          value: 50,
          text: "50%"
        },
        {
          value: 75,
          text: "75%"
        },
        {
          value: 100,
          text: "100%"
        },
        {
          value: 125,
          text: "125%"
        },
        {
          value: 150,
          text: "150%"
        },
        {
          value: 175,
          text: "175%"
        },
        {
          value: 200,
          text: "200%"
        }
      ],
      beforeUnloadTime: '',
      gapTime: '',
      menuAlign: "center",
      selectedValueTemp:100,
      aiSceneVisible:false,
      aiSceneLoading:false,
      aiScenePrompt:"",
      aiSceneMode:"replace",
      aiSceneWidth:1920,
      aiSceneHeight:1080,
      show2DTemplatePicker:false
    }
  },
  components: {
    JsonViewer,
    systemImageModel,
    ISM2DTemplatePicker,
  },
  watch:{
    configData: {
      handler(newVal, oldVal) {
        this.isCharge = true
        let _t = this
        if(!this.isHistoryOp){
          if(this._prevComponents && this._prevComponents.length > 0) {
            this.executeUndoStack({components: this._prevComponents})
          }
        }else{
          if(_t.isCreateTimer) {
            clearTimeout(_t.isCreateTimer)
          }

          _t.isCreateTimer = setTimeout(function () {
            _t.ClearOpFlag()
          }, 200)
        }
        // 保存当前 components 快照供下次 undo 使用（deep watch 下 oldVal === newVal，必须手动存）
        if(newVal && newVal.components) {
          this._prevComponents = deepCopy(newVal.components)
        }
      },
      deep: true
    }
  },
  props: [],
  methods: {
    ...mapMutations('ISMDisPlayEditorTool',[
      'setlayerZoom',
      'increaseCopyCount',
      'execute',
      'ClearSelectedComponent',
      'setCopySrcItems',
      'setFormatSrcItems',
      'undo',
      'clearSelectComponent',
      'ClearOpFlag',
      'executeUndoStack',
      'redo',
      'setLayerSelected',
      'SetFormatPainterState',
        'SetFormatPainterCell'
    ]),
    ...mapActions('ISMDisPlayEditorTool',[
      'saveLayerDataStruct',
      'setGroupList',
      'setLayerData',
      'SyncLayerData',
        'getLayerDataStruct',
        'updateAllLayerDataStruct'
    ]),
    openAiSceneModal(){
      this.aiSceneWidth = parseInt(this.configData.layer.width || 1920)
      this.aiSceneHeight = parseInt(this.configData.layer.height || 1080)
      if(!this.aiScenePrompt) {
        this.aiScenePrompt = "Create an industrial monitoring screen with title, equipment status, real-time data, alarms, and trend chart"
      }
      this.aiSceneVisible = true
    },
    handle2DTemplateSelect(tpl) {
      if(!this.ISMCavasContainer || !tpl || typeof tpl.build !== 'function') {
        this.$message.error(this.$t('displayConfig.Template2D.canvasNotReady'))
        return
      }
      const scene = normalizeISMScene(tpl.build(), {
        width: parseInt(this.configData.layer.width || 1920),
        height: parseInt(this.configData.layer.height || 1080),
        backColor: this.configData.layer.backColor || '#061321',
        allowedShapes: Array.from(new Set([
          ...this.availableAiShapes,
          'view-svg-text',
          'view-svg-time',
          'view-progress-bars',
          'dv-border-box1',
          'dv-border-box8',
          'ViewCanvasMoveLineArrow',
          'view-svg-electric1',
          'view-svg-electric2',
          'view-svg-electric3',
          'view-svg-electric4',
          'view-svg-electric5',
          'ism-view-device-status',
          'ism-view-real-data-smooth-chart-by-device',
          'ism-view-history-trend-chart',
          'ism-view-history-difference',
          'ism-view-real-table',
          'view-device-alarm-list',
          'view-device-real-data-table'
        ]))
      })
      if(!scene || !scene.components || !scene.components.cells.length) {
        this.$message.error(this.$t('displayConfig.Template2D.invalidTemplate'))
        return
      }
      this.configData.layer.width = scene.layer.width
      this.configData.layer.height = scene.layer.height
      this.configData.layer.backColor = scene.layer.backColor
      this.configData.layer.backgroundImage = scene.layer.backgroundImage
      this.ISMCavasContainer.resize(scene.layer.width, scene.layer.height)
      this.ISMCavasContainer.drawBackground({
        color: scene.layer.backColor,
        image: scene.layer.backgroundImage,
        size: '100% 100%',
        repeat: 'no-repeat',
        position:"center",
        quality:1,
      })
      const components = JSON.parse(JSON.stringify(scene.components))
      if (components.cells && Array.isArray(components.cells)) {
        components.cells = components.cells.filter(cell => cell && cell.shape)
      }
      this.ISMCavasContainer.fromJSON(components)
      this.setGroupList()
      this.isCharge = true
      this.show2DTemplatePicker = false
      this.$message.success(this.$t('displayConfig.Template2D.applySuccess'))
    },
    async handleAiSceneGenerate(){
      if(!this.ISMCavasContainer) {
        this.$message.error("Canvas is not ready")
        return
      }
      if(!this.aiScenePrompt || this.aiScenePrompt.trim().length === 0) {
        this.$message.warning("Please enter a prompt")
        return
      }
      this.aiSceneLoading = true
      try {
        const payload = {
          prompt: this.aiScenePrompt,
          canvas: {
            width: parseInt(this.aiSceneWidth || this.configData.layer.width || 1920),
            height: parseInt(this.aiSceneHeight || this.configData.layer.height || 1080)
          },
          mode: this.aiSceneMode,
          currentLayer: this.configData.layer,
          availableShapes: this.availableAiShapes
        }
        const scene = this.normalizeAiScene(await this.getAiSceneOrFallback(payload))
        if(!scene) {
          this.$message.error("AI scene data is invalid")
          return
        }
        if(this.aiSceneMode === "replace") {
          this.configData.layer.width = scene.layer.width
          this.configData.layer.height = scene.layer.height
          this.configData.layer.backColor = scene.layer.backColor
          this.configData.layer.backgroundImage = scene.layer.backgroundImage
          this.ISMCavasContainer.resize(scene.layer.width, scene.layer.height)
          this.ISMCavasContainer.drawBackground({
            color: scene.layer.backColor,
            image: scene.layer.backgroundImage,
            size: '100% 100%',
            repeat: 'no-repeat',
            position:"center",
            quality:1,
          })
          const components = JSON.parse(JSON.stringify(scene.components))
          if (components.cells && Array.isArray(components.cells)) {
            components.cells = components.cells.filter(cell => cell && cell.shape)
          }
          this.ISMCavasContainer.fromJSON(components)
        } else {
          scene.components.cells.forEach(cell => {
            this.ISMCavasContainer.addNode(cell)
          })
        }
        this.setGroupList()
        this.isCharge = true
        this.aiSceneVisible = false
        this.$message.success(scene.fromFallback ? "Generated with local fallback template" : "AI scene generated")
      } catch (e) {
        console.error(e)
        this.$message.error("AI scene generation failed")
      } finally {
        this.aiSceneLoading = false
      }
    },
    async getAiSceneOrFallback(payload) {
      try {
        const res = await GenerateISM2DScene(payload)
        const data = res && res.data && typeof res.data === "object" && res.data.data ? res.data.data : res.data
        const scene = this.normalizeAiScene(data)
        if(scene) {
          return scene
        }
      } catch (e) {
        console.warn("[AI Scene] backend unavailable, fallback to local generator", e)
      }
      const fallbackScene = buildLocalAiSceneGraph({
          prompt: payload.prompt,
          width: payload.canvas.width,
          height: payload.canvas.height,
          append: payload.mode === "append"
        })
      return {
        ...this.normalizeAiScene(fallbackScene),
        fromFallback:true
      }
    },
    normalizeAiScene(data) {
      if(!data || typeof data !== "object") {
        return null
      }
      const normalized = normalizeISMScene(data, {
        width: parseInt(this.aiSceneWidth || this.configData.layer.width || 1920),
        height: parseInt(this.aiSceneHeight || this.configData.layer.height || 1080),
        backColor: this.configData.layer.backColor || "#061321",
        allowedShapes: this.availableAiShapes
      })
      if(!normalized.components.cells.length) {
        return null
      }
      return {
        ...normalized,
        fromFallback: !!data.fromFallback
      }
    },
    SetIndexUp() {
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        // 获取当前节点的zIndex
        const currentZIndex = NodeInfo.getZIndex()
        // 上移一层
        NodeInfo.setZIndex(currentZIndex + 1)
      }
    },
    SetIndexDown() {
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        // 获取当前节点的zIndex
        const currentZIndex = NodeInfo.getZIndex()
// 下移一层
        NodeInfo.setZIndex(Math.max(-20, currentZIndex - 1))
      }
    },
    SetIndexFront(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        NodeInfo.toFront()
        console.log(NodeInfo.getZIndex())
      }
    },
    SetIndexBack(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        NodeInfo.toBack()
      }
    },
    FlipVertical(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        const tdata = NodeInfo.getData()
        if(tdata.detail.style.transform==-1099)
        {
          tdata.detail.style.transform=0
        }
        else
        {
          tdata.detail.style.transform = -1099
        }
        this.UpdateNodeDataFlag=!this.UpdateNodeDataFlag
        this.selectedNode.setData({
          UpdateNodeFlag:this.UpdateNodeDataFlag,
          detail:tdata.detail
        },{ overwrite: true })
      }
    },
    FlipHorizontally(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        const tdata = NodeInfo.getData()
        if(tdata.detail.style.transform==-1098)
        {
          tdata.detail.style.transform=0
        }
        else
        {
          tdata.detail.style.transform = -1098
        }
        this.UpdateNodeDataFlag=!this.UpdateNodeDataFlag
        this.selectedNode.setData({
          UpdateNodeFlag:this.UpdateNodeDataFlag,
          detail:tdata.detail
        },{ overwrite: true })
      }
    },
    revolve(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        let NodeAngle = NodeInfo.prop().angle
        NodeAngle = parseInt(NodeAngle)+90
        if(NodeAngle>=360)
        {
          NodeAngle=0
        }
        NodeInfo.rotate(parseInt(NodeAngle),{ absolute: true });
      }
    },
    reverse(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        let NodeAngle = NodeInfo.prop().angle
        NodeAngle = parseInt(NodeAngle)-90
        if(NodeAngle<=-360)
        {
          NodeAngle=0
        }
        NodeInfo.rotate(parseInt(NodeAngle),{ absolute: true });
      }
    },
    doWithUndo(){
      this.ISMCavasContainer.undo()
    },
    doWithRedo(){
      this.ISMCavasContainer.redo()
    },
    DelNode(){
      let _t = this
      this.$confirm({
        content: _t.$t('dataModel.deleteConfirm'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          _t.setLayerSelected(true)
          const selectedCells = _t.ISMCavasContainer.getSelectedCells();
          if(selectedCells.length == 0)
          {
            _t.ISMCavasContainer.removeCell(_t.UnSelectedComponent)
          }
          else {
            _t.ISMCavasContainer.removeCells(selectedCells);
          }
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {
          _t.setLayerSelected(true)
        },
      });
    },
    FormatDo(){
      this.isFormatPainterActiveTemp = !this.isFormatPainterActiveTemp
      this.SetFormatPainterState(this.isFormatPainterActiveTemp)
    },
    SelectAllComments(){ // 设定复制源
      this.ISMCavasContainer.select(this.ISMCavasContainer.getCells())
    },
    copyItem(){ // 设定复制源
      this.isCommentRightClick = false
      let items = [];
      for(let key in this.selectedComponentMap) {
        let item = deepCopy(this.selectedComponentMap[key]);
        items.push(item);
      }
      this.setCopySrcItems(items);
    },
    choiceModelJson(){
      this.$refs.filElem.dispatchEvent(new MouseEvent('click'))
    },
    showFile(input) {
      //支持chrome IE10
      if (window.FileReader) {
        let file = input.target.files[0];
        let reader = new FileReader();
        reader.onload=((event)=>{
          //显示文件
          const data = sm4DecryptBase64(event.target.result)
          let pageView = {}
          try{
            pageView = JSON.parse(data)
            const importedGraph = 'GraphView' in pageView
              ? pageView.GraphView
              : (pageView.components || pageView)
            const normalizedImport = normalizeISMScene({
              layer: pageView.GraphElements || this.configData.layer || {},
              components: importedGraph
            })
            this.setLayerData(normalizedImport.components)
            if ('GraphElements' in pageView) {
              this.ISMCavasContainer.drawBackground ( {
                color:pageView.GraphElements.backColor,
                image: pageView.GraphElements.backgroundImage,
                size: '100% 100%',    // 关键参数：填充整个容器
                repeat: 'no-repeat',
                position:"center",
                quality:1,
              });
              this.configData.layer.backColor = pageView.GraphElements.backColor
              this.configData.layer.backgroundImage = pageView.GraphElements.backgroundImage
              this.configData.layer.height = pageView.GraphElements.height
              this.configData.layer.width = pageView.GraphElements.width
              this.ISMCavasContainer.resize(pageView.GraphElements.width, pageView.GraphElements.height);
            }
          }catch (e) {
            this.$message.error("模版文件解析失败!")
          }

          // this.configData=event.target.result;
        })
        reader.readAsText(file);
      }
      else {
        alert("FileReader Not supported by your browser!");
      }
      input.target.value = '';
    },
    onZoomChange(event) {
      this.ISMCavasContainer.zoomTo(event/100.0,{ center: { x: 0, y: 0 }})
    },
    showSystemImageModel(showType){
      this.$refs.HeaderSystemImageModel.showModal(showType)
    },
    doSaveLayerData(uuid){
      let _t = this
      let params = {
        uuid:uuid,
        pageid:this.selectPageUuid,
        LayerData:null,
      }
      const LayerData = JSON.parse(JSON.stringify(this.configData))
      const normalizedLayerData = normalizeISMScene({
        layer: LayerData.layer,
        components: this.ISMCavasContainer.toJSON()
      })
      LayerData.layer = {
        ...LayerData.layer,
        ...normalizedLayerData.layer
      }
      LayerData.components = normalizedLayerData.components
      params.LayerData = LayerData
      this.saveLayerDataStruct(params).then(function (res){
        if(res.data.code == 200)
        {
          let uid = _t.$route.params.uid
          _t.updateAllLayerDataStruct({pageType:_t.isMobile,uuid:uid,cb:function (){}});
          // _t.SyncLayerData(LayerData)
          _t.isCharge=false
          _t.$message.success(_t.$t('displayModel.SaveDataSuccess'))
        }
        else
        {
          _t.$message.error(_t.$t('displayModel.SaveDataFailed'))
        }
      })
    },
    zoomOut(){
      this.selectedValueTemp = this.selectedValueTemp-25
      if(this.selectedValueTemp <25)
      {
        this.selectedValueTemp=25
      }
      this.ISMCavasContainer.zoomTo(this.selectedValueTemp/100.0,{ center: { x: 0, y: 0 } })
    },
    zoomIn(){
      this.selectedValueTemp=this.selectedValueTemp+25
      if(this.selectedValueTemp >200)
      {
        this.selectedValueTemp=200
      }
      this.ISMCavasContainer.zoomTo(this.selectedValueTemp/100.0,{ center: { x: 0, y: 0 }})
    },
    pageView() {
      try{
        localStorage.setItem(this.$route.params.uid,JSON.stringify(this.configData));
      }catch (e) {
        console.log(e)
      }

      let {href} = this.$router.resolve({
        path: '/AppRun/'+this.$route.params.uid,
        name: '组态预览',
        query: {
          sceneId: this.sceneId,
          sceneName: this.sceneName
        },
        params: {
          sceneId: this.sceneId,
          sceneName: this.sceneName
        }
      });
      window.open(href, '_blank');
    },
    exportSvg(){
      this.ISMCavasContainer.exportSVG(this.configData.name+".svg",{preserveBackground:true,copyStyles:false,preserveDimensions:true})
    },
    exportJson(){
      const normalizedLayerData = normalizeISMScene({
        layer: this.configData.layer,
        components: this.ISMCavasContainer.toJSON()
      })
      const dataJson = normalizedLayerData.components
      let exportJsonData={
        GraphView:dataJson,
        GraphElements:{
          height: normalizedLayerData.layer.height,
          width: normalizedLayerData.layer.width,
          backColor: normalizedLayerData.layer.backColor,
          backgroundImage: normalizedLayerData.layer.backgroundImage,
        }
      }
      let data = JSON.stringify(exportJsonData)
      data = sm4EncryptBase64(data)
      //encodeURIComponent解决中文乱码
      let uri = 'data:text/csv;charset=utf-8,\ufeff' + encodeURIComponent(data);
      //通过创建a标签实现
      let link = document.createElement("a");
      link.href = uri;
      //对下载的文件命名
      link.download =  this.configData.name+".json";
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
    beforeunloadHandler (e) {
      if(!this.isCharge)
      {
        return
      }
      this._beforeUnload_time = new Date().getTime()
      e = e || window.event
      if (e) {
        e.returnValue = '关闭提示'
      }
      return '关闭提示'
    },
    unloadHandler () {
      if(!this.isCharge)
      {
        return
      }
      this._gap_time = new Date().getTime() - this._beforeUnload_time
      // 判断是窗口关闭还是刷新
      if (this._gap_time <= 5) {
        console.log("关闭提示")
      }
    },
    saveContent(event) {
      if (event.keyCode === 83 && event.ctrlKey) {
        this.doSaveLayerData(this.$route.params.uid)
        event.preventDefault()
      }
    },
    onSelectImage(){},
    HeaderAlignNodesLeft(){
      const selectedNodes = this.ISMCavasContainer.getSelectedCells()
        .filter(cell => cell.isNode())
        .sort((a, b) => a.getBBox().x - b.getBBox().x);

      if (selectedNodes.length < 2) return;

      const baseX = selectedNodes[0].getBBox().x; // 获取最左侧节点的X坐标

      // 从第二个节点开始对齐（跳过基准节点）
      selectedNodes.slice(1).forEach(node => {
        node.setPosition(
          baseX,  // 统一X坐标
          node.getPosition().y // 保持原有Y坐标
        );
      });
    },
    HeaderAlignNodesRight(){
      // 获取选中节点并按X坐标降序排序
      const selectedNodes = this.ISMCavasContainer.getSelectedCells()
        .filter(cell => cell.isNode())
        .sort((a, b) => b.getBBox().x - a.getBBox().x);

      if (selectedNodes.length < 2) return;

      const baseX = selectedNodes[0].getBBox().x; // 获取最右侧节点的X坐标
      
      // 从第二个节点开始对齐（跳过基准节点）
      selectedNodes.slice(1).forEach(node => {
        node.setPosition(
          baseX,  // 统一X坐标
          node.getPosition().y // 保持原有Y坐标
        );
      });
    },
    HeaderAlignNodesTop(){
       // 获取选中节点并按Y坐标升序排序
        const selectedNodes = this.ISMCavasContainer.getSelectedCells()
          .filter(cell => cell.isNode())
          .sort((a, b) => a.getBBox().y - b.getBBox().y);

        if (selectedNodes.length < 2) return;

        const baseY = selectedNodes[0].getBBox().y; // 获取最上方节点的Y坐标
        
        // 从第二个节点开始对齐（跳过基准节点）
        selectedNodes.slice(1).forEach(node => {
          node.setPosition(
            node.getPosition().x, // 保持原有X坐标
            baseY  // 统一Y坐标
          );
        });
    },
    HeaderAlignNodesBottom(){
       // 获取选中节点并按底部坐标降序排序（Y坐标+高度）
        const selectedNodes = this.ISMCavasContainer.getSelectedCells()
          .filter(cell => cell.isNode())
          .sort((a, b) => {
            const bboxA = a.getBBox();
            const bboxB = b.getBBox();
            return (bboxB.y + bboxB.height) - (bboxA.y + bboxA.height);  // 按底部位置排序:ml-citation{ref="8" data="citationList"}
          });

        if (selectedNodes.length < 2) return;

        // 获取基准节点（最下方的节点）
        const baseNode = selectedNodes[0];
        const baseBBox = baseNode.getBBox();
        const baseBottom = baseBBox.y + baseBBox.height;  // 底部Y坐标:ml-citation{ref="1,8" data="citationList"}

        // 对齐其他节点（跳过基准节点）
        selectedNodes.slice(1).forEach(node => {
          const bbox = node.getBBox();
          const newY = baseBottom - bbox.height;  // 计算新位置:ml-citation{ref="8" data="citationList"}
          node.setPosition(bbox.x, newY);  // 保持X坐标不变，更新Y位置:ml-citation{ref="1" data="citationList"}
        });
    },
    HeaderSetCommentsAlign(Align){
      switch(Align) {
        case 'r':{
          this.HeaderAlignNodesRight()
          break
        }
        case 'l':{
          this.HeaderAlignNodesLeft()
          break
        }
        case 't':{
          this.HeaderAlignNodesTop()
          break
        }
        case 'b':{
          this.HeaderAlignNodesBottom()
          break
        }
        case 'Vertical':{
          this.HeaderArrangeNodesVertically()
          break
        }
        case 'Horizontal':{
          this.HeaderArrangeNodesHorizontally()
          break
        }
      }
    },
    //垂直等间距
    HeaderArrangeNodesVertically(){
      const selectedNodes = this.ISMCavasContainer.getSelectedCells()
          .filter(cell => cell.isNode())
          .sort((a, b) => a.getBBox().y - b.getBBox().y);

      if (selectedNodes.length < 2) return;

      const baseNode = selectedNodes[0]; // 保持Y最小的基准节点
      const maxY = selectedNodes[selectedNodes.length-1].getBBox().y;
      const minY = baseNode.getBBox().y;
      const deltaY = maxY - minY;

      // 计算总高度（排除基准节点）
      const totalHeight = selectedNodes.slice(1).reduce((sum, node) => {
        return sum + node.getSize().height;
      }, 0);

      // 计算间距 = (极差 - 总高度) / (节点数-1)
      const spacing = (deltaY - totalHeight) / (selectedNodes.length - 1);
      let currentY = minY + baseNode.getSize().height + spacing;

      // 从第二个节点开始定位（跳过基准节点）
      selectedNodes.slice(1).forEach(node => {
        const pos = node.getPosition()
        node.setPosition(
            pos.x, // X坐标对齐
            currentY
        );
        currentY += node.getSize().height + spacing;
      });
    },
    HeaderArrangeNodesHorizontally(){
         // 获取选中节点并按X坐标升序排序
      const selectedNodes = this.ISMCavasContainer.getSelectedCells()
        .filter(cell => cell.isNode())
        .sort((a, b) => a.getBBox().x - b.getBBox().x);

      if (selectedNodes.length < 2) return;

      const baseNode = selectedNodes[0]; // 保持X最小的基准节点
      const maxX = selectedNodes[selectedNodes.length-1].getBBox().x;
      const minX = baseNode.getBBox().x;
      const deltaX = maxX - minX;

      // 计算总宽度（排除基准节点）
      const totalWidth = selectedNodes.slice(1).reduce((sum, node) => {
        return sum + node.getSize().width;
      }, 0);

      // 计算间距 = (极差 - 总宽度) / (节点数-1)
      const spacing = (deltaX - totalWidth) / (selectedNodes.length - 1);
      let currentX = minX + baseNode.getSize().width + spacing;

      // 从第二个节点开始定位（跳过基准节点）
      selectedNodes.slice(1).forEach(node => {
        node.setPosition(
          currentX,
          baseNode.getPosition().y // Y坐标对齐
        );
        currentX += node.getSize().width + spacing;
      });
    },
    HeaderCenterNodesHorizontally(){
      const nodes = this.ISMCavasContainer.getSelectedCells().filter(cell => cell.isNode())
      if (nodes.length === 0) return

      // 计算节点组总宽度和最小X坐标
      let minX = Infinity
      let maxX = -Infinity
      nodes.forEach(node => {
        const bbox = node.getBBox()
        minX = Math.min(minX, bbox.x)
        maxX = Math.max(maxX, bbox.x + bbox.width)
      })
      const totalWidth = maxX - minX

      // 计算居中偏移量
      const viewport = this.ISMCavasContainer.getGraphArea()
      const offsetX = (viewport.width - totalWidth) / 2 - minX

      // 批量移动节点
      nodes.forEach(node => {
        const pos = node.getPosition()
        node.setPosition(pos.x + offsetX, pos.y)
      })
    },
    HeaderLockItem(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {

        let NodeInfo = selectedCells[i]
        const parent = NodeInfo.getParent();
        if(parent) {
          parent.getChildren().forEach(child => {
            let NodeInfoData = child.getData()
            NodeInfoData.locked = true
            child.setData(NodeInfoData, {overwrite: true})
          })
        }else {
          let NodeInfoData = NodeInfo.getData()
          NodeInfoData.locked = true
          NodeInfo.setData(NodeInfoData, {overwrite: true})
        }
      }
      this.ISMCavasContainer.cleanSelection();
    },
    HeaderUnlockItem(){
      const selectedCells = this.ISMCavasContainer.getCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        const parent = NodeInfo.getParent();
        if(parent) {
          parent.getChildren().forEach(child => {
            let NodeInfoData = child.getData()
            NodeInfoData.locked = false
            child.setData(NodeInfoData, {overwrite: true})
          })
        }else {
          let NodeInfoData = NodeInfo.getData()
          NodeInfoData.locked = false
          NodeInfo.setData(NodeInfoData, {overwrite: true})
        }
      }
      this.ISMCavasContainer.select(selectedCells)
    },
    HeaderCreateGroup(){
      let _t = this
      const selectedCells = _t.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        const Parent  = NodeInfo.getParent()
        if(Parent)
        {
          _t.$message.info('不能将节点添加到组中')
          return
        }
      }
      const bboxes = selectedCells.map(cell => cell.getBBox())
      const left = Math.min(...bboxes.map(b => b.x))
      const top = Math.min(...bboxes.map(b => b.y))
      const right = Math.max(...bboxes.map(b => b.x + b.width))
      const bottom = Math.max(...bboxes.map(b => b.y + b.height))

      const parent = _t.ISMCavasContainer.addNode({
        shape: 'view-ism-group-node',
        x: left,
        y: top,
        width: right-left,
        height: bottom-top,
        zIndex: -1000,
        data: {
          locked:false,
          UpdateNodeFlag:true,
          editMode: true,
          showDeviceUuid:"",
          IsToolBox:false,
          detail:{
            identifier :uuid.v1(),
            name:"节点组",
            "type": "image",
            isCanvas:true,
            "action": [],
            "dataBind":[],
            "active": [
              {
                id:"Forward",
                name:"component.ViewCanvasMoveLineArrow.Forward",
                result:"",
                isExpression:true,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "",
                  dataName: "",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                },
              },
              {
                id:"Reverse",
                name:"component.ViewCanvasMoveLineArrow.Reverse",
                result:"",
                isExpression:true,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "",
                  dataName: "",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                },
              },
            ],
            "animate": {
              "selected": [],
              "condition":{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
              "isExpression": false,
              "animateList": [

              ],
              "animateElement": [
                {
                  id: "blink",
                  elementList:[
                    {
                      "name":"component.public.animateSpeed",
                      "type":7,
                      "value":1,
                      "min":0.1,
                      "key":"blinkSpeed",
                    },
                  ]
                },
                {
                  id: "millcolorGrad",
                  elementList:[
                    {
                      "name": "component.public.startColor",
                      "type": 2,
                      "value": "#74f808",
                      "key": "startColor",
                    },
                    {
                      "name": "component.public.stopColor",
                      "type": 2,
                      "value": "#f30b0b",
                      "key": "stopColor",
                    },
                    {
                      "name":"component.public.animateSpeed",
                      "type":7,
                      "value":1,
                      "min":0.1,
                      "key":"animateSpeed",
                    },
                  ]
                },
                {
                  id: "animateSpin",
                  elementList:[
                    {
                      "name":"component.public.animateSpinSpeed",
                      "type":7,
                      "value":1,
                      "min":0.1,
                      "key":"spinSpeed",
                    },
                    {
                      name:"configComponent.bigScreen.border.border89Direction",
                      type:6,
                      value:0,
                      enumList:[
                        {
                          value:0,
                          option:"configComponent.bigScreen.border.border89DirectionForward"
                        },
                        {
                          value:1,
                          option:"configComponent.bigScreen.border.border89DirectionNegative"
                        }
                      ],
                      min:1,
                      key:"spinDirection",
                    }
                  ]
                },
              ],
            },
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": right-left,
                "h": bottom-top
              },
              "points": [],
              "visible":1,
              "zIndex": -1000,
              "transform": 0,
              "backColor": "",
              foreColor:"",
              borderWidth:2,
              BorderEdges:0,
              opacity:1,
              borderStyle:"solid",
              borderColor:"#13c2c2",
              "diy":[

              ]
            }
          }
        },
      })
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        parent.addChild(NodeInfo)
      }
    },
    HeaderSplitGroup() {
      let _t = this
      const selected = _t.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selected.length;i++)
      {
        let NodeInfo = selected[i]
        if(NodeInfo.prop().shape=="view-ism-group-node")
        {
          NodeInfo.setChildren(null)
          _t.ISMCavasContainer.removeNode(NodeInfo)
        }
        else
        {
          const parent = NodeInfo.getParent();
          if(parent) {
            parent.setChildren(null)
            _t.ISMCavasContainer.removeNode(parent)
          }
        }
      }
    },
  },
  computed: {
    ...mapState('setting', ['langList','isMobile','lang',]),
    ...mapState({
      isFormatPainterActive:state => store.state.ISMDisPlayEditorTool.isFormatPainterActive,
      UnSelectedComponent: state => store.state.ISMDisPlayEditorTool.UnSelectedComponent,
      configData: state => store.state.ISMDisPlayEditorTool.LayerData,
      selectPageUuid: state => store.state.ISMDisPlayEditorTool.selectPageUuid,
      selectedComponentMap: state => store.state.ISMDisPlayEditorTool.selectedComponentMap,
      curComponent: state => store.state.ISMDisPlayEditorTool.selectedComponent,
      copySrcItems: state => store.state.ISMDisPlayEditorTool.copySrcItems,
      copyCount: state => store.state.ISMDisPlayEditorTool.copyCount,
      isHistoryOp: state => store.state.ISMDisPlayEditorTool.isHistoryOp,
      selectedValue:state => store.state.ISMDisPlayEditorTool.selectedValue,
      selectedNode:state => store.state.ISMDisPlayEditorTool.selectedNode,
      ISMCavasContainer:state => store.state.ISMDisPlayEditorTool.ISMCavasContainer,
      toolBoxList: state => store.state.ISMDisPlayEditorTool.toolBoxList,
      PageCanVasList: state => store.state.ISMDisPlayEditorTool.PageCanVasList,
      MesComponentsList: state => store.state.ISMDisPlayEditorTool.MesComponentsList,
      DiyComponentsList: state => store.state.ISMDisPlayEditorTool.DiyComponentsList,
    }),
    availableAiShapes() {
      const groups = [
        ...(Array.isArray(this.toolBoxList) ? this.toolBoxList : []),
        ...(Array.isArray(this.PageCanVasList) ? this.PageCanVasList : []),
        ...(Array.isArray(this.MesComponentsList) ? this.MesComponentsList : []),
        ...(Array.isArray(this.DiyComponentsList) ? this.DiyComponentsList : [])
      ]
      const shapes = []
      groups.forEach(group => {
        const items = Array.isArray(group.items) ? group.items : []
        items.forEach(item => {
          const type = item && item.info && item.info.type
          if(type && !shapes.includes(type)) {
            shapes.push(type)
          }
        })
      })
      return shapes.length ? shapes : [
        "view-svg-text",
        "view-svg-time",
        "view-progress-bars",
        "dv-border-box1",
        "dv-border-box8",
        "ViewCanvasMoveLineArrow",
        "view-svg-electric1",
        "view-svg-electric2",
        "view-svg-electric3",
        "view-svg-electric4",
        "view-svg-electric5",
        "ism-view-device-status",
        "ism-view-real-data-smooth-chart-by-device",
        "ism-view-history-trend-chart",
        "ism-view-history-difference",
        "ism-view-real-table",
        "view-device-alarm-list",
        "view-device-real-data-table"
      ]
    },
    systemName () {
      return this.$store.state.setting.systemName
    },
    Isometric_row:{
      get(){
        return this.$store.state.ISMDisPlayEditorTool.Isometric_row;
      },
      set(v) {
        this.$store.state.ISMDisPlayEditorTool.Isometric_row = v
      }
    },
    Isometric_colu:{
      get(){
        return this.$store.state.ISMDisPlayEditorTool.Isometric_colu;
      },
      set(v) {
        this.$store.state.ISMDisPlayEditorTool.Isometric_colu = v
      }
    }
  },
  created() {

  },
  mounted () {
      document.addEventListener('keydown', this.saveContent)
      window.addEventListener('beforeunload', e => this.beforeunloadHandler(e))
      window.addEventListener('unload', e => this.unloadHandler(e))
      this.selectedValueTemp = this.selectedValue
      // 初始化 undo 快照：确保组件 remount 后首次编辑仍可撤销
      if(this.configData && this.configData.components && this.configData.components.length > 0) {
        this._prevComponents = deepCopy(this.configData.components)
      }
  },
  destroyed () {
    document.removeEventListener('keydown', this.saveContent)
    window.removeEventListener('beforeunload', e => this.beforeunloadHandler(e))
    window.removeEventListener('unload', e => this.unloadHandler(e))
  },

}
</script>

<style scoped>
::v-deep .ant-btn {
   padding: 0 6px;
}
.mymenu {
  padding: 0px 0;
  border-color: transparent;
  color: #090909;
  background: #fafafa;
}
.mymenu .menu-active {
  z-index: 9999999999999999;
  border-color: #b7d2ff;
  color: #000000;
  background: #eaf2ff;
}

.m-btn-downarrow{
  display: none;
}

.charge-tips{
  white-space: nowrap;
  font-size: 14px;
  padding: 6px;
  color: #a94442 !important;
  border-radius: 3px;
  cursor: pointer;
  float: right;
  margin-right: 30px;
}

.help-content{
  padding: 10px;
}
.help-content p{
  color: #0d1a26;
  font-weight: 500;
  margin-bottom: 20px;
  margin-top: 8px;
  font-family: Avenir,-apple-system,BlinkMacSystemFont,Segoe UI,PingFang SC,Hiragino Sans GB,Microsoft YaHei,Helvetica Neue,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol;
  font-size: 14px;
  font-variant: tabular-nums;
  line-height: 18px;
}

#toolbox .panel-header {
  border-width: 0px;
  border-style: solid;
}
</style>

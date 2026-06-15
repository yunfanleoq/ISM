<template>
  <a-spin :spinning="spinning" size="large" tip="Loading..." style="height: 100%">
    <div class="x6-container" :style="ContainerStyle">
      <div ref="rulerHorizontal" class="ruler-horizontal" v-if="UpdateRuler"></div>
      <div ref="rulerVertical" class="ruler-vertical" v-if="UpdateRuler"></div>
      <div ref="ISMContainer" class="graph-container"  :class="{ 'format-painter' : isFormatPainterActive,'animated':true,[`${configData.layer.animate}`]: true}"></div>
<!--      <div class="minimap-container" ref="refMiniMapContainer" />-->
    </div>
  </a-spin>
</template>

<script>
import { Graph, Shape } from '@antv/x6'
import Vue from 'vue'
import store from "../../store";
import Contextmenu from "vue-contextmenujs"
import { Transform } from '@antv/x6'
import { Snapline } from '@antv/x6'
import { Clipboard } from '@antv/x6'
import { Keyboard } from '@antv/x6'
import { History } from '@antv/x6'
import { Selection } from '@antv/x6'
import { Export } from '@antv/x6'
import { Scroller } from '@antv/x6'
import { MiniMap } from '@antv/x6'
import { mapActions, mapState, mapMutations } from 'vuex'
import Guides from "@scena/guides";
import { initDnd } from '@/utils/dnd'
import {register} from "@antv/x6-vue-shape";
import ISMGroupNode from "@/pages/ISMDisPlay/ISMGroupNode.vue";
import {uuid} from "vue-uuid";
import {ClickUnSelectedComponent} from "@/store/ISM/mutations";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
Vue.use(Contextmenu);
export default {
  mixins: [ISMChildAutoMixin],
  name: "ISMCanvas",
  i18n: require('../../i18n/language'),
  components: {

  },
  watch:{
    'configData.layer': {
      handler(newVal,oldvalue) {
        let _t=this
        this.ISMCavasContainer.resize(newVal.width,newVal.height)
        this.ISMCavasContainer.drawBackground ( {
          color:newVal.backColor,
          image: newVal.backgroundImage,
          size: '100% 100%',    // 关键参数：填充整个容器
          repeat: 'no-repeat',
          position:"center",
          quality:1,
        });
        this.UpdateRuler = false
        this.$refs.rulerHorizontal = null
        this.$refs.rulerVertical = null
        if(this.guidesX!=null) {
          this.guidesX.destroy()
        }
        if(this.guides!=null) {
          this.guides.destroy()
        }
        setTimeout(function (){
          _t.UpdateRuler = true
          setTimeout(function (){
            _t.initRuler(parseInt( _t.configData.layer.width),parseInt( _t.configData.layer.height))
          },300)
        },200)
      },
      deep: true
    },
    'configData.components': {
      handler(newVal,oldvalue) {
        let _t = this
        this.tempScale = this.ISMCavasContainer.transform.getZoom();
        const cells = this.ISMCavasContainer.getSelectedCells();
        this.ISMCavasContainer.cleanSelection();
        try{
          const components = JSON.parse(JSON.stringify(this.configData.components))
          if (components.cells && Array.isArray(components.cells)) {
            components.cells = components.cells.filter(cell => cell && cell.shape)
          }
          this.ISMCavasContainer.fromJSON(components)
        }catch (e){
          console.log(e)
        }
        this.ISMCavasContainer.zoomTo(this.tempScale,{ center: { x: 0, y: 0 }})
        this.ISMCavasContainer.select(cells);
      }
    },
  },
  data() {
    return {
      Connector:false,
      DrawPoints:[],
      beginDrawLine: false,
      UpdateRuler:false,
      spinning:false,
      tempScale:1,
      guidesX:null,
      guides:null,
      IsShowMinMap:false,
      IsShowMakerLine:false,
      isCommentRightClick:false,
      ISMCavasContainer:null,
      UpdateNodeDataFlag:true,
      dnd:null,
      drawingEdge:null,
    }
  },
  props: [],
  computed: {
    ...mapState('setting', ['langList','isMobile','lang',]),
    ...mapState({
      FormatPainterCell:state => store.state.ISMDisPlayEditorTool.FormatPainterCell,
      isFormatPainterActive:state => store.state.ISMDisPlayEditorTool.isFormatPainterActive,
      selectedComponents: state => store.state.ISMDisPlayEditorTool.selectedComponents,
      UnSelectedComponent: state => store.state.ISMDisPlayEditorTool.UnSelectedComponent,
      selectedComponentMap: state => store.state.ISMDisPlayEditorTool.selectedComponentMap,
      configData: state => store.state.ISMDisPlayEditorTool.LayerData,
      selectedNode: state => store.state.ISMDisPlayEditorTool.selectedNode,
    }),
    ContainerStyle() {
      return {
        "--rulerHeight":(this.configData.layer.height-20)+"px",
        "--rulerWidth":(this.configData.layer.width-20)+"px",
      };
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
  methods: {
    ...mapMutations('ISMDisPlayEditorTool',[
      'setCurrentISMCavasContainer',
        'setCurrentISMCavasDND',
        'setLayerSelected',
        'setCurrentSelectNode',
      'SetFormatPainterState',
        'ClickUnSelectedComponent'
    ]),
    ...mapActions('ISMDisPlayEditorTool',[
      'getLayerDataStruct',
      'setGroupList',
      'setLayerData',
      'SyncLayerComponents'
    ]),
    initRuler(width,height){
      const horizontalConfig = {
        // 标尺的位置，默认是end
        direction: 'end',
        // ruler 的背景色
        backgroundColor: 'transparent',
        lineColor:"#000",
        // 刻度文字的颜色
        textColor: '#000',
        // 刻度文字的偏移值
        textOffset: [-2, 7],
        width: width,
        height: 20,
        // 拖拽线的时候是否实时显示坐标位置
        displayDragPos: true,
        // 刻度的间隔值，默认50
        unit: 100,
        // guide 拖动结束后线的颜色，配合dragGuideStyle使用
        guideStyle: { background: '#05e2a8' },
        // guide 拖动中线的颜色 配合guideStyle使用
        dragGuideStyle: { background: '#05e2a8' },
        // 实时显示坐标的样式
        guidePosStyle: { color: '#05e2a8' }
      };
      const verticalConfig = {
        // 标尺的位置，默认是end
        direction: 'end',
        // ruler 的背景色
        backgroundColor: 'transparent',
        lineColor:"#000",
        height:height,
        width:20,
        // 刻度文字的颜色
        textColor: '#000',
        // 刻度文字的偏移值
        textOffset: [10, 2],
        // 拖拽线的时候是否实时显示坐标位置
        displayDragPos: true,
        // 刻度的间隔值，默认50
        unit: 100,
        // guide 拖动结束后线的颜色，配合dragGuideStyle使用
        guideStyle: { background: '#05e2a8' },
        // guide 拖动中线的颜色 配合guideStyle使用
        dragGuideStyle: { background: '#05e2a8' },
        // 实时显示坐标的样式
        guidePosStyle: { color: '#05e2a8' }
      };
      this.guidesX = new Guides(this.$refs.rulerHorizontal, {
        type: 'horizontal',
        ...horizontalConfig
      });
      this.guides =new Guides(this.$refs.rulerVertical, {
        type: 'vertical',
        ...verticalConfig
      });
    },
    createNode(type, x, y) {
      const baseNode = {
        x, y,
        width: 100,
        height: 40,
        attrs: {
          label: { text: type }
        }
      }
      switch(type) {
        case 'rectangle': return new Shape.Rect(baseNode)
        case 'circle': return new Shape.Circle({ ...baseNode, height: 100 })
        default:

      }

    },
    copyNode(){
      const cells = this.ISMCavasContainer.getSelectedCells()
      if (cells && cells.length) {
        this.ISMCavasContainer.copy(cells,{ deep: true,useLocalStorage:true })
        this.$message.success('复制成功')
      } else {
        this.$message.info('请先选中节点再复制')
      }
    },
    PasteNode(){
      if (this.ISMCavasContainer.isClipboardEmpty()) {
        this.$message.info('剪切板为空，不可粘贴')
      } else {
        const cells = this.ISMCavasContainer.paste({ offset: 32,useLocalStorage: true })
        cells.forEach(cell => {
          let cellData = cell.prop().data || {}
          cellData.detail.identifier = uuid.v1()
          cellData.detail.name = cellData.detail.name+"_copy"
          cell.setData({
            locked:false,
            editMode: true,
            showDeviceUuid:"",
            IsToolBox:false,
            UpdateNodeFlag: new Date().getTime(),
            detail: cellData.detail
          }, {overwrite: true})
        })
        this.ISMCavasContainer.cleanSelection()
        this.ISMCavasContainer.select(cells)
        this.$message.success('粘贴成功')
      }
    },
    ShowOrHideMinMap(){
      this.$refs.refMiniMapContainer.style.display =  this.IsShowMinMap ? 'block' : 'none'
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
    LockItem(){
      const selectedCells = this.ISMCavasContainer.getSelectedCells();
      for(let i=0;i<selectedCells.length;i++)
      {
        let NodeInfo = selectedCells[i]
        const Parent  = NodeInfo.getParent()
        if(Parent)
        {
          Parent.getChildren().forEach(child => {
            let NodeInfoData = child.getData()
            NodeInfoData.locked = true
            NodeInfo.setData(NodeInfoData, {overwrite: true})
          })
        }
        let NodeInfoData = NodeInfo.getData()
        NodeInfoData.locked = true
        NodeInfo.setData(NodeInfoData, {overwrite: true})
      }
      this.ISMCavasContainer.cleanSelection();
    },
    UnlockItem(node){
      let NodeInfo = node
      const Parent  = NodeInfo.getParent()
      if(Parent)
      {
        Parent.getChildren().forEach(child => {
          let NodeInfoData = child.getData()
          NodeInfoData.locked = false
          NodeInfo.setData(NodeInfoData, {overwrite: true})
        })
      }
      let NodeInfoData = NodeInfo.getData()
      NodeInfoData.locked = false
      NodeInfo.setData(NodeInfoData, {overwrite: true})
    },
    CreateGroup(){
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
    splitGroup() {
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
    alignNodesLeft(){
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
    alignNodesRight(){
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
    alignNodesTop(){
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
    alignNodesBottom(){
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
    setCommentsAlign(Align){
      switch(Align) {
        case 'r':{
          this.alignNodesRight()
          break
        }
        case 'l':{
          this.alignNodesLeft()
          break
        }
        case 't':{
          this.alignNodesTop()
          break
        }
        case 'b':{
          this.alignNodesBottom()
          break
        }
        case 'Vertical':{
          this.arrangeNodesVertically()
          break
        }
        case 'Horizontal':{
          this.arrangeNodesHorizontally()
          break
        }
      }
    },
    arrangeNodesVertically(){
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
    arrangeNodesHorizontally(){
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
    onContextLayerMenu(event) {
      let _t = this
      if(!this.isCommentRightClick) {
        this.$contextmenu({
          items: [
            {
              label: _t.$t('displayConfig.ToolBar.DrawMoveLine'),
              icon : "el-icon-guandao",
              onClick: () => {
                _t.beginDrawLine = true
              }
            },
            {
              label: _t.$t('displayConfig.ToolBar.Redo'),
              icon : "el-icon-redo",
              onClick: () => {
                _t.ISMCavasContainer.redo()
              }
            },
            {
              label: _t.$t('displayConfig.ToolBar.Undo'),
              icon: "el-icon-weibiaoti545",
              onClick: () => {
                _t.ISMCavasContainer.undo()
              }
            },
            {
              label: _t.$t('displayConfig.Canvas.menu.Paste'),
              icon: "el-icon-niantie",
              onClick: () => {
                _t.PasteNode()
              }
            },
            {
              label: _t.$t('displayConfig.Canvas.menu.Group'),
              icon: "el-icon-printer",
              children: [
                {
                  label: _t.$t('displayConfig.Canvas.menu.GroupTogether'),
                  onClick: () => {
                    _t.CreateGroup()
                  }
                },
                {
                  label: _t.$t('displayConfig.Canvas.menu.UnGroup'),
                  onClick: () => {
                    _t.splitGroup()
                  }
                }
              ]
            },
            {
              label: _t.Connector?_t.$t('displayConfig.Canvas.menu.ConnectHide'):_t.$t('displayConfig.Canvas.menu.ConnectShow'),
              icon: "el-icon-lianjiedian",
              onClick: () => {
                _t.Connector = !_t.Connector
              }
            },
            {
              label: _t.IsShowMakerLine?_t.$t('displayConfig.Canvas.menu.AuxiliaryLineClose'):_t.$t('displayConfig.Canvas.menu.AuxiliaryLine'),
              icon: "el-icon-fuzhuxian",
              onClick: () => {
                _t.IsShowMakerLine = !_t.IsShowMakerLine
                if(_t.IsShowMakerLine) {
                  _t.ISMCavasContainer.enableSnapline()
                }else{
                  _t.ISMCavasContainer.disableSnapline()
                }
              }
            },
            {
              label:  _t.IsShowMinMap?_t.$t('displayConfig.Canvas.menu.HideMinMap'):_t.$t('displayConfig.Canvas.menu.ShowMinMap'),
              icon: "el-icon-suolvetu",
              onClick: () => {
                _t.IsShowMinMap = !_t.IsShowMinMap
                _t.ShowOrHideMinMap()
              }
            },
            {
              label: _t.$t('displayConfig.Canvas.menu.PageAttr'),
              icon: "el-icon-queshengshuxing",
              onClick: () => {
                _t.isCommentRightClick = false
                _t.setLayerSelected(true);
              }
            },
          ],
          event, // 鼠标事件信息
          divided:true,
          customClass: "custom-class", // 自定义菜单 class
          zIndex: 10000, // 菜单样式 z-index
          minWidth: 230 // 主菜单最小宽度
        });
        return false;
      }
    },
    componentRightClick(node){
      let _t = this
      this.isCommentRightClick = true
      setTimeout(function (){
        _t.isCommentRightClick = false
      },500)
      this.$contextmenu({
        items: [
          {
            label: _t.$t('displayConfig.ToolBar.Redo'),
            icon: "el-icon-redo",
            onClick: () => {
              _t.ISMCavasContainer.redo()
            }
          },
          {
            label: _t.$t('displayConfig.ToolBar.Undo'),
            icon: "el-icon-weibiaoti545",
            onClick: () => {
              _t.ISMCavasContainer.undo()
            }
          },
          {
            label: _t.$t('displayConfig.Canvas.menu.Unlock'),
            icon: "el-icon-jiesuo",
            onClick: () => {
              _t.UnlockItem(node)
            }
          },
          {
            label: _t.$t('displayConfig.Canvas.menu.Lock'),
            icon: "el-icon-suoding",
            onClick: () => {
              _t.LockItem()
            }
          },
          {
            label: _t.$t('displayConfig.ToolBar.Rotate'),
            icon: "el-icon-printer",
            children: [
              {
                label: _t.$t('displayConfig.Canvas.menu.Revolve'),
                onClick: () => {
                  const selectedCells = _t.ISMCavasContainer.getSelectedCells();
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
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.Reverse'),
                onClick: () => {
                  const selectedCells = _t.ISMCavasContainer.getSelectedCells();
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
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.FlipVertical'),
                onClick: () => {
                  const selectedCells = _t.ISMCavasContainer.getSelectedCells();
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
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.FlipHorizontally'),
                onClick: () => {
                  const selectedCells = _t.ISMCavasContainer.getSelectedCells();
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
                }
              }
            ]
          },
          {
            label: _t.$t('displayConfig.ToolBar.Alignment'),
            icon: "el-icon-printer",
            children: [
              {
                label: _t.$t('displayConfig.Canvas.menu.topAlign'),
                onClick: () => {
                  _t.setCommentsAlign('t')
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.bottomAlign'),
                onClick: () => {
                  _t.setCommentsAlign('b')
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.rightAlign'),
                onClick: () => {
                  _t.setCommentsAlign('r')
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.leftAlign'),
                onClick: () => {
                  _t.setCommentsAlign('l')
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.VerticalEquidistant'),
                onClick: () => {
                  _t.setCommentsAlign('Vertical')
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.HorizontalEquidistant'),
                onClick: () => {
                  _t.setCommentsAlign('Horizontal')
                }
              }
            ]
          },
          {
            label: _t.$t('displayConfig.Canvas.menu.Group'),
            icon: "el-icon-printer",
            children: [
              {
                label: _t.$t('displayConfig.Canvas.menu.GroupTogether'),
                onClick: () => {
                  _t.CreateGroup()
                }
              },
              {
                label: _t.$t('displayConfig.Canvas.menu.UnGroup'),
                onClick: () => {
                  _t.splitGroup()
                }
              }
            ]
          },
          {
            label: _t.$t('displayConfig.Canvas.menu.Delete'),
            icon: "el-icon-shanchu",
            onClick: () => {
              _t.DelNode()
            }
          },
          {
            label: _t.$t('displayConfig.Canvas.menu.Copy'),
            icon: "el-icon-kaobei",
            onClick: () => {
              _t.copyNode()
            }
          },
          {
            label: _t.$t('displayConfig.Canvas.menu.Paste'),
            icon: "el-icon-niantie",
            onClick: () => {
              _t.PasteNode()
            }
          },
          {
            label: _t.$t('displayConfig.ToolBar.DrawMoveLine'),
            icon : "el-icon-guandao",
            onClick: () => {
              _t.beginDrawLine = true
            }
          },
          {
            label: _t.Connector?_t.$t('displayConfig.Canvas.menu.ConnectHide'):_t.$t('displayConfig.Canvas.menu.ConnectShow'),
            icon: "el-icon-lianjiedian",
            onClick: () => {
              _t.Connector = !_t.Connector
            }
          },
          {
            label: _t.IsShowMakerLine?_t.$t('displayConfig.Canvas.menu.AuxiliaryLineClose'):_t.$t('displayConfig.Canvas.menu.AuxiliaryLine'),
            icon: "el-icon-fuzhuxian",
            onClick: () => {
              _t.IsShowMakerLine = !_t.IsShowMakerLine
              if(_t.IsShowMakerLine) {
                _t.ISMCavasContainer.enableSnapline()
              }else{
                _t.ISMCavasContainer.disableSnapline()
              }
            }
          },
          {
            label:  _t.IsShowMinMap?_t.$t('displayConfig.Canvas.menu.HideMinMap'):_t.$t('displayConfig.Canvas.menu.ShowMinMap'),
            icon: "el-icon-suolvetu",
            onClick: () => {
              _t.IsShowMinMap = !_t.IsShowMinMap
              _t.ShowOrHideMinMap()
            }
          },
          {
            label: _t.$t('displayConfig.Canvas.menu.PageAttr'),
            icon: "el-icon-queshengshuxing",
            onClick: () => {
              _t.isCommentRightClick = false
              _t.setLayerSelected(true);
            }
          },
        ],
        event, // 鼠标事件信息
        divided:true,
        zIndex: 10000, // 菜单样式 z-index
        minWidth: 230 // 主菜单最小宽度
      });
      return false;
    },
    drawMove(e){
      if(!this.beginDrawLine)
      {
        return
      }
      const localPos = this.ISMCavasContainer.clientToLocal({
        x: e.clientX,
        y: e.clientY
      });
      this.ISMCavasContainer.disableSelection()
      if (this.drawingEdge) {
        this.drawingEdge.setTarget({ x:localPos.x, y:localPos.y });
      }
    },
    applyNodeProps(sourceNode,targetNode) {

      const copiedAttrs = { ...sourceNode.getAttrs() }
      // 复制自定义业务属性（data）
      const copiedData = { ...sourceNode.getData() }
      // （可选）复制节点尺寸（若需同步大小）
      const { width, height } = sourceNode.getSize()

      // 存储复制的属性（后续用于应用）
      this.copiedProps = {
        attrs: copiedAttrs,
        data: copiedData,
        size: { width, height } // 可选：尺寸同步
      }
      const { attrs, data, size } = this.copiedProps

      const targetData = { ...targetNode.getData() }
      targetData.detail.style.borderWidth=copiedData.detail.style.borderWidth
      targetData.detail.style.BorderEdges=copiedData.detail.style.BorderEdges
      targetData.detail.style.opacity=copiedData.detail.style.opacity
      targetData.detail.style.borderStyle=copiedData.detail.style.borderStyle
      targetData.detail.style.borderColor=copiedData.detail.style.borderColor
      targetData.detail.style.transform=copiedData.detail.style.transform
      targetData.detail.style.backColor=copiedData.detail.style.backColor
      targetData.detail.style.foreColor=copiedData.detail.style.foreColor
      targetData.detail.style.textAlign=copiedData.detail.style.textAlign
      targetData.detail.style.fontFamily=copiedData.detail.style.fontFamily
      targetData.detail.style.fontWeight=copiedData.detail.style.fontWeight
      targetData.detail.style.italic=copiedData.detail.style.italic
      targetData.detail.style.fontSize=copiedData.detail.style.fontSize
      // 1. 应用样式（覆盖目标节点的 attrs）
      targetNode.setAttrs(attrs)
      // 2. 应用业务属性（覆盖目标节点的 data）
      targetNode.setData(targetData)
      // 3. （可选）应用尺寸（同步节点大小）
      targetNode.setSize(size.width, size.height)

      this.$message.success('格式应用成功')
    },
    initCavasContainer(){
      let _t = this
      let startMovePos = {};
      register({
        shape: "view-ism-group-node",
        width: 100,
        height: 100,
        visible:true,
        zIndex: -1000,
        component: ISMGroupNode,
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
                "w": 600,
                "h": 50
              },
              "points": [],
              "visible":1,
              "zIndex": -1000,
              "transform": 0,
              "backColor": "#D9D9D9",
              foreColor:"#13c2c2",
              borderWidth:0,
              BorderEdges:0,
              opacity:1,
              borderStyle:"solid",
              borderColor:"#ccccff",
              "diy":[

              ]
            }
          }
        },
      })
      this.ISMCavasContainer = new Graph({
        container: this.$refs.ISMContainer,
        width: _t.configData.layer.width,
        height: _t.configData.layer.height,
        panning:false,
        virtual:false,
        async:false,
        grid: {
          visible: false,
          size: 2,
          type: 'doubleMesh',
          args: [
            {
              color: '#eee', // 主网格线颜色
              thickness: 1, // 主网格线宽度
            },
            {
              color: '#ddd', // 次网格线颜色
              thickness: 1, // 次网格线宽度
              factor: 4, // 主次网格线间隔
            },
          ],
        },
        background: {
          color: _t.configData.layer.backColor,   // 背景底色（可选）
          image: _t.configData.layer.backgroundImage,
          size: '100% 100%',
          repeat: 'no-repeat',
          position:"center",
          quality:1,
        },
        connecting: {
          router: 'manhattan',
          connector: {
            name: 'rounded',
            args: {
              radius: 8,
            },
          },
          anchor: 'center',
          connectionPoint: 'anchor',
          allowBlank: false,
          snap: {
            radius: 20,
          },
          createEdge() {
            return new Shape.Edge({
              attrs: {
                line: {
                  stroke: '#A2B1C3',
                  strokeWidth: 2,
                  targetMarker: {
                    name: 'block',
                    width: 12,
                    height: 8,
                  },
                },
              },
              zIndex: 0,
            })
          },
          validateConnection({ targetMagnet }) {
            return !!targetMagnet
          },
        },
        highlighting: {
          magnetAdsorbed: {
            name: 'stroke',
            args: {
              attrs: {
                fill: '#5F95FF',
                stroke: '#5F95FF',
              },
            },
          },
        },
        translating: {
          restrict: true,
        },
        interacting: (cellView) => {
          // 示例：根据节点数据中的 `locked` 字段判断是否可拖动
          const locked = cellView.cell.getData()?.locked;
          if (locked) {
            return {
              nodeMovable:false,
              selectable:true,
              magnetConnectable: false,  // 禁止节点移动
              edgeMovable: false,  // 禁止边移动
              edgeLabelMovable: false,  // 禁止箭头移动
              arrowheadMovable: false,  // 禁止顶点移动
              vertexMovable: false,  // 禁止添加顶点
              vertexAddable: false,  // 禁止删除顶点
              vertexDeletable: false,  // 禁止删除顶点
            };
          }
          return true; // 其他节点可拖动
        },
      })
      // 控制连接桩显示/隐藏
      const showPorts = (ports, show) => {
        for (let i = 0, len = ports.length; i < len; i += 1) {
          ports[i].style.visibility = show ? 'visible' : 'hidden'
        }
      }
      this.ISMCavasContainer.on('node:mouseenter', ({node}) => {
        startMovePos = node.position()
        let NodeInfoData = node.getData()
        if(NodeInfoData.locked)
        {
          node.addTools({
            name: 'boundary',
            args: {
              padding: 5,
              attrs: {
                stroke: 'red',
                strokeWidth: 2
              }
            }
          });
          return
        }
        else {
          if(this.Connector==true) {
            const container = this.$refs.ISMContainer
            const ports = container.querySelectorAll('.x6-port-body')
            showPorts(ports, true)
          }
        }
      })
      this.ISMCavasContainer.on('node:mouseleave', ({node}) => {
        node.removeTools();
        const container = this.$refs.ISMContainer
        const ports = container.querySelectorAll('.x6-port-body')
        showPorts(ports, false)
      })
      let startPoint = null;
      window.addEventListener('mousemove', this.drawMove)
      //===============================
      // 鼠标按下开始绘制
      this.ISMCavasContainer.on('blank:mousedown', ({ x, y }) => {
        if(!this.beginDrawLine)
        {
          return
        }
        this.ISMCavasContainer.disableSelection()
        startPoint = { x, y };
        if(this.drawingEdge==null) {
          this.drawingEdge = this.ISMCavasContainer.addEdge({
            source: startPoint,
            target: {x, y},
            position: {x, y},
            visible: true,
            size: {width: 100, height: 100},
            zIndex: 100,
            attrs: {
              line: {
                style: {
                  connection: true,
                  arrow: {size: 0},
                  strokeDasharray: "15 10",
                  stroke: "#13c2c2",
                  strokeWidth: 7,
                  strokeLinejoin: 'round',
                  rx: 6, // 圆角弧数
                  ry: 6, // 圆角弧数
                  animation: 'ant-line-forward 30s infinite linear',
                },
                sourceMarker: null,
                targetMarker: null,
              },
              wrap: {
                connection: true,
                stroke: "#D9D9D9",
                strokeWidth: 15,
                strokeLinejoin: 'round',
              }
            },
            router: 'normal',
            connector: "normal",
            data: {
              locked: false,
              UpdateNodeFlag: true,
              editMode: true,
              showDeviceUuid: "",
              IsToolBox: false,
              detail: {
                identifier: uuid.v1(),
                name: "线段",
                "type": "image",
                isCanvas: true,
                "action": [],
                "dataBind": [],
                "active": [
                  {
                    id: "Forward",
                    name: "component.ViewCanvasMoveLineArrow.Forward",
                    result: "",
                    isExpression: true,
                    condition: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                      operator: "",
                      OperatorValue: "",
                      OperatorMaxValue: "",
                    },
                  },
                  {
                    id: "Reverse",
                    name: "component.ViewCanvasMoveLineArrow.Reverse",
                    result: "",
                    isExpression: true,
                    condition: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                      operator: "",
                      OperatorValue: "",
                      OperatorMaxValue: "",
                    },
                  },
                ],
                "animate": {
                  "selected": [],
                  "condition": {
                    deviceSN: "",
                    selectVideoType: 0,
                    isBandDevice: false,
                    bandType: 1,
                    dataID: "",
                    dataName: "",
                    operator: "",
                    OperatorValue: "",
                    OperatorMaxValue: "",
                  },
                  "isExpression": false,
                  "animateList": [],
                  "animateElement": [
                    {
                      id: "blink",
                      elementList: [
                        {
                          "name": "component.public.animateSpeed",
                          "type": 7,
                          "value": 1,
                          "min": 0.1,
                          "key": "blinkSpeed",
                        },
                      ]
                    },
                    {
                      id: "millcolorGrad",
                      elementList: [
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
                          "name": "component.public.animateSpeed",
                          "type": 7,
                          "value": 1,
                          "min": 0.1,
                          "key": "animateSpeed",
                        },
                      ]
                    },
                    {
                      id: "animateSpin",
                      elementList: [
                        {
                          "name": "component.public.animateSpinSpeed",
                          "type": 7,
                          "value": 1,
                          "min": 0.1,
                          "key": "spinSpeed",
                        },
                        {
                          name: "configComponent.bigScreen.border.border89Direction",
                          type: 6,
                          value: 0,
                          enumList: [
                            {
                              value: 0,
                              option: "configComponent.bigScreen.border.border89DirectionForward"
                            },
                            {
                              value: 1,
                              option: "configComponent.bigScreen.border.border89DirectionNegative"
                            }
                          ],
                          min: 1,
                          key: "spinDirection",
                        }
                      ]
                    },
                  ],
                },
                "style": {
                  "position": {
                    "x": 0,
                    "y": 0,
                    "w": 600,
                    "h": 50
                  },
                  "points": [],
                  "visible": 1,
                  "zIndex": 100,
                  "transform": 0,
                  "backColor": "#D9D9D9",
                  foreColor: "#13c2c2",
                  borderWidth: 0,
                  BorderEdges: 0,
                  opacity: 1,
                  borderStyle: "solid",
                  borderColor: "#ccccff",
                  "diy": [
                    {
                      "name": "component.public.strokeLinejoin",
                      type: 6,
                      value: 1,
                      enumList: [
                        {
                          value: 0,
                          option: "component.public.strokeLineMiter"
                        },
                        {
                          value: 1,
                          option: "component.public.strokeLineRound"
                        },
                        {
                          value: 2,
                          option: "component.public.strokeLineBevel"
                        }
                      ],
                      "key": "strokeLinejoin",
                    },
                    {
                      "name": "component.public.strokeBgWidth",
                      "type": 1,
                      "value": 25,
                      "min": 1,
                      "key": "strokeBgWidth",
                    },
                    {
                      "name": "component.public.strokeWidth",
                      "type": 1,
                      "value": 10,
                      "min": 1,
                      "key": "strokeWidth",
                    },
                    {
                      "name": "component.public.strokeLength",
                      "type": 1,
                      "value": 15,
                      "min": 1,
                      "key": "strokeLength",
                    },
                    {
                      "name": "component.public.strokeSpace",
                      "type": 1,
                      "value": 10,
                      "min": 1,
                      "key": "strokeSpace",
                    },
                    {
                      "name": "component.public.animateSpeed",
                      "type": 1,
                      "value": 30,
                      "min": 1,
                      "key": "MoveBrokenLineInterval",
                    },
                    {
                      name: "displayConfig.ToolBox.Diagram.MoveBrokenLineConditionEnable",
                      type: 6,
                      value: 0,
                      enumList: [
                        {
                          value: 0,
                          option: "component.public.Forbidden"
                        },
                        {
                          value: 1,
                          option: "component.public.Enable"
                        }
                      ],
                      min: 1,
                      key: "MoveBrokenLineConditionEnable",
                    },
                    {
                      name: "configComponent.bigScreen.border.border89Direction",
                      type: 6,
                      value: 0,
                      enumList: [
                        {
                          value: 0,
                          option: "configComponent.bigScreen.border.border89DirectionForward"
                        },
                        {
                          value: 1,
                          option: "configComponent.bigScreen.border.border89DirectionNegative"
                        }
                      ],
                      min: 1,
                      key: "spinDirection",
                    }
                  ]
                }
              }
            },
          });
        }
        if(this.drawingEdge!=null) {
          this.drawingEdge.insertVertex({x, y});
        }
      });
      // 鼠标按下开始绘制
      this.ISMCavasContainer.on('cell:mousedown', ({ x, y }) => {
        if(!this.beginDrawLine)
        {
          return
        }
        this.ISMCavasContainer.disableSelection()
        startPoint = { x, y };
        if(this.drawingEdge==null) {
          this.drawingEdge = this.ISMCavasContainer.addEdge({
            source: startPoint,
            target: {x, y},
            position: {x, y},
            visible: true,
            size: {width: 100, height: 100},
            zIndex: 100,
            attrs: {
              line: {
                style: {
                  connection: true,
                  arrow: {size: 0},
                  strokeDasharray: "15 10",
                  stroke: "#13c2c2",
                  strokeWidth: 7,
                  strokeLinejoin: 'round',
                  rx: 6, // 圆角弧数
                  ry: 6, // 圆角弧数
                  animation: 'ant-line-forward 30s infinite linear',
                },
                sourceMarker: null,
                targetMarker: null,
              },
              wrap: {
                connection: true,
                stroke: "#D9D9D9",
                strokeWidth: 15,
                strokeLinejoin: 'round',
              }
            },
            router: 'normal',
            connector: "normal",
            data: {
              locked: false,
              UpdateNodeFlag: true,
              editMode: true,
              showDeviceUuid: "",
              IsToolBox: false,
              detail: {
                identifier: uuid.v1(),
                name: "线段",
                "type": "image",
                isCanvas: true,
                "action": [],
                "dataBind": [],
                "active": [
                  {
                    id: "Forward",
                    name: "component.ViewCanvasMoveLineArrow.Forward",
                    result: "",
                    isExpression: true,
                    condition: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                      operator: "",
                      OperatorValue: "",
                      OperatorMaxValue: "",
                    },
                  },
                  {
                    id: "Reverse",
                    name: "component.ViewCanvasMoveLineArrow.Reverse",
                    result: "",
                    isExpression: true,
                    condition: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                      operator: "",
                      OperatorValue: "",
                      OperatorMaxValue: "",
                    },
                  },
                ],
                "animate": {
                  "selected": [],
                  "condition": {
                    deviceSN: "",
                    selectVideoType: 0,
                    isBandDevice: false,
                    bandType: 1,
                    dataID: "",
                    dataName: "",
                    operator: "",
                    OperatorValue: "",
                    OperatorMaxValue: "",
                  },
                  "isExpression": false,
                  "animateList": [],
                  "animateElement": [
                    {
                      id: "blink",
                      elementList: [
                        {
                          "name": "component.public.animateSpeed",
                          "type": 7,
                          "value": 1,
                          "min": 0.1,
                          "key": "blinkSpeed",
                        },
                      ]
                    },
                    {
                      id: "millcolorGrad",
                      elementList: [
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
                          "name": "component.public.animateSpeed",
                          "type": 7,
                          "value": 1,
                          "min": 0.1,
                          "key": "animateSpeed",
                        },
                      ]
                    },
                    {
                      id: "animateSpin",
                      elementList: [
                        {
                          "name": "component.public.animateSpinSpeed",
                          "type": 7,
                          "value": 1,
                          "min": 0.1,
                          "key": "spinSpeed",
                        },
                        {
                          name: "configComponent.bigScreen.border.border89Direction",
                          type: 6,
                          value: 0,
                          enumList: [
                            {
                              value: 0,
                              option: "configComponent.bigScreen.border.border89DirectionForward"
                            },
                            {
                              value: 1,
                              option: "configComponent.bigScreen.border.border89DirectionNegative"
                            }
                          ],
                          min: 1,
                          key: "spinDirection",
                        }
                      ]
                    },
                  ],
                },
                "style": {
                  "position": {
                    "x": 0,
                    "y": 0,
                    "w": 600,
                    "h": 50
                  },
                  "points": [],
                  "visible": 1,
                  "zIndex": 100,
                  "transform": 0,
                  "backColor": "#D9D9D9",
                  foreColor: "#13c2c2",
                  borderWidth: 0,
                  BorderEdges: 0,
                  opacity: 1,
                  borderStyle: "solid",
                  borderColor: "#ccccff",
                  "diy": [
                    {
                      "name": "component.public.strokeLinejoin",
                      type: 6,
                      value: 1,
                      enumList: [
                        {
                          value: 0,
                          option: "component.public.strokeLineMiter"
                        },
                        {
                          value: 1,
                          option: "component.public.strokeLineRound"
                        },
                        {
                          value: 2,
                          option: "component.public.strokeLineBevel"
                        }
                      ],
                      "key": "strokeLinejoin",
                    },
                    {
                      "name": "component.public.strokeBgWidth",
                      "type": 1,
                      "value": 25,
                      "min": 1,
                      "key": "strokeBgWidth",
                    },
                    {
                      "name": "component.public.strokeWidth",
                      "type": 1,
                      "value": 10,
                      "min": 1,
                      "key": "strokeWidth",
                    },
                    {
                      "name": "component.public.strokeLength",
                      "type": 1,
                      "value": 15,
                      "min": 1,
                      "key": "strokeLength",
                    },
                    {
                      "name": "component.public.strokeSpace",
                      "type": 1,
                      "value": 10,
                      "min": 1,
                      "key": "strokeSpace",
                    },
                    {
                      "name": "component.public.animateSpeed",
                      "type": 1,
                      "value": 30,
                      "min": 1,
                      "key": "MoveBrokenLineInterval",
                    },
                    {
                      name: "displayConfig.ToolBox.Diagram.MoveBrokenLineConditionEnable",
                      type: 6,
                      value: 0,
                      enumList: [
                        {
                          value: 0,
                          option: "component.public.Forbidden"
                        },
                        {
                          value: 1,
                          option: "component.public.Enable"
                        }
                      ],
                      min: 1,
                      key: "MoveBrokenLineConditionEnable",
                    },
                    {
                      name: "configComponent.bigScreen.border.border89Direction",
                      type: 6,
                      value: 0,
                      enumList: [
                        {
                          value: 0,
                          option: "configComponent.bigScreen.border.border89DirectionForward"
                        },
                        {
                          value: 1,
                          option: "configComponent.bigScreen.border.border89DirectionNegative"
                        }
                      ],
                      min: 1,
                      key: "spinDirection",
                    }
                  ]
                }
              }
            },
          });
        }
        if(this.drawingEdge!=null) {
          this.drawingEdge.insertVertex({x, y});
        }
      });

// 监听顶点位置变化
      this.ISMCavasContainer.on('cell:vertex:mousemove', ({ cell }) => {
        console.log('顶点坐标更新:', cell.getVertices())
        // 可在此添加自定义逻辑（如限制移动范围等）
      })
      //=======================
      this.ISMCavasContainer.use( new Transform({
        resizing: {
          restrict: false,
          preserveAspectRatio: false,
          allowReverse: false,
          autoScroll:false,
          minWidth: 10,
          minHeight: 10,
          orthogonal:true,
          enabled(cellView) {
            return !cellView.getData()?.locked
          },
        },
        rotating: {
          enabled(cellView) {
            return !cellView.getData()?.locked&cellView.prop().shape!="view-ism-group-node"
          },
          grid: 1,
        },
      }))
      this.ISMCavasContainer.use(  new Snapline({
        enabled: true,
        clean: true,
        tolerance:20,
        resizing:true,
        className:"custom-snapline",
      }))
      this.ISMCavasContainer.use(   new Clipboard({
        enabled: true,
        useLocalStorage: true,
      }))
      this.ISMCavasContainer.use(  new Keyboard({
            enabled: true,
            global: true,
          }),
      )
      this.ISMCavasContainer.use(  new History({
            enabled: true,
            stackSize:20
          }),
      )
      this.ISMCavasContainer.use(   new Selection({
            enabled: true,
            multiple: true,
            rubberband: true,
            movable: true,
            strict:false,
            showNodeSelectionBox: true,
            showEdgeSelectionBox:true,
            filter: (cellView) => {
                 return !cellView.getData()?.locked
            },
          }),
      )
      this.ISMCavasContainer.use(new Export())
      // this.ISMCavasContainer.use(new MiniMap({
      //       container: this.$refs.refMiniMapContainer,
      //       width: 200,
      //       height: 160,
      //       padding: 10,
      //     }),
      // )
      //创建分组的父节点
      //小地图的显示和隐藏
      // this.ShowOrHideMinMap()
      this.ISMCavasContainer.on('cell:selected', ({ e, x, y, cell, view }) => {
        const viewShape = cell.prop().shape
        if(viewShape!="view-ism-group-node") {
          this.setCurrentSelectNode(cell)
        }
      })
      //检查是否启用辅助线
      if(_t.IsShowMakerLine) {
        _t.ISMCavasContainer.enableSnapline()
      }else{
        _t.ISMCavasContainer.disableSnapline()
      }
      this.ISMCavasContainer.on('scale', ( scale ) => {
        try {
          if (typeof this.guidesX.zoomTo !== 'undefined'&&typeof this.guides.zoomTo !== 'undefined') {
              this.guidesX.zoomTo(scale.sx)
              this.guides.zoomTo(scale.sy)
          }
        } catch (e) {

        }

        const width = _t.configData.layer.width *scale.sx
        const height =  _t.configData.layer.height*scale.sy;
        this.ISMCavasContainer.resize(width,height)
      });
      //页面空白处右键菜单
      this.ISMCavasContainer.on('blank:contextmenu', ({ e, x, y, cell, view }) => {
        if (this.drawingEdge&&this.drawingEdge!=null) {
          this.drawingEdge.insertVertex({ x, y });
          this.drawingEdge = null;
        }
        this.beginDrawLine = false
        this.onContextLayerMenu(e)
      })
      //节点菜单
      this.ISMCavasContainer.on('cell:contextmenu', ({ e, x, y, cell, view }) => {
        if (this.drawingEdge&&this.drawingEdge!=null) {
          this.drawingEdge.insertVertex({ x, y });
          this.drawingEdge = null;
        }
        this.beginDrawLine = false
        const viewShape = cell.prop().shape
        if(viewShape=="view-menu-nav")
        {
          this.ISMCavasContainer.unselect(cell)
        }
        else {
          if(viewShape!="view-ism-group-node") {
            this.setCurrentSelectNode(cell)
            this.ISMCavasContainer.select(cell)
          }
          this.componentRightClick(cell)
        }
      })
      this.ISMCavasContainer.on('blank:click', ({ e, x, y, cell, view }) => {
        this.setLayerSelected(true)
        if(this.beginDrawLine&&this.drawingEdge!=null)
        {
          this.drawingEdge.insertVertex({ x, y });
        }
      })
      this.ISMCavasContainer.on('cell:click', ({ e, x, y, cell, view }) => {
        if (this.isFormatPainterActive) {
          if(!cell.isEdge()) {
            this.applyNodeProps(this.selectedNode, cell)
          }
        }
        else {
          const viewShape = cell.prop().shape
          this.ClickUnSelectedComponent(cell)
          if (this.beginDrawLine && this.drawingEdge != null) {
            this.drawingEdge.insertVertex({x, y});
          } else {
            if (viewShape == "view-menu-nav") {
              this.ISMCavasContainer.unselect(cell)
            } else {
              this.setCurrentSelectNode(cell)
              this.ISMCavasContainer.select(cell)
            }
          }
        }
      })
      this.ISMCavasContainer.on('cell:dblclick', ({ e, x, y, cell, view }) => {
        if (this.drawingEdge&&this.drawingEdge!=null) {
          // 添加流动动画
          this.drawingEdge.insertVertex({ x, y });
          this.drawingEdge = null;
        }
        this.beginDrawLine = false
        this.ISMCavasContainer.enableSelection()
      })
      // 鼠标释放结束绘制
      this.ISMCavasContainer.on('blank:dblclick', ({ e, x, y }) => {
        if (this.drawingEdge&&this.drawingEdge!=null) {
          this.drawingEdge.insertVertex({ x, y });
          this.drawingEdge = null;
        }
        this.beginDrawLine = false
        this.ISMCavasContainer.enableSelection()
      });
      //绑定快捷键
      this.ISMCavasContainer.bindKey('ctrl+c', () => {
        const cells = this.ISMCavasContainer.getSelectedCells()
        if (cells.length) {
          this.ISMCavasContainer.copy(cells,{ deep: true,useLocalStorage:true })
        }
        return false
      })
      this.ISMCavasContainer.bindKey('esc', () => {
        this.drawingEdge = null;
        this.beginDrawLine = false
        this.SetFormatPainterState(false)
        return false
      })
      this.ISMCavasContainer.bindKey('ctrl+v', (e) => {
        if (!this.ISMCavasContainer.isClipboardEmpty()) {
          const cells = this.ISMCavasContainer.paste({ offset: 32,useLocalStorage: true })
          cells.forEach(cell => {
            let cellData = cell.prop().data || {}
            cellData.detail.identifier = uuid.v1()
            cellData.detail.name = cellData.detail.name+"_copy"
            cell.setData({
              locked:false,
              editMode: true,
              showDeviceUuid:"",
              IsToolBox:false,
              UpdateNodeFlag: new Date().getTime(),
              detail: cellData.detail
            }, {overwrite: true})
          })
          this.ISMCavasContainer.cleanSelection()
          this.ISMCavasContainer.select(cells)
        }
        return false
      })
      this.ISMCavasContainer.bindKey(['delete', 'backspace'], () => {
        this.DelNode()
        return false
      })
      this.ISMCavasContainer.bindKey('ctrl+z', () => {
        this.ISMCavasContainer.undo()
        return false
      })
      this.ISMCavasContainer.bindKey('ctrl+r', () => {
        this.ISMCavasContainer.redo()
        return false
      })
      // 左方向键
      this.ISMCavasContainer.bindKey(['arrowleft'], (e) => {
        e.preventDefault();
        const nodes = this.ISMCavasContainer.getSelectedCells();
        nodes.forEach(node => {
          node.translate(-1, 0); // x轴左移1像素
        });
      });

// 右方向键
      this.ISMCavasContainer.bindKey(['arrowright'], (e) => {
        e.preventDefault();
        const nodes = this.ISMCavasContainer.getSelectedCells();
        nodes.forEach(node => {
          node.translate(1, 0); // x轴右移1像素
        });
      });

// 上方向键
      this.ISMCavasContainer.bindKey(['arrowup'], (e) => {
        e.preventDefault();
        const nodes = this.ISMCavasContainer.getSelectedCells();
        nodes.forEach(node => {
          node.translate(0, -1); // y轴上移1像素
        });
      });

// 下方向键
      this.ISMCavasContainer.bindKey(['arrowdown'], (e) => {
        e.preventDefault();
        const nodes = this.ISMCavasContainer.getSelectedCells();
        nodes.forEach(node => {
          node.translate(0, 1); // y轴下移1像素
        });
      });

      //==============================
      //监听页面是否渲染完成，完成后就发广播给节点，告诉节点是在编辑模式
      this.ISMCavasContainer.on('render:done', () => {
        this.$EventBus.$emit('cell-editMode',{
          edit:true,
          toolbox:false
        })
        //由于节点内部无法访问vuex,只能通过广播把vuex传入节点
        this.$EventBus.$emit('cell-vuex',{
          PMapState:mapState,
          PMapActions:mapActions,
          PMapMutations:mapMutations,
          PStore:this.$store
        })
        this.setGroupList()
      });
      try{
        const components = JSON.parse(JSON.stringify(_t.configData.components))
        if (components.cells && Array.isArray(components.cells)) {
          components.cells = components.cells.filter(cell => cell && cell.shape)
        }
        _t.ISMCavasContainer.fromJSON(components)
      }catch (e){
        console.log(e)
      }
      this.dnd = initDnd(this.ISMCavasContainer, this.$refs.ISMContainer)
      this.setCurrentISMCavasContainer(this.ISMCavasContainer)
      this.setCurrentISMCavasDND(this.dnd)
      this.ISMCavasContainer.on('edge:dblclick', ({ cell }) => {
          this.setCurrentSelectNode(cell)
          this.ISMCavasContainer.select(cell)
      });
      this.ISMCavasContainer.on('edge:mouseenter', ({ cell }) => {
        if(this.beginDrawLine)
        {
          return
        }
        cell.addTools('vertices', 'onhover')
      })

      this.ISMCavasContainer.on('edge:mouseleave', ({ cell }) => {
        if(this.beginDrawLine)
        {
          return
        }
        if (cell.hasTools('onhover')) {
          cell.removeTools()
        }
      })
      this.ISMCavasContainer.on(['edge:connected'], ({ edge, currentCell }) => {
        edge.setProp({
          visible:true,
          position:{x:12, y:13},
          size:{width:100,height:100},
          connectType:"connected-edge",
          attrs: {
            line: {
              style: {
                stroke: "#A2B1C3",
                strokeWidth: 2,
              },
              sourceMarker:{
                name: "classic",
                width: 10,
                height:10
              } ,
              targetMarker: {
                name: "classic",
                width: 10,
                height:10
              },
            },
            wrap: {
              stroke: "#A2B1C3",
              strokeWidth: 2,
              strokeLinejoin: 'round',
            }
          },
          data: {
            locked:false,
            UpdateNodeFlag:true,
            editMode: true,
            showDeviceUuid:"",
            IsToolBox:false,
            detail:{
              identifier :uuid.v1(),
              name:"连接线段",
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
                  "w": 600,
                  "h": 50
                },
                "points": [],
                "visible":1,
                "zIndex": -1,
                "transform": 0,
                "backColor": "#A2B1C3",
                foreColor:"#A2B1C3",
                borderWidth:0,
                BorderEdges:0,
                opacity:1,
                borderStyle:"solid",
                borderColor:"#ccccff",
                "diy":[
                  {
                    "name":"component.public.strokeLineType",
                    type:6,
                    value:1,
                    enumList:[
                      {
                        value:0,
                        option:"component.public.strokeLineDashed"
                      },
                      {
                        value:1,
                        option:"component.public.strokeLineSolid"
                      }
                    ],
                    "key":"strokeLineType",
                  },
                  {
                    "name":"component.public.strokeLineMarker",
                    type:6,
                    value:1,
                    enumList:[
                      {
                        value:1,
                        option:"configComponent.ProgressBar.showInfoTrue"
                      },
                      {
                        value:0,
                        option:"configComponent.ProgressBar.showInfoFalse"
                      }
                    ],
                    "key":"strokeLineMarker",
                  },
                  {
                    "name":"component.public.strokeLineMarkerStyle",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"component.public.strokeLineMarkerClassic"
                      },
                      {
                        value:1,
                        option:"component.public.strokeLineMarkerDiamond"
                      },
                      {
                        value:2,
                        option:"component.public.strokeLineMarkerCross"
                      },
                      {
                        value:3,
                        option:"component.public.strokeLineMarkerCircle"
                      },
                      {
                        value:4,
                        option:"component.public.strokeLineMarkerCirclePlus"
                      },
                      {
                        value:5,
                        option:"component.public.strokeLineMarkerEllipse"
                      }
                    ],
                    "key":"strokeLineMarkerStyle",
                  },
                  {
                    "name":"component.public.strokeLineMarkerWidth",
                    "type":1,
                    "value":10,
                    "min":1,
                    "key":"strokeLineMarkerWidth",
                  },
                  {
                    "name":"component.public.strokeLineMarkerHeight",
                    "type":1,
                    "value":10,
                    "min":1,
                    "key":"strokeLineMarkerHeight",
                  },
                  {
                    "name": "component.public.strokeLineMarkerColor",
                    "type": 2,
                    "value": "#A2B1C3",
                    "key": "strokeLineMarkerColor",
                  },
                  {
                    "name":"component.public.strokeLinejoin",
                    type:6,
                    value:1,
                    enumList:[
                      {
                        value:0,
                        option:"component.public.strokeLineMiter"
                      },
                      {
                        value:1,
                        option:"component.public.strokeLineRound"
                      },
                      {
                        value:2,
                        option:"component.public.strokeLineBevel"
                      }
                    ],
                    "key":"strokeLinejoin",
                  },
                  {
                    "name":"component.public.strokeBgWidth",
                    "type":1,
                    "value":1,
                    "min":1,
                    "key":"strokeBgWidth",
                  },
                  {
                    "name":"component.public.strokeWidth",
                    "type":1,
                    "value":1,
                    "min":1,
                    "key":"strokeWidth",
                  },
                  {
                    "name":"component.public.strokeLength",
                    "type":1,
                    "value":1,
                    "min":1,
                    "key":"strokeLength",
                  },
                  {
                    "name":"component.public.strokeSpace",
                    "type":1,
                    "value":1,
                    "min":1,
                    "key":"strokeSpace",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":1,
                    "value":30,
                    "min":1,
                    "key":"MoveBrokenLineInterval",
                  },
                  {
                    name:"displayConfig.ToolBox.Diagram.MoveBrokenLineConditionEnable",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"component.public.Forbidden"
                      },
                      {
                        value:1,
                        option:"component.public.Enable"
                      }
                    ],
                    min:1,
                    key:"MoveBrokenLineConditionEnable",
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
              }
            }
          },
        })
      })

      this.ISMCavasContainer.on('node:mousemove', ({ node, options }) => {
        const selected = this.ISMCavasContainer.getSelectedCells();
        if (selected.length >1) {
          return
        }
        const currentPos = node.position();
        const delta = startMovePos ? {
          x: currentPos.x - startMovePos.x,
          y: currentPos.y - startMovePos.y
        } : { x: 0, y: 0 };
        startMovePos = currentPos;
        // 获取父节点及所有同级节点
        const parent = node.getParent();
        if(parent) {
          const ppos = parent.position()
          parent.position(ppos.x + delta.x, ppos.y + delta.y);
          const children = parent.getChildren()
          children.forEach(target => {
            if (target && target.id !== node.id) {
              const pos = target.position();
              target.position(pos.x + delta.x, pos.y + delta.y);
            }
          });
        }
      })

      this.ISMCavasContainer.on('cell:change:size', ({cell, current, previous}) => {
        const viewShape = cell.prop().shape
        if(viewShape=="view-ism-group-node") {
          const scaleX = current.width / previous.width
          const scaleY = current.height / previous.height
          cell.getChildren().forEach(child => {
            // 调整子节点位置
            const pos = child.getPosition({ relative: true })
            child.setPosition(pos.x * scaleX, pos.y * scaleY,{ relative: true })
            // child.resize(
            //     child.size().width * scaleX,
            //     child.size().height * scaleY
            // )
          })
        }
      });
      this.ISMCavasContainer.on('cell:added', ({cell, current, previous}) => {
        this.SyncLayerComponents()
        //由于节点内部无法访问vuex,只能通过广播把vuex传入节点
        this.$EventBus.$emit('cell-vuex',{
          PMapState:mapState,
          PMapActions:mapActions,
          PMapMutations:mapMutations,
          PStore:this.$store
        })
        cell.setVisible(true,{ silent: true });
        this.setGroupList()
        let dropData = cell.getData()
        if( typeof  dropData=="undefined")
        {
          if( cell.prop().shape=="edge") {
            cell.setProp({
              visible: true,
              position: {x: 12, y: 13},
              size: {width: 100, height: 100},
              connectType: "connected-edge",
              attrs: {
                line: {
                  style: {
                    stroke: "#A2B1C3",
                    strokeWidth: 2,
                  },
                  sourceMarker: {
                    name: "classic",
                    width: 10,
                    height: 10
                  },
                  targetMarker: {
                    name: "classic",
                    width: 10,
                    height: 10
                  },
                },
                wrap: {
                  stroke: "#A2B1C3",
                  strokeWidth: 2,
                  strokeLinejoin: 'round',
                }
              },
              data: {
                locked: false,
                UpdateNodeFlag: true,
                editMode: true,
                showDeviceUuid: "",
                IsToolBox: false,
                detail: {
                  identifier: uuid.v1(),
                  name: "连接线段",
                  "type": "image",
                  isCanvas: true,
                  "action": [],
                  "dataBind": [],
                  "active": [
                    {
                      id: "Forward",
                      name: "component.ViewCanvasMoveLineArrow.Forward",
                      result: "",
                      isExpression: true,
                      condition: {
                        deviceSN: "",
                        selectVideoType: 0,
                        isBandDevice: false,
                        bandType: 1,
                        dataID: "",
                        dataName: "",
                        operator: "",
                        OperatorValue: "",
                        OperatorMaxValue: "",
                      },
                    },
                    {
                      id: "Reverse",
                      name: "component.ViewCanvasMoveLineArrow.Reverse",
                      result: "",
                      isExpression: true,
                      condition: {
                        deviceSN: "",
                        selectVideoType: 0,
                        isBandDevice: false,
                        bandType: 1,
                        dataID: "",
                        dataName: "",
                        operator: "",
                        OperatorValue: "",
                        OperatorMaxValue: "",
                      },
                    },
                  ],
                  "animate": {
                    "selected": [],
                    "condition": {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                      operator: "",
                      OperatorValue: "",
                      OperatorMaxValue: "",
                    },
                    "isExpression": false,
                    "animateList": [],
                    "animateElement": [
                      {
                        id: "blink",
                        elementList: [
                          {
                            "name": "component.public.animateSpeed",
                            "type": 7,
                            "value": 1,
                            "min": 0.1,
                            "key": "blinkSpeed",
                          },
                        ]
                      },
                      {
                        id: "millcolorGrad",
                        elementList: [
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
                            "name": "component.public.animateSpeed",
                            "type": 7,
                            "value": 1,
                            "min": 0.1,
                            "key": "animateSpeed",
                          },
                        ]
                      },
                      {
                        id: "animateSpin",
                        elementList: [
                          {
                            "name": "component.public.animateSpinSpeed",
                            "type": 7,
                            "value": 1,
                            "min": 0.1,
                            "key": "spinSpeed",
                          },
                          {
                            name: "configComponent.bigScreen.border.border89Direction",
                            type: 6,
                            value: 0,
                            enumList: [
                              {
                                value: 0,
                                option: "configComponent.bigScreen.border.border89DirectionForward"
                              },
                              {
                                value: 1,
                                option: "configComponent.bigScreen.border.border89DirectionNegative"
                              }
                            ],
                            min: 1,
                            key: "spinDirection",
                          }
                        ]
                      },
                    ],
                  },
                  "style": {
                    "position": {
                      "x": 0,
                      "y": 0,
                      "w": 600,
                      "h": 50
                    },
                    "points": [],
                    "visible": 1,
                    "zIndex": -1,
                    "transform": 0,
                    "backColor": "#A2B1C3",
                    foreColor: "#A2B1C3",
                    borderWidth: 0,
                    BorderEdges: 0,
                    opacity: 1,
                    borderStyle: "solid",
                    borderColor: "#ccccff",
                    "diy": [
                      {
                        "name": "component.public.strokeLineType",
                        type: 6,
                        value: 1,
                        enumList: [
                          {
                            value: 0,
                            option: "component.public.strokeLineDashed"
                          },
                          {
                            value: 1,
                            option: "component.public.strokeLineSolid"
                          }
                        ],
                        "key": "strokeLineType",
                      },
                      {
                        "name": "component.public.strokeLineMarker",
                        type: 6,
                        value: 1,
                        enumList: [
                          {
                            value: 1,
                            option: "configComponent.ProgressBar.showInfoTrue"
                          },
                          {
                            value: 0,
                            option: "configComponent.ProgressBar.showInfoFalse"
                          }
                        ],
                        "key": "strokeLineMarker",
                      },
                      {
                        "name": "component.public.strokeLineMarkerStyle",
                        type: 6,
                        value: 0,
                        enumList: [
                          {
                            value: 0,
                            option: "component.public.strokeLineMarkerClassic"
                          },
                          {
                            value: 1,
                            option: "component.public.strokeLineMarkerDiamond"
                          },
                          {
                            value: 2,
                            option: "component.public.strokeLineMarkerCross"
                          },
                          {
                            value: 3,
                            option: "component.public.strokeLineMarkerCircle"
                          },
                          {
                            value: 4,
                            option: "component.public.strokeLineMarkerCirclePlus"
                          },
                          {
                            value: 5,
                            option: "component.public.strokeLineMarkerEllipse"
                          }
                        ],
                        "key": "strokeLineMarkerStyle",
                      },
                      {
                        "name": "component.public.strokeLineMarkerWidth",
                        "type": 1,
                        "value": 10,
                        "min": 1,
                        "key": "strokeLineMarkerWidth",
                      },
                      {
                        "name": "component.public.strokeLineMarkerHeight",
                        "type": 1,
                        "value": 10,
                        "min": 1,
                        "key": "strokeLineMarkerHeight",
                      },
                      {
                        "name": "component.public.strokeLineMarkerColor",
                        "type": 2,
                        "value": "#A2B1C3",
                        "key": "strokeLineMarkerColor",
                      },
                      {
                        "name": "component.public.strokeLinejoin",
                        type: 6,
                        value: 1,
                        enumList: [
                          {
                            value: 0,
                            option: "component.public.strokeLineMiter"
                          },
                          {
                            value: 1,
                            option: "component.public.strokeLineRound"
                          },
                          {
                            value: 2,
                            option: "component.public.strokeLineBevel"
                          }
                        ],
                        "key": "strokeLinejoin",
                      },
                      {
                        "name": "component.public.strokeBgWidth",
                        "type": 1,
                        "value": 1,
                        "min": 1,
                        "key": "strokeBgWidth",
                      },
                      {
                        "name": "component.public.strokeWidth",
                        "type": 1,
                        "value": 1,
                        "min": 1,
                        "key": "strokeWidth",
                      },
                      {
                        "name": "component.public.strokeLength",
                        "type": 1,
                        "value": 1,
                        "min": 1,
                        "key": "strokeLength",
                      },
                      {
                        "name": "component.public.strokeSpace",
                        "type": 1,
                        "value": 1,
                        "min": 1,
                        "key": "strokeSpace",
                      },
                      {
                        "name": "component.public.animateSpeed",
                        "type": 1,
                        "value": 30,
                        "min": 1,
                        "key": "MoveBrokenLineInterval",
                      },
                      {
                        name: "displayConfig.ToolBox.Diagram.MoveBrokenLineConditionEnable",
                        type: 6,
                        value: 0,
                        enumList: [
                          {
                            value: 0,
                            option: "component.public.Forbidden"
                          },
                          {
                            value: 1,
                            option: "component.public.Enable"
                          }
                        ],
                        min: 1,
                        key: "MoveBrokenLineConditionEnable",
                      },
                      {
                        name: "configComponent.bigScreen.border.border89Direction",
                        type: 6,
                        value: 0,
                        enumList: [
                          {
                            value: 0,
                            option: "configComponent.bigScreen.border.border89DirectionForward"
                          },
                          {
                            value: 1,
                            option: "configComponent.bigScreen.border.border89DirectionNegative"
                          }
                        ],
                        min: 1,
                        key: "spinDirection",
                      }
                    ]
                  }
                }
              },
            })
          }
        }
        else {
          cell.setZIndex(parseInt(dropData.detail.style.zIndex))
        }
      });
      this.ISMCavasContainer.on('cell:removed', ({cell, current, previous}) => {
        this.SyncLayerComponents()
        this.setGroupList()
      });
    }
  },
  mounted() {
    let _t = this
    let uid = this.$route.params.uid
    _t.spinning = true
    this.getLayerDataStruct({pageType:this.isMobile,uuid:uid,cb:function (){
        _t.initCavasContainer()
        _t.spinning = false
      }});
    _t.$EventBus.$on("ClickComponentsEvent", (id) => {
      const cell = _t.ISMCavasContainer.getCellById(id)
      _t.ISMCavasContainer.select(cell)
    })
  },
  beforeDestroy() {
    // 清理 DOM 事件
    window.removeEventListener('mousemove', this.drawMove)
    // 清理 EventBus
    // 清理 Guides 标尺
    if (this.guidesX) { this.guidesX.destroy(); this.guidesX = null }
    if (this.guides) { this.guides.destroy(); this.guides = null }
    // 清理 DnD
    if (this.dnd) { this.dnd = null }
    // 清理 X6 Graph
    if (this.ISMCavasContainer) {
      this.ISMCavasContainer.clearCells()
      this.ISMCavasContainer.off()
      this.ISMCavasContainer.dispose()
      this.ISMCavasContainer = null
    }
  }
}
</script>
<style>
.ant-spin-container{
  height:100%;
}
.normal {
  cursor: default;
}


.custom-snapline {
  stroke: #00FF00;
  stroke-width: 10px;
}
@keyframes ant-line-forward {
  to { stroke-dashoffset: -1000; }
}
@keyframes ant-line-inverse {
  to { stroke-dashoffset: 1000; }
}
</style>
<style  lang="less" scoped>
::v-deep .format-painter {
  cursor: crosshair ;
}
::v-deep  .x6-widget-transform > div:hover {
  background-color: #13e1eb;
}
::v-deep .x6-widget-transform-active-handle {
  background-color: #13e1eb;
}
::v-deep .x6-widget-transform-resize {
  border-radius: 0;
}
::v-deep .x6-widget-selection-inner {
  border: 2px solid #13e1eb;
}
::v-deep .x6-widget-transform {
  border: 0px dashed #13e1eb;
}
::v-deep  .x6-widget-selection-box {
  border: 2px dashed #13e1eb;
}
::v-deep .x6-widget-transform > div {
  background-color: #fff;
  border: 1px solid #13e1eb;
  transition: background-color 0.6s;
}

.minimap-container {
  position: absolute;
  right: 20px;
  bottom: 20px;
  cursor: move;
  z-index: 100;
  box-shadow: 0 0 10px rgba(0,0,0,0.2);
  border-radius: 4px;
}
.x6-container {
  overflow:scroll;
  position: relative;
  height: 100%;
  width: 100%;
}
.ruler-horizontal {
  position: absolute;
  top: 0;
  width:var(--rulerWidth);
  height: 20px;
  left: 20px;
}
.ruler-vertical {
  position: absolute;
  top: 20px;
  left: 0;
  height:var(--rulerHeight);
  width:20px;
}
.graph-container {
  position: absolute;
  height: 100%;
  width: 100%;
  margin-top: 20px;
  margin-left: 20px;
}
::-webkit-scrollbar {
  /*滚动条整体样式*/
  width : 3px;  /*高宽分别对应横竖滚动条的尺寸*/
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
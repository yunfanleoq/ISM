export default {
  LayerData: {
    name: '--',
    layer: {
      backColor: '',
      backgroundImage: '',
      widthHeightRatio: '',
      width: 300,
      height:600,
      autoSize:0,
      Padding:1
    },
    components: []
  }, //页面容器的页面数据
  LayerContainerData: {
    name: '--',
    layer: {
      backColor: '',
      backgroundImage: '',
      widthHeightRatio: '',
      width: 300,
      height:600,
      autoSize:0,
      Padding:1
    },
    components: []
  }, //当前场景的组态数据
  PCPageList:[],
  PhonePageList:[],
  selectPageUuid:"",
  selectPageContainerUuid:"",
  toolBoxList:[],
  PageCanVasList:[],
  MesComponentsList:[],
  DiyComponentsList:[],
  loggerList:[],
  loggerIndex:1,
  selectedValue:100,
  PopUpContainerConfigData:{//弹窗界面
    name: '--',
    layer: {
      backColor: '',
      backgroundImage: '',
      widthHeightRatio: '',
      width: 300,
      height:600,
      Padding:1
    },
    components: []
  },
  PopUpConfigData:{//弹窗界面
    name: '--',
    layer: {
      backColor: '',
      backgroundImage: '',
      widthHeightRatio: '',
      width: 300,
      height:600,
      Padding:1
    },
    components: []
  },
  selectedIsLayer: true, //当前选择的是不是layer层
  selectedComponent: null,//当前选择的单个组件--仅仅当只有一个组件选中有效，当有多个组件选中，则置为null
  selectedComponents: [], //当前选择的组件--只存identifier
  selectedComponentMap: {}, //当前选择的组件--key=identifier，本数据和selectedComponents同步，主要用于渲染判断
  copySrcItems: [],//当前是否使用了CTRL+C命令
  formatSrcItems:{},
  isFormat:false,
  copyCount: 0,//copy计数，对于同一个复制源，每次复制后计数+1
  undoStack: [],//
  isFormatPainterActive:false,
  FormatPainterCell:null,
  Isometric_row:"50",
  GroupList:[],
  Isometric_colu:"50",
  redoStack: [],//
  prePageUuid:"",
  isHistoryOp:false,
  HistoryOpTimer:null,
  ProtectedID:"",
  ISMCavasContainer:null,
  selectedComponentDiyData:null,
  selectedNodePops:null,
  ISMCavasDND:null,
  selectedNode:null,
  curPageUuid:"",
  UnSelectedComponent:null,
  isLocked: false,
  // 存储用户真实密码（实际项目需从接口获取或加密存储，此处仅为示例）
  realPassword: '123456' // 实际项目需替换为用户登录后的加密密码
}

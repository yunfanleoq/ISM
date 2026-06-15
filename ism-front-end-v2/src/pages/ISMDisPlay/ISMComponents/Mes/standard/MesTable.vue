<template>
  <div class="view-url">
    <div :ref="detail.identifier" :style="{'z-index':!editMode&&!IsToolBox?5:-10000,'position':'relative'}">
      基本信息
    </div>
  </div>
</template>
<script>
import BaseView from '../../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-mes-table',
  inject: ['getNode'],
  props: {

  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        if(this.editMode) {
          this.initComponents(newVal);
        }
      },
      deep: true
    }
  },
  components: {},
  data() {
    return {
      detail:{},
      IsToolBox:false,
      editMode:true,
      Url:"",
      width:600,
      height:600,
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      fillOpacity:1,
      strokeOpacity:1,
      animateType:"blink",
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      amisScoped:null,
      base:{
        "text": "configComponent.Mes.standard.edittable.title",
        "icon": "icon-kebianjibiaoge",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 440,
              "h": 330
            },
            "visible":1,
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"component.public.AMISConfig",
                "type":11,
                "value":"{\n" +
                    "  \"type\": \"page\",\n" +
                    "  \"body\": {\n" +
                    "    \"type\": \"service\",\n" +
                    "    \"api\": \"\",\n" +
                    "    \"body\": [\n" +
                    "      {\n" +
                    "        \"type\": \"table\",\n" +
                    "        \"title\": \"表格1\",\n" +
                    "        \"source\": \"$rows\",\n" +
                    "        \"columns\": [\n" +
                    "          {\n" +
                    "            \"name\": \"engine\",\n" +
                    "            \"label\": \"Engine\"\n" +
                    "          },\n" +
                    "          {\n" +
                    "            \"name\": \"version\",\n" +
                    "            \"label\": \"Version\"\n" +
                    "          }\n" +
                    "        ]\n" +
                    "      }\n" +
                    "    ]\n" +
                    "  }\n" +
                    "}",
                "rows":30,
                "key":"AMISConfig",
              }
            ]
          }
        }
      }
    }
  },
  methods: {
    initComponents(option){
      this.width = option.style.position.w
      this.height = option.style.position.h
      let i=0
      for(i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="AMISConfig")
        {
          this.AMISConfig=option.style.diy[i].value
        }
      }
      let refObj = this.detail.identifier
      let view = this.$refs[refObj]
      let AMISConfig = JSON.parse(this.AMISConfig)
      const amis = amisRequire('amis/embed')
      this.amisScoped = amis.embed(view,AMISConfig)
    }
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {


        })
        _t.$EventBus.$on(animateEvent, (data) => {
          _t.isStart = data
        })

    });
  },
  created() {
    let _t = this
    const node = this.getNode()
    node.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    node.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
    });
    this.detail = node.getData().detail
    this.editMode = node.getData().editMode
    this.showDeviceUuid = node.getData().showDeviceUuid
    this.IsToolBox = node.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
    })
  }
}
</script>

<style lang="less">
.view-url {
  position:absolute;
  z-index:1000;
}
</style>

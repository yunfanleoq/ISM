<template>
  <div>
    <div :id="detail.identifier" :style="{'width':detail.style.position.w+'px','height':detail.style.position.h+'px'}"></div>
    <system-image-model @onSelectImage="onSelectImage" :networkImageUrl="imageComponentData" ref="systemImageModel"></system-image-model>

    <a-modal
        :title="$t('configComponent.Map.IndoorMap.editTitle')"
        width="700px"
        height="600px"
        :footer="null"
        v-drag-modal
        :visible="attrVisible"
        @cancel="attrVisible=false"
    >
      <div>
        <a-tabs default-active-key="1">
        <a-tab-pane key="1" :tab="$t('configComponent.Map.IndoorMap.makerAttr')" >
          <div style="padding: 10px;height: auto">
            <a-form-item :label="$t('configComponent.Map.IndoorMap.editTitle')">
              <vue-hover-mask>
                <!-- 默认插槽 -->
                <img v-if="clickData.icon.url!=''"  :style="{width: clickData.icon.width+'px',height:clickData.icon.height+'px'}" :src="clickData.icon.url" />
                <div v-else :style="{width: '100px',height:'100px',cursor: 'pointer','background-color':'#F2F2F2'}"></div>
                <!-- action插槽 -->
                <template v-slot:action>
                  <span style="font-size: 14px" @click="showSystemImageModel">{{$t('component.systemImageModel.selectImage')}}</span>
                  <a-divider type="vertical" />
                  <span  style="font-size: 14px" @click="clickData.icon.url=''">{{$t('component.systemImageModel.delImage')}}</span>
                </template>
              </vue-hover-mask>
            </a-form-item>
            <a-form layout="inline">
              <a-form-item :label="$t('configComponent.Map.IndoorMap.IcoWidth')">
                <a-input

                    v-model="clickData.icon.width"
                />
              </a-form-item>
              <a-form-item :label="$t('configComponent.Map.IndoorMap.IcoHeight')">
                <a-input

                    v-model="clickData.icon.height"
                />
              </a-form-item>
            </a-form>
            <a-form layout="inline">
              <a-form-item :label="$t('configComponent.Map.IndoorMap.IcoOffsetX')">
                <a-input

                    v-model="clickData.icon.iconAnchorX"
                />
              </a-form-item>
              <a-form-item :label="$t('configComponent.Map.IndoorMap.IcoOffsetY')">
                <a-input

                    v-model="clickData.icon.iconAnchorY"
                />
              </a-form-item>
            </a-form>
          </div>
        </a-tab-pane>

        <a-tab-pane key="2" :tab="$t('configComponent.Map.IndoorMap.makerEvent')" >
          <div style="padding: 10px;height: auto">
            <a-form  :label-col="{ span: 7}" :wrapper-col="{ span: 16 }">
              <a-form-item :label="$t('displayConfig.Properties.linkType')">
                <a-select
                    v-model="clickData.link.linkType"
                    allowClear
                >
                  <a-select-option v-for="options in [{label:'displayConfig.Properties.linkInside',value:'Inside'},{label:'displayConfig.Properties.linkExternal',value:'External'}]" :key="options.value" :value="options.value">
                    {{  $t(options.label)}}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <div v-if="clickData.link.linkType=='Inside'">
                <a-form-item :label="$t('displayConfig.Properties.linkIAppUUID')">
                  <a-select
                      v-model="clickData.link.Inside.displayUUID"
                      allowClear
                  >
                    <a-select-option v-for="options in configurationModel" :key="options.uuid" :value="options.uuid">
                      {{ options.name}}
                    </a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.linkIAppPageUUID')">
                  <a-select
                      v-model="clickData.link.Inside.pageUUID"
                  >
                    <a-select-option v-for="options in generateTargetPage(clickData.link.Inside.displayUUID)" :key="options.value" :value="options.value">
                      {{ options.label}}
                    </a-select-option>
                  </a-select>
                </a-form-item>
              </div>
              <div v-else>
                <a-form-item :label="$t('displayConfig.Properties.linkExternalUrl')">

                  <a-input type="text" v-model="clickData.link.External" >

                  </a-input>
                </a-form-item>

                <a-form-item :label="$t('displayConfig.Properties.OpenLinkExternalType')">

                  <a-select v-model="clickData.link.OpenExternalType">
                    <a-select-option value="self">{{$t('displayConfig.Properties.OpenLinkExternalSelf')}}</a-select-option>
                    <a-select-option value="new">{{$t('displayConfig.Properties.OpenLinkExternalNew')}}</a-select-option>
                  </a-select>

                </a-form-item>
              </div>
              <a-form-item :label="$t('displayConfig.Properties.isLinkPopUp')">
                <a-checkbox :checked="clickData.link.isPopUp"  v-model="clickData.link.isPopUp">
                  {{ $t('displayConfig.Properties.isLinkPopUp') }}
                </a-checkbox>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.autoClose')" v-if="clickData.link.isPopUp">
                <a-checkbox :checked="clickData.link.autoClose"  v-model="clickData.link.autoClose">
                  {{ $t('displayConfig.Properties.autoClose') }}
                </a-checkbox>
              </a-form-item>
              <div v-if="clickData.link.isPopUp&&clickData.link.linkType!=='Inside'">
                <a-form-item :label="$t('displayConfig.Properties.linkExternalWidth')">

                  <a-input type="text" v-model="clickData.link.width" >

                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.linkExternalHeight')">

                  <a-input type="text" v-model="clickData.link.height" >

                  </a-input>
                </a-form-item>

                <a-form-item :label="$t('displayConfig.Properties.linkExternalTitle')">

                  <a-input type="text" v-model="clickData.link.title" >

                  </a-input>
                </a-form-item>
              </div>
            </a-form>
          </div>
        </a-tab-pane>

      </a-tabs>
      </div>
    </a-modal>

  </div>
</template>

<script>
import svgView from '../View';
import 'leaflet/dist/leaflet.css'
import 'leaflet/dist/leaflet'
import 'leaflet/dist/leaflet-src'
import 'leaflet/dist/leaflet-src.esm'
import * as L from 'leaflet'
import icon from 'leaflet/dist/images/marker-icon.png'
import iconShadow from 'leaflet/dist/images/marker-shadow.png'
import 'leaflet.pm'
import 'leaflet.pm/dist/leaflet.pm.css'
import systemImageModel from "@/components/systemImageModel/systemImageModel";

const DefaultIcon = L.icon({
  iconUrl: icon,
  shadowUrl: iconShadow,
})
L.Marker.prototype.options.icon = DefaultIcon
import VueHoverMask from "@/components/VueHoverMask/VueHoverMask"
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-map-indoor',
    inject: ['getNode'],
    i18n: require('@/i18n/language'),
    components: {
      systemImageModel,
      VueHoverMask
    },
    data() {
      return {
        detail:{},
        IsToolBox:false,
        editMode:true,
        pic: 0,
        displayPageList:new Map,
        mapUpdate:false,
        map: null,
        fillColor: '#2d75ff80',
        currMaker:null,
        bounds: [
          [0, 0],
          [0, 0],
        ], // 平面图大小
        w:23,
        h:25,
        imageComponentData:"",
        configurationModel:[],
        attrVisible:false,
        clickData:{
          icon:{
            url:"",
            width:50,
            height:50,
          },
          link:{
            linkType:"",
            isPopUp:false,
            width:50,
            height:50,
            title:"",
            autoClose:false,
            External:"",
            OpenExternalType:"new",
            Inside:{
              displayUUID:"",
              pageUUID:"",
            }
          }
        },
        DefaultZoom:13,
        MapPng:"",
        SaveData:"",
        MinZoom:0,
        MaxZoom:10,
        base:{
          text: "configComponent.Map.IndoorMap.title",
          "icon": "icon-map_marker",
          "isFontIcon": true,
          "info": {
            "type": "image",
            "action": [],
            "dataBind":[],
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
                {
                  id: "blink",
                  name: "component.public.animateBlink",
                },
                {
                  id: "Zoom",
                  name: "component.public.Zoom",
                },
                {
                  id: "animateSpin",
                  name: "component.public.animateSpin",
                },
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
                "w": 800,
                "h": 560
              },
              "visible":1,
              "zIndex": -1,
              "transform": 0,
              "diy":[
                {
                  "name":"configComponent.Map.IndoorMap.SaveData",
                  "type":15,
                  "value":"",
                  "key":"SaveData",
                },
                {
                  "name":"configComponent.Map.IndoorMap.MapPng",
                  "type":5,
                  "value":"/static/11.png",
                  "key":"MapPng",
                },
                {
                  "name":"configComponent.Map.IndoorMap.DefaultZoom",
                  "type":1,
                  "value":13,
                  "key":"DefaultZoom",
                },
                {
                  "name":"configComponent.Map.IndoorMap.MinZoom",
                  "type":1,
                  "value":0,
                  "key":"MinZoom",
                },
                {
                  "name":"configComponent.Map.IndoorMap.MaxZoom",
                  "type":1,
                  "value":10,
                  "key":"MaxZoom",
                },
              ]
            }
          }
        }
      }
    },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          if(this.editMode){
            if(newVal.style.diy[1].value!=this.MapPng)
            {
              this.initComponents(newVal);
            }
            else if(newVal.style.diy[2].value!=this.DefaultZoom)
            {
              this.initComponents(newVal);
            }
            else if(newVal.style.position.h!=this.h||newVal.style.position.w!=this.w)
            {
              this.initComponents(newVal);
            }
            this.detail.style.diy[0].value = this.SaveData
          }
        },
        immediate: true,
        deep: true
      },
      SaveData:{
        handler(newVal, oldVal) {
          if(!this.editMode){
            return
          }
          this.detail.style.diy[0].value = newVal
        },
        deep: true
      },
      clickData:{
        handler(newVal, oldVal) {
          if(!this.editMode){
            return
          }
          if(newVal.icon.url!="") {
            let greenIcon = L.icon({
              iconUrl: newVal.icon.url,
              iconSize: [newVal.icon.width, newVal.icon.height],  //图标宽高
              iconAnchor:[newVal.icon.iconAnchorX,newVal.icon.iconAnchorY]
            })
            this.currMaker.setIcon(greenIcon)
          }

          let res = JSON.parse(this.SaveData) || []
          // 没有值没有绘制,赋个值初始值
          if (!res)
          {
            return
          }
          for(let i=0;i<res.length;i++)
          {
            if((res[i]._leaflet_id!="undefined")&&(this.currMaker._leaflet_id==res[i]._leaflet_id)){
              res[i].icon = this.clickData.icon
              res[i].link = this.clickData.link
              break
            }
          }
          this.SaveData = JSON.stringify(res)
        },
        deep: true
      },
    },
    methods: {
      generateTargetPage (uuid) {
        return this.displayPageList.get(uuid)
      },
      GetDisplayPage(uuid){
        let params={
          muid:uuid
        }
        let _t = this
        getDisplayModelLayerData(params).then(function (res){
          if(res.data.code==0)
          {
            let pageLayer = res.data.layer
            if(pageLayer.length>0)
            {
              let displayArray=[]
              for(let i=0;i<pageLayer.length;i++)
              {
                if (pageLayer[i].IsLogin==1)
                {
                  continue
                }
                let pageInfo = {}
                pageInfo.label = pageLayer[i].PageName
                pageInfo.value = pageLayer[i].PageId
                pageInfo.pageType = pageLayer[i].PageType
                pageInfo.pageModelUuid = pageLayer[i].modelId
                displayArray.push(pageInfo)
              }
              _t.displayPageList.set(uuid,displayArray)
            }
          }
        })
      },
      getConfigurationModel(){
        this.configurationModel=[]
        let _t = this
        const  params= {
          DisplayType: 1
        }
        displayModelList(params).then(function (res){
          let tableData={}
          if(res.data.list!=null)
          {
            for(let i=0;i<res.data.list.length;i++)
            {
              tableData.name = res.data.list[i].name
              tableData.description = res.data.list[i].description
              tableData.uuid = res.data.list[i].displayUid
              _t.configurationModel.push(tableData)
              tableData={}
              _t.GetDisplayPage(res.data.list[i].displayUid)
            }
          }

        })
      },
      showModal() {
        this.attrVisible = true;
      },
      handleCancel(e) {
        this.attrVisible = false;
      },
      showSystemImageModel(){
        this.$refs.systemImageModel.showModal(0)
      },
      onSelectImage(url){
        this.clickData.icon.url = url
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        let i=0
        this.w = option.style.position.w
        this.h = option.style.position.h
        for(i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="SaveData")
          {
            this.SaveData=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MapPng")
          {
            this.MapPng=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="DefaultZoom")
          {
            this.DefaultZoom=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MinZoom")
          {
            this.MinZoom=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="MaxZoom")
          {
            this.MaxZoom=parseInt(option.style.diy[i].value)
          }
        }

        this.initMap()
        this.initDate()
        this.getConfigurationModel()
        // 窗口缩放,地图自适应缩放
        window.onresize = () => {
          this.map.invalidateSize(true)
          this.map.fitBounds(this.bounds)
        }
      },
      // 使用id为vue-leaflet的div容器初始化地图
      initMap() {
        if(this.map) {
          this.map.remove();
        }
        this.map = L.map(this.detail.identifier, {
          minZoom: this.MinZoom,
          zoom: this.DefaultZoom,
          maxZoom: this.MaxZoom,
          zoomSnap: 0.1, // 缩放步长
          attributionControl: false, // 右下角图例控件
          zoomControl: this.editMode, // 缩放控件
          crs: L.CRS.Simple, // 坐标系
          center: [0, 0], // 中心点
        })
        this.map.pm.setLang('zh') // 控件提示设置中文
        if(this.editMode) {
          this.map.pm.addControls({
            position: 'topleft', // 控件菜单位置
            drawPolygon: false, //绘制多边形
            drawMarker: true, //绘制标记点
            drawCircleMarker: false, //绘制圆形标记
            drawPolyline: false, //绘制线条
            drawRectangle: false, //绘制矩形
            drawCircle: false, //绘制圆圈
            editMode: false, //编辑多边形
            dragMode: false, //拖动多边形
            cutPolygon: false, //添加⼀个按钮以删除多边形⾥⾯的部分内容
            removalMode: true, //清除多边形
          })
        }
        // 全局图层样式
        this.map.pm.setPathOptions({
          fillOpacity: 0.5,
          fillColor: this.fillColor,
        })
        this.map.on('drag', this.mapDrag);
      },
      // 地图初始化
      initDate() {
        // 获取图片宽高,防止不同尺寸的图片回显在页面上有多余的留白问题
        let img = new Image()
        let that = this
        let url = this.MapPng
        img.src = url
        img.onload = () => {
          that.bounds = [
            [0, 0],
            [that.detail.style.position.h, that.detail.style.position.w],
          ]
          that.map.setMaxBounds(that.bounds)
          // 创建地图
          L.imageOverlay(url, that.bounds).addTo(that.map)
          // 设置地图图层区域
          that.map.fitBounds(that.bounds)
        }
        let res = []
        try{
          res = JSON.parse(this.SaveData)
        }catch (e){
          res=[]
          this.SaveData="[]"
        }
        // 获取本地存储的图层数据
        // 当前平面图下如果有绘制的图层,就进行回显
        if (res) {
          // 循环绘制图层
          for(let i=0;i<res.length;i++)
          {
            let e = res[i]
            if(e.type=="layer"&&typeof e.latlngs!="undefined")
            {
              let lay = L.polygon(e.latlngs).addTo(this.map).on('pm:edit', this.editLayClick).on('click', this.layClick).on('pm:dragend', this.dragendLayClick)
              // 这个id每次刷新,来回切换平面图,都会改变,所以要更新一下,如果不更新,下面删除,拖拽,编辑,的时候就找不到图层了
              res[i]._leaflet_id = lay._leaflet_id
            }
            else if(e.type=="maker")
            {
              if(typeof e.icon!="undefined"&&e.icon.url!="") {
                let greenIcon = L.icon({
                  iconUrl: e.icon.url,
                  iconSize: [e.icon.width, e.icon.height],  //图标宽高
                  iconAnchor:[e.icon.iconAnchorX,e.icon.iconAnchorY]
                })

                let lay = L.marker(e.latlngs, {icon: greenIcon}).addTo(this.map).on('pm:edit', this.editLayClick).on('click', this.layClick).on('pm:dragend', this.dragendLayClick)
                // 这个id每次刷新,来回切换平面图,都会改变,所以要更新一下,如果不更新,下面删除,拖拽,编辑,的时候就找不到图层了
                res[i]._leaflet_id = lay._leaflet_id
              }else{
                let lay = L.marker(e.latlngs).addTo(this.map).on('pm:edit', this.editLayClick).on('click', this.layClick).on('pm:dragend', this.dragendLayClick)
                // 这个id每次刷新,来回切换平面图,都会改变,所以要更新一下,如果不更新,下面删除,拖拽,编辑,的时候就找不到图层了
                res[i]._leaflet_id = lay._leaflet_id
              }
            }
          }
          this.SaveData = JSON.stringify(res)
        }
        // 给地图绑定绘制、删除的事件
        this.map.on('pm:create', this.createClick)
        this.map.on('pm:remove', this.removeClick)

        // 禁止背景图拖拽,默认是可拖拽的
        // this.map.dragging.disable()
        // 禁止双击缩放,默认是可双击缩放的
        // this.map.doubleClickZoom.disable()
        // 禁止滚动缩放,默认是可滚动缩放的
        // this.map.scrollWheelZoom.disable()
      },
      mapDrag(){
        this.map.panInsideBounds(this.bounds, { animate: false });
      },
      // 图层绘制完成
      createClick(e) {
        let res = JSON.parse(this.SaveData) || []
        // 没有值没有绘制,赋个值初始值
        if (!res) res = []
        let obj={}
        console.log(e)
        if(e.shape=="Marker")
        {
          obj = {
            type:"maker",
            _leaflet_id: e.layer._leaflet_id,
            latlngs: e.layer._latlng
          }
        }
        else {
          // 设置图层样式
          e.layer.setStyle({
            fillOpacity: 0.5,
            fillColor: this.fillColor,
            // color：线段颜色
            // weight：线段宽度
            // opacity：线段透明度
            // dashArray：虚线间隔
            // fill：是否填充内部(true/false)
            // fillColor:内部填充颜色，如不设置，默认为color颜色
            // fillOpacity：内部填充透明度
          })
          // 只要两个参数即可保存该图层信息,id和经纬度
          obj = {
            type:"layer",
            _leaflet_id: e.layer._leaflet_id,
            latlngs: e.layer._latlngs,
          }
        }

        obj.icon={
          url:"",
          width:50,
          height:50,
          iconAnchorX:0,
          iconAnchorY:0,
        }
        obj.link={
          linkType:"",
          isPopUp:false,
          width:50,
          height:50,
          title:"",
          autoClose:false,
          External:"",
          OpenExternalType:"new",
          Inside:{
            displayUUID:"",
            pageUUID:"",
          }
        }
        res.push(obj)
        this.SaveData = JSON.stringify(res)
        // 给图层绑定点击、拖拽、编辑事件
        e.layer.on('pm:edit', this.editLayClick).on('click', this.layClick).on('pm:dragend', this.dragendLayClick)
      },
      // 图层删除
      removeClick(e) {
        let res = JSON.parse(this.SaveData) || []
        // 没有值没有绘制,赋个值初始值
        if (!res)
        {
          return
        }
        for(let i=0;i<res.length;i++)
        {
          if( e.layer._leaflet_id==res[i]._leaflet_id){
            res.splice(i,1)
            break
          }
        }
        this.SaveData = JSON.stringify(res)
      },
      // 区域图层点击
      layClick(e) {
        let res = JSON.parse(this.SaveData) || []
        // 没有值没有绘制,赋个值初始值
        if (!res) {
          return
        }
        for (let i = 0; i < res.length; i++) {
          if ((res[i]._leaflet_id != "undefined") && (e.target._leaflet_id == res[i]._leaflet_id)) {
            if (typeof res[i].icon == "undefined") {
              this.clickData.icon = {
                url: "",
                width: 50,
                height: 50,
              }
            } else {
              this.clickData.icon = res[i].icon
            }

            if (typeof res[i].link == "undefined") {
              this.clickData.link = {
                linkType: "",
                isPopUp: false,
                width: 50,
                height: 50,
                title: "",
                autoClose: false,
                External: "",
                OpenExternalType: "new",
                Inside: {
                  displayUUID: "",
                  pageUUID: "",
                }
              }
            } else {
              this.clickData.link = res[i].link
            }

            break
          }
        }
        this.currMaker = e.target
        if(this.editMode) {

          this.attrVisible = true
        }else{
          let link = this.clickData.link
          let linkInfo = {
            IsPopUp :link.isPopUp,
            autoClose :link.autoClose,
            linkType :link.linkType,
            ModelId:link.Inside.displayUUID,
            PageUuid:link.Inside.pageUUID,
            width:link.width,
            height:link.height,
            External:link.External,
            title:link.title,
            OpenExternalType:link.OpenExternalType
          }
          this.$EventBus.$emit("GoPage", linkInfo);
        }

      },
      // 区域图层拖拽
      dragendLayClick(e) {

      },
      // 区域图层编辑
      editLayClick(e) {

      },
    },
    mounted() {
      let _t = this
      this.initComponents(this.detail)
      this.$nextTick(function(){
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {


        })
        _t.$EventBus.$on(animateEvent, (data) => {
          _t.isStart = data
        })

      });
    },
    created(){
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
        _t.initComponents(_t.detail)
      })
    }
}
</script>
<style scoped>
.leaflet-container {
   background: transparent;
  /*outline-offset: 1px;*/
}
</style>
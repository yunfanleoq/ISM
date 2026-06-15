<template>
  <div :style="animatedStyle" v-show="detail.style.visible==1 ||isStart? true:false">
    <div class="view-menu-panel" :class="{
          'animated':true,[`${detail.style.animate}`]: true,
          'color-animation':isStart&&animateType.includes('millcolorGrad')&&!IsToolBox,
          'blink-animation':isStart&&animateType.includes('blink')&&!IsToolBox,
          'scale-animation':isStart&&animateType.includes('Zoom')&&!IsToolBox,
          'rotate-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,
          'rotate-anti-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1
        }"
         :style="{
                                width: detail.style.position.w + 'px',
                                height: detail.style.position.h + 'px',
                                'background-color': detail.style.backColor,
                                'border-radius':detail.style.BorderEdges+'px',
                                opacity:detail.style.opacity,
                                borderWidth: detail.style.borderWidth + 'px',
                                borderStyle: detail.style.borderStyle,
                                borderColor: detail.style.borderColor,
                                transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
                            }">
          <a-menu
            ref="menuRef"
            class="view-menu-scroll"
            :class="{ 'view-menu-has-scroll': hasMenuScroll }"
            mode="inline"
            :style="styleVar"
            @click="handleMenuClick"
            @openChange="updateMenuScrollState"
          >
            <template v-for="item in detail.style.MenuConfig">
              <template v-if="!item.children || !Array.isArray(item.children) || item.children.length === 0">
                <a-menu-item
                  :key="item.key"
                  class="submenu-title-wrapper"
                  @contextmenu.prevent.stop="showContextMenu($event, item)"
                >
                  {{ item.title }}
                </a-menu-item>
              </template>
              <template v-else>
                <a-sub-menu
                  :key="item.key"
                  @contextmenu.prevent.stop="showContextMenu($event, item)"
                >
                  <div slot="title" class="submenu-title-wrapper">
                    {{ item.title }}
                  </div>
                  <template v-for="child in item.children">
                    <template v-if="!child.children || child.children.length === 0">
                      <a-menu-item
                        :key="child.key"
                        class="submenu-title-wrapper"
                        @contextmenu.prevent.stop="showContextMenu($event, child)"
                      >
                        {{ child.title }}
                      </a-menu-item>
                    </template>
                    <template v-else>
                      <a-sub-menu
                        :key="child.key"
                        @contextmenu.prevent.stop="showContextMenu($event, child)"
                      >
                        <div slot="title" class="submenu-title-wrapper">
                          {{ child.title }}
                        </div>
                        <template v-for="grandchild in child.children">
                          <template v-if="!grandchild.children || grandchild.children.length === 0">
                            <a-menu-item
                              :key="grandchild.key"
                              class="submenu-title-wrapper"
                              @contextmenu.prevent.stop="showContextMenu($event, grandchild)"
                            >
                              {{ grandchild.title }}
                            </a-menu-item>
                          </template>
                          <template v-else>
                            <a-sub-menu
                              :key="grandchild.key"
                              @contextmenu.prevent.stop="showContextMenu($event, grandchild)"
                            >
                              <div slot="title" class="submenu-title-wrapper">
                                {{ grandchild.title }}
                              </div>
                              <template v-for="greatgrandchild in grandchild.children">
                                <a-menu-item
                                  :key="greatgrandchild.key"
                                  class="submenu-title-wrapper"
                                  @contextmenu.prevent.stop="showContextMenu($event, greatgrandchild)"
                                >
                                  {{ greatgrandchild.title }}
                                </a-menu-item>
                              </template>
                            </a-sub-menu>
                          </template>
                        </template>
                      </a-sub-menu>
                    </template>
                  </template>
                </a-sub-menu>
              </template>
            </template>
          </a-menu>
          <a-modal :visible="PopUpDialog"
                   @ok="doAddMenu"
                   :title="DoType==1?$t('configComponent.Menu.addMenu'):$t('configComponent.Menu.editMenu')"
                   :width="500"
                   @cancel="PopUpDialog=false"
                   v-drag-modal
                   :destroyOnClose="true"
                   :maskClosable="false"
                   :mask="true">
        <div >
           <a-form  layout="vertical" >
             <a-form-item v-if="SelectType==2" :label="$t('configComponent.Menu.parentMenu')">
              <a-select
                  v-model="SelectKey"
                  :label-in-value="true"
                  placeholder="请选择父菜单"
              >
                <a-select-option v-for="item in flattenMenuList" :key="item.key" :value="item.key">
                  {{ item.title }}
                </a-select-option>
              </a-select>
             </a-form-item>

             <a-form-item :label="$t('configComponent.Menu.menuName')">
              <a-input v-model="MenuName">

              </a-input>
             </a-form-item>
             <a-form-item :label="$t('displayConfig.Properties.linkIAppUUID')">
              <a-select
                  v-model="displayUUID"
                  allowClear
              >
                <a-select-option v-for="options in configurationModel" :key="options.uuid" :value="options.uuid">
                  {{ options.name}}
                </a-select-option>
              </a-select>
             </a-form-item>

             <a-form-item :label="$t('displayConfig.Properties.linkIAppPageUUID')">
              <a-select
                  v-model="SelectPage"
              >
                <a-select-option v-for="options in generateTargetPage(displayUUID)" :key="options.value" :value="options.value">
                  {{ options.label}}
                </a-select-option>
              </a-select>
             </a-form-item>


             <a-form-item :label="$t('configComponent.Menu.IsPopUp')">
                <a-checkbox :checked="IsPopUp"  v-model="IsPopUp">
                  {{ $t('displayConfig.Properties.isLinkPopUp') }}
                </a-checkbox>
             </a-form-item>

           </a-form>
        </div>
      </a-modal>
    </div>
  </div>

</template>

<script>
import svgView from '../View';
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";
import Contextmenu from "vue-contextmenujs"
import Vue from 'vue'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
Vue.use(Contextmenu);
export default {
  mixins: [ISMChildAutoMixin],

    name: 'view-menu-nav',
    i18n: require('../../../../i18n/language'),
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:false,
        displayUUID:"",
        MenuName:"",
        configurationModel:[],
        PopUpDialog:false,
        SelectKey:"",
        SelectType:1,
        SelectPage:"",
        IsPopUp:false,
        displayPageList:new Map,
        Text:"",
        foreColor:"#000000",
        backColor:"#ffffff",
        hoverForeColor:"#000000a6",
        hoverBackColor:"#DAEAF6",
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
        fillOpacity:1,
        DoType:1,
        strokeOpacity:1,
        animateType:"blink",
        startColor:"#74f808",
        stopColor:"#74f808",
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        hasMenuScroll:false,
        MenuList:[],
        ClickType:0,
        placement:"bottomCenter",
        base:{
          text: "configComponent.Menu.title",
          "icon": "icon-daohangcaidan",
          "isFontIcon": true,
          "info": {
            "type": "image",
            "action": [],
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
            "active": [],
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 300,
                "h": 400
              },
              "visible":1,
              "backColor": "#ffffff",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              fontSize: 20,
              fontFamily: "Arial",
              italic:0,
              textAlign: "center",
              fontWeight:400,
              MenuConfig:[
                { key: '1', title: '首页' ,path:""},
                { key: '2', title: '产品中心',children:[
                    { key: '3', title: '菜单3', path:""},
                    { key: '4', title: '菜单4',path:"" }
                  ]}
              ],
              "diy":[
                {
                  "name":"configComponent.ComBox.hoverForce",
                  "type": 2,
                  "value": "#000000",
                  "key":"hoverForce",
                },
                {
                  "name":"configComponent.ComBox.hoverBack",
                  "type": 2,
                  "value": "#6ACBFF",
                  "key":"hoverBack",
                },
                {
                  "name":"configComponent.Menu.ClickType",
                  "type":6,
                  "enumList":[
                    {option:'configComponent.Menu.ClickTypeEmit',value:1},{option:'configComponent.Menu.ClickTypeJump',value:0}
                  ],
                  "value":0,
                  "key":"ClickType",
                },
              ]
            }
          }
        }
      }
    },
    computed: {
      animatedStyle(){
        return {
          "--tableBorderColor": this.tableBorderColor ,
          "--blinkSpeed":this.blinkSpeed+'s',
          "--stopColor":this.stopColor,
          "--startColor":this.startColor,
          "--animateSpeed":this.animateSpeed+'s',
          "--animateSpinSpeed":this.animateSpinSpeed+'s',
          "--height": this.detail.style.position.h+'px',
          "--width": this.detail.style.position.w+'px',
          "--foreColor": this.foreColor ,
          '--backColor':this.backColor,
          '--fontWeight':this.detail.style.fontWeight,
          '--fontSize':this.detail.style.fontSize+ 'px',
          '--fontFamily':this.detail.style.fontFamily,
          '--fontStyle':this.detail.style.italic?'oblique':'normal',
        }
      },
      textAlign: function(){
        if(this.detail.style.textAlign == undefined) {
          return "center";
        } else {
          return this.detail.style.textAlign;
        }
      },
      lineHeight: function() {
        if(this.detail.style.lineHeight == undefined) {
          return this.detail.style.position.h;
        }
        return this.detail.style.lineHeight;
      },
      flattenMenuList() {
        const result = []
        const flatten = (menus, prefix = '') => {
          for(const menu of menus) {
            result.push({
              key: menu.key,
              title: prefix + menu.title
            })
            if(menu.children && menu.children.length > 0) {
              flatten(menu.children, prefix + menu.title + ' / ')
            }
          }
        }
        if(this.detail && this.detail.style && this.detail.style.MenuConfig) {
          flatten(this.detail.style.MenuConfig)
        }
        return result
      },
      

      styleVar() {
        return {
          "--height": this.detail.style.position.h+'px',
          "--width": this.detail.style.position.w+'px',
          "--foreColor": this.foreColor ,
          '--backColor':this.backColor,
          "--hoverForeColor": this.hoverForeColor ,
          '--hoverBackColor':this.hoverBackColor,
          '--fontWeight':this.detail.style.fontWeight,
          '--fontSize':this.detail.style.fontSize+ 'px',
          '--fontFamily':this.detail.style.fontFamily,
          '--fontStyle':this.detail.style.italic?'oblique':'normal',
        };
      },
    },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          if(this.editMode) {
            this.initComponents(newVal);
          }
          this.updateMenuScrollState()
        },
        deep: true
      }
    },
    methods: {
      AddMenu(){
        if(this.DoType==2) {
          const menuData = this.findMenuByKey(this.detail.style.MenuConfig, this.SelectKey?.value)
          if (menuData) {
             this.MenuName = menuData.title ;
             this.IsPopUp = menuData.IsPopUp ;
             this.displayUUID = menuData.DisPlayID;
             this.SelectPage = menuData.path;
          }
        }
        else
        {
          this.MenuName = "" ;
          this.IsPopUp = false;
          this.displayUUID = "";
          this.SelectPage = "";
        }
        this.PopUpDialog = true
      },
      removeMenu(key) {
        const menuData = this.findParentsMenuByKey (this.detail.style.MenuConfig,key,this.detail.style.MenuConfig)
        for(let i=0;i<this.detail.style.MenuConfig.length;i++)
        {
          if(this.detail.style.MenuConfig[i].key==key)
          {
            this.detail.style.MenuConfig.splice(i,1)
            this.updateMenuScrollState()
            return
          }
        }
        if (menuData.children && menuData.children.length > 0)
        {
          for(let i=0;i<menuData.children.length;i++)
          {
            if(menuData.children[i].key==key)
            {
              menuData.children.splice(i,1)
              this.updateMenuScrollState()
            }
          }
        }
      },
      doAddMenu(item) {
        if(this.MenuName.length==0||this.SelectPage.length==0)
        {
          return
        }
        if(this.DoType==1)
        {
          let keyIndex = Date.now()

          if(this.SelectType==1)
          {
            this.detail.style.MenuConfig.push({ key: keyIndex,IsPopUp:this.IsPopUp,DisPlayID:this.displayUUID, path:this.SelectPage,title: this.MenuName });
          }
          else if(this.SelectType==2)
          {
            this.addSubMenu(this.SelectKey?.value,this.MenuName)
          }
        }
        else
        {
          const menuData = this.findMenuByKey (this.detail.style.MenuConfig,this.SelectKey?.value)
          if (menuData) {
            menuData.title = this.MenuName;
            menuData.IsPopUp = this.IsPopUp;
            menuData.DisPlayID = this.displayUUID;
            menuData.path = this.SelectPage;
          }
        }
        this.PopUpDialog=false
        this.updateMenuScrollState()
      },
      addSubMenu(addkey,name) {
        if(!addkey) {
          this.$message.error(this.$t('configComponent.Menu.selectParentMenu'))
          return
        }
        
        let keyIndex = Date.now()
        const menuData = this.findMenuByKey (this.detail.style.MenuConfig,addkey)
        
        if(!menuData) {
          this.$message.error(this.$t('configComponent.Menu.parentMenuNotFound'))
          return
        }
        
        if (!menuData.children) {
          menuData.children=[]
        }
        menuData.children.push({ key: keyIndex,IsPopUp:this.IsPopUp,DisPlayID:this.displayUUID, path:this.SelectPage, title: name })
        this.$forceUpdate()
        this.updateMenuScrollState()
      },
      findMenuByKey(menus, key) {
        for (const menu of menus) {
          if (menu.key === key) {
            return menu; // 找到匹配的菜单
          }
          if (menu.children && menu.children.length > 0) {
            const found = this.findMenuByKey(menu.children, key); // 递归查找子菜单
            if (found) {
              return found;
            }
          }
        }
        return null; // 未找到
      },
      findParentsMenuByKey(menus, key,Parents) {
        for (const menu of menus) {
          if (menu.key === key) {
            return Parents; // 找到匹配的菜单
          }
          if (menu.children && menu.children.length > 0) {
            const found = this.findParentsMenuByKey(menu.children, key,menu); // 递归查找子菜单
            if (found) {
              return menu;
            }
          }
        }
        return null; // 未找到
      },
      handleMenuClick(e) {
        this.updateMenuScrollState()
        const menudata = this.findMenuByKey (this.detail.style.MenuConfig,e.key)
        if(menudata!=null&&menudata.path!="")
        {
          let item={
            DisPlayID:menudata.DisPlayID,
            IsPopUp:menudata.IsPopUp,
            MenuName:menudata.title,
            PageID:menudata.path,
          }
          if(this.ClickType==1)
          {
            this.$EventBus.$emit("MenuConfigPage", item);
          }
          else
          {
            this.JumpPage(item)
          }
        }
      },
      generateTargetPage (uuid) {
        return this.displayPageList.get(uuid)
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
      showContextMenu($event,item) {
        let _t = this
        console.log('[ViewMenu] showContextMenu called, editMode:', this.editMode, 'item:', item)
        if(!this.editMode)
        {
          console.log('[ViewMenu] editMode is false, returning')
          return
        }
        $event.preventDefault()
        $event.stopPropagation()
        this.$contextmenu({
          items: [
            {
              label: _t.$t('configComponent.Menu.addMenu'),
              onClick: () => {
                _t.SelectKey = { value: item.key, label: item.title }
                _t.SelectType = 1
                _t.DoType = 1
                _t.AddMenu(item,1)
              }
            },
            {
              label: _t.$t('configComponent.Menu.addSubMenu'),
              onClick: () => {
                _t.SelectKey = { value: item.key, label: item.title }
                _t.SelectType = 2
                _t.DoType = 1
                _t.AddMenu(item,2)
              }
            },
            {
              label: _t.$t('configComponent.Menu.editMenu'),
              onClick: () => {
                _t.SelectKey = { value: item.key, label: item.title }
                _t.SelectType = 0
                _t.DoType = 2
                _t.AddMenu(item,2)
              }
            },
            {
              label: _t.$t('configComponent.Menu.removeMenu'),
              onClick: () => {
                _t.$confirm({
                  title: _t.$t('configComponent.Menu.deleteConfirm'),
                  content: _t.$t('configComponent.Menu.deleteConfirmContent'),
                  onOk() {
                    _t.removeMenu(item.key)
                  },
                  onCancel() {}
                })
              }
            }
          ],
          event: $event, // 鼠标事件信息
          divided:true,
          customClass: "custom-class", // 自定义菜单 class
          zIndex: 10000, // 菜单样式 z-index
          minWidth: 230 // 主菜单最小宽度
        });
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.backColor=this.detail.style.backColor
        this.foreColor=this.detail.style.foreColor
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="hoverBack")
          {
            this.hoverBackColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="hoverForce")
          {
            this.hoverForeColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ClickType")
          {
            this.ClickType=option.style.diy[i].value
          }
        }
        i=0
        this.animateType = option.animate.selected
        if(option.animate.isExpression)
        {
          this.isStart = false
        }
        else
        {
          this.isStart = true
        }
        for( i=0;i<option.animate.animateElement.length;i++)
        {
          if(option.animate.animateElement[i].id=="millcolorGrad")
          {
            for(let k =0;k<option.animate.animateElement[i].elementList.length;k++)
            {
              if(option.animate.animateElement[i].elementList[k].key=="startColor")
              {
                this.startColor=option.animate.animateElement[i].elementList[k].value
              }
              else if(option.animate.animateElement[i].elementList[k].key=="stopColor")
              {
                this.stopColor=option.animate.animateElement[i].elementList[k].value
              }
              else if(option.animate.animateElement[i].elementList[k].key=="animateSpeed")
              {
                this.animateSpeed=option.animate.animateElement[i].elementList[k].value
              }
            }
          }
          else if(option.animate.animateElement[i].id=="blink")
          {
            for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
              if (option.animate.animateElement[i].elementList[k].key == "blinkSpeed") {
                this.blinkSpeed = option.animate.animateElement[i].elementList[k].value
              }
            }
          }
          else if(option.animate.animateElement[i].id=="animateSpin")
          {
            for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
              if (option.animate.animateElement[i].elementList[k].key == "spinSpeed") {
                this.animateSpinSpeed = option.animate.animateElement[i].elementList[k].value
              }
              else if (option.animate.animateElement[i].elementList[k].key == "spinDirection") {
                this.spinDirection = option.animate.animateElement[i].elementList[k].value
              }
            }
          }
        }
      },
      JumpPage(item){
        this.$EventBus.$emit("ChargePage", item);
      },
      updateMenuScrollState() {
        this.$nextTick(() => {
          const update = () => {
            const menu = this.$refs.menuRef && (this.$refs.menuRef.$el || this.$refs.menuRef)
            if(!menu) {
              this.hasMenuScroll = false
              return
            }
            this.hasMenuScroll = menu.scrollHeight > menu.clientHeight + 1
          }
          update()
          clearTimeout(this._menuScrollTimer)
          this._menuScrollTimer = setTimeout(update, 350)
        })
      },
      bindContextMenuEvents() {
        const _t = this
        const container = this.$el
        
        if(!container) return
        
        container.addEventListener('contextmenu', _t._contextMenuHandler = function(e) {
          if(!_t.editMode) return
          
          e.preventDefault()
          e.stopPropagation()
          
          const target = e.target
          const menuItemEl = target.closest('.ant-menu-item, .ant-menu-submenu-title')
          
          if(menuItemEl) {
            const key = menuItemEl.getAttribute('data-key') || 
                        menuItemEl.getAttribute('key') ||
                        menuItemEl.dataset.key
            
            if(key) {
              const item = _t.findMenuByKey(_t.detail.style.MenuConfig, key)
              if(item) {
                _t.showContextMenu(e, item)
                return
              }
            }
          }
          
          _t.showBlankContextMenu(e)
        })
      },
      
      showBlankContextMenu(event) {
        let _t = this
        this.$contextmenu({
          items: [
            {
              label: _t.$t('configComponent.Menu.addMenu'),
              onClick: () => {
                _t.SelectKey = null
                _t.SelectType = 1
                _t.DoType = 1
                _t.AddMenu(null, 1)
              }
            },
            {
              label: _t.$t('configComponent.Menu.addSubMenu'),
              onClick: () => {
                _t.SelectKey = null
                _t.SelectType = 2
                _t.DoType = 1
                _t.AddMenu(null, 2)
              }
            }
          ],
          event: event,
          divided: true,
          customClass: "custom-class",
          zIndex: 10000,
          minWidth: 230
        });
      },
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
        this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, _t._activeEventHandler = (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          if(data.ID == "ControlStatus")
          {
            let temp  =  parseFloat(data.result)
            if(!isNaN(temp))
            {
              _t.StatusValue=temp
            }
          }
        })
        _t.$EventBus.$on(animateEvent, _t._animateEventHandler = (data) => {
            _t.isStart = data
          })

        _t.bindContextMenuEvents()
        _t.updateMenuScrollState()
      });
    },
    created(){
      this.getConfigurationModel()
      let _t = this
      const node = this.getNode()
      
      _t.$EventBus.$on('cell-editMode', _t._cellEditModeHandler = (data) => {
        if (!data || typeof data !== 'object') {
          return
        }
        if (data.source === 'ViewPagerContainer') {
          return
        }
        if ('edit' in data) {
          _t.editMode = data.edit
        }
        if ('toolbox' in data) {
          _t.IsToolBox = data.toolbox
        }
      })
      
      node.on('change:data', ({ current }) => {
        if(current) {
          _t.detail = current.detail
        }
      })
      node.on('change:size', ({ current }) => {
        _t.detail.style.position.w = current.width
        _t.detail.style.position.h = current.height
        _t.updateMenuScrollState()
      });
      
      this.detail = node.getData().detail
      this.editMode = node.getData().editMode
      this.showDeviceUuid = node.getData().showDeviceUuid
      this.IsToolBox = node.getData().IsToolBox
      
      this.initComponents(this.detail);
    },
    beforeDestroy() {
      // 清理 contextmenu 监听
      if (this._contextMenuHandler && this.$el) {
        this.$el.removeEventListener('contextmenu', this._contextMenuHandler)
      }
      // 清理 EventBus 监听
      if (this._activeEventHandler) {
        this.$EventBus.$off(this.detail.identifier + 'activeEvent', this._activeEventHandler)
      }
      if (this._animateEventHandler) {
        this.$EventBus.$off(this.detail.identifier + 'animateEvent', this._animateEventHandler)
      }
      if (this._cellEditModeHandler) {
        this.$EventBus.$off('cell-editMode', this._cellEditModeHandler)
      }
      clearTimeout(this._menuScrollTimer)
    }
}
</script>
<style lang="less" scoped>
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
::v-deep .ant-menu{
  background-color:  var(--backColor);
}
.view-menu-panel{
  overflow: hidden;
  box-sizing: border-box;
}
::v-deep .view-menu-scroll{
  width: 100%;
  height: 100%;
  overflow-y: hidden;
  overflow-x: hidden;
  box-sizing: border-box;
  overscroll-behavior: contain;
  scrollbar-width: thin;
  scrollbar-color: skyblue transparent;
}
::v-deep .view-menu-scroll.view-menu-has-scroll{
  overflow-y: auto;
}
::v-deep .view-menu-scroll::-webkit-scrollbar{
  width: 4px;
  height: 4px;
}
::v-deep .view-menu-scroll::-webkit-scrollbar-track{
  background: transparent;
  box-shadow: none;
}
::v-deep .view-menu-scroll::-webkit-scrollbar-thumb{
  border-radius: 10px;
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
::v-deep .view-menu-scroll.ant-menu-inline,
::v-deep .view-menu-scroll .ant-menu-inline{
  border-right: 0 !important;
}
::v-deep .view-menu-scroll .ant-menu-submenu:not(.ant-menu-submenu-open) > .ant-menu-sub{
  display: none !important;
}
.submenu-title-wrapper{
  font-size: var(--fontSize);
  font-family: var(--fontFamily);
  font-weight:var(--fontWeight);
  font-style:var(--fontStyle);
  color:  var(--foreColor);
}
::v-deep .ant-menu-horizontal{
  border-bottom: 1px solid #f0f0f0;
}
::v-deep .ant-menu-inline .ant-menu-item::after{
  border-right: 3px solid var(--hoverForeColor);
}
::v-deep .ant-menu-submenu-arrow::before {
  color: var(--hoverForeColor) !important;
}
::v-deep .ant-menu-submenu-arrow::after {
  color: var(--hoverForeColor) !important;
}
::v-deep .ant-menu-item-active {
  background-color:  var(--hoverBackColor) !important;
  color: var(--hoverForeColor) !important;
}
::v-deep .ant-menu-item-selected {
  background-color:  var(--hoverBackColor) !important;
  color: var(--hoverForeColor) !important;
}
::v-deep .ant-menu-submenu-selected {
  background-color: transparent !important;
  color: var(--hoverForeColor) !important;
}
::v-deep .ant-menu-submenu-selected > .ant-menu-submenu-title {
  background-color:  var(--hoverBackColor) !important;
  color: var(--hoverForeColor) !important;
}
::v-deep .ant-menu-submenu-active {
  background-color: transparent !important;
}
::v-deep .ant-menu-submenu-active > .ant-menu-submenu-title {
  background-color:  var(--hoverBackColor) !important;
  color: var(--hoverForeColor) !important;
}

</style>

<template>
  <div>
    <a-modal :visible="systemTempleteDialog"
            :title="$t('component.systemPageTemplete.title')"
             @cancel="systemTempleteDialog=false"
             v-drag-modal
            :resizable="true"
             :footer="null"
             width="700px"
             height="600px"
            :modal="false">

      <div>
        <div style="margin-top: 10px;height:530px;width: 650px;">
          <a-tabs
              :default-active-key="defaultIndex"
              tab-position="left"
              style="height:530px;"
          >
            <a-tab-pane  :key="index" :tab="$t(template.title)" v-for="(template,index) in templateList">
              <a-list :grid="{ gutter: 1, column: 2 }" size="small" :pagination="systemPagination" :data-source="template.items">
                <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
                  <vue-hover-mask >
                    <!-- 默认插槽 -->
                    <div class="companyImg" style="width: 220px;height: 220px;float: left;">
                      <img style="width: 220px;height: 220px;" :src="item.pagePicUrl">
                    </div>
                    <!-- action插槽 -->
                    <template v-slot:action  >
                      <span style="font-size: 20px" @click="selectTemplete(item.pageTemplete)"><a-icon type="export" />{{$t('component.systemPageTemplete.export')}}</span>
                      <span style="font-size: 20px;margin-left: 10px" @click="viewTemplate(item.pagePicUrl)"><a-icon type="eye" />{{$t('component.systemPageTemplete.view')}}</span>
                    </template>
                  </vue-hover-mask>
                </a-list-item>
              </a-list>
            </a-tab-pane>
          </a-tabs>
        </div>
      </div>
    </a-modal>
  </div>

</template>
<script>

import VueHoverMask from "@/components/VueHoverMask/VueHoverMask"
import {GetSystemPageTemplete} from "@/services/system";

export default {
  name: 'systemPageTemplete',
  i18n: require('../../i18n/language'),
  data() {
    return {
      pagination: {
        onChange: page => {
          console.log(page);
        },
        pageSize: 4,
      },
      systemTempleteDialog:false,
      defaultIndex:0,
      systemPagination: {
        onChange: page => {
          console.log(page);
        },
        pageSize: 4,
      },
      ShowIndex:0,
      templateList:[],
      radioValue:1,
      previewVisible: false,
      previewImage: '',
      visible: false,
      imageList:[],
      model3DList:[],
    };
  },
  components: {
    VueHoverMask,
  },
  activated(){
    this.$refs.systemTempleteDialog.close()
    this.getSystemTempleteList()
  },
  props: {

  },
  watch: {
    '$route'() {

    }
  },
  filters: {
    filterFun: function(value) {
      if (value && value.length > 8) {  //字符最大长度
        value = value.substring(0, 8) + "...";  //超过省略
      }

      return value;
    },
  },
  mounted(){
    this.getSystemTempleteList()
  },
  methods: {
    async viewTemplate(file) {
      window.open(file, '_blank');
    },
    getSystemTempleteList(){
      let _t = this
      GetSystemPageTemplete().then(function (res){
        if(res.data.Result!=null)
        {
          _t.templateList=[]
          for(let i=0;i<res.data.Result.length;i++)
          {
              let templateGroup={
                title: "component.systemImageModel.Button",
                icon: "icon-background",
                opened: false,
                items:[]
              }
              templateGroup.title = res.data.Result[i].groupName

              for(let k=0;k<res.data.Result[i].list.length;k++)
              {
                let tableData={}
                tableData.pagePicUrl = res.data.Result[i].list[k].pagePicUrl
                tableData.pageTemplete = res.data.Result[i].list[k].pageTemplete
                templateGroup.items.push(tableData)
              }
            _t.templateList.push(templateGroup)
          }
        }

      })
    },
    showModal() {
      this.systemTempleteDialog = true
    },
    selectTemplete(url){
      this.$emit("OnSelectTemplete", url);
      this.visible =false
    },
    handleCancel() {
      this.visible = false;
    },
  },
  created() {

//按钮

  }
}
</script>

<style lang="less" scoped>

::v-deep .ant-tabs-nav .ant-tabs-tab:hover {
  color: #13c2c2;
  right: 1px;
}

::v-deep .ant-tabs-nav .ant-tabs-tab-active {
 color: #13c2c2;
  background-color: #e6fffb;
 font-weight: 500;
}

.ant-avatar-lg {
  width: 48px;
  height: 48px;
  line-height: 48px;
}

.list-content-item {
  color: rgba(0, 0, 0, .45);
  display: inline-block;
  vertical-align: middle;
  font-size: 14px;
  margin-left: 40px;
  span {
    line-height: 20px;
  }
  p {
    margin-top: 4px;
    margin-bottom: 0;
    line-height: 22px;
  }
}
</style>

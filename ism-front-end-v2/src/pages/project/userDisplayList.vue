<template>
  <project-layout>
  <div class="card-list">
    <a-row>
      <a-col :span="isMobile?24:6" v-for="(item, index) in modelList" :key="index" @click="GoToRun(item.uuid)">
        <a-card id="displayCard" hoverable style="width: 300px;border-radius: 2px 2px 0 0; margin-bottom: 20px;">
          <template #cover>
            <!-- 默认插槽 -->
            <img  style="width: 298px;height: 220px;cursor: pointer" :src="item.DisplayImage==''?'/static/images/pcDefaultCover.jpg':item.DisplayImage" />
          </template>
          <a-card-meta :title="item.name" :description=" item.description">
          </a-card-meta>
        </a-card>
      </a-col>
    </a-row>
  </div>
  </project-layout>
</template>

<script>
import {
  GetUserDisplayList
} from "../../services/displayModel";
import ProjectLayout from "../../layouts/ProjectLayout";
import {DISPLAYIMAGEUPLOAD} from "@/services/api";
import {AUTH_TYPE, setAuthorization} from "@/utils/request";
import {mapState} from 'vuex'
export default {
  name: 'UserDisplayModelList',
  i18n: require('../../i18n/language'),
  data () {
    return {
      editIndex:0,
      isEditStatus:0,
      loading: false,
      uploadDisPlayUrl:DISPLAYIMAGEUPLOAD,
      visible:false,
      listIndex:-1,
      pagination: {
        hideOnSinglePage:true,
        showQuickJumper:true,
        pageSize: 4,
      },
      modelList:[]
    }
  },
  components: {ProjectLayout},
  computed: {
    ...mapState('setting', ['isMobile']),
  },
  mounted(){
    console.log("this.isMobile",this.isMobile)
  },
  activated(){

  },
  created(){
    this.getModelList()
  },
  watch: {
    '$route' () {
      this.modelList=[]

      this.getModelList()
    }
  },
  methods: {
    GoToRun(pageid){
      this.$router.push('/AppRun/'+pageid)
    },
    getModelList(){
      this.modelList=[]
      let _t = this
      const  params= {
        uuid:this.$route.params.uuid
      }
      this.loading = true
      GetUserDisplayList(params).then(function (res){
        if(res.data.list!=null)
        {
          if(res.data.list.length>1)
          {
            for(let i=0;i<res.data.list.length;i++)
            {
              let tableData={}
              tableData.name = res.data.list[i].name
              tableData.description = res.data.list[i].description
              tableData.uuid = res.data.list[i].displayUid
              if(res.data.list[i].DisplayImage=="")
              {
                tableData.DisplayImage = ""
              }else{
                tableData.DisplayImage =res.data.list[i].DisplayImage
              }

              _t.modelList.push(tableData)
              tableData={}
            }
          }
          else if(res.data.list.length==1)
          {
            setAuthorization({token: res.data.list[0].project_uuid},AUTH_TYPE.AUTH1)
            _t.$router.push('/AppRun/'+res.data.list[0].displayUid)
          }
        }
        _t.loading = false
      })
    },
  }
}
</script>

<style lang="less" scoped>

::v-deep #displayCard .ant-card-body {
  padding: 0px;
  margin-top: 5px;
  border-top: 1px solid #f0f0f0;

  margin-bottom: 10px;
  zoom: 1;
}

::v-deep  .vue-hover-mask {
  border-radius: 0px;
}
::v-deep #displayCard .ant-card-meta-detail   {
  margin-top: 5px;
  margin-left: 10px;
}
::v-deep #displayCard .ant-card-meta-detail > div:not(:last-child)  {
  margin-bottom: 0px;
}

::v-deep #displayCard .ant-card-head {
  min-height: 28px;
  margin-bottom: -1px;
  padding: 0 8px;
  color: #000;
  font-weight: 500;
  font-size: 13px;
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
  border-radius: 2px 2px 0 0;
  zoom: 1;
}
::v-deep .ant-card-actions{
  background-color: #fff;
}
</style>

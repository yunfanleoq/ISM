<template>
  <div>
    <a-modal
             :visible="systemImageModelDialog"
             @cancel="systemImageModelDialog=false"
             width="700px"
             height="600px"
             :footer="null"
             v-drag-modal
            :title="$t('component.systemImageModel.title')">

      <div>
        <a-tabs default-active-key="1">
          <a-tab-pane key="1" :tab="$t('component.systemImageModel.ProjectImages')" v-if="ShowIndex==0">
            <div style="padding: 10px;height: auto">
              <div style="margin-bottom: 5px">
                <a-radio-group v-model="radioValue" >
                  <a-radio :value="1">
                    {{$t('component.systemImageModel.localImage')}}
                  </a-radio>
                  <a-radio :value="2">
                    {{$t('component.systemImageModel.networkImage')}}
                  </a-radio>
                </a-radio-group>
                <a-upload v-if="radioValue==1"
                          name="file"
                          :multiple="true"
                          :action=uploadUrl
                          :showUploadList="false"
                          @change="afterUpload"
                >
                  <a-button> <a-icon type="upload" /> {{$t('displayConfig.Properties.upload')}} </a-button>
                </a-upload>
              </div>

              <div v-if="radioValue==1">
                <div >
                  <a-list :grid="{ gutter: 4, column: 4 }" size="small" :pagination="pagination" :data-source="imageList">
                    <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
                      <vue-hover-mask>
                        <!-- 默认插槽 -->
                        <div style="text-align: center;display: block;width: 128px;height: 128px;cursor: pointer">
                          <img  style="width: 100%;height: 100%;cursor: pointer" :src="item.imgurl" />
                        </div>
                        <!-- action插槽 -->
                        <template v-slot:action>
                          <span style="font-size: 14px" @click="selectImage(item.imgurl)">{{$t('component.systemImageModel.selectImage')}}</span>
                          <a-divider type="vertical" />
                          <span  style="font-size: 14px" @click="delImage(item.imgurl)">{{$t('component.systemImageModel.delImage')}}</span>
                        </template>
                      </vue-hover-mask>
                    </a-list-item>
                  </a-list>
                  <a-modal :visible="previewVisible" :footer="null" @cancel="viewImageCancel" style="min-height: 100px">
                    <img alt="预览" style="width: 100%;min-height: 100px" :src="previewImage" />
                  </a-modal>
                </div>
              </div>

              <div v-if="radioValue==2">
                <a-form  :label-col="{ span: 3}" :wrapper-col="{ span: 20 }">
                  <a-form-item :label="$t('component.systemImageModel.networkImageUrl')" >
                    <a-input   v-model="networkImageUrl" :default-value="networkImageUrl"/>
                  </a-form-item>
                  <div style="margin-top: 5px">
                    <a-button key="submit"  type="primary" @click="selectImage(networkImageUrl)">{{$t('component.systemImageModel.networkImageBtn')}}</a-button>
                  </div>
                </a-form>
              </div>
            </div>
          </a-tab-pane>

          <a-tab-pane key="2" :tab="$t('component.systemImageModel.SystemImages')" v-if="ShowIndex==0">
            <div style="margin-top: 10px">
              <a-tabs
                  :default-active-key="defaultIndex"
                  tab-position="left"
                  style="height:570px;"
              >
                <a-tab-pane  :key="index" :tab="$t(image.title)" v-for="(image,index) in SystemImageList">
                  <a-list :grid="{ gutter: 2, column: 6 }" size="small" :pagination="systemPagination" :data-source="image.items">
                    <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
                      <vue-hover-mask @click="selectImage(item)">
                        <!-- 默认插槽 -->
                        <div class="companyImg" style="width: 64px;height: 64px;float: left;">
                          <img style="width: 64px;height: 64px;" :src="item">
                        </div>
                        <!-- action插槽 -->
                        <template v-slot:action  >
                          <span style="font-size: 14px">{{$t('component.systemImageModel.selectImage')}}</span>
                        </template>
                      </vue-hover-mask>
                    </a-list-item>
                  </a-list>
                </a-tab-pane>
              </a-tabs>
            </div>
          </a-tab-pane>

          <a-tab-pane key="3" :tab="$t('component.systemImageModel.DiySystemImages')" v-if="ShowIndex==0&&DiySystemImageList.length>0">
            <div style="margin-top: 10px">
              <a-tabs
                  :default-active-key="defaultIndex"
                  tab-position="left"
                  style="height:570px;"
              >
                <a-tab-pane  :key="index" :tab="$t(image.title)" v-for="(image,index) in DiySystemImageList">
                  <a-list :grid="{ gutter: 2, column: 6 }" size="small" :pagination="systemPagination" :data-source="image.items">
                    <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
                      <vue-hover-mask @click="selectImage(item)">
                        <!-- 默认插槽 -->
                        <div class="companyImg" style="width: 64px;height: 64px;float: left;">
                          <img style="width: 64px;height: 64px;" :src="item">
                        </div>
                        <!-- action插槽 -->
                        <template v-slot:action  >
                          <span style="font-size: 14px">{{$t('component.systemImageModel.selectImage')}}</span>
                        </template>
                      </vue-hover-mask>
                    </a-list-item>
                  </a-list>
                </a-tab-pane>
              </a-tabs>
            </div>
          </a-tab-pane>

          <a-tab-pane key="4" :tab="$t('component.systemImageModel.Model3D')" v-if="ShowIndex==1">
            <div style="padding: 10px;height: auto">
              <div style="margin-bottom: 5px">
                <a-radio-group v-model="radioValue" >
                  <a-radio :value="1">
                    {{$t('component.systemImageModel.localImage')}}
                  </a-radio>
                  <a-radio :value="2">
                    {{$t('component.systemImageModel.networkImage')}}
                  </a-radio>
                </a-radio-group>
                <a-upload v-if="radioValue==1"
                          name="file"
                          :multiple="true"
                          :action=uploadUrl
                          :showUploadList="false"
                          @change="afterUpload"
                >
                  <a-button> <a-icon type="upload" /> {{$t('displayConfig.Properties.upload')}} </a-button>
                </a-upload>
              </div>

              <div v-if="radioValue==1">
                <div >
                  <a-list  size="small" :pagination="pagination" :data-source="model3DList">
                    <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
                      <a slot="actions" @click="selectImage(item.imgurl)">{{$t('component.systemImageModel.selectImage')}}</a>
                      <a slot="actions" @click="delImage(item.imgurl)">{{$t('component.systemImageModel.delImage')}}</a>
                      <a-list-item-meta>
                        <span slot="title" >{{ item.imgurl }}</span>
                      </a-list-item-meta>
                    </a-list-item>
                  </a-list>
                  <a-modal :visible="previewVisible" :footer="null" @cancel="viewImageCancel" style="min-height: 100px">
                    <img alt="预览" style="width: 100%;min-height: 100px" :src="previewImage" />
                  </a-modal>
                </div>
              </div>

              <div v-if="radioValue==2">
                <a-form  :label-col="{ span: 3}" :wrapper-col="{ span: 20 }">
                  <a-form-item :label="$t('component.systemImageModel.networkImageUrl')" >
                    <a-input   v-model="networkImageUrl" :default-value="networkImageUrl"/>
                  </a-form-item>
                  <div style="margin-top: 5px">
                    <a-button key="submit"  type="primary" @click="selectImage(networkImageUrl)">{{$t('component.systemImageModel.networkImageBtn')}}</a-button>
                  </div>
                </a-form>
              </div>
            </div>
          </a-tab-pane>

          <a-tab-pane key="5" :tab="$t('component.systemImageModel.ProjectDoc')" v-if="ShowIndex==2">
            <div style="padding: 10px;height: auto">
              <div style="margin-bottom: 5px">
                <a-upload
                          name="file"
                          :multiple="true"
                          :action=uploadUrl
                          :showUploadList="false"
                          @change="afterUpload"
                >
                  <a-button> <a-icon type="upload" /> {{$t('displayConfig.Properties.upload')}} </a-button>
                </a-upload>
              </div>

              <div >
                <div >
                  <a-list  size="small" :pagination="pagination" :data-source="modelDocList">
                    <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
                      <a slot="actions" @click="selectImage(item.imgurl)">{{$t('component.systemImageModel.selectImage')}}</a>
                      <a slot="actions" @click="delImage(item.imgurl)">{{$t('component.systemImageModel.delImage')}}</a>
                      <a-list-item-meta>
                        <span slot="title" >{{ item.name }}</span>
                      </a-list-item-meta>
                    </a-list-item>
                  </a-list>
                  <a-modal :visible="previewVisible" :footer="null" @cancel="viewImageCancel" style="min-height: 100px">
                    <img alt="预览" style="width: 100%;min-height: 100px" :src="previewImage" />
                  </a-modal>
                </div>
              </div>

            </div>
          </a-tab-pane>
        </a-tabs>

      </div>
    </a-modal>
  </div>

</template>
<script>

import {SYSTEMIMAGEUPLOAD} from "@/services/api";
import {systemImageList,systemImageDel} from "@/services/systemImages";
import VueHoverMask from "@/components/VueHoverMask/VueHoverMask"

const SystemImageButtonList = require.context("../../../public/static/ISM/systemImage/button/", false, /\.png$/)
const SystemImageSafeList = require.context("../../../public/static/ISM/systemImage/safe/", false, /\.png$/)
const SystemImageWasteWaterTreatmentList = require.context("../../../public/static/ISM/systemImage/WasteWaterTreatment/", false, /\.png$/)

const SystemImageSwitchList = require.context("../../../public/static/images/switch/", false, /\.*$/)

const SystemImageIconList = require.context("../../../public/static/ISM/systemImage/icon/", false, /\.svg$/)
const SystemImagePageBackgroundList = require.context('../../../public/static/ISM/page/background/', false, /\.jpg$/)
import {GetUserCustomPel} from "@/services/system";
export default {
  name: 'systemImageModel',
  i18n: require('../../i18n/language'),
  data() {
    return {
      systemImageModelDialog:false,
      pagination: {
        onChange: page => {
          console.log(page);
        },
        pageSize: 12,
      },
      defaultIndex:0,
      systemPagination: {
        onChange: page => {
          console.log(page);
        },
        pageSize: 36,
      },
      ShowIndex:0,
      SystemImageList:[],
      DiySystemImageList:[],
      radioValue:1,
      previewVisible: false,
      previewImage: '',
      uploadUrl:SYSTEMIMAGEUPLOAD,
      visible: false,
      imageList:[],
      model3DList:[],
      modelDocList:[],
    };
  },
  components: {
    VueHoverMask,
  },
  props: {
    networkImageUrl: {
      type: String,
      default: "",
      required: true
    },
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
    this.getSystemImageList()
    this.GetUserCustomPel()
  },
  methods: {
    GetUserCustomPel(){
      let _t = this
      this.DiySystemImageList = []
      GetUserCustomPel().then(function (res){
        if(res.data.code==0)
        {

          if((res.data.list!=null)&&(res.data.list.length>0))
          {
            for(let i=0;i<res.data.list.length;i++)
            {
              if((res.data.list[i].FilePath!=null)&&(res.data.list[i].FilePath.length>0))
              {
                let imagelist = {}
                imagelist.title =  res.data.list[i].DirName
                imagelist.items = []
                for(let k=0;k<res.data.list[i].FilePath.length;k++)
                {
                  imagelist.items.push(res.data.list[i].FilePath[k])
                }
                _t.DiySystemImageList.push(imagelist)
              }
            }
          }
        }

      })
    },
    viewImageCancel() {
      this.previewVisible = false;
    },
    async viewImage(file) {

      this.previewImage = file
      this.previewVisible = true
    },
    formatDateTime  (date) {
      let dateGet = new Date(date)
      let y = dateGet.getFullYear();
      let m = dateGet.getMonth() + 1;
      m = m < 10 ? ('0' + m) : m;
      let d = dateGet.getDate();
      d = d < 10 ? ('0' + d) : d;
      let h = dateGet.getHours();
      h=h < 10 ? ('0' + h) : h;
      let minute = dateGet.getMinutes();
      minute = minute < 10 ? ('0' + minute) : minute;
      let second=dateGet.getSeconds();
      second=second < 10 ? ('0' + second) : second;
      return y + '-' + m + '-' + d+' '+h+':'+minute+':'+second;
    },
    afterUpload(info) {
      if (info.file.status === 'done') {
        let result = info.file.response
        if(result.Code==2002) {
          this.getSystemImageList()
          this.$message.success(this.$t('component.systemImageModel.uploadSuccess'))
        }
        else
        {
          this.$message.error(this.$t('component.systemImageModel.uploadFailed'))
        }
      }
    },
    getSystemImageList(){
      let _t = this
      systemImageList().then(function (res){
        if(res.data.list!=null)
        {

          _t.imageList=[]
          _t.model3DList=[]
          _t.modelDocList=[]
          for(let i=0;i<res.data.list.length;i++)
          {
            let tableData={}
            tableData.imgurl = res.data.list[i].path
            tableData.name = res.data.list[i].name
            tableData.startAt = _t.formatDateTime(res.data.list[i].UpdatedAt)
            if(res.data.list[i].type==1) {
              _t.imageList.push(tableData)
            }else if(res.data.list[i].type==2) {
              _t.model3DList.push(tableData)
          }else if(res.data.list[i].type==3) {
              _t.modelDocList.push(tableData)
            }
            tableData={}
          }
        }

      })
    },
    showModal(index) {
      this.visible = true
      this.ShowIndex = index?index:0
      this.systemImageModelDialog = true
    },
    selectImage(url){
      this.$emit("onSelectImage", url);
      this.systemImageModelDialog =false
    },
    delImage(url){
      let _t = this
      const params={
        path:url
      }
      systemImageDel(params).then(function (res){
        if(res.data.code==200)
        {
          _t.getSystemImageList()
          _t.$message.success(_t.$t('component.systemImageModel.delImageSuccess'))
        }
        else
        {
          _t.$message.error(_t.$t('component.systemImageModel.delImageFailed'))
        }
      })
    },
    handleCancel() {
      this.systemImageModelDialog = false;
    },
  },
  created() {
    const SystemImageButtonArray={
      title: "component.systemImageModel.Button",
      icon: "icon-background",
      opened: false,
      items:[]
    }
    const SystemImageSafeArray={
      title: "component.systemImageModel.SafeImage",
      icon: "icon-background",
      opened: false,
      items:[]
    }
//按钮
    SystemImageButtonList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      const file = "/static/ISM/systemImage/button/"+fileName
      SystemImageButtonArray.items.push(file)
    })
    SystemImageSafeList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      const file = "/static/ISM/systemImage/safe/"+fileName
      SystemImageSafeArray.items.push(file)
    })
    const SystemImageWasteWaterTreatmentArray={
      title: "component.systemImageModel.WasteWaterTreatment",
      icon: "icon-background",
      opened: false,
      items:[]
    }
    SystemImageWasteWaterTreatmentList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      const file = "/static/ISM/systemImage/WasteWaterTreatment/"+fileName
      SystemImageWasteWaterTreatmentArray.items.push(file)
    })

    const SystemImageIcontArray={
      title: "component.systemImageModel.Icon",
      icon: "icon-background",
      opened: false,
      items:[]
    }
    SystemImageIconList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      const file = "/static/ISM/systemImage/icon/"+fileName
      SystemImageIcontArray.items.push(file)
    })
    const SystemImagePageBackgroundArray={
      title: "component.systemImageModel.Background",
      icon: "icon-background",
      opened: false,
      items:[]
    }
    const SystemImageSwitchArray={
      title: "component.systemImageModel.Switch",
      icon: "icon-background",
      opened: false,
      items:[]
    }
    SystemImagePageBackgroundList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      const file = "/static/ISM/page/background/"+fileName
      SystemImagePageBackgroundArray.items.push(file)
    })
    SystemImageSwitchList.keys().forEach(filePath => {
      const keyArr = filePath.split('/')
      const fileName = keyArr.pop()
      const file = "/static/images/switch/"+fileName
      SystemImageSwitchArray.items.push(file)
    })

    this.SystemImageList.push(SystemImageButtonArray)
    this.SystemImageList.push(SystemImageSafeArray)
    this.SystemImageList.push(SystemImageWasteWaterTreatmentArray)
    this.SystemImageList.push(SystemImageIcontArray)
    this.SystemImageList.push(SystemImagePageBackgroundArray)
    this.SystemImageList.push(SystemImageSwitchArray)
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

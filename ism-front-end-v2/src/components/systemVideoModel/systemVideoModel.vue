<template>
  <div>
    <a-modal  v-drag-modal  :footer="null" v-model="visible" :width="isMobile?'300px':'1000px'" :title="$t('configComponent.video.ModelTitle')" on-ok="handleOk" >
      <div style="padding: 10px;">
        <div v-if="!isMobile">
          <a-radio-group v-model="radioValue" @change="onChangeGroup">
            <a-radio :value="1">
              {{$t('configComponent.video.CamVideo')}}
            </a-radio>
            <a-radio :value="3">
              {{$t('configComponent.video.NetworkGB28281Video')}}
            </a-radio>
            <a-radio :value="2">
              {{$t('configComponent.video.NetworkVideo')}}
            </a-radio>
            <a-radio :value="4">
              {{$t('configComponent.video.LocalVideo')}}
            </a-radio>
          </a-radio-group>
          <a-divider dashed />
        </div>

        <div v-if="radioValue==1">
          <a-button v-if="!isMobile" @click="getSystemVideoList()"  :disabled="messageShowLoad" type="default">  {{$t('configComponent.video.VideoRefresh')}} </a-button>
          <div >

            <a-table rowKey='Name' :pagination="!isMobile?pagination:paginationPhone" :columns="isMobile?columnsPhone:columns" :data-source="videoList" >
              <template v-for="(item, index) in columns" :slot="item.slotName">
                <span :key="index">{{ $t(item.slotName) }}</span>
              </template>
              <span slot="Status" slot-scope="Status">
                  <div v-if="Status==1" style="color: #74f808">
                    {{$t('configComponent.video.VideoOnline')}}
                  </div>
                 <div v-else-if="Status==0" style="color: #ea1111">
                    {{$t('configComponent.video.VideoOffline')}}
                 </div>
            </span>
              <div slot="Opt" slot-scope="text, record">
                <a type="link"   @click="selectVideo(1,record.key,record.Name)" style="cursor: pointer;color: #13C2C2"><a-icon type="check" /><span style="margin-left: 2px;">{{$t('configComponent.video.SelectVideo')}}</span></a>
              </div>
            </a-table>

          </div>
        </div>

        <div v-if="radioValue==3">
          <a-button v-if="!isMobile" @click="getGB28281VideoList()"  :disabled="messageShowLoad" type="default">  {{$t('configComponent.video.VideoRefresh')}} </a-button>
          <div >
            <a-table rowKey="TableVideoIndex" :pagination="!isMobile?pagination:paginationPhone" :columns="isMobile?GBColumnsPhone:GBColumns" :data-source="GBVideoList" >
              <template v-for="(item, index) in GBColumns" :slot="item.slotName" v-if="!isMobile">
                <span :key="index">{{ $t(item.slotName) }}</span>
              </template>
              <template v-for="(item, index) in GBColumnsPhone" :slot="item.slotName" v-else>
                <span :key="index">{{ $t(item.slotName) }}</span>
              </template>
              <div slot="Opt" slot-scope="text, record">
                <a type="link"   @click="selectVideo(3,record.LiveSubSP,record.Name)" style="cursor: pointer;color: #13C2C2"><a-icon type="check" /><span style="margin-left: 2px;">{{$t('configComponent.video.SelectVideo')}}</span></a>
              </div>
            </a-table>
          </div>
        </div>

        <div v-if="radioValue==2&&!isMobile">
          <a-form  :label-col="{ span: 3}" :wrapper-col="{ span: 20 }">
            <a-form-item :label="$t('configComponent.video.TableVideoName')" >
              <a-input   v-model="networkVideo.name" :default-value="networkVideo.name"/>
            </a-form-item>
            <a-form-item :label="$t('configComponent.video.NetworkVideoUrl')" >
              <a-input   v-model="networkVideo.url" :default-value="networkVideo.url"/>
            </a-form-item>
            <a-form-item>
              <div style="margin-top: 5px;margin-left: 120px">
                <a-button key="submit"  type="primary" @click="selectVideo(0,networkVideo.url,networkVideo.name)">{{$t('component.systemImageModel.networkImageBtn')}}</a-button>
              </div>
            </a-form-item>

          </a-form>
        </div>

        <div v-if="radioValue==4">
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
                <a-list  size="small" :pagination="pagination" :data-source="LocalVideoList">
                  <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
                    <a slot="actions" @click="selectVideo(4,item.imgurl,item.imgurl)">{{$t('component.systemImageModel.selectImage')}}</a>
                    <a-list-item-meta>
                      <span slot="title" >{{ item.name }}</span>
                    </a-list-item-meta>
                  </a-list-item>
                </a-list>
              </div>
            </div>

          </div>
        </div>
      </div>
    </a-modal>
  </div>

</template>
<script>

import {videoList, videoDel, videoAdd, GetMonibucaVideoList} from "@/services/video";
import {mapState} from 'vuex'
import moment from 'moment'
import {systemImageList} from "@/services/systemImages";
import {SYSTEMIMAGEUPLOAD} from "@/services/api";
export default {
  name: 'SystemVideoModel',
  i18n: require('../../i18n/language'),
  data() {
    return {
      videoServer:"",
      isMobile:false,
      AddForm:this.$form.createForm(this),
      radioValue:1,
      GBVideoList:[],
      MonibucaServer:"",
      visible: false,
      addVisible:false,
      messageShowLoad:false,
      VideoServer:location.host,
      uploadUrl:SYSTEMIMAGEUPLOAD,
      videoList:[],
      GBColumns: [
        {
          width: '5%',
          slotName: 'configComponent.video.TableVideoIndex',
          scopedSlots: { customRender: 'TableVideoIndex', title: 'configComponent.video.TableVideoIndex' },
          dataIndex: 'TableVideoIndex',
        },
        {
          width: '10%',
          slotName: 'configComponent.video.TableVideoName',
          scopedSlots: { customRender: 'Name', title: 'configComponent.video.TableVideoName' },
          dataIndex: 'Name',
        },
        {
          slotName: 'configComponent.video.Manufacturer',
          width: '10%',
          scopedSlots: { customRender: 'Manufacturer', title: 'configComponent.video.Manufacturer' },
          dataIndex: 'Manufacturer',
        },
        {
          slotName: 'configComponent.video.Model',
          width: '15%',
          scopedSlots: { customRender: 'Model', title: 'configComponent.video.Model' },
          dataIndex: 'Model',
        },
        {
          slotName: 'configComponent.video.RegisterTime',
          width: '10%',
          scopedSlots: { customRender: 'RegisterTime',title: 'configComponent.video.RegisterTime'},
          dataIndex: 'RegisterTime',
        },
        {
          slotName: 'configComponent.video.UpdateTime',
          width: '10%',
          scopedSlots: { customRender: 'UpdateTime',title: 'configComponent.video.UpdateTime'},
          dataIndex: 'UpdateTime',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '5%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      columns: [
        {
          width: '10%',
          slotName: 'configComponent.video.TableVideoName',
          scopedSlots: { customRender: 'Name', title: 'configComponent.video.TableVideoName' },
          dataIndex: 'Name',
        },
        {
          slotName: 'configComponent.video.TableVideoIP',
          width: '10%',
          scopedSlots: { customRender: 'IP', title: 'configComponent.video.TableVideoIP' },
          dataIndex: 'ip',
        },
        {
          slotName: 'configComponent.video.TableVideoStatus',
          width: '5%',
          scopedSlots: { customRender: 'Status',title: 'configComponent.video.TableVideoStatus'},
          dataIndex: 'status',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '10%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      columnsPhone: [
        {
          width: '15%',
          slotName: 'configComponent.video.TableVideoName',
          scopedSlots: { customRender: 'Name', title: 'configComponent.video.TableVideoName' },
          dataIndex: 'Name',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '10%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      GBColumnsPhone:[
        {
          width: '10%',
          slotName: 'configComponent.video.TableVideoName',
          scopedSlots: { customRender: 'Name', title: 'configComponent.video.TableVideoName' },
          dataIndex: 'Name',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '5%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      LocalVideoList:[],
      paginationPhone:{
        pageSize:5,
        showSizeChanger:true
      },
    };
  },
  props: {
    networkVideo: {
      type: Object,
      required: true
    },
  },
  mounted(){
    this.getSystemVideoList()
  },
  computed: {

  },
  methods: {
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
    onChangeGroup(e){
        if(e.target.value==1)
        {
          this.getSystemVideoList()
        }
        else if(e.target.value==3)
        {
          this.getGB28281VideoList()
        }
        else if(e.target.value==4)
        {
          this.getSystemLocalVideoList()
        }
    },
    AddVideoAction(){
      let _t = this
      this.AddForm.validateFields((err) => {
        if (!err) {
          let params = {
            name:this.AddForm.getFieldValue('name'),
            ip:this.AddForm.getFieldValue('ip'),
            port:parseInt(this.AddForm.getFieldValue('port')),
            user:this.AddForm.getFieldValue('user'),
            password:this.AddForm.getFieldValue('password'),
            url:this.AddForm.getFieldValue('RtspUrl'),
          };
          videoAdd(params).then(function (res){
            if(res.data.code==0)
            {
              _t.getSystemVideoList()

              _t.addVisible = false;
            }else{
              _t.$message.success(_t.$t('alarm.trigger.AddFailed'))
            }
          })
        }
      })
    },
    getSystemVideoList(){
      let _t = this
      _t.videoList=[]
      _t.messageShowLoad=true
      videoList().then(function (res){
        if(res.data.list!=null)
        {
          let tableData={}

          for(let i=0;i<res.data.list.length;i++)
          {
            tableData.key = res.data.list[i].Uuid
            tableData.ip = res.data.list[i].Ip
            tableData.Name = res.data.list[i].Name
            tableData.port = res.data.list[i].Port
            tableData.password = res.data.list[i].Password
            tableData.user = res.data.list[i].User
            tableData.status = res.data.list[i].Status
            tableData.uuid = res.data.list[i].Uuid
            _t.videoList.push(tableData)
            tableData={}
          }
        }
        _t.messageShowLoad=false
      })
    },
    getGB28281VideoList(){
      let _t = this
      _t.messageShowLoad=true
      _t.GBVideoList=[]
      GetMonibucaVideoList().then(function (res){

        if(res.data.code==0)
        {
          _t.MonibucaServer = res.data.MonibucaServer
          let videoList = JSON.parse(res.data.data)
          let IndexNo=0
          for(let i=0;i<videoList.length;i++)
          {
            for(let k = 0;k<videoList[i].Channels.length;k++)
            {
              let tableData={}
              tableData.Manufacturer = videoList[i].Manufacturer
              tableData.RegisterTime = moment(videoList[i].RegisterTime).format('yyyy-MM-DD HH:mm:ss')
              tableData.UpdateTime = moment(videoList[i].UpdateTime).format('yyyy-MM-DD HH:mm:ss')
              tableData.TableVideoIndex=IndexNo
              tableData.Name = videoList[i].Channels[k].Name==""?videoList[i].Name:videoList[i].Channels[k].Name
              tableData.Model = videoList[i].Channels[k].Model
              tableData.Status = videoList[i].Channels[k].Status
              tableData.LiveSubSP = videoList[i].ID+"/"+videoList[i].Channels[k].DeviceID
              IndexNo++
              _t.GBVideoList.push(tableData)
              tableData={}
            }
          }
        }
        _t.messageShowLoad=false
      }).catch(function (error) {
        _t.messageShowLoad=false
      });
    },
    getSystemLocalVideoList(){
      let _t = this
      _t.LocalVideoList = []
      systemImageList().then(function (res){
        if(res.data.list!=null)
        {

          _t.imageList=[]
          _t.model3DList=[]
          for(let i=0;i<res.data.list.length;i++)
          {
            let tableData={}
            tableData.imgurl = res.data.list[i].path
            tableData.name = res.data.list[i].name
            tableData.startAt = _t.formatDateTime(res.data.list[i].UpdatedAt)
            if(res.data.list[i].type==4) {
              _t.LocalVideoList.push(tableData)
            }
            tableData={}
          }
        }

      })
    },
    showModal() {
      this.visible = true
    },
    selectVideo(type,key,Name){
      let url=""

      if(type==1)
      {
         url = "webrtc://"+this.videoServer+"/webrtcstream/"+key
      }
      else if(type==3)
      {
        url=this.MonibucaServer+"hdl/"+key+".flv"
      }
      else
      {
         url = key
      }

      this.$emit("onSelectVideo", {type:type,key:key,url:url,name:Name});
      this.visible =false
    },
    delVideo(name){
      let _t = this
      const params={
        uuid:name
      }
      videoDel(params).then(function (res){
        if(res.data.code==0)
        {
          _t.getSystemVideoList()
          _t.$message.success(_t.$t('component.systemImageModel.delImageSuccess'))
        }
        else
        {
          _t.$message.error(_t.$t('component.systemImageModel.delImageFailed'))
        }
      })
    },
    handleCancel() {
      this.visible = false;
    },
    afterUpload(info) {
      if (info.file.status === 'done') {
        let result = info.file.response
        if(result.Code==2002) {
          this.getSystemLocalVideoList()
          this.$message.success(this.$t('component.systemImageModel.uploadSuccess'))
        }
        else
        {
          this.$message.error(this.$t('component.systemImageModel.uploadFailed'))
        }
      }
    },
  },
  created(){
    this.$EventBus.$on('cell-vuex', (data) => {
      this.videoServer = data.PStore.state.setting.videoServer
      this.isMobile = data.PStore.state.setting.isMobile
    })
  }
}
</script>

<style lang="less" scoped>
::v-deep .ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
::v-deep .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  /*background: #f8f8f8;*/
  /*border-bottom: 1px solid #e8e8e8;*/
  transition: background .3s ease;
}
</style>

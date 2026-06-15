<template>
  <div>
    <a-modal  v-drag-modal  :footer="null" v-model="visible" :width="isMobile?'300px':'1000px'" :title="$t('configComponent.video.SnapHistoryVideoList')" on-ok="handleOk" >
      <div style="padding: 10px;">
        <a-form layout="inline" :label-col="{ span: 8}" :wrapper-col="{ span: 16 }">
          <a-form-item :label="$t('systemHistoryVideoModel.deviceList')" >
            <a-select v-model="selectDevice" @change="changeDeviceEvent" style="width: 200px" :dropdownStyle="{'z-index': 9999999}">
              <a-select-option  v-for="(deviceValue,index) in deviceList" :key="index" :value="deviceValue">
                {{ deviceValue }}
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item :label="$t('systemHistoryVideoModel.DateList')" >
            <a-select v-model="selectDate" @change="changeDeviceDateEvent" style="width: 200px" :dropdownStyle="{'z-index': 9999999}">
              <a-select-option  v-for="(dateValue,index) in deviceDateList" :key="index" :value="dateValue">
                {{ dateValue }}
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label=" ">
            <a-button v-if="!isMobile" @click="getSystemSnapHistoryVideoList"  :disabled="messageShowLoad" type="default">  {{$t('configComponent.video.VideoRefresh')}} </a-button>
          </a-form-item>

        </a-form>
        <div>
            <a-table rowKey='key' :pagination="pagination" :columns="columns" :data-source="deviceFileList" >
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
                <a type="link"   @click="selectVideo(1,record.key,record.Path)" style="cursor: pointer;color: #13C2C2"><a-icon type="check" /><span style="margin-left: 2px;">{{$t('configComponent.video.SelectHistoryVideo')}}</span></a>
              </div>
            </a-table>
        </div>
      </div>
    </a-modal>
  </div>

</template>
<script>

import {
  videoList,
  videoDel,
  videoAdd,
  GetMonibucaVideoList,
  historyListVideoList,
  snapHistoryListVideoList
} from "@/services/video";
import {mapState} from 'vuex'
import moment from 'moment'
export default {
  name: 'SystemSnapVideoModel',
  i18n: require('../../i18n/language'),
  data() {
    return {
      videoServer:"",
      isMobile:false,
      AddForm:this.$form.createForm(this),
      radioValue:1,
      deviceList:[],
      selectDevice:"",
      deviceDateList:[],
      deviceFileList:[],
      selectDate:"",
      GBVideoList:[],
      MonibucaServer:location.hostname,
      visible: false,
      addVisible:false,
      messageShowLoad:false,
      VideoServer:location.host,
      videoList:[],
      columns: [
        {
          width: '30%',
          slotName: 'configComponent.video.TableVideoName',
          scopedSlots: { customRender: 'Name', title: 'configComponent.video.TableVideoName' },
          dataIndex: 'Name',
        },
        {
          slotName: 'configComponent.video.TableVideoSize',
          width: '10%',
          scopedSlots: { customRender: 'Size',title: 'configComponent.video.TableVideoSize'},
          dataIndex: 'Size',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '10%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      pagination:{
        pageSize:5,
        showSizeChanger:true
      },
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
    this.getSystemSnapHistoryVideoList()
  },
  computed: {

  },
  methods: {
    changeDeviceEvent(device){
      this.selectDevice = device
      this.selectDate =""
      let tempDate = []
      for(let i=0;i<this.videoList.length;i++)
      {
        if(this.videoList[i].Device==device)
        {
          tempDate.push(this.videoList[i].Date)
        }
      }
      this.deviceDateList = tempDate.filter((value, index, array) => array.indexOf(value) === index);
    },
    changeDeviceDateEvent(date){
      this.selectDate = date
      this.deviceFileList=[]
      for(let i=0;i<this.videoList.length;i++)
      {
        if(this.videoList[i].Device==this.selectDevice&&(this.videoList[i].Date==date))
        {
          this.deviceFileList.push(this.videoList[i])
        }
      }
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
    getSystemSnapHistoryVideoList(){
      let _t = this
      _t.videoList=[]
      this.selectDevice=""
      this.selectDate=""
      this.deviceList=[]
      this.deviceDateList=[]
      this.deviceFileList=[]
      _t.messageShowLoad=true
      snapHistoryListVideoList().then(function (res){
        if(res.data.files!=null)
        {
          let tableData={}
          let list = res.data.files
          let deviceListTemp = []
          for(let i=0;i<list.length;i++)
          {
            let splitArray =  list[i].Path.split("/");
            tableData.key = list[i].Path
            tableData.Name = splitArray[4]
            tableData.Device = splitArray[2]
            tableData.Date = splitArray[3]
            tableData.FileName = splitArray[4]
            tableData.Path = list[i].Path
            tableData.Size = (list[i].Size/(1024*1024)).toFixed(2)
            _t.videoList.push(tableData)
            deviceListTemp.push(tableData.Device)
            tableData={}
          }
          _t.deviceList = deviceListTemp.filter((value, index, array) => array.indexOf(value) === index);
          if(_t.deviceList.length>0)
          {
            _t.changeDeviceEvent(_t.deviceList[0])
            if(_t.deviceDateList.length>0)
            {
              _t.changeDeviceDateEvent(_t.deviceDateList[0])
            }
          }
        }
        _t.messageShowLoad=false
      })
    },
    showModal() {
      this.visible = true
    },
    selectVideo(type,key,Name){
      this.$emit("onSelectVideo", {type:type,key:key,url:Name,name:Name});
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

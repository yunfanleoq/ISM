<template>
  <a-card>
    <a-button @click="addVisible=true" type="primary">  {{$t('configComponent.video.AddVideo')}} </a-button>
    <a-divider type="vertical" />
    <a-button @click="getSystemVideoList()" :disabled="messageShowLoad" type="default" >  {{$t('configComponent.video.VideoRefresh')}} </a-button>
    <div style="margin-top: 5px">
      <a-table rowKey="Name" :pagination="pagination" :columns="columns" :data-source="videoList" >
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
          <a type="link"   @click="goToEdit(record)" style="cursor: pointer;color: #13C2C2"><a-icon type="edit" /><span style="margin-left: 2px;">{{$t('monitor.nodeEdit')}}</span></a> |
          <a-popconfirm
              :title="$t('configComponent.video.DelVideoConfirm')"
              :ok-text="$t('component.systemImageModel.delImageYes')"
              :cancel-text="$t('component.systemImageModel.delImageNo')"
              @confirm="delVideo(record.uuid)"
          >
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('configComponent.video.DelVideo')}}</a> |
          </a-popconfirm>
          <a-popconfirm v-if="!record.isUsed"
              :title="$t('monitor.NotUseTips')"
              :ok-text="$t('component.systemImageModel.delImageYes')"
              :cancel-text="$t('component.systemImageModel.delImageNo')"
              @confirm="StopOrStartVideo(1,record.uuid)"
          >
            <a  type="link"   style="cursor: pointer;color: #F4A460"><a-icon type="stop" /><span style="margin-left: 2px;">{{$t('monitor.NotUse')}}</span></a>
          </a-popconfirm>
          <a v-if="record.isUsed" type="link"  @click="StopOrStartVideo(0,record.uuid)" style="cursor: pointer;color: #00FF00"><a-icon type="play-circle" /><span style="margin-left: 2px;">{{$t('monitor.HadUse')}}</span></a>
        </div>
      </a-table>
    </div>

    <a-modal v-model="addVisible" :title="$t('configComponent.video.AddVideo')" @ok="AddVideoAction">
      <a-form :form="AddForm" :label-col="{ span: 7 }" :wrapper-col="{ span: 12 }">
        <a-form-item
            :label="$t('configComponent.video.TableVideoName')"
        >
          <a-input
              v-decorator="[
                'name',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoIP')"
        >
          <a-input
              v-decorator="[
                'ip',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoPort')"
        >
          <a-input
              v-decorator="[
                'port',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.RtspUrl')"
        >
          <a-input
              v-decorator="[
                'RtspUrl',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoUser')"
        >
          <a-input
              v-decorator="[
                'user',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoPassword')"
        >
          <a-input type="password"
              v-decorator="[
                'password',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('configComponent.video.IsRecord')"
        >
          <a-checkbox
                   v-decorator="[
                'IsRecord',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('configComponent.video.RecordInter')"
        >
          <a-input
              v-decorator="[
                'RecordInter',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>

      </a-form>
    </a-modal>

    <a-modal v-model="editVisible"  :title="$t('configComponent.video.AddVideo')" @ok="editVideoAction">
      <a-form :form="EditForm" :label-col="{ span: 7 }" :wrapper-col="{ span: 12 }">
        <a-form-item
            :label="$t('configComponent.video.TableVideoName')"
        >
          <a-input
              v-decorator="[
                'name',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoIP')"
        >
          <a-input
              v-decorator="[
                'ip',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoPort')"
        >
          <a-input
              v-decorator="[
                'port',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.RtspUrl')"
        >
          <a-input
              v-decorator="[
                'RtspUrl',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoUser')"
        >
          <a-input
              v-decorator="[
                'user',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('configComponent.video.TableVideoPassword')"
        >
          <a-input type="password"
              v-decorator="[
                'password',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>


        <a-form-item
            :label="$t('configComponent.video.IsRecord')"
        >
          <a-checkbox v-model="IsRecord"
              v-decorator="[
                'IsRecord',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('configComponent.video.RecordInter')"
        >
          <a-input
              v-decorator="[
                'RecordInter',
                { rules: [{ required: false}] },
              ]"
          />
        </a-form-item>

      </a-form>
    </a-modal>
  </a-card>
</template>

<script>
import {videoList, videoDel, videoAdd, videoEdit, VideoStopOrStart} from "@/services/video";
import {mapState} from 'vuex'
export default {
  name: 'VideoManager',
  i18n: require('../../i18n/language'),
  components: {

  },
  data () {
    return {
      editVisible:false,
      AddForm:this.$form.createForm(this),
      EditForm:this.$form.createForm(this),
      radioValue:1,
      visible: false,
      addVisible:false,
      messageShowLoad:false,
      videoList:[],
      editUuid:"",
      IsRecord:false,
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
          slotName: 'configComponent.video.TableVideoPort',
          width: '5%',
          scopedSlots: { customRender: 'Port', title: 'configComponent.video.TableVideoPort' },
          dataIndex: 'port',
        },
        {
          slotName: 'configComponent.video.TableVideoUser',
          width: '10%',
          scopedSlots: { customRender: 'User',title: 'configComponent.video.TableVideoUser'},
          dataIndex: 'user',
        },
        {
          slotName: 'configComponent.video.TableVideoPassword',
          width: '10%',
          scopedSlots: { customRender: 'Password',title: 'configComponent.video.TableVideoPassword'},
          dataIndex: 'password',
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
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
    }
  },
  computed: {
    ...mapState('setting', ['videoServer']),
  },
  authorize: {

  },
  mounted(){
    this.getSystemVideoList()
  },
  activated(){

  },
  created(){
  },
  watch: {

  },
  methods: {
    goToEdit(item){
      this.editUuid = item.uuid
      let _t = this
      this.editVisible = true
      this.IsRecord = item.IsRecord
      this.$message.loading(this.$t("monitor.loading"), 0.5)
      setTimeout(function (){
        _t.EditForm.setFieldsValue(
            {
              name:item.Name,
              ip:item.ip,
              port:item.port.toString(),
              user:item.user,
              password:item.password,
              RtspUrl:item.url,
              RecordInter:item.RecordInter.toString(),
            })
      },500)
    },
    editVideoAction(){
      let _t = this
      this.EditForm.validateFields((err) => {
        if (!err) {
          let params = {
            uuid:_t.editUuid,
            data: {
              name: this.EditForm.getFieldValue('name'),
              ip: this.EditForm.getFieldValue('ip'),
              port: parseInt(this.EditForm.getFieldValue('port')),
              user: this.EditForm.getFieldValue('user'),
              password: this.EditForm.getFieldValue('password'),
              StreamURL: this.EditForm.getFieldValue('RtspUrl'),
              IsRecord:this.EditForm.getFieldValue('IsRecord')?1:2,
              RecordInter:parseInt(this.EditForm.getFieldValue('RecordInter')),
            }
          };
          videoEdit(params).then(function (res){
            if(res.data.code==0)
            {
              _t.getSystemVideoList()
              _t.$message.success(_t.$t('alarm.trigger.EditSuccess'))
              _t.editVisible = false;
            }else{
              _t.$message.success(_t.$t('VideoManager.EditFailed'))
            }
          })
        }
      })
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
            IsRecord:this.AddForm.getFieldValue('IsRecord')?1:2,
            RecordInter:parseInt(this.AddForm.getFieldValue('RecordInter')),
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
      _t.messageShowLoad=true
      _t.videoList=[]
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
            tableData.url = res.data.list[i].StreamURL
            tableData.isUsed = res.data.list[i].IsUsed
            tableData.uuid = res.data.list[i].Uuid
            tableData.IsRecord = res.data.list[i].IsRecord==1?true:false
            tableData.RecordInter = res.data.list[i].RecordInter
            _t.videoList.push(tableData)
            tableData={}
          }
        }
        _t.messageShowLoad=false
      }).catch(function (error) {
        _t.messageShowLoad=false
      });
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
    StopOrStartVideo(value,uuid){
      let _t = this
      let params = {
        uuid:uuid,
        data: {
          IsUsed:value
        }
      };
      VideoStopOrStart(params).then(function (res){
        if(res.data.code==0)
        {
          _t.getSystemVideoList()
          _t.$message.success(_t.$t('monitor.OptResultSuccess'))
        }else{
          _t.$message.success(_t.$t('monitor.OptResultFailed'))
        }
      })
    },
    handleCancel() {
      this.visible = false;
    },
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
::v-deep .ant-form-item {
  margin-bottom: 5px;

}
</style>

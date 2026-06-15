<template>
  <div>
    <a-card>
      <a-spin :spinning="messageShowLoad" >
        <a-tabs default-active-key="1" @change="callback">
          <a-tab-pane key="1">
           <span slot="tab" style="font-size: 20px">
              {{$t('ISMNode.General')}}
            </span>
            <a-form v-if="tabKey==1" :form="GeneralFrom" layout="vertical" style="margin: 0;">
              <a-form-item :label-col="{ span: 4}" :wrapper-col="wrapperCol" :label="$t('ISMNode.NodeName')">
                <a-input  autocomplete="autocomplete"
                          v-decorator="['ISMNodeName', {rules: [{ required: true, message: $t('ISMNode.NodeName'), whitespace: true}]}]"
                />

              </a-form-item>

              <a-form-item :label-col="{ span: 4}" :wrapper-col="wrapperCol" :label="$t('ISMNode.NodePort')">
                <a-input  autocomplete="autocomplete"
                          v-decorator="['ISMNodePort', {rules: [{ required: true, message: $t('ISMNode.NodePort'), whitespace: true}]}]"
                />

              </a-form-item>
              <a-divider orientation="left">{{$t('ISMNode.NodeInConfig')}}</a-divider>
              <a-form-item :label-col="{ span: 4}" :wrapper-col="wrapperCol" :label="$t('ISMNode.pingHeart')">
                <a-input  autocomplete="autocomplete"
                          v-decorator="['pingHeart', {rules: [{ required: true, message: $t('ISMNode.pingHeart'), whitespace: true}]}]"
                />

              </a-form-item>
              <a-form-item :label-col="{ span: 4}" :wrapper-col="wrapperCol" :label="$t('ISMNode.pingOutTime')">
                <a-input  autocomplete="autocomplete"
                          v-decorator="['pingOutTime', {rules: [{ required: true, message: $t('ISMNode.pingOutTime'), whitespace: true}]}]"
                />

              </a-form-item>
              <a-form-item :label-col="{ span: 4}" :wrapper-col="wrapperCol" :label="$t('ISMNode.pingOutTimeCount')">
                <a-input  autocomplete="autocomplete"
                          v-decorator="['pingOutTimeCount', {rules: [{ required: true, message: $t('ISMNode.pingOutTimeCount'), whitespace: true}]}]"
                />

              </a-form-item>
              <a-form-item style="margin-left: 100px" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
                <a-button type="primary" @click="SetNodeConfig">
                  {{$t('AlarmTips.Save')}}
                </a-button>
              </a-form-item>

            </a-form>
          </a-tab-pane>
          <a-tab-pane key="2" force-render>
          <span slot="tab" style="font-size: 20px">
             {{$t('ISMNode.OutConnect')}}
            </span>
            <div style="padding: 5px" v-if="tabKey==2">
              <a-space class="operator">
                <a-button @click="ModelVisible=true" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
                <a-button @click="GetConnectOut"  type="default" icon="sync" :loading="messageShowLoad">{{$t("dataModel.refModel")}}</a-button>
              </a-space>
              <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
                <a-table :pagination="pagination" :columns="columns" :data-source="dataSource" rowKey="OutConnectName" >
                  <template v-for="(item, index) in columns" :slot="item.slotName">
                    <span :key="index">{{ $t(item.slotName) }}</span>
                  </template>
                  <div slot="ConnectStatus" slot-scope="text,record">
                    <span v-if="record.IsEnable==2" style="margin-left: 2px;">{{$t('NetWork.NodeWork.Disable')}}</span>
                    <span v-else-if="text==2" style="margin-left: 2px;">{{$t('NetWork.NodeWork.ConnFailed')}}</span>
                    <span v-else-if="text==3" style="margin-left: 2px;">{{$t('NetWork.NodeWork.ConnSuccess')}}</span>
                    <span v-else style="margin-left: 2px;">{{$t('NetWork.NodeWork.ConnIng')}}</span>
                  </div>
                  <div slot="action" slot-scope="text, record">
                    <a type="link"   @click="gotoEditConnect(record)" style="cursor: pointer;color: #13C2C2"><a-icon type="edit" /><span style="margin-left: 2px;">{{$t('dataModel.modelDetail')}}</span></a> |
                    <a v-if="record.IsEnable==1" type="link"  @click="OptConnectOut(2,record)" style="cursor: pointer;color: #0af4d1"><a-icon type="stop" /><span style="margin-left: 2px;">{{$t('monitor.NotUse')}}</span></a>
                    <a v-else type="link"  @click="OptConnectOut(1,record)" style="cursor: pointer;color: #f10742"><a-icon type="play-circle" /><span style="margin-left: 2px;">{{$t('monitor.HadUse')}}</span></a> |
                    <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="DelConnectOut(record)">
                      <a-icon slot="icon" type="question-circle-o" style="color: red" />
                      <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
                    </a-popconfirm>
                  </div>
                </a-table>
              </a-spin>
            </div>
          </a-tab-pane>
          <a-tab-pane key="3">
           <span slot="tab" style="font-size: 20px">
              {{$t('ISMNode.inConnect')}}
            </span>
            <div style="padding: 5px" v-if="tabKey==3">
              <a-space class="operator">
                <a-button @click="GetConnectIn"  type="default" icon="sync" :loading="messageShowLoad">{{$t("dataModel.refModel")}}</a-button>
              </a-space>
              <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
                <a-table :pagination="pagination" :columns="Incolumns" :data-source="ConnectIn" rowKey="OutConnectName" >
                  <template v-for="(item, index) in Incolumns" :slot="item.slotName">
                    <span :key="index">{{ $t(item.slotName) }}</span>
                  </template>
                  <div slot="ConnectStatus" slot-scope="text,record">
                    <span v-if="record.IsEnable==2" style="margin-left: 2px;">{{$t('NetWork.NodeWork.Disable')}}</span>
                    <span v-else-if="text==2" style="margin-left: 2px;">{{$t('NetWork.NodeWork.ConnFailed')}}</span>
                    <span v-else-if="text==3" style="margin-left: 2px;">{{$t('NetWork.NodeWork.ConnSuccess')}}</span>
                    <span v-else style="margin-left: 2px;">{{$t('NetWork.NodeWork.ConnIng')}}</span>
                  </div>
                  <div slot="action" slot-scope="text, record">
                    <a type="link"   @click="gotoEditConnect(record)" style="cursor: pointer;color: #13C2C2"><a-icon type="edit" /><span style="margin-left: 2px;">{{$t('dataModel.modelDetail')}}</span></a> |
                    <a v-if="record.IsEnable==1" type="link"  @click="OptConnectOut(2,record)" style="cursor: pointer;color: #0af4d1"><a-icon type="stop" /><span style="margin-left: 2px;">{{$t('monitor.NotUse')}}</span></a>
                    <a v-else type="link"  @click="OptConnectOut(1,record)" style="cursor: pointer;color: #f10742"><a-icon type="play-circle" /><span style="margin-left: 2px;">{{$t('monitor.HadUse')}}</span></a> |
                    <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="DelConnectOut(record)">
                      <a-icon slot="icon" type="question-circle-o" style="color: red" />
                      <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
                    </a-popconfirm>
                  </div>
                </a-table>
              </a-spin>
            </div>
          </a-tab-pane>
        </a-tabs>
      </a-spin>
    </a-card>
    <a-modal v-model="ModelVisible" :title="$t('ISMNode.OutConnect')" @ok="AddConnectOut">
      <a-form :form="form" :label-col="{ span: 7 }" :wrapper-col="{ span: 15 }">
      <a-form-item
          :label="$t('ISMNode.OutConnectName')"
      >
        <a-input  autocomplete="autocomplete"

                  v-decorator="['OutConnectName', {rules: [{ required: true, message: $t('ISMNode.OutConnectName'), whitespace: true}]}]"
        />
      </a-form-item>

      <a-form-item
          :label="$t('ISMNode.IpAddress')"
      >
        <a-input  autocomplete="autocomplete"

                  v-decorator="['IpAddress', {rules: [{ required: true, message: $t('ISMNode.IpAddress'), whitespace: true}]}]"
        />
      </a-form-item>

        <a-form-item
            :label="$t('ISMNode.ConnectPort')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['ConnectPort', { initialValue:'8066',rules: [{ required: true, message: $t('ISMNode.ConnectPort'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.pingHeart')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['pingHeart', {initialValue:'1000',rules: [{ required: true, message: $t('ISMNode.pingHeart'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.pingOutTime')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['pingOutTime', {initialValue:'3000',rules: [{ required: true, message: $t('ISMNode.pingOutTime'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.pingOutTimeCount')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['pingOutTimeCount', {initialValue:'3',rules: [{ required: true, message: $t('ISMNode.pingOutTimeCount'), whitespace: true}]}]"
          />
        </a-form-item>


    </a-form>
    </a-modal>

    <a-modal v-model="ModelEditVisible" :title="$t('ISMNode.OutConnect')" @ok="EditConnectOut">
      <a-form :form="editform" :label-col="{ span: 7 }" :wrapper-col="{ span: 15 }">
        <a-form-item
            :label="$t('ISMNode.OutConnectName')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['OutConnectName', {rules: [{ required: true, message: $t('ISMNode.OutConnectName'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.IpAddress')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['IpAddress', {rules: [{ required: true, message: $t('ISMNode.IpAddress'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.ConnectPort')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['ConnectPort', { initialValue:'8066',rules: [{ required: true, message: $t('ISMNode.ConnectPort'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.pingHeart')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['pingHeart', {initialValue:'1000',rules: [{ required: true, message: $t('ISMNode.pingHeart'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.pingOutTime')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['pingOutTime', {initialValue:'3000',rules: [{ required: true, message: $t('ISMNode.pingOutTime'), whitespace: true}]}]"
          />
        </a-form-item>

        <a-form-item
            :label="$t('ISMNode.pingOutTimeCount')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['pingOutTimeCount', {initialValue:'3',rules: [{ required: true, message: $t('ISMNode.pingOutTimeCount'), whitespace: true}]}]"
          />
        </a-form-item>


      </a-form>
    </a-modal>
  </div>
</template>

<script>
import {
  GetNodeConfig,
  SetNodeConfig,
  GetConnectOut,
  AddOutConnect,
  EditOutConnect,
  DelOutConnect,
  OptOutConnect,
  GetConnectIn
} from "@/services/ismnetwork";
import moment from 'moment'
export default {
  name: "ISMNode",
  i18n: require('@/i18n/language'),
  data() {
    return {
      labelCol: { span: 2 },
      tabKey:"1",
      ModelVisible:false,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      ConnectIn:[],
      dataSource: [],
      EditUUid:"",
      ModelEditVisible:false,
      columns: [
        {
          width: '15%',
          slotName: 'ISMNode.OutConnectName',
          scopedSlots: { customRender: 'OutConnectName', title: 'ISMNode.OutConnectName' },
          dataIndex: 'OutConnectName',
        },
        {
          width: '10%',
          slotName: 'ISMNode.IpAddress',
          scopedSlots: { customRender: 'IpAddress', title: 'ISMNode.IpAddress' },
          dataIndex: 'IpAddress',
        },
        {
          slotName: 'ISMNode.ConnectPort',
          width: '10%',
          scopedSlots: { customRender: 'ConnectPort',title: 'ISMNode.ConnectPort'},
          dataIndex: 'ConnectPort',
        },
        {
          slotName: 'ISMNode.ConnectStatus',
          width: '10%',
          scopedSlots: { customRender: 'ConnectStatus',title: 'ISMNode.ConnectStatus'},
          dataIndex: 'ConnectStatus',
        },
        {
          slotName: 'ISMNode.CreateTime',
          width: '15%',
          scopedSlots: { customRender: 'UpdatedAt',title: 'ISMNode.CreateTime'},
          dataIndex: 'UpdatedAt',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          width: '15%',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      Incolumns: [
        {
          width: '15%',
          slotName: 'ISMNode.OutConnectName',
          scopedSlots: { customRender: 'OutConnectName', title: 'ISMNode.OutConnectName' },
          dataIndex: 'NodeName',
        },
        {
          width: '10%',
          slotName: 'ISMNode.IpAddress',
          scopedSlots: { customRender: 'IpAddress', title: 'ISMNode.IpAddress' },
          dataIndex: 'NodeIpAddress',
        },
      ],
      messageShowLoad:false,
      wrapperCol: { span: 7 },
      editform:this.$form.createForm(this),
      GeneralFrom: this.$form.createForm(this),
      form: this.$form.createForm(this),
    };
  },
  components: {

  },
  created(){
    if(this.tabKey==1) {
      this.GetNodeConfig()
    }else if(this.tabKey==2) {
      this.GetConnectOut()
    }
  },
  methods: {
    GetNodeConfig(){
      let _t = this
      this.messageShowLoad=true
      _t.$message.destroy()
      GetNodeConfig().then(function (res){
        if(res.data.code==0)
        {
          let NodeInfo = res.data.data
          _t.GeneralFrom.setFieldsValue({
            ISMNodeName:NodeInfo.NodeName,
            ISMNodePort:NodeInfo.NodePort.toString(),
            pingHeart:NodeInfo.NodePing.toString(),
            pingOutTime:NodeInfo.PingOutTime.toString(),
            pingOutTimeCount:NodeInfo.PingOutTimeCount.toString(),
          })
        }
        else
        {
          _t.$message.success(_t.$t('ISMNode.Success'), 3)
        }
        _t.messageShowLoad=false
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('ISMNode.Failed'), 3)
      })
    },
    SetNodeConfig(){
      let _t = this

      this.GeneralFrom.validateFields((err) => {
        if (!err) {
          let SetData={
            NodeName:this.GeneralFrom.getFieldValue('ISMNodeName'),
            NodePort:parseInt(this.GeneralFrom.getFieldValue('ISMNodePort')),
            NodePing:parseInt(this.GeneralFrom.getFieldValue('pingHeart')),
            PingOutTime:parseInt(this.GeneralFrom.getFieldValue('pingOutTime')),
            PingOutTimeCount:parseInt(this.GeneralFrom.getFieldValue('pingOutTimeCount')),
          }
          this.messageShowLoad=true
          _t.$message.destroy()
          SetNodeConfig(SetData).then(function (res){
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t('ISMNode.Success'), 3)
            }
            else
            {
              _t.$message.error(_t.$t('ISMNode.Failed'), 3)
            }
            _t.messageShowLoad=false
          }).catch(function(e){
            _t.messageShowLoad=false
            _t.$message.error(_t.$t('ISMNode.Failed'), 3)
          })
        }

      })
    },
    GetConnectOut(){
      let _t = this
      _t.dataSource=[]
      GetConnectOut().then(function (res){
        if(res.data.code == 0)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            res.data.list[i].CreatedAt = moment(res.data.list[i].CreatedAt).format('YYYY-MM-DD HH:mm:ss');
            res.data.list[i].UpdatedAt = moment(res.data.list[i].UpdatedAt).format('YYYY-MM-DD HH:mm:ss');
          }
          _t.dataSource = res.data.list
        }
        _t.messageShowLoad=false
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    gotoEditConnect(item){
      this.ModelEditVisible=true
      let _t = this
      this.EditUUid = item.uuid
      setTimeout(function (){
        _t.editform.setFieldsValue(
            {
              OutConnectName:item.OutConnectName,
              IpAddress:item.IpAddress,
              ConnectPort:item.ConnectPort.toString(),
              pingOutTime:item.PingOutTime.toString(),
              pingHeart:item.pingHeart.toString(),
              pingOutTimeCount:item.PingOutTimeCount.toString(),
            })
      },500)
    },
    DelConnectOut(item){
      let _t = this
      let param = {
        Uuid:item.uuid
      }

      DelOutConnect(param).then(function (res){
        if(res.data.code == 0)
        {
          _t.GetConnectOut()
        }
        _t.messageShowLoad=false
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    EditConnectOut(){
      let _t = this

      this.editform.validateFields((err) => {
        if (!err) {
          let SetData={
            OutConnectName:this.editform.getFieldValue('OutConnectName').toString(),
            IpAddress:this.editform.getFieldValue('IpAddress').toString(),
            ConnectPort:parseInt(this.editform.getFieldValue('ConnectPort')),
            PingOutTime:parseInt(this.editform.getFieldValue('pingOutTime')),
            pingHeart:parseInt(this.editform.getFieldValue('pingHeart')),
            PingOutTimeCount:parseInt(this.editform.getFieldValue('pingOutTimeCount')),
          }
          let Params = {
            Uuid:this.EditUUid,
            editData:SetData,
          }
          this.messageShowLoad=true
          EditOutConnect(Params).then(function (res){
            if(res.data.code==2002)
            {
              _t.ModelVisible = false
              _t.GetConnectOut()
              _t.ModelEditVisible=false
            }
            else
            {
              _t.$message.error(_t.$t('ISMNode.Failed'), 3)
            }
            _t.messageShowLoad=false
          }).catch(function(e){
            _t.messageShowLoad=false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }

      })
    },
    AddConnectOut(){
      let _t = this

      this.form.validateFields((err) => {
        if (!err) {
          let SetData={
            OutConnectName:this.form.getFieldValue('OutConnectName').toString(),
            IpAddress:this.form.getFieldValue('IpAddress').toString(),
            ConnectPort:parseInt(this.form.getFieldValue('ConnectPort')),
            PingOutTime:parseInt(this.form.getFieldValue('pingOutTime')),
            pingHeart:parseInt(this.form.getFieldValue('pingHeart')),
            PingOutTimeCount:parseInt(this.form.getFieldValue('pingOutTimeCount')),
          }
          this.messageShowLoad=true
          AddOutConnect(SetData).then(function (res){
            if(res.data.code==2002)
            {
              _t.ModelVisible = false
              _t.GetConnectOut()
            }
            else
            {
              _t.$message.error(_t.$t('ISMNode.Failed'), 3)
            }
            _t.messageShowLoad=false
          }).catch(function(e){
            _t.messageShowLoad=false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }

      })
    },
    OptConnectOut(st,item){
      let _t = this
      let Params = {
        Uuid:item.uuid,
        Status:st,
        Name:item.OutConnectName
      }
      this.messageShowLoad=true
      OptOutConnect(Params).then(function (res){
        if(res.data.code==2002)
        {
          _t.ModelVisible = false
          _t.GetConnectOut()
          _t.ModelEditVisible=false
        }
        else
        {
          _t.$message.error(_t.$t('ISMNode.Failed'), 3)
        }
        _t.messageShowLoad=false
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    callback(key) {
      this.tabKey = key
      if(key==1)
      {
        this.GetNodeConfig()
      }
      else if(key == 2)
      {
        this.GetConnectOut()
      }
      else if(key == 3)
      {
        this.GetConnectIn()
      }
    },
    GetConnectIn(){
      let _t = this
      this.messageShowLoad=true
      _t.$message.destroy()
      this.ConnectIn=[]
      GetConnectIn().then(function (res){
        if(res.data.code==0)
        {
          _t.ConnectIn = res.data.list
        }
        else
        {
          _t.$message.success(_t.$t('ISMNode.Success'), 3)
        }
        _t.messageShowLoad=false
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('ISMNode.Failed'), 3)
      })
    },
  },
}
</script>

<style scoped>
.ant-form-item{
  margin-bottom: 2px;
}
.operator{
  margin-bottom: 18px;
}
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

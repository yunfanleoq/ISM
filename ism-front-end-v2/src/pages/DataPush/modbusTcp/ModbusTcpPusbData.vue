<template>
  <div>
    <!--    寄存器组表格-->
    <a-card  style="min-height: 400px">
      <a-space class="operator">

        <a-button type="primary" @click="RegisterVisible=true;isEdit=false"> <a-icon type="plus" />
          {{$t("dataModel.DLT645Model.AddData")}}</a-button>



        <a-button type="link" @click="handleExport"> <a-icon type="export" />{{$t('dataModel.export')}}</a-button>

        <a-upload
            name="file"
            :multiple="false"
            :action=localUpgradeUrl
            :showUploadList="false"
            :beforeUpload="beforeUpload"
            @change="localUpgradeCharge"
        >
          <a-button type="link"> <a-icon type="import" />
            {{$t('dataModel.import')}}
          </a-button>
        </a-upload>

        <a-button type="link" @click="onBlackCLK()"> <a-icon type="backward" />
          {{$t("dataModel.opcuaModel.Back")}}</a-button>
      </a-space>
      <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
        <a-table  :pagination="pagination" rowKey="name" :columns="registerGroupColumns" :data-source="registerGroupDataSource" class="ant-table-tbody">
          <template v-for="(item, index) in registerGroupColumns" :slot="item.slotName">
            <span :key="index">{{ $t(item.slotName) }}</span>
          </template>
          <template slot="FunctionCode" slot-scope="text">
            <span v-if="text==1"> 01 {{$t('dataModel.modbusModel.FunctionReadCoils')}}</span>
            <span v-else-if="text==2"> 02 {{$t('dataModel.modbusModel.FunctionReadDisCrete')}}</span>
            <span v-else-if="text==3"> 03 {{$t('dataModel.modbusModel.FunctionReadHoldingRegisters')}}</span>
            <span v-else-if="text==4"> 04 {{$t('dataModel.modbusModel.FunctionReadInputRegisters')}}</span>
          </template>
          <template slot="NodeIDDataType" slot-scope="text">
            <span v-if="text=='Short'"> Signed</span>
            <span v-else-if="text=='Unsigned short'"> Unsigned</span>
            <span v-else-if="text=='Long'"> Long</span>
            <span v-else-if="text=='Float'"> Float</span>
          </template>
          <template slot="action" slot-scope="text, record">
            <div class="editable-row-operations">
              <span >
                <a  @click="() => edit(record)">
                  <a-icon type="edit" /> {{$t('dataModel.opcuaModel.NodeIDEdit')}}</a>
              </span>
              <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.uuid,record.muid)">
                <a-icon slot="icon" type="question-circle-o" style="color: red" />
                <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
              </a-popconfirm>
            </div>
          </template>
        </a-table>
      </a-spin>
    </a-card>
    <!--    添加节点-->

    <a-drawer
        :title="isEdit?$t('dataModel.opcuaModel.EditNodeID'):$t('dataModel.opcuaModel.AddNodeID')"
        :width="720"
        :zIndex="100"
        :visible="RegisterVisible"
        :body-style="{ paddingBottom: '80px' }"
        @close="onClose"
    >
      <a-spin style="padding: 1px;"  :spinning="ShowRegisterLoading" tip="Loading...">
        <a-form :form="RegisterForm" layout="vertical" @submit="AddNodeId">
          <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDName')"
              >
                <a-input autocomplete="autocomplete"
                         v-decorator="['NodeIDName', {rules: [{ required: true, validator: isValidateTxtNonSpec, message: $t('device.deviceNameVal'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.modbusModel.RegisterAddr')"
              >
                <a-input autocomplete="autocomplete"
                         v-decorator="['RegisterAddr', {rules: [{ required: true, message: $t('dataModel.modbusModel.RegisterAddr'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.modbusModel.RegisterFunction')"
              >
                <a-select  class="DataType" autocomplete="autocomplete"
                          v-decorator="['RegisterFunction', {initialValue:'1',rules: [{ required: true, message: $t('dataModel.modbusModel.RegisterFunction'), whitespace: true}]}]"
                >

                  <a-select-option value="1">
                    01 {{$t('dataModel.modbusModel.FunctionReadCoils')}}
                  </a-select-option>
                  <a-select-option value="2">
                    02 {{$t('dataModel.modbusModel.FunctionReadDisCrete')}}
                  </a-select-option>
                  <a-select-option value="3">
                    03 {{$t('dataModel.modbusModel.FunctionReadHoldingRegisters')}}
                  </a-select-option>
                  <a-select-option value="4">
                    04 {{$t('dataModel.modbusModel.FunctionReadInputRegisters')}}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item :label="$t('dataModel.editData.dataType')">
                <a-select   autocomplete="autocomplete"  @change="changeDataType"  v-decorator="[
                  'dataType',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataType') }],
                  },
                ]">
                  <a-select-option value="Short">Signed</a-select-option>
                  <a-select-option value="Unsigned short">Unsigned</a-select-option>
                  <a-select-option value="Long">Long</a-select-option>
                  <a-select-option value="Float">Float</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('ISMDataTemplete.TempleteBandData')"
              >
                <a-input  v-decorator="['TempleteBandData', {rules: [{ required: true, message: $t('ISMDataTemplete.TempleteBandData'), whitespace: true}]}]" >
                  <a-tooltip placement="top" slot="addonAfter">
                    <template slot="title">
                      <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
                    </template>
                    <icon-font @click="ShowDeviceDataModel()" type="icon-xuanzeshuju"  />
                  </a-tooltip>
                </a-input>
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="16">
            <a-col :span="24">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDDec')"
              >
                <a-textarea autocomplete="autocomplete"
                            v-decorator="['NodeIDDec', {rules: [{ required: false, message: $t('dataModel.opcuaModel.NodeIDDec'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
          </a-row>

        </a-form>
      </a-spin>
      <div
          :style="{
              position: 'absolute',
              right: 0,
              bottom: 0,
              width: '100%',
              borderTop: '1px solid #e9e9e9',
              padding: '10px 16px',
              background: '#fff',
              textAlign: 'right',
              zIndex: 1,
            }"
      >

        <a-button v-if="isEdit" key="submit" type="primary" :style="{ marginRight: '8px' }" @click="save">
          {{ $t('component.deviceDataModel.Edit')}}
        </a-button>
        <a-button v-else key="submit" type="primary" :style="{ marginRight: '8px' }" @click="AddNodeId">
          {{ $t('component.deviceDataModel.submit')}}
        </a-button>
        <a-button key="back" @click="RegisterVisible=false">
          {{$t('component.deviceDataModel.cancel')}}
        </a-button>

      </div>
    </a-drawer>
    <device-data-model @onSelectDataModel="onSelectData" ref="deviceDataModel" style="z-index: 999999999"></device-data-model>
  </div>
</template>

<script>
import {
  ModbusTcpDataPushAdd,
  ModbusTcpDataPushDel,
  ModbusTcpDataPushEdit,
  ModbusTcpDataPushList
} from "@/services/ModbusTcpDataPush";
import {exportExcel} from '@/services/excelExport'
import {ImportNodeID, LOCALUPGATEDATAMODBUSTCPDATA} from "@/services/api";
import deviceDataModel from "@/components/deviceDataModel/deviceDataModel.vue";
const dataSource= []
const loadingKey = 'updatable'
export default {
  name: 'ModbusTcpPushData',
  i18n: require('@/i18n/language'),
  data () {
    return {
      dataType:"Short",
      isEdit:false,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      DataRecordType:0,
      importUrl:ImportNodeID+"/"+this.$route.params.uid,
      registerGroupListTable:true,
      error: '',
      NodeIDAccessLevel:[],
      ShowRegisterLoading:false,
      alarmStatus:0,
      recordStatus:0,
      EditForm:this.$form.createForm(this),
      RegisterForm: this.$form.createForm(this),
      selectDataTableUuid:[],
      rowSelection:{
        onSelect:this.onDataTableSelect,
        onSelectAll:this.onDataTableSelectAll
      },
      RegisterVisible:false,
      editVisible:false,
      messageShowLoad:false,
      SelectDataCategory:1,
      DataCategoryYaoTiaoType:3,
      RegisterMessageShowLoad:false,
      json_fields_cn: {
        "数据名称": "name",    //常规字段
        "功能码": {
          field: "FunctionCode",
          //自定义回调函数
          callback: value => {
          switch (value){
            case 1:{
              return "01"
            }
            case 2:{
              return "02"
            }
            case 3:{
              return "03"
            }
            case 4:{
              return "04"
            }
          }
        }
        },
        "寄存器地址": "registerAddress",
        "绑定数据": "BandData",
        "类型": {
          field: "type",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 'Short':{
                return "Signed"
              }
              case 'Unsigned short':{
                return "Unsigned"
              }
              case 'Long':{
                return "Long"
              }
              case 'Float':{
                return "Float"
              }
            }
          }
        },
        "模版数据ID(勿修改)": "uuid"
      },
      json_meta: [
        [
          {
            " key ": " charset ",
            " value ": " utf- 8 "
          }
        ]
      ],
      exportName:"ModbusTCP数据点表模板",
      localUpgradeUrl:LOCALUPGATEDATAMODBUSTCPDATA+"/"+this.$route.params.uid,
      registerGroupColumns: [
        {
          slotName: this.$t("dataModel.opcuaModel.NodeIDName"),
          scopedSlots: {  customRender: 'name' ,title:this.$t("dataModel.opcuaModel.NodeIDName")},
          width: '15%',
          dataIndex: 'name'
        },
        {
          slotName:this.$t("dataModel.modbusModel.RegisterFunction"),
          scopedSlots: {  customRender: 'FunctionCode'  ,title:this.$t("dataModel.modbusModel.RegisterFunction")},
          width: '15%',
          align:"left",
          dataIndex: 'FunctionCode'
        },
        {
          slotName:this.$t("dataModel.modbusModel.RegisterAddr"),
          scopedSlots: {  customRender: 'registerAddress'  ,title:this.$t("dataModel.modbusModel.RegisterAddr")},
          width: '10%',
          align:"left",
          dataIndex: 'registerAddress'
        },
        {
          slotName:this.$t("ISMDataTemplete.TempleteBandData"),
          scopedSlots: {  customRender: 'BandData'  ,title:this.$t("ISMDataTemplete.TempleteBandData")},
          width: '15%',
          align:"left",
          dataIndex: 'BandData'
        },
        {
          slotName:this.$t("dataModel.opcuaModel.NodeIDDataType"),
          scopedSlots: {  customRender: 'NodeIDDataType' ,title:this.$t("dataModel.opcuaModel.NodeIDDataType") },
          width: '15%',
          align:"center",
          dataIndex: 'type',
        },
        {
          title: this.$t('dataModel.modelTableOpt'),
          width: '15%',
          scopedSlots: { customRender: 'action' }
        }
      ],
      dataSource,
      registerGroupDataSource:[],
      selectedRows: [],
    }
  },
  components: {
    deviceDataModel
  },
  created(){
    this.NodeIdList()
    this.registerGroupListTable=true
  },
  activated() {

  },
  mounted() {

  },
  computed: {

  },
  methods: {
    async handleExport(){
      await exportExcel(this.registerGroupDataSource, this.json_fields_cn, this.exportName)
    },
    changeDataType(value){
      this.dataType = value
    },
    ShowDeviceDataModel(){
      this.$refs.deviceDataModel.showDataModal()
    },
    onSelectData(selectData) {
      let data = selectData.DeviceName+"->"+selectData.name
      this.RegisterForm.setFieldsValue(
          {
            TempleteBandData:data,
          })
    },
    isSpec(s) {
      let pattern = /[~!@#$%^&*<>|'-]/gi
      return pattern.test(s)
    },
    isValidateTxtNonSpec (rule, value, callback) {
      if (value != null && value !== '') {
        let numStr = value.charAt(0);
        if ((this.isSpec(value)) || (value.indexOf(' ') !== -1)||(!isNaN(parseFloat(numStr)) && isFinite(numStr))) {
          callback(new Error('不能包含特殊字符或空格'))
        } else {
          callback()
        }
      } else {
        callback()
      }
    },
    chargeDataRecordType(value){
      this.DataRecordType = parseInt(value)
    },
    beforeUpload(){
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'),loadingKey,duration: 0 });
    },
    localUpgradeCharge(info){
      this.dataSource=[]
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy();
        if(result.Code==0) {
          this.$message.success(`${info.file.name} `+this.$t("dataModel.importSuccess"));
          this.NodeIdList()
        }
        else if(result.Code==-2)
        {
          this.$message.error(`${info.file.name} `+this.$t("dataModel.FormatError"));
        }
        else
        {
          this.$message.error(`${info.file.name} `+this.$t("SystemUpgrade.UpgradeFileSaveError"));
        }
      }
      else if (info.file.status === 'uploading') {
        //this.$message.success(`${info.file.name} `+this.$t("SystemUpgrade.BeginUpgradeUploading"));
      }

      this.messageShowLoad = false
    },
    onClose() {
      this.RegisterVisible = false;
    },
    edit(item) {
      let _t = this
      this.isEdit = true
      this.RegisterVisible=true
      this.editingKey = item.uuid
      this.ShowRegisterLoading = true
      setTimeout(function (){
        _t.RegisterForm.setFieldsValue(
            {
              NodeIDName:item.name,
              RegisterAddr:item.registerAddress.toString(),
              RegisterFunction:item.FunctionCode.toString(),
              ByteOrder:item.ByteOrder,
              dataType:item.type,
              TempleteBandData:item.BandData,
              NodeIDDec:item.Description,
            })
        _t.ShowRegisterLoading = false
      },500)
      this.editVisible = true;
    },
    save() {
      this.RegisterForm.validateFields((err) => {
        if (!err) {
          let params = {
            uuid:this.editingKey,
            muid:this.$route.params.uid,
            data: {
              name:this.RegisterForm.getFieldValue('NodeIDName'),
              registerAddress:parseInt(this.RegisterForm.getFieldValue('RegisterAddr')),
              FunctionCode:parseInt(this.RegisterForm.getFieldValue('RegisterFunction')),
              ByteOrder:this.RegisterForm.getFieldValue('ByteOrder'),
              type:this.RegisterForm.getFieldValue('dataType'),
              BandData:this.RegisterForm.getFieldValue('TempleteBandData'),
              Description: this.RegisterForm.getFieldValue('NodeIDDec'),
            }
          }
          let _t = this
          ModbusTcpDataPushEdit(params).then(function (res){
            if(res.data.code==2002)
            {
              const newData = [..._t.registerGroupDataSource];
              const target = newData.filter(item => _t.editingKey === item.uuid)[0];
              if (target) {
                target.name = _t.RegisterForm.getFieldValue('NodeIDName')
                target.registerAddress = _t.RegisterForm.getFieldValue('RegisterAddr')
                target.FunctionCode = _t.RegisterForm.getFieldValue('RegisterFunction')
                target.ByteOrder = _t.RegisterForm.getFieldValue('ByteOrder')
                target.type = _t.RegisterForm.getFieldValue('dataType')
                target.BandData=_t.RegisterForm.getFieldValue('TempleteBandData')
                target.Description=_t.RegisterForm.getFieldValue('NodeIDDec')
                _t.registerGroupDataSource = newData;
              }
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
              _t.RegisterVisible = false;
            }
            else
            {
              _t.$message.error(_t.$t("dataModel.saveFailed"));
            }
          })
        }
      })
    },
    NodeIdList(){
      this.messageShowLoad = true
      const params = {
        muid:this.$route.params.uid,
      }
      this.RegisterVisible = false;
      let _t = this
      _t.registerGroupDataSource = []
      ModbusTcpDataPushList(params).then(function (res){
        _t.messageShowLoad = false
        if(res.data.code==0)
        {
          _t.registerGroupDataSource = res.data.list
        }
      }).catch(function (){
        _t.messageShowLoad = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    AddNodeId(){
      this.RegisterForm.validateFields((err) => {
        if (!err) {
          this.messageShowLoad = true
          const params = {
            muid:this.$route.params.uid,
            name:this.RegisterForm.getFieldValue('NodeIDName'),
            registerAddress:parseInt(this.RegisterForm.getFieldValue('RegisterAddr')),
            FunctionCode:parseInt(this.RegisterForm.getFieldValue('RegisterFunction')),
            ByteOrder:this.RegisterForm.getFieldValue('ByteOrder'),
            type:this.RegisterForm.getFieldValue('dataType'),
            BandData:this.RegisterForm.getFieldValue('TempleteBandData'),
            Description: this.RegisterForm.getFieldValue('NodeIDDec'),
          }
          this.RegisterVisible = false;
          let _t = this
          ModbusTcpDataPushAdd(params).then(function (res){
            _t.messageShowLoad = false
            if(res.data.code==2002)
            {
              _t.NodeIdList()
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
            }
            else if(res.data.code==2001)
            {
              _t.$message.error(_t.$t("dataModel.IEC104Model.IsErrorExist"));
            }
            else
            {
              _t.$message.error(_t.$t("dataModel.saveFailed"));
            }
          }).catch(function (){
            _t.messageShowLoad = false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    deleteRecord(uuid,muid) {
      let _t = this
      let params={
        uuid:uuid,
        muid:muid
      }
      ModbusTcpDataPushDel(params).then(function (res) {
        if(res.data.code==200)
        {
          _t.$message.success(_t.$t("dataModel.deleteSuccess"));
          _t.NodeIdList()
        }
        else {
          _t.$message.error(_t.$t("dataModel.deleteFailed"));
        }
      })
    },
    onBlackCLK(){
      this.$router.push('/DataPush/ModbusDataTemplete')
    },
  }
}
</script>


<style lang="less" >
.search{
  margin-bottom: 54px;
}
.fold{
  width: calc(100% - 216px);
  display: inline-block
}
.operator{
  margin-bottom: 18px;
}
@media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}
.editable-row-operations a {
  margin-right: 8px;
}
.ant-table-tbody > tr > td {
  padding: 1px 1px;
  overflow-wrap: break-word;
}
.ant-form-item {
  margin-bottom: 5px;

}

.ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
.ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

.ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}

.DataType::-webkit-scrollbar {/*滚动条整体样式*/
  width:4px;/*高宽分别对应横竖滚动条的尺寸*/
  height:4px;
}

.DataType::-webkit-scrollbar-thumb {/*滚动条里面小方块*/
  /*滚动条里面小方块*/
  border-radius   : 10px;
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

.DataType::-webkit-scrollbar-track {/*滚动条里面轨道*/
  /*滚动条里面轨道*/
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}

</style>
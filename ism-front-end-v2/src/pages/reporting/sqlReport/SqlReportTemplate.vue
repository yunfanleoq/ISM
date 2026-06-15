<template>
  <a-card>
    <a-button @click="addVisible=true;textAreValue='';SelectDevice=[];SelectAlarmData=7;SelectPeriod=3" type="primary">  {{$t('configComponent.video.AddVideo')}} </a-button>
    <a-divider type="vertical" />
    <a-button @click="GetSQLReportTemplates()" :disabled="messageShowLoad" type="default" >  {{$t('configComponent.video.VideoRefresh')}} </a-button>
    <div style="margin-top: 5px">
      <a-table rowKey="Name" :pagination="pagination" :columns="columns" :data-source="dataSource" >
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

          <a-button type="link" :disabled="IsExport"  @click="goToEdit(record)" style="margin-right:2px;padding:0;cursor: pointer;color: #13C2C2"><a-icon type="edit" /><span style="margin-left: 2px;">{{$t('monitor.nodeEdit')}}</span></a-button>
          <a-divider type="vertical" />
          <a-upload
              name="file"
              :multiple="false"
              :action='upFileUrl+"/"+record.Uuid'
              :showUploadList="false"
              @change="localUpgradeCharge"
          >
            <a-button type="link"><icon-font type="icon-SQLjiaoben"  />
              {{$t('diyReportTemplete.SQLImportTemplate')}}
            </a-button>
          </a-upload>
          <a-divider type="vertical" />
          <a-button type="link"   @click="goToEditSql(record)" style="margin-right:2px;padding:0;cursor: pointer;color: #13C2C2"><a-icon type="export" /><span style="margin-left: 2px;">{{$t('diyReportTemplete.SQLQuery')}}</span></a-button>
          <a-divider type="vertical" />
          <a-button type="link"   @click="ExportReport(record)" style="margin-right:2px;padding:0;cursor: pointer;color: #13C2C2"><a-icon type="export" /><span style="margin-left: 2px;">{{$t('diyReportTemplete.ExportReport')}}</span></a-button>
          <a-divider type="vertical" />
          <a-popconfirm
              :disabled="IsExport"
              :title="$t('configComponent.video.DelVideoConfirm')"
              :ok-text="$t('component.systemImageModel.delImageYes')"
              :cancel-text="$t('component.systemImageModel.delImageNo')"
              @confirm="DelSQLTemplateAction(record.Uuid)"
          >
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('configComponent.video.DelVideo')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </div>
    <a-modal v-model="addVisible" :title="$t('configComponent.video.AddVideo')" @ok="SQLAddTemplateAction">
      <a-form :form="AddForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 19 }">
        <a-form-item
            :label="$t('diyReportTemplete.Name')"
        >
          <a-input
              v-decorator="[
                'Name',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('diyReportTemplete.Describe')"
        >
          <Mtextarea   v-model="textAreValue"
                       rows="4"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
                       v-decorator="['Describe', { rules: [{ required: true, message: $t('diyReportTemplete.Describe') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model="editVisible"  :title="$t('configComponent.video.AddVideo')" @ok="EditSQLTemplateAction">
      <a-form :form="EditForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 19 }">
        <a-form-item
            :label="$t('diyReportTemplete.Name')"
        >
          <a-input
              v-decorator="[
                'Name',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('diyReportTemplete.Describe')"
        >
          <Mtextarea   v-model="textAreValue"
                       rows="4"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
                       v-decorator="['Describe', { rules: [{ required: true, message: $t('diyReportTemplete.Describe') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal :visible="sqlVisible"
             :title="$t('displayConfig.Properties.CodeEditDig')"
             :width="600"
             :height="800"
             @cancel="sqlVisible=false"
             @ok="EditSQLTemplateQuery"
             :destroyOnClose="true"
             :maskClosable="false"
             :maskStyle="{}"
             :mask="false">
      <div >
        <code-editor
            v-if="sqlVisible"
            :value="codePopUpValue"
            language="javascript"
            @input="changeCodeTextarea($event)"
        >
        </code-editor>
      </div>
    </a-modal>
  </a-card>


</template>

<script>
import {formatDate} from '@/utils/common';
import Mtextarea from "@/components/textarea";
import {
  ExportReportTemplete,
  GetSQLReportTempletes,
  AddSQLReportTemplete,
  DelSQLReportTemplete,
  EditSQLReportTemplete,
  ExportReportTemplate
} from "@/services/SqlReportTemplete";
import {UPDATESQLREPORTTEMPLETE} from "@/services/api";
import codeEditor from '@/components/CodeEditor/index'
export default {
  name: 'SQLDataHistoryTemplete',
  i18n: require('../../../i18n/language'),
  components:{
    Mtextarea,
    codeEditor
  },
  data () {
    return {
      messageShowLoad:false,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      isFullscreen:false,
      codePopUpValue:"",
      sqlVisible:false,
      textAreValue:"",
      editVisible:false,
      addVisible:false,
      SelectAlarmData:7,
      SelectPeriod:3,
      IsExport:false,
      upFileUrl:UPDATESQLREPORTTEMPLETE,
      AddForm:this.$form.createForm(this),
      EditForm:this.$form.createForm(this),
      columns: [
        {
          slotName: 'diyReportTemplete.SQLReportName',
          width: '30%',
          scopedSlots: { customRender: 'Name', title: 'diyReportTemplete.SQLReportName' },
          dataIndex: 'Name',
        },
        {
          width: '30%',
          slotName: 'diyReportTemplete.SQLReportDes',
          scopedSlots: { customRender: 'Describe', title: 'diyReportTemplete.SQLReportDes' },
          dataIndex: 'Describe',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '40%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      dataSource: [],
      SelectDevice:[],
      deviceTreeData:[],
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    this.GetSQLReportTemplates()
  },
  activated(){

  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  created(){

  },
  watch: {

  },
  methods: {
    changeCodeTextarea(val) {
      this.codePopUpValue = val
    },
    localUpgradeCharge(info){
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        if(result.Code==0) {
          this.$message.success(`${info.file.name} `+this.$t("dataModel.importSuccess"));
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
    ExportReport(rc){
      let _t = this
      let params = {
        uuid:rc.Uuid
      };
      _t.$message.loading({ content: _t.$t('diyReportTemplete.ExportLoading'),duration: 0 });
      _t.IsExport=true
      ExportReportTemplate(params).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.destroy()
          const elink = document.createElement('a')
          elink.href = res.data.path
          elink.setAttribute('download', rc.Name+".zip")
          elink.style.display = 'none'
          document.body.appendChild(elink)
          setTimeout(() => {
            elink.click()
            document.body.removeChild(elink)
          }, 66)
        }else{
          _t.IsExport=false
          _t.$message.destroy()
          _t.$message.error(_t.$t('diyReportTemplete.HandExportFailed'))
        }
      }).finally(function (error) {
        _t.IsExport = false
      })
    },
    SQLAddTemplateAction(){
      let _t = this
      this.AddForm.validateFields((err) => {
        if (!err) {
          let params = {
            Name:this.AddForm.getFieldValue('Name'),
            Describe:this.textAreValue
          };
          AddSQLReportTemplete(params).then(function (res){
            if(res.data.code==0)
            {
              _t.GetSQLReportTemplates()
              _t.addVisible = false;
            }else{
              _t.$message.success(_t.$t('alarm.trigger.AddFailed'))
            }
          })
        }
      })
    },
    EditSQLTemplateAction(){
      let _t = this
      this.EditForm.validateFields((err) => {
        if (!err) {
          let params = {
            uuid:_t.editUuid,
            data: {
              Name:this.EditForm.getFieldValue('Name'),
              Describe:this.textAreValue
            }
          };
          EditSQLReportTemplete(params).then(function (res){
            if(res.data.code==0)
            {
              _t.GetSQLReportTemplates()
              _t.$message.success(_t.$t('alarm.trigger.EditSuccess'))
              _t.editVisible = false;
            }else{
              _t.$message.success(_t.$t('VideoManager.EditFailed'))
            }
          })
        }
      })
    },
    goToEdit(item){
      this.editUuid = item.Uuid
      this.textAreValue = item.Describe
      let _t = this
      this.editVisible = true
      this.$message.loading(this.$t("monitor.loading"), 0.5)
      setTimeout(function (){
        _t.EditForm.setFieldsValue(
            {
              Name:item.Name,
              Describe:item.Describe,
            })
      },200)
    },
    goToEditSql(item){
      this.editUuid = item.Uuid
      this.codePopUpValue = item.SqlScript
      let _t = this
      this.sqlVisible = true
    },
    DelSQLTemplateAction(uuid){
      let _t = this

      let params = {
        Uuid:uuid
      };
      DelSQLReportTemplete(params).then(function (res){
        if(res.data.code==0)
        {
          _t.GetSQLReportTemplates()
          _t.$message.success(_t.$t('component.systemImageModel.delImageSuccess'))
          _t.addVisible = false;
        }
        else
        {
          _t.$message.error(_t.$t('component.systemImageModel.delImageFailed'))
        }
      })
    },
    GetSQLReportTemplates(){
      let _t = this

      _t.loading = true
      this.dataSource=[]
      GetSQLReportTempletes().then(function (res){
        _t.loading = false
        if (res.data.code == 0) {
          if(res.data.list==null)
          {
            _t.dataSource=[]
          }
          else
          {
            _t.dataSource = res.data.list
          }
        }
      }).finally(function (error) {
        _t.loading = false
      })
    },
    EditSQLTemplateQuery(){
      let _t = this

      let params = {
        uuid:_t.editUuid,
        data: {
          SqlScript:this.codePopUpValue
        }
      };
      EditSQLReportTemplete(params).then(function (res){
        if(res.data.code==0)
        {
          _t.GetSQLReportTemplates()
          _t.$message.success(_t.$t('alarm.trigger.EditSuccess'))
          _t.sqlVisible = false;
        }else{
          _t.$message.success(_t.$t('VideoManager.EditFailed'))
        }
      })
    },
  }
}
</script>

<style lang="less" scoped>

::v-deep .ant-modal-body {
 padding: 0px;
}
::v-deep .plus-icon-enter-active{
  transition: opacity .5s;
}
::v-deep .plus-icon-enter{
  opacity: 0;
}
.plus-icon-leave-active{
  transition: opacity .5s;
}
.plus-icon-leave-to{
  opacity: 0;
}
.plus-icon-enter-to{
  opacity: 1;
}

::v-deep .code-box-actions {
  padding-top: 12px;
  text-align: center;
  opacity: .7;
  transition: opacity .3s;
}
::v-deep .code-box-meta .demo-description>h4, .code-box-meta>h4 {
  position: absolute;
  top: -14px;
  padding: 1px 8px;
  margin-left: 16px;
  color: #777;
  border-radius: 2px 2px 0 0;
  background: #fff;
  font-size: 14px;
  width: auto;
}
::v-deep .code-box {
  border: 1px solid #ebedf0;
  border-radius: 2px;
  display: inline-block;
  width: 100%;
  position: relative;
  margin: 0 0 16px;
  transition: all .2s;
}

::v-deep .search{
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
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
</style>

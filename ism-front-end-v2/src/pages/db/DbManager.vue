<template>
  <div class="page-header-index-wide">

      <a-card :bordered="false" :bodyStyle="{ padding: '16px', height: '100%' }" :style="{ height: '100%' }">
      <a-tabs default-active-key="3" @change="callback">
        <a-tab-pane key="3" :tab="$t('DbBack.DbConfig')">
          <a-spin :tip="$t('DbBack.Loading')" :spinning="messageShowLoad">
            <a-form :form="form" :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }">
              <a-form-item :label="$t('DbBack.DbType')">
                <a-select
                    @select="chargeDb"
                    v-decorator="['DbType', { rules: [{ required: true, message: $t('DbBack.DbType') }] }]"
                >
                  <a-select-option value="1">
                    Sqlite3
                  </a-select-option>
                  <a-select-option value="0">
                    Mysql
                  </a-select-option>
                  <a-select-option value="4">
                    OceanBase
                  </a-select-option>
                  <a-select-option value="3">
                    达梦数据库
                  </a-select-option>
                </a-select>
              </a-form-item>

              <div v-if="DbType==0||DbType==3||DbType==4">
                <a-form-item :label="$t('DbBack.DbServer')">
                  <a-input
                      v-decorator="['DbServer', { rules: [{ required: true, message: $t('DbBack.DbServer') }] }]"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbServerPort')">
                  <a-input
                      v-decorator="['DbServerPort', { rules: [{ required: true, message: $t('DbBack.DbServerPort') }] }]"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbName')" v-if="DbType==0||DbType==4">
                  <a-input
                      v-decorator="['DbName', { rules: [{ required: true, message: $t('DbBack.DbName') }] }]"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbUserName')">
                  <a-input
                      v-decorator="['DbUserName', { rules: [{ required: true, message: $t('DbBack.DbUserName') }] }]"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbPassword')">
                  <a-input type="password"
                      v-decorator="['DbPassword', { rules: [{ required: true, message: $t('DbBack.DbPassword') }] }]"
                  >
                  </a-input>
                </a-form-item>
              </div>
              <a-form-item :wrapper-col="{ span: 12, offset: 5 }">
                <a-button type="primary" html-type="submit" @click="SetDbConfig" :loading="Configuring">
                  {{$t('DbBack.Sure')}}
                </a-button>
              </a-form-item>
            </a-form>
          </a-spin>
        </a-tab-pane>
        <a-tab-pane key="1" v-if="DbType!=3" :tab="$t('DbBack.Backups')">
          <a-alert
            type="info"
            show-icon
            style="margin-bottom: 12px"
            message="海量历史表（device_record_*）建议排除或使用 scripts/ism_project_backup_core.py 分批备份；全表 SQL 导出可能超时或内存不足。"
          />
          <a-spin :tip="$t('DbBack.BackingUp')" :spinning="messageShowLoad">
          <a-checkbox-group
              style="margin-bottom: 10px"
              v-model="GetTables"
              @change="onChange"
          >

            <a-row>
              <a-col :span="8" v-for="(tableName,index) in SystemTables" :key="index">
                <a-checkbox :value="tableName" >
                  <span  style="font-size: 20px;font-weight: 400">{{ tableName }}</span>
                </a-checkbox>
              </a-col>
            </a-row>
          </a-checkbox-group>
          <a-button type="default" @click="BackUp" :loading="messageShowLoad">
            {{$t('DbBack.Backups')}}
          </a-button>
          <a-checkbox :indeterminate="indeterminate" :checked="checkAll" @change="onCheckAllChange" style="margin-left: 20px">
            {{$t('DbBack.CheckAll')}}
          </a-checkbox>
          </a-spin>
        </a-tab-pane>
        <a-tab-pane key="2" v-if="DbType!=3" :tab="$t('DbBack.Restore')" >
          <a-upload
              name="file"
              :multiple="false"
              :action=localUpgradeUrl
              :showUploadList="false"
              :beforeUpload="beforeUpload"
              @change="localUpgradeCharge"
          >
            <a-button type="default"> <a-icon type="upload" />
              {{$t('DbBack.BackupUpload')}}
            </a-button>
          </a-upload>
        <a-spin :tip="restoreTabSpinTip" :spinning="messageShowLoad">
        <a-table rowKey='FileName' :pagination="pagination" :columns="columns" :data-source="BackupList" >
          <template v-for="(item, index) in columns" :slot="item.slotName">
            <span :key="index">{{ $t(item.slotName) }}</span>
          </template>
          <div slot="Opt" slot-scope="text, record">
            <a type="link"   @click="DbRestore(record.FilePath)" style="cursor: pointer;color: #13C2C2"><a-icon type="reload" /><span style="margin-left: 2px;">{{$t('DbBack.Restore')}}</span></a>
            <a-divider type="vertical" />
            <a type="link"   @click="DbDown(record.FilePath)" style="cursor: pointer;color: #13C2C2"><a-icon type="download" /><span style="margin-left: 2px;">{{$t('DbBack.download')}}</span></a>
          </div>
        </a-table>
        </a-spin>
      </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>
<script>
import {DbBackup,DbDown, GetDbConfig,SetDbConfig,DbRestore, GetBackUpList, GetTablesList} from "@/services/dbbackup";
import {SQLUPLOAD} from "@/services/api";
export default {
  i18n: require('../../i18n/language'),
  data () {
    return {
      checkAll: true,
      messageShowLoad:false,
      // restore | download | loading — 还原页 spin tip，避免下载误显示「还原中」
      restoreTabSpinMode: 'loading',
      indeterminate: false,
      SystemTables:[],
      DbType:"1",
      DbConfig:{},
      localUpgradeUrl:SQLUPLOAD,
      GetTables:[],
      Configuring: false,
      form: this.$form.createForm(this),
      BackupList:[],
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      columns: [
        {
          width: '20%',
          slotName: 'DbBack.Name',
          scopedSlots: { customRender: 'FileName', title: 'DbBack.Name' },
          dataIndex: 'FileName',
        },
        {
          slotName: 'DbBack.BackUpTime',
          width: '20%',
          scopedSlots: { customRender: 'CreateTime', title: 'DbBack.BackUpTime' },
          dataIndex: 'CreateTime',
        },
        {
          slotName: 'DbBack.FileSize',
          width: '20%',
          scopedSlots: { customRender: 'FileSize', title: 'DbBack.FileSize' },
          dataIndex: 'FileSize',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '10%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
    }
  },
  components: {},
  computed: {
    restoreTabSpinTip () {
      if (this.restoreTabSpinMode === 'download') {
        return this.$t('DbBack.Downloading')
      }
      if (this.restoreTabSpinMode === 'restore') {
        return this.$t('DbBack.Restoring')
      }
      return this.$t('DbBack.Loading')
    }
  },
  mounted () {
    this.GetDbConfig()
  },
  methods: {
    chargeDb(v){
      this.DbType = v
    },
    onCheckAllChange(e) {
      Object.assign(this, {
        GetTables: e.target.checked ? this.SystemTables : [],
        indeterminate: false,
        checkAll: e.target.checked,
      });
    },
    onChange(checkedList) {
      this.indeterminate = !!checkedList.length && checkedList.length < this.SystemTables.length;
      this.checkAll = checkedList.length === this.SystemTables.length;
    },
    BackUp(){
      let _t = this
      if( this.GetTables.length==0)
      {
        _t.$message.error(_t.$t('DbBack.BackupsTips'))
        return
      }
      let params={
        tables:this.GetTables
      }
      this.messageShowLoad = true
      DbBackup(params).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t('DbBack.BackupsSuccess'))
        }
        else
        {
          _t.$message.error(_t.$t('DbBack.BackupsFailed'))
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
      })
    },
    GetDbConfig(){
      let _t = this
      let params={

      }
      this.messageShowLoad = true
      GetDbConfig(params).then(function (res){
        if(res.data.code==0)
        {
          _t.DbConfig = res.data.config
          _t.DbType = _t.DbConfig.DbType
          _t.form.setFieldsValue(
              {
                DbType:_t.DbConfig.DbType.toString(),
              })
          if (_t.DbType == 0||_t.DbType == 3||_t.DbType == 4)
          {
            setTimeout(function (){
              _t.form.setFieldsValue(
                  {
                    DbUserName:_t.DbType == 4 ? _t.DbConfig.Oceanbaseuser.toString() : _t.DbConfig.Mysqluser.toString(),
                    DbPassword:_t.DbType == 4 ? _t.DbConfig.Oceanbasepwd.toString() : _t.DbConfig.Mysqlpwd.toString(),
                    DbServer:_t.DbType == 4 ? _t.DbConfig.Oceanbasehost.toString() : _t.DbConfig.Mysqlhost.toString(),
                    DbServerPort:_t.DbType == 4 ? _t.DbConfig.Oceanbaseport.toString() : _t.DbConfig.Mysqlport.toString(),
                    DbName:_t.DbType == 4 ? _t.DbConfig.Oceanbasedbname.toString() : _t.DbConfig.Mysqldbname.toString(),
                  })
              _t.messageShowLoad = false
            },1000)

          }
          else
          {
            _t.messageShowLoad = false
          }
        }
      }).catch(function (error) {
        _t.messageShowLoad = false
      })
    },
    SetDbConfig(){
      let _t = this
      this.form.validateFields((err) => {
        if (!err){
          let params={
            DbType:parseInt(this.form.getFieldValue('DbType')),
          }
          if(params.DbType==0||params.DbType==3||params.DbType==4)
          {
            if (params.DbType == 4) {
              params.Oceanbaseuser   =   this.form.getFieldValue('DbUserName')
              params.Oceanbasepwd    =   this.form.getFieldValue('DbPassword')
              params.Oceanbasehost   =   this.form.getFieldValue('DbServer')
              params.Oceanbaseport   =   this.form.getFieldValue('DbServerPort')
              params.Oceanbasedbname =   this.form.getFieldValue('DbName')
            } else {
              params.Mysqluser   =   this.form.getFieldValue('DbUserName')
              params.Mysqlpwd    =   this.form.getFieldValue('DbPassword')
              params.Mysqlhost   =   this.form.getFieldValue('DbServer')
              params.Mysqlport   =   this.form.getFieldValue('DbServerPort')
              params.Mysqldbname =   this.form.getFieldValue('DbName')
            }
          }
          this.Configuring = true
          SetDbConfig(params).then(function (res){
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t('DbBack.ConfigureSuccess'))
              _t.$router.push("/login")
            }
            else if(res.data.code==-3)
            {
              _t.$message.error(_t.$t('DbBack.DbConnectError'))
            }
            else
            {
              _t.$message.error(_t.$t('DbBack.ConfigureFailed'))
            }
          }).finally(function (error) {
            _t.Configuring = false
          })
        }
      })
    },
    callback(key){
      if(key=="2")
      {
        this.GetBackUpList()
      }
      else if(key=="1")
      {
        this.GetTablesList()
      }
      else if(key=="3")
      {
        this.GetDbConfig()
      }
    },
    GetTablesList(){
      let _t = this
      this.messageShowLoad = true
      _t.SystemTables = []
      GetTablesList().then(function (res){
        if(res.data.code==null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            if(res.data.list[i]!="")
            {
              _t.SystemTables.push(res.data.list[i])
            }
          }
          _t.GetTables = _t.SystemTables
          _t.checkAll = true
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
      })
    },
    GetBackUpList(){
      let _t = this
      this.restoreTabSpinMode = 'loading'
      this.messageShowLoad = true
      _t.BackupList=[]
      GetBackUpList().then(function (res){
        if(res.data.code==0)
        {
          _t.BackupList = res.data.list
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
      })
    },
    DbRestore(name){
      let _t = this
      this.restoreTabSpinMode = 'restore'
      this.messageShowLoad = true
      let params={
        DbFilePath : name
      }
      DbRestore(params).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t('DbBack.RestoreSuccess'))
          _t.$router.push("/login")
        }
        else
        {
          _t.$message.error(_t.$t('DbBack.RestoreFailed'))
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
      })
    },
    beforeUpload(){
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'),duration: 0 });
    },
    localUpgradeCharge(info){
      this.dataSource=[]
      let _t = this
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy();
        if(result.Code==0) {
          this.$message.success(`${info.file.name} `+this.$t("SystemUpgrade.AuthUploadLoading"));
          _t.$store.state.setting.IsOEM = result.Oem
          // 刷新页面
          window.location.reload();
        }
        else if(result.Code==-2)
        {
          this.$message.error(`${info.file.name} `+this.$t("SystemUpgrade.UpgradeFileError"));
        }
        else
        {
          this.$message.error(`${info.file.name} `+this.$t("SystemUpgrade.UpgradeFileSaveError"));
        }
      }
      else if (info.file.status === 'uploading') {
        //this.$message.success(`${info.file.name} `+this.$t("SystemUpgrade.BeginUpgradeUploading"));
      }
      else if (info.file.status === 'error') {
        let _t = this
        _t.$message.loading(_t.$t("SystemUpgrade.BeginUpgradeLoading"),0)
        setTimeout(function (){
          location.reload()
        },10000)
      }

      this.messageShowLoad = false
    },
    DbDown(name){
      let _t = this
      let params={
        DbFilePath : name
      }
      _t.restoreTabSpinMode = 'download'
      _t.messageShowLoad = true
      DbDown(params).then(function (res){
        const data = res.data
        // 失败时后端仍可能返回 JSON blob
        if (data instanceof Blob && data.type && data.type.indexOf('application/json') !== -1) {
          const reader = new FileReader()
          reader.onload = function () {
            try {
              const json = JSON.parse(reader.result)
              if (json.code !== 0) {
                _t.$message.error(_t.$t('DbBack.DownloadFailed'))
              }
            } catch (e) {
              _t.$message.error(_t.$t('DbBack.DownloadFailed'))
            }
          }
          reader.readAsText(data)
          return
        }
        if (!(data instanceof Blob) || data.size === 0) {
          _t.$message.error(_t.$t('DbBack.DownloadFailed'))
          return
        }
        const fileName = (name ? String(name).replace(/\.[^.]+$/, '') : 'backup') + '.zip'
        const url = window.URL.createObjectURL(data)
        const link = document.createElement('a')
        link.href = url
        link.download = fileName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
      }).catch(function () {
        _t.$message.error(_t.$t('DbBack.DownloadFailed'))
      }).finally(function () {
        _t.messageShowLoad = false
      })
    }
  },
  watch: {

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
.page-header-index-wide{
  padding: 20px;
}
.account-settings-info-main {
  width: 100%;
  display: flex;
  height: 100%;
  overflow: auto;

  &.mobile {
    display: block;

    .account-settings-info-left {
      border-right: unset;
      border-bottom: 1px solid #e8e8e8;
      width: 100%;
      height: 50px;
      overflow-x: auto;
      overflow-y: scroll;
    }
    .account-settings-info-right {
      padding: 20px 40px;
    }
  }

  .account-settings-info-left {
    border-right: 1px solid #e8e8e8;
    width: 224px;
  }

  .account-settings-info-right {
    flex: 1 1;
    padding: 8px 40px;

    .account-settings-info-title {
      color: rgba(0,0,0,.85);
      font-size: 20px;
      font-weight: 500;
      line-height: 28px;
      margin-bottom: 12px;
    }
    .account-settings-info-view {
      padding-top: 12px;
    }
  }
}

</style>

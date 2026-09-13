<template>
  <div>
    <a-alert
      type="warning"
      show-icon
      style="margin-bottom: 12px"
      message="下面列出的是备份目录 ism_server_user/data/hisdbbackup/。还原会覆盖当前 TDengine 历史库 ISMHistoryDb。上传 zip ≠ 还原：上传只把备份解压到该目录，必须在下方点「还原」。请打包同时含 dbs.sql 和 taosdump- 文件夹的那一层。当前历史库若不是 TDengine，请先在「数据库配置」改成 TDengine 并确定，再点还原。"
    />
    <a-upload
        name="file"
        :multiple="false"
        accept=".zip"
        :showUploadList="false"
        :beforeUpload="beforeHisUpload"
    >
      <a-button type="default"> <a-icon type="upload" />
        上传历史库备份 zip
      </a-button>
    </a-upload>
    <a-button type="default" style="margin-left: 8px" @click="loadList">
      刷新列表
    </a-button>
    <a-spin :tip="hisRestoreSpinTip" :spinning="hisRestoreLoading">
      <a-table
        style="margin-top: 12px"
        rowKey="FileName"
        :pagination="hisPagination"
        :columns="hisColumns"
        :data-source="HisBackupList"
      >
        <template v-for="(item, index) in hisColumns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <span slot="ValidText" slot-scope="text">{{ text ? '可用' : '不完整' }}</span>
        <div slot="HisOpt" slot-scope="text, record">
          <a-popconfirm
            title="还原将覆盖当前 TDengine 历史库 ISMHistoryDb，原有历史数据会删除。确定继续？"
            :disabled="!record.Valid"
            @confirm="doHisRestore(record.FilePath)"
          >
            <a
              type="link"
              :style="{ cursor: record.Valid ? 'pointer' : 'not-allowed', color: record.Valid ? '#13C2C2' : '#bbb' }"
            ><a-icon type="reload" /><span style="margin-left: 2px;">{{$t('DbBack.Restore')}}</span></a>
          </a-popconfirm>
          <a-divider type="vertical" />
          <a-popconfirm :title="$t('DbBack.DeleteConfirm')" @confirm="doHisDelete(record.FilePath)">
            <a type="link" style="cursor: pointer;color: #eb2f96"><a-icon type="delete" /><span style="margin-left: 2px;">{{$t('DbBack.Delete')}}</span></a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
  </div>
</template>
<script>
import {GetHisBackUpList as requestGetHisBackUpList, HisDbRestore, HisDbDeleteBackup, HisDbBackupUpload} from "@/services/dbbackup";
export default {
  i18n: require('../../i18n/language'),
  data () {
    return {
      hisRestoreLoading:false,
      hisRestoreSpinMode:'loading',
      HisBackupList:[],
      hisPagination:{
        pageSize:15,
        showSizeChanger:true
      },
      hisColumns: [
        {
          width: '24%',
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
          width: '16%',
          scopedSlots: { customRender: 'FileSize', title: 'DbBack.FileSize' },
          dataIndex: 'FileSize',
        },
        {
          title: '完整性',
          width: '12%',
          dataIndex: 'Valid',
          scopedSlots: { customRender: 'ValidText' },
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '18%',
          scopedSlots: { customRender: 'HisOpt', title: 'configComponent.video.TableVideoOpt'}
        }
      ],
    }
  },
  mounted () {
    this.loadList()
  },
  computed: {
    hisRestoreSpinTip () {
      if (this.hisRestoreSpinMode === 'restore') {
        return this.$t('DbBack.Restoring')
      }
      return this.$t('DbBack.Loading')
    }
  },
  methods: {
    loadList(){
      let _t = this
      this.hisRestoreSpinMode = 'loading'
      this.hisRestoreLoading = true
      requestGetHisBackUpList({}).then(function (res){
        if(res.data && res.data.code==0) {
          _t.HisBackupList = res.data.list || []
        } else {
          _t.HisBackupList = []
          if(res.data && res.data.msg) {
            _t.$message.error(res.data.msg)
          }
        }
      }).catch(function (){
        _t.$message.error('获取历史库备份列表失败')
      }).finally(function (){
        _t.hisRestoreLoading = false
      })
    },
    doHisRestore(path){
      let _t = this
      this.hisRestoreSpinMode = 'restore'
      this.hisRestoreLoading = true
      HisDbRestore({DbFilePath: path}).then(function (res){
        if(res.data && res.data.code==0) {
          _t.$message.success(res.data.msg || _t.$t('DbBack.RestoreSuccess'))
        } else {
          _t.$message.error((res.data && res.data.msg) || _t.$t('DbBack.RestoreFailed'))
        }
      }).catch(function (){
        _t.$message.error(_t.$t('DbBack.RestoreFailed'))
      }).finally(function (){
        _t.hisRestoreLoading = false
      })
    },
    doHisDelete(path){
      let _t = this
      HisDbDeleteBackup({DbFilePath: path}).then(function (res){
        if(res.data && res.data.code === 0){
          _t.$message.success(_t.$t('DbBack.DeleteSuccess'))
          _t.loadList()
        } else {
          _t.$message.error((res.data && res.data.msg) || _t.$t('DbBack.DeleteFailed'))
        }
      }).catch(function (){
        _t.$message.error(_t.$t('DbBack.DeleteFailed'))
      })
    },
    beforeHisUpload(file){
      let _t = this
      const name = (file && file.name) ? String(file.name).toLowerCase() : ''
      if (!name.endsWith('.zip')) {
        this.$message.error('只接受 zip 文件')
        return false
      }
      const form = new FormData()
      form.append('file', file)
      this.hisRestoreSpinMode = 'loading'
      this.hisRestoreLoading = true
      HisDbBackupUpload(form).then(function (res){
        const ok = res.data && (res.data.code === 0 || res.data.Code === 0)
        if (ok) {
          _t.$message.success((res.data.msg || '已加入还原列表') + '，不会自动覆盖历史库')
          _t.loadList()
        } else {
          _t.$message.error((res.data && res.data.msg) || '上传失败')
        }
      }).catch(function (){
        _t.$message.error('上传失败')
      }).finally(function (){
        _t.hisRestoreLoading = false
      })
      return false
    },
  }
}
</script>

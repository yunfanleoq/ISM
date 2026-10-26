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
        <span slot="BackupStatus" slot-scope="text, record">
          <div>{{ record.StatusText || (record.Valid ? '备份完成' : '不完整') }}</div>
          <a-progress
            v-if="record.Status === 'running'"
            :percent="Number(record.Progress) || 0"
            :strokeWidth="8"
            status="active"
          />
        </span>
        <div slot="HisOpt" slot-scope="text, record">
          <a
            type="link"
            :style="{ cursor: record.Status === 'running' ? 'not-allowed' : 'pointer', color: record.Status === 'running' ? '#bbb' : '#13C2C2' }"
            @click="record.Status === 'running' ? null : doHisDownload(record)"
          >
            <a-icon type="download" /><span style="margin-left: 2px;">{{$t('DbBack.download')}}</span>
          </a>
          <a-divider type="vertical" />
          <a-popconfirm
            title="还原将覆盖当前 TDengine 历史库 ISMHistoryDb，原有历史数据会删除。确定继续？"
            :disabled="!record.Valid || record.Status === 'running'"
            @confirm="doHisRestore(record.FilePath)"
          >
            <a
              type="link"
              :style="{ cursor: record.Valid ? 'pointer' : 'not-allowed', color: record.Valid ? '#13C2C2' : '#bbb' }"
            ><a-icon type="reload" /><span style="margin-left: 2px;">{{$t('DbBack.Restore')}}</span></a>
          </a-popconfirm>
          <a-divider type="vertical" />
          <a-popconfirm :title="$t('DbBack.DeleteConfirm')" :disabled="record.Status === 'running'" @confirm="doHisDelete(record.FilePath)">
            <a
              type="link"
              :style="{ cursor: record.Status === 'running' ? 'not-allowed' : 'pointer', color: record.Status === 'running' ? '#bbb' : '#eb2f96' }"
            ><a-icon type="delete" /><span style="margin-left: 2px;">{{$t('DbBack.Delete')}}</span></a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
  </div>
</template>
<script>
import {GetHisBackUpList as requestGetHisBackUpList, HisDbRestore, HisDbDeleteBackup, HisDbDownload, HisDbBackupUpload} from "@/services/dbbackup";
export default {
  i18n: require('../../i18n/language'),
  data () {
    return {
      hisRestoreLoading:false,
      hisRestoreSpinMode:'loading',
      HisBackupList:[],
      watching:false,
      seenRunning:false,
      pollTimer:null,
      hisPagination:{
        pageSize:15,
        showSizeChanger:true
      },
      hisColumns: [
        {
          width: '22%',
          slotName: 'DbBack.Name',
          scopedSlots: { customRender: 'FileName', title: 'DbBack.Name' },
          dataIndex: 'FileName',
        },
        {
          slotName: 'DbBack.BackUpTime',
          width: '16%',
          scopedSlots: { customRender: 'CreateTime', title: 'DbBack.BackUpTime' },
          dataIndex: 'CreateTime',
        },
        {
          slotName: 'DbBack.FileSize',
          width: '12%',
          scopedSlots: { customRender: 'FileSize', title: 'DbBack.FileSize' },
          dataIndex: 'FileSize',
        },
        {
          title: '状态',
          width: '24%',
          dataIndex: 'StatusText',
          scopedSlots: { customRender: 'BackupStatus' },
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '26%',
          scopedSlots: { customRender: 'HisOpt', title: 'configComponent.video.TableVideoOpt'}
        }
      ],
    }
  },
  mounted () {
    this.loadList()
  },
  beforeDestroy () {
    this.stopWatch()
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
    stopWatch(){
      this.watching = false
      if (this.pollTimer) {
        clearTimeout(this.pollTimer)
        this.pollTimer = null
      }
    },
    hisBackupListRunning(){
      return (this.HisBackupList || []).some(function (item) {
        return item && item.Status === 'running'
      })
    },
    notifyRunning(){
      const running = this.hisBackupListRunning()
      this.$emit('running-change', running)
      if (running) this.startWatch(true)
    },
    startWatch(alreadyRunning){
      if (this.watching) return
      this.watching = true
      this.seenRunning = !!alreadyRunning
      this.pollOnce()
    },
    pollOnce(){
      if (!this.watching) return
      let _t = this
      this.loadList(true).finally(function () {
        if (!_t.watching) return
        const running = _t.hisBackupListRunning()
        if (running) _t.seenRunning = true
        if (running || !_t.seenRunning) {
          _t.pollTimer = setTimeout(function () { _t.pollOnce() }, running ? 3000 : 2000)
          return
        }
        _t.watching = false
      })
    },
    loadList(silent){
      let _t = this
      if (!silent) {
        this.hisRestoreSpinMode = 'loading'
        this.hisRestoreLoading = true
      }
      return requestGetHisBackUpList({}).then(function (res){
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
        if (!silent) _t.hisRestoreLoading = false
        _t.notifyRunning()
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
    doHisDownload(record){
      const path = record && record.FilePath
      const name = record && record.FileName
      let _t = this
      this.hisRestoreSpinMode = 'loading'
      this.hisRestoreLoading = true
      HisDbDownload({DbFilePath: path}).then(function (res){
        const data = res.data
        if (data instanceof Blob && data.type && data.type.indexOf('application/json') !== -1) {
          const reader = new FileReader()
          reader.onload = function () {
            let msg = '下载失败'
            try {
              const json = JSON.parse(reader.result)
              if (json && json.msg) msg = json.msg
            } catch (e) {}
            _t.$message.error(msg)
          }
          reader.readAsText(data)
          return
        }
        if (!(data instanceof Blob) || data.size === 0) {
          _t.$message.error('下载失败')
          return
        }
        const fileName = (name ? String(name) : 'history-backup') + '.zip'
        const url = window.URL.createObjectURL(data)
        const link = document.createElement('a')
        link.href = url
        link.download = fileName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
      }).catch(function (){
        _t.$message.error('下载失败')
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

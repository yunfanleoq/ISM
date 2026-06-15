<template>
  <a-card :title="$t('SystemUpgrade.about')" :headStyle="{'font-size': '20px','font-weight': 600}">
      <a-list item-layout="horizontal" >
        <a-list-item >
          <a-list-item-meta
              :description="currentVersionDes"
          >
            <a slot="title" >{{systemName}}</a>
            <a-avatar
                slot="avatar"
                :src="systemLogo"
            />
          </a-list-item-meta>
        </a-list-item>
      </a-list>
      <a-spin :spinning="messageShowLoad" >
        <a-card>

          <div style="white-space: pre-wrap;" v-if="remoteVersion!=systemVersion" >
            {{remoteVersionLog}}
          </div>
          <p style="margin-top: 20px">

            <a-upload
                name="file"
                :multiple="false"
                :action=localUpgradeUrl
                :showUploadList="false"
                :beforeUpload="beforeUpload"
                @change="localUpgradeCharge"
            >
              <a-button type="default"> <a-icon type="upload" />
                {{$t('SystemUpgrade.AuthUpload')}}
              </a-button>
            </a-upload>
          </p>
        </a-card>
      </a-spin>
  </a-card>
</template>

<script>

import {LocalUpgrade, OnlineCheckUpgrade, BeginUpgrade,} from "@/services/system";
import {mapState, mapMutations, mapGetters} from 'vuex'
import {AUTHUPLOAD} from "@/services/api";
import Swal from 'sweetalert2/dist/sweetalert2.js';
import '@sweetalert2/theme-bulma/bulma.min.css';
export default {
  name: "SystemParams",
  i18n: require('../../i18n/language'),
  data() {
    return {
      currentVersionDes:"",
      RemoteVersionDes:"",
      remoteVersion:"",
      localUpgradeUrl:AUTHUPLOAD,
      remoteDate:"",
      remoteVersionLog:"",
      remoteIsBete:false,
      messageShowLoad:false,
    };
  },
  components: {

  },
  created(){
    this.currentVersionDes = this.$t('SystemUpgrade.Version')+":"+this.systemVersion+"  ("+(this.versionBete?this.$t('SystemUpgrade.BeteVersion'):this.$t('SystemUpgrade.RCVersion'))+")"+
   "    "+   this.$t('SystemUpgrade.VersionDate')+":"+this.versionDate +"    (x64)"

    this.OnlineUpgrade()

  },
  computed: {
    ...mapState('setting', ['isMobile', 'theme', 'layout', 'footerLinks','systemName', 'copyright', 'fixedHeader', 'fixedSideBar',
      'fixedTabs', 'versionDate', 'versionBete', 'systemVersion']),
    systemLogo () {
      return this.$store.state.setting.SystemLogo
    },
    IsLicense () {
      return this.$store.state.setting.IsOEM
    },
  },
  methods: {
    beforeUpload(){
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'),duration: 0 });
    },
    BeginUpgrade(){
      let _t = this
      this.$confirm({
        title: this.$t('SystemUpgrade.UpgradeConfirm'),
        content: this.$t('SystemUpgrade.UpgradeContent'),
        okText: this.$t('SystemUpgrade.ConfirmUpgrade'),
        cancelText: this.$t('SystemUpgrade.ConfirmCancel'),
        onOk() {
          _t.$message.loading(_t.$t("SystemUpgrade.BeginUpgradeLoading"),0)
          BeginUpgrade().then(function (res){
            _t.$message.destroy()
            if(res.data.code==0)
            {
              _t.remoteVersion = res.data.Result.currentVersion
              _t.remoteIsBete = res.data.Result.IsBete
              _t.remoteDate = res.data.Result.versionData
              _t.remoteVersionLog = res.data.Result.currentVersionDescribeCH

              _t.RemoteVersionDes = _t.$t('SystemUpgrade.RemoteVersion')+":"+_t.remoteVersion+"  ("+(_t.remoteIsBete?_t.$t('SystemUpgrade.BeteVersion'):_t.$t('SystemUpgrade.RCVersion'))+")"+
                  "    "+   _t.$t('SystemUpgrade.VersionDate')+":"+_t.remoteDate+"    (x64)"
            }
            else if(res.data.code==4006)
            {
              _t.$message.error(_t.$t('SystemUpgrade.OnlineCheckNoAuth'))
            }
            else
            {
              _t.$message.error(_t.$t('SystemUpgrade.OnlineCheckError'))
            }
            setTimeout(function (){
              _t.$message.destroy()
            },5000)
            _t.messageShowLoad=false
          }).catch(function(e){
            _t.messageShowLoad=false
            _t.$message.success(_t.$t("SystemUpgrade.BeginUpgradeLoading"));
            setTimeout(function (){
              location.reload()
            },10000)
          })
        }
      });
    },
    OnlineUpgrade(){
      let _t = this
      const params = {
        DebugEnable:this.DebugEnable
      }
      this.messageShowLoad=true
      OnlineCheckUpgrade(params).then(function (res){
        if(res.data.code==0)
        {
          _t.remoteVersion = res.data.Result.currentVersion
          _t.remoteIsBete = res.data.Result.IsBete
          _t.remoteDate = res.data.Result.versionData
          _t.remoteVersionLog = res.data.Result.currentVersionDescribeCH
        }
        else
        {
          _t.$message.error(_t.$t('SystemUpgrade.OnlineCheckError'))
        }
        _t.messageShowLoad=false
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
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
    }
  },
}
</script>

<style scoped>
.ant-form-item{
  margin-bottom: 2px;
}
</style>

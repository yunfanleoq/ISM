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
      <div v-if="IsAuthTimeLimit" style="margin: 16px 0; padding: 20px; background: #fff5f5; border-radius: 8px; border: 1px solid #ffccc7;">
      <div style="display: flex; align-items: center;">
        <a-icon
            type="exclamation-circle"
            style="color: #ff4d4f; font-size: 20px; margin-right: 10px;"
        />
        <div>
          <h4 style="margin: 0 0 8px 0; font-size: 15px; font-weight: 600; color: #2d3748;">
            {{ $t('SystemUpgrade.AuthExpiredTitle') }}
          </h4>
          <p style="margin: 0; font-size: 14px; color: #e53e3e;">
            {{ $t('SystemUpgrade.AuthExpiredDesc') }}
          </p>
        </div>
      </div>
    </div>
      <!-- 授权剩余时间（美化版） -->
      <div v-if="IsAuthLimit&&!IsAuthTimeLimit" style="margin: 16px 0; padding: 20px; background: #f0f7ff; border-radius: 8px; border: 1px solid #dee2e6;">
          <div style="display: flex; align-items: center; margin-bottom: 16px;">
            <a-icon type="clock-circle" style="color: #40a9ff; font-size: 20px; margin-right: 10px;" />
            <h4 style="margin: 0; font-size: 15px; font-weight: 600; color: #2d3748;">
              {{ $t('SystemUpgrade.RemainingLicenseTime') }}
            </h4>
          </div>
          <!-- 剩余时间（天数 + 小时） -->
          <div  style="display: flex; gap: 30px; align-items: center;">
            <div style="text-align: center;">
              <p style="margin: 0; font-size: 32px; font-weight: 700; color: #1890ff;">
                {{ licenseRemainDays }}
              </p>
              <p style="margin: 0; font-size: 14px; color: #718096;">
                {{ $t('SystemUpgrade.Days') }}
              </p>
            </div>

            <span style="font-size: 24px; color: #cbd5e0;">:</span>
            <div style="text-align: center;">
              <p style="margin: 0; font-size: 32px; font-weight: 700; color: #1890ff;">
                {{ licenseRemainHours }}
              </p>
              <p style="margin: 0; font-size: 14px; color: #718096;">
                {{ $t('SystemUpgrade.Hours') }}
              </p>
            </div>

            <span style="font-size: 24px; color: #cbd5e0;">:</span>
            <div style="text-align: center;">
              <p style="margin: 0; font-size: 32px; font-weight: 700; color: #1890ff;">
                {{ licenseRemainMinutes }}
              </p>
              <p style="margin: 0; font-size: 14px; color: #718096;">
                {{ $t('SystemUpgrade.Minutes') }}
              </p>
            </div>
          </div>
      </div>
      <a-spin :spinning="messageShowLoad" >
        <a-card>
         <p style="font-size: 14px;
          line-height: 20px;
          font-weight: 600;" v-if="remoteVersion==systemVersion">{{$t('SystemUpgrade.LastVersion')}}</p>
          <p v-else>
           <a-button type="dashed" @click="BeginUpgrade"> {{RemoteVersionDes}}</a-button>
          </p>
          <div style="white-space: pre-wrap;" v-if="remoteVersion!=systemVersion" >
            {{remoteVersionLog}}
          </div>
          <p style="margin-top: 20px">
          <a-button type="primary" style="margin-right: 10px" icon="check" @click="OnlineUpgrade">
            {{$t('SystemUpgrade.OnlineCheck')}}
          </a-button>
            <a-upload
                name="file"
                :multiple="false"
                :action=localUpgradeUrl
                :showUploadList="false"
                :beforeUpload="beforeUpload"
                @change="localUpgradeCharge"
            >
              <a-button type="default"> <a-icon type="upload" />
                {{$t('SystemUpgrade.localUpgrade')}}
              </a-button>
            </a-upload>
          </p>
        </a-card>
      </a-spin>
  </a-card>
</template>

<script>

import {LocalUpgrade, OnlineCheckUpgrade, BeginUpgrade, GetSystemParams,} from "@/services/system";
import {mapState, mapMutations, mapGetters} from 'vuex'
import {LOCALUPGRADE} from "@/services/api";
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
      localUpgradeUrl:LOCALUPGRADE,
      remoteDate:"",
      remoteVersionLog:"",
      remoteIsBete:false,
      messageShowLoad:false,
      licenseRemainDays: 0, // 剩余天数
      licenseRemainHours: 0, // 剩余小时数
      licenseRemainMinutes: 0, // 剩余小时数
      licenseStatus: "", // 授权状态（例如：正常、过期、即将过期）
      IsAuthLimit:false,
      IsAuthTimeLimit:false
    };
  },
  components: {

  },
  created(){
    this.currentVersionDes = this.$t('SystemUpgrade.Version')+":"+this.systemVersion+"  ("+(this.versionBete?this.$t('SystemUpgrade.BeteVersion'):this.$t('SystemUpgrade.RCVersion'))+")"+
   "    "+   this.$t('SystemUpgrade.VersionDate')+":"+this.versionDate +"    (x64)"

    this.OnlineUpgrade()
    this.GetSystemParams()
  },
  computed: {
    ...mapState('setting', ['isMobile', 'theme', 'layout', 'footerLinks','systemName', 'copyright', 'fixedHeader', 'fixedSideBar',
      'fixedTabs', 'versionDate', 'versionBete', 'systemVersion']),
    systemLogo () {
      return this.$store.state.setting.SystemLogo
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

          _t.RemoteVersionDes = _t.$t('SystemUpgrade.RemoteVersion')+":"+_t.remoteVersion+" ("+(_t.remoteIsBete?_t.$t('SystemUpgrade.BeteVersion'):_t.$t('SystemUpgrade.RCVersion'))+")"+
              " "+ _t.$t('SystemUpgrade.VersionDate')+":"+_t.remoteDate+" (x64)"
          if(_t.remoteVersion!=_t.systemVersion) {
            Swal.fire({
              title: _t.$t('SystemUpgrade.HaveNewVersion'),
              text: _t.RemoteVersionDes,
              html:'<div style="white-space: pre-wrap;text-align:left" >'+_t.RemoteVersionDes+'\n\n'+_t.$t('SystemUpgrade.UpgradeLogger')+'\n'+_t.remoteVersionLog+'</div>',
              icon: 'success',
              showCancelButton: true,
              cancelButtonText:_t.$t('SystemUpgrade.ConfirmCancel'),
              confirmButtonText: _t.$t('SystemUpgrade.BeginUpgrade'),
              preConfirm: () => {
               _t.BeginUpgrade()
              }
            })
          }
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
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy();
        if(result.Code==0) {
          this.$message.success(`${info.file.name} `+this.$t("SystemUpgrade.BeginUpgradeLoading"));
          this.OnlineUpgrade()
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
    GetSystemParams(){
      let _t = this
      GetSystemParams().then(function (res) {
        _t.licenseRemainDays = res.data.list.AuthRemainingTimeDays; // 假设返回结果中有剩余天数
        _t.licenseRemainHours = res.data.list.AuthRemainingTimeHours; // 假设返回结果中有剩余小时数
        _t.licenseRemainMinutes = res.data.list.AuthRemainingTimeMinutes; // 假设返回结果中有剩余小时数
        _t.IsAuthLimit = res.data.list.IsAuthLimit; // 假设返回结果中有授权状态
        _t.IsAuthTimeLimit = res.data.list.IsAuthTimeLimit; // 假设返回结果中有授权状态
      }).catch(function(e){

      })
    },
  },
}
</script>

<style scoped>
.ant-form-item{
  margin-bottom: 2px;
}
</style>

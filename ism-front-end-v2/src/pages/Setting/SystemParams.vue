<template>
  <a-card>
    <a-spin :spinning="messageShowLoad" >
    <a-tabs default-active-key="0" @change="callback">
      <a-tab-pane key="0">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-shijian-copy" />
            {{$t('SystemParams.SystemTime.title')}}
          </span>
        <a-form layout="vertical" style="padding: 10px;" >

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.SystemTime.TimeZone')">
            <a-select v-model="time.timeZone">
              <a-select-option value="Etc/UTC">Etc/UTC</a-select-option>
              <a-select-option value="Asia/Shanghai">Asia/Shanghai</a-select-option>
              <a-select-option value="Asia/Tokyo" >Asia/Tokyo</a-select-option>
              <a-select-option value="Asia/Kolkata" >Asia/Kolkata</a-select-option>
              <a-select-option value="Asia/Dubai" >Asia/Dubai</a-select-option>
              <a-select-option value="Asia/Singapore" >Asia/Singapore</a-select-option>
              <a-select-option value="Asia/Bangkok" >Asia/Bangkok</a-select-option>
              <a-select-option value="Europe/London" >Europe/London</a-select-option>
              <a-select-option value="Europe/Paris" >Europe/Paris</a-select-option>
              <a-select-option value="Europe/Berlin" >Europe/Berlin</a-select-option>
              <a-select-option value="Europe/Moscow">Europe/Moscow</a-select-option>
              <a-select-option value="America/New_York">America/New_York</a-select-option>
              <a-select-option value="America/Chicago">America/Chicago</a-select-option>
              <a-select-option value="America/Denver">America/Denver</a-select-option>
              <a-select-option value="America/Los_Angeles">America/Los_Angeles</a-select-option>
              <a-select-option value="America/Caracas">America/Caracas</a-select-option>
              <a-select-option value="America/Argentina/Buenos_Aires">America/Argentina/Buenos_Aires</a-select-option>
              <a-select-option value="Pacific/Auckland" >Pacific/Auckland</a-select-option>
              <a-select-option value="Pacific/Fiji" >Pacific/Fiji</a-select-option>
              <a-select-option value="Australia/Sydney" >Australia/Sydney</a-select-option>
              <a-select-option value="Atlantic/South_Georgia" >Atlantic/South_Georgia</a-select-option>
              <a-select-option value="Atlantic/Azores" >Atlantic/Azores</a-select-option>
              <a-select-option value="Africa/Cairo" >Africa/Cairo</a-select-option>
              <a-select-option value="Africa/Johannesburg" >Africa/Johannesburg</a-select-option>
              <a-select-option value="Asia/Seoul" >Asia/Seoul</a-select-option>
              <a-select-option value="Asia/Hong_Kong" >Asia/Hong_Kong</a-select-option>
              <a-select-option value="Asia/Manila" >Asia/Manila</a-select-option>
              <a-select-option value="Asia/Jakarta" >Asia/Jakarta</a-select-option>
              <a-select-option value="Asia/Tehran" >Asia/Tehran</a-select-option>
              <a-select-option value="Asia/Baghdad" >Asia/Baghdad</a-select-option>
              <a-select-option value="Asia/Riyadh" >Asia/Riyadh</a-select-option>
              <a-select-option value="Europe/Istanbul" >Europe/Istanbul</a-select-option>
              <a-select-option value="Europe/Madrid" >Europe/Madrid</a-select-option>
              <a-select-option value="Europe/Rome" >Europe/Rome</a-select-option>
              <a-select-option value="Europe/Amsterdam" >Europe/Amsterdam</a-select-option>
              <a-select-option value="Europe/Brussels" >Europe/Brussels</a-select-option>
              <a-select-option value="Europe/Zurich" >Europe/Zurich</a-select-option>
              <a-select-option value="America/Toronto" >America/Toronto</a-select-option>
              <a-select-option value="America/Vancouver" >America/Vancouver</a-select-option>
              <a-select-option value="America/Mexico_City" >America/Mexico_City</a-select-option>
              <a-select-option value="America/Sao_Paulo" >America/Sao_Paulo</a-select-option>
              <a-select-option value="America/Bogota" >America/Bogota</a-select-option>
              <a-select-option value="America/Lima" >America/Lima</a-select-option>
              <a-select-option value="America/Santiago" >America/Santiago</a-select-option>
              <a-select-option value="Pacific/Honolulu" >Pacific/Honolulu</a-select-option>
              <a-select-option value="Pacific/Guam" >Pacific/Guam</a-select-option>
              <a-select-option value="Pacific/Port_Moresby" >Pacific/Port_Moresby</a-select-option>
              <a-select-option value="Pacific/Tongatapu" >Pacific/Tongatapu</a-select-option>
              <a-select-option value="Antarctica/South_Pole" >Antarctica/South_Pole</a-select-option>
              <a-select-option value="Antarctica/McMurdo" >Antarctica/McMurdo</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.SystemTime.DeviceTime')">
            <a-input v-model="time.deviceTime" disabled="disabled"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.SystemTime.CheckType')">
            <a-switch v-model="time.timeCheckType" :checked-children="$t('SystemParams.SystemTime.CheckNTP')" :un-checked-children="$t('SystemParams.SystemTime.CheckHandle')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.SystemTime.NtpServer')">
            <a-input v-model="time.timeNTPServer"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.SystemTime.NtpPort')">
            <a-input v-model="time.timeNTPPort"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.SystemTime.SetTime')">
            <a-row :gutter="24">
              <a-col :span="18">
                <a-date-picker style="width: 100%"
                    v-model="time.setDeviceTime"
                    show-time
                    :defaultValue="moment()"
                    valueFormat="YYYY-MM-DD HH:mm:ss"
                    format="YYYY-MM-DD HH:mm:ss"
                />
              </a-col>
              <a-col :span="4">
                <a-button @click="SyncLocal">{{$t('SystemParams.SystemTime.SyncLocal')}}</a-button>
              </a-col>
            </a-row>
          </a-form-item>

          <a-form-item :label-col="{ span: 4}" :wrapper-col="{ span: 7}" label="    " >
            <a-button type="primary" @click="SaveSysTime">
              {{$t('SystemParams.SystemTime.SaveTime')}}
            </a-button>
            <a-button type="default" @click="TestNtpServer" style="margin-left: 8px">
              {{$t('SystemParams.SystemTime.TestNtp')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="1">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-xitongshezhi" />
            {{$t('SystemParams.System.title')}}
          </span>
        <a-form layout="vertical" style="padding: 10px;" >
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.System.WebSocketPort')">
            <a-input v-model="WSPort"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.System.HttpPort')">
            <a-input v-model="HTTPPort"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.System.HttpsEnable')">
            <a-switch v-model="HTTPSEnable" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.System.HttpsPort')">
            <a-input v-model="HTTPS_Port"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.System.HttpsCert')">
            <a-input v-model="HttpsCert"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.System.HttpsKey')">
            <a-input v-model="HttpsKey">

            </a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 4}" :wrapper-col="{ span: 7}" label="    " >
            <a-button type="primary" @click="SaveSystemWeb">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="default" @click="RebootSystemFunc" style="margin-left: 8px">
              {{$t('NetWork.SystemNetwork.Reboot')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="2">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-mqttx" />
            {{$t('SystemParams.Mqtt.title')}}
          </span>
        <a-form  layout="vertical" style="padding: 10px;" >
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.mqttCloudPlat')">
            <a-input  :disabled="true" value="只支持内置Broker和EMQX"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.mqttEnable')">
            <a-switch v-model="mqtt.isEnable" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.BrokerHost')">
            <a-input v-model="mqtt.BrokerHost"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.BrokerPort')">
            <a-input v-model="mqtt.BrokerPort"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.BrokerUser')">
            <a-input v-model="mqtt.UserName"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.BrokerPassword')">
            <a-input v-model="mqtt.PassWord"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.BrokerClientID')">
            <a-input v-model="mqtt.ClientID"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.ttsEnable')">
            <a-switch v-model="mqtt.TLS" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.SubscribeTopic')">
            <a-input v-model="mqtt.SubscribeTopic"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Mqtt.PublishTopic')">
            <a-input v-model="mqtt.PublishTopic"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 4}" :wrapper-col="{ span: 7}" label="    " >
            <a-button type="primary" @click="SaveSystemMqtt">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="default" @click="RebootSystemFunc" style="margin-left: 8px">
              {{$t('NetWork.SystemNetwork.Reboot')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="3">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-DTUkongzhiqi" />
            {{$t('SystemParams.Modbus.title')}}
          </span>
        <a-form  layout="vertical" style="padding: 10px;" >


          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Modbus.DTUPort')">
            <a-input v-model="modbus.port"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 4}" :wrapper-col="{ span: 7}" label="    " >
            <a-button type="primary" @click="SaveSystemModbus">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="default" @click="RebootSystemFunc" style="margin-left: 8px">
              {{$t('NetWork.SystemNetwork.Reboot')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="4">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-zhaohuan" />
            {{$t('SystemParams.IEC104.title')}}
          </span>
        <a-form  layout="vertical" style="padding: 10px;" >


          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.IEC104.CallTime')">
            <a-input v-model="IECCallTime"></a-input>
          </a-form-item>

          <a-form-item :label-col="{ span: 4}" :wrapper-col="{ span: 7}" label="    " >
            <a-button type="primary" @click="SaveSystemIec104">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="default" @click="RebootSystemFunc" style="margin-left: 8px">
              {{$t('NetWork.SystemNetwork.Reboot')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="5">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-tiaozhuanjiedian" />
            {{$t('SystemParams.PageJump.title')}}
          </span>
        <a-form  layout="vertical" style="padding: 10px;" >


          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.PageJump.jump')">
            <a-switch v-model="jump.enable" :checked-children="$t('SystemParams.PageJump.jumpEnable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>

          <a-form-item :label-col="{ span: 4}" :wrapper-col="{ span: 7}" label="    " >
            <a-button type="primary" @click="saveJumpWindows">
              {{$t('AlarmTips.Save')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="6">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-opcUa" />
            {{$t('SystemParams.Opcua.title')}}
          </span>
        <a-form  layout="vertical" style="padding: 10px;" >
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.publish_enable')">
            <a-switch v-model="opcua.publish_enable" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.require_external_cert')">
            <a-switch v-model="opcua.require_external_cert" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.allow_anonymous')">
            <a-switch v-model="opcua.allow_anonymous" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.allow_username')">
            <a-switch v-model="opcua.allow_username" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.security_policy_none')">
            <a-switch v-model="opcua.security_policy_none" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
         <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.security_policies')">
            <a-select v-model="opcua.security_policies" mode="multiple" :placeholder="$t('SystemParams.Opcua.select_policies')">
              <a-select-option value="None">None</a-select-option>
              <a-select-option value="Basic128Rsa15">Basic128Rsa15</a-select-option>
              <a-select-option value="Basic256">Basic256</a-select-option>
              <a-select-option value="Basic256Sha256">Basic256Sha256</a-select-option>
              <a-select-option value="Aes128Sha256RsaOaep">Aes128Sha256RsaOaep</a-select-option>
              <a-select-option value="Aes256Sha256RsaPss">Aes256Sha256RsaPss</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.message_security_modes')">
            <a-select v-model="opcua.message_security_modes" mode="multiple" :placeholder="$t('SystemParams.Opcua.select_modes')">
              <a-select-option value="None">None</a-select-option>
              <a-select-option value="Sign">Sign</a-select-option>
              <a-select-option value="SignAndEncrypt">SignAndEncrypt</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.insecure_skip_verify')">
            <a-switch v-model="opcua.insecure_skip_verify" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.server_diagnostics')">
            <a-switch v-model="opcua.server_diagnostics" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.server_name')">
            <a-input v-model="opcua.server_name"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.product_URI')">
            <a-input v-model="opcua.product_URI"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.product_name')">
            <a-input v-model="opcua.product_name"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.manufacturer_name')">
            <a-input v-model="opcua.manufacturer_name"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.host')">
            <a-input v-model="opcua.host"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.port')">
            <a-input v-model="opcua.port"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.cert')">
            <a-input v-model="opcua.cert"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.key')">
            <a-input v-model="opcua.key"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.username')">
            <a-input v-model="opcua.username"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('SystemParams.Opcua.password')">
            <a-input v-model="opcua.password"></a-input>
          </a-form-item>
          <a-form-item :label-col="{ span: 4}" :wrapper-col="{ span: 7}" label="    " >
            <a-button type="primary" @click="SaveSystemOpcua">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="default" @click="RebootSystemFunc" style="margin-left: 8px">
              {{$t('NetWork.SystemNetwork.Reboot')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
    </a-tabs>
    </a-spin>
  </a-card>
</template>

<script>
import {
  RebootSystem,
  GetSystemWebParams,
  SaveSystemWebParams,
  SaveSystemMqttData,
  TestNTPConfig,
  GetSystemOpcuaData,
  SaveSystemOpcuaData,
} from "@/services/system";
import {SaveSystemTimeConfig,GetSystemModbusData, GetSystemMqttData, RebootISMSystem, SaveSystemModbusData} from "@/services/system";
import moment from 'moment';
import 'moment/locale/zh-cn';
import  'moment/locale/en-ie';
import  'moment/locale/zh-tw';
export default {
  name: "SystemParams",
  i18n: require('../../i18n/language'),
  data() {
    return {
      moment,
      jump:{
        enable:true,
      },
      time:{
        timeZone:"Asia/Shanghai",
        timeCheckType:false,
        timeNTPServer:"time.windows.com",
        timeNTPPort:123,
        timeNTPCheckTime:1440,
        deviceTime:"",
        setDeviceTime:""
      },
      labelCol: { span: 2 },
      tabKey:"1",
      HTTPPort:8081,
      HTTPS_Port:0,
      HTTPSEnable:false,
      DebugEnable:false,
      HttpsKey:"",
      HttpsCert:"",
      WSPort:10215,
      messageShowLoad:false,
      IECCallTime:0,
      mqtt:{
        isEnable:false,
        mqttCloudPlat:1,
        BrokerHost:"",
        BrokerPort:1883,
        UserName:"",
        PassWord:"",
        ClientID:"",
        TLS:false,
        certPath:"",
        SubscribeTopic:"",
        PublishTopic:"",
      },
      modbus:{
        port:3000
      },
      opcua:{
        publish_enable:true,
        require_external_cert:false,
        allow_anonymous:false,
        allow_username:true,
        security_policy_none:true,
        security_policies: ["None", "Basic128Rsa15", "Basic256", "Basic256Sha256", "Aes128Sha256RsaOaep", "Aes256Sha256RsaPss"],
        message_security_modes: ["None", "Sign", "SignAndEncrypt"],
        insecure_skip_verify:true,
        server_diagnostics:true,
        server_name:"ismopcua",
        product_URI:"https://www.ismctl.com",
        product_name:"ISM OPC UA Publisher2034",
        manufacturer_name:"ISM",
        host:"127.0.0.1",
        port:"4841",
        cert:"conf/192.168.199.120.crt",
        key:"conf/192.168.199.120.key",
        username:"ismopcua",
        password:"ism123456"
      },
      wrapperCol: { span: 10 },
    };
  },
  components: {

  },
  created(){
    this.GetSystemWeb()
  },
  methods: {
    SyncLocal(){
      this.time.setDeviceTime = moment().format('YYYY-MM-DD HH:mm:ss')
    },
    callback(key) {
      if(key==1||key==0)
      {
        this.GetSystemWeb()
      }
      else if(key==2)
      {
        this.GetSystemMqtt()
      }
      else if(key==3)
      {
        this.GetSystemModbus()
      }
      else if(key==4)
      {
        this.GetSystemIec104()
      }
      else if(key==6)
      {
        this.GetSystemOpcua()
      }
    },
    GetSystemWeb(){
      let _t = this
      const params = {
        DebugEnable:this.DebugEnable
      }
      this.messageShowLoad=true
      GetSystemWebParams(params).then(function (res){
        _t.messageShowLoad=false
        _t.WSPort = res.data.result.wsport
        _t.HTTPPort = res.data.result.httpport
        _t.HTTPS_Port = res.data.result.httpsport
        _t.HTTPSEnable = res.data.result.enablehttps
        _t.HttpsKey = res.data.result.httpskeyfile
        _t.HttpsCert = res.data.result.httpscertfile
        _t.time.deviceTime = res.data.result.systemtime
        _t.time.setDeviceTime = res.data.result.systemtime
        _t.time.timeCheckType = res.data.result.CheckType==1?true:false
        _t.time.timeNTPServer = res.data.result.NTPServer
        _t.time.timeNTPPort = res.data.result.NTPPort
        _t.time.timeZone = res.data.result.TimeZone
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    GetSystemMqtt(){
      let _t = this
      this.messageShowLoad=true
      GetSystemMqttData().then(function (res){
        _t.messageShowLoad=false
        _t.mqtt.isEnable = res.data.result.isEnable
        _t.mqtt.BrokerHost = res.data.result.BrokerHost
        _t.mqtt.BrokerPort = res.data.result.BrokerPort
        _t.mqtt.UserName = res.data.result.UserName
        _t.mqtt.PassWord = res.data.result.PassWord
        _t.mqtt.ClientID = res.data.result.ClientID
        _t.mqtt.TLS = res.data.result.TLS
        _t.mqtt.certPath = res.data.result.certPath
        _t.mqtt.SubscribeTopic = res.data.result.SubscribeTopic
        _t.mqtt.PublishTopic = res.data.result.PublishTopic
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    SaveSystemWeb(){
      let _t = this
      const params = {
        HTTPPort:parseInt(_t.HTTPPort),
        WSPort:parseInt(_t.WSPort),
        HTTPsPort:parseInt(_t.HTTPS_Port),
        HTTPSEnable:_t.HTTPSEnable,
        HttpsKey:_t.HttpsKey,
        HttpsCert:_t.HttpsCert
      }
      this.messageShowLoad=true
      SaveSystemWebParams(params).then(function (res){
        _t.messageShowLoad=false
        if(res.data.code==0){
          _t.$message.success(_t.$t('NetWork.SystemNetwork.SaveSuccess'), 3)
        }else{
          _t.$message.error(_t.$t('NetWork.SystemNetwork.SaveFailed'), 3)
        }
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    SaveSystemMqtt(){
      let _t = this
      if(_t.mqtt.SubscribeTopic.indexOf("${ClientID}")==-1||_t.mqtt.PublishTopic.indexOf("${ClientID}")==-1)
      {
        _t.$message.error(_t.$t('SystemParams.Mqtt.TopicError'), 3)
        return
      }
      const params = {
        isEnable:_t.mqtt.isEnable,
        BrokerHost:_t.mqtt.BrokerHost,
        BrokerPort:parseInt(_t.mqtt.BrokerPort),
        UserName:_t.mqtt.UserName,
        PassWord:_t.mqtt.PassWord,
        ClientID:_t.mqtt.ClientID,
        TLS:_t.mqtt.TLS,
        certPath:_t.mqtt.certPath,
        SubscribeTopic:_t.mqtt.SubscribeTopic,
        PublishTopic:_t.mqtt.PublishTopic,
      }
      this.messageShowLoad=true
      SaveSystemMqttData(params).then(function (res){
        _t.messageShowLoad=false
        if(res.data.code==0){
          _t.$message.success(_t.$t('NetWork.SystemNetwork.SaveSuccess'), 3)
        }else{
          _t.$message.error(_t.$t('NetWork.SystemNetwork.SaveFailed'), 3)
        }
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    GetSystemModbus(){
      let _t = this
      const params = {
        DebugEnable:this.DebugEnable
      }
      this.messageShowLoad=true
      GetSystemModbusData(params).then(function (res){
        _t.messageShowLoad=false
        _t.modbus.port = res.data.result.modbusServerPort
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    SaveSystemModbus(){
      let _t = this

      const params = {
        ModbusServerPort:parseInt(_t.modbus.port),
      }
      this.messageShowLoad=true
      SaveSystemModbusData(params).then(function (res){
        _t.messageShowLoad=false
        if(res.data.code==0){
          _t.$message.success(_t.$t('NetWork.SystemNetwork.SaveSuccess'), 3)
        }else{
          _t.$message.error(_t.$t('NetWork.SystemNetwork.SaveFailed'), 3)
        }
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    GetSystemIec104(){
      let _t = this
      const params = {
        DebugEnable:this.DebugEnable
      }
      this.messageShowLoad=true
      GetSystemModbusData(params).then(function (res){
        _t.messageShowLoad=false
        _t.IECCallTime = res.data.result.iec104calldelaytime
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    saveJumpWindows(){
      let setSpeech = {
        enable:this.jump.enable,
      }
      localStorage.setItem("JumpWindow",JSON.stringify(setSpeech))
      this.$message.success(this.$t('AlarmTips.SaveSuccess'))
    },
    SaveSystemIec104(){
      let _t = this

      const params = {
        ModbusServerPort:parseInt(_t.modbus.port),
        iec104calldelaytime:parseInt(_t.IECCallTime),
      }
      this.messageShowLoad=true
      SaveSystemModbusData(params).then(function (res){
        _t.messageShowLoad=false
        if(res.data.code==0){
          _t.$message.success(_t.$t('NetWork.SystemNetwork.SaveSuccess'), 3)
        }else{
          _t.$message.error(_t.$t('NetWork.SystemNetwork.SaveFailed'), 3)
        }
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    GetSystemOpcua(){
      let _t = this
      this.messageShowLoad=true
      GetSystemOpcuaData().then(function (res){
        console.log(res)
        _t.messageShowLoad=false
        _t.opcua.publish_enable = res.data.result.opcua_publish_enable
        _t.opcua.require_external_cert = res.data.result.opcua_publish_require_external_cert
        _t.opcua.allow_anonymous = res.data.result.opcua_publish_allow_anonymous
        _t.opcua.allow_username = res.data.result.opcua_publish_allow_username
        _t.opcua.security_policy_none = res.data.result.opcua_publish_security_policy_none
        _t.opcua.security_policies = res.data.result.opcua_publish_security_policies ? res.data.result.opcua_publish_security_policies.split(',') : []
        _t.opcua.message_security_modes = res.data.result.opcua_publish_message_security_modes ? res.data.result.opcua_publish_message_security_modes.split(',') : []
        _t.opcua.insecure_skip_verify = res.data.result.opcua_publish_insecure_skip_verify
        _t.opcua.server_diagnostics = res.data.result.opcua_publish_server_diagnostics
        _t.opcua.server_name = res.data.result.opcua_publish_server_name
        _t.opcua.product_URI = res.data.result.opcua_publish_product_URI
        _t.opcua.product_name = res.data.result.opcua_publish_product_name
        _t.opcua.manufacturer_name = res.data.result.opcua_publish_manufacturer_name
        _t.opcua.host = res.data.result.opcua_publish_host
        _t.opcua.port = res.data.result.opcua_publish_port
        _t.opcua.cert = res.data.result.opcua_publish_cert
        _t.opcua.key = res.data.result.opcua_publish_key
        _t.opcua.username = res.data.result.opcua_publish_username
        _t.opcua.password = res.data.result.opcua_publish_password
      }).catch(function(e){
        console.log(e)
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    SaveSystemOpcua(){
      let _t = this
      const params = {
        opcua_publish_enable:_t.opcua.publish_enable,
        opcua_publish_require_external_cert:_t.opcua.require_external_cert,
        opcua_publish_allow_anonymous:_t.opcua.allow_anonymous,
        opcua_publish_allow_username:_t.opcua.allow_username,
        opcua_publish_security_policy_none:_t.opcua.security_policy_none,
        opcua_publish_security_policies: Array.isArray(_t.opcua.security_policies) ? _t.opcua.security_policies.join(',') : _t.opcua.security_policies,
        opcua_publish_message_security_modes: Array.isArray(_t.opcua.message_security_modes) ? _t.opcua.message_security_modes.join(',') : _t.opcua.message_security_modes,
        opcua_publish_insecure_skip_verify:_t.opcua.insecure_skip_verify,
        opcua_publish_server_diagnostics:_t.opcua.server_diagnostics,
        opcua_publish_server_name:_t.opcua.server_name,
        opcua_publish_product_URI:_t.opcua.product_URI,
        opcua_publish_product_name:_t.opcua.product_name,
        opcua_publish_manufacturer_name:_t.opcua.manufacturer_name,
        opcua_publish_host:_t.opcua.host,
        opcua_publish_port:_t.opcua.port,
        opcua_publish_cert:_t.opcua.cert,
        opcua_publish_key:_t.opcua.key,
        opcua_publish_username:_t.opcua.username,
        opcua_publish_password:_t.opcua.password
      }
      this.messageShowLoad=true
      SaveSystemOpcuaData(params).then(function (res){
        _t.messageShowLoad=false
        if(res.data.code==0){
          _t.$message.success(_t.$t('NetWork.SystemNetwork.SaveSuccess'), 3)
        }else{
          _t.$message.error(_t.$t('NetWork.SystemNetwork.SaveFailed'), 3)
        }
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    SaveSysTime(){
      let _t = this
      const params = {
        TimeZone:_t.time.timeZone,
        CheckType:_t.time.timeCheckType?1:0,
        NTPServer:_t.time.timeNTPServer,
        NTPPort:parseInt(_t.time.timeNTPPort),
        NTPCheckTime:parseInt(_t.time.timeNTPCheckTime),
        SetTime:_t.time.setDeviceTime
      }
      this.messageShowLoad=true
      SaveSystemTimeConfig(params).then(function (res){
        _t.messageShowLoad=false
        if(res.data.code==0){
          _t.$message.success(_t.$t('NetWork.SystemNetwork.SaveSuccess'), 3)
        }else{
          _t.$message.error(_t.$t('NetWork.SystemNetwork.SaveFailed'), 3)
        }
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    TestNtpServer(){
      let _t = this
      const params = {
        NTPServer:_t.time.timeNTPServer,
        NTPPort:parseInt(_t.time.timeNTPPort),
      }
      this.messageShowLoad=true
      TestNTPConfig(params).then(function (res){
        _t.messageShowLoad=false
        if(res.data.code==0){
          _t.$message.success(_t.$t('SystemParams.SystemTime.TestNtpSuccess'), 3)
        }else{
          _t.$message.error(_t.$t('SystemParams.SystemTime.TestNtpError'), 3)
        }
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('SystemParams.SystemTime.TestNtpError'), 3)
      })
    },
    RebootSystemFunc(){
      let _t = this
      _t.messageShowLoad = true
      RebootISMSystem().then(function (res){
        if(res.data.code==0||res.data.code==-5)
        {
          _t.$message.success( _t.$t('NetWork.SystemNetwork.RebootSuccess'))
        }
        else
        {
          _t.$message.error( _t.$t('NetWork.SystemNetwork.RebootFailed'))
        }
      }).catch(function (error) {
        _t.messageShowLoad = false
        _t.$message.success( _t.$t('NetWork.SystemNetwork.RebootSuccess'))
      }).finally(function (error) {
        _t.messageShowLoad = false
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
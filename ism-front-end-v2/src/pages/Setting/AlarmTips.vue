<template>
  <a-card>
    <a-spin :spinning="messageShowLoad" >
    <a-tabs default-active-key="1" @change="callback">
      <a-tab-pane key="1">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-speech" />
            {{$t('AlarmTips.WEBSpeech.title')}}
          </span>
        <a-form v-if="tabKey==1" :form="SpeechFrom" layout="vertical" style="padding: 100px;" >
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('AlarmTips.WEBSpeech.title')">

            <a-switch v-model="WEBSpeechEnable" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.WEBSpeech.Speed')}}&nbsp;
              <a-tooltip :title="$t('AlarmTips.WEBSpeech.SpeedTips')">
                <a-icon type="question-circle-o" />
              </a-tooltip>
            </span>
            <a-slider  :marks="marks" v-model="WEBSpeechSpeed" :tooltip-visible="true"  :min="0.1" :max="2.01" :step="0.1" />
          </a-form-item>
          <a-form-item style="margin-left: 100px" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
            <a-button type="primary" @click="saveSpeech">
              {{$t('AlarmTips.Save')}}
            </a-button>
          </a-form-item>

        </a-form>
      </a-tab-pane>
      <a-tab-pane key="2" force-render>
        <span slot="tab" style="font-size: 20px">
             <a-icon type="mail" />
            {{$t('AlarmTips.Email.title')}}
          </span>
        <div style="padding: 20px" v-if="tabKey==2">
        <a-form  :form="EmailForm" layout="vertical" style="padding: 10px;" >

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" :label="$t('AlarmTips.Email.EmailNotice')">
            <a-switch  v-decorator="[
                  'MailEnable',
                  {
                    initialValue: true,valuePropName: 'checked',rules: [{ required: false, message: $t('AlarmTips.Email.EmailNotice') }],
                  },
                ]"  :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Email.MailServer')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['MailServerIP', {rules: [{ required: true, message: $t('AlarmTips.Email.MailServer'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Email.MailServerPort')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['MailServerPort', {rules: [{ required: true, message: $t('AlarmTips.Email.MailServerPort'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Email.MailUser')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['MailSendUser', {rules: [{ required: true, message: $t('AlarmTips.Email.MailUser'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Email.MailPassword')}}&nbsp;
            </span>
            <a-input-password   autocomplete="autocomplete"
                      v-decorator="['MailSendPassword', {rules: [{ required: true, message: $t('AlarmTips.Email.MailPassword'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item style="display: none" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" :label="$t('AlarmTips.Email.TLS')">
            <a-switch  v-decorator="[
                  'TLS',
                  {
                    initialValue: true,valuePropName: 'checked',rules: [{ required: false, message: $t('AlarmTips.Email.TLS') }],
                  },
                ]"  :checked-children="$t('AlarmTips.Email.TLSEnable')" :un-checked-children="$t('AlarmTips.Email.TLSDisable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Email.MailTo')}}&nbsp;
                <a-tooltip :title="$t('AlarmTips.Email.MailToTips')">
                <a-icon type="question-circle-o" />
              </a-tooltip>
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['MailTo', {rules: [{ required: false, message: $t('AlarmTips.Email.MailTo'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Email.MailSendUserName')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['MailSendUserName', {rules: [{ required: true, message: $t('AlarmTips.Email.MailSendUserName'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Email.MailSendSubject')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['MailSendSubject', {rules: [{ required: true, message: $t('AlarmTips.Email.MailSendSubject'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item style="margin-left: 100px" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
            <a-button type="primary" @click="UpdateAlarmNotice('Mail')">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="primary" style="margin-left: 10px" @click="TestEmail">
              {{$t('AlarmTips.Test')}}
            </a-button>
          </a-form-item>
        </a-form>
        </div>
      </a-tab-pane>
      <a-tab-pane key="3">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-duanxin" />
            {{$t('AlarmTips.Phone.title')}}
          </span>
        <a-form  :form="PhoneForm" layout="vertical" style="padding: 10px;" >

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" :label="$t('AlarmTips.Phone.NoticeEnable')">
            <a-switch  v-decorator="[
                  'PhoneNoticeEnable',
                  {
                    initialValue: true,valuePropName: 'checked',rules: [{ required: false, message: $t('AlarmTips.Phone.EmailNotice') }],
                  },
                ]"  :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.SMSFactory')}}&nbsp;
            </span>
            <a-select
                style="width: 100%"
                v-decorator="[
                  'SmsType',
                  {
                    initialValue: '1',rules: [{ required: false, message: $t('AlarmTips.Phone.SMSFactory') }],
                  },
                ]"
            >
              <a-select-option value="1">
                {{$t('AlarmTips.Phone.AliYunSms')}}
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.AccessKeyId')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['AccessKeyId', {rules: [{ required: true, message: $t('AlarmTips.Phone.AccessKeyId'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.AccessKeySecret')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['AccessKeySecret', {rules: [{ required: true, message: $t('AlarmTips.Phone.AccessKeySecret'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.SignName')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['SignName', {rules: [{ required: true, message: $t('AlarmTips.Phone.SignName'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.TemplateCode')}}&nbsp;
            </span>
            <a-input   autocomplete="autocomplete"
                                v-decorator="['TemplateCode', {rules: [{ required: true, message: $t('AlarmTips.Phone.TemplateCode'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.EveryDayCount')}}&nbsp;
            </span>
            <a-input  style="width: 100%" autocomplete="autocomplete"
                                v-decorator="['EveryDayCount', {rules: [{ required: true, message: $t('AlarmTips.Phone.EveryDayCount'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.SendAlarmLevel')}}&nbsp;
            </span>
            <a-select
                mode="multiple"
                :allowClear="true"
                style="width: 100%"
                v-model="showAlarmItems"
            >
              <a-select-option value='0'>{{$t('dataModel.alarm.Tips')}}</a-select-option>
              <a-select-option value='1'>{{$t('dataModel.alarm.Minor')}}</a-select-option>
              <a-select-option value='2'>{{$t('dataModel.alarm.Importance')}}</a-select-option>
              <a-select-option value='3'>{{$t('dataModel.alarm.Urgency')}}</a-select-option>
              <a-select-option value='4'>{{$t('dataModel.alarm.Deadly')}}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Phone.SendPhoneNumbers')}}&nbsp;
                <a-tooltip :title="$t('AlarmTips.Phone.SendPhoneNumbersTips')">
                <a-icon type="question-circle-o" />
              </a-tooltip>
            </span>
            <a-input  style="width: 100%" autocomplete="autocomplete"
                             v-decorator="['SendPhoneNumbers', {rules: [{ required: true, message: $t('AlarmTips.Phone.SendPhoneNumbers'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item style="margin-left: 100px" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
            <a-button type="primary" @click="UpdateAlarmNotice('Phone')">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="primary" style="margin-left: 10px" @click="TestSms">
              {{$t('AlarmTips.Test')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="4">
        <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-dingtalk" />
            {{$t('AlarmTips.dingTalk.title')}}
          </span>
        <a-button type="primary" icon="save" @click="UpdateAlarmNotice('dingTalk')">
          {{$t('AlarmTips.Save')}}
        </a-button>
        <a-divider type="vertical" />
        <a-button  icon="plus" @click="AddDingTalk">
          {{$t('AlarmTips.dingTalk.New')}}
        </a-button>
        <div style="padding: 20px" v-if="tabKey==4">
          <a-form  layout="inline" v-for="(item, index) in DingTalkList" :key="index">
            <a-form-item  :label="$t('AlarmTips.dingTalk.DingTalkEnable')">
              <a-switch  v-model="item.IsEnable"  :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
              </a-switch>
            </a-form-item>
            <a-form-item  >
             <span slot="label">
              {{$t('AlarmTips.dingTalk.Webhook')}}&nbsp;
            </span>
              <a-input  autocomplete="autocomplete" v-model="item.Webhook"

              />
            </a-form-item>
            <a-form-item  >
             <span slot="label">
              {{$t('AlarmTips.dingTalk.Secret')}}&nbsp;
            </span>
              <a-input  autocomplete="autocomplete" v-model="item.Secret"
              />
            </a-form-item>
            <a-form-item style=""  >
              <a-button type="primary" icon="arrow-right" style="margin-left: 10px" @click="TestDingTalk(index)">
                {{$t('AlarmTips.Test')}}
              </a-button>
              <a-popconfirm
                  :title="$t('AlarmTips.dingTalk.DeleteTips')"
                  @confirm="DelDingTalk(index)"
              >
                <a-button type="danger" icon="close" style="margin-left: 10px" >
                  {{$t('AlarmTips.dingTalk.Delete')}}
                </a-button>
              </a-popconfirm>

            </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>
      <a-tab-pane key="5">
        <span slot="tab" style="font-size: 20px">
             <a-icon type="wechat" />
            {{$t('AlarmTips.Wechat.title')}}
          </span>
        <a-form  :form="weChatForm" layout="vertical" style="padding: 10px;" >

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" :label="$t('AlarmTips.Wechat.Enable')">
            <a-switch  v-decorator="[
                  'WeChatEnable',
                  {
                    initialValue: true,valuePropName: 'checked',rules: [{ required: false, message: $t('AlarmTips.Email.EmailNotice') }],
                  },
                ]"  :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Wechat.EnterpriseID')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['EnterpriseID', {rules: [{ required: true, message: $t('AlarmTips.Wechat.EnterpriseID'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Wechat.AgentId')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['AgentId', {rules: [{ required: true, message: $t('AlarmTips.Wechat.AgentId'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Wechat.Secret')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['Secret', {rules: [{ required: true, message: $t('AlarmTips.Wechat.Secret'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item style="margin-left: 100px" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
            <a-button type="primary" @click="UpdateAlarmNotice('weChat')">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="primary" style="margin-left: 10px" @click="TestWeChat">
              {{$t('AlarmTips.Test')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="6">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-guzhangyuce" />
            {{$t('AlarmTips.AlarmWindows.title')}}
          </span>
        <a-form v-if="tabKey==6" :form="AlarmWindowsFrom" layout="vertical" style="padding: 100px;" >
          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('AlarmTips.AlarmWindows.title')">

            <a-switch v-model="AlarmWindowsEnable" :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.AlarmWindows.OpenLevel')}}&nbsp;
            </span>
            <a-select
                mode="multiple"
                :allowClear="true"
                style="width: 100%"
                v-model="showWindowsAlarmItems"
            >
              <a-select-option value='0'>{{$t('dataModel.alarm.Tips')}}</a-select-option>
              <a-select-option value='1'>{{$t('dataModel.alarm.Minor')}}</a-select-option>
              <a-select-option value='2'>{{$t('dataModel.alarm.Importance')}}</a-select-option>
              <a-select-option value='3'>{{$t('dataModel.alarm.Urgency')}}</a-select-option>
              <a-select-option value='4'>{{$t('dataModel.alarm.Deadly')}}</a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="wrapperCol" :label="$t('AlarmTips.AlarmWindows.AutoClose')">

            <a-switch v-model="AlarmWindowsAutoClose" :checked-children="$t('AlarmTips.AlarmWindows.AutoCloseOpen')" :un-checked-children="$t('AlarmTips.AlarmWindows.AutoCloseTrue')">
            </a-switch>
          </a-form-item>

          <a-form-item style="margin-left: 100px" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
            <a-button type="primary" @click="saveAlarmWindows">
              {{$t('AlarmTips.Save')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="7">
         <span slot="tab" style="font-size: 20px">
             <icon-font type="icon-shouji" />
            {{$t('AlarmTips.Voice.title')}}
          </span>
        <a-form  :form="VoiceForm" layout="vertical" style="padding: 10px;" >

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" :label="$t('AlarmTips.Voice.NoticeEnable')">
            <a-switch  v-decorator="[
                  'PhoneNoticeEnable',
                  {
                    initialValue: true,valuePropName: 'checked',rules: [{ required: false, message: $t('AlarmTips.Voice.EmailNotice') }],
                  },
                ]"  :checked-children="$t('AlarmTips.WEBSpeech.Enable')" :un-checked-children="$t('AlarmTips.WEBSpeech.Disable')">
            </a-switch>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Voice.SMSFactory')}}&nbsp;
            </span>
            <a-select
                style="width: 100%"
                v-decorator="[
                  'SmsType',
                  {
                    initialValue: '1',rules: [{ required: false, message: $t('AlarmTips.Voice.SMSFactory') }],
                  },
                ]"
            >
              <a-select-option value="1">
                {{$t('AlarmTips.Voice.AliYunSms')}}
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Voice.AccessKeyId')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['AccessKeyId', {rules: [{ required: true, message: $t('AlarmTips.Voice.AccessKeyId'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Voice.AccessKeySecret')}}&nbsp;
            </span>
            <a-input  autocomplete="autocomplete"
                      v-decorator="['AccessKeySecret', {rules: [{ required: true, message: $t('AlarmTips.Voice.AccessKeySecret'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Voice.TemplateCode')}}&nbsp;
            </span>
            <a-input   autocomplete="autocomplete"
                       v-decorator="['TemplateCode', {rules: [{ required: true, message: $t('AlarmTips.Voice.TemplateCode'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Voice.EveryDayCount')}}&nbsp;
            </span>
            <a-input  style="width: 100%" autocomplete="autocomplete"
                      v-decorator="['EveryDayCount', {rules: [{ required: true, message: $t('AlarmTips.Voice.EveryDayCount'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Voice.SendAlarmLevel')}}&nbsp;
            </span>
            <a-select
                mode="multiple"
                :allowClear="true"
                style="width: 100%"
                v-model="showAlarmVoiceItems"
            >
              <a-select-option value='0'>{{$t('dataModel.alarm.Tips')}}</a-select-option>
              <a-select-option value='1'>{{$t('dataModel.alarm.Minor')}}</a-select-option>
              <a-select-option value='2'>{{$t('dataModel.alarm.Importance')}}</a-select-option>
              <a-select-option value='3'>{{$t('dataModel.alarm.Urgency')}}</a-select-option>
              <a-select-option value='4'>{{$t('dataModel.alarm.Deadly')}}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
             <span slot="label">
              {{$t('AlarmTips.Voice.SendPhoneNumbers')}}&nbsp;
                <a-tooltip :title="$t('AlarmTips.Voice.SendPhoneNumbersTips')">
                <a-icon type="question-circle-o" />
              </a-tooltip>
            </span>
            <a-input  style="width: 100%" autocomplete="autocomplete"
                      v-decorator="['SendPhoneNumbers', {rules: [{ required: true, message: $t('AlarmTips.Voice.SendPhoneNumbers'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item style="margin-left: 100px" :label-col="{ span: 3}" :wrapper-col="{ span: 7}" >
            <a-button type="primary" @click="UpdateAlarmNotice('IhiyiVoice')">
              {{$t('AlarmTips.Save')}}
            </a-button>
            <a-button type="primary" style="margin-left: 10px" @click="TestVoice">
              {{$t('AlarmTips.Test')}}
            </a-button>
          </a-form-item>
        </a-form>
      </a-tab-pane>
    </a-tabs>
    </a-spin>
  </a-card>
</template>

<script>
import {TestSend,GetAlarmNoticeByType,UpdateAlarmNoticeByType} from "@/services/alarmNotice";

export default {
  name: "AlarmTips",
  i18n: require('../../i18n/language'),
  data() {
    return {
      labelCol: { span: 2 },
      WEBSpeechEnable:true,
      AlarmWindowsEnable:true,
      AlarmWindowsAutoClose:true,
      WEBSpeechSpeed:1,
      tabKey:"1",
      messageShowLoad:false,
      EmailParams:null,
      DingTalkList:[],
      showWindowsAlarmItems:[],
      showAlarmItems:[],
      showAlarmVoiceItems:[],
      marks: {
        0: '0.1',
        26: '1',
        37: '2.01',
        100: {
          style: {
            color: '#f50',
          },
        },
      },
      wrapperCol: { span: 7 },
      AlarmWindowsFrom: this.$form.createForm(this),
      SpeechFrom: this.$form.createForm(this),
      EmailForm: this.$form.createForm(this),
      weChatForm: this.$form.createForm(this),
      PhoneForm: this.$form.createForm(this),
      VoiceForm: this.$form.createForm(this),
    };
  },
  components: {

  },
  created(){
    if(this.tabKey==1)
    {
      let getSpeech = localStorage.getItem("Speech")
      if((getSpeech=="null")||(getSpeech==null)||(getSpeech==""))
      {
        this.WEBSpeechEnable = true
        this.WEBSpeechSpeed = 1
      }
      else
      {
        let getSpeechJson = JSON.parse(getSpeech)
        this.WEBSpeechEnable = getSpeechJson.enable
        this.WEBSpeechSpeed = getSpeechJson.speed
      }
    }

    let AlarmWindow = localStorage.getItem("AlarmWindow")

    if((AlarmWindow=="null")||(AlarmWindow==null)||(AlarmWindow==""))
    {
      this.AlarmWindowsEnable = true
      this.AlarmWindowsAutoClose= true
    }
    else
    {
      AlarmWindow = JSON.parse(AlarmWindow)
      this.AlarmWindowsEnable = AlarmWindow.enable
      this.showWindowsAlarmItems = AlarmWindow.Level
      this.AlarmWindowsAutoClose = AlarmWindow.isClose
    }
  },
  methods: {
    DelDingTalk(index){
      this.DingTalkList.splice(index,1)
    },
    AddDingTalk(){
      let talk={
        "IsEnable":true,
        "Webhook":"",
        "Secret":""
      }
      if ((this.DingTalkList instanceof Array)==false)
      {
        this.DingTalkList=[]
      }
      this.DingTalkList.push(talk)
    },
    GetAlarmNotice(AlarmNoticeType){
      let _t = this
      _t.dataSource = []
      const params = {
        type:AlarmNoticeType
      }
      this.messageShowLoad=true
      GetAlarmNoticeByType(params).then(function (res){
        if(res.data.code==0)
        {
          let EmailParams={}
          try{
            EmailParams =JSON.parse(res.data.list.AlarmNoticeParams)
            if(AlarmNoticeType=="Mail"){
              _t.EmailForm.setFieldsValue(
                  {
                    MailEnable:EmailParams.IsEnable,
                    MailServerIP:EmailParams.MailServerIP,
                    MailServerPort:EmailParams.MailServerPort.toString(),
                    MailSendUser:EmailParams.MailSendUser,
                    MailSendPassword:EmailParams.MailSendPassword,
                    MailSendAddress:EmailParams.MailSendAddress,
                    MailTo:EmailParams.MailTo,
                    TLS:EmailParams.TLS,
                    MailSendUserName:EmailParams.MailSendUserName,
                    MailSendSubject:EmailParams.MailSendSubject
                  })
            }
            else if(AlarmNoticeType=="weChat"){
              _t.weChatForm.setFieldsValue(
                  {
                    WeChatEnable:EmailParams.IsEnable,
                    EnterpriseID:EmailParams.EnterpriseID,
                    AgentId:EmailParams.AgentId.toString(),
                    Secret:EmailParams.Secret,
                  })
            }
            else if(AlarmNoticeType=="dingTalk"){
              _t.DingTalkList = EmailParams
            }
            else if(AlarmNoticeType=="Phone"){
              _t.PhoneForm.setFieldsValue(
                  {
                    PhoneNoticeEnable:EmailParams.IsEnable,
                    SendPhoneNumbers:EmailParams.PhoneNumbers,
                    EveryDayCount:EmailParams.EveryDayCount.toString(),
                    AccessKeyId:EmailParams.AliYunSms.AccessKeyId,
                    AccessKeySecret:EmailParams.AliYunSms.AccessKeySecret,
                    SignName:EmailParams.AliYunSms.SignName,
                    TemplateCode:EmailParams.AliYunSms.TemplateCode,
                  })
              _t.showAlarmItems = EmailParams.SendAlarmLevel
            }
            else if(AlarmNoticeType=="IhiyiVoice"){
              _t.VoiceForm.setFieldsValue(
                  {
                    PhoneNoticeEnable:EmailParams.IsEnable,
                    SendPhoneNumbers:EmailParams.PhoneNumbers,
                    EveryDayCount:EmailParams.EveryDayCount.toString(),
                    AccessKeyId:EmailParams.IhuyiVoice.AccessKeyId,
                    AccessKeySecret:EmailParams.IhuyiVoice.AccessKeySecret,
                    TemplateCode:EmailParams.IhuyiVoice.TemplateCode,
                  })
              _t.showAlarmVoiceItems = EmailParams.SendAlarmLevel
            }
          }catch (e){
            console.log(e)
          }
        }
        _t.messageShowLoad=false
      }).catch(function(e){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    UpdateAlarmNotice(AlarmNoticeType){
      let _t = this
      let params = {}
      if (AlarmNoticeType=="Mail"){
        this.EmailForm.validateFields((err) => {
          if (!err) {
            params.IsEnable = this.EmailForm.getFieldValue('MailEnable')
            params.MailServerIP = this.EmailForm.getFieldValue('MailServerIP')
            params.MailServerPort = parseInt(this.EmailForm.getFieldValue('MailServerPort'))?parseInt(this.EmailForm.getFieldValue('MailServerPort')):25
            params.MailSendUser = this.EmailForm.getFieldValue('MailSendUser')
            params.MailSendPassword = this.EmailForm.getFieldValue('MailSendPassword')
            params.MailSendAddress = this.EmailForm.getFieldValue('MailSendAddress')
            params.MailTo = this.EmailForm.getFieldValue('MailTo')
            params.TLS = this.EmailForm.getFieldValue('TLS')
            params.MailSendUserName = this.EmailForm.getFieldValue('MailSendUserName')
            params.MailSendSubject = this.EmailForm.getFieldValue('MailSendSubject')
          }
        })
      }
      else if(AlarmNoticeType=="weChat"){
        this.weChatForm.validateFields((err) => {
              if (!err) {
                params.IsEnable = this.weChatForm.getFieldValue('WeChatEnable')
                params.EnterpriseID = this.weChatForm.getFieldValue('EnterpriseID')
                params.AgentId = parseInt(this.weChatForm.getFieldValue('AgentId'))?parseInt(this.weChatForm.getFieldValue('AgentId')):0
                params.Secret = this.weChatForm.getFieldValue('Secret')
              }
        })
      }
      else if(AlarmNoticeType=="Phone"){
        this.PhoneForm.validateFields((err) => {
              if (!err) {
                params.IsEnable = this.PhoneForm.getFieldValue('PhoneNoticeEnable')
                params.SmsType = parseInt(this.PhoneForm.getFieldValue('SmsType'))?parseInt(this.PhoneForm.getFieldValue('SmsType')):1
                params.AliYunSms={
                  AccessKeyId : this.PhoneForm.getFieldValue('AccessKeyId'),
                  AccessKeySecret : this.PhoneForm.getFieldValue('AccessKeySecret'),
                  SignName : this.PhoneForm.getFieldValue('SignName'),
                  TemplateCode : this.PhoneForm.getFieldValue('TemplateCode'),
                }
                params.SendAlarmLevel = this.showAlarmItems
                params.PhoneNumbers = this.PhoneForm.getFieldValue('SendPhoneNumbers')
                params.EveryDayCount = parseInt(this.PhoneForm.getFieldValue('EveryDayCount'))?parseInt(this.PhoneForm.getFieldValue('EveryDayCount')):10
              }
        })
      }
      else if (AlarmNoticeType=="dingTalk"){
            params=this.DingTalkList
          }
      else if(AlarmNoticeType=="IhiyiVoice"){
        this.VoiceForm.validateFields((err) => {
          if (!err) {
            params.IsEnable = this.VoiceForm.getFieldValue('PhoneNoticeEnable')
            params.VoiceType = parseInt(this.VoiceForm.getFieldValue('SmsType'))?parseInt(this.VoiceForm.getFieldValue('SmsType')):1
            params.IhuyiVoice={
              AccessKeyId : this.VoiceForm.getFieldValue('AccessKeyId'),
              AccessKeySecret : this.VoiceForm.getFieldValue('AccessKeySecret'),
              TemplateCode : this.VoiceForm.getFieldValue('TemplateCode'),
            }
            params.SendAlarmLevel = this.showAlarmVoiceItems
            params.PhoneNumbers = this.VoiceForm.getFieldValue('SendPhoneNumbers')
            params.EveryDayCount = parseInt(this.VoiceForm.getFieldValue('EveryDayCount'))?parseInt(this.VoiceForm.getFieldValue('EveryDayCount')):10
          }
        })
      }
      const updateParams = {
        type:AlarmNoticeType,
        params:JSON.stringify(params)
      }
      this.messageShowLoad=true
      UpdateAlarmNoticeByType(updateParams).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t('AlarmTips.SaveSuccess'))
        }
        else
        {
          _t.$message.error(_t.$t('AlarmTips.SaveFailed'))
        }
        _t.messageShowLoad=false
      }).catch(function(e){

        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    TestEmail(){
      let _t = this

      this.EmailForm.validateFields((err) => {
        if (!err) {

          let params = {}

            params.IsEnable = this.EmailForm.getFieldValue('MailEnable')
            params.MailServerIP = this.EmailForm.getFieldValue('MailServerIP')
            params.MailServerPort = parseInt(this.EmailForm.getFieldValue('MailServerPort'))
            params.MailSendUser = this.EmailForm.getFieldValue('MailSendUser')
            params.MailSendPassword = this.EmailForm.getFieldValue('MailSendPassword')
            params.MailSendAddress = this.EmailForm.getFieldValue('MailSendAddress')
            params.MailTo = this.EmailForm.getFieldValue('MailTo')
            params.TLS = this.EmailForm.getFieldValue('TLS')
            params.MailSendUserName = this.EmailForm.getFieldValue('MailSendUserName')
            params.MailSendSubject = this.EmailForm.getFieldValue('MailSendSubject')


          const updateParams = {
            type:"Mail",
            params:JSON.stringify(params)
          }
          this.messageShowLoad=true
          TestSend(updateParams).then(function (res){
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t('AlarmTips.Email.SendSuccess'))
            }
            else
            {
              _t.$message.error(_t.$t('AlarmTips.Email.SendFailed'))
            }
            _t.messageShowLoad=false
          }).catch(function(e){

            _t.messageShowLoad=false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })


    },
    TestWeChat(){
      let _t = this

      this.weChatForm.validateFields((err) => {
        if (!err) {

          let params = {}

          params.IsEnable = this.weChatForm.getFieldValue('WeChatEnable')
          params.EnterpriseID = this.weChatForm.getFieldValue('EnterpriseID')
          params.AgentId = parseInt(this.weChatForm.getFieldValue('AgentId'))
          params.Secret = this.weChatForm.getFieldValue('Secret')

          const updateParams = {
            type:"weChat",
            params:JSON.stringify(params)
          }
          this.messageShowLoad=true
          TestSend(updateParams).then(function (res){
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t('AlarmTips.Email.SendSuccess'))
            }
            else if(res.data.code==-6)
            {
              _t.$message.error(res.data.msg)
            }
            else
            {
              _t.$message.error(_t.$t('AlarmTips.Email.SendFailed'))
            }
            _t.messageShowLoad=false
          }).catch(function(e){

            _t.messageShowLoad=false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })


    },
    TestSms(){
      let _t = this

      this.PhoneForm.validateFields((err) => {
        if (!err) {

          let params = {}

          params.IsEnable = this.PhoneForm.getFieldValue('PhoneNoticeEnable')
          params.SmsType = parseInt(this.PhoneForm.getFieldValue('SmsType'))
          params.AliYunSms={
            AccessKeyId : this.PhoneForm.getFieldValue('AccessKeyId'),
            AccessKeySecret : this.PhoneForm.getFieldValue('AccessKeySecret'),
            SignName : this.PhoneForm.getFieldValue('SignName'),
            TemplateCode : this.PhoneForm.getFieldValue('TemplateCode'),
          }
          params.SendAlarmLevel = this.showAlarmItems
          params.PhoneNumbers = this.PhoneForm.getFieldValue('SendPhoneNumbers')
          params.EveryDayCount = parseInt(this.PhoneForm.getFieldValue('EveryDayCount'))

          const updateParams = {
            type:"Phone",
            params:JSON.stringify(params)
          }
          this.messageShowLoad=true
          TestSend(updateParams).then(function (res){
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t('AlarmTips.Email.SendSuccess'))
            }
            else if(res.data.code==-3)
            {
              _t.$message.error(res.data.msg)
            }
            else
            {
              _t.$message.error(_t.$t('AlarmTips.Email.SendFailed'))
            }
            _t.messageShowLoad=false
          }).catch(function(e){

            _t.messageShowLoad=false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })


    },
    TestDingTalk(index){
      let _t = this
      const updateParams = {
        type:"dingTalk",
        index:index,
        params:JSON.stringify(this.DingTalkList)
      }
      this.messageShowLoad=true
      TestSend(updateParams).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t('AlarmTips.Email.SendSuccess'))
        }
        else if(res.data.code==-2)
        {
          _t.$message.error(res.data.msg)
        }else if(res.data.code==-3){
          _t.$message.error(_t.$t('AlarmTips.dingTalk.ParamError'))
        }else {
          _t.$message.error(_t.$t('AlarmTips.Email.SendFailed'))
        }
        _t.messageShowLoad=false
      }).catch(function(e){

        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    saveSpeech(){
      let setSpeech = {
        enable:this.WEBSpeechEnable,
        speed:this.WEBSpeechSpeed,
      }
      localStorage.setItem("Speech",JSON.stringify(setSpeech))
      this.$message.success(this.$t('AlarmTips.SaveSuccess'))
    },
    saveAlarmWindows(){
      let setSpeech = {
        enable:this.AlarmWindowsEnable,
        Level:this.showWindowsAlarmItems,
        isClose:this.AlarmWindowsAutoClose
      }
      localStorage.setItem("AlarmWindow",JSON.stringify(setSpeech))
      this.$message.success(this.$t('AlarmTips.SaveSuccess'))
    },
    callback(key) {
      this.tabKey = key
      if(key==1)
      {
        let getSpeech = localStorage.getItem("Speech")
        if((getSpeech=="null")||(getSpeech==null)||(getSpeech==""))
        {
          this.WEBSpeechEnable = true
          this.WEBSpeechSpeed = 1
        }
        else
        {
          let getSpeechJson = JSON.parse(getSpeech)
          this.WEBSpeechEnable = getSpeechJson.enable
          this.WEBSpeechSpeed = getSpeechJson.speed
        }
      }
      else if(key == 2)
      {
        this.GetAlarmNotice("Mail")
      }
      else if(key == 3)
      {
        this.GetAlarmNotice("Phone")
      }
      else if(key == 4){
        this.GetAlarmNotice("dingTalk")
      }
      else if(key == 5){
        this.GetAlarmNotice("weChat")
      }
      else if(key == 7){
        this.GetAlarmNotice("IhiyiVoice")
      }
    },
    TestVoice(){
      let _t = this
      this.VoiceForm.validateFields((err) => {
        if (!err) {
          let params = {}
          params.IsEnable = this.VoiceForm.getFieldValue('PhoneNoticeEnable')
          params.VoiceType = parseInt(this.VoiceForm.getFieldValue('SmsType'))?parseInt(this.VoiceForm.getFieldValue('SmsType')):1
          params.IhuyiVoice={
            AccessKeyId : this.VoiceForm.getFieldValue('AccessKeyId'),
            AccessKeySecret : this.VoiceForm.getFieldValue('AccessKeySecret'),
            TemplateCode : this.VoiceForm.getFieldValue('TemplateCode'),
          }
          params.SendAlarmLevel = this.showAlarmVoiceItems
          params.PhoneNumbers = this.VoiceForm.getFieldValue('SendPhoneNumbers')
          params.EveryDayCount = parseInt(this.VoiceForm.getFieldValue('EveryDayCount'))?parseInt(this.VoiceForm.getFieldValue('EveryDayCount')):10
          const updateParams = {
            type:"IhiyiVoice",
            params:JSON.stringify(params)
          }
          this.messageShowLoad=true
          TestSend(updateParams).then(function (res){
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t('AlarmTips.Email.SendSuccess'))
            }
            else if(res.data.code==-3)
            {
              _t.$message.error(res.data.msg)
            }
            else
            {
              _t.$message.error(_t.$t('AlarmTips.Email.SendFailed'))
            }
            _t.messageShowLoad=false
          }).catch(function(e){

            _t.messageShowLoad=false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
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

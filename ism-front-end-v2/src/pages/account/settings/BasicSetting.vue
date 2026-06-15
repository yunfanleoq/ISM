<template>
  <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
  <div class="account-settings-info-view">
    <a-row :gutter="16" type="flex" justify="center">
      <a-col :order="isMobile ? 2 : 1" :md="24" :lg="16">

        <a-form :form="EditForm"  layout="vertical">
          <a-form-item
              :label="$t('account.settings.basic.name')"
          >
            <a-input  v-decorator="[
                  'name',
                  {
                    rules: [{ required: true, message: $t('account.settings.basic.name-message') }],
                  },
                ]"
                      :placeholder="$t('account.settings.basic.name-message')" />
          </a-form-item>
          <a-form-item
              :label="$t('account.settings.basic.position')"
          >
            <a-input
                v-decorator="[
                  'job',
                  {
                    rules: [{ required: true, message: $t('account.settings.basic.position-message') }],
                  },
                ]"
                :placeholder="$t('account.settings.basic.position-message')" />
          </a-form-item>
          <a-form-item
              :label="$t('account.settings.basic.phone')"
          >
            <a-input
                v-decorator="[
                  'phone',
                  {
                    rules: [{ required: true, message: $t('account.settings.basic.phone-message') }],
                  },
                ]"
                :placeholder="$t('account.settings.basic.phone-message')" />
          </a-form-item>
          <a-form-item
            :label="$t('account.settings.basic.email')"
          >
            <a-input
                v-decorator="[
                  'email',
                  {
                    rules: [
                        { required: true, message: $t('account.settings.basic.email-message') },
                        {
                          type: 'email',
                          message:  $t('account.settings.basic.email-invalid'),
                        },

                        ],
                  },
                ]"
                placeholder="example@ism.com"/>
          </a-form-item>
          <a-form-item
              :label="$t('account.settings.basic.profile')"
          >
            <a-textarea rows="4"
                        v-decorator="[
                  'profile',
                  {
                    rules: [{ required: true, message: $t('account.settings.basic.profile-message') }],
                  },
                ]"
                        :placeholder="$t('account.settings.basic.profile-message')"/>
          </a-form-item>

          <a-form-item>
            <a-button type="primary" @click="SetUserInfo">{{ $t('account.settings.basic.update') }}</a-button>
          </a-form-item>
        </a-form>

      </a-col>
      <a-col :order="1" :md="24" :lg="8" :style="{ minHeight: '180px' }">
        <div class="ant-upload-preview" @click="$refs.modal.edit(1)" >
          <a-icon type="cloud-upload-o" class="upload-icon"/>
          <div class="mask">
            <a-icon type="plus" />
          </div>
          <img :src="option.img"/>
        </div>
      </a-col>
    </a-row>
    <avatar-modal ref="modal" @UploadAvatarFinish="setAvatar"/>
  </div>
  </a-spin>
</template>

<script>
import AvatarModal from './AvatarModal'
import {SetUserInfo,GetUserInfo} from "@/services/user";
import {mapGetters,mapMutations} from 'vuex'

export default {
  name: 'BasicSetting',
  i18n: require('../../../i18n/language'),
  components: {
    AvatarModal
  },
  data () {
    return {
      // cropper
      EditForm:this.$form.createForm(this),
      messageShowLoad:false,
      preview: {},
      isMobile:false,
      option: {
        img: '/avatar2.jpg',
        info: true,
        size: 1,
        outputType: 'jpeg',
        canScale: false,
        autoCrop: true,
        // 只有自动截图开启 宽度高度才生效
        autoCropWidth: 180,
        autoCropHeight: 180,
        fixedBox: true,
        // 开启宽度和高度比例
        fixed: true,
        fixedNumber: [1, 1]
      }
    }
  },
  computed: {
    ...mapGetters('account', ['user']),
  },
  created(){
    this.GetUserInfo()
  },
  methods: {
    ...mapMutations('account', ['setUser']),
    GetUserInfo(){
      let _t = this
      this.messageShowLoad = true
      GetUserInfo().then(function (res){
        if(res.data.code==200)
        {
          let item = res.data.info
          _t.option.img =item.avatar
          _t.EditForm.setFieldsValue(
              {
                job:item.job,
                name:item.name,
                phone:item.phone,
                profile:item.profile,
                email:item.email
              })
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
      })

    },
    SetUserInfo(){
      let _t = this
      this.EditForm.validateFields((err) => {
        if (!err) {
          let params = {
            job:this.EditForm.getFieldValue('job'),
            name:this.EditForm.getFieldValue('name'),
            phone:this.EditForm.getFieldValue('phone'),
            profile:this.EditForm.getFieldValue('profile'),
            email:this.EditForm.getFieldValue('email'),
          }
          _t.messageShowLoad = true
          SetUserInfo(params).then(function (res){
            if(res.data.code==200)
            {
              _t.$message.success( _t.$t('account.settings.basic.UserUpdateSuccess'))
            }
            else
            {
              _t.$message.error( _t.$t('account.settings.basic.UserUpdateFailed'))
            }
          }).catch(function (error) {
            _t.messageShowLoad = false
            _t.$message.error( _t.$t('account.settings.basic.UserUpdateFailed'))
          }).finally(function (error) {
            _t.messageShowLoad = false
          })
        }
      })
    },
    setAvatar (url) {
      this.option.img = url
      this.user.avatar = this.option.img
      this.setUser(this.user)
    }
  }
}
</script>

<style lang="less" scoped>
.ant-form-item{
  margin-bottom: 2px;
}
  .avatar-upload-wrapper {
    height: 200px;
    width: 100%;
  }

  .ant-upload-preview {
    position: relative;
    margin: 0 auto;
    width: 100%;
    max-width: 180px;
    border-radius: 50%;
    box-shadow: 0 0 4px #ccc;

    .upload-icon {
      position: absolute;
      top: 0;
      right: 10px;
      font-size: 1.4rem;
      padding: 0.5rem;
      background: rgba(222, 221, 221, 0.7);
      border-radius: 50%;
      border: 1px solid rgba(0, 0, 0, 0.2);
    }
    .mask {
      opacity: 0;
      position: absolute;
      background: rgba(0,0,0,0.4);
      cursor: pointer;
      transition: opacity 0.4s;

      &:hover {
        opacity: 1;
      }

      i {
        font-size: 2rem;
        position: absolute;
        top: 50%;
        left: 50%;
        margin-left: -1rem;
        margin-top: -1rem;
        color: #d6d6d6;
      }
    }

    img, .mask {
      width: 100%;
      max-width: 180px;
      height: 100%;
      border-radius: 50%;
      overflow: hidden;
    }
  }
</style>

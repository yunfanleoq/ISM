<template>
  <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
    <a-form :form="EditForm"  layout="vertical"  @submit="SetUserPassword">
    <a-form-item
        :label="$t('account.settings.security.OldPassword')"
    >
      <a-input-password   v-decorator="[
                  'OldPassword',
                  {
                    rules: [
                        { required: true, message: $t('account.settings.security.OldPassword') },
                        {
                          min: 8,
                          message: $t('account.settings.security.password-description'),
                        },
                       ],
                  },
                ]" />
    </a-form-item>
    <a-form-item
        :label="$t('account.settings.security.NewPassword')"
    >
      <a-input-password
        v-decorator="[
          'NewPassword',
          {
            rules: [
              {
                required: true,
                message: $t('account.settings.security.NewPassword'),
              },
              {
                min: 8,
                message: $t('account.settings.security.password-description'),
              },
              {
                validator: validateToNextPassword,
              },
            ],
          },
        ]"
                :placeholder="$t('account.settings.security.NewPassword')" />
    </a-form-item>
    <a-form-item
        :label="$t('account.settings.security.ConfirmNewPassword')"
    >
      <a-input-password
          v-decorator="[
          'ConfirmPassword',
          {
            rules: [
              {
                required: true,
                message: $t('account.settings.security.ConfirmNewPassword-description'),
              },
              {
                min: 8,
                message: $t('account.settings.security.password-description'),
              },
              {
                validator: compareToFirstPassword,
              },
            ],
          },
        ]"
          type="password"
          @blur="handleConfirmBlur"

          :placeholder="$t('account.settings.security.ConfirmNewPassword')" />
    </a-form-item>

    <a-form-item>
      <a-button type="primary" html-type="submit">{{ $t('account.settings.security.modify') }}</a-button>
    </a-form-item>
  </a-form>
  </a-spin>
</template>

<script>
import md5 from 'js-md5';

import {SetUserPassword} from "@/services/user";
export default {
  name: 'SecuritySetting',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      isMobile:false,
      messageShowLoad:false,
      EditForm:this.$form.createForm(this)
    }
  },
  methods: {
    SetUserPassword(){
      let _t = this
      this.EditForm.validateFields((err) => {
        if (!err) {
          let params = {
            oldPassword:md5(this.EditForm.getFieldValue('OldPassword')),
            newPassword:md5(this.EditForm.getFieldValue('ConfirmPassword')),
          }
          _t.messageShowLoad = true
          SetUserPassword(params).then(function (res){
            if(res.data.code==200)
            {
              _t.$message.success( _t.$t('account.settings.security.PasswordUpdateSuccess'))
            }
            else
            {
              _t.$message.error( _t.$t('account.settings.security.PasswordUpdateFailed'))
            }
          }).catch(function (error) {
            _t.messageShowLoad = false
            _t.$message.error( _t.$t('account.settings.security.PasswordUpdateFailed'))
          }).finally(function (error) {
            _t.messageShowLoad = false
          })
        }
      })
    },
    handleConfirmBlur(e) {
      const value = e.target.value;
      this.confirmDirty = this.confirmDirty || !!value;
    },
    compareToFirstPassword(rule, value, callback) {
      const form = this.EditForm;
      if (value && value !== form.getFieldValue('NewPassword')) {
        callback(this.$t('account.settings.security.PasswordNotSame'));
      } else {
        callback();
      }
    },
    validateToNextPassword(rule, value, callback) {
      const form = this.EditForm;
      if (value && this.confirmDirty) {
        form.validateFields(['ConfirmPassword'], { force: true });
      }
      callback();
    },
  }
}
</script>

<style lang="less" scoped>
.ant-form-item{
  margin-bottom: 20px;
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


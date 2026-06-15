<template>
  <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
  <div class="account-settings-info-view">
    <a-row :gutter="16" type="flex" justify="center">
      <a-col :order="isMobile ? 2 : 1" :md="24" :lg="16">
      <a-card>
        <a-form :form="EditForm"  layout="vertical">
          <a-form-item
              :label="$t('account.settings.basic.UserName')"
          >
            <a-input  v-decorator="[
                  'UserName',
                  {
                    rules: [{ required: true, message: $t('account.settings.basic.UserName-message') }],
                  },
                ]"
                      :placeholder="$t('account.settings.basic.name-message')" />
          </a-form-item>
          <a-form-item
              :label="$t('account.settings.basic.Role')"
          >
            <a-select
                v-decorator="['Role', {rules: [{ required: true, message: $t('account.settings.basic.Role'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in RoleList" :key="options.RoleId" :value="options.RoleId">
                <span v-if="options.RoleId=='Operator'">{{ $t("account.settings.UserList.RoleOperator") }}</span>
                <span v-if="options.RoleId=='User'">{{ $t("account.settings.UserList.RoleUser") }}</span>
              </a-select-option>

            </a-select>

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
            <a-button type="primary" @click="AddUser">{{ $t('account.settings.UserAdd.AddBtn') }}</a-button>
            <a-button type="default" @click="Back" style="margin-left: 10px">{{ $t('dataModel.back') }}</a-button>
          </a-form-item>
        </a-form>
      </a-card>
      </a-col>
    </a-row>
  </div>
  </a-spin>
</template>

<script>
import {AddSystemUser} from "@/services/user";
import {systemRolesList} from "@/services/system";
import md5 from 'js-md5';
import {mapState} from 'vuex'
export default {
  name: 'UserAdd',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      // cropper
      RoleList:[],
      EditForm:this.$form.createForm(this),
      messageShowLoad:false,
      preview: {},
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
  created(){
    this.SystemRolesList()
  },
  computed: {
    ...mapState('setting', ['isMobile', 'systemName'])
  },
  methods: {
    SystemRolesList(){
      let _t = this
      _t.RoleList = []
      this.messageShowLoad = true
      systemRolesList().then(function (res){
        if(res.data.code==0)
        {
          _t.RoleList = res.data.list
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
      })

    },
    AddUser(){
      let _t = this
      this.EditForm.validateFields((err) => {
        if (!err) {
          let params = {
            Username:this.EditForm.getFieldValue('UserName'),
            password:md5(this.EditForm.getFieldValue('NewPassword')),
            name:this.EditForm.getFieldValue('name'),
            phone:this.EditForm.getFieldValue('phone'),
            email:this.EditForm.getFieldValue('email'),
            job:this.EditForm.getFieldValue('job'),
            profile:this.EditForm.getFieldValue('profile'),
            role:this.EditForm.getFieldValue('Role'),
          }
          _t.messageShowLoad = true
          AddSystemUser(params).then(function (res){
            if(res.data.code==200)
            {
              _t.$message.success( _t.$t('account.settings.UserAdd.AddSuccess'))
              _t.$router.push('/Setting/UserManager')
            }
            else if(res.data.code==1009)
            {
              _t.$message.error( _t.$t('account.settings.UserAdd.AddUserExist'))
            }
            else{
              _t.$message.error( _t.$t('account.settings.UserAdd.AddFailed'))
            }
          }).catch(function (error) {
            _t.messageShowLoad = false
            _t.$message.error( _t.$t('account.settings.UserAdd.AddFailed'))
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
    Back(){
      this.$router.push('/Setting/UserManager')
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

<template>
  <div style="padding: 10px">
    <a-card class="project-list" :loading="loading" style="padding: 0px;min-height: 400px;" :bordered="false" :title="$t('displayModel3D.DisplayList')">
      <a slot="extra">
        <a-button type="primary" icon="plus" @click="visible=true;isEditStatus=0">
          {{$t('displayModel3D.AddModel')}}
        </a-button>
      </a>
      <a-list :grid="{ gutter: 16, column: 3 }" :dataSource="modelList" :pagination="pagination">
        <a-list-item slot="renderItem" slot-scope="item,index">
          <a-card id="display3dCard" hoverable style="width: 300px;border-radius: 2px 2px 0 0;">
            <template #cover>
              <vue-hover-mask>
                <img style="width: 298px;height: 220px;cursor: pointer" :src="item.DisplayImage==''?'/static/images/pcDefaultCover.jpg':item.DisplayImage" />
                <template v-slot:action>
                  <div style="margin-top: 20%" @click="updateImageIndex(index)">
                    <a-upload
                        progress="line"
                        name="file"
                        :multiple="true"
                        :action="uploadDisPlayUrl+'/'+item.uuid"
                        :showUploadList="false"
                        :beforeUpload="beforeUpload"
                        @change="afterUpload"
                    >
                      <a-row>
                        <span style="display:inline-block;height: 50px;width: 50px;border-radius: 50%;font-size: 25px;background-color: #e5e5e5"><a-icon style="color:#b2b2b2;margin-top: 10px" type="upload" /></span>
                      </a-row>
                      <a-row>
                        <span style="color:#fff;font-size: 15px">{{$t('displayModel3D.UploadImage')}}</span>
                      </a-row>
                    </a-upload>
                  </div>
                </template>
              </vue-hover-mask>
            </template>
            <template #actions>
              <router-link v-auth:role="`edit`" :to="`/DisPlay3DEditor/${item.uuid}`" target="_blank" style="font-size: 16px"><a-icon type="form" /> {{$t('displayModel3D.ModelDesign')}}</router-link>
              <router-link :to="`/DisPlay3DRunApp/${item.uuid}`" target="_blank" style="font-size: 16px"><a-icon type="play-circle" /> {{$t('displayModel3D.modelRun')}}</router-link>
              <a-dropdown>
                <a style="font-size: 16px" class="ant-dropdown-link" @click.prevent>
                  {{$t('displayModel3D.More')}}
                  <a-icon type="down" />
                </a>
                <template #overlay>
                  <a-menu>
                    <a-menu-item>
                      <a v-auth:role="`edit`" @click="ShowModelUser(item.uuid)" style="color: #565c64;font-size: 13px"><icon-font type="icon-authority" style="color: #565c64;font-size: 15px"/> {{$t('displayModel3D.Auth')}}</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a v-auth:role="`edit`" @click="copyFn(`/#/DisPlay3DRunApp/${item.uuid}`)" style="color: #565c64;font-size: 13px"><a-icon type="share-alt" /> {{$t('displayModel3D.ModelShare')}}</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a v-auth:role="`edit`" @click="showEditModel(item,index)" style="color: #565c64;font-size: 13px"><a-icon type="edit" /> {{$t('displayModel3D.ModelEdit')}}</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a v-auth:role="`delete`" @click="delModel(item.uuid,index)" style="color: #565c64;font-size: 13px"><a-icon type="delete" /> {{$t('displayModel3D.ModelDelete')}}</a>
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </template>
            <a-card-meta :title="item.name" :description="item.description">
            </a-card-meta>
          </a-card>
        </a-list-item>
      </a-list>
    </a-card>
    <a-modal v-model="visible" :title="isEditStatus===0 ? $t('displayModel3D.AddModel') : $t('displayModel3D.EditModelTitle')">
      <template slot="footer">
        <a-button type="primary" @click="doWhat">
          {{ $t('displayModel3D.ModelSure') }}
        </a-button>
        <a-button @click="handleCancel">
          {{ $t('displayModel3D.ModelCancel') }}
        </a-button>
      </template>
      <a-form :form="AddModelForm" :label-col="{ span: 3 }" :wrapper-col="{ span: 18 }">
        <a-form-item :label="$t('displayModel3D.ModelName')">
          <a-input
              v-decorator="['name', { rules: [{ required: true, message: $t('displayModel3D.ModelName') }] }]"
          />
        </a-form-item>
        <a-form-item :label="$t('displayModel3D.ModelDescription')">
          <a-textarea
              v-decorator="['description', { rules: [{ required: true, message: $t('displayModel3D.ModelDescription') }] }]"
          />
        </a-form-item>
        <a-form-item :label="$t('displayModel3D.ModelDescription')" style="display: none">
          <a-input
              v-decorator="['uuid', { rules: [{ required: false }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model="addAuthVisible" :footer="null" :title="$t('displayModel3D.AuthManage')" :bodyStyle="{height:'450px'}">
      <a-transfer
          :listStyle="{height:'400px',width:'200px'}"
          :titles="[$t('displayModel3D.UserList'), $t('displayModel3D.AuthUserList')]"
          :data-source="UserList"
          show-search
          :filter-option="filterOption"
          :target-keys="targetKeys"
          :render="item => item.title"
          @change="handleChange"
          @search="handleSearch"
      />
    </a-modal>
  </div>
</template>

<script>
import {
  displayModelAdd,
  displayModelList,
  displayModelDelete,
  displayModelEdit,
  DisplayModelAddUser,
  DisplayModelDelUser,
  GetDisplayModelUser
} from "@/services/displayModel";
import VueHoverMask from "@/components/VueHoverMask/VueHoverMask";
import { DISPLAYIMAGEUPLOAD } from "@/services/api";
import { SystemUserList } from "@/services/user";

export default {
  name: 'ISM3DAppList',
  i18n: require('@/i18n/language'),
  data () {
    return {
      editIndex: 0,
      isEditStatus: 0,
      UserList: [],
      mockData: [],
      targetKeys: [],
      selectDisplayUuid: "",
      showSearch: true,
      addAuthVisible: false,
      AddModelForm: this.$form.createForm(this),
      loading: false,
      uploadDisPlayUrl: DISPLAYIMAGEUPLOAD,
      visible: false,
      listIndex: -1,
      pagination: {
        hideOnSinglePage: true,
        showQuickJumper: true,
        pageSize: 3,
      },
      modelList: []
    }
  },
  components: {
    VueHoverMask,
  },
  created() {
    this.getModelList();
    this.SystemUserList();
  },
  watch: {
    '$route'() {
      this.modelList = [];
      this.SystemUserList();
      this.getModelList();
    }
  },
  methods: {
    filterOption(inputValue, option) {
      return option.description.indexOf(inputValue) > -1;
    },
    handleChange(targetKeys, direction, moveKeys) {
      this.targetKeys = [];
      let users = [];
      for (let i = 0; i < moveKeys.length; i++) {
        for (let j = 0; j < this.UserList.length; j++) {
          if (moveKeys[i] == this.UserList[j].key) {
            users.push(this.UserList[j]);
          }
        }
      }
      if (direction == "right") {
        this.AddModelUser(users);
      } else if (direction == "left") {
        this.DelModelUser(users);
      }
      this.targetKeys = targetKeys;
    },
    ShowModelUser(uuid) {
      this.targetKeys = [];
      this.selectDisplayUuid = uuid;
      this.addAuthVisible = true;
      this.GetDisplayModelUser();
    },
    handleSearch(dir, value) {
      console.log('search:', dir, value);
    },
    updateImageIndex(e) {
      this.listIndex = e;
    },
    beforeUpload() {
      this.$message.loading({ content: '上传中...', duration: 0 });
    },
    SystemUserList() {
      let _t = this;
      _t.UserList = [];
      SystemUserList().then(function (res) {
        if (res.data.code == 0) {
          for (let i = 0; i < res.data.List.length; i++) {
            const data = {
              key: i.toString(),
              title: res.data.List[i].name,
              uuid: res.data.List[i].uuid,
              Username: res.data.List[i].Username,
              chosen: Math.random() * 2 > 1,
            };
            _t.UserList.push(data);
          }
        }
      }).finally(function () {
        _t.messageShowLoad = false;
        _t.refIconLoading = false;
      });
    },
    copyFn(val) {
      val = location.protocol + "//" + location.hostname + ":" + location.port + val;
      let copyInput = document.createElement("input");
      copyInput.setAttribute("value", val);
      document.body.appendChild(copyInput);
      copyInput.select();
      try {
        let copyed = document.execCommand("copy");
        if (copyed) {
          document.body.removeChild(copyInput);
          this.$message.success(this.$t('device.CopySuccess'));
        }
      } catch {
        this.$message.error(this.$t('device.CopyFailed'));
      }
    },
    afterUpload(info) {
      if (info.file.status === 'done') {
        let result = info.file.response;
        this.$message.destroy();
        if (result.Code == 200) {
          this.modelList[this.listIndex].DisplayImage = result.Path;
          this.$message.success(this.$t('component.systemImageModel.uploadSuccess'));
        } else {
          this.$message.error(this.$t('component.systemImageModel.uploadFailed'));
        }
      }
    },
    doWhat(e) {
      if (this.isEditStatus == 0) {
        this.addModel(e);
      } else {
        this.editModel();
      }
    },
    addModel(e) {
      e.preventDefault();
      this.AddModelForm.validateFields((err) => {
        if (!err) {
          const params = {
            name: this.AddModelForm.getFieldValue('name'),
            description: this.AddModelForm.getFieldValue('description'),
            DisplayType: 2
          };
          let _t = this;
          displayModelAdd(params).then(function (res) {
            if (res.data.code == 4002) {
              _t.getModelList();
              _t.visible = false;
              _t.$message.success(_t.$t('displayModel3D.AddModelSuccess'), 3);
            } else if (res.data.code == 4001) {
              _t.$message.error(_t.$t('displayModel3D.ModelExist'), 3);
            } else if (res.data.code == 4003) {
              _t.$message.error(_t.$t('displayModel3D.AddModelFailed'), 3);
            } else if (res.data.code == 4005) {
              _t.$message.error(_t.$t('displayModel3D.AppOutCount'), 3);
            }
          });
        }
      });
    },
    showEditModel(item, index) {
      this.visible = true;
      this.isEditStatus = 1;
      this.editIndex = index;
      let _t = this;
      setTimeout(function () {
        _t.AddModelForm.setFieldsValue({
          name: item.name,
          uuid: item.uuid,
          description: item.description,
        });
      }, 300);
    },
    editModel() {
      this.AddModelForm.validateFields((err) => {
        if (!err) {
          const params = {
            uuid: this.AddModelForm.getFieldValue('uuid'),
            updateData: {
              name: this.AddModelForm.getFieldValue('name'),
              description: this.AddModelForm.getFieldValue('description'),
            }
          };
          let _t = this;
          displayModelEdit(params).then(function (res) {
            if (res.data.code == 200) {
              _t.visible = false;
              for (let i = 0; i < _t.modelList.length; i++) {
                if (_t.modelList[i].uuid == params.uuid) {
                  _t.modelList[i].name = params.updateData.name;
                  _t.modelList[i].description = params.updateData.description;
                  break;
                }
              }
              _t.$message.success(_t.$t('displayModel3D.EditModelSuccess'), 3);
            } else if (res.data.code == 500) {
              _t.$message.error(_t.$t('displayModel3D.EditModelFailed'), 3);
            }
          });
        }
      });
    },
    DelModelUser(users) {
      let _t = this;
      let uuid = [];
      for (let i = 0; i < users.length; i++) {
        uuid.push(users[i].uuid);
      }
      let params = {
        displayuuid: this.selectDisplayUuid,
        users: uuid
      };
      DisplayModelDelUser(params).then(function (res) {
        if (res.data.code == 0) {
        }
      }).finally(function () {
        _t.messageShowLoad = false;
        _t.refIconLoading = false;
      });
    },
    AddModelUser(users) {
      let _t = this;
      let params = {
        displayuuid: this.selectDisplayUuid,
        users: users
      };
      DisplayModelAddUser(params).then(function (res) {
        if (res.data.code == 0) {
        }
      }).finally(function () {
        _t.messageShowLoad = false;
        _t.refIconLoading = false;
      });
    },
    GetDisplayModelUser() {
      let _t = this;
      let params = {
        displayuuid: this.selectDisplayUuid
      };
      GetDisplayModelUser(params).then(function (res) {
        if (res.data.code == 200) {
          for (let i = 0; i < res.data.List.length; i++) {
            for (let k = 0; k < _t.UserList.length; k++) {
              if (res.data.List[i].suuid == _t.UserList[k].uuid) {
                _t.targetKeys.push(_t.UserList[k].key);
              }
            }
          }
        }
      }).finally(function () {
        _t.messageShowLoad = false;
        _t.refIconLoading = false;
      });
    },
    delModel(uuid, index) {
      let _t = this;
      this.$confirm({
        title: _t.$t('displayModel3D.DelModelConfirm'),
        okText: _t.$t('displayModel3D.ConfirmOk'),
        onOk() {
          const params = {
            displayUid: uuid
          };
          displayModelDelete(params).then(function (res) {
            if (res.data.code == 200) {
              _t.getModelList();
              _t.$message.success(_t.$t('displayModel3D.DelModelSuccess'), 3);
            } else if (res.data.code == 2004) {
              _t.$message.error(_t.$t('displayModel3D.modelHavedBind'), 3);
            } else {
              _t.$message.error(_t.$t('displayModel3D.DelModelFailed'), 3);
            }
          });
        },
        cancelText: _t.$t('displayModel3D.ConfirmCancel'),
        onCancel() {
        },
      });
    },
    handleCancel() {
      this.visible = false;
    },
    getModelList() {
      this.modelList = [];
      let _t = this;
      const params = {
        DisplayType: 2
      };
      this.loading = true;
      displayModelList(params).then(function (res) {
        if (res.data.list != null) {
          for (let i = 0; i < res.data.list.length; i++) {
            let tableData = {};
            tableData.name = res.data.list[i].name;
            tableData.description = res.data.list[i].description;
            tableData.uuid = res.data.list[i].displayUid;
            if (res.data.list[i].DisplayImage == "") {
              tableData.DisplayImage = "";
            } else {
              tableData.DisplayImage = res.data.list[i].DisplayImage;
            }
            _t.modelList.push(tableData);
            tableData = {};
          }
        }
        _t.loading = false;
      });
    },
  }
}
</script>

<style lang="less" scoped>

::v-deep #display3dCard .ant-card-body {
  padding: 0px;
  margin-top: 5px;
  border-top: 1px solid #f0f0f0;
  margin-bottom: 10px;
  zoom: 1;
}

::v-deep .vue-hover-mask {
  border-radius: 0px;
}

::v-deep #display3dCard .ant-card-meta-detail {
  margin-top: 5px;
  margin-left: 10px;
}

::v-deep #display3dCard .ant-card-meta-detail > div:not(:last-child) {
  margin-bottom: 0px;
}

::v-deep #display3dCard .ant-card-head {
  min-height: 28px;
  margin-bottom: -1px;
  padding: 0 8px;
  color: #000;
  font-weight: 500;
  font-size: 13px;
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
  border-radius: 2px 2px 0 0;
  zoom: 1;
}

::v-deep .ant-card-actions {
  background-color: #fff;
}
</style>

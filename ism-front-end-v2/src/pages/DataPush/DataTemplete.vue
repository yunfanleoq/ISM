<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="addVisible=true;isEdit=false" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
      <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
    </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="TempleteName" :pagination="pagination" :columns="columns" :data-source="dataSource">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="TempleteContent" slot-scope="text">
          <span class="">{{truncateString(text,0,120)}}<span v-if="text.length>120">.....</span></span>
        </div>
        <div slot="action" slot-scope="text, record">
          <a @click="GoToEdit(record)" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</a> |
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.TempleteUuid)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
    <a-drawer
        :title="isEdit?$t('ISMDataTemplete.editData'):$t('ISMDataTemplete.addData')"
        :width="800"
        :visible="addVisible"
        :body-style="{ paddingBottom: '80px' }"
        @close="onClose"
    >
      <a-form :form="PlanForm" layout="vertical" >
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item
                :label="$t('ISMDataTemplete.TempleteName')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['TempleteName', {rules: [{ required: true, message: $t('ISMDataTemplete.TempleteName'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item
                :label="$t('ISMDataTemplete.TempleteContent')"
            >
              <code-editor v-if="isCharge"
                           :dHeight="500"
                           :value="CodeContent"
                           language="javascript"
                           @input="changeTextarea"
              >
              </code-editor>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <div
          :style="{
          position: 'absolute',
          right: 0,
          bottom: 0,
          width: '100%',
          borderTop: '1px solid #e9e9e9',
          padding: '10px 16px',
          background: '#fff',
          textAlign: 'right',
          zIndex: 1,
        }"
      >
        <a-button  type="primary" :style="{ marginRight: '8px' }" v-if="!isEdit"  @click="AddScript()">
          {{$t('TaskPlan.TaskAdd')}}
        </a-button>
        <a-button  type="primary" :style="{ marginRight: '8px' }" v-if="isEdit"  @click="EditScript()">
          {{$t('TaskPlan.TaskEdit')}}
        </a-button>

        <a-button  @click="onClose">
          {{$t('device.CancelButton')}}
        </a-button>
      </div>
    </a-drawer>
  </a-card>
</template>

<script>
import codeEditor from '@/components/CodeEditor/index'
import moment from 'moment'
import {ExecSysScript,DisableSysScript} from "@/services/system";
import {formatDate} from "@/utils/common";
import {AddScript, CheckScript, DelScript, EditScript, GetScriptList} from "@/services/ismscripts";
import {AddTempleteData, DelTempleteData, EditTempleteData, GetTempleteDataList} from "../../services/datatemplete";
export default {
  name: 'DataTemplete',
  i18n: require('@/i18n/language'),
  data() {
    return {
      pagination: {
        pageSize: 15,
        showSizeChanger: true
      },
      ScriptType: 0,
      isCharge: true,
      CodeContent: "",
      isEdit: false,
      messageShowLoad: false,
      advanced: true,
      refIconLoading: false,
      columns: [
        {
          width: '10%',
          slotName: 'ISMDataTemplete.TempleteName',
          scopedSlots: {customRender: 'TempleteName', title: 'ISMDataTemplete.TempleteName'},
          dataIndex: 'TempleteName'
        },
        {
          slotName: 'ISMDataTemplete.TempleteContent',
          width: '40%',
          scopedSlots: {customRender: 'TempleteContent', title: 'ISMDataTemplete.TempleteContent'},
          dataIndex: 'TempleteContent',
        },
        {
          width: '10%',
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: {customRender: 'action', title: 'dataModel.modelTableOpt'}
        }
      ],
      dataSource: [],
      addVisible: false,
      error: '',
      editUuid: "",
      editVisible: false,
      PlanForm: this.$form.createForm(this),
      editForm: this.$form.createForm(this),
      textAreValue: "",
      that: this,
      value: 1
    }
  },
  components: {
    codeEditor,
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date, 'yyyy-MM-dd hh:mm:ss')
    }
  },
  mounted() {

  },
  activated() {
    this.GetTempleteList()
  },
  created() {

  },
  methods: {
    truncateString(str, start, length) {
      return str.substring(start, length);
    },
    changeTextarea(val) {
      this.CodeContent = val
    },
    onClose() {
      this.addVisible = false;
    },
    GetTempleteList() {
      let _t = this
      this.dataSource = []
      GetTempleteDataList().then(function (res) {
        _t.refIconLoading = false
        if (res.data.code == 200) {
          for(let i=0;i<res.data.list.length;i++){
            if(res.data.list[i].TempleteType==1||res.data.list[i].TempleteType==0)
            {
              _t.dataSource.push(res.data.list[i])
            }
          }
          _t.addVisible = false;
        } else if (res.data.code == 2001) {
          _t.$message.error(_t.$t('displayModel.ModelExist'), 3)
        } else if (res.data.code == 2003) {
          _t.$message.error(_t.$t('displayModel.AddModelFailed'), 3)
        }
      })
    },
    AddScript() {
      let _t = this
      this.PlanForm.validateFields((err) => {
        if (!err) {
          const params = {
            TempleteName: _t.PlanForm.getFieldValue('TempleteName'),
            TempleteContent: _t.CodeContent,
            TempleteType:1,
          }
          AddTempleteData(params).then(function (res) {
            if (res.data.code == 2002) {
              _t.GetTempleteList()
              _t.addVisible = false;
              _t.$message.success(_t.$t('ISMDataTemplete.AddSuccess'), 3)
            } else {
              _t.$message.error(_t.$t('ISMDataTemplete.AddFailed'), 3)
            }
          })
        }
      })
    },
    EditScript() {
      let _t = this
      this.PlanForm.validateFields((err) => {
        if (!err) {
          const params = {
            Uuid: _t.EditUUid,
            data: {
              TempleteName: _t.PlanForm.getFieldValue('TempleteName'),
              TempleteContent: _t.CodeContent,
            }
          }
          EditTempleteData(params).then(function (res) {
            if (res.data.code == 200) {
              _t.GetTempleteList()
              _t.addVisible = false;
              _t.$message.success(_t.$t('ISMDataTemplete.EditSuccess'), 3)
            } else {
              _t.$message.error(_t.$t('ISMDataTemplete.EditFailed'), 3)
            }
          })
        }
      })
    },
    GoToEdit(item) {
      let _t = this
      _t.isCharge = false
      this.isEdit = true
      this.addVisible = true
      _t.EditUUid = item.TempleteUuid
      _t.CodeContent = item.TempleteContent
      setTimeout(function () {
        _t.isCharge = true
        _t.PlanForm.setFieldsValue(
            {
              TempleteName: item.TempleteName,
            })
      }, 200)
    },
    refresh() {
      this.refIconLoading = true
      this.GetTempleteList()
    },
    deleteRecord(uuid) {
      let _t = this
      const params = {
        TempleteUuid: uuid
      }
      DelTempleteData(params).then(function (res) {
        if (res.data.code == 200) {
          _t.GetTempleteList()
          _t.addVisible = false;
          _t.$message.success(_t.$t('ISMDataTemplete.DelSuccess'), 3)
        } else {
          _t.$message.error(_t.$t('ISMDataTemplete.DelFailed'), 3)
        }
      })
    },
  }
}
</script>

<style lang="less" scoped>
::v-deep .search {
  margin-bottom: 54px;
}

::v-deep .ant-form-item {
  margin-bottom: 1px;
}

::v-deep .ant-row .ant-form-item {
  margin-bottom: 1px;
}

::v-deep .fold {
  width: calc(100% - 216px);
  display: inline-block
}
::v-deep .text-container {
  width: 10px; /* 设定容器宽度 */
  overflow: hidden; /* 隐藏溢出的内容 */
  text-overflow: ellipsis; /* 显示省略号 */
  white-space: nowrap; /* 防止文本换行 */
}
::v-deep .operator {
  margin-bottom: 18px;
}

@media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}

::v-deep .ant-form-item {
  margin-bottom: 10px;
}

::v-deep .ant-table-thead > tr > th {
  padding: 10px 10px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-thead > tr > th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
</style>

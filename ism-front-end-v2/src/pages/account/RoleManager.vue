<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="showAddModal" type="primary" icon="plus">{{ $t('roleManager.addRole') }}</a-button>
      <a-button @click="refresh()" type="default" icon="sync" :loading="refIconLoading">{{ $t('dataModel.refModel') }}</a-button>
    </a-space>
    <a-spin style="padding: 1px;" :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="RoleId" :pagination="pagination" :columns="columns" :data-source="roleList">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="RoleId" slot-scope="text">
          <a-tag :color="getRoleColor(text)">{{ text }}</a-tag>
        </div>
        <div slot="action" slot-scope="text, record">
          <a @click="showPermissionModal(record)"><a-icon type="safety" />{{ $t('roleManager.permission') }}</a>
          <a-divider type="vertical" />
          <a @click="showEditModal(record)"><a-icon type="edit" />{{ $t('roleManager.edit') }}</a>
          <a-divider type="vertical" />
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="RoleDel(record)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a style="color: #eb2f96"><a-icon type="delete" />{{ $t('dataModel.delete') }}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>

    <a-modal
      :title="isEdit ? $t('roleManager.editRole') : $t('roleManager.addRole')"
      :visible="roleModalVisible"
      @ok="handleRoleSubmit"
      @cancel="roleModalVisible = false"
      :confirm-loading="submitLoading"
    >
      <a-form :form="roleForm" layout="vertical">
        <a-form-item :label="$t('roleManager.roleId')">
          <a-input
            v-decorator="['RoleId', { rules: [{ required: true, message: $t('roleManager.roleIdRequired') }] }]"
            :disabled="isEdit"
            :placeholder="$t('roleManager.roleIdPlaceholder')"
          />
        </a-form-item>
        <a-form-item :label="$t('roleManager.roleName')">
          <a-input
            v-decorator="['RoleName', { rules: [{ required: true, message: $t('roleManager.roleNameRequired') }] }]"
            :placeholder="$t('roleManager.roleNamePlaceholder')"
          />
        </a-form-item>
        <a-form-item :label="$t('roleManager.roleDesc')">
          <a-textarea
            v-decorator="['Description']"
            :placeholder="$t('roleManager.roleDescPlaceholder')"
            :auto-size="{ minRows: 2, maxRows: 4 }"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      :title="$t('roleManager.permissionConfig') + ' - ' + currentRole.RoleName"
      :visible="permissionModalVisible"
      @ok="handlePermissionSubmit"
      @cancel="permissionModalVisible = false"
      width="600px"
      :confirm-loading="submitLoading"
    >
      <a-spin :spinning="permissionLoading">
        <a-tree
          v-model="checkedKeys"
          checkable
          :tree-data="permissionTree"
          :replaceFields="{ title: 'label', key: 'value', children: 'children' }"
          style="max-height: 400px; overflow-y: auto;"
        />
      </a-spin>
    </a-modal>
  </a-card>
</template>

<script>
import {
  systemRolesList,
  systemRoleAdd,
  systemRoleEdit,
  systemRoleDel,
  systemRolePermissions,
  systemRoleUpdatePermissions
} from '@/services/system'

const MENU_PERMISSIONS = [
  {
    label: '资源统计',
    value: '/dashboard',
  },
  {
    label: '数据仓库',
    value: '/DataWarehouse',
  },
  {
    label: '设备管理',
    value: '/DeviceLibraryConfig',
  },
  {
    label: '数据模型',
    value: '/DeviceModel',
    children: [
      { label: 'SNMP设备', value: '/DeviceModel/SnmpModel' },
      { label: 'Modbus设备', value: '/DeviceModel/ModbusModel' },
      { label: 'DLT645设备', value: '/DeviceModel/DLT645Model' },
      { label: 'CJT188设备', value: '/DeviceModel/CJT188Model' },
      { label: 'OPCUA设备', value: '/DeviceModel/OPCUAModel' },
      { label: 'MQTT设备', value: '/DeviceModel/MqttModel' },
      { label: '西门子S7设备', value: '/DeviceModel/SiemensS7Model' },
      { label: 'IEC104设备', value: '/DeviceModel/IEC104Model' },
      { label: 'IEC61850设备', value: '/DeviceModel/IEC61850Model' },
      { label: 'DNP3设备', value: '/DeviceModel/DNP3Model' },
      { label: 'BACnet设备', value: '/DeviceModel/BACnetModel' },
      { label: 'RESTFul设备', value: '/DeviceModel/RestFulModel' },
      { label: 'HJ212设备', value: '/DeviceModel/HJ212Model' },
      { label: '虚拟装置', value: '/DeviceModel/VirtualDevice' },
      { label: '自定义数据', value: '/DeviceModel/DeviceCustomData' },
      { label: '系统变量', value: '/DeviceModel/SystemData' },
    ]
  },
  {
    label: '应用管理',
    value: '/Application',
  },
  {
    label: '数字孪生',
    value: '/DigitalTwin',
  },
  {
    label: '视频管理',
    value: '/VideoManager',
    children: [
      { label: '分屏展示', value: '/VideoManager/videoScreen' },
      { label: '视频清单', value: '/VideoManager/videoList' },
      { label: '国标视频', value: '/VideoManager/GB28281List' },
    ]
  },
  {
    label: '实时告警',
    value: '/Real-timeAlarm',
  },
  {
    label: '告警策略',
    value: '/AlarmStrategy',
    children: [
      { label: '模型触发器', value: '/AlarmStrategy/ModelTrigger' },
      { label: '告警恢复', value: '/AlarmStrategy/AlarmRestoreMask' },
    ]
  },
  {
    label: '任务计划',
    value: '/TaskPlan',
  },
  {
    label: '系统脚本',
    value: '/ISMScripts',
  },
  {
    label: '数据接口',
    value: '/DataPush',
    children: [
      { label: '数据模版', value: '/DataPush/DataTemplete' },
      { label: 'IEC104数据模版', value: '/DataPush/IEC104DataTemplete' },
      { label: 'Modbus数据模版', value: '/DataPush/ModbusDataTemplete' },
      { label: '接口管理', value: '/DataPush/DataInterface' },
    ]
  },
  {
    label: '数据报表',
    value: '/Reporting',
    children: [
      { label: '告警报表', value: '/Reporting/AlarmHistory' },
      { label: '历史数据', value: '/Reporting/DataHistory' },
      { label: '自定义报表', value: '/Reporting/DiyReport' },
      { label: '报表模板', value: '/Reporting/DiyReportTemplete' },
    ]
  },
  {
    label: '网络中心',
    value: '/Network',
    children: [
      { label: '网络', value: '/Network/SystemNetwork' },
      { label: '组网', value: '/Network/ISMNetwork' },
    ]
  },
  {
    label: '数据库管理',
    value: '/DataBase',
    children: [
      { label: '实时数据库', value: '/DataBase/DbManager' },
      { label: '历史数据库', value: '/DataBase/HistoryManager' },
    ]
  },
  {
    label: '设置中心',
    value: '/Setting',
    children: [
      { label: '个人设置', value: '/Setting/Account' },
      { label: '用户管理', value: '/Setting/UserManager' },
      { label: '角色管理', value: '/Setting/RoleManager' },
      { label: '告警通知', value: '/Setting/AlarmTipsSetting' },
      { label: '系统参数', value: '/Setting/SystemParams' },
      { label: 'API令牌', value: '/Setting/AccessToken' },
    ]
  },
  {
    label: '系统日志',
    value: '/Journal',
    children: [
      { label: '操作日志', value: '/Journal/OperationJournal' },
    ]
  },
]

export default {
  name: 'RoleManager',
  i18n: require('@/i18n/language'),
  data() {
    return {
      pagination: { pageSize: 15, showSizeChanger: true },
      roleList: [],
      refIconLoading: false,
      messageShowLoad: false,
      submitLoading: false,
      roleModalVisible: false,
      permissionModalVisible: false,
      isEdit: false,
      currentRole: {},
      checkedKeys: [],
      permissionLoading: false,
      permissionTree: MENU_PERMISSIONS,
      roleForm: this.$form.createForm(this),
      columns: [
        {
          width: '15%',
          slotName: 'roleManager.roleId',
          scopedSlots: { customRender: 'RoleId', title: 'roleManager.roleId' },
          dataIndex: 'RoleId',
        },
        {
          width: '20%',
          slotName: 'roleManager.roleName',
          scopedSlots: { customRender: 'serial', title: 'roleManager.roleName' },
          dataIndex: 'RoleName',
        },
        {
          width: '35%',
          slotName: 'roleManager.roleDesc',
          scopedSlots: { customRender: 'serial', title: 'roleManager.roleDesc' },
          dataIndex: 'Description',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: { customRender: 'action', title: 'dataModel.modelTableOpt' },
        }
      ],
    }
  },
  created() {
    this.GetRoleList()
  },
  methods: {
    refresh() {
      this.refIconLoading = true
      this.GetRoleList()
    },
    GetRoleList() {
      let _t = this
      _t.roleList = []
      _t.messageShowLoad = true
      systemRolesList().then(function(res) {
        if (res.data.code == 0) {
          _t.roleList = res.data.list
        }
      }).finally(function() {
        _t.messageShowLoad = false
        _t.refIconLoading = false
      })
    },
    getRoleColor(roleId) {
      const colorMap = {
        'Admin': 'red',
        'Operator': 'blue',
        'User': 'green',
      }
      return colorMap[roleId] || 'orange'
    },
    showAddModal() {
      this.isEdit = false
      this.currentRole = {}
      this.roleModalVisible = true
      this.$nextTick(() => {
        this.roleForm.resetFields()
      })
    },
    showEditModal(record) {
      this.isEdit = true
      this.currentRole = { ...record }
      this.roleModalVisible = true
      this.$nextTick(() => {
        this.roleForm.setFieldsValue({
          RoleId: record.RoleId,
          RoleName: record.RoleName,
          Description: record.Description,
        })
      })
    },
    handleRoleSubmit() {
      let _t = this
      this.roleForm.validateFields((err, values) => {
        if (!err) {
          _t.submitLoading = true
          const params = {
            RoleId: values.RoleId,
            RoleName: values.RoleName,
            Description: values.Description || '',
          }
          const apiCall = _t.isEdit ? systemRoleEdit : systemRoleAdd
          apiCall(params).then(function(res) {
            if (res.data.code == 200 || res.data.code == 0) {
              _t.$message.success(_t.isEdit ? _t.$t('roleManager.editSuccess') : _t.$t('roleManager.addSuccess'))
              _t.roleModalVisible = false
              _t.GetRoleList()
            } else if (res.data.code == 1009) {
              _t.$message.error(_t.$t('roleManager.roleExist'))
            } else {
              _t.$message.error(_t.isEdit ? _t.$t('roleManager.editFailed') : _t.$t('roleManager.addFailed'))
            }
          }).catch(function() {
            _t.$message.error(_t.isEdit ? _t.$t('roleManager.editFailed') : _t.$t('roleManager.addFailed'))
          }).finally(function() {
            _t.submitLoading = false
          })
        }
      })
    },
    RoleDel(record) {
      let _t = this
      const params = { RoleId: record.RoleId }
      _t.messageShowLoad = true
      systemRoleDel(params).then(function(res) {
        if (res.data.code == 200 || res.data.code == 0) {
          _t.GetRoleList()
        } else {
          _t.$message.error(_t.$t('roleManager.delFailed'))
        }
      }).finally(function() {
        _t.messageShowLoad = false
      })
    },
    showPermissionModal(record) {
      this.currentRole = { ...record }
      this.checkedKeys = []
      this.permissionModalVisible = true
      this.permissionLoading = true
      let _t = this
      systemRolePermissions({ RoleId: record.RoleId }).then(function(res) {
        if (res.data.code == 0) {
          _t.checkedKeys = res.data.permissions || []
        }
      }).catch(function() {
        _t.$message.error(_t.$t('roleManager.loadPermissionFailed'))
      }).finally(function() {
        _t.permissionLoading = false
      })
    },
    handlePermissionSubmit() {
      let _t = this
      _t.submitLoading = true
      const params = {
        RoleId: _t.currentRole.RoleId,
        Permissions: _t.checkedKeys,
      }
      systemRoleUpdatePermissions(params).then(function(res) {
        if (res.data.code == 200 || res.data.code == 0) {
          _t.$message.success(_t.$t('roleManager.updatePermissionSuccess'))
          _t.permissionModalVisible = false
        } else {
          _t.$message.error(_t.$t('roleManager.updatePermissionFailed'))
        }
      }).catch(function() {
        _t.$message.error(_t.$t('roleManager.updatePermissionFailed'))
      }).finally(function() {
        _t.submitLoading = false
      })
    },
  }
}
</script>

<style lang="less" scoped>
.operator {
  margin-bottom: 18px;
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
  transition: background .3s ease;
}
</style>
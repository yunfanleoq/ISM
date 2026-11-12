<template>
  <a-card :bordered="false">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card title="字典类型" size="small">
          <div slot="extra">
            <a-button type="primary" size="small" icon="plus" @click="openType()">新增</a-button>
          </div>
          <a-table
            rowKey="ID"
            size="small"
            :loading="typeLoading"
            :columns="typeColumns"
            :data-source="types"
            :pagination="false"
            :customRow="typeRow"
            :rowClassName="typeRowClass"
          >
            <span slot="DictName" slot-scope="text, record">
              {{ text }}
              <a-tag v-if="record.Builtin === 1" color="blue">内置</a-tag>
            </span>
            <span slot="Status" slot-scope="text">
              <a-badge :status="text === 1 ? 'success' : 'default'" :text="text === 1 ? '启用' : '停用'" />
            </span>
            <span slot="action" slot-scope="text, record">
              <a @click.stop="openType(record)">编辑</a>
              <a-divider type="vertical" />
              <a-popconfirm v-if="record.Builtin !== 1" title="删除后，这个字典下的条目一并删除。确定吗？" @confirm="removeType(record)">
                <a style="color:#eb2f96">删除</a>
              </a-popconfirm>
              <span v-else style="color:#bbb">删除</span>
            </span>
          </a-table>
        </a-card>
      </a-col>
      <a-col :span="16">
        <a-card size="small">
          <span slot="title">{{ currentType ? currentType.DictName + '（' + currentType.DictCode + '）' : '字典项' }}</span>
          <div slot="extra">
            <a-button type="primary" size="small" icon="plus" :disabled="!currentType" @click="openItem()">新增条目</a-button>
          </div>
          <div v-if="currentType && currentType.Remark" style="margin-bottom: 8px; color: #888;">{{ currentType.Remark }}</div>
          <a-table
            rowKey="ID"
            size="small"
            :loading="itemLoading"
            :columns="itemColumns"
            :data-source="items"
            :pagination="false"
          >
            <span slot="ItemColor" slot-scope="text">
              <i v-if="text" :style="{display:'inline-block',width:'14px',height:'14px',borderRadius:'2px',background:text,verticalAlign:'middle',marginRight:'6px'}"></i>
              {{ text || '-' }}
            </span>
            <span slot="ItemLabel" slot-scope="text, record">
              <dict-tag :dict-code="record.DictCode" :value="record.ItemValue" />
            </span>
            <span slot="Status" slot-scope="text">
              <a-badge :status="text === 1 ? 'success' : 'default'" :text="text === 1 ? '启用' : '停用'" />
            </span>
            <span slot="action" slot-scope="text, record">
              <a @click="openItem(record)">编辑</a>
              <a-divider type="vertical" />
              <a-popconfirm v-if="record.Builtin !== 1" title="确定删除这条字典项吗？" @confirm="removeItem(record)">
                <a style="color:#eb2f96">删除</a>
              </a-popconfirm>
              <span v-else style="color:#bbb">删除</span>
            </span>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-modal :title="typeForm.ID ? '编辑字典' : '新增字典'" :visible="typeVisible" :confirm-loading="saving" @ok="saveType" @cancel="typeVisible=false">
      <a-form layout="vertical">
        <a-form-item label="字典编码" required>
          <a-input v-model="typeForm.DictCode" :disabled="typeForm.Builtin === 1" placeholder="字母开头，如 alarm_level" />
        </a-form-item>
        <a-form-item label="字典名称" required>
          <a-input v-model="typeForm.DictName" placeholder="例如 报警等级" />
        </a-form-item>
        <a-form-item label="说明">
          <a-textarea v-model="typeForm.Remark" :auto-size="{minRows:2,maxRows:4}" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch :checked="typeForm.Status === 1" @change="v => typeForm.Status = v ? 1 : 0" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal :title="itemForm.ID ? '编辑字典项' : '新增字典项'" :visible="itemVisible" :confirm-loading="saving" @ok="saveItem" @cancel="itemVisible=false">
      <a-form layout="vertical">
        <a-form-item label="取值" required>
          <a-input v-model="itemForm.ItemValue" :disabled="itemForm.Builtin === 1" :placeholder="currentType && currentType.DictCode === 'alarm_level' ? '整数，例如 5' : '业务里保存的值'" />
        </a-form-item>
        <a-form-item label="显示名" required>
          <a-input v-model="itemForm.ItemLabel" />
        </a-form-item>
        <a-form-item label="颜色">
          <div style="display:flex;align-items:center;gap:8px;">
            <input type="color" :value="colorPicker" @input="onColor" style="width:42px;height:32px;padding:0;border:none;background:transparent;" />
            <a-input v-model="itemForm.ItemColor" placeholder="#0099FF，可留空" style="width:160px" />
          </div>
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model="itemForm.Sort" :min="0" :max="9999" style="width:160px" />
        </a-form-item>
        <a-form-item label="说明">
          <a-input v-model="itemForm.Remark" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch :checked="itemForm.Status === 1" @change="v => itemForm.Status = v ? 1 : 0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script>
import {
  delSystemDictItem,
  delSystemDictType,
  getSystemDictItemList,
  getSystemDictTypeList,
  saveSystemDictItem,
  saveSystemDictType,
} from '@/services/dict'
import {loadSystemDict} from '@/utils/systemDict'

function blankType() {
  return {ID: 0, DictCode: '', DictName: '', Remark: '', Status: 1, Builtin: 0}
}
function blankItem(code) {
  return {ID: 0, DictCode: code || '', ItemValue: '', ItemLabel: '', ItemColor: '#0099FF', Sort: 0, Status: 1, Builtin: 0, Remark: ''}
}

export default {
  name: 'SystemDict',
  data() {
    return {
      types: [],
      items: [],
      currentType: null,
      typeLoading: false,
      itemLoading: false,
      saving: false,
      typeVisible: false,
      itemVisible: false,
      typeForm: blankType(),
      itemForm: blankItem(),
      typeColumns: [
        {title: '名称', dataIndex: 'DictName', scopedSlots: {customRender: 'DictName'}},
        {title: '编码', dataIndex: 'DictCode'},
        {title: '状态', dataIndex: 'Status', width: 80, scopedSlots: {customRender: 'Status'}},
        {title: '操作', width: 110, scopedSlots: {customRender: 'action'}},
      ],
      itemColumns: [
        {title: '排序', dataIndex: 'Sort', width: 70},
        {title: '取值', dataIndex: 'ItemValue', width: 100},
        {title: '显示', dataIndex: 'ItemLabel', scopedSlots: {customRender: 'ItemLabel'}},
        {title: '颜色', dataIndex: 'ItemColor', width: 140, scopedSlots: {customRender: 'ItemColor'}},
        {title: '状态', dataIndex: 'Status', width: 80, scopedSlots: {customRender: 'Status'}},
        {title: '操作', width: 110, scopedSlots: {customRender: 'action'}},
      ],
    }
  },
  computed: {
    colorPicker() {
      const c = (this.itemForm.ItemColor || '').trim()
      return /^#[0-9a-fA-F]{6}$/.test(c) ? c : '#0099FF'
    },
  },
  mounted() {
    this.loadTypes(true)
  },
  methods: {
    typeRow(record) {
      return {
        on: {click: () => this.selectType(record)},
      }
    },
    typeRowClass(record) {
      return this.currentType && record.ID === this.currentType.ID ? 'dict-type-active' : ''
    },
    async loadTypes(selectFirst) {
      this.typeLoading = true
      try {
        const res = await getSystemDictTypeList()
        const body = res.data || {}
        if (body.code !== 200) {
          this.$message.error(body.msg || '加载字典失败')
          return
        }
        this.types = body.list || []
        if (selectFirst && this.types.length && !this.currentType) {
          this.selectType(this.types[0])
        } else if (this.currentType) {
          const fresh = this.types.find(t => t.ID === this.currentType.ID || t.DictCode === this.currentType.DictCode)
          this.currentType = fresh || this.types[0] || null
          if (this.currentType) this.loadItems()
        }
      } finally {
        this.typeLoading = false
      }
    },
    selectType(record) {
      this.currentType = record
      this.loadItems()
    },
    async loadItems() {
      if (!this.currentType) {
        this.items = []
        return
      }
      this.itemLoading = true
      try {
        const res = await getSystemDictItemList({DictCode: this.currentType.DictCode})
        const body = res.data || {}
        this.items = body.code === 200 ? (body.list || []) : []
      } finally {
        this.itemLoading = false
      }
    },
    openType(record) {
      this.typeForm = record ? {...record} : blankType()
      this.typeVisible = true
    },
    async saveType() {
      this.saving = true
      try {
        const res = await saveSystemDictType(this.typeForm)
        const body = res.data || {}
        if (body.code !== 200) {
          this.$message.error(body.msg || '保存失败')
          return
        }
        this.typeVisible = false
        this.$message.success('已保存')
        await this.loadTypes(false)
        await loadSystemDict(true)
      } finally {
        this.saving = false
      }
    },
    async removeType(record) {
      const res = await delSystemDictType({ID: record.ID})
      const body = res.data || {}
      if (body.code !== 200) {
        this.$message.error(body.msg || '删除失败')
        return
      }
      if (this.currentType && this.currentType.ID === record.ID) {
        this.currentType = null
        this.items = []
      }
      await this.loadTypes(true)
      await loadSystemDict(true)
    },
    openItem(record) {
      if (!this.currentType) return
      this.itemForm = record ? {...record} : blankItem(this.currentType.DictCode)
      if (!record) {
        const maxSort = this.items.reduce((m, it) => Math.max(m, Number(it.Sort) || 0), -1)
        this.itemForm.Sort = maxSort + 1
      }
      this.itemVisible = true
    },
    onColor(e) {
      this.itemForm.ItemColor = e.target.value
    },
    async saveItem() {
      this.saving = true
      try {
        this.itemForm.DictCode = this.currentType.DictCode
        const res = await saveSystemDictItem(this.itemForm)
        const body = res.data || {}
        if (body.code !== 200) {
          this.$message.error(body.msg || '保存失败')
          return
        }
        this.itemVisible = false
        this.$message.success('已保存')
        await this.loadItems()
        await loadSystemDict(true)
      } finally {
        this.saving = false
      }
    },
    async removeItem(record) {
      const res = await delSystemDictItem({ID: record.ID})
      const body = res.data || {}
      if (body.code !== 200) {
        this.$message.error(body.msg || '删除失败')
        return
      }
      await this.loadItems()
      await loadSystemDict(true)
    },
  },
}
</script>

<style scoped>
/deep/ .dict-type-active td {
  background: #e6f7ff;
}
</style>

<template>

  <a-layout  style="height: 100%;">

    <a-layout-sider style="background: #fff;min-width: 300px">
      <device-tree @onSelect="onSelect" ref="deviceTree" style="min-height: 85vh"></device-tree>
    </a-layout-sider>

    <a-modal :width=modalWidth
             :confirmLoading="settingLoading"
             v-model="settingVisible"
             v-drag-modal
             :title="$t('monitor.Set')" @ok="setData"
    >
      <a-form :form="SetForm">
        <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 10px;" />

        <a-form-item :label="$t('monitor.SetValue')"
                     :labelCol="{span: 4}"
                     :wrapperCol="{span: 20}"
        >
              <a-input
                  v-decorator="[
                  'value',
                  {
                    rules: [{ required: true, message: $t('monitor.SetValue') }],
                  },
                ]"
              />
            </a-form-item>
      </a-form>
    </a-modal>

    <a-layout>
      <a-layout-content style="margin-left: 5px;margin-right: 5px">
        <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
          <a-card style="padding: 5px;min-height: 85vh" id="viewCard">
            <div v-if="deviceType==-1" class="ism-pageview-empty">
              <a-empty :description="$t('monitor.SelectDeviceTips')" />
            </div>
            <div v-else-if="deviceType==1">
<!--              <ISMRender :showUuid="showModelUuid" :showDeviceUuid="showDeviceUuid" v-if="!showRealData"/>-->

              <a-table rowKey="uuid" :pagination="pagination" :columns="columns" :data-source="tableDataSource" >
                <div
                    slot="filterDropdown"
                    slot-scope="{ setSelectedKeys, selectedKeys, confirm, clearFilters, column }"
                    style="padding: 8px"
                >
                  <a-input
                      v-ant-ref="c => (searchInput = c)"
                      :placeholder="`Search ${column.dataIndex}`"
                      :value="selectedKeys[0]"
                      style="width: 188px; margin-bottom: 8px; display: block;"
                      @change="e => setSelectedKeys(e.target.value ? [e.target.value] : [])"
                      @pressEnter="() => handleSearch(selectedKeys, confirm, column.dataIndex)"
                  />
                  <a-button
                      type="primary"
                      icon="search"
                      size="small"
                      style="width: 90px; margin-right: 8px"
                      @click="() => handleSearch(selectedKeys, confirm, column.dataIndex)"
                  >

                    {{$t('readData.Search')}}
                  </a-button>
                  <a-button size="small" style="width: 90px" @click="() => handleReset(clearFilters)">

                    {{$t('readData.Reset')}}
                  </a-button>
                </div>
                <a-icon
                    slot="filterIcon"
                    slot-scope="filtered"
                    type="search"
                    :style="{ color: filtered ? '#108ee9' : undefined }"
                />
                <template slot="customRender" slot-scope="text, record, index, column">
                    <span v-if="searchText && searchedColumn === column.dataIndex">
                      <template
                          v-for="(fragment, i) in String(text == null ? '' : text)
                          .split(new RegExp(`(${searchText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'i'))"
                      >
                        <mark
                            v-if="fragment.toLowerCase() === searchText.toLowerCase()"
                            :key="i"
                            class="highlight"
                        >{{ fragment }}</mark
                        >
                        <template v-else>{{ fragment }}</template>
                      </template>
                    </span>
                  <template v-else>
                    {{ text == null ? '' : text }}
                  </template>
                </template>
                <span slot="UpdateTime" slot-scope="UpdateTime">
                 {{UpdateTime|formatDate}}
                </span>
                <template v-for="(item, index) in columns" :slot="item.slotName">
                  <span :key="index">{{ $t(item.slotName) }}</span>
                </template>
                  <span slot="nodeType" slot-scope="nodeType">
                    <div v-if="nodeType==0" style="color: #FFCC00">
                      {{$t('monitor.ZoneType')}}
                    </div>
                   <div v-else-if="nodeType==1" style="color: #990000">
                      {{$t('monitor.DeivceType')}}
                   </div>
                </span>
                <div slot="action" slot-scope="text, record">
                  <a type="link"   @click="setting(record)" style="cursor: pointer;color: #13C2C2"><a-icon type="setting" /><span style="margin-left: 2px;">{{$t('monitor.Set')}}</span></a>
                </div>
              </a-table>
            </div>

            <div v-else-if="deviceType==0 && showModelUuid" class="ism-pageview">
              <ISMRender :showUuid="showModelUuid" showToken="" :showDeviceUuid="showDeviceUuid" />
            </div>
            <div v-else-if="deviceType==0" class="ism-pageview-empty">
              <a-empty :description="$t('monitor.SelectZoneTips')" />
            </div>
          </a-card>

        </a-spin>

      </a-layout-content>
    </a-layout>
  </a-layout>

</template>
<script>
import deviceTree from '../../components/deviceTree/DeviceTree'
import {getRealData, setData} from "../../services/device";
import {
  REAL_DATA_DEFAULT_PAGE_SIZE,
  REAL_DATA_PAGE_SIZE_OPTIONS,
  clampPageSize,
} from "@/utils/realDataBatch";
import ISMRender from '@/pages/ISMDisPlay/ISMRender';
import {formatDate} from '@/utils/common';
import {AUTH_TYPE, getAuthorization} from "@/utils/request";
import {ismDebug} from '@/utils/ismDebug';
export default {
  name: 'ISMMonitor',
  i18n: require('../../i18n/language'),
  data() {
    return {
      modalWidth:400,
      showModelUuid:"",
      showDeviceUuid:"",
      settingLoading:false,
      error: '',
      showRealData:true,
      getReadDataResponse:true,
      deviceType:-1,
      searchText: '',
      searchInput: null,
      searchedColumn: '',
      intervalId:null,
      setDataUuid:"",
      firstLoad : true,
      SetForm:this.$form.createForm(this),
      selectDeviceKey:"",
      messageShowLoad:false,
      realDataTotal:0,
      realDataPage:1,
      realDataPageSize:REAL_DATA_DEFAULT_PAGE_SIZE,
      _suppressPaginationEvent:false,
      settingVisible:false,
      columns: [
        {
          width: '20%',
          slotName: 'readData.tableName',
          scopedSlots: { filterDropdown: 'filterDropdown', filterIcon: 'filterIcon', customRender: 'customRender', title: 'readData.tableName' },
          dataIndex: 'name',
          onFilter: (value, record) =>
              record.name
                  .toString()
                  .toLowerCase()
                  .includes(value.toLowerCase()),
          onFilterDropdownVisibleChange: visible => {
            if (visible) {
              setTimeout(() => {
                this.searchInput.focus();
              }, 0);
            }
          },
        },
        {
          slotName: 'readData.tableValue',
          width: '12%',
          scopedSlots: { customRender: 'value', title: 'readData.tableValue' },
          dataIndex: 'value',
        },
        {
          slotName: 'readData.tableUnit',
          width: '6%',
          scopedSlots: { customRender: 'unit', title: 'readData.tableUnit' },
          dataIndex: 'unit',
        },
        {
          slotName: 'readData.UpdateTime',
          width: '10%',
          scopedSlots: { customRender: 'UpdateTime', title: 'readData.UpdateTime' },
          dataIndex: 'UpdateTime',
        },
        {
          slotName: 'readData.tableOpt',
          width: '10%',
          scopedSlots: { customRender: 'action',title: 'readData.tableOpt'}
        }
      ],
      pagination:{
        current:1,
        pageSize:REAL_DATA_DEFAULT_PAGE_SIZE,
        total:0,
        showSizeChanger:true,
        pageSizeOptions:REAL_DATA_PAGE_SIZE_OPTIONS,
        onChange:(page, pageSize)=>this.onRealDataPageChange(page, pageSize),
        onShowSizeChange:(current, size)=>this.onRealDataPageChange(1, size),
      },
      tableDataSource: [],
    };
  },
  components: {
    deviceTree,
    ISMRender,
  },
  mounted(){


  },
  // 多页签 keep-alive：切走再回来时恢复右侧内容，避免空白 / 一直 Loading
  activated() {
    this.messageShowLoad = false
    this.getReadDataResponse = true
    if (!this.selectNode) {
      return
    }
    const nodeValue = this.selectNode.value ? this.selectNode.value : this.selectNode
    if (nodeValue && nodeValue.type == 1 && this.selectDeviceKey) {
      this.deviceType = 1
      this.messageShowLoad = true
      this.getRealData(this.selectDeviceKey, this.realDataPage || 1, this.realDataPageSize)
      return
    }
    if (nodeValue && nodeValue.type == 0) {
      this.deviceType = 0
      this.showDeviceUuid = ''
      this.messageShowLoad = false
      const configUid = nodeValue.configUid || this.showModelUuid || ''
      // X6 画布在 keep-alive 隐藏后再显示常空白，强制重挂载 ISMRender
      if (configUid) {
        this.showModelUuid = ''
        this.$nextTick(() => {
          this.showModelUuid = configUid
        })
      } else {
        this.showModelUuid = ''
      }
    }
  },
  deactivated() {
    this.clear()
    this.messageShowLoad = false
    this.getReadDataResponse = true
  },
  watch: {
    '$route' (to) {
      // 仅离开数据仓库时清定时器；keep-alive 下切走再回来不要重置选中态
      if (!to || !to.path || to.path.indexOf('DataWarehouse') === -1) {
        this.firstLoad = true
        this.clear()
      }
    },
  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  methods: {
    handleSearch(selectedKeys, confirm, dataIndex) {
      confirm();
      this.searchText = selectedKeys[0];
      this.searchedColumn = dataIndex;
    },
    handleReset(clearFilters) {
      clearFilters();
      this.searchText = '';
    },
    setData(){
      let _t = this

      this.SetForm.validateFields((err) => {
        if (!err) {
          this.settingLoading=true
           let params = {
              deviceUuid:_t.selectDeviceKey,
              dataUuid:_t.setDataUuid,
              value:this.SetForm.getFieldValue('value'),
            };

          setData(params).then(function (res){
            _t.settingLoading=false
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t("readData.SetSuccess"))
              const newData = [..._t.tableDataSource];
              const target = newData.filter(item => _t.setDataUuid === item.mduid)[0];
              if (target) {
                target.value =_t.SetForm.getFieldValue('value')
                _t.registerGroupDataSource = newData;
              }
              _t.settingVisible = false
            }
            else
            {
              _t.$message.error(_t.$t("readData.SetFailed"))
            }
          }).catch(function (error) {
            _t.settingLoading = false
            _t.$message.error(_t.$t("readData.SetFailed"))
          }).finally(function (error) {
            _t.settingLoading = false
          })
        }
      })
    },
    setting(data){
      this.settingVisible = true
      this.settingLoading=false
      this.setDataUuid = data.mduid
    },
    syncRealDataPagination(page, pageSize, total) {
      this._suppressPaginationEvent = true
      this.pagination = Object.assign({}, this.pagination, {
        current: page,
        pageSize: pageSize,
        total: total,
      })
      this.$nextTick(() => {
        this._suppressPaginationEvent = false
      })
    },
    getRealData(uuid, page, pageSize){
      let _t = this
      const currentPage = page || _t.realDataPage || 1
      const currentSize = clampPageSize(pageSize || _t.realDataPageSize || REAL_DATA_DEFAULT_PAGE_SIZE)
      const nodeVal = (_t.selectNode && (_t.selectNode.value || _t.selectNode)) || {}
      const params = {
        uuid:uuid,
        page: currentPage,
        pageSize: currentSize,
        IsRemoveGW: !!nodeVal.IsRemoteGw,
        ProjectUuid: nodeVal.ProjectUUID || nodeVal.project_uuid || '',
      }
      ismDebug('DW.getRealData.enter', {
        uuid,
        page: currentPage,
        pageSize: currentSize,
        busy: this.getReadDataResponse !== true,
        deviceType: this.deviceType,
        messageShowLoad: this.messageShowLoad,
      })
      // 上一次请求未完成时直接 return，必须清掉 loading，否则会一直转圈
      if(this.getReadDataResponse!=true)
      {
        ismDebug('DW.getRealData.skipBusy', {uuid, page: currentPage})
        this.messageShowLoad = false
        return
      }
      if (currentPage === 1) {
        this.tableDataSource=[]
      }
      this.getReadDataResponse = false
      const reqSeq = (_t._realDataReqSeq = (_t._realDataReqSeq || 0) + 1)
      const finishLoad = function (tag, extra) {
        if (reqSeq !== _t._realDataReqSeq) {
          return
        }
        _t.messageShowLoad = false
        _t.getReadDataResponse = true
        if (_t._realDataLoadTimer) {
          clearTimeout(_t._realDataLoadTimer)
          _t._realDataLoadTimer = null
        }
        ismDebug(tag, Object.assign({
          deviceType: _t.deviceType,
          tableLen: (_t.tableDataSource || []).length,
          messageShowLoad: _t.messageShowLoad,
          reqSeq,
        }, extra || {}))
      }
      // 兜底：后端 gzip 双重压缩时 axios 可能永不 settle，强制关 loading
      if (_t._realDataLoadTimer) {
        clearTimeout(_t._realDataLoadTimer)
      }
      _t._realDataLoadTimer = setTimeout(function () {
        if (reqSeq === _t._realDataReqSeq && _t.messageShowLoad) {
          ismDebug('DW.getRealData.loadTimeout', {reqSeq, uuid})
          finishLoad('DW.getRealData.finally', {forcedTimeout: true})
        }
      }, 15000)
      getRealData(params).then(function (res){
        if (reqSeq !== _t._realDataReqSeq) {
          return
        }
        // 兼容 data 为字符串 / 被二次包装的情况
        let body = res && res.data
        if (typeof body === 'string') {
          try { body = JSON.parse(body) } catch (e) {
            ismDebug('DW.getRealData.parseError', {message: e && e.message, head: String(body).slice(0, 120)})
            finishLoad('DW.getRealData.finally', {parseError: true})
            _t.$message.error(_t.$t('readData.LoadFailed'))
            return
          }
        }
        const code = body ? body.code : null
        const rows = (body && body.realData) ? body.realData : []
        ismDebug('DW.getRealData.response', {
          code,
          rows: rows.length,
          total: body ? body.total : null,
          page: body ? body.page : null,
          pageSize: body ? body.pageSize : null,
          dataType: typeof (res && res.data),
        })
        if(body && body.code==0)
        {
            try {
              _t.realDataTotal = body.total != null ? body.total : rows.length
              _t.realDataPage = body.page || currentPage
              _t.realDataPageSize = body.pageSize || currentSize
              _t.syncRealDataPagination(_t.realDataPage, _t.realDataPageSize, _t.realDataTotal)
              const nextRows = []
              for(let i=0;i<rows.length;i++)
              {
                  nextRows.push({
                    key: rows[i].ID || rows[i].id || rows[i].uuid,
                    no: rows[i].ID || rows[i].id,
                    name: rows[i].name,
                    value: rows[i].value,
                    uuid: rows[i].uuid,
                    unit: rows[i].unit,
                    mduid: rows[i].mduid,
                    UpdateTime: rows[i].UpdatedAt || rows[i].updated_at,
                  })
              }
              _t.tableDataSource = nextRows
              _t.deviceType = 1
              ismDebug('DW.getRealData.rendered', {
                tableLen: nextRows.length,
                deviceType: _t.deviceType,
                messageShowLoad: false,
              })
            } catch (e) {
              ismDebug('DW.getRealData.mapError', {message: e && e.message, stack: e && e.stack})
              _t.$message.error(_t.$t('readData.LoadFailed'))
            }
        } else {
          ismDebug('DW.getRealData.badCode', {code, dataKeys: body ? Object.keys(body) : []})
        }
        finishLoad('DW.getRealData.finally')
      }).catch(function (err) {
        if (reqSeq !== _t._realDataReqSeq) {
          return
        }
        ismDebug('DW.getRealData.catch', {message: err && err.message})
        _t.$message.error(_t.$t('readData.LoadFailed'))
        finishLoad('DW.getRealData.finally', {fromCatch: true})
      })
    },
    onRealDataPageChange(page, pageSize) {
      if (this._suppressPaginationEvent) {
        return
      }
      if (!this.selectDeviceKey) {
        return
      }
      const nextSize = clampPageSize(pageSize)
      if (page === this.realDataPage && nextSize === this.realDataPageSize) {
        return
      }
      this.realDataPage = page
      this.realDataPageSize = nextSize
      this.messageShowLoad = true
      // 翻页只保留当前页数据，避免浏览器缓存堆积上万行
      this.tableDataSource = []
      this.getRealData(this.selectDeviceKey, page, nextSize)
    },
    onSelect(selectData) {
      const info = selectData.info
      const nodeRef = (info && info.dataRef) ? info.dataRef : info
      let _t = this
      this.clear()
      this.selectNode = nodeRef
      const nodeValue = nodeRef && nodeRef.value ? nodeRef.value : nodeRef
      this.showModelUuid = nodeValue && nodeValue.configUid
      this.deviceType=3
      this.messageShowLoad = true
      ismDebug('DW.onSelect', {
        key: selectData && selectData.key,
        type: nodeValue && nodeValue.type,
        configUid: nodeValue && nodeValue.configUid,
        uuid: nodeValue && nodeValue.uuid,
        hasDataRef: !!(info && info.dataRef),
      })
      if(nodeValue && nodeValue.type==1)
      {
        this.showDeviceUuid = nodeValue.uuid || nodeValue.configUid
        _t.$nextTick(() => {
          _t.deviceType=1
          _t.messageShowLoad = true
          _t.realDataPage = 1
          _t.selectDeviceKey = selectData.key
          ismDebug('DW.onSelect.deviceTick', {key: selectData.key, deviceType: _t.deviceType})
          _t.getRealData(selectData.key, 1, _t.realDataPageSize)
          _t.$EventBus.$off("readDataPush")
          _t.$EventBus.$off("StaticData")
          _t.$EventBus.$off("SystemData")
          _t.$EventBus.$on("readDataPush", (data) => {
            let realData = data
            if(realData.DeviceUuid==_t.selectDeviceKey) {
              if(realData.Data!=null) {
                for (let j = 0; j < _t.tableDataSource.length; j++) {
                  for (let k = 0; k < realData.Data.length; k++) {
                    if (_t.tableDataSource[j].uuid == realData.Data[k].Uuid) {
                      _t.tableDataSource[j].value = realData.Data[k].Value
                      _t.tableDataSource[j].UpdateTime = new Date()
						break
                    }
                  }
                }
              }
            }
          });

          _t.$EventBus.$on("StaticData", (data) => {
            let realData = data
            if(realData.DeviceUuid==_t.selectDeviceKey) {
              if(realData.Data!=null) {
                for (let j = 0; j < _t.tableDataSource.length; j++) {
                  for (let k = 0; k < realData.Data.length; k++) {
                    if (_t.tableDataSource[j].uuid == realData.Data[k].Uuid) {
                      _t.tableDataSource[j].value = realData.Data[k].Value
                      _t.tableDataSource[j].UpdateTime = new Date()
                    }
                  }
                }
              }
            }
          });

          _t.$EventBus.$on("SystemData", (data) => {

          });
          _t.$EventBus.$on("RealAlarm", (data) => {
            let realAlarmData = data
            if(realAlarmData.DeviceUuid==_t.selectDeviceKey) {
                for (let j = 0; j < _t.tableDataSource.length; j++) {
                    if (_t.tableDataSource[j].uuid == realAlarmData.DataUuid) {
                      _t.tableDataSource[j].value = realAlarmData.Value
                      _t.tableDataSource[j].UpdateTime = new Date()
                      break;
                    }
                }
            }
          });
        });
      }
      else
      {
        _t.$nextTick(() => {
          _t.firstLoad=true
          _t.showDeviceUuid=""
          _t.messageShowLoad = false
          _t.selectDeviceKey=''
          _t.deviceType=0
          // 区域未绑定组态页时给出明确提示，避免误以为页面坏了
          if (!(_t.showModelUuid)) {
            ismDebug('DW.onSelect.zoneNoConfig', {
              key: selectData && selectData.key,
              uuid: nodeValue && nodeValue.uuid,
            })
          }
          _t.clear()
        });
      }
    },
    startTimer(key) {
      // 计时器正在进行中，退出函数
      if (this.intervalId != null) {
        return;
      }
      this.intervalId = setInterval(()=>{
        this.getRealData(key, this.realDataPage, this.realDataPageSize)
      },1000)
    },
    clear() {
      clearInterval(this.intervalId); //清除计时器
      this.intervalId = null; //设置为null
    }
  },
};
</script>

<style lang="less">
.resize-table-th {
  position: relative;
.table-draggable-handle {
  height: 100% !important;
  bottom: 0;
  left: auto !important;
  right: -5px;
  cursor: col-resize;
  touch-action: none;
  position: absolute;
}
}
.ism-pageview {
  z-index: 999;
  height: 100%;
  width: 100%;
  overflow:scroll;
}
#components-layout-demo-side .logo {
  height: 32px;
  background: rgba(255, 255, 255, 0.2);
  margin: 16px;
}

.search{
  margin-bottom: 54px;
}
.fold{
  width: calc(100% - 216px);
  display: inline-block
}
.operator{
  margin-bottom: 18px;
}
.spin-content {
  border: 1px solid #91d5ff;
  background-color: #e6f7ff;
  padding: 30px;
}
#viewCard >.ant-card-body {
   padding: 0px;
  zoom: 1;
}
@media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}

.ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
.ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

.ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
</style>

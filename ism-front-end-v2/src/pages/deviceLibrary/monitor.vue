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
            <div v-if="deviceType==1">
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
                          v-for="(fragment, i) in text
                          .toString()
                          .split(new RegExp(`(?<=${searchText})|(?=${searchText})`, 'i'))"
                      >
                        <mark
                            v-if="fragment.toLowerCase() === searchText.toLowerCase()"
                            :key="i"
                            class="highlight"
                        >{{ $t(fragment) }}</mark
                        >
                        <template v-else>{{ $t(fragment) }}</template>
                      </template>
                    </span>
                  <template v-else>
                    {{ $t(text) }}
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

            <div v-else-if="deviceType==0" class="ism-pageview">
              <ISMRender :showUuid="showModelUuid" showToken="" :showDeviceUuid="showDeviceUuid" />
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
import ISMRender from '@/pages/ISMDisPlay/ISMRender';
import {formatDate} from '@/utils/common';
import {AUTH_TYPE, getAuthorization} from "@/utils/request";
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
      deviceType:0,
      searchText: '',
      searchInput: null,
      searchedColumn: '',
      intervalId:null,
      setDataUuid:"",
      firstLoad : true,
      SetForm:this.$form.createForm(this),
      selectDeviceKey:"",
      messageShowLoad:false,
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
        pageSize:15,
        showSizeChanger:true
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
  watch: {
    '$route' () {
      this.firstLoad = true
      this.clear();
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
    getRealData(uuid){
      let _t = this
      const params = {
        uuid:uuid,
        IsRemoveGW:this.selectNode.value.IsRemoteGw,
        ProjectUuid:this.selectNode.value.ProjectUUID,
      }
      if(this.getReadDataResponse!=true)
      {
        return
      }
      this.tableDataSource=[]
      this.getReadDataResponse = false
      getRealData(params).then(function (res){
        _t.getReadDataResponse = true
        _t.messageShowLoad = false
        if(res.data.code==0)
        {
            for(let i=0;i<res.data.realData.length;i++)
            {
                let temp = {}
                temp.key = res.data.realData[i].ID
                temp.no = res.data.realData[i].ID
                temp.name = res.data.realData[i].name
                temp.value = res.data.realData[i].value
                temp.uuid = res.data.realData[i].uuid
                temp.unit = res.data.realData[i].unit
                temp.mduid = res.data.realData[i].mduid
                temp.UpdateTime = res.data.realData[i].UpdatedAt
                _t.tableDataSource.push(temp)
                temp = {}
            }
        }
      })
    },
    onSelect(selectData) {
      const info = selectData.info
      let _t = this
      this.clear()
      this.selectNode = info
      this.showModelUuid =this.selectNode.value.configUid
      this.deviceType=3
      this.messageShowLoad = true
      if(this.selectNode.value.type==1)
      {
        this.showDeviceUuid =this.selectNode.value.configUid
        this.showDeviceUuid = this.selectNode.value.uuid
        this.$nextTick(function(){
          this.deviceType=1
          this.messageShowLoad = true
          this.selectDeviceKey = selectData.key
          this.getRealData(selectData.key)
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
        this.$nextTick(function(){
          this.firstLoad=true
          this.showDeviceUuid=""
          this.messageShowLoad = false
          this.selectDeviceKey=''
          this.deviceType=0
          this.clear()
        });
      }
    },
    startTimer(key) {
      // 计时器正在进行中，退出函数
      if (this.intervalId != null) {
        return;
      }
      this.intervalId = setInterval(()=>{
        this.getRealData(key)
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

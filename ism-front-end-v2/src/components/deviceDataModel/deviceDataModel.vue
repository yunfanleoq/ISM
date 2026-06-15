<template>
  <div ref="dataModel" >
    <a-modal
             :visible="dataModelDialog"
             @cancel="dataModelDialog=false"
             width="800px"
             :bodyStyle="{'min-height': '500px'}"
             :footer="null"
             v-drag-modal
            :title="$t('component.deviceDataModel.title')">

        <div style="margin: 6px;">
            <a-card size="small" :title="$t('component.deviceDataModel.DataFrom')">
              <a-form layout="inline" :label-col="{ span: 8}" :wrapper-col="{ span: 16 }">
                <a-form-item :label="$t('component.deviceDataModel.DataFrom')" >
                  <a-select @change="changeDataFromType" style="width: 200px" :dropdownStyle="{'z-index': 9999999}">
                    <a-select-option  v-for="(device,index) in FromList" :key="index" :value=device.value>
                      {{ $t(device.title) }}
                    </a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.ComponentBandDevice')" v-if="fromType==2" >
                  <a-tree-select
                      show-search
                      tree-node-filter-prop="title"
                      @select="SelectDevice"
                      style="width: 200px"
                      :dropdown-style="{ 'z-index': 9999999,maxHeight: '400px', overflow: 'auto' }"
                      :tree-data="deviceTreeData"
                      :replace-fields="{ value: 'key',title:'text'}"
                      placeholder="Please select"
                      tree-default-expand-all
                  >
                  </a-tree-select>
                </a-form-item>
                <a-form-item v-if="fromType==1" :label="$t('component.deviceDataModel.selectDeviceType')">
                  <a-select @change="changeDeviceType" style="width:230px" :dropdownStyle="{'z-index': 9999999}">
                    <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type>
                      {{ $t(device.name) }}
                    </a-select-option>
                  </a-select>
                </a-form-item>
              </a-form>
            </a-card>
        </div>

      <a-spin :spinning="messageShowLoad" style="min-height: 300px">
          <div style="margin: 6px;z-index: 9001">
            <a-layout >
              <a-layout-sider style="background: #fff;"  v-if="(fromType==1)&&(DeviceType!=6&&DeviceType!=7)"  :bodyStyle="{padding:'2px'}" :panelStyle="{width:'200px',minWidth:'100px'}">
                <a-card size="small" :hoverable="true" :title="$t('component.deviceDataModel.DataModelList')" style="max-height: 400px">
                  <a-directory-tree :treeData="modelList" @select="changeDeviceModel" style="max-height: 340px; overflow-y: auto;" :replace-fields="{ key: 'uuid',title:'text'}">

                  </a-directory-tree>
                </a-card>
              </a-layout-sider>
              <a-layout-content style="background: #fff;"  :bodyStyle="{padding:'2px','min-height':'300px'}" :panelStyle="{width:((fromType==1)&&(DeviceType!=6&&DeviceType!=7))?'500px':'700px',height:'100%',minWidth:'100px'}">

                <a-table
                    row-key="uuid"
                    :customRow="rowClick"
                    :columns="columns"
                    size="small"
                    :rowClassName="setRowClassName"
                    :data-source="dataSource"
                    :pagination="pagination"
                    :scroll="{ y: 340,x:0 }"
                >
                  <template slot="name" slot-scope="name">
                    {{ $t(name)}}
                  </template>

                  <div
                      slot="filterDropdown"
                      slot-scope="{ setSelectedKeys, selectedKeys, confirm, clearFilters, column }"
                      style="padding: 8px;"
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
                  <template v-for="(item, index) in columns" :slot="item.slotName">
                    <span :key="index">{{ $t(item.slotName) }}</span>
                  </template>
                </a-table>
              </a-layout-content>
            </a-layout>
          </div>
        </a-spin>
      <a-divider />
      <div class="dialog-button">

        <a-button key="submit" type="primary" @click="onSelectData">
          {{ $t('component.deviceDataModel.submit')}}
        </a-button>
        <a-button style="margin-left: 10px" key="back" @click="handleCancel">
          {{$t('component.deviceDataModel.cancel')}}
        </a-button>
      </div>
    </a-modal>
  </div>
</template>
<script>
import {vDragModal} from "@/utils/vmodalDrage"
import {getRealData, getSupportDeviceList} from "../../services/device";
import {getDatasByUuid, snmpModelList} from "@/services/snmpmodel";
import {getMonitorTree} from "@/services/device";
import {GetDataModelData, GetSystemData} from "@/services/system";
export default {
  name: 'deviceDataModel',
  i18n: require('../../i18n/language'),
  data() {
    return {
      dataModelDialog:false,
      visible: false,
      pageSize: 10,
      DeviceType:-1,
      searchText:"",
      searchedColumn: '',
      selectedRowsIndex:-1,
      selectNodeInfo:{},
      rowSelection:{
        onSelect:this.onDataTableSelect,
        onSelectAll:this.onDataTableSelectAll
      },
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      selectDataModelUUid:"",
      pagePosition: "bottom",
      pageOptions: [
        { value: "bottom", text: "Bottom" },
        { value: "top", text: "Top" },
        { value: "both", text: "Both" }
      ],
      columns :[
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
      ],
      deviceTreeData:[],
      fromType:2,
      FromList:[
        {
          value:2,
          title:"component.deviceDataModel.deviceData"
        },
        {
          value:1,
          title:"component.deviceDataModel.deviceModel"
        },
        {
          value:3,
          title:"component.deviceDataModel.SystemData"
        }
      ],
      deviceName:"",
      modelList:[],
      dataSource:[],
      IsDevice:false,
      DeviceSN:"",
      messageShowLoad:false,
      supportDeviceList:[],
      selectionDataInfo: {}
    };
  },
  watch: {
    '$route'() {

    }
  },
  mounted(){
    if(this.fromType==2)
    {
      this.getMonitorTree()
    }
    this.getSupportDevice()
  },
  methods: {
    setRowClassName(record,index){
      return this.selectedRowsIndex==index ? 'clickRowStyle' : '';
    },
    rowClick(record, index) {
      let _t = this
      return {
        on: {
          click: () => {
            _t.selectionDataInfo = record
            _t.selectedRowsIndex = index
          },
          dblclick: () => {
            _t.selectionDataInfo = record
            _t.selectedRowsIndex = index
          },
        }
      }
    },
    handleSearch(selectedKeys, confirm, dataIndex) {
      confirm();
      this.searchText = selectedKeys[0];
      this.searchedColumn = dataIndex;
    },
    handleReset(clearFilters) {
      clearFilters();
      this.searchText = '';
    },
    onDataTableSelect (record, selected, selectedRows)  {
      if(selected)
      {
        this.selectDataTableUuid.push(record.uuid)
      }
      else
      {
        const index=this.findDataTableSelectIndex(record.uuid)
        if(index!=-1)
        {
          this.selectDataTableUuid.splice(index,1)
        }
      }
    },
    onDataTableSelectAll(selected, selectedRows, changeRows) {
      if(selected)
      {
        for(let i=0;i<selectedRows.length;i++)
        {
          const index=this.findDataTableSelectIndex(selectedRows[i].uuid)
          if(index==-1)
          {
            this.selectDataTableUuid.push(selectedRows[i].uuid)
          }
        }

      }
      else
      {
        for(let i=0;i<changeRows.length;i++)
        {
          const index=this.findDataTableSelectIndex(changeRows[i].uuid)
          if(index!=-1)
          {
            this.selectDataTableUuid.splice(index,1)
          }
        }
      }
    },
    getRealData(uuid){
      let _t = this
      const params = {
        uuid:uuid,
        IsRemoveGW:this.selectNodeInfo.IsRemoteGw,
        ProjectUuid:this.selectNodeInfo.ProjectUUID,
      }
      this.dataSource = []
      this.messageShowLoad = true
      this.getReadDataResponse = false
      getRealData(params).then(function (res){
        _t.getReadDataResponse = true
        _t.messageShowLoad = false
        if(res.data.code==0)
        {
          for(let i=0;i<res.data.realData.length;i++)
          {
              let temp={}
              temp.id = i.toString()
              temp.name = res.data.realData[i].name
              temp.uuid = res.data.realData[i].mduid
              temp.unit = res.data.realData[i].unit
              _t.dataSource.push(temp)
              temp={}
          }
        }
      })
    },
    GetSystemData(uuid){
      let _t = this
      const params = {
        uuid:uuid
      }
      this.dataSource = []
      this.messageShowLoad = true
      this.getReadDataResponse = false
      GetSystemData(params).then(function (res){
        _t.getReadDataResponse = true
        _t.messageShowLoad = false
        if(res.data.code==0)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            let temp={}
            temp.id = i.toString()
            temp.name = res.data.list[i].name
            temp.uuid = res.data.list[i].uuid
            temp.unit = res.data.list[i].unit
            _t.dataSource.push(temp)
            temp={}
          }
        }
      })
    },
    GetDataModelData(type){
      let _t = this
      const params = {
        type:type
      }
      this.dataSource = []
      this.messageShowLoad = true
      this.getReadDataResponse = false
      GetDataModelData(params).then(function (res){
        _t.getReadDataResponse = true
        _t.messageShowLoad = false
        if(res.data.code==0)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            let temp={}
            temp.id = i.toString()
            temp.name = res.data.list[i].name
            temp.uuid = res.data.list[i].uuid
            temp.unit = res.data.list[i].unit
            _t.dataSource.push(temp)
            temp={}
          }
        }
      })
    },
    getMonitorTree(){
      let _t = this
      this.deviceTreeData=[]
      getMonitorTree().then(function (res){
        if(res.data.code==0)
        {
          _t.deviceTreeData =res.data.list==null?[]:res.data.list
        }
      })
    },
    getDeviceDataModel(uuid){
      let params={
        muid:typeof uuid=='undefined'?"":uuid,
        type:this.DeviceType
      }
      this.dataSource = []
      let _t = this
      this.messageShowLoad = true
      getDatasByUuid(params).then(function (res){
        _t.messageShowLoad = false
        if((res.data.mibs!=null)&&(res.data.mibs.length>0))
        {
          let temp={}
          for(let i=0;i<res.data.mibs.length;i++)
          {
            temp.id = i.toString()
            temp.name = res.data.mibs[i].name
            temp.uuid = res.data.mibs[i].uuid
            temp.unit = res.data.mibs[i].unit
            _t.dataSource.push(temp)
            temp={}
          }
        }
      })
    },
    changeDataFromType(value){
      this.dataSource=[]
      this.fromType = value
      if (value==2)
      {
        this.IsDevice = true
        this.getMonitorTree()
      }
      else if (value==1)
      {
        this.IsDevice = false
      }
      else if (value==3)
      {
        this.GetSystemData()
      }
    },
    changeDeviceType(value){
      this.dataSource=[]
      this.IsDevice=false
      this.DeviceSN=""
      this.DeviceType = value
      if(value!=6&&value!=7)
      {
        this.getModelList()
      }
      else
      {
        this.GetDataModelData(value)
      }
    },
    getModelList(){
      let _t = this
      this.modelList=[]
      const  params= {
        type:this.DeviceType
      }
      snmpModelList(params).then(function (res){
        let tempData={}
        if(res.data.list!=null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            tempData.id = i
            tempData.uuid = res.data.list[i].uuid
            tempData.text = res.data.list[i].name
            _t.modelList.push(tempData)
            tempData={}
          }
        }
      })
    },
    getSupportDevice(){
      let _t = this
      getSupportDeviceList().then(function (res){
        _t.supportDeviceList =res.data.list==null?[]:res.data.list
        let systemData={
          type:6,
          name:_t.$t('component.deviceDataModel.SystemVar'),
        }
        _t.supportDeviceList.push(systemData)
        let diyData={
          type:7,
          name:_t.$t('device.CustomDevice'),
        }
        _t.supportDeviceList.push(diyData)
      })
    },
    selectData(ev){

    },
    SelectDevice(ev,node){
      this.deviceName = node.$options.propsData.title
      this.IsDevice = true
      this.dataSource=[]
      this.DeviceSN = ev
      this.selectNodeInfo = node.$options.propsData.dataRef.value
      this.getRealData(ev)
    },
    changeDeviceModel(ev,event){
      this.dataSource=[]
      this.selectDataModelUUid = ev[0]
      this.getDeviceDataModel(ev[0])
    },
    handleCancel(){
      this.dataModelDialog = false
    },
    showDataModal() {
      this.selectedRowsIndex = -1
      this.dataModelDialog = true
      this.visible = true
    },
    onSelectData(){
      if (this.selectedRowsIndex==-1)
      {
        this.$message.error(this.$t('component.deviceDataModel.SelectError'), 3)
        return
      }
      if(this.IsDevice)
      {
        this.selectionDataInfo.DeviceName = this.deviceName
        this.selectionDataInfo.IsDevice = true
        this.selectionDataInfo.DeviceSN = this.DeviceSN
      }
      else
      {
        this.selectionDataInfo.DeviceType = this.DeviceType
        this.selectionDataInfo.IsDevice = false
        this.selectionDataInfo.DeviceSN = ""
      }
      if(this.fromType!=1)
      {
        this.selectionDataInfo.DeviceType=-1
        this.selectionDataInfo.selectDataModelUUid=""
      }
      else
      {
        this.selectionDataInfo.selectDataModelUUid = this.selectDataModelUUid
        this.selectionDataInfo.DeviceType = this.DeviceType
      }
      this.selectionDataInfo.name = this.$t(this.selectionDataInfo.name)
      this.$emit("onSelectDataModel", this.selectionDataInfo);
      this.dataModelDialog = false
    }
  },
}
</script>

<style scoped lang="less">

::v-deep .ant-dropdown  .ant-table-filter-dropdown{
  z-index: 9001;
}
::v-deep .ant-divider-horizontal {
  margin: 5px 0;
}
::v-deep .dialog-button{
  text-align: right;
}
::v-deep  .ant-card-body {
  padding: 1px;
}
::v-deep  .ant-card-head {
  padding: 0 10px;
}
::v-deep .clickRowStyle{
  background-color: #00ccff;
}
::v-deep .ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
::v-deep .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}
// 去除鼠标经过默认的背景颜色
::v-deep .ant-table-tbody>tr.ant-table-row:hover>td {
  background: none !important;
}
::v-deep  .ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}

::v-deep .ant-modal-body {
  padding: 5px;
  }
::v-deep .ant-select-dropdown {
  z-index: 900000;
}
  #DataGrid .datagrid-field-td{
    display: none;
  }
  #DataGrid .datagrid-header {
     border-style: unset;
  }
</style>

<template>
  <a-card>
    <div >
      <a-form layout="horizontal">
        <div :class="advanced ? null: 'fold'">
          <a-row >
            <a-col :md="8" :sm="24" >
              <a-form-item
                  :label="$t('reporting.AlarmHistory.DeviceList')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-tree-select
                    show-search
                    v-model="SelectDevice"
                    @change="SelectTreeDevice"
                    tree-node-filter-prop="title"
                    style="width: 100%"
                    tree-checkable
                    allow-clear
                    :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
                    :tree-data="deviceTreeData"
                    :replace-fields="{ value: 'key',title:'text'}"
                    :placeholder="$t('reporting.AlarmHistory.DeviceList')"
                    tree-default-expand-all
                >
                </a-tree-select>
              </a-form-item>
            </a-col>
            <a-col :md="8" :sm="24" >
              <a-form-item
                  :label="$t('reporting.AlarmHistory.DataList')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-select  @dropdownVisibleChange="GetDeviceModelDataList"
                           @popupScroll="handlePopupScroll"
                           @search="handleSearch" allowClear show-search optionFilterProp="children" mode="multiple"  style="width: 100%" :token-separators="[',']" v-model="SelectAlarmData">

                  <a-select-option v-for="(alarmItem,itemIndex) in frontDataZ" :key="itemIndex" :value=alarmItem.uuid>
                    {{ $t(alarmItem.name) }}
                  </a-select-option>
                </a-select>

              </a-form-item>
            </a-col>
            <a-col :md="2" :sm="24" >
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryAlarmList">{{$t('reporting.AlarmHistory.Query')}}</a-button>
              </span>
            </a-col>
          </a-row>
          <a-row >

          </a-row>
        </div>
      </a-form>
    </div>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table :pagination="pagination" :columns="columns" :data-source="dataSource" rowKey='DeviceName'>
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <span slot="UpdatedAt" slot-scope="UpdatedAt">
         {{UpdatedAt|formatDate}}
        </span>
        <span slot="name" slot-scope="name">
         {{$t(name)}}
        </span>
        <span slot="AlarmLevel" slot-scope="AlarmLevel">
          <a-tag style="background-color:#0099FF ;" v-if="AlarmLevel==0">
           {{ $t('dataModel.alarm.Tips')}}
         </a-tag>
          <a-tag style="background-color:#0099FF ;" v-if="AlarmLevel==1">
           {{ $t('dataModel.alarm.Minor')}}
         </a-tag>
          <a-tag style="background-color:yellow ;" v-if="AlarmLevel==2">
           {{ $t('dataModel.alarm.Importance')}}
         </a-tag>
          <a-tag style="background-color:orange ;" v-if="AlarmLevel==3">
           {{ $t('dataModel.alarm.Urgency')}}
         </a-tag>
          <a-tag style="background-color:red ;" v-else-if="AlarmLevel==4">
           {{ $t('dataModel.alarm.Deadly')}}
         </a-tag>
        </span>

        <div slot="action" slot-scope="text, record">
          <a-popconfirm :title="$t('alarm.current.RestoreTips')" @confirm="ShieldAlarm(record,0)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a style=""><a-icon type="alert" />{{$t('alarm.current.Restore')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>

  </a-card>
</template>

<script>
import {getMonitorTree} from "@/services/device";
import {GetDeviceModelDataList} from "@/services/device";
import {GetCurrentAlarmList,UpdateCurrentAlarm} from "@/services/alarm";
import moment from 'moment';

import {formatDate} from '@/utils/common';
import {GetCurrentShieldAlarmList} from "../../../services/alarm";

export default {
  name: 'TaskAndPlan',
  i18n: require('../../../i18n/language'),
  components:{

  },
  data () {
    return {
      dataZ:[],
      valueData: '',
      treePageSize: 100,
      scrollPage: 1,
      frontDataZ:[],
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      SelectDateType: 'Day',
      SelectDevice:[],
      SelectDateRange:"",
      SelectAlarmData:[],
      deviceTreeData:[],
      AlarmDataTree:[],
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      columns: [
        {
          slotName: 'reporting.AlarmHistory.DeviceName',
          width: '10%',
          scopedSlots: { customRender: 'DeviceName', title: 'reporting.AlarmHistory.DeviceName' },
          dataIndex: 'DeviceName',
        },
        {
          width: '10%',
          slotName: 'reporting.AlarmHistory.AlarmName',
          scopedSlots: { customRender: 'name', title: 'reporting.AlarmHistory.AlarmName' },
          dataIndex: 'name',
        },
        {
          width: '10%',
          slotName: 'alarm.current.ShieldTime',
          scopedSlots: { customRender: 'UpdatedAt', title: 'alarm.current.ShieldTime' },
          dataIndex: 'UpdatedAt',
        },
        {
          width: '10%',
          slotName: 'reporting.AlarmHistory.AlarmLevel',
          scopedSlots: { customRender: 'AlarmLevel', title: 'reporting.AlarmHistory.AlarmLevel' },
          dataIndex: 'alarmLevel',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          width: '10%',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      dataSource: [],
      conditionExpress:"",
      selectedRows: []
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    this.QueryAlarmList()
  },
  activated(){

  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  created(){
    this.getMonitorTree()
    this.GetDeviceModelDataList()
  },
  watch: {
    '$route' () {
      this.getMonitorTree()
    }
  },
  methods: {
    SelectTreeDevice(value,node,extera){
      this.GetDeviceModelDataList()
    },
    ShieldAlarm(item,Shield){
      let _t = this
      _t.dataSource = []
      const params = {
        type:2,
        update:{
          duid:item.duid,
          uuid:item.uuid,
          AlarmShield:Shield
        }
      }
      this.messageShowLoad=true
      UpdateCurrentAlarm(params).then(function (res){
        if(res.data.code==0)
        {
          _t.dataSource =res.data.list
        }
        _t.messageShowLoad=false
      }).catch(function(){
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    filterOption(input, option) {
      return (
          option.componentOptions.children[0].text.toLowerCase().indexOf(input.toLowerCase()) >= 0
      );
    },
    handleSearch (val) {
      this.valueData = val
      if (!val) {
        this.GetDeviceModelDataList()
      } else {
        this.frontDataZ = []
        this.scrollPage = 1
        this.dataZ.forEach(item => {
          if (item.name.indexOf(val) >= 0) {
            this.frontDataZ.push(item)
          }
        })
        this.allDataZ = this.frontDataZ
        this.frontDataZ = this.frontDataZ.slice(0, 100)
      }
    },
//下拉框下滑事件
    handlePopupScroll (e) {
      const { target } = e
      const scrollHeight = target.scrollHeight - target.scrollTop
      const clientHeight = target.clientHeight
      // 下拉框不下拉的时候
      if (scrollHeight === 0 && clientHeight === 0) {
        this.scrollPage = 1
        console.log(this.scrollPage)
      } else {
        // 当下拉框滚动条到达底部的时候
        if (scrollHeight < clientHeight + 5) {
          this.scrollPage = this.scrollPage + 1
          const scrollPage = this.scrollPage// 获取当前页
          const treePageSize = this.treePageSize * (scrollPage || 1)// 新增数据量
          const newData = [] // 存储新增数据
          let max = '' // max 为能展示的数据的最大条数
          if (this.dataZ.length > treePageSize) {
            // 如果总数据的条数大于需要展示的数据
            max = treePageSize
          } else {
            // 否则
            max = this.dataZ.length
          }
          // 判断是否有搜索
          if (this.valueData) {
            this.allDataZ.forEach((item, index) => {
              if (index < max) { // 当data数组的下标小于max时
                newData.push(item)
              }
            })
          } else {
            this.dataZ.forEach((item, index) => {
              if (index < max) { // 当data数组的下标小于max时
                newData.push(item)
              }
            })
          }

          this.frontDataZ = newData // 将新增的数据赋值到要显示的数组中
        }
      }
    },
    GetDeviceModelDataList(){
      let _t = this
      this.AlarmDataTree=[]
      _t.dataZ=[]
      _t.frontDataZ=[]
      const params ={
        SelectDevice:this.SelectDevice,
        getType:1
      }
      GetDeviceModelDataList(params).then(function (res){
        if(res.data.code==0)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            if(typeof res.data.list[i].DataList!="undefined"&&res.data.list[i].DataList!=null) {
              for (let j = 0; j < res.data.list[i].DataList.length; j++) {
                _t.dataZ.push(res.data.list[i].DataList[j])
                _t.AlarmDataTree.push(res.data.list[i].DataList[j])
              }
            }
          }
          _t.frontDataZ = _t.dataZ.slice(0, 100)
        }
      })
    },
    getMonitorTree(){
      let _t = this
      this.deviceTreeData=[]
      getMonitorTree().then(function (res){
        if(res.data.code==0)
        {
          _t.deviceTreeData =res.data.list
        }
      })
    },
    QueryAlarmList(){
      let _t = this


      _t.dataSource = []
      const params = {
        deviceList:this.SelectDevice,
        dataList:this.SelectAlarmData,
      }
      this.messageShowLoad=true
      GetCurrentShieldAlarmList(params).then(function (res){
        if(res.data.code==0)
        {
          _t.dataSource =res.data.list
        }
        _t.messageShowLoad=false
      }).catch(function(){
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    }
  }
}
</script>

<style lang="less">

.plus-icon-enter-active{
  transition: opacity .5s;
}
.plus-icon-enter{
  opacity: 0;
}
.plus-icon-leave-active{
  transition: opacity .5s;
}
.plus-icon-leave-to{
  opacity: 0;
}
.plus-icon-enter-to{
  opacity: 1;
}

.code-box-actions {
  padding-top: 12px;
  text-align: center;
  opacity: .7;
  transition: opacity .3s;
}
.code-box-meta .demo-description>h4, .code-box-meta>h4 {
  position: absolute;
  top: -14px;
  padding: 1px 8px;
  margin-left: 16px;
  color: #777;
  border-radius: 2px 2px 0 0;
  background: #fff;
  font-size: 14px;
  width: auto;
}
.code-box {
  border: 1px solid #ebedf0;
  border-radius: 2px;
  display: inline-block;
  width: 100%;
  position: relative;
  margin: 0 0 16px;
  transition: all .2s;
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
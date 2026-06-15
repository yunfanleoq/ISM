<template>
  <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
    <div class="page-header-index-wide">
      <a-card :bordered="false" :bodyStyle="{ padding: '16px 0', height: '100%' }" :style="{ height: '100%' }">
        <div class="account-settings-info-main" :class="{ 'mobile': isMobile }">
          <div class="account-settings-info-left">
            <a-menu
                :class="['avatar-menu']"
                :mode="isMobile ? 'horizontal' : 'inline'"
                :style="{ border: '0', width: isMobile ? '560px' : 'auto'}"
                type="inner"
            >
              <a-menu-item v-for="(net,index) in NetList" :class="openKeys==net.Name?'ant-menu-item ant-menu-item-selected':''" :key="index" @click="onOpenChange(net.Name,index)">
              {{net.Name}}
              </a-menu-item>
            </a-menu>

          </div>
          <div class="account-settings-info-right">
            <div class="account-settings-info-view">
              <a-row :gutter="16" type="flex" justify="center">
                <a-col :order="isMobile ? 2 : 1" :md="24" :lg="16">
                  <a-form :form="NetForm"  layout="vertical">
                    <a-form-item
                        :label="$t('NetWork.SystemNetwork.NetName')"
                    >
                      <a-input disabled v-decorator="[
                    'NetName',
                    {
                      rules: [{ required: true }],
                    },
                  ]"/>
                    </a-form-item>
                    <a-form-item
                        :label="$t('NetWork.SystemNetwork.NetMAC')"
                    >
                      <a-input
                          v-decorator="[
                    'NetMAC',
                    {
                      rules: [{ required: true, message: $t('NetWork.SystemNetwork.NetMAC') }],
                    },
                  ]"
                          :placeholder="$t('NetWork.SystemNetwork.NetMAC')" />
                    </a-form-item>
                    <a-form-item
                        :label="$t('NetWork.SystemNetwork.GateWay')"
                    >
                      <a-input
                          v-decorator="[
                    'GateWay',
                    {
                      rules: [{ required: true, message: $t('NetWork.SystemNetwork.GateWay') }],
                    },
                  ]"
                          :placeholder="$t('NetWork.SystemNetwork.GateWay')" />
                    </a-form-item>
                    <a-form-item
                        :label="$t('NetWork.SystemNetwork.IP')"
                    >
                      <a-button class="editable-add-btn" @click="handleAdd">
                        {{$t('NetWork.SystemNetwork.Add')}}
                      </a-button>
                      <a-table :columns="columns" :data-source="IPList" :pagination="false" :rowKey="(record,index)=>{return index}" bordered>
                        <template v-for="(item, index) in columns" :slot="item.slotName">
                          <span :key="index">{{ $t(item.slotName) }}</span>
                        </template>
                        <template
                            slot="IP"
                            slot-scope="text, record,index"
                        >

                            <a-input
                                style="margin: -5px 0"
                                v-model="record.IP"
                            />

                        </template>

                        <template
                            slot="Mask"
                            slot-scope="text, record,index"
                        >
                          <div >
                            <a-input
                                style="margin: -5px 0"
                                v-model="record.Mask"
                            />
                          </div>
                        </template>

                        <template slot="operation" slot-scope="text, record, index">
                          <div class="editable-row-operations">
                            <span >
                              <a-popconfirm :title="$t('NetWork.SystemNetwork.SureDel')" @confirm="() => DelList(index)">
                                <a>{{$t('NetWork.SystemNetwork.Del')}}</a>
                              </a-popconfirm>
                            </span>
                          </div>
                        </template>
                      </a-table>
                    </a-form-item>
                    <a-form-item>
                      <a-button type="primary" @click="SaveNetwork">{{ $t('NetWork.SystemNetwork.Save') }}</a-button>
                    </a-form-item>
                  </a-form>
                </a-col>
              </a-row>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </a-spin>
</template>

<script>

import {GetSystemNetwork,SaveSystemNetworkInfo,RebootSystem} from "@/services/system";

export default {
  i18n: require('@/i18n/language'),
  data () {
    return {
      // horizontal  inline
      mode: 'inline',
      columns:[
        {
          title: 'IP',
          dataIndex: 'IP',
          width: '25%',
          slotName: 'NetWork.SystemNetwork.SureDel',
          scopedSlots: { customRender: 'IP' ,title: 'dataModel.modelTableIndex' },
        },
        {
          title: 'Mask',
          dataIndex: 'Mask',
          width: '25%',
          slotName: 'NetWork.SystemNetwork.SureDel',
          scopedSlots: { customRender: 'Mask',title: 'dataModel.modelTableIndex' },
        },
        {
          title: 'operation',
          dataIndex: 'operation',
          slotName: 'NetWork.SystemNetwork.SureDel',
          width: '25%',
          scopedSlots: { customRender: 'operation',title: 'dataModel.modelTableIndex' },
        },
      ],
      messageShowLoad:false,
      page: {},
      NetForm:this.$form.createForm(this),
      NetList:[],
      openKeys: "basic",
      isMobile:false,
      selectIndex:0,
      IPList:[],
    }
  },
  components: {

  },
  mounted () {
    this.GetSystemNetwork()
  },
  methods: {
    validateIP(IP) {
      const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
      return ipRegex.test(IP);
    },
    SaveNetwork(){
      let _t = this
      _t.messageShowLoad = true
      let ipValid  = true
      this.NetList[this.selectIndex].MacAddress = this.NetForm.getFieldValue('NetMAC')
      this.NetList[this.selectIndex].GateWay = this.NetForm.getFieldValue('GateWay')

      for(let i=0;i<_t.NetList.length;i++)
      {
        for(let k=0;k<_t.NetList[i].IPv4.length;k++)
        {
            if(!_t.validateIP(_t.NetList[i].IPv4[k].IP))
            {
              ipValid = false
            }
            else if(!_t.validateIP(_t.NetList[i].IPv4[k].Mask))
            {
              ipValid = false
            }
        }
      }
      if(!ipValid)
      {
        this.$error({
          title: _t.$t('NetWork.SystemNetwork.IPError'),
          content: _t.$t('NetWork.SystemNetwork.IPErrorContent'),
        });
        _t.messageShowLoad = false
        return
      }
      SaveSystemNetworkInfo(_t.NetList).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success( _t.$t('NetWork.SystemNetwork.SaveSuccess'))

          _t.$confirm({
            title: _t.$t('NetWork.SystemNetwork.RebootTips'),
            content: _t.$t('NetWork.SystemNetwork.RebootContent'),
            okText:_t.$t('NetWork.SystemNetwork.Reboot'),
            cancelText:_t.$t('NetWork.SystemNetwork.RebootCancel'),
            onOk() {
              _t.RebootSystemFunc()
            },
            onCancel() {
              console.log('Cancel');
            },
            class: 'test',
          });
        } else if(res.data.code==-4){
          _t.$message.error( _t.$t('NetWork.SystemNetwork.SaveARMFailed'))
        }
        else
        {
          _t.$message.error( _t.$t('NetWork.SystemNetwork.SaveFailed'))
        }
      }).catch(function (error) {
        _t.messageShowLoad = false
        console.log(error)
        _t.$message.error( _t.$t('NetWork.SystemNetwork.SaveFailed'))
      }).finally(function (error) {
        _t.messageShowLoad = false
      })
    },
    DelList(index){
      this.IPList.splice(index,1)
    },
    handleAdd() {
      const newData = {
        IP: "192.168.1.2",
        Mask: "255.255.255.0",
      };
      this.IPList.push(newData)
    },
    onOpenChange (openKeys,index) {
      this.openKeys = openKeys
      let netInfo = this.NetList[index]
      if(this.selectIndex!=index) {
        this.NetList[this.selectIndex].MacAddress = this.NetForm.getFieldValue('NetMAC')
        this.NetList[this.selectIndex].GateWay = this.NetForm.getFieldValue('GateWay')
        this.NetList[this.selectIndex].IPv4 = this.IPList
      }
      this.IPList = this.NetList[index].IPv4
      this.NetForm.setFieldsValue(
          {
            NetName: netInfo.Name,
            NetMAC: netInfo.MacAddress,
            GateWay: netInfo.GateWay
          })
      this.selectIndex = index
    },
    GetSystemNetwork(){
      let _t = this
      _t.messageShowLoad = true
      _t.NetList=[]
      GetSystemNetwork().then(function (res){
        if(res.data.code==0)
        {
          for(let i=0;i< res.data.list.length;i++)
          {
            let single = res.data.list[i]
            if(single.IPv4==null)
            {
              single.IPv4=[]
            }
            if(single.IPv6==null)
            {
              single.IPv6=[]
            }
            _t.NetList.push(single)
          }
          _t.onOpenChange(_t.NetList[0].Name,0)
          _t.$message.success( _t.$t('NetWork.SystemNetwork.GetSuccess'))
        }
        else
        {
          _t.$message.error( _t.$t('NetWork.SystemNetwork.GetFailed'))
        }
      }).catch(function (error) {
        _t.messageShowLoad = false
        _t.$message.error( _t.$t('NetWork.SystemNetwork.GetFailed'))
      }).finally(function (error) {
        _t.messageShowLoad = false
      })
    },
    RebootSystemFunc(){
      let _t = this
      _t.messageShowLoad = true
      RebootSystem().then(function (res){
        if(res.data.code==0||res.data.code==-5)
        {
          _t.$message.success( _t.$t('NetWork.SystemNetwork.RebootSuccess'))
        }
        else
        {
          _t.$message.error( _t.$t('NetWork.SystemNetwork.RebootFailed'))
        }
      }).catch(function (error) {
        _t.messageShowLoad = false
        _t.$message.error( _t.$t('NetWork.SystemNetwork.GetFailed'))
      }).finally(function (error) {
        _t.messageShowLoad = false
      })
    },
  },
  watch: {

  }
}
</script>

<style lang="less" scoped>
.editable-add-btn{
  margin-bottom: 5px;
}
.page-header-index-wide{
  padding: 20px;
}
.account-settings-info-main {
  width: 100%;
  display: flex;
  height: 100%;
  overflow: auto;

  &.mobile {
    display: block;

    .account-settings-info-left {
      border-right: unset;
      border-bottom: 1px solid #e8e8e8;
      width: 100%;
      height: 50px;
      overflow-x: auto;
      overflow-y: scroll;
    }
    .account-settings-info-right {
      padding: 20px 40px;
    }
  }

  .account-settings-info-left {
    border-right: 1px solid #e8e8e8;
    width: 224px;
  }

  .account-settings-info-right {
    flex: 1 1;
    padding: 8px 40px;

    .account-settings-info-title {
      color: rgba(0,0,0,.85);
      font-size: 20px;
      font-weight: 500;
      line-height: 28px;
      margin-bottom: 12px;
    }
    .account-settings-info-view {
      padding-top: 12px;
    }
  }
}

</style>

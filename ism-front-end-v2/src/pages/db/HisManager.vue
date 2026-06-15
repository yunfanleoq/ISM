<template>
  <div class="page-header-index-wide">

      <a-card :bordered="false" :bodyStyle="{ padding: '16px', height: '100%' }" :style="{ height: '100%' }">
      <a-tabs default-active-key="3">
        <a-tab-pane key="1" :tab="$t('DbBack.DbConfig')">
          <a-spin :tip="$t('DbBack.Loading')" :spinning="messageShowLoad">
            <a-form :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }">
              <a-form-item :label="$t('DbBack.DbType')">
                <a-select
                    @select="chargeDb" v-model="DbType"
                >
                  <a-select-option value=1>
                    内置
                  </a-select-option>
                  <a-select-option value=2>
                    TDengine
                  </a-select-option>
                  <a-select-option value=3>
                    ClickHouse
                  </a-select-option>
                  <a-select-option value=4>
                    Influxdb
                  </a-select-option>
                  <a-select-option value=5>
                    PostgreSQL
                  </a-select-option>
                </a-select>
              </a-form-item>
<!--              //涛思-->
              <div v-if="DbType==2">
                <a-form-item :label="$t('DbBack.DbServer')">
                  <a-input v-model="td.TDengineHost">
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbServerPort')">
                  <a-input
                      v-model="td.TDenginePort"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbUserName')">
                  <a-input
                      v-model="td.TDengineUserName"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbPassword')">
                  <a-input
                      v-model="td.TDenginePassWord"
                  >
                  </a-input>
                </a-form-item>
              </div>
<!--              ChickHouse-->
              <div v-if="DbType==3">
                <a-form-item :label="$t('DbBack.DbServer')">
                  <a-input
                      v-model="ChickHouse.ChickHouseHost"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbServerPort')">
                  <a-input
                      v-model="ChickHouse.ChickHousePort"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbName')">
                  <a-input
                      v-model="ChickHouse.ChickHouseDataBase"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbUserName')">
                  <a-input
                      v-model="ChickHouse.ChickHouseUserName"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbPassword')">
                  <a-input
                      v-model="ChickHouse.ChickHousePassWord"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.ConnOutTime')">
                  <a-input 
                           v-model="ChickHouse.ChickHouseConnectTimeout"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.ReadOutTime')">
                  <a-input
                           v-model="ChickHouse.ChickHouseReadTimeout"
                  >
                  </a-input>
                </a-form-item>
              </div>
<!--              Influxdb-->
              <div v-if="DbType==4">
                <a-form-item :label="$t('DbBack.InfluxdbUrl')">
                  <a-input
                      v-model="Influxdb.InfluxdbUrl"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.InfluxdbToken')">
                  <a-input
                      v-model="Influxdb.InfluxdbToken"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.InfluxdbOrg')">
                  <a-input
                      v-model="Influxdb.InfluxdbOrg"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.InfluxdbBucket')">
                  <a-input
                      v-model="Influxdb.InfluxdbBucket"
                  >
                  </a-input>
                </a-form-item>
              </div>
<!--              PostgreSQL-->
              <div v-if="DbType==5">
                <a-form-item :label="$t('DbBack.DbServer')">
                  <a-input
                      v-model="pg.PGHost"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbServerPort')">
                  <a-input
                      v-model="pg.PGPort"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbName')">
                  <a-input
                      v-model="pg.PGDbName"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbUserName')">
                  <a-input
                      v-model="pg.PGUser"
                  >
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('DbBack.DbPassword')">
                  <a-input type="password"
                           v-model="pg.PGPassWord"
                  >
                  </a-input>
                </a-form-item>
              </div>
              <a-form-item :wrapper-col="{ span: 12, offset: 5 }">
                <a-button type="primary" html-type="submit" @click="SetDbConfig">
                  {{$t('DbBack.Sure')}}
                </a-button>
              </a-form-item>
            </a-form>
          </a-spin>
        </a-tab-pane>
        <a-tab-pane key="2" :tab="$t('DbBack.Partition')">
          <a-spin :tip="$t('DbBack.Loading')" :spinning="messageShowLoad">
            <a-form :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }">
              <a-form-item :label="$t('DbBack.OneWriteNumbers')">
                <a-input
                    v-model="OnceWriteHistoryCounts"
                >
                </a-input>
              </a-form-item>
              <a-form-item :label="$t('DbBack.PartitionType')" v-if="DbType==1||DbType==5">
                <a-select
                     v-model="PartitionType"
                >
                  <a-select-option value=0>
                   {{ $t('DbBack.PartitionNo')}}
                  </a-select-option>
                  <a-select-option value=1>
                    {{ $t('DbBack.PartitionYear')}}
                  </a-select-option>
                  <a-select-option value=2>
                    {{ $t('DbBack.PartitionMonth')}}
                  </a-select-option>
                  <a-select-option value=3>
                    {{ $t('DbBack.PartitionDay')}}
                  </a-select-option>
                  <a-select-option value=4>
                    {{ $t('DbBack.PartitionHour')}}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :wrapper-col="{ span: 12, offset: 5 }">
                <a-button type="primary" html-type="submit" @click="SetDbConfig">
                  {{$t('DbBack.Sure')}}
                </a-button>
              </a-form-item>
            </a-form>
          </a-spin>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>
<script>
import {GetSystemHistoryConfig, SaveSystemHistoryConfig} from "../../services/system";
export default {
  i18n: require('../../i18n/language'),
  data () {
    return {
      messageShowLoad:false,
      DbType:"1",
      OnceWriteHistoryCounts:100,
      PartitionType:1,
      td:{
        TDengineHost: "127.0.0.1",
        TDenginePassWord: "taosdata",
        TDenginePort: 6041,
        TDengineUserName: "root",
      },
      ChickHouse:{
        ChickHousePort:9000,
        ChickHouseHost:"",
        ChickHouseUserName:"",
        ChickHousePassWord:"",
        ChickHouseDataBase:"",
        ChickHouseConnectTimeout:"",
        ChickHouseReadTimeout:"",
      },
      Influxdb:{
        InfluxdbUrl:"",
        InfluxdbToken:"",
        InfluxdbOrg:"",
        InfluxdbBucket:"",
      },
      pg:{
        PGHost:"",
        PGPort:5432,
        PGDbName:"",
        PGUser:"",
        PGPassWord:"",
      }
    }
  },
  components: {},
  mounted () {
    this.GetDbConfig()
  },
  methods: {
    chargeDb(v){
      this.DbType = v
    },
    GetDbConfig(){
      let _t = this
      let params={

      }
      this.messageShowLoad = true
      GetSystemHistoryConfig(params).then(function (res){
        _t.messageShowLoad = false
        _t.OnceWriteHistoryCounts = res.data.result.OnceWriteHistoryCounts
        _t.DbType = res.data.result.HistoryRecordDbType
        _t.PartitionType = res.data.result.PartitionType
        _t.td.TDengineHost = res.data.result.TDengineHost
        _t.td.TDenginePassWord = res.data.result.TDenginePassWord
        _t.td.TDengineUserName = res.data.result.TDengineUserName
        _t.td.TDenginePort = res.data.result.TDenginePort

        _t.ChickHouse.ChickHousePort = res.data.result.ChickHousePort
        _t.ChickHouse.ChickHouseHost = res.data.result.ChickHouseHost
        _t.ChickHouse.ChickHouseUserName = res.data.result.ChickHouseUserName
        _t.ChickHouse.ChickHousePassWord = res.data.result.ChickHousePassWord
        _t.ChickHouse.ChickHouseDataBase = res.data.result.ChickHouseDataBase
        _t.ChickHouse.ChickHouseConnectTimeout = res.data.result.ChickHouseConnectTimeout
        _t.ChickHouse.ChickHouseReadTimeout = res.data.result.ChickHouseReadTimeout

        _t.Influxdb.InfluxdbUrl = res.data.result.InfluxdbUrl
        _t.Influxdb.InfluxdbToken = res.data.result.InfluxdbToken
        _t.Influxdb.InfluxdbOrg = res.data.result.InfluxdbOrg
        _t.Influxdb.InfluxdbBucket = res.data.result.InfluxdbBucket

        _t.pg.PGHost = res.data.result.PGHost
        _t.pg.PGPort = res.data.result.PGPort
        _t.pg.PGDbName = res.data.result.PGDbName
        _t.pg.PGUser = res.data.result.PGUser
        _t.pg.PGPassWord = res.data.result.PGPassWord
      }).catch(function (error) {
        _t.messageShowLoad = false
      })
    },
    SetDbConfig(){
      let _t = this
      let params={
        HistoryRecordDbType:_t.DbType,
        PartitionType:_t.PartitionType,
        OnceWriteHistoryCounts:_t.OnceWriteHistoryCounts,
        TDengineHost:_t.td.TDengineHost,
        TDenginePassWord:_t.td.TDenginePassWord,
        TDenginePort:_t.td.TDenginePort,
        TDengineUserName:_t.td.TDengineUserName,

        ChickHousePort:_t.ChickHouse.ChickHousePort,
        ChickHouseHost:_t.ChickHouse.ChickHouseHost,
        ChickHouseUserName:_t.ChickHouse.ChickHouseUserName,
        ChickHousePassWord:_t.ChickHouse.ChickHousePassWord,
        ChickHouseDataBase:_t.ChickHouse.ChickHouseDataBase,
        ChickHouseConnectTimeout:_t.ChickHouse.ChickHouseConnectTimeout,
        ChickHouseReadTimeout:_t.ChickHouse.ChickHouseReadTimeout,

        InfluxdbUrl:_t.Influxdb.InfluxdbUrl,
        InfluxdbToken:_t.Influxdb.InfluxdbToken,
        InfluxdbOrg:_t.Influxdb.InfluxdbOrg,
        InfluxdbBucket:_t.Influxdb.InfluxdbBucket,

        PGHost: _t.pg.PGHost,
        PGPort: _t.pg.PGPort,
        PGDbName:_t.pg.PGDbName,
        PGUser:_t.pg.PGUser,
        PGPassWord:_t.pg.PGPassWord,
      }
      this.Configuring = true
      SaveSystemHistoryConfig(params).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t('DbBack.ConfigureSuccess'))
        }
        else if(res.data.code==-3)
        {
          _t.$message.error(_t.$t('DbBack.DbConnectError'))
        }
        else
        {
          _t.$message.error(_t.$t('DbBack.ConfigureFailed'))
        }
      }).finally(function (error) {
        _t.Configuring = false
      })
    },
  },
  watch: {

  }
}
</script>

<style lang="less" scoped>
::v-deep .ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
::v-deep .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  /*background: #f8f8f8;*/
  /*border-bottom: 1px solid #e8e8e8;*/
  transition: background .3s ease;
}
::v-deep .ant-form-item {
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

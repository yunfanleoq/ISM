<template>
<div>
  <a-layout id="components-layout-demo-side" style="height: 100%;">
    <a-drawer
        :title="isDevice?$t('monitor.AddDevice'):$t('monitor.AddZone')"
        :width="720"
        :visible="visible"
        :body-style="{ paddingBottom: '80px' }"
        :after-visible-change="afterVisibleChange"
        @close="onClose"
    >
      <a-form :form="form" layout="vertical" v-if="isDevice">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('device.deviceName')">
              <a-input
                  v-decorator="[
                  'name',
                  {
                    rules: [{ required: true,  validator: isValidateTxtNonSpec,message: $t('device.deviceNameVal') }],
                  },
                ]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('monitor.DeviceType')">
              <a-select :disabled="edit" @change="changeDeviceType"
                        v-decorator="[
                  'DeviceType',
                  {
                    rules: [{ required: true, message: $t('monitor.DeviceType') }],
                  },
                ]"
              >
                <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type>
                  {{ $t(device.name)  }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('device.deviceModelName')">
              <a-select :disabled="edit" @select="changeModelType"
                        v-decorator="[
                  'model',
                  {
                    rules: [{ required: true, message: $t('device.deviceModelName') }],
                  },
                ]"
              >
                <a-select-option v-for="(model,index) in modelList" :key="index" :value="model.uuid">
                  {{ model.modelName }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>

          <a-col :span="12" style="display: none">
            <a-form-item :label="$t('device.deviceConfigurationModelName')">
              <a-select :disabled="edit" @select="GetDisplayPage"
                        v-decorator="[
                  'configurationModel',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationModelName') }],
                  },
                ]"
              >
                <a-select-option v-for="(model,index) in configurationModel" :key="index" :value="model.uuid">
                  {{ model.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12" style="display: none">
            <a-form-item :label="$t('device.deviceConfigurationPageName')">
              <a-select :disabled="edit"
                        v-decorator="[
                  'configurationPageUUID',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationPageName') }],
                  },
                ]"
              >
                <a-select-option v-for="(page,index) in displayPageList" :key="index" :value="page.value">
                  {{ page.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
<!--snmp设备-->
          <div v-if="DeviceType==1">
            <a-col :span="12" >
              <a-form-item :label="$t('device.deviceAgent')">
                <a-input
                    v-decorator="[
                  'agentIpaddress',
                  {
                    rules: [{ required: true, message: $t('device.deviceAgentEnd') }],
                  },
                ]"
                    style="width: 100%"
                />
              </a-form-item>
            </a-col>
          </div>
<!--modbus设备-->
          <a-col :span="12" v-if="DeviceType==2">
            <a-form-item  v-if="modbusConnectType=='TCPServer'">
               <span slot="label">
                {{$t('device.RegisterPack')}}&nbsp;
                 {{$t('device.RegisterPackByte')}}&nbsp;
                 <span style="cursor: pointer" @click="copyFn(RegisterPackByte)">{{RegisterPackByte}}</span>
              </span>
              <a-input-number @change= "RegisterPack($event,0)"
                  v-decorator="[
                  'RegisterPack',
                  {
                    rules: [{ required: true, message: $t('device.RegisterPack') }],
                  },
                ]"
                  style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <div v-if="modbusConnectType=='TCPClient'">
            <a-col :span="12" v-if="DeviceType==2">
              <a-form-item
                  :label="$t('dataModel.modbusModel.IpAddress')"
              >
                <a-input  autocomplete="autocomplete"

                          v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12" v-if="DeviceType==2">
              <a-form-item
                  :label="$t('dataModel.modbusModel.Port')"
              >
                <a-input  autocomplete="autocomplete"

                          v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
          </div>
          <a-col :span="12" v-if="DeviceType==2">
            <a-form-item :label="$t('device.DeviceAddress')">
              <a-input
                  v-decorator="[
                  'DeviceAddress',
                  {
                    rules: [{ required: true, message: $t('device.DeviceAddress') }],
                  },
                ]"
                  style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12" v-if="DeviceType==2">
            <a-form-item :label="$t('device.packTime')">
              <a-input
                  v-decorator="[
                  'packTime',
                  {
                    rules: [{ required: true, message: $t('device.packTime') }],
                  },
                ]"
                  style="width: 100%"
              />
            </a-form-item>
          </a-col>
<!--          OPCUA设备-->
          <a-col :span="12" v-if="DeviceType==3">
            <a-form-item :label="$t('device.OPCUAEndPoint')">
              <a-input  v-decorator="['OPCUAEndPoint', {rules: [{ required: true, message: $t('device.OPCUAEndPoint'), whitespace: true}]}]"/>
            </a-form-item>
          </a-col>
          <!--          西门子S7设备-->
          <div v-if="DeviceType==15">
            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.SimS7Model.IPAddress')">
                <a-input  v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.SimS7Model.Slot'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>
            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.SimS7Model.Slot')">
                <a-input  type="number" min="0" v-decorator="['Slot', {rules: [{ required: true, message: $t('dataModel.SimS7Model.Slot'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>
            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.SimS7Model.Rack')">
                <a-input type="number" min="0" v-decorator="['Rack', {rules: [{ required: true, message: $t('dataModel.SimS7Model.Rack'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>
          </div>
          <!--          Mqtt-->
          <div v-if="DeviceType==20">
            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.SimS7Model.ClientID')">
                <a-input  v-decorator="['ClientID', {rules: [{ required: true, message: $t('dataModel.SimS7Model.ClientID'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>
          </div>
          <!--          DLT645电表-->
          <div v-if="DeviceType==30">
            <a-col :span="12" v-if="DLT645ConnectType=='TCPServer'">
              <a-form-item  >
                 <span slot="label">
                  {{$t('device.RegisterPack')}}&nbsp;
                   {{$t('device.RegisterPackByte')}}&nbsp;
                   <span style="cursor: pointer" @click="copyFn(RegisterPackByte)">{{RegisterPackByte}}</span>
                </span>
                <a-input-number @change= "RegisterPack($event,0)"
                                v-decorator="[
                    'RegisterPack',
                    {
                      rules: [{ required: true, message: $t('device.RegisterPack') }],
                    },
                  ]"
                                style="width: 100%"
                />
              </a-form-item>
            </a-col>
            <div v-if="DLT645ConnectType=='TCPClient'">
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.IpAddress')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.Port')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
            </div>
            <a-col :span="12" >
              <a-form-item :label="$t('device.packTime')">
                <a-input
                    v-decorator="[
                  'packTime',
                  {
                    rules: [{ required: true, message: $t('device.packTime') }],
                  },
                ]"
                    style="width: 100%"
                />
              </a-form-item>
            </a-col>

            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.DLT645Model.ConnectAddress')">
                <a-input  v-decorator="['ConnectAddress', {rules: [{ required: true, message: $t('dataModel.DLT645Model.ConnectAddress'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>

            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.DLT645Model.OperatorCode')">
                <a-input  v-decorator="['OperatorCode', {rules: [{ required: true, message: $t('dataModel.DLT645Model.OperatorCode'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>

            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.DLT645Model.Password')">
                <a-input  v-decorator="['Password', {rules: [{ required: true, message: $t('dataModel.DLT645Model.Password'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>

            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.DLT645Model.BeforeCode')">
                <a-select  v-decorator="['BeforeCode', {rules: [{ required: true, message: $t('dataModel.DLT645Model.BeforeCode'), whitespace: true}]}]">
                  <a-select-option value="1">{{$t('dataModel.DLT645Model.Enable')}}</a-select-option>
                  <a-select-option value="0">{{$t('dataModel.DLT645Model.Disable')}}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>

          </div>
          <!--          IEC104协议-->
          <div v-if="DeviceType==40">
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.IpAddress')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.Port')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
            <a-col :span="12" >
              <a-form-item
                  :label="$t('dataModel.IEC104Model.DeviceAddress')"
              >
                <a-input  autocomplete="autocomplete"

                          v-decorator="['DeviceAddress', {rules: [{ required: true, message: $t('dataModel.IEC104Model.DeviceAddress'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
          </div>

          <!--          IEC61850设备-->
          <a-col :span="12" v-if="DeviceType==350">
            <a-form-item :label="$t('device.IEC61850EndPoint')">
              <a-input  v-decorator="['IEC61850EndPoint', {rules: [{ required: true, message: $t('device.IEC61850EndPoint'), whitespace: true}]}]"/>
            </a-form-item>
          </a-col>
          <a-col :span="12" v-if="DeviceType==350">
            <a-form-item
                :label="$t('dataModel.modbusModel.Port')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['IEC61850Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <!--          HJ212-2017设备-->
          <a-col :span="12" v-if="DeviceType==470">
            <a-form-item :label="$t('device.HJ212DeviceSN')">
              <a-input  v-decorator="['HJ212DeviceSN', {rules: [{ required: true, message: $t('device.HJ212DeviceSN'), whitespace: true}]}]"/>
            </a-form-item>
          </a-col>
          <a-col :span="12" v-if="DeviceType==470">
            <a-form-item
                :label="$t('device.HJ212PW')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['HJ212PW', {rules: [{ required: false, message: $t('device.HJ212PW'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <!--          CJT188仪表-->
          <div v-if="DeviceType==490">
            <a-col :span="12" v-if="CJT188ConnectType=='TCPServer'">
              <a-form-item  >
                 <span slot="label">
                  {{$t('device.RegisterPack')}}&nbsp;
                   {{$t('device.RegisterPackByte')}}&nbsp;
                   <span style="cursor: pointer" @click="copyFn(RegisterPackByte)">{{RegisterPackByte}}</span>
                </span>
                <a-input-number @change= "RegisterPack($event,0)"
                                v-decorator="[
                    'RegisterPack',
                    {
                      rules: [{ required: true, message: $t('device.RegisterPack') }],
                    },
                  ]"
                                style="width: 100%"
                />
              </a-form-item>
            </a-col>
            <div v-if="CJT188ConnectType=='TCPClient'">
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.IpAddress')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.Port')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
            </div>
            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.DLT645Model.ConnectAddress')">
                <a-input  v-decorator="['ConnectAddress', {rules: [{ required: true, message: $t('dataModel.DLT645Model.ConnectAddress'), whitespace: true}]}]"/>
              </a-form-item>
            </a-col>
          </div>
<!--          BACnet协议-->
          <div v-if="DeviceType==500">
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.IpAddress')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['BACnetIPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="12" >
                <a-form-item
                    :label="$t('dataModel.modbusModel.Port')"
                >
                  <a-input  autocomplete="autocomplete"

                            v-decorator="['BACnetPort', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                  />
                </a-form-item>
              </a-col>
          </div>
          <a-col :span="12">
            <a-form-item :label="$t('dataModel.longitude')">
              <a-input  @dblclick="mapVisible=true"  v-model="longitude"/>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('dataModel.latitude')">
              <a-input  @dblclick="mapVisible=true"  v-model="latitude"/>
            </a-form-item>
          </a-col>

          <div v-if="DeviceType!=20">
          <a-col :span="12" >
            <a-form-item :label="$t('dataModel.TimeOut')">
              <a-input type="number" v-decorator="['timeout', {initialValue:'1000',rules: [{ required: true, message: $t('dataModel.TimeOut'), whitespace: true,}]}]"/>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('dataModel.FailedTimes')">
            <a-input type="number" v-decorator="['failedTimes', {initialValue:'5',rules: [{ required: true, message: $t('dataModel.FailedTimes'), whitespace: true}]}]"/>
            </a-form-item>
          </a-col>
          <a-col :span="12" v-if="DeviceType!=5">
            <a-form-item :label="$t('dataModel.Interval')">
              <a-input type="number" v-decorator="['interval', {initialValue:'1000',rules: [{ required: true, message: $t('dataModel.Interval'), whitespace: true}]}]">
              </a-input>
            </a-form-item>
          </a-col>
          </div>
          <a-col :span="12">
            <a-form-item :label="$t('dataModel.offlineClear')" >
              <a-select
                        v-decorator="[
                  'offlineClear',
                  {
                    rules: [{ required: true, message: $t('dataModel.offlineClear') }],
                    initialValue: '2'
                  },
                ]"
              >
                <a-select-option value="1">
                  {{ $t('dataModel.offlineClearIs') }}
                </a-select-option>
                <a-select-option value="2">
                  {{ $t('dataModel.offlineClearNo') }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('dataModel.offlineClearDefault')" >
              <a-input  v-decorator="['offlineDefaultValue', {initialValue:'0',rules: [{ required: true, message: $t('dataModel.offlineClearDefault'), whitespace: true}]}]">
              </a-input>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item :label="$t('device.deviceDec')">
              <a-textarea
                  v-decorator="[
                  'description',
                  {
                    rules: [{ required: false, message: $t('device.deviceDec') }],
                  },
                ]"
                  :rows="4"
              />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <a-form :form="form" layout="vertical"  v-else-if="!isDevice">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('monitor.ZoneName')">
              <a-input
                  v-decorator="[
                  'name',
                  {
                    rules: [{ required: true, message: $t('device.ZoneName') }],
                  },
                ]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('device.deviceConfigurationModelName')">
              <a-select :disabled="edit"
                        v-decorator="[
                  'configurationModel',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationModelName') }],
                  },
                ]"
              >
                <a-select-option v-for="(model,index) in configurationModel" :key="index" :value="model.uuid">
                  {{ model.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('device.deviceConfigurationPageName')">
              <a-select :disabled="edit"
                        v-decorator="[
                  'configurationPageUUID',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationPageName') }],
                  },
                ]"
              >
                <a-select-option v-for="(page,index) in displayPageList" :key="index" :value="page.value">
                  {{ page.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item :label="$t('monitor.ZoneDes')">
              <a-textarea
                  v-decorator="[
                  'description',
                  {
                    rules: [{ required: false, message: $t('device.deviceDec') }],
                  },
                ]"
                  :rows="4"
              />
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
        <a-button v-if="edit" type="primary" :style="{ marginRight: '8px' }" @click="editDeviceOrZone">
          {{$t('device.EditButton')}}
        </a-button>
        <a-button v-else type="primary" :style="{ marginRight: '8px' }" @click="addDeviceOrZone">
          {{$t('device.AddButton')}}
        </a-button>

<!--        <a-button v-if="isDevice" :loading="isTesting" type="primary" :style="{ marginRight: '8px' }" @click="TestConnect('add')">-->
<!--          {{$t('device.TestConnect')}}-->
<!--        </a-button>-->

        <a-button  @click="onClose">
          {{$t('device.CancelButton')}}
        </a-button>
      </div>
    </a-drawer>

    <a-layout-sider  collapsedWidth="20" style="background: #fff;min-width: 300px;">
      <device-tree @onSelect="onSelect" @updateTree="updateTable" ref="deviceTree" style="min-height: 95vh"></device-tree>
    </a-layout-sider>

    <a-layout>
      <a-layout-content style="margin: 0 5px">
        <a-card style="padding: 0px">
          <a-space class="operator">
            <a-button @click="showDrawer('zone')" v-auth:role="`edit`" type="default" icon="plus">{{$t('monitor.newZone')}}</a-button>
            <a-button @click="showDrawer('device')" v-auth:role="`edit`" type="primary" icon="plus">{{$t('monitor.newDevice')}}</a-button>
            <a-button @click="deleteAllRecord" v-auth:role="`edit`" type="danger"><a-icon type="delete"/>{{$t('monitor.DelAllDevice')}}</a-button>
          </a-space>
<!--          <hr style="height:1px;border:none;border-top:1px dashed #13c2c2;" />-->
          <a-table :loading="false" rowKey="no" v-if="!editIsDevice" :row-selection="rowSelection" :pagination="pagination" :columns="columns" :data-source="tableDataSource">
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
            <span slot="deviceType" slot-scope="deviceType,record">
                  <div v-if="record.nodeType==1">
                      <div v-for="GetType in supportDeviceList">
                    <div v-if="GetType.type==deviceType">  {{ $t(GetType.name) }}</div>
                  </div>
                  </div>
                  <div v-else>-</div>

            </span>

            <span slot="Status" slot-scope="Status,record">
                <div v-if="record.nodeType==1" style="color: grey">
                    <div v-if="Status==0" style="color: grey">
                      {{$t('monitor.Offline')}}
                    </div>
                   <div v-else-if="Status==1" style="color: green">
                      {{$t('monitor.Online')}}
                   </div>
                   <div v-else-if="Status==2" style="color: #e60927">
                      {{$t('monitor.NonActivated')}}
                   </div>
                  <div v-else-if="Status==3" style="color: #e60927">
                      {{$t('monitor.deviceStop')}}
                  </div>
                </div>
                <div v-else>
                 -
                </div>
            </span>

            <div slot="action" slot-scope="text, record" >
              <a type="link"  v-auth:role="`edit`" @click="goToEdit(record.extra)" style="cursor: pointer;color: #13C2C2"><a-icon type="edit" /><span style="margin-left: 2px;">{{$t('monitor.nodeEdit')}}</span></a>
              <a-divider type="vertical" /><a type="link"  v-if="record.nodeType==1" v-auth:role="`edit`" @click="goToCopy(record.extra)" style="cursor: pointer;color: #13C2C2"><a-icon type="copy" />  <span style="margin-left: 2px;">{{$t('monitor.nodeCopy')}} </span></a>
              <a-popconfirm v-if="record.extra.pid!=0" :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.key,record.nodeType)" v-auth:role="`edit`">
                <a-divider type="vertical" /> <a-icon slot="icon" type="question-circle-o" style="color: red" />
                <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a> <a-divider type="vertical" />
              </a-popconfirm>

              <a-popconfirm v-auth:role="`edit`" v-if="record.IsEnable&&record.nodeType==1&&record.deviceType!=480"
                            :title="$t('monitor.StopDeviceTips')"
                            :ok-text="$t('component.systemImageModel.delImageYes')"
                            :cancel-text="$t('component.systemImageModel.delImageNo')"
                            @confirm="SetDeviceEnable(false,record.extra.uuid,record.deviceType,record.nodeName)"
              >
                <a  type="link"   style="cursor: pointer;color: #d73165"><a-icon type="stop" /><span style="margin-left: 2px;">{{$t('monitor.NodeStop')}}</span></a>
              </a-popconfirm>
              <a type="link" v-if="!record.IsEnable&&record.nodeType==1&&record.deviceType!=480" v-auth:role="`edit`" @click="SetDeviceEnable(true,record.extra.uuid,record.deviceType,record.nodeName)" style="cursor: pointer;color: #F4A460"><a-icon type="caret-right" /><span style="margin-left: 2px;">{{$t('monitor.NodeStart')}}</span></a>
            </div>
          </a-table>

          <div v-else-if="editIsDevice">
            <a-form :form="editForm" layout="vertical"  >

              <div  v-if="editType==1">
                <a-row :gutter="16" v-auth:role="`edit`">
                  <a-col :span="12">
                    <a-form-item :label="$t('device.deviceName')">
                      <a-input
                          v-decorator="[
                          'name',
                          {
                            rules: [{ required: true,validator: isValidateTxtNonSpec, message: $t('device.deviceNameVal') }],
                          },
                      ]"
                      />
                    </a-form-item>
                  </a-col>
                  <a-col :span="12">
                    <a-form-item :label="$t('dataModel.ZoneList')"  >
                      <a-tree-select
                          show-search
                          tree-node-filter-prop="title"
                          @select="SelectDevice"
                          :dropdown-style="{ 'z-index': 9999999,maxHeight: '400px', overflow: 'auto' }"
                          :tree-data="deviceTreeData"
                          v-model="selectNodeKey"
                          :replace-fields="{ value: 'key',title:'text'}"
                          placeholder="Please select"
                          tree-default-expand-all
                      >
                      </a-tree-select>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12">
                    <a-form-item :label="$t('monitor.DeviceType')">
                      <a-select :disabled=true @change="changeDeviceType"
                                v-decorator="[
                  'DeviceType',
                  {
                    rules: [{ required: true, message: $t('monitor.DeviceType') }],
                  },
                ]"
                      >
                        <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type>
                          {{ $t(device.name) }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12" v-if="DeviceType==5">
                    <a-form-item :label="$t('monitor.DeviceFlag')">
                        <a-input
                            v-decorator="[
                          'DeviceFlag',
                          {
                            rules: [{ required: false, message: $t('monitor.DeviceFlag') }],
                          },
                      ]"
                        />
                    </a-form-item>
                  </a-col>

                  <a-col :span="12">
                    <a-form-item :label="$t('device.deviceModelName')">
                      <a-select :disabled="true"
                                v-decorator="[
                  'model',
                  {
                    rules: [{ required: true, message: $t('device.deviceModelName') }],
                  },
                ]"
                      >
                        <a-select-option v-for="(model,index) in modelList" :key="index" :value="model.uuid">
                          {{ model.modelName }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12" style="display: none">
                    <a-form-item :label="$t('device.deviceConfigurationModelName')">
                      <a-select :disabled="edit"
                                @select="GetDisplayPage"
                                v-decorator="[
                  'configurationModel',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationModelName') }],
                  },
                ]"
                      >
                        <a-select-option v-for="(model,index) in configurationModel" :key="index" :value="model.uuid">
                          {{ model.name }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12" style="display: none">
                    <a-form-item :label="$t('device.deviceConfigurationPageName')">
                      <a-select :disabled="edit"
                                v-decorator="[
                  'configurationPageUUID',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationPageName') }],
                  },
                ]"
                      >
                        <a-select-option v-for="(page,index) in displayPageList" :key="index" :value="page.value">
                          {{ page.label }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </a-col>
                  <!--snmp设备-->
                  <div v-if="DeviceType==1">
                    <a-col :span="12" >
                      <a-form-item :label="$t('device.deviceAgent')">
                        <a-input
                            v-decorator="[
                  'agentIpaddress',
                  {
                    rules: [{ required: true, message: $t('device.deviceAgentEnd') }],
                  },
                ]"
                            style="width: 100%"
                        />
                      </a-form-item>
                    </a-col>
                  </div>
                  <!--modbus设备-->
                  <a-col :span="12" v-if="DeviceType==2&modbusConnectType=='TCPServer'">
                    <a-form-item>
                         <span slot="label">
                          {{$t('device.RegisterPack')}}&nbsp;
                            {{$t('device.RegisterPackByte')}}&nbsp;
                           <span style="cursor: pointer" @click="copyFn(RegisterPackByte)">{{RegisterPackByte}}</span>
                        </span>
                      <a-input-number @change= "RegisterPack($event,1)"
                          v-decorator="[
                  'RegisterPack',
                  {
                    rules: [{ required: true, message: $t('device.RegisterPack') }],
                  },
                ]"
                          style="width: 100%"
                      />
                    </a-form-item>
                  </a-col>
                  <div v-if="DeviceType==2&modbusConnectType=='TCPClient'">
                    <a-col :span="12">
                      <a-form-item
                          :label="$t('dataModel.modbusModel.IpAddress')"
                      >
                        <a-input  autocomplete="autocomplete"

                                  v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                        />
                      </a-form-item>
                    </a-col>
                    <a-col :span="12" v-if="DeviceType==2">
                      <a-form-item
                          :label="$t('dataModel.modbusModel.Port')"
                      >
                        <a-input  autocomplete="autocomplete"

                                  v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                        />
                      </a-form-item>
                    </a-col>
                  </div>
                  <a-col :span="12" v-if="DeviceType==2">
                    <a-form-item :label="$t('device.DeviceAddress')">
                      <a-input
                          v-decorator="[
                  'DeviceAddress',
                  {
                    rules: [{ required: true, message: $t('device.DeviceAddress') }],
                  },
                ]"
                          style="width: 100%"
                      />
                    </a-form-item>
                  </a-col>
                  <a-col :span="12" v-if="DeviceType==2">
                    <a-form-item :label="$t('device.packTime')">
                      <a-input
                          v-decorator="[
                  'packTime',
                  {
                    rules: [{ required: true, message: $t('device.packTime') }],
                  },
                ]"
                          style="width: 100%"
                      />
                    </a-form-item>
                  </a-col>

                  <!--          OPCUA设备-->
                  <a-col :span="12" v-if="DeviceType==3">
                    <a-form-item :label="$t('device.OPCUAEndPoint')">
                      <a-input  v-decorator="['OPCUAEndPoint', {rules: [{ required: true, message: $t('device.OPCUAEndPoint'), whitespace: true}]}]"/>
                    </a-form-item>
                  </a-col>

                  <!--          西门子S7设备-->
                  <div v-if="DeviceType==15">
                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.SimS7Model.IPAddress')">
                        <a-input  v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.SimS7Model.Slot'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>
                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.SimS7Model.Slot')">
                        <a-input type="number" v-decorator="['Slot', {rules: [{ required: true, message: $t('dataModel.SimS7Model.Slot'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>
                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.SimS7Model.Rack')">
                        <a-input type="number" v-decorator="['Rack', {rules: [{ required: true, message: $t('dataModel.SimS7Model.Rack'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>
                  </div>
                  <!--          Mqtt-->
                  <div v-if="DeviceType==20">
                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.SimS7Model.ClientID')">
                        <a-input  v-decorator="['ClientID', {rules: [{ required: true, message: $t('dataModel.SimS7Model.ClientID'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>
                  </div>
                  <!--          DLT645电表-->
                  <div v-if="DeviceType==30">
                    <a-col :span="12" v-if="DLT645ConnectType=='TCPServer'">
                      <a-form-item  >
                 <span slot="label">
                  {{$t('device.RegisterPack')}}&nbsp;
                   {{$t('device.RegisterPackByte')}}&nbsp;
                   <span style="cursor: pointer" @click="copyFn(RegisterPackByte)">{{RegisterPackByte}}</span>
                </span>
                        <a-input-number @change= "RegisterPack($event,1)"
                                        v-decorator="[
                    'RegisterPack',
                    {
                      rules: [{ required: true, message: $t('device.RegisterPack') }],
                    },
                  ]"
                                        style="width: 100%"
                        />
                      </a-form-item>
                    </a-col>
                    <div v-if="DLT645ConnectType=='TCPClient'">
                      <a-col :span="12" >
                        <a-form-item
                            :label="$t('dataModel.modbusModel.IpAddress')"
                        >
                          <a-input  autocomplete="autocomplete"

                                    v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                          />
                        </a-form-item>
                      </a-col>
                      <a-col :span="12" >
                        <a-form-item
                            :label="$t('dataModel.modbusModel.Port')"
                        >
                          <a-input  autocomplete="autocomplete"

                                    v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                          />
                        </a-form-item>
                      </a-col>
                    </div>

                    <a-col :span="12" >
                      <a-form-item :label="$t('device.packTime')">
                        <a-input
                            v-decorator="[
                  'packTime',
                  {
                    rules: [{ required: true, message: $t('device.packTime') }],
                  },
                ]"
                            style="width: 100%"
                        />
                      </a-form-item>
                    </a-col>

                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.DLT645Model.ConnectAddress')">
                        <a-input  v-decorator="['ConnectAddress', {rules: [{ required: true, message: $t('dataModel.DLT645Model.ConnectAddress'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>

                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.DLT645Model.OperatorCode')">
                        <a-input  v-decorator="['OperatorCode', {rules: [{ required: true, message: $t('dataModel.DLT645Model.OperatorCode'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>

                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.DLT645Model.Password')">
                        <a-input  v-decorator="['Password', {rules: [{ required: true, message: $t('dataModel.DLT645Model.Password'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>

                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.DLT645Model.BeforeCode')">
                        <a-select  v-decorator="['BeforeCode', {rules: [{ required: true, message: $t('dataModel.DLT645Model.BeforeCode'), whitespace: true}]}]">
                          <a-select-option value="1">{{$t('dataModel.DLT645Model.Enable')}}</a-select-option>
                          <a-select-option value="0">{{$t('dataModel.DLT645Model.Disable')}}</a-select-option>
                        </a-select>
                      </a-form-item>
                    </a-col>

                  </div>
                  <!--          IEC104协议-->
                  <div v-if="DeviceType==40">
                    <a-col :span="12" >
                      <a-form-item
                          :label="$t('dataModel.modbusModel.IpAddress')"
                      >
                        <a-input  autocomplete="autocomplete"

                                  v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                        />
                      </a-form-item>
                    </a-col>
                    <a-col :span="12" >
                      <a-form-item
                          :label="$t('dataModel.modbusModel.Port')"
                      >
                        <a-input  autocomplete="autocomplete"

                                  v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                        />
                      </a-form-item>
                    </a-col>
                    <a-col :span="12" >
                      <a-form-item
                          :label="$t('dataModel.IEC104Model.DeviceAddress')"
                      >
                        <a-input  autocomplete="autocomplete"

                                  v-decorator="['DeviceAddress', {rules: [{ required: true, message: $t('dataModel.IEC104Model.DeviceAddress'), whitespace: true}]}]"
                        />
                      </a-form-item>
                    </a-col>
                  </div>

                  <!--          IEC61850设备-->
                  <a-col :span="12" v-if="DeviceType==350">
                    <a-form-item :label="$t('device.IEC61850EndPoint')">
                      <a-input  v-decorator="['IEC61850EndPoint', {rules: [{ required: true, message: $t('device.IEC61850EndPoint'), whitespace: true}]}]"/>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12" v-if="DeviceType==350">
                    <a-form-item
                        :label="$t('dataModel.modbusModel.Port')"
                    >
                      <a-input  autocomplete="autocomplete"

                                v-decorator="['IEC61850Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                      />
                    </a-form-item>
                  </a-col>

                  <!--          HJ212-2017设备-->
                  <a-col :span="12" v-if="DeviceType==470">
                    <a-form-item :label="$t('device.HJ212DeviceSN')">
                      <a-input  v-decorator="['HJ212DeviceSN', {rules: [{ required: true, message: $t('device.HJ212DeviceSN'), whitespace: true}]}]"/>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12" v-if="DeviceType==470">
                    <a-form-item
                        :label="$t('device.HJ212PW')"
                    >
                      <a-input  autocomplete="autocomplete"

                                v-decorator="['HJ212PW', {rules: [{ required: true, message: $t('device.HJ212PW'), whitespace: true}]}]"
                      />
                    </a-form-item>
                  </a-col>
                  <!--          CJT188-->
                  <div v-if="DeviceType==490">
                    <a-col :span="12" v-if="CJT188ConnectType=='TCPServer'">
                      <a-form-item  >
                 <span slot="label">
                  {{$t('device.RegisterPack')}}&nbsp;
                   {{$t('device.RegisterPackByte')}}&nbsp;
                   <span style="cursor: pointer" @click="copyFn(RegisterPackByte)">{{RegisterPackByte}}</span>
                </span>
                        <a-input-number @change= "RegisterPack($event,1)"
                                        v-decorator="[
                    'RegisterPack',
                    {
                      rules: [{ required: true, message: $t('device.RegisterPack') }],
                    },
                  ]"
                                        style="width: 100%"
                        />
                      </a-form-item>
                    </a-col>
                    <div v-if="CJT188ConnectType=='TCPClient'">
                      <a-col :span="12" >
                        <a-form-item
                            :label="$t('dataModel.modbusModel.IpAddress')"
                        >
                          <a-input  autocomplete="autocomplete"

                                    v-decorator="['IPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                          />
                        </a-form-item>
                      </a-col>
                      <a-col :span="12" >
                        <a-form-item
                            :label="$t('dataModel.modbusModel.Port')"
                        >
                          <a-input  autocomplete="autocomplete"

                                    v-decorator="['Port', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                          />
                        </a-form-item>
                      </a-col>
                    </div>
                    <a-col :span="12" >
                      <a-form-item :label="$t('dataModel.DLT645Model.ConnectAddress')">
                        <a-input  v-decorator="['ConnectAddress', {rules: [{ required: true, message: $t('dataModel.DLT645Model.ConnectAddress'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>
                  </div>

                  <!--          BACnet协议-->
                  <div v-if="DeviceType==500">
                    <a-col :span="12" >
                      <a-form-item
                          :label="$t('dataModel.modbusModel.IpAddress')"
                      >
                        <a-input  autocomplete="autocomplete"

                                  v-decorator="['BACnetIPAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.IpAddress'), whitespace: true}]}]"
                        />
                      </a-form-item>
                    </a-col>
                    <a-col :span="12" >
                      <a-form-item
                          :label="$t('dataModel.modbusModel.Port')"
                      >
                        <a-input  autocomplete="autocomplete"

                                  v-decorator="['BACnetPort', {rules: [{ required: true, message: $t('dataModel.modbusModel.Port'), whitespace: true}]}]"
                        />
                      </a-form-item>
                    </a-col>
                  </div>

                  <div v-if="DeviceType!=20">
                    <a-col :span="12">
                      <a-form-item :label="$t('dataModel.TimeOut')">
                        <a-input  v-decorator="['timeout', {rules: [{ required: true, message: $t('dataModel.TimeOut'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>
                    <a-col :span="12">
                      <a-form-item :label="$t('dataModel.FailedTimes')">
                        <a-input  v-decorator="['failedTimes', {rules: [{ required: true, message: $t('dataModel.FailedTimes'), whitespace: true}]}]"/>
                      </a-form-item>
                    </a-col>
                    <a-col :span="12">
                      <a-form-item :label="$t('dataModel.Interval')" v-if="DeviceType!=5">
                        <a-input v-decorator="['interval', {rules: [{ required: true, message: $t('dataModel.Interval'), whitespace: true}]}]">
                        </a-input>
                      </a-form-item>
                    </a-col>
                  </div>

                  <a-col :span="12">
                    <a-form-item :label="$t('dataModel.offlineClear')" >
                      <a-select
                          v-decorator="[
                  'offlineClear',
                  {
                    rules: [{ required: true, message: $t('dataModel.offlineClear') }],
                    initialValue: '2'
                  },
                ]"
                      >
                        <a-select-option value="1">
                          {{ $t('dataModel.offlineClearIs') }}
                        </a-select-option>
                        <a-select-option value="2">
                          {{ $t('dataModel.offlineClearNo') }}
                        </a-select-option>
                      </a-select>

                    </a-form-item>
                  </a-col>
                  <a-col :span="12">
                    <a-form-item :label="$t('dataModel.offlineClearDefault')" >
                      <a-input  v-decorator="['offlineDefaultValue', {initialValue:'0',rules: [{ required: true, message: $t('dataModel.offlineClearDefault'), whitespace: true}]}]">
                      </a-input>
                    </a-form-item>
                  </a-col>

                  <a-col :span="12">
                    <a-form-item :label="$t('dataModel.longitude')">
                      <a-input  @dblclick="mapVisible=true"  v-model="longitude"/>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12">
                    <a-form-item :label="$t('dataModel.latitude')">
                      <a-input  @dblclick="mapVisible=true"  v-model="latitude"/>
                    </a-form-item>
                  </a-col>
                </a-row>
              </div>

              <div v-else-if="editType==0" v-auth:role="`edit`">
                <a-row :gutter="16" >
                  <a-col :span="12">
                    <a-form-item :label="$t('monitor.ZoneName')">
                      <a-input
                          v-decorator="[
                  'name',
                  {
                    rules: [{ required: true, message: $t('monitor.ZoneName') }],
                  },
                ]"
                      />
                    </a-form-item>
                  </a-col>
                  <a-col :span="12">
                    <a-form-item :label="$t('device.deviceConfigurationModelName')">
                      <a-select :disabled="edit"
                                v-decorator="[
                  'configurationModel',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationModelName') }],
                  },
                ]"
                      >
                        <a-select-option v-for="(model,index) in configurationModel" :key="index" :value="model.uuid">
                          {{ model.name }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12">
                    <a-form-item :label="$t('device.deviceConfigurationPageName')">
                      <a-select :disabled="edit"
                                v-decorator="[
                  'configurationPageUUID',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationPageName') }],
                  },
                ]"
                      >
                        <a-select-option v-for="(page,index) in displayPageList" :key="index" :value="page.value">
                          {{ page.label }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </a-col>
                </a-row>
              </div>

              <a-row :gutter="16" v-auth:role="`edit`">
                <a-col :span="24">
                  <a-form-item :label="$t('device.deviceDec')">
                    <a-textarea
                        v-decorator="[
                  'description',
                  {
                    rules: [{ required: false, message: $t('device.deviceDec') }],
                  },
                ]"
                        :rows="4"
                    />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-form-item style="margin-top: 24px" :wrapperCol="{span: 10, offset: 7}">
                <a-button type="primary" v-auth:role="`edit`" @click="editDeviceOrZone" htmlType="submit">{{$t('dataModel.edit')}}</a-button>
<!--                <a-button  style="margin-left: 8px" v-if="editIsDevice" :loading="isTesting" type="primary" :style="{ marginRight: '8px' }" @click="TestConnect('edit')">-->
<!--                  {{$t('device.TestConnect')}}-->
<!--                </a-button>-->
                <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(editUuid,1)" v-auth:role="`edit`">
                  <a-button style="margin-left: 8px" type="danger">{{$t('dataModel.delete')}}</a-button>
                </a-popconfirm>
                <a-button style="margin-left: 8px" @click="backZoneList()">{{$t('dataModel.back')}}</a-button>
              </a-form-item>
            </a-form>
          </div>
        </a-card>
      </a-layout-content>
    </a-layout>

  </a-layout>
  <a-modal
      title="拾取坐标"
      :footer="null"
      :zIndex=999999
      v-model="mapVisible"
      v-drag-modal
  >
    <baidu-map
        style="display:flex;flex-direction: column-reverse;"
        id="allmap"
        @ready="mapReady"
        @click="getLocation"
        :scroll-wheel-zoom="true"
    >
      <div style="display:flex;justify-content:center;margin:15px">
        <bm-auto-complete v-model="searchJingwei" :sugStyle="{zIndex: 999999}">
          <a-input v-model="searchJingwei" style="width:300px;margin-right:15px" placeholder="输入地址"></a-input>
        </bm-auto-complete>
        <a-button type="primary" @click="getBaiduMapPoint">搜索</a-button>
      </div>
      <bm-map-type :map-types="['BMAP_NORMAL_MAP', 'BMAP_HYBRID_MAP']" anchor="BMAP_ANCHOR_TOP_LEFT"></bm-map-type>
      <bm-marker v-if="infoWindowShow" :position="{lng: longitude, lat: latitude}">
        <bm-label content="" :labelStyle="{color: 'red', fontSize : '24px'}" :offset="{width: -35, height: 30}"/>
      </bm-marker>
      <bm-info-window :position="{lng: longitude, lat: latitude}" :show="infoWindowShow" @clickclose="infoWindowClose">
        <p>纬度:{{latitude}}</p>
        <p>经度:{{longitude}}</p>
      </bm-info-window>
    </baidu-map>
  </a-modal>
  <a-modal
      :title="$t('monitor.nodeCopyTitle')"
      :visible="CopyVisible"
      @ok="CopyHandleOk"
      @cancel="CopyVisible=false"
  >
    <a-form :form="form" :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }">
      <a-form-item :label="$t('monitor.nodeCopyNumber')">
        <a-input v-model="CopyNumber"/>
      </a-form-item>
      <a-form-item :label="$t('monitor.nodeCopyType')">
        <a-select :disabled=true v-model="CopyDeviceType">
          <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type>
            {{ $t(device.name) }}
          </a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</div>
</template>
<script>

import {snmpModelList} from "../../services/snmpmodel";
import {
  addMonitor, CopyDevices, delAllMonitor,
  delMonitor,
  editMonitor, getMonitorTree,
  getSupportDeviceList,
  SetDeviceStartOrStop,
  TestConnect
} from "../../services/device";
import deviceTree from '../../components/deviceTree/DeviceTree'
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";
import {modbusModelRegisterDel} from "@/services/modbusModel";
let CRC = {};

CRC._auchCRCHi = [
  0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0,
  0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41,
  0x00, 0xC1, 0x81, 0x40, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0,
  0x80, 0x41, 0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1, 0x81, 0x40,
  0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1,
  0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0, 0x80, 0x41,
  0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1,
  0x81, 0x40, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41,
  0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0,
  0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x00, 0xC1, 0x81, 0x40,
  0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1,
  0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1, 0x81, 0x40,
  0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0,
  0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x00, 0xC1, 0x81, 0x40,
  0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0,
  0x80, 0x41, 0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1, 0x81, 0x40,
  0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0,
  0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41,
  0x00, 0xC1, 0x81, 0x40, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0,
  0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41,
  0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0,
  0x80, 0x41, 0x00, 0xC1, 0x81, 0x40, 0x00, 0xC1, 0x81, 0x40,
  0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0, 0x80, 0x41, 0x00, 0xC1,
  0x81, 0x40, 0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41,
  0x00, 0xC1, 0x81, 0x40, 0x01, 0xC0, 0x80, 0x41, 0x01, 0xC0,
  0x80, 0x41, 0x00, 0xC1, 0x81, 0x40
];
CRC._auchCRCLo = [
  0x00, 0xC0, 0xC1, 0x01, 0xC3, 0x03, 0x02, 0xC2, 0xC6, 0x06,
  0x07, 0xC7, 0x05, 0xC5, 0xC4, 0x04, 0xCC, 0x0C, 0x0D, 0xCD,
  0x0F, 0xCF, 0xCE, 0x0E, 0x0A, 0xCA, 0xCB, 0x0B, 0xC9, 0x09,
  0x08, 0xC8, 0xD8, 0x18, 0x19, 0xD9, 0x1B, 0xDB, 0xDA, 0x1A,
  0x1E, 0xDE, 0xDF, 0x1F, 0xDD, 0x1D, 0x1C, 0xDC, 0x14, 0xD4,
  0xD5, 0x15, 0xD7, 0x17, 0x16, 0xD6, 0xD2, 0x12, 0x13, 0xD3,
  0x11, 0xD1, 0xD0, 0x10, 0xF0, 0x30, 0x31, 0xF1, 0x33, 0xF3,
  0xF2, 0x32, 0x36, 0xF6, 0xF7, 0x37, 0xF5, 0x35, 0x34, 0xF4,
  0x3C, 0xFC, 0xFD, 0x3D, 0xFF, 0x3F, 0x3E, 0xFE, 0xFA, 0x3A,
  0x3B, 0xFB, 0x39, 0xF9, 0xF8, 0x38, 0x28, 0xE8, 0xE9, 0x29,
  0xEB, 0x2B, 0x2A, 0xEA, 0xEE, 0x2E, 0x2F, 0xEF, 0x2D, 0xED,
  0xEC, 0x2C, 0xE4, 0x24, 0x25, 0xE5, 0x27, 0xE7, 0xE6, 0x26,
  0x22, 0xE2, 0xE3, 0x23, 0xE1, 0x21, 0x20, 0xE0, 0xA0, 0x60,
  0x61, 0xA1, 0x63, 0xA3, 0xA2, 0x62, 0x66, 0xA6, 0xA7, 0x67,
  0xA5, 0x65, 0x64, 0xA4, 0x6C, 0xAC, 0xAD, 0x6D, 0xAF, 0x6F,
  0x6E, 0xAE, 0xAA, 0x6A, 0x6B, 0xAB, 0x69, 0xA9, 0xA8, 0x68,
  0x78, 0xB8, 0xB9, 0x79, 0xBB, 0x7B, 0x7A, 0xBA, 0xBE, 0x7E,
  0x7F, 0xBF, 0x7D, 0xBD, 0xBC, 0x7C, 0xB4, 0x74, 0x75, 0xB5,
  0x77, 0xB7, 0xB6, 0x76, 0x72, 0xB2, 0xB3, 0x73, 0xB1, 0x71,
  0x70, 0xB0, 0x50, 0x90, 0x91, 0x51, 0x93, 0x53, 0x52, 0x92,
  0x96, 0x56, 0x57, 0x97, 0x55, 0x95, 0x94, 0x54, 0x9C, 0x5C,
  0x5D, 0x9D, 0x5F, 0x9F, 0x9E, 0x5E, 0x5A, 0x9A, 0x9B, 0x5B,
  0x99, 0x59, 0x58, 0x98, 0x88, 0x48, 0x49, 0x89, 0x4B, 0x8B,
  0x8A, 0x4A, 0x4E, 0x8E, 0x8F, 0x4F, 0x8D, 0x4D, 0x4C, 0x8C,
  0x44, 0x84, 0x85, 0x45, 0x87, 0x47, 0x46, 0x86, 0x82, 0x42,
  0x43, 0x83, 0x41, 0x81, 0x80, 0x40
];

CRC.CRC16 = function (buffer) {
  let hi = 0xff;
  let lo = 0xff;
  for (let i = 0; i < buffer.length; i++) {
    let idx = hi ^ buffer[i];
    hi = (lo ^ CRC._auchCRCHi[idx]);
    lo = CRC._auchCRCLo[idx];
  }
  return CRC.padLeft((hi << 8 | lo).toString(16).toUpperCase(), 4, '0');
};

CRC.isArray = function (arr) {
  return Object.prototype.toString.call(arr) === '[object Array]';
};

CRC.ToCRC16 = function (str) {
  return CRC.CRC16(CRC.isArray(str) ? str : CRC.strToByte(str));
};

CRC.ToModbusCRC16 = function (str) {
  return CRC.CRC16(CRC.isArray(str) ? str : CRC.strToHex(str));
};

CRC.strToByte = function (str) {
  let tmp = str.split(''), arr = [];
  for (let i = 0, c = tmp.length; i < c; i++) {
    let j = encodeURI(tmp[i]);
    if (j.length == 1) {
      arr.push(j.charCodeAt());
    } else {
      let b = j.split('%');
      for (let m = 1; m < b.length; m++) {
        arr.push(parseInt('0x' + b[m]));
      }
    }
  }
  return arr;
};

CRC.convertChinese = function (str) {
  let tmp = str.split(''), arr = [];
  for (let i = 0, c = tmp.length; i < c; i++) {
    let s = tmp[i].charCodeAt();
    if (s <= 0 || s >= 127) {
      arr.push(s.toString(16));
    }
    else {
      arr.push(tmp[i]);
    }
  }
  return arr;
};

CRC.filterChinese = function (str) {
  let tmp = str.split(''), arr = [];
  for (let i = 0, c = tmp.length; i < c; i++) {
    let s = tmp[i].charCodeAt();
    if (s > 0 && s < 127) {
      arr.push(tmp[i]);
    }
  }
  return arr;
};

CRC.strToHex = function (hex, isFilterChinese) {
  hex = isFilterChinese ? CRC.filterChinese(hex).join('') : CRC.convertChinese(hex).join('');

  //清除所有空格
  hex = hex.replace(/\s/g, "");
  //若字符个数为奇数，补一个空格
  hex += hex.length % 2 != 0 ? " " : "";

  let c = hex.length / 2, arr = [];
  for (let i = 0; i < c; i++) {
    arr.push(parseInt(hex.substr(i * 2, 2), 16));
  }
  return arr;
};

CRC.padLeft = function (s, w, pc) {
  if (pc == undefined) {
    pc = '0';
  }
  for (let i = 0, c = w - s.length; i < c; i++) {
    s = pc + s;
  }
  return s;
};
const loadingKey = 'updatable'
export default {
  name: 'ISMDeviceConfig',
  i18n: require('../../i18n/language'),
  data() {
    return {
      selectNodeKey:"",
      deviceTreeData:[],
      CopyNumber:2,
      CopyDeviceType:"",
      CopyType:0,
      CopyUuid:"",
      CopyVisible:false,
      searchJingwei:'',
      infoWindowShow:false,
      longitude:'',
      latitude:'',
      point:'',
      mapVisible:false,
      isTesting:false,
      selectNodeSid:-6535,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      rowSelection:{
        columnWidth: 10,
        columnAlign: 'center',
        onSelect:this.onDataTableSelect,
        onSelectAll:this.onDataTableSelectAll
      },
      RegisterPackByte:"",
      displayPageList:[],
      dataList:[],
      configurationModel:[],
      editType:0,
      selectNode:null,
      selectDataTableUuid:[],
      selectKey:null,
      form: this.$form.createForm(this),
      editForm: this.$form.createForm(this),
      edit:false,
      DeviceType:0,
      editIsDevice:false,
      editUuid:"",
      isDevice:false,
      visible: false,
      expandedKeys: [],
      supportDeviceList:[],
      modbusConnectType:"",
      DLT645ConnectType:"",
      CJT188ConnectType:"",
      columns: [
        {
          width: '5%',
          slotName: 'monitor.nodeIndex',
          scopedSlots: { customRender: 'nodeIndex', title: 'monitor.nodeIndex' },
          dataIndex: 'no'
        },
        {
          width: '15%',
          slotName: 'monitor.nodeName',
          scopedSlots: { customRender: 'serial', title: 'monitor.nodeName' },
          dataIndex: 'nodeName',
        },
        {
          slotName: 'monitor.nodeType',
          width: '7%',
          scopedSlots: { customRender: 'nodeType', title: 'monitor.nodeType' },
          dataIndex: 'nodeType',
        },
        {
          slotName: 'monitor.DeviceType',
          width: '10%',
          scopedSlots: { customRender: 'deviceType', title: 'monitor.DeviceType' },
          dataIndex: 'deviceType',
        },
        {
          slotName: 'monitor.Status',
          width: '5%',
          scopedSlots: { customRender: 'Status', title: 'monitor.Status' },
          dataIndex: 'Status',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          width: '20%',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      tableDataSource: [],
      selectedRows: [],

      searchValue: '',
      modelList: [],
      autoExpandParent: true,
      treeData:[]
    };
  },
  components: {
    deviceTree,
  },
  mounted(){
    this.getSupportDevice()
    this.getConfigurationModel()
    this.getMonitorTree()
  },
  watch: {
    '$route' () {
      this.getSupportDevice()
      this.getConfigurationModel()
    }
  },
  methods: {
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
    SelectDevice(ev,node){
      const info = node.$options.propsData.dataRef.value
      if(info.type!=0)
      {
        this.$message.error(this.$t("dataModel.ZoneListError"))
        this.selectNodeSid = -6535
      }
      else
      {
        this.selectNodeSid = info.sid
      }
    },
    deleteAllRecord() {
      let _t = this
      this.$confirm({
        content: _t.$t('dataModel.deleteConfirm'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          if(_t.selectDataTableUuid.length==0)
          {
            return
          }
          let params={
            uuid:_t.selectDataTableUuid
          }
          _t.$message.loading({ content: 'Waiting...',loadingKey,duration: 0 });
          delAllMonitor(params).then(function (res) {
            _t.$message.destroy()
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t("dataModel.deleteSuccess"));
              for(let i=0;i<_t.selectDataTableUuid.length;i++)
              {
                _t.tableDataSource = _t.tableDataSource.filter(item => item.key !== _t.selectDataTableUuid[i])
                _t.selectedRows = _t.selectedRows.filter(item => item.key !== _t.selectDataTableUuid[i])
              }
              _t.$refs.deviceTree.getMonitorTree()
              _t.selectDataTableUuid=[]
            }
            else {
              _t.$message.error(_t.$t("dataModel.deleteFailed"));
            }
          })
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });
    },
    findDataTableSelectIndex (key) {
      for(let i=0;i<this.selectDataTableUuid.length;i++)
      {
        if(this.selectDataTableUuid[i]==key)
        {
          return i
        }
      }
      return -1
    },
    onDataTableSelect (record, selected, selectedRows)  {
      if(record.nodeType==0)
      {
        return
      }

      if(selected)
      {
        this.selectDataTableUuid.push(record.key)
      }
      else
      {
        const index=this.findDataTableSelectIndex(record.key)
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
          if(selectedRows[i].nodeType==0)
          {
            continue
          }
          const index=this.findDataTableSelectIndex(selectedRows[i].key)
          if(index==-1)
          {
            this.selectDataTableUuid.push(selectedRows[i].key)
          }
        }
      }
      else
      {
        for(let i=0;i<changeRows.length;i++)
        {
          const index=this.findDataTableSelectIndex(changeRows[i].key)
          if(index!=-1)
          {
            this.selectDataTableUuid.splice(index,1)
          }
        }
      }
    },
    CopyHandleOk(){
      let _t =this
      let  params = {
        CopyType:parseInt(this.CopyType),
        CopyUuid:this.CopyUuid,
        CopyDeviceType:parseInt(this.CopyDeviceType),
        CopyCount:parseInt(this.CopyNumber),
      }
      if(this.$refs.deviceTree.checkIsEmpty()||this.selectNode==null)
      {
        params.pid=0
      }
      else
      {
        params.pid=this.selectNode.value.sid
      }
      this.$message.loading({ content: 'Waiting...',loadingKey,duration: 0 });
      CopyDevices(params).then(function (res){
        if(res.data.code==0)
        {
          _t.$refs.deviceTree.getMonitorTree()
          _t.CopyVisible = false;
          _t.$message.success(_t.$t("device.DeviceAddSuccess"));
        }
        else if(res.data.code==3001)
        {
          _t.$message.error(_t.$t("device.DeviceExist"));
        }
        else if(res.data.code==3003)
        {
          _t.$message.error(_t.$t("device.DeviceAddressExist"));
        }
        else
        {
          _t.$message.error(_t.$t("device.DeviceAddFailed"));
        }
        setTimeout(function (){
          _t.$message.destroy()
        },500)
      }).catch(function (){
        _t.$message.error(_t.$t("device.DeviceAddFailed"));
        setTimeout(function (){
          _t.$message.destroy()
        },500)
      })
    },
    isSpec(s) {
      let pattern = /[~!@#$%^&*<>|'-]/gi
      return pattern.test(s)
    },
    isValidateTxtNonSpec (rule, value, callback) {
      if (value != null && value !== '') {
        let numStr = value.charAt(0);
        if ((this.isSpec(value)) || (value.indexOf(' ') !== -1)||(!isNaN(parseFloat(numStr)) && isFinite(numStr))) {
          callback(new Error('不能包含特殊字符或空格'))
        } else {
          callback()
        }
      } else {
        callback()
      }
    },
    //地图初始化
    mapReady({ BMap, map }) {
      // 选择一个经纬度作为中心点
      this.point = new BMap.Point(113.27, 23.13);
      map.centerAndZoom(this.point, 12);
      this.BMap=BMap
      this.map=map
    },
    //点击获取经纬度
    getLocation(e){
      this.longitude=e.point.lng
      this.latitude=e.point.lat
      this.infoWindowShow=true
    },
    getBaiduMapPoint(){
      let that=this
      let myGeo = new this.BMap.Geocoder()
      myGeo.getPoint(this.searchJingwei,function(point){
        if(point){
          that.map.centerAndZoom(point,15)
          that.latitude=point.lat
          that.longitude=point.lng
          that.infoWindowShow=true
        }

      })
    },
    //信息窗口关闭
    infoWindowClose(){
      this.latitude=''
      this.longitude=''
      this.infoWindowShow=false
    },
    copyFn(val) {
      // createElement() 方法通过指定名称创建一个元素
      let copyInput = document.createElement("input");
      //val是要复制的内容
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
    RegisterPack(e,type){
      this.$nextTick(function (){
        let RegisterPackValue=""
        if(type==1)
        {
          RegisterPackValue = this.editForm.getFieldValue('RegisterPack')
        }
        else
        {
          RegisterPackValue = this.form.getFieldValue('RegisterPack')
        }
        this.RegisterPackByte = this.Value2Bytes(RegisterPackValue)
      })

    },
    Value2Bytes(str) {
      str = parseInt(str);
      str = str.toString(16);
      let hex = str;
      for(let len = (hex + "").length; len < 8; len = hex.length) {
        hex = "0" + hex;
      }
      hex ="FF"+hex
      hex = hex+CRC.ToModbusCRC16(hex, true)
      let byteStream = ""
      let i = 0;
      while (hex.length >= 2) {
        let x = hex.substring(0, 2);
        hex = hex.substring(2, hex.length);
        byteStream = byteStream+x+" ";
      }
      return byteStream
    },
    getSupportDevice(){
      let _t = this
      getSupportDeviceList().then(function (res){
        _t.supportDeviceList =res.data.list
      })
    },
    getConfigurationModel(){
      this.configurationModel=[]
      let _t = this
      const  params= {
        DisplayType: 1
      }
      displayModelList(params).then(function (res){
        let tableData={}
        if(res.data.list!=null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            tableData.name = res.data.list[i].name
            tableData.description = res.data.list[i].description
            tableData.uuid = res.data.list[i].displayUid
            _t.configurationModel.push(tableData)
            tableData={}
          }
        }

      })
    },
    GetDisplayPage(uuid){
      let params={
        muid:uuid
      }
      let _t = this
      getDisplayModelLayerData(params).then(function (res){
        _t.displayPageList = []
        if(res.data.code==0)
        {
          let pageLayer = res.data.layer
          if(pageLayer.length>0)
          {
            for(let i=0;i<pageLayer.length;i++)
            {
              let pageInfo = {}
              pageInfo.label = pageLayer[i].PageName
              pageInfo.value = pageLayer[i].PageId
              pageInfo.pageType = pageLayer[i].PageType
              pageInfo.pageModelUuid = pageLayer[i].modelId
              _t.displayPageList.push(pageInfo)
            }
          }
        }
      })
    },
    afterVisibleChange(val) {
      if(!val)
      {
        this.edit = false
      }
    },
    validateIPAddress  (rule, value, callback)  {
      let regexp = /^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$/;
      let valdata = value.split(',');
      let isCorrect = true;
      if (valdata.length) {
        for (let i = 0; i < valdata.length; i++) {
          if (regexp.test(valdata[i]) == false) {
            isCorrect = false;
          }
        }
      }

      if (value == '') {
        return callback(new Error('请输入iP地址'));
      } else if (!isCorrect) {
        callback(new Error('请输入正确对ip地址'));
      } else {
        callback()
      }
    },
    showDrawer(type) {
      if(type=='device')
      {
        this.isDevice=true
      }
      else
      {
        this.isDevice=false
      }
      this.visible = true;
    },
    onClose() {
      this.visible = false;
    },
    addDeviceOrZone(e){
      e.preventDefault()
      let _t = this
      this.form.validateFields((err) => {
        if (!err) {
          if((this.selectNode==null)||(this.selectNode.value.type==1))
          {
            this.$message.error(this.$t("monitor.SelectZoneTips"))
            return
          }
          let  params = {}
          if(this.isDevice)
          {
            params = {
              type:1,
              name:this.form.getFieldValue('name'),
              longitude:_t.longitude.toString(),
              latitude:_t.latitude.toString(),
              timeout:parseInt(this.form.getFieldValue('timeout')),
              interval:parseInt(this.form.getFieldValue('interval')),
              failedTimes:parseInt(this.form.getFieldValue('failedTimes')),
              muid:this.form.getFieldValue('model'),
              configUid:this.form.getFieldValue('configurationModel'),
              PageUUID:this.form.getFieldValue('configurationPageUUID'),
              deviceType:parseInt(this.form.getFieldValue('DeviceType')),
              description:this.form.getFieldValue('description'),
              offlineClear:parseInt(this.form.getFieldValue('offlineClear')),
              offlineDefaultValue:this.form.getFieldValue('offlineDefaultValue'),
            };
            if(params.deviceType==1)
            {
              params.extra = JSON.stringify({
                snmp:{
                  ipaddress:this.form.getFieldValue('agentIpaddress')
                }
              })
            }
            else if(params.deviceType==2)
            {
              if(typeof this.form.getFieldValue('RegisterPack')!="undefined")
              {
                if((parseInt(this.form.getFieldValue('RegisterPack'))<=0)||(parseInt(this.form.getFieldValue('RegisterPack'))>1000000))
                {
                  _t.$message.error(_t.$t("device.DeviceAddIDError"));
                  return
                }
              }
              params.extra = JSON.stringify({
                modbus:{
                  IPAddress:this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):null,
                  Port:this.form.getFieldValue('Port')?this.form.getFieldValue('Port'):null,
                  address:this.form.getFieldValue('DeviceAddress'),
                  packTime:this.form.getFieldValue('packTime')?parseInt(this.form.getFieldValue('packTime')):100,
                  RegisterPack:this.form.getFieldValue('RegisterPack')?parseInt(this.form.getFieldValue('RegisterPack')):-1
                }
              })
            }
            else if(params.deviceType==3)
            {
              params.extra = JSON.stringify({
                OpcuaExtraData:{
                  endpoint:this.form.getFieldValue('OPCUAEndPoint'),
                }
              })
            }
            else if(params.deviceType==15)
            {
              params.extra = JSON.stringify({
                SimS7:{
                  IpAddress:this.form.getFieldValue('IPAddress'),
                  Slot:this.form.getFieldValue('Slot'),
                  Rack:this.form.getFieldValue('Rack'),
                }
              })
            }
            else if(params.deviceType==20)
            {
              params.extra = JSON.stringify({
                Mqtt:{
                  ClientID:this.form.getFieldValue('ClientID'),
                }
              })
            }
            else if(params.deviceType==30)
            {
              params.extra = JSON.stringify({
                DLT645:{
                  IPAddress:this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):null,
                  Port:this.form.getFieldValue('Port')?this.form.getFieldValue('Port'):null,
                  RegisterPack:this.form.getFieldValue('RegisterPack')?parseInt(this.form.getFieldValue('RegisterPack')):-1,
                  ConnectAddress:this.form.getFieldValue('ConnectAddress'),
                  OperatorCode:this.form.getFieldValue('OperatorCode'),
                  Password:this.form.getFieldValue('Password'),
                  BeforeCode:this.form.getFieldValue('BeforeCode'),
                  packTime:this.form.getFieldValue('packTime')?parseInt(this.form.getFieldValue('packTime')):100,
                }
              })
            }
            else if(params.deviceType==40)
            {
              params.extra = JSON.stringify({
                IEC104:{
                  IPAddress:this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):"127.0.0.1",
                  Port:this.form.getFieldValue('Port')?this.form.getFieldValue('Port'):2404,
                  DeviceAddress:this.form.getFieldValue('DeviceAddress')?this.form.getFieldValue('DeviceAddress'):1,
                }
              })
            }
            else if(params.deviceType==350)
            {
              params.extra = JSON.stringify({
                IEC61850:{
                  IPAddress:this.form.getFieldValue('IEC61850EndPoint')?this.form.getFieldValue('IEC61850EndPoint'):"127.0.0.1",
                  Port:this.form.getFieldValue('IEC61850Port')?this.form.getFieldValue('IEC61850Port'):102,
                }
              })
            }
            else if(params.deviceType==470)
            {
              params.extra = JSON.stringify({
                HJ212:{
                  DeviceSN:this.form.getFieldValue('HJ212DeviceSN')?this.form.getFieldValue('HJ212DeviceSN'):"",
                  PW:this.form.getFieldValue('HJ212PW')?this.form.getFieldValue('HJ212PW'):"",
                }
              })
            }
            else if(params.deviceType==490)
            {
              params.extra = JSON.stringify({
                CJT188:{
                  IPAddress:this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):null,
                  Port:this.form.getFieldValue('Port')?this.form.getFieldValue('Port'):null,
                  RegisterPack:this.form.getFieldValue('RegisterPack')?parseInt(this.form.getFieldValue('RegisterPack')):-1,
                  ConnectAddress:this.form.getFieldValue('ConnectAddress'),
                }
              })
            }
            else if(params.deviceType==500)
            {
              params.extra = JSON.stringify({
                BACnet:{
                  IPAddress:this.form.getFieldValue('BACnetIPAddress')?this.form.getFieldValue('BACnetIPAddress'):"127.0.0.1",
                  Port:this.form.getFieldValue('BACnetPort')?parseInt(this.form.getFieldValue('BACnetPort')):48089,
                }
              })
            }
          }
          else  if(this.isDevice==0)
          {
            params = {
              type:0,
              name:this.form.getFieldValue('name'),
              description:this.form.getFieldValue('description'),
              configUid:this.form.getFieldValue('configurationModel'),
              PageUUID:this.form.getFieldValue('configurationPageUUID'),
            };
          }
          if(this.$refs.deviceTree.checkIsEmpty()||this.selectNode==null)
          {
            params.pid=0
          }
          else
          {
            params.pid=this.selectNode.value.sid
          }
          this.$message.loading({ content: 'Waiting...',loadingKey,duration: 0 });
          addMonitor(params).then(function (res){
            if(res.data.code==0)
            {
              _t.$refs.deviceTree.getMonitorTree()
              _t.getMonitorTree()
              _t.visible = false;
              _t.$message.success(_t.$t("device.DeviceAddSuccess"));
            }
            else if(res.data.code==3001)
            {
              _t.$message.error(_t.$t("device.DeviceExist"));
            }
            else if(res.data.code==3003)
            {
              _t.$message.error(_t.$t("device.DeviceAddressExist"));
            }
            else
            {
              _t.$message.error(_t.$t("device.DeviceAddFailed"));
            }
            setTimeout(function (){
              _t.$message.destroy()
            },500)
          }).catch(function (){
            _t.$message.error(_t.$t("device.DeviceAddFailed"));
            setTimeout(function (){
              _t.$message.destroy()
            },500)
          })
        }
      })
    },
    TestConnect(e){
      let _t = this

      let  params = {
        deviceType:e=="add"?parseInt(this.form.getFieldValue('DeviceType')):parseInt(this.editForm.getFieldValue('DeviceType'))
      }
      if (e=="add")
      {
        if(params.deviceType==1)
        {
          params.IP = this.form.getFieldValue('agentIpaddress')?this.form.getFieldValue('agentIpaddress'):null
        }
        else if(params.deviceType==2)
        {
          params.IP = this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==3)
        {
          params.IP = this.form.getFieldValue('OPCUAEndPoint')?this.form.getFieldValue('OPCUAEndPoint'):null
        }
        else if(params.deviceType==15)
        {
          params.IP = this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==30)
        {
          params.IP = this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==40)
        {
          params.IP = this.form.getFieldValue('IPAddress')?this.form.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==350)
        {
          params.IP = this.form.getFieldValue('IEC61850EndPoint')?this.form.getFieldValue('IEC61850EndPoint'):null
        }
        else
        {
          params.IP = null
        }
      }
      else
      {
        if(params.deviceType==1)
        {
          params.IP = this.editForm.getFieldValue('agentIpaddress')?this.editForm.getFieldValue('agentIpaddress'):null
        }
        else if(params.deviceType==2)
        {
          params.IP = this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==3)
        {
          params.IP = this.editForm.getFieldValue('OPCUAEndPoint')?this.editForm.getFieldValue('OPCUAEndPoint'):null
        }
        else if(params.deviceType==15)
        {
          params.IP = this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==30)
        {
          params.IP = this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==40)
        {
          params.IP = this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null
        }
        else if(params.deviceType==350)
        {
          params.IP = this.editForm.getFieldValue('IEC61850EndPoint')?this.editForm.getFieldValue('IEC61850EndPoint'):null
        }
        else
        {
          params.IP = null
        }
      }

      if(params.IP==""||params.IP==null)
      {
        _t.$message.error(_t.$t("device.DeviceNotIp"));
        return
      }
      this.isTesting=true
      TestConnect(params).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t("device.DeviceOnline"));
        }
        else
        {
          _t.$message.error(_t.$t("device.DeviceOffline"));
        }
        _t.isTesting=false
      }).catch(function (){
        _t.isTesting=false
        _t.$message.error(_t.$t("device.DeviceOffline"));
      })
    },
    editDeviceOrZone(e){
      e.preventDefault()

      this.editForm.validateFields((err) => {
        if (!err) {
          let  params = {}
          let _t = this
          if(this.isDevice)
          {
            params = {
              uuid:this.editUuid,
              editData: {
                type: 1,
				        pid:this.selectNodeSid,
                longitude:_t.longitude.toString(),
                latitude:_t.latitude.toString(),
                name: this.editForm.getFieldValue('name'),
                muid: this.editForm.getFieldValue('model'),
                configUid:this.editForm.getFieldValue('configurationModel'),
                PageUUID:this.editForm.getFieldValue('configurationPageUUID'),
                timeout:parseInt(this.editForm.getFieldValue('timeout')),
                interval:parseInt(this.editForm.getFieldValue('interval')),
                failedTimes:parseInt(this.editForm.getFieldValue('failedTimes')),
                deviceType: parseInt(this.editForm.getFieldValue('DeviceType')),
                description: this.editForm.getFieldValue('description'),
                offlineClear:parseInt(this.editForm.getFieldValue('offlineClear')),
                offlineDefaultValue:this.editForm.getFieldValue('offlineDefaultValue'),
              }
            };

            if(params.editData.deviceType==1)
            {
              params.editData.extra = JSON.stringify({
                snmp:{
                  ipaddress:this.editForm.getFieldValue('agentIpaddress')
                }
              })
            }
            else if(params.editData.deviceType==2)
            {
              if(typeof this.editForm.getFieldValue('RegisterPack')!="undefined")
              {
                if((parseInt(this.editForm.getFieldValue('RegisterPack'))<=0)||(parseInt(this.editForm.getFieldValue('RegisterPack'))>1000000))
                {
                  _t.$message.error(_t.$t("device.DeviceAddIDError"));
                  return
                }
              }
              params.editData.extra = JSON.stringify({
                modbus:{
                  IPAddress:this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null,
                  Port:this.editForm.getFieldValue('Port')?this.editForm.getFieldValue('Port'):null,
                  address:this.editForm.getFieldValue('DeviceAddress'),
                  packTime:this.editForm.getFieldValue('packTime')?parseInt(this.editForm.getFieldValue('packTime')):100,
                  RegisterPack:this.editForm.getFieldValue('RegisterPack')?parseInt(this.editForm.getFieldValue('RegisterPack')):-1
                }
              })
            }
            else if(params.editData.deviceType==3)
            {
              params.editData.extra = JSON.stringify({
                OpcuaExtraData:{
                  endpoint:this.editForm.getFieldValue('OPCUAEndPoint'),
                }
              })
            }
            else if(params.editData.deviceType==15)
            {
              params.editData.extra = JSON.stringify({
                SimS7:{
                  IpAddress:this.editForm.getFieldValue('IPAddress'),
                  Slot:this.editForm.getFieldValue('Slot'),
                  Rack:this.editForm.getFieldValue('Rack'),
                }
              })
            }
            else if(params.editData.deviceType==20)
            {
              params.editData.extra = JSON.stringify({
                Mqtt:{
                  ClientID:this.editForm.getFieldValue('ClientID'),
                }
              })
            }
            else if(params.editData.deviceType==30)
            {
              params.editData.extra = JSON.stringify({
                DLT645:{
                  IPAddress:this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null,
                  Port:this.editForm.getFieldValue('Port')?this.editForm.getFieldValue('Port'):null,
                  RegisterPack:this.editForm.getFieldValue('RegisterPack')?parseInt(this.editForm.getFieldValue('RegisterPack')):-1,
                  ConnectAddress:this.editForm.getFieldValue('ConnectAddress'),
                  OperatorCode:this.editForm.getFieldValue('OperatorCode'),
                  Password:this.editForm.getFieldValue('Password'),
                  BeforeCode:this.editForm.getFieldValue('BeforeCode'),
                  packTime:this.editForm.getFieldValue('packTime')?parseInt(this.editForm.getFieldValue('packTime')):100,
                }
              })
            }
            else if(params.editData.deviceType==40)
            {
              params.editData.extra = JSON.stringify({
                IEC104:{
                  IPAddress:this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null,
                  Port:this.editForm.getFieldValue('Port')?this.editForm.getFieldValue('Port'):null,
                  DeviceAddress:this.editForm.getFieldValue('DeviceAddress')?this.editForm.getFieldValue('DeviceAddress'):1,
                }
              })
            }
            else if(params.editData.deviceType==350)
            {
              params.editData.extra = JSON.stringify({
                IEC61850:{
                  IPAddress:this.editForm.getFieldValue('IEC61850EndPoint')?this.editForm.getFieldValue('IEC61850EndPoint'):'127.0.0.1',
                  Port:this.editForm.getFieldValue('IEC61850Port')?this.editForm.getFieldValue('IEC61850Port'):102
                }
              })
            }
            else if(params.editData.deviceType==470)
            {
              params.editData.extra = JSON.stringify({
                HJ212:{
                  DeviceSN:this.editForm.getFieldValue('HJ212DeviceSN')?this.editForm.getFieldValue('HJ212DeviceSN'):'',
                  PW:this.editForm.getFieldValue('HJ212PW')?this.editForm.getFieldValue('HJ212PW'):""
                }
              })
            }
            else if(params.editData.deviceType==490)
            {
              params.editData.extra = JSON.stringify({
                CJT188:{
                  IPAddress:this.editForm.getFieldValue('IPAddress')?this.editForm.getFieldValue('IPAddress'):null,
                  Port:this.editForm.getFieldValue('Port')?this.editForm.getFieldValue('Port'):null,
                  RegisterPack:this.editForm.getFieldValue('RegisterPack')?parseInt(this.editForm.getFieldValue('RegisterPack')):-1,
                  ConnectAddress:this.editForm.getFieldValue('ConnectAddress'),
                }
              })
            }
            else if(params.editData.deviceType==500)
            {
              params.editData.extra = JSON.stringify({
                BACnet:{
                  IPAddress:this.editForm.getFieldValue('BACnetIPAddress')?this.editForm.getFieldValue('BACnetIPAddress'):"127.0.0.1",
                  Port:this.editForm.getFieldValue('BACnetPort')?parseInt(this.editForm.getFieldValue('BACnetPort')):48089,
                }
              })
            }
          }
          else
          {
            params = {
              uuid:this.editUuid,
              editData: {
                type: 0,
                name: this.editForm.getFieldValue('name'),
                description: this.editForm.getFieldValue('description'),
                configUid:this.editForm.getFieldValue('configurationModel'),
                PageUUID:this.editForm.getFieldValue('configurationPageUUID'),
              }
            };
          }
          this.$message.loading({ content: 'Waiting...',loadingKey,duration: 0 });
          editMonitor(params).then(function (res){
            _t.$message.destroy()
            if(res.data.code==0)
            {
              _t.$refs.deviceTree.getMonitorTree()
              _t.getMonitorTree()
              _t.visible = false;
              _t.$message.success(_t.$t("monitor.EditSuccess"))
            }
            else
            {
              _t.$message.error(_t.$t("monitor.EditFailed"))
            }
            setTimeout(function (){
              _t.$message.destroy()
            },1000)

          })
        }
      })
    },
    onSelect(selectData) {
      const info = selectData.info
      const _t = this

      this.selectNode = info
      if(this.selectNode.value.type==1)
      {
        this.isDevice = true
        this.goToEdit(this.selectNode.value)
      }
      else
      {
        this.editIsDevice=false
        this.isDevice = false
        this.selectKey = selectData.key
        _t.tableDataSource=selectData.tableList
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
            tempData.uuid = res.data.list[i].uuid
            tempData.modelName = res.data.list[i].name
            //modbus的连接类型
            tempData.modbusConnectType = res.data.list[i].modbusConnectType
            tempData.DLT645ConnectType = res.data.list[i].DLT645ConnectType
            tempData.CJT188ConnectType = res.data.list[i].CJT188ConnectType
            _t.modelList.push(tempData)
            tempData={}
          }
        }
      })
    },
    goToEdit(data) {
      this.editUuid = data.uuid
      this.editType = data.type
      this.editIsDevice=true
      this.selectNodeKey = data.uuid
      this.selectNodeSid = data.pid
      if(data.configUid!="")
      {
        this.GetDisplayPage(data.configUid)
      }

      if(data.type==1)
      {
        this.isDevice = true
        let _t = this
        this.DeviceType = data.deviceType
        let extra = {}
        try{
           extra = JSON.parse(data.extra)
        }catch (e) {
          extra = {}
        }
        _t.longitude=data.longitude
        _t.latitude=data.latitude

        this.getModelList()
        this.$message
            .loading(this.$t("monitor.loading"), 0.5)

        setTimeout(function (){
          _t.editForm.setFieldsValue(
              {
                name:data.name,
                DeviceType:data.deviceType,
                model:data.muid,
                configurationModel:data.configUid,
                configurationPageUUID:data.PageUUID,
                description:data.description,
                timeout:data.timeout.toString(),
                interval:data.interval.toString(),
                failedTimes:data.failedTimes.toString(),
                offlineClear:data.offlineClear.toString(),
                offlineDefaultValue:data.offlineDefaultValue,
              })
            if(data.deviceType==1)
            {
              _t.editForm.setFieldsValue(
                  {
                    agentIpaddress:extra.snmp.ipaddress,
                  })
            }
            else if(data.deviceType==2){
              _t.editForm.setFieldsValue(
                  {
                    DeviceAddress:extra.modbus.address,
                    packTime:extra.modbus.packTime?extra.modbus.packTime:100,
                  })
              if(typeof extra.modbus.RegisterPack!='undefined'&&extra.modbus.RegisterPack!=-1)
              {
                _t.modbusConnectType="TCPServer"
                _t.$nextTick(function (){
                  _t.RegisterPackByte = _t.Value2Bytes(extra.modbus.RegisterPack)
                  _t.editForm.setFieldsValue(
                      {
                        RegisterPack:extra.modbus.RegisterPack,
                      })
                })
              }
              else if(extra.modbus.IPAddress!=null&&extra.modbus.Port!=null)
              {
                _t.modbusConnectType="TCPClient"
                _t.$nextTick(function (){
                  _t.editForm.setFieldsValue(
                      {
                        IPAddress:extra.modbus.IPAddress,
                        Port:extra.modbus.Port,
                      })
                })
              }
              else
              {
                _t.modbusConnectType=""
              }
            }
            else if(data.deviceType==3)
            {
            _t.editForm.setFieldsValue(
                {
                  OPCUAEndPoint:extra.OpcuaExtraData.endpoint,
                })
          }
            else if(data.deviceType==5)
            {
              _t.editForm.setFieldsValue(
                  {
                    DeviceFlag:data.uuid,
                  })
            }
            else if(data.deviceType==15)
            {
              _t.editForm.setFieldsValue(
                  {
                    IPAddress:extra.SimS7.IpAddress,
                    Slot:extra.SimS7.Slot,
                    Rack:extra.SimS7.Rack
                  })
            }
            else if(data.deviceType==20)
            {
              _t.editForm.setFieldsValue(
                  {
                    ClientID:extra.Mqtt.ClientID
                  })
            }
            else if(data.deviceType==30)
            {
              _t.editForm.setFieldsValue(
                  {
                    ConnectAddress:extra.DLT645.ConnectAddress,
                    OperatorCode:extra.DLT645.OperatorCode,
                    Password:extra.DLT645.Password,
                    BeforeCode:extra.DLT645.BeforeCode,
                    packTime:extra.DLT645.packTime?extra.DLT645.packTime:100,
                  })
              if(typeof extra.DLT645.RegisterPack!='undefined'&&extra.DLT645.RegisterPack!=-1)
              {
                _t.DLT645ConnectType="TCPServer"
                _t.$nextTick(function (){
                  _t.RegisterPackByte = _t.Value2Bytes(extra.DLT645.RegisterPack)
                  _t.editForm.setFieldsValue(
                      {
                        RegisterPack:extra.DLT645.RegisterPack,
                      })
                })
              }
              else if(extra.DLT645.IPAddress!=null&&extra.DLT645.Port!=null)
              {
                _t.DLT645ConnectType="TCPClient"
                _t.$nextTick(function (){
                  _t.editForm.setFieldsValue(
                      {
                        IPAddress:extra.DLT645.IPAddress,
                        Port:extra.DLT645.Port,
                      })
                })
              }
              else
              {
                _t.DLT645ConnectType=""
              }

            }
            else if(data.deviceType==40)
            {
              _t.editForm.setFieldsValue(
                  {
                    IPAddress:extra.IEC104.IPAddress,
                    Port:extra.IEC104.Port,
                    DeviceAddress:extra.IEC104.DeviceAddress,
                  })
            }
            else if(data.deviceType==350)
            {
              _t.editForm.setFieldsValue(
                  {
                    IEC61850EndPoint:extra.IEC61850.IPAddress,
                    IEC61850Port:extra.IEC61850.Port
                  })
            }
            else if(data.deviceType==470)
            {
              _t.editForm.setFieldsValue(
                  {
                    HJ212DeviceSN:extra.HJ212.DeviceSN,
                    HJ212PW:extra.HJ212.PW
                  })
            }
            else if(data.deviceType==490)
            {
              _t.editForm.setFieldsValue(
                  {
                    ConnectAddress:extra.CJT188.ConnectAddress,
                  })

              if(typeof extra.CJT188.RegisterPack!='undefined'&&extra.CJT188.RegisterPack!=-1)
              {
                _t.CJT188ConnectType="TCPServer"
                _t.$nextTick(function (){
                  _t.RegisterPackByte = _t.Value2Bytes(extra.CJT188.RegisterPack)
                  _t.editForm.setFieldsValue(
                      {
                        RegisterPack:extra.CJT188.RegisterPack,
                      })
                })
              }
              else if(extra.CJT188.IPAddress!=null&&extra.CJT188.Port!=null)
              {
                _t.CJT188ConnectType="TCPClient"
                _t.$nextTick(function (){
                  _t.editForm.setFieldsValue(
                      {
                        IPAddress:extra.CJT188.IPAddress,
                        Port:extra.CJT188.Port,
                      })
                })
              }
              else
              {
                _t.CJT188ConnectType=""
              }
            }
            else if(data.deviceType==500)
            {
              _t.editForm.setFieldsValue(
                  {
                    BACnetIPAddress:extra.BACnet.IPAddress,
                    BACnetPort:extra.BACnet.Port
                  })
            }
        },500)
      }
      else
      {
        this.isDevice = false
        let _t = this
        this.$message
            .loading(this.$t("monitor.loading"), 0.5)

        setTimeout(function (){
          _t.editForm.setFieldsValue(
              {
                name:data.name,
                configurationModel:data.configUid,
                configurationPageUUID:data.PageUUID,
                description:data.description
              })
        },500)
      }
    },
    goToCopy(data) {
      this.CopyUuid = data.uuid
      this.CopyType = data.type
      this.CopyDeviceType=data.deviceType
      this.CopyVisible = true
    },
    changeDeviceType(value){
      this.DeviceType = value
      this.getModelList()
    },
    changeModelType(uuid){
      for(let i=0;i<this.modelList.length;i++)
      {
        if(this.modelList[i].uuid==uuid)
        {
          this.modbusConnectType = this.modelList[i].modbusConnectType
          this.DLT645ConnectType = this.modelList[i].DLT645ConnectType
          this.CJT188ConnectType = this.modelList[i].CJT188ConnectType
          return
        }
      }
      this.CJT188ConnectType=""
      this.DLT645ConnectType=""
      this.modbusConnectType=""
    },
    deleteRecord(key,type){
      if(type!=1)
      {
        let havedDevice = this.$refs.deviceTree.checkZoneHavedDevice(key)
        if(havedDevice)
        {
          this.$message.error(this.$t("monitor.ZoneHavedDevice"))
          return
        }
      }
      let params={
        uuid:key
      }
      let _t = this
      this.$message.loading({ content: 'Waiting...',loadingKey,duration: 0 });
      delMonitor(params).then(function (res) {
        if(res.data.code==0)
        {
          _t.tableDataSource = _t.tableDataSource.filter(item => item.key !== key)
          _t.selectedRows = _t.selectedRows.filter(item => item.key !== key)
          _t.$refs.deviceTree.getMonitorTree()
        }
        _t.$message.destroy()
      })
    },
    backZoneList(){
      this.editIsDevice=false
      this.editUuid=""
    },
    updateTable(data){
        this.tableDataSource = data
    },
    SetDeviceEnable(e,editUuid,deviceType,name){
      let _t = this
      let params = {
        uuid:editUuid,
        editData: {
          IsEnable:e?1:0,
          name:name,
          deviceType:deviceType
        }
      };
      this.$message.loading({ content: 'Waiting...',loadingKey,duration: 0 });
      SetDeviceStartOrStop(params).then(function (res){
        _t.$message.destroy()
        if(res.data.code==0)
        {
          _t.$refs.deviceTree.getMonitorTree()
          _t.visible = false;
          _t.$message.success(_t.$t("monitor.EditSuccess"))
        }
        else
        {
          _t.$message.error(_t.$t("monitor.EditFailed"))
        }
        setTimeout(function (){
          _t.$message.destroy()
        },1000)

      })
    },
  },
};
</script>

<style scoped>
#allmap{
  height: 450px;
  width: 100%;
  margin: 10px 0;
}
/* 穿透作用域修改选择列宽度 */
::v-deep .ant-table-selection-column {
  min-width: 33px;  /* 最小宽度 */
  max-width: 33px;  /* 最大宽度 */
  width: 3%;        /* 百分比宽度（根据需求选择单位） */
}

/* 可选：调整复选框与边界的间距 */
::v-deep .ant-table-selection-column .ant-checkbox-wrapper {
  padding-left: 8px;
}
::v-deep #components-layout-demo-side .logo {
  height: 32px;
  background: rgba(255, 255, 255, 0.2);
  margin: 16px;
}
::v-deep .ant-card-body {
  padding: 10px;
}
::v-deep .search{
  margin-bottom: 54px;
}
::v-deep .fold{
  width: calc(100% - 216px);
  display: inline-block
}
::v-deep .operator{
  margin-bottom: 18px;
}
::v-deep .spin-content {
  border: 1px solid #91d5ff;
  background-color: #e6f7ff;
  padding: 30px;
}
::v-deep @media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}
::v-deep .ant-form-item {
  margin-bottom: 5px;
}
::v-deep .ant-row .ant-form-item {
  margin-bottom: 5px;
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

::v-deep .ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  /*background: #f8f8f8;*/
  /*border-bottom: 1px solid #e8e8e8;*/
  transition: background .3s ease;
}

</style>
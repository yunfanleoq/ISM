<template>
  <div style=" overflow: hidden;
            overflow-y: auto;">
    <template v-if="configObject != null&&isLayer == false">
      <a-tabs default-active-key="1" :tab-position="tabPosition" style="width:420px;max-width:600px;border-color: #95B8E7;">
        <a-tab-pane key="1" :tab="$t('displayConfig.Properties.TabHeaterStyle')" style="padding:5px;">
          <a-form  :label-col="{ span: 7}" :wrapper-col="{ span: 16 }">
            <div>
              <a-form-item :label="$t('displayConfig.Properties.ComponentName')">
                <a-input v-model="configObject.name" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentIsVisible')">
                <a-select
                    v-model="configObject.style.visible" @change="chargeGroupIDVisible"
                >
                  <a-select-option v-for="options in [{label:'True',value:1},{label:'False',value:0}]" :key="options.value" :value="options.value">
                    {{ options.label}}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </div>

            <div>
              <a-form-item :label="$t('displayConfig.Properties.ComponentX')">
                <a-input
                    @change="UpdateNode"
                    type="number"
                    suffix="px"
                    v-model="selectedNodePops.position.x"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentY')">
                <a-input
                    @change="UpdateNode"
                    type="number"
                    suffix="px"
                    v-model="selectedNodePops.position.y"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentWith')">
                <a-input
                    @change="UpdateNode"
                    type="number"
                    suffix="px"
                    v-model="selectedNodePops.size.width"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentHeight')">
                <a-input
                    @change="UpdateNode"
                    type="number"
                    suffix="px"
                    v-model="selectedNodePops.size.height"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentAnimate')">
                <a-select v-model="configObject.style.animate"  :allowClear="true" @change="UpdateNodeData">
                  <a-select-option :key="index" :value="item.name" v-for="(item, index) in animates">{{item.alias}}</a-select-option>
                </a-select>
              </a-form-item>
            </div>

            <div>
              <a-form-item :label="$t('displayConfig.Properties.BorderWidth')">
                <a-input  @change="UpdateNodeData"
                         type="number"
                         suffix="px"
                         v-model="configObject.style.borderWidth"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.BorderStyle')">
                <a-select @change="UpdateNodeData"
                    v-model="configObject.style.borderStyle"
                >

                  <a-select-option v-for="options in borderStyleOptions" :key="options" :value="options">
                    {{ options }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.BorderColor')">
                <a-input @change="UpdateNodeData"
                    type="color"
                    v-model="configObject.style.borderColor"
                >
                  <a-tooltip slot="suffix" :title="$t('displayConfig.Properties.ClearColor')" @click="configObject.style.borderColor='transparent';UpdateNodeData()">
                    <a-icon type="delete" style="color: rgba(0,0,0,.45)" />
                  </a-tooltip>
                </a-input>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.BorderEdges')">
                <a-input @change="UpdateNodeData"
                    type="number"
                    suffix="px"
                    min="1"
                    v-model="configObject.style.BorderEdges"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.opacity')">
                <a-input @change="UpdateNodeData"
                    type="number"
                    :min="0" :max="1" :step="0.1"
                    v-model="configObject.style.opacity"
                />
              </a-form-item>
            </div>

            <div>
              <a-form-item :label="$t('displayConfig.Properties.zIndex')">
                <a-input type="number" v-model="configObject.style.zIndex" @change="UpdateNode"/>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.Transform')" v-if="typeof selectedNodePops.angle!='undefined'">
                <a-input @change="UpdateNode"
                    type="number"
                    v-model="selectedNodePops.angle"
                    suffix="deg"
                />
              </a-form-item>

              <a-form-item :label="$t('displayConfig.Properties.ComponentBackColor')"  v-if="typeof configObject.style.backColor!='undefined'">
                <a-input @change="UpdateNodeData"
                    type="color"
                    v-model="configObject.style.backColor"
                >
                  <a-tooltip slot="suffix" :title="$t('displayConfig.Properties.ClearColor')" @click="configObject.style.backColor='transparent';UpdateNodeData()">
                    <a-icon type="delete" style="color: rgba(0,0,0,.45)" />
                  </a-tooltip>
                </a-input>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentForeColor')" v-if="typeof configObject.style.foreColor!='undefined'">
                <a-input @change="UpdateNodeData"
                    type="color"
                    v-model="configObject.style.foreColor"
                >
                  <a-tooltip slot="suffix" :title="$t('displayConfig.Properties.ClearColor')" @click="configObject.style.foreColor='transparent';UpdateNodeData()">
                    <a-icon type="delete" style="color: rgba(0,0,0,.45)" />
                  </a-tooltip>
                </a-input>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentUrl')" v-if="configObject.style.url != undefined && configObject.style.url != null">
                <a-input v-model="configObject.style.url" @change="UpdateNodeData"/>
              </a-form-item>

              <a-form-item :label="$t('displayConfig.Properties.ComponentText')" v-if="selectedNodePops.data.detail.style.text != undefined && selectedNodePops.data.detail.style.text != null">
                <a-textarea  v-model="selectedNodePops.data.detail.style.text" @change="UpdateNodeData"/>
              </a-form-item>

              <a-form-item :label="$t('displayConfig.Properties.ComponentTextAlign')" v-if="configObject.style.textAlign != undefined && configObject.style.textAlign != null">
                <a-select @change="UpdateNodeData"
                    v-model="configObject.style.textAlign"
                >
                  <a-select-option v-for="options in textAlignOptions" :key="options" :value="options">
                    {{ options }}
                  </a-select-option>
                </a-select>

              </a-form-item>

              <a-form-item :label="$t('displayConfig.Properties.ComponentFontFamily')" v-if="configObject.style.fontFamily != undefined && configObject.style.fontFamily != null">
                <a-select @change="UpdateNodeData"
                    v-model="configObject.style.fontFamily"
                >
                  <a-select-option v-for="options in fontFamilyOptions" :key="options" :value="options">
                    {{ options }}
                  </a-select-option>
                </a-select>

              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentFontBold')" v-if="configObject.style.fontWeight != undefined && configObject.style.fontWeight != null">
                <a-input @change="UpdateNodeData"
                    type="number"
                    min="100"
                    max="900"
                    v-model="configObject.style.fontWeight"
                />
              </a-form-item>


              <a-form-item :label="$t('displayConfig.Properties.ComponentLetterSpacing')" v-if="configObject.style.letterSpacing != undefined && configObject.style.letterSpacing != null">
                <a-input @change="UpdateNodeData"
                    type="number"
                    min="0"
                    max="100"
                    v-model="configObject.style.letterSpacing"
                />
              </a-form-item>

              <a-form-item :label="$t('displayConfig.Properties.ComponentFontItalic')" v-if="configObject.style.italic != undefined && configObject.style.italic != null">

                <a-select @change="UpdateNodeData"
                    v-model="configObject.style.italic"
                >
                  <a-select-option v-for="options in [{label:'True',value:1},{label:'False',value:0}]" :key="options.value" :value="options.value">
                    {{ options.label}}
                  </a-select-option>
                </a-select>
              </a-form-item>


              <a-form-item :label="$t('displayConfig.Properties.ComponentFontSize')" v-if="configObject.style.fontSize != undefined && configObject.style.fontSize != null">
                <a-input @change="UpdateNodeData"
                    type="number"
                    suffix="px"
                    v-model="configObject.style.fontSize"
                />
              </a-form-item>

              <a-form-item :label="$t('displayConfig.Properties.ComponentRadius')" v-if="configObject.style.radius != undefined && configObject.style.radius != null">
                <a-input @change="UpdateNodeData"
                    type="number"
                    v-model="configObject.style.radius"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentLineWidth')" v-if="configObject.style.lineWidth != undefined && configObject.style.lineWidth != null">
                <a-input @change="UpdateNodeData"
                    type="number"
                    suffix="px"
                    v-model="configObject.style.lineWidth"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentLineDash')" v-if="configObject.style.setLineDash != undefined && configObject.style.setLineDash != null">
                <a-input @change="UpdateNodeData" type="text" v-model="configObject.style.setLineDash" />
              </a-form-item>
            </div>

            <!--          自定义数据-->
            <template v-for="(diyData,diyIndex) in configObject.style.diy" >
              <div :key="diyIndex">
                <a-form-item :label="$t(diyData.name)"  v-if="diyData.type!=15">
                  <!--数字类型-->
                  <a-input @change="UpdateNodeData" type="number"  v-model="configObject.style.diy[diyIndex].value" :min="diyData.min" :max="diyData.max" :value="diyData.value" v-if="diyData.type==1"/>
                  <!--颜色类型-->
                  <a-input @change="UpdateNodeData" type="color"  v-model="configObject.style.diy[diyIndex].value" :value="diyData.value" v-if="diyData.type==2">
                    <a-tooltip slot="suffix" :title="$t('displayConfig.Properties.ClearColor')" @click="configObject.style.diy[diyIndex].value='transparent';UpdateNodeData()">
                      <a-icon type="delete" style="color: rgba(0,0,0,.45)" />
                    </a-tooltip>
                  </a-input>

                  <!--字体选择-->
                  <a-select @change="UpdateNodeData" v-model="configObject.style.diy[diyIndex].value" v-if="diyData.type==3">
                    <a-select-option v-for="options in fontFamilyOptions" :key="options" :value="options">
                      {{ options }}
                    </a-select-option>
                  </a-select>
                  <!--字符串-->
                  <a-input type="text" @change="UpdateNodeData"  v-model="configObject.style.diy[diyIndex].value" :value="diyData.value" v-if="diyData.type==4"/>
                  <!--图片类型-->
                  <div v-if="diyData.type==5">
                    <vue-hover-mask>
                      <!-- 默认插槽 -->
                      <img v-if="configObject.style.diy[diyIndex].value!=''"  style="width: 200px;height:200px;cursor: pointer" :src="configObject.style.diy[diyIndex].value" />
                      <div v-else :style="{width: '200px',height:'200px',cursor: 'pointer','background-color':'#F2F2F2'}"></div>
                      <!-- action插槽 -->
                      <template v-slot:action>
                        <span style="font-size: 14px" @click="showSystemImageModel(1,diyIndex)">{{$t('component.systemImageModel.selectImage')}}</span>
                        <a-divider type="vertical" />
                        <span  style="font-size: 14px" @click="configObject.style.diy[diyIndex].value=''">{{$t('component.systemImageModel.delImage')}}</span>
                      </template>
                    </vue-hover-mask>
                  </div>
                  <!--枚举类型-->
                  <a-select @change="UpdateNodeData" v-model="configObject.style.diy[diyIndex].value" v-if="diyData.type==6">
                    <a-select-option v-for="options in diyData.enumList" :key="options.value" :value="options.value">
                      {{ $t(options.option) }}
                    </a-select-option>
                  </a-select>
                  <!--浮点型-->
                  <a-input @change="UpdateNodeData" type="number"  v-model="configObject.style.diy[diyIndex].value" :step="0.1" :min="diyData.min" :max="diyData.max" :value="diyData.value" v-if="diyData.type==7"/>
                  <!--视频-->
                  <div v-if="diyData.type==8">
                    <a-input-search @change="UpdateNodeData"
                        v-model="configObject.style.diy[diyIndex].value.name"
                        :enter-button="$t('configComponent.video.SelectVideo')"
                        @search="showSystemVideoModel(1,diyIndex)"
                    />

                  </div>
                  <!--文本区域-->
                  <a-textarea  @change="UpdateNodeData" v-model="configObject.style.diy[diyIndex].value" :value="diyData.value" v-if="diyData.type==9" :rows="diyData.rows"/>
                  <!--3D模型-->
                  <div v-if="diyData.type==10">
                    <a-tooltip>
                      <template slot="title">
                        {{$t('configComponent.video.Model3DList')}}
                      </template>
                      <a-input-search @change="UpdateNodeData"
                          v-model="configObject.style.diy[diyIndex].value"
                          :enter-button="$t('configComponent.video.SelectVideo')"
                          @search="showSystemImageModel(1,diyIndex,1)"
                      />
                    </a-tooltip>
                  </div>
                  <!--code样式编辑-->
                  <div  v-if="diyData.type==11">
                    <a-button type="primary" icon="edit"  @click="CodeDbClick(diyIndex)">
                      {{$t('displayConfig.Properties.EditScript')}}
                    </a-button>
                  </div>

                  <!--文档类型-->
                  <div v-if="diyData.type==12">
                    <a-tooltip>
                      <template slot="title">
                        {{$t('configComponent.video.Model3DList')}}
                      </template>
                      <a-input-search @change="UpdateNodeData"
                          v-model="configObject.style.diy[diyIndex].value"
                          :enter-button="$t('configComponent.video.SelectVideo')"
                          @search="showSystemImageModel(1,diyIndex,2)"
                      />
                    </a-tooltip>
                  </div>
<!--                  画面选择-->
                  <div v-if="diyData.type==13">

                      <a-select @change="UpdateNodeData"
                          v-model="configObject.style.diy[diyIndex].Appid"
                          allowClear
                      >
                        <a-select-option v-for="options in configurationModel" :key="options.uuid" :value="options.uuid">
                          {{ options.name}}
                        </a-select-option>
                      </a-select>

                      <a-select @change="UpdateNodeData"
                          v-model="configObject.style.diy[diyIndex].value"
                      >
                        <a-select-option v-for="options in generateTargetPage(configObject.style.diy[diyIndex].Appid)" :key="options.value" :value="options.value">
                          {{ options.label}}
                        </a-select-option>
                      </a-select>

                  </div>
                  <!--JSON样式编辑-->
                  <div  v-if="diyData.type==14">
                    <a-button type="primary" icon="edit"  @click="CodeDbClick(diyIndex)">
                      {{$t('displayConfig.Properties.EditJSON')}}
                    </a-button>
                  </div>
                  <!--SQL样式编辑-->
                  <div  v-if="diyData.type==30">
                    <a-button type="primary" icon="edit"  @click="CodeDbClick(diyIndex)">
                      {{$t('displayConfig.Properties.EditSQL')}}
                    </a-button>
                  </div>
                  <!--脚本列表-->
                  <div  v-if="diyData.type==31">
                      <a-select @change="UpdateNodeData"
                          :allowClear="true"
                          v-model="configObject.style.diy[diyIndex].value"
                      >
                        <a-select-option v-for="options in ScriptDataSource" :key="options.ScriptUuid" :value="options.ScriptUuid">
                          {{ options.ScriptName}}
                        </a-select-option>
                      </a-select>
                  </div>
                  <!--SQL报表列表-->
                  <div  v-if="diyData.type==32">
                    <a-select @change="UpdateNodeData"
                              :allowClear="true"
                              v-model="configObject.style.diy[diyIndex].value"
                    >
                      <a-select-option v-for="options in SQLReportList" :key="options.Uuid" :value="options.Uuid">
                        {{ options.Name}}
                      </a-select-option>
                    </a-select>
                  </div>
                </a-form-item>
              </div>
            </template>

          </a-form>
        </a-tab-pane>
        <a-tab-pane key="3" :tab="$t('displayConfig.Properties.TabHeaterBehavior')" v-if="typeof(configObject.action)!='undefined'" style="padding:5px;">
          <div v-if="configObject && configObject.action">
            <div v-for="(event,index) in configObject.action" :key="index">
              <div  style="margin-top:5px;border:#13c2c2  solid 1px;">
                <div style="padding:5px;border-bottom:#13c2c2  solid 1px;">
                  {{$t('displayConfig.Properties.TabHeaterBehavior')}}-{{index+1}}
                  <a-tooltip placement="top">
                    <template slot="title">
                      <span>{{$t('displayConfig.Properties.ActionsTips')}}</span>
                    </template>
                    <a-icon type="delete" @click="delBindAction(index)" theme="twoTone" two-tone-color="#eb2f96" style="float: right;cursor:pointer;"/>
                  </a-tooltip>

                </div>
                <a-form  :label-col="{ span: 7}" :wrapper-col="{ span: 16 }">
                  <a-form-item :label="$t('displayConfig.Properties.ComponentEvent')">
                    <a-select
                        v-model="event.type" @change="UpdateNodeData"
                    >
                      <a-select-option v-for="options in [{label:'displayConfig.Properties.action.click',value:'click'},{label:'displayConfig.Properties.action.MouseDown',value:'mousedown'},{label:'displayConfig.Properties.action.MouseUp',value:'mouseup'},{label:'displayConfig.Properties.action.dbClick',value:'dblclick'},{label:'displayConfig.Properties.action.mouseenter',value:'mouseenter',cannotSelect: true},{label:'displayConfig.Properties.action.mouseleave',value:'mouseleave',cannotSelect: true}]" :key="options.value" :value="options.value">
                        {{ $t(options.label)}}
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ComponentActionAuth')">
                    <a-select @change="UpdateNodeData"
                        mode="multiple"
                        v-model="event.actionAuth"
                    >
                      <a-select-option v-for="options in RoleList" :key="options.RoleId" :value="options.RoleId">
                        <span v-if="options.RoleId=='Operator'">{{ $t("account.settings.UserList.RoleOperator") }}</span>
                        <span v-if="options.RoleId=='User'">{{ $t("account.settings.UserList.RoleUser") }}</span>
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ComponentActionPassword')">
                    <a-input @change="UpdateNodeData" type='text' v-model="event.ActionPassword" :style="{'text-security':'disc', '-webkit-text-security':'disc'}" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ComponentActionVoice')">
                    <a-input @change="UpdateNodeData" v-model="event.actionVoice">
                    </a-input>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ComponentAction')">
                    <a-select
                        v-model="event.action" @change="changeAction($event,index)"
                    >
                      <a-select-option v-for="options in [
                          {label:'displayConfig.Properties.action.openLink',value:'link',cannotSelect: true},
                          {label:'displayConfig.Properties.action.SetValue',value:'SetValue',cannotSelect: true},
                          {label:'displayConfig.Properties.action.Visible',value:'visible'},
                          {label:'displayConfig.Properties.action.SysScript',value:'SysScript'},
                          {label:'displayConfig.Properties.action.DeviceView',value:'DeviceView',cannotSelect: true},
                          {label:'displayConfig.Properties.action.RestApi',value:'RestApi',cannotSelect: true},
                          {label:'displayConfig.Properties.action.animation',value:'Animation',cannotSelect: true}
                          ]" :key="options.value" :value="options.value">
                        {{ $t(options.label)}}
                      </a-select-option>
                    </a-select>
                  </a-form-item>

                  <a-form-item :label="$t('displayConfig.Properties.SecondConfirm')">
                    <a-switch default-checked v-model="event.actionConfirm" @change="UpdateNodeData">
                      <a-icon slot="checkedChildren" type="check" />
                      <a-icon slot="unCheckedChildren" type="close" />
                    </a-switch>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.SetDelay')"  v-if="event.action=='SetValue'">
                    <a-input @change="UpdateNodeData" v-model="event.SetDelay" value="1000"></a-input>
                  </a-form-item>

                  <div v-if="event.action=='RestApi'">
                    <a-form-item :label="$t('displayConfig.Properties.RestApiName')">
                      <a-input  @change="UpdateNodeData" v-model="event.RestApi.Name"></a-input>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.RestApiFrom')">
                      <a-select @change="UpdateNodeData"
                          v-model="event.RestApi.IsSystem"
                      >
                        <a-select-option key="1" value="1">
                          {{ $t('displayConfig.Properties.ExternRestApi')}}
                        </a-select-option>
                        <a-select-option key="3" value="2">
                          {{ $t('displayConfig.Properties.SystemRestApi')}}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.RestApiType')">
                      <a-select @change="UpdateNodeData"
                          v-model="event.RestApi.Type"
                      >
                        <a-select-option key="Post" value="Post">
                          Post
                        </a-select-option>
                        <a-select-option key="Get" value="Get">
                          Get
                        </a-select-option>
                      </a-select>
                    </a-form-item>

                    <a-form-item :label="$t('displayConfig.Properties.RestApiUrl')">
                      <a-input @change="UpdateNodeData" v-model="event.RestApi.Url"></a-input>
                    </a-form-item>

                    <a-form-item :label="$t('displayConfig.Properties.RestApiParam')">
                      <a-textarea @change="UpdateNodeData" style="font-size: 16px" :rows="10" v-model="event.RestApi.Params"></a-textarea>
                    </a-form-item>
                  </div>
                  <div v-if="event.action=='SetValue'">
                    <div style="margin-top: 5px;text-align: center">
                      <a-button key="submit" block  style="width: 100px" @click="addBindSetValue(index)">{{$t('displayConfig.Properties.AddSetValue')}}</a-button>
                    </div>
                    <div v-if="typeof event.setValue !='undefined'" >
                      <div style="margin: 5px;border:#95B8E7 solid 1px;" v-for="(setValueItem,setValueIndex) in event.setValue" :key="setValueIndex">
                        <div style="padding:5px;border-bottom:#95B8E7 solid 1px;">
                          <a-tooltip placement="top">
                            <template slot="title">
                              <span>{{$t('displayConfig.Properties.SetValueTips')}}</span>
                            </template>
                            <a-icon type="close-circle" @click="delBindSetValue(index,setValueIndex)" theme="twoTone" two-tone-color="#eb2f96" style="float: right;cursor:pointer;"/>
                          </a-tooltip>
                          {{$t('displayConfig.Properties.SetValueInfo')}}-{{setValueIndex+1}}

                        </div>
                        <a-form-item :label="$t('displayConfig.Properties.IsManualLabel')">
                          <a-checkbox  v-model="setValueItem.IsManual" @change="UpdateNodeData">
                            <label > {{$t('displayConfig.Properties.IsManual')}}</label>
                          </a-checkbox>
                        </a-form-item>

                        <a-form-item :label="$t('displayConfig.Properties.ComponentFromBandDevice')" v-if="setValueItem.dataID!=''">
                          <a-checkbox disabled v-model="setValueItem.isBandDevice">
                            <label v-if="setValueItem.isBandDevice"> {{$t('displayConfig.Properties.ComponentIsBandDevice')}}</label>
                            <label v-else> {{$t('displayConfig.Properties.ComponentIsBandDeviceModel')}}</label>
                          </a-checkbox>
                        </a-form-item>

                        <a-form-item :label="$t('displayConfig.Properties.ComponentBandDevice')" v-if="setValueItem.isBandDevice">
                          <a-tree-select @change="UpdateNodeData"
                              disabled
                              v-model="setValueItem.deviceSN"
                              style="width: 100%"
                              :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
                              :tree-data="deviceTreeData"
                              :replace-fields="{ value: 'key',title:'text'}"
                              placeholder="Please select"
                              tree-default-expand-all
                          >
                          </a-tree-select>
                        </a-form-item>

                        <a-form-item :label="$t('displayConfig.Properties.ComponentBandData')">
                          <a-input  v-model="setValueItem.dataName" @change="UpdateNodeData">
                            <a-tooltip placement="top" slot="addonAfter">
                              <template slot="title">
                                <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
                              </template>
                              <icon-font @click="actionIndex=index;ShowDeviceDataModel(setValueIndex,'setValue')"  type="icon-xuanzeshuju"  />
                            </a-tooltip>
                          </a-input>
                        </a-form-item>

                        <a-form-item :label="$t('displayConfig.Properties.AutoSetValue')" v-if="!setValueItem.IsManual">
                          <a-input  @change="UpdateNodeData" v-model="setValueItem.AutoSetValue"></a-input>
                        </a-form-item>
                        <a-form-item :label="$t('displayConfig.Properties.SetPassword')">
                          <a-input  @change="UpdateNodeData" v-model="setValueItem.SetPassword"></a-input>
                        </a-form-item>
                      </div>
                    </div>

                  </div>
                  <div v-if="event.action=='visible'">
                    <a-form-item :label="$t('displayConfig.Properties.ComponentClickShow')">
                      <a-select :allowClear="true" @change="UpdateNodeData"
                                mode="multiple"
                                v-model="event.showItems"
                      >
                        <a-select-option v-for="options in GroupList" :key="options.ID" :value="options.ID">
                          {{ options.Name}}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.ComponentClickHide')">
                      <a-select :allowClear="true" @change="UpdateNodeData"
                                mode="multiple"
                                v-model="event.hideItems"
                      >
                        <a-select-option v-for="options in GroupList" :key="options.ID" :value="options.ID">
                          {{ options.Name}}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </div>
                  <div v-if="event.action=='link'">
                    <a-form-item :label="$t('displayConfig.Properties.linkType')">
                      <a-select @change="UpdateNodeData"
                          v-model="event.link.linkType"
                          allowClear
                      >
                        <a-select-option v-for="options in [{label:'displayConfig.Properties.linkInside',value:'Inside'},{label:'displayConfig.Properties.linkExternal',value:'External'}]" :key="options.value" :value="options.value">
                          {{  $t(options.label)}}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                    <div v-if="event.link.linkType=='Inside'">
                      <a-form-item :label="$t('displayConfig.Properties.linkIAppUUID')">
                        <a-select @change="changeLinkInsideDisplay(event.link)"
                            v-model="event.link.Inside.displayUUID"
                            allowClear
                        >
                          <a-select-option v-for="options in configurationModel" :key="options.uuid" :value="options.uuid">
                            {{ options.name}}{{ options.displayType == 2 ? '（数字孪生）' : '' }}
                          </a-select-option>
                        </a-select>
                      </a-form-item>
                      <a-form-item :label="$t('displayConfig.Properties.linkIAppPageUUID')">
                        <a-select @change="UpdateNodeData"
                            v-model="event.link.Inside.pageUUID"
                        >
                          <a-select-option v-for="options in generateTargetPage(event.link.Inside.displayUUID)" :key="options.value" :value="options.value">
                            {{ options.label}}
                          </a-select-option>
                        </a-select>
                      </a-form-item>
                    </div>
                    <div v-else>
                      <a-form-item label="网页来源">
                        <a-select :value="event.link.ExternalSource || 'Url'" @change="value => changeLinkExternalSource(event.link, value)">
                          <a-select-option value="Url">输入网址</a-select-option>
                          <a-select-option value="Page">选择页面</a-select-option>
                        </a-select>
                      </a-form-item>
                      <div v-if="event.link.ExternalSource === 'Page'">
                        <a-form-item label="应用">
                          <a-select @change="changeLinkExternalDisplay(event.link)"
                              v-model="event.link.ExternalPage.displayUUID"
                              allowClear
                          >
                            <a-select-option v-for="options in get2DConfigurationModel()" :key="options.uuid" :value="options.uuid">
                              {{ options.name }}
                            </a-select-option>
                          </a-select>
                        </a-form-item>
                        <a-form-item label="页面">
                          <a-select @change="changeLinkExternalPage(event.link)"
                              v-model="event.link.ExternalPage.pageUUID"
                              allowClear
                          >
                            <a-select-option v-for="options in generateTargetPage(event.link.ExternalPage.displayUUID)" :key="options.value" :value="options.value">
                              {{ options.label}}
                            </a-select-option>
                          </a-select>
                        </a-form-item>
                      </div>
                      <a-form-item :label="$t('displayConfig.Properties.linkExternalUrl')">

                        <a-input type="text" v-model="event.link.External" @change="UpdateNodeData" :disabled="event.link.ExternalSource === 'Page'">

                        </a-input>
                      </a-form-item>

                      <a-form-item :label="$t('displayConfig.Properties.OpenLinkExternalType')">

                        <a-select v-model="event.link.OpenExternalType" @change="UpdateNodeData">
                          <a-select-option value="self">{{$t('displayConfig.Properties.OpenLinkExternalSelf')}}</a-select-option>
                          <a-select-option value="new">{{$t('displayConfig.Properties.OpenLinkExternalNew')}}</a-select-option>
                        </a-select>

                      </a-form-item>
                    </div>
                    <a-form-item :label="$t('displayConfig.Properties.isLinkPopUp')">
                      <a-checkbox :checked="event.link.isPopUp"  v-model="event.link.isPopUp" @change="UpdateNodeData">
                        {{ $t('displayConfig.Properties.isLinkPopUp') }}
                      </a-checkbox>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.autoClose')" v-if="event.link.isPopUp">
                      <a-checkbox :checked="event.link.autoClose"  v-model="event.link.autoClose" @change="UpdateNodeData">
                        {{ $t('displayConfig.Properties.autoClose') }}
                      </a-checkbox>
                    </a-form-item>
                    <div v-if="event.link.isPopUp&&event.link.linkType!=='Inside'">
                      <a-form-item :label="$t('displayConfig.Properties.linkExternalWidth')">

                        <a-input type="text" v-model="event.link.width" @change="UpdateNodeData">

                        </a-input>
                      </a-form-item>
                      <a-form-item :label="$t('displayConfig.Properties.linkExternalHeight')">

                        <a-input type="text" v-model="event.link.height" @change="UpdateNodeData">

                        </a-input>
                      </a-form-item>

                      <a-form-item :label="$t('displayConfig.Properties.linkExternalTitle')">

                        <a-input type="text" v-model="event.link.title" @change="UpdateNodeData">

                        </a-input>
                      </a-form-item>
                    </div>
                  </div>
                  <div v-if="event.action=='SysScript'">
                    <a-form-item :label="$t('displayConfig.Properties.action.SysScript')">
                      <a-select @change="UpdateNodeData"
                          :allowClear="true"
                          mode="multiple"
                          v-model="event.ScriptList"
                      >
                        <a-select-option v-for="options in ScriptDataSource" :key="options.ScriptUuid" :value="options.ScriptUuid">
                          {{ options.ScriptName}}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </div>
                  <div v-if="event.action=='DeviceView'">
                    <a-form-item :label="$t('displayConfig.Properties.SelectDevice')">
                      <a-tree-select @change="UpdateNodeData"
                          show-search
                          tree-node-filter-prop="title"
                          @select="SelectTreeDevice"
                          @click="actionIndex=index"
                          :dropdown-style="{ 'z-index': 9999999,maxHeight: '400px', overflow: 'auto' }"
                          :tree-data="deviceTreeData"
                          v-model="event.DeviceView.key"
                          :replace-fields="{ value: 'key',title:'text'}"
                          placeholder="Please select"
                          tree-default-expand-all
                      >
                      </a-tree-select>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.isLinkPopUp')">
                      <a-checkbox :checked="event.DeviceView.isPopUp"  v-model="event.DeviceView.isPopUp" @change="UpdateNodeData">
                        {{ $t('displayConfig.Properties.isLinkPopUp') }}
                      </a-checkbox>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.isContainer')">
                      <a-checkbox :checked="event.DeviceView.isContainer"  v-model="event.DeviceView.isContainer" @change="UpdateNodeData">
                        {{ $t('displayConfig.Properties.isContainer') }}
                      </a-checkbox>
                    </a-form-item>
                  </div>
                  <div v-if="event.action=='Animation'">
                    <a-form-item :label="$t('displayConfig.Properties.action.animationList')">
                      <a-select v-model="configObject.animate.selected"   mode="multiple" @change="UpdateNodeData">
                        <a-select-option v-for="options in configObject.animate.animateList" :key="options.id" :value="options.id">
                          {{$t(options.name)}}
                        </a-select-option>
                        <a-select-option key="visible" value="visible">
                          {{$t('component.public.animateVisible')}}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.action.animationStatus')">
                      <a-select @change="UpdateNodeData"
                          :allowClear="true"
                          v-model="event.animationStatus"
                      >
                        <a-select-option v-for="options in [
                          {label:'displayConfig.Properties.action.animationStart',value:'start',cannotSelect: true},
                          {label:'displayConfig.Properties.action.animationStop',value:'stop',cannotSelect: true}
                          ]" :key="options.value" :value="options.value">
                          {{ $t(options.label)}}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </div>
                </a-form>
              </div>
            </div>
          </div>
          <div style="margin-top: 5px">
            <a-button key="submit" block type="primary" @click="addBindAction">{{$t('displayConfig.Properties.ComponentAddBandActionBtn')}}</a-button>
          </div>
        </a-tab-pane>
        <a-tab-pane key="4" :tab="$t('displayConfig.Properties.TabHeaterActive')"  v-if="typeof(configObject.active)!='undefined'"  style="padding:5px;">
          <template v-for="(activeItem,index) in configObject.active">
            <a-card  :key="index"  :title="$t(activeItem.name)" :bodyStyle="{padding:'2px'}" style="margin-top: 5px">
              <a slot="extra" >
                <div v-if="(typeof activeItem.isStatus!='undefined')&&(activeItem.isStatus)">
                  <a-tooltip placement="top" slot="addonAfter">
                    <template slot="title">
                      <span>{{$t('configComponent.status.AddStatus')}}</span>
                    </template>
                    <a-icon @click="addStatus" type="plus" style="line-height:30px" />
                  </a-tooltip>
                </div>

                <div v-else-if="(typeof activeItem.isMenu!='undefined')&&(activeItem.isMenu)">
                  <a-tooltip placement="top" slot="addonAfter">
                    <template slot="title">
                      <span>{{$t('configComponent.status.AddMenu')}}</span>
                    </template>
                    <a-icon @click="addMenu" type="plus" style="line-height:30px" />
                  </a-tooltip>
                </div>

                <a-icon @click="delBindActive(index)" v-else type="minus" style="line-height:30px" />
              </a>
              <a-form  :label-col="{ span: 5}" :wrapper-col="{ span: 19 }">
                <a-form-item :label="$t('displayConfig.Properties.ComponentFromBandDevice')" v-if="activeItem.condition.bandType==1&&activeItem.condition.dataID!=''">
                  <a-checkbox disabled v-model="activeItem.condition.isBandDevice">
                    <label v-if="activeItem.condition.isBandDevice"> {{$t('displayConfig.Properties.ComponentIsBandDevice')}}</label>
                    <label v-else> {{$t('displayConfig.Properties.ComponentIsBandDeviceModel')}}</label>
                  </a-checkbox>
                </a-form-item>

                <a-form-item :label="$t('displayConfig.Properties.ComponentBandDevice')" v-if="activeItem.condition.isBandDevice&&activeItem.condition.bandType==1">
                  <a-tree-select
                      disabled
                      v-model="activeItem.condition.deviceSN"
                      style="width: 100%"
                      :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
                      :tree-data="deviceTreeData"
                      :replace-fields="{ value: 'key',title:'text'}"
                      placeholder="Please select"
                      tree-default-expand-all
                  >
                  </a-tree-select>
                </a-form-item>

                <a-form-item :label="$t('displayConfig.Properties.ComponentBandData')" v-if="activeItem.condition.bandType==1">
                  <a-input  v-model="activeItem.condition.dataName" >
                    <a-tooltip placement="top" slot="addonAfter">
                      <template slot="title">
                        <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
                      </template>
                      <icon-font @click="ShowDeviceDataModel(index,'active')" type="icon-xuanzeshuju"  />
                    </a-tooltip>
                  </a-input>
                </a-form-item>

                <div v-if="(typeof activeItem.isExpression!='undefined')&&(activeItem.isExpression)">
                  <a-form-item :label="$t('component.public.Operator')">
                    <a-select v-model="activeItem.condition.operator">
                      <a-select-option value="==">
                        =
                      </a-select-option>
                      <a-select-option value=">">
                        >
                      </a-select-option>
                      <a-select-option value=">=">
                        >=
                      </a-select-option>
                      <a-select-option value="<">
                        &lt;
                      </a-select-option>
                      <a-select-option value="<=">
                        &lt;=
                      </a-select-option>
                      <a-select-option value="<>">
                        >=  * &lt;=
                      </a-select-option>
                      <a-select-option value="!=">
                        !=
                      </a-select-option>
                    </a-select>
                  </a-form-item>

                  <a-form-item :label="$t('component.public.OperatorValue')" >
                    <a-input   v-model="activeItem.condition.OperatorValue" >
                    </a-input>
                  </a-form-item>

                  <a-form-item :label="$t('component.public.OperatorMaxValue')" v-if="activeItem.condition.operator=='<>'">
                    <a-input   v-model="activeItem.condition.OperatorMaxValue" >
                    </a-input>
                  </a-form-item>
                </div>
                <div v-else-if="(typeof activeItem.isSwitch!='undefined')&&(activeItem.isSwitch)">
                  <a-form-item :label="$t('displayConfig.Properties.ComponentActionAuth')">
                    <a-select
                        mode="multiple"
                        v-model="activeItem.condition.actionAuth"
                    >
                      <a-select-option v-for="options in RoleList" :key="options.RoleId" :value="options.RoleId">
                        <span v-if="options.RoleId=='Operator'">{{ $t("account.settings.UserList.RoleOperator") }}</span>
                        <span v-if="options.RoleId=='User'">{{ $t("account.settings.UserList.RoleUser") }}</span>
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ComponentActionVoice')">
                    <a-input v-model="activeItem.condition.actionVoice"></a-input>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.SetPassword')">
                    <a-input  v-model="activeItem.condition.SetPassword"></a-input>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ConfirmationDialog')">
                    <a-select
                        v-model="activeItem.condition.ConfirmationDialog"
                    >
                      <a-select-option value="1">
                       Yes
                      </a-select-option>
                      <a-select-option value="0">
                        No
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t(setItem.name)" v-for="(setItem,ItemIndex) in activeItem.condition.AutoSet" :key="ItemIndex">
                    <a-input v-model="activeItem.condition.AutoSet[ItemIndex].value"></a-input>
                  </a-form-item>
                </div>
                <div v-else-if="(typeof activeItem.isStatus!='undefined')&&(activeItem.isStatus)">
                  <a-table v-if="activeItem.isTextStatus" :pagination="false" :loading="false"  row-key="uuid" :data-source="activeItem.condition.StatusList" :columns="StatusColumns" class="ant-table-tbody">
                    <template v-for="(item, index) in StatusColumns" :slot="item.slotName">
                      <span :key="index">{{ $t(item.slotName) }}</span>
                    </template>
                    <template slot="StatusValue" slot-scope="text, record">
                      <a-tooltip>
                        <template slot="title">
                          {{ $t('configComponent.status.StatusValueTips') }}
                        </template>
                        <a-input :value="text" v-model="record.value"/>
                      </a-tooltip>
                    </template>
                    <template slot="StatusValue2" slot-scope="text, record">
                      <a-tooltip>
                        <template slot="title">
                          {{ $t('configComponent.status.StatusValue2Tips') }}
                        </template>
                        <a-input :value="text" v-model="record.value2" :disabled="record.StatusOpt!='<>'&&record.StatusOpt!='<!>'"/>
                      </a-tooltip>
                    </template>
                    <template slot="StatusOpt" slot-scope="text, record">
                      <a-select v-model="record.StatusOpt">
                        <a-select-option value="==">
                          =
                        </a-select-option>
                        <a-select-option value=">">
                          >
                        </a-select-option>
                        <a-select-option value=">=">
                          >=
                        </a-select-option>
                        <a-select-option value="<">
                          &lt;
                        </a-select-option>
                        <a-select-option value="<=">
                          &lt;=
                        </a-select-option>
                        <a-select-option value="<>">
                          &&
                        </a-select-option>
                        <a-select-option value="<!>">
                          ||
                        </a-select-option>
                        <a-select-option value="!=">
                          !=
                        </a-select-option>
                      </a-select>
                    </template>
                    <template slot="StatusText" slot-scope="text,record">
                      <a-input :value="text" v-model="record.Text"/>
                    </template>
                    <template slot="StatusTextColor" slot-scope="text,record">
                      <a-input type="color" :value="text" v-model="record.TextColor"></a-input>
                    </template>
                    <template slot="OptStatus" slot-scope="text, record,index">
                      <div class="item">
                        <a-tooltip placement="top" slot="addonAfter">
                          <template slot="title">
                            <span>{{$t('configComponent.status.DeleteStatus')}}</span>
                          </template>
                          <a-icon @click="activeItem.condition.StatusList.splice(index,1)" type="delete" style="line-height:30px" />
                        </a-tooltip>
                      </div>
                    </template>
                  </a-table>
                  <a-table v-if="activeItem.isLineStatus" :pagination="false" :loading="false"  row-key="uuid" :data-source="activeItem.condition.StatusList" :columns="LineStatusColumns" class="ant-table-tbody">
                    <template v-for="(item, index) in LineStatusColumns" :slot="item.slotName">
                      <span :key="index">{{ $t(item.slotName) }}</span>
                    </template>
                    <template slot="StatusValue" slot-scope="text, record">
                      <a-tooltip>
                        <template slot="title">
                          {{ $t('configComponent.status.StatusValueTips') }}
                        </template>
                        <a-input :value="text" v-model="record.value"/>
                      </a-tooltip>
                    </template>
                    <template slot="StatusValue2" slot-scope="text, record">
                      <a-tooltip>
                        <template slot="title">
                          {{ $t('configComponent.status.StatusValue2Tips') }}
                        </template>
                        <a-input :value="text" v-model="record.value2" :disabled="record.StatusOpt!='<>'&&record.StatusOpt!='<!>'"/>
                      </a-tooltip>
                    </template>
                    <template slot="StatusOpt" slot-scope="text, record">
                      <a-select v-model="record.StatusOpt" style="width: 55px">
                        <a-select-option value="==">
                          =
                        </a-select-option>
                        <a-select-option value=">">
                          >
                        </a-select-option>
                        <a-select-option value=">=">
                          >=
                        </a-select-option>
                        <a-select-option value="<">
                          &lt;
                        </a-select-option>
                        <a-select-option value="<=">
                          &lt;=
                        </a-select-option>
                        <a-select-option value="<>">
                          &&
                        </a-select-option>
                        <a-select-option value="<!>">
                          ||
                        </a-select-option>
                        <a-select-option value="!=">
                          !=
                        </a-select-option>
                      </a-select>
                    </template>
                    <template slot="LineStatusBlink" slot-scope="text,record">
                      <a-select v-model="record.Blink" defaultValue="0" style="width: 65px">
                        <a-select-option value='0'>
                          False
                        </a-select-option>
                        <a-select-option value='1'>
                          True
                        </a-select-option>
                      </a-select>
                    </template>
                    <template slot="LineStatusBlinkSpeed" slot-scope="text,record">
                      <a-input :value="text" v-model="record.BlinkSpeed"/>
                    </template>
                    <template slot="StatusTextColor" slot-scope="text,record">
                      <a-input type="color" :value="text" v-model="record.TextColor"></a-input>
                    </template>
                    <template slot="OptStatus" slot-scope="text, record,index">
                      <div class="item">
                        <a-tooltip placement="top" slot="addonAfter">
                          <template slot="title">
                            <span>{{$t('configComponent.status.DeleteStatus')}}</span>
                          </template>
                          <a-icon @click="activeItem.condition.StatusList.splice(index,1)" type="delete" style="line-height:30px" />
                        </a-tooltip>
                      </div>
                    </template>
                  </a-table>
                  <a-table v-if="activeItem.isImageStatus" :pagination="false" :loading="false"  row-key="uuid" :data-source="activeItem.condition.StatusList" :columns="ImagesStatusColumns" class="ant-table-tbody">
                    <template v-for="(item, index) in ImagesStatusColumns" :slot="item.slotName">
                      <span :key="index">{{ $t(item.slotName) }}</span>
                    </template>
                    <template slot="StatusValue" slot-scope="text, record">
                      <a-tooltip>
                        <template slot="title">
                          {{ $t('configComponent.status.StatusValueTips') }}
                        </template>
                        <a-input :value="text" v-model="record.value"/>
                      </a-tooltip>
                    </template>
                    <template slot="StatusValue2" slot-scope="text, record">
                      <a-tooltip>
                        <template slot="title">
                          {{ $t('configComponent.status.StatusValue2Tips') }}
                        </template>
                        <a-input :value="text" v-model="record.value2" :disabled="record.StatusOpt!='<>'&&record.StatusOpt!='<!>'"/>
                      </a-tooltip>
                    </template>
                    <template slot="StatusOpt" slot-scope="text, record">
                      <a-select v-model="record.StatusOpt">
                        <a-select-option value="==">
                          =
                        </a-select-option>
                        <a-select-option value=">">
                          >
                        </a-select-option>
                        <a-select-option value=">=">
                          >=
                        </a-select-option>
                        <a-select-option value="<">
                          &lt;
                        </a-select-option>
                        <a-select-option value="<=">
                          &lt;=
                        </a-select-option>
                        <a-select-option value="<>">
                          &&
                        </a-select-option>
                        <a-select-option value="<!>">
                          ||
                        </a-select-option>
                        <a-select-option value="!=">
                          !=
                        </a-select-option>
                      </a-select>
                    </template>
                    <template slot="ImageTitle" slot-scope="text,record,index">
                      <div class="item">
                        <img @click="showSystemImageModel(2,index)" v-if="text!=''"  style="width: 32px;height:32px;cursor: pointer" :src="text" />
                        <img @click="showSystemImageModel(2,index)" v-else :style="{width: '32px',height:'32px',cursor: 'pointer'}" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAWZJREFUWEftlzFOAzEQRV9OAQUUCCooETXkBmmg5xS0iDvQ0wINNwBq6qQCJRJCAnEI0Ece5KxMNh57WQlhKYW16/nPf8beyYCex6BnfVIA68BmR2BPwHMcuwlwCxx0JG5h74ChTWIA7foxPNBLXQzb3BYgN+ZSoIdyYI6wMoU5LAe+Nhk78A/Q5sANsBdS8gCMHOlxp+DjB7Hce8QFcAUcAqfAfQDZD/Nr4CjDCRfAK7AC7ACTILYNjIE3YPXPA3hSsAFME864UqA4yxahauUE2AUugbMobYrjBtDitmMocbmloVtOx1p1oiK12ikCWFRrsbgEdTosdTFEJwApcYNtQpwHZ6p9CxaJpyDeawIsI56C0H1S7ECOeBNC8yIAj7hBFBdhiXjxPVAqXgRQQ9wNUEvcBVBTvBUg1ZZbG12rTbd4a8BLsyuOCTN6jOxXL4BjW/Xbf81mgH7fI7epzN5u24LeAT4B+PijITzWRRcAAAAASUVORK5CYII="/>
                      </div>
                    </template>
                    <template slot="OptStatus" slot-scope="text, record,index">
                      <div class="item">
                        <a-tooltip placement="top" slot="addonAfter">
                          <template slot="title">
                            <span>{{$t('configComponent.status.DeleteStatus')}}</span>
                          </template>
                          <a-icon @click="activeItem.condition.StatusList.splice(index,1)" type="delete" style="line-height:30px" />
                        </a-tooltip>
                      </div>
                    </template>
                  </a-table>
                  <a-table v-if="activeItem.isComBox" :pagination="false" :loading="false"  row-key="uuid" :data-source="activeItem.condition.StatusList" :columns="ComBoxStatusColumns" class="ant-table-tbody">
                    <template v-for="(item, index) in ComBoxStatusColumns" :slot="item.slotName">
                      <span :key="index">{{ $t(item.slotName) }}</span>
                    </template>
                    <template slot="ComBoxValue" slot-scope="text, record">
                      <a-input :value="text" v-model="record.value"/>
                    </template>
                    <template slot="StatusText" slot-scope="text,record">
                      <a-input :value="text" v-model="record.Text"/>
                    </template>
                    <template slot="OptStatus" slot-scope="text, record,index">
                      <div class="item">
                        <a-tooltip placement="top" slot="addonAfter">
                          <template slot="title">
                            <span>{{$t('configComponent.status.DeleteStatus')}}</span>
                          </template>
                          <a-icon @click="activeItem.condition.StatusList.splice(index,1)" type="delete" style="line-height:30px" />
                        </a-tooltip>
                      </div>
                    </template>
                  </a-table>
                </div>
                <div v-else-if="(typeof activeItem.isVoiceStatus!='undefined')&&(activeItem.isVoiceStatus)">
                  <DataGrid  style="height:250px" :data="activeItem.condition.StatusList" :clickToEdit="true" selectionMode="cell" editMode="cell">
                    <GridColumn field="value" :title="$t('configComponent.status.StatusValue')" align="center" :editable="true">
                      <template slot="edit" slot-scope="scope">
                        <TextBox  v-model="scope.row.value" :precision="2"></TextBox>
                      </template>
                    </GridColumn>
                    <GridColumn field="Text" :title="$t('configComponent.status.StatusText')" v-if="activeItem.isTextStatus" align="center" :editable="true">
                      <template slot="edit" slot-scope="scope">
                        <TextBox  v-model="scope.row.Text" :precision="2"></TextBox>
                      </template>
                    </GridColumn>
                    <GridColumn field="TextColor" :title="$t('configComponent.status.StatusTextColor')" v-if="activeItem.isTextStatus" align="center">
                      <template slot="body" slot-scope="scope">
                        <div class="item">
                          <a-input type="color" v-model="scope.row.TextColor"></a-input>
                        </div>
                      </template>
                    </GridColumn>
                    <GridColumn field="voiceUrl" :title="$t('configComponent.status.StatusVoiceUrl')" v-if="activeItem.isTextStatus" align="center" :editable="true">
                      <template slot="edit" slot-scope="scope">
                        <TextBox  v-model="scope.row.voiceUrl" :precision="2"></TextBox>
                      </template>
                    </GridColumn>

                    <GridColumn field="opt"   :title="$t('configComponent.status.OptStatus')" align="center">
                      <template slot="body" slot-scope="scope">
                        <div class="item">
                          <a-tooltip placement="top" slot="addonAfter">
                            <template slot="title">
                              <span>{{$t('configComponent.status.DeleteStatus')}}</span>
                            </template>
                            <a-icon @click="activeItem.condition.StatusList.splice(scope.rowIndex,1)" type="delete" style="line-height:30px" />
                          </a-tooltip>
                        </div>
                      </template>
                    </GridColumn>
                  </DataGrid>
                </div>
                <div v-else-if="(typeof activeItem.isMenu!='undefined')&&(activeItem.isMenu)">
                  <a-table :scroll="{ x: 600 }" :pagination="false" :loading="false"  row-key="uuid" :data-source="activeItem.condition.PageList" :columns="DropMenuColumns" class="ant-table-tbody">
                    <template v-for="(item, index) in DropMenuColumns" :slot="item.slotName">
                      <span :key="index">{{ $t(item.slotName) }}</span>
                    </template>
                    <template slot="MenuText" slot-scope="text, record">
                      <a-input :value="text" v-model="record.MenuName"/>
                    </template>
                    <template slot="APPListName" slot-scope="text,record">
                      <a-select
                                v-model="record.DisPlayID"
                                allowClear
                      >
                        <a-select-option v-for="options in configurationModel" :key="options.uuid" :value="options.uuid">
                          {{ options.name}}
                        </a-select-option>
                      </a-select>
                    </template>
                    <template slot="MenuPageName" slot-scope="text,record">
                      <a-select
                          v-model="record.PageID"
                      >
                        <a-select-option v-for="options in generateTargetPage(record.DisPlayID)" :key="options.value" :value="options.value">
                          {{ options.label}}
                        </a-select-option>
                      </a-select>
                    </template>
                    <template slot="MenuPagePop" slot-scope="text, record">
                      <a-switch default-checked v-model="record.IsPopUp">
                        <a-icon slot="checkedChildren" type="check" />
                        <a-icon slot="unCheckedChildren" type="close" />
                      </a-switch>
                    </template>
                    <template slot="OptStatus" slot-scope="text, record,index">
                      <div class="item">
                        <a-tooltip placement="top" slot="addonAfter">
                          <template slot="title">
                            <span>{{$t('configComponent.status.DeleteStatus')}}</span>
                          </template>
                          <a-icon @click="activeItem.condition.PageList.splice(index,1)" type="delete" style="line-height:30px" />
                        </a-tooltip>
                      </div>
                    </template>
                  </a-table>
                </div>
              </a-form>
            </a-card>
          </template>
          <div style="margin-top: 5px">
            <a-button key="submit" block type="primary" @click="addBindActive">{{$t('displayConfig.Properties.ComponentAddBandActiveBtn')}}</a-button>
          </div>
        </a-tab-pane >
        <a-tab-pane key="5" :tab="$t('displayConfig.Properties.TabHeaterAnimate')" v-if="typeof(configObject.animate)!='undefined'" style="padding:5px;">
          <a-divider>{{$t('component.public.animate')}}</a-divider>
          <template v-if="configObject && configObject.animate.animateList">
            <a-form  :label-col="{ span: 7}" :wrapper-col="{ span: 16 }" style="border:#f0f0f0 solid 1px;">
              <a-form-item :label="$t('component.public.animate')">
                <a-select v-model="configObject.animate.selected"  @change="UpdateNodeData" mode="multiple">
                  <a-select-option v-for="options in configObject.animate.animateList" :key="options.id" :value="options.id">
                    {{$t(options.name)}}
                  </a-select-option>
                  <a-select-option key="visible" value="visible">
                    {{$t('component.public.animateVisible')}}
                  </a-select-option>
                  <a-select-option key="animateMove" value="animateMove">
                    {{$t('component.public.animateMove')}}
                  </a-select-option>
                </a-select>
              </a-form-item>

              <div v-if="configObject.animate && configObject.animate.selected && configObject.animate.selected.includes('animateMove')">

                <a-form-item :label="$t('component.public.animateX')" >
                  <a-input  v-model="configObject.animate.move.x.dataName" @change="UpdateNodeData">
                    <a-tooltip placement="top" slot="addonAfter">
                      <template slot="title">
                        <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
                      </template>
                      <icon-font @click="ShowDeviceDataModel(0,'animateMoveX')" type="icon-xuanzeshuju"  />
                    </a-tooltip>
                  </a-input>
                </a-form-item>

                <a-form-item :label="$t('component.public.animateY')">
                  <a-input  v-model="configObject.animate.move.y.dataName" @change="UpdateNodeData">
                    <a-tooltip placement="top" slot="addonAfter">
                      <template slot="title">
                        <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
                      </template>
                      <icon-font @click="ShowDeviceDataModel(0,'animateMoveY')" type="icon-xuanzeshuju"  />
                    </a-tooltip>
                  </a-input>
                </a-form-item>
              </div>
              <div v-for="(elementList,index) in configObject.animate.animateElement" :key="elementList.id" >
                <div v-if="configObject.animate.selected && configObject.animate.selected.includes(elementList.id)">
                  <!--          自定义数据-->
                  <template v-for="(diyData,diyIndex) in configObject.animate.animateElement[index].elementList" >
                    <div :key="diyIndex">
                      <a-form-item :label="$t(diyData.name)" >
                        <!--数字类型-->
                        <a-input type="number" @change="UpdateNodeData" v-model="configObject.animate.animateElement[index].elementList[diyIndex].value" :min="diyData.min" :max="diyData.max" :value="diyData.value" v-if="diyData.type==1"/>
                        <!--颜色类型-->
                        <a-input type="color" @change="UpdateNodeData" v-model="configObject.animate.animateElement[index].elementList[diyIndex].value" :value="diyData.value" v-if="diyData.type==2">
                          <a-tooltip slot="suffix" :title="$t('displayConfig.Properties.ClearColor')" @click="configObject.animate.animateElement[index].elementList[diyIndex].value='transparent';UpdateNodeData()">
                            <a-icon type="delete" style="color: rgba(0,0,0,.45)" />
                          </a-tooltip>
                        </a-input>
                        <!--字体选择-->
                        <a-select @change="UpdateNodeData" v-model="configObject.animate.animateElement[index].elementList[diyIndex].value" v-if="diyData.type==3">
                          <a-select-option v-for="options in fontFamilyOptions" :key="options" :value="options">
                            {{ options }}
                          </a-select-option>
                        </a-select>
                        <!--字符串-->
                        <a-input type="text" @change="UpdateNodeData" v-model="configObject.animate.animateElement[index].elementList[diyIndex].value" :value="diyData.value" v-if="diyData.type==4"/>
                        <!--图片类型-->
                        <div v-if="diyData.type==5">
                          <vue-hover-mask>
                            <!-- 默认插槽 -->
                            <img v-if="configObject.animate.animateElement[index].elementList[diyIndex].value!=''"  style="width: 200px;height:200px;cursor: pointer" :src="configObject.animate.animateElement[index].elementList[diyIndex].value" />
                            <div v-else :style="{width: '200px',height:'200px',cursor: 'pointer','background-color':'#F2F2F2'}"></div>
                            <!-- action插槽 -->
                            <template v-slot:action>
                              <span style="font-size: 14px" @click="showSystemImageModel(1,diyIndex)">{{$t('component.systemImageModel.selectImage')}}</span>
                              <a-divider type="vertical" />
                              <span  style="font-size: 14px" @click="configObject.animate.animateElement[index].elementList[diyIndex].value=''">{{$t('component.systemImageModel.delImage')}}</span>
                            </template>
                          </vue-hover-mask>
                        </div>
                        <!--枚举类型-->
                        <a-select v-model="configObject.animate.animateElement[index].elementList[diyIndex].value" v-if="diyData.type==6">
                          <a-select-option v-for="options in diyData.enumList" :key="options.value" :value="options.value" @change="UpdateNodeData">
                            {{ $t(options.option) }}
                          </a-select-option>
                        </a-select>
                        <!--浮点型-->
                        <a-input type="number" @change="UpdateNodeData" v-model="configObject.animate.animateElement[index].elementList[diyIndex].value" :step="0.1" :min="diyData.min" :max="diyData.max" :value="diyData.value" v-if="diyData.type==7"/>
                      </a-form-item>
                    </div>
                  </template>
                </div>
              </div>

              <a-form-item :label="$t('component.public.isExpression')" v-if="configObject.animate && configObject.animate.selected && !configObject.animate.selected.includes('Forbidden')&&!configObject.animate.selected.includes('animateMove')">
                <a-radio-group v-model="configObject.animate.isExpression" @change="onExpressionChange">
                  <a-radio :value="true">
                    {{$t('component.public.Enable')}}
                  </a-radio>
                  <a-radio :value="false">
                    {{$t('component.public.Forbidden')}}
                  </a-radio>
                </a-radio-group>
              </a-form-item>
            </a-form>
          </template>
          <div v-if="configObject.animate.isExpression&&configObject.animate.selected!='Forbidden'">
            <a-divider>{{$t('displayConfig.Properties.Condition')}}</a-divider>
            <template >
              <div class="ant-form ant-form-horizontal" style="margin-top:5px;border:#f0f0f0 solid 1px;">
                <a-form  :label-col="{ span: 7}" :wrapper-col="{ span: 16 }">
                  <a-form-item :label="$t('displayConfig.Properties.ComponentFromBandDevice')" v-if="configObject.animate.condition.bandType==1&&configObject.animate.condition.dataID!=''">
                    <a-checkbox disabled v-model="configObject.animate.condition.isBandDevice" @change="UpdateNodeData">
                      <label v-if="configObject.animate.condition.isBandDevice"> {{$t('displayConfig.Properties.ComponentIsBandDevice')}}</label>
                      <label v-else> {{$t('displayConfig.Properties.ComponentIsBandDeviceModel')}}</label>
                    </a-checkbox>
                  </a-form-item>

                  <a-form-item :label="$t('displayConfig.Properties.ComponentBandDevice')" v-if="configObject.animate.condition.isBandDevice&&configObject.animate.condition.bandType==1">
                    <a-tree-select @change="UpdateNodeData"
                        disabled
                        v-model="configObject.animate.condition.deviceSN"
                        style="width: 100%"
                        :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
                        :tree-data="deviceTreeData"
                        :replace-fields="{ value: 'key',title:'text'}"
                        placeholder="Please select"
                        tree-default-expand-all
                    >
                    </a-tree-select>
                  </a-form-item>

                  <a-form-item :label="$t('displayConfig.Properties.ComponentBandData')" v-if="configObject.animate.condition.bandType==1">
                    <a-input  @change="UpdateNodeData" v-model="configObject.animate.condition.dataName" >
                      <a-tooltip placement="top" slot="addonAfter">
                        <template slot="title">
                          <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
                        </template>
                        <icon-font @click="ShowDeviceDataModel(index,'animate')" type="icon-xuanzeshuju"  />
                      </a-tooltip>
                    </a-input>
                  </a-form-item>

                  <a-form-item :label="$t('component.public.Operator')">
                    <a-select v-model="configObject.animate.condition.operator" @change="UpdateNodeData">
                      <a-select-option value="==">
                        =
                      </a-select-option>
                      <a-select-option value=">">
                        >
                      </a-select-option>
                      <a-select-option value=">=">
                        >=
                      </a-select-option>
                      <a-select-option value="<">
                        &lt;
                      </a-select-option>
                      <a-select-option value="<=">
                        &lt;=
                      </a-select-option>
                      <a-select-option value="<>">
                        &&
                      </a-select-option>
                      <a-select-option value="<!>">
                        ||
                      </a-select-option>
                      <a-select-option value="!=">
                        !=
                      </a-select-option>
                    </a-select>
                  </a-form-item>

                  <a-form-item :label="$t('component.public.OperatorValue')" >
                    <a-input   v-model="configObject.animate.condition.OperatorValue" @change="UpdateNodeData" >
                    </a-input>
                  </a-form-item>

                  <a-form-item :label="$t('component.public.OperatorMaxValue')" v-if="configObject.animate.condition.operator=='<>'||configObject.animate.condition.operator=='<!>'">
                    <a-input   v-model="configObject.animate.condition.OperatorMaxValue" @change="UpdateNodeData">
                    </a-input>
                  </a-form-item>

                </a-form>

              </div>
            </template>
          </div>
        </a-tab-pane>
      </a-tabs>
    </template>
    <template v-if="isLayer">
      <a-tabs default-active-key="1" :tab-position="tabPosition" style="width:420px;max-width:600px">
        <a-tab-pane key="1" :tab="$t('displayConfig.Properties.TabHeaterPageAttr')">
          <a-form  :label-col="{ span: 7}" :wrapper-col="{ span: 16 }">
            <div>
              <a-form-item :label="$t('displayConfig.Properties.PageName')">
                <a-input v-model="LayerData.name" />
              </a-form-item>
            </div>

            <div>
              <a-form-item :label="$t('displayConfig.Properties.PageBackColor')">
                <a-input
                    type="color"
                    v-model="LayerData.layer.backColor"
                >
                  <a-tooltip slot="suffix" :title="$t('displayConfig.Properties.ClearColor')" @click="LayerData.layer.backColor='transparent';UpdateNodeData()">
                    <a-icon type="delete" style="color: rgba(0,0,0,.45)" />
                  </a-tooltip>
                </a-input>
              </a-form-item>
            </div>

            <div>
              <a-form-item :label="$t('displayConfig.Properties.PageBackImage')">
                <vue-hover-mask>
                  <!-- 默认插槽 -->
                  <img v-if="LayerData.layer.backgroundImage!=''"  style="width: 280px;height:200px;cursor: pointer" :src="LayerData.layer.backgroundImage" />
                  <div v-else :style="{width: '280px',height:'200px',cursor: 'pointer','background-color':LayerData.layer.backColor}"></div>
                  <!-- action插槽 -->
                  <template v-slot:action>
                    <span style="font-size: 20px" @click="showSystemImageModel(0)">{{$t('component.systemImageModel.selectImage')}}</span>
                    <a-divider type="vertical" />
                    <span  style="font-size: 20px" @click="LayerData.layer.backgroundImage=''">{{$t('component.systemImageModel.delImage')}}</span>
                  </template>
                </vue-hover-mask>
              </a-form-item>
            </div>

            <div>
              <a-form-item :label="$t('displayConfig.Properties.PageResolution')">
                <a-select
                    v-model="layerWH"
                >
                  <a-select-option v-for="options in whOptions" :key="options" :value="options">
                    {{ options }}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </div>
            <div v-if="layerWH == 'Custom'">
              <div>
                <a-form-item :label="$t('displayConfig.Properties.PageWith')">
                  <a-input
                      type="number"
                      suffix="px"
                      v-model="LayerData.layer.width"
                  />
                </a-form-item>
              </div>
              <div>
                <a-form-item :label="$t('displayConfig.Properties.PageHeight')">
                  <a-input
                      type="number"
                      suffix="px"
                      v-model="LayerData.layer.height"
                  />
                </a-form-item>
              </div>
            </div>
            <div>
              <a-form-item :label="$t('displayConfig.Properties.AutoSize')">
                <a-select
                    v-model="LayerData.layer.autoSize"
                >
                  <a-select-option value="0">False</a-select-option>
                  <a-select-option value="1">True</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.Padding')">
                <a-select
                    v-model="LayerData.layer.Padding"
                >
                  <a-select-option value="0">False</a-select-option>
                  <a-select-option value="1">True</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentAnimate')">
                <a-select v-model="LayerData.layer.animate" :allowClear="true">
                  <a-select-option :key="index" :value="item.name" v-for="(item, index) in animates">{{item.alias}}</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.virtuallyKey')">
                <a-select
                    v-model="LayerData.layer.virtuallyKey"
                >
                  <a-select-option value="0">False</a-select-option>
                  <a-select-option value="1">True</a-select-option>
                </a-select>
              </a-form-item>
            </div>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </template>
    <device-data-model @onSelectDataModel="onSelectData" ref="deviceDataModel"></device-data-model>
    <system-image-model @onSelectImage="onSelectImage" :networkImageUrl="imageComponentData" ref="systemImageModel"></system-image-model>
    <system-video-model @onSelectVideo="onSelectVideo" :networkVideo="videoComponentData" ref="systemVideoModel"></system-video-model>

    <a-modal :visible="codePopUpDialog"
             :title="$t('displayConfig.Properties.CodeEditDig')"
             :width="700"
             :height="500"
             @cancel="codePopUpDialog=false"
             v-drag-modal
             :destroyOnClose="true"
             :maskClosable="false"
             :maskStyle="{}"
             :mask="false">
      <div >
          <code-editor
              v-if="codePopUpDialog"
              :value="codePopUpValue"
              language="javascript"
              @input="changeCodeTextarea($event)"
          >
          </code-editor>
          <a-button v-if="isFullscreen" class="fullscreen-btn" @click="toggleEditorFullscreen">
            退出全屏
          </a-button>
      </div>
      <template slot="footer">
        <a-button key="submit" type="primary" @click="toggleEditorFullscreen">
          全屏编辑
        </a-button>
      </template>
    </a-modal>
  </div>
</template>

<script>
import store from "../../store";

import deviceDataModel from '../../components/deviceDataModel/deviceDataModel'
import systemImageModel from '../../components/systemImageModel/systemImageModel'
import systemVideoModel from '../../components/systemVideoModel/systemVideoModel'
import {getMonitorTree} from "@/services/device";
import {systemFontsList, systemRolesList} from "@/services/system";
import {vDragModal} from "@/utils/vmodalDrage"
import { mapActions,mapState, mapMutations, mapGetters } from 'vuex'
import VueHoverMask from "@/components/VueHoverMask/VueHoverMask"
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";
import {GetScriptList} from "@/services/ismscripts";
import { uuid } from 'vue-uuid';
import codeEditor from '@/components/CodeEditor/index'
import {setGroupList} from "@/store/ISM/actions";
import { cloneDeep } from 'lodash-es'
import {GetSQLReportTempletes} from "@/services/SqlReportTemplete";
export default {
  name: 'ISMProperties',
  i18n: require('../../i18n/language'),
  components: {
    deviceDataModel,
    systemImageModel,
    systemVideoModel,
    VueHoverMask,
    codeEditor,
  },
  data () {
    return {
      isFullscreen:false,
      displayUUID:"",
      codePopUpValue:"",
      codePopUpValueIndex:0,
      codePopUpDialog:false,
      ScriptDataSource:[],
      SQLReportList:[],
      RoleList:[],
      BandType:"",
      actionIndex: 0,
      UpdateNodeDataFlag:true,
      configurationModel:[],
      displayPageList:new Map,
      imageComponentData:"",
      videoComponentData: {},
      deviceTreeData:[],
      selectBandDataIndex:0,
      selectedIndex:[0,1],
      selectBandDiyDataIndex:0,
      selectImageType:0,
      customStyle: 'background: #fff;border-radius: 4px;margin-bottom: 24px;border: 0;overflow: hidden',
      tabPosition:"top",
      tabIndex: 0,
      StatusColumns:[
        {
          slotName: this.$t("configComponent.status.StatusValue"),
          scopedSlots: {  customRender: 'StatusValue' ,title:this.$t("configComponent.status.StatusValue")},
          width: '10%',
          align:"center",
          dataIndex: 'value'
        },
        {
          slotName: this.$t("configComponent.status.StatusOpt"),
          scopedSlots: {  customRender: 'StatusOpt' ,title:this.$t("configComponent.status.StatusOpt")},
          width: '20%',
          align:"center",
          dataIndex: 'StatusOpt'
        },
        {
          slotName: this.$t("configComponent.status.StatusValue2"),
          scopedSlots: {  customRender: 'StatusValue2' ,title:this.$t("configComponent.status.StatusValue2")},
          width: '10%',
          align:"center",
          dataIndex: 'value2'
        },
        {
          slotName: this.$t("configComponent.status.StatusText"),
          scopedSlots: {  customRender: 'StatusText' ,title:this.$t("configComponent.status.StatusText")},
          width: '15%',
          align:"center",
          dataIndex: 'Text'
        },
        {
          slotName:this.$t("configComponent.status.StatusTextColor"),
          scopedSlots: {  customRender: 'StatusTextColor' ,title:this.$t("configComponent.status.StatusTextColor") },
          width: '8%',
          align:"center",
          dataIndex: 'TextColor',
        },
        {
          slotName:this.$t("configComponent.status.OptStatus"),
          scopedSlots: {  customRender: 'OptStatus'  ,title:this.$t("configComponent.status.OptStatus")},
          width: '10%',
          align:"center",
        }
      ],
      LineStatusColumns:[
        {
          slotName: this.$t("configComponent.status.StatusValue"),
          scopedSlots: {  customRender: 'StatusValue' ,title:this.$t("configComponent.status.StatusValue")},
          width: '10%',
          align:"center",
          dataIndex: 'value'
        },
        {
          slotName: this.$t("configComponent.status.StatusOpt"),
          scopedSlots: {  customRender: 'StatusOpt' ,title:this.$t("configComponent.status.StatusOpt")},
          width: '7%',
          align:"center",
          dataIndex: 'StatusOpt'
        },
        {
          slotName: this.$t("configComponent.status.StatusValue2"),
          scopedSlots: {  customRender: 'StatusValue2' ,title:this.$t("configComponent.status.StatusValue2")},
          width: '10%',
          align:"center",
          dataIndex: 'value2'
        },
        {
          slotName: this.$t("configComponent.status.LineStatusBlink"),
          scopedSlots: {  customRender: 'LineStatusBlink' ,title:this.$t("configComponent.status.LineStatusBlink")},
          width: '7%',
          align:"center",
          dataIndex: 'LineStatusBlink'
        },
        {
          slotName:this.$t("configComponent.status.LineStatusBlinkSpeed"),
          scopedSlots: {  customRender: 'LineStatusBlinkSpeed' ,title:this.$t("configComponent.status.LineStatusBlinkSpeed") },
          width: '12%',
          align:"center",
          dataIndex: 'BlinkSpeed',
        },
        {
          slotName:this.$t("configComponent.status.StatusTextColor"),
          scopedSlots: {  customRender: 'StatusTextColor' ,title:this.$t("configComponent.status.StatusTextColor") },
          width: '5%',
          align:"center",
          dataIndex: 'TextColor',
        },
        {
          slotName:this.$t("configComponent.status.OptStatus"),
          scopedSlots: {  customRender: 'OptStatus'  ,title:this.$t("configComponent.status.OptStatus")},
          width: '10%',
          align:"center",
        }
      ],
      ComBoxStatusColumns:[
        {
          slotName: this.$t("configComponent.status.ComBoxValue"),
          scopedSlots: {  customRender: 'ComBoxValue' ,title:this.$t("configComponent.status.ComBoxValue")},
          width: '20%',
          align:"center",
          dataIndex: 'value'
        },
        {
          slotName: this.$t("configComponent.status.StatusText"),
          scopedSlots: {  customRender: 'StatusText' ,title:this.$t("configComponent.status.StatusText")},
          width: '30%',
          align:"center",
          dataIndex: 'Text'
        },
        {
          slotName:this.$t("configComponent.status.OptStatus"),
          scopedSlots: {  customRender: 'OptStatus'  ,title:this.$t("configComponent.status.OptStatus")},
          width: '10%',
          align:"center",
        }
      ],
      DropMenuColumns:[
        {
          slotName: this.$t("configComponent.status.MenuText"),
          scopedSlots: {  customRender: 'MenuText' ,title:this.$t("configComponent.status.MenuText")},
          width: '20%',
          align:"center",
          dataIndex: 'MenuName'
        },
        {
          slotName: this.$t("configComponent.status.AppListName"),
          scopedSlots: {  customRender: 'APPListName' ,title:this.$t("configComponent.status.AppListName")},
          width: '20%',
          align:"center",
          dataIndex: 'DisplayUUID'
        },
        {
          slotName: this.$t("configComponent.status.MenuPageName"),
          scopedSlots: {  customRender: 'MenuPageName' ,title:this.$t("configComponent.status.MenuPageName")},
          width: '20%',
          align:"center",
          dataIndex: 'PageID'
        },
        {
          slotName: this.$t("configComponent.status.MenuPagePop"),
          scopedSlots: {  customRender: 'MenuPagePop' ,title:this.$t("configComponent.status.MenuPagePop")},
          width: '10%',
          align:"center",
          dataIndex: 'IsPopUp'
        },
        {
          slotName:this.$t("configComponent.status.OptStatus"),
          scopedSlots: {  customRender: 'OptStatus'  ,title:this.$t("configComponent.status.OptStatus")},
          width: '10%',
          align:"center",
        }
      ],
      ImagesStatusColumns:[
        {
          slotName: this.$t("configComponent.status.StatusValue"),
          scopedSlots: {  customRender: 'StatusValue' ,title:this.$t("configComponent.status.StatusValue")},
          width: '10%',
          align:"center",
          dataIndex: 'value'
        },
        {
          slotName: this.$t("configComponent.status.StatusOpt"),
          scopedSlots: {  customRender: 'StatusOpt' ,title:this.$t("configComponent.status.StatusOpt")},
          width: '17%',
          align:"center",
          dataIndex: 'StatusOpt'
        },
        {
          slotName: this.$t("configComponent.status.StatusValue2"),
          scopedSlots: {  customRender: 'StatusValue2' ,title:this.$t("configComponent.status.StatusValue2")},
          width: '10%',
          align:"center",
          dataIndex: 'value2'
        },
        {
          slotName: this.$t("configComponent.status.ImageTitle"),
          scopedSlots: {  customRender: 'ImageTitle' ,title:this.$t("configComponent.status.ImageTitle")},
          width: '20%',
          align:"center",
          dataIndex: 'Image'
        },
        {
          slotName:this.$t("configComponent.status.OptStatus"),
          scopedSlots: {  customRender: 'OptStatus'  ,title:this.$t("configComponent.status.OptStatus")},
          width: '10%',
          align:"center",
        }
      ],
      fontFamilyOptions: [
        "黑体",
        "楷体",
        "隶书",
        "宋体",
        "数字字体-1",
        "数字字体-2",
        "数字字体-3",
        "数字字体-4",
        "数字字体-5",
        "数字字体-6",
        "数字字体-7",
        "数字字体-8",
        "数字字体-9",
        "数字字体-10",
        "数字字体-11",
        "数字字体-12"
      ],
      textAlignOptions: ["left", "right", "center", "justify"],
      borderStyleOptions: ["solid", "dashed", "dotted"],
      whOptions: ['1024x768', '1366x768', '1280x800', '1440x900', '1600x900', '1920x1080', 'Custom'],
      layerWHTemp: '',
    }
  },
  watch: {},
  computed: {
    layerWH: {
      get: function(){
        return this.getLayerWH()
      },
      set: function(val){
        this.setLayerWH(val)
      }
    },
    ...mapState({
      LayerData: state => store.state.ISMDisPlayEditorTool.LayerData,
      GroupList:state => store.state.ISMDisPlayEditorTool.GroupList,
      selectedComponents: state => store.state.ISMDisPlayEditorTool.selectedComponents,
      selectedComponentMap: state => store.state.ISMDisPlayEditorTool.selectedComponentMap,
      isLayer: state => store.state.ISMDisPlayEditorTool.selectedIsLayer,
      configObject: state => store.state.ISMDisPlayEditorTool.selectedComponent,
      PCPageList: state => store.state.ISMDisPlayEditorTool.PCPageList,
      PhonePageList: state => store.state.ISMDisPlayEditorTool.PhonePageList,
      selectPageUuid: state => store.state.ISMDisPlayEditorTool.selectPageUuid,
      animates: state => store.state.setting.componentAnimates,
      selectedNode: state => store.state.ISMDisPlayEditorTool.selectedNode,
      selectedNodePops: state => store.state.ISMDisPlayEditorTool.selectedNodePops,
      ISMCavasContainer:state => store.state.ISMDisPlayEditorTool.ISMCavasContainer,
    }),
    ...mapGetters('ISMDisPlayEditorTool', [
      'getGroupList'
    ]),
    animations () {
      let items = [];
      if (this.configObject.type == 'dashed') {
        items = (this.configObject.direction && this.configObject.direction == 'vertical') ?
            [{ label: '向上', value: 'up' }, { label: '向下', value: 'down' }] : [{ label: '向右', value: 'right' }, { label: '向左', value: 'left' }];
      }
      return items;
    }
  },
  created(){
    this.getMonitorTree()
    this.SystemRolesList()
    this.GetScriptList()
    this.GetSQLReportTemplates()
  },
  methods: {
    ...mapActions('ISMDisPlayEditorTool',[
      'getLayerDataStruct',
      'setGroupList',
    ]),
    ...mapMutations('ISMDisPlayEditorTool',[
      'execute',
    ]),
    toggleEditorFullscreen () {
      this.isFullscreen = !this.isFullscreen
      const editorEl = document.querySelector('.CodeMirror')
      if(this.isFullscreen) {
        editorEl.classList.add('editor-fullscreen')
      }
      else{
        editorEl.classList.remove('editor-fullscreen')
      }
    },
    UpdateNode(){
      const tdata = this.selectedNodePops.data
      const popsData = this.selectedNodePops
      let NodeAngle = popsData.angle
      let Visible = popsData.visible
      this.selectedNode.setZIndex(parseInt(tdata.detail.style.zIndex));
      this.selectedNode.rotate(parseInt(NodeAngle),{ absolute: true });
      this.selectedNode.setVisible(Visible,{ silent: true });
      this.selectedNode.position(parseInt(popsData.position.x), parseInt(popsData.position.y))
      this.selectedNode.resize(parseInt(popsData.size.width), parseInt(popsData.size.height))
      tdata.detail.style.position.x = popsData.position.x
      tdata.detail.style.position.y = popsData.position.y
      tdata.detail.style.position.w = popsData.size.width
      tdata.detail.style.position.h = popsData.size.height
      this.UpdateNodeDataFlag=!this.UpdateNodeDataFlag
      let nodeData = this.selectedNode.getData()
      const nodeDataNew = cloneDeep(nodeData)
      nodeDataNew.UpdateNodeFlag = this.UpdateNodeDataFlag
      nodeDataNew.detail = tdata.detail
      this.ISMCavasContainer.batchUpdate(() => {
        this.selectedNode.setData(nodeDataNew, {overwrite: true})
      })
    },
    UpdateNodeData(){
      const tdata = this.selectedNodePops.data
      const popsData = this.selectedNodePops
      const Visible = tdata.detail.style.visible==1?true:false
      this.selectedNode.setZIndex(parseInt(tdata.detail.style.zIndex));
      this.selectedNode.setVisible(Visible,{ silent: false });
      const setDeltail = this.configObject
      const isEdge = this.selectedNode.isEdge()
      if(isEdge)
      {
        const option = setDeltail
        let i =0
        let strokeWidth=4
        let backColor=""
        let MoveBrokenLineInterval = 3
        let strokeSpace = ""
        let strokeLength = ""
        let strokeColor = ""
        let fillOpacity=1
        let strokeOpacity=1
        let spinDirection=1
        let MoveBrokenLineConditionEnable=0
        let strokeBgWidth=25
        let strokeLinejoin="round"
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="strokeWidth")
          {
            strokeWidth=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="MoveBrokenLineInterval")
          {
            MoveBrokenLineInterval=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="strokeLength")
          {
            strokeLength=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeSpace")
          {
            strokeSpace=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="spinDirection")
          {
            spinDirection=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="MoveBrokenLineConditionEnable")
          {
            MoveBrokenLineConditionEnable=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="strokeBgWidth")
          {
            strokeBgWidth=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="strokeLinejoin")
          {
            const strokeLinejoinType=parseInt(option.style.diy[i].value)
            if(strokeLinejoinType==0)
            {
              strokeLinejoin = "miter"
            }
            else if(strokeLinejoinType==2)
            {
              strokeLinejoin = "bevel"
            }
            else if(strokeLinejoinType==1)
            {
              strokeLinejoin = "round"
            }
          }
        }
        strokeColor = option.style.foreColor
        backColor   = option.style.backColor
        let animation = ""
        let strokeDasharray = strokeLength+" "+strokeSpace
        if(MoveBrokenLineConditionEnable)
        {
          animation=""
        }
        else if(spinDirection==0)
        {
          animation = 'ant-line-forward '+MoveBrokenLineInterval+'s infinite linear'
        }
        else{
          animation = 'ant-line-inverse '+MoveBrokenLineInterval+'s infinite linear'
        }
        const propdata = this.selectedNode.prop()
        if(typeof propdata.connectType!="undefined"&&propdata.connectType=="connected-edge")
        {
          let strokeLineType=1
          let strokeLineMarkerIsShow=0
          let strokeLineMarker = "classic"
          let strokeLineMarkerWidth = 10
          let strokeLineMarkerHeight = 10
          let strokeLineMarkerColor = "#06ad8e"
          for( i=0;i<option.style.diy.length;i++) {
            if (option.style.diy[i].key == "strokeLineType") {
              strokeLineType = parseInt(option.style.diy[i].value)
              if(strokeLineType==1)
              {
                strokeDasharray=""
              }
            }
            else if (option.style.diy[i].key == "strokeLineMarkerStyle") {
              const strokeLineMarkerStyleTemp = parseInt(option.style.diy[i].value)
              if(strokeLineMarkerStyleTemp==0)
              {
                strokeLineMarker = "classic"
              }
              else if(strokeLineMarkerStyleTemp==1)
              {
                strokeLineMarker = "diamond"
              }
              else if(strokeLineMarkerStyleTemp==2)
              {
                strokeLineMarker = "cross"
              }
              else if(strokeLineMarkerStyleTemp==3)
              {
                strokeLineMarker = "circle"
              }
              else if(strokeLineMarkerStyleTemp==4){
                strokeLineMarker = "circlePlus"
              }
              else if(strokeLineMarkerStyleTemp==5){
                  strokeLineMarker = "ellipse"
              }
            }
            else if (option.style.diy[i].key == "strokeLineMarker") {
              strokeLineMarkerIsShow = parseInt(option.style.diy[i].value)
            }
            else if (option.style.diy[i].key == "strokeLineMarkerColor") {
              strokeLineMarkerColor = option.style.diy[i].value
            }
            else if (option.style.diy[i].key == "strokeLineMarkerWidth") {
              strokeLineMarkerWidth = parseInt(option.style.diy[i].value)
            }
            else if (option.style.diy[i].key == "strokeLineMarkerHeight") {
              strokeLineMarkerHeight = parseInt(option.style.diy[i].value)
            }
          }
          if(strokeLineMarkerIsShow==0)
          {
            strokeLineMarker=""
          }
          this.selectedNode.setAttrs({
            line:{
              strokeOpacity: setDeltail.style.opacity,
              style: {
                strokeDasharray: strokeDasharray,
                stroke: strokeColor != "" ? strokeColor : "none",
                strokeWidth: strokeWidth,
                strokeLinejoin: strokeLinejoin,
                animation: animation,
              },
              sourceMarker:{
                name:strokeLineMarker,
                width:strokeLineMarkerWidth,
                fill:strokeLineMarkerColor,
                stroke:strokeLineMarkerColor,
                height:strokeLineMarkerHeight
              },
              targetMarker:{
                fill:strokeLineMarkerColor,
                stroke:strokeLineMarkerColor,
                name:strokeLineMarker,
                width:strokeLineMarkerWidth,
                height:strokeLineMarkerHeight
              }
            },
            wrap: {
              strokeOpacity: setDeltail.style.opacity,
              stroke: backColor!=""?backColor:"none",
              strokeWidth: strokeBgWidth,
              strokeLinejoin: strokeLinejoin,
            },
          });
        }
        else
        {
          this.selectedNode.setAttrs({
            line:{
              strokeOpacity: setDeltail.style.opacity,
              style: {
                strokeDasharray: strokeDasharray,
                stroke: strokeColor != "" ? strokeColor : "none",
                strokeWidth: strokeWidth,
                strokeLinejoin: strokeLinejoin,
                animation: animation,
              }
            },
            wrap: {
              strokeOpacity: setDeltail.style.opacity,
              stroke: backColor!=""?backColor:"none",
              strokeWidth: strokeBgWidth,
              strokeLinejoin: strokeLinejoin,
            },
          });
        }
      }
      else {
        this.configObject.style.position.x = popsData.position.x
        this.configObject.style.position.y = popsData.position.y
        this.configObject.style.position.w = popsData.size.width
        this.configObject.style.position.h = popsData.size.height
        this.UpdateNodeDataFlag = !this.UpdateNodeDataFlag
        let nodeData = this.selectedNode.getData()
        const nodeDataNew = cloneDeep(nodeData)
        nodeDataNew.UpdateNodeFlag = new Date().getTime()
        nodeDataNew.detail =setDeltail
        this.ISMCavasContainer.batchUpdate(() => {
          this.selectedNode.setData(nodeDataNew, {overwrite: true})
        })
      }
      this.setGroupList()
    },
    CodeDbClick(index) {
      this.codePopUpDialog = true
      this.codePopUpValue = this.configObject.style.diy[index].value
      this.codePopUpValueIndex = index
      this.UpdateNodeData()
    },
    chargeGroupIDVisible(value){
      let Visible = value==1?true:false
      this.UpdateNodeData()
      let getChildren = this.selectedNode.getChildren()
      if(getChildren!=null)
      {
        getChildren.forEach(child => {
          child.setVisible(Visible)
          let NodeInfoData = child.getData()
          NodeInfoData.UpdateNodeFlag = new Date().getTime()
          NodeInfoData.detail.style.visible = Visible
          child.setData(NodeInfoData, {overwrite: true})
        })
      }
    },
    chargeGroupID(ID){
      if(typeof ID=="undefined")
      {
        for(let key in this.selectedComponentMap) {
          for (let k = 0,componentsLen=this.LayerData.components.length; k < componentsLen; k++) {
            if((this.LayerData.components[k].identifier==key))
            {
              this.LayerData.components[k].groupID=""
              this.LayerData.components[k].GroupName=""
              break
            }
          }
          break
        }
      }
      else
      {
        let getGroupList = this.getGroupList
        let item = getGroupList.find(item => {
          return item.ID == ID;
        });
        if(typeof item=="undefined")
        {
          return
        }
        for(let key in this.selectedComponentMap) {
          for (let k = 0,componentsLen=this.LayerData.components.length; k < componentsLen; k++) {
            if((this.LayerData.components[k].identifier==key))
            {
              this.LayerData.components[k].groupID=ID
              this.LayerData.components[k].GroupName=item.Name
            }
          }
        }
      }
      this.setGroupList()
    },
    chargeGroupName(input){
      for(let key in this.selectedComponentMap) {
        const groupID = this.selectedComponentMap[key].groupID
        for (let k = 0,componentsLen=this.LayerData.components.length; k < componentsLen; k++) {
          if((typeof this.LayerData.components[k].groupID!="undefined")&&(this.LayerData.components[k].groupID!="")&&(this.LayerData.components[k].groupID == groupID))
          {
            this.LayerData.components[k].GroupName= input.target.value
          }
        }
        break
      }
      this.setGroupList()
    },
    changeCodeTextarea(val) {
      this.configObject.style.diy[this.codePopUpValueIndex].value = val
      this.UpdateNodeData()
    },
    GetScriptList(){
      let _t = this
      this.ScriptDataSource = []
      GetScriptList().then(function (res){
        _t.refIconLoading = false
        if (res.data.code == 200) {
          _t.ScriptDataSource = res.data.list
          _t.addVisible = false;
        }
        else if (res.data.code == 2001) {
          _t.$message.error(_t.$t('displayModel.ModelExist'), 3)
        }
        else if (res.data.code == 2003) {
          _t.$message.error(_t.$t('displayModel.AddModelFailed'), 3)
        }
      })
    },
    SelectTreeDevice(value, node, extray){
      let info = extray.selectedNodes[0].data.props.dataRef
      if (typeof this.configObject.action[this.actionIndex].DeviceView=='undefined'){
        this.configObject.action[this.actionIndex].DeviceView={
          key:"",
          showUUID:"",
          showPageUUID:"",
          type:"",
          selectKey:"",
          isPopUp:false,

        }
      }
      this.configObject.action[this.actionIndex].DeviceView.key=info.key
      this.configObject.action[this.actionIndex].DeviceView.showUUID=info.value.configUid
      this.configObject.action[this.actionIndex].DeviceView.showPageUUID=info.value.PageUUID
      this.configObject.action[this.actionIndex].DeviceView.type=info.value.type
      this.configObject.action[this.actionIndex].DeviceView.selectKey=info.key
      this.UpdateNodeData()
    },
    SystemRolesList(){
      let _t = this
      _t.RoleList = []
      this.messageShowLoad = true
      systemRolesList().then(function (res){
        if(res.data.code==0)
        {
          _t.RoleList = res.data.list
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
      })

    },
    GetDisplayPage(uuid){
      let params={
        muid:uuid
      }
      let _t = this
      getDisplayModelLayerData(params).then(function (res){
        if(res.data.code==0)
        {
          let pageLayer = res.data.layer
          if(pageLayer.length>0)
          {
            let displayArray=[]
            for(let i=0;i<pageLayer.length;i++)
            {
              if (pageLayer[i].IsLogin==1)
              {
                continue
              }
              let pageInfo = {}
              pageInfo.label = pageLayer[i].PageName
              pageInfo.value = pageLayer[i].PageId
              pageInfo.pageType = pageLayer[i].PageType
              pageInfo.pageModelUuid = pageLayer[i].modelId
              displayArray.push(pageInfo)
            }
            _t.displayPageList.set(uuid,displayArray)
          }
        }
      })
    },
    getConfigurationModel(){
      this.configurationModel=[]
      let _t = this
      ;[1, 2].forEach(function(displayType) {
        const params = {
          DisplayType: displayType
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
              tableData.displayType = displayType
              _t.configurationModel.push(tableData)
              tableData={}
              _t.GetDisplayPage(res.data.list[i].displayUid)
            }
          }
        })
      })
    },
    onExpressionChange(e){
      this.UpdateNodeData()
    },
    GetSystemFonts(){
      systemFontsList().then(function (res){

      })
    },
    addStatus(){
      let status={
        "Text":"Text",
        "Image":"",
        "value":0
      }
      this.configObject.active[0].condition.StatusList.push(status)
      this.UpdateNodeData()
    },
    addMenu(){
      let status={
        "MenuName":"菜单",
        "IsPopUp":false,
        "DisPlayID":this.$route.params.uid,
        "PageID":""
      }
      this.configObject.active[0].condition.PageList.push(status)
      this.UpdateNodeData()
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
    getLayerWH () {
      let wh = null
      if (!this.LayerData.layer.width || !this.LayerData.layer.height) {
        this.LayerData.layer.width = 1600;
        this.LayerData.layer.height = 900;
      }
      if (this.layerWHTemp == '') {
        wh = this.LayerData.layer.width + 'x' + this.LayerData.layer.height;
        if (this.whOptions.indexOf(wh, 0) == -1) {
          this.layerWHTemp = 'Custom';
        } else {
          this.layerWHTemp = wh;
        }
      } else {
        wh = this.LayerData.layer.width + 'x' + this.LayerData.layer.height;
        if (this.whOptions.indexOf(wh, 0) == -1) {
          this.layerWHTemp = 'Custom';
        }
      }
      return this.layerWHTemp;
    },
    setLayerWH(val) {
      this.layerWHTemp = val;
      if (val == 'Custom') {
        console.log(val)
      } else {
        var wh = val.split('x');
        this.LayerData.layer.width = parseInt(wh[0]);
        this.LayerData.layer.height = parseInt(wh[1]);
      }
    },
    initPage (configData) {
      this.configData = configData;
    },
    changeTab (tabIndex) {
      this.tabIndex = tabIndex;
    },
    generateTargetComponentOptions () {
      let options = [];
      this.LayerData.components.forEach(component => {
        // if (this.configObject.identifier != component.identifier) {
        options.push({
          label: component.name || component.type,
          value: component.identifier
        });
        // }
      });
      return options;
    },
    generateTargetPage (uuid) {
      return this.displayPageList.get(uuid) || []
    },
    get2DConfigurationModel() {
      return this.configurationModel.filter(item => item.displayType !== 2)
    },
    ensureLinkExternalPage(link) {
      if(!link.ExternalPage) {
        this.$set(link, 'ExternalPage', {
          displayUUID: "",
          pageUUID: ""
        })
      }
      return link.ExternalPage
    },
    buildExternalPageUrl(displayUUID, pageUUID) {
      if(!displayUUID || !pageUUID) {
        return ""
      }
      const basePath = window.location.origin + window.location.pathname
      return `${basePath}#/AppRun/${displayUUID}?pageId=${pageUUID}`
    },
    changeLinkExternalSource(link, value) {
      this.$set(link, 'ExternalSource', value)
      if(value === 'Page') {
        const externalPage = this.ensureLinkExternalPage(link)
        link.External = this.buildExternalPageUrl(externalPage.displayUUID, externalPage.pageUUID)
      }
      this.UpdateNodeData()
    },
    changeLinkExternalDisplay(link) {
      const externalPage = this.ensureLinkExternalPage(link)
      externalPage.pageUUID = ""
      link.External = ""
      this.UpdateNodeData()
    },
    changeLinkExternalPage(link) {
      const externalPage = this.ensureLinkExternalPage(link)
      link.External = this.buildExternalPageUrl(externalPage.displayUUID, externalPage.pageUUID)
      this.UpdateNodeData()
    },
    getLinkInsideDisplayType(uuid) {
      const target = this.configurationModel.find(item => item.uuid === uuid)
      return target ? target.displayType : 1
    },
    changeLinkInsideDisplay(link) {
      if(!link || !link.Inside) {
        this.UpdateNodeData()
        return
      }
      const displayType = this.getLinkInsideDisplayType(link.Inside.displayUUID)
      link.Inside.displayType = displayType
      if(displayType === 2) {
        link.Inside.pageUUID = ""
        link.isPopUp = false
      }
      this.UpdateNodeData()
    },
    onSelectData(selectData) {

      if(this.BandType=="data") {
        this.configObject.dataBind[this.selectBandDataIndex].DeviceName=selectData.DeviceName
        this.configObject.dataBind[this.selectBandDataIndex].isBandDevice = selectData.IsDevice
        this.configObject.dataBind[this.selectBandDataIndex].deviceSN = selectData.DeviceSN
        this.configObject.dataBind[this.selectBandDataIndex].dataName = selectData.name
        this.configObject.dataBind[this.selectBandDataIndex].dataID = selectData.uuid
      }
      else if(this.BandType=="active") {
        this.configObject.active[this.selectBandDataIndex].DeviceName=selectData.DeviceName
        this.configObject.active[this.selectBandDataIndex].condition.isBandDevice = selectData.IsDevice
        this.configObject.active[this.selectBandDataIndex].condition.deviceSN = selectData.DeviceSN
        this.configObject.active[this.selectBandDataIndex].condition.dataName = selectData.name
        this.configObject.active[this.selectBandDataIndex].condition.dataID = selectData.uuid
        this.configObject.active[this.selectBandDataIndex].condition.dataUnit = selectData.unit
        this.configObject.active[this.selectBandDataIndex].condition.DeviceName = selectData.DeviceName
      }
      else if(this.BandType=="animate") {
        this.configObject.animate.condition.DeviceName=selectData.DeviceName
        this.configObject.animate.condition.isBandDevice = selectData.IsDevice
        this.configObject.animate.condition.deviceSN = selectData.DeviceSN
        this.configObject.animate.condition.dataName = selectData.name
        this.configObject.animate.condition.dataID = selectData.uuid
      }
      else if(this.BandType=="animateMoveX") {
        this.configObject.animate.move.x.DeviceName=selectData.DeviceName
        this.configObject.animate.move.x.isBandDevice = selectData.IsDevice
        this.configObject.animate.move.x.deviceSN = selectData.DeviceSN
        this.configObject.animate.move.x.dataName = selectData.name
        this.configObject.animate.move.x.dataID = selectData.uuid
      }
      else if(this.BandType=="animateMoveY") {
        this.configObject.animate.move.y.DeviceName=selectData.DeviceName
        this.configObject.animate.move.y.isBandDevice = selectData.IsDevice
        this.configObject.animate.move.y.deviceSN = selectData.DeviceSN
        this.configObject.animate.move.y.dataName = selectData.name
        this.configObject.animate.move.y.dataID = selectData.uuid
      }
      else if(this.BandType=="setValue") {
        this.configObject.action[this.actionIndex].setValue[this.selectBandDataIndex].DeviceName=selectData.DeviceName
        this.configObject.action[this.actionIndex].setValue[this.selectBandDataIndex].isBandDevice = selectData.IsDevice
        this.configObject.action[this.actionIndex].setValue[this.selectBandDataIndex].deviceSN = selectData.DeviceSN
        this.configObject.action[this.actionIndex].setValue[this.selectBandDataIndex].dataName = selectData.name
        this.configObject.action[this.actionIndex].setValue[this.selectBandDataIndex].dataID = selectData.uuid
      }
      this.UpdateNodeData()
    },
    ShowDeviceDataModel(index,type){
      this.selectBandDataIndex = index;
      this.BandType = type
      this.$refs.deviceDataModel.showDataModal()
    },
    changeAction(v,index){
      if(v=='SetValue')
      {
        let setValue = {
          deviceSN:"",
          IsManual:false,
          AutoSetValue:"",
          isBandDevice:false,
          dataID: "",
          dataName: "",
        }
        if((typeof this.configObject.action[index].setValue=="undefined"))
        {
          this.configObject.action[index].setValue =[
            {
              deviceSN:"",
              IsManual:false,
              AutoSetValue:"",
              isBandDevice:false,
              dataID: "",
              dataName: "",
            }
          ]
        }
        if(this.configObject.action[index].setValue.length==0)
        {
          this.configObject.action[index].setValue.push(setValue)
        }
      }
      else if(v=='Animation')
      {
        this.configObject.animate.isExpression = true
      }
      else if(v=='visible')
      {
        this.setGroupList()
      }
      if (typeof this.configObject.action[index].DeviceView=='undefined'){
        this.configObject.action[index].DeviceView={
          key:"",
          showUUID:"",
          showPageUUID:"",
          type:"",
          selectKey:"",
          isPopUp:false,

        }
      }

      if (typeof this.configObject.action[index].link=='undefined'){
        this.configObject.action[index].link={
          linkType:"Inside",
          Inside:{
            displayUUID:"",
            pageUUID:"",
            displayType:1,
          },
          isPopUp:false,
          External:"",
          ExternalSource:"Url",
          ExternalPage:{
            displayUUID:"",
            pageUUID:""
          },
          OpenExternalType:"new"
        }
      }
      if (typeof this.configObject.action[index].RestApi=='undefined'){
        this.configObject.action[index].RestApi={
          Name:"",
          IsSystem:"1",
          Type:"Post",
          Url:"",
          Params:"{}",
        }
      }
      if (typeof this.configObject.action[index].setValue=='undefined'){
        this.configObject.action[index].setValue=[
          {
            deviceSN:"",
            IsManual:false,
            AutoSetValue:"",
            isBandDevice:false,
            dataID: "",
            dataName: "",
          }
        ]
      }
      if (typeof this.configObject.action[index].SysScript=='undefined'){
        this.configObject.action[index].ScriptList=[]
      }
      if (typeof this.configObject.action[index].Animation=='undefined'){
        this.configObject.animate.isExpression = true
        this.configObject.action[index].animationStatus=""
      }
      this.UpdateNodeData()
    },
    addBindSetValue(index){
      let setValue = {
        deviceSN:"",
        IsManual:false,
        AutoSetValue:"",
        isBandDevice:false,
        dataID: "",
        dataName: "",
      }
      if((typeof this.configObject.action[index].setValue=="undefined"))
      {
        this.configObject.action[index].setValue =[
          {
            deviceSN:"",
            IsManual:false,
            AutoSetValue:"",
            isBandDevice:false,
            dataID: "",
            dataName: "",
          }
        ]
      }
      else
      {
        this.configObject.action[index].setValue.push(setValue)
      }
      this.UpdateNodeData()
    },
    addBindActive(){
      let active = {
        id:uuid.v1(),
        name:"displayConfig.Properties.ComponentBandData",
        result:"",
        isExpression:false,
        condition:{
          deviceSN:"",
          selectVideoType:0,
          isBandDevice:false,
          bandType:1,
          dataID: "",
          dataUnit: "",
          dataName: "",
          operator:"",
          OperatorValue:"",
          OperatorMaxValue:"",
          actionAuth: [],
          actionVoice: "",
          SetPassword: "",
          StatusList: [],
          PageList: []
        },
      }
      this.configObject.active.push(active)
      this.UpdateNodeData()
    },
    addBindAction(){
      let action = {
        type: 'click',
        action: 'SetValue',
        actionAuth: [],
        ActionPassword: "",
        actionVoice: "",
        link:{
          linkType:"Inside",
          Inside:{
            displayUUID:"",
            pageUUID:"",
            displayType:1,
          },
          isPopUp:false,
          autoClose:false,
          External:"",
          ExternalSource:"Url",
          ExternalPage:{
            displayUUID:"",
            pageUUID:""
          },
          OpenExternalType:"new"
        },
        setValue:[
          {
            deviceSN:"",
            IsManual:false,
            AutoSetValue:"",
            isBandDevice:false,
            dataID: "",
            dataName: "",
          }
        ],
        RestApi:{
          Name:"",
          IsSystem:"1",
          Type:"Post",
          Url:"",
          Params:"{}",
        },
        ScriptList:[],
        DeviceView:{
          key:"",
          showUUID:"",
          showPageUUID:"",
          type:"",
          selectKey:"",
          isPopUp:false,
        },
        animationStatus:"",
        showItems: [],
        hideItems: [],
        SetDelay:"1000",
        actionConfirm:false
      };
      this.configObject.action.push(action)
      this.UpdateNodeData()
    },
    delBindAction(index){
      let _t = this
      this.$confirm({
        content: _t.$t('displayConfig.Properties.ComponentDelBandActionBtnTips'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          _t.configObject.action.splice(index, 1);
          _t.UpdateNodeData()
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });
    },
    delBindActive(index){
      let _t = this
      this.$confirm({
        content: _t.$t('displayConfig.Properties.ComponentDelBandActiveBtnTips'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          _t.configObject.active.splice(index, 1);
          _t.UpdateNodeData()
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });
    },
    delBindSetValue(index,setValueIndex){
      let _t = this
      this.$confirm({
        content: _t.$t('displayConfig.Properties.DeleteSetValueTips'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          _t.configObject.action[index].setValue.splice(setValueIndex, 1);
          _t.UpdateNodeData()
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });
    },
    showSystemImageModel(selectType,bindIndex,showType){
      this.$refs.systemImageModel.showModal(showType)
      this.selectImageType=selectType
      this.selectBandDiyDataIndex=bindIndex
      if(this.selectImageType==0)
      {
        this.imageComponentData = this.LayerData.layer.backgroundImage
      }
      else if(this.selectImageType==1)
      {
        this.imageComponentData = this.configObject.style.diy[this.selectBandDiyDataIndex].value
      }
      else if(this.selectImageType==2)
      {
        this.imageComponentData = this.configObject.active[0].condition.StatusList[bindIndex].Image
      }
      this.UpdateNodeData()
    },
    showSystemVideoModel(selectType,bindIndex){
      this.$refs.systemVideoModel.showModal()
      this.selectBandDiyDataIndex=bindIndex
      if(this.configObject.style.diy[this.selectBandDiyDataIndex].value.type==0)
      {
        this.videoComponentData = this.configObject.style.diy[this.selectBandDiyDataIndex].value
      }
      else
      {
        this.videoComponentData={}
      }
      this.UpdateNodeData()
    },
    onSelectImage(url){
      if(this.selectImageType==0)
      {
        this.LayerData.layer.backgroundImage=url
      }
      else if(this.selectImageType==1)
      {
        this.configObject.style.diy[this.selectBandDiyDataIndex].value=url
      }
      else if(this.selectImageType==2)
      {
        this.configObject.active[0].condition.StatusList[this.selectBandDiyDataIndex].Image=url
      }
      this.UpdateNodeData()
    },
    onSelectVideo(token){
      this.configObject.style.diy[this.selectBandDiyDataIndex].value=token
      this.UpdateNodeData()
    },
    GetSQLReportTemplates(){
      let _t = this
      this.SQLReportList=[]
      GetSQLReportTempletes().then(function (res){
        if (res.data.code == 0) {
          if(res.data.list==null)
          {
            _t.SQLReportList=[]
          }
          else
          {
            _t.SQLReportList = res.data.list
          }
        }
      }).finally(function (error) {

      })
    },
  },
  mounted () {
    this.getConfigurationModel()
  }
}
</script>

<style scoped>
::v-deep .fullscreen-btn{
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 9999;
}
::v-deep .editor-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  height: auto;
  z-index: 9000;
}
::v-deep .ant-modal-mask {
  background: transparent;
  pointer-events: none;
}
::v-deep .ant-modal-wrap {
  pointer-events: none;
}
::v-deep .ant-modal {
  pointer-events: all;
}
::v-deep .ant-table-tbody > tr > td {
  padding: 1px 1px;
  overflow-wrap: break-word;
}
::v-deep .ant-modal-body {
  padding: 1px;
}
::v-deep .ant-form-item {
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
//background: #f8f8f8;
//border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
::v-deep .ant-tabs-bar{
  margin: 0 0 3px 0;
}
::v-deep .ant-tabs-nav-wrap{
  background-color: #ffffff;
  background: -webkit-gradient(linear,left top, left bottom,color-stop(0, #ffffff),to(#ffffff));
  background: linear-gradient(to bottom, #ffffff 0, #ffffff 100%);
  font-size: 14px;
  font-weight: bold;
  color: #0E2D5F;
}
::v-deep .ant-form-item{
  margin-bottom: 0px;
}
::v-deep .ant-collapse-content > .ant-collapse-content-box {
  padding: 5px;
}

::-webkit-scrollbar {
  /*滚动条整体样式*/
  width : 5px;  /*高宽分别对应横竖滚动条的尺寸*/
  height: 9px;
}
::-webkit-scrollbar-thumb {
  /*滚动条里面小方块*/
  border-radius   : 10px;
  background-color: skyblue;
  background-image: -webkit-linear-gradient(
      45deg,
      rgba(255, 255, 255, 0.2) 25%,
      transparent 25%,
      transparent 50%,
      rgba(255, 255, 255, 0.2) 50%,
      rgba(255, 255, 255, 0.2) 75%,
      transparent 75%,
      transparent
  );
}
::-webkit-scrollbar-track {
  /*滚动条里面轨道*/
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}
</style>

<template>
  <project-layout>
      <a-row>
        <a-col style="padding: 0px;" :xl="24" :lg="24" :md="24" :sm="24" :xs="24">
          <a-card class="project-list" :loading="loading" style="padding: 0px;min-height: 400px;" :bordered="false" :title="$t('project.MyProject')" >
            <a slot="extra"> <a-button type="primary" style="width: 132px;height: 30px;font-size: 15px;
    font-weight: 500;" @click="projectVisible=true;isEditStatus=0">
              {{$t('project.AddProject')}}
            </a-button>
<!--              <a-button type="primary" style="width: 132px;height: 30px;font-size: 15px;margin-left: 20px;-->
<!--    font-weight: 500;" @click="ImportProject">-->
<!--                {{$t('project.ImportProject')}}-->
<!--              </a-button>-->
            </a>
            <a-row>
              <div v-if="projects.length>0">
                <a-col :xl="12"  :key="i" v-for="(item, i) in projects" style="padding-left: 10px;padding-right: 10px;" >
                  <a-card  @click="GoProject(item.ProjectInfo.uuid)" :hoverable="true" :bordered="true"  class="project-card">
                    <a-card-meta >
                      <div slot="title" class="projectTitle">
                        <p>{{item.ProjectInfo.name}}</p>
                      </div>
                      <div slot="description" class="description">
                        <p>{{item.ProjectInfo.description}}</p>
                      </div>
                    </a-card-meta>
                    <div  class="nav-box">
                      <div  class="nav-item">
                        <div  class="top">{{$t('project.ProjectDeviceCount')}}</div>
                        <div  class="bottom">{{item.DeviceCount}}</div>
                      </div>
                      <div  class="nav-item">
                      <div  class="top">{{$t('project.ProjectAppCount')}}</div>
                        <div  class="bottom">{{item.AppCount}}</div>
                      </div>
                      <div  class="nav-item">
                        <div  class="top">{{$t('project.ProjectDeviceOfflineCount')}}</div>
                        <div  class="bottom">{{item.DeviceOffCount}}</div>
                      </div>
                    </div>
                    <div  class="operate-box">
                      <div  class="create-time">
                        <label  for="create-time-label">{{$t('project.ProjectCreateTime')}}:</label>
                        <span  id="create-time-label">{{ item.ProjectInfo.CreatedAt|formatDate }}</span>
                      </div>
                      <div  class="create-time">
                        <label  for="create-time-label">{{$t('project.ProjectCreator')}}:</label>
                        <span  id="creator-label">{{ item.ProjectInfo.creator }}</span>
                      </div>
                      <div   class="operate-btn">
<!--                        <a-tooltip>-->
<!--                          <template slot="title">-->
<!--                            {{$t('project.ExportProject')}}-->
<!--                          </template>-->
<!--                          <button type="default" @click.stop="ExportProject($event,item.ProjectInfo)">-->
<!--                            <a-icon type="export" style="color: #ff4d4f"/>-->
<!--                          </button>-->
<!--                        </a-tooltip>-->

                        <a-tooltip>
                          <template slot="title">
                            {{$t('project.EditProject')}}
                          </template>
                          <button type="primary" @click.stop="editProject($event,item.ProjectInfo)">
                            <a-icon type="edit" />
                          </button>
                        </a-tooltip>

                        <a-tooltip>
                          <template slot="title">
                            {{$t('project.DelProject')}}
                          </template>
                          <button type="danger" @click.stop="delProject($event,item.ProjectInfo)">
                            <a-icon type="delete" style="color: #ff4d4f"/>
                          </button>
                        </a-tooltip>

                      </div>
                    </div>
                  </a-card>
                </a-col>
              </div>
              <div v-else>
                <a-empty >
                  <span slot="description"> {{$t('project.ProjectEmpty')}}</span>
                </a-empty>
              </div>
            </a-row>

          </a-card>
          <a-card style="margin-top: 20px" :bordered="true" :title="$t('project.ProjectHelp')">
            <a-steps  type="navigation" :current="4" direction="horizontal">
              <a-step status="finish">
                <template slot="title">
                  {{$t('project.ProjectHelpStep1')}}
                </template>
                <span slot="description">{{$t('project.ProjectHelpStep1Description')}}</span>
              </a-step>
              <a-step >
                <template slot="title">
                  {{$t('project.ProjectHelpStep2')}}
                </template>
                <span slot="description">{{$t('project.ProjectHelpStep2Description')}}</span>
              </a-step>
              <a-step >
                <template slot="title">
                  {{$t('project.ProjectHelpStep3')}}
                </template>
                <span slot="description">{{$t('project.ProjectHelpStep3Description')}}</span>
              </a-step>
              <a-step >
                <template slot="title">
                  {{$t('project.ProjectHelpStep4')}}
                </template>
                <span slot="description">{{$t('project.ProjectHelpStep4Description')}}</span>
              </a-step>
              <a-step status="finish">
                <template slot="title">
                  {{$t('project.ProjectHelpStep5')}}
                </template>
                <span slot="description">{{$t('project.ProjectHelpStep5Description')}}</span>
              </a-step>
            </a-steps>
          </a-card>
        </a-col>
<!--        <a-col style="padding: 0 10px 10px 15px" :xl="6" :lg="24" :md="24" :sm="24" :xs="24">-->
<!--          <a-card :title="$t('project.right.title')" style="margin-bottom: 10px;" :bordered="false" :body-style="{padding: 0}">-->
<!--            <a slot="extra" target="_blank" href="https://blog.csdn.net/hexinjun/category_12030477.html">更多</a>-->
<!--            <div class="right-guide-item">-->
<!--              <p>-->
<!--                <a href="https://blog.csdn.net/hexinjun/article/details/127066159" target="_blank" >ISM Web组态软件简介</a>-->
<!--              </p>-->
<!--              <p>-->
<!--                <a href="https://blog.csdn.net/hexinjun/article/details/125480702" target="_blank" >采集Modbus设备数据</a>-->
<!--              </p>-->
<!--              <p>-->
<!--                <a href="https://blog.csdn.net/hexinjun/article/details/125533642" target="_blank" >OPC UA设备的数据采集</a>-->
<!--              </p>-->
<!--              <p>-->
<!--                <a href="https://blog.csdn.net/hexinjun/article/details/125514049" target="_blank" >开发组态应用</a>-->
<!--              </p>-->
<!--              <p>-->
<!--                <a href="https://blog.csdn.net/hexinjun/article/details/125560403" target="_blank" >SNMP V3设备数据</a>-->
<!--              </p>-->
<!--            </div>-->
<!--          </a-card>-->
<!--          <a-card :title="$t('project.right.Dynamic')" style="margin-bottom: 5px;" :bordered="false" :body-style="{padding: 0}">-->
<!--            <div class="right-guide-dt">-->
<!--              <iframe :src="SystemDynamicUrl"-->
<!--                      width="100%"-->
<!--                      height="100%"-->
<!--                      frameborder="0"-->
<!--                      scrolling="no"> </iframe>-->
<!--            </div>-->
<!--          </a-card>-->
<!--        </a-col>-->
      </a-row>
      <a-modal lg v-model="projectVisible" :title="isEditStatus?$t('project.EditProject'):$t('project.AddProject')">
        <template slot="footer">
        <a-button type="primary"  @click="doAddProject">
          {{ $t('displayModel.ModelSure') }}
        </a-button>
        <a-button  @click="projectVisible=false">
          {{ $t('displayModel.ModelCancel') }}
        </a-button>
      </template>
      <a-form :form="AddModelForm" :label-col="{ span: 5 }" :wrapper-col="{ span: 16 }">
        <a-form-item :label="$t('project.ProjectName')">
          <a-input
              v-decorator="['ProjectName', { rules: [{ required: true, message: $t('project.ProjectName') }] }]"
          />
        </a-form-item>
        <a-form-item :label="$t('project.ProjectIndustry')">
          <a-select
              :getPopupContainer="
                triggerNode => {
                  return triggerNode.parentNode || document.body;
                }" style="width: 100%;height: 100%" :allowClear="true"  :dropdownStyle="{ maxHeight: '500px', overflow: 'auto' }"
              :dropdownMenuStyle="{ height:'100%',maxHeight: '500px', overflow: 'auto' }"
              v-decorator="['ProjectIndustry', { rules: [{ required: true, message: $t('project.ProjectIndustry') }] }]"
          >
            <a-select-option v-for="options in Industry" :key="options.value" :value="options.value">
              {{ $t(options.label)}}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="$t('project.ProjectDescription')">
          <Mtextarea   v-model="textAreValue"
                       rows="4"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
              v-decorator="['description', { rules: [{ required: true, message: $t('project.ProjectDescription') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </project-layout>
</template>

<script>
import Mtextarea from '@/components/textarea/index'
import {setAuthorization} from '@/utils/request'
import ProjectLayout from "../../layouts/ProjectLayout";
import {AUTH_TYPE} from "../../utils/request";
import {formatDate} from '@/utils/common';
import {
  ProjectAdd,
  ProjectList,
  ProjectDel,
  ProjectEdit,
  ExportProject,
  ImportProject,
  UpdateProject
} from "../../services/project";
const loadingKey = 'updatable'
export default {
  name: "Project",
  components: {ProjectLayout,Mtextarea},
  i18n: require('../../i18n/language'),
  data () {
    return {
      loading: false,
      textAreValue:"",
      OptLoading: false,
      isEditStatus:0,
      editProjectUUid:"",
      Industry:[
        {
          value:1,
          label:'project.IndustrialInternet',
        },
        {
          value:2,
          label:'project.SmartAgriculture',
        },
        {
          value:3,
          label:'project.Education',
        },
        {
          value:4,
          label:'project.NewEnergy',
        },
        {
          value:5,
          label:'project.IntelligentManufacturing',
        },
        {
          value:6,
          label:'project.SmartCity',
        },
        {
          value:7,
          label:'project.IntelligentMedicalTreatment',
        },
        {
          value:8,
          label:'project.OtherIndustry',
        }
      ],
      AddModelForm:this.$form.createForm(this),
      projectVisible:false,
      projects: [],
      form: this.$form.createForm(this)
    }
  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  computed: {
    SystemDynamicUrl () {
      return this.$store.state.setting.SystemDynamicUrl
    },
  },
  created(){
    this.GetProjectList()
  },
  methods: {
    ImportProject(){

    },
    GoProject(uuid){
      setAuthorization({token: uuid},AUTH_TYPE.AUTH1)
      this.$router.push('/dashboard')
    },
    editProject(e,item){
      this.projectVisible=true
      this.isEditStatus=1
      let _t = this
      this.editProjectUUid = item.uuid
      setTimeout(function (){
        _t.textAreValue = item.description
        _t.AddModelForm.setFieldsValue(
            {
              ProjectName:item.name,
              ProjectIndustry:item.industry,
              description:item.description,
            })
      },100)
    },
    ExportProject(e,item){
      let _t = this

      const params = {
        uuid:item.uuid
      };
      ExportProject(params).then(function (res){
        if (res.data.code == 200) {
          _t.GetProjectList()
          _t.$message.success(_t.$t('project.ProjectDelSuccess'), 3)
          setTimeout(function (){
            _t.$message.destroy();
          },300)
        }
        else {
          _t.$message.error(_t.$t('project.ProjectDelFailed'), 3)
          setTimeout(function (){
            _t.$message.destroy();
          },300)
        }
      })
    },
    delProject(e,item){
      let _t = this
      this.$confirm({
        title: item.name,
        content: _t.$t('project.DelProjectTips'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          _t.$message.loading({content: 'Loading...', loadingKey, duration: 0});
          const params = {
            uuid:item.uuid
          };
          ProjectDel(params).then(function (res){
            if (res.data.code == 200) {
              _t.GetProjectList()
              _t.$message.success(_t.$t('project.ProjectDelSuccess'), 3)
              setTimeout(function (){
                _t.$message.destroy();
              },300)
            }
            else {
              _t.$message.error(_t.$t('project.ProjectDelFailed'), 3)
              setTimeout(function (){
                _t.$message.destroy();
              },300)
            }
          })
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });
    },
    doAddProject(e){
      if(this.isEditStatus==0)
      {
        this.addProjectModel(e)
      }
      else
      {
        this.editProjectModel(e)
      }
    },
    editProjectModel(e){
      let _t = this
      e.preventDefault()
      this.AddModelForm.validateFields((err) => {
        if (!err) {
          _t.OptLoading = true
          const params = {
            uuid:_t.editProjectUUid,
            data:{
              name:this.AddModelForm.getFieldValue('ProjectName'),
              industry:parseInt(this.AddModelForm.getFieldValue('ProjectIndustry')),
              description:this.AddModelForm.getFieldValue('description'),
            }
          };
          ProjectEdit(params).then(function (res){
            _t.OptLoading = false
            if (res.data.code == 200) {
              _t.GetProjectList()
              _t.projectVisible = false;
              _t.$message.success(_t.$t('project.ProjectEditSuccess'), 3)
            }
            else {
              _t.$message.error(_t.$t('project.ProjectEditFailed'), 3)
            }
          })
        }
      })
    },
    addProjectModel(e){
      let _t = this
      e.preventDefault()
      this.AddModelForm.validateFields((err) => {
        if (!err) {
          _t.OptLoading = true
          const params = {
            name:this.AddModelForm.getFieldValue('ProjectName'),
            industry:parseInt(this.AddModelForm.getFieldValue('ProjectIndustry')),
            description:this.AddModelForm.getFieldValue('description'),
          };
          ProjectAdd(params).then(function (res){
            _t.OptLoading = false
            if (res.data.code == 4002) {
              _t.GetProjectList()
              _t.projectVisible = false;
              _t.$message.success(_t.$t('project.ProjectAddSuccess'), 3)
            }
            else {
              _t.$message.error(_t.$t('project.ProjectAddFailed'), 3)
            }
          })
        }
      })
    },
    GetProjectList(){
      let _t = this
      _t.loading = true
      this.projects=[]
      ProjectList().then(function (res){
        _t.loading = false
        if (res.data.code == 0) {
          if(res.data.list==null)
          {
            _t.projects=[]
          }
          else
          {
            _t.projects = res.data.list
          }

        }
      }).finally(function (error) {
        _t.loading = false
      })
    },
  }
}
</script>

<style scoped lang="less">
.right-guide-item  {
  margin-left: 20px;
  height: auto;
}
.right-guide-dt  {
  margin-left: 20px;
  height: 500px;
}
.right-guide-item p {
  display: block;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0px;
  margin-inline-end: 0px;
  font-size: 14px;
  font-weight: 400;
  color: #333;
  line-height: 20px;
  -webkit-font-feature-settings: "kern";
  font-feature-settings: "kern";
  cursor: pointer;
}
::v-deep  .ant-modal-header {
  padding: 16px 24px 5px;
  border-bottom: 0px solid #f0f0f0;
}
::v-deep .ant-form-item {
  margin-bottom: 15px;
}
::v-deep .ant-modal-body {
  padding: 20px 30px;
  max-height: 60vh;
  font-size: 14px;
  line-height: 1.5;
  word-wrap: break-word;
}
::v-deep .ant-modal .ant-modal-content{
  width: 640px;
}
::v-deep .ant-modal-footer {
  padding: 1px 16px 20px;
  border-top: 0px solid #f0f0f0;
}
.project-card .operate-box .operate-btn button i{
  width: 22px;
  height: 22px;
  font-size: 22px;
  color: @primary-color;
}
.project-card .operate-box .operate-btn button {
  width: 25px;
  height: 24px;
  border: 0;
  background: transparent;
  text-align: center;
  line-height: 24px;
  float: left;
  margin: 0 5px;
  cursor: pointer;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  padding-top: 2px;
}
::v-deep .ant-card-head-title {
  font-size: 16px;
  font-weight: 650;
  color: #333;
  line-height: 21px;
  -webkit-font-feature-settings: "kern";
  font-feature-settings: "kern";
}
.project-card .operate-box .create-time span{
  font-size: 14px;
  font-weight: 400;
  color: #333;
  line-height: normal;
  -webkit-font-feature-settings: "kern";
  font-feature-settings: "kern";
}
::v-deep .ant-card-hoverable:hover {
  border: 1px solid @primary-color;
  box-shadow: 10px 5px 5px rgba(0, 0, 0, 0.09);
}

::v-deep .ant-card-body {
  padding: 10px;
  //padding-left: 10px;
  //padding-right: 10px;
  zoom: 1;
}
::v-deep .ant-card-head {
  border-bottom: 0px solid #f0f0f0;
}
.project-card{
  background-color: #fff;
  border: 1px solid #e6e6e6;
  border-radius: 0;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  padding: 10px;
  height: 100%;
  cursor: pointer;
  margin-bottom: 10px;
}

.project-card .projectTitle{
  font-size: 20px;
  font-weight: 550;
  width: 100%;
  color: #333;
  height: 30px;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  z-index: 2;
}
.project-card .description{
  font-size: 20px;
  font-weight: 200;
  width: 100%;
  color: #333;
  height: 30px;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  z-index: 2;
}
.project-card .nav-box {
  margin-top: 20px;
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-pack: justify;
  -ms-flex-pack: justify;
  justify-content: space-between;
}
.project-card h2{
  font-size: 20px;
  font-weight: 400;
  color: #333;
  -webkit-font-feature-settings: "kern";
  font-feature-settings: "kern";
  line-height: 30px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.project-card .nav-box .nav-item{
  -webkit-box-flex: 1;
  -ms-flex: 1;
  flex: 1;
}
.project-card .nav-box .nav-item .top {
  font-size: 14px;
  font-weight: 400;
  color: #999;

  line-height: normal;
  -webkit-font-feature-settings: "kern";
  font-feature-settings: "kern";
}
.project-card .nav-box .nav-item .bottom {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  line-height: normal;
  -webkit-font-feature-settings: "kern";
  font-feature-settings: "kern";
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.project-card .operate-box {
  margin-top: 20px;
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-pack: justify;
  -ms-flex-pack: justify;
  justify-content: space-between;

}
.project-card .operate-box .create-time label {
  font-size: 14px;
  font-weight: 400;
  color: #999;
  line-height: normal;
  -webkit-font-feature-settings: "kern";
  font-feature-settings: "kern";
  margin-right: 5px;
}
</style>

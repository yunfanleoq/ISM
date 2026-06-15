<template>
  <a-card>
    <a-button @click="getSystemVideoList()" :disabled="messageShowLoad" type="default" >  {{$t('configComponent.video.VideoRefresh')}} </a-button>
    <div style="margin-top: 5px">
      <a-table rowKey="TableVideoIndex" :pagination="pagination" :columns="columns" :data-source="videoList" >
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="Opt" slot-scope="text, record">
          <a type="link"   @click="PlayModelShow(record)" style="cursor: pointer;color: #13C2C2"><a-icon type="youtube" /><span style="margin-left: 2px;">{{$t('configComponent.video.Play')}}</span></a>
        </div>
      </a-table>
    </div>

    <a-modal
        :visible="PlayModel"
        :title="$t('configComponent.video.VideoView')"
        :footer="null"
        :destroyOnClose="true"
        @cancel="handleCancel"
    >
      <div  class="root">
        <div id="container" ref="container"></div>
      </div>
    </a-modal>
  </a-card>
</template>

<script>
import { GetMonibucaVideoList} from "@/services/video";
import {mapState} from 'vuex'
import moment from 'moment'
import LivePlayer from '@liveqing/liveplayer'
export default {
  name: 'VideoManager',
  i18n: require('../../i18n/language'),
  components: {
    // LivePlayer
  },
  data () {
    return {
      jessibuca: null,
      version: '',
      wasm: false,
      vc: "ff",
      playing: false,
      quieting: true,
      loaded: false, // mute
      showOperateBtns: true,
      showBandwidth: true,
      err: "",
      speed: 0,
      performance: "",
      volume: 1,
      rotate: 0,
      useWCS: false,
      useMSE: true,
      useOffscreen: false,
      recording: false,
      recordType: 'webm',
      scale: 0,
      messageShowLoad:false,
      videoList:[],
      MonibucaServer:"",
      video_url:"",
      VideoTitle:"",
      PlayModel:false,
      editUuid:"",
      columns: [
        {
          width: '5%',
          slotName: 'configComponent.video.TableVideoIndex',
          scopedSlots: { customRender: 'TableVideoIndex', title: 'configComponent.video.TableVideoIndex' },
          dataIndex: 'TableVideoIndex',
        },
        {
          width: '10%',
          slotName: 'configComponent.video.TableVideoName',
          scopedSlots: { customRender: 'Name', title: 'configComponent.video.TableVideoName' },
          dataIndex: 'Name',
        },
        {
          slotName: 'configComponent.video.Manufacturer',
          width: '10%',
          scopedSlots: { customRender: 'Manufacturer', title: 'configComponent.video.Manufacturer' },
          dataIndex: 'Manufacturer',
        },
        {
          slotName: 'configComponent.video.Model',
          width: '15%',
          scopedSlots: { customRender: 'Model', title: 'configComponent.video.Model' },
          dataIndex: 'Model',
        },
        {
          slotName: 'configComponent.video.RegisterTime',
          width: '10%',
          scopedSlots: { customRender: 'RegisterTime',title: 'configComponent.video.RegisterTime'},
          dataIndex: 'RegisterTime',
        },
        {
          slotName: 'configComponent.video.UpdateTime',
          width: '10%',
          scopedSlots: { customRender: 'UpdateTime',title: 'configComponent.video.UpdateTime'},
          dataIndex: 'UpdateTime',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '5%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
    }
  },
  unmounted() {
    this.jessibuca.destroy();
  },
  computed: {
    ...mapState('setting', ['videoServer']),
  },
  authorize: {

  },
  mounted(){
    this.getSystemVideoList()
  },
  activated(){

  },
  created(){
  },
  watch: {

  },
  methods: {
    createVideo(options) {
      options = options || {};
      let _t = this
      this.jessibuca = new window.Jessibuca(
          Object.assign(
              {
                container: _t.$refs.container,
                videoBuffer: Number(6), // 缓存时长
                isResize: false,
                useWCS: this.useWCS,
                useMSE: this.useMSE,
                text: "",
                loadingText: "疯狂加载中...",
                hasAudio:true,
                debug: false,
                supportDblclickFullscreen: true,
                showBandwidth: this.showBandwidth, // 显示网速
                operateBtns: {
                  fullscreen: this.showOperateBtns,
                  screenshot: this.showOperateBtns,
                  play: this.showOperateBtns,
                  audio: this.showOperateBtns,
                },
                vod: this.vod,
                forceNoOffscreen: !this.useOffscreen,
                isNotMute: true,
                timeout: 10
              },
              options
          )
      );
    },
    handleCancel(){
      if (this.jessibuca) {
        this.jessibuca.destroy();
      }
      this.video_url=""
      this.PlayModel=false
    },
    PlayModelShow(videoInfo){
      this.PlayModel=true
      let _t = this
        setTimeout(function (){
          _t.createVideo()
          _t.video_url = _t.MonibucaServer+"hdl/"+videoInfo.LiveSubSP+".flv"
          console.log(_t.video_url)
          _t.jessibuca.play(_t.video_url);
        },500)
    },
    getSystemVideoList(){
      let _t = this
      _t.messageShowLoad=true
      _t.videoList=[]
      GetMonibucaVideoList().then(function (res){

        if(res.data.code==0)
        {
          _t.MonibucaServer = res.data.MonibucaServer
          let videoList = JSON.parse(res.data.data)
          let IndexNo=0
          for(let i=0;i<videoList.length;i++)
          {
            for(let k = 0;k<videoList[i].Channels.length;k++)
            {
              let tableData={}
              tableData.Manufacturer = videoList[i].Manufacturer
              tableData.RegisterTime = moment(videoList[i].RegisterTime).format('yyyy-MM-DD HH:mm:ss')
              tableData.UpdateTime = moment(videoList[i].UpdateTime).format('yyyy-MM-DD HH:mm:ss')
              tableData.TableVideoIndex=IndexNo
              tableData.Name = videoList[i].Channels[k].Name==""?videoList[i].Name:videoList[i].Channels[k].Name
              tableData.Model = videoList[i].Channels[k].Model
              tableData.Status = videoList[i].Channels[k].Status
              tableData.LiveSubSP = videoList[i].ID+"/"+videoList[i].Channels[k].DeviceID
              IndexNo++
              _t.videoList.push(tableData)
              tableData={}
            }
          }
        }
        _t.messageShowLoad=false
      }).catch(function (error) {
        _t.messageShowLoad=false
      });
    },
    selectVideo(type,key,Name){
      let url=""

      if(type==1)
      {
        url = "webrtc://"+this.videoServer+"/webrtcstream/"+key
      }
      else
      {
        url = key
      }

      this.$emit("onSelectVideo", {type:type,key:key,url:url,name:Name});
      this.visible =false
    },
  }
}
</script>

<style lang="less" scoped>
.root {
  display: flex;
  place-content: center;
}
::v-deep .ant-modal-body {
   padding: 2px;
}

#container {
  background: rgba(13, 14, 27, 0.7);
  width: 800px;
  height: 400px;
}

@media (max-width: 720px) {
  #container {
    width: 90vw;
    height: 52.7vw;
  }
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
::v-deep .ant-form-item {
  margin-bottom: 5px;

}
</style>

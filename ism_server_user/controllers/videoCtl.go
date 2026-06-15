/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:33
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

// 未完成寄存器组的修改，后续完成
import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	videoToWeb "ISMServer/protocol/videoServer"
	"ISMServer/utils/errmsg"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beevik/etree"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/mp4f"
	webrtc "github.com/deepch/vdk/format/webrtcv3"
	uuid "github.com/satori/go.uuid"
	"github.com/shiena/ansicolor"
	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	"github.com/use-go/onvif/xsd"
	"github.com/use-go/onvif/xsd/onvif"
	"golang.org/x/net/websocket"
)

type VideoController struct {
	beego.Controller
}
type JCodec struct {
	Type string
}
type JpgFileInfo struct {
	Path     string
	Size     int64
	Duration uint32
}

func (c *VideoController) AddRTSPStream() {

	type RTSPInfo struct {
		Name        string `form:"Name" json:"Name" binding:"required"`
		IP          string `form:"IP" json:"IP" binding:"required"`
		Port        uint32 `form:"port" json:"port" binding:"required"`
		User        string `form:"user" json:"user" `
		Password    string `form:"password" json:"password"`
		Url         string `form:"url" json:"url" binding:"required"`
		IsRecord    int    `form:"IsRecord" json:"IsRecord" binding:"required"`
		RecordInter int    `form:"RecordInter" json:"RecordInter" binding:"required"`
		Uuid        string
	}
	var rtsp RTSPInfo

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid == "" {
		result := map[string]interface{}{
			"code": -4,
			"msg":  "项目错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	//json数据封装到user对象中
	err := json.Unmarshal(data, &rtsp)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {

		//"rtsp://admin:wawjr19873110@16.100.1.108/cam/realmonitor?channel=1&subtype=0";
		var tmp videoToWeb.StreamST
		var saveTemp models.ProjectVideoList
		tmp.User = rtsp.User
		tmp.Password = rtsp.Password
		tmp.Ip = rtsp.IP
		tmp.Name = rtsp.Name
		tmp.Port = rtsp.Port
		tmp.StreamURL = rtsp.Url
		uuid := uuid.NewV4().String()
		tmp.Uuid = uuid
		rtsp.Uuid = uuid
		tmp.Cl = make(map[string]videoToWeb.Viwer)
		tmp.Status = 0

		if rtsp.User == "" || rtsp.Password == "" {
			tmp.URL = fmt.Sprintf("rtsp://%s:%d%s", tmp.Ip, tmp.Port, tmp.StreamURL)
		} else {
			tmp.URL = fmt.Sprintf("rtsp://%s:%s@%s:%d%s", tmp.User, tmp.Password, tmp.Ip, tmp.Port, tmp.StreamURL)
		}
		_, exist := videoToWeb.Config.Streams[rtsp.Uuid]
		if exist {
			result := map[string]interface{}{
				"code": -2,
				"msg":  "名称已经存在",
			}
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		saveTemp.User = tmp.User
		saveTemp.Name = tmp.Name
		saveTemp.Uuid = tmp.Uuid
		saveTemp.Password = tmp.Password
		saveTemp.Ip = tmp.Ip
		saveTemp.Port = tmp.Port
		saveTemp.StreamURL = tmp.StreamURL
		saveTemp.ProjectUuid = ProjectUuid
		saveTemp.IsUsed = 0
		saveTemp.IsRecord = rtsp.IsRecord
		saveTemp.RecordInter = rtsp.RecordInter
		code := models.VideoAdd(saveTemp)
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了视频"+tmp.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		if code == errmsg.SUCCSECODE {
			videoToWeb.Config.Streams[rtsp.Uuid] = tmp
		} else {
			result := map[string]interface{}{
				"code": code,
				"msg":  "添加失败",
			}
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}

	result := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
	}
	videoToWeb.VideoCloseChan()
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *VideoController) GetAllRTSPStreamList() {

	var getDbList []models.ProjectVideoList
	var returnData []models.ProjectVideoList
	var code = -1

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		getDbList, code = models.GetProjectVideoList(ProjectUuid)
		if code == errmsg.SUCCSECODE {
			for _, v := range getDbList {
				var tmp models.ProjectVideoList
				tmp.User = v.User
				tmp.Uuid = v.Uuid
				tmp.Key = v.Uuid
				tmp.Name = v.Name
				tmp.Password = v.Password
				tmp.Ip = v.Ip
				tmp.Port = v.Port
				tmp.IsUsed = v.IsUsed
				tmp.IsRecord = v.IsRecord
				tmp.RecordInter = v.RecordInter
				tmp.Status = videoToWeb.Config.GetVideoState(v.Uuid)
				tmp.StreamURL = v.StreamURL
				returnData = append(returnData, tmp)
			}
		}
	} else {
		code = -1
		returnData = nil
	}

	result := map[string]interface{}{
		"code": code,
		"list": returnData,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func readResponse(resp *http.Response) []byte {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
	}
	return b
}

func (c *VideoController) PtzControl() {

	type PtzControlStu struct {
		Ptzx      float64 `json:"Ptzx"`
		Ptzy      float64 `json:"Ptzy"`
		Ptzz      float64 `json:"Ptzz"`
		VideoUuid string  `json:"uuid"`
	}
	var getPtz PtzControlStu

	data := c.Ctx.Input.RequestBody

	var code = -1
	var onvifTokens = make(map[int]string)
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {

		//json数据封装到user对象中
		err := json.Unmarshal(data, &getPtz)
		if err != nil {
			result := map[string]interface{}{
				"code": -1,
			}
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}
		err1, info := models.VideoGetInfo(getPtz.VideoUuid)
		if err1 != 0 {
			result := map[string]interface{}{
				"code": -10,
			}
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}
		//Getting an camera instance
		dev, err := goonvif.NewDevice(goonvif.DeviceParams{
			// Xaddr:      "192.168.1.6:80",
			Xaddr:    info.Ip + ":80",
			Username: info.User,
			// Password:   "wawjr19873110$",
			Password:   info.Password,
			HttpClient: new(http.Client),
		})
		if err != nil {
			result := map[string]interface{}{
				"code": -2,
			}
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}

		mediaProfilesReq := media.GetProfiles{}
		mediaProfileResp, err := dev.CallMethod(mediaProfilesReq)
		if err != nil {
			result := map[string]interface{}{
				"code": -3,
			}
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}
		re := readResponse(mediaProfileResp)
		doc := etree.NewDocument()
		doc.ReadFromBytes(re)
		root := doc.SelectElement("Envelope")
		if root == nil {
			result := map[string]interface{}{
				"code": -4,
			}
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}
		token := root.FindElements("./Body/GetProfilesResponse/Profiles")
		if token == nil {
			result := map[string]interface{}{
				"code": -5,
			}
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}

		for k, res := range token {
			_token := res.SelectAttr("token").Value
			onvifTokens[k] = _token
		}
		ratio := math.Pow(10, 1)
		getPtz.Ptzx = math.Round(getPtz.Ptzx*ratio) / ratio
		getPtz.Ptzy = math.Round(getPtz.Ptzy*ratio) / ratio
		getPtz.Ptzz = math.Round(getPtz.Ptzz*ratio) / ratio
		ptzRelReq := ptz.AbsoluteMove{
			ProfileToken: onvif.ReferenceToken(onvifTokens[0]),
			Position: onvif.PTZVector{
				PanTilt: onvif.Vector2D{
					X:     getPtz.Ptzx,
					Y:     getPtz.Ptzy,
					Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/PanTiltSpaces/PositionGenericSpace"),
				},
				Zoom: onvif.Vector1D{
					X:     getPtz.Ptzz,
					Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/ZoomSpaces/PositionGenericSpace"),
				},
			},
			Speed: onvif.PTZSpeed{
				PanTilt: onvif.Vector2D{
					X:     0.1,
					Y:     0.1,
					Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/PanTiltSpaces/GenericSpeedSpace"),
				},
				Zoom: onvif.Vector1D{
					X:     0.1,
					Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/ZoomSpaces/ZoomGenericSpeedSpace"),
				},
			},
		}

		_, err67 := dev.CallMethod(ptzRelReq)
		if err67 != nil {
			result := map[string]interface{}{
				"code": -6,
			}
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}
	} else {
		code = -9
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *VideoController) DelRTSPStream() {

	type RTSPInfo struct {
		Uuid string `json:"uuid" binding:"required"`
	}
	var rtsp RTSPInfo

	data := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(data, &rtsp)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	key, exist := videoToWeb.Config.FindByUUID(rtsp.Uuid)
	if !exist {
		result := map[string]interface{}{
			"code": -2,
			"msg":  "名称不存在",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了视频", errmsg.JournalLevelInfo, c.Ctx.Input)
	code := models.VideoDel(rtsp.Uuid)
	if code == errmsg.SUCCSECODE {
		videoToWeb.Config.UpdatePthreadState(key, true)
		delete(videoToWeb.Config.Streams, key)
		result := map[string]interface{}{
			"code": 0,
			"msg":  "成功",
		}
		videoToWeb.VideoCloseChan()
		c.Data["json"] = result
		c.ServeJSON()
	} else {
		result := map[string]interface{}{
			"code": -5,
			"msg":  "失败",
		}
		c.Data["json"] = result
		c.ServeJSON()
	}
}
func (c *VideoController) EditRTSPStream() {

	type updateInfo struct {
		Uuid string                  `json:"uuid" binding:"required"`
		Data models.ProjectVideoList `json:"data" binding:"required"`
	}
	var rtsp updateInfo

	data := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(data, &rtsp)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	_, exist := videoToWeb.Config.FindByUUID(rtsp.Uuid)
	if !exist {
		result := map[string]interface{}{
			"code": -2,
			"msg":  "名称不存在",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了视频"+rtsp.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)

	code := models.VideoUpdate(rtsp.Uuid, rtsp.Data)
	if code == errmsg.SUCCSECODE {
		videoToWeb.Config.UpdateVideoData(rtsp.Uuid, rtsp.Data)
		videoToWeb.VideoCloseChan()
		result := map[string]interface{}{
			"code": 0,
			"msg":  "成功",
		}
		c.Data["json"] = result
		c.ServeJSON()
	} else {
		result := map[string]interface{}{
			"code": -5,
			"msg":  "失败",
		}
		c.Data["json"] = result
		c.ServeJSON()
	}
}

func (c *VideoController) WSLivestream() {
	c.EnableRender = false
	handler := websocket.Handler(ws)
	handler.ServeHTTP(c.Ctx.ResponseWriter, c.Ctx.Request)

}

func ws(ws *websocket.Conn) {
	defer ws.Close()

	suuid := ws.Request().FormValue("suuid")
	if !videoToWeb.Config.Ext(suuid) {
		logs.Error("Stream Not Found")
		return
	}
	ws.SetWriteDeadline(time.Now().Add(1 * time.Second))
	cuuid, ch := videoToWeb.Config.ClAd(suuid)
	defer videoToWeb.Config.ClDe(suuid, cuuid)
	codecs := videoToWeb.Config.CoGe(suuid)
	if codecs == nil {
		logs.Error("Codecs Error")
		return
	}
	for i, codec := range codecs {
		if codec.Type().IsAudio() && codec.Type() != av.AAC {
			logs.Error("Track", i, "Audio Codec Work Only AAC")
		}
	}
	muxer := mp4f.NewMuxer(nil)
	err := muxer.WriteHeader(codecs)
	if err != nil {
		logs.Error("muxer.WriteHeader", err)
		return
	}
	meta, init := muxer.GetInit(codecs)
	err = websocket.Message.Send(ws, append([]byte{9}, meta...))
	if err != nil {
		logs.Error("websocket.Message.Send", err)
		return
	}
	err = websocket.Message.Send(ws, init)
	if err != nil {
		return
	}
	var start bool
	go func() {
		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				ws.Close()
				return
			}
		}
	}()
	noVideo := time.NewTimer(60 * time.Second)
	var timeLine = make(map[int8]time.Duration)
	for {
		select {
		case <-noVideo.C:
			logs.Error("noVideo")
			return
		case pck := <-ch:
			if pck.IsKeyFrame {
				noVideo.Reset(60 * time.Second)
				start = true
			}
			if !start {
				continue
			}
			timeLine[pck.Idx] += pck.Duration
			pck.Time = timeLine[pck.Idx]
			ready, buf, _ := muxer.WritePacket(pck, false)
			if ready {
				err = ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					return
				}
				err := websocket.Message.Send(ws, buf)
				if err != nil {
					return
				}
			}
		}
	}
}

func (c *VideoController) Livestream() {
	c.EnableRender = false

	type RTSPInfo struct {
		Suuid string `json:"suuid"`
		Sdata string `json:"data"`
	}

	var rtsp RTSPInfo
	jsonData := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(jsonData, &rtsp)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	data := rtsp.Sdata
	suuid := rtsp.Suuid

	if !videoToWeb.Config.Ext(suuid) {
		logs.Error(suuid + " Stream Not Found")
		time.Sleep(time.Microsecond * 1000)
		return
	}

	codecs := videoToWeb.Config.CoGe(suuid)
	if codecs == nil {
		logs.Error(suuid + " Stream Codec Not Found")
		time.Sleep(time.Microsecond * 1000)
		return
	}
	var AudioOnly bool
	if len(codecs) == 1 && codecs[0].Type().IsAudio() {
		AudioOnly = true
	}
	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{ICEServers: videoToWeb.Config.GetICEServers(), ICEUsername: videoToWeb.Config.GetICEUsername(), ICECredential: videoToWeb.Config.GetICECredential(), PortMin: videoToWeb.Config.GetWebRTCPortMin(), PortMax: videoToWeb.Config.GetWebRTCPortMax()})
	answer, err := muxerWebRTC.WriteHeader(codecs, data)
	if err != nil {
		logs.Error(suuid+" WriteHeader", err)
		time.Sleep(time.Microsecond * 1000)
		return
	}
	err = c.Ctx.Output.Body([]byte(answer))
	if err != nil {
		logs.Error(suuid+" Write", err)
		time.Sleep(time.Microsecond * 1000)
		return
	}
	go func() {
		var sendFailedTimes = 0
		cid, ch := videoToWeb.Config.ClAd(suuid)
		defer videoToWeb.Config.ClDe(suuid, cid)
		defer muxerWebRTC.Close()
		var videoStart bool
		// noVideo := time.NewTimer(5 * time.Second)
		for {
			select {
			// case <-noVideo.C:
			// 	logs.Error(suuid + " noVideo")
			// 	return
			case pck := <-ch:
				if pck.IsKeyFrame || AudioOnly {
					// noVideo.Reset(5 * time.Second)
					videoStart = true
				}
				if !videoStart && !AudioOnly {
					continue
				}
				err = muxerWebRTC.WritePacket(pck)
				if err != nil {
					sendFailedTimes++
					if sendFailedTimes >= 5 {
						logs.Error(suuid+" WritePacket", err)
						sendFailedTimes = 0
						return
					}
				} else {
					sendFailedTimes = 0
				}
			}
		}
	}()

}

func (c *VideoController) WEBRTCLivestream() {
	c.EnableRender = false

	suuid := c.Ctx.Input.Param(":suuid")

	if !videoToWeb.Config.Ext(suuid) {
		logs.Error(suuid + " Stream Not Found")
		return
	}
	data := string(c.Ctx.Input.RequestBody)

	codecs := videoToWeb.Config.CoGe(suuid)
	if codecs == nil {
		logs.Error(suuid + " Stream Codec Not Found")
		return
	}

	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{ICEServers: videoToWeb.Config.GetICEServers(), ICEUsername: videoToWeb.Config.GetICEUsername(), ICECredential: videoToWeb.Config.GetICECredential(), PortMin: videoToWeb.Config.GetWebRTCPortMin(), PortMax: videoToWeb.Config.GetWebRTCPortMax()})

	answer, err := muxerWebRTC.WriteHeader(codecs, data)
	if err != nil {
		logs.Error(suuid+" WriteHeader", err)
		return
	}
	err = c.Ctx.Output.Body([]byte(answer))
	if err != nil {
		logs.Error(suuid+" Write", err)
		return
	}
	d := &videoToWeb.VideoCtl{WebRTC: muxerWebRTC, Key: suuid}
	go d.WebRTCPthread()

}

func (c *VideoController) Codec() {
	c.EnableRender = false
	if videoToWeb.Config.Ext(c.Ctx.Input.Param(":uuid")) {
		codecs := videoToWeb.Config.CoGe(c.Ctx.Input.Param(":uuid"))
		if codecs == nil {
			return
		}

		var tmpCodec []JCodec
		for _, codec := range codecs {
			if codec.Type() != av.H264 && codec.Type() != av.PCM_ALAW && codec.Type() != av.PCM_MULAW && codec.Type() != av.OPUS {
				log.Println("Codec Not Supported WebRTC ignore this track", codec.Type())
				continue
			}
			if codec.Type().IsVideo() {
				tmpCodec = append(tmpCodec, JCodec{Type: "video"})
			} else {
				tmpCodec = append(tmpCodec, JCodec{Type: "audio"})
			}
		}
		b, err := json.Marshal(tmpCodec)
		if err == nil {
			err = c.Ctx.Output.Body(b)
			if err != nil {
				log.Println("Write Codec Info error", err)
				return
			}
		}
	}
}

func (c *VideoController) GetVideoStatus() {
	c.EnableRender = false

	type RTSPInfo struct {
		Uuid string `json:"uuid" binding:"required"`
	}
	var rtsp RTSPInfo

	data := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(data, &rtsp)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	if videoToWeb.Config.Ext(rtsp.Uuid) {
		result := map[string]interface{}{
			"code":   0,
			"msg":    "成功",
			"result": videoToWeb.Config.Streams[rtsp.Uuid],
		}
		c.Data["json"] = result
		c.ServeJSON()
	} else {
		result := map[string]interface{}{
			"code": -2,
			"msg":  "没有此ID",
		}
		c.Data["json"] = result
		c.ServeJSON()
	}
}

func (c *VideoController) ExecCommand(name string, args ...string) {
	var ISMProcess *exec.Cmd
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	ISMProcess = exec.Command(name, args...) // 拼接参数与命令
	stdout, err3 := ISMProcess.StdoutPipe()
	if err3 != nil {
		fmt.Println("cmd.StdoutPipe: ", err3)
		return
	}
	var err error

	if err = ISMProcess.Start(); err != nil {
		log.Println(err)
	}

	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Fprint(w, line)
	}
}

func (c *VideoController) SetVideoStopOrStart() {

	type updateInfo struct {
		Uuid string                  `json:"uuid" binding:"required"`
		Data models.ProjectVideoList `json:"data" binding:"required"`
	}
	var rtsp updateInfo

	data := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(data, &rtsp)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	_, exist := videoToWeb.Config.FindByUUID(rtsp.Uuid)
	if !exist {
		result := map[string]interface{}{
			"code": -2,
			"msg":  "名称不存在",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	code := models.VideoStopOrStart(rtsp.Uuid, rtsp.Data)
	if code == errmsg.SUCCSECODE {
		videoToWeb.Config.UpdateVideoData(rtsp.Uuid, rtsp.Data)
		videoToWeb.VideoCloseChan()
		result := map[string]interface{}{
			"code": 0,
			"msg":  "成功",
		}
		c.Data["json"] = result
		c.ServeJSON()
	} else {
		result := map[string]interface{}{
			"code": -5,
			"msg":  "失败",
		}
		c.Data["json"] = result
		c.ServeJSON()
	}
}

func (c *VideoController) GetMonibucaVideoList() {

	method := "GET"
	result := map[string]interface{}{
		"code":           0,
		"MonibucaServer": "",
		"msg":            "失败",
		"data":           "",
	}
	MonibucaServer, getErr := config.String("MonibucaServer")
	if MonibucaServer == "" || getErr != nil {
		result["MonibucaServer"] = "http://127.0.0.1:18080/"
	} else {
		result["MonibucaServer"] = MonibucaServer
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, result["MonibucaServer"].(string)+"gb28181/api/list", nil)

	if err != nil {
		result["code"] = -3
		result["msg"] = err
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	req.Header.Add("User-Agent", "ISM/1.0.0 (https://ismctl.com.com)")

	res, err := client.Do(req)
	if err != nil {
		result["code"] = -2
		result["msg"] = err
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		result["code"] = -1
		result["msg"] = "失败"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	result["code"] = 0
	result["msg"] = "成功"
	result["data"] = string(body)

	c.Data["json"] = result
	c.ServeJSON()
}
func (c *VideoController) MonibucaHistoryVideoList() {

	var files []*JpgFileInfo
	filepath.WalkDir(protocol_common.RecordPath, func(path string, info os.DirEntry, err error) error {
		if !info.IsDir() {
			fileInfo, _ := info.Info()
			RecordPath := strings.ReplaceAll(protocol_common.RecordPath, "\\", "/")
			p := strings.ReplaceAll(path, "\\", "/")
			p = strings.TrimPrefix(p, RecordPath)
			//过滤抓拍的图片和视频
			if !strings.Contains(p, "SnapVideo") && !strings.Contains(p, "SnapImage") {
				files = append(files, &JpgFileInfo{
					Path: strings.TrimPrefix(p, "/"),
					Size: fileInfo.Size(),
				})
			}
		}

		return nil
	})
	result := map[string]interface{}{
		"code":  0,
		"msg":   "成功",
		"files": files,
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *VideoController) GetSnapTree() {
	var files []*JpgFileInfo
	filepath.WalkDir(protocol_common.RecordPath, func(path string, info os.DirEntry, err error) error {
		if !info.IsDir() {

			fileInfo, _ := info.Info()
			RecordPath := strings.ReplaceAll(protocol_common.RecordPath, "\\", "/")
			p := strings.ReplaceAll(path, "\\", "/")
			p = strings.TrimPrefix(p, RecordPath)
			p = "/record/" + p
			//过滤抓拍的图片和视频
			if strings.Contains(p, "SnapImage") {
				files = append(files, &JpgFileInfo{
					Path: strings.TrimPrefix(p, "/"),
					Size: fileInfo.Size(),
				})
			}
		}

		return nil
	})
	result := map[string]interface{}{
		"code":  0,
		"msg":   "成功",
		"files": files,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *VideoController) GetSnapVideoTree() {
	var files []*JpgFileInfo
	filepath.WalkDir(protocol_common.RecordPath, func(path string, info os.DirEntry, err error) error {
		if !info.IsDir() {
			fileInfo, _ := info.Info()
			RecordPath := strings.ReplaceAll(protocol_common.RecordPath, "\\", "/")
			p := strings.ReplaceAll(path, "\\", "/")
			p = strings.TrimPrefix(p, RecordPath)
			p = "/record/" + p
			//过滤抓拍的图片和视频
			if strings.Contains(p, "SnapVideo") {
				files = append(files, &JpgFileInfo{
					Path: strings.TrimPrefix(p, "/"),
					Size: fileInfo.Size(),
				})
			}
		}

		return nil
	})
	result := map[string]interface{}{
		"code":  0,
		"msg":   "成功",
		"files": files,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *VideoController) SnapRTSPStreamToJpg() {

	type RTSPInfo struct {
		Name     string `json:"Name"`
		Count    int    `json:"Count"`
		Interval int    `json:"Interval"`
	}
	var rtsp RTSPInfo

	data := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(data, &rtsp)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	Url, exist := videoToWeb.Config.FindByName(rtsp.Name)
	if !exist || Url == "" {
		result := map[string]interface{}{
			"code": -2,
			"msg":  "名称不存在",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	var fileName string
	if rtsp.Count != 1 {
		fileName = time.Now().Format("2006_01_02_15_04_05") + "_%03d.jpg"
	} else {
		fileName = time.Now().Format("2006_01_02_15_04_05") + ".jpg"
	}
	filePath := filepath.Join(protocol_common.RecordPath+"SnapImage/", rtsp.Name, time.Now().Format("2006-01-02"), fileName)

	if err = os.MkdirAll(filepath.Dir(filePath), 0766); err != nil {
		result := map[string]interface{}{
			"code": -3,
			"msg":  "目录创建失败",
		}
		c.Data["json"] = result
		c.ServeJSON()
	}
	if runtime.GOOS != "windows" {
		c.ExecCommand("./vendorBin/ffmpeg", "-i", Url, "-frames:v", fmt.Sprintf("%d", rtsp.Count), "-r", fmt.Sprintf("%d", rtsp.Interval), "-q:v", "2", filePath)
	} else {
		c.ExecCommand("vendorBin/ffmpeg.exe", "-i", Url, "-frames:v", fmt.Sprintf("%d", rtsp.Count), "-r", fmt.Sprintf("%d", rtsp.Interval), "-q:v", "2", filePath)
	}

	result := map[string]interface{}{
		"code": 0,
		"msg":  "抓拍成功",
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *VideoController) PlayVideo() {
	path := c.Ctx.Input.Param(":path")

	path = "d:\\" + "record\\" + path
	fmt.Println(path)
	data, err := os.ReadFile(string(path))
	if err != nil {
		fmt.Println(err)
	}
	if strings.HasSuffix(path, ".html") {
		c.Ctx.ResponseWriter.Header().Add("Content-Type", "text/html")
	} else if strings.HasSuffix(path, ".mp4") {
		c.Ctx.ResponseWriter.Header().Add("Content-Type", "video/mp4")
	}

	c.Ctx.ResponseWriter.Write(data)
}
func (c *VideoController) GetRecordVideoTree() {
	var files []*JpgFileInfo
	filepath.WalkDir(protocol_common.RecordPath, func(path string, info os.DirEntry, err error) error {
		if !info.IsDir() {
			fileInfo, _ := info.Info()
			p := strings.ReplaceAll(path, "\\", "/")
			files = append(files, &JpgFileInfo{
				Path: strings.TrimPrefix(p, "/"),
				Size: fileInfo.Size(),
			})
		}

		return nil
	})
	result := map[string]interface{}{
		"code":  0,
		"msg":   "成功",
		"files": files,
	}
	c.Data["json"] = result
	c.ServeJSON()
}

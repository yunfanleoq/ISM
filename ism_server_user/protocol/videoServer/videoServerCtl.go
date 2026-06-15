/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:01:12
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package videoserver

import (
	protocol_common "ISMServer/protocol/common"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/mp4"
	"github.com/deepch/vdk/format/rtspv2"
	webrtc "github.com/deepch/vdk/format/webrtcv3"
)

type VideoCtl struct {
	name              string
	Key               string
	url               string
	waitGroup         *sync.WaitGroup
	failedTimes       int
	webrtcFailedTimes int
	deviceStatus      int
	WebRTC            *webrtc.Muxer
	project_uuid      string
	rwMutex           sync.Mutex
	IsUsed            int
	Path              string
	BegenTime         int64
	Fragment          int64
	Recorder          *mp4.Muxer
	RecorderName      string
	RecordPath        string
	IsbegenRecord     bool
	RecordFw          *FileWriter
	RecordEnable      int
	RecordInter       int
}

type FileWr interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

type FileWriter struct {
	filePath string
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	bufw *bufio.Writer
}

var (
	ErrorStreamExitNoVideoOnStream = errors.New("Stream Exit No Video On Stream")
	ErrorStreamExitRtspDisconnect  = errors.New("Stream Exit Rtsp Disconnect")
	ErrorStreamExitNoViewer        = errors.New("Stream Exit On Demand No Viewer")
)

func (c *VideoCtl) InitVideoDevice() {
	c.Recorder = nil
	c.IsbegenRecord = false
	c.RecordPath = protocol_common.RecordPath
}
func (c *VideoCtl) InitVideoWebRTC(key, string, name string, url string) {
	c.name = name
	c.url = url
	c.Key = key
	c.failedTimes = 0
}

func (c *VideoCtl) videoReconnect() (session *rtspv2.RTSPClient) {
	name := c.name
	url := c.url
	uuid := c.Key

	rtspv2.Debug = true //Config.Server.Debug

	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: url, DisableAudio: true, DialTimeout: 10 * time.Second, ReadWriteTimeout: 180 * time.Second, Debug: false})

	if err != nil {
		logs.Error(name + " 设备连接失败 url:" + url + " 原因:" + fmt.Sprintf("%s", err))
		return nil
	}

	if RTSPClient.CodecData != nil {
		Config.coAd(uuid, RTSPClient.CodecData)
	}

	return RTSPClient
}
func (c *VideoCtl) CreateFileFn(filename string, append bool) (file FileWr, err error) {
	filePath := filepath.Join(c.RecordPath, filename)
	if err = os.MkdirAll(filepath.Dir(filePath), 0766); err != nil {
		return file, err
	}
	c.RecordFw = &FileWriter{filePath: filePath}

	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err == nil && !append {
		c.RecordFw.Reader = file
		c.RecordFw.Writer = file
		c.RecordFw.Seeker = file
		c.RecordFw.Closer = file
		return c.RecordFw, nil
	}
	return
}
func (c *VideoCtl) getFileName() (filename string) {
	filename = filepath.Join(filename, c.name, time.Now().Format("2006-01-02"), time.Now().Format("20060102T150405")+"_"+time.Now().Add(time.Second*time.Duration(c.RecordInter)).Format("20060102T150405"))
	return filename
}
func (c *VideoCtl) BegenRecord() {
	if c.RecordEnable == 1 {
		c.BegenTime = time.Now().Unix()
		c.RecorderName = c.getFileName() + ".mp4"
		logs.Info("开始录制,文件名：", c.RecorderName)
		fw, err := c.CreateFileFn(c.RecorderName, false)
		if err != nil {
			fmt.Println(err)
		} else {
			c.Recorder = mp4.NewMuxer(fw)
			c.Recorder.WriteHeader(Config.CoGe(c.Key))
			c.IsbegenRecord = true
		}
	} else {
		c.IsbegenRecord = false
	}
}
func (c *VideoCtl) RecordCut() {
	if c.RecordEnable == 1 {
		nowTime := time.Now().Unix()
		if nowTime-c.BegenTime >= int64(c.RecordInter) {
			if c.Recorder != nil {
				c.Recorder.WriteTrailer()
				c.RecordFw.Closer.Close()
			}
			c.Recorder = nil
			c.BegenRecord()
		}
	} else {
		c.IsbegenRecord = false
	}
}
func (c *VideoCtl) EndRecord() {
	if c.RecordEnable == 1 {
		if c.Recorder != nil {
			logs.Info("结束录制")
			c.Recorder.WriteTrailer()
			c.RecordFw.Closer.Close()
		}
	} else {
		c.IsbegenRecord = false
	}
}
func (c *VideoCtl) RecordWrite(packet av.Packet) {
	if c.RecordEnable == 1 {
		if c.IsbegenRecord {
			c.Recorder.WritePacket(packet)
			c.RecordCut()
		}
	} else {
		c.IsbegenRecord = false
	}
}
func (c *VideoCtl) VideoGatherPthread() {
	name := c.name
	uuid := c.Key
	c.failedTimes = 0
	var session *rtspv2.RTSPClient
	keyTest := time.NewTimer(20 * time.Second)

	var AlarmValue = 1
	var TempAlarmValue = 2
	var signleAlarm protocol_common.PushAlarm

	signleAlarm.DeviceUuid = uuid
	signleAlarm.ProjectUuid = c.project_uuid
	signleAlarm.DataUuid = "sys.suid.device.status"
	signleAlarm.ModelDataUuid = "sys.suid.device.status"

	signleAlarm.AlarmLevel = 3
	signleAlarm.AlarmClearMessage = "VideoManager.VideoOnline"
	signleAlarm.AlarmMessage = "VideoManager.VideoOffline"
	signleAlarm.DataName = "device.DeviceStatus"
	signleAlarm.DeviceName = name
	signleAlarm.HappenTime = time.Now()

	session = c.videoReconnect()

	for {

		//检测协程是否主动退出
		select {
		case <-GVideoChan:
			c.waitGroup.Done()
			c.EndRecord()
			Config.UpdateState(uuid, 0)
			logs.Error(name + " 主动退出")
			return
		default:
		}

		if session == nil {
			c.failedTimes++

			if c.failedTimes >= 5 {
				c.failedTimes = 0
				AlarmValue = 1
				Config.UpdateState(uuid, 0)
				signleAlarm.Value = fmt.Sprint(AlarmValue)
				signleAlarm.HappenTime = time.Now()
				c.failedTimes = 0
				TempAlarmValue = AlarmValue
				if c.IsUsed == 0 {
					protocol_common.GAlarmQueue.QueuePush(signleAlarm)
				}
				logs.Error(name + " 连接失败,重新连接")
			}
			session = c.videoReconnect()
			time.Sleep(5 * time.Second)
			continue
		} else {
			select {
			case <-keyTest.C:
				c.failedTimes++

				if c.failedTimes >= 5 {
					session.Close()
					c.failedTimes = 0
					session = nil
					AlarmValue = 1
					Config.UpdateState(uuid, 0)
					c.EndRecord()
					logs.Error(name + " 没有视频流")
				}

				time.Sleep(5 * time.Second)
			case signals := <-session.Signals:
				switch signals {
				case rtspv2.SignalCodecUpdate:
					Config.coAd(uuid, session.CodecData)
				case rtspv2.SignalStreamRTPStop:
					session.Close()
					session = nil
					Config.UpdateState(uuid, 0)
					AlarmValue = 1
					c.EndRecord()
					logs.Error(name + " 连接断开")
					time.Sleep(5 * time.Second)
				}
			case packetAV := <-session.OutgoingPacketQueue:
				if packetAV.IsKeyFrame {
					keyTest.Reset(20 * time.Second)
				}
				if AlarmValue == 1 {
					c.BegenRecord()
					logs.Info(name + " 连接成功")
				}

				c.RecordWrite(*packetAV)

				AlarmValue = 0
				Config.UpdateState(uuid, 1)
				Config.cast(uuid, *packetAV)
			}
		}
		if c.IsUsed == 0 {
			if AlarmValue != TempAlarmValue {
				signleAlarm.Value = fmt.Sprint(AlarmValue)
				signleAlarm.HappenTime = time.Now()
				c.failedTimes = 0
				TempAlarmValue = AlarmValue
				protocol_common.GAlarmQueue.QueuePush(signleAlarm)
			}
		}

		time.Sleep(time.Microsecond * 100)
	}
}

func (c *VideoCtl) WebRTCPthread() {
	uuid := c.Key
	name := c.name
	muxerWebRTC := c.WebRTC

	c.webrtcFailedTimes = 0
	cid, ch := Config.ClAd(uuid)
	defer Config.ClDe(uuid, cid)
	defer muxerWebRTC.Close()
	var videoStart bool
	var err error
	for {
		select {
		case pck := <-ch:
			if pck.IsKeyFrame {
				// noVideo.Reset(5 * time.Second)
				videoStart = true
			}
			if !videoStart {
				time.Sleep(time.Microsecond * 1000)
				continue
			}
			err = muxerWebRTC.WritePacket(pck)
			if err != nil {
				c.webrtcFailedTimes++
				if c.webrtcFailedTimes >= 5 {
					logs.Error(name+" WritePacket", err)
					c.webrtcFailedTimes = 0
					return
				}
			} else {
				c.webrtcFailedTimes = 0
			}
		}
		time.Sleep(time.Microsecond * 100)
	}
}

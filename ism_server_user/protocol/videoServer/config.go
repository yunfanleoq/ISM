/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:01:06
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package videoserver

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
)

// Config global
var Config = LoadConfig()
var (
	Success                         = "success"
	ErrorStreamNotFound             = errors.New("stream not found")
	ErrorStreamAlreadyExists        = errors.New("stream already exists")
	ErrorStreamChannelAlreadyExists = errors.New("stream channel already exists")
	ErrorStreamNotHLSSegments       = errors.New("stream hls not ts seq found")
	ErrorStreamNoVideo              = errors.New("stream no video")
	ErrorStreamNoClients            = errors.New("stream no clients")
	ErrorStreamRestart              = errors.New("stream restart")
	ErrorStreamStopCoreSignal       = errors.New("stream stop core signal")
	ErrorStreamStopRTSPSignal       = errors.New("stream stop rtsp signal")
	ErrorStreamChannelNotFound      = errors.New("stream channel not found")
	ErrorStreamChannelCodecNotFound = errors.New("stream channel codec not ready, possible stream offline")
	ErrorStreamsLen0                = errors.New("streams len zero")
)

// ConfigST struct
type ConfigST struct {
	mutex   sync.RWMutex
	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

type SaveConfigST struct {
	Server  ServerST                           `json:"server"`
	Streams map[string]models.ProjectVideoList `json:"streams"`
}

// ServerST struct
type ServerST struct {
	HTTPPort      string   `json:"http_port"`
	Debug         bool     `json:"debug"`
	ICEServers    []string `json:"ice_servers"`
	ICEUsername   string   `json:"ice_username"`
	ICECredential string   `json:"ice_credential"`
	WebRTCPortMin uint16   `json:"webrtc_port_min"`
	WebRTCPortMax uint16   `json:"webrtc_port_max"`
}

// Segment HLS cache section
type Segment struct {
	dur  time.Duration
	data []*av.Packet
}

// StreamST struct
type StreamST struct {
	Uuid             string `json:"uuid"`
	StreamURL        string `json:"StreamURL"`
	URL              string
	Ip               string `json:"ip"`
	User             string `json:"user"`
	Password         string `json:"password"`
	Port             uint32 `json:"port"`
	Name             string `json:"Name"`
	Status           int    `json:"status"`
	ProjectUuid      string `json:"ProjectUuid"`
	IsExit           bool
	RecordEnable     int
	Codecs           []av.CodecData
	Cl               map[string]Viwer
	IsUsed           int
	RecordInter      int
	hlsSegmentNumber int             `json:"-"`
	hlsSegmentBuffer map[int]Segment `json:"-"`
}

type Viwer struct {
	c chan av.Packet
}

func (element *ConfigST) GetICEServers() []string {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	return element.Server.ICEServers
}

func (element *ConfigST) GetICEUsername() string {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	return element.Server.ICEUsername
}

func (element *ConfigST) GetICECredential() string {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	return element.Server.ICECredential
}

func (element *ConfigST) GetWebRTCPortMin() uint16 {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	return element.Server.WebRTCPortMin
}

func (element *ConfigST) GetWebRTCPortMax() uint16 {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	return element.Server.WebRTCPortMax
}

func LoadConfig() *ConfigST {
	var tmp ConfigST
	data, err := ioutil.ReadFile(ConfigName)
	if err != nil {
		logs.Error(err)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		logs.Error(err)
	}
	tmp.Streams = make(map[string]StreamST)
	videoList, code := models.GetAllVideoList()
	if code == errmsg.SUCCSECODE {
		for _, v := range videoList {
			var videoInfo StreamST
			videoInfo.Cl = make(map[string]Viwer)
			videoInfo.Uuid = v.Uuid
			videoInfo.Ip = v.Ip
			videoInfo.Port = v.Port
			videoInfo.StreamURL = v.StreamURL
			videoInfo.User = v.User
			videoInfo.Password = v.Password
			videoInfo.Name = v.Name
			videoInfo.ProjectUuid = v.ProjectUuid
			videoInfo.IsUsed = v.IsUsed
			videoInfo.RecordEnable = v.IsRecord
			videoInfo.RecordInter = v.RecordInter

			videoInfo.hlsSegmentBuffer = make(map[int]Segment)

			if v.User == "" || v.Password == "" {
				videoInfo.URL = fmt.Sprintf("rtsp://%s:%d%s", v.Ip, v.Port, v.StreamURL)
			} else {
				videoInfo.URL = fmt.Sprintf("rtsp://%s:%s@%s:%d%s", v.User, v.Password, v.Ip, v.Port, v.StreamURL)
			}
			tmp.Streams[v.Uuid] = videoInfo
		}
	}
	return &tmp
}

func (element *ConfigST) cast(uuid string, pck av.Packet) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[uuid]
	if ok {
		for _, v := range element.Streams[uuid].Cl {
			if len(v.c) < cap(v.c) {
				v.c <- pck
			}
		}
	}
}

func (element *ConfigST) Ext(suuid string) bool {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	return ok
}
func (element *ConfigST) UpdateState(suuid string, value int) bool {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		t := element.Streams[suuid]
		t.Status = value
		element.Streams[suuid] = t
		return true
	} else {
		return false
	}
}

func (element *ConfigST) UpdateVideoData(suuid string, value models.ProjectVideoList) bool {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		t := element.Streams[suuid]
		t.Ip = value.Ip
		t.Port = value.Port
		t.StreamURL = value.StreamURL
		t.User = value.User
		t.Password = value.Password
		t.Name = value.Name
		t.IsUsed = value.IsUsed
		if value.User == "" || value.Password == "" {
			t.URL = fmt.Sprintf("rtsp://%s:%d%s", value.Ip, value.Port, value.StreamURL)
		} else {
			t.URL = fmt.Sprintf("rtsp://%s:%s@%s:%d%s", value.User, value.Password, value.Ip, value.Port, value.StreamURL)
		}
		element.Streams[suuid] = t
		return true
	} else {
		return false
	}
}
func (element *ConfigST) GetVideoState(suuid string) int {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		t := element.Streams[suuid]
		return t.Status
	} else {
		return 0
	}
}

func (element *ConfigST) UpdatePthreadState(suuid string, value bool) bool {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		t := element.Streams[suuid]
		t.IsExit = value
		element.Streams[suuid] = t
		return true
	} else {
		return false
	}
}

func (element *ConfigST) getPthreadStatus(suuid string) bool {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		return element.Streams[suuid].IsExit
	} else {
		return false
	}

}

func (element *ConfigST) coAd(suuid string, codecs []av.CodecData) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		t := element.Streams[suuid]
		t.Codecs = codecs
		element.Streams[suuid] = t
	}
}

func (element *ConfigST) CoGe(suuid string) []av.CodecData {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		return element.Streams[suuid].Codecs
	} else {
		return nil
	}
}

func (element *ConfigST) ClAd(suuid string) (string, chan av.Packet) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	cuuid := pseudoUUID()
	ch := make(chan av.Packet, 100)
	if ok {
		element.Streams[suuid].Cl[cuuid] = Viwer{c: ch}
		return cuuid, ch
	} else {
		return cuuid, ch
	}
}

func (element *ConfigST) list() (string, []string) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	var res []string
	var fist string
	for k := range element.Streams {
		if fist == "" {
			fist = k
		}
		res = append(res, k)
	}
	return fist, res
}
func (element *ConfigST) ClDe(suuid, cuuid string) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	if ok {
		delete(element.Streams[suuid].Cl, cuuid)
	}
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
func (element *ConfigST) FindByUUID(uuid string) (string, bool) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	var mapKey = ""
	var err = false
	for k, v := range element.Streams {
		if v.Uuid == uuid {
			mapKey = k
			err = true
			break
		}
	}
	return mapKey, err
}
func (element *ConfigST) FindByName(name string) (string, bool) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	var URL = ""
	var err = false
	for _, v := range element.Streams {
		if v.Name == name {
			URL = v.URL
			err = true
			break
		}
	}
	return URL, err
}

// StreamHLSAdd add hls seq to buffer
func (obj *ConfigST) StreamHLSAdd(uuid string, val []*av.Packet, dur time.Duration) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		tmp.hlsSegmentNumber++
		tmp.hlsSegmentBuffer[tmp.hlsSegmentNumber] = Segment{data: val, dur: dur}
		if len(tmp.hlsSegmentBuffer) >= 6 {
			delete(tmp.hlsSegmentBuffer, tmp.hlsSegmentNumber-6-1)
		}
		obj.Streams[uuid] = tmp
	}
}

// StreamHLSm3u8 get hls m3u8 list
func (obj *ConfigST) StreamHLSm3u8(uuid string) (string, int, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		var out string
		//TODO fix  it
		out += "#EXTM3U\r\n#EXT-X-TARGETDURATION:4\r\n#EXT-X-VERSION:4\r\n#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(tmp.hlsSegmentNumber) + "\r\n"
		var keys []int
		for k := range tmp.hlsSegmentBuffer {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		var count int
		for _, i := range keys {
			count++
			out += "#EXTINF:" + strconv.FormatFloat(tmp.hlsSegmentBuffer[i].dur.Seconds(), 'f', 1, 64) + ",\r\nsegment/" + strconv.Itoa(i) + "/file.ts\r\n"

		}
		return out, count, nil
	}
	return "", 0, ErrorStreamNotFound
}

// StreamHLSTS send hls segment buffer to clients
func (obj *ConfigST) StreamHLSTS(uuid string, seq int) ([]*av.Packet, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if buf, ok := tmp.hlsSegmentBuffer[seq]; ok {
			return buf.data, nil
		}
	}
	return nil, ErrorStreamNotFound
}

// StreamHLSFlush delete hls cache
func (obj *ConfigST) StreamHLSFlush(uuid string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		tmp.hlsSegmentBuffer = make(map[int]Segment)
		tmp.hlsSegmentNumber = 0
		obj.Streams[uuid] = tmp
	}
}

// stringToInt convert string to int if err to zero
func stringToInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

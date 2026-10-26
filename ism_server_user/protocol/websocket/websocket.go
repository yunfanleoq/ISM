/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:26
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-24 15:34:37
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ismWebsocket

import (
	"ISMServer/middleware"
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	alarmTask "ISMServer/task/RealData"
	SSEConnManager "ISMServer/utils/SSE"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-basic/uuid"
	"github.com/gorilla/websocket"
)

type WsServer struct {
	listener net.Listener
	addr     string
	upgrade  *websocket.Upgrader
}

var websocketConnArray sync.Map //make(map[string][]*WsConnection)

// realDataDropLog 节流：满通道 drop 时避免每条都刷屏拖垮进程
var realDataDropLog struct {
	mu         sync.Mutex
	lastLog    time.Time
	dropped    uint64
	suppressed uint64
	fileOnce   sync.Once
	file       *os.File
}

// 实时推送合并窗（毫秒）：同点位只保留最新值，降低 RealDataChanel 填满概率。
// 会议「调长周期」的正确落点；默认 2000ms，可用 app.conf RealDataPushMergeMs 覆盖。
var realDataPushMergeMs int64 = 2000

type realDataMergeBucket struct {
	deviceUuid  string
	deviceName  string
	projectUuid string
	points      map[string]protocol_common.UpdateStu
	lastFlush   time.Time
}

var realDataMerge struct {
	mu      sync.Mutex
	buckets map[string]*realDataMergeBucket
	started int32
}

func init() {
	// 供 models 等包推送 RealData，避免 models↔websocket 循环依赖。
	protocol_common.RealDataFrontendPush = func(msg protocol_common.PushRealDataWebData) {
		WSSend(msg, msg.ProjectUuid, 2)
	}
}

func initRealDataPushMerge() {
	if !atomic.CompareAndSwapInt32(&realDataMerge.started, 0, 1) {
		return
	}
	ms, err := config.Int("RealDataPushMergeMs")
	if err == nil && ms >= 0 {
		realDataPushMergeMs = int64(ms)
	}
	realDataMerge.buckets = make(map[string]*realDataMergeBucket)
	if realDataPushMergeMs <= 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			flushRealDataMergeBuckets(false)
		}
	}()
}

func flushRealDataMergeBuckets(force bool) {
	mergeMs := atomic.LoadInt64(&realDataPushMergeMs)
	if mergeMs <= 0 && !force {
		return
	}
	now := time.Now()
	realDataMerge.mu.Lock()
	pending := make([]*realDataMergeBucket, 0, len(realDataMerge.buckets))
	for key, b := range realDataMerge.buckets {
		if b == nil || len(b.points) == 0 {
			continue
		}
		if !force && mergeMs > 0 && now.Sub(b.lastFlush) < time.Duration(mergeMs)*time.Millisecond {
			continue
		}
		pending = append(pending, b)
		delete(realDataMerge.buckets, key)
	}
	realDataMerge.mu.Unlock()
	for _, b := range pending {
		msg := protocol_common.PushRealDataWebData{
			DeviceUuid:  b.deviceUuid,
			DeviceName:  b.deviceName,
			ProjectUuid: b.projectUuid,
			Cmd:         "RealData",
			Data:        make([]protocol_common.UpdateStu, 0, len(b.points)),
		}
		for _, p := range b.points {
			msg.Data = append(msg.Data, p)
		}
		pushToProjectConns(msg, b.projectUuid, 2)
	}
}

func enqueueMergedRealData(msg protocol_common.PushRealDataWebData) {
	initRealDataPushMerge()
	mergeMs := atomic.LoadInt64(&realDataPushMergeMs)
	if mergeMs <= 0 || len(msg.Data) == 0 {
		pushToProjectConns(msg, msg.ProjectUuid, 2)
		return
	}
	key := msg.ProjectUuid + "|" + msg.DeviceUuid
	realDataMerge.mu.Lock()
	b := realDataMerge.buckets[key]
	if b == nil {
		b = &realDataMergeBucket{
			deviceUuid:  msg.DeviceUuid,
			deviceName:  msg.DeviceName,
			projectUuid: msg.ProjectUuid,
			points:      make(map[string]protocol_common.UpdateStu, len(msg.Data)),
			lastFlush:   time.Now(),
		}
		realDataMerge.buckets[key] = b
	}
	if msg.DeviceName != "" {
		b.deviceName = msg.DeviceName
	}
	for _, p := range msg.Data {
		id := p.Uuid
		if id == "" {
			id = p.DataName
		}
		if id == "" {
			continue
		}
		b.points[id] = p
	}
	shouldFlush := time.Since(b.lastFlush) >= time.Duration(mergeMs)*time.Millisecond
	var flushCopy *realDataMergeBucket
	if shouldFlush {
		flushCopy = b
		delete(realDataMerge.buckets, key)
	}
	realDataMerge.mu.Unlock()
	if flushCopy != nil {
		out := protocol_common.PushRealDataWebData{
			DeviceUuid:  flushCopy.deviceUuid,
			DeviceName:  flushCopy.deviceName,
			ProjectUuid: flushCopy.projectUuid,
			Cmd:         "RealData",
			Data:        make([]protocol_common.UpdateStu, 0, len(flushCopy.points)),
		}
		for _, p := range flushCopy.points {
			out.Data = append(out.Data, p)
		}
		pushToProjectConns(out, out.ProjectUuid, 2)
	}
}

func openDropLogFile() *os.File {
	realDataDropLog.fileOnce.Do(func() {
		_ = os.MkdirAll("logs", 0o755)
		path := filepath.Join("logs", "ws_realdata_drop.log")
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Printf("[WS] open drop log failed: %v\n", err)
			return
		}
		realDataDropLog.file = f
	})
	return realDataDropLog.file
}

func noteRealDataChanelDrop(project string, ch chan interface{}, connCount int) {
	realDataDropLog.mu.Lock()
	realDataDropLog.dropped++
	now := time.Now()
	if now.Sub(realDataDropLog.lastLog) < time.Second {
		realDataDropLog.suppressed++
		realDataDropLog.mu.Unlock()
		return
	}
	total := realDataDropLog.dropped
	suppressed := realDataDropLog.suppressed
	realDataDropLog.lastLog = now
	realDataDropLog.suppressed = 0
	realDataDropLog.mu.Unlock()

	chLen, chCap := 0, 0
	if ch != nil {
		chLen = len(ch)
		chCap = cap(ch)
	}
	line := fmt.Sprintf("%s [WS] RealDataChanel full, drop (total=%d, suppressed=%d/s, project=%s, conn=%d, chan=%d/%d, mergeMs=%d)\n",
		now.Format("2006-01-02 15:04:05.000"), total, suppressed, project, connCount, chLen, chCap, atomic.LoadInt64(&realDataPushMergeMs))
	fmt.Print(line)
	if f := openDropLogFile(); f != nil {
		_, _ = f.WriteString(line)
	}
}

// tryPushChanel 非阻塞入队；满则立即丢弃，禁止在持锁路径上 sleep。
func tryPushChanel(ch chan interface{}, message interface{}, isRealData bool, project string, connCount int) {
	select {
	case ch <- message:
	default:
		if isRealData {
			noteRealDataChanelDrop(project, ch, connCount)
		}
	}
}

func realDataDeviceKey(message interface{}) string {
	if realMsg, ok := message.(protocol_common.PushRealDataWebData); ok {
		if realMsg.DeviceUuid != "" {
			return realMsg.DeviceUuid
		}
		return realMsg.ProjectUuid + "|" + realMsg.DeviceName
	}
	return ""
}

// tryPushRealDataChanel 高水位时按设备 latest-wins 压缩队列，再入队；否则非阻塞入队。
func tryPushRealDataChanel(ch chan interface{}, message interface{}, project string, connCount int) {
	if ch == nil {
		return
	}
	capN := cap(ch)
	if capN > 0 && len(ch)*5 >= capN*4 {
		latest := make(map[string]interface{}, 64)
		order := make([]string, 0, 64)
		anon := 0
		drainOne := func(m interface{}) {
			key := realDataDeviceKey(m)
			if key == "" {
				anon++
				key = fmt.Sprintf("_anon_%d", anon)
			}
			if _, exists := latest[key]; !exists {
				order = append(order, key)
			}
			latest[key] = m
		}
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				drainOne(m)
			default:
				goto requeue
			}
		}
	requeue:
		drainOne(message)
		for _, key := range order {
			select {
			case ch <- latest[key]:
			default:
				noteRealDataChanelDrop(project, ch, connCount)
				return
			}
		}
		return
	}
	tryPushChanel(ch, message, true, project, connCount)
}

type WsConnection struct {
	connId               string
	ws                   *websocket.Conn
	RealDataChanel       chan interface{}
	RealAlarmChanel      chan interface{}
	RealSystemDataChanel chan interface{}
	project              string
	mutex                sync.Mutex
	RwMutex              sync.Mutex
	isClosed             bool
}

func NewWsConnection(conn *websocket.Conn) *WsConnection {
	ws := &WsConnection{}
	ws.ws = conn
	ws.RealDataChanel = make(chan interface{}, protocol_common.RealDataChanelCache)
	ws.RealAlarmChanel = make(chan interface{}, protocol_common.RealDataChanelCache)
	ws.RealSystemDataChanel = make(chan interface{}, protocol_common.RealDataChanelCache)
	return ws
}

func NewWsServer() *WsServer {
	WSPort, wserr := config.Int("WSPort")
	if wserr != nil {
		WSPort = 10215
	}
	ws := new(WsServer)
	ws.addr = "0.0.0.0:" + fmt.Sprintf("%d", WSPort)
	ws.upgrade = &websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   8192,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				fmt.Println("method is not GET")
				return false
			}
			if r.URL.Path != "/ws" {
				fmt.Println("path error")
				return false
			}
			return true
		},
	}
	return ws
}

func (c *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	project := r.URL.Query().Get("project")
	if r.URL.Path != "/ws" {
		httpCode := http.StatusInternalServerError
		reasePhrase := http.StatusText(httpCode)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	if token == "" || project == "" {
		httpCode := http.StatusForbidden
		reasePhrase := http.StatusText(httpCode)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	result, _, _, _, _ := middleware.JwtToken(token)
	if result != errmsg.SUCCSE {
		httpCode := http.StatusUnauthorized
		reasePhrase := http.StatusText(httpCode)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	//查询项目ID是否存在
	_, code := models.ProjectSingleModel(project)
	if code == -1 {
		httpCode := http.StatusUnauthorized
		reasePhrase := http.StatusText(httpCode)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	conn, err := c.upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	ws := NewWsConnection(conn)
	ws.project = project
	ws.connId = uuid.New()
	v, ok := websocketConnArray.Load(project)
	if ok {
		conn := v.([]*WsConnection)
		connNewList := append(conn, ws)
		websocketConnArray.Store(project, connNewList)
	} else {
		var conn []*WsConnection
		connNewList := append(conn, ws)
		websocketConnArray.Store(project, connNewList)
	}
	go ws.connHandleRealData()
	go ws.connHandleRealAlarm()
	go ws.connHandleRealSystemData()
	go ws.connHandleHeart()
}
func (conn *WsConnection) Close() {
	conn.ws.Close() //线程安全的
	if !conn.isClosed {
		//一个chan只能关闭一次，保证此代码只执行一次
		close(conn.RealDataChanel)
		close(conn.RealAlarmChanel)
		close(conn.RealSystemDataChanel)
		// conn.sendpoll.Release()
		conn.isClosed = true
	}
}
func (conn *WsConnection) WriteToClient(msg any) error {
	defer conn.RwMutex.Unlock()
	conn.RwMutex.Lock()
	if conn.isClosed || conn.ws == nil {
		return fmt.Errorf("connection closed")
	}
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("解析错误")
	}
	_ = conn.ws.SetWriteDeadline(time.Now().Add(3 * time.Second))
	err = conn.ws.WriteMessage(websocket.TextMessage, jsonBytes)
	_ = conn.ws.SetWriteDeadline(time.Time{})
	return err
}

func (c *WsConnection) dropDeadConn() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isClosed {
		return
	}
	DeleteStringElement(c.connId, c.project)
	c.Close()
}

func (c *WsConnection) connHandleRealData() {
	if c.project == "" {
		return
	}
	for {
		if c.isClosed {
			return
		}
		for realDataMsg := range c.RealDataChanel {
			if err := c.WriteToClient(realDataMsg); err != nil {
				c.dropDeadConn()
				return
			}
		}
	}
}
func (c *WsConnection) connHandleHeart() {
	if c.project == "" {
		return
	}
	for {
		if c.isClosed {
			return
		}

		ee := c.WriteToClient("ping")
		if ee != nil {
			c.dropDeadConn()
			return
		}
		time.Sleep(5 * time.Second)
	}
}
func (c *WsConnection) connHandleRealAlarm() {
	if c.project == "" {
		return
	}
	for {
		if c.isClosed {
			return
		}
		for realDataMsg := range c.RealAlarmChanel {
			if err := c.WriteToClient(realDataMsg); err != nil {
				c.dropDeadConn()
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
func (c *WsConnection) connHandleRealSystemData() {
	if c.project == "" {
		return
	}
	for {
		if c.isClosed {
			return
		}
		for realDataMsg := range c.RealSystemDataChanel {
			if err := c.WriteToClient(realDataMsg); err != nil {
				c.dropDeadConn()
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
func DeleteStringElement(uuid string, project string) []*WsConnection {
	result := make([]*WsConnection, 0)
	connList, ok := websocketConnArray.Load(project)
	if ok && connList != nil {
		getConn := connList.([]*WsConnection)
		for _, v := range getConn {
			if v.connId != uuid {
				result = append(result, v)
			}
		}
		websocketConnArray.Store(project, result)
	}
	return result
}
func SendToISMNode(project string, message interface{}, mType int) {
	SSEConnManager.GlobalConnManager.PushToAll(message, 1000)
	if GetConn, IsTrue := protocol_common.ISMNodeProjectConn.Load(project); IsTrue {
		Conn, ok := GetConn.([]protocol_common.ISMNodeProjectConnStu)
		if !ok {
			return
		}
		type ISMNetNodeFormatCmd struct {
			Cmd         string `json:"Cmd"`
			PackIndex   int64  `json:"PackIndex"`
			ProjectUuid string `json:"ProjectUuid"`
			NodeName    string `json:"NodeName"`
			Data        struct {
				MessageType   int `json:"MessageType"`
				RequestParams any `json:"RequestParams"`
			}
		}
		var RealPush ISMNetNodeFormatCmd
		RealPush.Cmd = "PushRealData"
		RealPush.ProjectUuid = project
		RealPush.PackIndex = time.Now().UnixMilli()
		RealPush.Data.MessageType = mType
		RealPush.Data.RequestParams = message
		jsonBytes, err := json.Marshal(RealPush)
		if err != nil {
			return
		}
		for _, itemConn := range Conn {
			if itemConn.ChanelConn != nil {
				itemConn.ConnRwMutex.Lock()
				itemConn.ChanelConn.WriteMessage(websocket.TextMessage, jsonBytes)
				itemConn.ConnRwMutex.Unlock()
			}
		}
	}
}
func pushToProjectConns(message interface{}, project string, msgType int) {
	connList, ok := websocketConnArray.Load(project)
	if !ok || connList == nil {
		return
	}
	getConn := connList.([]*WsConnection)
	connCount := len(getConn)
	for _, connect := range getConn {
		connect.mutex.Lock()
		if connect.isClosed {
			connect.mutex.Unlock()
			continue
		}
		switch msgType {
		case 1:
			tryPushChanel(connect.RealAlarmChanel, message, false, project, connCount)
		case 2:
			tryPushRealDataChanel(connect.RealDataChanel, message, project, connCount)
		case 3:
			tryPushChanel(connect.RealSystemDataChanel, message, false, project, connCount)
		}
		connect.mutex.Unlock()
	}
}

func WSSend(message interface{}, project string, msgType int) {
	SSEConnManager.GlobalConnManager.PushToAll(message, 1000)
	if msgType == 2 {
		if realMsg, ok := message.(protocol_common.PushRealDataWebData); ok {
			if realMsg.ProjectUuid == "" {
				realMsg.ProjectUuid = project
			}
			enqueueMergedRealData(realMsg)
		} else {
			pushToProjectConns(message, project, msgType)
		}
	} else {
		pushToProjectConns(message, project, msgType)
	}
	if _, IsTrue := protocol_common.ISMNodeProjectConn.Load(project); IsTrue {
		tempPushData := new(protocol_common.ISMNodePushDataStu)
		tempPushData.ProjectUuid = project
		tempPushData.Message = message
		tempPushData.MsgType = msgType
		select {
		case protocol_common.NetworkNodePushDataChanel <- tempPushData:
		default:
		}
		tempPushData = nil
	}
}
func WSSendISMNode(message interface{}, project string, msgType int) {
	SSEConnManager.GlobalConnManager.PushToAll(message, 1000)
	if msgType == 2 {
		if realMsg, ok := message.(protocol_common.PushRealDataWebData); ok {
			if realMsg.ProjectUuid == "" {
				realMsg.ProjectUuid = project
			}
			enqueueMergedRealData(realMsg)
			return
		}
	}
	pushToProjectConns(message, project, msgType)
}
func WSSendAlarmOrOther(message interface{}, project string, msgType int) {
	SSEConnManager.GlobalConnManager.PushToAll(message, 1000)
	// RealData 强制走合并窗，禁止旁路直灌 RealDataChanel。
	if msgType == 2 {
		if realMsg, ok := message.(protocol_common.PushRealDataWebData); ok {
			if realMsg.ProjectUuid == "" {
				realMsg.ProjectUuid = project
			}
			enqueueMergedRealData(realMsg)
			return
		}
	}
	pushToProjectConns(message, project, msgType)
}

func (w *WsServer) Start() (err error) {
	w.listener, err = net.Listen("tcp", w.addr)
	if err != nil {
		fmt.Println("net listen error:", err)
		return
	}
	err = http.Serve(w.listener, w)
	if err != nil {
		fmt.Println("http serve error:", err)
		return
	}
	return nil
}

func RunWebSocketServer() {
	ws := NewWsServer()
	ws.Start()
}
func PthreadSendAlarmQueue() {
	for {
		data, code := protocol_common.PushGAlarmQueue.QueuePull()
		if code == 0 && data != nil {
			project, ok := data.(protocol_common.PushAlarm)
			if ok {
				happen := project.HappenTime
				if happen.IsZero() {
					happen = time.Now()
				}
				msg := project.AlarmMessage
				if project.Cmd != "RealAlarm" && project.AlarmClearMessage != "" {
					msg = project.AlarmClearMessage
				}
				logs.Info("push alarm: time=%s cmd=%s device=%s point=%s value=%s real=%s level=%d msg=%s deviceUuid=%s dataUuid=%s project=%s",
					happen.Format("2006-01-02 15:04:05"),
					project.Cmd,
					project.DeviceName,
					project.DataName,
					project.Value,
					project.RealValue,
					project.AlarmLevel,
					msg,
					project.DeviceUuid,
					project.DataUuid,
					project.ProjectUuid,
				)
				WSSendAlarmOrOther(data, project.ProjectUuid, 1)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}
func PthreadSendDataQueue() {
	for {
		processed := 0
		// 批量 drain，减少固定 sleep 带来的积压延时
		for processed < 32 {
			data, code := protocol_common.GGatherDataQueue.QueuePull()
			if code != 0 || data == nil {
				break
			}
			processed++
			project, ok := data.(protocol_common.PushRealDataWebData)
			if !ok {
				continue
			}
			// 队列仅写库；前端推送由协议侧 WSSend / NotifyRealDataFrontend 单入口负责，避免双推。
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("[RECOVER] WriteRealDataFunc panic: %v\n", r)
					}
				}()
				alarmTask.WriteRealDataFunc(project)
			}()
		}
		if processed == 0 {
			time.Sleep(time.Millisecond * 20)
		}
	}
}
func PthreadSendSystemDataQueue() {
	for {
		data, code := protocol_common.GGatherSystemDataQueue.QueuePull()
		if code == 0 && data != nil {
			project, ok := data.(protocol_common.PushSystemDataWebData)
			if ok {
				WSSendAlarmOrOther(data, project.ProjectUuid, 3)
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func PthreadSendNodeDataQueue() {
	NetorkPushDelay, err1 := config.Int("NetorkPushDelay")
	if err1 != nil {
		NetorkPushDelay = 100
	}
	for {

		for realDataMsg := range protocol_common.NetworkNodePushDataChanel {
			pushData, ok := realDataMsg.(*protocol_common.ISMNodePushDataStu)
			if ok {
				SendToISMNode(pushData.ProjectUuid, pushData.Message, pushData.MsgType)
				time.Sleep(time.Millisecond * time.Duration(NetorkPushDelay))
			}
		}
	}
}

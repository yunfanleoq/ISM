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
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/go-basic/uuid"
	"github.com/gorilla/websocket"
)

type WsServer struct {
	listener net.Listener
	addr     string
	upgrade  *websocket.Upgrader
}

var websocketConnArray sync.Map //make(map[string][]*WsConnection)

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
		ReadBufferSize:    100,
		WriteBufferSize:   500,
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
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("解析错误")
	}
	return conn.ws.WriteMessage(websocket.TextMessage, jsonBytes)
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
			c.WriteToClient(realDataMsg)
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
			c.mutex.Lock()
			DeleteStringElement(c.connId, c.project)
			c.Close()
			c.mutex.Unlock()
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
			c.WriteToClient(realDataMsg)
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
			c.WriteToClient(realDataMsg)
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
func WSSend(message interface{}, project string, msgType int) {
	SSEConnManager.GlobalConnManager.PushToAll(message, 1000)
	connList, ok := websocketConnArray.Load(project)
	if ok && connList != nil {
		getConn := connList.([]*WsConnection)
		for _, connect := range getConn {
			connect.mutex.Lock()
			if connect.isClosed {
				connect.mutex.Unlock()
				continue
			}
			if msgType == 1 {
				select {
				case connect.RealAlarmChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			} else if msgType == 2 {
				select {
				case connect.RealDataChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			} else if msgType == 3 {
				select {
				case connect.RealSystemDataChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			}
			connect.mutex.Unlock()
		}
	}
	if _, IsTrue := protocol_common.ISMNodeProjectConn.Load(project); IsTrue {
		tempPushData := new(protocol_common.ISMNodePushDataStu)
		tempPushData.ProjectUuid = project
		tempPushData.Message = message
		tempPushData.MsgType = msgType
		select {
		case protocol_common.NetworkNodePushDataChanel <- tempPushData:
		case <-time.After(1 * time.Millisecond):
			break
		}
		tempPushData = nil
	}
}
func WSSendISMNode(message interface{}, project string, msgType int) {
	SSEConnManager.GlobalConnManager.PushToAll(message, 1000)
	connList, ok := websocketConnArray.Load(project)
	if ok && connList != nil {
		getConn := connList.([]*WsConnection)
		for _, connect := range getConn {
			connect.mutex.Lock()
			if connect.isClosed {
				connect.mutex.Unlock()
				continue
			}
			if msgType == 1 {
				select {
				case connect.RealAlarmChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			} else if msgType == 2 {
				select {
				case connect.RealDataChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			} else if msgType == 3 {
				select {
				case connect.RealSystemDataChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			}
			connect.mutex.Unlock()
		}
	}
}
func WSSendAlarmOrOther(message interface{}, project string, msgType int) {
	SSEConnManager.GlobalConnManager.PushToAll(message, 1000)
	connList, ok := websocketConnArray.Load(project)
	if ok && connList != nil {
		getConn := connList.([]*WsConnection)
		for _, connect := range getConn {
			connect.mutex.Lock()
			if connect.isClosed {
				connect.mutex.Unlock()
				continue
			}
			if msgType == 1 {
				select {
				case connect.RealAlarmChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			} else if msgType == 2 {
				select {
				case connect.RealDataChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			} else if msgType == 3 {
				select {
				case connect.RealSystemDataChanel <- message:
				case <-time.After(1 * time.Millisecond):
					break
				}
			}
			connect.mutex.Unlock()
		}
	}
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
				fmt.Println("push alarm Data")
				WSSendAlarmOrOther(data, project.ProjectUuid, 1)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}
func PthreadSendDataQueue() {
	for {
		data, code := protocol_common.GGatherDataQueue.QueuePull()
		// fmt.Println("实时数据队列长度:", protocol_common.GGatherDataQueue.QueueLength())
		if code == 0 && data != nil {
			project, ok := data.(protocol_common.PushRealDataWebData)
			if ok {
				// 串行化 SQLite 写入：避免并发 goroutine 竞争 MaxOpenConns=1
				// 使用 recover 防止单次写入 panic 拖垮整个队列处理循环
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("[RECOVER] WriteRealDataFunc panic: %v\n", r)
						}
					}()
					alarmTask.WriteRealDataFunc(project)
				}()
				WSSendAlarmOrOther(data, project.ProjectUuid, 2)
			}
		}
		time.Sleep(time.Millisecond * 50)
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

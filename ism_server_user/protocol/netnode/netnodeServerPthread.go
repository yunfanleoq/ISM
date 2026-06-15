/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-12 17:40:31
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ismnetnode

import (
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/thedevsaddam/gojsonq"
)

var ProjectDataCollects sync.Map
var ProjectDeviceRealDataCollects sync.Map
var ISMNodeList sync.Map

var ISMNodeResponseMap sync.Map

type NodeInfoStu struct {
	NodeName        string
	NodeIpAddress   string
	NodeProjectUUID string
	NodeStatus      int
}

func (c *ISMNetNodeServerHandlerCtl) isServerChanClose() bool {
	select {
	case _, received := <-c.GCloseChan:
		return !received
	default:
	}
	return false
}
func (c *ISMNetNodeServerHandlerCtl) NodeServerCloseChan() {

	isOpen := c.isServerChanClose()
	if !isOpen && c.GCloseChan != nil {
		close(c.GCloseChan)
	}
}
func (c *ISMNetNodeServerCtl) NetNodeServerInit() {

}
func (c *ISMNetNodeServerHandlerCtl) NetNodeServerHandlerInit() {
	c.jsonResRWPoll = &sync.Pool{
		New: func() interface{} {
			return new(ISMNetNodeCmdResponse)
		},
	}
	c.jsonReqRWPoll = &sync.Pool{
		New: func() interface{} {
			return new(ISMNetNodeFormatCmd)
		},
	}
}
func (c *ISMNetNodeServerHandlerCtl) NetNodeSendApi(SProjectUuid, projectUUID string, APIName string, params any) map[string]interface{} {
	var register_cmd = new(ISMNetNodeFormatCmd)
	result := map[string]interface{}{
		"code":     -7,
		"realData": nil,
	}
	defer func() { register_cmd = nil }()
	register_cmd.Cmd = "ResquestApi"
	register_cmd.ProjectUuid = projectUUID
	register_cmd.PackIndex = time.Now().UnixMilli()
	register_cmd.Data.InterfaceName = APIName
	register_cmd.Data.RequestParams = params
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	register_write, err := json.Marshal(register_cmd)
	if err == nil {
		if conn, OK := protocol_common.ISMNodeServerProjectConn.Load(projectUUID); OK {
			WConn := conn.(*websocket.Conn)
			if WConn != nil {
				WConn.WriteMessage(1, register_write)
			} else {
				result["code"] = -5
				return result
			}
			var timeout int = 0
			for {
				responseMsg, isResponse := ISMNodeResponseMap.Load(register_cmd.PackIndex)
				if isResponse {
					responseMsgRead := responseMsg.(*ISMNetNodeCmdResponse)
					result["code"] = responseMsgRead.ResCode
					result["realData"] = responseMsgRead.ResData
					ISMNodeResponseMap.Delete(register_cmd.PackIndex)
					break
				} else {
					timeout++
					if timeout >= 20*2 {
						result["code"] = -8
						break
					}
					time.Sleep(500 * time.Millisecond)
				}
			}
		} else {
			result["code"] = -3
		}
	}
	return result
}
func (c *ISMNetNodeServerHandlerCtl) DealWithOffLine() {
	treeDeviceList, isTrue := ProjectDataCollects.Load(c.ProjectUUID)
	if isTrue {
		if NodeInfo, isOK := ISMNodeList.Load(c.NodeName); isOK {
			NodeInfoData := NodeInfo.(NodeInfoStu)
			treeDeviceListArray := treeDeviceList.([]string)
			var updateTreeListArray []string
			for _, item := range treeDeviceListArray {
				if item != NodeInfoData.NodeProjectUUID {
					updateTreeListArray = append(updateTreeListArray, item)
				}
			}
			ProjectDataCollects.Store(c.ProjectUUID, updateTreeListArray)
		}
	}
	ISMNodeList.Delete(c.NodeName)
}
func (c *ISMNetNodeServerHandlerCtl) dealWithNodedata(msg []byte) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	node_response := c.jsonResRWPoll.Get().(*ISMNetNodeCmdResponse)
	node_request := c.jsonReqRWPoll.Get().(*ISMNetNodeFormatCmd)
	var err error
	c.PingOutTimeSfCount = 0
	readBufRecv := string(msg)

	//检测心跳
	if readBufRecv == "ping" {
		c.ChanelConn.WriteMessage(1, msg)
	} else {
		jsonXP := gojsonq.New().FromString(readBufRecv)
		isResCmd := jsonXP.Find("ResCmd")
		if isResCmd != nil {
			err = json.Unmarshal(msg, &node_response)
			if err != nil {
				return
			}
			//如果是回复包，不用再判断直接把返回值存储到map
			ISMNodeResponseMap.Store(node_response.ResPackIndex, node_response)
		} else {
			err = json.Unmarshal(msg, &node_request)
			if err != nil {
				return
			}
			node_response.ResCmd = node_request.Cmd
			node_response.ResProjectUuid = node_request.ProjectUuid
			node_response.ResNodeName = node_request.NodeName
			node_response.ResPackIndex = node_request.PackIndex

			if node_request.Cmd == "Register" {
				if _, OK := ISMNodeList.Load(node_request.NodeName); OK {
					node_response.ResMsg = "重复的节点名称"
					node_response.ResCode = -1
				} else {

					var singleNode NodeInfoStu
					singleNode.NodeName = node_request.NodeName
					singleNode.NodeIpAddress = c.ChanelConn.RemoteAddr().String()
					singleNode.NodeProjectUUID = node_request.ProjectUuid
					c.NodeName = node_request.NodeName
					ISMNodeList.Store(node_request.NodeName, singleNode)
					node_response.ResMsg = "成功"
					node_response.ResCode = 0
					protocol_common.ISMNodeServerProjectConn.Store(node_request.ProjectUuid, c.ChanelConn)
				}
				jsonBytes, _ := json.Marshal(node_response)
				c.ChanelConn.WriteMessage(websocket.TextMessage, jsonBytes)
			} else if node_request.Cmd == "UpdateDataCollect" {
				treeDeviceList, isTrue := ProjectDataCollects.Load(c.ProjectUUID)
				if !isTrue {
					var treeListArray []string
					treeListArray = append(treeListArray, node_request.ProjectUuid)
					ProjectDataCollects.Store(c.ProjectUUID, treeListArray)
				} else {
					treeDeviceListArray := treeDeviceList.([]string)
					var isfind bool = false
					for _, item := range treeDeviceListArray {
						if node_request.ProjectUuid == item {
							isfind = true
							break
						}
					}
					if !isfind {
						treeDeviceListArray = append(treeDeviceListArray, node_request.ProjectUuid)
					}
					ProjectDataCollects.Store(c.ProjectUUID, treeDeviceListArray)
				}
				protocol_common.ISMNodeProjectDataCollects.Store(node_request.ProjectUuid, &node_request.Data.PushTreeList)

				node_response.ResMsg = "更新成功"
				node_response.ResCode = 0
				jsonBytes, _ := json.Marshal(node_response)
				c.ChanelConn.WriteMessage(websocket.TextMessage, jsonBytes)
			} else if node_request.Cmd == "PushNodeData" {
				node_response.ResMsg = "更新成功"
				node_response.ResCode = 0
				jsonBytes, _ := json.Marshal(node_response)
				c.ChanelConn.WriteMessage(websocket.TextMessage, jsonBytes)
			} else if node_request.Cmd == "PushRealData" {
				ismWebsocket.WSSendISMNode(node_request.Data.RequestParams, c.ProjectUUID, node_request.Data.MessageType)
				node_response.ResMsg = "更新成功"
				node_response.ResCode = 0
				jsonBytes, _ := json.Marshal(node_response)
				c.ChanelConn.WriteMessage(websocket.TextMessage, jsonBytes)
			} else if node_request.Cmd == "SyncDevicesDatas" {

				for _, item := range node_request.Data.ReadDeviceData {
					ProjectDeviceRealDataCollects.Store(item.DeviceUUID, item.RealData)
					for _, data := range item.RealData {
						protocol_common.ISMNodeDeviceRealDataMapByUUID.Store(data.ModelDataUuid, data)
						protocol_common.DeviceRealDataMap.Store(data.DeviceName+"->"+data.Name, data.Value)
					}
				}
				node_response.ResMsg = "更新成功"
				node_response.ResCode = 0

				jsonBytes, _ := json.Marshal(node_response)
				c.ChanelConn.WriteMessage(websocket.TextMessage, jsonBytes)
			}
		}
	}
}
func (c *ISMNetNodeServerHandlerCtl) NetNodeTcpServerRecvHandler(conn *websocket.Conn) {

	for {
		c.rwMutex.Lock()
		// timeout := time.Now().Add(time.Duration(c.PingOutTime) * time.Millisecond)
		// c.ChanelConn.SetReadDeadline(timeout)
		_, msg, err := c.ChanelConn.ReadMessage()
		if err != nil {
			c.PingOutTimeSfCount++
			if c.PingOutTimeSfCount >= c.PingOutTimeCount {
				c.ChanelConn.Close()
				c.PingOutTimeSfCount = 0
				c.NodeServerCloseChan()
				c.DealWithOffLine()
				logs.Error(c.NodeName, " 断开连接")
				c.rwMutex.Unlock()
				break
			}
		} else {
			c.dealWithNodedata(msg)
		}
		c.rwMutex.Unlock()
		// time.Sleep(50 * time.Millisecond)
	}
}

// 创建WsServer结构体
type ISMNetNodeServerCtl struct {
	listener net.Listener
	addr     string
	upgrade  *websocket.Upgrader

	waitGroup          *sync.WaitGroup
	ListenPort         int
	PingHeart          int
	PingOutTime        int
	PingOutTimeCount   int
	NodeName           string
	ProjectUUID        string
	PingOutTimeSfCount int
}

// 创建WsServer结构体
type ISMNetNodeServerHandlerCtl struct {
	waitGroup          *sync.WaitGroup
	ChanelConn         *websocket.Conn
	GCloseChan         chan bool
	PingHeart          int
	PingOutTime        int
	PingOutTimeCount   int
	NodeName           string
	ProjectUUID        string
	PingOutTimeSfCount int
	jsonResRWPoll      *sync.Pool
	jsonReqRWPoll      *sync.Pool
	rwMutex            *sync.Mutex
}

// 初始化WsServer
func (c *ISMNetNodeServerCtl) NewWsServer() {

	c.addr = fmt.Sprintf("0.0.0.0:%d", c.ListenPort)
	c.upgrade = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 500,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				fmt.Println("method is not GET")
				return false
			}
			if r.URL.Path != "/ismnode" {
				fmt.Println("path error")
				return false
			}
			return true
		},
	}
	c.Start()
}

// 处理WebSocket连接
func (c *ISMNetNodeServerCtl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ismnode" {
		httpCode := http.StatusInternalServerError
		reasePhrase := http.StatusText(httpCode)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	conn, err := c.upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	d := &ISMNetNodeServerHandlerCtl{waitGroup: c.waitGroup, ProjectUUID: c.ProjectUUID, PingHeart: c.PingHeart, PingOutTime: c.PingOutTime, PingOutTimeCount: c.PingOutTimeCount}
	d.rwMutex = &sync.Mutex{}
	d.NetNodeServerHandlerInit()
	d.PingOutTimeSfCount = 0
	d.ChanelConn = conn
	d.GCloseChan = make(chan bool)
	logs.Info("ISM Node Connect:", conn.RemoteAddr())
	go d.NetNodeTcpServerRecvHandler(conn)
}

// 启动WebSocket服务器
func (w *ISMNetNodeServerCtl) Start() (err error) {
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

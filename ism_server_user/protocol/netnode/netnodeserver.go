/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-30 10:00:09
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ismnetnode

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"net"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

var NetNodeWg sync.WaitGroup
var NetClientNodeWg sync.WaitGroup
var GNodeClientChan chan bool

func NetNodeTcpServerHandler(conn net.Conn) {

}

func isClientChanClose() bool {
	select {
	case _, received := <-GNodeClientChan:
		return !received
	default:
	}
	return false
}
func NodeClientCloseChan() {

	isOpen := isClientChanClose()
	if !isOpen && GNodeClientChan != nil {
		close(GNodeClientChan)
	}
}

func NetNodeTcpServer() {
	var getProjectLists []models.ProjectLists
	models.Db.Model(&models.ProjectLists{}).Select("*").Where("ID >0").Find(&getProjectLists)

	NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")
	if err != nil {
		logs.Error("不能打开ISMNodeConfig文件")
		return
	}
	for _, item := range getProjectLists {
		NodePort, geterr := NodeConfigConf.Int(item.Uuid + "::NodePort")
		if geterr != nil {
			logs.Info("没有配置项目%s的组网端口", item.Name)
			continue
		}
		PingHeart, geterr := NodeConfigConf.Int(item.Uuid + "::NodePing")
		if geterr != nil {
			logs.Info("没有配置项目%s的心跳间隔", item.Name)
			continue
		}
		PingOutTime, geterr := NodeConfigConf.Int(item.Uuid + "::PingOutTime")
		if geterr != nil {
			logs.Info("没有配置项目%s的心跳超时", item.Name)
			continue
		}
		PingOutTimeCount, geterr := NodeConfigConf.Int(item.Uuid + "::PingOutTimeCount")
		if geterr != nil {
			logs.Info("没有配置项目%s的心跳超时次数", item.Name)
			continue
		}

		d := &ISMNetNodeServerCtl{waitGroup: &NetNodeWg, ProjectUUID: item.Uuid, ListenPort: NodePort, PingHeart: PingHeart, PingOutTime: PingOutTime, PingOutTimeCount: PingOutTimeCount}
		d.NetNodeServerInit()
		go d.NewWsServer()
	}
}

func waitForNodeClientUpdate() {
	//do nothing
	for {
		select {
		case <-GNodeClientChan:
			NetClientNodeWg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func NetNodeTcpClient() {
	var getOutConnectList []models.OutConnectList
	var is_starting = 0

	NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")
	if err != nil {
		logs.Error("不能打开ISMNodeConfig文件")
		return
	}
	for {
		if is_starting == 1 {
			NetClientNodeWg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}

		NodeClientCloseChan()
		GNodeClientChan = make(chan bool)
		is_starting = 0
		models.Db.Model(&models.OutConnectList{}).Select("*").Where("ID >0").Find(&getOutConnectList)

		for _, item := range getOutConnectList {
			NodeName, geterr := NodeConfigConf.String(item.ProjectUuid + "::NodeName")
			if geterr != nil {
				logs.Info("没有配置项目节点名称", item.IpAddress)
			}
			if item.IsEnable == 1 {
				d := &ISMNetNodeClientCtl{Uuid: item.Uuid, waitGroup: &NetClientNodeWg, ConnectPort: item.ConnectPort, ConnectAddr: item.IpAddress, ProjectUuid: item.ProjectUuid, NodeName: NodeName}
				d.NetNodeServerInit(item.PingHeart, item.PingOutTime, item.PingOutTimeCount, item.OutConnectName)
				go d.NetNodeConnectServer()
				NetClientNodeWg.Add(1)
				is_starting = 1
			}
		}
		if is_starting == 0 {
			go waitForNodeClientUpdate()
			NetClientNodeWg.Add(1)
			is_starting = 1
		}
	}

}

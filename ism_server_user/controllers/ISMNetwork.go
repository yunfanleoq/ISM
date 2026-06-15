/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:26
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	ismnode "ISMServer/protocol/netnode"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type ISMNetworkController struct {
	beego.Controller
}

func (c *ISMNetworkController) GetNodeConfig() {

	type NodeConfigStu struct {
		NodeName         string `json:"NodeName"`
		NodePort         int    `json:"NodePort"`
		NodePing         int    `json:"NodePing"`
		PingOutTime      int    `json:"PingOutTime"`
		PingOutTimeCount int    `json:"PingOutTimeCount"`
	}
	var code int
	var NodeConfig NodeConfigStu
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")
		if err != nil {
			code = 0
			NodeConfig.NodeName = "ISMNode001"
			NodeConfig.NodePort = 8066
			NodeConfig.NodePing = 1000
			NodeConfig.PingOutTime = 3000
			NodeConfig.PingOutTimeCount = 3
		} else {
			NodeConfig.NodeName, err = NodeConfigConf.String(ProjectUuid + "::NodeName")
			if err != nil || NodeConfig.NodeName == "" {
				NodeConfig.NodeName = "-"
			}

			NodeConfig.NodePort, err = NodeConfigConf.Int(ProjectUuid + "::NodePort")
			if err != nil {
				NodeConfig.NodePort = 0
			}

			NodeConfig.NodePing, err = NodeConfigConf.Int(ProjectUuid + "::NodePing")
			if err != nil {
				NodeConfig.NodePing = 1000
			}

			NodeConfig.PingOutTime, err = NodeConfigConf.Int(ProjectUuid + "::PingOutTime")
			if err != nil {
				NodeConfig.PingOutTime = 3000
			}
			NodeConfig.PingOutTimeCount, err = NodeConfigConf.Int(ProjectUuid + "::PingOutTimeCount")
			if err != nil {
				NodeConfig.PingOutTimeCount = 3
			}
			code = 0
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
		"data": NodeConfig,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMNetworkController) SetNodeConfig() {
	type NodeConfigStu struct {
		NodeName         string `json:"NodeName"`
		NodePort         int    `json:"NodePort"`
		NodePing         int    `json:"NodePing"`
		PingOutTime      int    `json:"PingOutTime"`
		PingOutTimeCount int    `json:"PingOutTimeCount"`
	}
	var code int
	var NodeConfig NodeConfigStu
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid == "" {
		code = -1
		result := map[string]interface{}{
			"code": code,
		}

		c.Data["json"] = result
		c.ServeJSON() //返回json格式
		return
	}
	err := json.Unmarshal(data, &NodeConfig)
	if err != nil {
		code = -2
		result := map[string]interface{}{
			"code": code,
		}

		c.Data["json"] = result
		c.ServeJSON() //返回json格式
		return
	}

	NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")

	if err != nil || NodeConfig.NodePort == 0 {
		code = -3
	} else {
		NodeConfigConf.Set(ProjectUuid+"::NodeName", NodeConfig.NodeName)
		NodeConfigConf.Set(ProjectUuid+"::NodePort", fmt.Sprintf("%d", NodeConfig.NodePort))
		NodeConfigConf.Set(ProjectUuid+"::NodePing", fmt.Sprintf("%d", NodeConfig.NodePing))
		NodeConfigConf.Set(ProjectUuid+"::PingOutTime", fmt.Sprintf("%d", NodeConfig.PingOutTime))
		NodeConfigConf.Set(ProjectUuid+"::PingOutTimeCount", fmt.Sprintf("%d", NodeConfig.PingOutTimeCount))

		NodeConfigConf.SaveConfigFile("conf/ISMNodeConfig.conf")
		code = 0
	}

	result := map[string]interface{}{
		"code": code,
		"data": NodeConfig,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
	if code == 0 {
		go func() {
			timer := time.NewTimer(3 * time.Second)
			<-timer.C
			os.Exit(0)
		}()

	}
}

func (c *ISMNetworkController) AddOutConnect() {
	var addConnect models.OutConnectList
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addConnect)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			addConnect.ProjectUuid = ProjectUuid
			addConnect.IsEnable = 1
			addConnect.PushTime = 1000
			addConnect.Uuid = uuid.New()

			code = models.AddConnect(addConnect)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "成功添加了出站连接"+addConnect.OutConnectName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了出站连接失败"+addConnect.OutConnectName, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	ismnode.NodeClientCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMNetworkController) EditOutConnect() {
	type updateStu struct {
		Uuid       string                `json:"uuid"`
		UpdateData models.OutConnectList `json:"editData"`
	}
	var code int
	var message string
	var editConnect updateStu
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &editConnect)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			code = models.EditConnect(editConnect.Uuid, editConnect.UpdateData)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "成功编辑了出站连接"+editConnect.UpdateData.OutConnectName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了出站连接失败"+editConnect.UpdateData.OutConnectName, errmsg.JournalLevelInfo, c.Ctx.Input)
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	ismnode.NodeClientCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMNetworkController) DelOutConnect() {
	var code int

	type DelStu struct {
		Uuid string `json:"Uuid"`
	}

	var DelConnect DelStu

	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &DelConnect)
		if err != nil {
			code = -2
		} else {
			code = models.DelConnectOut(ProjectUuid, DelConnect.Uuid)
		}

	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
	}
	ismnode.NodeClientCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMNetworkController) OptOutConnect() {
	type updateStu struct {
		Uuid   string `json:"uuid"`
		Status int    `json:"Status"`
		Name   string `json:"Name"`
	}
	var code int
	var message string
	var editConnect updateStu
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &editConnect)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			code = models.EditConnectStatus(editConnect.Uuid, editConnect.Status)
			var Estring string = ""
			if editConnect.Status == 1 {
				Estring = "使能"
			} else {
				Estring = "禁止"
			}
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, Estring+"出站连接"+editConnect.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "使能出站连接失败"+editConnect.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	ismnode.NodeClientCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMNetworkController) GetConnectOut() {
	var code int
	var list []models.OutConnectList
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		list, code = models.GetConnectOutList(ProjectUuid)
		for key, item := range list {
			if st, OK := ismnode.ISMNodeConnStatus.Load(item.Uuid); OK {
				list[key].ConnectStatus = st.(int)
			}
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
		"list": list,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMNetworkController) GetInConnectList() {
	var code int
	var list []any
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {

		ismnode.ISMNodeList.Range(func(k, v interface{}) bool {
			list = append(list, v)
			return true
		})

	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
		"list": list,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

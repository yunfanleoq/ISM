/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:46
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

type ISMJournal struct {
	beego.Controller
}

func (c *ISMJournal) GetJournalList() {

	var data interface{}
	var params = make(map[string]interface{})
	var code int

	rawData := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(rawData, &params)
	if err != nil {
		code = -1
	} else {
		data, code = models.GetJournalList(params)
	}
	result := map[string]interface{}{
		"code": code,
		"list": data,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func WriteOperationJournal(Authorization string, projectUuid string, content string, JournalLevel int, client *context.BeegoInput) {

	var params models.SystemJournal
	var Name string = ""
	var UserName string = ""

	var tokenCode int
	tokenCode, UserName, _, Name, _ = middleware.JwtToken(Authorization)
	if tokenCode != errmsg.SUCCSE || projectUuid == "" {
		return
	}
	params.Content = content
	params.JournalType = errmsg.JournalOperationType
	params.JournalLevel = JournalLevel
	params.Time = time.Now()
	params.Operator = Name
	params.UserName = UserName
	params.ClientInfo = "协议:" + client.Protocol() + ",请求地址:" + client.URL() + ",方法:" + client.Method() + ",来源:" + client.IP() + ",浏览器:" + client.UserAgent()
	params.ProjectUuid = projectUuid
	models.WriteJournalModel(params)
}
func WriteSystemJournal(username string, Operator string, content string, JournalLevel int, client *context.BeegoInput) {

	var params models.SystemJournal

	params.Content = content
	params.ClientInfo = "协议:" + client.Protocol() + ",请求地址:" + client.URL() + ",方法:" + client.Method() + ",来源:" + client.IP() + ",浏览器:" + client.UserAgent()
	params.JournalType = errmsg.JournalSystemType
	params.JournalLevel = JournalLevel
	params.Time = time.Now()
	params.UserName = username
	params.Operator = Operator
	models.WriteJournalModel(params)
}

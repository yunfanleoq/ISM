/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:49
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	memRis "ISMServer/ISMVendor/memdb"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type MemDbController struct {
	beego.Controller
}

func (c *MemDbController) GetMemData() {
	var Params = make(map[string]interface{})
	var value interface{}
	var code = 0
	var err1 bool
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		data := c.Ctx.Input.RequestBody
		err := json.Unmarshal(data, &Params)
		if err != nil {
			code = errmsg.NOTJSON
		}
		value, err1 = memRis.Memdb.Get(Params["key"].(string))
		if !err1 {
			code = -1
		}
	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code":  code,
		"value": value,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *MemDbController) SetMemData() {

	var Params = make(map[string]interface{})
	var value interface{}
	var code = 0
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		data := c.Ctx.Input.RequestBody
		err := json.Unmarshal(data, &Params)
		if err != nil {
			code = errmsg.NOTJSON
		}
		memRis.Memdb.Put(Params["key"].(string), Params["value"])
	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code":  code,
		"value": value,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

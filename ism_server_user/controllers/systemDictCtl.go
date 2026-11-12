package controllers

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type SystemDictController struct {
	beego.Controller
}

func (c *SystemDictController) write(code int, msg string, list interface{}) {
	c.Data["json"] = map[string]interface{}{
		"code": code,
		"msg":  msg,
		"list": list,
	}
	c.ServeJSON()
}

func (c *SystemDictController) GetSystemDictAll() {
	c.write(errmsg.SUCCSE, "成功", models.LoadSystemDictSnapshot())
}

func (c *SystemDictController) GetSystemDictTypeList() {
	c.write(errmsg.SUCCSE, "成功", models.ListSystemDictTypes())
}

func (c *SystemDictController) SaveSystemDictType() {
	var row models.SystemDictType
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &row); err != nil {
		c.write(-1, "参数错误", nil)
		return
	}
	code, err := models.SaveSystemDictType(row)
	if err != nil {
		c.write(-1, err.Error(), nil)
		return
	}
	c.write(errmsg.SUCCSE, "成功", code)
}

func (c *SystemDictController) DelSystemDictType() {
	var req struct {
		ID uint `json:"ID"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil || req.ID == 0 {
		c.write(-1, "参数错误", nil)
		return
	}
	if err := models.DeleteSystemDictType(req.ID); err != nil {
		c.write(-1, err.Error(), nil)
		return
	}
	c.write(errmsg.SUCCSE, "成功", nil)
}

func (c *SystemDictController) GetSystemDictItemList() {
	var req struct {
		DictCode string `json:"DictCode"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.write(-1, "参数错误", nil)
		return
	}
	c.write(errmsg.SUCCSE, "成功", models.ListSystemDictItems(req.DictCode))
}

func (c *SystemDictController) SaveSystemDictItem() {
	var row models.SystemDictItem
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &row); err != nil {
		c.write(-1, "参数错误", nil)
		return
	}
	if err := models.SaveSystemDictItem(row); err != nil {
		c.write(-1, err.Error(), nil)
		return
	}
	c.write(errmsg.SUCCSE, "成功", nil)
}

func (c *SystemDictController) DelSystemDictItem() {
	var req struct {
		ID uint `json:"ID"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil || req.ID == 0 {
		c.write(-1, "参数错误", nil)
		return
	}
	if err := models.DeleteSystemDictItem(req.ID); err != nil {
		c.write(-1, err.Error(), nil)
		return
	}
	c.write(errmsg.SUCCSE, "成功", nil)
}

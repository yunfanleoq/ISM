package controllers

import (
	"ISMServer/models"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

// XunzhiJfWorkorderController dcim-workorder Mock（JSON 存储 · 占位 API）
type XunzhiJfWorkorderController struct {
	beego.Controller
}

func (c *XunzhiJfWorkorderController) GetWorkorder() {
	woId := c.Ctx.Input.Param(":woId")
	st, code := models.LoadXunzhiJfWorkorderStore()
	if code != 1000 {
		c.Data["json"] = map[string]interface{}{"code": code, "message": "run scripts/xunzhi_jf_workorder_store.py init"}
		c.ServeJSON()
		return
	}
	for _, wo := range st.Workorders {
		if wo.WoId == woId {
			c.Data["json"] = map[string]interface{}{"code": 1000, "workorder": wo}
			c.ServeJSON()
			return
		}
	}
	c.Data["json"] = map[string]interface{}{"code": 404}
	c.ServeJSON()
}

func (c *XunzhiJfWorkorderController) CreateFromAlarm() {
	var req struct {
		AlarmId  string `json:"alarm_id"`
		Assignee string `json:"assignee"`
		Title    string `json:"title"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{"code": -3}
		c.ServeJSON()
		return
	}
	st, code := models.LoadXunzhiJfWorkorderStore()
	if code != 1000 {
		c.Data["json"] = map[string]interface{}{"code": code}
		c.ServeJSON()
		return
	}
	// 最小占位：仅校验告警存在，完整建单请用 Python store demo
	found := false
	for _, a := range st.Alarms {
		if id, ok := a["alarm_id"].(string); ok && id == req.AlarmId {
			found = true
			break
		}
	}
	if !found {
		c.Data["json"] = map[string]interface{}{"code": 404, "message": "alarm not found"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{
		"code":    501,
		"message": "use scripts/xunzhi_jf_workorder_store.py for mock lifecycle; HTTP create pending",
	}
	c.ServeJSON()
}

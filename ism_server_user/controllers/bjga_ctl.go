package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type BjgaController struct {
	beego.Controller
}

func (c *BjgaController) actor() string {
	code, user, _, name, _ := middleware.JwtToken(c.Ctx.Request.Header.Get("Authorization"))
	if code != errmsg.SUCCSE {
		return "未知"
	}
	if name != "" {
		return name
	}
	return user
}

func (c *BjgaController) ok(data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["code"] = 1000
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *BjgaController) fail(code int, msg string) {
	c.Data["json"] = map[string]interface{}{"code": code, "message": msg}
	c.ServeJSON()
}

func (c *BjgaController) readJSON(v interface{}) bool {
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, v); err != nil {
		c.fail(-3, "bad json")
		return false
	}
	return true
}

func (c *BjgaController) Bootstrap() {
	c.ok(map[string]interface{}{"dcim": models.BjgaBootstrap()})
}

func (c *BjgaController) AssetFlow() {
	var req struct {
		No     string `json:"no"`
		Action string `json:"action"`
		Why    string `json:"why"`
	}
	if !c.readJSON(&req) {
		return
	}
	flow, msg := models.BjgaFlowAct(req.No, req.Action, req.Why, c.actor())
	if msg != "" {
		c.fail(422, msg)
		return
	}
	c.ok(map[string]interface{}{"flow": flow})
}

func (c *BjgaController) StocktakeList() {
	var rows []models.BjgaStockDiff
	models.Db.Order("id").Find(&rows)
	c.ok(map[string]interface{}{"diffs": rows})
}

func (c *BjgaController) StocktakeAct() {
	var req struct {
		ID uint `json:"id"`
	}
	if !c.readJSON(&req) {
		return
	}
	row, msg := models.BjgaVerifyDiff(req.ID, c.actor())
	if msg != "" {
		c.fail(404, msg)
		return
	}
	c.ok(map[string]interface{}{"diff": row})
}

func (c *BjgaController) StockBills() {
	var rows []models.BjgaStockBill
	models.Db.Order("id desc").Find(&rows)
	c.ok(map[string]interface{}{"bills": rows})
}

func (c *BjgaController) StockBillCreate() {
	var req struct {
		Kind string `json:"kind"`
		Sn   string `json:"sn"`
		Note string `json:"note"`
	}
	if !c.readJSON(&req) {
		return
	}
	if req.Kind == "" {
		req.Kind = "入库"
	}
	c.ok(map[string]interface{}{"bill": models.BjgaAddBill(req.Kind, req.Sn, req.Note, c.actor())})
}

func (c *BjgaController) MaintList() {
	var rows []models.BjgaMaintPlan
	models.Db.Order("id").Find(&rows)
	c.ok(map[string]interface{}{"plans": rows})
}

func (c *BjgaController) CapacityRecalc() {
	cap, err := models.BjgaRecalcCapacity()
	if err != nil {
		c.fail(500, err.Error())
		return
	}
	models.BjgaAudit(c.actor(), "容量重算", "全站", "", "已按台账重算")
	c.ok(map[string]interface{}{"capacity": cap})
}

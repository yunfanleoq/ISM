package controllers

import "ISMServer/models"

func (c *BjgaController) PowerSettings() {
	if c.Ctx.Input.IsPost() {
		var req struct {
			ID      uint    `json:"id"`
			Current float64 `json:"current"`
		}
		if !c.readJSON(&req) {
			return
		}
		if msg := models.BjgaUpdateSetting(req.ID, req.Current, c.actor()); msg != "" {
			c.fail(404, msg)
			return
		}
	}
	c.ok(map[string]interface{}{"settings": models.BjgaSettings()})
}

func (c *BjgaController) PowerQuality() {
	c.ok(map[string]interface{}{"points": models.BjgaQuality()})
}

func (c *BjgaController) PowerComtrade() {
	if c.Ctx.Input.IsPost() {
		var req struct {
			Name string `json:"name"`
			Text string `json:"text"`
		}
		if !c.readJSON(&req) {
			return
		}
		out := models.BjgaParseComtradeText(req.Name, req.Text)
		if code, ok := out["code"].(int); ok && code != 1000 {
			c.Data["json"] = out
			c.ServeJSON()
			return
		}
		out["code"] = 1000
		c.Data["json"] = out
		c.ServeJSON()
		return
	}
	c.ok(models.BjgaComtradeGet(0))
}

func (c *BjgaController) PowerTickets() {
	if c.Ctx.Input.IsPost() {
		var req struct {
			No     string `json:"no"`
			Action string `json:"action"`
			Seq    int    `json:"seq"`
		}
		if !c.readJSON(&req) {
			return
		}
		item, msg := models.BjgaTicketAct(req.No, req.Action, c.actor(), req.Seq)
		if msg != "" {
			c.fail(422, msg)
			return
		}
		c.ok(map[string]interface{}{"ticket": item, "tickets": models.BjgaTickets()})
		return
	}
	c.ok(map[string]interface{}{"tickets": models.BjgaTickets()})
}

func (c *BjgaController) PowerSim() {
	if c.Ctx.Input.IsPost() {
		var req struct {
			Kva float64 `json:"kva"`
		}
		if !c.readJSON(&req) {
			return
		}
		if req.Kva <= 0 {
			c.fail(422, "容量需要大于 0")
			return
		}
		c.ok(models.BjgaSimUpdate(req.Kva, c.actor()))
		return
	}
	c.ok(models.BjgaSim())
}

package controllers

import "ISMServer/models"

func (c *BjgaController) EnergyPUE() {
	c.ok(models.BjgaPUE())
}

func (c *BjgaController) EnergyMeter() {
	var req struct {
		ID     uint    `json:"id"`
		Value  float64 `json:"value"`
		Reason string  `json:"reason"`
	}
	if !c.readJSON(&req) {
		return
	}
	meter, msg := models.BjgaUpdateMeter(req.ID, req.Value, req.Reason, c.actor())
	if msg != "" {
		c.fail(404, msg)
		return
	}
	c.ok(map[string]interface{}{"meter": meter, "pue": models.BjgaPUE()})
}

func (c *BjgaController) EnergyAirflow() {
	c.ok(models.BjgaAirflow())
}

func (c *BjgaController) EnergyCooling() {
	c.ok(models.BjgaCooling())
}

func (c *BjgaController) EnergyCarbon() {
	c.ok(models.BjgaCarbon())
}

package controllers

import "ISMServer/models"

func (c *BjgaController) OpsBundle() {
	c.ok(models.BjgaOpsBundle())
}

func (c *BjgaController) OpsHaSwitch() {
	var nodes []models.BjgaHaNode
	models.Db.Order("id").Find(&nodes)
	for i := range nodes {
		role := nodes[i].Role
		if role == "主" {
			nodes[i].Role = "备"
		} else if role == "备" {
			nodes[i].Role = "主"
		}
		models.Db.Model(&nodes[i]).Update("role", nodes[i].Role)
	}
	models.BjgaAudit(c.actor(), "热备切换演练", "应用/数据库", "主备", "角色对调")
	models.Db.Order("id").Find(&nodes)
	c.ok(map[string]interface{}{"ha": nodes})
}

func (c *BjgaController) OpsPeace() {
	var req struct {
		To    string                 `json:"to"`
		Items []models.BjgaPeaceItem `json:"items"`
	}
	if !c.readJSON(&req) {
		return
	}
	row, msg := models.BjgaSubmitPeace(c.actor(), req.To, req.Items)
	if msg != "" {
		c.fail(422, msg)
		return
	}
	c.ok(map[string]interface{}{"peace": row, "records": models.BjgaPeaceRecords(12)})
}

func (c *BjgaController) OpsChange() {
	var req struct {
		Title string `json:"title"`
	}
	if !c.readJSON(&req) {
		return
	}
	if req.Title == "" {
		req.Title = "未命名变更"
	}
	row := models.BjgaChange{Title: req.Title, Status: "登记", At: models.BjgaNow(), Owner: c.actor()}
	models.Db.Create(&row)
	models.BjgaAudit(c.actor(), "需求变更", req.Title, "", "登记")
	c.ok(map[string]interface{}{"change": row})
}

func (c *BjgaController) AuditList() {
	var rows []models.BjgaAuditLog
	models.Db.Order("id desc").Limit(100).Find(&rows)
	c.ok(map[string]interface{}{"logs": rows})
}

func (c *BjgaController) AuditPolicy() {
	var p models.BjgaSecPolicy
	if err := models.Db.First(&p).Error; err != nil {
		p = models.BjgaSecPolicy{MinLen: 8, MinClasses: 3, MaxAgeDays: 90, HistoryN: 5, FailLimit: 5, LockMin: 30}
	}
	if c.Ctx.Input.IsPost() {
		var req models.BjgaSecPolicy
		if !c.readJSON(&req) {
			return
		}
		before := p
		p.MinLen, p.MinClasses, p.MaxAgeDays = req.MinLen, req.MinClasses, req.MaxAgeDays
		p.HistoryN, p.FailLimit, p.LockMin = req.HistoryN, req.FailLimit, req.LockMin
		models.Db.Save(&p)
		models.BjgaAudit(c.actor(), "安全策略", "口令", "更新前", "已保存")
		_ = before
	}
	var comp []models.BjgaCompliance
	var lic []models.BjgaLicense
	models.Db.Order("id").Find(&comp)
	models.Db.Order("id").Find(&lic)
	c.ok(map[string]interface{}{"policy": p, "compliance": comp, "licenses": lic})
}

func (c *BjgaController) GovernList() {
	var rules []models.BjgaQualityRule
	var ifaces []models.BjgaInterface
	models.Db.Order("id").Find(&rules)
	models.Db.Order("id").Find(&ifaces)
	c.ok(map[string]interface{}{"rules": rules, "interfaces": ifaces})
}

func (c *BjgaController) SiteList() {
	var videos []models.BjgaVideo
	var doors []models.BjgaDoor
	var comms []models.BjgaComm
	models.Db.Order("id").Find(&videos)
	models.Db.Order("id").Find(&doors)
	models.Db.Order("id").Find(&comms)
	c.ok(map[string]interface{}{"video": videos, "doors": doors, "comms": comms})
}

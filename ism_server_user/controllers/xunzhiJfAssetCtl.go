package controllers

import (
	"ISMServer/models"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

// XunzhiJfAssetController dcim-asset Mock API（读写在 data/xunzhi-jf/dcim_asset_store.json）
type XunzhiJfAssetController struct {
	beego.Controller
}

func (c *XunzhiJfAssetController) RackUPositions() {
	rackId := c.Ctx.Input.Param(":rackId")
	st, code := models.LoadXunzhiJfAssetStore()
	if code != 1000 {
		c.Data["json"] = map[string]interface{}{"code": code, "message": "store missing; run scripts/xunzhi_jf_asset_store.py init"}
		c.ServeJSON()
		return
	}
	var rows []models.XunzhiJfUPosition
	for _, p := range st.UPositions {
		if p.RackId == rackId {
			rows = append(rows, p)
		}
	}
	c.Data["json"] = map[string]interface{}{"code": 1000, "rack_id": rackId, "u_positions": rows}
	c.ServeJSON()
}

func (c *XunzhiJfAssetController) UPositionBind() {
	var req struct {
		AssetId  string `json:"asset_id"`
		RackId   string `json:"rack_id"`
		UStart   int    `json:"u_start"`
		UHeight  int    `json:"u_height"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{"code": -3, "message": "bad json"}
		c.ServeJSON()
		return
	}
	st, code := models.LoadXunzhiJfAssetStore()
	if code != 1000 {
		c.Data["json"] = map[string]interface{}{"code": code}
		c.ServeJSON()
		return
	}
	if req.UHeight < 1 {
		req.UHeight = 1
	}
	end := req.UStart + req.UHeight - 1
	for _, p := range st.UPositions {
		if p.RackId == req.RackId && p.UIndex >= req.UStart && p.UIndex <= end && p.Status == "occupied" {
			c.Data["json"] = map[string]interface{}{"code": 422, "message": "u conflict"}
			c.ServeJSON()
			return
		}
	}
	var asset *models.XunzhiJfAssetRecord
	for i := range st.Assets {
		if st.Assets[i].AssetId == req.AssetId {
			asset = &st.Assets[i]
			break
		}
	}
	if asset == nil {
		c.Data["json"] = map[string]interface{}{"code": 404, "message": "asset not found"}
		c.ServeJSON()
		return
	}
	for i := range st.UPositions {
		p := &st.UPositions[i]
		if p.RackId == req.RackId && p.UIndex >= req.UStart && p.UIndex <= end {
			p.Status = "occupied"
			id := req.AssetId
			p.AssetId = &id
		}
	}
	asset.Status = "mounted"
	asset.RackId = req.RackId
	asset.UStart = req.UStart
	asset.UHeight = req.UHeight
	code = models.SaveXunzhiJfAssetStore(st)
	c.Data["json"] = map[string]interface{}{"code": code, "asset": asset}
	c.ServeJSON()
}

func (c *XunzhiJfAssetController) UPositionUnbind() {
	var req struct {
		AssetId string `json:"asset_id"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{"code": -3}
		c.ServeJSON()
		return
	}
	st, code := models.LoadXunzhiJfAssetStore()
	if code != 1000 {
		c.Data["json"] = map[string]interface{}{"code": code}
		c.ServeJSON()
		return
	}
	var asset *models.XunzhiJfAssetRecord
	for i := range st.Assets {
		if st.Assets[i].AssetId == req.AssetId {
			asset = &st.Assets[i]
			break
		}
	}
	if asset == nil || asset.RackId == "" {
		c.Data["json"] = map[string]interface{}{"code": 404}
		c.ServeJSON()
		return
	}
	end := asset.UStart + asset.UHeight - 1
	for i := range st.UPositions {
		p := &st.UPositions[i]
		if p.RackId == asset.RackId && p.UIndex >= asset.UStart && p.UIndex <= end {
			p.Status = "empty"
			p.AssetId = nil
		}
	}
	asset.Status = "in_stock"
	asset.RackId = ""
	asset.UStart = 0
	code = models.SaveXunzhiJfAssetStore(st)
	c.Data["json"] = map[string]interface{}{"code": code, "asset": asset}
	c.ServeJSON()
}

package controllers

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"
)

func (c *DeviceLibraryController) GetDeviceArchive() {
	type recv struct {
		Uuid []string `json:"uuid"`
	}
	var params recv
	code := errmsg.SUCCSECODE
	var list []models.DeviceArchive
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
		code = errmsg.NOTJSON
		list = []models.DeviceArchive{}
	} else {
		var qerr int
		list, qerr = models.GetDeviceArchives(params.Uuid)
		if qerr != errmsg.SUCCSECODE {
			code = qerr
		}
		if list == nil {
			list = []models.DeviceArchive{}
		}
	}
	c.Data["json"] = map[string]interface{}{
		"code": code,
		"list": list,
	}
	c.ServeJSON()
}

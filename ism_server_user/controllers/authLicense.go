/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:10
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	common "ISMServer/protocol/common"
	"encoding/json"
	"io/ioutil"

	license "ISMServer/license"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthLicenseController struct {
	beego.Controller
}

func (c *AuthLicenseController) SaveLicense() {

	var err_code = -1

	type authStu struct {
		AuthCode string `json:"authCode"`
	}

	var paramsJson authStu

	var mac_code = license.GetPhysicalID()
	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &paramsJson)
	if err != nil {
		err_code = -3
	} else {
		de_auth_code, de_err := license.DePwdCode(paramsJson.AuthCode)
		if de_err != nil {
			err_code = -2
		} else {
			if mac_code != string(de_auth_code[:]) {
				err_code = -4
			} else {
				err2 := ioutil.WriteFile("conf/license.lic", []byte(paramsJson.AuthCode), 0666) //写入文件(字节数组)
				if err2 != nil {
					err_code = -5
				} else {
					common.IsLicense = true
					err_code = 0
				}
			}
		}
	}

	result := map[string]interface{}{
		"code": err_code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *AuthLicenseController) GetLicense() {

	result := map[string]interface{}{
		"islicense": common.IsLicense,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

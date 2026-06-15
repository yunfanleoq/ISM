/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-17 14:42:08
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type DisplayModelController struct {
	beego.Controller
}

func (c *DisplayModelController) ModelList() {

	var getLists []models.DisplayModels
	var code int
	type GetModelList struct {
		DisplayType int `json:"DisplayType"`
	}
	var getModel GetModelList

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	data := c.Ctx.Input.RequestBody
	if ProjectUuid != "" {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		//json数据封装到user对象中
		err := json.Unmarshal(data, &getModel)
		if err != nil {
			code = -1
		} else {
			getLists, code = models.DisplayModelList(ProjectUuid, getModel.DisplayType)
		}

	} else {
		code = -1
		getLists = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": getLists,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) GetUserModelList() {

	var getLists []models.DisplayModels
	var code int

	Authorization := c.Ctx.Request.Header.Get("Authorization")

	tokenCode, UserName, _, _, _ := middleware.JwtToken(Authorization)

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" && tokenCode == errmsg.SUCCSE {
		getLists, code = models.DisplayUserModelList(UserName, ProjectUuid)
	} else {
		code = -1
		getLists = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": getLists,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) ModelAdd() {

	var addModel models.DisplayModels
	var code int

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		//json数据封装到user对象中
		err := json.Unmarshal(data, &addModel)
		if err != nil {
			code = -1
		} else {
			addModel.DisplayModelUid = uuid.New()
			addModel.ProjectUuid = ProjectUuid
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了显示模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
			code = models.DisplayModelAdd(addModel)
		}
	} else {
		code = -4
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) ModelDel() {

	var delModel models.DisplayModels
	var code int

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.DisplayModelDel(delModel.DisplayModelUid)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了显示模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) ModelGet() {

	var getModelJson models.DevicesModel
	var getModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &getModelJson)
	if err != nil {
		code = -1
	} else {
		getModel, code = models.SnmpModelGet(getModelJson.Uuid)
	}

	result := map[string]interface{}{
		"code": code,
		"data": getModel,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) ModelEdit() {

	type updateJson struct {
		Uuid string               `json:"uuid"`
		Data models.DisplayModels `json:"updateData"`
	}
	var update updateJson
	var code int

	dataJson := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &update)
	if err != nil {
		code = -1
	} else {
		code = models.DisplayModelUpdate(update.Uuid, update.Data)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了显示模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) ModelLayerPagerGet() {

	var code int
	var getModel models.DisplayModelLayer

	var getModelJson = make(map[string]interface{})

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &getModelJson)
	if err != nil {
		code = -1
	} else {
		getModel, code = models.DisplayModelLayerPageGet(fmt.Sprintf("%s", getModelJson["pageid"]))
	}

	tempComponents, deErr := base64.StdEncoding.DecodeString(getModel.Components)
	if deErr == nil {
		getModel.Components = string(tempComponents)
	}

	tempLayers, layErr := base64.StdEncoding.DecodeString(getModel.Layer)
	if layErr == nil {
		getModel.Layer = string(tempLayers)
	}

	result := map[string]interface{}{
		"code":  code,
		"layer": getModel,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)

	defer func() {
		gw.Close()
		result = nil
		getModelJson = nil
	}()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func (c *DisplayModelController) ModelLayerGet() {

	var code int
	var getModel []models.DisplayModelLayer
	var getDisplayInfo models.DisplayModels

	var getModelJson = make(map[string]interface{})

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &getModelJson)
	if err != nil {
		code = -1
	} else {
		getModel, code = models.DisplayModelLayerGet(fmt.Sprintf("%s", getModelJson["muid"]))
		getDisplayInfo, _ = models.DisplayModelGet(fmt.Sprintf("%s", getModelJson["muid"]))
	}
	for key, _ := range getModel {
		tempComponents, deErr := base64.StdEncoding.DecodeString(getModel[key].Components)
		if deErr == nil {
			getModel[key].Components = string(tempComponents)
		}

		tempLayers, layErr := base64.StdEncoding.DecodeString(getModel[key].Layer)
		if layErr == nil {
			getModel[key].Layer = string(tempLayers)
		}
		tempLayers = nil
		tempComponents = nil
	}
	result := map[string]interface{}{
		"code":    code,
		"layer":   getModel,
		"Display": getDisplayInfo,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)

	defer func() {
		gw.Close()
		result = nil
		getModel = nil
		getModelJson = nil
	}()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func (c *DisplayModelController) ModelLayerGetByToken() {

	var code int
	var getModel []models.DisplayModelLayer
	var getDisplayInfo models.DisplayModels

	var getModelJson = make(map[string]interface{})

	data := c.Ctx.Input.RequestBody
	var expireAt int64
	var tokenGet string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &getModelJson)
	if err != nil {
		code = -1
	} else {
		var getToken models.UserApiAccessToken
		token := fmt.Sprintf("%s", getModelJson["token"])
		if !protocolCommon.IsLicense {
			code = -5
		} else {
			err := models.Db.Model(&models.UserApiAccessToken{}).Where("access_token = ?", token).First(&getToken).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				code = -4
			} else {
				getModel, code = models.DisplayModelLayerGet(fmt.Sprintf("%s", getModelJson["muid"]))
				getDisplayInfo, _ = models.DisplayModelGet(fmt.Sprintf("%s", getModelJson["muid"]))
				tokenGet, expireAt, _ = middleware.SetToken("admin", "admin", "超级管理员", "662929b7-1d42-41fa-baa2-08b1988e4f54")

				for key, _ := range getModel {
					tempComponents, deErr := base64.StdEncoding.DecodeString(getModel[key].Components)
					if deErr == nil {
						getModel[key].Components = string(tempComponents)
					}

					tempLayers, layErr := base64.StdEncoding.DecodeString(getModel[key].Layer)
					if layErr == nil {
						getModel[key].Layer = string(tempLayers)
					}
				}
			}
		}
	}

	result := map[string]interface{}{
		"code":     code,
		"layer":    getModel,
		"Display":  getDisplayInfo,
		"expireAt": expireAt,
		"token":    tokenGet,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func (c *DisplayModelController) ModelLayerCopy() {

	var code int

	type updateJson struct {
		ModelUuid string `json:"muid"`
		PageUuid  string `json:"pageid"`
	}

	var update updateJson

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &update)
	if err != nil {
		code = -1
	} else {

		code = models.DisplayModelLayerCopy(update.ModelUuid, update.PageUuid)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) ModelLayerSave() {

	var code int

	type updateJson struct {
		ModelUuid string                 `json:"muid"`
		PageUuid  string                 `json:"pageid"`
		Data      map[string]interface{} `json:"saveData"`
	}

	var update updateJson
	var updateStu models.DisplayModelLayer

	update.Data = make(map[string]interface{})

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &update)
	if err != nil {
		code = -1
	} else {
		layer, _ := json.Marshal(update.Data["layer"])
		components, _ := json.Marshal(update.Data["components"])

		updateStu.Layer = base64.StdEncoding.EncodeToString(layer)
		updateStu.Components = base64.StdEncoding.EncodeToString(components)
		code = models.DisplayModelLayerUpdate(update.ModelUuid, update.PageUuid, updateStu)
	}

	result := map[string]interface{}{
		"code": code,
	}
	defer func() {

		result = nil
		update.Data = nil
	}()
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) ModelPageAdd() {

	var code int

	type pageInfoStu struct {
		Uuid     string `json:"modelUuid"`
		Name     string `json:"name"`
		Size     string `json:"size"`
		PageType int    `json:"pageType"`
		IsLogin  int    `json:"isLogin"`
	}

	var pageInfo pageInfoStu

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &pageInfo)
	if err != nil {
		code = -1
	} else {

		code = models.DisplayModelPageAdd(pageInfo.Uuid, pageInfo.Name, pageInfo.Size, pageInfo.PageType, pageInfo.IsLogin)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了显示模型的页面"+pageInfo.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) ModelPageDel() {

	var code int

	type delPageStu struct {
		DisplayerUuid string `json:"modelUuid"`
		PageUuid      string `json:"pageId"`
	}

	var update delPageStu

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &update)
	if err != nil {
		code = -1
	} else {
		code = models.DisplayModelPageDel(update.DisplayerUuid, update.PageUuid)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了显示模型的页面", errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) ModelPageEdit() {

	var code int

	type editPageStu struct {
		DisplayerUuid string `json:"modelUuid"`
		PageUuid      string `json:"pageUuid"`
		Name          string `json:"name"`
	}

	var update editPageStu

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &update)
	if err != nil {
		code = -1
	} else {
		code = models.DisplayModelPageEdit(update.DisplayerUuid, update.PageUuid, update.Name)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了显示模型的页面"+update.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) ModelPageSetHome() {

	var code int

	type updateJson struct {
		ModelUuid string `json:"muid"`
		PageUuid  string `json:"pageid"`
		PageType  int    `json:"PageType"`
	}

	var update updateJson

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &update)
	if err != nil {
		code = -1
	} else {
		code = models.DisplayModelPageSetHome(update.ModelUuid, update.PageUuid, update.PageType)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) ModelGetLoginPage() {

	var code int
	var getModel models.DisplayModelLayer
	var getDisplayInfo models.DisplayModels

	var getModelJson = make(map[string]interface{})

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &getModelJson)
	if err != nil {
		code = -1
	} else {
		getModel, code = models.DisplayModelLayerGetLogin(fmt.Sprintf("%s", getModelJson["muid"]), int(getModelJson["pageType"].(float64)))
		getDisplayInfo, _ = models.DisplayModelGet(fmt.Sprintf("%s", getModelJson["muid"]))

		tempComponents, deErr := base64.StdEncoding.DecodeString(getModel.Components)
		if deErr == nil {
			getModel.Components = string(tempComponents)
		}

		tempLayers, layErr := base64.StdEncoding.DecodeString(getModel.Layer)
		if layErr == nil {
			getModel.Layer = string(tempLayers)
		}
	}

	result := map[string]interface{}{
		"code":    code,
		"layer":   getModel,
		"Display": getDisplayInfo,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DisplayModelController) DisplayImageUpload() {

	type UploadResult struct {
		Code int
		Path string
	}
	var reponse_result UploadResult
	suuid := c.Ctx.Input.Param(":suuid")
	if suuid == "" {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".png":  true,
		".gif":  true,
		".jpg":  true,
		".jpeg": true,
		".bmp":  true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	//创建目录
	uploadDir := models.SystemImagePath
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fileName := h.Filename
	fpath := uploadDir + fileName
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	var updateInfo models.DisplayModels
	updateInfo.DisplayImage = fpath

	reponse_result.Path = fpath
	reponse_result.Code = models.DisplayModelUpdate(suuid, updateInfo)
	c.Data["json"] = reponse_result
	c.ServeJSON()
}
func (c *DisplayModelController) DisplayTempleteList() {

	var getLists []models.SystemDataModel
	var code int64
	result := map[string]interface{}{
		"code":   code,
		"Result": getLists,
	}
	upgradeServer, upgradeServererr := config.String("upgradeServer")
	if upgradeServererr != nil || upgradeServer == "" {
		upgradeServer = "http://www.ismctl.com/ism"
	}
	url := upgradeServer + "/displayTemplete.json"
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		result["code"] = -2
		c.Data["json"] = result
		c.ServeJSON() //返回json格式
		return
	}
	//得到返回结果
	ResBody, _ := ioutil.ReadAll(res.Body)
	//对返回的json数据做解析
	var dataAttr []map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(ResBody, &dataAttr); err == nil {
		result["Result"] = dataAttr
	} else {
		fmt.Println(err)
		result["code"] = -1
	}
	res.Body.Close()
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) DisplayTempleteGet() {

	var params map[string]interface{}
	var code int64
	result := map[string]interface{}{
		"code":   code,
		"Result": nil,
	}
	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//json数据封装到对象中
	err := json.Unmarshal(data, &params)
	if err != nil {
		code = -1
	} else {
		req, _ := http.NewRequest("GET", params["url"].(string), nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			result["code"] = -2
			c.Data["json"] = result
			c.ServeJSON() //返回json格式
			return
		}
		//得到返回结果
		ResBody, _ := ioutil.ReadAll(res.Body)

		body := bytes.TrimPrefix(ResBody, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}
		//对返回的json数据做解析
		var dataAttr map[string]interface{}

		if err := json.Unmarshal(body, &dataAttr); err == nil {
			result["Result"] = dataAttr
		} else {
			fmt.Println(err)
			result["code"] = -1
		}
		res.Body.Close()
	}
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) ModelAddUser() {

	var code int
	type UserStu struct {
		Username string `json:"Username"`
		Name     string `json:"title"`
		Uuid     string `json:"uuid"`
	}
	type AddUserStu struct {
		Displayuuid string    `json:"displayuuid"`
		Users       []UserStu `json:"users"`
	}
	var addUserInfo AddUserStu
	var insertUser []models.DisplayModelsUserList
	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addUserInfo)
		if err != nil {
			code = -1
		} else {
			for _, item := range addUserInfo.Users {
				var insertUsersSingle models.DisplayModelsUserList
				insertUsersSingle.ProjectUuid = ProjectUuid
				insertUsersSingle.DisplayModelUid = addUserInfo.Displayuuid
				insertUsersSingle.User = item.Username
				insertUsersSingle.UserName = item.Name
				insertUsersSingle.UserUuid = item.Uuid
				insertUser = append(insertUser, insertUsersSingle)
			}
			if len(insertUser) > 0 {
				code = models.DisplayModelAddUser(insertUser)
			} else {
				code = -3
			}

		}
	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) ModelDelUser() {

	var code int

	type AddUserStu struct {
		Displayuuid string   `json:"displayuuid"`
		Users       []string `json:"users"`
	}

	data := c.Ctx.Input.RequestBody
	var addUserInfo AddUserStu
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addUserInfo)
		if err != nil {
			code = -1
		} else {
			code = models.DisplayModelDelUser(addUserInfo.Users, addUserInfo.Displayuuid)
		}
	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DisplayModelController) GetModelUsers() {

	var code int

	type AddUserStu struct {
		Displayuuid string `json:"displayuuid"`
	}
	var userList []models.DisplayModelsUserList
	data := c.Ctx.Input.RequestBody
	var addUserInfo AddUserStu
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addUserInfo)
		if err != nil {
			code = -1
		} else {
			code, userList = models.DisplayModelGetUser(addUserInfo.Displayuuid)
		}
	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
		"List": userList,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-26 18:25:54
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

// user模型

package models

import (
	protocolCommon "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"encoding/base64"
	"encoding/json"
	"time"

	createUuid "github.com/go-basic/uuid"
	"gorm.io/gorm"
)

type DisplayModels struct {
	gorm.Model

	Name            string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	ProjectUuid     string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	Description     string `gorm:"type:varchar(250);not null" json:"description" validate:"required,min=2,max=250" label:"描述"`
	DisplayModelUid string `gorm:"index;type:varchar(250);not null" json:"displayUid" validate:"required,min=2,max=250" label:"标识"`
	DisplayImage    string `gorm:"type:varchar(250);" json:"DisplayImage" validate:"required,min=2,max=250" label:"显示的图片"`
	DisplayUserList string `gorm:"type:longtext;" json:"DisplayUserList" validate:"required,min=2,max=250" label:"授权的用户名"`
	DisplayType     int    `gorm:"type:int;not null default:1" json:"DisplayType" validate:"required,min=2,max=250" label:"显示类型"`
}
type DisplayModelsUserList struct {
	gorm.Model
	ProjectUuid     string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	UserUuid        string `gorm:"index;type:varchar(250);not null"  json:"suuid" validate:"required,min=4,max=250" label:"用户UUID"`
	User            string `gorm:"index;type:varchar(250);not null"  json:"User" validate:"required,min=4,max=250" label:"用户账号"`
	UserName        string `gorm:"index;type:varchar(250);not null" json:"UserName" validate:"required" label:"用户名称"`
	DisplayModelUid string `gorm:"index;type:varchar(250);not null" json:"displayUid" validate:"required,min=2,max=250" label:"标识"`
}
type DisplayModelLayer struct {
	gorm.Model

	ModelId           string `gorm:"index;index:idx_tpl_kind_model,priority:1;type:varchar(250);not null"  json:"modelId" validate:"required,min=4,max=250" label:"模型的ID标识"`
	PageName          string `gorm:"type:varchar(250);not null"  json:"PageName" validate:"required,min=4,max=250" label:"页面名称"`
	PageId            string `gorm:"index;type:varchar(250);not null"  json:"PageId" validate:"required,min=4,max=250" label:"页名称"`
	IsHome            int    `gorm:"index;type:int;not null"  json:"IsHome" validate:"required,min=4,max=250" label:"是否为首页"`
	IsLogin           int    `gorm:"index;type:int;"  json:"IsLogin" validate:"required,min=4,max=250" label:"是否为登录页"`
	PageType          int    `gorm:"index;type:int;not null"  json:"PageType" validate:"required,min=4,max=250" label:"页面类型,0：手机 1：PC"`
	Layer             string `gorm:"column:layer;type:longtext;not null" json:"layer" validate:"required" label:"图层信息"`
	Components        string `gorm:"column:components;type:longtext;not null" json:"components" validate:"required" label:"图层中的组件信息"`
	TemplateKind      string `gorm:"index:idx_tpl_kind_model,priority:2;type:varchar(64);default:''" json:"templateKind" label:"模板层级 home|zone|room|cabinet|device"`
	TemplateModelUuid string `gorm:"index:idx_tpl_kind_model,priority:3;type:varchar(250);default:''" json:"templateModelUuid" label:"设备物模型覆盖(仅 device)"`
}

type layerStu struct {
	BackColor        string `json:"backColor"`
	BackgroundImage  string `json:"backgroundImage"`
	WidthHeightRatio string `json:"widthHeightRatio"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
}

// 模型添加
func DisplayModelAdd(params DisplayModels) int {

	var addDisplayModel DisplayModels
	var addDisplayModelLayer DisplayModelLayer
	// var appCount int64

	// err := Db.Model(&DisplayModels{}).Count(&appCount).Error
	// if err != nil {
	// 	appCount = 0
	// }
	// if appCount >= int64(protocolCommon.ConfigAppCount) {
	// 	return errmsg.DISPLAY_MODEL_OUT
	// }
	layerInit := layerStu{BackColor: "#eee", BackgroundImage: "", WidthHeightRatio: "", Width: 867, Height: 765}

	Db.Where("name = ? and project_uuid = ?", params.Name, params.ProjectUuid).First(&addDisplayModel)
	if addDisplayModel.ID != 0 {
		return errmsg.DISPLAY_MODEL_EXIST
	}

	addDisplayModelLayer.ModelId = params.DisplayModelUid
	addDisplayModelLayer.PageId = createUuid.New()
	addDisplayModelLayer.PageName = "demo"
	addDisplayModelLayer.IsHome = 1
	addDisplayModelLayer.PageType = 1
	layer, jsonErr := json.Marshal(layerInit)
	addDisplayModelLayer.Layer = string(layer)

	addDisplayModelLayer.Components = `{"cells": []}`

	if jsonErr != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}

	result := Db.Model(&DisplayModels{}).Create(&params)

	if result.Error != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}
	Db.Model(&DisplayModelLayer{}).Create(&addDisplayModelLayer)

	return errmsg.DISPLAY_MODEL_ADD_SUCCSE

}

// 模型删除（软删除：仅标记 deleted_at，可在回收站恢复，不做物理删除）
func DisplayModelDel(key string) int {

	var delModels DisplayModels
	// var getBandDevice MonitorList
	var delDisplayModelLayer DisplayModelLayer

	// err1 := Db.Model(&MonitorList{}).Where("configuration_uid = ?", key).First(&getBandDevice)
	// if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
	// 	return errmsg.MODEL_HAVED_BAND
	// }

	// 去掉 Unscoped()，GORM 会自动写入 deleted_at 实现软删除，默认查询会自动过滤
	err := Db.Model(&DisplayModels{}).Where("display_model_uid = ?", key).Delete(&delModels).Error
	if err != nil {
		return errmsg.ERROR
	}

	err = Db.Model(&DisplayModelLayer{}).Where("model_id = ?", key).Delete(&delDisplayModelLayer).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// DeletedDisplayModel 回收站列表项（仅暴露展示与回溯所需字段）
type DeletedDisplayModel struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DisplayModelUid string `json:"displayUid"`
	DisplayImage    string `json:"DisplayImage"`
	DisplayType     int    `json:"DisplayType"`
	DeletedAt       string `json:"deletedAt"`
}

// 回收站：列出某项目下已软删除的显示模型（deleted_at 非空）
func DisplayModelDeletedList(ProjectUuid string) ([]DeletedDisplayModel, int) {

	var list []DeletedDisplayModel
	// 必须 Unscoped() 才能查到被软删除的记录
	err := Db.Unscoped().Model(&DisplayModels{}).
		Select("name, description, display_model_uid, display_image, display_type, deleted_at").
		Where("project_uuid = ? AND deleted_at IS NOT NULL", ProjectUuid).
		Order("deleted_at DESC").
		Scan(&list).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return list, errmsg.SUCCSE
}

// 回收站：恢复已软删除的显示模型（清空 deleted_at，连带恢复其图层页面）
func DisplayModelRestore(key string) int {

	// 清空 display_models 的 deleted_at
	err := Db.Unscoped().Model(&DisplayModels{}).
		Where("display_model_uid = ?", key).
		Update("deleted_at", nil).Error
	if err != nil {
		return errmsg.ERROR
	}

	// 连带恢复其所有图层页面
	err = Db.Unscoped().Model(&DisplayModelLayer{}).
		Where("model_id = ?", key).
		Update("deleted_at", nil).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 回收站：彻底删除（物理删除，不可恢复，需回收站二次确认后调用）
func DisplayModelForceDel(key string) int {

	var delModels DisplayModels
	var delDisplayModelLayer DisplayModelLayer

	err := Db.Unscoped().Model(&DisplayModels{}).Where("display_model_uid = ?", key).Delete(&delModels).Error
	if err != nil {
		return errmsg.ERROR
	}

	err = Db.Unscoped().Model(&DisplayModelLayer{}).Where("model_id = ?", key).Delete(&delDisplayModelLayer).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 模型更新
func DisplayModelUpdate(key string, data DisplayModels) int {
	err := Db.Model(&DisplayModels{}).Where("display_model_uid = ?", key).Updates(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 模型获取
func DisplayModelGet(key string) (DisplayModels, int) {

	var getDisplayModels DisplayModels
	err := Db.Model(&DisplayModels{}).Select("*").Where("display_model_uid = ?", key).First(&getDisplayModels).Error
	if err != nil {
		return getDisplayModels, errmsg.ERROR
	}
	return getDisplayModels, errmsg.SUCCSE
}

// 模型获取
func DisplayModelList(ProjectUuid string, GetType int) ([]DisplayModels, int) {

	var getDisplayModels []DisplayModels
	var total int = 0
	if GetType == 1 {
		Db.Model(&DisplayModels{}).Select("*").Where("project_uuid = ? and (display_type = ? or display_type = 0 or display_type is null)", ProjectUuid, GetType).Find(&getDisplayModels)
	} else {
		Db.Model(&DisplayModels{}).Select("*").Where("project_uuid = ? and display_type = ? ", ProjectUuid, GetType).Find(&getDisplayModels)
	}

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil, 0
	}

	return getDisplayModels, total
}

// 模型获取
func DisplayUserModelList(UserName, ProjectUuid string) ([]DisplayModels, int) {

	var getDisplayModels []DisplayModels
	var getDisplayUserModels []DisplayModelsUserList
	var duuid []string
	var total int = 0

	Db.Model(&DisplayModelsUserList{}).Select("*").Where("user = ? ", UserName).Find(&getDisplayUserModels)
	if len(getDisplayUserModels) >= 0 {
		for _, model := range getDisplayUserModels {
			duuid = append(duuid, model.DisplayModelUid)
		}
		Db.Model(&DisplayModels{}).Select("*").Where("display_model_uid in ? and project_uuid = ? ", duuid, ProjectUuid).Find(&getDisplayModels)

		var dberr error
		if dberr == gorm.ErrRecordNotFound {
			return nil, 0
		}

	} else {
		return nil, 0
	}

	return getDisplayModels, total
}

// 模型获取
func DisplayUser11ModelList(ProjectUuid string, uuid string) ([]DisplayModels, int) {

	var getDisplayModels []DisplayModels
	var total int = 0

	// db.Model(&User{}).Select("name, sum(age) as total").Where("name LIKE ?", "group%").Group("name").First(&result)
	// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name` LIMIT 1

	Db.Model(&DisplayModels{}).Select("*").Where("DisplayUserList LIKE %?% and project_uuid = ? ", uuid, ProjectUuid).Find(&getDisplayModels)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil, 0
	}

	return getDisplayModels, total
}

// 模型获取
func DisplayModelLayerGet(muid string) ([]DisplayModelLayer, int) {

	var getDisplayModelLayer []DisplayModelLayer

	result := Db.Raw("SELECT id, created_at, updated_at, deleted_at, model_id, page_name, page_id, is_home, is_login, page_type, layer, components, COALESCE(template_kind,'') AS template_kind, COALESCE(template_model_uuid,'') AS template_model_uuid FROM display_model_layer WHERE model_id = ? AND deleted_at IS NULL", muid).Scan(&getDisplayModelLayer)
	
	if result.Error != nil {
		return getDisplayModelLayer, 0
	}

	return getDisplayModelLayer, 0
}

// 模型获取(仅元数据)：只对首页返回 components，其余页面 components 置空，
// 避免一次性传输全部页面(大屏页面极多时单次响应可达数十/数百 MB)。
// 配合前端按需加载：首屏渲染首页，下钻时按 page_id 单独拉取目标页。
func DisplayModelLayerGetMeta(muid string) ([]DisplayModelLayer, int) {

	var getDisplayModelLayer []DisplayModelLayer

	result := Db.Raw("SELECT id, created_at, updated_at, deleted_at, model_id, page_name, page_id, is_home, is_login, page_type, layer, CASE WHEN is_home = 1 THEN components ELSE '' END AS components, COALESCE(template_kind,'') AS template_kind, COALESCE(template_model_uuid,'') AS template_model_uuid FROM display_model_layer WHERE model_id = ? AND deleted_at IS NULL ORDER BY is_home DESC, id ASC", muid).Scan(&getDisplayModelLayer)

	if result.Error != nil {
		return getDisplayModelLayer, 0
	}

	return getDisplayModelLayer, 0
}

// 单个页面获取
func DisplayModelLayerPageGet(pageid string) (DisplayModelLayer, int) {

	var getDisplayModelLayer DisplayModelLayer
	var total int = 0

	Db.Model(&DisplayModelLayer{}).Select("*").Where("page_id = ? AND deleted_at IS NULL", pageid).First(&getDisplayModelLayer)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return getDisplayModelLayer, 0
	}

	return getDisplayModelLayer, total
}

func DisplayModelLayerGetLogin(muid string, pageType int) (DisplayModelLayer, int) {

	var getDisplayModelLayer DisplayModelLayer

	result := Db.Table("display_model_layer").Where("is_login = ? and model_id = ? and page_type = ? AND deleted_at IS NULL", 1, muid, pageType).First(&getDisplayModelLayer)
	if result.Error != nil {
		return getDisplayModelLayer, -1
	}

	return getDisplayModelLayer, 1
}

// 模型更新
func DisplayModelLayerUpdate(key string, pageid string, data DisplayModelLayer) int {

	err := Db.Model(&DisplayModelLayer{}).Where("model_id = ? and page_id=? AND deleted_at IS NULL", key, pageid).Updates(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 模型拷贝
func DisplayModelLayerCopy(key string, pageid string) int {

	var getPage DisplayModelLayer
	var ComponentsMap map[string]interface{}
	// var ConfigPageCount int64

	// err1 := Db.Model(&DisplayModelLayer{}).Count(&ConfigPageCount).Error
	// if err1 != nil {
	// 	ConfigPageCount = 0
	// }
	// if ConfigPageCount >= int64(protocolCommon.ConfigPageCount) {
	// 	return errmsg.DISPLAY_MODEL_OUT
	// }
	err := Db.Model(&DisplayModelLayer{}).Where("model_id = ? and page_id=? AND deleted_at IS NULL", key, pageid).First(&getPage).Error
	if err != nil {
		return errmsg.ERROR
	}
	getPage.ID = 0
	getPage.IsHome = 0
	getPage.PageId = createUuid.New()
	getPage.PageName = getPage.PageName + " copy"
	getPage.TemplateKind = ""
	getPage.TemplateModelUuid = ""
	getPage.CreatedAt = time.Now()
	getPage.UpdatedAt = time.Now()

	copyTempComponents := getPage.Components
	tempComponents, deErr := base64.StdEncoding.DecodeString(getPage.Components)
	if deErr == nil {
		copyTempComponents = string(tempComponents)
	}
	json.Unmarshal([]byte(copyTempComponents), &ComponentsMap)
	// fmt.Println(ComponentsMap)
	// cells := ComponentsMap["cells"].([]map[string]interface{})
	// for key, _ := range cells {
	// 	cells[key]["id"] = createUuid.New()
	// }
	Components, _ := json.Marshal(ComponentsMap)
	getPage.Components = string(Components)
	result := Db.Model(&DisplayModelLayer{}).Create(&getPage)
	if result.Error != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}
	return errmsg.SUCCSE
}

// 模型页面添加
// templateKind / templateModelUuid 可选：用于层级模板页（home|zone|room|cabinet|device）
// 同 model 下同 (templateKind, templateModelUuid) 已存在时返回已有页信息码 DISPLAY_MODEL_EXIST，并由调用方读 PageId
func DisplayModelPageAdd(uuid string, name string, size string, pageType, islogin int, templateKind, templateModelUuid string) (int, string) {
	var addDisplayModelLayer DisplayModelLayer

	if islogin == 1 && !protocolCommon.IsLicense {
		return errmsg.LOGIN_NO_AUTH, ""
	}

	templateKind = normalizeTemplateKind(templateKind)
	templateModelUuid = normalizeTemplateModelUuid(templateKind, templateModelUuid)

	layerInit := layerStu{BackColor: "#eee", BackgroundImage: "", WidthHeightRatio: "", Width: 867, Height: 765}

	if size == "1" {
		layerInit.Width = 1920
		layerInit.Height = 1080
	} else if size == "2" {
		layerInit.Width = 1060
		layerInit.Height = 520
	} else if size == "3" {
		layerInit.Width = 1440
		layerInit.Height = 900
	} else if size == "4" {
		layerInit.Width = 1366
		layerInit.Height = 768
	} else if size == "5" {
		layerInit.Width = 320
		layerInit.Height = 568
	} else if size == "6" {
		layerInit.Width = 414
		layerInit.Height = 896
	} else if size == "7" {
		layerInit.Width = 390
		layerInit.Height = 844
	} else if size == "8" {
		layerInit.Width = 768
		layerInit.Height = 1024
	} else if size == "9" {
		layerInit.Width = 1024
		layerInit.Height = 1366
	} else if size == "10" {
		layerInit.Width = 820
		layerInit.Height = 1180
	} else if size == "11" {
		layerInit.Width = 768
		layerInit.Height = 1024
	} else if size == "12" {
		layerInit.Width = 375
		layerInit.Height = 667
	} else if size == "13" {
		layerInit.Width = 414
		layerInit.Height = 736
	} else if size == "14" {
		layerInit.Width = 375
		layerInit.Height = 812
	} else if size == "15" {
		layerInit.Width = 360
		layerInit.Height = 740
	} else if size == "16" {
		layerInit.Width = 412
		layerInit.Height = 915
	} else if size == "17" {
		layerInit.Width = 412
		layerInit.Height = 914
	} else {
		layerInit.Width = 1024
		layerInit.Height = 1366
	}

	Db.Model(&DisplayModelLayer{}).Where("page_name = ? and model_id=? AND deleted_at IS NULL", name, uuid).First(&addDisplayModelLayer)
	if addDisplayModelLayer.ID != 0 {
		return errmsg.DISPLAY_MODEL_EXIST, addDisplayModelLayer.PageId
	}

	if templateKind != "" {
		var existTpl DisplayModelLayer
		Db.Model(&DisplayModelLayer{}).Where(
			"model_id = ? AND template_kind = ? AND COALESCE(template_model_uuid,'') = ? AND deleted_at IS NULL",
			uuid, templateKind, templateModelUuid,
		).First(&existTpl)
		if existTpl.ID != 0 {
			return errmsg.DISPLAY_MODEL_EXIST, existTpl.PageId
		}
	}

	addDisplayModelLayer = DisplayModelLayer{}
	addDisplayModelLayer.ModelId = uuid
	addDisplayModelLayer.PageId = createUuid.New()
	addDisplayModelLayer.PageName = name
	addDisplayModelLayer.IsHome = 0
	addDisplayModelLayer.PageType = pageType
	addDisplayModelLayer.IsLogin = islogin
	addDisplayModelLayer.TemplateKind = templateKind
	addDisplayModelLayer.TemplateModelUuid = templateModelUuid
	layer, jsonErr := json.Marshal(layerInit)
	addDisplayModelLayer.Layer = base64.StdEncoding.EncodeToString(layer)

	addDisplayModelLayer.Components = base64.StdEncoding.EncodeToString([]byte(`{"cells": []}`))

	if jsonErr != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED, ""
	}

	result := Db.Model(&DisplayModelLayer{}).Create(&addDisplayModelLayer)

	if result.Error != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED, ""
	}

	return errmsg.DISPLAY_MODEL_ADD_SUCCSE, addDisplayModelLayer.PageId

}

func normalizeTemplateKind(kind string) string {
	switch kind {
	case "home", "zone", "room", "cabinet", "device":
		return kind
	default:
		return ""
	}
}

func normalizeTemplateModelUuid(kind, modelUuid string) string {
	if kind != "device" {
		return ""
	}
	return modelUuid
}

// DisplayModelPageBindTemplate 绑定/改绑/解绑页面的模板角色。
// templateKind 为空表示解除模板角色。force=true 时覆盖同键已有页。
func DisplayModelPageBindTemplate(modelId, pageId, templateKind, templateModelUuid string, force bool) (int, string) {
	templateKind = normalizeTemplateKind(templateKind)
	templateModelUuid = normalizeTemplateModelUuid(templateKind, templateModelUuid)

	var page DisplayModelLayer
	err := Db.Model(&DisplayModelLayer{}).Where("model_id = ? AND page_id = ? AND deleted_at IS NULL", modelId, pageId).First(&page).Error
	if err != nil {
		return errmsg.ERROR, ""
	}

	if templateKind == "" {
		err = Db.Model(&DisplayModelLayer{}).Where("model_id = ? AND page_id = ? AND deleted_at IS NULL", modelId, pageId).
			Updates(map[string]interface{}{"template_kind": "", "template_model_uuid": ""}).Error
		if err != nil {
			return errmsg.ERROR, ""
		}
		return errmsg.SUCCSE, pageId
	}

	var conflict DisplayModelLayer
	Db.Model(&DisplayModelLayer{}).Where(
		"model_id = ? AND template_kind = ? AND COALESCE(template_model_uuid,'') = ? AND page_id <> ? AND deleted_at IS NULL",
		modelId, templateKind, templateModelUuid, pageId,
	).First(&conflict)
	if conflict.ID != 0 {
		if !force {
			return errmsg.DISPLAY_MODEL_EXIST, conflict.PageId
		}
		_ = Db.Model(&DisplayModelLayer{}).Where("id = ?", conflict.ID).
			Updates(map[string]interface{}{"template_kind": "", "template_model_uuid": ""}).Error
	}

	err = Db.Model(&DisplayModelLayer{}).Where("model_id = ? AND page_id = ? AND deleted_at IS NULL", modelId, pageId).
		Updates(map[string]interface{}{
			"template_kind":        templateKind,
			"template_model_uuid": templateModelUuid,
		}).Error
	if err != nil {
		return errmsg.ERROR, ""
	}
	return errmsg.SUCCSE, pageId
}

// DisplayModelTemplateMap 返回大屏各层级模板页映射
func DisplayModelTemplateMap(muid string) map[string]interface{} {
	var rows []DisplayModelLayer
	Db.Model(&DisplayModelLayer{}).
		Select("page_id, page_name, template_kind, template_model_uuid, is_home").
		Where("model_id = ? AND deleted_at IS NULL AND COALESCE(template_kind,'') <> ''", muid).
		Find(&rows)

	out := map[string]interface{}{
		"home":          "",
		"zone":          "",
		"room":          "",
		"cabinet":       "",
		"floor":         "",
		"deviceDefault": "",
		"deviceByModel":  map[string]string{},
		"pages":         []map[string]interface{}{},
	}
	deviceByModel := map[string]string{}
	pages := make([]map[string]interface{}, 0, len(rows))

	for _, r := range rows {
		kind := r.TemplateKind
		pageInfo := map[string]interface{}{
			"pageId":             r.PageId,
			"pageName":           r.PageName,
			"templateKind":       kind,
			"templateModelUuid":  r.TemplateModelUuid,
			"isHome":             r.IsHome,
		}
		pages = append(pages, pageInfo)
		switch kind {
		case "home":
			out["home"] = r.PageId
		case "zone":
			out["zone"] = r.PageId
		case "room":
			out["room"] = r.PageId
		case "cabinet":
			out["cabinet"] = r.PageId
		case "floor":
			out["floor"] = r.PageId
		case "device":
			if r.TemplateModelUuid == "" {
				out["deviceDefault"] = r.PageId
			} else {
				deviceByModel[r.TemplateModelUuid] = r.PageId
			}
		}
	}
	out["deviceByModel"] = deviceByModel
	out["pages"] = pages
	return out
}

// 模型页面删除
func DisplayModelPageDel(uuid string, pageid string) int {
	var delDisplayModelLayer DisplayModelLayer

	err := Db.Unscoped().Model(&DisplayModelLayer{}).Where("model_id = ? and page_id=? AND deleted_at IS NULL", uuid, pageid).Delete(&delDisplayModelLayer).Error
	if err != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}

	return errmsg.DISPLAY_MODEL_ADD_SUCCSE

}

// 模型页面更新
func DisplayModelPageEdit(uuid string, pageid string, name string) int {

	err := Db.Model(&DisplayModelLayer{}).Where("model_id = ? and page_id=? AND deleted_at IS NULL", uuid, pageid).Update("page_name", name).Error
	if err != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}

	return errmsg.DISPLAY_MODEL_ADD_SUCCSE

}

// 模型首页页面更新
func DisplayModelPageSetHome(uuid string, pageid string, pageType int) int {

	err := Db.Model(&DisplayModelLayer{}).Where("model_id = ? and page_type=? AND deleted_at IS NULL", uuid, pageType).Update("is_home", 0).Error
	if err != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}

	err = Db.Model(&DisplayModelLayer{}).Where("model_id = ? and page_id=? and page_type=? AND deleted_at IS NULL", uuid, pageid, pageType).Update("is_home", 1).Error
	if err != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}

	return errmsg.DISPLAY_MODEL_ADD_SUCCSE

}

// 模型用户增加
func DisplayModelAddUser(UserList []DisplayModelsUserList) int {

	err := Db.Model(&DisplayModelsUserList{}).CreateInBatches(&UserList, 10).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE

}

// 模型用户删除
func DisplayModelDelUser(uuid []string, pageid string) int {

	var delDisplayModelUser DisplayModelsUserList
	err := Db.Model(&DisplayModelsUserList{}).Unscoped().Where("user_uuid in ? and display_model_uid=?", uuid, pageid).Delete(&delDisplayModelUser).Error
	if err != nil {
		return errmsg.DISPLAY_MODEL_ADD_FAILED
	}

	return errmsg.DISPLAY_MODEL_ADD_SUCCSE
}

// 模型用户获取
func DisplayModelGetUser(pageid string) (int, []DisplayModelsUserList) {

	var getDisplayModelUser []DisplayModelsUserList
	err := Db.Model(&DisplayModelsUserList{}).Where("display_model_uid=?", pageid).Find(&getDisplayModelUser).Error
	if err != nil {
		return errmsg.ERROR, getDisplayModelUser
	}

	return errmsg.SUCCSE, getDisplayModelUser
}

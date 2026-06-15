/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-05 15:34:10
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package DataInterface

import (
	"ISMServer/models"
	ISMScript "ISMServer/task/ISMScript/func"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	ISMServer "github.com/beego/beego/v2/server/web"
)

var GUrlInterfaceChan chan bool
var urlWg sync.WaitGroup

type UrlInterfaceData struct {
	InterfaceName        string
	ProjectUuid          string
	InterfaceUuid        string
	InterfaceType        int
	InterfaceDataUuid    string
	InterfaceUrl         string
	InterfaceUrlType     string
	InterfaceContent     string
	InterfaceDataContent []string
}

var GetUrlInterfaceData []UrlInterfaceData

func isChanClose() bool {
	select {
	case _, received := <-GUrlInterfaceChan:
		return !received
	default:
	}
	return false
}
func UrlInterfaceCloseChan() {
	isOpen := isChanClose()
	if !isOpen && GUrlInterfaceChan != nil {
		close(GUrlInterfaceChan)
	}
}
func findData(url string) (int, []string) {
	for _, v := range GetUrlInterfaceData {
		if v.InterfaceUrl == url {
			return 1, v.InterfaceDataContent
		}
	}
	return -1, nil
}
func startUrlServer(url []UrlInterfaceData) {
	for _, item := range url {
		// 动态注册自定义 Handler
		ISMServer.Handler(item.InterfaceUrl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := map[string]string{
				"message": "数据模版找不到",
				"status":  "error",
			}
			// 定义正则表达式，匹配 {{}} 中间的字符串
			re := regexp.MustCompile(`\{\{(.*?)\}\}`)
			rdata, djson := findData(r.URL.Path)
			if rdata == -1 || len(djson) == 0 {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(data)
			} else {
				if len(djson[0]) == 0 {
					data["message"] = "数据模版内容为空"
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(data)
				} else {
					// 使用函数动态替换匹配的内容
					replacedText := re.ReplaceAllStringFunc(djson[0], func(match string) string {
						// 去掉 {{ 和 }}，只保留中间的内容
						inner := strings.TrimPrefix(match, "{{")
						inner = strings.TrimSuffix(inner, "}}")

						value := ISMScript.GetDeviceRealData(inner)
						if value == nil {
							value = 0
						}
						_, ok := value.(string)
						if ok {
							return value.(string)
						} else {
							tempValue, ok := value.(float64)
							if ok {
								truncated := math.Trunc(tempValue)
								diff := tempValue - truncated
								istrue := math.Abs(diff) > 1e-9
								if !istrue {
									return fmt.Sprintf("%d", int64(tempValue))
								} else {
									return fmt.Sprintf("%.5f", float64(tempValue))
								}
							} else if tempValue1, ok1 := value.(int); ok1 {
								return fmt.Sprintf("%d", int64(tempValue1))
							} else {
								return fmt.Sprintf("%d", int64(tempValue1))
							}
						}
					})
					w.Header().Set("Content-Type", item.InterfaceUrlType)
					w.Write([]byte(replacedText))
				}
			}
		}))
	}
	//检测协程是否主动退出
	<-GUrlInterfaceChan
	urlWg.Done()
	logs.Info("重新加载API接口")
}
func UrlInterfaceStart() {

	var is_starting = 0
	type UrlData struct {
		Url  string `json:"url"`
		Type string `json:"type"`
	}
	for {
		var TempGetUrlInterfaceData []UrlInterfaceData
		if is_starting == 1 {
			urlWg.Wait()
		}
		UrlInterfaceCloseChan()
		GUrlInterfaceChan = make(chan bool)
		var getData []models.SystemDataInterface
		err := models.Db.Model(&models.SystemDataInterface{}).Where("interface_type = 1").Select("*").Find(&getData).Error
		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		for _, v := range getData {
			if v.InterfaceStatus == 0 {
				continue
			}
			var SUrlInterfaceData UrlInterfaceData
			var edata UrlData
			SUrlInterfaceData.ProjectUuid = v.InterfaceUuid
			SUrlInterfaceData.InterfaceName = v.InterfaceName
			SUrlInterfaceData.InterfaceType = v.InterfaceType
			SUrlInterfaceData.InterfaceDataUuid = v.InterfaceDataUuid
			SUrlInterfaceData.InterfaceContent = v.InterfaceContent
			err := json.Unmarshal([]byte(SUrlInterfaceData.InterfaceContent), &edata)
			if err != nil {
				continue
			}
			SUrlInterfaceData.InterfaceUrl = edata.Url
			SUrlInterfaceData.InterfaceUrlType = edata.Type
			d := strings.Split(v.InterfaceDataUuid, ",")
			for _, uuid := range d {
				var getDataContent models.SystemDataTemplete
				err := models.Db.Model(&models.SystemDataTemplete{}).Where("templete_uuid = ?", uuid).Select("*").Find(&getDataContent).Error
				if err == nil {
					SUrlInterfaceData.InterfaceDataContent = append(SUrlInterfaceData.InterfaceDataContent, getDataContent.TempleteContent)
				}
			}
			TempGetUrlInterfaceData = append(TempGetUrlInterfaceData, SUrlInterfaceData)
		}
		if len(TempGetUrlInterfaceData) <= 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		GetUrlInterfaceData = TempGetUrlInterfaceData
		go startUrlServer(GetUrlInterfaceData)
		urlWg.Add(1)
		is_starting = 1
	}

}

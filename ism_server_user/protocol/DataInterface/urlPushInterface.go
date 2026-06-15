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
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

var GUrlPushInterfaceChan chan bool
var urlPushWg sync.WaitGroup

type UrlPushCtl struct {
	waitGroup   *sync.WaitGroup
	PushUrlData UrlPushInterfaceData
}

type urlPushDataHeader struct {
	HeaderName  string `json:"HeaderName"`
	HeaderValue string `json:"HeaderValue"`
}
type urlPushDataCookie struct {
	CookieName  string `json:"CookieName"`
	CookieValue string `json:"CookieValue"`
}

func (c *UrlPushCtl) InitUrlPush() {

}
func (c *UrlPushCtl) UrlPushPthread() {

	for {
		//检测协程是否主动退出
		select {
		case <-GUrlPushInterfaceChan:
			c.waitGroup.Done()
			return
		default:
			time.Sleep(1 * time.Millisecond) // 降低空转频率
		}
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		Method := "POST"
		//GET
		if c.PushUrlData.InterfaceUrlMethod == 1 {
			Method = http.MethodGet
		} else if c.PushUrlData.InterfaceUrlMethod == 2 {
			//POST
			Method = http.MethodPost
		} else {
			time.Sleep(10 * time.Second)
			continue
		}

		for _, data := range c.PushUrlData.InterfaceDataContent {

			// 定义正则表达式，匹配 {{}} 中间的字符串
			re := regexp.MustCompile(`\{\{(.*?)\}\}`)

			// 使用函数动态替换匹配的内容
			replacedText := re.ReplaceAllStringFunc(data, func(match string) string {
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
			reader := bytes.NewReader([]byte(replacedText))
			req, err := http.NewRequest(Method, c.PushUrlData.InterfaceUrl, reader)
			if err != nil {
				continue
			}
			req.Header.Add("Content-type", c.PushUrlData.InterfaceUrlType)
			// 添加请求头
			for _, v := range c.PushUrlData.InterfaceUrlHeader {
				req.Header.Add(v.HeaderName, v.HeaderValue)
			}

			// 添加cookie
			for _, v := range c.PushUrlData.InterfaceUrlCookie {
				cookie := &http.Cookie{
					Name:  v.CookieName,
					Value: v.CookieValue,
				}
				req.AddCookie(cookie)
			}
			// 发送请求
			_, rerr := client.Do(req)
			if rerr != nil {
				logs.Info("请求%s出错,%v", c.PushUrlData.InterfaceUrl, rerr)
			}
		}
		time.Sleep(time.Duration(c.PushUrlData.InterfaceDataInterval) * time.Millisecond)
	}
}

type UrlPushInterfaceData struct {
	InterfaceName         string
	InterfaceUrl          string
	InterfaceUrlType      string
	InterfaceUrlMethod    int
	InterfaceUrlHeader    []urlPushDataHeader
	InterfaceUrlCookie    []urlPushDataCookie
	InterfaceDataContent  []string
	InterfaceDataInterval int
}

func isUrlPushChanClose() bool {
	select {
	case _, received := <-GUrlPushInterfaceChan:
		return !received
	default:
	}
	return false
}
func UrlPushInterfaceCloseChan() {
	isOpen := isUrlPushChanClose()
	if !isOpen && GUrlPushInterfaceChan != nil {
		close(GUrlPushInterfaceChan)
	}
}
func UrlPushInterfaceStart() {

	var is_starting = 0
	type urlPushData struct {
		Url      string              `json:"url"`
		Type     string              `json:"type"`
		Method   int                 `json:"method"`
		Interval int                 `json:"interval"`
		Header   []urlPushDataHeader `json:"header"`
		Cookie   []urlPushDataCookie `json:"cookie"`
	}
	for {
		if is_starting == 1 {
			urlPushWg.Wait()
		}
		UrlPushInterfaceCloseChan()
		GUrlPushInterfaceChan = make(chan bool)
		var getData []models.SystemDataInterface
		err := models.Db.Model(&models.SystemDataInterface{}).Where("interface_type = 2").Select("*").Find(&getData).Error
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		for _, v := range getData {
			if v.InterfaceStatus == 0 {
				continue
			}
			var SUrlInterfaceData UrlPushInterfaceData
			var edata urlPushData
			SUrlInterfaceData.InterfaceName = v.InterfaceName
			err := json.Unmarshal([]byte(v.InterfaceContent), &edata)
			if err != nil {
				fmt.Println(err)
				continue
			}
			SUrlInterfaceData.InterfaceUrl = edata.Url
			SUrlInterfaceData.InterfaceUrlType = edata.Type
			SUrlInterfaceData.InterfaceUrlMethod = edata.Method
			SUrlInterfaceData.InterfaceUrlHeader = edata.Header
			SUrlInterfaceData.InterfaceUrlCookie = edata.Cookie
			SUrlInterfaceData.InterfaceDataInterval = edata.Interval
			d := strings.Split(v.InterfaceDataUuid, ",")
			for _, uuid := range d {
				var getDataContent models.SystemDataTemplete
				err := models.Db.Model(&models.SystemDataTemplete{}).Where("templete_uuid = ?", uuid).Select("*").Find(&getDataContent).Error
				if err == nil {
					SUrlInterfaceData.InterfaceDataContent = append(SUrlInterfaceData.InterfaceDataContent, getDataContent.TempleteContent)
				}
			}
			if len(SUrlInterfaceData.InterfaceDataContent) == 0 {
				continue
			}
			durl := &UrlPushCtl{waitGroup: &urlPushWg, PushUrlData: SUrlInterfaceData}
			go durl.UrlPushPthread()
			urlPushWg.Add(1)
			is_starting = 1
		}
		if is_starting == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
	}

}

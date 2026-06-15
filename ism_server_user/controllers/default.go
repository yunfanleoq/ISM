/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:27
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
)

var casServerUrl = "https://ids.hfit.edu.cn/authserver"

// model/cas.go
type CasServiceResponse struct {
	XMLName xml.Name `xml:"serviceResponse"`
	Data    struct {
		SFRZH      string `xml:"user"`
		Attributes struct {
			Uid      string `xml:"uid"`
			UserName string `xml:"userName"`
		} `xml:"attributes"`
	} `xml:"authenticationSuccess"`
}

func IsAuthentication(w http.ResponseWriter, r *MainController, casServerUrl string) (bool, *CasServiceResponse) {
	if !hasTicket(r) {
		return false, nil
	}
	localUrl := getLocalUrl(r.Ctx.Request)
	ok, _, res := validateTicket(localUrl, casServerUrl, r)
	if !ok {
		return false, nil
	}
	return true, res
}
func redirectToCasServer(w http.ResponseWriter, r *MainController, casServerUrl string) {
	casServerUrl = casServerUrl + "/login?service=" + getLocalUrl(r.Ctx.Request)
	http.Redirect(w, r.Ctx.Request, casServerUrl, http.StatusFound)
}

/*
获取ticket
*/
func getTicket(r *MainController) string {
	return r.GetString("ticket")
}

/*
判断是否有ticket
*/
func hasTicket(r *MainController) bool {
	t := getTicket(r)
	return len(t) != 0
}

/*
处理并确保路径中只有一个ticket参数
*/
func ensureOneTicketParam(urlParams string) string {
	if len(urlParams) == 0 || !strings.Contains(urlParams, "ticket") {
		return urlParams
	}
	sep := "&"
	params := strings.Split(urlParams, sep)
	newParams := ""
	ticket := ""
	for _, value := range params {
		if strings.Contains(value, "ticket") {
			ticket = value
			continue
		}
		if len(newParams) == 0 {
			newParams = value
		} else {
			newParams = newParams + sep + value
		}
	}
	newParams = newParams + sep + ticket
	return newParams
}

/*
从请求中获取访问路径
*/
func getLocalUrl(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	url := strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
	slice := strings.Split(url, "?")
	if len(slice) > 1 {
		localUrl := slice[0]
		urlParamStr := ensureOneTicketParam(slice[1])
		url = localUrl + "?" + urlParamStr
	}
	CasAppUrl, err := config.String("CasAppUrl")
	if err == nil && len(CasAppUrl) > 0 {
		url = CasAppUrl
	}
	return url
}

func validateTicket(localUrl, casServerUrl string, c *MainController) (bool, error, *CasServiceResponse) {
	casServerUrl = casServerUrl + "/serviceValidate?service=" + localUrl + "&ticket=" + getTicket(c)
	res, err := http.Get(casServerUrl)
	if err != nil {
		return false, err, nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err, nil
	}
	casRes, err := ParseCasUserInfo(data)
	if err != nil {
		return false, err, nil
	}
	if casRes.Data.SFRZH == "" {
		return false, errors.New("authentication failed"), nil
	}
	return true, nil, casRes
}
func ParseCasUserInfo(data []byte) (*CasServiceResponse, error) {
	var casResponse CasServiceResponse
	if err := xml.Unmarshal(data, &casResponse); err != nil {
		return nil, err
	}
	return &casResponse, nil
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	CasAuthUrl, err := config.String("CasAuthUrl")
	if err == nil && len(CasAuthUrl) > 0 {
		casServerUrl = CasAuthUrl
		isAuth, casResponse := IsAuthentication(c.Ctx.ResponseWriter, c, casServerUrl)
		if isAuth {
			code, userInfo := models.CheckLoginAdminUser("admin") //这里必须要写死，后面可以配置一下
			if code == errmsg.LOGIN_SUCCSE {
				token, expireAt, _ := middleware.SetToken(casResponse.Data.Attributes.UserName, "Admin", casResponse.Data.Attributes.Uid, userInfo.Uuid)
				WriteSystemJournal(casResponse.Data.Attributes.UserName, casResponse.Data.Attributes.Uid, "account.Journal.LoginSuccess&"+c.Ctx.Input.IP(), errmsg.JournalLevelInfo, c.Ctx.Input)
				c.Ctx.SetCookie("Authorization", token, expireAt)
				c.TplName = "index.html"
			}
		} else {
			redirectToCasServer(c.Ctx.ResponseWriter, c, casServerUrl)
		}
	} else {
		c.TplName = "index.html"
	}
}

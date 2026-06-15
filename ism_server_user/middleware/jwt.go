/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:58:14
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package middleware

import (
	"ISMServer/utils/errmsg"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/dgrijalva/jwt-go"
)

var jk = "ISMSlat"
var JwtKey = []byte(jk)

type MyClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name" `
	Uuid     string `json:"uuid"`
	jwt.StandardClaims
}

// 生成token
func SetToken(username string, role string, name string, user_uuid string) (string, int64, int) {

	dbtype, typeErr := config.Int("tokenexpirestime")
	if typeErr != nil {
		dbtype = 10
	}
	expireTime := time.Now().Add(time.Duration(dbtype) * time.Hour) //过期时间10小时
	setClaim := MyClaims{
		Username: username,
		Role:     role,
		Name:     name,
		Uuid:     user_uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间（时间戳）
			Issuer:    "ginblog",         //发行者
		},
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaim) //生成token
	token, err := reqClaim.SignedString(JwtKey)                     //转换为字符串
	if err != nil {
		return "", 0, errmsg.ERROR
	}
	return token, expireTime.Unix(), errmsg.SUCCSE
}

// 验证token
func CheckToken(token string) (*MyClaims, int) {
	setToken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if setToken != nil {
		if key, _ := setToken.Claims.(*MyClaims); setToken.Valid {
			return key, errmsg.SUCCSE
		} else {
			return nil, errmsg.ERROR
		}
	} else {
		return nil, errmsg.ERROR
	}
}

// jwt中间件
func JwtToken(Token string) (int, string, string, string, string) {

	key, tCode := CheckToken(Token)
	if tCode == errmsg.ERROR {
		return errmsg.ERROR_TOKEN_WRONG, "", "", "", "" //token验证失败
	}
	if time.Now().Unix() > key.ExpiresAt {
		return errmsg.ERROR_TOKEN_RUNTIME, "", "", "", ""
	}
	return errmsg.SUCCSE, key.Username, key.Role, key.Name, key.Uuid
}

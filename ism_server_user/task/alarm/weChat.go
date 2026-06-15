/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:01:58
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	//发送消息使用导的url
	SendUrl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//获取token使用导的url
	GetToken = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// 定义一个简单的文本消息格式
type weChatMessage struct {
	Touser  string            `json:"touser"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}
type WeChat struct {
	Corpid     string `json:"corpid"`
	Corpsecret string `json:"corpsecret"`
	Agentid    int    `json:"agentid"`
}

type sendMsgError struct {
	ErrCode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 发送消息.msgbody 必须是 API支持的类型
func (t *WeChat) SendMessage(s string, at ...string) (int, string) {
	var sendMessage weChatMessage
	sendMessage.Touser = "@all"
	sendMessage.Msgtype = "text"
	sendMessage.Agentid = t.Agentid
	sendMessage.Text = map[string]string{"content": s}

	sendbuf, _ := json.Marshal(sendMessage)

	resp, err := http.Post(SendUrl+t.GetToken(), "application/json", bytes.NewBuffer(sendbuf))
	if err != nil {
		return -1, ""
	}
	if resp.StatusCode != 200 {
		return -3, ""
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e sendMsgError
	err1 := json.Unmarshal(buf, &e)
	if err1 != nil {
		return -5, ""
	}
	if e.ErrCode != 0 && e.Errmsg != "ok" {
		return -6, e.Errmsg
	}
	return 0, e.Errmsg
}

// 通过corpid 和 corpsecret 获取token
func (t *WeChat) GetToken() string {
	resp, err := http.Get(GetToken + t.Corpid + "&corpsecret=" + t.Corpsecret)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ""
	}
	var at accessToken
	buf, _ := ioutil.ReadAll(resp.Body)
	err1 := json.Unmarshal(buf, &at)
	if err1 != nil || at.AccessToken == "" {
		return ""
	}
	return at.AccessToken
}

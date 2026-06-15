/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:01:46
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DingTalk struct {
	Webhook  string
	Secret   string
	EnableAt bool
	AtAll    bool
}

// SendMessage Function to send message
//
//goland:noinspection GoUnhandledErrorResult
func (t *DingTalk) SendMessage(s string, at ...string) (int, string) {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": s,
		},
	}
	var getResult map[string]interface{}

	if t.EnableAt {
		if t.AtAll {
			if len(at) > 0 {
				return -1, ""
			}
			msg["at"] = map[string]interface{}{
				"isAtAll": t.AtAll,
			}
		} else {
			msg["at"] = map[string]interface{}{
				"atMobiles": at,
				"isAtAll":   t.AtAll,
			}
		}
	} else {
		if len(at) > 0 {
			return -2, ""
		}
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return -3, ""
	}
	resp, err := http.Post(t.getURL(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return -4, ""
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	re := json.Unmarshal(result, &getResult)
	if err != nil && re != nil {
		return -5, ""
	}
	if int(getResult["errcode"].(float64)) != 0 {
		return int(getResult["errcode"].(float64)), getResult["errmsg"].(string)
	}
	return 0, ""
}

func (t *DingTalk) hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (t *DingTalk) getURL() string {
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, t.Secret)
	sign := t.hmacSha256(stringToSign, t.Secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", t.Webhook, timestamp, sign)
	return url
}

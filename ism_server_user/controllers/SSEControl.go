package controllers

import (
	SSEConnManager "ISMServer/utils/SSE"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type SSEController struct {
	web.Controller
}

func (c *SSEController) SSEStream() {
	w := c.Ctx.ResponseWriter
	r := c.Ctx.Request
	// 1. 正确设置 SSE 响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 前端域名
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)

	clientID := c.GetString("client_id")
	if clientID == "" {
		// 方式2：自动生成（备用）
		clientID = fmt.Sprintf("client_%d_%s", time.Now().UnixNano(), r.RemoteAddr)
	}

	// 2. 监听客户端断开
	ctx := c.Ctx.Request.Context()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	groupIDs := []string{"222"}
	// 4. 创建客户端并注册到管理器
	client := SSEConnManager.NewSSEClient(clientID, groupIDs, 100) // 缓冲 100 条消息
	SSEConnManager.GlobalConnManager.Add(client)
	defer SSEConnManager.GlobalConnManager.Remove(clientID) // 连接关闭时清理

	// 3. 循环推送数据（正确格式）
	for {
		select {
		case msg := <-client.MsgChan:
			// 发送消息（SSE 格式：data: 内容\n\n）
			if _, err := fmt.Fprintf(w, "data: %s\n\n", msg); err != nil {
				return
			}
			w.Flush() // 立即推送

		case <-client.Quit:
			// 客户端被移除（如超时）
			return

		case <-ctx.Done():
			// 客户端主动断开连接
			return
		}
	}
}

// 发送带事件类型的 SSE 消息
func (c *SSEController) sendEvent(flusher http.Flusher, eventType string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(c.Ctx.ResponseWriter, "event: %s\ndata: %s\n\n", eventType, dataBytes)
	if err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

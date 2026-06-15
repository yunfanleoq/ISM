// utils/conn_manager.go
package sseConnManager

import (
	"encoding/json"
	"sync"
	"time"
)

// SSEClient 表示一个 SSE 客户端连接
type SSEClient struct {
	UserID   string        // 用户唯一标识
	GroupIDs []string      // 所属分组
	MsgChan  chan []byte   // 消息缓冲通道（避免阻塞）
	LastAct  time.Time     // 最后活动时间（用于超时清理）
	Quit     chan struct{} // 退出信号
}

// NewSSEClient 创建客户端实例
func NewSSEClient(userID string, groupIDs []string, bufSize int) *SSEClient {
	return &SSEClient{
		UserID:   userID,
		GroupIDs: groupIDs,
		MsgChan:  make(chan []byte, bufSize), // 缓冲大小根据业务调整（如 100）
		LastAct:  time.Now(),
		Quit:     make(chan struct{}),
	}
}

// ConnManager 管理所有 SSE 连接（并发安全）
type ConnManager struct {
	mu           sync.RWMutex                     // 读写锁（读多写少优化）
	userClients  map[string]*SSEClient            // userID → SSEClient
	groupClients map[string]map[string]*SSEClient // groupID → userID → SSEClient
	cleanTicker  *time.Ticker                     // 定时清理超时连接
}

var GlobalConnManager *ConnManager

// 初始化连接管理器（5分钟清理一次超时连接）
func init() {
	GlobalConnManager = &ConnManager{
		userClients:  make(map[string]*SSEClient),
		groupClients: make(map[string]map[string]*SSEClient),
		cleanTicker:  time.NewTicker(5 * time.Minute),
	}
	// 启动定时清理协程
	go GlobalConnManager.cleanExpiredClients(10080 * time.Minute) // 30分钟无活动则清理
}

// Add 注册客户端
func (m *ConnManager) Add(client *SSEClient) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 关联用户
	m.userClients[client.UserID] = client

	// 关联分组
	for _, gID := range client.GroupIDs {
		if _, ok := m.groupClients[gID]; !ok {
			m.groupClients[gID] = make(map[string]*SSEClient)
		}
		m.groupClients[gID][client.UserID] = client
	}
}

// Remove 移除客户端
func (m *ConnManager) Remove(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, ok := m.userClients[userID]
	if !ok {
		return
	}

	// 从用户映射删除
	delete(m.userClients, userID)

	// 从分组映射删除
	for _, gID := range client.GroupIDs {
		if users, ok := m.groupClients[gID]; ok {
			delete(users, userID)
			if len(users) == 0 {
				delete(m.groupClients, gID)
			}
		}
	}

	// 关闭通道，触发客户端退出
	close(client.Quit)
}

// PushToUser 向单个用户推送消息
func (m *ConnManager) PushToUser(userID string, msg []byte) {
	m.mu.RLock()
	client, ok := m.userClients[userID]
	m.mu.RUnlock()

	if ok {
		m.pushToClient(client, msg)
	}
}

// PushToGroup 向分组推送消息（并发推送，避免阻塞）
func (m *ConnManager) PushToGroup(groupID string, dataAny interface{}) {
	m.mu.RLock()
	msg, err := json.Marshal(dataAny)
	if err != nil {
		m.mu.RUnlock()
		return // 序列化失败，放弃推送
	}
	// 复制分组内客户端，避免持有锁期间阻塞
	clients := make(map[string]*SSEClient)
	if groupUsers, ok := m.groupClients[groupID]; ok {
		for uid, c := range groupUsers {
			clients[uid] = c
		}
	}
	m.mu.RUnlock()

	// 并发推送（每个客户端一个协程）
	var wg sync.WaitGroup
	for _, c := range clients {
		wg.Add(1)
		go func(client *SSEClient) {
			defer wg.Done()
			m.pushToClient(client, msg)
		}(c)
	}
	wg.Wait()
}

// PushToAll 全量广播（分批处理，控制并发量）
func (m *ConnManager) PushToAll(dataAny interface{}, batchSize int) {
	m.mu.RLock()
	msg, err := json.Marshal(dataAny)
	if err != nil {
		m.mu.RUnlock()
		return // 序列化失败，放弃推送
	}
	// 复制所有客户端
	allClients := make([]*SSEClient, 0, len(m.userClients))
	for _, c := range m.userClients {
		allClients = append(allClients, c)
	}
	m.mu.RUnlock()

	// 分批推送，避免瞬间创建过多协程
	for i := 0; i < len(allClients); i += batchSize {
		end := i + batchSize
		if end > len(allClients) {
			end = len(allClients)
		}
		batch := allClients[i:end]

		var wg sync.WaitGroup
		wg.Add(len(batch))
		for _, c := range batch {
			go func(client *SSEClient) {
				defer wg.Done()
				m.pushToClient(client, msg)
			}(c)
		}
		wg.Wait() // 等待当前批次完成
	}
}

// 向单个客户端推送消息（非阻塞，避免通道满导致阻塞）
func (m *ConnManager) pushToClient(client *SSEClient, msg []byte) {
	select {
	case client.MsgChan <- msg:
		client.LastAct = time.Now() // 更新活动时间
	default:
		// 消息通道满，丢弃或记录日志（根据业务容忍度）
		// log.Printf("client %s msg channel full", client.UserID)
	}
}

// 清理超时客户端
func (m *ConnManager) cleanExpiredClients(timeout time.Duration) {
	for range m.cleanTicker.C {
		m.mu.RLock()
		// 收集超时用户ID
		expiredUsers := make([]string, 0)
		for uid, c := range m.userClients {
			if time.Since(c.LastAct) > timeout {
				expiredUsers = append(expiredUsers, uid)
			}
		}
		m.mu.RUnlock()

		// 批量删除
		for _, uid := range expiredUsers {
			m.Remove(uid)
		}
	}
}

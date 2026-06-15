// ISM MCP Server - 协议核心
// 支持 stdio 和 SSE 传输，实现 MCP 2024-11-05 协议
package mcpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// MCP 协议常量
const (
	ProtocolVersion = "2024-11-05"
	ServerName      = "ism-mcp-server"
	ServerVersion   = "1.0.0"
)

// JSON-RPC 消息类型
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP 协议能力
type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// 工具定义
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema ToolInputSchema `json:"inputSchema"`
	Handler     ToolHandler
}

type ToolInputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

type ToolHandler func(params json.RawMessage) (interface{}, error)

// MCP Server
type Server struct {
	tools    map[string]*Tool
	in       io.Reader
	out      io.Writer
	logger   *log.Logger
}

func NewServer(in io.Reader, out io.Writer) *Server {
	return &Server{
		tools:  make(map[string]*Tool),
		in:     in,
		out:    out,
		logger: log.New(os.Stderr, "[mcp] ", log.LstdFlags),
	}
}

func (s *Server) RegisterTool(tool *Tool) {
	s.tools[tool.Name] = tool
}

func (s *Server) GetTools() []*Tool {
	result := make([]*Tool, 0, len(s.tools))
	for _, t := range s.tools {
		result = append(result, t)
	}
	return result
}

// Run stdio 模式主循环
func (s *Server) Run() error {
	decoder := json.NewDecoder(s.in)
	encoder := json.NewEncoder(s.out)

	for {
		var msg JSONRPCMessage
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				return nil
			}
			s.logger.Printf("decode error: %v", err)
			continue
		}

		response := s.handleMessage(&msg)
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				s.logger.Printf("encode error: %v", err)
				return err
			}
		}
	}
}

func (s *Server) handleMessage(msg *JSONRPCMessage) *JSONRPCMessage {
	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg)
	case "notifications/initialized":
		return nil // 无需响应
	case "tools/list":
		return s.handleToolsList(msg)
	case "tools/call":
		return s.handleToolsCall(msg)
	default:
		s.logger.Printf("unknown method: %s", msg.Method)
		return s.errorResponse(msg.ID, -32601, fmt.Sprintf("Method not found: %s", msg.Method))
	}
}

func (s *Server) handleInitialize(msg *JSONRPCMessage) *JSONRPCMessage {
	result := map[string]interface{}{
		"protocolVersion": ProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]bool{},
		},
		"serverInfo": map[string]string{
			"name":    ServerName,
			"version": ServerVersion,
		},
	}
	resultBytes, _ := json.Marshal(result)
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  resultBytes,
	}
}

func (s *Server) handleToolsList(msg *JSONRPCMessage) *JSONRPCMessage {
	tools := make([]map[string]interface{}, 0, len(s.tools))
	for _, t := range s.GetTools() {
		tools = append(tools, map[string]interface{}{
			"name":        t.Name,
			"description": t.Description,
			"inputSchema": map[string]interface{}{
				"type":       t.InputSchema.Type,
				"properties": t.InputSchema.Properties,
				"required":   t.InputSchema.Required,
			},
		})
	}
	result := map[string]interface{}{"tools": tools}
	resultBytes, _ := json.Marshal(result)
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  resultBytes,
	}
}

func (s *Server) handleToolsCall(msg *JSONRPCMessage) *JSONRPCMessage {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		return s.errorResponse(msg.ID, -32602, "Invalid params")
	}

	tool, ok := s.tools[params.Name]
	if !ok {
		return s.errorResponse(msg.ID, -32602, fmt.Sprintf("Tool not found: %s", params.Name))
	}

	result, err := tool.Handler(params.Arguments)
	if err != nil {
		return s.errorResponse(msg.ID, -32000, err.Error())
	}

	content := []map[string]interface{}{
		{
			"type": "text",
			"text": fmt.Sprintf("%v", result),
		},
	}
	// If result is a string that looks like JSON, try to present it as structured
	if str, ok := result.(string); ok {
		var js interface{}
		if json.Unmarshal([]byte(str), &js) == nil {
			resultBytes, _ := json.Marshal(js)
			content[0]["text"] = string(resultBytes)
		}
	}

	finalResult := map[string]interface{}{
		"content": content,
	}
	resultBytes, _ := json.Marshal(finalResult)
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  resultBytes,
	}
}

func (s *Server) errorResponse(id interface{}, code int, message string) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
}

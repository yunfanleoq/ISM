// ISM MCP Server 入口
// 支持 stdio 模式运行
package main

import (
	mcpserver "ISMServer/mcp-server"
	"os"
)

func main() {
	// stdio 模式：从 stdin 读取，写入 stdout
	server := mcpserver.NewServer(os.Stdin, os.Stdout)

	// 注册所有工具
	mcpserver.RegisterAllTools(server)

	// 启动主循环
	if err := server.Run(); err != nil {
		os.Exit(1)
	}
}

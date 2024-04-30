package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// 设置路由和WebSocket处理函数
	http.HandleFunc("/ws", func(writer http.ResponseWriter, req *http.Request) {
		log.Println("url", req.URL.String())
		log.Println("headers", req.Header)

		ws, err := upgrader.Upgrade(writer, req, nil)
		if err != nil {
			fmt.Println("WebSocket升级失败:", err)
			return
		}
		//defer ws.Close()

		go handleConnection(ws)
	})

	// 启动HTTP服务器监听指定端口（例如8080）
	fmt.Println("WebSocket服务启动，监听端口8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil { // 使用默认的HTTP服务器处理请求和WebSocket连接
		fmt.Println("HTTP服务器启动失败:", err) // 如果启动失败，打印错误信息给用户看或记录日志中。这里只是简单打印到控制台。实际生产环境中需要更复杂的错误处理逻辑。
	}
}

// handleConnection 处理WebSocket连接
func handleConnection(ws *websocket.Conn) {
	defer ws.Close() // 确保连接关闭时执行清理操作
	for {
		// 读取客户端消息
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("读取消息失败:", err)
			break // 如果读取失败，退出循环
		}
		fmt.Printf("收到消息: %s\n", p)

		// 向客户端发送消息
		err = ws.WriteMessage(messageType, []byte("I am WebSocket Test Server: 已收到消息 \""+string(p)+"\""))
		if err != nil {
			fmt.Println("发送消息失败:", err)
			break // 如果发送失败，退出循环
		}
	}
}

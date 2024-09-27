package server

import (
	checkinMsg "checkin/server/msg"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type CmdType string

const (
	CmdReg         = "reg"
	CmdSendlog     = "sendlog"
	CmdSenduser    = "senduser"
	CmdGetuserlist = "getuserlist"
	CmdGetuserinfo = "getuserinfo"
	CmdSetuserinfo = "setuserinfo"
	CmdDeleteuser  = "deleteuser"
	CmdEnableuser  = "enableuser"
)

// 创建一个 WebSocket 升级器
var upgrader = websocket.Upgrader{
	// 允许跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 命令处理函数类型
type cmdHandlerFunc func(conn *websocket.Conn, msg []byte)

type cmdReceiverFunc func(conn *websocket.Conn, msg []byte)

// 定义不同命令的处理函数
var handlers = map[string]cmdHandlerFunc{
	"reg":      receiveReg,
	"sendlog":  receiveSendlog,
	"senduser": receiveSenduser,
}

var receives = map[string]cmdReceiverFunc{
	"getuserlist": receiveGetuserlist,
	"getuserinfo": receiveGetuserinfo,
	"setuserinfo": receiveSetuserinfo,
	"deleteuser":  receiveDeleteuser,
	"enableuser":  receiveEnableuser,
}

func sendData(conn *websocket.Conn, data interface{}) {
	log.Debugf("发送数据: %+v", data)
	if err := conn.WriteJSON(data); err != nil {
		log.Println("Write error:", err)
		conn.Close()
		sn := clientsByConn[conn]
		delete(clientsBySn, sn)
		delete(clientsByConn, conn)
	}
}

// WebSocket 处理函数
var clientsByConn = make(map[*websocket.Conn]string)
var clientsBySn = make(map[string]*websocket.Conn)

// WebSocket 处理函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close() // 在函数返回前关闭连接
	// clients[conn] = true
	log.Println("Client connected")

	// 监听来自客户端的消息
	for {
		// 读取客户端发送的消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// 打印收到的消息
		fmt.Printf("Received message: %s\n", message)
		var baseMsg checkinMsg.BaseMessage
		if err := json.Unmarshal(message, &baseMsg); err != nil {
			log.Println("Base JSON unmarshal error:", err)
			continue
		}
		if baseMsg.Cmd != "" {
			// 调用处理函数，根据 cmd 字段执行不同逻辑
			if handler, ok := handlers[baseMsg.Cmd]; ok {
				handler(conn, message)
			} else {
				log.Printf("Unknown command: %s", baseMsg.Cmd)
				sendData(conn, checkinMsg.WSResponse{
					Ret:    "failure",
					Result: false,
					Reason: 1002, // 未知命令错误代码
				})
			}
		} else if baseMsg.Ret != "" {
			log.Println("处理考勤机回复")
			if receiver, ok := receives[baseMsg.Ret]; ok {
				receiver(conn, message)
			} else {
				log.Printf("Unknown command: %s", baseMsg.Ret)
			}
		}

	}
}

func Run() {
	http.HandleFunc("/user/enable", UserStatusHandle)
	http.HandleFunc("/user", UserHandle)
	// 设置路由，定义 WebSocket 处理路径
	http.HandleFunc("/", wsHandler)
	// 启动 HTTP 服务器并监听端口
	log.Println("Server started at :7788")

	if err := http.ListenAndServe(":7788", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

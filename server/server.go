package server

import (
	"checkin/query"
	"checkin/query/model"
	checkinMsg "checkin/server/msg"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CmdType string

// 创建一个 WebSocket 升级器
var upgrader = websocket.Upgrader{
	// 允许跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 命令处理函数类型
type cmdHandlerFunc func(conn *websocket.Conn, msg []byte)

// 定义不同命令的处理函数
var handlers = map[string]cmdHandlerFunc{
	"reg":      handleReg,
	"sendlog":  handleSendlog,
	"senduser": handleSenduser,
}

// 处理 "reg" 命令
func handleReg(conn *websocket.Conn, msg []byte) {
	// log.Printf("Received registration from device: %s", msg.Sn)
	var regMsg checkinMsg.RegMessage
	if err := json.Unmarshal(msg, &regMsg); err != nil {
		log.Println("RegMessage unmarshal error:", err)
		// 返回成功响应
		sendResponse(conn, checkinMsg.WSResponse{
			Ret:    "reg",
			Result: false,
			Reason: 1,
		})
		return
	}
	checkinMachine, err := query.UserCheckinMachine.WithContext(context.Background()).
		Where(query.UserCheckinMachine.Sn.Eq(regMsg.Sn)).
		First()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
		jsonData, err := json.Marshal(regMsg.Devinfo)
		if err != nil {
			return
		}
		query.UserCheckinMachine.WithContext(context.Background()).
			Create(&model.UserCheckinMachine{
				Sn:      regMsg.Sn,
				Devinfo: string(jsonData),
			})
		sendResponse(conn, checkinMsg.WSResponse{
			Ret:       "reg",
			Result:    true,
			Cloudtime: time.Now().Format(time.DateTime),
		})
		return
	}
	jsonData, err := json.Marshal(regMsg.Devinfo)
	if err != nil {
		return
	}
	query.UserCheckinMachine.WithContext(context.Background()).Where(query.UserCheckinMachine.Sn.Eq(checkinMachine.Sn)).Update(query.UserCheckinMachine.Devinfo, jsonData)
	sendResponse(conn, checkinMsg.WSResponse{
		Ret:       "reg",
		Result:    true,
		Cloudtime: time.Now().Format(time.DateTime),
	})
}

func handleSendlog(conn *websocket.Conn, msg []byte) {
	var regMsg checkinMsg.SendlogMessage
	if err := json.Unmarshal(msg, &regMsg); err != nil {
		log.Println("RegMessage unmarshal error:", err)
		// 返回成功响应
		sendResponse(conn, checkinMsg.WSResponse{
			Ret:    "sendlog",
			Result: false,
			Reason: 1,
		})
		return
	}
	sendResponse(conn, checkinMsg.WSResponse{
		Ret:       "sendlog",
		Result:    true,
		Count:     1,
		Logindex:  0,
		Cloudtime: time.Now().Format(time.DateTime),
		Access:    1,
	})
}

func handleSenduser(conn *websocket.Conn, msg []byte) {
	var senduserMsg checkinMsg.SenduserMessage
	if err := json.Unmarshal(msg, &senduserMsg); err != nil {
		log.Println("RegMessage unmarshal error:", err)
		// 返回成功响应
		sendResponse(conn, checkinMsg.WSResponse{
			Ret:    "senduser",
			Result: false,
			Reason: 1,
		})
		return
	}

	err := query.UserCheckinMachineInfo.WithContext(context.Background()).Create(&model.UserCheckinMachineInfo{
		Sn:        senduserMsg.Sn,
		Enrollid:  senduserMsg.Enrollid,
		Name:      senduserMsg.Name,
		Backupnum: senduserMsg.Backupnum,
		Record:    senduserMsg.Record,
	})
	if err != nil {
		sendResponse(conn, checkinMsg.WSResponse{
			Ret:    "senduser",
			Result: false,
			Reason: 1,
		})
		return
	}
	sendResponse(conn, checkinMsg.WSResponse{
		Ret:       "senduser",
		Result:    true,
		Cloudtime: time.Now().Format(time.DateTime),
	})
}

// 发送 WebSocket 响应
func sendResponse(conn *websocket.Conn, response checkinMsg.WSResponse) {
	if err := conn.WriteJSON(response); err != nil {
		log.Println("Write error:", err)
	}
}

// WebSocket 处理函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	defer conn.Close() // 在函数返回前关闭连接

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
		// 调用处理函数，根据 cmd 字段执行不同逻辑
		if handler, ok := handlers[baseMsg.Cmd]; ok {
			handler(conn, message)
		} else {
			log.Printf("Unknown command: %s", baseMsg.Cmd)
			sendResponse(conn, checkinMsg.WSResponse{
				Ret:    "failure",
				Result: false,
				Reason: 1002, // 未知命令错误代码
			})
		}

	}
}

func Run() {
	// 设置路由，定义 WebSocket 处理路径
	http.HandleFunc("/", wsHandler)

	// 启动 HTTP 服务器并监听端口
	log.Println("Server started at :7788")
	if err := http.ListenAndServe(":7788", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package server

import (
	"checkin/query"
	checkinMsg "checkin/server/msg"
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func handleGetuserlist(conn *websocket.Conn, stn bool) {
	sendData(conn, checkinMsg.GetuserlistMessage{
		Cmd: CmdGetuserlist,
		Stn: stn,
	})
}

func receiveGetuserlist(conn *websocket.Conn, msg []byte) {
	var response checkinMsg.GetuserlistResponse
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return
	}
	if !response.Result {
		return
	}

	if response.Count == 0 {
		log.Println("获取用户列表结束")
		return
	}

	handleGetuserlist(conn, false)

}

func handleGetuserinfo(conn *websocket.Conn, msg checkinMsg.GetuserinfoMessage) {
	sendData(conn, msg)
}

func receiveGetuserinfo(conn *websocket.Conn, msg []byte) {

}

// HandleSetUserInfoAll 向所有设备下发
func handleSetUserInfoAll(msg checkinMsg.SetuserinfoMessage) {

	machines, err := query.UserCheckinMachine.WithContext(context.Background()).Find()
	if err != nil {
		log.Errorf("Error query machines: %v", err)
		return
	}
	for _, machine := range machines {
		conn, exists := clientsBySn[machine.Sn]
		if !exists {
			log.Warnf("下发用户失败，设备[%s]未连接", machine.Sn)
			continue
		}
		// clientsBySn
		handleSetuserinfo(conn, msg)
	}
}

// handleSetuserinfo设置
func handleSetuserinfo(conn *websocket.Conn, msg checkinMsg.SetuserinfoMessage) {
	sendData(conn, msg)
}

func receiveSetuserinfo(conn *websocket.Conn, msg []byte) {

	var response checkinMsg.WSResponse
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return
	}

	if !response.Result {
		if sn := clientsByConn[conn]; sn != "" {
			log.Warnf("对设备[%s]下发用户信息失败: %v", sn, response.Msg)
		} else {
			log.Println("Error set user info:", response.Msg)
		}
	} else {
		log.Printf("对设备[%s]下发用户信息[%d]成功", response.Sn, response.Enrollid)
	}

}

// HandleSetUserInfoAll 向所有设备下发
func handleDeleteuserAll(msg checkinMsg.DeleteuserMessage) {

	machines, err := query.UserCheckinMachine.WithContext(context.Background()).Find()
	if err != nil {
		log.Errorf("Error query machines: %v", err)
		return
	}
	for _, machine := range machines {
		conn, exists := clientsBySn[machine.Sn]
		if !exists {
			log.Warnf("删除用户失败，设备[%s]未连接", machine.Sn)
			continue
		}
		// clientsBySn
		handleDeleteuser(conn, msg)
	}
}

// handleDeleteuser 处理删除用户信息
func handleDeleteuser(conn *websocket.Conn, msg checkinMsg.DeleteuserMessage) {
	sendData(conn, msg)
}

func receiveDeleteuser(conn *websocket.Conn, msg []byte) {
	var response checkinMsg.DeleteuserResponse
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return
	}
	sn := clientsByConn[conn]
	if !response.Result {
		log.Errorf("设备[%s]删除用户信息失败, 原因:%d", sn, response.Reason)
	} else {
		log.Printf("设备[%s]删除用户信息[%d]成功", sn, response.Enrollid)
		if response.Backupnum == 13 {
			// 全部删除
			query.UserCheckinMachineInfo.WithContext(context.Background()).
				Where(query.UserCheckinMachineInfo.Enrollid.Eq(response.Enrollid)).Delete()
		} else {
			query.UserCheckinMachineInfo.WithContext(context.Background()).
				Where(
					query.UserCheckinMachineInfo.Enrollid.Eq(response.Enrollid),
					query.UserCheckinMachineInfo.Backupnum.Eq(response.Backupnum),
				).Delete()
		}
	}
}

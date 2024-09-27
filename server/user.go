package server

import (
	"checkin/query"
	"checkin/query/model"
	checkinMsg "checkin/server/msg"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func handleGetuserlistRandomDevice() {
	device, err := query.CheckinDevice.WithContext(context.Background()).First()
	if err != nil {
		log.Debugf("Error query device info")
		return
	}
	conn, exists := clientsBySn[device.Sn]
	if exists {
		handleGetuserlist(conn, true)
	}
}

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

// HandleSetUserInfoAll 向所有设备下发,
// 对于没有用户信息，目前为先添加数据，处理失败则删除未登记成功的数据(status为-1)，成功则更新状态，存在用户信息，更新信息
// 现在更新信息会先更新数据库
// TODO 后续数据保存于缓存/Redis，得到下发处理结果后再进行数据库更新
func handleSetUserInfoAll(msg checkinMsg.SetuserinfoMessage) {

	machines, err := query.CheckinDevice.WithContext(context.Background()).Find()
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
	sn, exists := clientsByConn[conn]
	if !exists {
		log.Warn("连接已断开")
		return
	}
	ctx := context.Background()
	userInfo, err := query.CheckinDeviceUser.WithContext(ctx).
		Where(query.CheckinDeviceUser.Sn.Eq(sn),
			query.CheckinDeviceUser.Enrollid.Eq(msg.Enrollid),
			query.CheckinDeviceUser.Backupnum.Eq(msg.Backupnum),
		).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("Error query user info")

		} else {
			return
		}
	}
	recordStr := ""
	switch v := msg.Record.(type) {
	case float64:
		recordStr = fmt.Sprintf("%.f", v) // 将数字转换为字符串
	case string:
		recordStr = v // 直接使用字符串
	default:
		recordStr = "" // 或者处理其他情况
	}
	if userInfo == nil {

		err = query.CheckinDeviceUser.WithContext(ctx).Create(&model.CheckinDeviceUser{
			Sn:        sn,
			Enrollid:  msg.Enrollid,
			Name:      msg.Name,
			Backupnum: msg.Backupnum,
			Status:    -1,
			Record:    recordStr,
		})
		if err != nil {
			log.Errorf("Error create user info: %v", err)
		}
	} else {
		query.DB.Model(&userInfo).Updates(model.CheckinDeviceUser{
			Name:   msg.Name,
			Record: recordStr,
		})
	}
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
			query.CheckinDeviceUser.WithContext(context.Background()).Where(
				query.CheckinDeviceUser.Sn.Eq(response.Sn),
				query.CheckinDeviceUser.Enrollid.Eq(response.Enrollid),
				query.CheckinDeviceUser.Backupnum.Eq(response.Backupnum),
				query.CheckinDeviceUser.Status.Eq(-1),
			).Delete()
		} else {
			log.Println("Error set user info:", response.Msg)
		}
	} else {
		log.Printf("对设备[%s]下发用户信息[%d]成功", response.Sn, response.Enrollid)
		query.CheckinDeviceUser.WithContext(context.Background()).Where(
			query.CheckinDeviceUser.Sn.Eq(response.Sn),
			query.CheckinDeviceUser.Enrollid.Eq(response.Enrollid),
			query.CheckinDeviceUser.Backupnum.Eq(response.Backupnum),
			query.CheckinDeviceUser.Status.Eq(-1),
		).Update(query.CheckinDeviceUser.Status, 1)
	}

}

// HandleSetUserInfoAll 向所有设备下发
func handleDeleteuserAll(msg checkinMsg.DeleteuserMessage) {

	machines, err := query.CheckinDevice.WithContext(context.Background()).Find()
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
			query.CheckinDeviceUser.WithContext(context.Background()).
				Where(query.CheckinDeviceUser.Enrollid.Eq(response.Enrollid)).Delete()
		} else {
			query.CheckinDeviceUser.WithContext(context.Background()).
				Where(
					query.CheckinDeviceUser.Enrollid.Eq(response.Enrollid),
					query.CheckinDeviceUser.Backupnum.Eq(response.Backupnum),
				).Delete()
		}
	}
}

// handleEnableuserAll 对所有设备更新用户状态
func handleEnableuserAll(msg checkinMsg.EnableuserMessage) {
	devices, err := query.CheckinDevice.WithContext(context.Background()).Find()
	if err != nil {
		log.Errorf("Error query machines: %v", err)
		return
	}
	for _, device := range devices {
		conn, exists := clientsBySn[device.Sn]
		if !exists {
			log.Warnf("设置用户状态失败，设备[%s]未连接", device.Sn)
			continue
		}
		// clientsBySn
		handleEnableuser(conn, msg)
	}
}

func handleEnableuser(conn *websocket.Conn, msg checkinMsg.EnableuserMessage) {
	sendData(conn, msg)
}

func receiveEnableuser(conn *websocket.Conn, msg []byte) {
	var response checkinMsg.EnableuserResponse
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return
	}
	sn := clientsByConn[conn]
	if !response.Result {
		log.Errorf("设备[%s]设置用户状态失败, 原因:%d", sn, response.Reason)
	} else {
		log.Printf("设备[%s]设置用户状态成功", sn)
	}

}

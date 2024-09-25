package server

import (
	"checkin/query"
	"checkin/query/model"
	checkinMsg "checkin/server/msg"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 处理 "reg" 命令
func receiveReg(conn *websocket.Conn, msg []byte) {
	// log.Printf("Received registration from device: %s", msg.Sn)
	var regMsg checkinMsg.RegMessage
	if err := json.Unmarshal(msg, &regMsg); err != nil {
		log.Println("RegMessage unmarshal error:", err)
		// 返回成功响应
		sendData(conn, checkinMsg.WSResponse{
			Ret:    "reg",
			Result: false,
			Reason: 1,
		})
		return
	}
	// 记录连接信息
	clientsBySn[regMsg.Sn] = conn
	clientsByConn[conn] = regMsg.Sn

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
		sendData(conn, checkinMsg.WSResponse{
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
	sendData(conn, checkinMsg.WSResponse{
		Ret:       "reg",
		Result:    true,
		Cloudtime: time.Now().Format(time.DateTime),
	})
}

func receiveSendlog(conn *websocket.Conn, msg []byte) {
	var regMsg checkinMsg.SendlogMessage
	if err := json.Unmarshal(msg, &regMsg); err != nil {
		log.Println("RegMessage unmarshal error:", err)
		// 返回成功响应
		sendData(conn, checkinMsg.WSResponse{
			Ret:    "sendlog",
			Result: false,
			Reason: 1,
		})
		return
	}
	sendData(conn, checkinMsg.WSResponse{
		Ret:       "sendlog",
		Result:    true,
		Count:     1,
		Logindex:  0,
		Cloudtime: time.Now().Format(time.DateTime),
		Access:    1,
	})
}

func receiveSenduser(conn *websocket.Conn, msg []byte) {
	var senduserMsg checkinMsg.SenduserMessage
	if err := json.Unmarshal(msg, &senduserMsg); err != nil {
		log.Println("RegMessage unmarshal error:", err)
		// 返回成功响应
		sendData(conn, checkinMsg.WSResponse{
			Ret:    "senduser",
			Result: false,
			Reason: 1,
		})
		return
	}
	recordStr := ""
	switch v := senduserMsg.Record.(type) {
	case float64:
		recordStr = fmt.Sprintf("%.f", v) // 将数字转换为字符串
	case string:
		recordStr = v // 直接使用字符串
	default:
		recordStr = "" // 或者处理其他情况
	}
	fmt.Println("接收Message:", senduserMsg)
	ctx := context.Background()
	userInfo, err := query.UserCheckinMachineInfo.WithContext(ctx).
		Where(query.UserCheckinMachineInfo.Enrollid.Eq(senduserMsg.Enrollid),
			query.UserCheckinMachineInfo.Backupnum.Eq(senduserMsg.Backupnum),
			query.UserCheckinMachineInfo.Sn.Eq(senduserMsg.Sn),
		).First()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			sendData(conn, checkinMsg.WSResponse{
				Ret:    "senduser",
				Result: false,
				Reason: 1,
			})
			return
		}
	}
	// 用户信息未记录
	if userInfo == nil {
		log.Debugf("用户信息未登记: %+v", senduserMsg)
		err = query.UserCheckinMachineInfo.WithContext(ctx).Create(&model.UserCheckinMachineInfo{
			Sn:        senduserMsg.Sn,
			Enrollid:  senduserMsg.Enrollid,
			Name:      senduserMsg.Name,
			Backupnum: senduserMsg.Backupnum,
			Record:    recordStr,
		})
		if err != nil {
			log.Debugf("Error create user: %v", err)
			sendData(conn, checkinMsg.WSResponse{
				Ret:    "senduser",
				Result: false,
				Reason: 1,
			})
			return
		}
	}

	sendData(conn, checkinMsg.WSResponse{
		Ret:       "senduser",
		Result:    true,
		Cloudtime: time.Now().Format(time.DateTime),
	})
}

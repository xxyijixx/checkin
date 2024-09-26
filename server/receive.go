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

	checkinMachine, err := query.CheckinDevice.WithContext(context.Background()).
		Where(query.CheckinDevice.Sn.Eq(regMsg.Sn)).
		First()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
		jsonData, err := json.Marshal(regMsg.Devinfo)
		if err != nil {
			return
		}
		query.CheckinDevice.WithContext(context.Background()).
			Create(&model.CheckinDevice{
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
	query.CheckinDevice.WithContext(context.Background()).Where(query.CheckinDevice.Sn.Eq(checkinMachine.Sn)).Update(query.CheckinDevice.Devinfo, jsonData)
	sendData(conn, checkinMsg.WSResponse{
		Ret:       "reg",
		Result:    true,
		Cloudtime: time.Now().Format(time.DateTime),
	})
}

func receiveSendlog(conn *websocket.Conn, msg []byte) {
	var sendlogMsg checkinMsg.SendlogMessage
	if err := json.Unmarshal(msg, &sendlogMsg); err != nil {
		log.Println("RegMessage unmarshal error:", err)
		// 返回成功响应
		sendData(conn, checkinMsg.WSResponse{
			Ret:    "sendlog",
			Result: false,
			Reason: 1,
		})
		return
	}
	ctx := context.Background()

	logRecord := make([]*model.CheckinDeviceRecord, sendlogMsg.Count)
	for i, record := range sendlogMsg.Record {
		reportTime, err := time.Parse(time.DateTime, record.Time)
		if err != nil {
			reportTime = time.Now()
		}
		logRecord[i] = &model.CheckinDeviceRecord{
			Sn:         sendlogMsg.Sn,
			Mode:       record.Mode,
			Event:      record.Event,
			Inout:      record.Inout,
			ReportTime: reportTime,
			Enrollid:   record.Enrollid,
		}
	}
	err := query.CheckinDeviceRecord.WithContext(ctx).Create(logRecord...)
	if err != nil {
		log.Errorf("记录考勤数据失败: %v", err)
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
	userInfo, err := query.CheckinDeviceUser.WithContext(ctx).
		Where(query.CheckinDeviceUser.Enrollid.Eq(senduserMsg.Enrollid),
			query.CheckinDeviceUser.Backupnum.Eq(senduserMsg.Backupnum),
			query.CheckinDeviceUser.Sn.Eq(senduserMsg.Sn),
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
		err = query.CheckinDeviceUser.WithContext(ctx).Create(&model.CheckinDeviceUser{
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

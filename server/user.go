// 处理与用户相关的逻辑
// 1. 获取用户列表
// 2. 获取用户信息
// 3. 下发用户信息
// 4. 禁用\启用用户

package server

import (
	"checkin/query"
	"checkin/server/common"
	checkinMsg "checkin/server/msg"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// handleGeuserlistRandomDevice 从已连接的设备中随机请求设备获取用户数据
func handleGetuserlistRandomDevice() {

	deviceSn := make([]string, len(clientsBySn))
	for sn := range clientsBySn {
		deviceSn = append(deviceSn, sn)
	}
	randomKey := rand.Intn(len(deviceSn))
	conn := clientsBySn[deviceSn[randomKey]]

	handleGetuserlist(conn, true)
}

// handleGetuserlist 处理获取用户数据
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

// HandleSetUserInfoAll 向所有设备下发用户信息,
func handleSetUserInfoAll(msg checkinMsg.SetuserinfoMessage) *common.RetMessage[checkinMsg.RetSetuserinfo] {

	devices, err := query.CheckinDevice.WithContext(context.Background()).Find()
	if err != nil {
		log.Errorf("Error query machines: %v", err)
		return common.Error[checkinMsg.RetSetuserinfo]("获取设备信息失败")
	}
	routingKey := fmt.Sprintf("setuserinfo-%d-%d", msg.Enrollid, msg.Backupnum)
	GlobalCache.Set(routingKey, len(devices), CacheDefaultExpiration)
	for _, device := range devices {
		conn, exists := clientsBySn[device.Sn]
		if !exists {
			log.Warnf("下发用户失败，设备[%s]未连接", device.Sn)
			continue
		}
		// clientsBySn
		go handleSetuserinfo(conn, msg)
	}
	// 等待处理
	response, err := waitForResponses[checkinMsg.RetDeviceSetuserinfo](routingKey, len(devices), CacheDefaultExpiration)
	if err != nil {
		return common.Error[checkinMsg.RetSetuserinfo]("上传失败")
	}
	data := checkinMsg.RetSetuserinfo{
		Result: true,
	}
	for _, res := range response {
		if res.Ret != 1 {
			data.Reason = res.Reason
			data.Result = false
			return common.ErrorWithData(res.Msg, data)
		}
	}

	return common.SuccessWithData(data)
}

// handleSetuserinfo设置 仅对设备下发，不记录
func handleSetuserinfo(conn *websocket.Conn, msg checkinMsg.SetuserinfoMessage) {
	_, exists := clientsByConn[conn]
	if !exists {
		log.Warn("连接已断开")
		return
	}
	sendData(conn, msg)
}

func receiveSetuserinfo(conn *websocket.Conn, msg []byte) {

	var response checkinMsg.WSResponse
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return
	}

	ret := checkinMsg.RetDeviceSetuserinfo{}

	// 处理设备响应
	if !response.Result {
		if sn := clientsByConn[conn]; sn != "" {
			log.Warnf("对设备[%s]下发用户信息失败: %v", sn, response.Msg)
			ret.Msg = response.Msg
			ret.Reason = response.Reason
			ret.Ret = 0
		} else {
			log.Println("Error set user info:", response.Msg)
		}
	} else {
		log.Printf("对设备[%s]下发用户信息[%d]成功", response.Sn, response.Enrollid)
		ret.Msg = "success"
		ret.Ret = 1
	}

	jsonData, err := json.Marshal(ret)
	if err != nil {
		log.Debugf("Marshal json error %v", err)
	}

	// 添加消息
	MessagesChan <- RetMessage{
		RoutingKey: fmt.Sprintf("setuserinfo-%d-%d", response.Enrollid, response.Backupnum),
		Data:       string(jsonData),
		RetryCount: 0,
	}
}

// HandleSetUserInfoAll 向所有设备下发
func handleDeleteuserAll(msg checkinMsg.DeleteuserMessage) *common.RetMessage[checkinMsg.RetSetuserinfo] {

	devices, err := query.CheckinDevice.WithContext(context.Background()).Find()
	if err != nil {
		log.Errorf("Error query devices: %v", err)
		return common.Error[checkinMsg.RetSetuserinfo]("获取设备信息失败")
	}
	routingKey := fmt.Sprintf("deleteuser-%d-%d", msg.Enrollid, msg.Backupnum)
	for _, device := range devices {
		conn, exists := clientsBySn[device.Sn]
		if !exists {
			log.Warnf("删除用户失败，设备[%s]未连接", device.Sn)
			continue
		}
		// clientsBySn
		handleDeleteuser(conn, msg)
	}
	// 等待处理
	response, err := waitForResponses[checkinMsg.RetDeviceSetuserinfo](routingKey, len(devices), CacheDefaultExpiration)
	if err != nil {
		return common.Error[checkinMsg.RetSetuserinfo]("处理失败")
	}
	data := checkinMsg.RetSetuserinfo{
		Result: true,
	}
	for _, res := range response {
		if res.Ret != 1 {
			data.Reason = res.Reason
			data.Result = false
			return common.ErrorWithData(res.Msg, data)
		}
	}

	return common.SuccessWithData(data)
}

// handleDeleteuser 处理删除用户信息
func handleDeleteuser(conn *websocket.Conn, msg checkinMsg.DeleteuserMessage) {
	sendData(conn, msg)
}

// receiveDeleteser 接收设备删除用户命令的响应
func receiveDeleteuser(conn *websocket.Conn, msg []byte) {
	var response checkinMsg.DeleteuserResponse
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return
	}
	ret := checkinMsg.RetDeviceSetuserinfo{}
	sn := clientsByConn[conn]
	if !response.Result {
		log.Errorf("设备[%s]删除用户信息失败, 原因:%d", sn, response.Reason)
		ret.Msg = "error"
		ret.Reason = response.Reason
		ret.Ret = 0
	} else {
		log.Printf("设备[%s]删除用户信息[%d]成功", sn, response.Enrollid)
		ret.Msg = "success"
		ret.Ret = 1
	}
	jsonData, err := json.Marshal(ret)
	if err != nil {
		log.Debugf("Marshal json error %v", err)
	}

	MessagesChan <- RetMessage{
		RoutingKey: fmt.Sprintf("deleteuser-%d-%d", response.Enrollid, response.Backupnum),
		Data:       string(jsonData),
		RetryCount: 0,
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

// handleEnableuser 处理用户禁用\启用命令
func handleEnableuser(conn *websocket.Conn, msg checkinMsg.EnableuserMessage) {
	sendData(conn, msg)
}

// receiveEnableuser 接收设备处理用户禁用\启用命令的响应
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

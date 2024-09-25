package server

import (
	"checkin/server/msg"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type SetUserInfo struct {
	Enrollid  int         `json:"enrollid,omitempty"`
	Name      string      `json:"name,omitempty"`
	Backupnum int         `json:"backupnum,omitempty"`
	Admin     int         `json:"admin,omitempty"`
	Record    interface{} `json:"record"` // 为密码的时候为number，其他时候为字符串
}

type DeleteUserParams struct {
	Enrollid  int `json:"enrollid"`
	Backupnum int `json:"backupnum"`
}

func UserHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listUserHandle(w, r)
	case http.MethodPost:
		setUserHandle(w, r)
	case http.MethodDelete:
		deleteUserHandle(w, r)
	default:
		// 如果请求方法不支持，返回405 Method Not Allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listUserHandle(w http.ResponseWriter, r *http.Request) {

}

func setUserHandle(w http.ResponseWriter, r *http.Request) {
	var userInfo SetUserInfo
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	logrus.Infof("处理设置用户: %+v", userInfo)
	if userInfo.Name == "" || userInfo.Enrollid == 0 || userInfo.Backupnum == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	handleSetUserInfoAll(msg.SetuserinfoMessage{
		Cmd:       CmdSetuserinfo,
		Name:      userInfo.Name,
		Enrollid:  userInfo.Enrollid,
		Backupnum: userInfo.Backupnum,
		Admin:     userInfo.Admin,
		Record:    userInfo.Record,
	})
}

func deleteUserHandle(w http.ResponseWriter, r *http.Request) {
	var params DeleteUserParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if params.Enrollid == 0 || (params.Backupnum != 50 && (params.Backupnum < 0 || params.Backupnum > 13)) {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// handleDeleteuser()
	handleDeleteuserAll(msg.DeleteuserMessage{
		Cmd:       CmdDeleteuser,
		Enrollid:  params.Enrollid,
		Backupnum: params.Backupnum,
	})

}

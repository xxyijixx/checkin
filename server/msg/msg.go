package msg

type BaseMessage struct {
	Cmd string `json:"cmd"`
}

// RegMessage 登录注册消息
type RegMessage struct {
	Cmd     string `json:"cmd"`
	Sn      string `json:"sn"`
	Devinfo struct {
		Modelname  string `json:"modelname"`
		Usersize   int    `json:"usersize"`
		Fpsize     int    `json:"fpsize"`
		Cardsize   int    `json:"cardsize"`
		Pwdsize    int    `json:"pwdsize"`
		Logsize    int    `json:"logsize"`
		Useduser   int    `json:"useduser"`
		Usedfp     int    `json:"usedfp"`
		Usedcard   int    `json:"usedcard"`
		Usedpwd    int    `json:"usedpwd"`
		Usedlog    int    `json:"usedlog"`
		Usednewlog int    `json:"usednewlog"`
		Fpalgo     string `json:"fpalgo"`
		Firmware   string `json:"firmware"`
		Time       string `json:"time"`
	} `json:"devinfo"`
}

// SendlogMessage 上传考勤记录消息
type SendlogMessage struct {
	Cmd      string `json:"cmd"`
	Count    int    `json:"count"`
	Sn       string `json:"sn"`
	Logindex int    `json:"logindex"`
	Record   []struct {
		Enrollid   int    `json:"enrollid"`
		Name       string `json:"name"`
		Time       string `json:"time,omitempty"`
		Mode       int    `json:"mode"`
		Verifymode int    `json:"verifymode"`
		Image      string `json:"image"`
		Time0      string `json:"time ,omitempty"`
		Inout      int    `json:"inout,omitempty"`
	} `json:"record"`
}

// SenduserMessage 发送用户信息消息
type SenduserMessage struct {
	Cmd       string `json:"cmd"`
	Sn        string `json:"sn"`
	Enrollid  int    `json:"enrollid"`
	Name      string `json:"name"`
	Backupnum int    `json:"backupnum"`
	Admin     int    `json:"admin"`
	Record    string `json:"record"`
}

// 响应
type WSResponse struct {
	Ret       string `json:"ret"`
	Result    bool   `json:"result"`
	Count     int    `json:"count,omitempty"`
	Logindex  int    `json:"logindex,omitempty"`
	Cloudtime string `json:"cloudtime,omitempty"`
	Reason    int    `json:"reason,omitempty"`
	Access    int    `json:"access,omitempty"` // 拓展功能，用于指示这个用户是否可以进门。1，可以进门， 0不可以进门
}

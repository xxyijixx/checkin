package msg

type RetSetuserinfo struct {
	Result bool `json:"result"`
	Reason int  `json:"reason"`
}

type RetDeviceSetuserinfo struct {
	Ret    int    `json:"ret"`
	Sn     string `json:"sn"`
	Reason int    `json:"reason"`
	Msg    string `json:"msg"`
}

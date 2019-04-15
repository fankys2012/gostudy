package cmodel

const (
	UserOnline = iota //在线
	UserOffline		  //离线
)

type User struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userNane"`
	UserState int   `json:"userState"`//用户在线状态
}

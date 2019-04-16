package message

import "github.com/fankys2012/gostudy/chatroom/common/cmodel"

const (
	LoginMesType                 = "LoginMes"
	LoginResMesType              = "LoginResMes"
	RegisterMesType              = "RegisterMesType"
	UserExitsMesType             = "UserExitsMesType"
	StandardResponseMesType      = "StandardResponseMesType"
	NotifyUserOnlineStateMesType = "NotifyUserOnlineStateMesType"
	CharMessageMesType           = "CharMessageType"
)

//聊天消息类型
const (
	PrivateMesType = "privateMesType" //私聊
	GroupMesType   = "groupMesType"   //群聊
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息内容的类型
}

type LoginMes struct {
	UserId   int    `json:"userid"`
	UserPwd  string `json:"userpwd"`
	UserName string `json:"username"`
}

type LoginResMes struct {
	Code     int            `json:"code"`     //返回状态码  500 未注册 200 登陆成功
	UserList map[int]string `json:"userList"` //在线用户列表
	Error    string         `json:"error"`    //返回错误信息
	User     *cmodel.User   `json:"user"`
}

//服务器端标准返回消息体
type StandardResponseMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //返回错误信息
}

//注册消息体
type RegisterMes struct {
	User cmodel.User
}

//用户状态通知消息
type NotifyUserOnlineStateMes struct {
	UserId    int    `json:"userId"`
	UserState int    `json:"userState"`
	UserName  string `json:"userName"`
}

//聊天消息
type ChatMessageMes struct {
	Data    string       `json:"data"`
	MesType string       `json:"mesType"` //消息类型  group:群聊，private:私聊
	Id      int          `json:"id"`      //群组ID/用户ID
	User    *cmodel.User `json:"user"`    //发送消息用户的用户基本信息
}

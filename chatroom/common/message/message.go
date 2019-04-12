package message

const (
	LoginMesType     = "LoginMes"
	LoginResMesType  = "LoginResMes"
	RegisterMesType  = "RegisterMesType"
	UserExitsMesType = "UserExitsMesType"
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
	Code  int    `json:"code"`  //返回状态码  500 未注册 200 登陆成功
	Error string `json:"error"` //返回错误信息
}

//服务器端标准返回消息体
type StandardResponseMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //返回错误信息
}

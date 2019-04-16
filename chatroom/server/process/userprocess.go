package process

import (
	"encoding/json"
	"fmt"
	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"net"

	"github.com/fankys2012/gostudy/chatroom/server/model"

	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
)

type UserProcess struct {
	Conn     net.Conn
	UserId   int
	UserName string
}

//工厂方法 -- 实例化对象
func NewUserPorcess(conn net.Conn) (userprocess *UserProcess) {
	return &UserProcess{
		Conn: conn,
	}
}

//登陆逻辑
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1 从mes 取出 mes.data 并反序列化
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("反序列化失败 err=", err)
		return
	}

	//响应消息体
	var loginResMes message.LoginResMes

	//从redis/db 获取用户信息
	redisConn := model.MyUserDao.RedisPool.Get()
	user, err := model.MyUserDao.GetUserById(redisConn, loginMes.UserId)
	if err != nil {
		fmt.Println("获取用户信息失败,err:", err)
		loginResMes.Code = 500
		loginResMes.Error = "用户不存在，请重新登录"
	} else {
		if loginMes.UserId == user.UserId && loginMes.UserPwd == user.UserPwd {
			loginResMes.Code = 200
			loginResMes.User = user
			//登录成功处理逻辑...
			//将登录用户添加到在线用户列表中
			this.UserId = loginMes.UserId
			this.UserName = user.UserName
			userMg.AddOnlineUser(this)

			//将已上线的用户返回给当前登录用户
			loginResMes.UserList = make(map[int]string)
			for id, user := range userMg.onlineUsers {
				loginResMes.UserList[id] = user.UserName
			}

			//通知其他用户我上线了
			this.NotifyOnlineState(user.UserId, cmodel.UserOnline, user.UserName)

		} else {
			//登陆失败
			loginResMes.Code = 500
			loginResMes.Error = "密码错误，请重新登录"
		}
	}

	//将 响应消息体序列化
	data, err := json.Marshal(loginResMes)
	fmt.Println("响应数据==", loginResMes)
	if err != nil {
		fmt.Println("序列化失败,err = ", err)
		return
	}

	//返回消息体
	resMes := message.Message{
		Type: message.LoginResMesType,
		Data: string(data), //切片转字符串
	}

	fmt.Println("响应数据==", resMes.Data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化失败,err = ", err)
		return
	}
	// 发送响应消息
	transfer := &utils.Transfer{
		Conn: this.Conn,
	}
	err = transfer.WritePkg(data)
	return
}

//校验用户是否存在
func (this *UserProcess) ServerCheckUserExitsById(mes *message.Message) (err error) {
	// 1 从mes 取出 mes.data 并反序列化
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		return
	}
	ok, err := model.MyUserDao.ExistsById(loginMes.UserId)
	if err != nil {
		return
	}

	var resposeMes message.StandardResponseMes
	if ok {
		resposeMes.Code = 200
		resposeMes.Error = model.ERROR_USER_EXISTS.Error()
	} else {
		resposeMes.Code = 40
	}

	data, err := json.Marshal(resposeMes)
	if err != nil {
		fmt.Println("package response message failed ", err)
		return
	}

	smes := message.Message{
		Type: message.UserExitsMesType,
		Data: string(data),
	}
	data, err = json.Marshal(smes)
	if err != nil {
		fmt.Println("package response message failed ", err)
		return
	}

	// 发送响应消息
	transfer := &utils.Transfer{
		Conn: this.Conn,
	}
	err = transfer.WritePkg(data)
	if err != nil {
		fmt.Println("ServerCheckUserExitsById response faield err :", err)
	}
	return
}

func (this *UserProcess) ServerRegister(mes *message.Message) (err error) {
	var registerUser message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerUser)
	if err != nil {
		return
	}
	//直接调用model register 方法
	err = model.MyUserDao.Register(&registerUser.User)
	if err != nil {
		return
	}
	resposeMes := message.StandardResponseMes{
		Code: 200,
	}

	resData, err := json.Marshal(resposeMes)
	if err != nil {
		return
	}

	resMes := message.Message{
		Type: message.StandardResponseMesType,
		Data: string(resData),
	}
	resMesJson, err := json.Marshal(resMes)
	if err != nil {
		return
	}

	// 发送响应消息
	transfer := &utils.Transfer{
		Conn: this.Conn,
	}
	err = transfer.WritePkg(resMesJson)

	return
}

//通知其他在线用户我的状态
func (this *UserProcess) NotifyOnlineState(userId, state int, userName string) {
	for id, uprocess := range userMg.onlineUsers {
		//过滤自己
		if id == this.UserId {
			continue
		}
		uprocess.notifyState(userId, state, userName)
	}
}

func (this *UserProcess) notifyState(userId, state int, userName string) {
	//通知消息体
	userState := message.NotifyUserOnlineStateMes{
		UserId:    userId,
		UserState: state,
		UserName:  userName,
	}

	data, err := json.Marshal(userState)
	if err != nil {
		return
	}

	mes := message.Message{
		Type: message.NotifyUserOnlineStateMesType,
		Data: string(data),
	}
	data, err = json.Marshal(mes)
	if err != nil {
		return
	}

	transfer := &utils.Transfer{
		Conn: this.Conn,
	}
	err = transfer.WritePkg(data)
	if err != nil {
		return
	}

}

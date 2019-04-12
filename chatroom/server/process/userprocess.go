package process

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/fankys2012/gostudy/chatroom/server/model"

	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
)

type UserProcess struct {
	Conn net.Conn
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

	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//响应消息体
	var loginResMes message.LoginResMes

	//伪登陆
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		loginResMes.Code = 200
	} else {
		//登陆失败
		loginResMes.Code = 500
		loginResMes.Error = "login failed"
	}

	//将 响应消息体序列化
	data, err := json.Marshal(loginResMes)
	fmt.Println("响应数据==", loginResMes)
	if err != nil {
		fmt.Println("序列化失败,err = ", err)
		return
	}

	//4 将data 赋值给mes.Data
	resMes.Data = string(data) //切片转字符串
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
	// 发送响应消息
	transfer := &utils.Transfer{
		Conn: this.Conn,
	}
	err = transfer.WritePkg(data)
	return
}

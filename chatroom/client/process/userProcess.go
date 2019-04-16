package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fankys2012/gostudy/chatroom/client/model"
	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
	"net"
)

type UserProcess struct {
	conn net.Conn
}

var (
	MyTransfer *utils.Transfer
)

func NewUserPorcess() (userProcess *UserProcess, err error) {
	//链接服务器
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("connect server failed ", err)
		return
	}
	userProcess = &UserProcess{
		conn: conn,
	}
	MyTransfer = &utils.Transfer{
		Conn: conn,
	}
	return
}

//用户登录客服端部分
func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	//链接服务器
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()

	//2 通过conn给服务器发送消息
	var mes message.Message
	mes.Type = message.LoginMesType
	//3 创建LoginMes 结构体
	// var loginMes message.LoginMes
	// loginMes.userId = userId
	// loginMes.userPwd = userPwd
	loginMes := message.LoginMes{
		UserId:  userId,
		UserPwd: userPwd,
	}
	//4 将loginMes 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal failed ", err)
		return
	}

	mes.Data = string(data)

	//将mes 序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("message json failed ", err)
		return
	}

	//处理服务器端响应消息
	transfer := &utils.Transfer{
		Conn: conn,
	}

	err = transfer.WritePkg(data)
	if err != nil {
		fmt.Println("发送登录信息失败 err:", err)
	}

	mes, err = transfer.ReadPkg()
	if err != nil {
		fmt.Println(err)
		return
	}

	//解析效应消息
	var reponseMes message.LoginResMes
	//json.Unmarshal 第二个参数必须是地址 坑啊坑啊 。。。。
	err = json.Unmarshal([]byte(mes.Data), &reponseMes)
	//fmt.Println(reponseMes)
	if reponseMes.Code == 200 {

		fmt.Println("当前在线用户列表:")
		for id, v := range reponseMes.UserList {
			if id == userId {
				continue
			}
			fmt.Printf("[%d:%s]\r\n", id, v)

			//本地维护的用户列表
			user := &cmodel.User{
				UserId:    id,
				UserName:  v,
				UserState: cmodel.UserOnline,
			}
			onlineUserList[id] = user
		}

		//保存当前用户信息
		CurUserInfo = &model.UserInfo{
			Conn: this.conn,
		}
		CurUserInfo.UserId = reponseMes.User.UserId
		CurUserInfo.UserName = reponseMes.User.UserName
		CurUserInfo.UserState = cmodel.UserOnline

		/**
		 * 启动一个协程，该协程保持与服务器端的通讯，如果服务器推送消息给客服端
		 * 则接收并显示在终端
		 */
		go serverProcessMes(conn)
		//显示登录成功后的界面
		showMenu(conn)

	} else if reponseMes.Code == 500 {
		fmt.Println(reponseMes.Error)
		err = errors.New(reponseMes.Error)
		return err
	}
	return nil
}

func (this *UserProcess) userExitsCheck(id int) (err error) {

	//fmt.Println("userExitsCheck start...")
	loginMes := message.LoginMes{
		UserId: id,
	}
	//序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("package heck user data failed", err)
		return
	}
	sendData := message.Message{
		Type: message.UserExitsMesType,
		Data: string(data),
	}
	data, err = json.Marshal(sendData)
	if err != nil {
		fmt.Println("Marshal err :", err)
		return
	}

	//实例化utils 包
	//transfer := &utils.Transfer{
	//	Conn: this.conn,
	//}
	transfer := MyTransfer
	//向服务器发送验证用户是否存在的消息
	err = transfer.WritePkg(data)
	//处理服务器端响应消息
	mes, err := transfer.ReadPkg()
	if err != nil {
		fmt.Println("userExitsCheck read response failed. err:", err)
		return
	}

	//解析效应消息
	var resposeMes message.StandardResponseMes
	//json.Unmarshal 第二个参数必须是地址 坑啊坑啊 。。。。
	err = json.Unmarshal([]byte(mes.Data), &resposeMes)
	//fmt.Println("userExitsCheck server response ",resposeMes,err)
	if err != nil {
		return
	}
	if resposeMes.Code == 200 {
		err = errors.New(resposeMes.Error)
	}
	//fmt.Println("userExitsCheck end...")
	return
}

func (this *UserProcess) PreRegister() (user *cmodel.User, err error) {
	var userId int
	var userPwd, userName string
	fmt.Println("请输入用户ID")
	for {
		fmt.Scanf("%d\n", &userId)
		if userId == 0 {
			fmt.Println("用户ID无效请重新输入")
		} else {
			err := this.userExitsCheck(userId)
			if err == nil {
				break
			}
			fmt.Println("用户已存在请重新输入,err == ", err)
		}
	}
	fmt.Println("请输入密码")
	for {
		fmt.Scanf("%s\n", &userPwd)
		if userPwd == "" {
			fmt.Println("密码不能为空，请重新输入")
		} else {
			break
		}
	}
	fmt.Println("请输入昵称")
	for {
		fmt.Scanf("%s\n", &userName)
		if userName == "" {
			fmt.Println("用户名不能为空，请重新输入")
		} else {
			break
		}
	}

	myUser := cmodel.User{
		UserId:   userId,
		UserPwd:  userPwd,
		UserName: userName,
	}
	user = &myUser

	fmt.Println("输入结束")
	return
}

func (this *UserProcess) Register() (err error) {
	user, err := this.PreRegister()
	if err != nil {
		return
	}
	regUser := message.RegisterMes{
		User: *user,
	}
	fmt.Println("register user :", regUser)
	userJsonData, err := json.Marshal(regUser)
	if err != nil {
		fmt.Println("user json failed ", err)
	}

	mes := message.Message{
		Type: message.RegisterMesType,
		Data: string(userJsonData),
	}
	mesData, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("package register data failed ", err)
		return
	}
	//实例化utils 包
	//transfer := &utils.Transfer{
	//	Conn: this.conn,
	//}
	transfer := MyTransfer
	err = transfer.WritePkg(mesData)
	if err != nil {
		fmt.Println("send register message failed ", err)
		return
	}
	//处理服务器端响应消息
	mes, err = transfer.ReadPkg()
	if err != nil {
		return
	}
	var response message.StandardResponseMes
	err = json.Unmarshal([]byte(mes.Data), &response)
	if err != nil {
		return
	}
	if response.Code == 200 {
		fmt.Println("注册成功，请登录")
		return
	} else {
		fmt.Println("注册失败")
	}
	return
}

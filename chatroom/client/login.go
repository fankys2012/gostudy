package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
)

func login(userId int, userPwd string) (err error) {
	// fmt.Printf("userId = %d userPwd = %s ", userId, userPwd)
	// return nil

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

	//发送消息长度
	// 先获取 data 的长度-> 转换成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) //len -> byte
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("字符长度发送失败 ", err)
		return
	}
	//发送内容
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("消息内容发送失败 ", err)
		return
	}

	//处理服务器端响应消息
	mes, err = utils.ReadPkg(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	//解析效应消息
	var reponseMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), reponseMes)
	if reponseMes.Code == 200 {
		fmt.Printf("登陆成功")
	} else if reponseMes.Code == 500 {
		fmt.Println(reponseMes.Error)
	}
	return nil
}

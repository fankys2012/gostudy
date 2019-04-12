package process

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
)

type UserProcess struct {
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
	transfer := &utils.Transfer{
		Conn: conn,
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
	fmt.Println(reponseMes)
	if reponseMes.Code == 200 {

		/**
		 * 启动一个协程，该协程保持与服务器端的通讯，如果服务器推送消息给客服端
		 * 则接收并显示在终端
		 */
		go serverProcessMes(conn)
		//显示登录成功后的界面
		showMenu()

	} else if reponseMes.Code == 500 {
		fmt.Println(reponseMes.Error)
		err = errors.New(reponseMes.Error)
		return err
	}
	return nil
}

func (this *UserProcess) PreRegister() (err error) {
	var userId int
	fmt.Println("请输入用户ID")
	for {
		fmt.Scanf("%d\n", &userId)
		if userId == 0 {
			fmt.Println("用户ID无效请重新输入")
		} else {

		}

		if userId == 10 {
			break
		}
		fmt.Println("用户已存在请重新输入")
	}
	fmt.Println("输入结束")
	return nil
}

func (this *UserProcess) Register(user *cmodel.User) (err error) {
	return nil
}

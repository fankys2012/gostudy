package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/fankys2012/gostudy/chatroom/common/message"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	//conn.Read 只有在conn 没有被关闭的情况才会一直阻塞
	//如果客户端关闭了，则
	_, err = conn.Read(buf[:4])
	if err != nil {
		fmt.Printf("数据读取失败err =%s", err)
		// err = errors.New("read pkg header error")
		return
	}

	//根据buf[:4] 转成uint32
	pkgLen := binary.BigEndian.Uint32(buf[0:4])

	//根据pkgLen 读取消息内容  从conn 中读取内容存入buf中
	n, err := conn.Read(buf[:pkgLen])
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("读取消息失败,err ", err)
		// err = errors.New("read message body failed")
		return
	}

	//反序列化message  mes必须传地址; mes 在返回参数中已声明，不用重复声明
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("消息解析失败 ", err)
	}
	return
}

//登陆逻辑
func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
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
	if err != nil {
		fmt.Println("序列化失败,err = ", err)
		return
	}

	//4 将data 赋值给mes.Data
	resMes.Data = string(data) //切片转字符串

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化失败,err = ", err)
		return
	}
	// 发送响应消息
	err = writePkg(conn, data)
	return
}

func writePkg(conn net.Conn, data []byte) (err error) {
	//发送消息内容长度
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

	//发送消息体内容
	//发送内容
	n, err = conn.Write(data)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("消息内容发送失败 ", err)
		return
	}
	return

}

//根据不同的消息 处理不同的逻辑
func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//登陆
		err = serverProcessLogin(conn, mes)
	case message.RegisterMesType:
		//注册
	default:
		fmt.Println("消息类型不存在")

	}
	return
}

func process(conn net.Conn) {
	defer conn.Close()

	//读取客服端消息
	for {
		mes, err := readPkg(conn)
		if err != nil {
			fmt.Println("消息读取失败", err)
			return
		}
		err = serverProcessMes(conn, &mes)
		if err != nil {
			return
		}
		fmt.Println("mes=", mes)
	}
}

func main() {

	fmt.Println("服务器在8888端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("listen err = ", err)
		return
	}

	//监听
	for {
		fmt.Println("等待客服端链接")
		conn, err := listen.Accept()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客服端退出")
			} else {
				fmt.Println("listen accept err =", err)
			}
			return

		}
		//链接成功，则启动一个协程和客服端保持通讯
		go process(conn)
	}
}

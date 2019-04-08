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

func process(conn net.Conn) {
	defer conn.Close()

	//读取客服端消息
	for {
		mes, err := readPkg(conn)
		if err != nil {
			fmt.Println("消息读取失败", err)
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

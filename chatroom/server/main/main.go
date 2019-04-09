package main

import (
	"fmt"
	"io"
	"net"
)

func dispatch(conn net.Conn) {
	defer conn.Close()
	//调用processer.Do
	processor := &Processer{
		Conn: conn,
	}
	err := processor.Do()
	if err != nil {
		return
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
		go dispatch(conn)
	}
}

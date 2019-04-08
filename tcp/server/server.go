package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close()

	for {
		fmt.Printf("服务器等待客服端%s发送信息\n", conn.RemoteAddr().String())
		buf := make([]byte, 1024)
		//等待客服端通过conn 发送信息 如果客服端没有write 则一直阻塞在这里
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("服务器端读取失败 err :", err)
			return
		}
		//显示内容
		content := string(buf[:n])
		content = "服务端收到内容 ：" + content
		fmt.Print(content)

		_, errs := conn.Write([]byte(content))
		if errs != nil {
			fmt.Println("conn.write err=\n", err)
		}
	}
}

func main() {
	fmt.Println("服务器开始监听...")

	listen, err := net.Listen("tcp", "127.0.0.1:8888")

	if err != nil {
		fmt.Println("listen err = ", err)
	}

	defer listen.Close()
	//等待客服端链接
	for {
		fmt.Println("等待客服端链接...")
		// fmt.Println("Addr =%v", listen.Addr())
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept() err ", err)
		} else {
			fmt.Printf("conn=%v,Ip=%s\n", conn, conn.RemoteAddr().String())
		}
		go process(conn)
	}

	fmt.Println("listen succ = %v", listen)
}

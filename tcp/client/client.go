package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("链接失败", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("reader error : %v", err)
		}
		//如果用户输入exit 退出
		if strings.Trim(line, " \r\n") == "exit" {
			fmt.Print("客户端退出...\n")
			break
		}

		n, err := conn.Write([]byte(line))
		if err != nil {
			fmt.Println("conn.write err=\n,line", err, n)
		}

		buf := make([]byte, 1024)
		//等待客服端通过conn 发送信息 如果客服端没有write 则一直阻塞在这里
		n, errs := conn.Read(buf)
		if errs != nil {
			fmt.Println("服务器端读取失败 err :", errs)
			return
		}
		fmt.Println("服务端响应：", string(buf[:n]))
	}

}

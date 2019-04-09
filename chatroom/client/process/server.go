package process

import (
	"fmt"
	"net"
	"os"

	"github.com/fankys2012/gostudy/chatroom/common/utils"
)

//登录后界面

func showMenu() {
	for {

		fmt.Println("-----------登录成功-----------")
		fmt.Println("\t\t\t 1 显示用户列表")
		fmt.Println("\t\t\t 2 发送消息")
		fmt.Println("\t\t\t 3 信息列表")
		fmt.Println("\t\t\t 4 退出系统")
		fmt.Println("\t\t\t 请选择（1-4）")

		var key int

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("用户列表")
		case 2:
			fmt.Println("发送消息")
		case 3:
			fmt.Println("查看列表")
		case 4:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入错误请重新输入")

		}
	}
}

func serverProcessMes(conn net.Conn) {
	//创建transfer 实例 ,不停读取消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			return
		}
		fmt.Println("读取到服务器消息 mes=", mes)
	}
}

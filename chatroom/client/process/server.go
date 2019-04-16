package process

import (
	"encoding/json"
	"fmt"
	"github.com/fankys2012/gostudy/chatroom/common/message"
	"net"
	"os"

	"github.com/fankys2012/gostudy/chatroom/common/utils"
)

//登录后界面

func showMenu(conn net.Conn) {
	for {

		fmt.Printf("\t\t-----------欢迎[%d:%s]回来-----------\t", CurUserInfo.UserId, CurUserInfo.UserName)
		fmt.Println("\t\t\t")
		fmt.Println("\t\t\t 1 显示用户列表")
		fmt.Println("\t\t\t 2 发送消息")
		fmt.Println("\t\t\t 3 信息列表")
		fmt.Println("\t\t\t 4 退出系统")
		fmt.Println("\t\t\t 请选择（1-4）")

		var key int

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			showUserList()
		case 2:
			err := sendMessage(conn)
			if err != nil {
				fmt.Println("发送消息失败,err:", err)
			}
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
		switch mes.Type {
		case message.NotifyUserOnlineStateMesType: //用户上下线通知
			var notifyUserState message.NotifyUserOnlineStateMes
			json.Unmarshal([]byte(mes.Data), &notifyUserState)
			updateOnlineUserList(&notifyUserState)
		case message.CharMessageMesType: //消息通知
			showMessage(&mes)
		default:
			fmt.Printf("未知消息类型%s", mes.Type)
		}
	}
}

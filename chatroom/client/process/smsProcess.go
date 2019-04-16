package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
	"net"
	"os"
)

func sendMessage(conn net.Conn) (err error) {
	fmt.Println("请输入聊天方式：1 私聊 / 2 群聊")
	var chatFlag int
	var charType, mes string

	fmt.Scanf("%d\n", &chatFlag)
	if chatFlag == 1 {
		charType = message.PrivateMesType
	} else {
		charType = message.GroupMesType
	}

	fmt.Println("请输入用户ID:")
	var userId int
	fmt.Scanf("%d\n", &userId)

	flag := false

	for id, _ := range onlineUserList {
		if id == userId {
			flag = true
			break
		}
	}
	if flag == false {
		fmt.Println("用户不存在")
		return
	}
	fmt.Println("请输入聊天内容:")
	//fmt.Scanf 若遇到空则认为是下一次输入
	//fmt.Scanf("%s\n", &mes)
	reader := bufio.NewReader(os.Stdin)
	rdata, _, _ := reader.ReadLine()
	mes = string(rdata)

	sendMes := message.ChatMessageMes{
		Data:    mes,
		MesType: charType,
		Id:      userId,
		User: &cmodel.User{
			UserId:   CurUserInfo.UserId,
			UserName: CurUserInfo.UserName,
		},
	}

	data, err := json.Marshal(sendMes)
	if err != nil {
		return
	}

	sendBody := message.Message{
		Type: message.CharMessageMesType,
		Data: string(data),
	}
	data, err = json.Marshal(sendBody)
	if err != nil {
		return
	}

	transfer := &utils.Transfer{
		Conn: conn,
	}
	err = transfer.WritePkg(data)
	if err != nil {
		return
	}
	return
}

func showMessage(mes *message.Message) {
	var chatMes message.ChatMessageMes
	err := json.Unmarshal([]byte(mes.Data), &chatMes)
	if err != nil {
		return
	}
	fmt.Printf("[%d:%s] 说：%s\r\n", chatMes.User.UserId, chatMes.User.UserName, chatMes.Data)
	//fmt.Println("收到新的消息：", chatMes)
}

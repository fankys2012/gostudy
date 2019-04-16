package process

import (
	"encoding/json"
	"fmt"
	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
	"net"
)

type ChatProcess struct {
	conn net.Conn
}

func NewChatProcess(conn net.Conn) (cp *ChatProcess) {
	cp = &ChatProcess{
		conn: conn,
	}
	return
}

//发送聊天信息
func (this *ChatProcess) SendChatMessage(mes *message.Message) (err error) {
	var chatMes message.ChatMessageMes
	err = json.Unmarshal([]byte(mes.Data), &chatMes)
	if err != nil {
		return
	}
	fmt.Println("服务端收到的消息：", chatMes)

	switch chatMes.MesType {
	case message.PrivateMesType: //私聊
		err = this.sendPrivateMes(chatMes.User, []byte(chatMes.Data))
	case message.GroupMesType: //群聊
		err = this.sendGroupMes(chatMes.User, []byte(chatMes.Data))
	default:
		fmt.Println("错误的消息类型", chatMes.MesType)
	}
	return
}

//发送私聊信息
func (this *ChatProcess) sendPrivateMes(user *cmodel.User, data []byte) (err error) {
	//假设好友列表存储在 userMg.onlineUsers
	for id, uprocess := range userMg.onlineUsers {
		if id == user.UserId {
			err = this.sendMes(data, uprocess.Conn, user)
			break
		}
	}
	return
}

//发送群聊信息
//userId 发送消息用户ID
func (this *ChatProcess) sendGroupMes(user *cmodel.User, data []byte) (err error) {
	for id, uprocess := range userMg.onlineUsers {
		if id == user.UserId {
			continue
		}
		err = this.sendMes(data, uprocess.Conn, user)
		if err != nil {
			return
		}
	}
	return
}

func (this *ChatProcess) sendMes(data []byte, conn net.Conn, user *cmodel.User) (err error) {

	chatMes := message.ChatMessageMes{
		Data: string(data),
		User: user,
	}
	sendData, err := json.Marshal(chatMes)
	if err != nil {
		return
	}

	mes := message.Message{
		Type: message.CharMessageMesType,
		Data: string(sendData),
	}
	sendData, err = json.Marshal(mes)
	if err != nil {
		return
	}
	transfer := &utils.Transfer{
		Conn: conn,
	}
	err = transfer.WritePkg(sendData)
	if err != nil {
		return
	}
	return
}

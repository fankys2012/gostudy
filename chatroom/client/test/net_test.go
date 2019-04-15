package test

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/fankys2012/gostudy/chatroom/common/utils"

	"github.com/fankys2012/gostudy/chatroom/common/message"
)

func TestNet(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		t.Error("connect server failed ", err)
		return
	}

	mes := message.Message{
		Type: "wrongMesType",
		Data: "this is my test",
	}
	mesData, err := json.Marshal(mes)

	transfer2 := &utils.Transfer{
		Conn: conn,
	}

	//发送消息体内容
	//发送内容
	err = transfer2.WritePkg(mesData)
	if err != nil {
		t.Error("消息内容发送失败 ", err)
		return
	}

	res1, err := transfer2.ReadPkg()
	t.Log(res1)

	mes2 := message.Message{
		Type: "wrongMesType2",
		Data: "this is my test",
	}
	mesData, err = json.Marshal(mes2)
	transfer := &utils.Transfer{
		Conn: conn,
	}

	//发送内容
	err = transfer.WritePkg(mesData)
	if err != nil {
		t.Error("消息内容发送失败 ", err)
		return
	}
	res2, err := transfer.ReadPkg()
	if err != nil {
		t.Error("读取消息失败 ", err)
		return
	}
	t.Log(res2)

	t.Log("success ", err)
}

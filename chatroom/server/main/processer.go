package main

import (
	"fmt"
	"io"
	"net"

	"github.com/fankys2012/gostudy/chatroom/common/message"
	"github.com/fankys2012/gostudy/chatroom/common/utils"
	"github.com/fankys2012/gostudy/chatroom/server/process"
)

type Processer struct {
	Conn net.Conn
}

//根据不同的消息 处理不同的逻辑
func (this *Processer) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//登陆
		userprocess := &process.UserProcess{
			Conn: this.Conn,
		}
		err = userprocess.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//注册
	case message.UserExitsMesType:
		userprocess := process.NewUserPorcess(this.Conn)
		err = userprocess.ServerCheckUserExitsById(mes)
	default:
		fmt.Println("消息类型不存在")

	}
	return
}

func (this *Processer) Do() (err error) {
	//读取客服端消息
	for {
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客服端退出连接")
				return err
			} else {
				fmt.Println("读取客服端消息失败", err)
				return err
			}
		}
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
		fmt.Println("mes=", mes)
	}
}

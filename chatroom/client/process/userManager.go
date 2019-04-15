package process

import (
	"fmt"
	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
)

//初始化在线用户列表
var onlineUserList map[int]*cmodel.User = make(map[int]*cmodel.User,20)

func showUserList()  {
	fmt.Println("当前在线用户列表:")
	for id,v := range onlineUserList{
		var onlineState string
		if v.UserState == cmodel.UserOnline {
			onlineState = "在线"
		} else if v.UserState == cmodel.UserOffline {
			onlineState = "离线"
		}
		fmt.Printf("[%d:%s/%s]",id,v.UserName,onlineState)
	}
	fmt.Println("\r\r")
}
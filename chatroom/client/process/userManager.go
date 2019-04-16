package process

import (
	"fmt"
	"github.com/fankys2012/gostudy/chatroom/client/model"
	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"github.com/fankys2012/gostudy/chatroom/common/message"
)

//初始化在线用户列表
var onlineUserList map[int]*cmodel.User = make(map[int]*cmodel.User, 20)

//存储当前用户信息
var CurUserInfo *model.UserInfo

func showUserList() {
	fmt.Println("当前在线用户列表:")
	for id, v := range onlineUserList {
		var onlineState string
		if v.UserState == cmodel.UserOnline {
			onlineState = "在线"
		} else if v.UserState == cmodel.UserOffline {
			onlineState = "离线"
		}
		fmt.Printf("[%d:%s/%s]\r\n", id, v.UserName, onlineState)
	}
	fmt.Println("\r\n")

}

//更新本地用户列表
func updateOnlineUserList(notifyUser *message.NotifyUserOnlineStateMes) {
	user, ok := onlineUserList[notifyUser.UserId]
	if ok {
		user.UserState = notifyUser.UserState
	} else {
		user = &cmodel.User{
			UserId:    notifyUser.UserId,
			UserName:  notifyUser.UserName,
			UserState: notifyUser.UserState,
		}
	}
	onlineUserList[notifyUser.UserId] = user
}

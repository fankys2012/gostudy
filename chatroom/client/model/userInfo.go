package model

import (
	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"net"
)

//存在当前用户的基本信息
type UserInfo struct {
	Conn net.Conn
	cmodel.User
}

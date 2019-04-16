package process

var (
	userMg *UserManager
)

type UserManager struct {
	onlineUsers map[int]*UserProcess //用户在线列表 实际存储的是UserProcess实例
}

func init() {
	userMg = &UserManager{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//添加
func (this *UserManager) AddOnlineUser(uprocess *UserProcess) {
	this.onlineUsers[uprocess.UserId] = uprocess
}

//删除
func (this *UserManager) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

func (this *UserManager) GetOnlineUserList() map[int]*UserProcess {
	return this.onlineUsers
}

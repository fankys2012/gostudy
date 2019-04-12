package model

import (
	"strconv"

	"github.com/fankys2012/gostudy/chatroom/common/cmodel"
	"github.com/garyburd/redigo/redis"
)

var (
	//声明一个全局实例，在服务启动后就立刻初始化
	MyUserDao *UserDao
)

type UserDao struct {
	pool *redis.Pool
}

//工厂模式
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	return &UserDao{
		pool: pool,
	}
}
func (this *UserDao) getUserById(conn redis.Conn, id int) (myuser *cmodel.User, err error) {
	key := "user:" + strconv.Itoa(id)
	value, err := redis.StringMap(conn.Do("hgetall", key))
	if err != nil {
		//在redis 中没有获取到值
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
			return
		}
		return
	}
	//存在err == nil 但是未获取到数据的情况
	if _, ok := value["id"]; !ok {
		err = ERROR_USER_NOTEXISTS
		return
	}

	var user cmodel.User
	myuser = &user
	// string -> int
	uid, err := strconv.Atoi(value["id"])
	user.UserId = uid
	user.UserName = value["name"]
	user.UserPwd = value["pwd"]
	return
}

func (this *UserDao) Login(id int, pwd string) (user *cmodel.User, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	user, err = this.getUserById(conn, id)
	if err != nil {
		return
	}

	if user.UserPwd != pwd {
		err = ERROR_USER_PWD_ERROR
		return
	}
	return
}

//通过ID判断用户是否存在
//false 不存在 ；true 存在
func (this *UserDao) ExistsById(id int) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	key := "user:" + strconv.Itoa(id)
	_, err := redis.Int(conn.Do("hget", key, "id"))
	if err != nil {
		//在redis 中没有获取到值
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (this *UserDao) Register(user *cmodel.User) (err error) {

	//用户ID是否存在
	exists, err := this.ExistsById(user.UserId)
	if err != nil {
		return
	} else if exists {
		err = ERROR_USER_EXISTS
		return
	}
	conn := this.pool.Get()
	defer conn.Close()
	key := "user:" + strconv.Itoa(user.UserId)
	_, err = conn.Do("hmset", key, "id", user.UserId, "name", user.UserName, "pwd", user.UserPwd)
	return

}

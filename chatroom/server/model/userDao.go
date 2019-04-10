package model

import (
	"encoding/json"

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
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		//在redis 中没有获取到值
		if err == redis.ErrNil {

		}
		return
	}

	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		return
	}
	return
}

func (this *UserDao) Login(id int, pwd string) (user *User, err error) {
	conn := this.pool.Get()

	defer conn.Close()

	user, err = this.getUserById(conn, id)
	if err != nil {
		return
	}

	if user.UserPwd != pwd {
		return
	}
	return
}

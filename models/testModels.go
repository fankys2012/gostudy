package models

import (
	"github.com/astaxie/beego/orm"     //引入beego的orm
	_ "github.com/go-sql-driver/mysql" //引入beego的mysql驱动
)

type User struct {
	Id     int
	Uname  string
	Upwd   string `orm:"size(255)"`
	Utoken string `orm:"size(1000)"`
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)                                                           //注册数据库驱动
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(192.168.56.120:3306)/test?charset=utf8") //注册一个别名为default的数据库
	orm.SetMaxIdleConns("default", 30)                                                                 //设置数据库最大空闲连接
	orm.SetMaxOpenConns("default", 30)                                                                 //设置数据库最大连接数
	orm.RegisterModel(new(User))
}

func (this *User) Add(title, body string, token string) (int64, error) {
	o := orm.NewOrm()
	arc := User{Uname: title, Upwd: body, Utoken: token}
	id, err := o.Insert(&arc)
	return id, err
}

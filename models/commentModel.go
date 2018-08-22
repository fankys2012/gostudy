package models

import (
	"fmt"
	"time"

	"github.com/fankys2012/gostudy/orm"
)

func init() {
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(192.168.56.120:3306)/test?charset=utf8") //注册一个别名为default的数据库
	fmt.Println("commentModel Init ...")
	orm.RegisterModel("", new(Comment))
}

type Comment struct {
	Comment_id      int    `orm:"auto"`
	Comment_content string `orm:"size(1000)"`
	User_id         int
	Comment_title   string `orm:"size(200)"`
	CreateTime      time.Time
}

func (this *Comment) Add(content, title string, user_id int) {
	o := orm.NewOrm()
	var datetime = time.Now()
	datetime.Format("2006-01-02 15:04:05")
	arc := &Comment{Comment_content: content, User_id: user_id, Comment_title: title, CreateTime: datetime}
	id, err := o.Insert(arc)
	fmt.Println(id, err)
}

func (c *Comment) GetOne(id int, content string) {
	o := orm.NewOrm()
	whereCols := []string{"Comment_id", "Comment_content"}

	arc := &Comment{Comment_id: id, Comment_content: content}
	o.Read(arc, whereCols)
	fmt.Println(arc)
}

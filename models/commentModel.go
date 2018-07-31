package models

import (
	"fmt"

	"github.com/fankys2012/gostudy/orm"
)

func init() {
	fmt.Println("commentModel Init ...")
	orm.RegisterModel(new(Comment))
}

type Comment struct {
	Id              int
	Comment_content string `orm:"size(1000)"`
	User_id         int
	Comment_title   string `orm:"size(200)"`
}

func (this *Comment) Add(content, title string, user_id int) {
	o := orm.NewOrm()
	arc := &Comment{Comment_content: content, User_id: user_id, Comment_title: title}
	id, err := o.Insert(arc)
	fmt.Println(id, err)
}

package main

import (
	"fmt"

	"github.com/fankys2012/gostudy/models"
)

// type User struct {
// 	id   int64
// 	name string
// }

func main() {
	var arc models.User
	arc.Add("title", "body", "token")
	// user := &User{id: 1, name: "zhansan"}
	// ormMysql := orm.NewBaseMysql()
	// ormMysql.Insert(user)
	fmt.Print("end")

}

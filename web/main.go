package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fankys2012/gostudy/models"
	// "github.com/astaxie/beego/orm"
)

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	var comt models.Comment
	comt.GetOne(3, "content")
	// commentModel := &models.Comment{}
	// commentModel.Add("content", "title", 11)

	// orm.RegisterDriver("mysql", orm.DRMySQL)                                                           //注册数据库驱动

	fmt.Println("starting")
	// path := getCurrentPath()
	// fmt.Println(path)
	// gostudy.Static["/static"] = "./views"
	// server := gostudy.NewServer("localhost", 9000)
	// server.Router(&controllers.UserController{})
	// server.Run()
}

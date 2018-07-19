package main

import (
	"github.com/fankys2012/gostudy"
	"github.com/fankys2012/gostudy/controllers"
)

func main() {
	server := gostudy.NewServer("localhost", 9000)
	server.Router(&controllers.UserController{})
	server.Run()
}

package controllers

import (
	"fmt"

	"github.com/fankys2012/gostudy"
)

type UserController struct {
	gostudy.App
}

func (this *UserController) AddAction() {
	fmt.Fprintf(this.W(), "User Add")
}

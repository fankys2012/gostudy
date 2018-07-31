package controllers

import (
	"fmt"

	"github.com/fankys2012/gostudy"
	"github.com/fankys2012/gostudy/models"
)

type UserController struct {
	gostudy.App
}

var arc models.Archives

func (this *UserController) AddAction() {
	arc.Add("title", "body", 32)
	fmt.Fprintf(this.W(), "User Add")
}

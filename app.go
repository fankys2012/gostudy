package gostudy

import (
	"net/http"
)

type IContext interface {
	Config(w http.ResponseWriter, r *http.Request)
}

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func (this *Context) Config(w http.ResponseWriter, r *http.Request) {
	this.w = w
	this.r = r
}

func (this *Context) CW() http.ResponseWriter {
	return this.w
}

type IApp interface {
	Init(ctx *Context)
}

type App struct {
	ctx  *Context
	Data map[string]interface{}
}

func (this *App) Init(ctx *Context) {
	this.ctx = ctx
	this.Data = make(map[string]interface{})
}
func (this *App) W() http.ResponseWriter {
	return this.ctx.w
}

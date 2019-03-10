package api

import (
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// ShowLoginPage dispatch request for Login Page
func (r *Router) ShowLoginPage(ctx *routing.Context) error {
	fasthttp.ServeFile(ctx.RequestCtx, "./static/login_form.html")
	return nil
}

// DoLogin dispatch request to perform Login action
func (r *Router) DoLogin(ctx *routing.Context) error {
	login := ctx.FormValue("username")
	pass := ctx.FormValue("password")

	fmt.Fprintf(ctx, "You logged in :) \nLogin: %s\nPassword: %s", login, pass)
	return nil
}

package routers

import (
	"git-me/controller"
	"git-me/middleware"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowHeaders:     []string{"Access-Control-Allow-Origin"},
		//AllowAllOrigins: true,
		//AllowOrigins: []string{"http://asai.com", "http://localhost:8080"},
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/r",
			beego.NSRouter("/login", &controller.UserController{}, "post:Login"),
			beego.NSRouter("/register", &controller.UserController{}, "post:Register"),
		),
		beego.NSNamespace("/user",
			beego.NSBefore(middleware.AuthLogin),
			beego.NSRouter("/logout", &controller.UserController{}, "get:Logout"),
			beego.NSRouter("/info", &controller.UserController{}, "post:UpdateInfo"),
			beego.NSRouter("/pass", &controller.UserController{}, "post:UpdatePass"),
		),
		beego.NSNamespace("/download",
			beego.NSRouter("/vid", &controller.DownloaderController{}, "post:ParseVideo"),
		),
	)

	beego.AddNamespace(ns)
}

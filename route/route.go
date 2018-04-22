package route

import (
	"git-me/controller"
	"git-me/middleware"

	"github.com/astaxie/beego"
)

func init() {
	beego.NSNamespace("/v1",
		beego.NSNamespace("/r",
			beego.NSRouter("/login", &controller.UserController{}, "post:Login"),
			beego.NSRouter("/register", &controller.UserController{}, "post:Register"),
		),
		beego.NSNamespace("/user",
			beego.NSBefore(middleware.AuthLogin),
			beego.NSRouter("/logout", &controller.UserController{}, "post:Logout"),
			beego.NSRouter("/info", &controller.UserController{}, "post:UpdateInfo"),
			beego.NSRouter("/pass", &controller.UserController{}, "post:UpdatePass"),
		),
	)

}

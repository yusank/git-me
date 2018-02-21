package route

import (
	"git-me/controller"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controller.UserController{}, "login")
}

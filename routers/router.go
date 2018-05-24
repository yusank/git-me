/*
 * MIT License
 *
 * Copyright (c) 2018 Yusan Kurban <yusankurban@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/04/01        Yusan Kurban
 */

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
			beego.NSRouter("/vid", &controller.DownloaderController{}, "get:ParseVideo"),
		),

		beego.NSNamespace("history",
			beego.NSBefore(middleware.AuthLogin),
			beego.NSRouter("/list", &controller.HistoryController{}, "get:List"),
		),

		beego.NSNamespace("collect",
			beego.NSBefore(middleware.AuthLogin),
			beego.NSRouter("/list", &controller.CollectController{}, "get:List"),
			beego.NSRouter("/add", &controller.CollectController{}, "get:AddCollect"),
		),
		beego.NSNamespace("/task",
			beego.NSBefore(middleware.AuthLogin),
			beego.NSRouter("/list/:id", &controller.TaskController{}, "get:ListTask"),
			beego.NSRouter("/add", &controller.TaskController{}, "post:AddTask"),
			beego.NSRouter("/update", &controller.TaskController{}, "post:UpdateTask"),
			beego.NSRouter("/del", &controller.TaskController{}, "post:DelTask"),
		),

		beego.NSNamespace("/inner",
			beego.NSRouter("/task-list", &controller.InnerController{}, "post:ListUserTasks"),
			beego.NSRouter("/task-upload", &controller.InnerController{}, "post:HandleEvent"),
		),
	)

	beego.AddNamespace(ns)
}

package main

import (
	"git-me/db"
	"github.com/astaxie/beego"
)

func main() {
	//cmd.Execute()

	if err := db.InitDB(); err != nil {
		panic(err)
	}

	if err := Init(); err != nil {
		panic(err)
	}

	beego.Run()
}

package main

import (
	"git-me/db"
	"git-me/env"
	"log"
	"os"

	"github.com/astaxie/beego"
)

func main() {
	if len(os.Args) < 2 {
		panic("env param is nil")
	}

	envParam := os.Args[1]
	log.Println("run env:", envParam)

	// Init Env
	if err := env.InitEnv(envParam); err != nil {
		panic(err)
	}

	// InitDB
	if err := db.InitDB(); err != nil {
		panic(err)
	}

	if err := Init(); err != nil {
		panic(err)
	}

	beego.Run(":17717")
}

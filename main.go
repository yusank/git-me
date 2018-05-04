package main

import (
	"os"

	"git-me/db"
	_ "git-me/routers"

	"git-me/extractors"
	"git-me/utils"
	"github.com/astaxie/beego"
)

func main() {
	envParam := ""
	if len(os.Args) == 2 {
		envParam = os.Args[1]
	}
	//logger.Printf("environment is [%s]", envParam)

	if envParam == "prod" {
		beego.BConfig.RunMode = "prod"
	} else {
		beego.BConfig.RunMode = "dev"
	}

	// InitDB
	if err := db.InitDB(); err != nil {
		panic(err)
	}

	if err := Init(); err != nil {
		panic(err)
	}

	// init downloader
	extractors.BeforeRun()
	utils.InitHttpClient()

	beego.Run(":17717")
}

package db

import "github.com/astaxie/beego"

var (
	Mongo *MongoDB
)

func InitDB() (err error) {
	if Mongo, err = NewMongoDB(beego.AppConfig.String("mongo_url"), "git-me");err != nil {
		return
	}

	return
}

func CloseDB() {
	Mongo.Close()
}
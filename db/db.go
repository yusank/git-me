package db

import (
	"fmt"

	"github.com/astaxie/beego"
)

var (
	// Mongo connection
	Mongo *MongoDB

	// Redis redis connection
	Redis *RedisDB
)

// InitDB init db connection
func InitDB() (err error) {
	// init mongo
	if Mongo, err = NewMongoDB(beego.AppConfig.String("mongo_url"), "git-me"); err != nil {
		return
	}

	// init redis
	optionCache := NewRedisOption("tcp", fmt.Sprintf("%s:%s", beego.AppConfig.String("redis_host"), beego.AppConfig.String("redis_port")), beego.AppConfig.String("redis_password"))
	optionCache.SetMaxAge(0)
	optionCache.SetMaxLength(0)
	if Redis, err = NewRedisDBWithOption(optionCache); err != nil {
		err = fmt.Errorf("redis err : %s", err)
		return
	}
	return
}

// CloseDB - close db
func CloseDB() {
	Mongo.Close()
	Redis.Option.Close()
}

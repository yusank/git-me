package db

import (
	"fmt"

	"git-me/env"
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
	if Mongo, err = NewMongoDB(env.MongoConnection, "gitMe"); err != nil {
		return
	}

	// init redis
	optionCache := NewRedisOption("tcp", fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort), env.RedisPassword)
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

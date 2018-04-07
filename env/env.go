package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

var (
	MongoConnection string
	RedisHost       string
	RedisPort       string
	RedisPassword   string

	Env string
)

func InitEnv(env string) (err error) {
	Env = env
	var envFile *os.File
	if envFile, err = os.Open(fmt.Sprintf("conf/%s.json", env)); err != nil {
		return err
	}
	defer envFile.Close()

	fileReader := bufio.NewReader(envFile)
	decoder := json.NewDecoder(fileReader)
	var jsonMap map[string]string
	if err = decoder.Decode(&jsonMap); err != nil {
		return err
	}

	MongoConnection = jsonMap["mongoConnection"]
	RedisHost = jsonMap["redisHost"]
	RedisPort = jsonMap["redisPort"]
	RedisPassword = jsonMap["redisPassword"]

	return
}

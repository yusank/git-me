package db

import (
	"gopkg.in/mgo.v2"
)

// modelIniters 实体初始化函数集合
var modelIniters = []*func(mongo *MongoDB){}

// FieldMappings 属性map
var FieldMappings map[string]map[string]string

// MongoDB 数据库实例结构.
type MongoDB struct {
	Session *mgo.Session
	DBName  string
}

func (m *MongoDB) Close() {
	m.Session.Close()
}

// NewMongoDB 创建mongo数据库连接实例.
func NewMongoDB(url, dbname string) (db *MongoDB, err error) {
	var session *mgo.Session
	session, err = mgo.Dial(url)
	if err != nil {
		return
	}

	db = &MongoDB{
		Session: session,
		DBName:  dbname,
	}
	session.SetMode(mgo.Strong, true)
	FieldMappings = make(map[string]map[string]string)
	for i := 0; i < len(modelIniters); i++ {
		(*modelIniters[i])(db)
	}
	return
}

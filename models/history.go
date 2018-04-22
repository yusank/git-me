package models

import (
	"git-me/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type History struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	UserID  bson.ObjectId `json:"userID" bson:"userID"`
	Site    string        `json:"site" bson:"site"`
	URL     string        `json:"url" bson:"url"`
	Type    int           `json:"type" bson:"type"`
	LastUse int64         `json:"lastUse" bson:"lastUse"`
}

const (
	HistoryC = "history"
)

var HistoryCollection *mgo.Collection

func PrepareHistory() error {
	HistoryCollection = db.Mongo.Session.DB(db.Mongo.DBName).C(HistoryC)

	if err := HistoryCollection.EnsureIndexKey("userID"); err != nil {
		return err
	}

	return HistoryCollection.EnsureIndexKey("site")
}

func (his *History) Insert() error {
	his.Id = bson.NewObjectId()
	return HistoryCollection.Insert(his)
}

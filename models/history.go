package models

import (
	"git-me/db"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type History struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	UserID    bson.ObjectId `json:"userId" bson:"userID"`
	Site      string        `json:"site" bson:"site"`
	URL       string        `json:"url" bson:"url"`
	Size      int64         `json:"size" bson:"size"`
	Quality   string        `json:"quality" bson:"quality"`
	Type      int           `json:"type" bson:"type"`
	LastUse   int64         `json:"lastUse" bson:"lastUse"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
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
	his.CreatedAt = time.Now().Unix()
	return HistoryCollection.Insert(his)
}

func (his *History) Update() error {
	his.LastUse = time.Now().Unix()
	return HistoryCollection.UpdateId(his.Id, his)
}

func GetHistory(userId, url string) (his *History, err error) {
	query := bson.M{"userID": bson.ObjectIdHex(userId), "url": url}
	his = new(History)

	err = HistoryCollection.Find(query).One(his)
	if err == mgo.ErrNotFound {
		err = nil
	}

	return
}

func ListHistory(userId string, page, size int) (list []*History, err error) {
	query := bson.M{"userID": bson.ObjectIdHex(userId)}
	list = make([]*History, 0)

	err = HistoryCollection.Find(query).Sort("createdAt").Skip((page - 1) * size).Limit(size).All(&list)
	return
}

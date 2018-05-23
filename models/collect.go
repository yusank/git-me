package models

import (
	"git-me/db"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CollectInfo struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	UserId    bson.ObjectId `json:"userId" bson:"userId"`
	URL       string        `json:"url" bson:"url"`
	Site      string        `json:"site" bson:"site"`
	Size      int64         `json:"size" bson:"size"`
	Quality   string        `json:"quality" bson:"quality"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
}

const (
	CollectC = "collect"
)

var CollectCollection *mgo.Collection

func PrepareCollect() error {
	CollectCollection = db.Mongo.Session.DB(db.Mongo.DBName).C(CollectC)

	if err := CollectCollection.EnsureIndexKey("userId"); err != nil {
		return err
	}

	return CollectCollection.EnsureIndexKey("site")
}

func (col *CollectInfo) Insert() error {
	col.Id = bson.NewObjectId()
	col.CreatedAt = time.Now().Unix()
	return CollectCollection.Insert(col)
}

func GetCollectByUserID(userId, url string) (col *CollectInfo, err error) {
	query := bson.M{"userId": bson.ObjectIdHex(userId), "url": url}
	col = new(CollectInfo)

	err = CollectCollection.Find(query).One(col)
	if err == mgo.ErrNotFound {
		err = nil
	}

	return
}

func ListCollect(userId string, page, size int) (list []*CollectInfo, err error) {
	list = make([]*CollectInfo, 0)
	err = CollectCollection.Find(bson.M{"userId": bson.ObjectIdHex(userId)}).Sort("createdAt").Skip((page - 1) * size).Limit(size).All(&list)
	return
}

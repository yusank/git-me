package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"git-me/db"
)

type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name"`
	Password  string        `json:"-"`
	HeadImg   string        `json:"headImg" bson:"headImg"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt"`
}

const (
	CollectionNameUser = "user"
)

var UserCollection *mgo.Collection

func PrepareUser() error {
	UserCollection = db.Mongo.Session.DB(db.Mongo.DBName).C(CollectionNameUser)

	return UserCollection.EnsureIndexKey("name")
}

func (u *User) Insert() error {
	u.Id = bson.NewObjectId()
	return UserCollection.Insert(u)
}

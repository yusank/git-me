package models

import (
	"git-me/db"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TaskInfo struct {
	Id        bson.ObjectId `bson:"_id"`
	UserId    bson.ObjectId `bson:"userId"`
	URL       string        `bson:"url"`
	Status    int           `bson:"status"`
	Sort      int           `bson:"sort"`
	Tp        int           `bson:"tp"`
	CreatedAt int64         `bson:"createdAt"`
	UpdateAt  int64         `bson:"updatedAt"`
}

var (
	TaskInfoCollection *mgo.Collection
)

const (
	CollectionTask = "task"
)

func PrepateTaskInfo() error {
	TaskInfoCollection = db.Mongo.Session.DB(db.Mongo.DBName).C(CollectionTask)

	idx := mgo.Index{
		Key: []string{"userId"},
	}

	return TaskInfoCollection.EnsureIndex(idx)
}

func (task *TaskInfo) Insert() error {
	task.Id = bson.NewObjectId()
	task.CreatedAt = time.Now().Unix()
	return TaskInfoCollection.Insert(task)
}

func (task *TaskInfo) Update() error {
	task.UpdateAt = time.Now().Unix()
	return TaskInfoCollection.UpdateId(task.Id, task)
}

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
	Type      int           `bson:"type"`
	CreatedAt int64         `bson:"createdAt"`
	UpdateAt  int64         `bson:"updatedAt"`
}

var (
	TaskInfoCollection *mgo.Collection
)

const (
	CollectionTask = "task"
)

const (
	TaskStatusDefault = iota
	TaskStatusDownlaoding
	TaskStatusFail
	TaskStatusFinish
)

func PrepareTaskInfo() error {
	TaskInfoCollection = db.Mongo.Session.DB(db.Mongo.DBName).C(CollectionTask)

	idx := mgo.Index{
		Key: []string{"userId"},
	}

	if err := TaskInfoCollection.EnsureIndex(idx); err != nil {
		return err
	}

	idx = mgo.Index{
		Key: []string{"url"},
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

func (task *TaskInfo) Delete() error {
	return TaskInfoCollection.RemoveId(task.Id)
}

func ListTaskInfo(userId string, page, size int) (list []*TaskInfo, err error) {
	list = make([]*TaskInfo, 0)
	err = TaskInfoCollection.FindId(bson.ObjectIdHex(userId)).Sort("createdAt").Skip((page - 1) * size).Limit(size).All(&list)

	return
}

func ListUnFinishedTaskInfo(userId string) (list []*TaskInfo, err error) {
	list = make([]*TaskInfo, 0)
	query := bson.M{
		"userId": bson.ObjectIdHex(userId),
		"status": bson.M{"$ne": TaskStatusFinish},
	}
	err = TaskInfoCollection.FindId(query).All(&list)

	return
}

func GetTaskInfoByUserAndUrl(userId, url string) (t *TaskInfo, err error) {
	t = new(TaskInfo)
	err = TaskInfoCollection.Find(bson.M{"userId": bson.ObjectIdHex(userId), "url": url}).One(t)
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	return
}

func GetTaskInfoById(id string) (t *TaskInfo, err error) {
	t = new(TaskInfo)
	err = TaskInfoCollection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(t)
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	return
}

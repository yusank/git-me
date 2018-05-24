/*
 * MIT License
 *
 * Copyright (c) 2018 Yusan Kurban <yusankurban@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/04/01        Yusan Kurban
 */

package models

import (
	"github.com/yusank/git-me/db"

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

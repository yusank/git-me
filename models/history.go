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

type History struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	UserID    bson.ObjectId `json:"userId" bson:"userID"`
	Site      string        `json:"site" bson:"site"`
	Title     string        `json:"title" bson:"title"`
	URL       string        `json:"url" bson:"url"`
	Size      int64         `json:"size" bson:"size"`
	Quality   string        `json:"quality" bson:"quality"`
	Type      string        `json:"type" bson:"type"`
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
		his = nil
	}

	return
}

func ListHistory(userId string, page, size int) (list []*History, err error) {
	query := bson.M{"userID": bson.ObjectIdHex(userId)}
	list = make([]*History, 0)

	err = HistoryCollection.Find(query).Sort("createdAt").Skip((page - 1) * size).Limit(size).All(&list)
	return
}

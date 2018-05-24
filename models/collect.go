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

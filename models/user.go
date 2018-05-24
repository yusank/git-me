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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"math/rand"
	"time"

	"github.com/yusank/git-me/db"
)

type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Nickname  string        `json:"nickname" bson:"nickname"`
	Email     string        `json:"email" bson:"email"`
	Password  string        `json:"-" bson:"password"`
	HeadImg   string        `json:"headImg" bson:"headImg"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt"`
}

type UserRegister struct {
	Name  string `json:"name" valid:"Required"`
	Email string `json:"email" valid:"Required"`
	Pass  string `json:"pass" valid:"Required"`
}

type UserLogin struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type UpdatePass struct {
	Name    string `json:"name" valid:"Required"`
	OldPass string `json:"oldPass" valid:"Required"`
	NewPass string `json:"newPass" valid:"Required;MinSize(6)"`
}

const (
	CollectionNameUser = "user"
)

var UserCollection *mgo.Collection

func PrepareUser() error {
	UserCollection = db.Mongo.Session.DB(db.Mongo.DBName).C(CollectionNameUser)

	idx := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	}

	if err := UserCollection.EnsureIndex(idx); err != nil {
		return err
	}

	idx = mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}

	return UserCollection.EnsureIndex(idx)
}

func (u *User) Insert() error {
	u.Id = bson.NewObjectId()
	u.Nickname = RandNickname()
	u.CreatedAt = time.Now().Unix()
	return UserCollection.Insert(u)
}

func (u *User) Update() error {
	u.UpdatedAt = time.Now().Unix()
	return UserCollection.UpdateId(u.Id, u)
}

func (u *User) Get() (*User, error) {
	query := bson.M{"name": u.Name}
	if u.Id != "" {
		query = bson.M{"_id": u.Id}
	}

	if err := UserCollection.Find(query).One(u); err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserById(id string) (u *User, err error) {
	u = new(User)
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	err = UserCollection.Find(query).One(u)
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	return
}

func RandNickname() string {
	rand.Seed(time.Now().Unix())
	length := 8
	data := make([]byte, length)
	var num int

	for i := 0; i < length; i++ {
		num = rand.Intn(57) + 65
		for {
			if num > 90 && num < 97 {
				num = rand.Intn(57) + 65
			} else {
				break
			}
		}
		data[i] = byte(num)
	}
	return string(data)
}

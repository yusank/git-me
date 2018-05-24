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

package db

import (
	"gopkg.in/mgo.v2"
)

// modelIniters 实体初始化函数集合
var modelIniters = []*func(mongo *MongoDB){}

// FieldMappings 属性map
var FieldMappings map[string]map[string]string

// MongoDB 数据库实例结构.
type MongoDB struct {
	Session *mgo.Session
	DBName  string
}

func (m *MongoDB) Close() {
	m.Session.Close()
}

// NewMongoDB 创建mongo数据库连接实例.
func NewMongoDB(url, dbname string) (db *MongoDB, err error) {
	var session *mgo.Session
	session, err = mgo.Dial(url)
	if err != nil {
		return
	}

	db = &MongoDB{
		Session: session,
		DBName:  dbname,
	}
	session.SetMode(mgo.Strong, true)
	FieldMappings = make(map[string]map[string]string)
	for i := 0; i < len(modelIniters); i++ {
		(*modelIniters[i])(db)
	}
	return
}

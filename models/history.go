package models

import "gopkg.in/mgo.v2/bson"

type History struct {
	Id      bson.ObjectId   `json:"id" bson:"_id"`
	UserID  bson.ObjectId   `json:"userID" bson:"userID"`
	History []*SingleRecord `json:"history" bosn:"history"`
}

type SingleRecord struct {
	URL     string `json:"url" bson:"url"`
	Type    int    `json:"type" bson:"type"`
	LastUse int64  `json:"lastUse" bson:"lastUse"`
}

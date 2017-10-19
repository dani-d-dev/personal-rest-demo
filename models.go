package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Player struct {
	ID bson.ObjectId	`json:"id" bson:"_id,omitempty"`
	Firstname string   	`json:"first_name,omitempty"`
	Lastname  string   	`json:"last_name,omitempty"`
	NickName  string   	`json:"nick_name,omitempty"`
	AvatarUrl  string   `json:"avatar_url,omitempty"`
}

type Players []Player

type Match struct {
	ID bson.ObjectId	`json:"id" bson:"_id,omitempty"`
	Startime  time.Time `json:"start_time"`
	Endtime  time.Time 	`json:"end_time"`
	Player1 string 		`json:"player_1"`
	Player2 string 		`json:"player_2"`
	Games []Game		`json:"games"`
}

type Game struct {
	result [2]int
}

type matches []Match

package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Provider struct {
	Platform string		`json:"platform"`
	Token string		`json:"token"`
}

type Player struct {
	ID string				`json:"id" bson:"uid"`
	Token string			`json:"token"`
	FirstName string		`json:"first_name"`
	LastName string			`json:"last_name",omitempty`
	NickName string			`json:"name",omitempty`
	Avatar interface{}		`json:"picture",omitempty`
	Location string			`json:"location",omitempty`
	IsLeftHanded bool		`json:"is_left_handed"`
	IsGripShakeHand bool 	`json:"is_grip_shakehand"`
}

type Players []Player

type Match struct {
	ID bson.ObjectId	`json:"id" bson:"_id,omitempty"`
	Startime  time.Time `json:"start_time"`
	Endtime  time.Time 	`json:"end_time"`
	Player1 string 		`json:"player_1"`
	Player2 string 		`json:"player_2"`
	Winner string 		`json:"winner,omitempty"`
	Loser string 		`json:"loser,omitempty"`
	Games []int			`json:"games"`
}

type matches []Match

type Team struct {
	ID bson.ObjectId		`json:"id" bson:"_id,omitempty"`
	Name string				`json:"name"`
	City string				`json:"city"`
	Description string		`json:"description"`
	Members Players			`json:"members"`
}

type Teams []Team

type Message struct {
	ID bson.ObjectId		`json:"id" bson:"_id"`
	SenderID	string		`json:"sender_id" bson:"sender_id"`
	ReceiverID	string		`json:"receiver_id" bson:"receiver_id"`
	Text		string		`json:"text"`
	Reason		ReasonType	`json:"reason"`
}

type Messages []Message

type ReasonType int16

const (
	JOIN_REQUEST ReasonType = 0
	CHALLENGE ReasonType = 1
)



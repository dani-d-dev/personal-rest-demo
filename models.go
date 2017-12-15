package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Provider struct {
	Platform string `json:"platform"`
	Token    string `json:"token"`
}

type Player struct {
	ID              string      `json:"id" bson:"uid"`
	Token           string      `json:"token"`
	FirstName       string      `json:"first_name" bson:"first_name"`
	LastName        string      `json:"last_name" bson:"last_name"`
	NickName        string      `json:"name" bson:"nick_name"`
	Avatar          interface{} `json:"picture"`
	Location        string      `json:"location"`
	IsLeftHanded    bool        `json:"is_left_handed" bson:"is_left_handed"`
	IsGripShakeHand bool        `json:"is_grip_shakehand" bson:"is_grip_shakehand"`
}

type Players []Player

type Match struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	StarTime time.Time     `json:"start_time" bson:"start_time"`
	EndTime  time.Time     `json:"end_time" bson:"end_time"`
	Player1  string        `json:"player_1" bson:"player_1"`
	Player2  string        `json:"player_2" bson:"player_2"`
	Winner   string        `json:"winner"`
	Loser    string        `json:"loser"`
	Games    []int         `json:"games"`
}

type matches []Match

type Team struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string        `json:"name"`
	City         string        `json:"city"`
	Description  string        `json:"description"`
	Members      Players       `json:"members"`
	JoinRequests []string      `json:"join_requests, omitempty" bson:"join_requests"`
}

type Teams []Team

type Message struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	SenderID   string        `json:"sender_id" bson:"sender_id"`
	ReceiverID string        `json:"receiver_id" bson:"receiver_id"`
	Text       string        `json:"text"`
	Reason     ReasonType    `json:"reason"`
	Declined   bool          `json:"declined"` // Whether has declined or accepted
}

type Messages []Message

type ReasonType int16

const (
	JOIN_REQUEST ReasonType = 0
	CHALLENGE    ReasonType = 1
)

package main

import "time"

type Player struct {
	ID        string   	`json:"id,omitempty"`
	Firstname string   	`json:"first_name,omitempty"`
	Lastname  string   	`json:"last_name,omitempty"`
	NickName  string   	`json:"nick_name,omitempty"`
	AvatarUrl  string   `json:"avatar_url,omitempty"`
}

type Players []Player

type Match struct {
	Startime  time.Time `json:"start_time"`
	Endtime  time.Time 	`json:"end_time"`
	Winner string 		`json:"winner"`
	Loser string 		`json:"loser"`
}

type matches []Match

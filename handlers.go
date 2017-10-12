package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

// Handlers

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main page goes here..."))
}

func PlayersList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(players)
}

func PlayerShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range players {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Player{})
}

func PlayerAdd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var player Player
	_ = json.NewDecoder(r.Body).Decode(&player)
	player.ID = params["id"]
	players = append(players, player)
	json.NewEncoder(w).Encode(players)
}

func PlayerDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range players {
		if item.ID == params["id"] {
			players = append(players[:index], players[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(players)
	}
}
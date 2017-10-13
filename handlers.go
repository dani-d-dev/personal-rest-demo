package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"fmt"
)

// Handlers

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main page goes here..."))
}

func PlayersList(w http.ResponseWriter, r *http.Request) {

	var results Players
	err := collection.Find(nil).All(&results)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
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

	decoder := json.NewDecoder(r.Body)

	var player Player
	err := decoder.Decode(&player)

	if err != nil {
		panic(err)
		fmt.Printf("Post handler failed")
	}

	defer r.Body.Close()

	collection.Insert(player)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(player)
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
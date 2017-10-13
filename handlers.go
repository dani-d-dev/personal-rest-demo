package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"gopkg.in/mgo.v2"
)

// Handlers

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main page goes here..."))
}

func PlayersList(w http.ResponseWriter, r *http.Request) {

	var results Players
	err := collection.Find(nil).Sort("-_id").All(&results)

	if err != nil {
		ErrorWithJSON(w, "Error finding record", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, results, 200)
}

func PlayerShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	player_id := params["id"]

	if !bson.IsObjectIdHex(player_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(player_id)

	var results Player
	err := collection.FindId(oid).One(&results)

	if err != nil {
		ErrorWithJSON(w, "Player not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, results, 200)
}

func PlayerInsert(w http.ResponseWriter, r *http.Request) {

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)

	if err != nil {
		ErrorWithJSON(w, "Cannot insert record on db", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	err = collection.Insert(player)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, player, 200)
}

func PlayerUpdate(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	player_id := params["id"]

	if !bson.IsObjectIdHex(player_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(player_id)

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)

	if err != nil {
		ErrorWithJSON(w, "Failed decoding json into object", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	document := bson.M{"_id":oid}
	change := bson.M{"$set":player}
	err = collection.Update(document, change)

	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed updating player: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "Record not found", http.StatusNotFound)
			return
		}
	}

	ResponseWithJSON(w, player, 200)
}

func PlayerDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	player_id := params["id"]

	if !bson.IsObjectIdHex(player_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(player_id)
	err := collection.RemoveId(oid)

	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed removing player: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "Player not found", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func ResponseWithJSON(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(result)
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}
package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
)

// Handlers

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main page goes here..."))
}

// Player CRUD

func PlayersList(w http.ResponseWriter, _ *http.Request) {

	var result Players
	err := playerCollection.Find(nil).Sort("-_id").All(&result)

	if err != nil || len(result) == 0{
		ErrorWithJSON(w, "Players not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, result, http.StatusOK)
}

func PlayerShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	player_id := params["id"]

	if !bson.IsObjectIdHex(player_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(player_id)

	var result Player
	err := playerCollection.FindId(oid).One(&result)

	if err != nil {
		ErrorWithJSON(w, "Player not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, result, http.StatusOK)
}

func PlayerInsert(w http.ResponseWriter, r *http.Request) {

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)

	if err != nil {
		ErrorWithJSON(w, "Cannot insert record on db", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	err = playerCollection.Insert(player)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, player, http.StatusOK)
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
	err = playerCollection.Update(document, change)

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

	ResponseWithJSON(w, player, http.StatusOK)
}

func PlayerDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	player_id := params["id"]

	if !bson.IsObjectIdHex(player_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(player_id)
	err := playerCollection.RemoveId(oid)

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

// Matches CRUD

func MatchList(w http.ResponseWriter, _ *http.Request) {

	var result matches
	err := matchCollection.Find(nil).Sort("-_id").All(&result)

	if err != nil || len(result) == 0 {
		ErrorWithJSON(w,"Matches not found",http.StatusNotFound)
		return
	}

	// TODO : Populate results with player data refered by id's

	//fmt.Printf("First element's id: ", bson.M{"_id":result[0].ID})

	ResponseWithJSON(w, result, http.StatusOK)
}

func MatchShow(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	match_id := params["id"]

	if !bson.IsObjectIdHex(match_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(match_id)

	var match Match
	err := matchCollection.FindId(oid).One(&match)

	if err != nil {
		ErrorWithJSON(w, "Match not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, match, http.StatusOK)
}

func MatchInsert(w http.ResponseWriter, r *http.Request) {

	var match Match
	err := json.NewDecoder(r.Body).Decode(&match)

	if err != nil {
		fmt.Printf("Error : ", err)
		ErrorWithJSON(w, "Cannot insert record on db", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	err = matchCollection.Insert(match)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, match, http.StatusOK)
}

// TODO : MatchUpdate

func MatchDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	match_id := params["id"]

	if !bson.IsObjectIdHex(match_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(match_id)
	err := matchCollection.RemoveId(oid)

	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed removing match: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "Match not found", http.StatusNotFound)
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
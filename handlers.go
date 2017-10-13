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
	err := collection.Find(nil).All(&results)

	if err != nil {
		ErrorWithJSON(w, "Error finding record", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, results, 200)
}

func PlayerShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var player Player
	err := collection.Find(bson.M{"id": id}).One(&player)

	if err != nil {
		ErrorWithJSON(w, "Error finding record", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, player, 200)
}

func PlayerAdd(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var player Player
	err := decoder.Decode(&player)

	if err != nil {
		ErrorWithJSON(w, "Cannot insert record on db", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	err = collection.Update(bson.M{"id": player.ID}, &player)
	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed update book: ", err)
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
	id := params["id"]

	err := collection.Remove(bson.M{"id":id})

	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed delete book: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "Book not found", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func ResponseWithJSON(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}
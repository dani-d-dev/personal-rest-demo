package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	fb "github.com/huandu/facebook"
	"crypto/sha256"
	"encoding/base64"
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

	// Pre-calculate winner & looser
	winner_id, loser_id := CalculateScore(match.Player1, match.Player2, match.Games)
	match.Winner = winner_id
	match.Loser = loser_id

	err = matchCollection.Insert(match)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, match, http.StatusOK)
}

// Get winner and looser tuple id's given a game array

func CalculateScore(player_id_1 string, player_id_2 string, games []int) (string, string) {

	var player1_victories = 0
	var player2_victories = 0

	for i := 0; i < len(games) + 1; i++ {
		if i % 2 == 0 {
			if i > 0 {
				var p1 = games[i-2]
				var p2 = games[i-1]
				if p1 > p2 {
					player1_victories +=1
				} else {
					player2_victories +=1
				}

				fmt.Printf("[Player 1: %d || Player 2: %d]\n", p1, p2)
			}
		}
	}

	if player1_victories > player2_victories {
		return player_id_1, player_id_2
	}
	return player_id_2, player_id_1
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

func Login(w http.ResponseWriter, r *http.Request) {

	var provider Provider
	err := json.NewDecoder(r.Body).Decode(&provider)

	if err != nil {
		ErrorWithJSON(w, "Error while parsing provider data", http.StatusNotAcceptable)
		return
	}

	res, err := fb.Get("/me", fb.Params{
		"fields": "id, name, first_name, last_name, picture.type(large)",
		"access_token": provider.Token,
	})

	if err != nil {
		ErrorWithJSON(w, "An Facebook API error has ocurred", http.StatusUnauthorized)
		return
	}

	var user FBUser
	res.Decode(&user)

	user.Token = encryptToken(provider.Token)
	user.Avatar = res.Get("picture.data.url")

	fmt.Printf("%v", user)

	// Save user
	err = userPlayerCollection.Insert(user)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, user, http.StatusOK)
}

func encryptToken(token string) string {
	tk := sha256.New()
	tk.Write([]byte(token))
	b := tk.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	ResponseWithJSON(w, "this should be done", http.StatusOK)
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
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	fb "github.com/huandu/facebook"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"os/user"
)

// Handlers

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main page goes here..."))
}

// Team Minimalistic CRUD

func TeamShow(w http.ResponseWriter, r *http.Request) {

	teamId := mux.Vars(r)["id"]

	team, err := FindByID(teamId, teamCollection)

	if err != nil {
		ErrorWithJSON(w, "Team not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, team, http.StatusOK)
}

func TeamList(w http.ResponseWriter, _ *http.Request) {

	teams, err := FindAll(teamCollection)

	if err != nil || len(teams) == 0 {
		ErrorWithJSON(w, "Teams not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, teams, http.StatusOK)
}

func TeamInsert(w http.ResponseWriter, r *http.Request) {

	var team Team
	err := json.NewDecoder(r.Body).Decode(&team)

	if err != nil {
		ErrorWithJSON(w, "Cannot insert record on db", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	// Check if an owner is provided on the teams array

	if len(team.Members) != 1 {
		ErrorWithJSON(w, "An owner of the team should be provided", http.StatusNotImplemented)
		return
	}

	err = teamCollection.Insert(team)

	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, team, http.StatusOK)
}

func TeamJoin(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	team_id := params["id"]

	// Search for a valid ID format and if it exists in DB

	if !bson.IsObjectIdHex(team_id) {
		ErrorWithJSON(w, "Identifier field not in hex format", http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(team_id)

	var team Team
	err := teamCollection.FindId(oid).One(&team)

	if err != nil {
		ErrorWithJSON(w, "Team not found", http.StatusNotFound)
		return
	}

	// Parse player from body request

	var player Player

	error := json.NewDecoder(r.Body).Decode(&player)

	if error != nil {
		ErrorWithJSON(w, "Cannot decode player json dictionary", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	// Check that player exists in db

	if !userExists(player.ID) {
		ErrorWithJSON(w, "Check that user exists", http.StatusNotFound)
		return
	}

	// Add new player to the team's array and update DB

	members := append(team.Members, player)
	team.Members = members

	teamCollection.UpdateId(team.ID, team)

	ResponseWithJSON(w, team, http.StatusOK)
}

func TeamAsk(w http.ResponseWriter, r *http.Request) {

	teamId := mux.Vars(r)["id"]

	var team Team
	er := FindEntityByID(teamId, teamCollection, &team)

	if er != nil {
		ErrorWithJSON(w, "Team not found", http.StatusNotFound)
		return
	}

	// Check if userId already exists before adding it

	for _, s := range team.JoinRequests {
		if s == teamId {
			ErrorWithJSON(w, "The provided user id was already requested", http.StatusNotAcceptable)
			return
		}
	}

	// Update team by adding player id to pending join requests

	joinRequests := append(team.JoinRequests, teamId)
	team.JoinRequests = joinRequests

	teamCollection.UpdateId(teamId, team)

	ResponseWithJSON(w, team, http.StatusOK)
}

// Message CRUD

func MessageList(w http.ResponseWriter, _ *http.Request) {

	messages, err := FindAll(messageCollection)

	if err != nil || len(messages) == 0 {
		ErrorWithJSON(w, "Messages not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, messages, http.StatusOK)
}

func MessageSend(w http.ResponseWriter, r *http.Request) {

	receiverId := mux.Vars(r)["id"]

	if !userExists(receiverId) {
		ErrorWithJSON(w, "The provided user id was not found, check if user exists", http.StatusNotFound)
		return
	}

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		ErrorWithJSON(w, "Cannot parse json into an object model", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	err = messageCollection.Insert(message)

	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, message, http.StatusOK)
}

// Player CRUD

func PlayersList(w http.ResponseWriter, _ *http.Request) {

	players, err := FindAll(playerCollection)

	if err != nil || len(players) == 0 {
		ErrorWithJSON(w, "Players not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, players, http.StatusOK)
}

func PlayerShow(w http.ResponseWriter, r *http.Request) {

	playerId := mux.Vars(r)["id"]
	player, err := FindByUID(playerId, playerCollection)

	if err != nil {
		ErrorWithJSON(w, "Player not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, player, http.StatusOK)
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

	document := bson.M{"_id": oid}
	change := bson.M{"$set": player}
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

	matches, err := FindAll(matchCollection)

	if err != nil || len(matches) == 0 {
		ErrorWithJSON(w, "Matches not found", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, matches, http.StatusOK)
}

func MatchShow(w http.ResponseWriter, r *http.Request) {

	matchId := mux.Vars(r)["id"]

	match, err := FindByID(matchId, matchCollection)

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
		"fields":       "id, name, first_name, last_name, picture.type(large)",
		"access_token": provider.Token,
	})

	if err != nil {
		ErrorWithJSON(w, "An Facebook API error has ocurred", http.StatusUnauthorized)
		return
	}

	var user Player
	res.Decode(&user)

	user.Token = encryptToken(provider.Token)
	user.NickName = "Your NickName"
	user.Avatar = res.Get("picture.data.url")
	user.IsLeftHanded = false
	user.IsGripShakeHand = true
	info, err := playerCollection.Upsert(bson.M{"uid": user.ID}, bson.M{"$set": user})
	log.Println("Update info:", info)

	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Insertion failed with error :", err)
		return
	}

	ResponseWithJSON(w, user, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	uid := r.Header.Get("X-User")
	pwd := r.Header.Get("X-Password")
	usr, err := authUser(uid, pwd)

	if err != nil {
		http.Error(w, "User not founded", http.StatusNotFound)
		return
	}

	// TODO : Do a soft delete from user (put a boolean flag or something)

	ResponseWithJSON(w, usr, http.StatusOK)
}

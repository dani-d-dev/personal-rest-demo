package main

import (
	"net/http"
	"log"
	"gopkg.in/mgo.v2/bson"
)

// auth middleware
func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Get user_id and token from headers

	uid := r.Header.Get("X-User")
	pwd := r.Header.Get("X-Password")

	// Query for a db user with the given credentials

	usr, err := getUser(uid, pwd)

	if err != nil {
		http.Error(rw, "Not Authorized", 401)
		log.Println("Failed with error: ", err)
		return
	}

	log.Println("User with id: ", usr.ID)
	next(rw, r)
}

func getUser(uid string, token string) (Player, error) {
	query := bson.M{"uid":uid, "token":token}
	var user Player
	err := playerCollection.Find(query).One(&user)

	if err != nil {
		return Player{}, err
	}

	return user, err
}
package main

import (
	"net/http"
	"log"
	"gopkg.in/mgo.v2/bson"
)

// auth middleware
func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Logging on the way there...")

	// Get user_id and token from headers

	usrId := r.Header.Get("X-User")
	pwd := r.Header.Get("X-Password")

	// Query for a db user with the given credentials

	query := bson.M{"user_id":usrId, "token":pwd}
	user := playerCollection.Find(query)

	if user != nil {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}

	/*
	if usrId == "usr1" && pwd == "secret123" {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}*/

	log.Println("Logging on the way back...")
}

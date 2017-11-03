package main

import (
	"net/http"
	"log"
	//"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

// auth middleware
func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	log.Println("Logging on the way there...")

	// Get user_id and token from headers

	usrId := r.Header.Get("X-User")
	pwd := r.Header.Get("X-Password")

	// Query for a db user with the given credentials

	query := bson.M{"uid":usrId, "token":pwd}
	var user FBUser
	err := userPlayerCollection.Find(query).One(&user)

	if err != nil {
		http.Error(rw, "Not Authorized", 401)
		log.Println("Failed with error: ", err)
		log.Println("Logging on the way back...")
		return
	}

	fmt.Printf("User with id: %s", user.ID)
	log.Println("Ok, authenticated")
	next(rw, r)

	/*
	if usrId == "usr1" && pwd == "secret123" {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}
	*/
}

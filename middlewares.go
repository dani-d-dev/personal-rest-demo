package main

import (
	"log"
	"net/http"
)

// auth middleware
func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Get user_id and token from headers

	uid := r.Header.Get("X-User")
	pwd := r.Header.Get("X-Password")

	// Query for a db user with the given credentials

	usr, err := authUser(uid, pwd)

	if err != nil {
		http.Error(rw, "Not Authorized", 401)
		log.Println("Failed with error: ", err)
		return
	}

	log.Println("User with id: ", usr.ID)
	next(rw, r)
}

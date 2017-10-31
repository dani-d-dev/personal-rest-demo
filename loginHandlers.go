package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
	"net/http"
	"encoding/json"
	"log"
	//"gopkg.in/mgo.v2/bson"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var provider Provider
	err := json.NewDecoder(r.Body).Decode(&provider)

	if err != nil {
		ErrorWithJSON(w, "Error while parsing provider data", http.StatusNotAcceptable)
		return
	}

	res, err := fb.Get("/me", fb.Params{
		"fields": "first_name, last_name, picture",
		"access_token": provider.Token,
	})

	if err != nil {
		ErrorWithJSON(w, "An Facebook API error has ocurred", http.StatusUnauthorized)
		return
	}

	var user FBUser
	res.Decode(&user)

	fmt.Println("print first_name in struct:", user.FirstName)
	fmt.Println("print first_name in struct:", user.LastName)

	ResponseWithJSON(w, user, http.StatusOK)
}

// auth middleware
func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Logging on the way there...")

	usrId := r.Header.Get("X-User")
	pwd := r.Header.Get("X-Password")
	//query := bson.M{"user_id":usrId, "token":pwd}
	//user := playerCollection.Find(query)

	/*
	if user != nil {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}
	*/


	if usrId == "usr1" && pwd == "secret123" {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}

	log.Println("Logging on the way back...")
}


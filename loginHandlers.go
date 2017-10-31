package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
	"net/http"
	"encoding/json"
	"log"
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

	if r.Header.Get("password")== "secret123" {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}

	log.Println("Logging on the way back...")
}


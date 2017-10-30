package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
	"net/http"
	"encoding/json"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var provider Provider
	err := json.NewDecoder(r.Body).Decode(&provider)

	if err != nil {
		ErrorWithJSON(w, "Error while parsing provider data", http.StatusNotAcceptable)
		return
	}

	res, err := fb.Get("/me", fb.Params{
		"fields": "first_name",
		"access_token": provider.Token,
	})

	if err != nil {
		ErrorWithJSON(w, "An Facebook API error has ocurred", http.StatusUnauthorized)
		return
	}

	var user FBUser
	res.Decode(&user)
	fmt.Println("print first_name in struct:", user.FirstName)
	ResponseWithJSON(w, user, http.StatusOK)
}


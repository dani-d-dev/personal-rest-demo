package main

import (
	"log"
	"net/http"
	"os"
)

var players Players

func main() {
	setupMockedData()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":5050", router))
}

func setupMockedData() {
	players = append(players, Player{"1", "Ma", "Long", "The Dragon", "-"})
	players = append(players, Player{"2", "Timo", "Boll", "The Nice Guy", "-"})
}

func port() string {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
		return ""
	}

	return port
}



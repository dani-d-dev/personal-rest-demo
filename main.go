package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
)

var players Players
var collection = getSession().DB("godata").C("user")

func main() {

	setupMockedData()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+port(), router))
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

func getSession() *mgo.Session {
	sess, err := mgo.Dial("mongodb://mlab-dani:dani1234@ds117935.mlab.com:17935/godata")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	return sess
}



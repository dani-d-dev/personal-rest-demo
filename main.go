package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
)

var playerCollection = getSession().DB("godata").C("user")
var matchCollection = getSession().DB("godata").C("match")

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+port(), router))
}

func port() string {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
		return ""
	}

	return port

}

func mongoURL() string {
	url := os.Getenv("MONGO_URL")

	if url == "" {
		log.Fatal("$MONGO_URL must be set")
		return ""
	}

	return url
}

func getSession() *mgo.Session {
	sess, err := mgo.Dial(mongoURL())
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	return sess
}

//mongodb://mlab-dani:dani1234@ds117935.mlab.com:17935/godata



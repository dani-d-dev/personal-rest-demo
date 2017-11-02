package main

import (
	"os"
	"log"
	"gopkg.in/mgo.v2"
	"fmt"
)

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
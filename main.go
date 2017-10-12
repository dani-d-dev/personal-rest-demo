package main

import (
	"log"
	"net/http"
	"os"
)

var people []Person

func main() {
	setupMockedData()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+port(), router))
}

func setupMockedData() {
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
}

func port() string {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
		return ""
	}

	return port
}



package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var players Players

func main() {

	setupMockedData()

	uri:="mongodb://mlab-dani:dani1234@ds117935.mlab.com:17935/godata"

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connected to db")

	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB("godata").C("user")

	/*
	err = collection.Insert(&Player{"1", "Ma", "Long", "The Dragon", "-"},
							&Player{"2", "Timo", "Boll", "The Nice Guy", "-"})
	if err != nil {
		log.Fatal("Problem inserting data: ", err)
		return
	}
	*/

	result := Player{}
	err = collection.Find(bson.M{"firstname": "Ma"}).One(&result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	fmt.Println("Last name:", result.Lastname)

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



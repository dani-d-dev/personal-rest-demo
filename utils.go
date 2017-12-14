package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
)

func port() string {

	os.Setenv("PORT", "5050")

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
		return ""
	}

	return port
}

func mongoURL() string {

	os.Setenv("MONGO_URL", "mongodb://mlab-dani:dani1234@ds117935.mlab.com:17935/godata")

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

// Token encryption

func encryptToken(token string) string {
	tk := sha256.New()
	tk.Write([]byte(token))
	b := tk.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}

// Get winner and looser tuple id's given a game array

func CalculateScore(player_id_1 string, player_id_2 string, games []int) (string, string) {

	var player1_victories = 0
	var player2_victories = 0

	for i := 0; i < len(games)+1; i++ {
		if i%2 == 0 {
			if i > 0 {
				var p1 = games[i-2]
				var p2 = games[i-1]
				if p1 > p2 {
					player1_victories += 1
				} else {
					player2_victories += 1
				}

				fmt.Printf("[Player 1: %d || Player 2: %d]\n", p1, p2)
			}
		}
	}

	if player1_victories > player2_victories {
		return player_id_1, player_id_2
	}
	return player_id_2, player_id_1
}

func ResponseWithJSON(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(result)
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

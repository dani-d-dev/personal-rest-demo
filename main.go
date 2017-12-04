package main

import (
	//"log"
	//"net/http"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	//"goji.io/middleware"
	"net/http"
	"log"
)

var teamCollection = getSession().DB("godata").C("team")
var playerCollection = getSession().DB("godata").C("player")
var matchCollection = getSession().DB("godata").C("match")

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	// Auth

	auth := router.PathPrefix("/auth").Subrouter()
	auth.Path("/login").HandlerFunc(Login)
	auth.Path("/logout").HandlerFunc(Logout)

	// Api

	api := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)
	api.HandleFunc("/player/all", PlayersList).Methods("GET")
	api.HandleFunc("/player/{id}",PlayerShow).Methods("GET")
	api.HandleFunc("/player", PlayerInsert).Methods("POST")
	api.HandleFunc("/player/{id}", PlayerUpdate).Methods("PUT")
	api.HandleFunc("/player/{id}", PlayerDelete).Methods("DELETE")
	api.HandleFunc("/match/all", MatchList).Methods("GET")
	api.HandleFunc("/match/{id}", MatchShow).Methods("GET")
	api.HandleFunc("/match", MatchInsert).Methods("POST")
	api.HandleFunc("/match/{id}", MatchDelete).Methods("DELETE")
	api.HandleFunc("/team/all", TeamList).Methods("GET")
	api.HandleFunc("/team/", TeamInsert).Methods("POST")
	api.HandleFunc("/team/{id}/join", TeamJoin).Methods("POST")

	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(AuthMiddleware),
		negroni.Wrap(api),
	))

	n:= negroni.Classic()
	n.UseHandler(router)
	log.Fatal(http.ListenAndServe(":"+port(), n))
}



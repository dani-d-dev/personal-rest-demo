package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(route.HandleFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Players",
		"GET",
		"/api/v1/player/all",
		PlayersList,
	},
	Route{
		"PlayerShow",
		"GET",
		"/api/v1/player/{id}",
		PlayerShow,
	},
	Route{
		"PlayerInsert",
		"POST",
		"/api/v1/player",
		PlayerInsert,
	},
	Route{
		"PlayerUpdate",
		"PUT",
		"/api/v1/player/{id}",
		PlayerUpdate,
	},
	Route{
		"PlayerDelete",
		"DELETE",
		"/api/v1/player/{id}",
		PlayerDelete,
	},
	Route{
		"MatchList",
		"GET",
		"/api/v1/match/all",
		MatchList,
	},
	Route{
		"MatchShow",
		"GET",
		"/api/v1/match/{id}",
		MatchShow,
	},
	Route{
		"MatchInsert",
		"POST",
		"/api/v1/match",
		MatchInsert,
	},
	Route{
		"MatchDelete",
		"DELETE",
		"/api/v1/match/{id}",
		MatchDelete,
	},
	Route{
		"Login",
		"POST",
		"api/v1/login",
		Login,
	},
}

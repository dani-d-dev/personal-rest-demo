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
		"PlayerAdd",
		"POST",
		"/api/v1/player",
		PlayerAdd,
	},
	Route{
		"PlayerDelete",
		"DELETE",
		"/api/v1/player/{id}",
		PlayerDelete,
	},
}

package main

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

//JSON Sample
/*
{
"id" : "3",
"firstname" : "Ben",
"lastname" : "El Gordo",
"address" : {
"city" : "Derry",
"state" : "Minnessota"
}
}*/
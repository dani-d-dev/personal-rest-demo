package main

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FindAll(collection *mgo.Collection, result interface{}) error {

	return collection.Find(nil).Sort("-_id").All(result)
}

func FindByID(id string, collection *mgo.Collection, result interface{}) error {

	if !bson.IsObjectIdHex(id) {
		return errors.New("The provided id is not in hex format")
	}

	oid := bson.ObjectIdHex(id)

	return collection.FindId(oid).One(result)
}

func FindByUID(uid string, collection *mgo.Collection, result interface{}) error {

	return collection.Find(bson.M{"uid": uid}).One(result)
}

func authUser(uid string, token string) (Player, error) {
	query := bson.M{"uid": uid, "token": token}
	var user Player
	err := playerCollection.Find(query).One(&user)

	if err != nil {
		return Player{}, err
	}

	return user, err
}

func userExists(uid string) bool {

	var user Player
	err := playerCollection.Find(bson.M{"uid": uid}).One(&user)

	if err != nil {
		return false
	}

	return true
}

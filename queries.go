package main

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FindAll(collection *mgo.Collection) ([]interface{}, error) {

	var result []interface{}
	err := collection.Find(nil).Sort("-_id").All(&result)

	return result, err
}

func FindByID(id string, collection *mgo.Collection) (interface{}, error) {

	var result interface{}

	if !bson.IsObjectIdHex(id) {
		return result, errors.New("The provided id is not in hex format")
	}

	oid := bson.ObjectIdHex(id)
	err := collection.FindId(oid).One(&result)

	return result, err
}

func FindEntityByID(id string, collection *mgo.Collection, result interface{}) error {

	if !bson.IsObjectIdHex(id) {
		return errors.New("The provided id is not in hex format")
	}

	oid := bson.ObjectIdHex(id)
	err := collection.FindId(oid).One(result)

	return err
}

func FindByUID(uid string, collection *mgo.Collection) (interface{}, error) {

	var result interface{}
	err := collection.Find(bson.M{"uid": uid}).One(&result)

	return result, err
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

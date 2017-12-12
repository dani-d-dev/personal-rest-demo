package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FindAll(collection *mgo.Collection) ([]interface{}, error){

	var result []interface{}
	err := collection.Find(nil).Sort("-_id").All(&result)

	return result, err
}

func FindByID(id string, collection *mgo.Collection) (interface{}, error) {

	var result interface{}
	err := collection.Find(bson.M{"_id": id}).One(&result)

	return result, err
}

func FindByUID(uid string, collection *mgo.Collection) (interface{}, error) {

	var result interface{}
	err := collection.Find(bson.M{"uid": uid}).One(&result)

	return result, err
}



func getUser(uid string, token string) (Player, error) {
	query := bson.M{"uid":uid, "token":token}
	var user Player
	err := playerCollection.Find(query).One(&user)

	if err != nil {
		return Player{}, err
	}

	return user, err
}

func userExists(uid string) (bool) {

	var user Player
	err := playerCollection.Find(bson.M{"uid":uid}).One(&user)

	if err != nil {
		return false
	}

	return true
}
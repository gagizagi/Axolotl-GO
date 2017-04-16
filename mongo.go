package main

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

//DBanimeList is a database collection "animeList"
var DBanimeList *mgo.Collection

func dbConn() {
	//Set url
	url := "localhost"
	if os.Getenv("OPENSHIFT_MONGODB_DB_URL") != "" {
		url = os.Getenv("OPENSHIFT_MONGODB_DB_URL")
	}

	//Connect to url
	dbSession, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}

	//authenticate
	if os.Getenv("OPENSHIFT_MONGODB_DB_USERNAME") != "" &&
		os.Getenv("OPENSHIFT_MONGODB_DB_PASSWORD") != "" {
		creds := mgo.Credential{
			Username: os.Getenv("OPENSHIFT_MONGODB_DB_USERNAME"),
			Password: os.Getenv("OPENSHIFT_MONGODB_DB_PASSWORD")}

		dbSession.Login(&creds)
	}
	dbSession.SetMode(mgo.Monotonic, true)

	//Set collections
	DBanimeList = dbSession.DB("axolotl").C("animeList")

	//Success
	log.Printf("Connected to MongoDB on %s\n", url)
}

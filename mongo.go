package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

//DBanimeList is a database collection "animeList"
var DBanimeList *mgo.Collection

func dbConn() {
	//Set url
	url := fmt.Sprintf(
		"mongodb://%s:%s@ds141434.mlab.com:41434/axolotl",
		os.Getenv("MLAB_USER"),
		os.Getenv("MLAB_PASS"))

	//Connect to url
	dbSession, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}

	dbSession.SetMode(mgo.Monotonic, true)

	//Set collections
	DBanimeList = dbSession.DB("axolotl").C("animeList")

	//Success
	log.Println("Connected to MongoDB")
}

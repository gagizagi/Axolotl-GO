package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

// DBanimeList is a database collection "animeList"
var DBanimeList *mgo.Collection

// DBserverList is a database collection "serverList"
var DBserverList *mgo.Collection

func init() {
	requireEnvVars("AX_DATABASE_HOST", "AX_DATABASE_PORT", "AX_DATABASE_DB")
}

func dbConn() {
	// Make the mongo connection URL
	// mongodb://<user>:<password>@<host>:<port>/<database>
	url := "mongodb://"
	if os.Getenv("AX_DATABASE_USERNAME") != "" &&
		os.Getenv("AX_DATABASE_PASSWORD") != "" {

		url += fmt.Sprintf("%s:%s@",
			os.Getenv("AX_DATABASE_USERNAME"),
			os.Getenv("AX_DATABASE_PASSWORD"))
	}
	url += fmt.Sprintf("%s:%s/%s",
		os.Getenv("AX_DATABASE_HOST"),
		os.Getenv("AX_DATABASE_PORT"),
		os.Getenv("AX_DATABASE_DB"))

	// Connect to database URL
	dbSession, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}

	// Set session mode
	// https://godoc.org/labix.org/v2/mgo#Session.SetMode
	dbSession.SetMode(mgo.Monotonic, true)

	// Set collections
	dbName := os.Getenv("AX_DATABASE_DB")
	DBanimeList = dbSession.DB(dbName).C("animeList")
	DBserverList = dbSession.DB(dbName).C("serverList")

	// Success
	log.Println(fmt.Sprintf("Connected to MongoDB %s on port %s",
		os.Getenv("AX_DATABASE_HOST"),
		os.Getenv("AX_DATABASE_PORT")))
}

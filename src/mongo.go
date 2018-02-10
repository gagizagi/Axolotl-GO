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
	requireEnvVars("DATABASE_HOST", "DATABASE_PORT", "DATABASE_DB")
}

func dbConn() {
	// Make the mongo connection URL
	// mongodb://<user>:<password>@<host>:<port>/<database>
	url := "mongodb://"
	if os.Getenv("DATABASE_USERNAME") != "" &&
		os.Getenv("DATABASE_PASSWORD") != "" {

		url += fmt.Sprintf("%s:%s@",
			os.Getenv("DATABASE_USERNAME"),
			os.Getenv("DATABASE_PASSWORD"))
	}
	url += fmt.Sprintf("%s:%s/%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_DB"))

	// Connect to database URL
	dbSession, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}

	// Set session mode
	// https://godoc.org/labix.org/v2/mgo#Session.SetMode
	dbSession.SetMode(mgo.Monotonic, true)

	// Set collections
	dbName := os.Getenv("DATABASE_DB")
	DBanimeList = dbSession.DB(dbName).C("animeList")
	DBserverList = dbSession.DB(dbName).C("serverList")

	// Success
	log.Println(fmt.Sprintf("Connected to MongoDB %s on port %s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT")))
}

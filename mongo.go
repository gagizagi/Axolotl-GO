package main

import(
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
)

var DBsession *mgo.Session
var DBanimeList *mgo.Collection

func init() {
	var err error
	url := "localhost"
	if os.Getenv("OPENSHIFT_MONGODB_DB_URL") != "" {
		url = os.Getenv("OPENSHIFT_MONGODB_DB_URL")
	}

	DBsession, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	if os.Getenv("OPENSHIFT_MONGODB_DB_USERNAME") != "" && os.Getenv("OPENSHIFT_MONGODB_DB_PASSWORD") != "" {
		creds := mgo.Credential{
			Username:os.Getenv("OPENSHIFT_MONGODB_DB_USERNAME"),
			Password:os.Getenv("OPENSHIFT_MONGODB_DB_PASSWORD")}
		
		DBsession.Login(&creds)
	}
	DBsession.SetMode(mgo.Monotonic, true)
	
	DBanimeList = DBsession.DB("axolotl").C("animeList")
	fmt.Printf("Connected to MongoDB on %s\n", url)
}
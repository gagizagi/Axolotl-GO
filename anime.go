package main

import(
	"fmt"
	"time"
	"math/rand"
	"gopkg.in/mgo.v2/bson"
)

type Anime struct {
	Id 			string 		`bson:"id"`//done
	Name 		string 		`bson:"name"`//done
	Href 		string 		`bson:"href"`//TODO
	Episode 	int 		`bson:"ep"`//done
	Subs 		[]string 	`bson:"subs"`//TODO
	LastUpdate 	time.Time 	`bson:"lastUpdate"`//TODO
}

//Gets every anime in animeList db and
//returns it as []Anime
func Get_anime_list() ([]Anime) {
	var result []Anime
	err := DBanimeList.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	
	return result
}

//Inserts a new anime entry to db
//Generates unique id if it doesent exist
func (a *Anime) Insert() (error) {
	if a.Id == "" {
		a.GenId()
	}
	a.LastUpdate = time.Now()
	err := DBanimeList.Insert(a)
	if err != nil {
		return err
	}
	return nil
}

//Updates the db entry with up-to-date
//episode number
func (a *Anime) UpdateEp() {
	updateQuery := bson.M{
		"$set":bson.M{
			"ep":a.Episode,
			"lastUpdate":time.Now(),
		},
	}
	DBanimeList.Update(bson.M{"name":a.Name}, updateQuery)
}

//Generates a unique 3char alphanumeric ID
//Not case sensitive
func (a *Anime) GenId() {
	var id,byteList string
	byteList = "0123456789abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 3; i++ {
		id += string(byteList[rand.Intn(len(byteList))])
	}
	
	if n, _ := DBanimeList.Find(bson.M{"id":id}).Count(); n == 0 {
		a.Id = id
	}else{
		a.GenId()
	}
}

//Checks if there is already an entry in db
//with the same id OR the same name
//returns true if it already exists
func (a Anime) Exists() (bool) {
	query := bson.M{
		"$or": []interface{}{
			bson.M{"id":a.Id},
			bson.M{"name":a.Name},
		},
	}
	if n, _ := DBanimeList.Find(query).Count(); n == 0 {
		return false
	} else {
		return true
	}
}

//Checks if episode # already exists in db
//Returns true if episode in db is outdated and 
//needs to be updated and false if db is already
//up to date
func (a Anime) NewEpisode() (bool) {
	query := bson.M{
		"name":a.Name,
		"ep":bson.M{
			"$lt":a.Episode,
		},
	}
	
	if n, _ := DBanimeList.Find(query).Count(); n == 0 {
		return false
	} else {
		return true
	}
}
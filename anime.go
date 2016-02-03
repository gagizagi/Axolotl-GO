package main

import(
	"strings"
	"fmt"
	"time"
	"math/rand"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"github.com/PuerkitoBio/goquery"
)

type AnimeList []Anime

type Anime struct {
	Id 			string 		`bson:"id"`
	Name 		string 		`bson:"name"`
	Href 		string 		`bson:"href"`
	Episode 	int 		`bson:"ep"`
	Subs 		[]string 	`bson:"subs"`
	LastUpdate 	time.Time 	`bson:"lastUpdate"`
}

//Gets every anime in animeList db and returns it as AnimeList
func Get_anime_list() (result AnimeList) {
	err := DBanimeList.Find(nil).Sort("lastUpdate").All(&result)
	if err != nil {
		fmt.Println(err)
	}	
	return
}

//db maintanance called every interval time.Duration
//deletes entries over 22 days old
//gets urls for entries that don't have them
func Maintain_anime_list(interval time.Duration) {
	const LIMIT = 22 * 24 * time.Hour
	for _ = range time.Tick(interval) {
		animeList := Get_anime_list()
		now := time.Now()
		for _, a := range animeList {
			if len(a.Href) < 5 {
				a.GetHref()
			}
			if now.Sub(a.LastUpdate) > LIMIT {
				a.Remove()
			}
		}
		fmt.Println("AUTO-MAINTANANCE: animeList updated at", now)
	}
}

//Inserts a new anime entry to db
//Generates unique id if it doesent exist
func (a *Anime) Insert() {
	if a.Id == "" {
		a.GenId()
	}
	a.LastUpdate = time.Now()
	DBanimeList.Insert(a)
}

//Remove anime from db by Anime.Name or Anime.Id
func (a Anime) Remove() {
	if a.Name != "" {
		DBanimeList.Remove(bson.M{"name":a.Name})
	} else if a.Id != "" {
		DBanimeList.Remove(bson.M{"id":a.Id})	
	}
}

//Updates the db entry with up-to-date episode number
func (a *Anime) UpdateEp() {
	updateQuery := bson.M{
		"$set":bson.M{
			"ep":a.Episode,
			"lastUpdate":time.Now(),
		},
	}
	change := mgo.Change{
		Update: updateQuery,
		Upsert: false,
		Remove: false,
		ReturnNew: true,
	}
	DBanimeList.Find(bson.M{"name":a.Name}).Apply(change, a)
}

//Adds new sub Name to the db entry of Anime.Id 
func (a *Anime) AddSub(sub string) {
	updateQuery := bson.M{
		"$addToSet":bson.M{
			"subs":sub,
		},
	}
	change := mgo.Change{
		Update: updateQuery,
		Upsert: false,
		Remove: false,
		ReturnNew: true,
	}
	DBanimeList.Find(bson.M{"id":a.Id}).Apply(change, a)
}

//Removes the sub from the db entry of Anime.Id
func (a *Anime) RemoveSub(sub string) {
	updateQuery := bson.M{
		"$pull":bson.M{
			"subs":sub,
		},
	}
	change := mgo.Change{
		Update: updateQuery,
		Upsert: false,
		Remove: false,
		ReturnNew: true,
	}
	DBanimeList.Find(bson.M{"id":a.Id}).Apply(change, a)
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

//Gets href for Anime.Name
func (a *Anime) GetHref() {
	doc, err := goquery.NewDocument("http://horriblesubs.info/current-season/")
	if err != nil {
		fmt.Println(err)
	} else {
		doc.Find(".ind-show.linkful").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Find("a").Attr("title")
			url, _ := s.Find("a").Attr("href")
			if strings.ToLower(name) == strings.ToLower(a.Name) {
				newHref := fmt.Sprintf("http://horriblesubs.info%s",url)
				updateQuery := bson.M{
					"$set":bson.M{
						"href":newHref,
						"lastUpdate":time.Now(),
					},
				}
				DBanimeList.Update(bson.M{"name":a.Name}, updateQuery)
				return
			}
		})
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

//returns length of AnimeList
//used for sort interface
func (a AnimeList) Len() (int) {
	return len(a)
}

//Checks if index i should sort before index j
//used for sort interface
func (a AnimeList) Less(i, j int) (bool) {
	if len(a[i].Subs) > len(a[j].Subs) {
		return true
	}
	return false
}

//Swaps the values of i and j indexes
//used for sort interface
func (a AnimeList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
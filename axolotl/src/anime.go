package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type animeList []anime

type anime struct {
	ID         string    `bson:"id"`
	Name       string    `bson:"name"`
	Href       string    `bson:"href"`
	Episode    int       `bson:"ep"`
	Subs       []string  `bson:"subs"`
	LastUpdate time.Time `bson:"lastUpdate"`
	Show       bool      `bson:"show"`
}

//LIMIT is a time constant for 15 Days
const LIMIT = 15 * 24 * time.Hour

//Gets every anime in animeList db and returns it as AnimeList type
func getAnimeList() (result animeList) {
	defer panicRecovery()

	err := DBanimeList.Find(nil).Sort("lastUpdate").All(&result)
	if err != nil {
		panic(fmt.Sprintf("Error querying MongoDB in function: %s - %s", "getAnimeList", err))
	}

	return
}

//Get every anime this user.id is subscribed to
func getAnimeListForUser(userID string) (result animeList) {
	defer panicRecovery()

	err := DBanimeList.Find(bson.M{"subs": userID}).All(&result)
	if err != nil {
		panic(fmt.Sprintf("Error querying MongoDB in function: %s - %s", "getAnimeListForUser", err))
	}

	return
}

// maintainAnimeList db maintanance function
// changes the 'show' database field of all entries over LIMIT days old to false
// tries to update the 'href' database field for entries that don't have one yet
// TODO: Potentially merge this with rssReader
func maintainAnimeList() {
	// Number of hidden and/or updated entries
	hidden, updated := 0, 0
	newAnimeList := getAnimeList()
	now := time.Now()

	for _, a := range newAnimeList {
		if len(a.Href) < 5 {
			updated += a.GetHref() //TODO: Limit HS scraping to maximum 1 per maintanance, instead of 1 per empty href per maintanance
		}
		if now.Sub(a.LastUpdate) > LIMIT {
			hidden += a.Hide()
		}
	}
	log.Printf("AUTO-MAINTANANCE: animeList updated! (hidden: %d | updated: %d)\n",
		hidden, updated)
}

// Insert inserts a new anime entry to db
// generates a unique ID if it isn't present in the object yet
func (a *anime) Insert() {
	if a.ID == "" {
		a.GenID()
	}
	a.LastUpdate = time.Now()
	a.Show = true
	DBanimeList.Insert(a)
}

// Hide changes database field 'show' for this anime to false
func (a anime) Hide() (success int) {
	defer panicRecovery()

	success = 0

	if a.Name != "" {
		err := DBanimeList.Update(bson.M{"name": a.Name}, bson.M{"show": false})

		if err != nil {
			panic(fmt.Sprintf("Error updating MongoDB document in anime method: %s - %s", "Hide", err))
		}

		success = 1
	} else if a.ID != "" {
		err := DBanimeList.Update(bson.M{"id": a.ID}, bson.M{"show": false})

		if err != nil {
			panic(fmt.Sprintf("Error updating MongoDB document in anime method: %s - %s", "Hide", err))
		}

		success = 1
	}

	return
}

// UpdateEp updates the db entry with the new episode number
func (a *anime) UpdateEp() {
	updateQuery := bson.M{
		"$set": bson.M{
			"ep":         a.Episode,
			"lastUpdate": time.Now(),
		},
	}
	change := mgo.Change{
		Update:    updateQuery,
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	DBanimeList.Find(bson.M{"name": a.Name}).Apply(change, a)
}

//Adds new sub Name to the db entry of Anime.Id
func (a *anime) AddSub(sub string) {
	updateQuery := bson.M{
		"$addToSet": bson.M{
			"subs": sub,
		},
	}
	change := mgo.Change{
		Update:    updateQuery,
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	DBanimeList.Find(bson.M{"id": a.ID}).Apply(change, a)
}

//Removes the sub from the db entry of Anime.Id
func (a *anime) RemoveSub(sub string) {
	updateQuery := bson.M{
		"$pull": bson.M{
			"subs": sub,
		},
	}
	change := mgo.Change{
		Update:    updateQuery,
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	DBanimeList.Find(bson.M{"id": a.ID}).Apply(change, a)
}

// GenID generates a unique 3char alphanumeric ID
// not case sensitive
func (a *anime) GenID() {
	var id, byteList string
	byteList = "0123456789abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 3; i++ {
		id += string(byteList[rand.Intn(len(byteList))])
	}

	if n, _ := DBanimeList.Find(bson.M{"id": id}).Count(); n == 0 {
		a.ID = id
	} else {
		a.GenID()
	}
}

//Gets href for Anime.Name
func (a *anime) GetHref() (success int) {
	success = 0
	//NOTE: Cloudflare scraping not needed for now
	//scrapper := "http://scraper-422.rhcloud.com/?href="
	target := "http://horriblesubs.info/current-season/"

	doc, err := goquery.NewDocument( /*scrapper + */ target)
	if err != nil {
		log.Println(err)
	} else {
		doc.Find(".ind-show.linkful").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Find("a").Attr("title")
			url, _ := s.Find("a").Attr("href")
			if strings.ToLower(name) == strings.ToLower(a.Name) {
				newHref := fmt.Sprintf("http://horriblesubs.info%s", url)
				updateQuery := bson.M{
					"$set": bson.M{
						"href":       newHref,
						"lastUpdate": time.Now(),
					},
				}
				DBanimeList.Update(bson.M{"name": a.Name}, updateQuery)
				success = 1
			}
		})
	}
	return
}

// Exists checks if there is already an entry in db
// with the same id OR the same name
// Returns true if it already exists
func (a anime) Exists() bool {
	query := bson.M{
		"$or": []interface{}{
			bson.M{"id": a.ID},
			bson.M{"name": a.Name},
		},
	}
	if n, _ := DBanimeList.Find(query).Count(); n == 0 {
		return false
	}
	return true
}

// NewEpisode checks if episode # already exists in db
// Returns true if episode in db is outdated and
// needs to be updated and false if db is already
// up to date
func (a anime) NewEpisode() bool {
	query := bson.M{
		"name": a.Name,
		"ep": bson.M{
			"$lt": a.Episode,
		},
	}

	if n, _ := DBanimeList.Find(query).Count(); n == 0 {
		return false
	}
	return true
}

//returns length of AnimeList
//used for sort interface
func (a animeList) Len() int {
	return len(a)
}

//Checks if index i should sort before index j
//used for sort interface
func (a animeList) Less(i, j int) bool {
	if len(a[i].Subs) > len(a[j].Subs) {
		return true
	}
	return false
}

//Swaps the values of i and j indexes
//used for sort interface
func (a animeList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

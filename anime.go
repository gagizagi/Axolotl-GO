package main

import(
	"fmt"
	"time"
)

type AnimeList []Anime

type Anime struct {
	Id 			string 		`bson:"id,omitempty"`
	Name 		string 		`bson:"name,omitempty"`
	Href 		string 		`bson:"href,omitempty"`
	Episode 	int 		`bson:"ep,omitempty"`
	Subs 		[]string 	`bson:"subs"`
	LastUpdate 	time.Time 	`bson:"lastUpdate,omitempty"`
}

func Get_anime_list() (AnimeList) {
	var result AnimeList
	err := DBanimeList.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	
	return result
}

func (a Anime) Insert() (error){
	err := DBanimeList.Insert(&a)
	if err != nil {
		return err
	}
	return nil
}
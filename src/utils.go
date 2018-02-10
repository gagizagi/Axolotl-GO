package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	botStartTime time.Time
	botResponses int
	botMessages  int
)

func init() {
	rand.Seed(time.Now().Unix())
	botStartTime = time.Now()
}

func getUptime() string {
	return time.Now().Sub(botStartTime).String()
}

//getInfo returns the bot infromation as a discord formated string
func getInfo() (result string) {
	result += "```"
	result += fmt.Sprintf("Name: %s\n", discordCfg.Name)
	result += fmt.Sprintf("Uptime: %s\n", getUptime())
	result += fmt.Sprintf("Guilds: %d\n", len(discordCfg.Guilds))
	result += fmt.Sprintf("Anime channels: %d\n", len(discordCfg.AnimeChannels))
	result += fmt.Sprintf("Unique subscribers: %d\n", getUniqueSubs())
	result += fmt.Sprintf("Messages read: %d\n", botMessages)
	result += fmt.Sprintf("Message responses: %d\n", botResponses)
	result += "```"

	return
}

//getUniqueSubs returs the amount of unique subscribers in the database
func getUniqueSubs() int {
	var uniqueSubs []string
	err := DBanimeList.Find(bson.M{}).Distinct("subs", &uniqueSubs)
	if err != nil {
		log.Print("MongoDB error: ", err)
	}

	return len(uniqueSubs)
}

//appendUnique is a function for appending a string to string array
//only appends the string if it doesn't already exsist in the array
func appendUnique(slice []string, id string) []string {
	for _, s := range slice {
		if s == id {
			return slice
		}
	}
	return append(slice, id)
}

//removeItem removes a string from the array and returns the new array
//first param is the array to remove from
//second param is the string to remove from array (removes all instances)
func removeItem(slice []string, item string) []string {
	for i, value := range slice {
		if value == item {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

//checks if any of the provided required ENV variables is missing
//logs the error and exits the program if any of them is missing
func requireEnvVars(args ...string) {
	for _, v := range args {
		if os.Getenv(v) == "" {
			log.Fatal("Missing ENV variable: ", v)
		}
	}
}

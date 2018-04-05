package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2/bson"
)

var (
	botStartTime time.Time
	botReads     int
	botMessages  int
	botResponses int
)

func init() {
	rand.Seed(time.Now().Unix())
	botStartTime = time.Now()
}

func getUptime() string {
	return time.Now().Sub(botStartTime).String()
}

// getUniqueSubs
// Returns the number of unique subscribers in the database
func getUniqueSubs() int {
	defer panicRecovery()

	var uniqueSubs []string
	err := DBanimeList.Find(bson.M{}).Distinct("subs", &uniqueSubs)
	if err != nil {
		panic(fmt.Sprintf("Error querying MongoDB in function: %s - %s", "getUniqueSubs", err))
	}

	return len(uniqueSubs)
}

// appendUnique is a function for appending a string to string array
// only appends the string if it doesn't already exsist in the array
func appendUnique(slice []string, id string) []string {
	for _, s := range slice {
		if s == id {
			return slice
		}
	}
	return append(slice, id)
}

// removeItem removes a string from the array and returns the new array
// first param is the array to remove from
// second param is the string to remove from array (removes all instances)
func removeItem(slice []string, item string) []string {
	for i, value := range slice {
		if value == item {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// requireEnvVars checks if any of the required ENV variables are missing
// logs the error and exits the program if any of them are missing
func requireEnvVars(args ...string) {
	for _, v := range args {
		if os.Getenv(v) == "" {
			log.Fatal("Missing ENV variable: ", v)
		}
	}
}

// tickerHelper funs a function in intervals
// d is the duration betwen function execution
// f is the executing function
// runFirst true will run a function once before starting a ticker
func tickerHelper(d time.Duration, f func(), runFirst bool) {
	ticker := time.NewTicker(d)

	if runFirst {
		f()
	}

	for range ticker.C {
		f()
	}
}

// panicRecovery will recover from any panics during function execution
// will log the panic message
// try to stick to the following format:
// Error <doing something>: [variables] - [error]
func panicRecovery() {
	if r := recover(); r != nil {
		log.Println(r)
	}
}

// scrapeHS scrapes the HorribleSubs website HTML
// Returns a goquery document or nil if it fails
func scrapeHS() *goquery.Document {
	defer panicRecovery()

	//NOTE: Cloudflare scraping not needed for now
	//scrapper := "http://scraper-422.rhcloud.com/?href="
	target := "http://horriblesubs.info/shows/"

	doc, err := goquery.NewDocument( /*scrapper + */ target)
	if err != nil {
		panic(fmt.Sprintf("Error trying to scrape website: %s - %s", target, err))
	}

	return doc
}

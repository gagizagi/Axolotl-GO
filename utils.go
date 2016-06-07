package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	owm "github.com/briandowns/openweathermap"
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

func getWeather(location string) (result string) {
	w, err := owm.NewCurrent("C", "en")
	if err != nil {
		log.Print("getWeather() => newCurrent() error:\n", err)
		return "error"
	}

	w.CurrentByName(location)
	result += fmt.Sprintf("***Weather for %s (%s)***\n\n", w.Name, w.Sys.Country)
	result += fmt.Sprintf("```Temperature: %.1f°C\n", w.Main.Temp)
	result += fmt.Sprintf("Humidity: %d%%\n", w.Main.Humidity)
	for _, item := range w.Weather {
		result += fmt.Sprintf("%s: %s\n", item.Main, item.Description)
	}
	result += fmt.Sprintf("Wind speed: %.1fm/s\n", w.Wind.Speed)
	result += fmt.Sprintf("Clouds: %d%%```", w.Clouds.All)
	return
}

//getInfo returns the bot infromation as a discord formated string
func getInfo() (result string) {
	result += "```"
	result += fmt.Sprintf("Name: %s\n", discordCfg.Name)
	result += fmt.Sprintf("Uptime: %s\n", getUptime())
	result += fmt.Sprintf("Guilds: %d\n", len(discordCfg.Guilds))
	result += fmt.Sprintf("Anime channels: %d\n", len(discordCfg.AnimeChannels))
	result += fmt.Sprintf("Messages read: %d\n", botMessages)
	result += fmt.Sprintf("Message responses: %d\n", botResponses)
	result += "```"

	return
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

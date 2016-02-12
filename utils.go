package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	owm "github.com/briandowns/openweathermap"
)

var botStartTime time.Time

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

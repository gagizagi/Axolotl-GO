package main

import (
	"fmt"
	"time"
	"math/rand"
	owm "github.com/briandowns/openweathermap"
)

var startTime time.Time

func init() {
	rand.Seed(time.Now().Unix())
	startTime = time.Now()
}

func GetUptime() (string) {
	return time.Now().Sub(startTime).String()
}

func GetWeather(location string) (result string) {
	w, err := owm.NewCurrent("C", "en")
	if err != nil {
		fmt.Println(err)
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
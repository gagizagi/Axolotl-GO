package main

import (
	"os"
	"time"
)

var ircClient = ircConfig{
	Server:   "irc.rizon.net:6667",
	Channels: []string{"#422", "#HORRIBLESUBS"},
	Username: "Axolotl",
	Nickname: "Axolotl",
	Verbose:  false,
	Debug:    false,
}

var discordClient = discordConfig{
	Boss:     "110846867473973248",
	Username: os.Getenv("DISCORD_USERNAME"),
	Password: os.Getenv("DISCORD_PASSWORD"),
	Debug:    false,
}

func main() {
	dbConn()

	go ircConnStart(&ircClient)
	go discordConnStart(&discordClient)
	go maintainAnimeList(10 * time.Hour)

	webServer() //Last
}

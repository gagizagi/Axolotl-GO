package main

import (
	"os"
	"time"
)

var ircClient = ircConfig{
	Server:   "irc.rizon.net:6667",
	Channels: []string{"#422", "#HORRIBLESUBS"},
	Nickname: "Axolotl-moe",
	Verbose:  false,
	Debug:    false,
}

var discordClient = discordConfig{
	Boss:  os.Getenv("DISCORD_BOSS"),
	Token: os.Getenv("DISCORD_TOKEN"),
	Debug: false,
}

func main() {
	dbConn()

	go ircConnStart(&ircClient)
	go discordConnStart(&discordClient)
	go maintainAnimeListProcess(10 * time.Hour)

	webServer() //Last
}

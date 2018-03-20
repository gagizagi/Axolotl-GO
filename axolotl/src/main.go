package main

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var ircClient = ircConfig{
	Server:   "irc.rizon.net:6667",
	Channels: []string{"#422", "#HORRIBLESUBS"},
	Nickname: os.Getenv("IRC_NICKNAME"),
	Verbose:  false,
	Debug:    false,
}

var discordClient = discordConfig{
	Boss:         os.Getenv("DISCORD_BOSS"),
	Token:        os.Getenv("DISCORD_TOKEN"),
	AnimeChannel: "anime",
	Debug:        false,
}

func main() {
	defer webServer()
	dbConn()

	go ircConnStart(&ircClient)
	go discordStart(&discordClient)
	go maintainAnimeListProcess(10 * time.Hour)
}

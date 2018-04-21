package main

import (
	"os"
	"time"
)

var discordClient = discordConfig{
	Boss:      os.Getenv("DISCORD_BOSS"),
	Token:     os.Getenv("DISCORD_TOKEN"),
	AvatarURL: "https://axolotl.gazzy.online/static/axolotl.jpg",
	Debug:     false,
}

func init() {
	requireEnvVars("DATABASE_HOST", "DATABASE_PORT", "DATABASE_DB", "DISCORD_BOSS", "DISCORD_TOKEN")
}

func main() {
	defer webServer()
	dbConn()

	go discordStart(&discordClient)
	go tickerHelper(10*time.Minute, rssReader, false)
	go tickerHelper(10*time.Hour, maintainAnimeList, true)
}

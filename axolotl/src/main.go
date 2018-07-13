package main

import (
	"os"
	"time"
)

var discordClient = discordConfig{
	Boss:      os.Getenv("DISCORD_BOSS"),
	Token:     os.Getenv("DISCORD_TOKEN"),
	AvatarURL: "https://cdn.discordapp.com/attachments/450681214781882379/450777891626811393/75ecbeedecd376ae8a32b8d7b0a5cfc3.png",
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

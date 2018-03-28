package main

import (
	"os"
	"time"
)

var discordClient = discordConfig{
	Boss:         os.Getenv("DISCORD_BOSS"),
	Token:        os.Getenv("DISCORD_TOKEN"),
	AnimeChannel: "anime",
	Debug:        false,
}

func init() {
	requireEnvVars("DATABASE_HOST", "DATABASE_PORT", "DATABASE_DB", "DISCORD_BOSS", "DISCORD_TOKEN")
}

func main() {
	defer webServer()
	dbConn()

	go rssReader()
	go discordStart(&discordClient)
	go maintainAnimeListProcess(10 * time.Hour)
}

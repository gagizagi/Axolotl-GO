package main

import (
	"os"
	"time"
)

var discordClient = discordConfig{
	Boss:         os.Getenv("DISCORD_BOSS"),
	Token:        os.Getenv("DISCORD_TOKEN"),
	AvatarURL:    "https://camo.githubusercontent.com/c40a9a73cc03b760a567df127d7fcebb59724580/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f617661746172732f3138353137373835313739393031313332392f37306336653365396135373633626564396664663336353130653831323733612e6a7067",
	AnimeChannel: "anime",
	Debug:        false,
}

func init() {
	requireEnvVars("DATABASE_HOST", "DATABASE_PORT", "DATABASE_DB", "DISCORD_BOSS", "DISCORD_TOKEN")
}

func main() {
	defer webServer()
	dbConn()

	go tickerHelper(10*time.Minute, rssReader, true)
	go discordStart(&discordClient)
	go tickerHelper(10*time.Hour, maintainAnimeList, true)
}

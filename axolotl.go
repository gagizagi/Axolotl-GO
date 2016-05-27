package main

import "os"

var ircClient = ircConfig{
	Server:   "irc.rizon.net:6667",
	Channels: []string{"#422"},
	Nickname: "Axolotl-Dev",
	Verbose:  false,
	Debug:    false,
}

var discordClient = discordConfig{
	Boss:  os.Getenv("DISCORD_BOSS"),
	Token: os.Getenv("DISCORD_TOKEN"),
	Debug: false,
}

func main() {
	//dbConn()

	go ircConnStart(&ircClient)
	go discordConnStart(&discordClient)
	//go maintainAnimeListProcess(10 * time.Hour)

	webServer() //Last
}

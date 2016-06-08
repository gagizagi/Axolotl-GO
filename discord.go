package main

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

//discordConfig is a struct containing user configuration for discord connection
type discordConfig struct {
	Boss          string
	Name          string
	AnimeChannel  string
	AnimeChannels []string
	Guilds        []string
	Token         string
	Debug         bool
}

//Declare constants
const (
	//Help string lists all discord chat commands for this bot
	discordHelp string = "***LIST OF BOT COMMANDS***\n" +
		"Fields in [] are optional\n" +
		"Fields in <> are mandatory\n\n" +
		"```" +
		"!help [bot] (Lists all bot commands)\n" +
		"!uptime [bot](Prints current bot uptime)\n" +
		"!sub <id> (Subscribe to anime and get notified when a new episode is released)\n" +
		"!unsub <id> (Unsubscribe from anime)\n" +
		"!w <location> (Prints current weather)\n" +
		"!info [bot] (Prints bot information)" +
		"```\n" +
		"Find anime list at http://gazzy.space/anime"
)

//Declare variables
var (
	discord       *discordgo.Session           //Discord client
	discordCfg    *discordConfig               //Discord options
	relevantRegex = regexp.MustCompile(`^!\w`) //Discord regex msg parser
)

//Starts discord connection and sets handlers and behavior
func discordStart(c *discordConfig) {
	//Updates discordCfg variable with c parameter data
	//For future use outside this function
	discordCfg = c

	//Connect to discord with token
	var err error
	discord, err = discordgo.New(c.Token)
	if err != nil {
		log.Fatal("Error initializing discord in discordStart() function!\n", err)
	}

	//Set behavior & assign handlers
	discord.Debug = c.Debug
	discord.AddHandler(discordMsgHandler)
	discord.AddHandler(discordReadyHandler)
	discord.AddHandler(discordNewGuildHandler)
	discord.AddHandler(discordLeaveGuildHandler)
	discord.AddHandler(discordNewChannelHandler)
	discord.AddHandler(discordLeaveChannelHandler)
	discord.AddHandler(discordChannelUpdateHandler)
	discord.Open() //Opens discord connection
}

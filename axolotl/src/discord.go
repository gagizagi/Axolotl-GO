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

type msgObject struct {
	Message string
	Channel string
}

//Declare variables
var (
	discord       *discordgo.Session           //Discord client
	discordCfg    *discordConfig               //Discord options
	relevantRegex = regexp.MustCompile(`^!\w`) //Discord regex msg parser
	msgChan       chan msgObject
)

//Starts discord connection and sets handlers and behavior
func discordStart(c *discordConfig) {
	msgChan = make(chan msgObject)
	//Updates discordCfg variable with c parameter data
	//For future use outside this function
	discordCfg = c

	//Connect to discord with token
	var err error
	discord, err = discordgo.New("Bot " + c.Token)
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

	go discordMsgDispatcher(msgChan)
}

// discordMsgDispather receives discord message strings through a chan
// and sends them to appropriate discord text channels
func discordMsgDispatcher(c <-chan msgObject) {
	for msg := range c {
		_, err := discord.ChannelMessageSend(msg.Channel, msg.Message)
		if err != nil {
			log.Printf("\nError sending discord message on channel %s: %s - %s",
				msg.Channel, msg.Message, err)
		}
	}
}

package main

import (
	"log"

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

type discordCommandHandler func([]string, *discordgo.MessageCreate)

type msgObject struct {
	Message string
	Channel string
	Embed   *discordgo.MessageEmbed
}

var (
	discord     *discordgo.Session               // Discord client
	discordCfg  *discordConfig                   // Discord options
	msgChan     chan msgObject                   // Channel for dispatching discord messages
	commandList map[string]discordCommandHandler // Maps command strings to appropriate handler functions
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
		log.Fatal("Error initializing discord connection!\n", err)
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

	commandList = make(map[string]discordCommandHandler)
	commandList["HELP"] = helpCommand
	commandList["SUB"] = subCommand
	commandList["UNSUB"] = unsubCommand
	commandList["MYSUBS"] = mySubs
	commandList["UPTIME"] = uptime
	commandList["STATUS"] = setStatus
	commandList["INFO"] = botInfo
	commandList["GUILDS"] = guilds

	go discordMsgDispatcher(msgChan)
}

// discordMsgDispather receives discord message strings through a chan
// and sends them to appropriate discord text channels
func discordMsgDispatcher(c <-chan msgObject) {
	for msg := range c {
		// Increment the number of messages this bot has responded to
		botResponses++
		var err error

		if msg.Embed != nil {
			_, err = discord.ChannelMessageSendEmbed(msg.Channel, msg.Embed)
		} else if msg.Message != "" {
			_, err = discord.ChannelMessageSend(msg.Channel, msg.Message)
		}
		if err != nil {
			botResponses--
			log.Printf("\nError sending discord message on channel %s: %s - %s",
				msg.Channel, msg.Message, err)
		}
	}
}

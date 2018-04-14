package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// discordConfig is a struct containing user configuration for discord connection
type discordConfig struct {
	Boss      string
	Name      string
	AvatarURL string
	Token     string
	Debug     bool
}

type discordCommandHandler func([]string, *discordgo.MessageCreate)

// msgObject contains the information for sending discord messages
// Message is the message string to send (If embed is present serves as a short description of embed)
// Channel is the channel to which the message is sent
// Embed is the embed object to send (Replaces Message string if present)
// Author is the user ID that triggered this message
type msgObject struct {
	Message string
	Channel string
	Embed   *discordgo.MessageEmbed
	Author  string
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
	discord.AddHandler(discordLeaveGuildHandler)
	discord.AddHandler(discordLeaveChannelHandler)
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
	commandList["NOTIFYHERE"] = updateAnimeChannelCommand
	commandList["PREFIX"] = updatePrefixCommand
	//FIXME: add 'mode' command

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

		// Handle errors
		if err != nil && strings.Contains(err.Error(), "Missing Permissions") {
			// If bot is missing permissions to send text messages on this channel
			// sends a private message to the author of the message instead
			if msg.Author != "" {
				privCh, _ := discord.UserChannelCreate(msg.Author)
				if msg.Embed != nil {
					_, err = discord.ChannelMessageSendEmbed(privCh.ID, msg.Embed)
				} else if msg.Message != "" {
					_, err = discord.ChannelMessageSend(privCh.ID, msg.Message)
				}
				msg.Message = "FALLBACK PRIVATE MESSAGE" // Set error message in case it failed
			}
		}

		if err != nil {
			botResponses--
			log.Printf("Error sending discord message on channel %s: %s - %s",
				msg.Channel, msg.Message, err)
		}
	}
}

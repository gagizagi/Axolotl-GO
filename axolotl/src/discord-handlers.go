package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//discordReadyHandler sets required data after a successful connection
func discordReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	//Sets this bots name in discordCfg.Name
	discordCfg.Name = strings.ToUpper(r.User.Username)

	//Iterates through guilds
	for _, guild := range r.Guilds {
		//Gets a list of channels for this guild
		channels, _ := s.GuildChannels(guild.ID)
		//Iterates through a list of channels for this guild
		for _, channel := range channels {
			//If it finds a channel with name 'anime'
			//it will add it to discordCfg.AnimeChannels array
			//which is used when sending new episode messages in discord
			if channel.Name == discordCfg.AnimeChannel {
				discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, channel.ID)
			}
		}
	}

	//Logs successful connection to discord was established
	log.Println("Connected to discord as", discordCfg.Name)
}

//discordNewChannelHandler handles joining new channels
func discordNewChannelHandler(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Name == discordCfg.AnimeChannel {
		discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, c.ID)
	}
}

//discordLeaveChannelHandler handles leaving channels
func discordLeaveChannelHandler(s *discordgo.Session, c *discordgo.ChannelDelete) {
	if c.Name == discordCfg.AnimeChannel {
		discordCfg.AnimeChannels = removeItem(discordCfg.AnimeChannels, c.ID)
	}
}

//discordChannelUpdateHandler handles channels being updated
func discordChannelUpdateHandler(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	if c.Name == discordCfg.AnimeChannel {
		discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, c.ID)
	} else {
		discordCfg.AnimeChannels = removeItem(discordCfg.AnimeChannels, c.ID)
	}
}

//discordNewGuildHandler handlers joining a guild
func discordNewGuildHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	discordCfg.Guilds = appendUnique(discordCfg.Guilds, g.Name)

	for _, c := range g.Channels {
		if c.Name == discordCfg.AnimeChannel {
			discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, c.ID)
		}
	}
}

//discordLeaveGuildHandler handles leaving/being kicked from a guild
//FIXME: remove this anime channel from config when kicked
func discordLeaveGuildHandler(s *discordgo.Session, g *discordgo.GuildDelete) {
	discordCfg.Guilds = removeItem(discordCfg.Guilds, g.Name)
}

// discordMsgHandler is a handler function for incoming discord messages
func discordMsgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Increment the number of messages this bot has read
	botReads++

	// Split received message into different parts (space seperated)
	// prefix is the first char of the first part
	// command is the rest of first part (not including first char)
	// args is an array of all parts after first one
	var prefix string
	var command string
	var args []string
	args = strings.Fields(m.Content)
	if len(args) > 0 {
		prefix = args[0][0:1]
		command = args[0][1:]
		args = strings.Fields(m.Content)[1:]
	}

	// If bot needs to parse this message, parse it
	// otherwise ignore it and get out of the function
	if prefix == "!" {
		// FIXME: dev only
		fmt.Printf("\nCommand received:\nPrefix: %s\nCommand: %s\nargs: ",
			prefix, command)
		fmt.Print(args)

		// Increment the number of messages this bot has parsed
		botMessages++

		if f := mapCommand(commandList, command); f != nil {
			f(args, m)
		}
	}
}

// mapCommand takes discord command string and maps it to a function with matching key
// Returns the discordCommandHandler function if there is a match
// Returns nil if there is no matching function found
func mapCommand(m map[string]discordCommandHandler, c string) discordCommandHandler {
	for k, v := range m {
		if strings.ToUpper(k) == strings.ToUpper(c) {
			return v
		}
	}

	return nil
}

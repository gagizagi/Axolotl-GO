package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// discordReadyHandler sets required data after a successful connection
func discordReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	// Sets this bots name in discordCfg.Name
	discordCfg.Name = strings.ToUpper(r.User.Username)

	// Iterate through guilds
	for _, guild := range r.Guilds {
		server := server{}
		server.ID = guild.ID
		server.fetch()

		// TODO: remove eventually
		// Attempt to get anime channel default for guilds that are missing it
		// Temporary measure until all guilds migrate to the new model
		if server.AnimeChannel == "" {
			// Get a list of channels for this guild
			channels, _ := s.GuildChannels(guild.ID)
			// Iterate through channels for this guild
			for _, channel := range channels {
				// If it finds a channel with name 'anime' in its name
				// it will update the db to use this channel until changed by guild admin
				if channel.Name == "anime" {
					server.updateAnimeChannel(channel.ID)
				}
			}
		}

		if server.GuildName == "" {
			server.updateGuildName(guild.Name)
		}
	}

	// Logs successful connection to discord was established
	log.Println("Connected to discord as", discordCfg.Name)
}

// discordLeaveChannelHandler handles leaving channels
func discordLeaveChannelHandler(s *discordgo.Session, c *discordgo.ChannelDelete) {
	guild, _ := s.Guild(c.GuildID)
	server := server{}
	server.ID = guild.ID
	server.fetch()

	if c.ID == server.AnimeChannel {
		server.updateAnimeChannel("")
	}
}

// discordLeaveGuildHandler handles leaving/being kicked from a guild
func discordLeaveGuildHandler(s *discordgo.Session, g *discordgo.GuildDelete) {
	server := server{}
	server.ID = g.ID
	err := server.delete()
	if err != nil {
		log.Println(err)
	}
}

// discordMsgHandler is a handler function for incoming discord messages
func discordMsgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer panicRecovery()

	// Increment the number of messages this bot has read
	botReads++

	// Ignore bots own messages
	if m.Author.ID == discord.State.User.ID {
		return
	}

	var server server

	if !isPrivateMessage(m.Message) {
		_, guild, err := getChannelGuildInfo(m.Message)
		if err != nil {
			panic(err)
		}
		server.ID = guild.ID
		err = server.fetch()
		if err != nil {
			panic(err)
		}
	}

	// Split received message into different parts (space seperated)
	// prefix is the first char of the first part
	// command is the rest of first part (not including first char)
	// args is an array of all parts after first one
	var prefix string
	var command string
	var args []string
	args = strings.Fields(m.Content)
	if len(args) > 0 {
		if len(args[0]) > len(server.Prefix) {
			// Make sure there are no index out of bounds runtime errors
			prefix = args[0][0:len(server.Prefix)]
			command = args[0][len(server.Prefix):]
		}
		args = strings.Fields(m.Content)[1:]
	}

	// If bot needs to parse this message, parse it
	// otherwise ignore it and get out of the function
	if prefix == server.Prefix {
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

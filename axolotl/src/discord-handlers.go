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

//discordMsgHandler is a handler function for incomming discord messages
func discordMsgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	botMessages++

	//Split messages into arguments
	args := strings.Fields(m.Content)
	//Check if author is admin
	boss := m.Author.ID == discordCfg.Boss
	//Check if second argument is this bots name
	botcheck := (len(args) > 1 && strings.ToUpper(args[1]) == discordCfg.Name)
	//Check if message is relevant to the bot
	//i.e message starts with '!' followed by a word
	relevant := relevantRegex.MatchString(m.Content)

	//If message is relevant process it otherwise leave this function
	//TODO: Refactor this into something more readable
	if relevant {
		botResponses++
		switch strings.ToUpper(args[0]) {

		// !HELP
		case "!HELP":
			helpCommand(m.ChannelID)

		// !SUB anime.id
		// anime.id is the id of the anime (string with length of 3 chars)
		case "!SUB":
			subCommand(args, m.Author.ID, m.ChannelID)

		//!UNSUB anime.id
		//anime.id is the id of the anime (string with length of 3 chars)
		case "!UNSUB":
			unsubCommand(args, m.Author.ID, m.ChannelID)

		//!MYSUBS
		//Lists all series this user is subscriberd to in format "SeriesName(id)"
		case "!MYSUBS":
			subs := fmt.Sprintf("<@%s> is subscribed to: ", m.Author.ID)
			animeArray := getAnimeListForUser(m.Author.ID)

			for i, anime := range animeArray {
				if i > 0 {
					subs += ", "
				}
				subs += fmt.Sprintf("%s(%s)", anime.Name, anime.ID)
			}

			s.ChannelMessageSend(m.ChannelID, subs)

		//!UPTIME [string]
		//Can optionally include bots name as second argument
		case "!UPTIME":
			if len(args) > 1 {
				if strings.ToUpper(args[1]) == discordCfg.Name {
					s.ChannelMessageSend(m.ChannelID, "Current uptime is "+getUptime())
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Current uptime is "+getUptime())
			}

		//!P string
		//Sets the 'currently playing' state of the bot
		//will only work for admin of the bot
		case "!P":
			if boss {
				if len(args) > 1 {
					var game string
					for i := 1; i < len(args); i++ {
						game += args[i] + " "
					}
					game = game[:len(game)-1]
					discord.UpdateStatus(0, game)
				} else {
					discord.UpdateStatus(0, "")
				}
			}

		//!INFO [string]
		//Can optionally include bots name as second argument
		//Lists bot usage and general information
		case "!INFO":
			if len(args) > 1 {
				if strings.ToUpper(args[1]) == discordCfg.Name {
					s.ChannelMessageSend(m.ChannelID, getInfo())
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, getInfo())
			}

		//!GUILDS [string]
		//Can optionally include bots name as second argument
		//Lists all the guilds this bot is a part of
		//will only work for admin of the bot
		//FIXME: https://github.com/gagizagi/Axolotl-GO/issues/6
		case "!GUILDS":
			if boss && (botcheck || len(args) == 1) {
				s.ChannelMessageSend(m.ChannelID,
					fmt.Sprintf("Currently in %d guilds: %s",
						len(discordCfg.Guilds), strings.Join(discordCfg.Guilds, ", ")))
			}

		case "!TEST":
			msgChan <- msgObject{
				Message: "hello",
				Channel: m.ChannelID,
			}
		}
	}
}

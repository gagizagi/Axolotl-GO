package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//discordMsgHandler is a handler function for incomming discord messages
func discordMsgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	botMessages++

	/*TODO:
	Add counter for !w command to track number of calls.
	Limit for this OWM api key is 600 calls per 10min AND 50k calls a day.
	*/

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
	//TODO: refactor this into something more readable
	if relevant {
		botResponses++
		switch strings.ToUpper(args[0]) {

		//!HELP [string]
		//Can optionally include bots name as second argument
		case "!HELP":
			if len(args) > 1 {
				if strings.ToUpper(args[1]) == discordCfg.Name {
					s.ChannelMessageSend(m.ChannelID, discordHelp)
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, discordHelp)
			}

		//!SUB anime.id
		//anime.id is the id of the anime (string with length of 3 chars)
		case "!SUB":
			if len(args) >= 2 {
				newAnime := anime{ID: strings.ToLower(args[1])}
				if newAnime.Exists() {
					newAnime.AddSub(m.Author.ID)
					s.ChannelMessageSend(m.ChannelID,
						"Successfully subscribed to "+newAnime.Name)
				} else {
					s.ChannelMessageSend(m.ChannelID, "Invalid ID")
				}
			}

		//!UNSUB anime.id
		//anime.id is the id of the anime (string with length of 3 chars)
		case "!UNSUB":
			if len(args) >= 2 {
				newAnime := anime{ID: strings.ToLower(args[1])}
				if newAnime.Exists() {
					newAnime.RemoveSub(m.Author.ID)
					s.ChannelMessageSend(m.ChannelID,
						"Successfully unsubscribed from "+newAnime.Name)
				} else {
					s.ChannelMessageSend(m.ChannelID, "Invalid ID")
				}
			}

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

		//!W string
		//Looks up weather at the the location string
		case "!W":
			if len(args) > 1 {
				var location string
				for i := 1; i < len(args); i++ {
					location += args[i] + " "
				}
				location = location[:len(location)-1]
				s.ChannelMessageSend(m.ChannelID, getWeather(location))
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
		//FIXME:
		case "!GUILDS":
			if boss && (botcheck || len(args) == 1) {
				s.ChannelMessageSend(m.ChannelID,
					fmt.Sprintf("Currently in %d guilds: %s",
						len(discordCfg.Guilds), strings.Join(discordCfg.Guilds, ", ")))
			}
		}
	}
}

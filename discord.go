package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//discordConfig is a struct containing user configuration for discord connection
type discordConfig struct {
	Boss          string
	Name          string
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
		"!sub <id> (Subscribe to anime)\n" +
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
	discord.Open() //Opens discord connection
}

//discordReadyHandler sets required data after a successful connection
func discordReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	//Sets this bots name in discordCfg.Name
	discordCfg.Name = strings.ToUpper(r.User.Username)

	//Iterates through guilds
	for _, guild := range r.Guilds {
		discordCfg.Guilds = appendUnique(discordCfg.Guilds, guild.Name)
		//Gets a list of channels for this guild
		channels, _ := s.GuildChannels(guild.ID)
		//Iterates through a list of channels for this guild
		for _, channel := range channels {
			//If it finds a channel with name 'anime'
			//it will add it to discordCfg.AnimeChannels array
			//which is used when sending new episode messages in discord
			if channel.Name == "anime" {
				discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, channel.ID)
			}
		}
	}

	//Logs successful connection to discord was established
	log.Println("Connected to discord as", discordCfg.Name)
}

//discordMsgHandler is a handler function for incomming discord messages
func discordMsgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	botMessages++
	//Check if author is admin
	boss := m.Author.ID == discordCfg.Boss
	//Check if message is relevant to the bot
	//i.e message starts with '!' followed by a word
	relevant := relevantRegex.MatchString(m.Content)

	//If message is relevant process it otherwise leave this function
	if relevant {
		botResponses++
		args := strings.Fields(m.Content)
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
		//FIXME
		case "!GUILDS":
			if boss {
				if len(args) > 1 {
					if strings.ToUpper(args[1]) == discordCfg.Name {
						s.ChannelMessageSend(m.ChannelID, strings.Join(discordCfg.Guilds, ","))
					}
				} else {
					log.Println(discordCfg.Guilds)
					s.ChannelMessageSend(m.ChannelID, strings.Join(discordCfg.Guilds, ","))
				}
			}
		}
	}
}

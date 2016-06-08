package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

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
	discordCfg.Guilds = append(discordCfg.Guilds, g.Name)
}

//discordLeaveGuildHandler handles leaving/being kicked from a guild
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
			if boss && (botcheck || len(args) == 1) {
				s.ChannelMessageSend(m.ChannelID,
					fmt.Sprintf("Currently in %d guilds: %s",
						len(discordCfg.Guilds), strings.Join(discordCfg.Guilds, ", ")))
			}
		}
	}
}

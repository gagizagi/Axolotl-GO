package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type discordConfig struct {
	Boss          string
	Name          string
	AnimeChannels []string
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
		"```\n" +
		"Find anime list at http://gazzy.space/anime"
)

//Declare variables
var (
	discord       *discordgo.Session           //Discord client
	discordCfg    *discordConfig               //Discord options
	relevantRegex = regexp.MustCompile(`^!\w`) //Discord regex msg parser
)

//Starts discord client
func discordStart(c *discordConfig) {
	discordCfg = c

	var err error
	discord, err = discordgo.New(c.Token)
	if err != nil {
		log.Fatal("discordConnStart() => New() error:\t", err)
	}

	discord.Debug = c.Debug
	discord.AddHandler(discordMsgHandler)
	discord.AddHandler(discordReadyHandler)
	discord.Open()
}

func discordReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	discordCfg.Name = strings.ToUpper(r.User.Username) //Sets bot name
	//Iterates through guilds
	for _, guild := range r.Guilds {
		//Gets a list of channels for this guild
		channels, _ := s.GuildChannels(guild.ID)
		//Iterates through a list of channels for this guild
		for _, channel := range channels {
			if channel.Name == "anime" {
				//Sets anime channel id
				discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, channel.ID)
			}
		}
	}
	log.Println("Connected to discord as", discordCfg.Name)
}

func appendUnique(slice []string, id string) []string {
	for _, s := range slice {
		if s == id {
			return slice
		}
	}
	return append(slice, id)
}

//discord incomming message handler
func discordMsgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	boss := m.Author.ID == discordCfg.Boss
	relevant := relevantRegex.MatchString(m.Content)

	if relevant {
		args := strings.Fields(m.Content)
		switch strings.ToUpper(args[0]) {
		case "!HELP":
			if len(args) > 1 {
				if strings.ToUpper(args[1]) == discordCfg.Name {
					s.ChannelMessageSend(m.ChannelID, discordHelp)
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, discordHelp)
			}
		case "!SUB":
			if len(args) >= 2 {
				anime := anime{ID: strings.ToLower(args[1])}
				if anime.Exists() {
					anime.AddSub(m.Author.ID)
					s.ChannelMessageSend(m.ChannelID,
						"Successfully subscribed to "+anime.Name)
				} else {
					s.ChannelMessageSend(m.ChannelID, "Invalid ID")
				}
			}
		case "!UNSUB":
			if len(args) >= 2 {
				anime := anime{ID: strings.ToLower(args[1])}
				if anime.Exists() {
					anime.RemoveSub(m.Author.ID)
					s.ChannelMessageSend(m.ChannelID,
						"Successfully unsubscribed from "+anime.Name)
				} else {
					s.ChannelMessageSend(m.ChannelID, "Invalid ID")
				}
			}
		case "!UPTIME":
			if len(args) > 1 {
				if strings.ToUpper(args[1]) == discordCfg.Name {
					s.ChannelMessageSend(m.ChannelID, "Current uptime is "+getUptime())
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Current uptime is "+getUptime())
			}
		case "!W":
			if len(args) > 1 {
				var location string
				for i := 1; i < len(args); i++ {
					location += args[i] + " "
				}
				location = location[:len(location)-1]
				s.ChannelMessageSend(m.ChannelID, getWeather(location))
			}
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
		}
	}
}

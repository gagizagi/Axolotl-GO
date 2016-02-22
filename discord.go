package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type discordConfig struct {
	Boss         string
	Name         string
	ChannelAnime string
	Username     string
	Password     string
	Debug        bool
}

const (
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
		"Find anime list at http://axolotl-422.rhcloud.com/anime"
)

var (
	discordConn   *discordgo.Session           //Discord client
	discordCfg    *discordConfig               //Discord options
	relevantRegex = regexp.MustCompile(`^!\w`) //Discord regex msg parser
)

//Starts discord client
func discordConnStart(c *discordConfig) {
	discordCfg = c

	var err error
	discordConn, err = discordgo.New(c.Username, c.Password)
	if err != nil {
		log.Fatal("discordConnStart() => New() error:\t", err)
	}

	discordConn.Debug = c.Debug
	discordConn.AddHandler(discordMsgHandler)
	discordConn.AddHandler(discordReadyHandler)
	discordConn.Open()
}

func discordReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	discordCfg.Name = strings.ToUpper(r.User.Username) //Sets bot name
	for _, guild := range r.Guilds {
		if guild.Name == "422" {
			for _, channel := range guild.Channels {
				if channel.Name == "anime" {
					discordCfg.ChannelAnime = channel.ID //Sets anime channel id
				}
			}
		}
	}
	log.Println("Connected to discord as", discordCfg.Name)
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
					discordConn.UpdateStatus(0, game)
				} else {
					discordConn.UpdateStatus(0, "")
				}
			}
		}
	}
}

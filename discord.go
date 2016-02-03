package main

import (
	"os"
	"regexp"
	"strings"
	"fmt"
	"github.com/bwmarrin/discordgo"
)


const (
	DISCORD_BOSS string = "110846867473973248"
	DISCORD_HELP string = "***LIST OF BOT COMMANDS***\n"+
	"Fields in [] are optional\n"+
	"Fields in <> are mandatory\n\n"+
	"```"+
	"!help [bot] (Lists all bot commands)\n"+
	"!uptime [bot](Prints current bot uptime)\n"+
	"!sub <id> (Subscribe to anime)\n"+
	"!unsub <id> (Unsubscribe from anime)\n"+
	"!w <location> (Prints current weather)\n"+
	"```\n"+
	"Find anime list at http://axolotl-422.rhcloud.com/anime"
)

var (
	DISCORD_NAME string
	DISCORD_ANIME string
	DISCORD_USERNAME string = os.Getenv("DISCORD_USERNAME")
	DISCORD_PASSWORD string = os.Getenv("DISCORD_PASSWORD")
	relevantRegex *regexp.Regexp = regexp.MustCompile(`^!\w`)
	discord *discordgo.Session//discord session
)

//initializes discord bot
func init() {
	var err error
	discord, err = discordgo.New(DISCORD_USERNAME, DISCORD_PASSWORD)
	if err != nil {
		panic(err)
	}
	discord.Debug = true
	discord.OnMessageCreate = DiscordMsgHandler
	discord.OnReady = DiscordReadyHandler
}

func DiscordReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	DISCORD_NAME = strings.ToUpper(r.User.Username)//Sets bot name
	for _, guild := range r.Guilds {
		if guild.Name == "422" {
			for _, channel := range guild.Channels {
				if channel.Name == "anime" {
					DISCORD_ANIME = channel.ID//Sets anime channel id
				}
			}
		}
	}
	fmt.Println("Connected to discord as", DISCORD_NAME)
}

//discord incomming message handler
func DiscordMsgHandler(s *discordgo.Session, m *discordgo.Message) {
	boss := m.Author.ID == DISCORD_BOSS
	relevant := relevantRegex.MatchString(m.Content)

	if relevant {
		args := strings.Fields(m.Content)
		switch strings.ToUpper(args[0]) {
			case "!HELP":
				if len(args) > 1 {
					if strings.ToUpper(args[1]) == DISCORD_NAME {
						s.ChannelMessageSend(m.ChannelID, DISCORD_HELP)
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, DISCORD_HELP)
				}
			case "!SUB":
				if len(args) >= 2 {
					anime := Anime{Id: strings.ToLower(args[1])}
					if anime.Exists() {
						anime.AddSub(m.Author.ID)
						s.ChannelMessageSend(m.ChannelID, "Successfully subscribed to " + anime.Name)
					} else {
						s.ChannelMessageSend(m.ChannelID, "Invalid ID")
					}
				}
			case "!UNSUB":
				if len(args) >= 2 {
					anime := Anime{Id: strings.ToLower(args[1])}
					if anime.Exists() {
						anime.RemoveSub(m.Author.ID)
						s.ChannelMessageSend(m.ChannelID, "Successfully unsubscribed from " + anime.Name)
					} else {
						s.ChannelMessageSend(m.ChannelID, "Invalid ID")
					}
				}
			case "!UPTIME":
				if len(args) > 1 {
					if strings.ToUpper(args[1]) == DISCORD_NAME {
						s.ChannelMessageSend(m.ChannelID, "Current uptime is " + GetUptime())
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "Current uptime is " + GetUptime())
				}
			case "!W":
				if len(args) > 1 {
					var location string
					for i := 1; i < len(args); i++ {
						location += args[i] + " "
					}
					location = location[:len(location)-1]
					s.ChannelMessageSend(m.ChannelID, GetWeather(location))
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
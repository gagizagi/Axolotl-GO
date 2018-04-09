package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// helpCommand responds to discord text command 'help'
// will respond on discord channel c (same channel as the received message)
// response text is the a discord embed
func helpCommand(args []string, m *discordgo.MessageCreate) {
	desc := "Notifies you of any new anime releases as soon as they are available.\n"
	desc += "Subscribe to any anime you want to receive @mentions for with the `sub ID` command.\n\n"
	desc += "[List of commands](https://github.com/gagizagi/Axolotl-GO#bot-commands)\n"
	desc += "[List of anime](https://axolotl.gazzy.online/)\n\n"

	embed := &discordgo.MessageEmbed{
		Color:       0xB1F971,
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "Axolotl anime bot",
		URL:         "https://github.com/gagizagi/Axolotl-GO",
		Description: desc,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: discordCfg.AvatarURL,
			Text:    "GazZy#5249",
		},
	}

	msgChan <- msgObject{
		Channel: m.ChannelID,
		Message: "HELP EMBED",
		Embed:   embed,
		Author:  m.Author.ID,
	}
}

// subCommand is called when the user tries to subscribe to an anime
// via the 'sub <ID>' text command
// will subscribe the user if the anime exists and return a message
// will return an error message to the user if the anime does not exist
// response is sent on the same channel as the received message
// TODO: add responses when format of the command is wrong
// TODO: make it possible to sub to an array of ids with one command
// or user is already a sub
func subCommand(args []string, m *discordgo.MessageCreate) {
	if len(args) > 0 {
		newAnime := anime{ID: strings.ToLower(args[0])}
		if newAnime.Exists() {
			newAnime.AddSub(m.Author.ID)
			msgChan <- msgObject{
				Channel: m.ChannelID,
				Message: "Successfully subscribed to " + newAnime.Name,
				Author:  m.Author.ID,
			}
		} else {
			msgChan <- msgObject{
				Channel: m.ChannelID,
				Message: "Invalid ID",
				Author:  m.Author.ID,
			}
		}
	}
}

// subCommand is called when the user tries to unsubscribe from an anime
// via the 'unsub <ID>' text command
// will unsubscribe the user if the anime exists and return a message
// will return an error message to the user if the anime does not exist
// response is sent on the same channel as the received message
// TODO: add responses when format of the command is wrong
// TODO: make it possible to unsub from an array of ids with one command
// or user is not a sub
func unsubCommand(args []string, m *discordgo.MessageCreate) {
	if len(args) > 0 {
		newAnime := anime{ID: strings.ToLower(args[0])}
		if newAnime.Exists() {
			newAnime.RemoveSub(m.Author.ID)
			msgChan <- msgObject{
				Channel: m.ChannelID,
				Message: "Successfully unsubscribed from " + newAnime.Name,
				Author:  m.Author.ID,
			}
		} else {
			msgChan <- msgObject{
				Channel: m.ChannelID,
				Message: "Invalid ID",
				Author:  m.Author.ID,
			}
		}
	}
}

// mySubs responds with a list of all the anime the requesting user
// is subscribed to
// response is sent on the same channel as the received message
func mySubs(args []string, m *discordgo.MessageCreate) {
	var subs []string
	animeArray := getAnimeListForUser(m.Author.ID)

	for _, a := range animeArray {
		subs = append(subs, fmt.Sprintf("%s(%s)", a.Name, a.ID))
	}

	result := fmt.Sprintf("<@%s> is subscribed to %d series: %s",
		m.Author.ID, len(animeArray), strings.Join(subs, ", "))

	msgChan <- msgObject{
		Message: result,
		Channel: m.ChannelID,
		Author:  m.Author.ID,
	}
}

// uptime responds with the current uptime of the bot
// response is sent on the same channel as the received message
func uptime(args []string, m *discordgo.MessageCreate) {
	msgChan <- msgObject{
		Message: "Current bot uptime is " + getUptime(),
		Channel: m.ChannelID,
		Author:  m.Author.ID,
	}
}

// setStatus sets the bots status
// admin only command
// sending empty command (!p) will set it to no status
func setStatus(args []string, m *discordgo.MessageCreate) {
	// Only proceed if the message sender is discord admin
	if len(args) > 0 && m.Author.ID == discordCfg.Boss {
		game := strings.Join(args[0:len(args)], " ")
		discord.UpdateStatus(0, game)
	} else {
		discord.UpdateStatus(0, "")
	}
}

// botInfo responds with the different bot statistics
// response is sent on the same channel as the received message
func botInfo(args []string, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Color:     0xB1F971,
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     "Axolotl anime bot",
		URL:       "https://github.com/gagizagi/Axolotl-GO",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Uptime",
				Value:  getUptime(),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Guilds",
				Value:  strconv.Itoa(len(discordCfg.Guilds)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Unique Subscribers",
				Value:  strconv.Itoa(getUniqueSubs()),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Messages Read",
				Value:  strconv.Itoa(botReads),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Messages Parsed",
				Value:  strconv.Itoa(botMessages),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Messages Sent",
				Value:  strconv.Itoa(botResponses),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: discordCfg.AvatarURL,
			Text:    "GazZy#5249",
		},
	}

	msgChan <- msgObject{
		Message: "INFO EMBED",
		Embed:   embed,
		Channel: m.ChannelID,
		Author:  m.Author.ID,
	}
}

// guilds responds with a list of all the guilds this bot is currently in
// admin only command because it is potentially multiple messages long
func guilds(args []string, m *discordgo.MessageCreate) {
	// Only proceed if the message sender is discord admin
	if m.Author.ID != discordCfg.Boss {
		return
	}

	result := fmt.Sprintf("Bot is currently in %d guilds: ",
		len(discordCfg.Guilds))

	for i, guild := range discordCfg.Guilds {
		if i == 0 {
			result += " " + guild
		} else if (len(result) + len(guild) + 2) <= 2000 {
			result += ", " + guild
		} else {
			msgChan <- msgObject{
				Message: result,
				Channel: m.ChannelID,
				Author:  m.Author.ID,
			}
			result = ""
		}
	}

	if result != "" {
		msgChan <- msgObject{
			Message: result,
			Channel: m.ChannelID,
			Author:  m.Author.ID,
		}
	}
}

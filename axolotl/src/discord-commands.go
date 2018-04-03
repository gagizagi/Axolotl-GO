package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HELPMESSAGE string lists all discord chat commands for this bot
// already formatted as a discord string
const HELPMESSAGE string = "***LIST OF BOT COMMANDS***\n" +
	"Fields in [] are optional\n" +
	"Fields in <> are mandatory\n\n" +
	"```" +
	"!help [bot] (Lists all bot commands)\n" +
	"!uptime [bot](Prints current bot uptime)\n" +
	"!sub <id> (Subscribe to anime and get notified when a new episode is released)\n" +
	"!unsub <id> (Unsubscribe from anime)\n" +
	"!mysubs (List all the anime you are subscribed to)\n" +
	"!info [bot] (Prints bot information)" +
	"```\n" +
	"Full anime list at https://axolotl.gazzy.online/ \n" +
	"For issues and suggestions go to https://github.com/gagizagi/Axolotl-GO"

// helpCommand responds to discord text command 'help'
// will respond on discord channel c (same channel as the received message)
// response text is the HELPMESSAGE constant (list of all bot commands)
func helpCommand(args []string, m *discordgo.MessageCreate) {
	msgChan <- msgObject{
		Channel: m.ChannelID,
		Message: HELPMESSAGE,
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
			}
		} else {
			msgChan <- msgObject{
				Channel: m.ChannelID,
				Message: "Invalid ID",
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
			}
		} else {
			msgChan <- msgObject{
				Channel: m.ChannelID,
				Message: "Invalid ID",
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
	}
}

// uptime responds with the current uptime of the bot
// response is sent on the same channel as the received message
func uptime(args []string, m *discordgo.MessageCreate) {
	msgChan <- msgObject{
		Message: getUptime(),
		Channel: m.ChannelID,
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
	result := "```"
	result += fmt.Sprintf("Name: %s\n", discordCfg.Name)
	result += fmt.Sprintf("Uptime: %s\n", getUptime())
	result += fmt.Sprintf("Guilds: %d\n", len(discordCfg.Guilds))
	result += fmt.Sprintf("Anime channels: %d\n", len(discordCfg.AnimeChannels))
	result += fmt.Sprintf("Unique subscribers: %d\n", getUniqueSubs())
	result += fmt.Sprintf("Messages read: %d\n", botReads)
	result += fmt.Sprintf("Messages parsed: %d\n", botMessages)
	result += fmt.Sprintf("Message responses: %d\n", botResponses+1)
	result += "```"

	msgChan <- msgObject{
		Message: result,
		Channel: m.ChannelID,
	}
}

// guilds responds with a list of all the guilds this bot is currently in
// admin only command because it is potentially multiple messages long
// TODO: test this command on production instance of the bot
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
			}
			result = ""
		}
	}

	if result != "" {
		msgChan <- msgObject{
			Message: result,
			Channel: m.ChannelID,
		}
	}
}

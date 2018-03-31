package main

import (
	"fmt"
	"strings"
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
func helpCommand(c string) {
	msgChan <- msgObject{
		Channel: c,
		Message: HELPMESSAGE,
	}
}

// subCommand is called when the user tries to subscribe to an anime
// via the 'sub <ID>' text command
// will subscribe the user if the anime exists and return a message
// will return an error message to the user if the anime does not exist
// response is sent on the same channel as the received message
// TODO: add responses when format of the command is wrong
// or user is already a sub
func subCommand(args []string, authorID string, channelID string) {
	if len(args) >= 2 {
		newAnime := anime{ID: strings.ToLower(args[1])}
		if newAnime.Exists() {
			newAnime.AddSub(authorID)
			msgChan <- msgObject{
				Channel: channelID,
				Message: "Successfully subscribed to " + newAnime.Name,
			}
		} else {
			msgChan <- msgObject{
				Channel: channelID,
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
// or user is not a sub
func unsubCommand(args []string, authorID string, channelID string) {
	if len(args) >= 2 {
		newAnime := anime{ID: strings.ToLower(args[1])}
		if newAnime.Exists() {
			newAnime.RemoveSub(authorID)
			msgChan <- msgObject{
				Channel: channelID,
				Message: "Successfully unsubscribed from " + newAnime.Name,
			}
		} else {
			msgChan <- msgObject{
				Channel: channelID,
				Message: "Invalid ID",
			}
		}
	}
}

// mySubs responds with a list of all the anime the requesting user
// is subscribed to
// response is sent on the same channel as the received message
func mySubs(authorID string, channelID string) {
	var subs []string
	animeArray := getAnimeListForUser(authorID)

	for _, a := range animeArray {
		subs = append(subs, fmt.Sprintf("%s(%s)", a.Name, a.ID))
	}

	result := fmt.Sprintf("<@%s> is subscribed to %d series: %s",
		authorID, len(animeArray), strings.Join(subs, ", "))

	msgChan <- msgObject{
		Message: result,
		Channel: channelID,
	}
}

// uptime responds with the current uptime of the bot
// response is sent on the same channel as the received message
func uptime(channelID string) {
	msgChan <- msgObject{
		Message: getUptime(),
		Channel: channelID,
	}
}

// setStatus sets the bots status
// admin only command
// sending empty command (!p) will set it to no status
func setStatus(args []string) {
	if len(args) > 1 {
		game := strings.Join(args[1:len(args)], " ")
		discord.UpdateStatus(0, game)
	} else {
		discord.UpdateStatus(0, "")
	}
}

// botInfo responds with the different bot statistics
// response is sent on the same channel as the received message
func botInfo(channelID string) {
	result := "```"
	result += fmt.Sprintf("Name: %s\n", discordCfg.Name)
	result += fmt.Sprintf("Uptime: %s\n", getUptime())
	result += fmt.Sprintf("Guilds: %d\n", len(discordCfg.Guilds))
	result += fmt.Sprintf("Anime channels: %d\n", len(discordCfg.AnimeChannels))
	result += fmt.Sprintf("Unique subscribers: %d\n", getUniqueSubs())
	result += fmt.Sprintf("Messages read: %d\n", botMessages)
	result += fmt.Sprintf("Message responses: %d\n", botResponses)
	result += "```"

	msgChan <- msgObject{
		Message: result,
		Channel: channelID,
	}
}

// guilds responds with a list of all the guilds this bot is currently in
// admin only command because it is potentially multiple messages long
// TODO: test this command on production instance of the bot
func guilds(channelID string) {
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
				Channel: channelID,
			}
			result = ""
		}
	}

	if result != "" {
		msgChan <- msgObject{
			Message: result,
			Channel: channelID,
		}
	}
}

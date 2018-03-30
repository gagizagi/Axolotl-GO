package main

import "strings"

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

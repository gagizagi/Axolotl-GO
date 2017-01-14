package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	ircevent "github.com/thoj/go-ircevent"
)

//irc client options struct
type ircConfig struct {
	Server   string
	Channels []string
	Nickname string
	Verbose  bool
	Debug    bool
}

var (
	ircConn *ircevent.Connection //irc connection
	ircCfg  *ircConfig           //irc config
	//regex parse string for newEpisode()
	releaseWatch = regexp.MustCompile(
		`(?i)release:.+\[horriblesubs\] (.+) - ([0-9]{1,3}) \[(1080p|720p|480p)\]`)
)

//Starts the irc client
func ircConnStart(c *ircConfig) {
	ircCfg = c

	ircConn = ircevent.IRC(c.Nickname, c.Nickname)
	ircConn.Debug = c.Debug
	ircConn.VerboseCallbackHandler = c.Verbose
	ircConn.AddCallback("PRIVMSG", ircMsgHandler)
	ircConn.AddCallback("001", ircWelcomeHandler)

	err := ircConn.Connect(c.Server)
	if err != nil {
		log.Fatal("ircConnStart() => Connect() error:\t", err)
	}
	ircConn.Loop()
}

//irc client incomming message handler function
func ircMsgHandler(e *ircevent.Event) {
	//Samples new anime string:
	//Release: [Anime] [HorribleSubs] Ushio to Tora - 22 [480p].mkv
	//Release: [Anime] [HorribleSubs] Hackadoll the Animation - 09 [720p].mkv

	if releaseWatch.MatchString(e.Message()) == true {
		newEpisode(releaseWatch.FindStringSubmatch(e.Message()))
	} else if e.Arguments[0] == ircClient.Nickname && e.Nick != ircClient.Nickname {
		ircConn.Privmsg(e.Nick, "Discord BOT - https://github.com/gagizagi/Axolotl-GO")
		ircConn.Privmsg(e.Nick, "Current uptime is "+getUptime())
	}
}

//irc client onWelcome handler function
func ircWelcomeHandler(e *ircevent.Event) {
	for _, channel := range ircCfg.Channels {
		ircConn.Join(channel)
	}
}

//new episode handler function
//args
// 0 = full message
// 1 = anime name
// 2 = episode (1-3 length integers only)
// 3 = resolution (1080p | 720p | 480p)
func newEpisode(args []string) {
	epnum, _ := strconv.Atoi(args[2])
	episode := anime{Name: args[1], Episode: epnum}
	if episode.Exists() {
		if episode.NewEpisode() {
			episode.UpdateEp()
			if len(episode.Subs) > 0 {
				resultstr := fmt.Sprintf(
					"**New episode of %s released - Episode %d**\n",
					episode.Name, episode.Episode)

				//Add mentions for every subbed user
				for _, person := range episode.Subs {
					resultstr += fmt.Sprintf("<@%s>", person)
				}

				//Add download link
				if episode.Href != "" {
					resultstr += fmt.Sprintf("\nDownload at %s\n", episode.Href)
				} else {
					resultstr += fmt.Sprint("\nDownload at http://horriblesubs.info/\n")
				}

				//Add subscribe ID
				resultstr += fmt.Sprintf(
					"To subscribe to this anime type \"!sub %s\"",
					episode.ID)

				//Send update message to all anime channels
				//TODO make a function
				for _, channel := range discordCfg.AnimeChannels {
					discord.ChannelMessageSend(channel, resultstr)
				}
			}
		}
	} else {
		//Insert series into DB
		episode.Insert()

		//Announce new series to all anime channels
		resultstr := fmt.Sprintf("**New series started: %s - Episode %d**\n",
			episode.Name, episode.Episode)
		resultstr += fmt.Sprintf("To subscribe to this anime type \"!sub %s\"",
			episode.ID)

		//TODO make a function
		for _, channel := range discordCfg.AnimeChannels {
			discord.ChannelMessageSend(channel, resultstr)
		}
	}
}

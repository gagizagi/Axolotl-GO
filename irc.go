package main

import (
	"fmt"
	"strconv"
	"regexp"
	irc "github.com/thoj/go-ircevent"
)

const (
	IRC_NAME string = "Axolotl"
	IRC_SERVER string = "irc.rizon.net:6667"
	IRC_DEBUG bool = false
)

var (
	ircobj *irc.Connection//irc client
	ReleaseWatch *regexp.Regexp = regexp.MustCompile(`(?i)release:.+\[horriblesubs\] (.+) - ([0-9]{1,3}) \[(1080p|720p|480p)\]`)
)

//Initializing of the IRC client
func init() {
	ircobj = irc.IRC(IRC_NAME, IRC_NAME)
	err := ircobj.Connect(IRC_SERVER)
	if err != nil {
		panic(err)
	}
	ircobj.Debug = IRC_DEBUG
	ircobj.AddCallback("PRIVMSG", IrcMsgHandler)
	ircobj.Join("#HORRIBLESUBS")
	ircobj.Join("#422")
}

//irc client incomming message handler function
func IrcMsgHandler(event *irc.Event) {
	//Samples new anime string:
	//Release: [Anime] [HorribleSubs] Ushio to Tora - 22 [480p].mkv
	//Release: [Anime] [HorribleSubs] Hackadoll the Animation - 09 [720p].mkv
	ReleaseWatch := regexp.MustCompile(`(?i)release:.+\[horriblesubs\] (.+) - ([0-9]{1,3}) \[(1080p|720p|480p)\]`)
	
	if ReleaseWatch.MatchString(event.Message()) == true {
		newEpisode(ReleaseWatch.FindStringSubmatch(event.Message()))
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
	episode := Anime{Name:args[1], Episode:epnum}
	if episode.Exists() {
		if episode.NewEpisode() {
			episode.UpdateEp()
			if len(episode.Subs) > 0 {
				resultstr := fmt.Sprintf("**New episode of %s released - Episode %d**\n", episode.Name, episode.Episode)
				for _, person := range episode.Subs {
					resultstr += fmt.Sprintf("<@%s>", person)
				}
				resultstr += fmt.Sprintf("\nDownload at %s\n", episode.Href)
				resultstr += fmt.Sprintf("To subscribe to this anime type \"!sub %s\"", episode.Id)
				discord.ChannelMessageSend(DISCORD_ANIME, resultstr)
			}
		}
	} else {
		episode.Insert()
	}
}
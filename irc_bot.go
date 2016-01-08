package main

import (
	"strconv"
	"regexp"
	irc "github.com/thoj/go-ircevent"
)

var ircobj *irc.Connection

func init() {
	ircobj = irc.IRC("GO-Axolotl-GO", "GO-Axolotl-GO")
	ircobj.Debug = false
	err := ircobj.Connect("irc.rizon.net:6667")
	if err != nil {
		panic(err)
	}
	
	ircobj.Join("#HORRIBLESUBS")
	ircobj.AddCallback("PRIVMSG", IrcMsgHandler)
}

func IrcMsgHandler(event *irc.Event) {
	//Samples new anime string:
	//Release: [Anime] [HorribleSubs] Ushio to Tora - 22 [480p].mkv
	//Release: [Anime] [HorribleSubs] Hackadoll the Animation - 09 [720p].mkv
	ReleaseWatch := regexp.MustCompile(`(?i)release:.+\[horriblesubs\] (.+) - ([0-9]{1,3}) \[(1080p|720p|480p)\]`)
	
	if ReleaseWatch.MatchString(event.Message()) == true {
		newEpisode(ReleaseWatch.FindStringSubmatch(event.Message()))
	}
}

func newEpisode(args []string) {
	//args
	// 0 = full message
	// 1 = anime name
	// 2 = episode (1-3 length integers only)
	// 3 = resolution (1080p | 720p | 480p)
	//epnum, _ := strconv.Atoi(args[2])
	epnum, _ := strconv.Atoi(args[2])
	episode := Anime{Name:args[1], Episode:epnum}
	if episode.Exists() {
		if episode.NewEpisode() {
			episode.UpdateEp()
		}
	} else {
		episode.Insert()
	}
}
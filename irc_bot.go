package main

import (
	"strconv"
	"fmt"
	"regexp"
	"time"
	irc "github.com/thoj/go-ircevent"
	"gopkg.in/mgo.v2/bson"
)

var ircobj *irc.Connection

func init() {
	ircobj = irc.IRC("GO-Axolotl-GO", "GO-Axolotl-GO")
	ircobj.Debug = true
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
		HandleNewAnimeEpisode(ReleaseWatch.FindStringSubmatch(event.Message()))
	}
}

func HandleNewAnimeEpisode(args []string) (error) {
	if n, _ := DBanimeList.Find(bson.M{"name": args[1]}).Count();n == 0 {
		fmt.Println(n)
		newEpisode, _ := strconv.Atoi(args[2])
		newAnime := Anime{
			Name		: args[1],
			Episode		: newEpisode,
			LastUpdate	: time.Now()}
		err := newAnime.Insert()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
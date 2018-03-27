package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	cutoff     time.Time
	titleRegex = regexp.MustCompile(
		`(?i)\[horriblesubs\] (.+) - ([0-9]{1,4}) \[(1080p|720p|480p)\]`)
	rssURL = "http://horriblesubs.info/rss.php?res=sd"
)

func rssReader() {
	defer rssReaderCleanup()

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		panic(fmt.Sprintf("Error trying to parse RSS feed URL: %s - %s", rssURL, err))
	}

	//Iterates through the RSS feed items in reverse order
	for i := len(feed.Items) - 1; i >= 0; i-- {
		//True if the publish time/date of this feed item is after the cutoff time/date
		relevantDate := feed.Items[i].PublishedParsed.After(cutoff)
		//True if the title of this feed item matches the regular expression
		relevantTitle := titleRegex.MatchString(feed.Items[i].Title)

		if relevantTitle && relevantDate {
			regexArray := titleRegex.FindStringSubmatch(feed.Items[i].Title)
			newUpdate(regexArray)
		}
	}
}

func rssReaderCleanup() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

// new episode or series handler function
// args
// 0 = full message
// 1 = anime name
// 2 = episode (1-4 length integers only)
// 3 = resolution (1080p | 720p | 480p)
func newUpdate(args []string) {
	epnum, _ := strconv.Atoi(args[2])
	entry := anime{Name: args[1], Episode: epnum}

	if entry.Exists() && entry.NewEpisode() {
		// If this series already exists in the database
		// but needs to update the episode number
		entry.UpdateEp()

		if len(entry.Subs) > 0 {
			newMessage := fmt.Sprintf(
				"**New episode of %s released - Episode %d**\n",
				entry.Name, entry.Episode)

			// Add mentions for every subbed user
			for _, person := range entry.Subs {
				newMessage += fmt.Sprintf("<@%s>", person)
			}

			// Add downloads link
			if entry.Href != "" {
				newMessage += fmt.Sprintf("\nDownload at %s\n", entry.Href)
			} else {
				newMessage += fmt.Sprint("\nDownload at http://horriblesubs.info/\n")
			}

			// Add subscribe ID
			newMessage += fmt.Sprintf(
				"To subscribe to this anime type \"!sub %s\"",
				entry.ID)

			// Send update message to all anime channels
			// TODO: make a function
			for _, channel := range discordCfg.AnimeChannels {
				discord.ChannelMessageSend(channel, newMessage)
			}
		}
	} else if !entry.Exists() {
		// If this series does not exist in the database yet
		// Insert it into the database
		entry.Insert()

		// Announce new series to all anime channels
		newMessage := fmt.Sprintf("**New series started: %s - Episode %d**\n",
			entry.Name, entry.Episode)
		newMessage += fmt.Sprintf("To subscribe to this anime type \"!sub %s\"",
			entry.ID)

		// TODO: make a function
		for _, channel := range discordCfg.AnimeChannels {
			discord.ChannelMessageSend(channel, newMessage)
		}
	}
}

/*
func timer(interval time.Duration, function func()) {
	function()

	for range time.Tick(interval) {
		function()
	}
}
*/

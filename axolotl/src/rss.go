package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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
	defer panicRecovery()

	// Parse the RSS URL
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		panic(fmt.Sprintf("Error trying to parse RSS feed URL: %s - %s", rssURL, err))
	}

	// Array of updated objects in this function interval
	var titleUpdates []string

	// Iterate through the RSS feed items in reverse order
	for i := len(feed.Items) - 1; i >= 0; i-- {
		// True if the publish time/date of this feed item is after the cutoff time/date
		relevantDate := feed.Items[i].PublishedParsed.After(cutoff)
		// True if the title of this feed item matches the regular expression
		relevantTitle := titleRegex.MatchString(feed.Items[i].Title)

		// If there is a new RSS entry published since last update date
		// handle it with newUpdate() function
		// if there is a new entry that doesn't match the regex log it
		// update cutoff time with the latest update time
		if relevantTitle && relevantDate {
			regexArray := titleRegex.FindStringSubmatch(feed.Items[i].Title)
			ok := newUpdate(regexArray)

			if ok {
				titleUpdates = appendUnique(titleUpdates, regexArray[1])
			}
		} else if relevantDate && !relevantTitle {
			log.Println("Error trying to match this RSS feed item title:", feed.Items[i].Title)
		}

		if relevantDate {
			// Update the cutoff time for old feed entries
			cutoff = *feed.Items[i].PublishedParsed
		}
	}

	if len(titleUpdates) > 0 {
		log.Printf("Updated %d anime entries: %s", len(titleUpdates), strings.Join(titleUpdates, ", "))
	}
}

// new episode or series handler function
// args
// 0 = full message
// 1 = anime name
// 2 = episode (1-4 length integers only)
// 3 = resolution (1080p | 720p | 480p)
func newUpdate(args []string) bool {
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

		return true
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

		return true
	}

	return false
}

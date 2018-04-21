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
		if entry.Href == "" {
			entry.Href = "http://horriblesubs.info"
		}

		// Discord message format for new episode release notification
		messageBuilder := "**New episode of %s released - Episode %d**\n"
		messageBuilder += "Download at %s\n"
		messageBuilder += "To subscribe to this anime type `%ssub %s` - Total users subscribed `%d`\n"

		// Range over all the guilds bot is in
		gg := discord.State.Guilds
		for _, g := range gg {
			// Fetch this guilds settings from the db
			guild := server{}
			guild.ID = g.ID
			guild.fetch()

			if guild.AnimeChannel == "" {
				continue
			}

			// Handle sending the message based on the guilds notification mode
			switch guild.Mode {
			// Won't send any notifications
			case "ignore":
				continue

			// Will always send notifications, but will never include mentions for subscribers
			case "always":
				msgChan <- msgObject{
					Channel: guild.AnimeChannel,
					Message: fmt.Sprintf(messageBuilder, entry.Name, entry.Episode, entry.Href, guild.Prefix, entry.ID, len(entry.Subs)),
				}

			// Will always send notifications, and also include mentions for subscribers
			case "alwaysplus":
				messageBuilderPlus := messageBuilder + "Subscribers in this guild: %s"
				mentions := ""

				mm := g.Members
				for _, m := range mm {
					if contains(entry.Subs, m.User.ID) {
						mentions += fmt.Sprintf("<@%s> ", m.User.ID)
					}
				}

				msgChan <- msgObject{
					Channel: guild.AnimeChannel,
					Message: fmt.Sprintf(messageBuilderPlus, entry.Name, entry.Episode, entry.Href, guild.Prefix, entry.ID, len(entry.Subs), mentions),
				}

			// subonly mode and default will send notifications only if someone on the server is subscribed to the anime
			default:
				if len(entry.Subs) > 0 {
					relevantGuild := false
					messageBuilderPlus := messageBuilder + "Subscribers in this guild: %s"
					mentions := ""

					mm := g.Members
					for _, m := range mm {
						if contains(entry.Subs, m.User.ID) {
							mentions += fmt.Sprintf("<@%s> ", m.User.ID)
							relevantGuild = true
						}
					}

					if relevantGuild {
						msgChan <- msgObject{
							Channel: guild.AnimeChannel,
							Message: fmt.Sprintf(messageBuilderPlus, entry.Name, entry.Episode, entry.Href, guild.Prefix, entry.ID, len(entry.Subs), mentions),
						}
					}
				}
			}
		}

		return true
	} else if !entry.Exists() {
		// If this series does not exist in the database yet
		// Insert it into the database
		entry.Insert()

		// Discord message format for new episode release notification
		messageBuilder := "**New series started: %s - Episode %d**\n"
		messageBuilder += "Download at http://horriblesubs.info\n"
		messageBuilder += "To subscribe to this anime type `%ssub %s`\n"

		// Range over all the guilds bot is in
		gg := discord.State.Guilds
		for _, g := range gg {
			// Fetch this guilds settings from the db
			guild := server{}
			guild.ID = g.ID
			guild.fetch()

			if guild.AnimeChannel == "" {
				continue
			}
			if guild.Mode == "ignore" {
				continue
			}

			msgChan <- msgObject{
				Channel: guild.AnimeChannel,
				Message: fmt.Sprintf(messageBuilder, entry.Name, entry.Episode, guild.Prefix, entry.ID),
			}
		}

		return true
	}

	return false
}

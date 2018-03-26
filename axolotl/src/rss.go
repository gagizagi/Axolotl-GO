package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	cutoff *time.Time
	rssURL = "http://horriblesubs.info/rss.php?res=sd"
	regex  = regexp.MustCompile(
		`(?i)\[horriblesubs\] (.+) - ([0-9]{1,4}) \[(1080p|720p|480p)\]`)
)

func rssReader() {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range feed.Items {
		if item.Title == "[HorribleSubs] Citrus - 12 [480p].mkv" {
			cutoff = item.PublishedParsed
		}
	}

	for _, item := range feed.Items {
		if item.PublishedParsed.After(*cutoff) {
			strArr := regex.FindStringSubmatch(item.Title)

			for i, strSub := range strArr {
				fmt.Println(i, strSub)
			}
		}
	}
}

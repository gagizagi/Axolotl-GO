package main

import (
	"time"
)

func main() {

	go IrcJoinChannels()
	go Maintain_anime_list(3 * time.Hour)
		
	Web_server()//Last
}
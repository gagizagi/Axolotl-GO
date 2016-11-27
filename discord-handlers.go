package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//discordReadyHandler sets required data after a successful connection
func discordReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	//Sets this bots name in discordCfg.Name
	discordCfg.Name = strings.ToUpper(r.User.Username)

	//Iterates through guilds
	for _, guild := range r.Guilds {
		//Gets a list of channels for this guild
		channels, _ := s.GuildChannels(guild.ID)
		//Iterates through a list of channels for this guild
		for _, channel := range channels {
			//If it finds a channel with name 'anime'
			//it will add it to discordCfg.AnimeChannels array
			//which is used when sending new episode messages in discord
			if channel.Name == discordCfg.AnimeChannel {
				discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, channel.ID)
			}
		}
	}

	//Logs successful connection to discord was established
	log.Println("Connected to discord as", discordCfg.Name)
}

//discordNewChannelHandler handles joining new channels
func discordNewChannelHandler(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Name == discordCfg.AnimeChannel {
		discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, c.ID)
	}
}

//discordLeaveChannelHandler handles leaving channels
func discordLeaveChannelHandler(s *discordgo.Session, c *discordgo.ChannelDelete) {
	if c.Name == discordCfg.AnimeChannel {
		discordCfg.AnimeChannels = removeItem(discordCfg.AnimeChannels, c.ID)
	}
}

//discordChannelUpdateHandler handles channels being updated
func discordChannelUpdateHandler(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	if c.Name == discordCfg.AnimeChannel {
		discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, c.ID)
	} else {
		discordCfg.AnimeChannels = removeItem(discordCfg.AnimeChannels, c.ID)
	}
}

//discordNewGuildHandler handlers joining a guild
func discordNewGuildHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	discordCfg.Guilds = appendUnique(discordCfg.Guilds, g.Name)

	for _, c := range g.Channels {
		if c.Name == discordCfg.AnimeChannel {
			discordCfg.AnimeChannels = appendUnique(discordCfg.AnimeChannels, c.ID)
		}
	}
}

//discordLeaveGuildHandler handles leaving/being kicked from a guild
//FIXME remove this anime channel from config when kicked
func discordLeaveGuildHandler(s *discordgo.Session, g *discordgo.GuildDelete) {
	discordCfg.Guilds = removeItem(discordCfg.Guilds, g.Name)
}

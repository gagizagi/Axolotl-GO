# Axolotl-GO

![alt text](https://cdn.discordapp.com/avatars/185177851799011329/70c6e3e9a5763bed9fdf36510e81273a.jpg "Axolotl so moe")

### Introduction

Update 14.4.2018: Almost done with the big update. All that's needed now is some testing and more polish. Hopefully live 
within the next week.

Discord and IRC bot focused primarily on providing current season airing anime episode updates as discord @mentions.
Bot has 99.9% uptime and is currently in daily use on over 200 servers.

### Navigation

* [TOP](#axolotl-go)
* [Introduction](#introduction)
* [Navigation](#navigation)
* [Usage](#usage)
* [Bot commands](#bot-commands)
* [Contact](#contact)
* [Credits](#credits)

### Usage

###### Adding the bot to your guild:
* Authorize it to join your guild [here](https://discordapp.com/oauth2/authorize?client_id=185177389163085824&scope=bot&permissions=19456) NOTE: You need to have sufficient permissions in guild and be logged in to authorize the bot.
* Use the `!notifyhere` command in the channel you want to use for anime updates.
* Optionally change the prefix for bot commands from '!' to anything you want with the `!prefix newPrefix` command.
* Optionally change the notification mode of the bot with the `!mode` command. [Read more under commands](#bot-commands)
* Feel free to [contact me](#contact) if you have any issues getting the above to work, if you have any cool suggestion or found any bugs.

### Bot Commands

|Command|Description|Example|Extra|
---|---|---|---
help|Returns a list of all available commands in Discord chat.|!help
uptime|Returns current uptime of the bot.|!uptime
sub *id*|Subscribe to the anime series to receive @mentions whenever a new episode is released.|!sub yls|Get full list of series [here](https://axolotl.gazzy.online/)
unsub *id*|Unsubscribe from the anime series|!unsub 6aj|Get full list of series [here](https://axolotl.gazzy.online/)
mysubs|Lists all the anime you are subscribed to|!mysubs
info|Returns information about the bot|!info
prefix *new prefix*|Sets a new prefix for this bot on this server|!prefix ?|Server owner only
notifyhere|Set this channel for anime notifications|!notifyhere|Server owner only
mode *mode*|Set the notification mode of the bot|!mode always|`ignore` - Bot won't send any anime notifications<br>`subonly`(default) - Bot will only send notifications if someone from the server is subscribed to the anime<br>`always` - Bot will always send notifications, but won't mention anyone<br>`alwaysplus` - Bot will always send notifications and mention everyone subbed to the anime
status *string*|Sets the game bot is "playing"|!status minecraft|Admin only
guilds|Returns a list of all the guilds this bot is in|!guilds|Admin only



### Contact

PM me on discord @GazZy#5249

or

Visit us in our Discord guild at [422 Discord](https://discord.gg/0cO53vq26OlEojNi)

### Credits

###### Anime updates - HorribleSubs
* <http://horriblesubs.info/>
* \#HorribleSubs on irc.rizon.net

###### Discord free VoIP - <https://discordapp.com/>

###### Discord GO package - [discordgo](https://github.com/bwmarrin/discordgo#discordgo-) by [bwmarrin](https://github.com/bwmarrin)
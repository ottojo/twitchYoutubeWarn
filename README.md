# Twitch Chat Bot

> Twitch Chat Bot that warns if troll youtube links are found

## Authors
 
 - [Jonas Otto](https://github.com/ottojo)

## Setup

``` golang
# Install dependencies
go get github.com/gempir/go-twitch-irc
go get google.golang.org/api/option
go get google.golang.org/api/youtube/v3

# Start project
go build
./main
```

## Add to Discord

 1. [https://discordapp.com/developers/applications](https://discordapp.com/developers/applications)
 2. Select "**New Application**"
 3. Set a name to your bot
 4. Copy **CLIENT ID**
 5. Select "**Bot**"
 6. Select "**Add Bot**"
 7. Select "**Yes, do it!**"
 8. Add ClienteID to url: https://discordapp.com/api/oauth2/authorize?client_id=**CLIENTID**&scope=bot
 9. Open url in browser
 10. Add bot to your server
 11. Select Authorize
 12. Done!! 

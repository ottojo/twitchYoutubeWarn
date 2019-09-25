package main

import (
	"context"
	"flag"
	"github.com/gempir/go-twitch-irc"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"strings"
	"time"
)

func init() {
}

const commandPrefix = "!warn"

var interceptPrefixes = [...]string{"!sr", "!songrequest"}
var botOauthString string
var botUsername string
var youtubeApiKey string

var twitchChatClient *twitch.Client
var youtubeService *youtube.Service

func init() {
	flag.StringVar(&botOauthString, "oauth", "oauth:123123123", "oAuth String to authenticate to twitch chat")
	flag.StringVar(&botUsername, "botname", "ingressres_bot", "Name of bot account")
	flag.StringVar(&youtubeApiKey, "youtubeApiKey", "", "Youtube API key")
	flag.Parse()

	twitchChatClient = twitch.NewClient(botUsername, botOauthString)
	ctx := context.Background()
	var err error
	youtubeService, err = youtube.NewService(ctx, option.WithScopes(youtube.YoutubeReadonlyScope), option.WithAPIKey(youtubeApiKey))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	twitchChatClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		receiveTime := time.Now()
		log.Printf("Message | in | %s | %s | \"%s\"\n", message.Channel, message.User.DisplayName, message.Message)
		if strings.HasPrefix(message.Message, commandPrefix) {
			handleCommand(message)
		}

		for _, interceptedPrefix := range interceptPrefixes {
			if strings.HasPrefix(message.Message, interceptedPrefix) {
				intercept(message, interceptedPrefix)
			}
		}
		log.Printf("Message | Done in %04dms | %s | %s | \"%s\"\n", time.Now().Sub(receiveTime).Milliseconds(),
			message.Channel, message.User.DisplayName, message.Message)
	})

	twitchChatClient.Join(botUsername)

	err := twitchChatClient.Connect()
	if err != nil {
		panic(err)
	}

	select {}
}

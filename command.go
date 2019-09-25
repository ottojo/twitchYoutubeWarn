package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc"
	"log"
	"strings"
)

func handleCommand(message twitch.PrivateMessage) {
	log.Printf("Command | %s | %s | \"%s\"\n", message.Channel, message.User.DisplayName, message.Message)
	if strings.EqualFold(message.Channel, botUsername) {
		// Message in chat of bot user
		words := strings.Split(message.Message, " ")
		if len(words) < 2 || !strings.EqualFold(words[0], commandPrefix) {
			sendUsage(message.Channel)
			return
		}

		switch words[1] {
		case "join":
			log.Printf("Command | join | %s \n", message.User)
			twitchChatClient.Join(message.User.Name)
			twitchChatClient.Say(botUsername, "Joining "+message.User.Name)
			break
		case "leave":
			log.Printf("Command | leave | %s \n", message.User)
			twitchChatClient.Depart(message.User.Name)
			twitchChatClient.Say(botUsername, "Leaving "+message.User.Name)
			break
		}
	} else {
		// Message in different channel
		// TODO commands to add to ban list, add words to ban list
	}
}

func sendUsage(channel string) {
	twitchChatClient.Say(channel, fmt.Sprintf(
		"Usage: \"%[1]s join\", \"%[1]s leave\", for more configuration type \"%[1]s\" in your own chat.",
		commandPrefix))
}

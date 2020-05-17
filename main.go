package main

import (
	"./cmd"
	"./config"

	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
/*
TODO:
	- add writeup
	- Add new roles when connecting
	- remove user from users.json when he leaves discord (it breaks user.mention())
	- modify PrintBoxInfo to get box name as parameter, so we can print it in the Shoutbox function

	- see https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds 
	
	- add challs in progress
	- create just one write .json function and all goroutines that send to channel (see end of verifyCommand)
	- add roles 
	- add retired challs
*/	

	// Discord Bot
	bot, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		panic(err.Error())
	}

	// Register handlers
	bot.AddHandler(cmd.Ready)
	bot.AddHandler(cmd.CommandHandler)
	bot.AddHandler(cmd.ReactionsHandler)
	
	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	defer bot.Close()
	if err != nil {
		fmt.Println("Could not connect to discord", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Discord bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	
	// Cleanly close down the Discord session.
	fmt.Println("Closing connection")
	bot.Close()
}


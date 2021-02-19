package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func main() {

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Missing token environment variable")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(sendMessage)
	dg.AddHandler(commandHandle)
	//dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// WIP
}

func commandHandle(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignores messages by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Checks if the message has a command prefix
	if !strings.HasPrefix(m.Content, "!") {
		return
	}
	args := strings.Split(m.Content, "!")[1:]

	switch args[0] {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "help":
		s.ChannelMessageSend(m.ChannelID, "WIP")
	case "source":
		s.ChannelMessageSend(m.ChannelID, "Źródło dla bota można znaleźć na https://github.com/BenekDoesHorses/snowfall-discord")
	default:
		s.ChannelMessageSend(m.ChannelID, "Nieznana komenda\nWpisz !help, aby uzyskać liste poprawnych komend!")
	}

	fmt.Println(args)

}

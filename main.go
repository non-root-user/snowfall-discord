package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/streadway/amqp"
)

func main() {

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Missing token environment variable")
	}
	//URL formated such as "amqp://guest:guest@server.localhost:5672/"
	rabbit := os.Getenv("RABBIT_URL")
	if token == "" {
		panic("Missing rabbit url environment variable")
	}

	dg, err := discordgo.New("Bot " + token)
	failOnError(err, "error creating Discord session")

	conn, err := amqp.Dial(rabbit)
	failOnError(err, "error connecting to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"SendDiscordValidationMessage", // name
		false,                          // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

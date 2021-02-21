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

var resendMem = make(map[string]string)

func main() {

	token := os.Getenv("BOT_TOKEN")
	guild := os.Getenv("GUILD_ID")
	channel := os.Getenv("CHANNEL_ID")
	//URL formated such as "amqp://guest:guest@server.localhost:5672/"
	rabbit := os.Getenv("RABBIT_URL")
	for k, v := range map[string]string{
		"token": token, "guild": guild, "channel": channel, "rabbit": rabbit} {
		if v == "" {
			panic("Missing " + k + " environment variable.")
		}
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
			userID := findUsersID(dg, guild, strings.Split(string(d.Body), ":")[0])

			resendMem[userID] = strings.Split(string(d.Body), ":")[1]

			directMessage, err := dg.UserChannelCreate(userID)
			if err != nil {
				fmt.Println("Problem with sending a message, or the recipient does not exist")
				continue
			}

			dg.ChannelMessageSend(directMessage.ID, "Beep boop. Twój magiczny kod aktywacyjny to: `` "+strings.Split(string(d.Body), ":")[1]+" ``")
			dg.ChannelMessageSend(directMessage.ID, "Jeśli ta wiadomość nie powinna tutaj dotrzeć, koniecznie powiadom organizatorów wydarzenia!!!")
			dg.ChannelMessageSend(channel, "<@"+userID+"> Wysłałem ci wiadomość aktywacyjną!\nJeśli wiadomość nie została przez Ciebie otrzymana, włącz prywatne wiadomości na discordzie i użyj komendy ``!resend``")
		}
	}()

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
	case "resend":
		if val, ok := resendMem[m.Author.ID]; ok {
			directMessage, err := s.UserChannelCreate(m.Author.ID)
			if err != nil {
				fmt.Println("Problem with sending a message, or the recipient does not exist")
				return
			}
			if m.ChannelID != directMessage.ID {
				s.ChannelMessageSend(m.ChannelID, "psst, nie chce tutaj zdradzać sekretów, napisz mi prywatną wiadomość :wink:")
				return
			}
			s.ChannelMessageSend(directMessage.ID, "Twój kod aktywacyjny to: `` "+val+" ``, miłej zabawy!")
			delete(resendMem, m.Author.ID)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Wybacz, nie mam twojego kodu, na pewno się zarejestrowałeś?")
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Nieznana komenda\nWpisz !help, aby uzyskać liste poprawnych komend!")
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func findUsersID(s *discordgo.Session, guildID string, username string) string {
	guildMembers, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, m := range guildMembers {
		if m.User.Username+"#"+m.User.Discriminator == username {
			return m.User.ID
		}
	}
	return ""
}

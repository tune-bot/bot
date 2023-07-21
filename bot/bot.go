package bot

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/tune-bot/database"
	"github.com/tune-bot/discord/command"
	"github.com/tune-bot/discord/data"
)

func Start(session *discordgo.Session) {
	// Register the onMessage callback handler, only receive server messages
	session.Identify.Intents = discordgo.IntentsGuildMessages
	session.AddHandler(onMessage)

	// Open a connection to Discord and begin listening
	err := session.Open()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Started Discord session")

	// Open a connection to the database
	if err = database.Connect(); err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Connected to database")

	// Wait until interrupt signal is received, then exit
	log.Printf("%s is running...\n(Ctrl+C to kill)\n", data.TITLE)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Ending Discord session")
	database.Disconnect()
	session.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from  Tunebot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Parse a command
	if strings.HasPrefix(m.Content, data.CMD_PREFIX) {
		tokens := strings.Split(m.Content[strings.LastIndex(m.Content, data.CMD_PREFIX)+1:], " ")
		cmd := tokens[0]
		args := tokens[1:]

		if callback, ok := command.Commands[cmd]; ok {
			s.ChannelMessageSend(m.ChannelID, callback(args))
		}
	}
}

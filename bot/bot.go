package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	command "github.com/tune-bot/discord/command"
	data "github.com/tune-bot/discord/data"
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

	// Wait until interrupt signal is received, then exit
	fmt.Printf("%s is running...\n(Ctrl+C to kill)\n", data.TITLE)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	session.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from  Tunebot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Parse a command
	if strings.HasPrefix(m.Content, data.CMD_PREFIX) {
		cmd := m.Content[strings.LastIndex(m.Content, data.CMD_PREFIX)+1:]
		args := strings.Split(cmd, " ")

		if callback, ok := command.Commands[args[0]]; ok {
			result := "Failure!"

			if callback(m.ChannelID, args[1:]) {
				result = "Success!"
			}

			s.ChannelMessageSend(m.ChannelID, result)
		}
	}
}

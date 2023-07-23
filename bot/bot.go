package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"

	"github.com/tune-bot/database"
	"github.com/tune-bot/discord/data"
)

// Wrap a Discord session and other bot data
type Bot struct {
	*discordgo.Session

	running bool
}

// The bot singleton
var Tunebot *Bot

// Initialize the bot
func init() {
	var session *discordgo.Session
	var err error

	if session, err = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_TOKEN"))); err != nil {
		log.Fatal(err)
		return
	}

	Tunebot = &Bot{Session: session}
}

func (b *Bot) Start() {
	// Register the onMessage callback handler, only receive server messages
	b.Identify.Intents = discordgo.IntentsGuildMessages
	b.AddHandler(onMessage)

	// Open a connection to Discord and begin listening
	err := b.Open()
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

	// Wait until the bot is told to close
	var wg sync.WaitGroup
	wg.Add(1)
	b.running = true

	go func() {
		defer func() {
			log.Println("Ending Discord session")
			database.Disconnect()
			b.Close()
			wg.Done()
		}()

		log.Printf("%s ready!\n", data.TITLE)
		for b.running {
		}
	}()

	wg.Wait()
}

func (b *Bot) Stop() {
	b.running = false
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

		if callback, ok := commands[cmd]; ok {
			s.ChannelMessageSend(m.ChannelID, callback(Tunebot, m, args))
		}
	}
}

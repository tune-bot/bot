package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/tune-bot/core"
)

// Wrap a Discord session and other bot data
type Bot struct {
	*discordgo.Session

	running bool
}

// The bot singleton
var Tunebot *Bot

// A Reaction is a function and a set of args to pass to that function
// that runs after a react response occurs on a bot message
type Reaction struct {
	action reactionCallback
	args   []string
}

// Store message IDs that are awaiting reaction response
// Messages map to an action that must be run after receiving a response
var pendingResponse = make(map[string]*Reaction)

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

func (b *Bot) start() {
	// Listen for message create and react events
	b.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions
	b.AddHandler(onChannelMessage)
	b.AddHandler(onMessageReact)

	// Open a connection to Discord and begin listening
	err := b.Open()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Started Discord session")

	// Open a connection to the database
	if err = core.Connect(); err != nil {
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
			core.Disconnect()
			b.Close()
			wg.Done()
		}()

		log.Printf("%s ready!\n", Title)
		for b.running {
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

func (b *Bot) stop() {
	b.running = false
}

func onChannelMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from  Tunebot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, CmdPrefix) {
		// Parse a command
		tokens := strings.Split(m.Content[strings.LastIndex(m.Content, CmdPrefix)+1:], " ")
		cmd := tokens[0]
		args := tokens[1:]

		if cmdCallback, ok := commands[cmd]; ok {
			output, rxnCallback := cmdCallback(Tunebot, m, args)

			// Send the response, store this message's ID if it requires a reaction response
			if msg, err := s.ChannelMessageSend(m.ChannelID, output); err == nil && rxnCallback != nil {
				reaction := &Reaction{action: rxnCallback}
				reaction.args = append(reaction.args, args...)
				pendingResponse[msg.ID] = reaction
			}
		}
	}
}

func onMessageReact(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// If this message is pending a response, run the callback on this response
	if rxnCallback, ok := pendingResponse[r.MessageID]; ok {
		reaction := rxnCallback.action
		s.ChannelMessageSend(r.ChannelID, reaction(Tunebot, r, rxnCallback.args))
	}
}

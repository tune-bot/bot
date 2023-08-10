package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"

	core "github.com/tune-bot/core/src"
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

// Store usernames whose next message will be the answer to a question
// and map them to callbacks that will be executed on the answer
var questions = make(map[string]answerCallback)

// Initialize the bot
func init() {
	var session *discordgo.Session
	var err error

	if session, err = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_TOKEN"))); err != nil {
		core.PrintError(err.Error())
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
		core.PrintError(err.Error())
		return
	}
	core.PrintSuccess("Started Discord session")

	// Open a connection to the database
	if err = core.Connect(); err != nil {
		core.PrintError(err.Error())
		return
	}
	core.PrintSuccess("Connected to database")

	// Wait until the bot is told to close
	var wg sync.WaitGroup
	wg.Add(1)
	b.running = true

	go func() {
		defer func() {
			core.PrintSuccess("Ending Discord session")
			core.Disconnect()
			b.Close()
			wg.Done()
		}()

		core.PrintSuccess(fmt.Sprintf("%s ready!", Title))
		for b.running {
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

func (b *Bot) stop() {
	b.running = false
}

func (b *Bot) askQuestion(c *discordgo.Channel, recipient, prompt string, callback answerCallback, emojiOption string, emojiCallback reactionCallback) {
	if msg, err := b.ChannelMessageSend(c.ID, prompt); err == nil {
		b.MessageReactionAdd(c.ID, msg.ID, emojiOption)
		questions[recipient] = callback
		pendingResponse[msg.ID] = &Reaction{action: emojiCallback}
	}
}

func onChannelMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Process messages from Tunebot separately
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Process commands
	if strings.HasPrefix(m.Content, CmdPrefix) {
		tokens := strings.Split(m.Content[1:], " ")
		cmd := tokens[0]
		args := tokens[1:]

		if cmdCallback, ok := commands[cmd]; ok {
			output, emojiOptions, reaction := cmdCallback(Tunebot, m, args)

			// Send the response, apply and postActions, and store this message's ID if it requires a reaction response
			if msg, err := s.ChannelMessageSend(m.ChannelID, output); err == nil {
				if reaction != nil {
					reaction := &Reaction{action: reaction}
					reaction.args = append(reaction.args, args...)
					pendingResponse[msg.ID] = reaction
				}

				for _, emoji := range emojiOptions {
					Tunebot.MessageReactionAdd(msg.ChannelID, msg.ID, emoji)
				}
			}
		}

		// Non-command messages
	} else {
		// Check pending message actions
		if user, err := getAccount(m.Author.Username); err == nil {
			if action, found := playlistActions[user.Id]; found {
				if msg, ok := action(Tunebot, m.Message); ok {
					s.ChannelMessageSend(m.ChannelID, msg)
					delete(playlistActions, user.Id)
				}
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Check pending questions
		if _, ok := questions[m.Author.Username]; ok {
			if user, err := getAccount(m.Author.Username); err == nil {
				response := questions[m.Author.Username](user, m.Content)
				s.ChannelMessageSend(m.ChannelID, response)
			} else {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
		}
	}

	// After each message, if there is a pending song and pending playlist loaded, fire off the add function
	if pendingAddPlaylist != nil && pendingAddSong != nil {
		AddPendingSong(Tunebot, m.ChannelID)
	}
}

func onMessageReact(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// If this message is pending a response, run the callback on this response
	if rxnCallback, ok := pendingResponse[r.MessageID]; ok {
		reaction := rxnCallback.action
		response, ok := reaction(Tunebot, r, rxnCallback.args)
		s.ChannelMessageSend(r.ChannelID, response)

		if ok {
			delete(pendingResponse, r.MessageID)
		}
	}
}

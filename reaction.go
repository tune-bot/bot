package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/tune-bot/core"
)

// Callbacks for reactions on bot messages
// A reaction callback takes a reference to the bot and the reaction event
// and returns a response to send to the server
type reactionCallback func(b *Bot, r *discordgo.MessageReactionAdd, args []string) string

func browseSourceCallback(b *Bot, r *discordgo.MessageReactionAdd, _ []string) string {
	if dm, err := b.UserChannelCreate(r.UserID); err == nil {
		msg := fmt.Sprintf("Browse the source code for %s v%s here:\nhttps://github.com/tune-bot/discord", Title, version())
		b.ChannelMessageSend(dm.ID, msg)
	}
	return ""
}

func linkUserCallback(_ *Bot, r *discordgo.MessageReactionAdd, args []string) string {
	var userId string
	var err error

	if userId, err = core.FindUser(args[0]); err != nil {
		return err.Error()
	}

	discordUser := core.Discord{
		Name:   r.Member.User.Username,
		UserId: userId,
	}

	if err = discordUser.Link(); err != nil {
		return err.Error()
	}

	delete(pendingResponse, r.MessageID)
	return fmt.Sprintf("Linked Discord user %s to %s account: %s!", discordUser.Name, Title, args[0])
}

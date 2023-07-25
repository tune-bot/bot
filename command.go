package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/tune-bot/database"
)

// A command callback takes a reference to the bot singleton and
// the list of args it was passed and returns a reply to send the server
// as well as a reaction callback if the message is awaiting a response
type commandCallback func(b *Bot, m *discordgo.MessageCreate, args []string) (string, reactionCallback)

var commands = map[string]commandCallback{
	"help": cmdHelp,
	"kill": cmdKill,
	"link": cmdLink,
}

func cmdHelp(_ *Bot, _ *discordgo.MessageCreate, _ []string) (string, reactionCallback) {
	rsp := fmt.Sprintf(":notes: Hello, I'm %s!\n", Title)
	rsp += fmt.Sprintf(":notes: To get started, create an account on the %s app.\n", Title)
	rsp += fmt.Sprintf(":notes: Then, link your Discord account using the following command: `%slink <%s username>`\n", CmdPrefix, Title)
	rsp += fmt.Sprintf("\n%s v%s\nReact to this message to browse the source code!", Title, version())
	return rsp, browseSourceCallback
}

func cmdKill(b *Bot, _ *discordgo.MessageCreate, _ []string) (string, reactionCallback) {
	b.stop()
	return "", nil
}

func cmdLink(b *Bot, m *discordgo.MessageCreate, args []string) (string, reactionCallback) {
	var err error
	var rsp string

	if len(args) < 1 {
		rsp = fmt.Sprintf("No %s username supplied", Title)
		return rsp, nil
	}

	discordUser := database.Discord{Name: m.Author.Username}
	tunebotUser := database.User{Username: args[0]}

	// Check if the Tunebot user exists before trying to link a Discord account to it
	if _, err := database.FindUser(tunebotUser.Username); err != nil {
		return err.Error(), nil
	}

	// Check if this Discord account is already linked
	if tunebotUser, err = discordUser.GetUser(); err == nil {
		rsp = fmt.Sprintf("Your Discord account is already linked to %s user %s", Title, tunebotUser.Username)
		return rsp, nil
	}

	// If there isn't an account linked, GetUser will return a specific error
	if err == database.ErrNoDiscordUser {
		rsp = fmt.Sprintf("%s\nTo link your Discord account, react to this message", err)
		return rsp, linkUserCallback
	}

	// Any remaining errors that could have been returned from GetUser will be reported here
	return err.Error(), nil
}

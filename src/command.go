package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	core "github.com/tune-bot/core/src"
)

// A command callback takes a reference to the bot singleton and
// the list of args it was passed and returns a reply to send the server
// as well as a reaction callback if the message is awaiting a response
type commandCallback func(*Bot, *discordgo.MessageCreate, []string) (string, []string, reactionCallback)

var commands = map[string]commandCallback{
	"help":     cmdHelp,
	"kill":     cmdKill,
	"link":     cmdLink,
	"playlist": cmdPlaylist,
	"search":   cmdSearch,
}

func cmdHelp(_ *Bot, _ *discordgo.MessageCreate, _ []string) (string, []string, reactionCallback) {
	rsp := fmt.Sprintf(":notes: Hello, I'm %s!\n", Title)
	rsp += fmt.Sprintf(":calling: To get started, create an account on the %s app.\n", Title)
	rsp += fmt.Sprintf(":link: Then, link your Discord account using the following command: `%slink <%s username>`\n", CmdPrefix, Title)
	rsp += fmt.Sprintf(":play_pause: Manage your playlists with the `%splaylist` command\n", CmdPrefix)
	rsp += fmt.Sprintf(":mag: To add songs to your active playlist, use this command: `%ssearch <song query>`\n", CmdPrefix)
	rsp += fmt.Sprintf("\n%s %s\nReact to this message to browse the source code!", Title, version())
	return rsp, nil, browseSourceCallback
}

func cmdKill(b *Bot, _ *discordgo.MessageCreate, _ []string) (string, []string, reactionCallback) {
	b.stop()
	return "", nil, nil
}

func cmdLink(b *Bot, m *discordgo.MessageCreate, args []string) (string, []string, reactionCallback) {
	var err error
	var rsp string

	if len(args) < 1 {
		rsp = fmt.Sprintf("No %s username supplied", Title)
		return rsp, nil, nil
	}

	discordUser := core.Discord{Name: m.Author.Username}
	tunebotUser := core.User{Username: args[0]}

	// Check if the Tunebot user exists before trying to link a Discord account to it
	if _, err := core.FindUser(tunebotUser.Username); err != nil {
		return err.Error(), nil, nil
	}

	// Check if this Discord account is already linked
	if tunebotUser, err = discordUser.GetUser(); err == nil {
		rsp = fmt.Sprintf("Your Discord account is already linked to %s user %s", Title, tunebotUser.Username)
		return rsp, nil, nil
	}

	// If there isn't an account linked, GetUser will return a specific error
	if err == core.ErrNoDiscordUser {
		rsp = fmt.Sprintf("%s\nTo link your Discord account, react to this message", err)
		return rsp, nil, linkUserCallback
	}

	// Any remaining errors that could have been returned from GetUser will be reported here
	return err.Error(), nil, nil
}

func cmdPlaylist(b *Bot, m *discordgo.MessageCreate, args []string) (string, []string, reactionCallback) {
	var user core.User
	var err error

	if user, err = getAccount(m.Author.Username); err != nil {
		return err.Error(), nil, nil
	}

	playlists, count := ListPlaylists(user)
	msg := "**Playlist Management**\nYour playlists:\n"
	msg += playlists
	msg += "\n"
	nOptions := 0

	for i, emoji := range playlistOptionsEmojis {
		if i > 0 && count == 0 {
			continue
		}

		nOptions++
		msg += fmt.Sprintf("%s: %s a playlist\n", emoji, playlistOptions[i])
	}

	return msg, playlistOptionsEmojis[:nOptions], choosePlaylistOptionCallback
}

func cmdSearch(b *Bot, m *discordgo.MessageCreate, args []string) (string, []string, reactionCallback) {
	if len(args) == 0 {
		return "No search query provided!", nil, nil
	}

	b.MessageReactionAdd(m.ChannelID, m.ID, "🔎")
	query := strings.Join(args, " ")
	results, found := getCachedResults(query)

	// If the search results aren't already cached, get the results then cache them
	if !found {
		results = core.Search(query, 5)
		addToSearchCache(query, results)
	}

	msg := "**Search Results**\n\n"
	for i, result := range results {
		msg += fmt.Sprintf("%s %s\n", searchOptionsEmojis[i], result)
	}

	return msg, searchOptionsEmojis, chooseSearchOptionCallback
}

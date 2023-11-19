package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	core "github.com/tune-bot/core/src"
)

// Callbacks for reactions on bot messages
// A reaction callback takes a reference to the bot and the reaction event
// and returns a response to send to the server with a success indicator
type reactionCallback func(*Bot, *discordgo.MessageReactionAdd, []string) (string, bool)

func browseSourceCallback(b *Bot, r *discordgo.MessageReactionAdd, _ []string) (string, bool) {
	if dm, err := b.UserChannelCreate(r.UserID); err == nil {
		msg := fmt.Sprintf("Browse the source code for %s %s here:\nhttps://github.com/tune-bot/discord", Title, version())
		b.ChannelMessageSend(dm.ID, msg)
	}

	return "", true
}

func linkUserCallback(_ *Bot, r *discordgo.MessageReactionAdd, args []string) (string, bool) {
	var userId string
	var err error

	if userId, err = core.FindUser(args[0]); err != nil {
		return err.Error(), false
	}

	discordUser := core.Discord{
		Name:   r.Member.User.Username,
		UserId: userId,
	}

	if err = discordUser.Link(); err != nil {
		return err.Error(), false
	}

	return fmt.Sprintf("Linked Discord user %s to %s account: %s!", discordUser.Name, Title, args[0]), true
}

func chooseSearchOptionCallback(b *Bot, r *discordgo.MessageReactionAdd, args []string) (string, bool) {
	if r.Member.User.ID == b.State.User.ID {
		return "", false
	}

	var user core.User
	var err error

	if user, err = getAccount(r.Member.User.Username); err != nil {
		return err.Error(), false
	}

	// Find the search results for this query and add the chosen song
	query := strings.Join(args, " ")
	for i, emoji := range searchOptionsEmojis {
		if r.Emoji.Name == emoji {
			if results, found := getCachedResults(query); found {
				pendingAddSong = &results[i]
				break
			} else {
				return fmt.Sprintf("Error getting search results for \"%s\"", query), false
			}
		}
	}

	playlists, count := ListPlaylists(user)
	if count == 0 {
		pendingAddSong = nil
		return fmt.Sprintf("You have no playlists! Create one with `%splaylist`", CmdPrefix), false
	}

	if count == 1 {
		pendingAddPlaylist = &user.Playlists[0]
		AddPendingSong(b, r.ChannelID)
	} else {
		if channel, err := b.Channel(r.ChannelID); err == nil {
			suffix := ""
			emoji := ""

			// Check if we've added to any playlists yet
			if recentPlaylist, ok := mostRecentPlaylist[user.Username]; ok {
				suffix = fmt.Sprintf(" (%s to add to \"%s\")", EmojiAddToLastPlaylist, recentPlaylist.Name)
				emoji = EmojiAddToLastPlaylist
			}

			addPrompt := fmt.Sprintf("Add to which playlist?%s", suffix)
			b.askQuestion(channel, r.Member.User.Username, addPrompt, answerWhichPlaylist, emoji, addToLastPlaylistOption)
			b.ChannelMessageSend(r.ChannelID, playlists)
		} else {
			pendingAddSong = nil
			return err.Error(), false
		}
	}

	return "", false
}

func choosePlaylistOptionCallback(b *Bot, r *discordgo.MessageReactionAdd, args []string) (string, bool) {
	if r.Member.User.ID == b.State.User.ID {
		return "", false
	}

	var user core.User
	var err error

	if user, err = getAccount(r.Member.User.Username); err != nil {
		return err.Error(), false
	}

	switch r.Emoji.Name {
	case playlistOptionsEmojis[0]:
		playlistActions[user.Id] = PlaylistCreate
		return "Enter a name for the new playlist", true
	case playlistOptionsEmojis[1]:
		playlistActions[user.Id] = PlaylistRemove
		return "Delete which playlist?", true
	case playlistOptionsEmojis[2]:
		playlistActions[user.Id] = PlaylistShow
		return "Show which playlist?", true
	case playlistOptionsEmojis[3]:
		playlistActions[user.Id] = PlaylistToggle
		return "Toggle which playlist?", true
	}

	return "Not an option", false
}

func addToLastPlaylistOption(b *Bot, r *discordgo.MessageReactionAdd, args []string) (string, bool) {
	if user, err := getAccount(r.Member.User.Username); err == nil {
		if playlist, ok := mostRecentPlaylist[user.Username]; ok {
			answerWhichPlaylist(user, playlist.Name)
			AddPendingSong(b, r.ChannelID)
		}
	} else {
		return err.Error(), false
	}

	return "", false
}

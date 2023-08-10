package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	core "github.com/tune-bot/core/src"
)

// Fields and functions for letting users manage their playlists

var (
	// Add the pending song to the pending playlist
	pendingAddPlaylist *core.Playlist
	pendingAddSong     *core.Song

	// Keep track of what the most recently added-to playlist was for each user
	mostRecentPlaylist = make(map[string]*core.Playlist)
)

func AddPendingSong(b *Bot, channelID string) {
	if err := pendingAddSong.AddToPlaylist(pendingAddPlaylist.Id); err != nil {
		b.ChannelMessageSend(channelID, fmt.Sprintf("Error adding \"%s\" to playlist \"%s\":\n%s", pendingAddSong.Title.String, pendingAddPlaylist.Name, err.Error()))
	} else {
		b.ChannelMessageSend(channelID, fmt.Sprintf("Added \"%s\" to playlist \"%s\"", pendingAddSong.Title.String, pendingAddPlaylist.Name))
	}

	pendingAddPlaylist = nil
	pendingAddSong = nil
}

func GetPlaylist(user core.User, name string) *core.Playlist {
	for _, playlist := range user.Playlists {
		if playlist.Name == name {
			return &playlist
		}
	}

	return nil
}

func ListPlaylists(user core.User) (string, int) {
	var playlists string = ""
	var count int = len(user.Playlists)

	for i, playlist := range user.Playlists {
		disabled := ""

		if !playlist.Enabled {
			disabled = " (disabled)"
		}

		playlists += fmt.Sprintf("%d. %s%s\n", i, playlist.Name, disabled)
	}

	if count == 0 {
		playlists = "You have no playlists\n"
	}

	return playlists, count
}

func ListSongs(playlist core.Playlist) string {
	songs := fmt.Sprintf("**Songs in Playlist \"%s\"**\n", playlist.Name)
	for _, song := range playlist.Songs {
		songs += fmt.Sprintf("*%s* by %s\n", song.Title.String, song.Artist.String)
	}

	return songs
}

// Playlist action request queue:
// Map user id to playlist actions
// When a user in the list sends their next message, the action will run on it
var playlistActions = make(map[string]playlistAction)

// An action that is run on a pending user's message, return a response message and success status
type playlistAction func(*Bot, *discordgo.Message) (string, bool)

func PlaylistCreate(b *Bot, m *discordgo.Message) (string, bool) {
	var user core.User
	var err error

	if user, err = getAccount(m.Author.Username); err != nil {
		return fmt.Sprintf("Error getting account for %s", m.Author.Username), false
	}

	playlist := core.Playlist{
		Name:    strings.TrimSpace(m.Content),
		Enabled: true,
	}

	if err = playlist.Create(user.Id); err != nil {
		return err.Error(), true
	}

	return fmt.Sprintf("Created playlist \"%s\"", playlist.Name), true
}

func PlaylistRemove(b *Bot, m *discordgo.Message) (string, bool) {
	var user core.User
	var err error
	var playlist *core.Playlist

	if user, err = getAccount(m.Author.Username); err != nil {
		return fmt.Sprintf("Error getting account for %s", m.Author.Username), false
	}

	if playlist = GetPlaylist(user, m.Content); playlist == nil {
		return fmt.Sprintf("Could not find playlist \"%s\"", m.Content), false
	}

	if err = playlist.Delete(); err != nil {
		return err.Error(), true
	}

	return fmt.Sprintf("Deleted playlist \"%s\"", playlist.Name), true
}

func PlaylistShow(b *Bot, m *discordgo.Message) (string, bool) {
	var user core.User
	var err error
	var playlist *core.Playlist

	if user, err = getAccount(m.Author.Username); err != nil {
		return fmt.Sprintf("Error getting account for %s", m.Author.Username), false
	}

	if playlist = GetPlaylist(user, m.Content); playlist == nil {
		return fmt.Sprintf("Could not find playlist \"%s\"", m.Content), false
	}

	if dm, err := b.UserChannelCreate(m.Author.ID); err == nil {
		msg := ListSongs(*playlist)
		b.ChannelMessageSend(dm.ID, msg)
		return fmt.Sprintf("The songs in playlist \"%s\" have been sent to your DMs", playlist.Name), true
	}
	return "", true
}

func PlaylistToggle(b *Bot, m *discordgo.Message) (string, bool) {
	var user core.User
	var err error
	var playlist *core.Playlist
	var state string = "En"

	if user, err = getAccount(m.Author.Username); err != nil {
		return fmt.Sprintf("Error getting account for %s", m.Author.Username), false
	}

	if playlist = GetPlaylist(user, m.Content); playlist == nil {
		return fmt.Sprintf("Could not find playlist \"%s\"", m.Content), false
	}

	playlist.Enabled = !playlist.Enabled
	if !playlist.Enabled {
		state = "Dis"
	}

	if err = playlist.Update(); err != nil {
		return err.Error(), true
	}

	return fmt.Sprintf("%sabled playlist \"%s\"", state, playlist.Name), true
}

package main

import (
	"fmt"

	core "github.com/tune-bot/core/src"
)

// An answerCallback is a function that runs on a string response to a question asked by the bot to a user
type answerCallback func(core.User, string) string

func answerWhichPlaylist(user core.User, answer string) string {
	for _, playlist := range user.Playlists {
		if playlist.Name == answer {
			pendingAddPlaylist = &playlist
			mostRecentPlaylist[user.Username] = &playlist
			return ""
		}
	}

	return fmt.Sprintf("You have no playlist called \"%s\"", answer)
}

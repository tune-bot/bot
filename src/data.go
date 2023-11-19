package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	core "github.com/tune-bot/core/src"
)

// Declare constants and data management resources

const (
	Title     = "TuneBot"
	CmdPrefix = "?"
	EmojiAddToLastPlaylist = "📋"
)

var (
	cleanKeyRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

var searchOptionsEmojis []string = []string{
	"🔴",
	"🟠",
	"🟢",
	"🔵",
	"🟣",
}

var playlistOptionsEmojis []string = []string{
	"📜",
	"🗑️",
	"📂",
	"⚙️",
}

var playlistOptions []string = []string{
	"Add",
	"Remove",
	"Show",
	"Toggle",
}

// Keep a cache of search results, mapping cleaned query string to songs
var searchCache map[string][]core.Song = make(map[string][]core.Song)

func cacheKey(query string) string {
	return strings.ToLower(cleanKeyRegex.ReplaceAllString(query, ""))
}

func addToSearchCache(query string, songs []core.Song) {
	key := cacheKey(query)
	if _, found := searchCache[key]; found {
		return
	}

	searchCache[key] = songs
}

func getCachedResults(query string) ([]core.Song, bool) {
	key := cacheKey(query)
	if songs, found := searchCache[key]; found {
		return songs, true
	}

	return nil, false
}

func getAccount(discordName string) (core.User, error) {
	var user core.User
	var err error

	// Find the user associated with this account
	discordUser := core.Discord{Name: discordName}
	if user, err = discordUser.GetUser(); err != nil {
		return user, fmt.Errorf("No %s account associated with Discord user: %s\n%s", Title, discordUser.Name, err.Error())
	}

	return user, nil
}

func version() string {
	if v, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output(); err == nil {
		return strings.TrimSpace(string(v))
	}

	return ""
}

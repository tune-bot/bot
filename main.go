package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/tune-bot/discord/bot"
)

func main() {
	if session, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_TOKEN"))); err == nil {
		bot.Start(session)
	}
}

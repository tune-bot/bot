package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"

	bot "github.com/tune-bot/discord/bot"
	data "github.com/tune-bot/discord/data"
)

func main() {
	if session, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv(data.TOKEN_VAR))); err == nil {
		bot.Start(session)
	}
}

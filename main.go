package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + token)

	if discord.StateEnabled {

	}

	if err != nil {
		return
	}
}

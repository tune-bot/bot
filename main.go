package main

import (
	"github.com/tune-bot/discord/bot"
)

func main() {
	if bot.Tunebot != nil {
		bot.Tunebot.Start()
	}
}

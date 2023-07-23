package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/tune-bot/database"
)

// A command callback takes a reference to the bot singleton and
// the list of args it was passed and returns a reply to send the server
type commandCallback func(b *Bot, m *discordgo.MessageCreate, args []string) string

var commands = map[string]commandCallback{
	"kill": cmdKill,
	"link": cmdLink,
}

func cmdKill(b *Bot, _ *discordgo.MessageCreate, _ []string) string {
	b.Stop()
	return ""
}

func cmdLink(_ *Bot, m *discordgo.MessageCreate, _ []string) string {
	var err error

	dbUser := database.User{Username: m.Author.Username}
	discordUser := database.Discord{Name: dbUser.Username}

	if err = dbUser.Create(); err != nil {
		log.Println(err)
		return err.Error()
	}

	discordUser.UserId = dbUser.Id
	if err = discordUser.Link(); err != nil {
		log.Println(err)
		return err.Error()
	}

	if dbUser, err = discordUser.GetUser(); err != nil {
		log.Println(err)
		return err.Error()
	}

	return fmt.Sprintf("Linked Discord user %s", discordUser.Name)
}

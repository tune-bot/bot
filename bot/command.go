package bot

// A command callback takes a reference to the bot singleton and
// the list of args it was passed and returns a reply to send the server
type commandCallback func(b *Bot, args []string) string

var commands = map[string]commandCallback{
	"exit": exitCallback,
}

func exitCallback(b *Bot, args []string) string {
	b.Stop()
	return ""
}

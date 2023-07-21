package command

import (
	"fmt"
	"strings"
)

// A command callback takes the list of args it was passed and returns the reply
// to send the server
type commandCallback func(args []string) string

var Commands = map[string]commandCallback{
	"test": testCallback,
}

func testCallback(args []string) string {
	return fmt.Sprintf("Hello, world!\n%s\n", strings.Join(args, ", "))
}

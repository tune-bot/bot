package command

import (
	"fmt"
	"strings"
)

type commandCallback func(channelID string, args []string) bool

var Commands = map[string]commandCallback{
	"test": testCallback,
}

func testCallback(channelID string, args []string) bool {
	fmt.Printf("Hello, world!\n%s\n", strings.Join(args, ", "))

	return true
}

package main

import "os/exec"

const (
	Title     = "TuneBot"
	CmdPrefix = "?"
)

func version() string {
	if v, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output(); err == nil {
		return string(v)
	}

	return ""
}

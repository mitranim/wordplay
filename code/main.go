package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitranim/cmd"
)

var commands = cmd.Map{}

func init() {
	time.Local = nil
	spew.Config.Indent = "  "
}

func main() {
	defer cmd.Report()
	commands.Get()()
}

package main

import (
	"time"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

var commands = cmd.Map{`norm`: cmdNorm}

func init() { time.Local = nil }

func main() {
	defer gg.Fatal()
	commands.Get()()
}

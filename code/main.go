package main

import (
	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

func main() {
	defer gg.Fatal()

	cmd.Map{
		`norm`:             CmdNorm,
		`discord_download`: CmdDiscordDownload,
	}.Get()()
}

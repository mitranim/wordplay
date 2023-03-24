package main

import "os"

func CmdDiscordDownload() {
	var tar DiscordDownload
	tar.Log = os.Stderr
	tar.OutPath = `testdata/download.md`
	tar.AfterMsgId = ConfGlobal.ReqDiscordAfterMsgId()
	tar.Download()
}

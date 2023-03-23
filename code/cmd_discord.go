package main

func CmdDiscordDownload() {
	var tar DiscordDownload
	tar.OutPath = `download.md`
	tar.AfterMsgId = ConfGlobal.ReqDiscordAfterMsgId()
	tar.Download()
}

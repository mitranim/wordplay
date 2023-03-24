package main

import (
	"time"
	u "unsafe"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/mitranim/gg"
	"github.com/mitranim/gt"
)

func init() {
	time.Local = nil
	gg.TraceBaseDir = gg.Cwd()
}

var (
	DiscordApiUrl  = gg.ParseTo[gt.NullUrl](`https://discord.com/api/v9`)
	DiscordSelfUrl = gg.ParseTo[gt.NullUrl](`https://discord.com/channels/@me`)
	ConfGlobal     = gg.With((*Conf).Init)
)

type Conf struct {
	DiscordChan            string       `json:"-" env:"DISCORD_CHAN"`
	DiscordCookie          string       `json:"-" env:"DISCORD_COOKIE"`
	DiscordAuthorization   string       `json:"-" env:"DISCORD_AUTHORIZATION"`
	DiscordSuperProperties string       `json:"-" env:"DISCORD_SUPER_PROPERTIES"`
	DiscordAfterMsgId      DiscordMsgId `json:"-" env:"DISCORD_AFTER_MSG_ID"`
}

func (self *Conf) Init() {
	// Ideally, this would be called in `init`. Placed here for technical reasons.
	gg.Try(godotenv.Load(`.env.properties`))

	gg.Try(envdecode.StrictDecode(self))
}

func (self Conf) ReqDiscordChan() string {
	return ReqField[string](self, u.Offsetof(self.DiscordChan))
}

func (self Conf) ReqDiscordCookie() string {
	return ReqField[string](self, u.Offsetof(self.DiscordCookie))
}

func (self Conf) ReqDiscordAuthorization() string {
	return ReqField[string](self, u.Offsetof(self.DiscordAuthorization))
}

func (self Conf) ReqDiscordSuperProperties() string {
	return ReqField[string](self, u.Offsetof(self.DiscordSuperProperties))
}

func (self Conf) ReqDiscordAfterMsgId() DiscordMsgId {
	return ReqField[DiscordMsgId](self, u.Offsetof(self.DiscordAfterMsgId))
}

func (self Conf) DiscordApiUrlMsgs() gt.NullUrl {
	return DiscordApiUrl.AddPath(`channels`, self.ReqDiscordChan(), `messages`)
}

func (self Conf) DiscordReferer() gt.NullUrl {
	return DiscordSelfUrl.AddPath(self.ReqDiscordChan())
}

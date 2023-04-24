package main

import (
	"io"
	"net/url"
	"path/filepath"
	"time"
	u "unsafe"

	"github.com/mitranim/gg"
	"github.com/mitranim/gr"
	"github.com/mitranim/gt"
)

const logPrefixDiscordDownload = `[discord_download] `

type DiscordDownload struct {
	Log        io.Writer
	OutPath    string
	Pages      DiscordMsgPages
	AfterMsgId DiscordMsgId
	MaxMsg     DiscordMsg
}

func (self DiscordDownload) ReqAfterMsgId() DiscordMsgId {
	return gg.Or(
		self.MaxMsg.Pk(),
		ReqField[DiscordMsgId](self, u.Offsetof(self.AfterMsgId)),
	)
}

func (self DiscordDownload) Req() *gr.Req {
	return DiscordReq().Query(url.Values{
		`after`: {gg.String(self.ReqAfterMsgId())},
	})
}

// Must be defined as pointer method for compatibility with `defer`.
func (self *DiscordDownload) Flush() {
	path := self.OutPath
	if gg.IsZero(path) {
		return
	}

	gg.MkdirAll(filepath.Dir(path))
	gg.WriteFile(path, self.Pages.String())
}

func (self *DiscordDownload) Download() {
	defer self.Flush()

	const limit = 128
	for ind := range gg.Iter(limit) {
		var page DiscordMsgPage
		self.Req().Res().Ok().Json(&page)

		if gg.IsEmpty(page) {
			Fprintln(self.Log, logPrefixDiscordDownload, `found empty page, done; output path: `, self.OutPath)
			return
		}
		gg.Append(&self.Pages, page)
		self.SetMaxMsg(page.Max())

		Fprintln(self.Log, logPrefixDiscordDownload, `downloaded page `, ind, `, sleeping`)
		time.Sleep(time.Second)
		Fprintln(self.Log, logPrefixDiscordDownload, `continuing`)
	}

	Fprintln(self.Log, logPrefixDiscordDownload, `exceeded maximum number of iterations: `, limit)
}

func (self *DiscordDownload) SetMaxMsg(next DiscordMsg) {
	prev := self.MaxMsg

	if gg.IsZero(next) {
		panic(gg.Errv(`unexpected missing next msg`))
	}
	if !prev.Less(next) {
		panic(gg.Errf(`expected prev msg to be lesser than next msg; prev timestamp: %v, next timestamp: %v`, prev.Timestamp, next.Timestamp))
	}

	self.MaxMsg = next
}

type DiscordMsgPages []DiscordMsgPage

func (self DiscordMsgPages) MinPk() DiscordMsgId { return self.Min().MinPk() }
func (self DiscordMsgPages) MaxPk() DiscordMsgId { return self.Max().MaxPk() }
func (self DiscordMsgPages) MinMsg() DiscordMsg  { return self.Min().Min() }
func (self DiscordMsgPages) MaxMsg() DiscordMsg  { return self.Max().Max() }
func (self DiscordMsgPages) Min() DiscordMsgPage { return gg.Min(self...) }
func (self DiscordMsgPages) Max() DiscordMsgPage { return gg.Max(self...) }

// Suboptimal, but frees us from knowing or caring about msg ordering in Discord
// API responses.
func (self DiscordMsgPages) String() string {
	tar := gg.Concat(self...)
	tar.Norm()
	return tar.String()
}

type DiscordMsgPage []DiscordMsg

func (self DiscordMsgPage) MinPk() DiscordMsgId { return self.Min().Pk() }
func (self DiscordMsgPage) MaxPk() DiscordMsgId { return self.Max().Pk() }
func (self DiscordMsgPage) Min() DiscordMsg     { return gg.Min(self...) }
func (self DiscordMsgPage) Max() DiscordMsg     { return gg.Max(self...) }
func (self DiscordMsgPage) Norm()               { SortReverse(self) }

func (self DiscordMsgPage) Less(val DiscordMsgPage) bool {
	if gg.IsEmpty(val) {
		return true
	}
	if gg.IsEmpty(self) {
		return false
	}

	// return self.Max().Less(val.Min())
	return gg.Head(self).Less(gg.Head(val))
}

func (self DiscordMsgPage) Properties() string {
	return self.Max().Properties()
}

func (self DiscordMsgPage) String() string {
	return JoinLinesSparse(self.Properties(), self.Content())
}

func (self DiscordMsgPage) Content() string {
	return JoinLinesSparse(gg.Map(self, DiscordMsg.String)...)
}

type DiscordAuthorId string

type DiscordAuthor struct {
	Id            DiscordAuthorId `json:"id"`
	Username      string          `json:"username"`
	Discriminator string          `json:"discriminator"`
}

type DiscordMsgId string

type DiscordMsg struct {
	Id        DiscordMsgId          `json:"id"`
	Content   string                `json:"content"`
	Author    gg.Zop[DiscordAuthor] `json:"author"`
	Timestamp gt.NullTime           `json:"timestamp"`
}

func (self DiscordMsg) Pk() DiscordMsgId { return self.Id }

func (self DiscordMsg) Less(val DiscordMsg) bool {
	return self.Timestamp.Before(val.Timestamp)
}

func (self DiscordMsg) String() string {
	if gg.IsZero(self.Content) {
		return ``
	}
	return gg.Str(self.Content, ` © `, self.Author.Val.Username)
}

func (self DiscordMsg) Properties() string {
	return gg.Str(`DISCORD_AFTER_MSG_ID=`, self.Id)
}

func DiscordReq() *gr.Req {
	return new(gr.Req).
		Url(ConfGlobal.DiscordApiUrlMsgs().Maybe()).
		Get().
		HeadAdd(`authority`, `discord.com`).
		HeadAdd(`x-super-properties`, ConfGlobal.ReqDiscordSuperProperties()).
		HeadAdd(`x-discord-locale`, `en-US`).
		HeadAdd(`x-debug-options`, `bugReporterEnabled`).
		HeadAdd(`accept-language`, `en-US`).
		HeadAdd(`authorization`, ConfGlobal.ReqDiscordAuthorization()).
		HeadAdd(`user-agent`, `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.273 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36`).
		HeadAdd(`sec-fetch-site`, `same-origin`).
		HeadAdd(`sec-fetch-mode`, `cors`).
		HeadAdd(`sec-fetch-dest`, `empty`).
		HeadAdd(`referer`, ConfGlobal.DiscordReferer().String()).
		HeadAdd(`cookie`, ConfGlobal.ReqDiscordCookie())
}

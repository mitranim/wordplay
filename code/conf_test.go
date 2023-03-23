package main

import (
	"testing"

	"github.com/mitranim/gg/gtest"
)

func Test_ConfGlobal(t *testing.T) {
	t.Skip(`enable on demand`)

	defer gtest.Catch(t)

	gtest.NotZero(ConfGlobal.DiscordApiUrlMsgs())
	gtest.NotZero(ConfGlobal.ReqDiscordSuperProperties())
	gtest.NotZero(ConfGlobal.ReqDiscordAuthorization())
	gtest.NotZero(ConfGlobal.DiscordReferer().String())
	gtest.NotZero(ConfGlobal.ReqDiscordCookie())
}

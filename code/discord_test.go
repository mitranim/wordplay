package main

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
	"github.com/mitranim/gt"
)

func Test_DiscordDownload_Flush(t *testing.T) {
	defer gtest.Catch(t)

	const PATH = `testdata/test.md`

	var tar DiscordDownload
	tar.OutPath = PATH

	tar.Flush()
	gtest.Eq(gg.ReadFile[string](PATH), `DISCORD_AFTER_MSG_ID=`)

	tar.Pages = testDiscordMsgPages()
	tar.Flush()
	gtest.Eq(gg.ReadFile[string](PATH), TEST_DISCORD_MSG_PAGES_STRING)
}

func Test_DiscordMsgPages_String(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(
		testDiscordMsgPages().String(),
		TEST_DISCORD_MSG_PAGES_STRING,
	)
}

func testDiscordMsgPages() DiscordMsgPages {
	return DiscordMsgPages{
		{
			{
				Id:        `id_0`,
				Content:   `one`,
				Author:    gg.ZopVal(DiscordAuthor{Username: `two`}),
				Timestamp: gt.NullDateUTC(1000, 1, 1),
			},
			{
				Id:        `id_2`,
				Content:   `five`,
				Author:    gg.ZopVal(DiscordAuthor{Username: `six`}),
				Timestamp: gt.NullDateUTC(1000, 3, 1),
			},
		},
		{
			{
				Id:        `id_1`,
				Content:   `three`,
				Author:    gg.ZopVal(DiscordAuthor{Username: `four`}),
				Timestamp: gt.NullDateUTC(1000, 2, 1),
			},
			{
				Id:        `id_4`,
				Content:   ``,
				Author:    gg.ZopVal(DiscordAuthor{Username: `nine`}),
				Timestamp: gt.NullDateUTC(1000, 5, 1),
			},
			{
				Id:        `id_3`,
				Content:   `seven`,
				Author:    gg.ZopVal(DiscordAuthor{Username: `eight`}),
				Timestamp: gt.NullDateUTC(1000, 4, 1),
			},
		},
	}
}

const TEST_DISCORD_MSG_PAGES_STRING = `DISCORD_AFTER_MSG_ID=id_4

seven © eight

five © six

three © four

one © two`

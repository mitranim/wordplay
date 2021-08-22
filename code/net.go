package main

import (
	"log"
	h "net/http"

	"github.com/google/go-github/github"
	x "github.com/mitranim/gax"
	"github.com/mitranim/goh"
	"github.com/mitranim/rout"
	"github.com/mitranim/try"
	e "github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type CommitHash string

const (
	ACCESS_TOKEN = ""
	REPO_OWNER   = "mitranim"
	REPO_NAME    = "wordplay"
	SERVER_PORT  = "57830"
)

var tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ACCESS_TOKEN})

type IndexDat struct {
	Version CommitHash
	Grouped EntryMap
	Entries Entries
}

func main() {
	log.Printf("Starting server on http://localhost:%v", SERVER_PORT)
	try.To(h.ListenAndServe(":"+SERVER_PORT, h.HandlerFunc(respond)))
}

func respond(rew Rew, req *Req) {
	goh.ErrHandler(rew, req, false, rout.Route(rew, req, routes))
}

func routes(r rout.R) { r.Get(`^/$`, routeIndex) }

func routeIndex(rew Rew, req *Req) {
	goh.Respond(rew, req, func() h.Handler {
		entries, version := ReadAndParseBackingFile(req.Context())

		return goh.BytesOk(RenderIndex(IndexDat{
			Version: version,
			Grouped: GroupEntriesByAuthor(entries),
			Entries: entries,
		}))
	})
}

func RenderIndex(dat IndexDat) x.Bui {
	return RenderHtml(func(b x.Bui) {
		if dat.Version != "" {
			b.E(`div`, nil, `Version `+dat.Version)
		}
		for _, author := range dat.Grouped.Keys {
			b.E(`h1`, nil, author)
			for _, entry := range dat.Grouped.Map[author] {
				b.E(`div`, nil, formatEntryPhrase(entry), ` `, formatEntryMeanings(entry))
			}
		}
	})
}

func RenderHtml(children ...interface{}) x.Bui {
	return x.F(
		x.Str(x.Doctype),
		E(`html`, nil,
			E(`head`, nil,
				E(`meta`, AP(`charset`, `utf-8`)),
				E(`meta`, AP(`http-equiv`, `X-UA-Compatible`, `content`, `IE=edge`)),
				E(`meta`, AP(`name`, `viewport`, `content`, `width=device-width, initial-scale=1`)),
				E(`link`, AP(`rel`, `icon`, `href`, `data:;base64,=`)),
				E(`title`, nil, `wordplay`),
				E(`style`, nil, style),
			),
			E(`body`, nil, children),
		),
	)
}

func formatEntryPhrase(entry Entry) string {
	return bytesToStringAlloc(entry.AppendPhrase(nil))
}

func formatEntryMeanings(entry Entry) string {
	return bytesToStringAlloc(entry.AppendMeanings(nil))
}

/*
In principle, we might want to deduplicate concurrent instances of this request,
returning the result of a single request to multiple callers that have been
waiting on it. In reality, we probably won't have concurrent requests.
*/
func ReadBackingFile(ctx Ctx) (string, CommitHash) {
	client := github.NewClient(oauth2.NewClient(ctx, tokenSource))

	repoContent, _, err := client.Repositories.GetReadme(ctx, REPO_OWNER, REPO_NAME, nil)
	if err != nil {
		panic(e.Wrap(err, "failed to fetch backing file from repository"))
	}

	// This call performs an unnecessary bytes-to-string conversion, which we
	// then have to reverse. Waste of CPU cycles. Of course, we have no business
	// complaining when our own parser allocates strings instead of reslicing.
	content, err := repoContent.GetContent()
	if err != nil {
		panic(e.Wrap(err, "failed to decode content of backing file"))
	}

	return content, CommitHash(repoContent.GetSHA())
}

func ReadAndParseBackingFile(ctx Ctx) (Entries, CommitHash) {
	defer try.Detail(`failed to read and parse backing file`)
	content, version := ReadBackingFile(ctx)
	entries := ParseEntries(content)
	return entries, version
}

// nolint:deadcode
func WriteBackingFile(ctx Ctx, content []byte, version CommitHash) {
	client := github.NewClient(oauth2.NewClient(ctx, tokenSource))
	msg := "(automatic)"

	_, _, err := client.Repositories.UpdateFile(
		ctx,
		REPO_OWNER,
		REPO_NAME,
		"readme.md",
		&github.RepositoryContentFileOptions{
			Message: &msg,
			Content: content,
			SHA:     (*string)(&version),
			Author:  &github.CommitAuthor{},
		},
	)
	if err != nil {
		panic(e.Wrap(err, "failed to update backing file"))
	}
}

var style = `
:root {
	background-color: hsl(200, 5%, 10%);
	color: white;
}

input, textarea {
	background-color: inherit;
	color: inherit;
}

.block {
	display: block;
}

.gaps-v-1 > :not(:first-child) {
	margin-top: 1rem;
}

.gaps-v-0x5 > :not(:first-child) {
	margin-top: 0.5rem;
}
`

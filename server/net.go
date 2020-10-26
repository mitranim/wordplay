package main

import (
	"bytes"
	"context"
	"fmt"
	ht "html/template"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Rew = http.ResponseWriter
type Req = *http.Request
type Ctx = context.Context
type CommitHash = string

const ACCESS_TOKEN = "<redacted>"
const REPO_OWNER = "mitranim"
const REPO_NAME = "wordplay"
const SERVER_PORT = "57830"

var tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ACCESS_TOKEN})

func main() {
	log.Printf("Starting server on http://localhost:%v", SERVER_PORT)
	err := http.ListenAndServe(":"+SERVER_PORT, http.HandlerFunc(routes))
	if err != nil {
		panic(err)
	}
}

func routes(rew Rew, req Req) {
	switch req.URL.Path {
	case "/":
		indexRoute(rew, req)
	default:
		http.NotFound(rew, req)
	}
}

func indexRoute(rew Rew, req Req) {
	if req.Method != http.MethodGet {
		methodNotAllowed(rew, req)
		return
	}

	entries, version, err := ReadAndParseBackingFile(req.Context())
	if err != nil {
		rew.WriteHeader(http.StatusInternalServerError)
		rew.Write(stringToBytesAlloc(err.Error()))
		return
	}

	var buf bytes.Buffer

	err = templates.ExecuteTemplate(&buf, "index.html", struct {
		Version CommitHash
		Grouped EntryMap
		Entries []Entry
	}{
		Version: version,
		Grouped: GroupEntriesByAuthor(entries),
		Entries: entries,
	})
	if err != nil {
		rew.WriteHeader(http.StatusInternalServerError)
		rew.Write(stringToBytesAlloc(err.Error()))
		return
	}

	rew.Write(buf.Bytes())
}

var templates = ht.Must(ht.New("").
	Funcs(ht.FuncMap{
		"formatEntryPhrase":   formatEntryPhrase,
		"formatEntryMeanings": formatEntryMeanings,
	}).
	ParseGlob("templates/*"))

func methodNotAllowed(rew Rew, req Req) {
	rew.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(rew, `unsupported method: %v %v`, req.Method, req.URL.Path)
}

func formatEntryPhrase(entry Entry) string {
	return bytesToStringAlloc(entry.FormatAppendPhrase(nil))
}

func formatEntryMeanings(entry Entry) string {
	return bytesToStringAlloc(entry.FormatAppendMeanings(nil))
}

/*
In principle, we might want to deduplicate concurrent instances of this request,
returning the result of a single request to multiple callers that have been
waiting on it. In reality, we probably won't have concurrent requests.
*/
func ReadBackingFile(ctx Ctx) ([]byte, CommitHash, error) {
	client := github.NewClient(oauth2.NewClient(ctx, tokenSource))

	repoContent, _, err := client.Repositories.GetReadme(ctx, REPO_OWNER, REPO_NAME, nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to fetch backing file from repository")
	}

	// This call performs an unnecessary bytes-to-string conversion, which we
	// then have to reverse. Waste of CPU cycles. Of course, we have no business
	// complaining when our own parser allocates strings instead of reslicing.
	content, err := repoContent.GetContent()
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to decode content of backing file")
	}

	return stringToBytesAlloc(content), repoContent.GetSHA(), nil
}

func ReadAndParseBackingFile(ctx Ctx) ([]Entry, CommitHash, error) {
	content, version, err := ReadBackingFile(ctx)
	if err != nil {
		return nil, "", err
	}

	entries, err := ParseEntries(content)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse content of backing file")
	}

	return entries, version, nil
}

func WriteBackingFile(ctx Ctx, content []byte, version CommitHash) error {
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
			SHA:     &version,
			Author:  &github.CommitAuthor{},
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed to update backing file")
	}

	return nil
}

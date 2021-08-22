package main

import (
	"log"
	h "net/http"

	"github.com/mitranim/try"
)

func init() { commands[`srv`] = cmdSrv }

func cmdSrv() {
	log.Printf("Starting server on http://localhost:%v", SERVER_PORT)
	try.To(h.ListenAndServe(":"+SERVER_PORT, h.HandlerFunc(respond)))
}

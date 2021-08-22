package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mitranim/try"
	"github.com/pkg/errors"
)

func main() {
	time.Local = nil

	err := try.Catch(runMain)
	if err != nil {
		fmt.Printf("%T: %+v\n", err, err)
		os.Exit(1)
	}
}

func runMain() {
	cmd := os.Args[1]

	switch cmd {
	case "srv":
		cmdSrv()
	case "times":
		cmdTimes()
	case "reorder":
		cmdReorder()
	default:
		panic(errors.Errorf(`unknown cmd %q`, cmd))
	}
}

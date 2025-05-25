package main

import (
	"github.com/oyaro-tech/dispatchgo/internal/app"
	"log"
	"runtime/debug"
)

func main() {
	a := app.New()

	if err := a.Run(); err != nil {
		if a.IsDebug() {
			debug.PrintStack()
		}

		log.Fatal(err)
	}
}

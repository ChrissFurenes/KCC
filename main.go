package main

import (
	"log"
	"os"

	"github.com/chrissfurenes/kcc/app"
)

var version = "0.8.1-beta.4"

func main() {
	app := app.NewApp(version)

	handled, err := app.HandleArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if handled {
		return
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

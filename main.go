package main

import (
	"log"
	"os"

	"github.com/chrissfurenes/kcc/app"
)

func main() {
	app := app.NewApp()

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

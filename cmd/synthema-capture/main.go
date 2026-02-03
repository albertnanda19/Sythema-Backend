package main

import (
	"log"

	"synthema/internal/bootstrap"
)

func main() {
	app, err := bootstrap.BootstrapCapture()
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info("bootstrap complete")
}

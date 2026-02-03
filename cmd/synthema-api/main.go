package main

import (
	"log"
	"strconv"

	"synthema/internal/bootstrap"
)

func main() {
	app, err := bootstrap.BootstrapAPI()
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info("bootstrap complete")

	if err := app.App.Listen(":" + strconv.Itoa(app.Config.API.Port)); err != nil {
		log.Fatal(err)
	}
}

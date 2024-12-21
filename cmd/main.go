package main

import (
	"log"

	"github.com/syirnik/GO_Yandex/internal/application"
)

func main() {
	app := application.New()
	//app.Run()
	if err := app.RunServer(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

package main

import (
	"shopify/internal/app"
	"shopify/internal/config"
)

func main() {
	config := config.NewConfig()
	application := app.NewApplication(config)

	if err := application.Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"user-management-api/internal/app"
	"user-management-api/internal/config"
)

func main() {
	config := config.NewConfig()
	application := app.NewApplication(config)

	if err := application.Run(); err != nil {
		panic(err)
	}
}

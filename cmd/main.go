package main

import (
	"review-rotator/internal/app"
)

func main() {
	app := app.NewApp()
	app.Run()
}
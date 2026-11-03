package main

import (
	"elcomware/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}

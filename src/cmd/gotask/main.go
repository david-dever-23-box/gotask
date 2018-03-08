package main

import (
	"os"

	"gotask/cli"
)

func main() {
	app := cli.NewApp()
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

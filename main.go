package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
	"os"
)

var (
	Log *Logger
)

func init() {
	Log = NewLogger()
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmdRun(),
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

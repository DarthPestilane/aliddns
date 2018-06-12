package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"os"
)

var (
	Log *Logger
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("load env failed: %v", err))
	}
	Log = NewLogger()
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmdRun(),
	}
	app.Run(os.Args)
}

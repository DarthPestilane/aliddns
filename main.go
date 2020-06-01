package main

import (
	"github.com/DarthPestilane/aliddns/app/cmd"
	"github.com/DarthPestilane/aliddns/bootstrap"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	bootstrap.Boot()
	cmd.Run()
}

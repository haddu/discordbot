package main

import (
	"log"
	"os"
	"path"
	"fmt"

	cli "github.com/jawher/mow.cli"

	"github.com/c2nc/discordbot/discordbot/bot"
)

const (
	appDescr = "This application doing something"
)

var (
	appVer string
	appName = path.Base(os.Args[0])
)

func main() {
	app := cli.App(appName, appDescr)
	app.Version("V version", fmt.Sprintf("%s %s", appName, appVer))

	bot.Init(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("%s failed: %v", appName, err)
	}
}
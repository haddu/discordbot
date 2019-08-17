package bot

import (
	"log"

	cli "github.com/jawher/mow.cli"
)

//const token string = "NjExOTEwOTM3NDUxOTU0MTk2.XVas4g.mOOx7hUvwivq510BYfGIeKVw6xo"

func Init(app *cli.Cli) {
	// Start cli command
	app.Command("start", "start bot session", StartSession)
}

func StartSession(cmd *cli.Cmd) {
	var (
		token  = cmd.StringOpt("t token", "", "authority token")
		rooms = cmd.IntOpt("r rooms", 20,  "rooms capacity")

		client *client
		err    error
	)

	cmd.Before = func() {
		if client, err = New(*token, *rooms); err != nil {
			log.Fatalf("client failed: %v", err)
		}
	}

	cmd.Action = func() {
		client.Start()
		log.Printf("Bot is now running. Press CTRL-C to exit.")
		waitForInterupt()
		client.Close()
	}
}

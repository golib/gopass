package main

import (
	"os"

	"github.com/golib/cli"
)

var (
	app *cli.App
)

func main() {
	app = cli.NewApp()
	app.Name = "gopass"
	app.Version = "1.0.0"
	app.Author = "Spring MC"
	app.Commands = []cli.Command{
		{
			Name:  "chrome",
			Usage: "convert credentials from csv file to Chrome Passwords",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "csv",
					Usage: "specify csv `FILE` of credentials, exported from lastpass or 1Password",
				},
			},
			Action: chromeImporter(),
		},
		{
			Name:  "keychain",
			Usage: "import credentials from csv file to keychain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "csv",
					Usage: "specify csv `FILE` of credentials",
				},
				cli.BoolFlag{
					Name:  "update",
					Usage: "update old credential if existed",
				},
			},
			Action: keychainImporter(),
		},
	}

	app.Run(os.Args)
}

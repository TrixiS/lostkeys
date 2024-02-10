package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/TrixiS/lostkeys/pkg/command_context"
	"github.com/TrixiS/lostkeys/pkg/commands"
	"github.com/boltdb/bolt"
	"github.com/urfave/cli/v2"
)

const LostKeysDirname = ".lostkeys"
const LostKeysDBFilename = "lostkeys.db"

func main() {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	lostKeysDirpath := path.Join(homeDir, LostKeysDirname)

	if err := os.MkdirAll(lostKeysDirpath, 0777); err != nil {
		log.Fatal(err)
	}

	dbFilepath := path.Join(lostKeysDirpath, LostKeysDBFilename)

	contextProvider := &command_context.CommandContextProvider{
		DBFactory: func() *bolt.DB {
			db, err := bolt.Open(dbFilepath, 0600, nil)

			if err != nil {
				log.Fatal(fmt.Errorf("db open %w", err))
			}

			return db
		},
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "[key] [value]",
				Args:    true,
				Action:  contextProvider.Wraps(commands.Set),
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "[key]",
				Args:    true,
				Action:  contextProvider.Wraps(commands.Get),
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list all key/value pairs",
				Action:  contextProvider.Wraps(commands.List),
			},
			{
				Name:      "clipboard",
				Aliases:   []string{"cb"},
				Usage:     "[key]",
				UsageText: "copy key value to clipboard",
				Action:    contextProvider.Wraps(commands.Clipboard),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

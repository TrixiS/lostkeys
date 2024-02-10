package command_context

import (
	"github.com/boltdb/bolt"
	"github.com/urfave/cli/v2"
)

type CommandContextProvider struct {
	DBFactory func() *bolt.DB
}

func (provider *CommandContextProvider) Wraps(f CommandFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		commandContext := CommandContext{
			Provider:   provider,
			CLIContext: ctx,
		}

		return f(&commandContext)
	}
}

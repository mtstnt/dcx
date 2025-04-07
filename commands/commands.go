package commands

import (
	"github.com/urfave/cli/v3"
)

func GetCommands() []*cli.Command {
	return []*cli.Command{
		Start,
	}
}

package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// This is the command to start the Master & Worker nodes.
var Start = &cli.Command{
	Name:  "start",
	Usage: "start the Master & Worker nodes",
	Commands: []*cli.Command{
		{
			Name:   "master",
			Usage:  "start the Master node",
			Action: onStartMaster,
		},
		{
			Name:   "worker",
			Usage:  "start the Worker node",
			Action: onStartWorker,
		},
	},
}

func onStartMaster(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Starting Master node...")
	return nil
}

func onStartWorker(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Starting Worker node...")
	return nil
}

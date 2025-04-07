package main

import (
	"context"
	"log"
	"os"

	"github.com/mtstnt/launch/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "dcx",
		Usage:       "Stuff works",
		Description: "lorem ipsum dolor sit amet consectetur adipiscing elit",
		Commands:    commands.GetCommands(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

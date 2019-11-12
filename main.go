package main

import (
	"fmt"
	"github.com/tmazeika/train-planner/cmd"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "train-planner"
	app.Usage = "Amtrak train planner"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:      "list",
			Aliases:   []string{"l"},
			ArgsUsage: "<from station> <to station> [date]",
			Usage: "Lists all trains between two stations for a given date " +
					"(format '1/13/19') or today if empty",
			Action: cmd.List,
		},
		{
			Name:      "plan",
			Aliases:   []string{"p"},
			ArgsUsage: "<from station> <to station> [date]",
			Usage: "Show journey plans between two stations for a given date " +
					"(format '1/13/19') or today if empty",
			Action: cmd.Plan,
		},
		{
			Name:      "save",
			Aliases:   []string{"s"},
			ArgsUsage: "<from station> <to station> [date]",
			Usage: "Saves the HTML of the fetched page for a trip between two " +
					"stations for a given date (format '1/13/19') or today if empty",
			Action: cmd.Save,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

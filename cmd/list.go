package cmd

import (
	"fmt"
	"github.com/tmazeika/train-planner/web"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	query, err := web.NewQuery(c)
	if err != nil {
		return err
	}

	trains, err := web.Scrape(query)
	if err != nil {
		return err
	}

	fmt.Printf("Trains from %s to %s on %s:\n",
		query.FromStation, query.ToStation, query.DateStr())
	for _, train := range trains {
		fmt.Printf("- [%s -> %s] %s\n", train.FromTimeStr(), train.ToTimeStr(), train.Name)
	}
	return nil
}

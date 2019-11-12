package cmd

import (
	"fmt"
	"github.com/tmazeika/train-planner/fetch"
	"github.com/urfave/cli"
	"sort"
)

func List(c *cli.Context) error {
	q, trains, err := fetch.FromArgs(c)
	if err != nil {
		return err
	}
	sort.Sort(fetch.SortByFromTime(trains))

	fmt.Printf("Trains from %s to %s on %s:\n",
		q.FromStation, q.ToStation, q.DateStr())
	for _, train := range trains {
		fmt.Printf("- [%s -> %s] %s\n", train.FromTimeStr(), train.ToTimeStr(), train.Name)
	}
	return nil
}

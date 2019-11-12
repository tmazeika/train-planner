package cmd

import (
	"github.com/tmazeika/train-planner/fetch"
	"github.com/urfave/cli"
	"os"
)

func Save(c *cli.Context) error {
	query, err := fetch.NewQuery(c)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("trains.html", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	return fetch.RawScrapeTo(file, query)
}

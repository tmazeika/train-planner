package cmd

import (
	"github.com/tmazeika/train-planner/web"
	"github.com/urfave/cli"
	"os"
)

func Save(c *cli.Context) error {
	query, err := web.NewQuery(c)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("trains.html", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	return web.RawScrape(file, query)
}

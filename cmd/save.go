package cmd

import (
	"github.com/tmazeika/train-planner/fetch"
	"github.com/urfave/cli"
	"os"
)

func Save(c *cli.Context) error {
	file, err := os.OpenFile("trains.html", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	return fetch.Raw(file, c)
}

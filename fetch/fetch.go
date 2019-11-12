package fetch

import (
	"github.com/urfave/cli"
	"io"
	"time"
)

func From(fromStation, toStation string, when time.Time) (Query, []Train, error) {
	q := Query{
		FromStation: fromStation,
		ToStation:   toStation,
		When:        when,
	}
	trains, err := fetch(q)
	return q, trains, err
}

func FromArgs(c *cli.Context) (Query, []Train, error) {
	q, err := newQueryFromArgs(c)
	if err != nil {
		return Query{}, nil, err
	}
	trains, err := fetch(q)
	return q, trains, err
}

func Raw(dst io.Writer, c *cli.Context) error {
	q, err := newQueryFromArgs(c)
	if err != nil {
		return err
	}
	return rawScrapeTo(dst, q)
}

func fetch(q Query) ([]Train, error) {
	trains, err := getCached(q)
	if err == cacheColdMissErr {
		trains, err = scrape(q)
	}
	if err != nil {
		return nil, err
	}
	return trains, nil
}

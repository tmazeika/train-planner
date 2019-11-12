package fetch

import (
	"errors"
	"github.com/urfave/cli"
	"strings"
	"time"
)

type Query struct {
	FromStation string
	ToStation   string
	When        time.Time
}

func (q Query) DateStr() string {
	return q.When.Format("Mon Jan 2, 2006")
}

func NewQuery(c *cli.Context) (*Query, error) {
	fromStation := strings.ToUpper(c.Args().Get(0))
	toStation := strings.ToUpper(c.Args().Get(1))
	date := c.Args().Get(2)
	if fromStation == "" {
		return nil, errors.New("missing <from station>")
	}
	if toStation == "" {
		return nil, errors.New("missing <to station>")
	}
	if date == "" {
		date = time.Now().Format("1/2/06")
	}
	when, err := time.Parse("1/2/06", date)
	if err != nil {
		return nil, errors.New("malformed [date]: format is '1/13/19'")
	}
	return &Query{
		FromStation: fromStation,
		ToStation:   toStation,
		When:        when,
	}, nil
}

func (q Query) newTrain(name string, fromTime time.Time, duration time.Duration) *Train {
	return &Train{
		Name:        name,
		FromStation: q.FromStation,
		FromTime:    fromTime,
		ToStation:   q.ToStation,
		Duration:    duration,
	}
}

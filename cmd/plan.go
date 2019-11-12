package cmd

import (
	"errors"
	"fmt"
	"github.com/tmazeika/train-planner/fetch"
	"github.com/urfave/cli"
	"sort"
)

type journey struct {
	originTrain fetch.Train
	train fetch.Train
}

func allJourneys(originTrains []fetch.Train, trains []fetch.Train) []journey {
	var journeys []journey
	originLoop:
	for _, origin := range originTrains {
		for _, t := range trains {
			if origin.Name == t.Name {
				journeys = append(journeys, journey{
					originTrain: origin,
					train:       t,
				})
				continue originLoop
			}
		}
	}
	return journeys
}

func Plan(c *cli.Context) error {
	q, trains, err := fetch.FromArgs(c)
	if err != nil {
		return err
	}

	if q.FromStation != "BBY" && q.FromStation != "PHL" {
		return errors.New("unimplemented <from station> " + q.FromStation)
	}

	var originTrains []fetch.Train
	switch q.FromStation {
	case "BBY":
		_, originTrains, err = fetch.From("BOS", q.ToStation, q.When)
	case "PHL":
		_, originTrains, err = fetch.From("WAS", q.ToStation, q.When)
	}
	if err != nil {
		return err
	}

	journeys := allJourneys(originTrains, trains)
	sort.Sort(sortByFromTime(journeys))

	fmt.Printf("Journeys from %s to %s on %s:\n\n",
		q.FromStation, q.ToStation, q.DateStr())
	for _, journey := range journeys {
		fmt.Printf(" - [%s -> %s] [%s] %s\n\t(priority tickets from %s available %.0f minutes before %s)\n\n",
			journey.train.FromTimeStr(),
			journey.train.ToTimeStr(),
			journey.train.DurationStr(),
			journey.train.Name,
			journey.originTrain.FromStation,
			journey.train.FromTime.Sub(journey.originTrain.FromTime).Minutes() + 60,
			journey.train.FromStation)
	}
	return nil
}

type sortByFromTime []journey

func (s sortByFromTime) Len() int {
	return len(s)
}
func (s sortByFromTime) Less(i, j int) bool {
	return s[i].train.FromTime.Before(s[j].train.FromTime)
}
func (s sortByFromTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

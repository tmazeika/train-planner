package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type train struct {
	name        string
	fromStation string
	fromTime    time.Time
	toStation   string
	duration    time.Duration
}

func (t train) toTime() time.Time {
	return t.fromTime.Add(t.duration)
}

func (t train) fromTimeStr() string {
	return t.fromTime.Format("03:04pm")
}

func (t train) toTimeStr() string {
	return t.toTime().Format("03:04pm")
}

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
			Usage:     "Lists all trains between two stations for a given date " +
				"(format '1/13/19') or today if empty",
			Action: func(c *cli.Context) error {
				fromStation := strings.ToUpper(c.Args().Get(0))
				toStation := strings.ToUpper(c.Args().Get(1))
				date := c.Args().Get(2)
				if fromStation == "" {
					return errors.New("missing <from station>")
				}
				if toStation == "" {
					return errors.New("missing <to station>")
				}
				if date == "" {
					date = time.Now().Format("1/2/06")
				}
				when, err := time.Parse("1/2/06", date)
				if err != nil {
					return errors.New("malformed [date]: format is '1/13/19'")
				}

				body, err := getHtml(when, fromStation, toStation)
				if err != nil {
					return fmt.Errorf("failed to get HTML: %s", err)
				}

				trains, err := htmlToTrains(body, fromStation, toStation)
				if err != nil {
					return fmt.Errorf("failed to extract trains from HTML: %s", err)
				}

				fmt.Printf("Trains from %s to %s on %s:\n", fromStation, toStation, date)
				for _, train := range trains {
					fmt.Printf("- [%s -> %s] %s\n", train.fromTimeStr(), train.toTimeStr(), train.name)
				}
				return nil
			},
		},
	}
	// app.Flags = []cli.Flag {
	// 	cli.StringFlag{
	// 		Name: "after",
	// 		Usage: "only find trains after `TIME` (format '3:04pm')",
	// 	},
	// }

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func getHtml(when time.Time, fromStation, toStation string) (io.ReadCloser, error) {
	date := when.Format("01/02/2006")
	form := url.Values{}

	form.Add("wdf_origin", fromStation)
	form.Add("wdf_destination", toStation)
	form.Add("/sessionWorkflow/productWorkflow[@product='Rail']/tripRequirements/journeyRequirements[1]/departDate.usdate", date)
	form.Add("xwdf_person_type1", "/sessionWorkflow/productWorkflow[@product='Rail']/tripRequirements/allJourneyRequirements/person[1]/personType")
	form.Add("wdf_person_type1", "Adult")
	form.Add("_handler=amtrak.presentation.handler.request.rail.farefamilies.AmtrakRailFareFamiliesSearchRequestHandler/_xpath=/sessionWorkflow/productWorkflow[@product='Rail'].x", "62")
	form.Add("xwdf_origin", "/sessionWorkflow/productWorkflow[@product='Rail']/travelSelection/journeySelection[1]/departLocation/search")
	form.Add("xwdf_destination", "/sessionWorkflow/productWorkflow[@product='Rail']/travelSelection/journeySelection[1]/arriveLocation/search")

	req, err := http.NewRequest(http.MethodPost,
		"https://tickets.amtrak.com/itd/amtrak", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US, en;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie",
		"ADRUM=s=1573183464194&r=https%3A%2F%2Ftickets.amtrak.com%2Fitd%2Famtrak%3F0")
	req.Header.Set("Referer", "https://www.amtrak.com/home.html")
	req.Header.Set("User-Agent", "train-planner")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func htmlToTrains(body io.ReadCloser, fromStation, toStation string) ([]train, error) {
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var trains []train
	doc.Find(".newFareFamilyTable:not(:has(.transfer_copy)) td:has(#_depart_time_span)").
			Each(func(i int, s *goquery.Selection) {
				name := s.Find("#_service_span").Text()
				durationStr := s.Find("#_duration_span").Text()
				fromTimeStr := s.Find("#_depart_time_span").Text()

				if name == "" {
					log.Println("found empty name")
					return
				}

				duration, err := time.ParseDuration(strings.ReplaceAll(durationStr, " ", ""))
				if err != nil {
					log.Printf("found invalid duration '%s': %s\n", durationStr, err)
					return
				}

				fromTime, err := time.Parse("2006-01-02T15:04:05.999-07:00", fromTimeStr)
				if err != nil {
					log.Printf("found invalid departure time '%s': %s\n", fromTimeStr, err)
					return
				}

				trains = append(trains, train{
					name:        name,
					fromStation: fromStation,
					fromTime:    fromTime,
					toStation:   toStation,
					duration:    duration,
				})
			})
	return trains, nil
}

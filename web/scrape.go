package web

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RawScrape(dst io.Writer, query *Query) error {
	reader, err := getHtml(query)
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(dst, reader)
	return err
}

func Scrape(query *Query) ([]*Train, error) {
	reader, err := getHtml(query)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return htmlToTrains(reader, query)
}

func getHtml(query *Query) (io.ReadCloser, error) {
	date := query.When.Format("01/02/2006")
	form := url.Values{}

	form.Add("wdf_origin", query.FromStation)
	form.Add("wdf_destination", query.ToStation)
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

func htmlToTrains(body io.Reader, query *Query) ([]*Train, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var trains []*Train
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

				trains = append(trains, query.newTrain(name, fromTime, duration))
			})
	return trains, nil
}

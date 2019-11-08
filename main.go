package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	html := getHtml(time.Now().Add(time.Hour), "PHL", "BBY")

	// b, err := ioutil.ReadAll(html)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(b))
	// return

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		panic(err)
	}

	fmt.Println("Today's trains:")
	doc.Find(".newFareFamilyTable:not(:has(.transfer_copy)) td:has(#_depart_time_span)").
			Each(func(i int, s *goquery.Selection) {
				train := s.Find("#_service_span").Text()
				duration, err := time.ParseDuration(strings.ReplaceAll(s.Find("#_duration_span").Text(), " ", ""))
				if err != nil {
					panic(err)
				}
				from, err := time.Parse("2006-01-02T15:04:05.999-07:00",
					s.Find("#_depart_time_span").Text())
				if err != nil {
					panic(err)
				}
				fmt.Printf("- [%s -> %s] %s\n", from.Format("03:04pm"), from.Add(duration).Format("03:04pm"), train)
			})
}

func getHtml(when time.Time, from, to string) io.ReadCloser {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	knownStations := map[string]string{
		"BBY": "Boston, MA",
		"BOS": "Boston, MA",
		"PHL": "Philadelphia, PA",
		"WAS": "Washington, DC",
	}
	date := when.Format("01/02/2006")

	if _, ok := knownStations[from]; !ok {
		panic("invalid from station: " + from)
	}
	if _, ok := knownStations[to]; !ok {
		panic("invalid to station: " + from)
	}

	form := url.Values{}
	form.Add("wdf_origin", from)
	form.Add("wdf_origin", knownStations[from])
	form.Add("wdf_destination", to)
	form.Add("wdf_destination", knownStations[to])
	form.Add("/sessionWorkflow/productWorkflow[@product='Rail']/tripRequirements/journeyRequirements[1]/departDate.usdate", date)
	form.Add("xwdf_person_type1", "/sessionWorkflow/productWorkflow[@product=\"Rail\"]/tripRequirements/allJourneyRequirements/person[1]/personType")
	form.Add("wdf_person_type1", "Adult")
	form.Add("_handler=amtrak.presentation.handler.request.rail.farefamilies.AmtrakRailFareFamiliesSearchRequestHandler/_xpath=/sessionWorkflow/productWorkflow[@product='Rail'].x", "62")
	form.Add("xwdf_origin", "/sessionWorkflow/productWorkflow[@product='Rail']/travelSelection/journeySelection[1]/departLocation/search")
	form.Add("xwdf_destination", "/sessionWorkflow/productWorkflow[@product='Rail']/travelSelection/journeySelection[1]/arriveLocation/search")

	req, err := http.NewRequest(http.MethodPost, "https://tickets.amtrak.com/itd/amtrak", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US, en;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "ADRUM=s=1573183464194&r=https%3A%2F%2Ftickets.amtrak.com%2Fitd%2Famtrak%3F0")
	req.Header.Set("Referer", "https://www.amtrak.com/home.html")
	req.Header.Set("User-Agent", "train-planner")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp.Body
}

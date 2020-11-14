package trains

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func cleanResult(s string) (clean string) {
	s = strings.Replace(s, "\n", "", -1)
	clean = strings.Trim(s, " ")
	return
}

func getStation(doc *goquery.Document) *string {
	station := doc.Find(`h2 span.to`).Text()
	return &station
}

func getDepartures(doc *goquery.Document) (Departures, error) {
	table := doc.Find(`div.results.trains div.tbl-cont tr`)
	depts := Departures{}
	if table.Length() == 0 {
		return nil, fmt.Errorf("no results found")
	}
	for i := range table.Nodes {
		item := table.Eq(i)
		columns := item.Find(`td`)
		dept := Departure{}
		for col := range columns.Nodes {
			selCol := columns.Eq(col)
			switch col {
			case 0:
				dept.Due = cleanResult(selCol.Text())
			case 1:
				dept.Dest = cleanResult(selCol.Text())
			case 2:
				dept.Status = cleanResult(selCol.Text())
			case 3:
				dept.Platform = cleanResult(selCol.Text())
			}
		}
		if dept.Dest == "" && dept.Due == "" {
			continue
		}
		depts = append(depts, &dept)
	}
	//fmt.Println(depts[2].Dest)
	//os.Exit(1)
	return depts, nil
}

func getResponse(code, destination string) (*http.Response, error) {
	switch destination {
	case "":
		return http.Get(fmt.Sprintf("http://ojp.nationalrail.co.uk/service/ldbboard/dep/%s/", code))
	default:
		return http.Get(fmt.Sprintf("http://ojp.nationalrail.co.uk/service/ldbboard/dep/%s/%s/To", code, destination))
	}
}

// StationInfo gets the station info for a station and destination
func StationInfo(code string, destination string) (*StationDepartures, error) {
	resp, err := getResponse(code, destination)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	stationName := getStation(doc)
	departures, _ := getDepartures(doc)

	sd := StationDepartures{
		Name:       *stationName,
		Departures: departures,
	}

	return &sd, nil
}

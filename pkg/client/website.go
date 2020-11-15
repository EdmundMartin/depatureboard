package client

import (
	"fmt"
	"net/http"
	"strings"

	t "github.com/EdmundMartin/depatureboard/pkg/trains"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// WebsiteClient client to access the website data
type WebsiteClient struct{}

// NewWebsiteClient returns an instance of website client
func NewWebsiteClient() Client {
	return WebsiteClient{}
}

// CleanString cleans a scraped string to a straight string
func CleanString(s string) (clean string) {
	clean = strings.Replace(s, "\n", "", -1)
	clean = strings.Trim(clean, " ")
	clean = strings.Replace(clean, "\t", "", -1)
	clean = html.UnescapeString(clean)
	clean = strings.Join(strings.Fields(clean), " ")
	return
}

func getStation(doc *goquery.Document) *string {
	station := doc.Find(`h2 span.to`).Text()
	return &station
}

func getDepartures(doc *goquery.Document) (t.Departures, error) {
	table := doc.Find(`div.results.trains div.tbl-cont tr`)
	depts := t.Departures{}
	if table.Length() == 0 {
		return nil, fmt.Errorf("no results found")
	}
	for i := range table.Nodes {
		item := table.Eq(i)
		columns := item.Find(`td`)
		dept := t.Departure{}
		for col := range columns.Nodes {
			selCol := columns.Eq(col)
			switch col {
			case 0:
				dept.Due = CleanString(selCol.Text())
			case 1:
				dept.Dest = CleanString(selCol.Text())
			case 2:
				dept.Status = CleanString(selCol.Text())
			case 3:
				dept.Platform = CleanString(selCol.Text())
			}
		}
		if dept.Dest == "" && dept.Due == "" {
			continue
		}
		depts = append(depts, &dept)
	}
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
func (c WebsiteClient) StationInfo(code string, destination string) (*t.StationDepartures, error) {
	resp, err := getResponse(code, destination)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	stationName := getStation(doc)
	departures, err := getDepartures(doc)
	if err != nil {
		return nil, err
	}

	sd := t.StationDepartures{
		Name:       *stationName,
		Departures: departures,
	}

	return &sd, nil
}

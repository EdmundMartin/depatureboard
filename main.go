package main

import (
	"flag"
	"log"
	"time"

	client "github.com/EdmundMartin/depatureboard/pkg/client"
	t "github.com/EdmundMartin/depatureboard/pkg/trains"
)

func poll(station, destination string, refresh int, departures chan *t.StationDepartures) {
	c := client.NewWebsiteClient()
	for {
		res, err := c.StationInfo(station, destination)
		if err == nil {
			departures <- res
		}
		time.Sleep(time.Duration(refresh) * time.Second)
	}
}

func main() {
	var station, destination string
	var refresh, maxResults int
	flag.StringVar(&station, "station", "", "enter Departure station code")
	flag.StringVar(&destination, "destination", "", "enter a destination station code")
	flag.IntVar(&refresh, "refresh", 30, "frequency of data refresh")
	flag.IntVar(&maxResults, "results", 30, "number of results to display")
	flag.Parse()
	if station == "" || len(station) != 3 {
		log.Fatalf("failed to provide valid station code")
	}
	if destination != "" && len(destination) != 3 {
		log.Fatalf("failed to enter a valid destination code")
	}
	ch := make(chan *t.StationDepartures)
	go poll(station, destination, refresh, ch)
	for {
		res := <-ch
		res.Departures = res.Departures.Limit(maxResults)
		res.Display()
	}
}

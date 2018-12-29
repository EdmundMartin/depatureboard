package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type depature struct {
	Due      string
	Dest     string
	Status   string
	Platform string
}

func cleanResult(s string) (clean string) {
	s = strings.Replace(s, "\n", "", -1)
	clean = strings.Trim(s, " ")
	return
}

func getDepatures(doc *goquery.Document) ([]*depature, error) {
	table := doc.Find(`div.results.trains div.tbl-cont tr`)
	depts := []*depature{}
	if table.Length() == 0 {
		return nil, fmt.Errorf("no results found")
	}
	for i := range table.Nodes {
		item := table.Eq(i)
		columns := item.Find(`td`)
		dept := &depature{}
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
		depts = append(depts, dept)
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

func stationInfo(code, destination string) ([]*depature, error) {
	resp, err := getResponse(code, destination)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	return getDepatures(doc)
}

func clearScreen(std *os.File) *exec.Cmd {
	runtime := runtime.GOOS
	if runtime == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = std
		return cmd
	} else if runtime == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		return cmd
	}
	return nil
}

func writeOutput(w *tabwriter.Writer, depts []*depature, maxResults int, cmd *exec.Cmd) {
	format := "%s\t%s\t%s\t%s\t"
	fmt.Fprintln(w, fmt.Sprintf(format, "Due", "Destination", "Status", "Platform"))
	for i, dept := range depts {
		if i < maxResults {
			fmt.Fprintln(w, fmt.Sprintf(format, dept.Due, dept.Dest, dept.Status, dept.Platform))
		} else {
			break
		}
	}
	if cmd != nil {
		cmd.Run()
	}
	w.Flush()
}

func poll(station, destination string, refresh int, depatures chan []*depature) {
	for {
		res, err := stationInfo(station, destination)
		if err == nil {
			depatures <- res
		}
		time.Sleep(time.Duration(refresh) * time.Minute)
	}
}

func main() {
	var station, destination string
	var refresh, maxResults int
	flag.StringVar(&station, "station", "", "enter depature station code")
	flag.StringVar(&destination, "destination", "", "enter a destination station code")
	flag.IntVar(&refresh, "refresh", 1, "frequency of data refresh")
	flag.IntVar(&maxResults, "results", 30, "number of results to display")
	flag.Parse()
	if station == "" || len(station) != 3 {
		log.Fatalf("failed to provide valid station code")
	}
	if destination != "" && len(destination) != 3 {
		log.Fatalf("failed to enter a valid destination code")
	}
	ch := make(chan []*depature)
	go poll(station, destination, refresh, ch)
	std := os.Stdout
	cmd := clearScreen(std)
	w := tabwriter.NewWriter(std, 0, 0, 10, ' ', 0)
	for {
		res := <-ch
		writeOutput(w, res, maxResults, cmd)
	}
}

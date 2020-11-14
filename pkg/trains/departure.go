package trains

import (
	"strings"
	"time"

	tm "github.com/buger/goterm"
	tw "github.com/olekukonko/tablewriter"
)

type StationDepartures struct {
	Departures
	Name string
}

// Departures a slice of Departure
type Departures []*Departure

// Departure definition of a departure
type Departure struct {
	Due      string
	Dest     string
	Status   string
	Platform string
}

// Limit reduce the amount of departures
func (deps Departures) Limit(limit int) Departures {
	if (len(deps)) > limit {
		return deps[:limit]
	}
	return deps
}

// Display writes a table of departures to terminal
func (deps *Departures) Display() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Printf("Last updated at: %v\n\n", time.Now().Format("15:04:05"))
	t := &strings.Builder{}
	table := tw.NewWriter(t)
	table.SetHeader([]string{"Due", "Destination", "Status", "Platform"})
	table.SetBorder(false)
	cNormal := []tw.Colors{{}, {tw.Normal, tw.FgCyanColor}, {tw.Bold, tw.FgWhiteColor}, {}}
	cLate := []tw.Colors{{}, {tw.Normal, tw.FgRedColor}, {tw.Bold, tw.FgWhiteColor}, {}}
	whiteForeRedBack := tw.Colors{tw.Bold, tw.FgWhiteColor, tw.BgRedColor}
	cCancelled := []tw.Colors{whiteForeRedBack, whiteForeRedBack, whiteForeRedBack, whiteForeRedBack}
	for _, d := range *deps {
		data := []string{d.Due, strings.ReplaceAll(d.Dest, " ", " "), d.Status, d.Platform}
		c := cNormal
		if strings.Contains(d.Status, "late") {
			c = cLate
		}
		if strings.Contains(d.Status, "Cancelled") {
			c = cCancelled
		}
		table.Rich(data, c)
	}

	table.Render()
	tm.Print(t.String())
	tm.Flush()
}

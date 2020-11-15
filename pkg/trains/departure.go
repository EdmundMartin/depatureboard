package trains

import (
	"strings"
	"time"

	tm "github.com/buger/goterm"
	tw "github.com/olekukonko/tablewriter"
)

// StationDepartures struct contianing the depatures and station name
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
func (s StationDepartures) Display() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Printf("Last updated at: %v\n\n", time.Now().Format("15:04:05"))
	tm.Printf("%v\n", strings.Repeat("-", len(s.Name)))
	tm.Printf("%v\n", s.Name)
	tm.Printf("%v\n", strings.Repeat("-", len(s.Name)))
	tm.Println()
	t := &strings.Builder{}
	table := tw.NewWriter(t)
	table.SetHeader([]string{"Due", "Destination", "Status", "Platform"})
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	cNormal := []tw.Colors{{}, {tw.Normal, tw.FgCyanColor}, {tw.Bold, tw.FgWhiteColor}, {}}
	cLate := []tw.Colors{{}, {tw.Normal, tw.FgRedColor}, {tw.Bold, tw.FgWhiteColor}, {}}
	whiteForeRedBack := tw.Colors{tw.Bold, tw.FgWhiteColor, tw.BgRedColor}
	cCancelled := []tw.Colors{whiteForeRedBack, whiteForeRedBack, whiteForeRedBack, whiteForeRedBack}
	for _, d := range s.Departures {
		data := []string{d.Due, d.Dest, d.Status, d.Platform}
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

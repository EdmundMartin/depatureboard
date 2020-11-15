package client

import t "github.com/EdmundMartin/depatureboard/pkg/trains"

// Client for accessing station info
type Client interface {
	StationInfo(string, string) (*t.StationDepartures, error)
}

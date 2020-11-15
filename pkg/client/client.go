package client

import t "github.com/EdmundMartin/depatureboard/pkg/trains"

type Client interface {
	StationInfo(string, string) (*t.StationDepartures, error)
}

package status

import (
	"time"
)

// Client ...
type Client struct {
	CommonName         string     `json:"name"`
	RealAddress        string     `json:"addr"`
	VirtualAddress     string     `json:"vip"`
	VirtualIPv6Address string     `json:"vip6,omitempty"`
	BytesReceived      int        `json:"recv"`
	BytesSent          int        `json:"sent"`
	ConnectedSince     *time.Time `json:"since,omitempty"`
	ConnectedUTC       int        `json:"stamp,omitempty"`
	Username           string     `json:"-"`
	ClientID           int        `json:"-"`
	PeerID             int        `json:"-"`
}

// Routing ...
type Routing struct {
	VirtualAddress string     `json:"vaddr"`
	CommonName     string     `json:"name"`
	RealAddress    string     `json:"addr"`
	LastRef        *time.Time `json:"lastRef,omitempty"`
	LastRefUTC     int        `json:"lastStamp,omitempty"`
}

// GlobalStats ...
type GlobalStats struct {
	MaxBcastMcastQueueLength int `json:"maxBMQueueLen,omitempty"`
}

// Status ...
type Status struct {
	Title        string      `json:"title,omitempty"`
	ClientList   []Client    `json:"clients"`
	RoutingTable []Routing   `json:"routings"`
	GlobalStats  GlobalStats `json:"stats"`
	Result       string      `json:"result,omitempty"`
	Label        string      `json:"label"`
}

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
	MaxBcastMcastQueueLength int
}

// Status ...
type Status struct {
	Title        string      `json:"title"`
	ClientList   []Client    `json:"clients"`
	RoutingTable []Routing   `json:"routings"`
	GlobalStats  GlobalStats `json:"stats"`
	Result       string      `json:"result,omitempty"`
}

var clientListHeaderColumns = []string{
	"HEADER",
	"CLIENT_LIST",
	"Common Name",
	"Real Address",
	"Virtual Address",
	"Virtual IPv6 Address",
	"Bytes Received",
	"Bytes Sent",
	"Connected Since",
	"Connected Since (time_t)",
	"Username",
	"Client ID",
	"Peer ID",
}

var routingTableHeadersColumns = []string{
	"HEADER",
	"ROUTING_TABLE",
	"Virtual Address",
	"Common Name",
	"Real Address",
	"Last Ref",
	"Last Ref (time_t)",
}

// judge header data matched or not
const (
	clientListHeaders = 1 << iota
	routingTableHeaders
	globalStatsHeaders
)

func checkClientListHeader(headers []string) bool {
	for i, v := range headers {
		if v != clientListHeaderColumns[i] {
			return false
		}
	}
	return true
}

func checkRoutingTableHeader(headers []string) bool {
	for i, v := range headers {
		if v != routingTableHeadersColumns[i] {
			return false
		}
	}
	return true
}

func checkHeaders(headers []string) int {
	if checkClientListHeader(headers) {
		return clientListHeaders
	} else if checkRoutingTableHeader(headers) {
		return routingTableHeaders
	} else {
		return 0
	}
}

// Time Parse
const dateLayout = "Mon Jan _2 15:04:05 2006"

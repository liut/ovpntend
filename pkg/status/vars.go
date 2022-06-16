package status

// Time Parse
const dateLayout = "2006-01-02 15:04:05"

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
	if len(headers) < len(clientListHeaderColumns) {
		return false
	}
	for i, v := range clientListHeaderColumns {
		if v != headers[i] {
			return false
		}
	}
	return true
}

func checkRoutingTableHeader(headers []string) bool {
	if len(headers) < len(routingTableHeadersColumns) {
		return false
	}
	for i, v := range routingTableHeadersColumns {
		if v != headers[i] {
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

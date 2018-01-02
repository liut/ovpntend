package status

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	CommonName         string
	RealAddress        string
	VirtualAddress     string
	VirtualIPv6Address string
	BytesReceived      int
	BytesSent          int
	ConnectedSince     time.Time
	ConnectedUTC       int
	Username           string
	ClientID           int
	PeerID             int
}

type Routing struct {
	VirtualAddress string
	CommonName     string
	RealAddress    string
	LastRef        time.Time
	LastRefUTC     int
}

type GlobalStats struct {
	MaxBcastMcastQueueLength int
}

type Status struct {
	ClientList   []Client
	RoutingTable []Routing
	GlobalStats  GlobalStats
	Result       string
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

const (
	MaxLen              = 32768 // 32*1024
	firstSentenceMinLen = 75
	stepMinReadLen      = 256
)

// ParseAddr func parse OpenVPN `status 2` from connect write and return status struct、error
func ParseAddr(management string) (*Status, error) {
	conn, err := net.DialTimeout("tcp", management, 2*time.Second)
	if err != nil {
		return &Status{Result: "connect to open server false"}, err
	}

	// read first sentence
	buf := make([]byte, firstSentenceMinLen)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return &Status{Result: "conn read false"}, err
	}

	//write `status 2`
	var size int
	size, err = conn.Write([]byte("status 2\n"))
	if err != nil {
		return &Status{Result: "connect write false"}, err
	}
	log.Printf("wrote: %d bytes", size)

	// read `status 2` result
	buf = make([]byte, MaxLen)
	if _, err := io.ReadAtLeast(conn, buf, stepMinReadLen); err != nil {
		return &Status{Result: "connect client message read false"}, err
	}

	out := bytes.NewBuffer(buf)
	defer conn.Close()
	return Parse(out)
}

//ParseFile func parse OpenVPN status from log file and return status struct、error
func ParseFile(file string) (*Status, error) {
	rd, err := os.Open(file)
	if err != nil {
		return &Status{Result: "open false"}, err
	}
	defer rd.Close()
	return Parse(rd)
}

//Parse func parse OpenVPN status from io.Reader and return status struct、error
func Parse(rd io.Reader) (*Status, error) {
	scanner := bufio.NewScanner(rd)
	scanner.Split(bufio.ScanLines)

	var (
		err                   error
		clients               []Client
		routingTable          []Routing
		maxBcastMcastQueueLen int
		isEmpty               bool
		judgeFileType         int
	)
	judgeFileType = 0
	isEmpty = true

	for scanner.Scan() {
		isEmpty = false

		var (
			ct time.Time
			rt time.Time
		)

		fields := strings.Split(scanner.Text(), ",")
		if fields[0] == "TITLE" {

		} else if fields[0] == "TIME" {

		} else if fields[0] == "" {

		} else if checkHeaders(fields) == clientListHeaders {
			judgeFileType = clientListHeaders
		} else if judgeFileType == clientListHeaders && len(fields) == len(clientListHeaderColumns)-1 {
			ct, _ = time.Parse(dateLayout, fields[7])
			clients = append(clients, Client{
				fields[1], fields[2], fields[3], fields[4], Atoi(fields[5]), Atoi(fields[6]), ct, Atoi(fields[8]), fields[9], Atoi(fields[10]), Atoi(fields[11])})
		} else if fields[0] == "" {

		} else if checkHeaders(fields) == routingTableHeaders {
			judgeFileType = routingTableHeaders
		} else if judgeFileType == routingTableHeaders && len(fields) == len(routingTableHeadersColumns)-1 {
			rt, _ = time.Parse(dateLayout, fields[4])
			routingTable = append(routingTable, Routing{fields[1], fields[2], fields[3], rt, Atoi(fields[5])})
		} else if fields[0] == "GLOBAL_STATS" {
			if fields[1] == "Max bcast/mcast queue length" {
				i, err := strconv.Atoi(fields[2])
				if err == nil {
					maxBcastMcastQueueLen = i
				}
			}
		} else if fields[0] == "END" {
			if len(fields) == 1 {
				break
			}
		} else {
			return &Status{Result: "Unable to Parse Status file"}, err
		}
	}

	if isEmpty {
		return &Status{Result: "file is empty"}, err
	}

	return &Status{
		ClientList:   clients,
		RoutingTable: routingTable,
		GlobalStats:  GlobalStats{maxBcastMcastQueueLen},
		Result:       "OK",
	}, nil
}

func Atoi(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("Type transform err: %v", err)
	}
	return i
}

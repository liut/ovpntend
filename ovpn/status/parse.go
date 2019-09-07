package status

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

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
				fields[1], fields[2], fields[3], fields[4], Atoi(fields[5]), Atoi(fields[6]), &ct, Atoi(fields[8]), fields[9], Atoi(fields[10]), Atoi(fields[11])})
		} else if fields[0] == "" {

		} else if checkHeaders(fields) == routingTableHeaders {
			judgeFileType = routingTableHeaders
		} else if judgeFileType == routingTableHeaders && len(fields) == len(routingTableHeadersColumns)-1 {
			rt, _ = time.Parse(dateLayout, fields[4])
			routingTable = append(routingTable, Routing{fields[1], fields[2], fields[3], &rt, Atoi(fields[5])})
		} else if fields[0] == "GLOBAL_STATS" {
			if fields[1] == "Max bcast/mcast queue length" {
				var i int
				i, err = strconv.Atoi(fields[2])
				if err == nil {
					maxBcastMcastQueueLen = i
				} else {
					log.Printf("strconv ERR %s", err)
				}
			}
		} else if fields[0] == "END" {
			if len(fields) == 1 {
				break
			}
		} else {
			return &Status{Result: "Unable to Parse Status "}, err
		}
	}

	if isEmpty {
		return &Status{Result: "data is empty"}, err
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
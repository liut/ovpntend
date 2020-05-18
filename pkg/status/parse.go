package status

import (
	"bufio"
	"io"
	"net"
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
		title                 string
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
			err error
			ct  time.Time
			rt  time.Time
		)

		fields := strings.Split(scanner.Text(), ",")
		logger().Debugw("got line", "fields", fields)
		if fields[0] == "TITLE" {
			title = fields[1]
		} else if fields[0] == "TIME" {

		} else if fields[0] == "" {

		} else if checkHeaders(fields) == clientListHeaders {
			judgeFileType = clientListHeaders
		} else if judgeFileType == clientListHeaders && len(fields) == len(clientListHeaderColumns)-1 {
			ct, err = time.ParseInLocation(dateLayout, fields[7], time.Local)
			if err != nil {
				logger().Infow("parse time fail", "err", err)
			} else {
				logger().Debugw("parsed", "t", ct, "since", time.Since(ct))
			}
			host, port, _ := net.SplitHostPort(fields[2])
			clients = append(clients, Client{
				fields[1], HostPort{host, port},
				fields[3], fields[4], Atoi(fields[5]),
				Atoi(fields[6]), &ct, Atoi(fields[8]),
				fields[9], Atoi(fields[10]), Atoi(fields[11])})
		} else if fields[0] == "" {

		} else if checkHeaders(fields) == routingTableHeaders {
			judgeFileType = routingTableHeaders
		} else if judgeFileType == routingTableHeaders && len(fields) == len(routingTableHeadersColumns)-1 {
			rt, err = time.ParseInLocation(dateLayout, fields[4], time.Local)
			if err != nil {
				logger().Infow("parse time fail", "err", err)
			}
			routingTable = append(routingTable, Routing{fields[1], fields[2], fields[3], &rt, Atoi(fields[5])})
		} else if fields[0] == "GLOBAL_STATS" {
			if fields[1] == "Max bcast/mcast queue length" {
				var i int
				i, err = strconv.Atoi(fields[2])
				if err == nil {
					maxBcastMcastQueueLen = i
				} else {
					logger().Infow("strconv fail", "err", err)
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
		Title:        title,
		ClientList:   clients,
		RoutingTable: routingTable,
		GlobalStats:  GlobalStats{maxBcastMcastQueueLen},
		Result:       "OK",
	}, nil
}

func Atoi(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		logger().Infow("type transform fail", "err", err)
	}
	return i
}

package openVPNstatus

import (
	"bufio"
	"io/ioutil"
	"os"
	"time"
)

// ,,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID

// kc-liut,103.43.184.114:19269,10.168.252.167,pv6,218578,473516,Thu Dec 28 10:20:49 2017,1514427649,UNDEF,41,1
type Blob []byte

// cur := time.UTC
// timestamp := cur.UnixNano() / 1000000 //UnitNano获取的是纳秒，除以1000000获取毫秒级的时间戳

type Client struct {
	CommonName         string
	RealAddress        string
	VirtualAddress     string
	VirtualIPv6Address string
	BytesReceived      Blob
	BytesSent          Blob
	ConnectedSince     time.Time
	ConnectedUTC       time.UTC
	Username           string
	ClientID           int
	PeerID             int
}

// HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
// ROUTING_TABLE,b6:f6:18:08:da:4b,kc-box-00006,221.223.54.60:45505,Thu Dec 28 15:40:41 2017,1514446841

type Routing struct {
	VirtualAddress string
	CommonName     string
	RealAddress    string
	LastRef        time.Time
	LastRefUTC     time.UTC
}

// GLOBAL_STATS,Max bcast/mcast queue length,4

type GlobalStats struct {
	MaxBcastMcastQueueLength int
}

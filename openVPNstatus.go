package openVPNstatus

import (
	"bufio"
	"io/ioutil"
	"os"
)

HEADER,CLIENT_LIST,,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID

CLIENT_LIST,kc-liut,103.43.184.114:19269,10.168.252.167,,218578,473516,Thu Dec 28 10:20:49 2017,1514427649,UNDEF,41,1

type Client struct{
	CommonName string
	
}
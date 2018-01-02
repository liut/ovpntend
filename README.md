# go-openvpn-status

## Usage

```go
status, _ := status.ParseFile("examples/log_status.txt")

fmt.Println(status.RoutingTable)

fmt.Println(status.ClientList)

fmt.Println(status.GlobalStats)

```

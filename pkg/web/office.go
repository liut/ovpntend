package web

import (
	"log"
	"net"
	"sort"

	"github.com/liut/ovpntend/pkg/settings"
)

type officeEntry struct {
	single net.IP
	cidr   *net.IPNet
	label  string
}

var officeEntries []officeEntry

func init() {
	ReloadOfficeIPs()
}

// ReloadOfficeIPs (re)parses settings.Current.OfficeIPs into typed entries.
// Intended for production startup and test setup.
func ReloadOfficeIPs() {
	officeEntries = parseOfficeIPs(settings.Current.OfficeIPs)
}

func parseOfficeIPs(m map[string]string) []officeEntry {
	if len(m) == 0 {
		return nil
	}
	out := make([]officeEntry, 0, len(m))
	for ip, label := range m {
		if _, cidr, err := net.ParseCIDR(ip); err == nil {
			out = append(out, officeEntry{cidr: cidr, label: label})
			continue
		}
		if parsed := net.ParseIP(ip); parsed != nil {
			out = append(out, officeEntry{single: parsed, label: label})
			continue
		}
		log.Printf("office: ignoring invalid entry %q", ip)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return cidrSpecificity(out[i].cidr) > cidrSpecificity(out[j].cidr)
	})
	return out
}

func cidrSpecificity(n *net.IPNet) int {
	if n == nil {
		return -1
	}
	ones, _ := n.Mask.Size()
	return ones
}

// LookupOfficeLabel returns the configured office label for ip, if any.
func LookupOfficeLabel(ip string) (string, bool) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", false
	}
	for i := range officeEntries {
		e := &officeEntries[i]
		if e.single != nil && e.single.Equal(parsed) {
			return e.label, true
		}
		if e.cidr != nil && e.cidr.Contains(parsed) {
			return e.label, true
		}
	}
	return "", false
}

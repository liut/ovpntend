package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPlace_OfficeSingleIPMatch(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"1.2.3.4": "Beijing"})
	defer func() { officeEntries = nil }()
	assert.Equal(t, "Beijing", FindPlace("1.2.3.4"))
}

func TestFindPlace_OfficeCIDRMatch(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"10.0.0.0/24": "OfficeNet"})
	defer func() { officeEntries = nil }()
	assert.Equal(t, "OfficeNet", FindPlace("10.0.0.5"))
}

func TestFindPlace_NoMatchFallsBack(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"1.2.3.4": "Beijing"})
	defer func() { officeEntries = nil }()
	// ipip data is not loaded in test env, so miss falls through to "[未知地区]"
	assert.Equal(t, "[未知地区]", FindPlace("8.8.8.8"))
}

func TestFindPlace_EmptyConfigUnchanged(t *testing.T) {
	officeEntries = nil
	// ipip data is not loaded in test env, so behavior is identical to pre-change
	assert.Equal(t, "[未知地区]", FindPlace("8.8.8.8"))
}

func TestFindPlace_HitShortCircuitsIpip(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"8.8.8.8": "GoogleDNS"})
	defer func() { officeEntries = nil }()
	// If the hit did not short-circuit, FindPlace would call ipip.FindCity
	// which returns "" in test env, producing "[未知地区]" instead.
	assert.Equal(t, "GoogleDNS", FindPlace("8.8.8.8"))
}

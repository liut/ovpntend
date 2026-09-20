package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOfficeIPs_Empty(t *testing.T) {
	assert.Empty(t, parseOfficeIPs(nil))
	assert.Empty(t, parseOfficeIPs(map[string]string{}))
}

func TestParseOfficeIPs_SingleIP(t *testing.T) {
	entries := parseOfficeIPs(map[string]string{"1.2.3.4": "Beijing"})
	if assert.Len(t, entries, 1) {
		assert.NotNil(t, entries[0].single)
		assert.Nil(t, entries[0].cidr)
		assert.Equal(t, "Beijing", entries[0].label)
	}
}

func TestParseOfficeIPs_CIDR(t *testing.T) {
	entries := parseOfficeIPs(map[string]string{"10.0.0.0/24": "OfficeNet"})
	if assert.Len(t, entries, 1) {
		assert.Nil(t, entries[0].single)
		assert.NotNil(t, entries[0].cidr)
		assert.Equal(t, "OfficeNet", entries[0].label)
	}
}

func TestParseOfficeIPs_Invalid(t *testing.T) {
	entries := parseOfficeIPs(map[string]string{"not-an-ip": "Foo"})
	assert.Empty(t, entries)
}

func TestParseOfficeIPs_IPv6(t *testing.T) {
	entries := parseOfficeIPs(map[string]string{"::1": "Loopback"})
	if assert.Len(t, entries, 1) {
		assert.NotNil(t, entries[0].single)
	}
}

func TestParseOfficeIPs_Mixed(t *testing.T) {
	entries := parseOfficeIPs(map[string]string{
		"1.2.3.4":     "Beijing",
		"10.0.0.0/24": "OfficeNet",
		"bad-entry":   "Skip",
	})
	assert.Len(t, entries, 2)
}

func TestLookupOfficeLabel_SingleIPMatch(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"1.2.3.4": "Beijing"})
	defer func() { officeEntries = nil }()
	label, ok := LookupOfficeLabel("1.2.3.4")
	assert.True(t, ok)
	assert.Equal(t, "Beijing", label)
}

func TestLookupOfficeLabel_CIDRMatch(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"10.0.0.0/24": "OfficeNet"})
	defer func() { officeEntries = nil }()
	label, ok := LookupOfficeLabel("10.0.0.5")
	assert.True(t, ok)
	assert.Equal(t, "OfficeNet", label)
}

func TestLookupOfficeLabel_NoMatch(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"1.2.3.4": "Beijing"})
	defer func() { officeEntries = nil }()
	_, ok := LookupOfficeLabel("8.8.8.8")
	assert.False(t, ok)
}

func TestLookupOfficeLabel_EmptyConfig(t *testing.T) {
	officeEntries = nil
	_, ok := LookupOfficeLabel("1.2.3.4")
	assert.False(t, ok)
}

func TestLookupOfficeLabel_InvalidIP(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{"1.2.3.4": "Beijing"})
	defer func() { officeEntries = nil }()
	_, ok := LookupOfficeLabel("not-an-ip")
	assert.False(t, ok)
}

func TestLookupOfficeLabel_MostSpecificWinsOnOverlap(t *testing.T) {
	officeEntries = parseOfficeIPs(map[string]string{
		"10.0.0.0/16": "Outer",
		"10.0.0.0/24": "Inner",
	})
	defer func() { officeEntries = nil }()
	label, ok := LookupOfficeLabel("10.0.0.5")
	assert.True(t, ok)
	assert.Equal(t, "Inner", label)
}

func TestParseOfficeIPs_OverlappingCIDRsSortedBySpecificity(t *testing.T) {
	entries := parseOfficeIPs(map[string]string{
		"10.0.0.0/16": "Outer",
		"10.0.0.0/24": "Inner",
	})
	if assert.Len(t, entries, 2) {
		ones0, _ := entries[0].cidr.Mask.Size()
		ones1, _ := entries[1].cidr.Mask.Size()
		assert.Greater(t, ones0, ones1, "more specific CIDR first")
		assert.Equal(t, "Inner", entries[0].label)
		assert.Equal(t, "Outer", entries[1].label)
	}
}

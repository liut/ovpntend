package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadFile(t *testing.T) {
	s, _ := ParseFile("examples/badFile.txt")
	assert.Equal(t, s.Result, "Unable to Parse Status ")
}

func TestEmptyFile(t *testing.T) {
	s, _ := ParseFile("examples/emptyFile.txt")
	assert.Equal(t, s.Result, "data is empty")
}

func TestUnableParse(t *testing.T) {
	s, _ := ParseFile("examples/unableParse.txt")
	assert.Equal(t, s.Result, "Unable to Parse Status ")
}

func TestLogStatus(t *testing.T) {
	s, _ := ParseFile("examples/log_status.txt")
	assert.Equal(t, s.Result, "OK")
}

func TestOpenFalse(t *testing.T) {
	s, _ := ParseFile("examples/notExistFile.txt")
	assert.Equal(t, s.Result, "open false")
}

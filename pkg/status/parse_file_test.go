package status

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	slog.SetLogLoggerLevel(slog.LevelDebug)
	os.Exit(m.Run())
}

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
	s, err := ParseFile("examples/status-24.txt")
	assert.NoError(t, err)
	assert.Equal(t, s.Result, "OK")
	s, err = ParseFile("examples/status-25.txt")
	assert.NoError(t, err)
	assert.Equal(t, s.Result, "OK")
}

func TestOpenFalse(t *testing.T) {
	s, _ := ParseFile("examples/notExistFile.txt")
	assert.Equal(t, s.Result, "open false")
}

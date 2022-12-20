package query

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ui-kreinhard/boca-status-readout/web"
)

func TestReadOut(t *testing.T) {
	go web.NewBocaMockServer("localhost:8081").Start()
	time.Sleep(1 * time.Second)
	t.Run("Reads out printer with no errors", func(t *testing.T) {
		s, err := FetchStatus("localhost:8081/ok")
		assert.Nil(t, err)
		assert.True(t, s.Ready)
		assert.False(t, s.PaperOut)
		assert.False(t, s.PaperJam)
		assert.False(t, s.CutterJam)
	})
	t.Run("Reads out printer with paper out", func(t *testing.T) {
		s, err := FetchStatus("localhost:8081/emptyPaper")
		assert.Nil(t, err)
		assert.False(t, s.Ready)
		assert.True(t, s.PaperOut)
		assert.False(t, s.PaperJam)
		assert.False(t, s.CutterJam)
	})
	t.Run("Reads out ticket count", func(t *testing.T) {
		s, err := FetchStatus("localhost:8081/ok")
		assert.Nil(t, err)
		assert.Equal(t, s.TicketCount, 2914)
	})
}

func TestError(t *testing.T) {
	t.Run("Should result in a immediate error", func(t *testing.T) {
		_, err := FetchStatus("localhost:8082")
		assert.Equal(t, err.Error(), "Get \"http://localhost:8082/realtime.htm\": dial tcp 127.0.0.1:8082: connect: connection refused")
	})
	t.Run("Should result in a connection timeout", func(t *testing.T) {
		_, err := FetchStatusWithTimeout("1.2.3.4", 500*time.Millisecond)
		assert.Error(t, err)
		assert.Equal(t, err, errors.New("timeout"))
	})
}

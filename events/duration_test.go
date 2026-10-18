package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationSecondsKeepsFraction(t *testing.T) {
	var d DurationSeconds
	require.NoError(t, json.Unmarshal([]byte("1.5"), &d))
	require.Equal(t, 1500*time.Millisecond, time.Duration(d))

	require.NoError(t, json.Unmarshal([]byte("36"), &d))
	require.Equal(t, 36*time.Second, time.Duration(d))
}

func TestDurationMinutesKeepsFraction(t *testing.T) {
	var d DurationMinutes
	require.NoError(t, json.Unmarshal([]byte("1.5"), &d))
	require.Equal(t, 90*time.Second, time.Duration(d))

	require.NoError(t, json.Unmarshal([]byte("60"), &d))
	require.Equal(t, 60*time.Minute, time.Duration(d))
}

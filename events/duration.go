package events

import (
	"encoding/json"
	"math"
	"time"
)

type DurationSeconds time.Duration

func (duration *DurationSeconds) UnmarshalJSON(data []byte) error {
	var seconds float64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return err
	}

	*duration = DurationSeconds(time.Duration(seconds) * time.Second)
	return nil
}

func (duration DurationSeconds) MarshalJSON() ([]byte, error) {
	seconds := time.Duration(duration).Seconds()
	return json.Marshal(int64(math.Ceil(seconds)))
}

type DurationMinutes time.Duration

func (duration *DurationMinutes) UnmarshalJSON(data []byte) error {
	var minutes float64
	if err := json.Unmarshal(data, &minutes); err != nil {
		return err
	}

	*duration = DurationMinutes(time.Duration(minutes) * time.Minute)
	return nil
}

func (duration DurationMinutes) MarshalJSON() ([]byte, error) {
	minutes := time.Duration(duration).Minutes()
	return json.Marshal(int64(math.Ceil(minutes)))
}

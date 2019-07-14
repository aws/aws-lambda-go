package events

import (
	"encoding/json"
	"time"
)

type CodeBuildTime time.Time

const codeBuildTimeFormat = "Jan 2, 2006 3:04:05 PM"

func (t CodeBuildTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(codeBuildTimeFormat))
}

func (t *CodeBuildTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	ts, err := time.Parse(codeBuildTimeFormat, s)
	if err != nil {
		return err
	}

	*t = CodeBuildTime(ts)
	return nil
}

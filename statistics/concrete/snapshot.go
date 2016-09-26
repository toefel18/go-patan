package concrete

import "github.com/toefel18/go-patan/statistics"

type StatsSnapshot struct {
	TimestampTaken int64
	Durations      map[string]statistics.Distribution `json:"durations"`
	Counters       map[string]int64                   `json:"counters"`
	Samples        map[string]statistics.Distribution `json:"samples"`
}

func (sh *StatsSnapshot) GetTimestampTaken() int64 {
	return sh.TimestampTaken
}

func (sh *StatsSnapshot) GetDurations() map[string]statistics.Distribution {
	return sh.Durations
}

func (sh *StatsSnapshot) GetCounters() map[string]int64 {
	return sh.Counters
}

func (sh *StatsSnapshot) GetSamples() map[string]statistics.Distribution {
	return sh.Samples
}

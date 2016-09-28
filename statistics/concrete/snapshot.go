package concrete

import "github.com/toefel18/go-patan/statistics/api"

type StatsSnapshot struct {
	TimestampCreated  int64                       `json:"timestampTaken"`
	DurationsSnapshot map[string]api.Distribution `json:"durations"`
	CountersSnapshot  map[string]int64            `json:"counters"`
	SamplesSnapshot   map[string]api.Distribution `json:"samples"`
}

func (sh *StatsSnapshot) TimestampTaken() int64 {
	return sh.TimestampCreated
}

func (sh *StatsSnapshot) Durations() map[string]api.Distribution {
	return sh.DurationsSnapshot
}

func (sh *StatsSnapshot) Counters() map[string]int64 {
	return sh.CountersSnapshot
}

func (sh *StatsSnapshot) Samples() map[string]api.Distribution {
	return sh.SamplesSnapshot
}

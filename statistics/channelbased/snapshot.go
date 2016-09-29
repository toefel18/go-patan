package channelbased

import "github.com/toefel18/go-patan/statistics/api"

type Snapshot struct {
	TimestampCreated  int64                       `json:"timestampTaken"`
	DurationsSnapshot map[string]api.Distribution `json:"durations"`
	CountersSnapshot  map[string]int64            `json:"counters"`
	SamplesSnapshot   map[string]api.Distribution `json:"samples"`
}

func (sh *Snapshot) TimestampTaken() int64 {
	return sh.TimestampCreated
}

func (sh *Snapshot) Durations() map[string]api.Distribution {
	return sh.DurationsSnapshot
}

func (sh *Snapshot) Counters() map[string]int64 {
	return sh.CountersSnapshot
}

func (sh *Snapshot) Samples() map[string]api.Distribution {
	return sh.SamplesSnapshot
}

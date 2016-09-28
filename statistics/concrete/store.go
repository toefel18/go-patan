package concrete

import "github.com/toefel18/go-patan/statistics/api"

type Update struct {
	key    string
	sample int64
}

type Store struct {
	durations map[string]*Distribution
	counters  map[string]int64
	samples   map[string]*Distribution

	durationUpdates chan Update
	counterUpdates  chan Update
	sampleUpdates   chan Update
    snapshotRequest chan chan api.Snapshot
}

func NewStore() *Store {
	return &Store{
		durations: make(map[string]*Distribution),
		counters:  make(map[string]int64),
		samples:   make(map[string]*Distribution),
		durationUpdates: make(chan Update, 8),
		counterUpdates:  make(chan Update, 8),
		sampleUpdates:   make(chan Update, 8),
        snapshotRequest: make(chan chan api.Snapshot, 8),
	}
}

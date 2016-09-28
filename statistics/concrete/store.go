package concrete

import "github.com/toefel18/go-patan/statistics"

type StatStore struct {
    durations map[string]*Distribution
    counters  map[string]int64
    samples   map[string]*Distribution
}

func NewStore() *StatStore {
    return &StatStore{
        durations: make(map[string]*Distribution),
        counters:  make(map[string]int64),
        samples:   make(map[string]*Distribution),
    }
}

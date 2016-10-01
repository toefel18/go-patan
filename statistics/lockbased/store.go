package lockbased

import (
	"github.com/toefel18/go-patan/statistics/api"
	"log"
	"math"
	"sync"
	"time"
)

type Store struct {
	// durations, counters and samples should never be modified by anything else than the StoreUpdater method!
	durations map[string]*Distribution
	counters  map[string]int64
	samples   map[string]*Distribution

	lock *sync.Mutex
}

// Creates a new store and starts a go-routine that listens for requests on the channels.
// Don't forget to call store.Close() when throwing away the store!
func NewStore() *Store {
	store := &Store{
		durations: make(map[string]*Distribution),
		counters:  make(map[string]int64),
		samples:   make(map[string]*Distribution),
		lock:      &sync.Mutex{},
	}
	log.Println("[STATISTICS] created new lockbased store")
	return store
}

func (store *Store) addSample(key string, value int64) {
	store.addToStore(store.samples, key, value)
}

func (store *Store) addDuration(key string, value int64) {
	store.addToStore(store.durations, key, value)
}

func (store *Store) addToCounter(key string, value int64) {
	store.lock.Lock()
	defer store.lock.Unlock()
	store.counters[key] = store.counters[key] + value
}

func (store *Store) addToStore(destination map[string]*Distribution, key string, value int64) {
	store.lock.Lock()
	defer store.lock.Unlock()
	distribution, exists := destination[key]
	if !exists {
		distribution = NewDistribution()
		destination[key] = distribution
	}
	distribution.addSample(value)
}

func (store *Store) Snapshot() api.Snapshot {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.doGetSnapshot()
}

func (store *Store) doGetSnapshot() api.Snapshot {
	durationsCopy := deepCopy(store.durations)
	countersCopy := shallowCopy(store.counters)
	samplesCopy := deepCopy(store.samples)

	return &Snapshot{
		TimestampCreated:  time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
		DurationsSnapshot: durationsCopy,
		CountersSnapshot:  countersCopy,
		SamplesSnapshot:   samplesCopy,
	}
}

func (store *Store) SnapshotAndReset() api.Snapshot {
	store.lock.Lock()
	defer store.lock.Unlock()
	snapshot := store.doGetSnapshot()
	store.doReset()
	return snapshot
}

func (store *Store) Reset() {
	store.lock.Lock()
	defer store.lock.Unlock()
	store.doReset()
}

func (store *Store) doReset() {
	store.durations = make(map[string]*Distribution)
	store.counters = make(map[string]int64)
	store.samples = make(map[string]*Distribution)
}

func deepCopy(source map[string]*Distribution) map[string]api.Distribution {
	distMapCopy := make(map[string]api.Distribution)
	for key, distribution := range source {
		valueCopy := *distribution // dereference pointer to get a copy of the struct
		if math.IsNaN(valueCopy.StdDeviation) {
			valueCopy.StdDeviation = -1.0
		}
		distMapCopy[key] = &valueCopy
	}
	return distMapCopy
}

func shallowCopy(source map[string]int64) map[string]int64 {
	intMapCopy := make(map[string]int64)
	for key, counter := range source {
		intMapCopy[key] = counter
	}
	return intMapCopy
}

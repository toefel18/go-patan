/*
 *
 *     Copyright 2016 Christophe Hesters
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 *
 */

package lockbased

import (
	"log"
	"sync"

	"github.com/toefel18/go-patan/metrics/api"
	"github.com/toefel18/go-patan/metrics/common"
)

//Store holds the state of the lockbased implementation
type Store struct {
	timestampStarted int64

	durations map[string]*common.Distribution
	counters  map[string]int64
	samples   map[string]*common.Distribution

	lock sync.Mutex
}

// NewStore creates a new store and starts a go-routine that listens for requests on the channels.
func NewStore() *Store {
	store := &Store{
		timestampStarted: common.CurrentTimeMillis(),
		durations:        make(map[string]*common.Distribution),
		counters:         make(map[string]int64),
		samples:          make(map[string]*common.Distribution),
	}
	log.Println("[METRICS] created new lockbased store")
	return store
}

func (store *Store) addSample(key string, value float64) {
	store.addToStore(store.samples, key, value)
}

func (store *Store) addDuration(key string, value float64) {
	store.addToStore(store.durations, key, value)
}

func (store *Store) addToCounter(key string, value int64) {
	store.lock.Lock()
	store.counters[key] = store.counters[key] + value
	store.lock.Unlock()
}

func (store *Store) addToStore(destination map[string]*common.Distribution, key string, value float64) {
	store.lock.Lock()
	distribution, exists := destination[key]
	if !exists {
		distribution = common.NewDistribution()
		destination[key] = distribution
	}
	distribution.AddSample(value)
	store.lock.Unlock()
}

//Snapshot creates a new snapshot of the current state
func (store *Store) Snapshot() api.Snapshot {
	store.lock.Lock()
	snapshot := store.doGetSnapshot()
	store.lock.Unlock()
	return snapshot
}

func (store *Store) doGetSnapshot() api.Snapshot {
	durationsCopy := deepCopy(store.durations)
	countersCopy := shallowCopy(store.counters)
	samplesCopy := deepCopy(store.samples)

	return &common.Snapshot{
		TimestampStarted:  store.timestampStarted,
		TimestampCreated:  common.CurrentTimeMillis(),
		DurationsSnapshot: durationsCopy,
		CountersSnapshot:  countersCopy,
		SamplesSnapshot:   samplesCopy,
	}
}

// SnapshotAndReset creates a snapshot and clears the recorded counters, durations and samples
func (store *Store) SnapshotAndReset() api.Snapshot {
	store.lock.Lock()
	snapshot := store.doGetSnapshot()
	store.doReset()
	store.lock.Unlock()
	return snapshot
}

// Reset clears the recorded counters, durations and samples
func (store *Store) Reset() {
	store.lock.Lock()
	store.doReset()
	store.lock.Unlock()
}

func (store *Store) doReset() {
	store.timestampStarted = common.CurrentTimeMillis()
	store.durations = make(map[string]*common.Distribution)
	store.counters = make(map[string]int64)
	store.samples = make(map[string]*common.Distribution)
}

func deepCopy(source map[string]*common.Distribution) map[string]api.Distribution {
	distMapCopy := make(map[string]api.Distribution)
	for key, distribution := range source {
		valueCopy := *distribution // dereference pointer to get a copy of the struct
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

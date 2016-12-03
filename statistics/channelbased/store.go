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

package channelbased

import (
	"log"
	"github.com/toefel18/go-patan/statistics/api"
	"github.com/toefel18/go-patan/statistics/common"
)

// NewMeasurement contains one new measurement to process
type NewMeasurement struct {
	key   string
	value float64
}

// Request contains a query/reset request for the store
type Request struct {
	resetStore            bool
	createSnapshot        bool
	snapshotReturnChannel chan api.Snapshot
}

// Store contains the channelbased store implementation
type Store struct {
	// durations, counters and samples should never be modified by anything else than the StoreUpdater method!
	durations        map[string]*common.Distribution
	counters         map[string]int64
	samples          map[string]*common.Distribution
	timestampStarted int64 // current_time_millis

	durationUpdates chan NewMeasurement
	counterUpdates  chan NewMeasurement
	sampleUpdates   chan NewMeasurement
	requests        chan Request
	stop            chan bool
}

// NewStore creates a new store and starts a go-routine that listens for requests on the channels.
// Don't forget to call store.Close() when throwing away the store!
func NewStore() *Store {
	store := &Store{
		timestampStarted: common.CurrentTimeMillis(),
		durations:        make(map[string]*common.Distribution),
		counters:         make(map[string]int64),
		samples:          make(map[string]*common.Distribution),
		durationUpdates:  make(chan NewMeasurement, 50),
		counterUpdates:   make(chan NewMeasurement, 50),
		sampleUpdates:    make(chan NewMeasurement, 50),
		requests:         make(chan Request, 50),
		stop:             make(chan bool),
	}
	go store.storeUpdater()
	log.Println("[STATISTICS] created new store and started goroutine to respond to updates")
	return store
}

// StoreUpdater The only method that may handle store mutations and snapshot requests.
// start 1, and only 1, goroutine at this function for each store. The facade
// sends store mutations to the appropriate channel. This set-up requires
// no explicit locking
func (store *Store) storeUpdater() {
	for {
		select {
		case durationUpdate := <-store.durationUpdates:
			store.addSample(store.durations, durationUpdate)
		case counterUpdate := <-store.counterUpdates:
			store.counters[counterUpdate.key] = store.counters[counterUpdate.key] + int64(counterUpdate.value)
		case sampleUpdate := <-store.sampleUpdates:
			store.addSample(store.samples, sampleUpdate)
		case request := <-store.requests:
			store.handleRequest(request)
		case <-store.stop:
			log.Println("[STATISTICS] stopping store updater")
			close(store.stop)
			return
		}
	}
}

func (store *Store) addSample(destination map[string]*common.Distribution, sampleToAdd NewMeasurement) {
	distribution, exists := destination[sampleToAdd.key]
	if !exists {
		distribution = common.NewDistribution()
		destination[sampleToAdd.key] = distribution
	}
	distribution.AddSample(sampleToAdd.value)
}

func (store *Store) handleRequest(request Request) {
	var snapshot api.Snapshot
	if request.createSnapshot {
		snapshot = store.snapshot()
	}
	if request.resetStore {
		log.Println("[STATISTICS] clearing all the collected statistics")
		store.timestampStarted = common.CurrentTimeMillis()
		store.durations = make(map[string]*common.Distribution)
		store.counters = make(map[string]int64)
		store.samples = make(map[string]*common.Distribution)
	}
	if request.createSnapshot {
		if request.snapshotReturnChannel != nil {
			request.snapshotReturnChannel <- snapshot
		} else {
			log.Println("[STATISTICS] WARN snapshot requested but return channel was nil")
		}
	}
}

// Close closes the store and frees up the event processing go-routine. Has no effect when store is already closed.
func (store *Store) Close() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[STATISTICS] WARN store was already closed")
		}
	}()
	store.stop <- true
}

func (store *Store) snapshot() api.Snapshot {
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

func deepCopy(source map[string]*common.Distribution) map[string]api.Distribution {
	distMapCopy := make(map[string]api.Distribution)
	for key, distribution := range source {
		distributionCopy := *distribution // dereference pointer to get a copy of the struct
		distMapCopy[key] = &distributionCopy
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

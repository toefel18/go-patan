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
	"github.com/toefel18/go-patan/statistics/api"
	"log"
	"time"
	"math"
)

type NewMeasurement struct {
	key   string
	value int64
}

type Request struct {
	resetStore            bool
	createSnapshot        bool
	snapshotReturnChannel chan api.Snapshot
}

type Store struct {
	// durations, counters and samples should never be modified by anything else than the StoreUpdater method!
	durations       map[string]*Distribution
	counters        map[string]int64
	samples         map[string]*Distribution

	durationUpdates chan NewMeasurement
	counterUpdates  chan NewMeasurement
	sampleUpdates   chan NewMeasurement
	requests        chan Request
	stop            chan bool
}

// Creates a new store and starts a go-routine that listens for requests on the channels.
// Don't forget to call store.Close() when throwing away the store!
func NewStore() *Store {
	store := &Store{
		durations:        make(map[string]*Distribution),
		counters:         make(map[string]int64),
		samples:          make(map[string]*Distribution),
		durationUpdates:  make(chan NewMeasurement, 20),
		counterUpdates:   make(chan NewMeasurement, 20),
		sampleUpdates:    make(chan NewMeasurement, 20),
		requests: make(chan Request, 20),
		stop:             make(chan bool),
	}
	go store.StoreUpdater()
	log.Println("[STATISTICS] created new store and started goroutine to respond to updates")
	return store
}

// The only method that may handle store mutations and snapshot requests.
// start 1, and only 1, goroutine at this function for each store. The facade
// sends store mutations to the appropriate channel. This set-up requires
// no explicit locking
func (store *Store) StoreUpdater() {
	for {
		select {
		case durationUpdate := <-store.durationUpdates:
			store.addSample(store.durations, durationUpdate)
		case counterUpdate := <-store.counterUpdates:
			store.counters[counterUpdate.key] = store.counters[counterUpdate.key] + counterUpdate.value
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

func (store *Store) addSample(destination map[string]*Distribution, sampleToAdd NewMeasurement) {
	distribution, exists := destination[sampleToAdd.key]
	if !exists {
		distribution = NewDistribution()
		destination[sampleToAdd.key] = distribution
	}
	distribution.addSample(sampleToAdd.value)
}

func (store *Store) handleRequest(request Request) {
	var snapshot api.Snapshot
	if request.createSnapshot {
		snapshot = store.Snapshot()
	}
	if request.resetStore {
		log.Println("[STATISTICS] clearing all the collected statistics")
		store.durations = make(map[string]*Distribution)
		store.counters = make(map[string]int64)
		store.samples = make(map[string]*Distribution)
	}
	if request.createSnapshot {
		if request.snapshotReturnChannel != nil {
			request.snapshotReturnChannel <- snapshot
		} else {
			log.Println("[STATISTICS] snapshot requested but return channel was nil")
		}
	}
}

// Closes the store and frees up the event processing goroutine. Has no effect when store is already closed.
func (store *Store) Close() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[STATISTICS] store was already closed")
		}
	}()
	store.stop <- true
}

func (store *Store) Snapshot() api.Snapshot {
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

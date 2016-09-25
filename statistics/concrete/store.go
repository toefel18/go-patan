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

func (store *StatStore) StartStopwatch() Stopwatch {
	return startNewStopwatch()
}

func (store *StatStore) FindDuration(key string) (statistics.Distribution, bool) {
	return getOrEmpty(store.durations, key)
}

func (store *StatStore) FindCounter(key string) (int64, bool) {
	return *store.counters[key]
}

func (store *StatStore) FindSample(key string) (statistics.Distribution, bool) {
	return getOrEmpty(store.samples, key)
}

func (store *StatStore) RecordElapsedTime(key string, stopwatch Stopwatch) int64 {
	millis := stopwatch.ElapsedMillis()
	if distribution, exists := store.durations[key]; exists {
		distribution.addSample(millis)
	} else {
		newDistribution := &Distribution{}
		newDistribution.addSample(millis)
		store.durations[key] = newDistribution
	}
	return millis
}

func (store *StatStore) MeasureFunc(key string, subject func()) int64 {
	return 0
}

func (store *StatStore) MeasureFuncWithReturn(key string, subject func() interface{}) (int64, interface{}) {
	return 0, nil
}

func (store *StatStore) IncrementCounter(key string) {

}

func (store *StatStore) DecrementCounter(key string) {

}

func (store *StatStore) AddToCounter(key string, value int) {

}

func (store *StatStore) AddSample(key string, value int64) {

}

func (store *StatStore) Reset() {

}

func (store *StatStore) Snapshot() statistics.Snapshot {
	return nil
}

func (store *StatStore) SnapshotAndReset() statistics.Snapshot {
	return nil
}

// Fetches the distribution and returns a copy, if the distribution does not exist, an empty one is created.
// The public API of patan does not use pointers and requires copies to be returned.
func getOrEmpty(distributionsByKey map[string]*Distribution, key string) (statistics.Distribution, bool) {
	distribution, present := distributionsByKey[key]
	if !present {
		distribution = &Distribution{}
	}
	return *distribution
}
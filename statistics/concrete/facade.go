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

import "github.com/toefel18/go-patan/statistics/api"

type Facade struct {
    store *StatStore
}

func NewFacade(store *StatStore) *Facade {
    if store == nil {
        panic("store = nil, Facade needs a store")
    }
    return &Facade{store}
}

func (facade Facade) StartStopwatch() api.Stopwatch {
    return startNewStopwatch()
}

func (facade Facade) FindDuration(key string) (api.Distribution, bool) {
    return NewDistribution(), false
    //return getOrEmpty(facade.store.durations, key)
}

func (facade Facade) FindCounter(key string) (int64, bool) {
    return 0, false
    //return *facade.store.counters[key]
}

func (facade Facade) FindSample(key string) (api.Distribution, bool) {
    return NewDistribution(), false
    //return getOrEmpty(facade.store.samples, key)
}

func (facade Facade) RecordElapsedTime(key string, stopwatch api.Stopwatch) int64 {
    return 0
    //millis := stopwatch.ElapsedMillis()
    //if distribution, exists := facade.store.durations[key]; exists {
    //    distribution.addSample(millis)
    //} else {
    //    newDistribution := &Distribution{}
    //    newDistribution.addSample(millis)
    //    facade.store.durations[key] = newDistribution
    //}
    //return millis
}

func (facade Facade) MeasureFunc(key string, subject func()) int64 {
    return 0
}

func (facade Facade) MeasureFuncWithReturn(key string, subject func() interface{}) (int64, interface{}) {
    return 0, nil
}

func (facade Facade) IncrementCounter(key string) {

}

func (facade Facade) DecrementCounter(key string) {

}

func (facade Facade) AddToCounter(key string, value int) {

}

func (facade Facade) AddSample(key string, value int64) {

}

func (facade Facade) Reset() {

}

func (facade Facade) Snapshot() api.Snapshot {
    return &StatsSnapshot{}
}

func (facade Facade) SnapshotAndReset() api.Snapshot {
    return &StatsSnapshot{}
}

// Fetches the distribution and returns a copy, if the distribution does not exist, an empty one is created.
// The public API of patan does not use pointers and requires copies to be returned.
func getOrEmpty(distributionsByKey map[string]*Distribution, key string) (Distribution, bool) {
    distribution, present := distributionsByKey[key]
    if !present {
        distribution = &Distribution{}
    }
    return *distribution, false
}


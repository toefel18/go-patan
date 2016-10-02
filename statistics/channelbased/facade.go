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

import "github.com/toefel18/go-patan/statistics/api"

type Facade struct {
    store *Store
}

func NewFacade(store *Store) *Facade {
    if store == nil {
        panic("store = nil, Facade needs a store")
    }
    return &Facade{store}
}

func (facade *Facade) StartStopwatch() api.Stopwatch {
    return startNewStopwatch()
}

func (facade *Facade) RecordElapsedTime(key string, stopwatch api.Stopwatch) int64 {
    millis := stopwatch.ElapsedMillis()
    facade.store.durationUpdates <- NewMeasurement{key, millis}
    return millis
}

func (facade *Facade) MeasureFunc(key string, subject func()) int64 {
    sw := startNewStopwatch()
    subject()
    return facade.RecordElapsedTime(key, sw)
}

func (facade *Facade) IncrementCounter(key string) {
    facade.AddToCounter(key, 1)
}

func (facade *Facade) DecrementCounter(key string) {
    facade.AddToCounter(key, -1)
}

func (facade *Facade) AddToCounter(key string, value int64) {
    facade.store.counterUpdates <- NewMeasurement{key, value}
}

func (facade *Facade) AddSample(key string, value int64) {
    facade.store.sampleUpdates <- NewMeasurement{key, value}
}

func (facade *Facade) Reset() {
    facade.store.requests <- Request{resetStore: true, createSnapshot:false}
}

// USAGE NOTE: Invocation order does not necessarily reflect the processing order! Users should
// not depend on that. Consider the following example:
//
// stopwatch := statistics.StartStopwatch()
// ... heavy work for 3 seconds
// statistics.RecordElapsedTime("my.heavy.operation", stopwatch)   // A
// snapshot := statistics.Snapshot()                               // B
//
// It is possible (and even likely) that snapshot doesn't have my.heavy.operation yet, meaning that
// B is executed earlier than A! This differs from the java version of Patan and is a consequence
// of the non-blocking setup with channels. This is OK because patan is meant to give insight in
// the distribution of data over a longer period of time, not for individual measurements.
func (facade *Facade) Snapshot() api.Snapshot {
    returnChannel := make(chan api.Snapshot)
    facade.store.requests <- Request{resetStore:false, createSnapshot:true, snapshotReturnChannel: returnChannel}
    return <-returnChannel
}

func (facade *Facade) SnapshotAndReset() api.Snapshot {
    returnChannel := make(chan api.Snapshot)
    facade.store.requests <- Request{resetStore:true, createSnapshot:true, snapshotReturnChannel: returnChannel}
    return <-returnChannel
}

func (facade *Facade) Close() {
    facade.store.Close()
}
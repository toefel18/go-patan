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
	"github.com/toefel18/go-patan/metrics/api"
	"github.com/toefel18/go-patan/metrics/common"
)

// Facade provides a nice API on the lockbased store
type Facade struct {
	store *Store
}

// NewFacade creates a new facade initialized with the store
func NewFacade(store *Store) *Facade {
	if store == nil {
		panic("store = nil, Facade needs a store")
	}
	return &Facade{store}
}

// StartStopwatch starts a new stopwatch
func (facade *Facade) StartStopwatch() api.Stopwatch {
	return common.StartNewStopwatch()
}

// RecordElapsedTime records the elapsed time of the stopwatch under the distribution identified with key
func (facade *Facade) RecordElapsedTime(key string, stopwatch api.Stopwatch) float64 {
	millis := stopwatch.ElapsedMillis()
	facade.store.addDuration(key, millis)
	return millis
}

// MeasureFunc runs the subject function and records it's execution duration under the distribution identified with key
func (facade *Facade) MeasureFunc(key string, subject func()) float64 {
	sw := common.StartNewStopwatch()
	subject()
	return facade.RecordElapsedTime(key, sw)
}

// IncrementCounter increments the counter identified by key by 1
func (facade *Facade) IncrementCounter(key string) {
	facade.AddToCounter(key, 1)
}

// DecrementCounter decrements the counter identified by key by 1
func (facade *Facade) DecrementCounter(key string) {
	facade.AddToCounter(key, -1)
}

// AddToCounter adds value to the counter identified by key, value can be negative
func (facade *Facade) AddToCounter(key string, value int64) {
	facade.store.addToCounter(key, value)
}

// AddSample adds a sample to the distribution identified by value, if the distribution doesn't
// exist, it will be created
func (facade *Facade) AddSample(key string, value float64) {
	facade.store.addSample(key, value)
}

// Reset clears the store
func (facade *Facade) Reset() {
	facade.store.Reset()
}

// Snapshot returns a snapshot of all the counters, durations and samples recorded
// since creation or the last reset.
func (facade *Facade) Snapshot() api.Snapshot {
	return facade.store.Snapshot()
}

// SnapshotAndReset returns a snapshot of all the counters, durations and samples recorded
// since creation or the last reset, and then clears the internal state
func (facade *Facade) SnapshotAndReset() api.Snapshot {
	return facade.store.SnapshotAndReset()
}

// Close closes the underlying store
func (facade *Facade) Close() {
	// nop, the internal store has no resources to clean-up
}

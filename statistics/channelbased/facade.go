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

// Package channelbased is DEPRECATED, mainly due to the Reset() behaviour. It's unclear
// when, after a call to Reset(), the state is actually reset. This CAN lead
// to unexpected behaviour. Clients are advised to use the lockbased implementation.
package channelbased

import (
	"github.com/toefel18/go-patan/statistics/api"
	"github.com/toefel18/go-patan/statistics/common"
)

// Facade provides a nice API on the channelbased store
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
	facade.store.durationUpdates <- NewMeasurement{key, millis}
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
	facade.store.counterUpdates <- NewMeasurement{key, float64(value)}
}

// AddSample adds a sample to the distribution identified by value, if the distribution doesn't
// exist, it will be created
func (facade *Facade) AddSample(key string, value float64) {
	facade.store.sampleUpdates <- NewMeasurement{key, value}
}

// Reset clears the store
func (facade *Facade) Reset() {
	facade.store.requests <- Request{resetStore: true, createSnapshot: false}
}

// Snapshot returns a snapshot of all the counters, durations and samples recorded
// since creation or the last reset.
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
// the distribution of data over a longer period of time, not for individual measurements. For immediate
// consistency, use the lockbased implementation
func (facade *Facade) Snapshot() api.Snapshot {
	returnChannel := make(chan api.Snapshot)
	facade.store.requests <- Request{resetStore: false, createSnapshot: true, snapshotReturnChannel: returnChannel}
	return <-returnChannel
}

// SnapshotAndReset creates a snapshot and resets it.
func (facade *Facade) SnapshotAndReset() api.Snapshot {
	returnChannel := make(chan api.Snapshot)
	facade.store.requests <- Request{resetStore: true, createSnapshot: true, snapshotReturnChannel: returnChannel}
	return <-returnChannel
}

// Close closes the underlying store
func (facade *Facade) Close() {
	facade.store.Close()
}

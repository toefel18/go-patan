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

// Package metrics contains a ready to use instance of patan. This instance can
// be used as the sole active instance of patan within the application. Clients
// are advised to use this. example:
// metrics.AddSample("key", 123)
package metrics

import (
	"log"

	"github.com/toefel18/go-patan/metrics/api"
	"github.com/toefel18/go-patan/metrics/lockbased"
)

// Standard instance of patan, ready to use
var std api.Facade

// Initializes a global instance of patan.
func init() {
	log.Println("[METRICS] intializing global instance of patan")
	std = lockbased.NewFacade(lockbased.NewStore())
	log.Println("[METRICS] global version of patan initialized")
}

// New returns a new API facade with a new and empty underlying store.
func New() api.Facade {
	return lockbased.NewFacade(lockbased.NewStore())
}

// The methods below are equal to those of api.Facade and operate on the
// global instance of api.Facade that is ready to use

// StartStopwatch starts a new stopwatch
func StartStopwatch() api.Stopwatch {
	return std.StartStopwatch()
}

// RecordElapsedTime records the elapsed time of the stopwatch under the distribution identified with key
func RecordElapsedTime(key string, stopwatch api.Stopwatch) float64 {
	return std.RecordElapsedTime(key, stopwatch)
}

// MeasureFunc runs the subject function and records it's execution duration under the distribution identified with key
func MeasureFunc(key string, subject func()) float64 {
	return std.MeasureFunc(key, subject)
}

// MeasureFuncCanPanic runs the subject function and records it's execution duration under the distribution identified
// with key. When subject() panics, the measurement is recored under the same key with .panic appended. This function
// itself will panic with the same error as the inner function.
func MeasureFuncCanPanic(key string, subject func()) float64 {
	return std.MeasureFuncCanPanic(key, subject)
}

// IncrementCounter increments the counter identified by key by 1
func IncrementCounter(key string) {
	std.IncrementCounter(key)
}

// DecrementCounter decrements the counter identified by key by 1
func DecrementCounter(key string) {
	std.DecrementCounter(key)
}

// AddToCounter adds value to the counter identified by key, value can be negative
func AddToCounter(key string, value int64) {
	std.AddToCounter(key, value)
}

// AddSample adds a sample to the distribution identified by value, if the distribution doesn't
// exist, it will be created
func AddSample(key string, value float64) {
	std.AddSample(key, value)
}

// Reset clears the store
func Reset() {
	std.Reset()
}

// Snapshot returns a snapshot of all the counters, durations and samples recorded
// since creation or the last reset.
func Snapshot() api.Snapshot {
	return std.Snapshot()
}

// SnapshotAndReset returns a snapshot of all the counters, durations and samples recorded
// since creation or the last reset, and then clears the internal state
func SnapshotAndReset() api.Snapshot {
	return std.SnapshotAndReset()
}

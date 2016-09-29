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
package api

type Stopwatch interface {
	ElapsedMillis() int64
	ElapsedNanos() int64
}

// Models a statistical distribution
type Distribution interface {
	SampleCount() int64
	Min() int64
	Max() int64
	Avg() float64
	Variance() float64
	StdDev() float64
}

type Snapshot interface {
	TimestampTaken() int64
	Durations() map[string]Distribution
	Counters() map[string]int64
	Samples() map[string]Distribution
}

type Facade interface {
	StartStopwatch() Stopwatch

	// Records the elapsed time of the stopwatch and adds that to the distribution identified by key.
	// Returns the recorded millis
	RecordElapsedTime(key string, stopwatch Stopwatch) int64

	// Records duration of the subject function and adds that to the distribution identified by key.
	// Returns the recorded millis
	MeasureFunc(key string, subject func()) int64

	// Increments the counter identified with key by 1. If the counter does not yet exist, it will be created
	// with initial value of 1
	IncrementCounter(key string)

	// Decrements the counter identified with key by 1. If the counter does not yet exist, it will be created
	// with initial value of -1
	DecrementCounter(key string)

	// Adds value to the counter identified with key with key, if the counter does not yet exist, it will be created
	// and initialized to value. Value can be negative.
	AddToCounter(key string, value int64)

	// Adds a value to the sample distribution identified by key. If the distribution does not yet exist, value will be it's initial value.
	AddSample(key string, value int64)

	// Clears all durations, counters and samples
	Reset()

	// Creates a snapshot containing all currently registered durations, counters and samples

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
	Snapshot() Snapshot

	// Creates a snapshot and then calls Reset()
	// Also, see documentation above
	SnapshotAndReset() Snapshot

	// Free up resources
	Close();
}

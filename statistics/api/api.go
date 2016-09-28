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

// Models a sampling distribution
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

	// Records duration of the subject function and adds that to the distribution identified by key.
	// Returns the recorded millis and the returned value of the subject function
	MeasureFuncWithReturn(key string, subject func() interface{}) (int64, interface{})

	// Increments the counter identified with key by 1. If the counter does not yet exist, it will be created
	// with initial value of 1
	IncrementCounter(key string)

	// Decrements the counter identified with key by 1. If the counter does not yet exist, it will be created
	// with initial value of -1
	DecrementCounter(key string)

	// Adds value to the counter identified with key with key, if the counter does not yet exist, it will be created
	// and initialized to value. Value can be negative.
	AddToCounter(key string, value int)

	// Adds a value to the sample distribution identified by key. If the distribution does not yet exist, value will be it's initial value.
	AddSample(key string, value int64)

	// Clears all durations, counters and samples
	Reset()

	// Creates a snapshot containing all currently registered durations, counters and samples
	Snapshot() Snapshot

	// Creates a snapshot and then calls Reset()
	SnapshotAndReset() Snapshot
}

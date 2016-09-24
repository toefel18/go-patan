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
package statistics

type Stopwatch interface {
    ElapsedMillis() int64
    ElapsedNanos() int64
}

type Distribution interface {
    GetSampleCount() int64
    GetMinimum() int64
    GetMaximum() int64
    GetSampleAverage() float64
    GetSampleVariance() float64
    GetSampleStdDeviation() float64
}

type Snapshot struct {
    TimestampTaken int64
    Durations      map[string]Distribution `json:"durations"`
    Counters       map[string]int64        `json:"counters"`
    Samples        map[string]Distribution `json:"samples"`
}

type Store interface {
    StartStopwatch() Stopwatch

    // Finds a duration and returns (true, duration distribution) if found else (false, distribution with all fields set to 0)
    FindDuration(key string) (bool, Distribution)

    // Finds a counter and returns (true, counter value) if found, else (false, 0)
    FindCounter(key string) (bool, int64)

    // Finds a sample and returns (true, sample distribution) if found else (false, distribution with all fields set to 0)
    FindSample(key string) (bool, Distribution)

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

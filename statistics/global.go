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

import (
	"github.com/toefel18/go-patan/statistics/api"
	"github.com/toefel18/go-patan/statistics/channelbased"
	"log"
)

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

// Standard instance of patan, ready to use
var std api.Facade

// Initializes a global instance of patan.
func init() {
	log.Println("[STATISTICS] intializing global instance of patan")
	std = channelbased.NewFacade(channelbased.NewStore())
	log.Println("[STATISTICS] global version of patan initialized")
}

// The methods below are equal to those of api.Facade and operate on the
// global instance of api.Facade that is ready to use

func StartStopwatch() api.Stopwatch {
	return std.StartStopwatch()
}

func RecordElapsedTime(key string, stopwatch api.Stopwatch) int64 {
	return std.RecordElapsedTime(key, stopwatch)
}

func MeasureFunc(key string, subject func()) int64 {
	return std.MeasureFunc(key, subject)
}

func IncrementCounter(key string) {
	std.IncrementCounter(key)
}

func DecrementCounter(key string) {
	std.DecrementCounter(key)
}

func AddToCounter(key string, value int64) {
	std.AddToCounter(key, value)
}

func AddSample(key string, value int64) {
	std.AddSample(key, value)
}

func Reset() {
	std.Reset()
}

func Snapshot() api.Snapshot {
	return std.Snapshot()
}

func SnapshotAndReset() api.Snapshot {
	return std.SnapshotAndReset()
}

func Close() {
	std.Close()
}

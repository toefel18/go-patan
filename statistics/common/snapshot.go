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

package common

import "github.com/toefel18/go-patan/statistics/api"

// Snapshot contains a copy of all the measurements recorded between TimestampStarted and TimestampCreated
type Snapshot struct {
	TimestampStarted  int64                       `json:"timestampStarted"`
	TimestampCreated  int64                       `json:"timestampTaken"`
	DurationsSnapshot map[string]api.Distribution `json:"durations"`
	CountersSnapshot  map[string]int64            `json:"counters"`
	SamplesSnapshot   map[string]api.Distribution `json:"samples"`
}

// CreatedTimestamp returns the timestamp (millis since epoch) on which the snapshot was created
func (sh *Snapshot) CreatedTimestamp() int64 {
	return sh.TimestampCreated
}

// StartedTimestamp returns the timestamp (millis since epoch) on which recordings started
func (sh *Snapshot) StartedTimestamp() int64 {
	return sh.TimestampStarted
}

// Durations returns the map of recorded durations
func (sh *Snapshot) Durations() map[string]api.Distribution {
	return sh.DurationsSnapshot
}

// Counters returns the map of recorded counters
func (sh *Snapshot) Counters() map[string]int64 {
	return sh.CountersSnapshot
}

//Samples returns the map of recorded Samples
func (sh *Snapshot) Samples() map[string]api.Distribution {
	return sh.SamplesSnapshot
}

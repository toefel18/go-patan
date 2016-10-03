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

type Snapshot struct {
	TimestampCreated  int64                       `json:"timestampTaken"`
	DurationsSnapshot map[string]api.Distribution `json:"durations"`
	CountersSnapshot  map[string]int64            `json:"counters"`
	SamplesSnapshot   map[string]api.Distribution `json:"samples"`
}

func (sh *Snapshot) TimestampTaken() int64 {
	return sh.TimestampCreated
}

func (sh *Snapshot) Durations() map[string]api.Distribution {
	return sh.DurationsSnapshot
}

func (sh *Snapshot) Counters() map[string]int64 {
	return sh.CountersSnapshot
}

func (sh *Snapshot) Samples() map[string]api.Distribution {
	return sh.SamplesSnapshot
}

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

import (
	"testing"
	"github.com/toefel18/go-patan/metrics/api"
	"reflect"
)

func TestSnapshot(t *testing.T) {
	durations := make(map[string]api.Distribution)
	counters := make(map[string]int64)
	samples := make(map[string]api.Distribution)
	snapshot := &Snapshot{
		TimestampStarted:  10000,
		TimestampCreated:  20000,
		DurationsSnapshot: durations,
		CountersSnapshot:  counters,
		SamplesSnapshot:   samples,
	}
	if snapshot.StartedTimestamp() != 10000 {
		t.Errorf("Started Timestamp = %v, expected %v", snapshot.StartedTimestamp(), 10000)
	}
	if snapshot.CreatedTimestamp() != 20000 {
		t.Errorf("Started Timestamp = %v, expected %v", snapshot.CreatedTimestamp(), 20000)
	}
	if !reflect.DeepEqual(snapshot.Durations(), durations) {
		t.Error("Durations returns a different instance than expected")
	}
	if !reflect.DeepEqual(snapshot.Samples(), durations) {
		t.Error("Samples returns a different instance than expected")
	}
	if !reflect.DeepEqual(snapshot.Counters(), counters){
		t.Error("Counters returns a different instance than expected")
	}
}
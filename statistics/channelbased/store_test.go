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

import (
	"github.com/toefel18/go-patan/statistics/api"
	"math"
	"testing"
	"time"
)

const (
	ProcessingTime = 30 * time.Millisecond
	Duration = 1
	Sample = 2
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Error("new store returnes nil")
		t.FailNow()
	}
	store.Close()
}

func TestMultipleCloseDoesNotHang(t *testing.T) {
	store := NewStore()
	store.Close()
	finished := make(chan bool, 2)
	go func() {
		store.Close()
		finished <- true
	}()
	go func() {
		store.Close()
		finished <- true
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-finished:
		case <-time.After(100 * time.Millisecond):
			t.Error("closing the store in multiple goroutines at once causes hang!")
		}
	}
}

func TestAddDurationUpdate(t *testing.T) {
	store := NewStore()
	defer store.Close()
	var snapshot api.Snapshot
	var dist api.Distribution
	// add first item and check that it's added
	snapshot, dist = addSample(store, Duration, "test.duration", 10, t)
	assertDistributionHasValues(dist, 1, 10, 10, 10, 0.0, math.NaN(), t)
	assertCountersEmpty(snapshot, t)
	assertSamplesEmpty(snapshot, t)
	// add second sample and check that a new snapshot contains the update
	snapshot, dist = addSample(store, Duration, "test.duration", 20, t)
	assertDistributionHasValues(dist, 2, 10, 20, 15.0, 50.0, 7.0710, t)
	assertCountersEmpty(snapshot, t)
	assertSamplesEmpty(snapshot, t)

	// add a third sample and check again
	snapshot, dist = addSample(store, Duration, "test.duration", 0, t)
	assertDistributionHasValues(dist, 3, 0, 20, 10.0, 50.0, 10.0, t)
	assertCountersEmpty(snapshot, t)
	assertSamplesEmpty(snapshot, t)
}

func TestAddSampleUpdate(t *testing.T) {
	store := NewStore()
	defer store.Close()
	var snapshot api.Snapshot
	var dist api.Distribution
	// add first item and check that it's added
	snapshot, dist = addSample(store, Sample, "test.sample", 10, t)
	assertDistributionHasValues(dist, 1, 10, 10, 10, 0.0, math.NaN(), t)
	assertCountersEmpty(snapshot, t)
	assertDurationsEmpty(snapshot, t)
	// add second sample and check that a new snapshot contains the update
	snapshot, dist = addSample(store, Sample, "test.sample", 20, t)
	assertDistributionHasValues(dist, 2, 10, 20, 15.0, 50.0, 7.0710, t)
	assertCountersEmpty(snapshot, t)
	assertDurationsEmpty(snapshot, t)
	// add a third sample and check again
	snapshot, dist = addSample(store, Sample, "test.sample", 0, t)
	assertDistributionHasValues(dist, 3, 0, 20, 10.0, 50.0, 10.0, t)
	assertCountersEmpty(snapshot, t)
	assertDurationsEmpty(snapshot, t)
}

func TestAddCounter(t *testing.T) {
	store := NewStore()
	defer store.Close()
	store.counterUpdates <- NewMeasurement{"active-sessions", 10}
	time.Sleep(30 * time.Millisecond)
	snapshot := getSnapshot(store, t)
	if snapshot.Counters()["active-sessions"] != 10 {
		t.Error("Counter active-sessions should be 10 but was", snapshot.Counters()["active-sessions"])
	}
}

func TestSnapshotsAreDisconnectedFromStore(t *testing.T) {
	store := NewStore()
	defer store.Close()

	// add first item and check that it's added
	snapshot1, dist1 := addSample(store, Duration, "test.duration", 10, t)
	assertDistributionHasValues(dist1, 1, 10, 10, 10, 0.0, math.NaN(), t)

	// add second sample and check that a new snapshot contains the update
	snapshot2, dist2 := addSample(store, Duration, "test.duration", 20, t)
	assertDistributionHasValues(dist1, 1, 10, 10, 10, 0.0, math.NaN(), t)
	assertDistributionHasValues(dist2, 2, 10, 20, 15.0, 50.0, 7.0710, t)

	// add a third sample and check again
	snapshot3, dist3 := addSample(store, Duration, "another.duration", 10, t)
	assertDistributionHasValues(dist1, 1, 10, 10, 10, 0.0, math.NaN(), t)
	assertDistributionHasValues(dist2, 2, 10, 20, 15.0, 50.0, 7.0710, t)
	assertDistributionHasValues(dist3, 1, 10, 10, 10, 0.0, math.NaN(), t)
	if len(snapshot1.Durations()) != len(snapshot2.Durations()) || len(snapshot1.Durations()) != 1 {
		t.Error("snapshot 1 and 2 use the same duration key, so should only have 1 duration but have", len(snapshot1.Durations()), len(snapshot2.Durations()))
	}
	if len(snapshot3.Durations()) != 2 {
		t.Error("snapshot 3 should contain 2 durations but has", len(snapshot3.Durations()))
	}
}

func addSample(store *Store, durationOrSample int, key string, value int64, t *testing.T) (api.Snapshot, api.Distribution) {
	// send the update to the channel
	if durationOrSample == Duration {
		store.durationUpdates <- NewMeasurement{key, value}
	} else {
		store.sampleUpdates <- NewMeasurement{key, value}
	}

	// because channels are selected by the processing thread randomly, the request for a snapshot below
	// might actually get scheduled before the addition. This sleep gives it enough time to pickup the update first.
	time.Sleep(ProcessingTime)

	snapshot := getSnapshot(store, t)

	// check that key is present
	var dist api.Distribution
	if durationOrSample == Duration {
		dist = snapshot.Durations()[key]
	} else {
		dist = snapshot.Samples()[key]
	}
	if dist == nil {
		t.Error(key, "key was not present in snapshot")
		t.FailNow()
	}
	return snapshot, dist
}

func getSnapshot(store *Store, t *testing.T) api.Snapshot {
	respondWithSnapshotTo := make(chan api.Snapshot)
	store.requests <- Request{resetStore: false, createSnapshot: true, snapshotReturnChannel: respondWithSnapshotTo}
	snapshot := <-respondWithSnapshotTo
	nowMillis := time.Now().UnixNano() / time.Millisecond.Nanoseconds()
	if nowMillis - snapshot.TimestampTaken() > 500 {
		t.Error("TimestampTaken of snapshot is older than half a second, possibly an error")
	}
	return snapshot
}

func assertDurationsEmpty(snapshot api.Snapshot, t *testing.T) {
	if len(snapshot.Durations()) != 0 {
		t.Error("snapshot contains", len(snapshot.Durations()), "durations, expected none")
	}
}

func assertCountersEmpty(snapshot api.Snapshot, t *testing.T) {
	if len(snapshot.Counters()) != 0 {
		t.Error("snapshot contains", len(snapshot.Counters()), "counters, expected none")
	}
}

func assertSamplesEmpty(snapshot api.Snapshot, t *testing.T) {
	if len(snapshot.Samples()) != 0 {
		t.Error("snapshot contains", len(snapshot.Samples()), "samples, expected none")
	}
}

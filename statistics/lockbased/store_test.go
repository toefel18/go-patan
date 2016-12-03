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

package lockbased

import (
	"testing"
	"time"

	"github.com/toefel18/go-patan/statistics/api"
	"github.com/toefel18/go-patan/statistics/common/commontest"
	"github.com/toefel18/go-patan/statistics/common"
)

const (
	ProcessingTime = 30 * time.Millisecond
	Duration       = 1
	Sample         = 2
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Error("new store returnes nil")
		t.FailNow()
	}
}

func TestAddDurationUpdate(t *testing.T) {
	store := NewStore()
	var snapshot api.Snapshot
	var dist api.Distribution
	// add first item and check that it's added
	snapshot, dist = addSample(store, Duration, "test.duration", 10, t)
	commontest.AssertDistributionHasValues(dist, 1, 10, 10, 10, 0.0, t)
	assertCountersEmpty(snapshot, t)
	assertSamplesEmpty(snapshot, t)
	// add second sample and check that a new snapshot contains the update
	snapshot, dist = addSample(store, Duration, "test.duration", 20, t)
	commontest.AssertDistributionHasValues(dist, 2, 10, 20, 15.0, 7.0710, t)
	assertCountersEmpty(snapshot, t)
	assertSamplesEmpty(snapshot, t)
	// add a third sample and check again
	snapshot, dist = addSample(store, Duration, "test.duration", 0, t)
	commontest.AssertDistributionHasValues(dist, 3, 0, 20, 10.0, 10.0, t)
	assertCountersEmpty(snapshot, t)
	assertSamplesEmpty(snapshot, t)
}

func TestAddSampleUpdate(t *testing.T) {
	store := NewStore()
	var snapshot api.Snapshot
	var dist api.Distribution
	// add first item and check that it's added
	snapshot, dist = addSample(store, Sample, "test.sample", 10, t)
	commontest.AssertDistributionHasValues(dist, 1, 10, 10, 10, 0.0, t)
	assertCountersEmpty(snapshot, t)
	assertDurationsEmpty(snapshot, t)
	// add second sample and check that a new snapshot contains the update
	snapshot, dist = addSample(store, Sample, "test.sample", 20, t)
	commontest.AssertDistributionHasValues(dist, 2, 10, 20, 15.0, 7.0710, t)
	assertCountersEmpty(snapshot, t)
	assertDurationsEmpty(snapshot, t)
	// add a third sample and check again
	snapshot, dist = addSample(store, Sample, "test.sample", 0, t)
	commontest.AssertDistributionHasValues(dist, 3, 0, 20, 10.0, 10.0, t)
	assertCountersEmpty(snapshot, t)
	assertDurationsEmpty(snapshot, t)
}

func TestAddCounter(t *testing.T) {
	store := NewStore()
	store.addToCounter("active-sessions", 10)
	snapshot := store.doGetSnapshot()
	if snapshot.Counters()["active-sessions"] != 10 {
		t.Error("Counter active-sessions should be 10 but was", snapshot.Counters()["active-sessions"])
	}
}

func TestStore_Reset(t *testing.T) {
	store := NewStore()
	store.addToCounter("active-sessions", 10)
	store.addSample("sample", 123.0)
	store.addDuration("duration", 1674)
	snapshot := store.doGetSnapshot()
	time.Sleep(100 * time.Millisecond)
	store.Reset()
	expStartedTime := common.CurrentTimeMillis()
	snapshotAfterReset := store.doGetSnapshot()
	if len(snapshotAfterReset.Counters()) > 0 || len(snapshotAfterReset.Durations()) > 0 || len(snapshotAfterReset.Samples()) > 0 {
		t.Error("After a Reset() the store still contains items!")
	}
	if snapshotAfterReset.StartedTimestamp() == snapshot.StartedTimestamp() {
		t.Error("After a reset, the snapshots should differ")
	}
	if !commontest.CloseTo(expStartedTime, snapshotAfterReset.StartedTimestamp(), 30) {
		t.Errorf("expStartTime(%v) is not close enough to snapshotAfterReset(%v)", expStartedTime, snapshotAfterReset.StartedTimestamp())
	}
}

func TestSnapshotsAreDisconnectedFromStore(t *testing.T) {
	store := NewStore()

	// add first item and check that it's added
	snapshot1, dist1 := addSample(store, Duration, "test.duration", 10, t)
	commontest.AssertDistributionHasValues(dist1, 1, 10, 10, 10, 0.0, t)

	// add second sample and check that a new snapshot contains the update
	snapshot2, dist2 := addSample(store, Duration, "test.duration", 20, t)
	commontest.AssertDistributionHasValues(dist1, 1, 10, 10, 10, 0.0, t)
	commontest.AssertDistributionHasValues(dist2, 2, 10, 20, 15.0, 7.0710, t)

	// add a third sample and check again
	snapshot3, dist3 := addSample(store, Duration, "another.duration", 10, t)
	commontest.AssertDistributionHasValues(dist1, 1, 10, 10, 10, 0.0, t)
	commontest.AssertDistributionHasValues(dist2, 2, 10, 20, 15.0, 7.0710, t)
	commontest.AssertDistributionHasValues(dist3, 1, 10, 10, 10, 0.0, t)
	if len(snapshot1.Durations()) != len(snapshot2.Durations()) || len(snapshot1.Durations()) != 1 {
		t.Error("snapshot 1 and 2 use the same duration key, so should only have 1 duration but have", len(snapshot1.Durations()), len(snapshot2.Durations()))
	}
	if len(snapshot3.Durations()) != 2 {
		t.Error("snapshot 3 should contain 2 durations but has", len(snapshot3.Durations()))
	}
}

func addSample(store *Store, durationOrSample int, key string, value float64, t *testing.T) (api.Snapshot, api.Distribution) {
	if durationOrSample == Duration {
		store.addDuration(key, value)
	} else {
		store.addSample(key, value)
	}
	snapshot := store.doGetSnapshot()

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

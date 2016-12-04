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
	"math"
	"testing"
	"time"

	"github.com/toefel18/go-patan/metrics/api"
	"github.com/toefel18/go-patan/metrics/common/commontest"
)

func TestNewFacade(t *testing.T) {
	testStore := NewStore()
	facade := NewFacade(testStore)
	if facade == nil {
		t.Error("NewFacade returned nil")
	}
}

func TestNewFacadeWithNilStore(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewFacade should panic when store=nil, but no panic")
		}
	}()

	NewFacade(nil)
}

func TestFacadeImplementsApiInterface(t *testing.T) {
	var facade = NewFacade(NewStore())
	var apiFacade api.Facade = facade
	if apiFacade.StartStopwatch() == nil {
		t.Error("lockbased.Facade has problems implementing the api.Facade")
	}
}

func TestFacadeHappyFlow(t *testing.T) {
	facade := NewFacade(NewStore())

	//add a sample and reset
	facade.AddSample("some.sample", 10)
	facade.Reset()

	// add some test data
	facade.AddSample("some.sample", 10)
	facade.MeasureFunc("some.duration", func() { time.Sleep(100 * time.Millisecond) })
	facade.IncrementCounter("some.counter")

	// allow the other goroutine to process the results
	time.Sleep(50 * time.Millisecond)

	// assert expected values
	var snapshot api.Snapshot
	snapshot = facade.Snapshot()
	commontest.AssertDistributionHasValues(snapshot.Samples()["some.sample"], 1, 10, 10, 10, 0.0, t)
	assertDurationWithin(snapshot.Durations()["some.duration"], 1, 100, 100, 100, t, 30)
	assertCounter(snapshot, "some.counter", 1, t)

	// add more samples
	facade.AddSample("some.sample", 15)
	facade.DecrementCounter("some.counter")
	sw := facade.StartStopwatch()
	time.Sleep(100 * time.Millisecond)
	facade.RecordElapsedTime("some.duration", sw)
	time.Sleep(50 * time.Millisecond)

	// assert expected values
	snapshot = facade.SnapshotAndReset()
	commontest.AssertDistributionHasValues(snapshot.Samples()["some.sample"], 2, 10, 15, 12.5, 3.5355339, t)
	assertDurationWithin(snapshot.Durations()["some.duration"], 2, 100, 100, 100, t, 30)
	assertCounter(snapshot, "some.counter", 0, t)

	snapshot = facade.Snapshot()
	if len(snapshot.Durations()) > 0 || len(snapshot.Counters()) > 0 || len(snapshot.Samples()) > 0 {
		t.Error("Snapshot and reset should have cleared the data, everything should be empty")
	}

	// same test as before to assert that everything still works after reset
	facade.AddSample("some.sample", 10)
	facade.MeasureFunc("some.duration", func() { time.Sleep(100 * time.Millisecond) })
	facade.IncrementCounter("some.counter")

	// allow the other goroutine to process the results
	time.Sleep(50 * time.Millisecond)

	// assert expected values
	snapshot = facade.Snapshot()
	commontest.AssertDistributionHasValues(snapshot.Samples()["some.sample"], 1, 10, 10, 10, 0.0, t)
	assertDurationWithin(snapshot.Durations()["some.duration"], 1, 100, 100, 100, t, 30)
	assertCounter(snapshot, "some.counter", 1, t)
}

// this test is replicated from the distribution and is useful as an integration test.
func TestDistributionAddSample1To10(t *testing.T) {
	facade := NewFacade(NewStore())
	for i := 1; i <= 10; i++ {
		facade.AddSample("sample", float64(i))
	}
	expDeviation := math.Sqrt((2*4.5*4.5 + 2*3.5*3.5 + 2*2.5*2.5 + 2*1.5*1.5 + 2*0.5*0.5) / 9)
	time.Sleep(100 * time.Millisecond)
	snapshot := facade.Snapshot()
	commontest.AssertDistributionHasValues(snapshot.Samples()["sample"], 10, 1.0, 10, 5.5, expDeviation, t)
}

func TestDistributionAddSample10To1(t *testing.T) {
	facade := NewFacade(NewStore())
	for i := 10; i >= 1; i-- {
		facade.AddSample("sample", float64(i))
	}
	expDeviation := math.Sqrt((2*4.5*4.5 + 2*3.5*3.5 + 2*2.5*2.5 + 2*1.5*1.5 + 2*0.5*0.5) / 9)
	time.Sleep(100 * time.Millisecond)
	snapshot := facade.Snapshot()
	commontest.AssertDistributionHasValues(snapshot.Samples()["sample"], 10, 1.0, 10, 5.5, expDeviation, t)
}

func assertCounter(snapshot api.Snapshot, key string, value int64, t *testing.T) {
	if snapshot.Counters()[key] != value {
		t.Error("Expected counter with value", value, "but got", snapshot.Counters()[key])
	}
}

func assertDurationWithin(dist api.Distribution, sampleCount int64, min, max, avg float64, t *testing.T, within float64) {
	if dist.SampleCount() != sampleCount {
		t.Errorf("expected sample count to be %v but was %v", sampleCount, dist.SampleCount())
	}
	if math.Abs(dist.Min()-min) > within {
		t.Errorf("expected minimum to be +-%v from %v but was %v", within, min, dist.Min())
	}
	if math.Abs(dist.Max()-max) > within {
		t.Errorf("expected maximum to be +-%v from %v but was %v", within, max, dist.Max())
	}
	if math.Abs(dist.Avg()-float64(avg)) > within {
		t.Errorf("expected sample average to be +- %v from %v but was %v", within, avg, dist.Avg())
	}
	// variance and stddev are tested elsewere
}

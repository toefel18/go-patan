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
package metrics

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/toefel18/go-patan/metrics/common"
	"github.com/toefel18/go-patan/metrics/lockbased"
)

func Example() {
	sw := StartStopwatch()

	IncrementCounter("active.connections")
	AddSample("shopping.basket.total", 143.26)
	AddSample("shopping.basket.total", 167.16)
	AddSample("shopping.basket.total", 23.98)
	DecrementCounter("active.connections")
	RecordElapsedTime("order.processing.duration", sw)

	snapshot := SnapshotAndReset() //Returns a snapshot of all recorded counters, samples and durations and clears the store!

	_ = snapshot.Counters()["active.connections"] // 0    (1 increment, 1 decrement, counters start at 0)

	basketTotal := snapshot.Samples()["shopping.basket.total"] // baskedTotal = a distribution with a min, max, avg and stddev
	basketTotal.SampleCount()                                  // 3
	basketTotal.Avg()                                          // 111,46
	basketTotal.StdDev()                                       // ...

	processingDuration := snapshot.Durations()["order.processing.duration"] //  returns a distribution as well
	processingDuration.SampleCount()                                        // 1
	processingDuration.Avg()                                                // ...
	//...

	timeWindow := snapshot.CreatedTimestamp() - snapshot.StartedTimestamp() // time window in which measurements were recorded in millis
	_ = timeWindow
}

func TestConcurrency(t *testing.T) {
	runConcurrencyTest(10, 20000)
	runConcurrencyTest(100, 20000)
	runConcurrencyTest(500, 10000)
	runConcurrencyTest(10, 50000)
	runConcurrencyTest(100, 50000)
}

func TestHapppyFlow(t *testing.T) {
	Reset()
	RecordElapsedTime("record.elapsed.time", StartStopwatch())
	MeasureFunc("measure.func", func() {
		time.Sleep(200 * time.Millisecond)
	})
	MeasureFuncCanPanic("measure.func.safe", func() {
		time.Sleep(200 * time.Millisecond)
	})
	otherUnused := New()
	IncrementCounter("inc")
	DecrementCounter("dec")
	AddToCounter("add.to.counter", 9541)
	AddSample("sample", 5.0)

	ssNew := otherUnused.SnapshotAndReset() // SHOULD DO NOTHING, otherUnused is unrelated
	if len(ssNew.Counters())+len(ssNew.Durations())+len(ssNew.Samples()) != 0 {
		t.Error("an unrelated instance should not affect the global instance")
	}

	AddSample("sample", 10.0)

	ss := Snapshot()
	if len(ss.Counters()) != 3 {
		t.Error("expected 3 counters but got ", len(ss.Counters()))
	}
	if len(ss.Durations()) != 3 {
		t.Error("expected 2 durations but got ", len(ss.Durations()))
	}
	if len(ss.Samples()) != 1 {
		t.Error("expected 2 durations but got ", len(ss.Samples()))
	}
	xx := SnapshotAndReset()
	if len(xx.Counters()) == 0 || len(xx.Durations()) == 0 || len(ss.Samples()) == 0 {
		result, _ := json.Marshal(ss)
		t.Error("SnapshotAndReset should return the latest data inside the repository, but some items appear cleared state: " + string(result))
	}
	ss = Snapshot() // this snapshot should be cleared
	if len(ss.Counters())+len(ss.Durations())+len(ss.Samples()) != 0 {
		result, _ := json.Marshal(ss)
		t.Error("Snapshot after SnapshotAndReset() should be empty, but got: " + string(result))
	}

}

func BenchmarkTimer(b *testing.B) {
	timer := StartStopwatch()
	for i := 0; i < b.N; i++ {
		RecordElapsedTime("some.duration", timer)
	}
}

func BenchmarkAddToCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AddToCounter("some.counter", int64(i))
	}
}

func BenchmarkMeasureFuncCanPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MeasureFuncCanPanic("some.func", func() {
			return
		})
	}
}

func BenchmarkMeasureFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MeasureFunc("some.func", func() {
			return
		})
	}
}

func BenchmarkIncrementCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IncrementCounter("some.counter")
	}
}

func BenchmarkDecrementCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecrementCounter("some.counter")
	}
}

func BenchmarkSample(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AddSample("some.sample", float64(i))
	}
}

func BenchmarkAll(b *testing.B) {
	timer := StartStopwatch()
	for i := 0; i < b.N; i++ {
		MeasureFunc("some.func", func() {
			DecrementCounter("some.counter")
			RecordElapsedTime("some.duration", timer)
			AddSample("some.sample", float64(i))
			IncrementCounter("some.counter")
		})
	}
}

func BenchmarkAllParallel(b *testing.B) {
	timer := StartStopwatch()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			MeasureFunc("some.func", func() {
				DecrementCounter("some.counter")
				RecordElapsedTime("some.duration", timer)
				AddSample("some.sample", 333.333)
				IncrementCounter("some.counter")
			})
		}
	})
}

func runConcurrencyTest(threads int64, itemsPerThread int64) {
	millisStart := common.CurrentTimeMillis()
	wg := sync.WaitGroup{}
	subject := lockbased.NewFacade(lockbased.NewStore())
	for i := int64(0); i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sw := subject.StartStopwatch()
			defer subject.RecordElapsedTime("goroutine.duration", sw)
			for i := int64(0); i < itemsPerThread; i++ {
				subject.IncrementCounter("concurrency.counter")
				subject.AddSample("concurrency.sample", float64(i))
			}
		}()
	}
	wg.Wait()
	snapshot := subject.Snapshot()
	expectedItems := threads * itemsPerThread
	if snapshot.Counters()["concurrency.counter"] != expectedItems {
		panic(fmt.Sprint(expectedItems, "counters expected, but got", snapshot.Counters()["concurrency.counter"]))
	}
	if snapshot.Durations()["goroutine.duration"].SampleCount() != threads {
		panic(fmt.Sprint("There should be", threads, "durations registered but got", snapshot.Durations()["goroutine.duration"].SampleCount()))
	}
	if snapshot.Samples()["concurrency.sample"].SampleCount() != threads*itemsPerThread {
		panic(fmt.Sprint(expectedItems, "samples expected but got", snapshot.Samples()["concurrency.sample"].SampleCount()))
	}
	millisEnd := common.CurrentTimeMillis()
	fmt.Println(threads, "threads with", itemsPerThread, "items took", (millisEnd - millisStart))
}

func TestEmptyCounters(t *testing.T) {
	Reset()
	snapshot := Snapshot()
	if snapshot.Counters()["nonexistingcounter"] != 0 {
		t.Error("nonexisting counter should be 0")
	}
	_, exists := snapshot.Durations()["nonexistingduration"]
	if exists == true {
		t.Error("nonexisting duration should not exist")
	}
	_, exists2 := snapshot.Samples()["nonexistingsample"]
	if exists2 == true {
		t.Error("nonexisting duration should not exist")
	}
}

func TestJsonize(t *testing.T) {
	Reset()
	sw := StartStopwatch()
	IncrementCounter("json.counter")
	AddSample("json.sample", 15)
	AddSample("json.sample", 25)
	time.Sleep(20 * time.Millisecond)
	RecordElapsedTime("json.duration", sw)
	time.Sleep(20 * time.Millisecond)
	snapshot := Snapshot()
	json, err := json.Marshal(snapshot)
	if err == nil {
		jsonString := string(json)
		if !strings.Contains(jsonString, "json.duration") ||
			!strings.Contains(jsonString, "timestampStarted") ||
			!strings.Contains(jsonString, "timestampTaken") ||
			!strings.Contains(jsonString, "sampleCount") ||
			!strings.Contains(jsonString, "minimum") ||
			!strings.Contains(jsonString, "maximum") ||
			!strings.Contains(jsonString, "mean") ||
			!strings.Contains(jsonString, "stdDeviation") {
			t.Error("the json output does not contain some of the expected values")
		}
	} else {
		t.Error("json marshalling failed", err)
	}
}

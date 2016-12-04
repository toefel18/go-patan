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

func TestConcurrency(t *testing.T) {
	Benchmrk(10, 20000)
	Benchmrk(100, 20000)
	Benchmrk(500, 10000)
	Benchmrk(10, 50000)
	Benchmrk(100, 50000)
}

func Benchmrk(threads int64, itemsPerThread int64) {
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
	fmt.Println(string(json))
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

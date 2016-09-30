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
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {
	IncrementCounter("counter.1")
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sw := StartStopwatch()
			defer RecordElapsedTime("goroutine.duration", sw)
			for i := 0; i < 200000; i++ {
				IncrementCounter("concurrency.counter")
				AddSample("concurrency.sample", int64(i))
			}
			fmt.Println("done!")
		}()
	}
	wg.Wait()
	fmt.Println("goroutines done")
	time.Sleep(20 * time.Millisecond)
	snapshot := Snapshot()
	if snapshot.Counters()["concurrency.counter"] != 2000000 {
		t.Error("Counter should be 2000000 but was", snapshot.Counters()["concurrency.counter"])
	}
	if snapshot.Durations()["goroutine.duration"].SampleCount() != 10 {
		t.Error("There should be 10 durations registered but was", snapshot.Durations()["goroutine.duration"].SampleCount())
	}
	if snapshot.Samples()["concurrency.sample"].SampleCount() != 2000000 {
		t.Error("There should be 2000000 samples but got", snapshot.Samples()["concurrency.sample"].SampleCount())
	}
	Reset()
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
			!strings.Contains(jsonString, "sampleCount") ||
			!strings.Contains(jsonString, "minimum") ||
			!strings.Contains(jsonString, "maximum") ||
			!strings.Contains(jsonString, "average") ||
			!strings.Contains(jsonString, "standardDeviation") {
			t.Error("the json output does not contain some of the expected values")
		}
	} else {
		t.Error("json marshalling failed", err)
	}
}

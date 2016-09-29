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
package patan

import (
	"fmt"
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
}

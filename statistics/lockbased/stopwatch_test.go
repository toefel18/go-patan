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
	"fmt"
	"testing"
	"time"
	"github.com/toefel18/go-patan/statistics/api"
)

func TestStartStopwatch(t *testing.T) {
	sw := startNewStopwatch()
	elapsedMillis := sw.ElapsedMillis()
	elapsedNanos := sw.ElapsedNanos()
	if elapsedMillis > 10 {
		t.Errorf("newly created stopwatch already has already %v millis elapsed, could indicate programming error", elapsedMillis)
	} else if elapsedNanos > (10 * time.Millisecond).Nanoseconds() {
		t.Errorf("newly created stopwatch already has already %v nanos elapsed, could indicate programming error", elapsedNanos)
	}

	time.Sleep(100 * time.Millisecond)

	elapsedMillis = sw.ElapsedMillis()
	elapsedNanos = sw.ElapsedNanos()

	if elapsedMillis > 105 || elapsedMillis < 100 {
		t.Errorf("stopwatch elapsed %v millis, expected to be in range 100 to 105 millis", elapsedMillis)
	} else if elapsedNanos > (105*time.Millisecond).Nanoseconds() || elapsedNanos < (100*time.Millisecond).Nanoseconds() {
		t.Errorf("stopwatch elapsed %v nanos, expected to be in range 100000000 - 105000000", elapsedNanos)
	}
	fmt.Println(elapsedMillis, elapsedNanos)
}

func TestStopwatchImplementsApiInterface(t *testing.T) {
	var channelbasedSw *Stopwatch = startNewStopwatch()
	var apiSw api.Stopwatch = channelbasedSw
	if apiSw.ElapsedMillis() > 100 {
		t.Error("channelbased.Stopwatch has problems implementing api.Stopwatch")
	}
}
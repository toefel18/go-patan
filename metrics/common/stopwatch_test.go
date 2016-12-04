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
	"time"

	"github.com/toefel18/go-patan/metrics/api"
)

func TestStartStopwatch(t *testing.T) {
	sw := StartNewStopwatch()
	elapsedMillis := sw.ElapsedMillis()
	if elapsedMillis > 1 {
		t.Errorf("newly created stopwatch already has already %v millis elapsed, could indicate programming error", elapsedMillis)
	}

	time.Sleep(100 * time.Millisecond)

	elapsedMillis = sw.ElapsedMillis()

	if elapsedMillis > 105.0 || elapsedMillis < 100.0 {
		t.Errorf("stopwatch elapsed %v millis, expected to be in range 100 to 105 millis", elapsedMillis)
	}
}

func TestStopwatchImplementsApiInterface(t *testing.T) {
	sw := StartNewStopwatch()
	var apiSw api.Stopwatch = sw
	if apiSw.ElapsedMillis() > 100 {
		t.Error("channelbased.Stopwatch has problems implementing api.Stopwatch")
	}
}

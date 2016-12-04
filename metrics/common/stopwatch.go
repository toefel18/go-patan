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
	"time"
)

// Stopwatch records how much time has elapsed since it's creation.
type Stopwatch struct {
	time.Time
}

// StartNewStopwatch creates a new Stopwatch. Stopwatches start immediatlly once created.
func StartNewStopwatch() *Stopwatch {
	return &Stopwatch{time.Now()}
}

// ElapsedMillis contains the milliseconds elapsed since it's creation. The return value is
// a float, which has nanosecond accuracy.
func (sw *Stopwatch) ElapsedMillis() float64 {
	return float64(time.Now().Sub(sw.Time).Nanoseconds()) / float64(time.Millisecond.Nanoseconds())
}

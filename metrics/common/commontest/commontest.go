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

package commontest

import (
	"math"
	"testing"

	"github.com/toefel18/go-patan/metrics/api"
)

//AssertDistributionHasValues checks if the distribution contains the expected values
func AssertDistributionHasValues(dist api.Distribution, sampleCount int64, min, max, avg, stdDev float64, t *testing.T) {
	if dist == nil {
		t.Error("dist is nil, the value is not present (look for a Reset() or SnapshotAndReset() on the channelbased)")
		return
	}
	if dist.SampleCount() != sampleCount {
		t.Errorf("expected sample count to be %v but was %v", sampleCount, dist.SampleCount())
	}
	if !FloatEquals(dist.Min(), min) {
		t.Errorf("expected minimum to be %v but was %v", min, dist.Min())
	}
	if !FloatEquals(dist.Max(), max) {
		t.Errorf("expected maximum to be %v but was %v", max, dist.Max())
	}
	if !FloatEquals(dist.Avg(), avg) {
		t.Errorf("expected sample average to be %v but was %v", avg, dist.Avg())
	}
	if !FloatEquals(dist.StdDev(), stdDev) {
		t.Errorf("expected sample variance to be %v but was %v", stdDev, dist.StdDev())
	}
}

// Epsilon is the required precision in tests
var Epsilon = 0.001

// FloatEquals compares two floats and returns true if they are close enougth
func FloatEquals(a, b float64) bool {
	return math.Abs(b-a) < Epsilon
}

// CloseTo tests if b is close to a with the given offset
func CloseTo(a, b, offset int64) bool {
	return b > a-offset && b < a+offset
}

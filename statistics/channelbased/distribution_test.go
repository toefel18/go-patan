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
package channelbased

import (
	"math"
	"testing"
	"github.com/toefel18/go-patan/statistics/api"
)

func TestNewDistribution(t *testing.T) {
	var dist *Distribution
	dist = NewDistribution()
	assertDistributionHasValues(dist, 0, math.MaxInt64, math.MinInt64, 0.0, 0.0, 0.0, t)
}

func TestAddSample(t *testing.T) {
	dist := NewDistribution()
	dist.addSample(10)
	assertDistributionHasValues(dist, 1, 10, 10, 10.0, 0.0, math.NaN(), t)
	dist.addSample(20)
	assertDistributionHasValues(dist, 2, 10, 20, 15.0, 50.0, 7.0710, t)
	dist.addSample(0)
	assertDistributionHasValues(dist, 3, 0, 20, 10.0, 50.0, 10.0, t)
}

func assertDistributionHasValues(dist api.Distribution, sampleCount, min, max int64, avg, variance, stdDev float64, t *testing.T) {
	if dist.SampleCount() != sampleCount {
		t.Errorf("expected sample count to be %v but was %v", sampleCount, dist.SampleCount())
	}
	if dist.Min() != min {
		t.Errorf("expected minimum to be %v but was %v", min, dist.Min())
	}
	if dist.Max() != max {
		t.Errorf("expected maximum to be %v but was %v", max, dist.Max())
	}
	if !floatEquals(dist.Avg(), avg) {
		t.Errorf("expected sample average to be %v but was %v", avg, dist.Avg)
	}
	if !floatEquals(dist.Variance(), variance) {
		t.Errorf("expected sample variance to be %v but was %v", variance, dist.Variance())
	}
	if !(math.IsNaN(stdDev) && math.IsNaN(dist.StdDev())) && !floatEquals(stdDev, dist.StdDev()) {
		t.Errorf("expected sample std dev to be %v but was %v", stdDev, dist.StdDev())
	}
}

func TestDistributionImplementsInterface(t *testing.T) {
	var apiDist api.Distribution
	var channelbasedDist *Distribution = NewDistribution()
	apiDist = channelbasedDist
	if apiDist.SampleCount() != 0 {
		t.Error("channelbased.Distribution has problems implementing api.Distribution interface")
	}
}

var EPSILON float64 = 0.0001

func floatEquals(a, b float64) bool {
	if (b - a) < EPSILON {
		return true
	}
	return false
}

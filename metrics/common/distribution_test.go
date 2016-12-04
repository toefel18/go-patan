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
	"math"
	"testing"

	"github.com/toefel18/go-patan/metrics/api"
	"github.com/toefel18/go-patan/metrics/common/commontest"
)

func TestNewDistribution(t *testing.T) {
	var dist *Distribution
	dist = NewDistribution()
	commontest.AssertDistributionHasValues(dist, 0, math.MaxFloat64, math.SmallestNonzeroFloat64, 0.0, 0.0, t)
}

func TestDistributionAddSample1To10(t *testing.T) {
	dist := NewDistribution()
	for i := 1; i <= 10; i++ {
		dist.AddSample(float64(i))
	}
	expDeviation := math.Sqrt((2*4.5*4.5 + 2*3.5*3.5 + 2*2.5*2.5 + 2*1.5*1.5 + 2*0.5*0.5) / 9)
	commontest.AssertDistributionHasValues(dist, 10, 1.0, 10, 5.5, expDeviation, t)
}

func TestDistributionAddSample10To1(t *testing.T) {
	dist := NewDistribution()
	for i := 10; i >= 1; i-- {
		dist.AddSample(float64(i))
	}
	expDeviation := math.Sqrt((2*4.5*4.5 + 2*3.5*3.5 + 2*2.5*2.5 + 2*1.5*1.5 + 2*0.5*0.5) / 9)
	commontest.AssertDistributionHasValues(dist, 10, 1.0, 10, 5.5, expDeviation, t)
}

func TestDistributionImplementsInterface(t *testing.T) {
	var apiDist api.Distribution
	dist := NewDistribution()
	apiDist = dist
	if apiDist.SampleCount() != 0 {
		t.Error("common.Distribution has problems implementing api.Distribution interface")
	}
}

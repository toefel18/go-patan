package concrete

import (
	"github.com/toefel18/go-patan/statistics"
	"math"
	"testing"
)

func TestNewDistribution(t *testing.T) {
	var dist statistics.Distribution
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

func assertDistributionHasValues(dist statistics.Distribution, sampleCount, min, max int64, avg, variance, stdDev float64, t *testing.T) {
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

var EPSILON float64 = 0.0001

func floatEquals(a, b float64) bool {
	if (b - a) < EPSILON {
		return true
	}
	return false
}

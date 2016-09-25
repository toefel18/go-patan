package concrete

import (
    "math"
    "testing"
)

func TestNewDistribution(t *testing.T) {
    dist := NewDistribution()
    if (dist.Minimum != math.MinInt64){
        t.Errorf("expected minimum to be MIN_INT but was %v", dist.Minimum)
    }
    if (dist.Maximum != math.MaxInt64){
        t.Errorf("expected maximum to be MAX_INT but was %v", dist.Maximum)
    }
    if (dist.SampleCount != 0){
        t.Errorf("expected sample count to be 0 but was %v", dist.SampleCount)
    }
    if (dist.SampleAverage != 0.0){
        t.Errorf("expected sample average to be 0.0 but was %v", dist.SampleAverage)
    }
    if (dist.SampleVariance != 0.0){
        t.Errorf("expected sample variance to be 0.0 but was %v", dist.SampleVariance)
    }
    if (dist.SampleStdDeviation != 0.0){
        t.Errorf("expected sample std dev to be 0.0 but was %v", dist.SampleStdDeviation)
    }
}
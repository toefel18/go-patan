package channelbased

import (
	"math"
)

type Distribution struct {
	sampleCount int64   `json:"sampleCount"`
	min         int64   `json:"minimum"`
	max         int64   `json:"maximum"`
	avg         float64 `json:"sampleAverage"`
	variance    float64 `json:"sampleVariance"`
	stdDev      float64 `json:"sampleStdDeviation"`
}

func NewDistribution() *Distribution {
	return &Distribution{min: math.MaxInt64, max: math.MinInt64}
}

func (dist *Distribution) SampleCount() int64 {
	return dist.sampleCount
}

func (dist *Distribution) Min() int64 {
	return dist.min
}

func (dist *Distribution) Max() int64 {
	return dist.max
}

func (dist *Distribution) Avg() float64 {
	return dist.avg
}

func (dist *Distribution) Variance() float64 {
	return dist.variance
}

func (dist *Distribution) StdDev() float64 {
	return dist.stdDev
}

func (dist *Distribution) addSample(value int64) {
	updatedSampleCount := dist.sampleCount + 1
	updatedMin := min(value, dist.min)
	updatedMax := max(value, dist.max)
	updatedAvg := dist.avg + ((float64(value) - dist.avg) / float64(updatedSampleCount))
	updatedVar := dist.variance + ((float64(value) - dist.avg) * (float64(value) - updatedAvg))
	updatedStdDev := math.Sqrt(updatedVar / float64(dist.sampleCount))
	dist.sampleCount = updatedSampleCount
	dist.min = updatedMin
	dist.max = updatedMax
	dist.avg = updatedAvg
	dist.variance = updatedVar
	dist.stdDev = updatedStdDev
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

package channelbased

import (
	"math"
)

type Distribution struct {
	Samples       int64   `json:"sampleCount"`
	Minimum       int64   `json:"minimum"`
	Maximum       int64   `json:"maximum"`
	Average       float64 `json:"average"`
	totalVariance float64  // since this value is not useful to expose
	StdDeviation  float64 `json:"standardDeviation"`
}

func NewDistribution() *Distribution {
	return &Distribution{Minimum: math.MaxInt64, Maximum: math.MinInt64}
}

func (dist *Distribution) SampleCount() int64 {
	return dist.Samples
}

func (dist *Distribution) Min() int64 {
	return dist.Minimum
}

func (dist *Distribution) Max() int64 {
	return dist.Maximum
}

func (dist *Distribution) Avg() float64 {
	return dist.Average
}

func (dist *Distribution) Variance() float64 {
	return dist.totalVariance
}

func (dist *Distribution) StdDev() float64 {
	return dist.StdDeviation
}

func (dist *Distribution) addSample(value int64) {
	updatedSampleCount := dist.Samples + 1
	updatedMin := min(value, dist.Minimum)
	updatedMax := max(value, dist.Maximum)
	updatedAvg := dist.Average + ((float64(value) - dist.Average) / float64(updatedSampleCount))
	updatedVar := dist.totalVariance + ((float64(value) - dist.Average) * (float64(value) - updatedAvg))
	updatedStdDev := math.Sqrt(updatedVar / float64(dist.Samples))
	dist.Samples = updatedSampleCount
	dist.Minimum = updatedMin
	dist.Maximum = updatedMax
	dist.Average = updatedAvg
	dist.totalVariance = updatedVar
	dist.StdDeviation = updatedStdDev
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

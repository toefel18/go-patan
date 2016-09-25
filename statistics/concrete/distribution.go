package concrete

import "math"

type Distribution struct {
    SampleCount        int64   `json:"sampleCount"`
    Minimum            int64   `json:"minimum"`
    Maximum            int64   `json:"maximum"`
    SampleAverage      float64 `json:"sampleAverage"`
    SampleVariance     float64 `json:"sampleVariance"`
    SampleStdDeviation float64 `json:"sampleStdDeviation"`
}

func (dist *Distribution) GetSampleCount() int64 {
    return dist.SampleCount
}

func (dist *Distribution) GetMinimum() int64 {
    return dist.Minimum
}

func (dist *Distribution) GetMaximum() int64 {
    return dist.Maximum
}

func (dist *Distribution) GetSampleAverage() float64 {
    return dist.SampleAverage
}

func (dist *Distribution) GetSampleVariance() float64 {
    return dist.SampleVariance
}

func (dist *Distribution) GetSampleStdDeviation() float64 {
    return dist.SampleStdDeviation
}

func (dist *Distribution) addSample(value int64) {

}

func NewDistribution() *Distribution {
    return &Distribution{Minimum: math.MinInt64,Maximum: math.MaxInt64}
}
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
)

// Distribution contains a summarized view of a statistical distribution
type Distribution struct {
	Samples       int64   `json:"sampleCount"`
	Minimum       float64 `json:"minimum"`
	Maximum       float64 `json:"maximum"`
	Mean          float64 `json:"mean"`
	totalVariance float64 // since this value is not useful to expose
	StdDeviation  float64 `json:"stdDeviation"`
}

// NewDistribution creates a new initialized distribution with
func NewDistribution() *Distribution {
	return &Distribution{Minimum: math.MaxFloat64, Maximum: math.SmallestNonzeroFloat64}
}

// SampleCount returns the number of samples taken
func (dist *Distribution) SampleCount() int64 {
	return dist.Samples
}

// Min returns the minimum value recorded
func (dist *Distribution) Min() float64 {
	return dist.Minimum
}

// Max returns the maximum value recorded
func (dist *Distribution) Max() float64 {
	return dist.Maximum
}

// Avg returns the average value (json marshalled as mean)
func (dist *Distribution) Avg() float64 {
	return dist.Mean
}

// StdDev returns the standad deviation of the distribution.
func (dist *Distribution) StdDev() float64 {
	return dist.StdDeviation
}

// AddSample updates the distribution to contain the value
func (dist *Distribution) AddSample(value float64) {
	updatedSampleCount := dist.Samples + 1
	updatedMin := min(value, dist.Minimum)
	updatedMax := max(value, dist.Maximum)
	updatedAvg := dist.Mean + ((float64(value) - dist.Mean) / float64(updatedSampleCount))
	updatedVar := dist.totalVariance + ((float64(value) - dist.Mean) * (float64(value) - updatedAvg))
	updatedStdDev := math.Sqrt(updatedVar / float64(dist.Samples))

	dist.Samples = updatedSampleCount
	dist.Minimum = updatedMin
	dist.Maximum = updatedMax
	dist.Mean = updatedAvg
	dist.totalVariance = updatedVar
	if math.IsNaN(updatedStdDev) {
		updatedStdDev = 0.0 //if there's only one value, the stdDeviation should be 0
	} else {
		dist.StdDeviation = updatedStdDev
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

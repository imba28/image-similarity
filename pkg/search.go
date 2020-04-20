package pkg

import (
	"math"
	"sort"
)

type ImageDistance struct {
	Image    Image
	Distance float64
}

type ByDistance []ImageDistance

func (a ByDistance) Len() int {
	return len(a)
}
func (a ByDistance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByDistance) Less(i, j int) bool {
	return a[i].Distance < a[j].Distance
}

func CalculateDistances(referencePath Image, vectors []ImageDescriptor) ([]ImageDistance, error) {
	referenceVector, err := FeatureVector(referencePath)
	if err != nil {
		return nil, err
	}

	var d []ImageDistance

	for _, vector := range vectors {
		d = append(d, ImageDistance{
			Image:    vector.Image,
			Distance: chi2Distance(referenceVector, vector.Features),
		})
	}

	sort.Sort(ByDistance(d))

	return d, nil
}

func chi2Distance(v1, v2 []float64) float64 {
	d := 0.
	min := int(math.Min(float64(len(v1)), float64(len(v2))))

	for i := 0; i < min; i++ {
		distance := math.Pow(v1[i]-v2[i], 2) / (v1[i] + v2[i] + 1e-10)
		d += distance
	}

	return d * 0.5
}

package pkg

import (
	"errors"
	"math"
	"runtime"
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

func calculateRange(reference *Image, images []*Image, from int, to int, dc chan<- *ImageDistance) {
	for i := from; i < to; i++ {
		dc <- &ImageDistance{
			Image:    *images[i],
			Distance: chi2Distance(reference.Features, images[i].Features),
		}
	}
}

func CalculateDistances(reference Image, images []*Image) ([]ImageDistance, error) {
	var d []ImageDistance
	distanceChannel := make(chan *ImageDistance, 20)
	defer close(distanceChannel)

	availableCores := float64(runtime.NumCPU())
	logLen := math.Ceil(math.Log10(float64(len(images))))
	routineCount := int(math.Min(logLen, availableCores))
	routineCalculationRange := len(images) / routineCount

	for i := 0; i < routineCount; i++ {
		go calculateRange(&reference, images, i*routineCalculationRange, (i+1)*routineCalculationRange, distanceChannel)
	}

	for distance := range distanceChannel {
		d = append(d, *distance)
		if len(d) == len(images) {
			sort.Sort(ByDistance(d))
			return d, nil
		}
	}

	return nil, errors.New("shit, this should not have happened")
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

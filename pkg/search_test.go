package pkg

import (
	"sort"
	"testing"
)

func TestUnitCalculateDistances(t *testing.T) {
	images := []*Image{
		{
			Id:       "1",
			Guid:     1,
			Path:     "test1.png",
			Name:     "test1.png",
			Features: []float64{1, 1, 1, 1, 1},
		},
		{
			Id:       "2",
			Guid:     2,
			Path:     "test2.png",
			Name:     "test2.png",
			Features: []float64{1, 1, 1, 1, 0},
		},
		{
			Id:       "3",
			Guid:     3,
			Path:     "test3.png",
			Name:     "test3.png",
			Features: []float64{1, 0, 0, 0, 0},
		},
		{
			Id:       "4",
			Guid:     4,
			Path:     "test4.png",
			Name:     "test4.png",
			Features: []float64{0, 0, 0, 0, 1},
		},
		{
			Id:       "5",
			Guid:     5,
			Path:     "test5.png",
			Name:     "test5.png",
			Features: []float64{1, 0, 0, 0, 1},
		},
	}

	ds, err := CalculateDistances(*images[0], images)
	if err != nil {
		t.Fatalf("should not return error, got: %v", err)
	}
	if ds[0].Image.Id != images[0].Id {
		t.Errorf("expected the reference image to have the lowest distance to itself, want: %s, got: %s", images[0].Id, ds[0].Image.Id)
	}
	if ds[1].Image.Id != images[1].Id {
		t.Errorf("expected the second image to have the lowest distance to the reference image, want: %s, got: %s", images[1].Id, ds[1].Image.Id)
	}

	var distances []float64
	for i := range ds {
		distances = append(distances, ds[i].Distance)
	}

	if !sort.Float64sAreSorted(distances) || (distances != nil && distances[0] > distances[len(distances)-1]) {
		t.Error("image distance slice should be sorted in ascending order")
	}
}

func TestIntegrationCalculateDistances__calculate_vector(t *testing.T) {
	images := []*Image{
		{
			Path: "../test_sets/sea/beach-2179183_640.jpg",
		},
		{
			Path: "../test_sets/sea/beach-1852945_640.jpg",
		},
		{
			Path: "../test_sets/salzburg/salzburg-116768_640.jpg",
		},
	}
	for i := range images {
		f, _ := FeatureVector(*images[i])
		images[i].Features = f
	}

	ds, err := CalculateDistances(*images[0], images)
	if err != nil {
		t.Fatalf("should not return error and calucate feature vectors, got: %v", err)
	}

	if len(ds) != len(images) {
		t.Errorf("Wrong length of image distance slice, want: %d, got: %d", len(images), len(ds))
	}

	for i := range images {
		// all images following the first (=reference) should have a distance greater than 0
		if i > 0 && ds[i].Distance <= 0 {
			t.Errorf("incorrect image distance, want : > 0, got: %f", ds[i].Distance)
		}
	}
}

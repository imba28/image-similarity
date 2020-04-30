package pkg

import "testing"

func TestIntegrationFeatureVector(t *testing.T) {
	f, err := FeatureVector(Image{
		Path: "../test_sets/bikes/bike-4500339_640.jpg",
		Name: "bike-4500339_640.jpg",
	})
	if err != nil {
		t.Fatalf("should calulate feature vector but got error instead, got: %v", err)
	}
	if len(f) != hueBins*saturationBins*brightnessBins*5 {
		t.Errorf("Incorrect length of feature vector, want: %d, got: %d", hueBins*saturationBins*brightnessBins*5, len(f))
	}

	_, err = FeatureVector(Image{
		Path: "../i/dont/exist.jpg",
		Name: "exist.jpg",
	})
	if err == nil {
		t.Errorf("should return error if image does not exist, got: %v", err)
	}
}

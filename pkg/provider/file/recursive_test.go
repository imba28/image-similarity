package file

import (
	"imba28/images/pkg"
	"testing"
)

func TestUnitRecursiveImageProvider_Images(t *testing.T) {
	provider := RecursiveImageProvider{Dir: "../../../test_sets"}

	files, err := provider.Images()
	if err != nil {
		t.Error("SingleLevelImageProvider should not return error")
	}
	if len(files) != 7 {
		t.Errorf("Length of files in dir %s incorrect, got: %d, want: %d", provider.Dir, len(files), 7)
	}
}

func TestUnitRecursiveImageProvider_Get(t *testing.T) {
	provider := RecursiveImageProvider{Dir: "../../../test_sets"}
	if image := provider.Get("foobar.png"); image != nil {
		t.Errorf("should not return an image that does not exist, got: %v, want: %v", image, nil)
	}
	if image := provider.Get("README.md"); image == nil {
		t.Errorf("should return an image that exists, got: %v, want: %v", image, pkg.Image{
			Path: "../../../test_sets/README.md",
			Name: "README.md",
		})
	}
	if image := provider.Get("bikes/bike-4500339_640.jpg"); image == nil {
		t.Errorf("should return an image that lives inside a subdirectory, got: %v, want: %v", image, pkg.Image{
			Path: "../../../test_sets/bikes/bike-4500339_640.jpg",
			Name: "bikes/bike-4500339_640.jpg",
		})
	}
}

func TestUnitRecursiveImageProvider_Persist(t *testing.T) {
	provider := RecursiveImageProvider{Dir: "../../../test_sets"}

	if err := provider.Persist(&pkg.Image{
		Id:       "",
		Guid:     0,
		Path:     "",
		Name:     "",
		Features: nil,
	}); err != nil {
		t.Error("File provider should never return an error as it is a noop.")
	}

	if err := provider.Persist(&pkg.Image{
		Id:       "123",
		Guid:     1,
		Path:     "123",
		Name:     "123",
		Features: nil,
	}); err != nil {
		t.Error("File provider should never return an error as it is a noop.")
	}

	if err := provider.Persist(&pkg.Image{
		Features: []float64{1, 2, 3, 4, 5},
	}); err != nil {
		t.Error("File provider should never return an error as it is a noop.")
	}

	if err := provider.Persist(nil); err != nil {
		t.Error("File provider should never return an error as it is a noop.")
	}
}

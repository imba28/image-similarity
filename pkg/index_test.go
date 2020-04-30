package pkg

import (
	"reflect"
	"testing"
)

type testImageProvider struct {
}

func (t testImageProvider) Images() ([]*Image, error) {
	return []*Image{
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
	}, nil
}

func (t testImageProvider) Get(id string) *Image {
	return &Image{
		Id:       id,
		Guid:     1,
		Path:     "test.png",
		Name:     "test.png",
		Features: []float64{1, 2, 3, 4, 5},
	}
}

func (t testImageProvider) Persist(*Image) error {
	return nil
}

var _ ImageProvider = (*testImageProvider)(nil)

func TestUnitImageIndex(t *testing.T) {
	index, err := NewIndex(testImageProvider{})
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	if len(index.images) != 3 {
		t.Errorf("Expected length of ndex to match length of its provider, expected : %d, got: %d", 3, len(index.images))
	}
}

func TestUnitImageIndex_Add(t *testing.T) {
	index, err := NewIndex(testImageProvider{})
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	err = index.Add(&Image{
		Id:       "4",
		Guid:     4,
		Path:     "test4.png",
		Name:     "test4.png",
		Features: []float64{1, 2, 3, 4, 5},
	})
	if err != nil {
		t.Errorf("Add should not return error, got: %v", err)
	}
	if index.images[len(index.images)-1].Id != "4" || len(index.images) != 4 {
		t.Errorf("Expected Add to append the image to the list")
	}

	if err = index.Add(nil); err != nil {
		t.Errorf("Add should not return if argument is nil")
	}
}

func TestUnitImageIndex_Add__existing(t *testing.T) {
	index, err := NewIndex(testImageProvider{})
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	err = index.Add(&Image{
		Id:       "1",
		Guid:     1,
		Path:     "test1.png",
		Name:     "test1.png",
		Features: []float64{1, 2, 3, 4, 5},
	})
	if err != nil {
		t.Error("Add should not return error, because image already exists")
	}
	if len(index.images) != 3 {
		t.Error("Expected Add to not append the image to the list because it is already contained")
	}
}

func TestUnitImageIndex_Add__calculate_feature_vector(t *testing.T) {
	index, err := NewIndex(testImageProvider{})
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	image := Image{
		Id:   "6",
		Guid: 6,
		Path: "../test_sets/bikes/dirt-bike-828644_640.jpg",
		Name: "dirt-bike-828644_640.jpg",
	}
	err = index.Add(&image)
	if err != nil {
		t.Error("should not return error")
	}
	if image.Features == nil {
		t.Error("should update feature vector after adding the image")
	} else if len(image.Features) == 0 {
		t.Errorf("Incorrect length of feature vectors, got: %d", len(image.Features))
	}

	err = index.Add(&Image{
		Id:   "7",
		Guid: 7,
		Path: "foobar/test7.png",
		Name: "test7.png",
	})
	if err == nil {
		t.Error("should return error, because the feature vector of the image cannot be calculated")
	}
}

func TestUnitImageIndex_Load(t *testing.T) {
	p := testImageProvider{}
	index, err := NewIndex(p)
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	i := index.Load("2")
	expectedI := p.Get("2")

	if !reflect.DeepEqual(i, expectedI) {
		t.Error("index should load images from its provider")
	}
}

func TestUnitImageIndex_Images(t *testing.T) {
	p := testImageProvider{}
	index, err := NewIndex(p)
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	expectedImages, _ := p.Images()
	images := index.Images()

	if !reflect.DeepEqual(expectedImages, images) {
		t.Error("Index should load images from its provider")
	}
}

func TestUnitImageIndex_Has(t *testing.T) {
	p := testImageProvider{}
	index, err := NewIndex(p)
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	images, _ := p.Images()
	for _, image := range images {
		if !index.Has(image.Id) {
			t.Errorf("index should contain image with id %s", image.Id)
		}
	}
}

func TestUnitImageIndex_Search(t *testing.T) {
	images, _ := testImageProvider{}.Images()
	index, err := NewIndex(testImageProvider{})
	if err != nil {
		t.Fatalf("Index constructor should not return error, got: %v", err)
	}

	ds, err := index.Search(*images[0], 2, 5)
	if len(ds) == 0 {
		t.Fatalf("index search should return at least one result. got: 0")
	}
	if !reflect.DeepEqual(*images[0], ds[0].Image) {
		t.Errorf("the reference image should have the smallest distance and therefor be placed at index 0")
	}
	if !reflect.DeepEqual(*images[1], ds[1].Image) {
		t.Errorf("image 2 should be the most similar to image 1, want: %d, got: %s", 2, ds[1].Image.Id)
	}

	ds, err = index.Search(*images[2], 1, 5)
	if len(ds) == 0 {
		t.Errorf("index search with image 3 as reference and low threshold should return one result. got: %d, want: 1", len(ds))
	}

	ds, err = index.Search(*images[2], 10, 5)
	if len(ds) != 3 {
		t.Errorf("index search with image 3 as reference and high threshold should return all results. got: %d, want: %d", len(ds), len(images))
	}
}

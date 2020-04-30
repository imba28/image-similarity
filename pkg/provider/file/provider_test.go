package file

import (
	"imba28/images/pkg"
	"testing"
)

func TestUnitImageProvider_Get(t *testing.T) {
	provider := New("../../../template")
	image := provider.Get("index.html")

	if image == nil {
		t.Errorf("should have returned image, got: %v", nil)
		return
	}

	if image.Name != "index.html" {
		t.Errorf("incorrect image name, got: %q, want: %q", image.Name, "index.html")
	}
	if image.Path != "../../../template/index.html" {
		t.Errorf("incorrect image path, got: %q, want: %q", image.Path, "../../../template/index.html")
	}
	if image.Guid != 0 {
		t.Errorf("incorrect image guid, got: %d, want: %d", image.Guid, 0)
	}
	if image.Features != nil {
		t.Errorf("incorrect image feature vector, got: %v, want: %v", image.Features, nil)
	}
	if image.Id != "index.html" {
		t.Errorf("incorrect image name, got: %q, want: %q", image.Name, "index.html")
	}
}

func TestUnitImageProvider_Images(t *testing.T) {
	provider := New("../../../template")

	files, err := provider.Images()
	if err != nil {
		t.Error("ImageProvider should not return error")
	}
	if len(files) != 2 {
		t.Errorf("Length of files incorrect, got: %d, want: %d", len(files), 2)
	}
}

func TestUnitImageProvider_Images_not_existing(t *testing.T) {
	provider := New("../foobar")
	_, err := provider.Images()
	if err == nil {
		t.Errorf("if directory does not exist Images() should return an error")
	}
}

func TestUnitImageProvider_Images_name(t *testing.T) {
	provider := New("../../../template")

	files, err := provider.Images()
	expectedNames := []string{"index.html", "similar.html"}

	if err != nil {
		t.Error("ImageProvider should not return error")
	}
	for i := range files {
		if files[i].Name != expectedNames[i] {
			t.Errorf("Name of %dth file incorrect, got: %s, want %s", i, files[i].Name, expectedNames[i])
		}
	}
}

func TestUnitImageProvider_Images_path(t *testing.T) {
	provider := New("../../../template")

	files, err := provider.Images()
	expectedPaths := []string{"../../../template/index.html", "../../../template/similar.html"}

	if err != nil {
		t.Error("ImageProvider should not return error")
	}
	for i := range files {
		if files[i].Path != expectedPaths[i] {
			t.Errorf("Path of %dth file incorrect, got: %s, want %s", i, files[i].Path, expectedPaths[i])
		}
	}
}

func TestUnitImageProvider_Images_hidden_files(t *testing.T) {
	provider := New("../../../") // project root dir

	files, err := provider.Images()
	if err != nil {
		t.Error("ImageProvider should not return error")
	}
	if len(files) != 5 {
		t.Errorf("Length of files non hidden files incorrect, got: %d, want: %d", len(files), 5)
	}
}

func TestUnitImageProvider_Persist(t *testing.T) {
	provider := New("../../../") // project root dir
	image := pkg.Image{
		Id:       "2",
		Guid:     2,
		Path:     "foo/bar.png",
		Name:     "",
		Features: []float64{1, 2, 3},
	}
	if err := provider.Persist(&image); err != nil {
		t.Errorf("file provider should never return an error as this is a noop operation, got: %q", err)
	}
}

func TestUnitNewImage(t *testing.T) {
	image := NewImage("/locations/foo/hello-world.png")
	if image.Name != "hello-world.png" {
		t.Errorf("incorrect image name, got: %q, want: %q", image.Name, "hello-world.png")
	}
	if image.Path != "/locations/foo/hello-world.png" {
		t.Errorf("incorrect image path, got: %q, want: %q", image.Path, "/locations/foo/hello-world.png")
	}
	if image.Guid != 0 {
		t.Errorf("incorrect image guid, got: %d, want: %d", image.Guid, 0)
	}
	if image.Features != nil {
		t.Errorf("incorrect image feature vector, got: %v, want: %v", image.Features, nil)
	}
	if image.Id != "hello-world.png" {
		t.Errorf("incorrect image name, got: %q, want: %q", image.Name, "hello-world.png")
	}
}

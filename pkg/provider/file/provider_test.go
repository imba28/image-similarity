package file

import "testing"

func TestImageProvider_Images(t *testing.T) {
	provider := New("../../../template")

	files, err := provider.Images()
	if err != nil {
		t.Error("ImageProvider should not return error")
	}
	if len(files) != 2 {
		t.Errorf("Length of files incorrect, got: %d, want: %d", len(files), 2)
	}
}

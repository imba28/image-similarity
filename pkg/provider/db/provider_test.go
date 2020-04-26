package dbprovider

import (
	"imba28/images/pkg"
	"testing"
)
import "github.com/DATA-DOG/go-sqlmock"

func TestImageProvider_Images(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "guid", "name", "path"}).
		AddRow(1, "10", "first photo", "/images/10.png").
		AddRow(2, "11", "second photo", "/images/11.png").
		AddRow(3, "12", "third photo", "/images/12.png")

	mock.ExpectQuery("SELECT id, guid, name, path FROM photos ORDER BY id DESC").WillReturnRows(rows)
	provider := NewFromDb(db)

	images, err := provider.Images()
	if err != nil {
		t.Error("ImageProvider should not return error")
	}
	if len(images) != 3 {
		t.Errorf("Length of images incorrect, got: %d, want: %d", len(images), 3)
	}

	nameTests := []string{"first photo", "second photo", "third photo"}
	idTests := []string{"1", "2", "3"}
	guidTests := []int{10, 11, 12}
	pathTest := []string{"/images/10.png", "/images/11.png", "/images/12.png"}

	for i, image := range images {
		if image.Name != nameTests[i] {
			t.Errorf("Name of %dth image incorrect, got: %s, want: %s", i, image.Name, nameTests[i])
		}
		if image.Guid != guidTests[i] {
			t.Errorf("Guid of %dth image incorrect, got: %d, want: %d", i, image.Guid, guidTests[i])
		}
		if image.Id != idTests[i] {
			t.Errorf("Id of %dth image incorrect, got: %s, want: %s", i, image.Id, idTests[i])
		}
		if image.Path != pathTest[i] {
			t.Errorf("Path of %dth image incorrect, got: %s, want: %s", i, image.Path, pathTest[i])
		}
		if image.Features != nil {
			t.Errorf("Feature vector of %dth image incorrect, got: %v, want: %v", i, image.Features, nil)
		}
	}
}

func TestImageProvider_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "guid", "name", "path", "vector"}).
		AddRow(1, "10", "first photo", "/images/10.png", nil)

	mock.ExpectQuery("SELECT id, guid, name, path, vector FROM photos WHERE guid = ?").WithArgs("10").WillReturnRows(rows)

	provider := NewFromDb(db)

	image := provider.Get("10")
	if image == nil {
		t.Error("Expected method to find the image object with guid 10")
		return
	}

	if image.Name != "first photo" {
		t.Errorf("Name of image incorrect, got: %s, want: %s", image.Name, "first photo")
	}
	if image.Guid != 10 {
		t.Errorf("Guid of image incorrect, got: %d, want: %d", image.Guid, 10)
	}
	if image.Id != "1" {
		t.Errorf("Id of image incorrect, got: %s, want: %s", image.Id, "1")
	}
	if image.Path != "/images/10.png" {
		t.Errorf("Path of image incorrect, got: %s, want: %s", image.Path, "/images/10.png")
	}
	if image.Features != nil {
		t.Errorf("Feature vector of image incorrect, got: %v, want: %v", image.Features, nil)
	}
}

func TestImageProvider_Get_features(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "guid", "name", "path", "vector"}).
		AddRow(1, "10", "first photo", "/images/10.png", "{0,1,2,3,4,5,6,7,8,9}")

	mock.ExpectQuery("SELECT id, guid, name, path, vector FROM photos WHERE guid = ?").WithArgs("10").WillReturnRows(rows)

	provider := NewFromDb(db)

	image := provider.Get("10")
	if image == nil {
		t.Error("Expected method to find the image object with guid 10")
		return
	}
	if image.Features == nil {
		t.Error("Expected method to load feature vector from database")
		return
	}

	if len(image.Features) != 10 {
		t.Errorf("Incorrect feature vector, got: %v, want %v", image.Features, []float64{0, 1, 2, 3, 4, 5, 6, 7, 9})
	}
	for i := 0; i < 10; i++ {
		if image.Features[i]-float64(i) >= 1e-9 || float64(i)-image.Features[i] >= 1e-9 {
			t.Errorf("Incorrect %dth feature vector item, got: %f, want: %f", i, image.Features[i], float64(i))
		}
	}
}

func TestImageProvider_Persist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE photos SET guid = \\$1, name = \\$2, path = \\$3, vector = \\$4 WHERE id = \\$5").WithArgs(42, "foobar", "/foo/bar.png", nil, "1").WillReturnResult(sqlmock.NewResult(-1, 1))

	image := &pkg.Image{
		Id:   "1",
		Guid: 42,
		Path: "/foo/bar.png",
		Name: "foobar",
	}

	provider := NewFromDb(db)
	err = provider.Persist(image)
	if err != nil {
		t.Errorf("persist should not return error: %q", err)
	}
}

func TestImageProvider_Persist__id(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO photos \\(guid, name, path, vector\\) VALUES \\((.+)\\)  RETURNING id").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	image := &pkg.Image{
		Guid: 42,
		Path: "/foo/bar.png",
		Name: "foobar",
	}

	provider := NewFromDb(db)
	err = provider.Persist(image)
	if err != nil {
		t.Errorf("persist should not return error: %q", err)
	}
	if image.Id != "10" {
		t.Errorf("Invalid id set after persisting object, got: %v, want: %v", image.Id, 10)
	}
}

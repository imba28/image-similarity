package dbprovider

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"imba28/images/pkg"
)

type ImageProvider struct {
	dataSourceName string
}

func (i ImageProvider) Images() ([]pkg.Image, error) {
	db, err := sql.Open("postgres", i.dataSourceName)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, name, image as path FROM locations_photo WHERE visibility = 3 ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []pkg.Image
	for rows.Next() {
		var image pkg.Image
		err = rows.Scan(&image.Id, &image.Name, &image.Path)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func (i ImageProvider) Get(id string) *pkg.Image {
	db, err := sql.Open("postgres", i.dataSourceName)
	if err != nil {
		return nil
	}

	defer db.Close()

	var image pkg.Image
	row := db.QueryRow("SELECT id, name, image as path FROM locations_photo WHERE visibility = 3 AND id = $1", id)
	err = row.Scan(&image.Id, &image.Name, &image.Path)
	if err != nil {
		return nil
	}
	return &image
}

func New(dataSourceString string) ImageProvider {
	return ImageProvider{
		dataSourceName: dataSourceString,
	}
}

func NewFromCredentials(host string, username string, password string, port uint, database string) ImageProvider {
	return ImageProvider{
		dataSourceName: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database),
	}
}

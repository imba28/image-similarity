package dbprovider

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"imba28/images/pkg"
)

type ImageProvider struct {
	dataSourceName string
	db             *sql.DB
}

func (i *ImageProvider) Images() ([]pkg.Image, error) {
	if i.db == nil {
		if err := i.connect(); err != nil {
			fmt.Printf("could not connect to db %q", err)
			return nil, err
		}
	}

	rows, err := i.db.Query("SELECT id, name, image as path FROM locations_photo WHERE visibility = 3 ORDER BY id DESC")
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

func (i *ImageProvider) Get(id string) *pkg.Image {
	if i.db == nil {
		if err := i.connect(); err != nil {
			fmt.Printf("could not connect to db %q", err)
			return nil
		}
	}

	var image pkg.Image
	row := i.db.QueryRow("SELECT id, name, image as path FROM locations_photo WHERE visibility = 3 AND id = $1", id)
	err := row.Scan(&image.Id, &image.Name, &image.Path)
	if err != nil {
		return nil
	}
	return &image
}

func (i *ImageProvider) connect() error {
	db, err := sql.Open("postgres", i.dataSourceName)
	if err != nil {
		return err
	}
	i.db = db

	return nil
}

func New(dataSourceString string) *ImageProvider {
	return &ImageProvider{
		dataSourceName: dataSourceString,
	}
}

func NewFromCredentials(host string, username string, password string, port uint, database string) *ImageProvider {
	return &ImageProvider{
		dataSourceName: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database),
	}
}

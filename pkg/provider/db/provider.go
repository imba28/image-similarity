package dbprovider

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"imba28/images/pkg"
	"log"
)

type ImageProvider struct {
	dataSourceName string
	db             *sql.DB
}

func (i *ImageProvider) Images() ([]*pkg.Image, error) {
	if i.db == nil {
		if err := i.connect(); err != nil {
			log.Printf("could not connect to db %q", err)
			return nil, err
		}
	}

	rows, err := i.db.Query("SELECT id, guid, name, path FROM photos ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*pkg.Image
	for rows.Next() {
		var image pkg.Image
		err = rows.Scan(&image.Id, &image.Guid, &image.Name, &image.Path)
		if err != nil {
			return nil, err
		}

		images = append(images, &image)
	}

	return images, nil
}

func (p *ImageProvider) Persist(i *pkg.Image) error {
	if p.db == nil {
		if err := p.connect(); err != nil {
			log.Printf("could not connect to db %q", err)
			return nil
		}
	}

	if len(i.Id) > 0 {
		_, err := p.db.Exec("UPDATE photos SET guid = $1, name = $2, path = $3, vector = $4 WHERE id = $5", i.Guid, i.Name, i.Path, pq.Array(&i.Features), i.Id)
		return err
	} else {
		err := p.db.QueryRow("INSERT INTO photos (guid, name, path, vector) VALUES ($1, $2, $3, $4)  RETURNING id", i.Guid, i.Name, i.Path, pq.Array(&i.Features)).Scan(&i.Id)
		if err != nil {
			return err
		}
		return nil
	}
}

func (i *ImageProvider) Get(guid string) *pkg.Image {
	if i.db == nil {
		if err := i.connect(); err != nil {
			log.Printf("could not connect to db %q", err)
			return nil
		}
	}

	var image pkg.Image
	row := i.db.QueryRow("SELECT id, guid, name, path FROM photos WHERE guid = $1", guid)
	err := row.Scan(&image.Id, &image.Guid, &image.Name, &image.Path)
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
	return i.migrate()
}

func (i *ImageProvider) migrate() error {
	log.Println("Running migrations")

	driver, err := postgres.WithInstance(i.db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func New(dataSourceString string) *ImageProvider {
	return &ImageProvider{
		dataSourceName: dataSourceString,
	}
}

func NewFromDb(db *sql.DB) *ImageProvider {
	return &ImageProvider{
		db: db,
	}
}

func NewFromCredentials(host string, username string, password string, port uint, database string) *ImageProvider {
	return &ImageProvider{
		dataSourceName: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database),
	}
}

package db

import (
	"database/sql"
	"errors"
	"imba28/images/pkg"
	dbprovider "imba28/images/pkg/provider/db"
	"imba28/images/pkg/provider/file"
)

func CreateImageFixtures(db *sql.DB, dir string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM photos").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("database already contains entries. You should drop it first")
	}

	p := file.New(dir)
	dbp := dbprovider.NewFromDb(db)

	images, err := p.Images()
	if err != nil {
		return err
	}

	for i := range images {
		features, err := pkg.FeatureVector(*images[i])
		if err != nil {
			return err
		}
		images[i].Guid = i
		images[i].Features = features
		images[i].Id = ""
		err = dbp.Persist(images[i])
		if err != nil {
			return err
		}
	}

	return nil
}

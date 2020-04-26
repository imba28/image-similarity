package db

import (
	"database/sql"
	"imba28/images/pkg"
	dbprovider "imba28/images/pkg/provider/db"
	"imba28/images/pkg/provider/file"
)

func CreateImageFixtures(db *sql.DB, dir string) error {
	p := file.New(dir)
	dbp := dbprovider.NewFromDb(db)

	images, err := p.Images()
	if err != nil {
		return err
	}

	for i := range images {
		features, err := pkg.FeatureVector(images[i])
		if err != nil {
			return err
		}
		images[i].Guid = i
		images[i].Features = features
		images[i].Id = ""
		err = dbp.Persist(&images[i])
		if err != nil {
			return err
		}
	}

	return nil
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	db2 "imba28/images/internal/db"
	"log"
	"os"
)

func main() {
	dir := flag.String("directory", "test_sets", "directory containing images")
	flag.Parse()

	db, err := sql.Open("postgres", dsn())
	if err != nil {
		log.Fatal(err)
	}

	err = db2.CreateImageFixtures(db, *dir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ok")
}

func dsn() string {
	if len(os.Getenv("DATABASE_URL")) > 0 {
		return os.Getenv("DATABASE_URL")
	} else {
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), 5432, os.Getenv("POSTGRES_DB"))
	}
}

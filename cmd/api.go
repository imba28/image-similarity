package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"imba28/images/pkg"
	"net/http"
	"strings"
)

const (
	distanceThreshold  = 10
	maxResultSetLength = 10
)

func main() {
	dir := flag.String("directory", "images", "Directory that contains the images set")
	flag.Parse()

	indexTemplate := template.Must(template.ParseFiles("template/index.html"))
	similarTemplate := template.Must(template.ParseFiles("template/similar.html"))

	fmt.Println("Building index...")
	index, err := pkg.NewIndex(*dir)
	if err != nil {
		fmt.Printf("could not open image directory %q\n", *dir)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/similar/", func(w http.ResponseWriter, r *http.Request) {
		photo := getPhoto(r)

		fmt.Printf("Executing query for image %q", photo)

		if len(photo) == 0 {
			w.WriteHeader(400)
			return
		}

		images, err := index.Search(photo, distanceThreshold, maxResultSetLength)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		b, err := json.Marshal(images)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	mux.HandleFunc("/similar/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.Split(r.URL.Path, "/")
		if len(p) != 3 {
			w.WriteHeader(400)
			return
		}

		images, err := index.Search(p[2], distanceThreshold, maxResultSetLength)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(200)
		similarTemplate.Execute(w, images)
	})

	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(*dir))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		descriptors := index.Descriptors()
		w.WriteHeader(200)
		indexTemplate.Execute(w, descriptors)
	})

	fmt.Println("Listening on port 8080...")
	panic(http.ListenAndServe(":8080", mux))
}

func getPhoto(r *http.Request) string {
	p := strings.Split(r.URL.Path, "/")
	if len(p) != 5 {
		return ""
	}

	return p[4]
}

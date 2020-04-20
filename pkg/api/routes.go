package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"imba28/images/pkg"
	"imba28/images/pkg/provider/file"
	"net/http"
	"strings"
)

const (
	distanceThreshold  = 10
	maxResultSetLength = 10
)

func SimilarPhotosJsonHandler(index *pkg.ImageIndex, photoDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		if len(urlParts) != 5 {
			w.WriteHeader(400)
			return
		}
		photo := urlParts[4]
		fmt.Printf("Executing query for image %q", photo)

		images, err := index.Search(file.NewImage(photoDir+"/"+photo), distanceThreshold, maxResultSetLength)
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
	}
}

func SimilarPhotosHandler(index *pkg.ImageIndex, photoDir string) http.HandlerFunc {
	similarTemplate := template.Must(template.ParseFiles("template/similar.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		p := strings.Split(r.URL.Path, "/")
		if len(p) != 3 {
			w.WriteHeader(400)
			return
		}

		images, err := index.Search(file.NewImage(photoDir+"/"+p[2]), distanceThreshold, maxResultSetLength)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(200)
		similarTemplate.Execute(w, images)
	}
}

func IndexHandler(index *pkg.ImageIndex) http.HandlerFunc {
	indexTemplate := template.Must(template.ParseFiles("template/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		descriptors := index.Descriptors()
		w.WriteHeader(200)
		indexTemplate.Execute(w, descriptors)
	}
}

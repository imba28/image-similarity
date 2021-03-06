package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"imba28/images/pkg"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	distanceThreshold  = 10
	maxResultSetLength = 10
	itemsPerPage       = 15
	maxPageItems       = 15
)

func SimilarPhotosJsonHandler(index *pkg.ImageIndex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		if len(urlParts) != 5 {
			w.WriteHeader(400)
			return
		}
		photo := index.Load(urlParts[4])
		if photo == nil {
			w.WriteHeader(404)
			return
		}

		fmt.Printf("Executing query for image %q", photo.Id)
		images, err := index.Search(*photo, distanceThreshold, maxResultSetLength)
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

func SimilarPhotosHandler(index *pkg.ImageIndex) http.HandlerFunc {
	similarTemplate := initTemplate("template/similar.html")

	return func(w http.ResponseWriter, r *http.Request) {
		p := strings.Split(r.URL.Path, "/")
		if len(p) != 3 {
			w.WriteHeader(400)
			return
		}

		image := index.Load(p[2])
		if image == nil {
			w.WriteHeader(404)
			return
		}

		imageDistances, err := index.Search(*image, distanceThreshold, maxResultSetLength)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		if len(imageDistances) > 0 && imageDistances[0].Image.Id == image.Id {
			imageDistances = imageDistances[1:]
		}

		view := struct {
			Images    []pkg.ImageDistance
			Image     pkg.Image
			MediaRoot string
			Random    int
		}{
			Images: imageDistances,
			Image:  *image,
			Random: rand.Intn(len(index.Images())),
		}

		similarTemplate.Execute(w, view)
	}
}

func IndexHandler(index *pkg.ImageIndex) http.HandlerFunc {
	indexTemplate := initTemplate("template/index.html")

	return func(w http.ResponseWriter, r *http.Request) {
		images := index.Images()
		pageParam := r.URL.Query().Get("page")
		page := 1
		if len(pageParam) > 0 {
			p, err := strconv.Atoi(pageParam)
			if err == nil && p > 0 {
				page = p
			}
		}

		maxPageNumber := int(math.Ceil(float64(len(index.Images())) / itemsPerPage))
		if page > maxPageNumber {
			page = maxPageNumber
		}

		lowerBound := (page - 1) * itemsPerPage
		upperBound := lowerBound + itemsPerPage
		if upperBound >= len(index.Images()) {
			upperBound = len(index.Images())
		}

		pages := make([]int, int(math.Min(float64(maxPageNumber), maxPageItems)))
		for i := range pages {
			offset := 0
			if page-1 > maxPageItems/2 {
				offset = (page - 1) - maxPageItems/2
			}
			pages[i] = offset + (i + 1)

			if page+maxPageItems > maxPageNumber && pages[i] > maxPageNumber {
				pages = pages[:i]
				break
			}
		}

		view := struct {
			Images    []*pkg.Image
			Page      int
			Total     int
			Pages     []int
			MediaRoot string
			Random    int
		}{
			Images: images[lowerBound:upperBound],
			Page:   page,
			Pages:  pages,
			Total:  len(index.Images()),
			Random: rand.Intn(len(index.Images())),
		}

		indexTemplate.Execute(w, view)
	}
}

func initTemplate(path string) *template.Template {
	t := template.New(filepath.Base(path)).Funcs(template.FuncMap{
		"FormatDistance": func(distance float64) float64 {
			return math.Ceil(distance*100) / 100
		},
	})

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	return template.Must(t.Parse(string(contents)))
}

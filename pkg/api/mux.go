package api

import (
	"imba28/images/pkg"
	"net/http"
)

func New(index *pkg.ImageIndex, staticFolder string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/similar/", SimilarPhotosJsonHandler(index, staticFolder))
	mux.HandleFunc("/similar/", SimilarPhotosHandler(index, staticFolder))
	mux.Handle("/"+staticFolder+"/", http.StripPrefix("/"+staticFolder+"/", http.FileServer(http.Dir(staticFolder))))
	mux.HandleFunc("/", IndexHandler(index))

	return mux
}

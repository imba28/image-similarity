package api

import (
	"imba28/images/pkg"
	"net/http"
)

func New(index *pkg.ImageIndex, staticFolder string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/similar/", LoggingMiddleware(SimilarPhotosJsonHandler(index)))
	mux.Handle("/similar/", LoggingMiddleware(SimilarPhotosHandler(index)))
	mux.Handle("/"+staticFolder+"/", http.StripPrefix("/"+staticFolder+"/", http.FileServer(http.Dir(staticFolder))))
	mux.Handle("/", LoggingMiddleware(IndexHandler(index)))

	return mux
}

package api

import (
	"imba28/images/pkg"
	"log"
	"net/http"
	"strings"
)

func New(index *pkg.ImageIndex, mediaRoot string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/similar/", LoggingMiddleware(SimilarPhotosJsonHandler(index)))
	mux.Handle("/similar/", LoggingMiddleware(SimilarPhotosHandler(index)))
	mux.Handle("/media/", http.StripPrefix("/media/"+strings.Trim(mediaRoot, "/"), http.FileServer(http.Dir(mediaRoot))))
	mux.Handle("/", LoggingMiddleware(IndexHandler(index)))

	log.Printf("Serving static files from %q", mediaRoot)

	return mux
}

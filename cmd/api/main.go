package main

import (
	"flag"
	"fmt"
	"imba28/images/pkg"
	"imba28/images/pkg/api"
	"imba28/images/pkg/provider/db"
	"net/http"
	"os"
	"strconv"
)

func main() {
	dir := flag.String("directory", "locations", "Directory that contains the images set")
	port := flag.Uint("port", 8080, "Port to bind http server to")
	flag.Parse()

	if len(os.Getenv("PORT")) > 0 {
		p, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			fmt.Printf("cannot convert %q to number!\n", os.Getenv("PORT"))
		}
		pp := uint(p)
		port = &pp
	}

	fmt.Println("Building index...")
	index, err := pkg.NewIndex(dbprovider.New(os.Getenv("DATABASE_URL")))
	if err != nil {
		fmt.Printf("could not connect to database %q\n", *dir)
		return
	}

	mux := api.New(index, *dir)

	address := fmt.Sprintf(":%d", *port)
	fmt.Printf("Listening on %s", address)

	panic(http.ListenAndServe(address, mux))
}
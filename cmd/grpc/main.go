package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"imba28/images/pkg"
	"imba28/images/pkg/pb"
	dbprovider "imba28/images/pkg/provider/db"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	port := flag.Uint("port", 8080, "Port to bind grpc server to")
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

	index, err := pkg.NewIndex(getImageProvider())
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterImageSimilarityServiceServer(grpcServer, pb.NewImageSimilarityService(index))

	fmt.Printf("Listening on port :%d\n", *port)
	panic(grpcServer.Serve(listener))
}

func getImageProvider() pkg.ImageProvider {
	if len(os.Getenv("DATABASE_URL")) > 0 {
		return dbprovider.New(os.Getenv("DATABASE_URL"))
	} else {
		return dbprovider.NewFromCredentials(os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), 5432, os.Getenv("POSTGRES_DB"))
	}
}

# Image similarity

This application extracts feature vectors from images, builds an index and queries the index to find similar images. 

My initial inspiration came from reading one of [Adrian Rosebrock's blog articles](https://www.pyimagesearch.com/2014/12/01/complete-guide-building-image-search-engine-python-opencv/), where he describes a simple image search algorithm written in Python.  Though for the most part the calculation of feature vectors is based on Adrian's implementation, I am looking forward to extending it as I learn more about image processing and Golang.

I am working on this module as part of my master's project at the University of Applied Sciences Salzburg.
## Requirements
- OpenCv 4
- Go 1.13

## Install
```shell script
go mod download
```

## Run
**Index & Web Interface**:

```shell script
go run cmd/api/main.go
# or
go build -o api cmd/api/main.go 
./api -directory ./directory/containing/images
```

**gRPC API**
```shell script
go run cmd/grpc/main.go
# or
go build -o api cmd/grpc/main.go 
DATABASE_URL=postgres://user:password@host/database ./api
```


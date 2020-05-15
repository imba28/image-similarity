# Image similarity

This application extracts feature vectors from images, builds an index and queries the index to find similar images. 

My initial inspiration came from reading one of [Adrian Rosebrock's blog articles](https://www.pyimagesearch.com/2014/12/01/complete-guide-building-image-search-engine-python-opencv/), where he describes a simple image search algorithm written in Python.  
The calculation of the feature vectors is based on Adrian's implementation, and I am looking forward to extending it as I learn more about image processing.

I am working on this module as part of my master's project at the University of Applied Sciences Salzburg.

## How it looks like
![demo](doc/demo.png)

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
PORT=8080 ./api -directory ./directory/containing/images
```

To import images from a directory and save them to a database run the following commands:

```shell script
go build -o api cmd/api/main.go 
go build -o fixtures cmd/fixtures/main.go 

docker-compose up -d
DATABASE_URL=postgres://user:1234@127.0.0.1/imgsim ./fixtures -directory ./directory/containing/images
./api
```

**gRPC API**
```shell script
go run cmd/grpc/main.go
# or
go build -o api cmd/grpc/main.go 
DATABASE_URL=postgres://user:password@host/database ./api
```

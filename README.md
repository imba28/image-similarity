# Image similarity

This application extracts feature vectors from images, builds an index and queries the index to find similar images. 

I was inspired by [Adrian Rosebrock's blog article](https://www.pyimagesearch.com/2014/12/01/complete-guide-building-image-search-engine-python-opencv/), where he describes a simple image search algorithm written in Python. For the most part, the calculation of feature vectors is based on Adrian's implementation. 

I am working on this module as part of my master's project at the University of Applied Sciences Salzburg.

## Requirements
- OpenCv 4
- Go 1.13

## Install
```shell script
go mod download
```

## Run
```shell script
go run cmd/api/main.go
# or
go build -o api cmd/api/main.go 
./api -directory ./directory/containing/images
```


# Image similarity

This application extracts feature vectors from images, builds an index and queries the index to find similar images.

## Requirements
- OpenCv 4
- Go 1.13

## Install
```shell script
go mod download
```

## Run
```shell script
go run cmd/cli.go
# or
go build cmd/cli.go
./cli -directory imageDirectory
```


all: proto build

test:
	go test ./...

build:
	go build -o api ./cmd/grpc/main.go

.PHONY: proto
proto:
	 protoc -I proto proto/image.proto --go_out=plugins=grpc:pkg/pb

migrate:
	migrate -database ${DATABASE_URL} -path db/migrations up

cover:
	go test -cover ./...

profile:
	go test -coverprofile=./.profile.out ./...
	go tool cover -html=./.profile.out && rm ./.profile.out


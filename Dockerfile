FROM imba28/gocv:1.13-builder as builder

ENV APP_HOME /app

RUN apk update
RUN apk add protoc

RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go get -v github.com/golang/protobuf/protoc-gen-go
RUN protoc -I proto proto/image.proto --go_out=plugins=grpc:pkg/pb
RUN go build -o api ./cmd/grpc/main.go

FROM imba28/gocv:1.13

COPY --from=builder /app/api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
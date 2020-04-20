FROM imba28/gocv:1.13-builder as builder

ENV APP_HOME /app

RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o api ./cmd/grpc/main.go

FROM imba28/gocv:1.13

COPY --from=builder /app/api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
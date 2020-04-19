FROM imba28/gocv:1.13-builder as builder

ENV APP_HOME /app

RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o app ./cmd/api.go

FROM imba28/gocv:1.13

COPY --from=builder /app/app /app
COPY --from=builder /app/template /template

EXPOSE 8080
ENTRYPOINT ["/app", "-directory"]
CMD ["/images"]
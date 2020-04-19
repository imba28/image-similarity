FROM imba28/gocv:1.13-builder as builder

ENV APP_HOME /app

RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY . $APP_HOME
RUN go build -o app ./cmd/cli.go

FROM imba28/gocv:1.13

COPY --from=builder /app/app /app
ENTRYPOINT ["/app", "-directory"]
CMD ["/images"]
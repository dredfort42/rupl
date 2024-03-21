FROM golang:latest AS env

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
RUN go mod download

FROM env AS build

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o ./web /app/cmd/web
RUN if [ ! -e "/app/local.cfg" ]; then \
    touch "/app/local.cfg"; \
    fi

FROM scratch
COPY --from=build /app/web /app/web
COPY --from=build /app/local.cfg /bin/local.cfg

EXPOSE 80
CMD ["/app/web"]
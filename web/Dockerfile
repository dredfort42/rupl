FROM golang:latest AS build

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o ./web /app/cmd/web
RUN if [ ! -e "/app/local.cfg" ]; then \
    touch "/app/local.cfg"; \
    fi

FROM scratch
COPY --from=build /app/web /bin/web
COPY --from=build /app/local.cfg /app/local.cfg

EXPOSE 80
CMD ["/bin/web"]
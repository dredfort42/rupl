FROM golang:latest

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o ./web /app/cmd/web

EXPOSE 80
CMD ["/app/web"]
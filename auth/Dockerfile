FROM golang:latest

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o ./auth /app/cmd/auth

EXPOSE 4242
CMD ["/app/auth"]
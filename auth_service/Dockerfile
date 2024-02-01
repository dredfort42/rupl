FROM golang:latest

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN go build -o /bin/auth /app/cmd/auth

EXPOSE 4242

CMD ["/bin/manager"]
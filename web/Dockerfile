FROM golang:latest AS build

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/web /app/cmd/web

FROM scratch
COPY --from=build /app/web /bin/web

EXPOSE 80
CMD ["/bin/web"]
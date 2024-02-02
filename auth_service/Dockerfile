FROM golang:latest AS build

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/auth /app/cmd/auth

FROM scratch
COPY --from=build /app/auth /bin/auth

EXPOSE 4242

CMD ["/bin/auth"]
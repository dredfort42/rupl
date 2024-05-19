FROM golang:latest AS env

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
RUN go mod download

FROM env AS build

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o ./sessions_receiver /app/cmd/sessions_receiver
RUN if [ ! -e "/app/local.cfg" ]; then \
    touch "/app/local.cfg"; \
    fi

FROM scratch
COPY --from=build /app/sessions_receiver /app/sessions_receiver
COPY --from=build /app/local.cfg /app/local.cfg

EXPOSE 4444
CMD ["/app/sessions_receiver"]
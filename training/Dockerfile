FROM golang:latest AS env

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
RUN go mod download

FROM env AS build

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o ./training /app/cmd/training
RUN if [ ! -e "/app/local.cfg" ]; then \
    touch "/app/local.cfg"; \
    fi

FROM scratch
COPY --from=build /app/training /app/training
COPY --from=build /app/local.cfg /app/local.cfg

EXPOSE 4646
CMD ["/app/training"]
FROM golang:latest AS env

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
RUN go mod tidy

FROM env AS build

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o ./profile_manager /app/cmd/profile_manager
RUN if [ ! -e "/app/local.cfg" ]; then \
    touch "/app/local.cfg"; \
    fi

FROM scratch
COPY --from=build /app/profile_manager /app/profile_manager
COPY --from=build /app/local.cfg /app/local.cfg

EXPOSE 4848
CMD ["/app/profile_manager"]
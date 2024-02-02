FROM golang:latest

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o ./auth /app/cmd/auth
RUN if [ ! -e "/app/local.cfg" ]; then \
    touch "/app/local.cfg"; \
    fi

FROM scratch
COPY --from=build /app/auth /app/auth
COPY --from=build /app/local.cfg /app/local.cfg

EXPOSE 4242
CMD ["/app/auth"]
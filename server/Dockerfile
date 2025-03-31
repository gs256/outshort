FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/dist/server

COPY .env.template .env
RUN mkdir -p env

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN /go/bin/goose up


FROM debian:stable-slim
WORKDIR /app

COPY --from=build /app/dist/server /app
COPY --from=build /app/env /app/env

CMD ["/app/server"]
EXPOSE 8249

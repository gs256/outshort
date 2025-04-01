FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/dist/server


FROM debian:stable-slim
WORKDIR /app
COPY --from=build /app/dist/server /app
CMD ["/app/server"]
EXPOSE 8249

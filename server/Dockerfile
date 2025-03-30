FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o dist/server
CMD ["./dist/server"]
EXPOSE 8249

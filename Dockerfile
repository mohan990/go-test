FROM golang:1.22-alpine

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]

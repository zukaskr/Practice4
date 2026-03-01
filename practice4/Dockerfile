FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o myapp ./cmd/api/main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/myapp /myapp
EXPOSE 8080
ENTRYPOINT ["/myapp"]